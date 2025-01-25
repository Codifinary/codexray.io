package errlogs

import (
	"codexray/clickhouse"
	"codexray/model"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"time"

	"k8s.io/klog"
)

const defaultLimit = 100

type View struct {
	Status  model.Status `json:"status"`
	Message string       `json:"message"`
	Errors  []ErrorLogs  `json:"errors"`
	Limit   int          `json:"limit"`
}

type Query struct {
	Limit int `json:"limit"`
}
type ErrorLogs struct {
	ErrorName    string    `json:"error_name"`
	EventCount   uint64    `json:"event_count"`
	UserImpacted uint64    `json:"user_impacted"`
	LastReported time.Time `json:"last_reported"`
	Category     string    `json:"category"`
}

func RenderErrors(w *model.World, ctx context.Context, ch *clickhouse.Client, query url.Values, serviceName string) *View {
	v := &View{}

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

	from := w.Ctx.From.ToStandard()
	to := w.Ctx.To.ToStandard()

	// Fetch performance data
	rows, err := ch.GetErrorLogs(ctx, &from, &to, serviceName)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Clickhouse error: %s", err)
		return v
	}

	var errors []ErrorLogs
	for _, row := range rows {
		errors = append(errors, ErrorLogs{
			ErrorName:    row.ErrorName,
			EventCount:   row.EventCount,
			UserImpacted: row.UserImpacted,
			LastReported: row.LastReported,
			Category:     row.Category,
		})
	}

	// Sort by error count
	sort.Slice(errors, func(i, j int) bool {
		return errors[i].EventCount > errors[j].EventCount
	})

	v.Status = model.OK
	v.Errors = errors
	v.Limit = q.Limit

	return v
}
