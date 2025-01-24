package clickhouse

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ErrorRow struct {
	ErrorName    string
	EventCount   uint64
	UserImpacted uint64
	LastReported time.Time
	Category     string
}

func (c *Client) GetErrorLogs(ctx context.Context, from, to *time.Time, serviceName string) ([]ErrorRow, error) {
	fmt.Println("from:", from)
	fmt.Println("to:", to)
	fmt.Println("serviceName:", serviceName)

	query := `
SELECT 
    e.ErrorName,
    count(e.ErrorName) AS eventCount,
    countDistinct(e.UserId) AS userImpacted,
    max(e.Timestamp) AS lastReported,
    e.Category
FROM
    err_log_data e`

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

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	query += `
GROUP BY
    e.ErrorName, e.Category
	ORDER BY
		eventCount DESC
	`

	fmt.Println("query:", query)
	fmt.Println("args:", args)

	rows, err := c.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ErrorRow
	for rows.Next() {
		var row ErrorRow
		if err := rows.Scan(&row.ErrorName, &row.EventCount, &row.UserImpacted, &row.LastReported, &row.Category); err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	fmt.Println("result:", result)

	return result, nil
}
