package clickhouse

import (
	"codexray/timeseries"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type PerfRow struct {
	PagePath           string
	AvgLoadPageTime    float64
	JsErrorPercentage  float64
	ApiErrorPercentage float64
	ImpactedUsers      uint64
	Requests           uint64
}

func (c *Client) GetPerformanceOverview(ctx context.Context, from, to *time.Time, serviceName string) ([]PerfRow, error) {
	// Build the base query
	query := `
SELECT 
    p.PageName AS PagePath, 
    avg(p.LoadPageTime) AS avgLoadPageTime,
    countIf(e.Category = 'js') * 100.0 / count() AS jsErrorPercentage,
    countIf(e.Category = 'api') * 100.0 / count() AS apiErrorPercentage,
    countDistinct(e.UserId) AS impactedUsers,
    count(p.PageName) AS Requests
FROM 
    perf_data p
LEFT JOIN 
    err_log_data e 
ON 
    p.PageName = e.PagePath`

	// Conditionally add time range and service name filtering
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

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	query += `
GROUP BY 
    p.PageName`

	// Execute the query
	rows, err := c.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process results
	var results []PerfRow
	for rows.Next() {
		var row PerfRow
		if err := rows.Scan(&row.PagePath, &row.AvgLoadPageTime, &row.JsErrorPercentage, &row.ApiErrorPercentage, &row.ImpactedUsers, &row.Requests); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, nil
}

func (c *Client) GetPerformanceTimeSeries(ctx context.Context, serviceName, pageName string, from, to timeseries.Time, step timeseries.Duration) (map[string]*timeseries.TimeSeries, error) {
	query := fmt.Sprintf(`
    SELECT
        toUnixTimestamp(toStartOfInterval(p.Timestamp, INTERVAL %d SECOND)) * 1000 AS ts,
        avg(p.LoadPageTime) AS loadTime,
        avg(p.ResTime) AS responseTime,
        sum(CASE WHEN e.Category = 'js' THEN 1 ELSE 0 END) AS jsErrors,
        sum(CASE WHEN e.Category = 'api' THEN 1 ELSE 0 END) AS apiErrors,
        countDistinct(CASE WHEN e.UserId IS NOT NULL THEN e.UserId ELSE NULL END) AS usersImpacted
    FROM
        perf_data p
    LEFT JOIN
        err_log_data e
    ON
        p.PageName = e.PagePath
        AND p.ServiceName = e.ServiceName
    WHERE
        p.ServiceName = @serviceName
        AND p.PageName = @pageName
        AND p.Timestamp BETWEEN @from AND @to
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

	rows, err := c.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	loadTimeSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	responseTimeSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	jsErrorsSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	apiErrorsSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	usersImpactedSeries := timeseries.New(from, int(to.Sub(from)/step), step)

	for rows.Next() {
		var timestamp uint64
		var loadTime, responseTime float64
		var jsErrors, apiErrors, usersImpacted uint64
		if err := rows.Scan(&timestamp, &loadTime, &responseTime, &jsErrors, &apiErrors, &usersImpacted); err != nil {
			return nil, err
		}
		ts := timeseries.Time(timestamp / 1000)
		loadTimeSeries.Set(ts, float32(loadTime))
		responseTimeSeries.Set(ts, float32(responseTime))
		jsErrorsSeries.Set(ts, float32(jsErrors))
		apiErrorsSeries.Set(ts, float32(apiErrors))
		usersImpactedSeries.Set(ts, float32(usersImpacted))
	}

	return map[string]*timeseries.TimeSeries{
		"loadTime":      loadTimeSeries,
		"responseTime":  responseTimeSeries,
		"jsErrors":      jsErrorsSeries,
		"apiErrors":     apiErrorsSeries,
		"usersImpacted": usersImpactedSeries,
	}, nil
}
