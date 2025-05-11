package overview

import (
	"context"
	"sort"
	"time"

	"codexray/clickhouse"
	"codexray/model"
	"codexray/timeseries"

	"golang.org/x/exp/maps"
	"k8s.io/klog"
)

type DashboardView struct {
	Status             model.Status        `json:"status"`
	Message            string              `json:"message"`
	Applications       ApplicationOverview `json:"applications"`
	EumOverview        EumOverview         `json:"eumOverview"`
	Nodes              NodeOverview        `json:"nodes"`
	Incidents          IncidentOverview    `json:"incidents"`
	AppStatsChart      *model.EChart       `json:"appStatsChart,omitempty"`
	IncidentStatsChart *model.EChart       `json:"incidentStatsChart,omitempty"`
}

type ApplicationOverview struct {
	ApplicationStatus ApplicationsStats  `json:"applicationsStats"`
	Applications      []ApplicationTable `json:"applicationTable"`
}

type EumOverview struct {
	EumApps   []EumTable `json:"eumOverview"`
	BadgeView EumBadge   `json:"badgeView"`
}

type NodeOverview struct {
	NodeStats NodeStats    `json:"nodeStats"`
	Nodes     []NodesTable `json:"nodesTable"`
}

type IncidentOverview struct {
	IncidentStats IncidentStats   `json:"incidentStats"`
	Incidents     []IncidentTable `json:"incidentTable"`
}

type ApplicationTable struct {
	ID                   string       `json:"id"`
	Status               model.Status `json:"status"`
	TransactionPerSecond float32      `json:"transactionPerSecond"`
	ResponseTime         float32      `json:"responseTime"`
	Errors               float32      `json:"errors"`
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
	BrowserApps int `json:"browserApps" ch:"browserApps"`
	MobileApps  int `json:"mobileApps" ch:"mobileApps"`
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
	TotalIncidents    int `json:"totalIncidents"`
	CriticalIncidents int `json:"criticalIncidents"`
	WarningIncidents  int `json:"warningIncidents"`
	ClosedIncidents   int `json:"closedIncidents"`
}

type IncidentTable struct {
	Icon            string    `json:"icon"`
	ApplicationName string    `json:"applicationName"`
	OpenIncidents   int64     `json:"openIncidents"`
	LastOccurence   time.Time `json:"lastOccurrence"`
}

func renderDashboard(ctx context.Context, ch *clickhouse.Client, w *model.World) *DashboardView {
	v := &DashboardView{}

	applicationOverview := getApplications(w)

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

	nodesOverview := renderNode(w)

	incidentOverview := renderIncidents(w)

	report := model.NewAuditReport(nil, w.Ctx, nil, model.AuditReportPerformance, true)

	appStatsChart := renderApplicationsStatsDonutChart(applicationOverview.ApplicationStatus, report)
	incidentStatsChart := renderIncidentStatsDonutChart(incidentOverview.IncidentStats, report)

	v.Applications = applicationOverview
	v.EumOverview = eumOverview
	v.Nodes = nodesOverview
	v.Incidents = incidentOverview
	v.AppStatsChart = appStatsChart
	v.IncidentStatsChart = incidentStatsChart
	v.Status = model.OK
	return v
}

func getApplications(w *model.World) ApplicationOverview {
	applicationStatuses := renderApplications(w)

	var applicationTable []ApplicationTable
	var totalCount, goodCount, fairCount, poorCount int64

	for _, appStatus := range applicationStatuses {
		totalRequests, totalLatency, totalErrors := calculateMetrics(w.GetApplication(appStatus.Id).Downstreams)

		applicationTable = append(applicationTable, ApplicationTable{
			ID:                   appStatus.Id.Name,
			Status:               appStatus.Status,
			TransactionPerSecond: totalRequests,
			ResponseTime:         totalLatency,
			Errors:               totalErrors,
		})

		totalCount++
		switch appStatus.Status {
		case model.OK:
			goodCount++
		case model.WARNING:
			fairCount++
		case model.INFO:
			fairCount++
		case model.CRITICAL:
			poorCount++

		}
	}

	sort.Slice(applicationTable, func(i, j int) bool {
		return applicationTable[i].TransactionPerSecond > applicationTable[j].TransactionPerSecond
	})

	if len(applicationTable) > 5 {
		applicationTable = applicationTable[:5]
	}

	return ApplicationOverview{
		ApplicationStatus: ApplicationsStats{
			Total: totalCount,
			Good:  goodCount,
			Fair:  fairCount,
			Poor:  poorCount,
		},
		Applications: applicationTable,
	}
}

