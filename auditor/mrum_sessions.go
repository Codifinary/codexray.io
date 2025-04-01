package auditor

import (
	"codexray/clickhouse"
	"codexray/model"
	"codexray/timeseries"
	"context"
	"fmt"
)

func GenerateMrumSessionsReport(w *model.World, ch *clickhouse.Client, from, to timeseries.Time) *model.AuditReport {
	report := model.NewAuditReport(nil, w.Ctx, nil, model.AuditReportMobileSessions, true)
	report.Status = model.OK

	sessionTrendsGroup := report.GetOrCreateChartGroup("Session Activity Trends", nil)

	now := timeseries.Now()
	sevenDays := now.Add(-7 * 24 * 60 * 60)
	oneHourStep := timeseries.Duration(3600)

	countryColors := []string{
		"#4285F4", // Google Blue
		"#EA4335", // Google Red
		"#FBBC05", // Google Yellow
	}

	deviceColors := []string{
		"#34A853", // Google Green
		"#9C27B0", // Purple
		"#FF9800", // Orange
	}

	osColors := []string{
		"#2196F3", // Light Blue
		"#F44336", // Red
		"#4CAF50", // Green
	}

	sessionsByCountryChart := sessionTrendsGroup.GetOrCreateChart("Sessions by Country")
	sessionsByCountryData, err := ch.GetSessionsByCountryTrendChart(context.Background(), sevenDays, now, oneHourStep)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println("Error getting sessions by country data:", err)
		return report
	}

	countryIndex := 0
	for country, timeSeries := range sessionsByCountryData {
		color := countryColors[countryIndex%len(countryColors)]
		sessionsByCountryChart.AddSeriesWithFill(country, timeSeries, color, true)
		countryIndex++
	}

	sessionsByDeviceChart := sessionTrendsGroup.GetOrCreateChart("Sessions by Device")
	sessionsByDeviceData, err := ch.GetSessionsByDeviceTrendChart(context.Background(), sevenDays, now, oneHourStep)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println("Error getting sessions by device data:", err)
		return report
	}

	deviceIndex := 0
	for device, timeSeries := range sessionsByDeviceData {
		color := deviceColors[deviceIndex%len(deviceColors)]
		sessionsByDeviceChart.AddSeriesWithFill(device, timeSeries, color, true)
		deviceIndex++
	}

	sessionsByOSChart := sessionTrendsGroup.GetOrCreateChart("Sessions by Operating System")
	sessionsByOSData, err := ch.GetSessionsByOSTrendChart(context.Background(), sevenDays, now, oneHourStep)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println("Error getting sessions by OS data:", err)
		return report
	}

	osIndex := 0
	for os, timeSeries := range sessionsByOSData {
		color := osColors[osIndex%len(osColors)]
		sessionsByOSChart.AddSeriesWithFill(os, timeSeries, color, true)
		osIndex++
	}

	return report
}
