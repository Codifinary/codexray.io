package errlogs

import (
	"codexray/clickhouse"
	"codexray/model"
	"context"
	"encoding/json"
	"fmt"

	"net/url"
	"time"

	"k8s.io/klog"
)

type ErrorsView struct {
	Status  model.Status `json:"status"`
	Message string       `json:"message"`
	Errors  []Error      `json:"errors"`
	Limit   int          `json:"limit"`
}

type Error struct {
	EventID      string    `json:"event_id"`
	UserID       string    `json:"user_id"`
	Device       string    `json:"device"`
	OS           string    `json:"os"`
	Browser      string    `json:"browser"`
	LastReported time.Time `json:"last_reported"`
}

func Errors(w *model.World, ctx context.Context, ch *clickhouse.Client, query url.Values, serviceName string, errorName string) *ErrorsView {
	v := &ErrorsView{}

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

	rows, err := ch.GetErrors(ctx, &from, &to, serviceName, errorName)
	if err != nil {
		klog.Errorln(err)
		v.Status = model.WARNING
		v.Message = fmt.Sprintf("Clickhouse error: %s", err)
		return v
	}

	var errors []Error
	for _, row := range rows {
		errors = append(errors, Error{
			EventID:      row.EventID,
			UserID:       row.UserID,
			Device:       row.Device,
			OS:           row.OS,
			Browser:      row.Browser,
			LastReported: row.LastReported,
		})
	}

	v.Status = model.OK
	v.Errors = errors
	v.Limit = q.Limit

	return v

}
