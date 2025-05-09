package auditor

import (
	"context"

	"codexray/clickhouse"
	"codexray/model"
)

func GeneratePerformanceReport(w *model.World, serviceName, pageName string, ch *clickhouse.Client) *model.AuditReport {
	report := model.NewAuditReport(nil, w.Ctx, nil, model.AuditReportPerformance, true)
	report.Status = model.OK

	// Fetch metrics from ClickHouse
	metrics, err := ch.GetPerformanceTimeSeries(context.Background(), serviceName, pageName, w.Ctx.From, w.Ctx.To, w.Ctx.Step)
	if err != nil {
		report.Status = model.UNKNOWN
		return report
	}

	//
	// Create & populate charts
	//

	// 1) Load Chart

	loadChart := report.GetOrCreateChart("Load requests", nil).Stacked()
	loadChart.AddSeries("Page Loaded", metrics["requests"], "light-blue")

	// 2) Response Time Chart
	responseTimeChart := report.GetOrCreateChart("Response Time, ms", nil).Stacked()
	responseTimeChart.AddSeries("Response Time", metrics["loadTime"], "green")

	// 3) Users Impacted Chart
	usersImpactedChart := report.GetOrCreateChart("Users Impacted", nil).Stacked()
	usersImpactedChart.AddSeries("Users Impacted", metrics["usersImpacted"], "red")

	// 4) Chart Group for JS and API Errors

	errorChart := report.GetOrCreateChart("Errors", nil).Stacked()
	errorChart.AddSeries("JS Errors", metrics["jsErrors"], "orange")
	errorChart.AddSeries("API Errors", metrics["apiErrors"], "purple")

	// 5) User-Centric Metrics Chart
	userCentric := model.NewChart(w.Ctx, "User-Centric Metrics").Stacked()
	userCentric.AddSeries("DNS Time", metrics["dnsTime"], "cyan")
	userCentric.AddSeries("TCP Time", metrics["tcpTime"], "magenta")
	userCentric.AddSeries("SSL Time", metrics["sslTime"], "yellow")
	userCentric.AddSeries("DOM Analysis Time", metrics["domAnalysisTime"], "lime")
	userCentric.AddSeries("DOM Ready Time", metrics["domReadyTime"], "pink")
	userCentric.AddSeries("First Pack Time", metrics["firstPackTime"], "teal")
	userCentric.AddSeries("FMP Time", metrics["fmpTime"], "brown")
	userCentric.AddSeries("FPT Time", metrics["fptTime"], "navy")
	userCentric.AddSeries("Redirect Time", metrics["redirectTime"], "olive")
	userCentric.AddSeries("TTFB Time", metrics["ttfbTime"], "maroon")
	userCentric.AddSeries("TTL Time", metrics["ttlTime"], "gray")
	userCentric.AddSeries("Trans Time", metrics["transTime"], "black")
	userCentric.AddSeries("Response Time", metrics["responseTime"], "blue")
	userCentricWidget := &model.Widget{
		Chart: userCentric,
		Width: "100%",
	}
	report.AddWidget(userCentricWidget)
	return report
}
