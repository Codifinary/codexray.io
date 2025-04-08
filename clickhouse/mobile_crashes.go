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
	DeviceType     string
	StackTrace     string
	CrashTimestamp timeseries.Time
	AffectedUser   string
	Application    string
	CrashReason    string
	AppVersion     string
	MemoryUsage    int64
}

func (c *Client) GetMobileCrashesResults(ctx context.Context, from, to timeseries.Time, service string) (MobileCrashesResults, error) {
	query := `
	SELECT
		count() AS TotalCrashes
	FROM mobile_crash_reports
	WHERE 
		Timestamp BETWEEN @from AND @to
		AND Service = @service
	`

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.Named("service", service),
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

func (c *Client) GetCrashesReasonwiseOverview(ctx context.Context, from, to timeseries.Time, limit int, service string) ([]CrashesReasonwiseOverview, error) {
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
		AND cr.Service = @service
	GROUP BY 
		cr.CrashReason
	ORDER BY TotalCrashes DESC
	LIMIT @limit
	`

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.Named("limit", uint64(limit)),
		clickhouse.Named("service", service),
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

func (c *Client) GetTopDevicesByCrashCount(ctx context.Context, from, to timeseries.Time, service string) ([]struct {
	Name  string
	Value uint64
}, error) {
	query := `
	SELECT
		DeviceInfo as Name,
		count() as Value
	FROM
		mobile_crash_reports
	WHERE
		Timestamp BETWEEN @from AND @to
		AND DeviceInfo != ''
		AND Service = @service
	GROUP BY
		DeviceInfo
	ORDER BY
		Value DESC
	LIMIT 5
	`

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.Named("service", service),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		Name  string
		Value uint64
	}

	for rows.Next() {
		var result struct {
			Name  string
			Value uint64
		}
		if err := rows.Scan(&result.Name, &result.Value); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (c *Client) GetTopOSByCrashCount(ctx context.Context, from, to timeseries.Time, service string) ([]struct {
	Name  string
	Value uint64
}, error) {
	query := `
	SELECT
		Os as Name,
		count() as Value
	FROM
		mobile_crash_reports
	WHERE
		Timestamp BETWEEN @from AND @to
		AND Os != ''
		AND Service = @service
	GROUP BY
		Os
	ORDER BY
		Value DESC
	LIMIT 5
	`

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.Named("service", service),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		Name  string
		Value uint64
	}

	for rows.Next() {
		var result struct {
			Name  string
			Value uint64
		}
		if err := rows.Scan(&result.Name, &result.Value); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (c *Client) GetTopAppVersionsByCrashCount(ctx context.Context, from, to timeseries.Time, service string) ([]struct {
	Name  string
	Value uint64
}, error) {
	query := `
	SELECT
		ServiceVersion as Name,
		count() as Value
	FROM
		mobile_crash_reports
	WHERE
		Timestamp BETWEEN @from AND @to
		AND ServiceVersion != ''
		AND Service = @service
	GROUP BY
		Name
	ORDER BY
		Value DESC
	LIMIT 5
	`

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.Named("service", service),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		Name  string
		Value uint64
	}

	for rows.Next() {
		var result struct {
			Name  string
			Value uint64
		}
		if err := rows.Scan(&result.Name, &result.Value); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (c *Client) GetCrashesByDeviceTrendChart(ctx context.Context, from, to timeseries.Time, step timeseries.Duration, service string) (map[string]*timeseries.TimeSeries, error) {
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
		AND Service = @service
		AND DeviceInfo IN (
			SELECT DeviceInfo
			FROM mobile_crash_reports
			WHERE Timestamp BETWEEN @from AND @to
			AND DeviceInfo != ''
			AND Service = @service
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
		clickhouse.Named("service", service),
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
		return make(map[string]*timeseries.TimeSeries), nil
	}

	return result, nil
}

func (c *Client) GetCrashReasonData(ctx context.Context, crashReason string, from, to timeseries.Time, limit int, service string) ([]CrashReasonData, error) {
	query := `
	SELECT
		cr.UniqueId AS CrashId,
		cr.DeviceInfo AS DeviceType,
		cr.CrashStackTrace AS StackTrace,
		cr.CrashTime AS CrashTimestamp,
		msd.UserId AS AffectedUser,
		cr.Service AS Application,
		cr.CrashReason AS CrashReason,
		cr.ServiceVersion AS AppVersion,
		IFNULL(cr.MemoryUsage, 0) AS MemoryUsage
	FROM mobile_crash_reports cr
	LEFT JOIN mobile_session_data msd ON cr.SessionId = msd.SessionId
	WHERE 
		cr.Timestamp BETWEEN @from AND @to
		AND msd.Timestamp BETWEEN @from AND @to
		AND cr.CrashReason = @crashReason
		AND cr.Service = @service
	GROUP BY 
		cr.UniqueId, cr.DeviceInfo, cr.CrashStackTrace, cr.CrashTime, msd.UserId, cr.Service, cr.CrashReason, cr.ServiceVersion, cr.MemoryUsage
	ORDER BY CrashTimestamp DESC
	LIMIT @limit
	`

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.Named("crashReason", crashReason),
		clickhouse.Named("limit", uint64(limit)),
		clickhouse.Named("service", service),
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
			&result.DeviceType,
			&result.StackTrace,
			&crashTimestamp,
			&result.AffectedUser,
			&result.Application,
			&result.CrashReason,
			&result.AppVersion,
			&result.MemoryUsage,
		); err != nil {
			return nil, err
		}

		result.CrashTimestamp = timeseries.Time(crashTimestamp)
		results = append(results, result)
	}

	return results, nil
}
