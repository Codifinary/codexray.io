package overview

import (
	"context"
	"fmt"
	"sort"
	"time"

	"codexray/clickhouse"
	"codexray/model"
	"codexray/timeseries"

	"k8s.io/klog"
)

const defaultLimit = 100

type PerfView struct {
	Status    model.Status   `json:"status"`
	Message   string         `json:"message"`
	Overviews []PerfOverview `json:"overviews"`
	Charts    []*model.Chart `json:"charts,omitempty"`
	Limit     int            `json:"limit"`
}

type PerfQuery struct {
	From  *timeseries.Time `json:"from"`
	To    *timeseries.Time `json:"to"`
	Step  int64            `json:"step"`
	Limit int              `json:"limit"`
}

type PerfOverview struct {
	PagePath           string  `json:"pagePath"`
	AvgLoadPageTime    float64 `json:"avgLoadPageTime"`
	JsErrorPercentage  float64 `json:"jsErrorPercentage"`
	ApiErrorPercentage float64 `json:"apiErrorPercentage"`
	ImpactedUsers      uint64  `json:"impactedUsers"`
}

func renderPerfs(ctx context.Context, ch *clickhouse.Client, w *model.World, query string) *PerfView {
	v := &PerfView{}

	var q PerfQuery

	if q.Limit <= 0 {
		q.Limit = defaultLimit
	}

	if ch == nil {
		v.Status = model.UNKNOWN
		v.Message = "Clickhouse integration is not configured"
		return v
	}

	// Fetch performance data
	rows, err := ch.GetPerformanceOverview(ctx, &time.Time{}, &time.Time{})
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Clickhouse error: %s", err)
		return v
	}

	var overviews []PerfOverview
	for _, row := range rows {
		overviews = append(overviews, PerfOverview{
			PagePath:           row.PagePath,
			AvgLoadPageTime:    row.AvgLoadPageTime,
			JsErrorPercentage:  row.JsErrorPercentage,
			ApiErrorPercentage: row.ApiErrorPercentage,
			ImpactedUsers:      row.ImpactedUsers,
		})
	}

	sort.Slice(overviews, func(i, j int) bool {
		return overviews[i].PagePath < overviews[j].PagePath
	})

	v.Status = model.OK
	v.Overviews = overviews
	v.Limit = q.Limit

	return v
}
