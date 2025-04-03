package overview

import (
	"codexray/clickhouse"
	"codexray/model"
	"context"
)

type MrumOverviewView struct {
	Status  model.Status       `json:"status"`
	Message string             `json:"message"`
	Summary []MrumOverviewData `json:"summary"`
}

type MrumOverviewData struct {
	ServiceName   string `json:"serviceName"`
	TotalRequests uint64 `json:"totalRequests"`
	TotalErrors   uint64 `json:"totalErrors"`
	TotalUsers    uint64 `json:"totalUsers"`
	TotalSessions uint64 `json:"totalSessions"`
}

func RenderMrumOverview(ctx context.Context, ch *clickhouse.Client, w *model.World) *MrumOverviewView {
	v := &MrumOverviewView{}

	if ch == nil {
		v.Status = model.WARNING
		v.Message = "Clickhouse not available"
		return v
	}

	services, err := ch.GetMobileOverview(ctx, w.Ctx.From, w.Ctx.To)
	if err != nil {
		v.Status = model.WARNING
		v.Message = "Failed to get mobile overview: " + err.Error()
		return v
	}

	if len(services) == 0 {
		v.Status = model.OK
		v.Summary = []MrumOverviewData{}
		return v
	}

	var results []MrumOverviewData
	for _, service := range services {
		results = append(results, MrumOverviewData{
			ServiceName:   service.ServiceName,
			TotalRequests: service.TotalRequests,
			TotalErrors:   service.TotalErrors,
			TotalUsers:    service.TotalUsers,
			TotalSessions: service.TotalSessions,
		})
	}

	v.Summary = results
	v.Status = model.OK
	return v
}
