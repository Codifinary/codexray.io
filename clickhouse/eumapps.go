package clickhouse

import (
	"context"
	"time"
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
    (? IS NULL OR Timestamp >= parseDateTimeBestEffort(?)) 
    AND (? IS NULL OR Timestamp <= parseDateTimeBestEffort(?))
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
    (? IS NULL OR Timestamp >= parseDateTimeBestEffort(?)) 
    AND (? IS NULL OR Timestamp <= parseDateTimeBestEffort(?))
GROUP BY 
    ServiceName
`

	// Format time values or pass nil
	var fromStr, toStr interface{}
	if from != nil {
		fromStr = from.Format("2006-01-02 15:04:05") // ClickHouse format
	}
	if to != nil {
		toStr = to.Format("2006-01-02 15:04:05")
	}

	args := []any{
		fromStr, fromStr,
		toStr, toStr,
	}

	// Execute performance query
	perfRows, err := c.executePerfQuery(ctx, perfQuery, args)
	if err != nil {
		return nil, err
	}

	// Execute error query
	errorRows, err := c.executeErrorQuery(ctx, errorQuery, args)
	if err != nil {
		return nil, err
	}

	// Combine results
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

	// Convert map to slice
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
