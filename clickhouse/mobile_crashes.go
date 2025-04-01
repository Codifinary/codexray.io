package clickhouse

import (
	"codexray/timeseries"
	"context"
	"fmt"

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

type CrashReasonData struct {
	CrashId        string
	DeviceId       string
	StackTrace     string
	CrashTimestamp timeseries.Time
	AffectedUser   string
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

func (c *Client) GetCrashesReasonwiseOverview(ctx context.Context, from, to timeseries.Time, limit int) ([]CrashesReasonwiseOverview, error) {
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
	LIMIT @limit
	`

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.Named("limit", uint64(limit)),
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

func (c *Client) GetCrashesByDeviceTrendChart(ctx context.Context, from, to timeseries.Time, step timeseries.Duration) (map[string]*timeseries.TimeSeries, error) {
	optimizedQuery := fmt.Sprintf(`
	SELECT
		toUnixTimestamp(toStartOfInterval(Timestamp, INTERVAL %d SECOND)) * 1000 as interval_start,
		DeviceInfo as device,
		count() as crash_count
	FROM
		mobile_crash_reports
	WHERE
		Timestamp BETWEEN @from AND @to
		AND DeviceInfo != ''
		AND DeviceInfo IN (
			SELECT DeviceInfo
			FROM mobile_crash_reports
			WHERE Timestamp BETWEEN @from AND @to
			AND DeviceInfo != ''
			GROUP BY DeviceInfo
			ORDER BY count() DESC
			LIMIT 5
		)
	GROUP BY
		interval_start, device
	ORDER BY
		interval_start, device
	`, step)

	rows, err := c.Query(ctx, optimizedQuery,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]*timeseries.TimeSeries)

	for rows.Next() {
		var intervalStart uint64
		var device string
		var crashCount uint64

		if err := rows.Scan(&intervalStart, &device, &crashCount); err != nil {
			return nil, err
		}

		ts, exists := result[device]
		if !exists {
			ts = timeseries.New(from, int(to.Sub(from)/step), step)
			result[device] = ts
		}

		ts.Set(timeseries.Time(intervalStart/1000), float32(crashCount))
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}

func (c *Client) GetCrashReasonData(ctx context.Context, crashReason string, from, to timeseries.Time, limit int) ([]CrashReasonData, error) {
	query := `
	SELECT
		cr.UniqueId AS CrashId,
		cr.DeviceInfo AS DeviceId,
		cr.CrashStackTrace AS StackTrace,
		cr.CrashTime AS CrashTimestamp,
		msd.UserId AS AffectedUser
	FROM mobile_crash_reports cr
	LEFT JOIN mobile_session_data msd ON cr.SessionId = msd.SessionId
	WHERE 
		cr.Timestamp BETWEEN @from AND @to
		AND msd.Timestamp BETWEEN @from AND @to
		AND cr.CrashReason = @crashReason
	GROUP BY 
		cr.UniqueId, cr.DeviceInfo, cr.CrashStackTrace, cr.CrashTime, msd.UserId
	ORDER BY CrashTimestamp DESC
	LIMIT @limit
	`

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.Named("crashReason", crashReason),
		clickhouse.Named("limit", uint64(limit)),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []CrashReasonData
	for rows.Next() {
		var result CrashReasonData
		var crashTimestamp int64

		if err := rows.Scan(
			&result.CrashId,
			&result.DeviceId,
			&result.StackTrace,
			&crashTimestamp,
			&result.AffectedUser,
		); err != nil {
			return nil, err
		}

		result.CrashTimestamp = timeseries.Time(crashTimestamp)
		results = append(results, result)
	}

	return results, nil
}
