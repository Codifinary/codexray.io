package clickhouse

import (
	"codexray/timeseries"
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type MobileUserResult struct {
	TotalUsers          uint64
	NewUsers            uint64
	ReturningUsers      uint64
	DailyActiveUsers    uint64
	WeeklyActiveUsers   uint64
	DailyTrend          float64
	CrashFreePercentage float64
}

type MobileUsersData struct {
	UserID    string
	Country   string
	StartTime string
	EndTime   string
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
			user_activity_windows AS (
				SELECT 
					countDistinctIf(UserId, Timestamp >= now() - INTERVAL 24 HOUR) AS daily_active_users,
					countDistinctIf(UserId, Timestamp BETWEEN now() - INTERVAL 48 HOUR AND now() - INTERVAL 24 HOUR) AS previous_daily_users,
					countDistinctIf(UserId, Timestamp >= now() - INTERVAL 7 DAY) AS weekly_active_users
				FROM mobile_event_data
				WHERE Timestamp >= now() - INTERVAL 7 DAY
			),
			daily_trend AS (
				SELECT
					CASE
						WHEN uaw.previous_daily_users = 0 AND uaw.daily_active_users > 0 THEN 100.0
						WHEN uaw.previous_daily_users > 0 THEN (uaw.daily_active_users - uaw.previous_daily_users) * 100.0 / uaw.previous_daily_users
						ELSE 0
					END as trend
				FROM user_activity_windows uaw
			),
			user_crashes AS (
				SELECT DISTINCT med.UserId
				FROM mobile_crash_reports mcr
				JOIN mobile_event_data med ON mcr.SessionId = med.SessionId
				WHERE mcr.Timestamp BETWEEN @from AND @to
			),
			crash_free_users AS (
				SELECT count(DISTINCT a.UserId) as count
				FROM active_users a
				LEFT JOIN user_crashes uc ON a.UserId = uc.UserId
				WHERE uc.UserId IS NULL
			)
		SELECT
			toUInt64(count(DISTINCT a.UserId)) as total_users,
			toUInt64((SELECT count FROM new_users)) as new_users,
			toUInt64(count(DISTINCT a.UserId) - (SELECT count FROM new_users)) as returning_users,
			toUInt64(any(uaw.daily_active_users)) as daily_active_users,
			toUInt64(any(uaw.weekly_active_users)) as weekly_active_users,
			toFloat64((SELECT trend FROM daily_trend)) as daily_trend,
			CASE
				WHEN count(DISTINCT a.UserId) > 0 THEN
					toFloat64(COALESCE((SELECT count FROM crash_free_users), 0) * 100.0 / count(DISTINCT a.UserId))
				ELSE 100.0
			END as crash_free_percentage
		FROM active_users a
		CROSS JOIN user_activity_windows uaw
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
			&result.CrashFreePercentage,
		); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (c *Client) GetUserBreakdown(ctx context.Context, from, to timeseries.Time, step timeseries.Duration) (map[string]*timeseries.TimeSeries, error) {
	newUsersSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	returningUsersSeries := timeseries.New(from, int(to.Sub(from)/step), step)

	query := fmt.Sprintf(`
	WITH
		active_users AS (
			SELECT
				toUnixTimestamp(toStartOfInterval(Timestamp, INTERVAL %[1]d SECOND)) * 1000 as interval_start,
				UserId
			FROM
				mobile_event_data
			WHERE
				Timestamp BETWEEN @from AND @to
		),
		active_users_agg AS (
			SELECT
				interval_start,
				count(DISTINCT UserId) as total_users
			FROM
				active_users
			GROUP BY
				interval_start
		),
		new_users AS (
			SELECT
				toUnixTimestamp(toStartOfInterval(med.Timestamp, INTERVAL %[1]d SECOND)) * 1000 as interval_start,
				med.UserId
			FROM
				mobile_event_data med
			INNER JOIN
				mobile_user_registration mur ON med.UserId = mur.UserId
			WHERE
				med.Timestamp BETWEEN @from AND @to
				AND mur.RegistrationTime BETWEEN @from AND @to
		),
		new_users_agg AS (
			SELECT
				interval_start,
				count(DISTINCT UserId) as new_users
			FROM
				new_users
			GROUP BY
				interval_start
		)
	SELECT
		au.interval_start,
		COALESCE(nu.new_users, 0) as new_users,
		(au.total_users - COALESCE(nu.new_users, 0)) as returning_users
	FROM
		active_users_agg au
	LEFT JOIN
		new_users_agg nu ON au.interval_start = nu.interval_start
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
		var newUsers uint64
		var returningUsers int64

		if err := rows.Scan(&intervalStart, &newUsers, &returningUsers); err != nil {
			return nil, err
		}

		if returningUsers < 0 {
			returningUsers = 0
		}

		newUsersSeries.Set(timeseries.Time(intervalStart/1000), float32(newUsers))
		returningUsersSeries.Set(timeseries.Time(intervalStart/1000), float32(returningUsers))
	}

	return map[string]*timeseries.TimeSeries{
		"newUsers":       newUsersSeries,
		"returningUsers": returningUsersSeries,
	}, nil
}

func (c *Client) GetMobileUsersData(ctx context.Context, from, to timeseries.Time) ([]MobileUsersData, error) {
	query := `
		SELECT
			msd.UserId AS UserID,
			mur.Country AS Country,
			toString(msd.StartTime) AS StartTime,
			toString(msd.EndTime) AS EndTime
		FROM
			mobile_session_data msd
		INNER JOIN
			mobile_user_registration mur ON msd.UserId = mur.UserId
		WHERE
			msd.StartTime >= @from AND
			msd.StartTime <= @to AND
			msd.EndTime IS NOT NULL
		ORDER BY
			msd.StartTime DESC
	`

	var result []MobileUsersData
	err := c.conn.Select(ctx, &result, query, clickhouse.Named("from", from), clickhouse.Named("to", to))
	if err != nil {
		return nil, fmt.Errorf("error querying mobile users data: %w", err)
	}

	return result, nil
}
