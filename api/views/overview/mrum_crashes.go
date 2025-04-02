package overview

import (
	"codexray/auditor"
	"codexray/clickhouse"
	"codexray/model"
	"codexray/timeseries"
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog"
)

type CrashQuery struct {
	CrashReason string `json:"crash_reason"`
	Limit       int    `json:"limit"`
}

type MrumCrashesView struct {
	Status                   model.Status                           `json:"status"`
	Message                  string                                 `json:"message"`
	Summary                  MrumCrashesData                        `json:"summary"`
	Report                   *model.AuditReport                     `json:"report"`
	EchartReport             *model.AuditReport                     `json:"echartReport"`
	CrashReasonWiseOverviews []clickhouse.CrashesReasonwiseOverview `json:"crashReasonWiseOverview"`
	CrashDatabyCrashReason   []clickhouse.CrashReasonData           `json:"crashDatabyCrashReason"`
}

type MrumCrashesData struct {
	TotalCrashes uint64 `json:"totalCrashes"`
}

type CrashPieChartConfig struct {
	Title      string
	ChartType  string
	SeriesName string
	Colors     []string
}

func RenderMrumCrashes(ctx context.Context, ch *clickhouse.Client, w *model.World, query string) *MrumCrashesView {
	v := &MrumCrashesView{}

	if ch == nil {
		v.Status = model.WARNING
		v.Message = "clickhouse not available"
		return v
	}

	q := parseCrashQuery(query, w.Ctx)

	switch {
	case q.CrashReason != "":
		crashData, err := ch.GetCrashReasonData(ctx, q.CrashReason, w.Ctx.From, w.Ctx.To, q.Limit)
		if err != nil {
			klog.Errorln(err)
			v.Status = model.WARNING
			v.Message = fmt.Sprintf("Clickhouse error: %s", err)
			return v
		}
		v.CrashDatabyCrashReason = crashData

	default:
		rows, err := ch.GetMobileCrashesResults(ctx, w.Ctx.From, w.Ctx.To)
		if err != nil {
			klog.Errorln(err)
			v.Status = model.WARNING
			v.Message = fmt.Sprintf("Clickhouse error: %s", err)
			return v
		}

		v.Summary = MrumCrashesData{
			TotalCrashes: rows.TotalCrashes,
		}

		crashReasonWiseOverviews, err := ch.GetCrashesReasonwiseOverview(ctx, w.Ctx.From, w.Ctx.To, q.Limit)
		if err != nil {
			klog.Errorln(err)
			v.Status = model.WARNING
			v.Message = fmt.Sprintf("Clickhouse error: %s", err)
			return v
		}
		v.CrashReasonWiseOverviews = crashReasonWiseOverviews

		v.Report = auditor.GenerateMrumCrashesReport(w, ch, w.Ctx.From, w.Ctx.To)

		commonColors := []string{"#4169E1", "#6495ED", "#1E90FF", "#00BFFF", "#87CEEB"}

		deviceConfig := CrashPieChartConfig{
			Title:      "Top Devices by Crash Count",
			ChartType:  "device",
			SeriesName: "Devices",
			Colors:     commonColors,
		}

		osConfig := CrashPieChartConfig{
			Title:      "Top OS by Crash Count",
			ChartType:  "os",
			SeriesName: "OS",
			Colors:     commonColors,
		}

		appVersionConfig := CrashPieChartConfig{
			Title:      "Top App Versions by Crash Count",
			ChartType:  "appVersion",
			SeriesName: "App Versions",
			Colors:     commonColors,
		}

		v.EchartReport, err = createCrashPieCharts(ctx, ch, w, []CrashPieChartConfig{deviceConfig, osConfig, appVersionConfig})
		if err != nil {
			klog.Errorln(err)
			v.Status = model.WARNING
			v.Message = fmt.Sprintf("Failed to create crash pie charts: %s", err)
			return v
		}
	}

	v.Status = model.OK
	return v
}

func parseCrashQuery(query string, ctx timeseries.Context) CrashQuery {
	var res CrashQuery
	res.Limit = 10
	if query != "" {
		if err := json.Unmarshal([]byte(query), &res); err != nil {
			klog.Warningln(err)
		}
	}

	return res
}

func createCrashPieCharts(ctx context.Context, ch *clickhouse.Client, w *model.World, configs []CrashPieChartConfig) (*model.AuditReport, error) {
	echartReport := model.NewAuditReport(nil, w.Ctx, nil, model.AuditReportMobileCrashes, true)

	for _, config := range configs {
		var data []struct {
			Name  string
			Value uint64
		}
		var err error

		switch config.ChartType {
		case "device":
			data, err = ch.GetTopDevicesByCrashCount(ctx, w.Ctx.From, w.Ctx.To)
			if err != nil {
				return nil, fmt.Errorf("failed to get top devices by crash count: %w", err)
			}
		case "os":
			data, err = ch.GetTopOSByCrashCount(ctx, w.Ctx.From, w.Ctx.To)
			if err != nil {
				return nil, fmt.Errorf("failed to get top OS by crash count: %w", err)
			}
		case "appVersion":
			data, err = ch.GetTopAppVersionsByCrashCount(ctx, w.Ctx.From, w.Ctx.To)
			if err != nil {
				return nil, fmt.Errorf("failed to get top app versions by crash count: %w", err)
			}
		default:
			return nil, fmt.Errorf("unknown chart type: %s", config.ChartType)
		}

		chartData := make([]model.DataPoint, len(data))
		for i, item := range data {
			chartData[i] = model.DataPoint{
				Value: int(item.Value),
				Name:  item.Name,
			}
		}

		donutChart := echartReport.GetOrCreateEChart(config.Title, nil)
		donutChart.Title = model.TextTitle{
			Text: config.Title,
			TextStyle: &model.TextStyle{
				FontSize:   16,
				FontWeight: "normal",
			},
		}
		donutChart.Tooltip = model.Tooltip{Trigger: "item"}
		donutChart.Legend = model.Legend{Right: "5%", Top: "10%"}
		donutChart.Color = config.Colors
		donutChart.SetPieChartSeries(config.SeriesName, "pie", []string{"40%", "70%"}, chartData)

		donutChart.Graphic = &model.Graphic{
			Type: "text",
			Left: "center",
			Top:  "center",
			Style: model.GraphicStyle{
				FontSize:   26,
				FontWeight: "bold",
				Fill:       "#333",
				TextAlign:  "center",
			},
		}
	}

	return echartReport, nil
}
