package auditor

import (
	"context"
	"fmt"

	"codexray/clickhouse"
	"codexray/model"
)

func GenerateTracesReportByService(w *model.World, serviceName string, ch *clickhouse.Client) *model.AuditReport {
	report := model.NewAuditReport(nil, w.Ctx, nil, model.AuditReportTraces, true)
	report.Status = model.OK

	// Fetch metrics from ClickHouse
	traces, err := ch.GetTracesChartsByServiceName(context.Background(), serviceName, w.Ctx.From, w.Ctx.To, w.Ctx.Step)
	if err != nil {
		report.Status = model.UNKNOWN
		fmt.Println(err)
		return report
	}

	// Create & populate charts
	loadChart := model.NewChart(w.Ctx, "Requests Trend")
	loadChart.AddSeries("Requests Trend", traces["requestsTrend"], "light-blue")

	errorTrendChart := model.NewChart(w.Ctx, "Error Trend")
	errorTrendChart.AddSeries("Error Trend", traces["errorTrend"], "red")

	latencyTrendChart := model.NewChart(w.Ctx, "Latency Trend")
	latencyTrendChart.AddSeries("p25", traces["p25Latency"], "cyan")
	latencyTrendChart.AddSeries("p50", traces["p50Latency"], "magenta")
	latencyTrendChart.AddSeries("p75", traces["p75Latency"], "yellow")
	latencyTrendChart.AddSeries("p95", traces["p95Latency"], "pink")
	latencyTrendChart.AddSeries("p99", traces["p99Latency"], "lime")

	// Create widgets and add each chart to a separate widget
	requestsWidget := &model.Widget{
		Chart: loadChart,
		Width: "33%",
	}
	errorWidget := &model.Widget{
		Chart: errorTrendChart,
		Width: "33%",
	}
	latencyWidget := &model.Widget{
		Chart: latencyTrendChart,
		Width: "33%",
	}

	// Add the widgets to the report
	report.AddWidget(requestsWidget)
	report.AddWidget(errorWidget)
	report.AddWidget(latencyWidget)

	return report
}
