package clickhouse

import (
	"codexray/timeseries"
	"context"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type MobileCrashesResults struct {
	TotalCrashes uint64
}

type CrashesReasonwiseOverview struct {
	CrashReason   string
	TotalCrashes  uint64
	ImpactedUsers uint64
	LastOccurance timeseries.Time
}

func (c *Client) GetMobileCrashesResults(ctx context.Context, from, to timeseries.Time) (MobileCrashesResults, error) {
	query := `
	SELECT
		count() AS TotalCrashes
	FROM mobile_crash_reports
	WHERE 
		Timestamp BETWEEN @from AND @to
	`

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	)
	if err != nil {
		return MobileCrashesResults{}, err
	}
	defer rows.Close()

	var result MobileCrashesResults
	if rows.Next() {
		if err := rows.Scan(&result.TotalCrashes); err != nil {
			return MobileCrashesResults{}, err
		}
	}

	return result, nil
}

func (c *Client) GetCrashesReasonwiseOverview(ctx context.Context, from, to timeseries.Time) ([]CrashesReasonwiseOverview, error) {
	query := `
	SELECT
		cr.CrashReason,
		count(DISTINCT cr.SessionId) AS TotalCrashes,
		count(DISTINCT msd.UserId) AS ImpactedUsers,
		max(cr.CrashTime) AS LastCrashTimestamp
	FROM mobile_crash_reports cr
	LEFT JOIN mobile_session_data msd ON cr.SessionId = msd.SessionId
	WHERE 
		cr.Timestamp BETWEEN @from AND @to
		AND msd.Timestamp BETWEEN @from AND @to
	GROUP BY 
		cr.CrashReason
	ORDER BY TotalCrashes DESC
	`

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []CrashesReasonwiseOverview
	for rows.Next() {
		var result CrashesReasonwiseOverview
		var lastCrashTimestamp int64

		if err := rows.Scan(
			&result.CrashReason,
			&result.TotalCrashes,
			&result.ImpactedUsers,
			&lastCrashTimestamp,
		); err != nil {
			return nil, err
		}

		result.LastOccurance = timeseries.Time(lastCrashTimestamp)
		results = append(results, result)
	}

	return results, nil
}
