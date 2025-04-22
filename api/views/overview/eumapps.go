package overview

import (
	"context"
	"fmt"
	"sort"
	"time"

	"codexray/clickhouse"
	"codexray/model"

	"k8s.io/klog"
)

type EumView struct {
	Status       model.Status      `json:"status"`
	Message      string            `json:"message"`
	Overviews    []ServiceOverview `json:"overviews"`
	BadgeView    Badge             `json:"badgeView"`
	Report       model.AuditReport `json:"report"`
	EchartReport model.AuditReport `json:"Echartreport"`
	Limit        int               `json:"limit"`
}

type ServiceOverview struct {
	ServiceName        string  `json:"serviceName"`
	Pages              uint64  `json:"pages"`
	AvgLoadPageTime    float64 `json:"avgLoadPageTime"`
	JsErrorPercentage  float64 `json:"jsErrorPercentage"`
	ApiErrorPercentage float64 `json:"apiErrorPercentage"`
	ImpactedUsers      uint64  `json:"impactedUsers"`
	AppType            string  `json:"appType"`
	Requests           uint64  `json:"requests"`
}

type Badge struct {
	TotalApplications uint64  `json:"totalApplications"`
	TotalPages        uint64  `json:"totalPages"`
	AvgLatency        float64 `json:"avgLatency"`
	TotalErrors       uint64  `json:"totalError"`
	ErrorPerSec       float64 `json:"errorPerSec"`
	ErrorTrend        float64 `json:"errorTrend"`
}

func renderEumApps(ctx context.Context, ch *clickhouse.Client, w *model.World, query string) *EumView {
	v := &EumView{}

	from := w.Ctx.From.ToStandard()
	to := w.Ctx.To.ToStandard()

	overviews, badge, err := getServiceOverviews(ctx, ch, from, to)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Clickhouse error: %s", err)
		return v
	}

	v.Overviews = overviews
	v.BadgeView = badge

	echartReport, err := createECharts(w, ctx, ch, from, to, overviews)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Failed to create echart: %s", err)
		return v
	}

	lineChartReport, err := createLineChart(w, ch, overviews)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Failed to create line charts: %s", err)
		return v
	}

	v.Report = *lineChartReport
	v.EchartReport = *echartReport
	v.Status = model.OK
	return v
}

func getServiceOverviews(ctx context.Context, ch *clickhouse.Client, from, to time.Time) ([]ServiceOverview, Badge, error) {
	rows, err := ch.GetServiceOverviews(ctx, &from, &to)
	if err != nil {
		return nil, Badge{}, err
	}

	var overviews []ServiceOverview
	var totalApplications, totalPages uint64
	var totalLatency float64

	for _, row := range rows {
		overviews = append(overviews, ServiceOverview{
			ServiceName:        row.ServiceName,
			Pages:              row.Pages,
			AvgLoadPageTime:    row.AvgLoadPageTime,
			JsErrorPercentage:  row.JsErrorPercentage,
			ApiErrorPercentage: row.ApiErrorPercentage,
			ImpactedUsers:      row.ImpactedUsers,
			Requests:           row.Requests,
			AppType:            row.AppType,
		})
		totalApplications++
		totalPages += row.Pages
		totalLatency += row.AvgLoadPageTime
	}

	sort.Slice(overviews, func(i, j int) bool {
		return overviews[i].ServiceName < overviews[j].ServiceName
	})

	// Calculate average latency
	avgLatency := totalLatency / float64(totalApplications)

	// Get total errors and error trend
	totalErrors, err := ch.GetTotalErrors(ctx, &from, &to, "", "")
	if err != nil {
		return nil, Badge{}, err
	}

	var errorTrend float64
	duration := to.Sub(from)
	if duration <= 60*time.Minute {
		previousFrom := from.Add(-duration)
		previousTotalErrors, err := ch.GetTotalErrors(ctx, &previousFrom, &from, "", "")
		if err != nil {
			return nil, Badge{}, err
		}

		if previousTotalErrors == 0 {
			errorTrend = float64(totalErrors) * 100
		} else {
			errorTrend = float64(int64(totalErrors)-int64(previousTotalErrors)) / float64(previousTotalErrors) * 100
		}

	} else {
		errorTrend = 0
	}

	badge := Badge{
		TotalApplications: totalApplications,
		TotalPages:        totalPages,
		AvgLatency:        avgLatency,
		TotalErrors:       totalErrors,
		ErrorTrend:        errorTrend,
	}

	return overviews, badge, nil
}

