package overview

import (
	"codexray/auditor"
	"codexray/clickhouse"
	"codexray/model"
	"codexray/utils"
	"context"
	"fmt"
	"strings"

	"k8s.io/klog"
)

type MrumPerfView struct {
	Status  model.Status       `json:"status"`
	Message string             `json:"message"`
	Summary MrumPerfData       `json:"summary"`
	Report  *model.AuditReport `json:"report"`
}

type MrumPerfData struct {
	TotalRequests          uint64  `json:"totalRequests"`
	RequestsPerSecond      float64 `json:"requestsPerSecond"`
	RequestsTrend          float64 `json:"requestsTrend"`
	TotalErrors            uint64  `json:"totalErrors"`
	ErrorsPerSecond        float64 `json:"errorsPerSecond"`
	ErrorsTrend            float64 `json:"errorsTrend"`
	UsersImpacted          uint64  `json:"usersImpacted"`
	UsersImpactedPerSecond float64 `json:"usersImpactedPerSecond"`
	UsersImpactedTrend     float64 `json:"usersImpactedTrend"`
}

func RenderMrumPerf(ctx context.Context, ch *clickhouse.Client, w *model.World, query string) *MrumPerfView {
	v := &MrumPerfView{}
	
	// Handle empty query case by setting default interval to "now-1hr"
	if query == "" {
		query = "now-1hr"
	}
	
	parts := strings.Split(query, "-")
	fromStdTime := utils.ParseTime(w.Ctx.To, parts[0], w.Ctx.From)
	toStdTime := utils.ParseTime(w.Ctx.To, parts[1], w.Ctx.To)
	rows, err := ch.GetMobilePerfResults(ctx, &fromStdTime, &toStdTime)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Clickhouse error: %s", err)
		return v
	}

	v.Summary = MrumPerfData{
		TotalRequests:          rows.TotalRequests,
		RequestsPerSecond:      rows.RequestsPerSecond,
		RequestsTrend:          rows.RequestsTrend,
		TotalErrors:            rows.TotalErrors,
		ErrorsPerSecond:        rows.ErrorsPerSecond,
		ErrorsTrend:            rows.ErrorsTrend,
		UsersImpacted:          rows.UsersImpacted,
		UsersImpactedPerSecond: rows.UsersImpactedPerSecond,
		UsersImpactedTrend:     rows.UsersImpactedTrend,
	}

	v.Report = auditor.GenerateMrumPerfReport(w, ch, &fromStdTime, &toStdTime)
	v.Status = model.OK
	return v
}
