package overview

import (
	"codexray/auditor"
	"codexray/clickhouse"
	"codexray/model"
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog"
)

type MrumUsersView struct {
	Status         model.Status                 `json:"status"`
	Message        string                       `json:"message"`
	Summary        MrumUsersData                `json:"summary"`
	Report         json.RawMessage              `json:"report"`
	MobileUserData []clickhouse.MobileUsersData `json:"mobileUserData"`
}

type MrumUsersData struct {
	TotalUsers          uint64  `json:"totalUsers"`
	NewUsers            uint64  `json:"newUsers"`
	ReturningUsers      uint64  `json:"returningUsers"`
	DailyActiveUsers    uint64  `json:"dailyActiveUsers"`
	WeeklyActiveUsers   uint64  `json:"weeklyActiveUsers"`
	DailyTrend          float64 `json:"dailyTrend"`
	CrashFreePercentage float64 `json:"crashFreePercentage"`
	UserTrend           float64 `json:"userTrend"`
	NewUserTrend        float64 `json:"newUserTrend"`
	ReturningUserTrend  float64 `json:"returningUserTrend"`
}

func RenderMrumUsers(ctx context.Context, ch *clickhouse.Client, w *model.World, query string, service string) *MrumUsersView {
	v := &MrumUsersView{}

	if ch == nil {
		v.Status = model.WARNING
		v.Message = "clickhouse not available"
		return v
	}

	rows, err := ch.GetMobileUserResults(ctx, w.Ctx.From, w.Ctx.To, service)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Clickhouse error: %s", err)
		return v
	}

	v.Summary = MrumUsersData{
		TotalUsers:          rows.TotalUsers,
		NewUsers:            rows.NewUsers,
		ReturningUsers:      rows.ReturningUsers,
		DailyActiveUsers:    rows.DailyActiveUsers,
		WeeklyActiveUsers:   rows.WeeklyActiveUsers,
		DailyTrend:          rows.DailyTrend,
		CrashFreePercentage: rows.CrashFreePercentage,
		UserTrend:           rows.UserTrend,
		NewUserTrend:        rows.NewUserTrend,
		ReturningUserTrend:  rows.ReturningUserTrend,
	}

	originalReport := auditor.GenerateMrumUsersReport(w, ch, w.Ctx.From, w.Ctx.To, service)

	reportJSON, err := convertReportWithChartArray(originalReport)
	if err != nil {
		klog.Errorln("Error converting report to JSON:", err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Error converting report: %s", err)
		return v
	}
	v.Report = reportJSON

	mobileUserData, err := ch.GetMobileUsersData(ctx, w.Ctx.From, w.Ctx.To, service)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Clickhouse error: %s", err)
		return v
	}

	v.MobileUserData = mobileUserData
	v.Status = model.OK
	return v
}

func convertReportWithChartArray(report *model.AuditReport) (json.RawMessage, error) {
	originalJSON, err := json.Marshal(report)
	if err != nil {
		return nil, err
	}

	var reportMap map[string]interface{}
	if err := json.Unmarshal(originalJSON, &reportMap); err != nil {
		return nil, err
	}

	widgets, ok := reportMap["widgets"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("widgets not found or not an array")
	}

	for _, w := range widgets {
		widget, ok := w.(map[string]interface{})
		if !ok {
			continue
		}

		if chart, exists := widget["chart"]; exists && chart != nil {
			widget["chart"] = []interface{}{chart}
		}
	}

	modifiedJSON, err := json.Marshal(reportMap)
	if err != nil {
		return nil, err
	}

	return modifiedJSON, nil
}