func createECharts(w *model.World, ctx context.Context, ch *clickhouse.Client, from, to time.Time, overviews []ServiceOverview) (*model.AuditReport, error) {
	// Create a new audit report for ECharts
	echartReport := model.NewAuditReport(nil, w.Ctx, nil, model.AuditReportPerformance, true)

	// Fetch top browsers from perf_data
	topBrowsers, err := ch.GetTopBrowser(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get top browsers: %w", err)
	}

	topBrowsersData := make([]model.DataPoint, len(topBrowsers))
	for i, browser := range topBrowsers {
		topBrowsersData[i] = model.DataPoint{
			Value: int(browser.Value),
			Name:  browser.Name,
		}
	}

	// for local testing
	// topBrowsersData := []model.DataPoint{
	// 	{Value: 40, Name: "Chrome"},
	// 	{Value: 30, Name: "Firefox"},
	// 	{Value: 20, Name: "Safari"},
	// 	{Value: 5, Name: "Edge"},
	// 	{Value: 5, Name: "Opera"},
	// }

	// Create the donut chart for top 5 browsers
	donutChart1 := echartReport.GetOrCreateEChart("Top Services by Browsers", nil)
	donutChart1.Title = model.TextTitle{
		Text: "Top Services by Browsers",
		TextStyle: &model.TextStyle{
			FontSize:   16,
			FontWeight: "normal",
		},
	}
	donutChart1.Tooltip = model.Tooltip{Trigger: "item"}
	donutChart1.Legend = model.Legend{Bottom: "0"}
	donutChart1.SetPieChartSeries("Browsers", "pie", []string{"40%", "70%"}, topBrowsersData)
	donutChart1.Graphic = &model.Graphic{
		Type: "text",
		Left: "center",
		Top:  "center",
		Style: model.GraphicStyle{
			Text:       fmt.Sprintf("%d", len(topBrowsersData)),
			FontSize:   26,
			FontWeight: "bold",
			Fill:       "#333",
			TextAlign:  "center",
		},
	}

	// Create the donut chart for top 5 services by impacted users
	sort.Slice(overviews, func(i, j int) bool {
		return overviews[i].ImpactedUsers > overviews[j].ImpactedUsers
	})
	topServicesByUsers := make([]model.DataPoint, 0, 5)
	for i, overview := range overviews {
		if i >= 5 {
			break
		}
		topServicesByUsers = append(topServicesByUsers, model.DataPoint{
			Value: int(overview.ImpactedUsers),
			Name:  overview.ServiceName,
		})
	}
	donutChart2 := echartReport.GetOrCreateEChart("Top Services by Impacted Users", nil)
	donutChart2.Title = model.TextTitle{
		Text: "Top Services by Impacted Users",
		TextStyle: &model.TextStyle{
			FontSize:   16,
			FontWeight: "normal",
		},
	}
	donutChart2.Tooltip = model.Tooltip{Trigger: "item"}
	donutChart2.Legend = model.Legend{Bottom: "0"}

	donutChart2.SetPieChartSeries("Services", "pie", []string{"40%", "70%"}, topServicesByUsers)

	donutChart2.Graphic = &model.Graphic{
		Type: "text",
		Left: "center",
		Top:  "center",
		Style: model.GraphicStyle{
			Text:       fmt.Sprintf("%d", len(topServicesByUsers)),
			FontSize:   26,
			FontWeight: "bold",
			Fill:       "#333",
			TextAlign:  "center",
		},
	}

	// Create the bar chart for top 10 services by load
	sort.Slice(overviews, func(i, j int) bool {
		return overviews[i].Requests > overviews[j].Requests
	})
	topServicesByLoad := make([]model.DataPoint, 0, 10)
	for i, overview := range overviews {
		if i >= 10 {
			break
		}
		topServicesByLoad = append(topServicesByLoad, model.DataPoint{
			Value: int(overview.Requests),
			Name:  overview.ServiceName,
		})
	}

	sort.Slice(topServicesByLoad, func(i, j int) bool {
		return topServicesByLoad[i].Value < topServicesByLoad[j].Value
	})

	barChart := echartReport.GetOrCreateEChart("Top Services by Load", nil)
	barChart.Title = model.TextTitle{
		Text: "Top Services by Load",
		TextStyle: &model.TextStyle{
			FontSize:   16,
			FontWeight: "normal",
		},
	}
	barChart.Tooltip = model.Tooltip{Trigger: "axis"}
	barChart.Grid = &model.Grid{ContainLabel: true, Top: 30,
		Bottom: 40,
		Left:   10,
		Right:  20}
	barChart.Legend = model.Legend{Bottom: "0"}
	barChart.XAxis = &model.Axis{Type: "value"}
	barChart.YAxis = &model.Axis{
		Type: "category",
		Data: extractServiceNames(topServicesByLoad),
		AxisLabel: &model.AxisLabel{
			Show:      true,
			FontSize:  12,
			Formatter: "{value}",
			Rotate:    0,
		},
	}
	barChart.SetSeries("Load", "bar", topServicesByLoad)
	barChart.Series.BarWidth = "40%"

	return echartReport, nil
}