func calculateMetrics(connections []*model.Connection) (float32, float32, float32) {
	var totalRequests float32
	var totalLatency float32
	var totalErrors float32

	requests := model.GetConnectionsRequestsSum(connections, nil).Last()
	latency := model.GetConnectionsRequestsLatency(connections, nil).Last()
	errors := model.GetConnectionsErrorsSum(connections, nil).Last()

	if !timeseries.IsNaN(requests) {
		totalRequests += requests
	}
	if !timeseries.IsNaN(latency) {
		totalLatency += latency
	}
	if !timeseries.IsNaN(errors) {
		totalErrors += errors
	}

	return totalRequests, totalLatency, totalErrors
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
	browserApps, mobileApps, err := ch.GetAppCounts(ctx, &from, &to)
	if err != nil {
		return nil, EumBadge{}, err
	}

	badge := EumBadge{
		BrowserApps: browserApps,
		MobileApps:  mobileApps,
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

func renderIncidents(w *model.World) IncidentOverview {
	var totalIncidents, criticalIncidents, warningIncidents, closedIncidents int
	var incidentTable []IncidentTable

	for _, app := range w.Applications {
		if len(app.Incidents) == 0 {
			continue
		}

		switch {
		case app.IsK8s():
		case app.Id.Kind == model.ApplicationKindNomadJobGroup:
		case !app.IsStandalone():
		default:
			continue
		}

		sort.Slice(app.Incidents, func(i, j int) bool {
			return app.Incidents[i].OpenedAt.Before(app.Incidents[j].OpenedAt)
		})

		criticalCount := 0
		warningCount := 0
		closedCount := 0
		for _, incident := range app.Incidents {
			if incident.Resolved() {
				closedCount++
				continue
			}
			switch incident.Severity {
			case model.CRITICAL:
				criticalCount++
			case model.WARNING:
				warningCount++
			}
		}

		totalIncidents += len(app.Incidents)
		criticalIncidents += criticalCount
		warningIncidents += warningCount
		closedIncidents += closedCount

		if len(app.Incidents) > 0 {
			icon := getApplicationIcon(app)

			incidentTable = append(incidentTable, IncidentTable{
				Icon:            icon,
				ApplicationName: app.Id.Name,
				OpenIncidents:   int64(len(app.Incidents) - closedCount),
				LastOccurence:   app.Incidents[len(app.Incidents)-1].OpenedAt.ToStandard(),
			})
		}
	}

	sort.Slice(incidentTable, func(i, j int) bool {
		return incidentTable[i].OpenIncidents > incidentTable[j].OpenIncidents
	})

	if len(incidentTable) > 5 {
		incidentTable = incidentTable[:5]
	}

	return IncidentOverview{
		IncidentStats: IncidentStats{
			TotalIncidents:    totalIncidents,
			CriticalIncidents: criticalIncidents,
			WarningIncidents:  warningIncidents,
			ClosedIncidents:   closedIncidents,
		},
		Incidents: incidentTable,
	}
}

func getApplicationIcon(app *model.Application) string {
	types := maps.Keys(app.ApplicationTypes())
	if len(types) == 0 {
		return ""
	}

	var t model.ApplicationType
	if len(types) == 1 {
		t = types[0]
	} else {
		sort.Slice(types, func(i, j int) bool {
			ti, tj := types[i], types[j]
			tiw, tjw := ti.Weight(), tj.Weight()
			if tiw == tjw {
				return ti < tj
			}
			return tiw < tjw
		})
		t = types[0]
	}

	return t.Icon()
}
func renderApplicationsStatsDonutChart(appStats ApplicationsStats, report *model.AuditReport) *model.EChart {
	data := []model.DataPoint{
		{Value: int(appStats.Good), Name: "Good"},
		{Value: int(appStats.Fair), Name: "Fair"},
		{Value: int(appStats.Poor), Name: "Poor"},
	}

	colors := []string{"#66BB6A", "#F99737", "#E7514E"}

	chart := report.GetOrCreateEChart("Node Applications", nil)
	chart.Title = model.TextTitle{
		Text: "Node Applications",
		TextStyle: &model.TextStyle{
			FontSize:   16,
			FontWeight: "normal",
		},
	}
	chart.Tooltip = model.Tooltip{Trigger: "item"}
	chart.Legend = model.Legend{Bottom: "0"}
	chart.SetDonutChartSeries("Applications", data, colors)

	return chart
}

func renderIncidentStatsDonutChart(incidentStats IncidentStats, report *model.AuditReport) *model.EChart {
	data := []model.DataPoint{
		{Value: int(incidentStats.CriticalIncidents), Name: "Critical"},
		{Value: int(incidentStats.WarningIncidents), Name: "Warning"},
		{Value: int(incidentStats.ClosedIncidents), Name: "Closed"},
	}

	colors := []string{"#EF5350", "#FFA726", "#66BB6A"}

	chart := report.GetOrCreateEChart("Incidents Summary", nil)
	chart.Title = model.TextTitle{
		Text: "Incident Summary",
		TextStyle: &model.TextStyle{
			FontSize:   16,
			FontWeight: "normal",
		},
	}
	chart.Tooltip = model.Tooltip{Trigger: "item"}
	chart.Legend = model.Legend{Bottom: "0"}
	chart.SetDonutChartSeries("Incidents", data, colors)

	return chart
}
