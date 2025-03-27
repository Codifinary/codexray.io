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

	colors := []string{
		"#4285F4", // Google Blue
		"#EA4335", // Google Red
		"#FBBC05", // Google Yellow
	}

	sessionsByCountryChart := sessionTrendsGroup.GetOrCreateChart("Sessions by Country")
	sessionsByCountryData, err := ch.GetSessionsByCountryTrendChart(context.Background(), sevenDays, now, oneHourStep)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println("Error getting sessions by country data:", err)
		return report
	}

	if len(sessionsByCountryData) == 0 {
		emptyTs := timeseries.New(sevenDays, int(now.Sub(sevenDays)/oneHourStep), oneHourStep)
		emptyTs.Set(sevenDays, 0)
		emptyTs.Set(now-1, 0)
		sessionsByCountryChart.AddSeries("No data available", emptyTs, "#CCCCCC")
	} else {
		colorIndex := 0
		for country, timeSeries := range sessionsByCountryData {
			color := colors[colorIndex%len(colors)]
			sessionsByCountryChart.AddSeries(country, timeSeries, color)
			colorIndex++
		}
	}

	sessionsByDeviceChart := sessionTrendsGroup.GetOrCreateChart("Sessions by Device")
	sessionsByDeviceData, err := ch.GetSessionsByDeviceTrendChart(context.Background(), sevenDays, now, oneHourStep)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println("Error getting sessions by device data:", err)
		return report
	}

	if len(sessionsByDeviceData) == 0 {
		emptyTs := timeseries.New(sevenDays, int(now.Sub(sevenDays)/oneHourStep), oneHourStep)
		emptyTs.Set(sevenDays, 0)
		emptyTs.Set(now-1, 0)
		sessionsByDeviceChart.AddSeries("No data available", emptyTs, "#CCCCCC")
	} else {
		colorIndex := 0
		for device, timeSeries := range sessionsByDeviceData {
			color := colors[colorIndex%len(colors)]
			sessionsByDeviceChart.AddSeries(device, timeSeries, color)
			colorIndex++
		}
	}

	sessionsByOSChart := sessionTrendsGroup.GetOrCreateChart("Sessions by Operating System")
	sessionsByOSData, err := ch.GetSessionsByOSTrendChart(context.Background(), sevenDays, now, oneHourStep)
	if err != nil {
		report.Status = model.WARNING
		fmt.Println("Error getting sessions by OS data:", err)
		return report
	}

	if len(sessionsByOSData) == 0 {
		emptyTs := timeseries.New(sevenDays, int(now.Sub(sevenDays)/oneHourStep), oneHourStep)
		emptyTs.Set(sevenDays, 0)
		emptyTs.Set(now-1, 0)
		sessionsByOSChart.AddSeries("No data available", emptyTs, "#CCCCCC")
	} else {
		colorIndex := 0
		for os, timeSeries := range sessionsByOSData {
			color := colors[colorIndex%len(colors)]
			sessionsByOSChart.AddSeries(os, timeSeries, color)
			colorIndex++
		}
	}

	return report
}
