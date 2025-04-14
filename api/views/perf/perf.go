package perf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"time"

	"codexray/clickhouse"
	"codexray/model"

	"k8s.io/klog"
)

const defaultLimit = 100

type View struct {
	Status       model.Status      `json:"status"`
	Message      string            `json:"message"`
	Overviews    []PerfOverview    `json:"overviews"`
	BadgeView    Badge             `json:"badgeView"`
	Report       model.AuditReport `json:"report"`
	EchartReport model.AuditReport `json:"echartReport"`
	BrowserStats []BrowserStats    `json:"browserStats"`
	Limit        int               `json:"limit"`
}

type Query struct {
	Limit int `json:"limit"`
}

type PerfOverview struct {
	PagePath           string  `json:"pagePath"`
	AvgLoadPageTime    float64 `json:"avgLoadPageTime"`
	JsErrorPercentage  float64 `json:"jsErrorPercentage"`
	ApiErrorPercentage float64 `json:"apiErrorPercentage"`
	ImpactedUsers      uint64  `json:"impactedUsers"`
	Requests           uint64  `json:"requests"`
}

type Badge struct {
	TotalRequests     uint64  `json:"totalRequests"`
	TotalErrors       uint64  `json:"totalErrors"`
	RequestTrend      float64 `json:"requestTrend"`
	ErrorTrend        float64 `json:"errorTrend"`
	RequestsPerSecond float64 `json:"requestsPerSecond"`
	ErrorsPerSecond   float64 `json:"errorsPerSecond"`
}

type BrowserStats struct {
	Name         string  `json:"name"`
	Requests     uint64  `json:"requests"`
	ResponseTime float64 `json:"responseTime"`
	Errors       uint64  `json:"errors"`
}

