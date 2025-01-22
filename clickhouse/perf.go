package clickhouse

import (
	"context"
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
	Browser            string
}

func (c *Client) GetPerformanceOverview(ctx context.Context, from, to *time.Time) ([]PerfRow, error) {
	// Build the base query
	query := `
SELECT 
    p.PageName AS PagePath, 
    avg(p.LoadPageTime) AS avgLoadPageTime,
    countIf(e.Category = 'js') * 100.0 / count() AS jsErrorPercentage,
    countIf(e.Category = 'api') * 100.0 / count() AS apiErrorPercentage,
    countDistinct(e.UserId) AS impactedUsers
FROM 
    perf_data p
LEFT JOIN 
    err_log_data e 
ON 
    p.PageName = e.PagePath`

	// Conditionally add time range filtering
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
		if err := rows.Scan(&row.PagePath, &row.AvgLoadPageTime, &row.JsErrorPercentage, &row.ApiErrorPercentage, &row.ImpactedUsers); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, nil
}
