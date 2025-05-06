package overview

import (
	"context"
	"time"

	"codexray/clickhouse"
	"codexray/model"

	"k8s.io/klog"
)

type DashboardView struct {
	Status       model.Status        `json:"status"`
	Message      string              `json:"message"`
	Applications ApplicationOverview `json:"applications"`
	EumOverview  EumOverview         `json:"eumOverview"`
	Nodes        NodeOverview        `json:"nodes"`
	Incidents    IncidentOverview    `json:"incidents"`
}

type ApplicationOverview struct {
	ApplicationStatus ApplicationsStats  `json:"applicationsStats"`
	Applications      []ApplicationTable `json:"applications"`
}

type EumOverview struct {
	EumApps   []EumTable `json:"eumOverview"`
	BadgeView EumBadge   `json:"badgeView"`
}

type NodeOverview struct {
	NodeStats NodeStats    `json:"nodeStats"`
	Nodes     []NodesTable `json:"nodes"`
}

type IncidentOverview struct {
	IncidentStats IncidentStats   `json:"incidentStats"`
	Incidents     []IncidentTable `json:"incidents"`
}

type ApplicationTable struct {
	ID                   string  `json:"id"`
	Status               string  `json:"status"`
	TransactionPerSecond float64 `json:"transactionPerSecond"`
	ResponseTime         float64 `json:"responseTime"`
	Errors               uint64  `json:"errors"`
}

type ApplicationsStats struct {
	Total int64 `json:"total"`
	Good  int64 `json:"good"`
	Fair  int64 `json:"fair"`
	Poor  int64 `json:"poor"`
}

type EumTable struct {
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

type NodeStats struct {
	TotalNodes     int64   `json:"totalNodes"`
	UpNodes        int64   `json:"upNodes"`
	DownNodes      int64   `json:"downNodes"`
	AvgCPUUsage    float64 `json:"avgCpuUsage"`
	AvgMemoryUsage float64 `json:"avgMemoryUsage"`
	AvgDiskUsage   float64 `json:"avgDiskUsage"`
}

type NodesTable struct {
	NodeName    string  `json:"nodeName"`
	NodeStatus  string  `json:"nodeStatus"`
	CpuUsage    float64 `json:"cpuUsage"`
	MemoryUsage float64 `json:"memoryUsage"`
	DiskUsage   float64 `json:"diskUsage"`
}

type IncidentStats struct {
	TotalIncidents    int64 `json:"totalIncidents"`
	CriticalIncidents int64 `json:"criticalIncidents"`
	WarningIncidents  int64 `json:"warningIncidents"`
	ClosedIncidents   int64 `json:"closedIncidents"`
}

type IncidentTable struct {
	ApplicationName string    `json:"applicationName"`
	OpenIncidents   int64     `json:"openIncidents"`
	LastOccurence   time.Time `json:"lastOccurrence"`
}

func renderDashboard(ctx context.Context, ch *clickhouse.Client, w *model.World) *DashboardView {
	v := &DashboardView{}

	from := w.Ctx.From.ToStandard()
	to := w.Ctx.To.ToStandard()

	// EUM Overview
	var eumOverview EumOverview
	eumApps, badge, err := getEumOverviews(ctx, ch, from, to)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = err.Error()
		return v
	}
	eumOverview.EumApps = eumApps
	eumOverview.BadgeView = badge
	v.EumOverview = eumOverview
	v.Status = model.OK
	return v
}

func getEumOverviews(ctx context.Context, ch *clickhouse.Client, from, to time.Time) ([]EumTable, EumBadge, error) {
	rows, err := ch.GetEUMOverview(ctx, &from, &to)
	if err != nil {
		return nil, EumBadge{}, err
	}

	var eumTable []EumTable
	durationSeconds := to.Sub(from).Seconds()

	for _, row := range rows {
		var requestsPerSecond float64
		if durationSeconds > 0 {
			requestsPerSecond = float64(row.Requests) / durationSeconds
		}

		eumTable = append(eumTable, EumTable{
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
		BrowserApps: appCounts.BrowserApps,
		MobileApps:  appCounts.MobileApps,
	}

	return eumTable, badge, nil
}
