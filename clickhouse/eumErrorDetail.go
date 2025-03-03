package clickhouse

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	// Remove the quotes around the JSON string
	s := strings.Trim(string(b), `"`)
	// Parse the time as milliseconds since the Unix epoch
	ms, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	ct.Time = time.Unix(0, ms*int64(time.Millisecond))
	return nil
}

type ErrorDetail struct {
	Message     string       `json:"errorName"`
	Detail      string       `json:"message"`
	URL         string       `json:"errorurl"`
	Category    string       `json:"category"`
	App         string       `json:"service"`
	AppVersion  string       `json:"serviceVersion"`
	Timestamp   CustomTime   `json:"timestamp"`
	Level       string       `json:"grade"`
	Stack       string       `json:"stack"`
	Breadcrumbs []Breadcrumb `json:"breadcrumbs"`
}

type Breadcrumb struct {
	Type        string     `json:"type"`
	Category    string     `json:"category"`
	Level       string     `json:"level"`
	Description string     `json:"message"`
	Timestamp   CustomTime `json:"timestamp"`
	Data        []KeyValue `json:"data"`
}

type KeyValue struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func (c *Client) GetErrorDetail(ctx context.Context, uniqueId string) (ErrorDetail, error) {
	query := `
    SELECT
        e.Timestamp,
        e.RawData
    FROM 
        err_log_data e
    `

	var filters []string
	var args []any
	if uniqueId != "" {
		filters = append(filters, "e.UniqueId = @uniqueId")
		args = append(args, clickhouse.Named("uniqueId", uniqueId))
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	data, err := c.Query(ctx, query, args...)
	if err != nil {
		return ErrorDetail{}, err
	}
	defer data.Close()

	var timestamp time.Time
	var rawData string
	if data.Next() {
		if err := data.Scan(&timestamp, &rawData); err != nil {
			return ErrorDetail{}, err
		}
	}

	var errorDetail ErrorDetail
	if err := json.Unmarshal([]byte(rawData), &errorDetail); err != nil {
		return ErrorDetail{}, err
	}

	// Set the Timestamp field
	errorDetail.Timestamp = CustomTime{Time: timestamp}

	return errorDetail, nil
}

func (c *Client) GetBreadcrumb(ctx context.Context, uniqueId, breadcrumbType string) ([]Breadcrumb, error) {
	query := `
    SELECT
        e.Timestamp,
        e.RawData
    FROM 
        err_log_data e
    `

	var filters []string
	var args []any
	if uniqueId != "" {
		filters = append(filters, "e.UniqueId = @uniqueId")
		args = append(args, clickhouse.Named("uniqueId", uniqueId))
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	data, err := c.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	var timestamp time.Time
	var rawData string
	if data.Next() {
		if err := data.Scan(&timestamp, &rawData); err != nil {
			return nil, err
		}
	}

	var errorDetail ErrorDetail
	if err := json.Unmarshal([]byte(rawData), &errorDetail); err != nil {
		return nil, err
	}

	var filteredBreadcrumbs []Breadcrumb
	if breadcrumbType == "all" {
		filteredBreadcrumbs = errorDetail.Breadcrumbs
	} else {
		filteredBreadcrumbs = filterBreadcrumbsByType(errorDetail.Breadcrumbs, breadcrumbType)
	}

	// Handle navigation and HTTP breadcrumbs
	for i, breadcrumb := range filteredBreadcrumbs {
		if (breadcrumb.Type == "Navigation" || breadcrumb.Type == "HTTP") && breadcrumb.Description == "" {
			var dataDescription []string
			for _, kv := range breadcrumb.Data {
				// Convert numeric values to strings properly
				switch v := kv.Value.(type) {
				case float64:
					dataDescription = append(dataDescription, fmt.Sprintf("%s: %d", kv.Key, int(v)))
				default:
					dataDescription = append(dataDescription, fmt.Sprintf("%s: %v", kv.Key, kv.Value))
				}
			}
			filteredBreadcrumbs[i].Description = strings.Join(dataDescription, ", ")
		}
	}

	return filteredBreadcrumbs, nil
}

func filterBreadcrumbsByType(breadcrumbs []Breadcrumb, breadcrumbType string) []Breadcrumb {
	var filtered []Breadcrumb
	for _, breadcrumb := range breadcrumbs {
		if breadcrumb.Type == breadcrumbType {
			filtered = append(filtered, breadcrumb)
		}
	}
	return filtered
}
