package auditor

import (
	"codexray/clickhouse"
	"codexray/model"
	"codexray/timeseries"
	"context"
	"fmt"
)

func GenerateMrumCrashesReport(w *model.World, ch *clickhouse.Client, from, to timeseries.Time) *model.AuditReport {
	report := model.NewAuditReport(nil, w.Ctx, nil, model.AuditReportMobileCrashes, true)
	report.Status = model.OK

	now := timeseries.Now()
	sevenDays := now.Add(-7 * 24 * 60 * 60)
	oneHourStep := timeseries.Duration(3600)

	deviceColors := []string{
		"#E53935", // Red
		"#8E24AA", // Purple
		"#FB8C00", // Orange
		"#43A047", // Green
		"#1E88E5", // Blue
	}

	crashesByDeviceChart := report.GetOrCreateChart("Crashes by Device", nil)
	crashesByDeviceData, err := ch.GetCrashesByDeviceTrendChart(context.Background(), sevenDays, now, oneHourStep)
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

	return report
}
