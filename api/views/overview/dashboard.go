package overview

import (
	"context"
	"time"

	"codexray/clickhouse"
	"codexray/model"

	"k8s.io/klog"
)

type DashboardView struct {
	Status    model.Status  `json:"status"`
	Message   string        `json:"message"`
	EumApps   []EumOverview `json:"eumOverview"`
	BadgeView EumBadge      `json:"badgeView"`
}

type EumOverview struct {
	ServiceName       string  `json:"serviceName" ch:"ServiceName"`
	AppType           string  `json:"appType" ch:"AppType"`
	RequestsPerSecond float64 `json:"requestsPerSecond" ch:"requestsPerSecond"`
	ResponseTime      float64 `json:"responseTime" ch:"responseTime"`
	Errors            uint64  `json:"errors" ch:"errors"`
	AffectedUsers     uint64  `json:"affectedUsers" ch:"affectedUsers"`
}

type EumBadge struct {
	BrowserApps uint64 `json:"browserApps" ch:"browserApps"`
	MobileApps  uint64 `json:"mobileApps" ch:"mobileApps"`
}

func renderDashboard(ctx context.Context, ch *clickhouse.Client, w *model.World) *DashboardView {
	v := &DashboardView{}

	from := w.Ctx.From.ToStandard()
	to := w.Ctx.To.ToStandard()

	eumApps, badge, err := getEumOverviews(ctx, ch, from, to)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = err.Error()
		return v
	}

	v.EumApps = eumApps
	v.BadgeView = badge
	v.Status = model.OK
	return v
}

func getEumOverviews(ctx context.Context, ch *clickhouse.Client, from, to time.Time) ([]EumOverview, EumBadge, error) {
	rows, err := ch.GetEUMOverview(ctx, &from, &to)
	if err != nil {
		return nil, EumBadge{}, err
	}

	var eumOverviews []EumOverview
	durationSeconds := to.Sub(from).Seconds()

	for _, row := range rows {
		var requestsPerSecond float64
		if durationSeconds > 0 {
			requestsPerSecond = float64(row.Requests) / durationSeconds
		}

		eumOverviews = append(eumOverviews, EumOverview{
			ServiceName:       row.ServiceName,
			AppType:           row.AppType,
			RequestsPerSecond: requestsPerSecond,
			ResponseTime:      row.ResponseTime,
			Errors:            row.Errors,
			AffectedUsers:     row.AffectedUsers,
		})

	}
	appCounts, err := ch.GetAppCounts(ctx, &from, &to)
	if err != nil {
		return nil, EumBadge{}, err
	}

	badge := EumBadge{
		BrowserApps: appCounts.browserApps,
		MobileApps:  appCounts.mobileApps,
	}

	return eumOverviews, badge, nil
}
