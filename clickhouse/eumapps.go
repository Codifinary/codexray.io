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
	Browser            string
}

func (c *Client) GetServiceOverviews(ctx context.Context, from, to *time.Time) ([]ServiceOverview, error) {
	query := `
SELECT 
    p.ServiceName, 
    countDistinct(p.PageName) AS pages,
    avg(p.LoadPageTime) AS avgLoadPageTime,
    round(countIf(e.Category = 'js') * 100.0 / count(), 2) AS jsErrorPercentage,
    round(countIf(e.Category = 'api') * 100.0 / count(), 2) AS apiErrorPercentage,
    countDistinct(if(e.UserId != '', e.UserId, NULL)) AS impactedUsers
	p.Browser AS Browser
FROM 
    perf_data p
LEFT JOIN 
    err_log_data e 
ON 
    p.PageName = e.PagePath AND p.ServiceName = e.ServiceName
WHERE 
    (? IS NULL OR p.Timestamp >= parseDateTimeBestEffort(?)) 
    AND (? IS NULL OR p.Timestamp <= parseDateTimeBestEffort(?))
GROUP BY 
    p.ServiceName
ORDER BY 
    pages DESC
`

	// Format time values or pass nil
	var fromStr, toStr interface{}
	if from != nil {
		fromStr = from.Format("2006-01-02 15:04:05") //clickHouse  format
	}
	if to != nil {
		toStr = to.Format("2006-01-02 15:04:05")
	}

	args := []any{
		fromStr, fromStr,
		toStr, toStr,
	}

	rows, err := c.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ServiceOverview
	for rows.Next() {
		var row ServiceOverview
		if err := rows.Scan(&row.ServiceName, &row.Pages, &row.AvgLoadPageTime, &row.JsErrorPercentage, &row.ApiErrorPercentage, &row.ImpactedUsers, &row.Browser); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, nil
}
