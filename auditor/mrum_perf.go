package auditor

import (
	"codexray/clickhouse"
	"codexray/model"
	"codexray/timeseries"
	"context"
	"fmt"
)

func GenerateMrumPerfReport(w *model.World, ch *clickhouse.Client, from, to *timeseries.Time) *model.AuditReport {
	report := model.NewAuditReport(nil, w.Ctx, nil, model.AuditReportMrumPerf, true)
	report.Status = model.OK

	requestByTimeSliceChartData, err := ch.GetRequestsByTimeSliceChart(context.Background(), from, to)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println(err)
		return report
	}

	requestByTimeSliceChart := model.NewChart(w.Ctx, "Requests by Time Slice")
	requestByTimeSliceChart.AddSeries("Requests by Time Slice", requestByTimeSliceChartData, "light-blue")

	errorRateTrendByTimeChartData, err := ch.GetErrorRateTrendByTimeChart(context.Background(), from, to)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println(err)
		return report
	}

	errorRateTrendByTimeChart := model.NewChart(w.Ctx, "Error Rate Trend by Time")
	errorRateTrendByTimeChart.AddSeries("Error Rate Trend by Time", errorRateTrendByTimeChartData, "red")

	userImptactedByErrorsByTimeChartData, err := ch.GetUserImptactedByErrorsByTimeChart(context.Background(), from, to)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println(err)
		return report
	}

	userImptactedByErrorsByTimeChart := model.NewChart(w.Ctx, "User Impacted by Errors by Time")
	userImptactedByErrorsByTimeChart.AddSeries("User Impacted by Errors by Time", userImptactedByErrorsByTimeChartData, "orange")
	return report
}
