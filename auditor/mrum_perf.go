package auditor

import (
	"codexray/clickhouse"
	"codexray/model"
	"codexray/timeseries"
	"context"
	"fmt"
)

func GenerateMrumPerfReport(w *model.World, ch *clickhouse.Client, from, to timeseries.Time, service string) *model.AuditReport {
	report := model.NewAuditReport(nil, w.Ctx, nil, model.AuditReportMrumPerf, true)
	report.Status = model.OK
	requestByTimeSliceChartData, err := ch.GetRequestsByTimeSliceChart(context.Background(), from, to, w.Ctx.Step, service)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println(err)
		return report
	}

	requestByTimeSliceChart := report.GetOrCreateChart("Requests by Time Slice", nil)
	requestByTimeSliceChart.AddSeries("Requests by Time Slice", requestByTimeSliceChartData, "light-blue")

	errorRateTrendByTimeChartData, err := ch.GetErrorRateTrendByTimeChart(context.Background(), from, to, w.Ctx.Step, service)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println(err)
		return report
	}

	errorRateTrendByTimeChart := report.GetOrCreateChart("Error Rate Trend by Time", nil)
	errorRateTrendByTimeChart.AddSeries("Error Rate Trend by Time", errorRateTrendByTimeChartData, "red")

	userImptactedByErrorsByTimeChartData, err := ch.GetUserImptactedByErrorsByTimeChart(context.Background(), from, to, w.Ctx.Step, service)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println(err)
		return report
	}

	userImptactedByErrorsByTimeChart := report.GetOrCreateChart("User Impacted by Errors by Time", nil)
	userImptactedByErrorsByTimeChart.AddSeries("User Impacted by Errors by Time", userImptactedByErrorsByTimeChartData, "orange")

	sq := clickhouse.SpanQuery{
		Ctx: w.Ctx,
	}
	histogram, err := ch.GetHttpResponsePerfHistogram(context.Background(), sq, service)
	if err == nil && len(histogram) > 1 {
		heatmapWidget := &model.Widget{
			Heatmap: model.NewHeatmap(w.Ctx, "HTTP Response Latency & Errors heatmap, requests per second"),
		}

		for _, h := range model.HistogramSeries(histogram[1:], 0, 0) {
			heatmapWidget.Heatmap.AddSeries(h.Name, h.Title, h.Data, h.Threshold, h.Value)
		}
		heatmapWidget.Heatmap.AddSeries("errors", "errors", histogram[0].TimeSeries, "", "err")

		report.Widgets = append(report.Widgets, heatmapWidget)
	} else {
		fmt.Println(err, len(histogram))
	}
	return report
}
