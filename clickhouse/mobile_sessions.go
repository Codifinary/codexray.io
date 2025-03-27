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
	GeoMapColorCode   string
}

type SessionHistoricData struct {
	SessionID       string
	UserID          string
	Country         string
	NoOfRequest     uint64
	SessionDuration int64
	LastPage        string
	StartTime       timeseries.Time
	GeoMapColorCode string
}

func (c *Client) GetMobileSessionResults(ctx context.Context, from, to timeseries.Time) (*MobileSessionResult, error) {

	fromTime := from.ToStandard()
	toTime := to.ToStandard()

	windowDuration := toTime.Sub(fromTime)
	prevToTime := fromTime
	prevFromTime := prevToTime.Add(-windowDuration)

	prevFrom := timeseries.Time(prevFromTime.Unix())
	prevTo := timeseries.Time(prevToTime.Unix())

	query := `
	WITH
    current AS (
        SELECT
            uniqExact(SessionId) AS totalSessions,
            uniqExact(UserId) AS TotalUsers
        FROM mobile_session_data
        WHERE 
            Timestamp BETWEEN @from AND @to
    ),
    previous AS (
        SELECT
            uniqExact(SessionId) AS totalSessions,
            uniqExact(UserId) AS TotalUsers
        FROM mobile_session_data
        WHERE 
            Timestamp BETWEEN @prevFrom AND @prevTo
    ),
    current_duration AS (
        SELECT
            avg(dateDiff('second', StartTime, EndTime)) AS avgDuration 
        FROM mobile_session_data
        WHERE 
            Timestamp BETWEEN @from AND @to
            AND StartTime IS NOT NULL 
            AND EndTime IS NOT NULL  
    ),
    previous_duration AS (
        SELECT
            avg(dateDiff('second', StartTime, EndTime)) AS avgDuration 
        FROM mobile_session_data
        WHERE 
            Timestamp BETWEEN @prevFrom AND @prevTo
            AND StartTime IS NOT NULL 
            AND EndTime IS NOT NULL  
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
	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("prevFrom", prevFrom.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("prevTo", prevTo.ToStandard(), clickhouse.NanoSeconds),
	)
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

func (c *Client) GetSessionsByCountryTrendChart(ctx context.Context, from, to timeseries.Time, step timeseries.Duration) (map[string]*timeseries.TimeSeries, error) {
	optimizedQuery := fmt.Sprintf(`
	SELECT
		toUnixTimestamp(toStartOfInterval(Timestamp, INTERVAL %d SECOND)) * 1000 as interval_start,
		Country as country,
		count(DISTINCT SessionId) as session_count
	FROM
		mobile_session_data
	WHERE
		Timestamp BETWEEN @from AND @to
		AND Country != ''
		AND Country IN (
			SELECT Country
			FROM mobile_session_data
			WHERE Timestamp BETWEEN @from AND @to
			AND Country != ''
			GROUP BY Country
			ORDER BY count(DISTINCT SessionId) DESC
			LIMIT 3
		)
	GROUP BY
		interval_start, country
	ORDER BY
		interval_start, country
	`, step)

	rows, err := c.Query(ctx, optimizedQuery,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
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

func (c *Client) GetSessionsByDeviceTrendChart(ctx context.Context, from, to timeseries.Time, step timeseries.Duration) (map[string]*timeseries.TimeSeries, error) {
	optimizedQuery := fmt.Sprintf(`
	SELECT
		toUnixTimestamp(toStartOfInterval(Timestamp, INTERVAL %d SECOND)) * 1000 as interval_start,
		Device as device,
		count(DISTINCT SessionId) as session_count
	FROM
		mobile_session_data
	WHERE
		Timestamp BETWEEN @from AND @to
		AND Device != ''
		AND Device IN (
			SELECT Device
			FROM mobile_session_data
			WHERE Timestamp BETWEEN @from AND @to
			AND Device != ''
			GROUP BY Device
			ORDER BY count(DISTINCT SessionId) DESC
			LIMIT 3
		)
	GROUP BY
		interval_start, device
	ORDER BY
		interval_start, device
	`, step)

	rows, err := c.Query(ctx, optimizedQuery,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
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

func (c *Client) GetSessionsByOSTrendChart(ctx context.Context, from, to timeseries.Time, step timeseries.Duration) (map[string]*timeseries.TimeSeries, error) {
	optimizedQuery := fmt.Sprintf(`
	SELECT
		toUnixTimestamp(toStartOfInterval(Timestamp, INTERVAL %d SECOND)) * 1000 as interval_start,
		OS as os,
		count(DISTINCT SessionId) as session_count
	FROM
		mobile_session_data
	WHERE
		Timestamp BETWEEN @from AND @to
		AND OS != ''
		AND OS IN (
			SELECT OS
			FROM mobile_session_data
			WHERE Timestamp BETWEEN @from AND @to
			AND OS != ''
			GROUP BY OS
			ORDER BY count(DISTINCT SessionId) DESC
			LIMIT 3
		)
	GROUP BY
		interval_start, os
	ORDER BY
		interval_start, os
	`, step)

	rows, err := c.Query(ctx, optimizedQuery,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
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

func (c *Client) GetSessionLiveData(ctx context.Context, from, to timeseries.Time) ([]SessionLiveData, error) {

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
	GROUP BY s.SessionId, s.UserId, s.StartTime, s.Country
	ORDER BY NoOfRequest DESC
	`

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	)
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

		switch {
		case session.NoOfRequest <= 20:
			session.GeoMapColorCode = "#5BBC7A" // Green
		case session.NoOfRequest <= 80:
			session.GeoMapColorCode = "#F1AB47" // Yellow
		default: // > 80%
			session.GeoMapColorCode = "#EF5350" // Red
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

func (c *Client) GetSessionHistoricData(ctx context.Context, from, to timeseries.Time) ([]SessionHistoricData, error) {
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
	GROUP BY s.SessionId, s.UserId, s.StartTime, s.EndTime, s.Country
	ORDER BY NoOfRequest DESC
	`

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	)
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

		switch {
		case session.NoOfRequest <= 20:
			session.GeoMapColorCode = "#5BBC7A" // Green
		case session.NoOfRequest <= 80:
			session.GeoMapColorCode = "#F1AB47" // Yellow
		default: // > 80%
			session.GeoMapColorCode = "#EF5350" // Red
		}

		session.StartTime = timeseries.Time(startTime.Unix())
		result = append(result, session)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
