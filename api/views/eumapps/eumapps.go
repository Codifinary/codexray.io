package eumapps

import (
	"context"
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
	Status    model.Status      `json:"status"`
	Message   string            `json:"message"`
	Overviews []ServiceOverview `json:"overviews"`
	Limit     int               `json:"limit"`
}

type Query struct {
	Limit int `json:"limit"`
}

type ServiceOverview struct {
	ServiceName        string  `json:"serviceName"`
	Pages              uint64  `json:"pages"`
	AvgLoadPageTime    float64 `json:"avgLoadPageTime"`
	JsErrorPercentage  float64 `json:"jsErrorPercentage"`
	ApiErrorPercentage float64 `json:"apiErrorPercentage"`
	ImpactedUsers      uint64  `json:"impactedUsers"`
}

func Render(ctx context.Context, ch *clickhouse.Client, query url.Values, from, to *time.Time) *View {
	v := &View{}

	// Default time range
	if from == nil || to == nil {
		now := time.Now()
		if from == nil {
			defaultFrom := now.Add(-24 * time.Hour)
			from = &defaultFrom
		}
		if to == nil {
			defaultTo := now
			to = &defaultTo
		}
	}

	rows, err := ch.GetServiceOverviews(ctx, from, to)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Clickhouse error: %s", err)
		return v
	}

	var overviews []ServiceOverview
	for _, row := range rows {
		overviews = append(overviews, ServiceOverview{
			ServiceName:        row.ServiceName,
			Pages:              row.Pages,
			AvgLoadPageTime:    row.AvgLoadPageTime,
			JsErrorPercentage:  row.JsErrorPercentage,
			ApiErrorPercentage: row.ApiErrorPercentage,
			ImpactedUsers:      row.ImpactedUsers,
		})
	}

	sort.Slice(overviews, func(i, j int) bool {
		return overviews[i].ServiceName < overviews[j].ServiceName
	})

	v.Status = model.OK
	v.Overviews = overviews
	return v
}