func createLineChart(w *model.World, ch *clickhouse.Client, overviews []ServiceOverview) (*model.AuditReport, error) {
	sort.Slice(overviews, func(i, j int) bool {
		return overviews[i].Requests > overviews[j].Requests
	})

	topServices := overviews
	if len(overviews) > 5 {
		topServices = overviews[:5]
	}

	lineChartReport := model.NewAuditReport(nil, w.Ctx, nil, model.AuditReportPerformance, true)
	lineChartReport.Status = model.OK

	loadChart := lineChartReport.GetOrCreateChart("Top 5 Services by Load", nil)
	for _, service := range topServices {
		loadTimeSeries, err := ch.GetLoadTimeSeries(context.Background(), service.ServiceName, w.Ctx.From, w.Ctx.To, w.Ctx.Step)
		if err != nil {
			return nil, fmt.Errorf("failed to get load time series for service %s: %w", service.ServiceName, err)
		}
		loadChart.AddSeries(service.ServiceName, loadTimeSeries)
	}

	responseChart := lineChartReport.GetOrCreateChart("Top 5 Services by Response Time", nil)
	for _, service := range topServices {
		responseTimeSeries, err := ch.GetResponseTimeSeries(context.Background(), service.ServiceName, w.Ctx.From, w.Ctx.To, w.Ctx.Step)
		if err != nil {
			return nil, fmt.Errorf("failed to get response time series for service %s: %w", service.ServiceName, err)
		}
		responseChart.AddSeries(service.ServiceName, responseTimeSeries)
	}

	errorChart := lineChartReport.GetOrCreateChart("Top 5 Services by Errors", nil)
	for _, service := range topServices {
		errorTimeSeries, err := ch.GetErrorTimeSeries(context.Background(), service.ServiceName, w.Ctx.From, w.Ctx.To, w.Ctx.Step)
		if err != nil {
			return nil, fmt.Errorf("failed to get error time series for service %s: %w", service.ServiceName, err)
		}
		errorChart.AddSeries(service.ServiceName, errorTimeSeries)
	}

	usersImpactedChart := lineChartReport.GetOrCreateChart("Top 5 Services by Users Impacted", nil)
	for _, service := range topServices {
		usersImpactedTimeSeries, err := ch.GetUsersImpactedTimeSeries(context.Background(), service.ServiceName, w.Ctx.From, w.Ctx.To, w.Ctx.Step)
		if err != nil {
			return nil, fmt.Errorf("failed to get users impacted time series for service %s: %w", service.ServiceName, err)
		}
		usersImpactedChart.AddSeries(service.ServiceName, usersImpactedTimeSeries)
	}

	return lineChartReport, nil
}
func extractServiceNames(dataPoints []model.DataPoint) []string {
	names := make([]string, len(dataPoints))
	for i, dp := range dataPoints {
		names[i] = dp.Name
	}
	return names
}
