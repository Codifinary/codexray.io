package auditor

import (
	"codexray/clickhouse"
	"codexray/model"
	"codexray/timeseries"
	"context"
	"fmt"
)

func GenerateMrumCrashesReport(w *model.World, ch *clickhouse.Client, from, to timeseries.Time, service string) *model.AuditReport {
	report := model.NewAuditReport(nil, w.Ctx, nil, model.AuditReportMobileCrashes, true)
	report.Status = model.OK

	now := timeseries.Now()
	sevenDays := now.Add(-7 * 24 * 60 * 60)
	oneHourStep := timeseries.Duration(3600)

	deviceColors := []string{
		"#f44336", // Red
		"#6a329f", // Purple
		"#f3820c", // Orange
		"#1f693f", // Green
		"#0b5394", // Blue
	}

	crashesByDeviceChart := model.NewChart(w.Ctx, "Crashes by Device")
	crashesByDeviceData, err := ch.GetCrashesByDeviceTrendChart(context.Background(), sevenDays, now, oneHourStep, service)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println("Error getting crashes by device data:", err)
		return report
	}

	deviceIndex := 0
	for device, timeSeries := range crashesByDeviceData {
		color := deviceColors[deviceIndex%len(deviceColors)]
		crashesByDeviceChart.AddSeries(device, timeSeries, color)
		deviceIndex++
	}

	crashByDeviceWidget := &model.Widget{
		Chart: crashesByDeviceChart,
		Width: "100%",
	}
	report.AddWidget(crashByDeviceWidget)

	return report
}
