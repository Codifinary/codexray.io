package clickhouse

import (
	"codexray/timeseries"
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type MobileSessionResult struct {
	TotalSessions       uint64
	SessionTrend        float64
	TotalUsers          uint64
	UserTrend           float64
	AverageSession      float64
	AverageSessionTrend float64
}

type SessionLiveData struct {
	SessionID         string
	UserID            string
	Country           string
	NoOfRequest       uint64
	LastPageTimestamp timeseries.Time
	LastPage          string
	StartTime         timeseries.Time
}

type SessionHistoricData struct {
	SessionID       string
	UserID          string
	Country         string
	NoOfRequest     uint64
	SessionDuration int64
	LastPage        string
	StartTime       timeseries.Time
}

type SessionGeoMapData struct {
	Country         string
	Count           uint64
	GeoMapColorCode string
}

func (c *Client) GetMobileSessionResults(ctx context.Context, from, to timeseries.Time, service string) (*MobileSessionResult, error) {

	fromTime := from.ToStandard()
	toTime := to.ToStandard()

	windowDuration := toTime.Sub(fromTime)
	prevToTime := fromTime
	prevFromTime := prevToTime.Add(-windowDuration)

	prevFrom := timeseries.Time(prevFromTime.Unix())
	prevTo := timeseries.Time(prevToTime.Unix())

	serviceFilter := "AND Service = @service"

	query := `
	WITH
    current AS (
        SELECT
            uniqExact(msd.SessionId) AS totalSessions,
            uniqExact(msd.UserId) AS TotalUsers
        FROM mobile_session_data msd
        LEFT JOIN mobile_event_data med ON msd.SessionId = med.SessionId
        WHERE 
            msd.Timestamp BETWEEN @from AND @to
            ` + serviceFilter + `
    ),
    previous AS (
        SELECT
            uniqExact(msd.SessionId) AS totalSessions,
            uniqExact(msd.UserId) AS TotalUsers
        FROM mobile_session_data msd
        LEFT JOIN mobile_event_data med ON msd.SessionId = med.SessionId
        WHERE 
            msd.Timestamp BETWEEN @prevFrom AND @prevTo
            ` + serviceFilter + `
    ),
    current_duration AS (
        SELECT
            avg(dateDiff('second', msd.StartTime, msd.EndTime)) AS avgDuration 
        FROM mobile_session_data msd
        LEFT JOIN mobile_event_data med ON msd.SessionId = med.SessionId
        WHERE 
            msd.Timestamp BETWEEN @from AND @to
            AND msd.StartTime IS NOT NULL 
            AND msd.EndTime IS NOT NULL
            ` + serviceFilter + `
    ),
    previous_duration AS (
        SELECT
            avg(dateDiff('second', msd.StartTime, msd.EndTime)) AS avgDuration 
        FROM mobile_session_data msd
        LEFT JOIN mobile_event_data med ON msd.SessionId = med.SessionId
        WHERE 
            msd.Timestamp BETWEEN @prevFrom AND @prevTo
            AND msd.StartTime IS NOT NULL 
            AND msd.EndTime IS NOT NULL
            ` + serviceFilter + `
    )
SELECT 
    current.totalSessions,
    least(100, greatest(-100, (current.totalSessions - previous.totalSessions) / greatest(1, toFloat64(previous.totalSessions)) * 100 )) AS sessionTrend,
    current.TotalUsers,
    least(100, greatest(-100, (current.TotalUsers - previous.TotalUsers) / greatest(1, toFloat64(previous.TotalUsers)) * 100 )) AS userTrend,
    current_duration.avgDuration AS avgSession,
    least(100, greatest(-100, (ifNull(current_duration.avgDuration, 0) - ifNull(previous_duration.avgDuration, 0)) / greatest(1, toFloat64(ifNull(previous_duration.avgDuration, 1))) * 100 )) 
	AS avgSessionTrend
FROM current, previous, current_duration, previous_duration
`
	params := []interface{}{
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("prevFrom", prevFrom.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("prevTo", prevTo.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.Named("service", service),
	}

	rows, err := c.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result MobileSessionResult
	if rows.Next() {
		if err := rows.Scan(
			&result.TotalSessions,
			&result.SessionTrend,
			&result.TotalUsers,
			&result.UserTrend,
			&result.AverageSession,
			&result.AverageSessionTrend,
		); err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (c *Client) GetSessionsByCountryTrendChart(ctx context.Context, from, to timeseries.Time, step timeseries.Duration, service string) (map[string]*timeseries.TimeSeries, error) {

	query := fmt.Sprintf(`
	SELECT
		toUnixTimestamp(toStartOfInterval(msd.Timestamp, INTERVAL %d SECOND)) * 1000 as interval_start,
		msd.Country as country,
		count(DISTINCT msd.SessionId) as session_count
	FROM
		mobile_session_data msd
	LEFT JOIN mobile_event_data med ON msd.SessionId = med.SessionId
	WHERE
		msd.Timestamp BETWEEN @from AND @to
		AND msd.Country != ''
		AND med.Service = @service
		AND msd.Country IN (
			SELECT msd2.Country
			FROM mobile_session_data msd2
			LEFT JOIN mobile_event_data med2 ON msd2.SessionId = med2.SessionId
			WHERE msd2.Timestamp BETWEEN @from AND @to
			AND msd2.Country != ''
			AND med2.Service = @service
			GROUP BY msd2.Country
			ORDER BY count(DISTINCT msd2.SessionId) DESC
			LIMIT 3
		)
	GROUP BY
		interval_start, country
	ORDER BY
		interval_start, country
	`, step)

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.Named("service", service),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]*timeseries.TimeSeries)

	for rows.Next() {
		var intervalStart uint64
		var country string
		var sessionCount uint64

		if err := rows.Scan(&intervalStart, &country, &sessionCount); err != nil {
			return nil, err
		}

		ts, exists := result[country]
		if !exists {
			ts = timeseries.New(from, int(to.Sub(from)/step), step)
			result[country] = ts
		}

		ts.Set(timeseries.Time(intervalStart/1000), float32(sessionCount))
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}

func (c *Client) GetSessionsByDeviceTrendChart(ctx context.Context, from, to timeseries.Time, step timeseries.Duration, service string) (map[string]*timeseries.TimeSeries, error) {

	query := fmt.Sprintf(`
	SELECT
		toUnixTimestamp(toStartOfInterval(msd.Timestamp, INTERVAL %d SECOND)) * 1000 as interval_start,
		msd.Device as device,
		count(DISTINCT msd.SessionId) as session_count
	FROM
		mobile_session_data msd
	LEFT JOIN mobile_event_data med ON msd.SessionId = med.SessionId
	WHERE
		msd.Timestamp BETWEEN @from AND @to
		AND msd.Device != ''
		AND med.Service = @service
		AND msd.Device IN (
			SELECT msd2.Device
			FROM mobile_session_data msd2
			LEFT JOIN mobile_event_data med2 ON msd2.SessionId = med2.SessionId
			WHERE msd2.Timestamp BETWEEN @from AND @to
			AND msd2.Device != ''
			AND med2.Service = @service
			GROUP BY msd2.Device
			ORDER BY count(DISTINCT msd2.SessionId) DESC
			LIMIT 3
		)
	GROUP BY
		interval_start, device
	ORDER BY
		interval_start, device
	`, step)

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.Named("service", service),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]*timeseries.TimeSeries)

	for rows.Next() {
		var intervalStart uint64
		var device string
		var sessionCount uint64

		if err := rows.Scan(&intervalStart, &device, &sessionCount); err != nil {
			return nil, err
		}

		ts, exists := result[device]
		if !exists {
			ts = timeseries.New(from, int(to.Sub(from)/step), step)
			result[device] = ts
		}

		ts.Set(timeseries.Time(intervalStart/1000), float32(sessionCount))
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}

func (c *Client) GetSessionsByOSTrendChart(ctx context.Context, from, to timeseries.Time, step timeseries.Duration, service string) (map[string]*timeseries.TimeSeries, error) {

	query := fmt.Sprintf(`
	SELECT
		toUnixTimestamp(toStartOfInterval(msd.Timestamp, INTERVAL %d SECOND)) * 1000 as interval_start,
		msd.OS as os,
		count(DISTINCT msd.SessionId) as session_count
	FROM
		mobile_session_data msd
	JOIN
		mobile_event_data med ON msd.SessionId = med.SessionId
	WHERE
		msd.Timestamp BETWEEN @from AND @to
		AND msd.OS != ''
		AND med.Service = @service
		AND msd.OS IN (
			SELECT msd2.OS
			FROM mobile_session_data msd2
			JOIN mobile_event_data med2 ON msd2.SessionId = med2.SessionId
			WHERE msd2.Timestamp BETWEEN @from AND @to
			AND msd2.OS != ''
			AND med2.Service = @service
			GROUP BY msd2.OS
			ORDER BY count(DISTINCT msd2.SessionId) DESC
			LIMIT 3
		)
	GROUP BY
		interval_start, os
	ORDER BY
		interval_start, os
	`, step)

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.Named("service", service),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]*timeseries.TimeSeries)

	for rows.Next() {
		var intervalStart uint64
		var os string
		var sessionCount uint64

		if err := rows.Scan(&intervalStart, &os, &sessionCount); err != nil {
			return nil, err
		}

		ts, exists := result[os]
		if !exists {
			ts = timeseries.New(from, int(to.Sub(from)/step), step)
			result[os] = ts
		}

		ts.Set(timeseries.Time(intervalStart/1000), float32(sessionCount))
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}

func (c *Client) GetSessionLiveData(ctx context.Context, from, to timeseries.Time, limit int, service string) ([]SessionLiveData, error) {

	query := `
	SELECT
		s.SessionId,
		s.UserId,
		s.StartTime,
		COUNT(p.EndpointName) AS NoOfRequest,
		argMax(p.EndpointName, p.Timestamp) AS LastPage,
		argMax(p.Timestamp, p.Timestamp) AS LastPageTimestamp,
		s.Country
	FROM mobile_session_data s
	LEFT JOIN mobile_perf_data p ON s.SessionId = p.SessionId
	WHERE s.StartTime BETWEEN @from AND @to
	AND s.EndTime IS NULL
	AND p.Service = @service
	GROUP BY s.SessionId, s.UserId, s.StartTime, s.Country
	ORDER BY NoOfRequest DESC
	LIMIT @limit
	`

	params := []interface{}{
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.Named("service", service),
		clickhouse.Named("limit", uint64(limit)),
	}

	rows, err := c.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []SessionLiveData

	for rows.Next() {
		var session SessionLiveData
		var startTime, lastPageTimestamp time.Time

		if err := rows.Scan(
			&session.SessionID,
			&session.UserID,
			&startTime,
			&session.NoOfRequest,
			&session.LastPage,
			&lastPageTimestamp,
			&session.Country,
		); err != nil {
			return nil, err
		}

		session.StartTime = timeseries.Time(startTime.Unix())
		session.LastPageTimestamp = timeseries.Time(lastPageTimestamp.Unix())
		result = append(result, session)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) GetSessionHistoricData(ctx context.Context, from, to timeseries.Time, limit int, service string) ([]SessionHistoricData, error) {

	query := `
	SELECT
		s.SessionId,
		s.UserId,
		s.StartTime,
		COUNT(p.EndpointName) AS NoOfRequest,
		argMax(p.EndpointName, p.Timestamp) AS LastPage,
		TIMESTAMPDIFF(SECOND, s.StartTime, s.EndTime) AS SessionDuration,
		s.Country
	FROM mobile_session_data s
	LEFT JOIN mobile_perf_data p ON s.SessionId = p.SessionId
	WHERE s.StartTime BETWEEN @from AND @to
	AND s.EndTime IS NOT NULL
	AND p.Service = @service
	GROUP BY s.SessionId, s.UserId, s.StartTime, s.EndTime, s.Country
	ORDER BY NoOfRequest DESC
	LIMIT @limit
	`

	params := []interface{}{
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.Named("service", service),
		clickhouse.Named("limit", uint64(limit)),
	}

	rows, err := c.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []SessionHistoricData

	for rows.Next() {
		var session SessionHistoricData
		var startTime time.Time

		if err := rows.Scan(
			&session.SessionID,
			&session.UserID,
			&startTime,
			&session.NoOfRequest,
			&session.LastPage,
			&session.SessionDuration,
			&session.Country,
		); err != nil {
			return nil, err
		}

		session.StartTime = timeseries.Time(startTime.Unix())
		result = append(result, session)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) GetSessionGeoMapData(ctx context.Context, from, to timeseries.Time, service string) ([]SessionGeoMapData, error) {
	query := `
	SELECT
		mpd.Country as country,
		count() as request_count
	FROM
		mobile_perf_data mpd
	WHERE
		mpd.Timestamp BETWEEN @from AND @to
		AND mpd.Country != ''
		AND mpd.Service = @service
	GROUP BY
		mpd.Country
	ORDER BY
		request_count DESC
	`

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.Named("service", service),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []SessionGeoMapData

	for rows.Next() {
		var data SessionGeoMapData

		if err := rows.Scan(&data.Country, &data.Count); err != nil {
			return nil, err
		}

		switch {
		case data.Count <= 20:
			data.GeoMapColorCode = "#EF5350" // red
		case data.Count <= 80:
			data.GeoMapColorCode = "#F1AB47" // Yellow
		default: // > 80
			data.GeoMapColorCode = "#5BBC7A" // Green
		}

		result = append(result, data)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
