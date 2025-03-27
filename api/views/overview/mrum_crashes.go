package overview

import (
	"codexray/clickhouse"
	"codexray/model"
	"context"
	"fmt"

	"k8s.io/klog"
)

type MrumCrashesView struct {
	Status                   model.Status                           `json:"status"`
	Message                  string                                 `json:"message"`
	Summary                  MrumCrashesData                        `json:"summary"`
	Report                   *model.AuditReport                     `json:"report"`
	CrashReasonWiseOverviews []clickhouse.CrashesReasonwiseOverview `json:"crashReasonWiseOverview"`
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
	v.Status = model.OK
	return v
}
