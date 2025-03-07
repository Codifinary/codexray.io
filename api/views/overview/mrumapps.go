package overview

import (
	"codexray/clickhouse"
	"codexray/model"
	"context"
	"fmt"
	"strconv"

	"k8s.io/klog"
)

type MrumView struct {
	Status  model.Status `json:"status"`
	Message string       `json:"message"`
	Data    MrumData     `json:"data"`
}

type MrumData struct {
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

func renderMrumApps(ctx context.Context, ch *clickhouse.Client, w *model.World, query string) *MrumView {
	v := &MrumView{}

	from := w.Ctx.From.ToStandard()
	to := w.Ctx.To.ToStandard()

	rows, err := ch.GetMobilePerfResults(ctx, &from, &to)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Clickhouse error: %s", err)
		return v
	}

	v.Data = MrumData{
		TotalRequests:          rows.TotalRequests,
		RequestsPerSecond:      float64(mustParseFloat(fmt.Sprintf("%.3f", rows.RequestsPerSecond))),
		RequestsTrend:          float64(mustParseFloat(fmt.Sprintf("%.3f", rows.RequestsTrend))),
		TotalErrors:            rows.TotalErrors,
		ErrorsPerSecond:        float64(mustParseFloat(fmt.Sprintf("%.3f", rows.ErrorsPerSecond))),
		ErrorsTrend:            float64(mustParseFloat(fmt.Sprintf("%.3f", rows.ErrorsTrend))),
		UsersImpacted:          rows.UsersImpacted,
		UsersImpactedPerSecond: float64(mustParseFloat(fmt.Sprintf("%.3f", rows.UsersImpactedPerSecond))),
		UsersImpactedTrend:     float64(mustParseFloat(fmt.Sprintf("%.3f", rows.UsersImpactedTrend))),
	}

	v.Status = model.OK
	return v
}

func mustParseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
