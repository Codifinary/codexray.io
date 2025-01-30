package clickhouse

import (
	"context"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

// ErrorRow struct is a row of the error log table.
type ErrorsRow struct {
	EventID      string    `json:"event_id"`
	UserID       string    `json:"user_id"`
	Device       string    `json:"device"`
	OS           string    `json:"os"`
	Browser      string    `json:"browser"`
	LastReported time.Time `json:"last_reported"`
}

func (c *Client) GetErrors(ctx context.Context, from, to *time.Time, serviceName, errorName string) ([]ErrorsRow, error) {

	query := `
	SELECT
		e.UniqueId AS eventID,
		e.UserId,
		e.Device,
		e.OS,
		e.Browser,
		max(e.Timestamp) AS lastReported
	FROM
		err_log_data e
	
	 
	`

	// Conditionally add time range and service name filtering
	var filters []string
	var args []any
	if from != nil {
		filters = append(filters, "e.Timestamp >= @from")
		args = append(args, clickhouse.Named("from", *from))
	}
	if to != nil {
		filters = append(filters, "e.Timestamp <= @to")
		args = append(args, clickhouse.Named("to", *to))
	}
	if serviceName != "" {
		filters = append(filters, "e.ServiceName = @serviceName")
		args = append(args, clickhouse.Named("serviceName", serviceName))
	}
	if errorName != "" {
		filters = append(filters, "e.ErrorName = @errorName")
		args = append(args, clickhouse.Named("errorName", errorName))
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	query += `
	GROUP BY
		 e.UniqueId, e.UserId, e.Device, e.OS, e.Browser
		ORDER BY
			lastReported DESC
		`

	rows, err := c.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ErrorsRow
	for rows.Next() {
		var row ErrorsRow
		if err := rows.Scan(&row.EventID, &row.UserID, &row.Device, &row.OS, &row.Browser, &row.LastReported); err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, nil
}
