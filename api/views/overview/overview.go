package overview

import (
	"context"

	"codexray/api/views/incident"
	"codexray/clickhouse"
	"codexray/model"
)

type Overview struct {
	Applications []*ApplicationStatus        `json:"applications"`
	Health       []*ApplicationStatus        `json:"health"`
	Incidents    []incident.Summary          `json:"incidents"`
	Map          []*Application              `json:"map"`
	Nodes        *model.Table                `json:"nodes"`
	Deployments  []*Deployment               `json:"deployments"`
	Traces       *Traces                     `json:"traces"`
	TracesView   *TracesView                 `json:"traces_view"`
	Costs        *Costs                      `json:"costs"`
	Categories   []model.ApplicationCategory `json:"categories"`
	EumApps      *EumView                    `json:"eumapps"`
	// Perfs        *PerfView                   `json:"perfs"`
}

func Render(ctx context.Context, ch *clickhouse.Client, w *model.World, view, query string) *Overview {
	v := &Overview{
		Categories: w.Categories,
	}

	switch view {
	case "applications":
		v.Applications = renderApplications(w)
	case "incidents":
		for _, app := range w.Applications {
			switch {
			case app.IsK8s():
			case app.Id.Kind == model.ApplicationKindNomadJobGroup:
			case !app.IsStandalone():
			default:
				continue
			}
			for _, i := range app.Incidents {
				v.Incidents = append(v.Incidents, incident.CalcSummary(w, app, i))
			}
		}
	case "map":
		v.Map = renderServiceMap(w)
	case "nodes":
		v.Nodes = renderNodes(w)
	case "deployments":
		v.Deployments = renderDeployments(w)
	case "traces":
		v.TracesView = renderTraceOverview(ctx, ch, w, query)
	case "costs":
		v.Costs = renderCosts(w)
	case "eumapps":
		v.EumApps = renderEumApps(ctx, ch, w, query)
		// case "perfs":
		// 	v.Perfs = renderPerfs(ctx, ch, w, query)
	}
	return v
}
