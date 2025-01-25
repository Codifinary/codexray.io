package auditor

import (
	"context"
	"time"

	"codexray/clickhouse"
	"codexray/model"
	"codexray/timeseries"
)

func AuditPerf(w *model.World, serviceName, pageName string, ch *clickhouse.Client) *model.AuditReport {
	report := model.NewAuditReport(nil, w.Ctx, nil, model.AuditReportPerformance, true)

	if ch == nil {
		report.Status = model.UNKNOWN
		report.GetOrCreateTable("Error").AddRow(
			model.NewTableCell("Clickhouse integration is not configured."),
		)
		return report
	}

	from := w.Ctx.From.ToStandard()
	to := w.Ctx.To.ToStandard()
	stepSeconds := int64(w.Ctx.Step) / int64(time.Second)

	chartData, err := ch.GetChartData(context.Background(), serviceName, pageName, &from, &to, stepSeconds)
	if err != nil {
		report.Status = model.WARNING
		report.GetOrCreateTable("Error").AddRow(
			model.NewTableCell("Failed to fetch performance data").SetUnit(err.Error()),
		)
		return report
	}

	// If no data
	if len(chartData) == 0 {
		report.Status = model.WARNING
		report.GetOrCreateTable("Error").AddRow(
			model.NewTableCell("No performance data available for the given service and page."),
		)
		return report
	}

	perfChart := report.GetOrCreateChart("Performance Metrics", nil)

	loadTimeSeries := timeseries.New(w.Ctx.From, len(chartData), w.Ctx.Step)
	responseTimeSeries := timeseries.New(w.Ctx.From, len(chartData), w.Ctx.Step)
	jsErrorsSeries := timeseries.New(w.Ctx.From, len(chartData), w.Ctx.Step)
	apiErrorsSeries := timeseries.New(w.Ctx.From, len(chartData), w.Ctx.Step)
	usersImpactedSeries := timeseries.New(w.Ctx.From, len(chartData), w.Ctx.Step)

	for _, row := range chartData {
		ts := timeseries.Time(row.Timestamp / 1000)
		loadTimeSeries.Set(ts, float32(row.LoadTime))
		responseTimeSeries.Set(ts, float32(row.ResponseTime))
		jsErrorsSeries.Set(ts, float32(row.JsErrors))
		apiErrorsSeries.Set(ts, float32(row.ApiErrors))
		usersImpactedSeries.Set(ts, float32(row.UsersImpacted))
	}

	perfChart.AddSeries("Page Load Time", loadTimeSeries)
	perfChart.AddSeries("Response Time", responseTimeSeries)
	perfChart.AddSeries("JS Errors", jsErrorsSeries)
	perfChart.AddSeries("API Errors", apiErrorsSeries)
	perfChart.AddSeries("Users Impacted", usersImpactedSeries)

	report.Status = model.OK
	return report
}
