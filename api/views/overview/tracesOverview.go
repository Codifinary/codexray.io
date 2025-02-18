package overview

import (
	"codexray/clickhouse"
	"codexray/model"
	"context"
)

type TracesView struct {
	Message string              `json:"message"`
	Error   string              `json:"error"`
	Summary model.TracesSummary `json:"summary"`
	Traces  []model.TraceView   `json:"traces"`
}

func renderTraceOverview(ctx context.Context, ch *clickhouse.Client, w *model.World, serviceName string) *TracesView {
	v := &TracesView{}

	if ch == nil {
		v.Message = "no_clickhouse"
		return v
	}

	summary, err := ch.GetAggregatedTracesSummary(ctx, w)
	if err != nil {
		v.Error = err.Error()
		v.Message = "error"
		return v
	}
	v.Summary = *summary

	traces, err := ch.GetServiceSpecificTraceViews(ctx, w)
	if err != nil {
		v.Error = err.Error()
		v.Message = "error"
		return v
	}
	v.Traces = traces

	return v
}
