package errlogs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"codexray/clickhouse"
	"codexray/model"

	"k8s.io/klog"
)

type ErrorDetailsView struct {
	Status  model.Status `json:"status"`
	Message string       `json:"message"`
	Detail  ErrorDetail  `json:"detail"`
}

type BreadcrumbView struct {
	Status      model.Status `json:"status"`
	Message     string       `json:"message"`
	Breadcrumbs []Breadcrumb `json:"breadcrumbs"`
}

type Breadcrumb struct {
	Type        string    `json:"type"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Level       string    `json:"level"`
	Timestamp   time.Time `json:"timestamp"`
}

type ErrorDetail struct {
	Message    string    `json:"message"`
	Detail     string    `json:"detail"`
	URL        string    `json:"url"`
	Category   string    `json:"category"`
	App        string    `json:"app"`
	AppVersion string    `json:"app_version"`
	Timestamp  time.Time `json:"timestamp"`
	Level      string    `json:"level"`
	Stack      string    `json:"stack"`
}

func ErrorDetails(w *model.World, ctx context.Context, ch *clickhouse.Client, query url.Values, eventID string) *ErrorDetailsView {
	v := &ErrorDetailsView{}

	var q Query
	if s := query.Get("query"); s != "" {
		if err := json.Unmarshal([]byte(s), &q); err != nil {
			klog.Warningln(err)
		}
	}

	if q.Limit <= 0 {
		q.Limit = defaultLimit
	}

	// Check Clickhouse client
	if ch == nil {
		v.Status = model.UNKNOWN
		v.Message = "Clickhouse integration is not configured"
		return v
	}

	// Get error detail from Clickhouse
	result, err := ch.GetErrorDetail(ctx, eventID)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Clickhouse error: %s", err)
		return v
	}

	// Populate the ErrorDetailsView with the result
	v.Status = model.OK
	v.Detail = ErrorDetail{
		Message:    result.Message,
		Detail:     result.Detail,
		URL:        result.URL,
		Category:   result.Category,
		App:        result.App,
		AppVersion: result.AppVersion,
		Timestamp:  result.Timestamp.Time,
		Level:      result.Level,
		Stack:      result.Stack,
	}

	return v
}

func ErrorDetailBreadcrumb(w *model.World, ctx context.Context, ch *clickhouse.Client, query url.Values, eventID, breadcrumbType string) *BreadcrumbView {
	v := &BreadcrumbView{}

	var q Query
	if s := query.Get("query"); s != "" {
		if err := json.Unmarshal([]byte(s), &q); err != nil {
			klog.Warningln(err)
		}
	}

	if q.Limit <= 0 {
		q.Limit = defaultLimit
	}

	// Check Clickhouse client
	if ch == nil {
		v.Status = model.UNKNOWN
		v.Message = "Clickhouse integration is not configured"
		return v
	}

	// Get breadcrumbs from Clickhouse
	rows, err := ch.GetBreadcrumb(ctx, eventID, breadcrumbType)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Clickhouse error: %s", err)
		return v
	}

	var breadcrumbs []Breadcrumb
	for _, row := range rows {
		breadcrumbs = append(breadcrumbs, Breadcrumb{
			Type:        row.Type,
			Category:    row.Category,
			Level:       row.Level,
			Description: row.Description,
			Timestamp:   row.Timestamp.Time,
		})
	}

	v.Status = model.OK
	v.Breadcrumbs = breadcrumbs

	return v
}