func Render(w *model.World, ctx context.Context, ch *clickhouse.Client, query url.Values, serviceName string) *View {
	v := &View{}

	var q Query
	if s := query.Get("query"); s != "" {
		if err := json.Unmarshal([]byte(s), &q); err != nil {
			klog.Warningf("Failed to parse query: %v", err)
		}
	}
	if q.Limit <= 0 {
		q.Limit = defaultLimit
	}

	if ch == nil {
		v.Status = model.UNKNOWN
		v.Message = "ClickHouse integration is not configured"
		return v
	}

	from := w.Ctx.From.ToStandard()
	to := w.Ctx.To.ToStandard()

	rows, badge, browserStats, err := getPerformanceData(ctx, ch, serviceName, from, to)
	if err != nil {
		klog.Errorf("Failed to get performance data: %v", err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("ClickHouse error: %s", err)
		return v
	}

	v.Overviews = rows
	v.BadgeView = badge
	v.BrowserStats = browserStats
	v.Limit = q.Limit

	lineChartReport, err := createLineCharts(w, ctx, ch, serviceName)
	if err != nil {
		klog.Errorf("Failed to create line charts: %v", err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Failed to create line charts: %s", err)
		return v
	}

	echartReport, err := createECharts(w, ctx, ch, serviceName, from, to, rows)
	if err != nil {
		klog.Errorf("Failed to create ECharts: %v", err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Failed to create ECharts: %s", err)
		return v
	}

	v.Status = model.OK
	v.Report = *lineChartReport
	v.EchartReport = *echartReport
	return v
}

func getPerformanceData(ctx context.Context, ch *clickhouse.Client, serviceName string, from, to time.Time) ([]PerfOverview, Badge, []BrowserStats, error) {
	rows, err := ch.GetPerformanceOverview(ctx, &from, &to, serviceName)
	if err != nil {
		return nil, Badge{}, nil, err
	}

	// // Get browser stats
	perfBrowsers, err := ch.GetBrowserStats(ctx, serviceName, &from, &to)
	if err != nil {
		return nil, Badge{}, nil, fmt.Errorf("failed to get browser stats: %w", err)
	}
	browserStats := make([]BrowserStats, len(perfBrowsers))
	for i, stat := range perfBrowsers {
		browserStats[i] = BrowserStats{
			Name:         stat.Name,
			Requests:     stat.Requests,
			ResponseTime: stat.ResponseTime,
			Errors:       stat.Errors,
		}
	}
	// browserStats := []BrowserStats{
	// 	{"Chrome", 12000, 1.23, 150},
	// 	{"Firefox", 8500, 1.45, 120},
	// 	{"Safari", 7300, 1.30, 100},
	// 	{"Edge", 6500, 1.50, 90},
	// 	{"Opera", 4000, 1.35, 60},
	// }

	var overviews []PerfOverview
	var totalRequests, totalErrors uint64
	for _, row := range rows {
		overview := PerfOverview{
			PagePath:           row.PagePath,
			AvgLoadPageTime:    row.AvgLoadPageTime,
			JsErrorPercentage:  row.JsErrorPercentage,
			ApiErrorPercentage: row.ApiErrorPercentage,
			ImpactedUsers:      row.ImpactedUsers,
			Requests:           row.Requests,
		}
		overviews = append(overviews, overview)
		totalRequests += row.Requests
		totalErrors += uint64(row.JsErrorPercentage*float64(row.Requests)/100 + row.ApiErrorPercentage*float64(row.Requests)/100)
	}

	sort.Slice(overviews, func(i, j int) bool {
		return overviews[i].PagePath < overviews[j].PagePath
	})

	durationSeconds := to.Sub(from).Seconds()
	var requestsPerSecond, errorsPerSecond float64
	if durationSeconds > 0 {
		requestsPerSecond = float64(totalRequests) / durationSeconds
		errorsPerSecond = float64(totalErrors) / durationSeconds
	}

	var requestTrend, errorTrend float64
	if duration := to.Sub(from); duration == time.Hour {
		prevFrom := from.Add(-duration)
		prevRows, err := ch.GetPerformanceOverview(ctx, &prevFrom, &from, serviceName)
		if err != nil {
			return nil, Badge{}, nil, err
		}
		var prevRequests, prevErrors uint64
		for _, row := range prevRows {
			prevRequests += row.Requests
			prevErrors += uint64(row.JsErrorPercentage*float64(row.Requests)/100 + row.ApiErrorPercentage*float64(row.Requests)/100)
		}
		if prevRequests > 0 {
			requestTrend = float64(totalRequests-prevRequests) / float64(prevRequests) * 100
		}
		if prevErrors > 0 {
			errorTrend = float64(totalErrors-prevErrors) / float64(prevErrors) * 100
		}
	}

	badge := Badge{
		TotalRequests:     totalRequests,
		TotalErrors:       totalErrors,
		RequestTrend:      requestTrend,
		ErrorTrend:        errorTrend,
		RequestsPerSecond: requestsPerSecond,
		ErrorsPerSecond:   errorsPerSecond,
	}

	return overviews, badge, browserStats, nil
}

func createECharts(w *model.World, ctx context.Context, ch *clickhouse.Client, serviceName string, from, to time.Time, overviews []PerfOverview) (*model.AuditReport, error) {
	report := model.NewAuditReport(nil, w.Ctx, nil, model.AuditReportPerformance, true)
	// Impacted users pie chart
	sort.Slice(overviews, func(i, j int) bool {
		return overviews[i].ImpactedUsers > overviews[j].ImpactedUsers
	})
	impactedData := make([]model.DataPoint, len(overviews))
	var totalImpacted uint64
	for i, overview := range overviews {
		impactedData[i] = model.DataPoint{
			Value: int(overview.ImpactedUsers),
			Name:  overview.PagePath,
		}
		totalImpacted += overview.ImpactedUsers
	}

	impactedChart := report.GetOrCreateEChart("Impacted Users by Page", nil)
	impactedChart.Title = model.TextTitle{Text: "Impacted Users by Page"}
	impactedChart.Tooltip = model.Tooltip{Trigger: "item"}
	impactedChart.Legend = model.Legend{Bottom: "0"}
	impactedChart.SetPieChartSeries("Users", "pie", []string{"40%", "70%"}, impactedData)
	impactedChart.Graphic = &model.Graphic{
		Type: "text",
		Left: "center",
		Top:  "center",
		Style: model.GraphicStyle{Text: fmt.Sprintf("%d", totalImpacted), FontSize: 26,
			FontWeight: "bold",
			Fill:       "#333",
			TextAlign:  "center"},
	}

	return report, nil
}

func createLineCharts(w *model.World, ctx context.Context, ch *clickhouse.Client, serviceName string) (*model.AuditReport, error) {
	report := model.NewAuditReport(nil, w.Ctx, nil, model.AuditReportPerformance, true)
	report.Status = model.OK

	loadChart := report.GetOrCreateChart("Load Time", nil).Stacked()
	loadSeries, err := ch.GetLoadTimeSeries(ctx, serviceName, w.Ctx.From, w.Ctx.To, w.Ctx.Step)
	if err != nil {
		return nil, fmt.Errorf("failed to get load time series: %w", err)
	}
	loadChart.AddSeries("Load Time", loadSeries, "blue")

	respChart := report.GetOrCreateChart("Response Time", nil).Stacked()
	respSeries, err := ch.GetResponseTimeSeries(ctx, serviceName, w.Ctx.From, w.Ctx.To, w.Ctx.Step)
	if err != nil {
		return nil, fmt.Errorf("failed to get response time series: %w", err)
	}
	respChart.AddSeries("Response Time", respSeries, "green")

	metrics, err := ch.GetErrorAndUsersImpactedSeries(context.Background(), serviceName, w.Ctx.From, w.Ctx.To, w.Ctx.Step)
	if err != nil {
		report.Status = model.UNKNOWN
		return report, err
	}

	errorChartGroup := report.GetOrCreateChartGroup("Errors <selector>", nil)
	errorChartGroup.GetOrCreateChart("Errors").Stacked().
		AddSeries("JS Errors", metrics["jsErrors"], "orange").
		AddSeries("API Errors", metrics["apiErrors"], "purple")

	usersImpactedChart := report.GetOrCreateChart("Users Impacted", nil).Stacked().Column()
	usersImpactedChart.AddSeries("Total Users", metrics["totalUsers"], "green")
	usersImpactedChart.AddSeries("Users Impacted", metrics["usersImpacted"], "red")

	return report, nil
}
