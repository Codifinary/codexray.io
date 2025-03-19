package clickhouse

import (
	"codexray/timeseries"
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type MobileUserResult struct {
	TotalUsers        uint64
	NewUsers          uint64
	ReturningUsers    uint64
	DailyActiveUsers  uint64
	WeeklyActiveUsers uint64
	DailyTrend        float64
}

func (c *Client) GetMobileUserResults(ctx context.Context, from, to timeseries.Time) (*MobileUserResult, error) {
	q := `
		WITH 
			active_users AS (
				SELECT DISTINCT UserId
				FROM mobile_event_data
				WHERE Timestamp BETWEEN @from AND @to
			),
			new_registrations AS (
				SELECT DISTINCT UserId
				FROM mobile_user_registration
				WHERE RegistrationTime BETWEEN @from AND @to
			),
			new_users AS (
				SELECT count(DISTINCT a.UserId) as count
				FROM active_users a
				INNER JOIN new_registrations n ON a.UserId = n.UserId
			),
			current_daily AS (
				SELECT count(DISTINCT UserId) as count
				FROM mobile_event_data
				WHERE Timestamp >= now() - INTERVAL 24 HOUR
			),
			previous_daily AS (
				SELECT count(DISTINCT UserId) as count
				FROM mobile_event_data
				WHERE Timestamp BETWEEN now() - INTERVAL 48 HOUR AND now() - INTERVAL 24 HOUR
			),
			weekly_active AS (
				SELECT count(DISTINCT UserId) as count
				FROM mobile_event_data
				WHERE Timestamp >= now() - INTERVAL 7 DAY
			),
			daily_trend AS (
				SELECT
					CASE
						WHEN prev.count = 0 AND curr.count > 0 THEN 100.0
						WHEN prev.count > 0 THEN (curr.count - prev.count) * 100.0 / prev.count
						ELSE 0
					END as trend
				FROM current_daily curr
				CROSS JOIN previous_daily prev
			)
		SELECT
			toUInt64(count(DISTINCT a.UserId)) as total_users,
			toUInt64((SELECT count FROM new_users)) as new_users,
			toUInt64(count(DISTINCT a.UserId) - (SELECT count FROM new_users)) as returning_users,
			toUInt64((SELECT count FROM current_daily)) as daily_active_users,
			toUInt64((SELECT count FROM weekly_active)) as weekly_active_users,
			toFloat64((SELECT trend FROM daily_trend)) as daily_trend
		FROM active_users a
	`

	rows, err := c.Query(ctx, q,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := &MobileUserResult{}

	if rows.Next() {
		if err := rows.Scan(
			&result.TotalUsers,
			&result.NewUsers,
			&result.ReturningUsers,
			&result.DailyActiveUsers,
			&result.WeeklyActiveUsers,
			&result.DailyTrend,
		); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (c *Client) GetUserTrendByTimeChart(ctx context.Context, from, to timeseries.Time, step timeseries.Duration) (*timeseries.TimeSeries, error) {
	ts := timeseries.New(from, int(to.Sub(from)/step), step)

	query := fmt.Sprintf(`
	SELECT
		toUnixTimestamp(toStartOfInterval(Timestamp, INTERVAL %d SECOND)) * 1000 as interval_start,
		count(DISTINCT UserId) as total_users
	FROM
		mobile_event_data
	WHERE
		Timestamp BETWEEN @from AND @to
	GROUP BY
		interval_start
	ORDER BY
		interval_start
	`, step)

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var intervalStart uint64
		var totalUsers uint64

		if err := rows.Scan(&intervalStart, &totalUsers); err != nil {
			return nil, err
		}

		ts.Set(timeseries.Time(intervalStart/1000), float32(totalUsers))
	}

	return ts, nil
}

func (c *Client) GetNewUsersTrendByTimeChart(ctx context.Context, from, to timeseries.Time, step timeseries.Duration) (*timeseries.TimeSeries, error) {
	ts := timeseries.New(from, int(to.Sub(from)/step), step)

	query := fmt.Sprintf(`
	SELECT
		toUnixTimestamp(toStartOfInterval(med.Timestamp, INTERVAL %d SECOND)) * 1000 as interval_start,
		count(DISTINCT med.UserId) as new_users
	FROM
		mobile_event_data med
	INNER JOIN
		mobile_user_registration mur ON med.UserId = mur.UserId
	WHERE
		med.Timestamp BETWEEN @from AND @to
		AND mur.RegistrationTime BETWEEN @from AND @to
	GROUP BY
		interval_start
	ORDER BY
		interval_start
	`, step)

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var intervalStart uint64
		var newUsers uint64

		if err := rows.Scan(&intervalStart, &newUsers); err != nil {
			return nil, err
		}

		ts.Set(timeseries.Time(intervalStart/1000), float32(newUsers))
	}

	return ts, nil
}

func (c *Client) GetReturningUsersTrendByTimeChart(ctx context.Context, from, to timeseries.Time, step timeseries.Duration) (*timeseries.TimeSeries, error) {
	ts := timeseries.New(from, int(to.Sub(from)/step), step)

	query := fmt.Sprintf(`
	WITH
		active_users AS (
			SELECT
				toUnixTimestamp(toStartOfInterval(Timestamp, INTERVAL %[1]d SECOND)) * 1000 as interval_start,
				count(DISTINCT UserId) as total_users
			FROM
				mobile_event_data
			WHERE
				Timestamp BETWEEN @from AND @to
			GROUP BY
				interval_start
		),
		new_users AS (
			SELECT
				toUnixTimestamp(toStartOfInterval(med.Timestamp, INTERVAL %[1]d SECOND)) * 1000 as interval_start,
				count(DISTINCT med.UserId) as new_users
			FROM
				mobile_event_data med
			INNER JOIN
				mobile_user_registration mur ON med.UserId = mur.UserId
			WHERE
				med.Timestamp BETWEEN @from AND @to
				AND mur.RegistrationTime BETWEEN @from AND @to
			GROUP BY
				interval_start
		)
	SELECT
		au.interval_start,
		(au.total_users - COALESCE(nu.new_users, 0)) as returning_users
	FROM
		active_users au
	LEFT JOIN
		new_users nu ON au.interval_start = nu.interval_start
	ORDER BY
		au.interval_start
	`, step)

	rows, err := c.Query(ctx, query,
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var intervalStart uint64
		var returningUsers int64

		if err := rows.Scan(&intervalStart, &returningUsers); err != nil {
			return nil, err
		}

		if returningUsers < 0 {
			returningUsers = 0
		}
		ts.Set(timeseries.Time(intervalStart/1000), float32(returningUsers))
	}

	return ts, nil
}
