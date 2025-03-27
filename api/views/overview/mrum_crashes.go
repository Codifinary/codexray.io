package overview

import (
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
}

type MrumCrashesView struct {
	Status                   model.Status                           `json:"status"`
	Message                  string                                 `json:"message"`
	Summary                  MrumCrashesData                        `json:"summary"`
	Report                   *model.AuditReport                     `json:"report"`
	CrashReasonWiseOverviews []clickhouse.CrashesReasonwiseOverview `json:"crashReasonWiseOverview"`
	CrashDatabyCrashReason   []clickhouse.CrashReasonData           `json:"crashDatabyCrashReason"`
}

type MrumCrashesData struct {
	TotalCrashes uint64 `json:"totalCrashes"`
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
		crashData, err := ch.GetCrashReasonData(ctx, q.CrashReason, w.Ctx.From, w.Ctx.To)
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

		crashReasonWiseOverviews, err := ch.GetCrashesReasonwiseOverview(ctx, w.Ctx.From, w.Ctx.To)
		if err != nil {
			klog.Errorln(err)
			v.Status = model.WARNING
			v.Message = fmt.Sprintf("Clickhouse error: %s", err)
			return v
		}
		v.CrashReasonWiseOverviews = crashReasonWiseOverviews
	}

	v.Status = model.OK
	return v
}

func parseCrashQuery(query string, ctx timeseries.Context) CrashQuery {
	var res CrashQuery
	if query != "" {
		if err := json.Unmarshal([]byte(query), &res); err != nil {
			klog.Warningln(err)
		}
	}

	return res
}
