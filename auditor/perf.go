package auditor

import (
	"context"

	"codexray/clickhouse"
	"codexray/model"
)

type PerformanceReport struct {
	Status             string            `json:"status"`
	RequestChart       *model.Chart      `json:"request_chart"`
	ResponseTimeChart  *model.Chart      `json:"response_time_chart"`
	UsersImpactedChart *model.Chart      `json:"users_impacted_chart"`
	ErrorChartGroup    *model.ChartGroup `json:"error_chart_group"`
	UserCentricChart   *model.Chart      `json:"user_centric_chart"`
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
		report.Status = "warning"
		return report
	}

	// Create separate charts for load time, response time, and users impacted
	requestChart := model.NewChart(w.Ctx, "Load")
	responseTimeChart := model.NewChart(w.Ctx, "Response Time")
	usersImpactedChart := model.NewChart(w.Ctx, "Users Impacted")

	// Create a chart group for JS and API errors
	errorChartGroup := model.NewChartGroup(w.Ctx, "Errors")

	// Add series to charts
	requestChart.AddSeries("Page Loaded", metrics["requests"])
	responseTimeChart.AddSeries("Response Time", metrics["loadTime"])
	usersImpactedChart.AddSeries("Users Impacted", metrics["usersImpacted"])

	// Add JS and API errors to the error chart group
	errorChartGroup.GetOrCreateChart("JS Errors").AddSeries("JS Errors", metrics["jsErrors"])
	errorChartGroup.GetOrCreateChart("API Errors").AddSeries("API Errors", metrics["apiErrors"])

	// Create a user-centric chart and add all extra metrics to it
	userCentricChart := model.NewChart(w.Ctx, "User-Centric Metrics")
	userCentricChart.AddSeries("DNS Time", metrics["dnsTime"])
	userCentricChart.AddSeries("TCP Time", metrics["tcpTime"])
	userCentricChart.AddSeries("SSL Time", metrics["sslTime"])
	userCentricChart.AddSeries("DOM Analysis Time", metrics["domAnalysisTime"])
	userCentricChart.AddSeries("DOM Ready Time", metrics["domReadyTime"])
	userCentricChart.AddSeries("First Pack Time", metrics["firstPackTime"])
	userCentricChart.AddSeries("FMP Time", metrics["fmpTime"])
	userCentricChart.AddSeries("FPT Time", metrics["fptTime"])
	userCentricChart.AddSeries("Redirect Time", metrics["redirectTime"])
	userCentricChart.AddSeries("TTFB Time", metrics["ttfbTime"])
	userCentricChart.AddSeries("TTL Time", metrics["ttlTime"])
	userCentricChart.AddSeries("Trans Time", metrics["transTime"])
	userCentricChart.AddSeries("Response Time", metrics["responseTime"])

	report.RequestChart = requestChart
	report.ResponseTimeChart = responseTimeChart
	report.UsersImpactedChart = usersImpactedChart
	report.ErrorChartGroup = errorChartGroup
	report.UserCentricChart = userCentricChart
	report.Status = "ok"

	return report
}
