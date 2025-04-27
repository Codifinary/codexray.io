package clickhouse

import (
	"codexray/timeseries"
	"context"
	"fmt"

	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type BrowserStats struct {
	Name         string  `json:"browser_name"`
	Requests     uint64  `json:"requests"`
	ResponseTime float64 `json:"response_time"`
	Errors       uint64  `json:"errors"`
}

type BrowserDataPoint struct {
	Value uint64 `json:"value"`
	Name  string `json:"name"`
}

type PerfRow struct {
	PagePath           string
	AvgLoadPageTime    float64
	JsErrorPercentage  float64
	ApiErrorPercentage float64
	ImpactedUsers      uint64
	Requests           uint64
}
type Totals struct {
	TotalRequests uint64 `json:"totalRequests"`
	TotalErrors   uint64 `json:"totalErrors"`
}

type PagePerformanceMetrics struct {
	MedianLoadTime float64
	P90LoadTime    float64
	AvgLoadTime    float64
	UniqueUsers    uint64
	LoadCount      uint64
}

type PageExperienceScore struct {
	PageLoadingClass       string
	PageInteractivityClass string
	PageRenderingClass     string
	ResourceLoadingClass   string
	ServerPerformanceClass string
}

type rawPageStats struct {
	ServiceName                             string
	PageName                                string
	FptTime, FmpTime, DomReadyTime          float64
	LoadPageTime, TtlTime, DomAnalysisTime  float64
	ResTime, TransTime                      float64
	RedirectTime, DnsTime, TcpTime, SslTime float64
	TtfbTime                                float64
}

func (c *Client) GetPerformanceOverview(ctx context.Context, from, to *time.Time, serviceName string) ([]PerfRow, error) {
	// Query for performance data from perf_data table
	perfQuery := `
SELECT 
    PageName AS PagePath, 
    avg(LoadPageTime) AS avgLoadPageTime,
    count(PageName) AS Requests
FROM 
    perf_data
WHERE 
    Timestamp >= @from
    AND Timestamp <= @to
    AND ServiceName = @serviceName
GROUP BY 
    PageName`

	// Query for error data from err_log_data table
	errorQuery := `
SELECT 
    PagePath, 
    countIf(Category = 'js') * 100.0 / count() AS jsErrorPercentage,
    countIf(Category = 'api') * 100.0 / count() AS apiErrorPercentage,
    countDistinct(UserId) AS impactedUsers
FROM 
    err_log_data
WHERE 
    Timestamp >= @from
    AND Timestamp <= @to
    AND ServiceName = @serviceName
GROUP BY 
    PagePath`

	// Execute the queries and combine results
	var perfRows []PerfRow
	perfData := make(map[string]*PerfRow)

	// Execute performance query
	perfResults, err := c.executePerfQueryForPerformance(ctx, perfQuery, from, to, serviceName)
	if err != nil {
		return nil, err
	}
	for _, row := range perfResults {
		perfData[row.PagePath] = &row
	}

	// Execute error query
	errorResults, err := c.executeErrorQueryForPerformance(ctx, errorQuery, from, to, serviceName)
	if err != nil {
		return nil, err
	}
	for _, row := range errorResults {
		if perfData[row.PagePath] != nil {
			perfData[row.PagePath].JsErrorPercentage = row.JsErrorPercentage
			perfData[row.PagePath].ApiErrorPercentage = row.ApiErrorPercentage
			perfData[row.PagePath].ImpactedUsers = row.ImpactedUsers
		}
	}

	// Combine results into a single slice
	for _, row := range perfData {
		perfRows = append(perfRows, *row)
	}

	return perfRows, nil
}

// Helper function to execute performance query
func (c *Client) executePerfQueryForPerformance(ctx context.Context, query string, from, to *time.Time, serviceName string) ([]PerfRow, error) {
	args := []any{
		clickhouse.Named("from", *from),
		clickhouse.Named("to", *to),
		clickhouse.Named("serviceName", serviceName),
	}

	rows, err := c.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []PerfRow
	for rows.Next() {
		var row PerfRow
		if err := rows.Scan(&row.PagePath, &row.AvgLoadPageTime, &row.Requests); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, nil
}

// Helper function to execute error query
func (c *Client) executeErrorQueryForPerformance(ctx context.Context, query string, from, to *time.Time, serviceName string) ([]PerfRow, error) {
	args := []any{
		clickhouse.Named("from", *from),
		clickhouse.Named("to", *to),
		clickhouse.Named("serviceName", serviceName),
	}

	rows, err := c.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []PerfRow
	for rows.Next() {
		var row PerfRow
		if err := rows.Scan(&row.PagePath, &row.JsErrorPercentage, &row.ApiErrorPercentage, &row.ImpactedUsers); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, nil
}

func (c *Client) GetBrowserStats(ctx context.Context, serviceName string, from, to *time.Time) ([]BrowserStats, error) {
	query := `
WITH browser_requests AS (
    SELECT 
        Browser AS browser_name,
        COUNT(*) AS requests,
        AVG(ResTime) AS response_time
    FROM perf_data
    WHERE ServiceName = @serviceName
      AND Timestamp BETWEEN @from AND @to
    GROUP BY Browser
),
browser_errors AS (
    SELECT 
        Browser AS browser_name,
        COUNT(*) AS errors
    FROM err_log_data
    WHERE ServiceName = @serviceName
      AND Timestamp BETWEEN @from AND @to
    GROUP BY Browser
),
browser_metrics AS (
    SELECT 
        r.browser_name,
        r.requests,
        r.response_time,
        COALESCE(e.errors, 0) AS errors
    FROM browser_requests r
    LEFT JOIN browser_errors e ON r.browser_name = e.browser_name
),
top_browsers AS (
    SELECT 
        browser_name,
        requests,
        response_time,
        errors
    FROM browser_metrics
    ORDER BY requests DESC
    LIMIT 5
),
remaining_browsers AS (
    SELECT 
        browser_name,
        requests,
        response_time,
        errors
    FROM browser_metrics
    WHERE browser_name NOT IN (SELECT browser_name FROM top_browsers)
),
others AS (
    SELECT 
        'Others' AS browser_name,
        SUM(requests) AS requests,
        AVG(response_time) AS response_time,
        SUM(errors) AS errors
    FROM remaining_browsers
    GROUP BY 'Others'
)
SELECT * FROM top_browsers
UNION ALL
SELECT * FROM others
WHERE requests > 0 
ORDER BY requests DESC;
`

	args := []any{
		clickhouse.Named("serviceName", serviceName),
		clickhouse.DateNamed("from", *from, clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", *to, clickhouse.NanoSeconds),
	}

	rows, err := c.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []BrowserStats
	for rows.Next() {
		var stats BrowserStats
		if err := rows.Scan(&stats.Name, &stats.Requests, &stats.ResponseTime, &stats.Errors); err != nil {
			return nil, err
		}
		results = append(results, stats)
	}
	return results, nil
}
func (c *Client) GetPerformanceTimeSeries(ctx context.Context, serviceName, pageName string, from, to timeseries.Time, step timeseries.Duration) (map[string]*timeseries.TimeSeries, error) {
	perfQuery := fmt.Sprintf(`
    SELECT
        toUnixTimestamp(toStartOfInterval(Timestamp, INTERVAL %d SECOND)) * 1000 AS ts,
        avg(LoadPageTime) AS loadTime,
        avg(ResTime) AS responseTime,
        avg(DnsTime) AS dnsTime,
        avg(TcpTime) AS tcpTime,
        avg(SslTime) AS sslTime,
        avg(DomAnalysisTime) AS domAnalysisTime,
        avg(DomReadyTime) AS domReadyTime,
        avg(FirstPackTime) AS firstPackTime,
        avg(FmpTime) AS fmpTime,
        avg(FptTime) AS fptTime,
        avg(RedirectTime) AS redirectTime,
        avg(TtfbTime) AS ttfbTime,
        avg(TtlTime) AS ttlTime,
        avg(TransTime) AS transTime,
        count(PageName) AS requests
    FROM
        perf_data
    WHERE
        ServiceName = @serviceName
        AND PageName = @pageName
        AND Timestamp BETWEEN @from AND @to
    GROUP BY
        ts
    ORDER BY
        ts ASC;
    `, step)

	errorQuery := fmt.Sprintf(`
    SELECT
        toUnixTimestamp(toStartOfInterval(Timestamp, INTERVAL %d SECOND)) * 1000 AS ts,
        sum(CASE WHEN Category = 'js' THEN 1 ELSE 0 END) AS jsErrors,
        sum(CASE WHEN Category = 'api' THEN 1 ELSE 0 END) AS apiErrors,
        countDistinct(UserId) AS usersImpacted
    FROM
        err_log_data
    WHERE
        ServiceName = @serviceName
        AND PagePath = @pageName
        AND Timestamp BETWEEN @from AND @to
    GROUP BY
        ts
    ORDER BY
        ts ASC;
    `, step)

	args := []any{
		clickhouse.Named("serviceName", serviceName),
		clickhouse.Named("pageName", pageName),
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	}

	perfRows, err := c.Query(ctx, perfQuery, args...)
	if err != nil {
		return nil, err
	}
	defer perfRows.Close()

	errorRows, err := c.Query(ctx, errorQuery, args...)
	if err != nil {
		return nil, err
	}
	defer errorRows.Close()

	loadTimeSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	responseTimeSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	dnsTimeSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	tcpTimeSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	sslTimeSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	domAnalysisTimeSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	domReadyTimeSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	firstPackTimeSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	fmpTimeSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	fptTimeSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	redirectTimeSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	ttfbTimeSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	ttlTimeSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	transTimeSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	requestsSeries := timeseries.New(from, int(to.Sub(from)/step), step)

	for perfRows.Next() {
		var timestamp uint64
		var loadTime, responseTime, dnsTime, tcpTime, sslTime, domAnalysisTime, domReadyTime, firstPackTime, fmpTime, fptTime, redirectTime, ttfbTime, ttlTime, transTime float64
		var requests uint64
		if err := perfRows.Scan(&timestamp, &loadTime, &responseTime, &dnsTime, &tcpTime, &sslTime, &domAnalysisTime, &domReadyTime, &firstPackTime, &fmpTime, &fptTime, &redirectTime, &ttfbTime, &ttlTime, &transTime, &requests); err != nil {
			return nil, err
		}
		ts := timeseries.Time(timestamp / 1000)
		loadTimeSeries.Set(ts, float32(loadTime))
		responseTimeSeries.Set(ts, float32(responseTime))
		dnsTimeSeries.Set(ts, float32(dnsTime))
		tcpTimeSeries.Set(ts, float32(tcpTime))
		sslTimeSeries.Set(ts, float32(sslTime))
		domAnalysisTimeSeries.Set(ts, float32(domAnalysisTime))
		domReadyTimeSeries.Set(ts, float32(domReadyTime))
		firstPackTimeSeries.Set(ts, float32(firstPackTime))
		fmpTimeSeries.Set(ts, float32(fmpTime))
		fptTimeSeries.Set(ts, float32(fptTime))
		redirectTimeSeries.Set(ts, float32(redirectTime))
		ttfbTimeSeries.Set(ts, float32(ttfbTime))
		ttlTimeSeries.Set(ts, float32(ttlTime))
		transTimeSeries.Set(ts, float32(transTime))
		requestsSeries.Set(ts, float32(requests))
	}

	jsErrorsSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	apiErrorsSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	usersImpactedSeries := timeseries.New(from, int(to.Sub(from)/step), step)

	for errorRows.Next() {
		var timestamp uint64
		var jsErrors, apiErrors, usersImpacted uint64
		if err := errorRows.Scan(&timestamp, &jsErrors, &apiErrors, &usersImpacted); err != nil {
			return nil, err
		}
		ts := timeseries.Time(timestamp / 1000)
		jsErrorsSeries.Set(ts, float32(jsErrors))
		apiErrorsSeries.Set(ts, float32(apiErrors))
		usersImpactedSeries.Set(ts, float32(usersImpacted))
	}

	return map[string]*timeseries.TimeSeries{
		"loadTime":        loadTimeSeries,
		"responseTime":    responseTimeSeries,
		"dnsTime":         dnsTimeSeries,
		"tcpTime":         tcpTimeSeries,
		"sslTime":         sslTimeSeries,
		"domAnalysisTime": domAnalysisTimeSeries,
		"domReadyTime":    domReadyTimeSeries,
		"firstPackTime":   firstPackTimeSeries,
		"fmpTime":         fmpTimeSeries,
		"fptTime":         fptTimeSeries,
		"redirectTime":    redirectTimeSeries,
		"ttfbTime":        ttfbTimeSeries,
		"ttlTime":         ttlTimeSeries,
		"transTime":       transTimeSeries,
		"requests":        requestsSeries,
		"jsErrors":        jsErrorsSeries,
		"apiErrors":       apiErrorsSeries,
		"usersImpacted":   usersImpactedSeries,
	}, nil
}

func (c *Client) GetTotalRequests(ctx context.Context, from, to *time.Time, serviceName, pagePath string) (uint64, error) {
	query := `
		SELECT
			count(*) AS totalRequests
		FROM
			perf_data p
		LEFT JOIN
			err_log_data e
		ON
			p.PageName = e.PagePath
			AND p.ServiceName = e.ServiceName`

	var filters []string
	var args []any

	// Add time range filters
	if from != nil {
		filters = append(filters, "p.Timestamp >= @from")
		args = append(args, clickhouse.Named("from", *from))
	}
	if to != nil {
		filters = append(filters, "p.Timestamp <= @to")
		args = append(args, clickhouse.Named("to", *to))
	}

	// Add service name and page path filters
	if serviceName != "" {
		filters = append(filters, "p.ServiceName = @serviceName")
		args = append(args, clickhouse.Named("serviceName", serviceName))
	}
	if pagePath != "" {
		filters = append(filters, "p.PageName = @pagePath OR e.PagePath = @pagePath")
		args = append(args, clickhouse.Named("pagePath", pagePath))
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	// Execute the query for total requests
	row := c.conn.QueryRow(ctx, query, args...)
	var totalRequests uint64
	if err := row.Scan(&totalRequests); err != nil {
		return 0, err
	}

	return totalRequests, nil
}

func (c *Client) GetTotalErrors(ctx context.Context, from, to *time.Time, serviceName, pagePath string) (uint64, error) {
	query := `
		SELECT
			count(*) AS totalErrors
		FROM
			err_log_data e`
	var filters []string
	var args []any

	// Add time range filters
	if from != nil {
		filters = append(filters, "e.Timestamp >= @from")
		args = append(args, clickhouse.Named("from", *from))
	}
	if to != nil {
		filters = append(filters, "e.Timestamp <= @to")
		args = append(args, clickhouse.Named("to", *to))
	}

	// Add service name and page path filters
	if serviceName != "" {
		filters = append(filters, "e.ServiceName = @serviceName")
		args = append(args, clickhouse.Named("serviceName", serviceName))
	}
	if pagePath != "" {
		filters = append(filters, "e.PagePath = @pagePath")
		args = append(args, clickhouse.Named("pagePath", pagePath))
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	// Execute the query for total errors
	row := c.conn.QueryRow(ctx, query, args...)
	var totalErrors uint64
	if err := row.Scan(&totalErrors); err != nil {
		return 0, err
	}

	return totalErrors, nil
}

func (c *Client) GetTopBrowser(ctx context.Context, from, to time.Time) ([]BrowserDataPoint, error) {
	query := `
        SELECT
            Browser AS name,
            count(*) AS value
        FROM
            perf_data
        WHERE
            Timestamp >= @from AND Timestamp <= @to
        GROUP BY
            Browser
        ORDER BY
            value DESC
        LIMIT 5`

	args := []any{
		clickhouse.Named("from", from),
		clickhouse.Named("to", to),
	}

	rows, err := c.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topBrowsers []BrowserDataPoint
	for rows.Next() {
		var dp BrowserDataPoint
		if err := rows.Scan(&dp.Name, &dp.Value); err != nil {
			return nil, err
		}
		topBrowsers = append(topBrowsers, dp)
	}

	return topBrowsers, nil
}

func (c *Client) GetLoadTimeSeries(ctx context.Context, serviceName string, from, to timeseries.Time, step timeseries.Duration) (*timeseries.TimeSeries, error) {
	query := fmt.Sprintf(`
        SELECT
            toUnixTimestamp(toStartOfInterval(p.Timestamp, INTERVAL %d SECOND)) * 1000 AS ts,
            avg(p.LoadPageTime) AS loadTime
        FROM
            perf_data p
        WHERE
            p.ServiceName = @serviceName
            AND p.Timestamp BETWEEN @from AND @to
        GROUP BY
            ts;
    `, step)

	args := []any{
		clickhouse.Named("serviceName", serviceName),
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	}

	rows, err := c.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	loadTimeSeries := timeseries.New(from, int(to.Sub(from)/step), step)

	for rows.Next() {
		var timestamp uint64
		var loadTime float64
		if err := rows.Scan(&timestamp, &loadTime); err != nil {
			return nil, err
		}
		ts := timeseries.Time(timestamp / 1000)
		loadTimeSeries.Set(ts, float32(loadTime))
	}

	return loadTimeSeries, nil
}

func (c *Client) GetResponseTimeSeries(ctx context.Context, serviceName string, from, to timeseries.Time, step timeseries.Duration) (*timeseries.TimeSeries, error) {
	query := fmt.Sprintf(`
        SELECT
            toUnixTimestamp(toStartOfInterval(p.Timestamp, INTERVAL %d SECOND)) * 1000 AS ts,
            avg(p.ResTime) AS responseTime
        FROM
            perf_data p
        WHERE
            p.ServiceName = @serviceName
            AND p.Timestamp BETWEEN @from AND @to
        GROUP BY
            ts;
    `, step)

	args := []any{
		clickhouse.Named("serviceName", serviceName),
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	}

	rows, err := c.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	responseTimeSeries := timeseries.New(from, int(to.Sub(from)/step), step)

	for rows.Next() {
		var timestamp uint64
		var responseTime float64
		if err := rows.Scan(&timestamp, &responseTime); err != nil {
			return nil, err
		}
		ts := timeseries.Time(timestamp / 1000)
		responseTimeSeries.Set(ts, float32(responseTime))
	}

	return responseTimeSeries, nil
}

func (c *Client) GetErrorTimeSeries(ctx context.Context, serviceName string, from, to timeseries.Time, step timeseries.Duration) (*timeseries.TimeSeries, error) {
	query := fmt.Sprintf(`
        SELECT
            toUnixTimestamp(toStartOfInterval(e.Timestamp, INTERVAL %d SECOND)) * 1000 AS ts,
            count(*) AS totalErrors
        FROM
            err_log_data e
        WHERE
            e.ServiceName = @serviceName
            AND e.Timestamp BETWEEN @from AND @to
        GROUP BY
            ts;
    `, step)

	args := []any{
		clickhouse.Named("serviceName", serviceName),
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	}

	rows, err := c.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	errorTimeSeries := timeseries.New(from, int(to.Sub(from)/step), step)

	for rows.Next() {
		var timestamp uint64
		var totalErrors uint64
		if err := rows.Scan(&timestamp, &totalErrors); err != nil {
			return nil, err
		}
		ts := timeseries.Time(timestamp / 1000)
		errorTimeSeries.Set(ts, float32(totalErrors))
	}

	return errorTimeSeries, nil
}

func (c *Client) GetUsersImpactedTimeSeries(ctx context.Context, serviceName string, from, to timeseries.Time, step timeseries.Duration) (*timeseries.TimeSeries, error) {
	query := fmt.Sprintf(`
        SELECT
            toUnixTimestamp(toStartOfInterval(e.Timestamp, INTERVAL %d SECOND)) * 1000 AS ts,
            countDistinct(e.UserId) AS usersImpacted
        FROM
            err_log_data e
        WHERE
            e.ServiceName = @serviceName
            AND e.Timestamp BETWEEN @from AND @to
        GROUP BY
            ts;
    `, step)

	args := []any{
		clickhouse.Named("serviceName", serviceName),
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	}

	rows, err := c.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usersImpactedTimeSeries := timeseries.New(from, int(to.Sub(from)/step), step)

	for rows.Next() {
		var timestamp uint64
		var usersImpacted uint64
		if err := rows.Scan(&timestamp, &usersImpacted); err != nil {
			return nil, err
		}
		ts := timeseries.Time(timestamp / 1000)
		usersImpactedTimeSeries.Set(ts, float32(usersImpacted))
	}

	return usersImpactedTimeSeries, nil
}

func (c *Client) GetErrorAndUsersImpactedSeries(ctx context.Context, serviceName string, from, to timeseries.Time, step timeseries.Duration) (map[string]*timeseries.TimeSeries, error) {
	perfQuery := fmt.Sprintf(`
    SELECT
        toUnixTimestamp(toStartOfInterval(Timestamp, INTERVAL %d SECOND)) * 1000 AS ts,
        countDistinct(UserId) AS totalUsers
    FROM
        perf_data
    WHERE
        ServiceName = @serviceName
        AND Timestamp BETWEEN @from AND @to
    GROUP BY
        ts
    ORDER BY
        ts ASC;
    `, step)

	errorQuery := fmt.Sprintf(`
    SELECT
        toUnixTimestamp(toStartOfInterval(Timestamp, INTERVAL %d SECOND)) * 1000 AS ts,
        sum(CASE WHEN Category = 'js' THEN 1 ELSE 0 END) AS jsErrors,
        sum(CASE WHEN Category = 'api' THEN 1 ELSE 0 END) AS apiErrors,
        countDistinct(UserId) AS usersImpacted
    FROM
        err_log_data
    WHERE
        ServiceName = @serviceName
        AND Timestamp BETWEEN @from AND @to
    GROUP BY
        ts
    ORDER BY
        ts ASC;
    `, step)

	args := []any{
		clickhouse.Named("serviceName", serviceName),
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	}

	perfRows, err := c.Query(ctx, perfQuery, args...)
	if err != nil {
		return nil, err
	}
	defer perfRows.Close()

	errorRows, err := c.Query(ctx, errorQuery, args...)
	if err != nil {
		return nil, err
	}
	defer errorRows.Close()

	jsErrorsSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	apiErrorsSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	usersImpactedSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	totalUsersSeries := timeseries.New(from, int(to.Sub(from)/step), step)

	totalUsersMap := make(map[uint64]uint64)
	for perfRows.Next() {
		var timestamp uint64
		var totalUsers uint64
		if err := perfRows.Scan(&timestamp, &totalUsers); err != nil {
			return nil, err
		}
		totalUsersMap[timestamp] = totalUsers
		ts := timeseries.Time(timestamp / 1000)
		totalUsersSeries.Set(ts, float32(totalUsers))
	}

	for errorRows.Next() {
		var timestamp uint64
		var jsErrors, apiErrors, usersImpacted uint64
		if err := errorRows.Scan(&timestamp, &jsErrors, &apiErrors, &usersImpacted); err != nil {
			return nil, err
		}
		ts := timeseries.Time(timestamp / 1000)
		jsErrorsSeries.Set(ts, float32(jsErrors))
		apiErrorsSeries.Set(ts, float32(apiErrors))
		usersImpactedSeries.Set(ts, float32(usersImpacted))

		if _, exists := totalUsersMap[timestamp]; !exists {
			totalUsersSeries.Set(ts, 0)
		}
	}

	return map[string]*timeseries.TimeSeries{
		"jsErrors":      jsErrorsSeries,
		"apiErrors":     apiErrorsSeries,
		"usersImpacted": usersImpactedSeries,
		"totalUsers":    totalUsersSeries,
	}, nil
}

func (c *Client) GetPagePerformanceMetrics(ctx context.Context, from, to *time.Time, serviceName, pageName string) (PagePerformanceMetrics, error) {
	query := `
    SELECT
        quantile(0.5)(p.LoadPageTime) AS medianLoadTime,
        quantile(0.9)(p.LoadPageTime) AS p90LoadTime,
        avg(p.LoadPageTime) AS avgLoadTime,
        countDistinct(p.UserId) AS uniqueUsers,
        count(*) AS loadCount
    FROM
        perf_data p`

	// Filters
	var filters []string
	var args []any
	if from != nil {
		filters = append(filters, "p.Timestamp >= @from")
		args = append(args, clickhouse.Named("from", *from))
	}
	if to != nil {
		filters = append(filters, "p.Timestamp <= @to")
		args = append(args, clickhouse.Named("to", *to))
	}
	if serviceName != "" {
		filters = append(filters, "p.ServiceName = @serviceName")
		args = append(args, clickhouse.Named("serviceName", serviceName))
	}
	if pageName != "" {
		filters = append(filters, "p.PageName = @pageName")
		args = append(args, clickhouse.Named("pageName", pageName))
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	query += `
    GROUP BY
        p.ServiceName, p.PageName
    ORDER BY
        medianLoadTime DESC`

	rows, err := c.Query(ctx, query, args...)
	if err != nil {
		return PagePerformanceMetrics{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var row PagePerformanceMetrics
		if err := rows.Scan(
			&row.MedianLoadTime,
			&row.P90LoadTime,
			&row.AvgLoadTime,
			&row.UniqueUsers,
			&row.LoadCount,
		); err != nil {
			return PagePerformanceMetrics{}, err
		}
		return row, nil
	}

	return PagePerformanceMetrics{}, fmt.Errorf("no data found")
}

var metricRanges = map[string]struct{ Min, Max float64 }{
	"FptTime":         {1.8, 3.5},
	"FmpTime":         {2.0, 4.0},
	"DomReadyTime":    {2.0, 4.5},
	"LoadPageTime":    {2.0, 4.5},
	"TtlTime":         {3.0, 6.0},
	"DomAnalysisTime": {0.5, 1.5},
	"ResTime":         {1.0, 3.0},
	"TransTime":       {1.0, 3.0},
	"RedirectTime":    {0.1, 0.5},
	"DnsTime":         {0.1, 0.5},
	"TcpTime":         {0.1, 0.5},
	"SslTime":         {0.1, 0.5},
	"TtfbTime":        {0.2, 0.8},
}

func classify(score float64) string {
	switch {
	case score > 0.6:
		return "Good"
	case score >= 0.3:
		return "Moderate"
	default:
		return "Poor"
	}
}

func normalize(value, min, max float64) float64 {
	if value <= min {
		return 1.0
	}
	if value >= max {
		return 0.0
	}
	return (max - value) / (max - min) // Invert so lower times = higher scores
}

func (c *Client) GetPageExperienceScores(ctx context.Context, from, to *time.Time, serviceName, pageName string) (PageExperienceScore, error) {
	query := `
    SELECT
        avg(FptTime),
        avg(FmpTime),
        avg(DomReadyTime),
        avg(LoadPageTime),
        avg(TtlTime),
        avg(DomAnalysisTime),
        avg(ResTime),
        avg(TransTime),
        avg(RedirectTime),
        avg(DnsTime),
        avg(TcpTime),
        avg(SslTime),
        avg(TtfbTime)
    FROM perf_data p
    `
	var filters []string
	var args []any
	if from != nil {
		filters = append(filters, "p.Timestamp >= @from")
		args = append(args, clickhouse.Named("from", *from))
	}
	if to != nil {
		filters = append(filters, "p.Timestamp <= @to")
		args = append(args, clickhouse.Named("to", *to))
	}
	if serviceName != "" {
		filters = append(filters, "p.ServiceName = @serviceName")
		args = append(args, clickhouse.Named("serviceName", serviceName))
	}
	if pageName != "" {
		filters = append(filters, "p.PageName = @pageName")
		args = append(args, clickhouse.Named("pageName", pageName))
	}
	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	rows, err := c.Query(ctx, query, args...)
	if err != nil {
		return PageExperienceScore{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var r rawPageStats
		err := rows.Scan(
			&r.FptTime, &r.FmpTime, &r.DomReadyTime,
			&r.LoadPageTime, &r.TtlTime, &r.DomAnalysisTime,
			&r.ResTime, &r.TransTime,
			&r.RedirectTime, &r.DnsTime, &r.TcpTime, &r.SslTime,
			&r.TtfbTime,
		)
		if err != nil {
			return PageExperienceScore{}, err
		}
		n := func(field string, val float64) float64 {
			r := metricRanges[field]
			return normalize(val, r.Min, r.Max)
		}

		// Weighted averages with rationale-based weights
		pl := 0.2*n("FptTime", r.FptTime) + 0.2*n("FmpTime", r.FmpTime) +
			0.2*n("DomReadyTime", r.DomReadyTime) + 0.3*n("LoadPageTime", r.LoadPageTime) +
			0.1*n("TtlTime", r.TtlTime)
		pi := 0.4*n("DomReadyTime", r.DomReadyTime) + 0.3*n("FmpTime", r.FmpTime) +
			0.3*n("LoadPageTime", r.LoadPageTime)
		pr := 0.3*n("FptTime", r.FptTime) + 0.3*n("FmpTime", r.FmpTime) +
			0.4*n("DomAnalysisTime", r.DomAnalysisTime)
		rl := 0.5*n("ResTime", r.ResTime) + 0.5*n("TransTime", r.TransTime)
		sp := 0.15*n("RedirectTime", r.RedirectTime) + 0.2*n("DnsTime", r.DnsTime) +
			0.2*n("TcpTime", r.TcpTime) + 0.2*n("SslTime", r.SslTime) +
			0.25*n("TtfbTime", r.TtfbTime)
		return PageExperienceScore{
			PageLoadingClass:       classify(pl),
			PageInteractivityClass: classify(pi),
			PageRenderingClass:     classify(pr),
			ResourceLoadingClass:   classify(rl),
			ServerPerformanceClass: classify(sp),
		}, nil
	}

	return PageExperienceScore{}, fmt.Errorf("no data found")
}
