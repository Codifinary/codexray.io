package overview

import (
	"codexray/auditor"
	"codexray/clickhouse"
	"codexray/model"
	"context"
	"fmt"

	"k8s.io/klog"
)

type MrumSessionsView struct {
	Status              model.Status                     `json:"status"`
	Message             string                           `json:"message"`
	Summary             MrumSessionsData                 `json:"summary"`
	Report              *model.AuditReport               `json:"report"`
	SessionLiveData     []clickhouse.SessionLiveData     `json:"sessionLiveData"`
	SessionHistoricData []clickhouse.SessionHistoricData `json:"sessionHistoricData"`
}

type MrumSessionsData struct {
	TotalSessions       uint64  `json:"totalSessions"`
	SessionTrend        float64 `json:"sessionTrend"`
	TotalUsers          uint64  `json:"totalUsers"`
	UserTrend           float64 `json:"userTrend"`
	AverageSession      float64 `json:"avgSession"`
	AverageSessionTrend float64 `json:"avgSessionTrend"`
}

func RenderMrumSessions(ctx context.Context, ch *clickhouse.Client, w *model.World, query string) *MrumSessionsView {
	v := &MrumSessionsView{}

	if ch == nil {
		v.Status = model.WARNING
		v.Message = "clickhouse not available"
		return v
	}

	rows, err := ch.GetMobileSessionResults(ctx, w.Ctx.From, w.Ctx.To)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Clickhouse error: %s", err)
		return v
	}

	v.Summary = MrumSessionsData{
		TotalSessions:       rows.TotalSessions,
		SessionTrend:        rows.SessionTrend,
		TotalUsers:          rows.TotalUsers,
		UserTrend:           rows.UserTrend,
		AverageSession:      rows.AverageSession,
		AverageSessionTrend: rows.AverageSessionTrend,
	}

	v.Report = auditor.GenerateMrumSessionsReport(w, ch, w.Ctx.From, w.Ctx.To)

	sessionLiveData, err := ch.GetSessionLiveData(ctx, w.Ctx.From, w.Ctx.To)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Clickhouse error: %s", err)
		return v
	}

	sessionHistoricData, err := ch.GetSessionHistoricData(ctx, w.Ctx.From, w.Ctx.To)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Clickhouse error: %s", err)
		return v
	}
	v.SessionLiveData = sessionLiveData
	v.SessionHistoricData = sessionHistoricData
	v.Status = model.OK
	return v
}
