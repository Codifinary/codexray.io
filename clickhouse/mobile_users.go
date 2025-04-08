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
	UserTrend           float64
	NewUserTrend        float64
	ReturningUserTrend  float64
}

type MobileUsersData struct {
	UserID    string
	Country   string
	StartTime string
	EndTime   string
}

func (c *Client) GetMobileUserResults(ctx context.Context, from, to timeseries.Time, service string) (*MobileUserResult, error) {
	fromTime := from.ToStandard()
	toTime := to.ToStandard()
	windowDuration := toTime.Sub(fromTime)
	prevToTime := fromTime
	prevFromTime := prevToTime.Add(-windowDuration)
	prevFrom := timeseries.Time(prevFromTime.Unix())
	prevTo := timeseries.Time(prevToTime.Unix())

	serviceFilter := "AND Service = @service"

	q := `
		WITH 
		active_users AS (
			SELECT DISTINCT UserId
			FROM mobile_event_data
			WHERE Timestamp BETWEEN @from AND @to
			` + serviceFilter + `
		),
		new_registrations AS (
			SELECT DISTINCT UserId
			FROM mobile_user_registration
			WHERE RegistrationTime BETWEEN @from AND @to
			` + serviceFilter + `
		),
		current_period AS (
			SELECT
				count(DISTINCT au.UserId) AS total_users,
				count(DISTINCT nr.UserId) AS new_users
			FROM active_users au
			LEFT JOIN new_registrations nr ON au.UserId = nr.UserId
		),
		prev_active_users AS (
			SELECT DISTINCT UserId
			FROM mobile_event_data
			WHERE Timestamp BETWEEN @prevFrom AND @prevTo
			` + serviceFilter + `
		),
		prev_new_registrations AS (
			SELECT DISTINCT UserId
			FROM mobile_user_registration
			WHERE RegistrationTime BETWEEN @prevFrom AND @prevTo
			` + serviceFilter + `
		),
		previous_period AS (
			SELECT
				count(DISTINCT pau.UserId) AS total_users,
				count(DISTINCT pnr.UserId) AS new_users
			FROM prev_active_users pau
			LEFT JOIN prev_new_registrations pnr ON pau.UserId = pnr.UserId
		),
		activity_windows AS (
			SELECT 
				countDistinctIf(UserId, Timestamp >= now() - INTERVAL 24 HOUR) AS daily_active_users,
				countDistinctIf(UserId, Timestamp BETWEEN now() - INTERVAL 48 HOUR AND now() - INTERVAL 24 HOUR) AS previous_daily_users,
				countDistinctIf(UserId, Timestamp >= now() - INTERVAL 7 DAY) AS weekly_active_users
			FROM mobile_event_data
			WHERE Timestamp >= now() - INTERVAL 7 DAY
			` + serviceFilter + `
		),
		user_crashes AS (
			SELECT DISTINCT med.UserId
			FROM mobile_event_data med
			JOIN mobile_crash_reports mcr ON mcr.SessionId = med.SessionId
			WHERE mcr.Timestamp BETWEEN @from AND @to
			AND med.Timestamp BETWEEN @from AND @to
			` + serviceFilter + `
		),
		crash_metrics AS (
			SELECT
				(SELECT count(DISTINCT UserId) FROM active_users) AS total_users,
				(SELECT count(DISTINCT UserId) FROM active_users) - 
				(SELECT count(DISTINCT UserId) FROM user_crashes) AS crash_free_users
		)
		SELECT
			toUInt64(cp.total_users) AS total_users,
			toUInt64(cp.new_users) AS new_users,
			toUInt64(cp.total_users - cp.new_users) AS returning_users,
			toUInt64(aw.daily_active_users) AS daily_active_users,
			toUInt64(aw.weekly_active_users) AS weekly_active_users,
			toFloat64(
				CASE
					WHEN aw.previous_daily_users = 0 AND aw.daily_active_users > 0 THEN 100.0
					WHEN aw.previous_daily_users > 0 THEN 
						(aw.daily_active_users - aw.previous_daily_users) * 100.0 / aw.previous_daily_users
					ELSE 0
				END
			) AS daily_trend,
			toFloat64(
				CASE
					WHEN cm.total_users > 0 THEN cm.crash_free_users * 100.0 / cm.total_users
					ELSE 100.0
				END
			) AS crash_free_percentage,
			least(100, greatest(-100, 
				CASE 
					WHEN pp.total_users > 0 THEN (cp.total_users - pp.total_users) * 100.0 / pp.total_users 
					WHEN cp.total_users > 0 THEN 100.0
					ELSE 0.0 
				END
			)) AS user_trend,
			least(100, greatest(-100, 
				CASE 
					WHEN pp.new_users > 0 THEN (cp.new_users - pp.new_users) * 100.0 / pp.new_users 
					WHEN cp.new_users > 0 THEN 100.0
					ELSE 0.0 
				END
			)) AS new_user_trend,
			least(100, greatest(-100, 
				CASE 
					WHEN (pp.total_users - pp.new_users) > 0 THEN 
						((cp.total_users - cp.new_users) - (pp.total_users - pp.new_users)) * 100.0 / (pp.total_users - pp.new_users) 
					WHEN (cp.total_users - cp.new_users) > 0 THEN 100.0
					ELSE 0.0 
				END
			)) AS returning_user_trend
		FROM current_period cp
		CROSS JOIN previous_period pp
		CROSS JOIN activity_windows aw
		CROSS JOIN crash_metrics cm
	`

	params := []interface{}{
		clickhouse.DateNamed("from", from.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("prevFrom", prevFrom.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("prevTo", prevTo.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.Named("service", service),
	}

	rows, err := c.Query(ctx, q, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute mobile user results query: %w", err)
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
			&result.UserTrend,
			&result.NewUserTrend,
			&result.ReturningUserTrend,
		); err != nil {
			return nil, fmt.Errorf("failed to scan mobile user results: %w", err)
		}
	}

	return result, nil
}

func (c *Client) GetUserBreakdown(ctx context.Context, from, to timeseries.Time, step timeseries.Duration, service string) (map[string]*timeseries.TimeSeries, error) {
	sevenDaysFrom := to.Add(-7 * timeseries.Day)

	newUsersSeries := timeseries.New(from, int(to.Sub(from)/step), step)
	returningUsersSeries := timeseries.New(from, int(to.Sub(from)/step), step)

	numIntervals := int(to.Sub(from) / step)
	for i := 0; i <= numIntervals; i++ {
		timePoint := from.Add(timeseries.Duration(i) * step)
		newUsersSeries.Set(timePoint, 0)
		returningUsersSeries.Set(timePoint, 0)
	}

	query := `
	WITH
		date_range AS (
			SELECT
				toUnixTimestamp(toStartOfInterval(toDateTime(@from), INTERVAL @step SECOND)) + (number * @step) as interval_start
			FROM numbers(@numPoints)
		),
		active_users AS (
			SELECT
				toUnixTimestamp(toStartOfInterval(Timestamp, INTERVAL @step SECOND)) as interval_start,
				UserId
			FROM
				mobile_event_data med
			WHERE
				Timestamp BETWEEN @from AND @to
				AND med.Service = @service
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
				toUnixTimestamp(toStartOfInterval(med.Timestamp, INTERVAL @step SECOND)) as interval_start,
				med.UserId
			FROM
				mobile_event_data med
			INNER JOIN
				mobile_user_registration mur ON med.UserId = mur.UserId
			WHERE
				med.Timestamp BETWEEN @from AND @to
				AND mur.RegistrationTime BETWEEN @from AND @to
				AND med.Service = @service
				AND mur.Service = @service
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
		dr.interval_start * 1000 as interval_start,
		COALESCE(nu.new_users, 0) as new_users,
		COALESCE(au.total_users, 0) - COALESCE(nu.new_users, 0) as returning_users
	FROM
		date_range dr
	LEFT JOIN
		active_users_agg au ON dr.interval_start = au.interval_start
	LEFT JOIN
		new_users_agg nu ON dr.interval_start = nu.interval_start
	ORDER BY
		dr.interval_start
	`

	params := []interface{}{
		clickhouse.DateNamed("from", sevenDaysFrom.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.DateNamed("to", to.ToStandard(), clickhouse.NanoSeconds),
		clickhouse.Named("service", service),
		clickhouse.Named("step", int(step)),
		clickhouse.Named("numPoints", numIntervals+1),
	}

	rows, err := c.Query(ctx, query, params...)
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

func (c *Client) GetMobileUsersData(ctx context.Context, from, to timeseries.Time, service string) ([]MobileUsersData, error) {
	serviceFilter := "AND mur.Service = @service"

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
			` + serviceFilter + `
		ORDER BY
			msd.StartTime DESC
	`

	var result []MobileUsersData
	err := c.conn.Select(ctx, &result, query,
		clickhouse.Named("from", from),
		clickhouse.Named("to", to),
		clickhouse.Named("service", service))
	if err != nil {
		return nil, fmt.Errorf("error querying mobile users data: %w", err)
	}

	return result, nil
}
