package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ServiceMetric struct {
	ServiceName   string
	AppType       string
	Requests      float64
	ResponseTime  float64
	Errors        uint64
	AffectedUsers uint64
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
            SUM(Requests) / NULLIF(toUnixTimestamp(@to) - toUnixTimestamp(@from), 0) AS RequestsPerSecond,
            SUM(ResponseTime) / NULLIF(toUnixTimestamp(@to) - toUnixTimestamp(@from), 0) AS ResponseTimePerSecond,
            COUNTIF(Status = 0) AS Errors,
            COUNTIF(Status = 0, DISTINCT UserID) AS AffectedUsers
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
            SUM(Requests) / NULLIF(toUnixTimestamp(@to) - toUnixTimestamp(@from), 0) AS RequestsPerSecond,
            SUM(ResponseTime) / NULLIF(toUnixTimestamp(@to) - toUnixTimestamp(@from), 0) AS ResponseTimePerSecond,
            0 AS Errors,
            0 AS AffectedUsers
        FROM perf_data
        WHERE (@from IS NULL OR Timestamp >= @from)
          AND (@to IS NULL OR Timestamp <= @to)
        GROUP BY ServiceName
        ORDER BY RequestsPerSecond DESC
        LIMIT 5
    `

	args := make([]any, 0, 2)
	if from != nil {
		args = append(args, clickhouse.Named("from", from.Format("2006-01-02 15:04:05")))
	} else {
		args = append(args, clickhouse.Named("from", nil))
	}
	if to != nil {
		args = append(args, clickhouse.Named("to", to.Format("2006-01-02 15:04:05")))
	} else {
		args = append(args, clickhouse.Named("to", nil))
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

	args := make([]any, 0, 2)
	if from != nil {
		args = append(args, clickhouse.Named("from", from.Format("2006-01-02 15:04:05")))
	} else {
		args = append(args, clickhouse.Named("from", nil))
	}
	if to != nil {
		args = append(args, clickhouse.Named("to", to.Format("2006-01-02 15:04:05")))
	} else {
		args = append(args, clickhouse.Named("to", nil))
	}

	var metrics []ServiceMetric
	if err := c.conn.Select(ctx, &metrics, query, args...); err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return metrics, nil
}

func (c *Client) GetAppCounts(ctx context.Context, from, to *time.Time) (int, int, error) {
	query := `
        SELECT 
            COUNT(DISTINCT ServiceName) AS BrowserServiceCount
        FROM perf_data
        WHERE (@from IS NULL OR Timestamp >= @from)
          AND (@to IS NULL OR Timestamp <= @to);

        SELECT 
            COUNT(DISTINCT Service) AS MobileServiceCount
        FROM mobile_perf_data
        WHERE (@from IS NULL OR Timestamp >= @from)
          AND (@to IS NULL OR Timestamp <= @to);
    `

	args := make([]any, 0, 2)
	if from != nil {
		args = append(args, clickhouse.Named("from", from.Format("2006-01-02 15:04:05")))
	} else {
		args = append(args, clickhouse.Named("from", nil))
	}
	if to != nil {
		args = append(args, clickhouse.Named("to", to.Format("2006-01-02 15:04:05")))
	} else {
		args = append(args, clickhouse.Named("to", nil))
	}

	var browserServiceCount, mobileServiceCount int
	if err := c.conn.QueryRow(ctx, query, args...).Scan(&browserServiceCount, &mobileServiceCount); err != nil {
		return 0, 0, fmt.Errorf("failed to execute query: %w", err)
	}

	return browserServiceCount, mobileServiceCount, nil
}
