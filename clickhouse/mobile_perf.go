package clickhouse

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type MobilePerfResult struct {
	TotalRequests          uint64
	RequestsPerSecond      float64
	RequestsTrend          float64
	TotalErrors            uint64
	ErrorsPerSecond        float64
	ErrorsTrend            float64
	UsersImpacted          uint64
	UsersImpactedPerSecond float64
	UsersImpactedTrend     float64
}

func (c *Client) GetMobilePerfResults(ctx context.Context, from, to *time.Time) (*MobilePerfResult, error) {
	// Calculate the previous time window (1 hour before the query window)
	prevFrom := from.Add(-24 * time.Hour)
	prevTo := to.Add(-1 * time.Hour)

	// Build the query with trend calculations
	query := `
WITH 
    current AS (
        SELECT 
            count() AS totalRequests,
            count() / (toUInt64(toUnixTimestamp(@to)) - toUInt64(toUnixTimestamp(@from))) AS requestsPerSecond,
            countIf(Status = 0) AS totalErrors,
            countIf(Status = 0) / (toUInt64(toUnixTimestamp(@to)) - toUInt64(toUnixTimestamp(@from))) AS errorsPerSecond,
            uniqExact(UserID) AS usersImpacted,
            uniqExact(UserID) / (toUInt64(toUnixTimestamp(@to)) - toUInt64(toUnixTimestamp(@from))) AS usersImpactedPerSecond
        FROM 
            mobile_perf_data
        WHERE 
            Timestamp BETWEEN @from AND @to
    ),
    previous AS (
        SELECT 
            count() AS totalRequests,
            countIf(Status = 0) AS totalErrors,
            uniqExact(UserID) AS usersImpacted
        FROM 
            mobile_perf_data
        WHERE 
            Timestamp BETWEEN @prevFrom AND @prevTo
    )
SELECT 
    current.totalRequests,
    current.requestsPerSecond,
    multiIf(previous.totalRequests > 0, (current.totalRequests - previous.totalRequests) / toFloat64(previous.totalRequests) * 100, 0) AS requestsTrend,
    current.totalErrors,
    current.errorsPerSecond,
    multiIf(previous.totalErrors > 0, (current.totalErrors - previous.totalErrors) / toFloat64(previous.totalErrors) * 100, 0) AS errorsTrend,
    current.usersImpacted,
    current.usersImpactedPerSecond,
    multiIf(previous.usersImpacted > 0, (current.usersImpacted - previous.usersImpacted) / toFloat64(previous.usersImpacted) * 100, 0) AS usersImpactedTrend
FROM 
    current, previous`

	// Execute the query
	rows, err := c.Query(ctx, query,
		clickhouse.Named("from", *from),
		clickhouse.Named("to", *to),
		clickhouse.Named("prevFrom", prevFrom),
		clickhouse.Named("prevTo", prevTo),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process results
	var result MobilePerfResult
	if rows.Next() {
		if err := rows.Scan(
			&result.TotalRequests,
			&result.RequestsPerSecond,
			&result.RequestsTrend,
			&result.TotalErrors,
			&result.ErrorsPerSecond,
			&result.ErrorsTrend,
			&result.UsersImpacted,
			&result.UsersImpactedPerSecond,
			&result.UsersImpactedTrend,
		); err != nil {
			return nil, err
		}
	}

	return &result, nil
}
