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
	Totals    Total             `json:"totals"`
	Trends    Trend             `json:"trends"`
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

type Total struct {
	TotalRequests uint64  `json:"totalRequests"`
	TotalErrors   uint64  `json:"totalErrors"`
	RequestPerSec float64 `json:"requestPerSec"`
	ErrorPerSec   float64 `json:"errorPerSec"`
}

type Trend struct {
	RequestUpTrend bool    `json:"requestUpTrend"`
	RequestTrend   float64 `json:"requestTrend"`
	ErrorUpTrend   bool    `json:"errorUpTrend"`
	ErrorTrend     float64 `json:"errorTrend"`
}

func renderEumApps(ctx context.Context, ch *clickhouse.Client, w *model.World, query string) *EumView {
	v := &EumView{}

	from := w.Ctx.From.ToStandard()
	to := w.Ctx.To.ToStandard()

	// Get service overviews
	overviews, err := getServiceOverviews(ctx, ch, from, to)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Clickhouse error: %s", err)
		return v
	}

	v.Overviews = overviews

	// Calculate totals and trends
	totals, trends, err := calculateTotalsAndTrends(ctx, ch, from, to)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Clickhouse error: %s", err)
		return v
	}

	v.Totals = totals
	v.Trends = trends

	v.Status = model.OK
	return v
}

func getServiceOverviews(ctx context.Context, ch *clickhouse.Client, from, to time.Time) ([]ServiceOverview, error) {
	rows, err := ch.GetServiceOverviews(ctx, &from, &to)
	if err != nil {
		return nil, err
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
			Requests:           row.Requests,
			AppType:            row.AppType,
		})
	}

	sort.Slice(overviews, func(i, j int) bool {
		return overviews[i].ServiceName < overviews[j].ServiceName
	})

	return overviews, nil
}

func calculateTotalsAndTrends(ctx context.Context, ch *clickhouse.Client, from, to time.Time) (Total, Trend, error) {
	requests, err := ch.GetTotalRequests(ctx, &from, &to, "", "")
	if err != nil {
		return Total{}, Trend{}, err
	}
	errors, err := ch.GetTotalErrors(ctx, &from, &to, "", "")
	if err != nil {
		return Total{}, Trend{}, err
	}

	duration := to.Sub(from)
	totals := Total{
		TotalRequests: requests,
		TotalErrors:   errors,
		RequestPerSec: float64(requests) / duration.Seconds(),
		ErrorPerSec:   float64(errors) / duration.Seconds(),
	}

	var trends Trend
	if duration == time.Hour {
		previousHour := from.Add(-1 * time.Hour)
		previousTotalRequests, err := ch.GetTotalRequests(ctx, &previousHour, &from, "", "")
		if err != nil {
			return Total{}, Trend{}, err
		}

		previousTotalErrors, err := ch.GetTotalErrors(ctx, &previousHour, &from, "", "")
		if err != nil {
			return Total{}, Trend{}, err
		}

		trends = Trend{
			RequestUpTrend: requests > previousTotalRequests,
			RequestTrend:   float64(requests-previousTotalRequests) / float64(previousTotalRequests) * 100,
			ErrorUpTrend:   errors > previousTotalErrors,
			ErrorTrend:     float64(errors-previousTotalErrors) / float64(previousTotalErrors) * 100,
		}
	}

	return totals, trends, nil
}
