package overview

import (
	"codexray/auditor"
	"codexray/clickhouse"
	"codexray/model"
	"context"
	"fmt"

	"k8s.io/klog"
)

type MrumUsersView struct {
	Status         model.Status                 `json:"status"`
	Message        string                       `json:"message"`
	Summary        MrumUsersData                `json:"summary"`
	Report         *model.AuditReport           `json:"report"`
	MobileUserData []clickhouse.MobileUsersData `json:"mobileUserData"`
}

type MrumUsersData struct {
	TotalUsers        uint64  `json:"totalUsers"`
	NewUsers          uint64  `json:"newUsers"`
	ReturningUsers    uint64  `json:"returningUsers"`
	DailyActiveUsers  uint64  `json:"dailyActiveUsers"`
	WeeklyActiveUsers uint64  `json:"weeklyActiveUsers"`
	DailyTrend        float64 `json:"dailyTrend"`
}

func RenderMrumUsers(ctx context.Context, ch *clickhouse.Client, w *model.World, query string) *MrumUsersView {
	v := &MrumUsersView{}

	if ch == nil {
		v.Status = model.WARNING
		v.Message = "clickhouse not available"
		return v
	}

	rows, err := ch.GetMobileUserResults(ctx, w.Ctx.From, w.Ctx.To)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Clickhouse error: %s", err)
		return v
	}

	v.Summary = MrumUsersData{
		TotalUsers:        rows.TotalUsers,
		NewUsers:          rows.NewUsers,
		ReturningUsers:    rows.ReturningUsers,
		DailyActiveUsers:  rows.DailyActiveUsers,
		WeeklyActiveUsers: rows.WeeklyActiveUsers,
		DailyTrend:        rows.DailyTrend,
	}

	v.Report = auditor.GenerateMrumUsersReport(w, ch, w.Ctx.From, w.Ctx.To)
	mobileUserData, err := ch.GetMobileUsersData(ctx, w.Ctx.From, w.Ctx.To)
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
