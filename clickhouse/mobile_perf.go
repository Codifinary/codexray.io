package clickhouse

import (
	"codexray/timeseries"
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

func (c *Client) GetMobilePerfResults(ctx context.Context, from, to *timeseries.Time) (*MobilePerfResult, error) {
	// Convert timeseries.Time to standard time.Time for calculations
	fromTime := from.ToStandard()
	toTime := to.ToStandard()

	// Calculate the previous time window (1 hour before the query window)
	prevFromTime := fromTime.Add(-24 * time.Hour)
	prevToTime := toTime.Add(-1 * time.Hour)

	// Convert back to timeseries.Time
	prevFrom := timeseries.Time(prevFromTime.Unix())
	prevTo := timeseries.Time(prevToTime.Unix())

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
		clickhouse.Named("from", fromTime),
		clickhouse.Named("to", toTime),
		clickhouse.Named("prevFrom", prevFrom.ToStandard()),
		clickhouse.Named("prevTo", prevTo.ToStandard()),
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

func (c *Client) GetRequestsByTimeSliceChart(ctx context.Context, from, to *timeseries.Time) (*timeseries.TimeSeries, error) {
	// Use timeseries.Time directly
	tsFrom := *from
	tsTo := *to

	// Calculate the interval duration (divide the total time range into 6 equal parts)
	totalDuration := tsTo.Sub(tsFrom)
	intervalDuration := totalDuration / 6

	// Create a new TimeSeries with 6 points
	ts := timeseries.New(tsFrom, 6, intervalDuration)

	// Build the query to get request counts for each interval
	query := `
	WITH
		toUnixTimestamp(@from) as start_time,
		toUnixTimestamp(@to) as end_time,
		(end_time - start_time) / 6 as interval_seconds
	SELECT
		toUnixTimestamp(toStartOfInterval(Timestamp, interval_seconds * number + start_time, 'second')) as interval_start,
		count() as request_count
	FROM
		mobile_perf_data
	WHERE
		Timestamp BETWEEN @from AND @to
	GROUP BY
		interval_start
	ORDER BY
		interval_start
	`

	// Execute the query
	rows, err := c.Query(ctx, query,
		clickhouse.Named("from", from.ToStandard()),
		clickhouse.Named("to", to.ToStandard()),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process results and populate the time series
	for rows.Next() {
		var intervalStart int64
		var requestCount uint64

		if err := rows.Scan(&intervalStart, &requestCount); err != nil {
			return nil, err
		}

		// Set the data point in the time series
		ts.Set(timeseries.Time(intervalStart), float32(requestCount))
	}

	return ts, nil
}

func (c *Client) GetErrorRateTrendByTimeChart(ctx context.Context, from, to *timeseries.Time) (*timeseries.TimeSeries, error) {
	// Use timeseries.Time directly
	tsFrom := *from
	tsTo := *to

	// Calculate the interval duration (divide the total time range into 6 equal parts)
	totalDuration := tsTo.Sub(tsFrom)
	intervalDuration := totalDuration / 6

	// Create a new TimeSeries with 6 points
	ts := timeseries.New(tsFrom, 6, intervalDuration)

	// Build the query to get error counts for each interval
	// Error is when Status = 0 in mobile_perf_data
	query := `
	WITH
		toUnixTimestamp(@from) as start_time,
		toUnixTimestamp(@to) as end_time,
		(end_time - start_time) / 6 as interval_seconds
	SELECT
		toUnixTimestamp(toStartOfInterval(Timestamp, interval_seconds * number + start_time, 'second')) as interval_start,
		countIf(Status = 0) as error_count,
		count() as total_count,
		if(total_count > 0, error_count / total_count * 100, 0) as error_rate
	FROM
		mobile_perf_data
	WHERE
		Timestamp BETWEEN @from AND @to
	GROUP BY
		interval_start
	ORDER BY
		interval_start
	`

	// Execute the query
	rows, err := c.Query(ctx, query,
		clickhouse.Named("from", from.ToStandard()),
		clickhouse.Named("to", to.ToStandard()),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process results and populate the time series
	for rows.Next() {
		var intervalStart int64
		var errorCount uint64
		var totalCount uint64
		var errorRate float32

		if err := rows.Scan(&intervalStart, &errorCount, &totalCount, &errorRate); err != nil {
			return nil, err
		}

		// Set the data point in the time series (error rate as percentage)
		ts.Set(timeseries.Time(intervalStart), errorRate)
	}

	return ts, nil
}

func (c *Client) GetUserImptactedByErrorsByTimeChart(ctx context.Context, from, to *timeseries.Time) (*timeseries.TimeSeries, error) {
	// Use timeseries.Time directly
	tsFrom := *from
	tsTo := *to

	// Calculate the interval duration (divide the total time range into 6 equal parts)
	totalDuration := tsTo.Sub(tsFrom)
	intervalDuration := totalDuration / 6

	// Create a new TimeSeries with 6 points
	ts := timeseries.New(tsFrom, 6, intervalDuration)

	// Build the query to get unique users affected by errors for each interval
	// Error is when Status = 0 in mobile_perf_data
	query := `
	WITH
		toUnixTimestamp(@from) as start_time,
		toUnixTimestamp(@to) as end_time,
		(end_time - start_time) / 6 as interval_seconds
	SELECT
		toUnixTimestamp(toStartOfInterval(Timestamp, interval_seconds * number + start_time, 'second')) as interval_start,
		uniqExact(UserID) as unique_users_with_errors
	FROM
		mobile_perf_data
	WHERE
		Timestamp BETWEEN @from AND @to
		AND Status = 0
	GROUP BY
		interval_start
	ORDER BY
		interval_start
	`

	// Execute the query
	rows, err := c.Query(ctx, query,
		clickhouse.Named("from", from.ToStandard()),
		clickhouse.Named("to", to.ToStandard()),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process results and populate the time series
	for rows.Next() {
		var intervalStart int64
		var uniqueUsersWithErrors uint64

		if err := rows.Scan(&intervalStart, &uniqueUsersWithErrors); err != nil {
			return nil, err
		}

		// Set the data point in the time series
		ts.Set(timeseries.Time(intervalStart), float32(uniqueUsersWithErrors))
	}

	return ts, nil
}
