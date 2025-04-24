package clickhouse

import (
	"codexray/timeseries"
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type MobileOverview struct {
	Service       string `json:"service"`
	TotalRequests uint64 `json:"totalRequests"`
	TotalErrors   uint64 `json:"totalErrors"`
	TotalUsers    uint64 `json:"totalUsers"`
	TotalSessions uint64 `json:"totalSessions"`
}

func (c *Client) GetMobileOverview(ctx context.Context, from, to timeseries.Time) ([]MobileOverview, error) {
	query := `
		WITH 
		service_requests AS (
			SELECT 
				Service,
				COUNT(*) as total_requests,
				COUNT(DISTINCT SessionId) as total_sessions,
				COUNT(DISTINCT UserID) as total_users
			FROM mobile_perf_data
			WHERE Timestamp BETWEEN @from AND @to AND Service != ''
			GROUP BY Service
		),
		service_crashes AS (
			SELECT 
				Service,
				COUNT(*) as total_crashes
			FROM mobile_crash_reports
			WHERE Timestamp BETWEEN @from AND @to AND Service != ''
			GROUP BY Service
		)
		SELECT 
			sr.Service,
			sr.total_requests,
			ifNull(sc.total_crashes, 0) as total_errors,
			sr.total_users,
			sr.total_sessions
		FROM service_requests sr
		LEFT JOIN service_crashes sc ON sr.Service = sc.Service
		ORDER BY sr.total_requests DESC
	`

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute service overview query: %w", err)
	}
	defer rows.Close()

	var results []MobileOverview
	for rows.Next() {
		var overview MobileOverview

		if err := rows.Scan(
			&overview.Service,
			&overview.TotalRequests,
			&overview.TotalErrors,
			&overview.TotalUsers,
			&overview.TotalSessions,
		); err != nil {
			return nil, fmt.Errorf("failed to scan service overview: %w", err)
		}

		results = append(results, overview)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating service overview rows: %w", err)
	}

	if len(results) == 0 {
		return []MobileOverview{}, nil
	}

	return results, nil
}
