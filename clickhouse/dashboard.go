package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ServiceMetric struct {
	ServiceName   string  `ch:"ServiceName"`
	AppType       string  `ch:"AppType"`
	Requests      float64 `ch:"requests"`
	ResponseTime  float64 `ch:"responseTime"`
	Errors        uint64  `ch:"errors"`
	AffectedUsers uint64  `ch:"affectedUsers"`
}

func (c *Client) GetEUMOverview(ctx context.Context, from, to *time.Time) ([]ServiceMetric, error) {
	query := `
SELECT 
    p.ServiceName,
    p.AppType,
    p.requests,
    p.responseTime,
    ifNull(e.errors, 0) AS errors,
    ifNull(e.affectedUsers, 0) AS affectedUsers
FROM (
    SELECT 
        ServiceName,
        AppType,
        count() AS requests,
        avg(if(AppType = 'Browser', LoadPageTime, ResponseTime)) AS responseTime
    FROM (
        SELECT ServiceName, AppType, LoadPageTime
        FROM perf_data
        WHERE (@from IS NULL OR Timestamp >= @from) 
            AND (@to IS NULL OR Timestamp <= @to)
        UNION ALL
        SELECT Service AS ServiceName, AppType, ResponseTime
        FROM mobile_perf_data
        WHERE (@from IS NULL OR Timestamp >= @from) 
            AND (@to IS NULL OR Timestamp <= @to)
    ) combined
    GROUP BY ServiceName, AppType
) p
LEFT JOIN (
    SELECT 
        ServiceName,
        count() AS errors,
        countDistinct(if(UserId != '', UserId, NULL)) AS affectedUsers
    FROM err_log_data
    WHERE (@from IS NULL OR Timestamp >= @from) 
        AND (@to IS NULL OR Timestamp <= @to)
    GROUP BY ServiceName
) e ON p.ServiceName = e.ServiceName
ORDER BY requests DESC
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

	var results []ServiceMetric
	if err := c.conn.Select(ctx, &results, query, args...); err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}

	return results, nil
}

type AppCounts struct {
	BrowserApps uint64 `ch:"browserApps"`
	MobileApps  uint64 `ch:"mobileApps"`
}

func (c *Client) GetAppCounts(ctx context.Context, from, to *time.Time) (AppCounts, error) {
	query := `
SELECT 
    sum(if(AppType = 'Browser', app_count, 0)) AS browserApps,
    sum(if(AppType = 'mobile', app_count, 0)) AS mobileApps
FROM (
    SELECT countDistinct(ServiceName) AS app_count, AppType
    FROM perf_data
    WHERE (@from IS NULL OR Timestamp >= @from)
        AND (@to IS NULL OR Timestamp <= @to)
        AND AppType = 'Browser'
    GROUP BY AppType
    UNION ALL
    SELECT countDistinct(Service) AS app_count, AppType
    FROM mobile_perf_data
    WHERE (@from IS NULL OR Timestamp >= @from)
        AND (@to IS NULL OR Timestamp <= @to)
        AND AppType = 'mobile'
    GROUP BY AppType
) combined
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

	// Execute query
	var result AppCounts
	if err := c.conn.Select(ctx, &result, query, args...); err != nil {
		return AppCounts{}, err
	}

	return result, nil
}
