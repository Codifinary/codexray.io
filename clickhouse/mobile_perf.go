package clickhouse

import (
	"codexray/model"
	"codexray/timeseries"
	"context"
	"fmt"
	"strings"
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

type MobilePerfCountrywiseOverview struct {
	Country             string
	Requests            uint64
	Errors              uint64
	ErrorRatePercentage float64
	AvgResponseTime     float64
	GeoMapColorCode     string
}

func (c *Client) GetMobilePerfResults(ctx context.Context, from, to timeseries.Time) (*MobilePerfResult, error) {

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

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var intervalStart uint64
		var requestCount uint64

		if err := rows.Scan(&intervalStart, &requestCount); err != nil {
			return nil, err
		}

		ts.Set(timeseries.Time(intervalStart/1000), float32(requestCount))
	}

	return ts, nil
}

func (c *Client) GetErrorRateTrendByTimeChart(ctx context.Context, from, to timeseries.Time, step timeseries.Duration) (*timeseries.TimeSeries, error) {
	ts := timeseries.New(from, int(to.Sub(from)/step), step)

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

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var intervalStart uint64
		var errorCount uint64
		var totalCount uint64
		var errorRate float64

		if err := rows.Scan(&intervalStart, &errorCount, &totalCount, &errorRate); err != nil {
			return nil, err
		}

		ts.Set(timeseries.Time(intervalStart/1000), float32(errorRate))
	}

	return ts, nil
}

func (c *Client) GetUserImptactedByErrorsByTimeChart(ctx context.Context, from, to timeseries.Time, step timeseries.Duration) (*timeseries.TimeSeries, error) {
	ts := timeseries.New(from, int(to.Sub(from)/step), step)

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

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var intervalStart uint64
		var uniqueUsersWithErrors uint64

		if err := rows.Scan(&intervalStart, &uniqueUsersWithErrors); err != nil {
			return nil, err
		}

		ts.Set(timeseries.Time(intervalStart/1000), float32(uniqueUsersWithErrors))
	}

	return ts, nil
}

func (c *Client) GetMobilePerfCountrywiseOverviews(ctx context.Context, from, to timeseries.Time) ([]MobilePerfCountrywiseOverview, error) {

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

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

		switch {
		case overview.ErrorRatePercentage <= 20:
			overview.GeoMapColorCode = "#FFEBCD" // Lightest orange (BlanchedAlmond)
		case overview.ErrorRatePercentage <= 40:
			overview.GeoMapColorCode = "#FFA07A" // Light orange (LightSalmon)
		case overview.ErrorRatePercentage <= 60:
			overview.GeoMapColorCode = "#FF8C00" // Medium orange (DarkOrange)
		case overview.ErrorRatePercentage <= 80:
			overview.GeoMapColorCode = "#FF4500" // Dark orange (OrangeRed)
		default: // > 80%
			overview.GeoMapColorCode = "#FF0000" // Darkest orange/red (Red)
		}

		results = append(results, overview)
	}

	return results, nil
}

func (c *Client) GetHttpResponsePerfHistogram(ctx context.Context, q SpanQuery) ([]model.HistogramBucket, error) {
	filter, filterArgs := q.RootSpansFilter()
	return c.getHttpResponsePerfHistogram(ctx, q, filter, filterArgs)
}

func (c *Client) getHttpResponsePerfHistogram(ctx context.Context, q SpanQuery, filters []string, filterArgs []any) ([]model.HistogramBucket, error) {
	step := q.Ctx.Step
	from := q.Ctx.From
	to := q.Ctx.To.Add(step)

	tsFilter := "Timestamp BETWEEN @from AND @to"
	filters = append(filters, tsFilter)
	filterArgs = append(filterArgs,
		clickhouse.Named("step", step),
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	)

	query := "SELECT toStartOfInterval(Timestamp, INTERVAL @step second) as timeInterval, count(1) as requestCount, countIf(Status = 0) as errorCount"
	query += " FROM mobile_perf_data"
	query += " WHERE " + strings.Join(filters, " AND ")
	query += " GROUP BY timeInterval"
	query += " ORDER BY timeInterval"

	rows, err := c.Query(ctx, query, filterArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var t time.Time
	var requestCount, errorCount uint64
	totalRequests := timeseries.New(from, int(to.Sub(from)/step), step)
	errorRequests := timeseries.New(from, int(to.Sub(from)/step), step)

	var maxRequestCount uint64 = 0

	for rows.Next() {
		if err = rows.Scan(&t, &requestCount, &errorCount); err != nil {
			return nil, err
		}

		ts := timeseries.Time(t.Unix())
		totalRequests.Set(ts, float32(requestCount))
		errorRequests.Set(ts, float32(errorCount))

		if requestCount > maxRequestCount {
			maxRequestCount = requestCount
		}
	}

	if maxRequestCount == 0 {
		return nil, nil
	}

	numBuckets := 10
	bucketSize := float32(maxRequestCount) / float32(numBuckets)

	res := []model.HistogramBucket{
		{TimeSeries: errorRequests},
	}

	for i := 1; i <= numBuckets; i++ {
		bucketLimit := bucketSize * float32(i)
		res = append(res, model.HistogramBucket{
			Le:         bucketLimit,
			TimeSeries: totalRequests,
		})
	}

	return res, nil
}
