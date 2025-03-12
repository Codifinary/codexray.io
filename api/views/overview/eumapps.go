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
	Status    model.Status      `json:"status"`
	Message   string            `json:"message"`
	Overviews []ServiceOverview `json:"overviews"`
	BadgeView Badge             `json:"badgeView"`
	Report    model.AuditReport `json:"report"`
	Limit     int               `json:"limit"`
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

	// Create a new audit report
	report := model.NewAuditReport(nil, w.Ctx, nil, model.AuditReportPerformance, true)

	// Create and add the ECharts
	eCharts, err := createECharts(ctx, ch, from, to, overviews)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Failed to create ECharts: %s", err)
		return v
	}
	for _, eChart := range eCharts {
		report.AddEChartWidget(eChart, nil)
	}

	// Create and add the line charts
	// lineChartWidgets := createLineChartWidgets(ctx, from, to)
	// for _, widget := range lineChartWidgets {
	// 	report.AddWidget(&widget)
	// }

	v.Report = *report
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
	if duration == time.Hour {
		previousFrom := from.Add(-duration)
		previousTotalErrors, err := ch.GetTotalErrors(ctx, &previousFrom, &from, "", "")
		if err != nil {
			return nil, Badge{}, err
		}
		errorTrend = float64(totalErrors-previousTotalErrors) / float64(previousTotalErrors) * 100
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

func extractServiceNames(dataPoints []model.DataPoint) []string {
	names := make([]string, len(dataPoints))
	for i, dp := range dataPoints {
		names[i] = dp.Name
	}
	return names
}

func createECharts(ctx context.Context, ch *clickhouse.Client, from, to time.Time, overviews []ServiceOverview) ([]*model.EChart, error) {
	var eCharts []*model.EChart

	// Fetch top browsers from perf_table
	topBrowsers, err := ch.GetTopBrowser(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get top browsers: %w", err)
	}

	topBrowsersData := make([]model.DataPoint, len(topBrowsers))
	for i, browser := range topBrowsers {
		topBrowsersData[i] = model.DataPoint{
			Value: browser.Value,
			Name:  browser.Name,
		}
	}

	// for local testing
	// topBrowsers := []model.DataPoint{
	//     {Value: 40, Name: "Chrome"},
	//     {Value: 30, Name: "Firefox"},
	//     {Value: 20, Name: "Safari"},
	//     {Value: 5, Name: "Edge"},
	//     {Value: 5, Name: "Opera"},
	// }

	// Create the donut chart for top 5 browsers
	donutChart1 := model.NewEChart("Top 5 Browsers")
	donutChart1.Tooltip = model.Tooltip{Trigger: "item"}
	donutChart1.Legend = model.Legend{Bottom: "0"}
	donutChart1.SetSeries("Browsers", "pie", topBrowsersData)
	eCharts = append(eCharts, donutChart1)

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
	donutChart2 := model.NewEChart("Top 5 Services by Impacted Users")
	donutChart2.Tooltip = model.Tooltip{Trigger: "item"}
	donutChart2.Legend = model.Legend{Bottom: "0"}
	donutChart2.SetSeries("Services", "pie", topServicesByUsers)
	eCharts = append(eCharts, donutChart2)

	// Create the bar chart for top 10 services by load
	sort.Slice(overviews, func(i, j int) bool {
		return overviews[i].AvgLoadPageTime > overviews[j].AvgLoadPageTime
	})
	topServicesByLoad := make([]model.DataPoint, 0, 10)
	for i, overview := range overviews {
		if i >= 10 {
			break
		}
		topServicesByLoad = append(topServicesByLoad, model.DataPoint{
			Value: int(overview.AvgLoadPageTime),
			Name:  overview.ServiceName,
		})
	}
	barChart := model.NewEChart("Top 10 Services by Load")
	barChart.Tooltip = model.Tooltip{Trigger: "axis"}
	barChart.Legend = model.Legend{Bottom: "0"}
	barChart.XAxis = &model.Axis{Type: "value"}
	barChart.YAxis = &model.Axis{Type: "category", Data: extractServiceNames(topServicesByLoad)}
	barChart.SetSeries("Load", "bar", topServicesByLoad)
	barChart.Series.BarWidth = "40%"
	eCharts = append(eCharts, barChart)

	return eCharts, nil
}
