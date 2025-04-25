package clickhouse

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ServiceOverview struct {
	ServiceName        string
	Pages              uint64
	AvgLoadPageTime    float64
	JsErrorPercentage  float64
	ApiErrorPercentage float64
	ImpactedUsers      uint64
	AppType            string
	Requests           uint64
}

func (c *Client) GetServiceOverviews(ctx context.Context, from, to *time.Time) ([]ServiceOverview, error) {
	// Query for performance data from perf_data table
	perfQuery := `
SELECT 
    ServiceName, 
    countDistinct(PageName) AS pages,
    avg(LoadPageTime) AS avgLoadPageTime,
    count(ServiceName) AS requests,
    AppType
FROM 
    perf_data
WHERE 
    (@from IS NULL OR Timestamp >= @from) 
    AND (@to IS NULL OR Timestamp <= @to)
GROUP BY 
    ServiceName, AppType
ORDER BY 
    pages DESC
`

	// Query for error data from err_log_data table
	errorQuery := `
SELECT 
    ServiceName, 
    round(countIf(Category = 'js') * 100.0 / count(), 2) AS jsErrorPercentage,
    round(countIf(Category = 'api') * 100.0 / count(), 2) AS apiErrorPercentage,
    countDistinct(if(UserId != '', UserId, NULL)) AS impactedUsers
FROM 
    err_log_data
WHERE 
    (@from IS NULL OR Timestamp >= @from) 
    AND (@to IS NULL OR Timestamp <= @to)
GROUP BY 
    ServiceName
`

	var args []any
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

	perfRows, err := c.executePerfQuery(ctx, perfQuery, args)
	if err != nil {
		return nil, err
	}

	errorRows, err := c.executeErrorQuery(ctx, errorQuery, args)
	if err != nil {
		return nil, err
	}

	perfData := make(map[string]*ServiceOverview)
	for _, row := range perfRows {
		perfData[row.ServiceName] = &row
	}

	for _, row := range errorRows {
		if perfData[row.ServiceName] != nil {
			perfData[row.ServiceName].JsErrorPercentage = row.JsErrorPercentage
			perfData[row.ServiceName].ApiErrorPercentage = row.ApiErrorPercentage
			perfData[row.ServiceName].ImpactedUsers = row.ImpactedUsers
		}
	}

	var results []ServiceOverview
	for _, row := range perfData {
		results = append(results, *row)
	}

	return results, nil
}

// Helper function to execute performance query
func (c *Client) executePerfQuery(ctx context.Context, query string, args []any) ([]ServiceOverview, error) {
	rows, err := c.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ServiceOverview
	for rows.Next() {
		var row ServiceOverview
		if err := rows.Scan(&row.ServiceName, &row.Pages, &row.AvgLoadPageTime, &row.Requests, &row.AppType); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, nil
}

// Helper function to execute error query
func (c *Client) executeErrorQuery(ctx context.Context, query string, args []any) ([]ServiceOverview, error) {
	rows, err := c.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ServiceOverview
	for rows.Next() {
		var row ServiceOverview
		if err := rows.Scan(&row.ServiceName, &row.JsErrorPercentage, &row.ApiErrorPercentage, &row.ImpactedUsers); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, nil
}
