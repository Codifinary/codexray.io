package auditor

import (
	"codexray/clickhouse"
	"codexray/model"
	"codexray/timeseries"
	"context"
	"fmt"
)

func GenerateMrumUsersReport(w *model.World, ch *clickhouse.Client, from, to timeseries.Time) *model.AuditReport {
	report := model.NewAuditReport(nil, w.Ctx, nil, model.AuditReportMobileUsers, true)
	report.Status = model.OK

	now := timeseries.Now().Truncate(timeseries.Duration(24 * 60 * 60))
	sevenDaysAgo := now.Add(-6 * 24 * 60 * 60)
	oneDayStep := timeseries.Duration(24 * 60 * 60)

	userBreakdownData, err := ch.GetUserBreakdown(context.Background(), sevenDaysAgo, now, oneDayStep)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println(err)
		return report
	}

	userBreakdownChart := report.GetOrCreateChart("New vs Returning Users (Last 7 Days)", nil)

	userBreakdownChart.Ctx = timeseries.Context{
		From:    sevenDaysAgo,
		To:      now,
		Step:    oneDayStep,
		RawStep: oneDayStep,
	}

	userBreakdownChart.Column()

	userBreakdownChart.AddSeries("New Users", userBreakdownData["newUsers"], "#AB47BC")
	userBreakdownChart.AddSeries("Returning Users", userBreakdownData["returningUsers"], "#42A5F5")

	return report
}
