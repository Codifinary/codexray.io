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

	userTrendsGroup := report.GetOrCreateChartGroup("User Activity Trends", nil)

	userTrendData, err := ch.GetUserTrendByTimeChart(context.Background(), from, to, w.Ctx.Step)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println(err)
		return report
	}

	userTrendChart := userTrendsGroup.GetOrCreateChart("Total Active Users")
	userTrendChart.AddSeries("Total Users", userTrendData, "#FFA726")

	newUsersTrendData, err := ch.GetNewUsersTrendByTimeChart(context.Background(), from, to, w.Ctx.Step)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println(err)
		return report
	}

	newUsersTrendChart := userTrendsGroup.GetOrCreateChart("New Users")
	newUsersTrendChart.AddSeries("New Users", newUsersTrendData, "#AB47BC")

	returningUsersTrendData, err := ch.GetReturningUsersTrendByTimeChart(context.Background(), from, to, w.Ctx.Step)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println(err)
		return report
	}

	returningUsersTrendChart := userTrendsGroup.GetOrCreateChart("Returning Users")
	returningUsersTrendChart.AddSeries("Returning Users", returningUsersTrendData, "#42A5F5")

	userBreakdownGroup := report.GetOrCreateChartGroup("User Breakdown", nil)

	now := timeseries.Now()
	sevenDaysAgo := now.Add(-7 * 24 * 60 * 60)
	oneDayStep := timeseries.Duration(24 * 60 * 60)

	userBreakdownData, err := ch.GetUserBreakdown(context.Background(), sevenDaysAgo, now, oneDayStep)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println(err)
		return report
	}

	userBreakdownChart := userBreakdownGroup.GetOrCreateChart("New vs Returning Users (Last 7 Days)")
	userBreakdownChart.Column()

	userBreakdownChart.AddSeries("New Users", userBreakdownData["newUsers"], "#AB47BC")
	userBreakdownChart.AddSeries("Returning Users", userBreakdownData["returningUsers"], "#42A5F5")

	fmt.Println(report)
	return report
}
