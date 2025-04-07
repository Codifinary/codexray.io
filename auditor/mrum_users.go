package auditor

import (
	"codexray/clickhouse"
	"codexray/model"
	"codexray/timeseries"
	"context"
	"fmt"
)

func GenerateMrumUsersReport(w *model.World, ch *clickhouse.Client, from, to timeseries.Time, service string) *model.AuditReport {
	report := model.NewAuditReport(nil, w.Ctx, nil, model.AuditReportMobileUsers, true)
	report.Status = model.OK

	now := timeseries.Now().Truncate(timeseries.Duration(24 * 60 * 60))
	sevenDaysAgo := now.Add(-6 * 24 * 60 * 60)
	oneDayStep := timeseries.Duration(24 * 60 * 60)

	userBreakdownData, err := ch.GetUserBreakdown(context.Background(), sevenDaysAgo, now, oneDayStep, service)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println(err)
		return report
	}

	userBreakdownChart := model.NewChart(w.Ctx, "New vs Returning Users (Last 7 Days)")

	userBreakdownChart.Ctx = timeseries.Context{
		From:    sevenDaysAgo,
		To:      now,
		Step:    oneDayStep,
		RawStep: oneDayStep,
	}

	userBreakDownWidget := &model.Widget{
		Chart: userBreakdownChart,
		Width: "100%",
	}
	report.AddWidget(userBreakDownWidget)

	userBreakdownChart.Column()

	userBreakdownChart.AddSeries("New Users", userBreakdownData["newUsers"], "#AB47BC")
	userBreakdownChart.AddSeries("Returning Users", userBreakdownData["returningUsers"], "#42A5F5")

	return report
}
