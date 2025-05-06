package overview

import (
	"context"
	"sort"
	"time"

	"codexray/clickhouse"
	"codexray/model"
	"codexray/timeseries"

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
	AvgCPUUsage    float32 `json:"avgCpuUsage"`
	AvgMemoryUsage float32 `json:"avgMemoryUsage"`
	AvgDiskUsage   float32 `json:"avgDiskUsage"`
}

type NodesTable struct {
	NodeName    string  `json:"nodeName"`
	NodeStatus  string  `json:"nodeStatus"`
	CpuUsage    float32 `json:"cpuUsage"`
	MemoryUsage float32 `json:"memoryUsage"`
	DiskUsage   float32 `json:"diskUsage"`
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

func renderNode(w *model.World) NodeOverview {
	var totalNodes, upNodes, downNodes int64
	var totalCPU, totalMemory, totalDisk float32
	var cpuCount, memoryCount, diskCount int64

	var nodesTable []NodesTable

	for _, n := range w.Nodes {
		totalNodes++

		name := n.GetName()
		if name == "" {
			klog.Warningln("empty node name for", n.Id)
			continue
		}

		status := "up"
		switch {
		case !n.IsAgentInstalled():
			status = "no agent installed"
		case n.IsDown():
			status = "down"
			downNodes++
		default:
			upNodes++
		}

		var cpuUsage float32
		if l := n.CpuUsagePercent.Last(); !timeseries.IsNaN(l) {
			cpuUsage = l
			totalCPU += cpuUsage
			cpuCount++
		}

		var memoryUsage float32
		if total := n.MemoryTotalBytes.Last(); !timeseries.IsNaN(total) {
			if avail := n.MemoryAvailableBytes.Last(); !timeseries.IsNaN(avail) {
				memoryUsage = 100 - (avail / total * 100)
				totalMemory += memoryUsage
				memoryCount++
			}
		}

		var diskUsage float32
		var totalDiskSpace, usedDiskSpace float32

		for _, i := range n.Instances {
			for _, v := range i.Volumes {
				if capacity := v.CapacityBytes.Last(); !timeseries.IsNaN(capacity) {
					totalDiskSpace += capacity
				}
				if used := v.UsedBytes.Last(); !timeseries.IsNaN(used) {
					usedDiskSpace += used
				}
			}
		}

		if totalDiskSpace > 0 {
			diskUsage = (usedDiskSpace / totalDiskSpace) * 100
			totalDisk += diskUsage
			diskCount++
		}

		nodesTable = append(nodesTable, NodesTable{
			NodeName:    name,
			NodeStatus:  status,
			CpuUsage:    cpuUsage,
			MemoryUsage: memoryUsage,
			DiskUsage:   diskUsage,
		})
	}

	sort.Slice(nodesTable, func(i, j int) bool {
		return nodesTable[i].CpuUsage > nodesTable[j].CpuUsage
	})

	topNodes := nodesTable
	if len(nodesTable) > 5 {
		topNodes = nodesTable[:5]
	}

	var avgCPU, avgMemory, avgDisk float32
	if cpuCount > 0 {
		avgCPU = totalCPU / float32(cpuCount)
	}
	if memoryCount > 0 {
		avgMemory = totalMemory / float32(memoryCount)
	}
	if diskCount > 0 {
		avgDisk = totalDisk / float32(diskCount)
	}

	return NodeOverview{
		NodeStats: NodeStats{
			TotalNodes:     totalNodes,
			UpNodes:        upNodes,
			DownNodes:      downNodes,
			AvgCPUUsage:    avgCPU,
			AvgMemoryUsage: avgMemory,
			AvgDiskUsage:   avgDisk,
		},
		Nodes: topNodes,
	}
}
