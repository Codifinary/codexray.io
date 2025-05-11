package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ServiceMetric struct {
	ServiceName       string
	AppType           string
	RequestsPerSecond float64
	ResponseTime      float64
	Errors            uint64
	AffectedUsers     uint64
}

func (c *Client) GetEUMOverview(ctx context.Context, from, to *time.Time) ([]ServiceMetric, error) {
	top5, err := c.getTop5Applications(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get top 5 applications: %w", err)
	}

	var browserServiceNames []string
	for _, metric := range top5 {
		if metric.AppType == "Browser" {
			browserServiceNames = append(browserServiceNames, metric.ServiceName)
		}
	}

	browserMetrics, err := c.getBrowserErrorsAndAffectedUsers(ctx, browserServiceNames, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get browser errors and affected users: %w", err)
	}

	for i, metric := range top5 {
		if metric.AppType == "Browser" {
			for _, browserMetric := range browserMetrics {
				if metric.ServiceName == browserMetric.ServiceName {
					top5[i].Errors = browserMetric.Errors
					top5[i].AffectedUsers = browserMetric.AffectedUsers
					break
				}
			}
		}
	}

	return top5, nil
}

func (c *Client) getTop5Applications(ctx context.Context, from, to *time.Time) ([]ServiceMetric, error) {
	query := `
        SELECT 
            Service AS ServiceName,
            'Mobile' AS AppType,
            COUNT(*) / NULLIF(toUnixTimestamp(@to) - toUnixTimestamp(@from), 0) AS RequestsPerSecond,
            SUM(ResponseTime) / NULLIF(toUnixTimestamp(@to) - toUnixTimestamp(@from), 0) AS ResponseTime,
            SUM(CASE WHEN Status = 0 THEN 1 ELSE 0 END) AS Errors,
            COUNT(DISTINCT CASE WHEN Status = 0 THEN UserID END) AS AffectedUsers
        FROM mobile_perf_data
        WHERE (@from IS NULL OR Timestamp >= @from)
          AND (@to IS NULL OR Timestamp <= @to)
        GROUP BY Service
        ORDER BY RequestsPerSecond DESC
        LIMIT 5

        UNION ALL

        SELECT 
            ServiceName,
            'Browser' AS AppType,
            COUNT(*) / NULLIF(toUnixTimestamp(@to) - toUnixTimestamp(@from), 0) AS RequestsPerSecond,
            SUM(ResTime) / NULLIF(toUnixTimestamp(@to) - toUnixTimestamp(@from), 0) AS ResponseTime,
            0 AS Errors,
            0 AS AffectedUsers
        FROM perf_data
        WHERE (@from IS NULL OR Timestamp >= @from)
          AND (@to IS NULL OR Timestamp <= @to)
        GROUP BY ServiceName
        ORDER BY RequestsPerSecond DESC
        LIMIT 5
    `

	var args []any
	if from != nil {
		args = append(args, clickhouse.Named("from", *from))
	}
	if to != nil {
		args = append(args, clickhouse.Named("to", *to))
	}

	var metrics []ServiceMetric
	if err := c.conn.Select(ctx, &metrics, query, args...); err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return metrics, nil
}

func (c *Client) getBrowserErrorsAndAffectedUsers(ctx context.Context, serviceNames []string, from, to *time.Time) ([]ServiceMetric, error) {
	if len(serviceNames) == 0 {
		return nil, nil
	}

	inClause := "("
	for i, name := range serviceNames {
		if i > 0 {
			inClause += ", "
		}
		inClause += fmt.Sprintf("'%s'", name)
	}
	inClause += ")"

	query := fmt.Sprintf(`
        SELECT 
            ServiceName AS ServiceName,
            'Browser' AS AppType,
            COUNT(*) AS Errors,
            COUNT(DISTINCT UserId) AS AffectedUsers
        FROM err_log_data
        WHERE ServiceName IN %s
          AND (@from IS NULL OR Timestamp >= @from)
          AND (@to IS NULL OR Timestamp <= @to)
        GROUP BY ServiceName
    `, inClause)

	var args []any
	if from != nil {
		args = append(args, clickhouse.Named("from", *from))
	}
	if to != nil {
		args = append(args, clickhouse.Named("to", *to))
	}

	var metrics []ServiceMetric
	if err := c.conn.Select(ctx, &metrics, query, args...); err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return metrics, nil
}

func (c *Client) GetAppCounts(ctx context.Context, from, to *time.Time) (uint64, uint64, error) {
	browserQuery := `
        SELECT 
            COUNT(DISTINCT ServiceName) AS BrowserServiceCount
        FROM perf_data
        WHERE (@from IS NULL OR Timestamp >= @from)
          AND (@to IS NULL OR Timestamp <= @to)
    `

	mobileQuery := `
        SELECT 
            COUNT(DISTINCT Service) AS MobileServiceCount
        FROM mobile_perf_data
        WHERE (@from IS NULL OR Timestamp >= @from)
          AND (@to IS NULL OR Timestamp <= @to)
    `

	var args []any
	if from != nil {
		args = append(args, clickhouse.Named("from", *from))
	}
	if to != nil {
		args = append(args, clickhouse.Named("to", *to))
	}

	var browserServiceCount uint64
	if err := c.conn.QueryRow(ctx, browserQuery, args...).Scan(&browserServiceCount); err != nil {
		return 0, 0, fmt.Errorf("failed to execute browser query: %w", err)
	}

	var mobileServiceCount uint64
	if err := c.conn.QueryRow(ctx, mobileQuery, args...).Scan(&mobileServiceCount); err != nil {
		return 0, 0, fmt.Errorf("failed to execute mobile query: %w", err)
	}

	return browserServiceCount, mobileServiceCount, nil
}
