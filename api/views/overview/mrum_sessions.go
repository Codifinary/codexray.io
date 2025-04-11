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

type MrumSessionsView struct {
	Status              model.Status                     `json:"status"`
	Message             string                           `json:"message"`
	Summary             MrumSessionsData                 `json:"summary"`
	Report              *model.AuditReport               `json:"report"`
	SessionLiveData     []clickhouse.SessionLiveData     `json:"sessionLiveData"`
	SessionHistoricData []clickhouse.SessionHistoricData `json:"sessionHistoricData"`
	SessionGeoMapData   []clickhouse.SessionGeoMapData   `json:"sessionGeoMapData"`
}

type MrumSessionsData struct {
	TotalSessions       uint64  `json:"totalSessions"`
	SessionTrend        float64 `json:"sessionTrend"`
	TotalUsers          uint64  `json:"totalUsers"`
	UserTrend           float64 `json:"userTrend"`
	AverageSession      float64 `json:"avgSession"`
	AverageSessionTrend float64 `json:"avgSessionTrend"`
}

type SessionQuery struct {
	SessionType string `json:"session_type"`
	Limit       int    `json:"limit"`
}

func RenderMrumSessions(ctx context.Context, ch *clickhouse.Client, w *model.World, query string, service string) *MrumSessionsView {
	v := &MrumSessionsView{}

	if ch == nil {
		v.Status = model.WARNING
		v.Message = "clickhouse not available"
		return v
	}

	q := parseSessionQuery(query, w.Ctx)

	rows, err := ch.GetMobileSessionResults(ctx, w.Ctx.From, w.Ctx.To, service)
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

	v.Report = auditor.GenerateMrumSessionsReport(w, ch, w.Ctx.From, w.Ctx.To, service)
	sessionGeoMapData, err := ch.GetSessionGeoMapData(ctx, w.Ctx.From, w.Ctx.To, service)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Clickhouse error: %s", err)
		return v
	}
	v.SessionGeoMapData = sessionGeoMapData

	switch {
	case q.SessionType != "":
		sessionHistoricData, err := ch.GetSessionHistoricData(ctx, w.Ctx.From, w.Ctx.To, q.Limit, service)
		if err != nil {
			klog.Errorln(err)
			v.Status = model.WARNING
			v.Message = fmt.Sprintf("Clickhouse error: %s", err)
			return v
		}
		v.SessionHistoricData = sessionHistoricData

	default:
		sessionLiveData, err := ch.GetSessionLiveData(ctx, w.Ctx.From, w.Ctx.To, q.Limit, service)
		if err != nil {
			klog.Errorln(err)
			v.Status = model.WARNING
			v.Message = fmt.Sprintf("Clickhouse error: %s", err)
			return v
		}
		v.SessionLiveData = sessionLiveData
	}

	v.Status = model.OK
	return v
}

func parseSessionQuery(query string, ctx timeseries.Context) SessionQuery {
	var res SessionQuery
	res.Limit = 10
	if query != "" {
		if err := json.Unmarshal([]byte(query), &res); err != nil {
			klog.Warningln(err)
		}
	}
	return res
}
