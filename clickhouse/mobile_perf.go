package clickhouse

import (
	"codexray/timeseries"
	"context"
	"fmt"

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

type MobilePerfCountrywiseOverview struct {
	Country             string
	Requests            uint64
	Errors              uint64
	ErrorRatePercentage float64
	AvgResponseTime     float64
}

func (c *Client) GetMobilePerfResults(ctx context.Context, from, to timeseries.Time) (*MobilePerfResult, error) {
	// Convert timeseries.Time to standard time.Time for calculations
	fromTime := from.ToStandard()
	toTime := to.ToStandard()

	// Calculate the previous time window (same duration as query window, but shifted back in time)
	windowDuration := toTime.Sub(fromTime)
	prevToTime := fromTime                          // Previous window ends where current window starts
	prevFromTime := prevToTime.Add(-windowDuration) // Previous window has same duration

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
    (current.totalRequests - previous.totalRequests) / greatest(1, toFloat64(previous.totalRequests)) * 100 AS requestsTrend,
    current.totalErrors,
    current.errorsPerSecond,
    (current.totalErrors - previous.totalErrors) / greatest(1, toFloat64(previous.totalErrors)) * 100 AS errorsTrend,
    current.usersImpacted,
    current.usersImpactedPerSecond,
    (current.usersImpacted - previous.usersImpacted) / greatest(1, toFloat64(previous.usersImpacted)) * 100 AS usersImpactedTrend
FROM 
    current, previous`

	// Execute the query
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

func (c *Client) GetRequestsByTimeSliceChart(ctx context.Context, from, to timeseries.Time, step timeseries.Duration) (*timeseries.TimeSeries, error) {
	ts := timeseries.New(from, int(to.Sub(from)/step), step)

	// Build the query to get request counts for each interval
	query := fmt.Sprintf(`
	SELECT
		toUnixTimestamp(toStartOfInterval(Timestamp, INTERVAL %d SECOND)) * 1000 as interval_start,
		count() as request_count
	FROM
		mobile_perf_data
	WHERE
		Timestamp BETWEEN @from AND @to
	GROUP BY
		interval_start
	ORDER BY
		interval_start
	`, step)

	// Execute the query
	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process results and populate the time series
	for rows.Next() {
		var intervalStart uint64
		var requestCount uint64

		if err := rows.Scan(&intervalStart, &requestCount); err != nil {
			return nil, err
		}

		// Set the data point in the time series
		// Convert back from milliseconds to seconds for consistency with GetPerformanceTimeSeries
		ts.Set(timeseries.Time(intervalStart/1000), float32(requestCount))
	}

	return ts, nil
}

func (c *Client) GetErrorRateTrendByTimeChart(ctx context.Context, from, to timeseries.Time, step timeseries.Duration) (*timeseries.TimeSeries, error) {
	ts := timeseries.New(from, int(to.Sub(from)/step), step)

	// Build the query to get error counts for each interval
	// Error is when Status = 0 in mobile_perf_data
	query := fmt.Sprintf(`
	SELECT
		toUnixTimestamp(toStartOfInterval(Timestamp, INTERVAL %d SECOND)) * 1000 as interval_start,
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
	`, step)

	// Execute the query
	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process results and populate the time series
	for rows.Next() {
		var intervalStart uint64
		var errorCount uint64
		var totalCount uint64
		var errorRate float64

		if err := rows.Scan(&intervalStart, &errorCount, &totalCount, &errorRate); err != nil {
			return nil, err
		}

		// Set the data point in the time series (error rate as percentage)
		// Convert back from milliseconds to seconds for consistency with GetPerformanceTimeSeries
		ts.Set(timeseries.Time(intervalStart/1000), float32(errorRate))
	}

	return ts, nil
}

func (c *Client) GetUserImptactedByErrorsByTimeChart(ctx context.Context, from, to timeseries.Time, step timeseries.Duration) (*timeseries.TimeSeries, error) {
	ts := timeseries.New(from, int(to.Sub(from)/step), step)

	// Build the query to get unique users affected by errors for each interval
	// Error is when Status = 0 in mobile_perf_data
	query := fmt.Sprintf(`
	SELECT
		toUnixTimestamp(toStartOfInterval(Timestamp, INTERVAL %d SECOND)) * 1000 as interval_start,
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
	`, step)

	// Execute the query
	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process results and populate the time series
	for rows.Next() {
		var intervalStart uint64
		var uniqueUsersWithErrors uint64

		if err := rows.Scan(&intervalStart, &uniqueUsersWithErrors); err != nil {
			return nil, err
		}

		// Set the data point in the time series
		// Convert back from milliseconds to seconds for consistency with GetPerformanceTimeSeries
		ts.Set(timeseries.Time(intervalStart/1000), float32(uniqueUsersWithErrors))
	}

	return ts, nil
}

func (c *Client) GetMobilePerfCountrywiseOverviews(ctx context.Context, from, to timeseries.Time) ([]MobilePerfCountrywiseOverview, error) {
	// Build the query to get country-wise performance metrics
	query := `
	SELECT
		Country,
		count() as requests,
		countIf(Status = 0) as errors,
		if(count() > 0, countIf(Status = 0) / count() * 100, 0) as error_rate_percentage,
		avg(ResponseTime) as avg_response_time
	FROM
		mobile_perf_data
	WHERE
		Timestamp BETWEEN @from AND @to
		AND Country != ''
	GROUP BY
		Country
	ORDER BY
		requests DESC
	`

	// Execute the query
	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process results
	var results []MobilePerfCountrywiseOverview
	for rows.Next() {
		var overview MobilePerfCountrywiseOverview
		if err := rows.Scan(
			&overview.Country,
			&overview.Requests,
			&overview.Errors,
			&overview.ErrorRatePercentage,
			&overview.AvgResponseTime,
		); err != nil {
			return nil, err
		}
		results = append(results, overview)
	}

	return results, nil
}
