package perf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"time"

	"codexray/clickhouse"
	"codexray/model"

	"k8s.io/klog"
)

const defaultLimit = 100

type View struct {
	Status    model.Status   `json:"status"`
	Message   string         `json:"message"`
	Overviews []PerfOverview `json:"overviews"`
	Limit     int            `json:"limit"`
}

type Query struct {
	Limit int `json:"limit"`
}

type PerfOverview struct {
	PagePath           string  `json:"pagePath"`
	AvgLoadPageTime    float64 `json:"avgLoadPageTime"`
	JsErrorPercentage  float64 `json:"jsErrorPercentage"`
	ApiErrorPercentage float64 `json:"apiErrorPercentage"`
	ImpactedUsers      uint64  `json:"impactedUsers"`
}

func Render(ctx context.Context, ch *clickhouse.Client, query url.Values, from, to *time.Time) *View {
	v := &View{}

	var q Query
	if s := query.Get("query"); s != "" {
		if err := json.Unmarshal([]byte(s), &q); err != nil {
			klog.Warningln(err)
		}
	}
	if q.Limit <= 0 {
		q.Limit = defaultLimit
	}

	// Check Clickhouse client
	if ch == nil {
		v.Status = model.UNKNOWN
		v.Message = "Clickhouse integration is not configured"
		return v
	}

	// Fetch performance data
	rows, err := ch.GetPerformanceOverview(ctx, from, to)
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
