package clickhouse

import (
	"codexray/timeseries"
	"context"

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
