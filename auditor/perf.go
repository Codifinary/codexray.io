package auditor

import (
	"context"
	"fmt"

	"codexray/clickhouse"
	"codexray/model"
)

type PerformanceReport struct {
	Status             string            `json:"status"`
	LoadTimeChart      *model.Chart      `json:"load_time_chart"`
	ResponseTimeChart  *model.Chart      `json:"response_time_chart"`
	UsersImpactedChart *model.Chart      `json:"users_impacted_chart"`
	ErrorChartGroup    *model.ChartGroup `json:"error_chart_group"`
}

func GeneratePerformanceReport(w *model.World, serviceName, pageName string, ch *clickhouse.Client) *PerformanceReport {
	report := &PerformanceReport{
		Status: "unknown",
	}

	if ch == nil {
		report.Status = "unknown"
		return report
	}

	from := w.Ctx.From
	to := w.Ctx.To
	step := w.Ctx.Step

	// Fetch performance metrics directly from ClickHouse
	metrics, err := ch.GetPerformanceTimeSeries(context.Background(), serviceName, pageName, from, to, step)
	if err != nil {
		fmt.Println("Error fetching performance metrics:", err)
		report.Status = "warning"
		return report
	}

	// Create separate charts for load time, response time, and users impacted
	loadTimeChart := model.NewChart(w.Ctx, "Page Load Time")
	responseTimeChart := model.NewChart(w.Ctx, "Response Time")
	usersImpactedChart := model.NewChart(w.Ctx, "Users Impacted")

	// Create a chart group for JS and API errors
	errorChartGroup := model.NewChartGroup(w.Ctx, "Errors")

	// Add series to charts
	loadTimeChart.AddSeries("Page Load Time", metrics["loadTime"])
	responseTimeChart.AddSeries("Response Time", metrics["responseTime"])
	usersImpactedChart.AddSeries("Users Impacted", metrics["usersImpacted"])

	// Add JS and API errors to the error chart group
	errorChartGroup.GetOrCreateChart("JS Errors").AddSeries("JS Errors", metrics["jsErrors"])
	errorChartGroup.GetOrCreateChart("API Errors").AddSeries("API Errors", metrics["apiErrors"])

	report.LoadTimeChart = loadTimeChart
	report.ResponseTimeChart = responseTimeChart
	report.UsersImpactedChart = usersImpactedChart
	report.ErrorChartGroup = errorChartGroup
	report.Status = "ok"

	fmt.Println("Generated performance report:", report)

	return report
}
