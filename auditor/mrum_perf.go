package auditor

import (
	"codexray/clickhouse"
	"codexray/model"
	"codexray/timeseries"
	"context"
	"fmt"
)

func GenerateMrumPerfReport(w *model.World, ch *clickhouse.Client, from, to timeseries.Time) *model.AuditReport {
	report := model.NewAuditReport(nil, w.Ctx, nil, model.AuditReportMrumPerf, true)
	report.Status = model.OK
	requestByTimeSliceChartData, err := ch.GetRequestsByTimeSliceChart(context.Background(), from, to, w.Ctx.Step)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println(err)
		return report
	}

	requestByTimeSliceChart := report.GetOrCreateChart("Requests by Time Slice", nil).Stacked()
	requestByTimeSliceChart.AddSeries("Requests by Time Slice", requestByTimeSliceChartData, "light-blue")

	errorRateTrendByTimeChartData, err := ch.GetErrorRateTrendByTimeChart(context.Background(), from, to, w.Ctx.Step)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println(err)
		return report
	}

	errorRateTrendByTimeChart := report.GetOrCreateChart("Error Rate Trend by Time", nil).Stacked()
	errorRateTrendByTimeChart.AddSeries("Error Rate Trend by Time", errorRateTrendByTimeChartData, "red")

	userImptactedByErrorsByTimeChartData, err := ch.GetUserImptactedByErrorsByTimeChart(context.Background(), from, to, w.Ctx.Step)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println(err)
		return report
	}

	userImptactedByErrorsByTimeChart := report.GetOrCreateChart("User Impacted by Errors by Time", nil).Stacked()
	userImptactedByErrorsByTimeChart.AddSeries("User Impacted by Errors by Time", userImptactedByErrorsByTimeChartData, "orange")
	return report
}
