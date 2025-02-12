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

	//
	// Create & populate charts
	//

	loadChart := report.GetOrCreateChart("Requests Trend", nil)
	loadChart.AddSeries("Requests Trend", traces["requestsTrend"], "light-blue")

	// 2) error Trend Chart
	responseTimeChart := report.GetOrCreateChart("Error Trend", nil)
	responseTimeChart.AddSeries("Error Trend", traces["errorTrend"], "red")

	// 3) latency Chart

	quantile := report.GetOrCreateChart("Latency Trend", nil)
	quantile.AddSeries("p25", traces["p25Latency"], "cyan")
	quantile.AddSeries("p50", traces["p50Latency"], "magenta")
	quantile.AddSeries("p75", traces["p75Latency"], "yellow")
	quantile.AddSeries("p95", traces["p95Latency"], "pink")
	quantile.AddSeries("p99", traces["p99Latency"], "lime")

	return report
}
