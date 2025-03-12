package overview

import (
	"context"
	"fmt"
	"sort"
	"time"

	"codexray/clickhouse"
	"codexray/model"

	"k8s.io/klog"
)

type EumView struct {
	Status    model.Status      `json:"status"`
	Message   string            `json:"message"`
	Overviews []ServiceOverview `json:"overviews"`
	BadgeView Badge             `json:"badgeView"`
	Limit     int               `json:"limit"`
}

type ServiceOverview struct {
	ServiceName        string  `json:"serviceName"`
	Pages              uint64  `json:"pages"`
	AvgLoadPageTime    float64 `json:"avgLoadPageTime"`
	JsErrorPercentage  float64 `json:"jsErrorPercentage"`
	ApiErrorPercentage float64 `json:"apiErrorPercentage"`
	ImpactedUsers      uint64  `json:"impactedUsers"`
	AppType            string  `json:"appType"`
	Requests           uint64  `json:"requests"`
}

type Badge struct {
	TotalApplications uint64  `json:"totalApplications"`
	TotalPages        uint64  `json:"totalPages"`
	AvgLatency        float64 `json:"avgLatency"`
	TotalErrors       uint64  `json:"totalError"`
	ErrorPerSec       float64 `json:"errorPerSec"`
	ErrorTrend        float64 `json:"errorTrend"`
}

func renderEumApps(ctx context.Context, ch *clickhouse.Client, w *model.World, query string) *EumView {
	v := &EumView{}

	from := w.Ctx.From.ToStandard()
	to := w.Ctx.To.ToStandard()

	overviews, badge, err := getServiceOverviews(ctx, ch, from, to)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Clickhouse error: %s", err)
		return v
	}

	v.Overviews = overviews
	v.BadgeView = badge

	v.Status = model.OK
	return v
}

func getServiceOverviews(ctx context.Context, ch *clickhouse.Client, from, to time.Time) ([]ServiceOverview, Badge, error) {
	rows, err := ch.GetServiceOverviews(ctx, &from, &to)
	if err != nil {
		return nil, Badge{}, err
	}

	var overviews []ServiceOverview
	var totalApplications, totalPages uint64
	var totalLatency float64

	for _, row := range rows {
		overviews = append(overviews, ServiceOverview{
			ServiceName:        row.ServiceName,
			Pages:              row.Pages,
			AvgLoadPageTime:    row.AvgLoadPageTime,
			JsErrorPercentage:  row.JsErrorPercentage,
			ApiErrorPercentage: row.ApiErrorPercentage,
			ImpactedUsers:      row.ImpactedUsers,
			Requests:           row.Requests,
			AppType:            row.AppType,
		})
		totalApplications++
		totalPages += row.Pages
		totalLatency += row.AvgLoadPageTime
	}

	sort.Slice(overviews, func(i, j int) bool {
		return overviews[i].ServiceName < overviews[j].ServiceName
	})

	// Calculate average latency
	avgLatency := totalLatency / float64(totalApplications)

	// Get total errors and error trend
	totalErrors, err := ch.GetTotalErrors(ctx, &from, &to, "", "")
	if err != nil {
		return nil, Badge{}, err
	}

	var errorTrend float64
	duration := to.Sub(from)
	if duration == time.Hour {
		previousFrom := from.Add(-duration)
		previousTotalErrors, err := ch.GetTotalErrors(ctx, &previousFrom, &from, "", "")
		if err != nil {
			return nil, Badge{}, err
		}
		errorTrend = float64(totalErrors-previousTotalErrors) / float64(previousTotalErrors) * 100
	} else {
		errorTrend = 0
	}

	badge := Badge{
		TotalApplications: totalApplications,
		TotalPages:        totalPages,
		AvgLatency:        avgLatency,
		TotalErrors:       totalErrors,
		ErrorTrend:        errorTrend,
	}

	return overviews, badge, nil
}
