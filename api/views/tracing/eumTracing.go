package tracing

import (
	"codexray/clickhouse"
	"codexray/model"
	"codexray/timeseries"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"k8s.io/klog"
)

type TraceView struct {
	Status  model.Status   `json:"status"`
	Message string         `json:"message"`
	Heatmap *model.Heatmap `json:"heatmap"`
	Traces  []Trace        `json:"traces"`
	Limit   int            `json:"limit"`
}

type Query struct {
	Limit int `json:"limit"`
}

type Trace struct {
	Service    string                `json:"service"`
	TraceId    string                `json:"trace_id"`
	Id         string                `json:"id"`
	ParentId   string                `json:"parent_id"`
	Name       string                `json:"name"`
	Timestamp  int64                 `json:"timestamp"`
	Duration   float64               `json:"duration"`
	Status     model.TraceSpanStatus `json:"status"`
	Attributes map[string]string     `json:"attributes"`
	Events     []Event               `json:"events"`
}

func EumTraces(w *model.World, ch *clickhouse.Client, ctx context.Context, query url.Values, serviceName string) *TraceView {
	v := &TraceView{}

	var q Query
	if s := query.Get("query"); s != "" {
		if err := json.Unmarshal([]byte(s), &q); err != nil {
			klog.Warningln(err)
		}
	}

	if q.Limit <= 0 {
		q.Limit = 100 // default limit
	}

	// Check Clickhouse client
	if ch == nil {
		v.Status = model.UNKNOWN
		v.Message = "Clickhouse integration is not configured"
		return v
	}

	from := w.Ctx.From.ToStandard()
	to := w.Ctx.To.ToStandard()

	heatmap, err := fetchHistogramData(ctx, ch, serviceName, from, to)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Error fetching histogram data: %s", err)
		return v
	}

	traces, err := fetchOtelTraces(ctx, ch, serviceName, from, to, q.Limit)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Error fetching traces: %s", err)
		return v
	}

	v.Status = model.OK
	v.Message = "Data fetched successfully"
	v.Heatmap = heatmap
	v.Traces = traces
	v.Limit = q.Limit

	return v
}

func fetchHistogramData(ctx context.Context, ch *clickhouse.Client, serviceName string, from, to time.Time) (*model.Heatmap, error) {
	q := clickhouse.SpanQuery{
		Ctx: timeseries.Context{
			From: timeseries.Time(from.Unix()),
			To:   timeseries.Time(to.Unix()),
			Step: 60, // Set a positive step value (e.g., 60 seconds)
		},
		Filters: []clickhouse.SpanFilter{
			clickhouse.NewSpanFilter("ServiceName", "=", serviceName),
		},
	}

	histogram, err := ch.GetRootSpansHistogram(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(histogram) < 2 {
		return nil, fmt.Errorf("insufficient histogram data")
	}

	tsCtx := timeseries.Context{
		From: timeseries.Time(from.Unix()),
		To:   timeseries.Time(to.Unix()),
		Step: 60, // Set a positive step value (e.g., 60 seconds)
	}
	heatmap := model.NewHeatmap(tsCtx, "Latency & Errors heatmap, requests per second")
	for _, h := range model.HistogramSeries(histogram[1:], 0, 0) {
		heatmap.AddSeries(h.Name, h.Title, h.Data, h.Threshold, h.Value)
	}
	heatmap.AddSeries("errors", "errors", histogram[0].TimeSeries, "", "err")

	return heatmap, nil
}

func fetchOtelTraces(ctx context.Context, ch *clickhouse.Client, serviceName string, from, to time.Time, limit int) ([]Trace, error) {
	traces, err := ch.GetTracesByServiceName(ctx, serviceName, from, to, limit)
	if err != nil {
		return nil, err
	}

	var result []Trace
	for _, t := range traces {
		for _, s := range t.Spans {
			trace := Trace{
				Service:    s.ServiceName,
				TraceId:    s.TraceId,
				Id:         s.SpanId,
				ParentId:   s.ParentSpanId,
				Name:       s.Name,
				Timestamp:  s.Timestamp.UnixMilli(),
				Duration:   s.Duration.Seconds() * 1000,
				Status:     s.Status(),
				Attributes: map[string]string{},
				Events:     []Event{},
			}
			for name, value := range s.ResourceAttributes {
				trace.Attributes[name] = value
			}
			for name, value := range s.SpanAttributes {
				trace.Attributes[name] = value
			}
			for _, e := range s.Events {
				trace.Events = append(trace.Events, Event{
					Timestamp:  e.Timestamp.UnixMilli(),
					Name:       e.Name,
					Attributes: e.Attributes,
				})
			}
			result = append(result, trace)
		}
	}

	return result, nil
}
