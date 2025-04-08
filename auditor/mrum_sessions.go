package auditor

import (
	"codexray/clickhouse"
	"codexray/model"
	"codexray/timeseries"
	"context"
	"fmt"
)

func GenerateMrumSessionsReport(w *model.World, ch *clickhouse.Client, from, to timeseries.Time, service string) *model.AuditReport {
	report := model.NewAuditReport(nil, w.Ctx, nil, model.AuditReportMobileSessions, true)
	report.Status = model.OK

	now := timeseries.Now()
	sevenDays := now.Add(-7 * 24 * 60 * 60)
	oneHourStep := timeseries.Duration(3600)

	countryColors := []string{
		"#b4a7d6", // Purple
		"#ffd966", // Orange
		"#bcbcbc", // Grey
	}

	deviceColors := []string{
		"#6fa8dc", // Light Blue
		"#2986cc", // Blue
		"#bcbcbc", // Grey
	}

	osColors := []string{
		"#f44336", // Red
		"#f1c232", // Yellow
		"#bcbcbc", // Grey
	}

	width := "33%"

	sessionsByCountryData, err := ch.GetSessionsByCountryTrendChart(context.Background(), sevenDays, now, oneHourStep, service)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println("Error getting sessions by country data:", err)
		return report
	}

	sessionsByCountryChart := model.NewChart(w.Ctx, "Sessions by Country")
	countryIndex := 0
	for country, timeSeries := range sessionsByCountryData {
		color := countryColors[countryIndex%len(countryColors)]
		sessionsByCountryChart.AddSeriesWithFill(country, timeSeries, color, true)
		countryIndex++
	}

	sessionsByCountryWidget := &model.Widget{
		Chart: sessionsByCountryChart,
		Width: width,
	}
	report.AddWidget(sessionsByCountryWidget)

	sessionsByDeviceData, err := ch.GetSessionsByDeviceTrendChart(context.Background(), sevenDays, now, oneHourStep, service)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println("Error getting sessions by device data:", err)
		return report
	}

	sessionsByDeviceChart := model.NewChart(w.Ctx, "Sessions by Device")
	deviceIndex := 0
	for device, timeSeries := range sessionsByDeviceData {
		color := deviceColors[deviceIndex%len(deviceColors)]
		sessionsByDeviceChart.AddSeriesWithFill(device, timeSeries, color, true)
		deviceIndex++
	}

	sessionsByDeviceWidget := &model.Widget{
		Chart: sessionsByDeviceChart,
		Width: width,
	}
	report.AddWidget(sessionsByDeviceWidget)

	sessionsByOSData, err := ch.GetSessionsByOSTrendChart(context.Background(), sevenDays, now, oneHourStep, service)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println("Error getting sessions by OS data:", err)
		return report
	}

	sessionsByOSChart := model.NewChart(w.Ctx, "Sessions by Operating System")
	osIndex := 0
	for os, timeSeries := range sessionsByOSData {
		color := osColors[osIndex%len(osColors)]
		sessionsByOSChart.AddSeriesWithFill(os, timeSeries, color, true)
		osIndex++
	}

	sessionsByOSWidget := &model.Widget{
		Chart: sessionsByOSChart,
		Width: width,
	}
	report.AddWidget(sessionsByOSWidget)

	return report
}
