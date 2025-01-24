package clickhouse

import (
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
type ChartRow struct {
	Timestamp     uint64
	LoadTime      float64
	ResponseTime  float64
	JsErrors      uint64
	ApiErrors     uint64
	UsersImpacted uint64
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

func (c *Client) GetChartData(ctx context.Context, serviceName, pageName string, from, to *time.Time, step int64) ([]ChartRow, error) {

	query := `
SELECT
    toUnixTimestamp(toStartOfInterval(p.Timestamp, INTERVAL 60 SECOND)) * 1000 AS ts,
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
    p.ServiceName = 'health-care'
    AND p.PageName = 'home'
    AND p.Timestamp BETWEEN toDateTime64('2025-01-20 00:00:00', 6) AND toDateTime64('2025-01-25 01:00:00', 6)
GROUP BY
    ts
ORDER BY
    ts ASC;
	`
	// Fallback to default time range if `from` or `to` is nil
	defaultFrom := time.Now().Add(-24 * time.Hour)
	defaultTo := time.Now()

	formattedFrom := formatTimeForClickHouse(defaultFrom)
	formattedTo := formatTimeForClickHouse(defaultTo)

	if from != nil {
		formattedFrom = formatTimeForClickHouse(*from)
	}
	if to != nil {
		formattedTo = formatTimeForClickHouse(*to)
	}
	fmt.Println("perf details ", serviceName, pageName, formattedFrom, formattedTo)

	rows, err := c.Query(ctx, query, step, serviceName, pageName, formattedFrom, formattedTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chartRows []ChartRow
	for rows.Next() {
		var row ChartRow
		if err := rows.Scan(&row.Timestamp, &row.LoadTime, &row.ResponseTime, &row.JsErrors, &row.ApiErrors, &row.UsersImpacted); err != nil {
			return nil, err
		}
		chartRows = append(chartRows, row)
	}

	return chartRows, nil
}
func formatTimeForClickHouse(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000")
}
