package collector

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/ClickHouse/ch-go"
	chproto "github.com/ClickHouse/ch-go/proto"
	"k8s.io/klog"
)

type ErrLogPayload struct {
	UniqueId       string `json:"uniqueId"`
	Service        string `json:"service"`
	ServiceVersion string `json:"serviceVersion"`
	PagePath       string `json:"pagePath"`
	Category       string `json:"category"`
	Grade          string `json:"grade"`
	ErrorUrl       string `json:"errorUrl"`
	Line           int64  `json:"line"`
	Col            int64  `json:"col"`
	Message        string `json:"message"`
	Stack          string `json:"stack"`
	Timestamp      int64  `json:"timestamp"`
	UserId         string `json:"userId"`
	// Include other fields as necessary
}

type DataPointErr struct {
	UniqueId    string
	Timestamp   time.Time
	ServiceName string
	PagePath    string
	Category    string
	Grade       string
	ErrorUrl    string
	Line        int64
	Col         int64
	Message     string
	Stack       string
	UserId      string
}

// ErrLogBatch handles batching of error logs for insertion into ClickHouse.
type ErrLogBatch struct {
	limit int
	exec  func(query ch.Query) error

	lock sync.Mutex
	done chan struct{}

	UniqueId    *chproto.ColStr
	Timestamp   *chproto.ColDateTime64
	ServiceName *chproto.ColLowCardinality[string]
	PagePath    *chproto.ColStr
	Category    *chproto.ColStr
	Grade       *chproto.ColStr
	ErrorUrl    *chproto.ColStr
	Line        *chproto.ColInt64
	Col         *chproto.ColInt64
	Message     *chproto.ColStr
	Stack       *chproto.ColStr
	UserId      *chproto.ColStr
	RawData     *chproto.ColStr
}

func NewErrLogBatch(limit int, timeout time.Duration, exec func(query ch.Query) error) *ErrLogBatch {
	b := &ErrLogBatch{
		limit:       limit,
		exec:        exec,
		done:        make(chan struct{}),
		UniqueId:    new(chproto.ColStr),
		Timestamp:   new(chproto.ColDateTime64).WithPrecision(chproto.PrecisionNano),
		ServiceName: new(chproto.ColStr).LowCardinality(),
		PagePath:    new(chproto.ColStr),
		Category:    new(chproto.ColStr),
		Grade:       new(chproto.ColStr),
		ErrorUrl:    new(chproto.ColStr),
		Line:        new(chproto.ColInt64),
		Col:         new(chproto.ColInt64),
		Message:     new(chproto.ColStr),
		Stack:       new(chproto.ColStr),
		UserId:      new(chproto.ColStr),
		RawData:     new(chproto.ColStr),
	}
	go func() {
		ticker := time.NewTicker(timeout)
		defer ticker.Stop()
		for {
			select {
			case <-b.done:
				return
			case <-ticker.C:
				b.lock.Lock()
				b.save()
				b.lock.Unlock()
			}
		}
	}()
	return b
}

func (b *ErrLogBatch) Close() {
	b.done <- struct{}{}
	b.lock.Lock()
	defer b.lock.Unlock()
	b.save()
}

func (b *ErrLogBatch) Add(dataPoint DataPointErr, raw string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.UniqueId.Append(dataPoint.UniqueId)
	b.Timestamp.Append(dataPoint.Timestamp)
	b.ServiceName.Append(dataPoint.ServiceName)
	b.PagePath.Append(dataPoint.PagePath)
	b.Category.Append(dataPoint.Category)
	b.Grade.Append(dataPoint.Grade)
	b.ErrorUrl.Append(dataPoint.ErrorUrl)
	b.Line.Append(dataPoint.Line)
	b.Col.Append(dataPoint.Col)
	b.Message.Append(dataPoint.Message)
	b.Stack.Append(dataPoint.Stack)
	b.UserId.Append(dataPoint.UserId)
	b.RawData.Append(raw)
	if b.Timestamp.Rows() >= b.limit {
		b.save()
	}
}

func (b *ErrLogBatch) save() {
	if b.Timestamp.Rows() == 0 {
		return
	}
	input := chproto.Input{
		{Name: "UniqueId", Data: b.UniqueId},
		{Name: "Timestamp", Data: b.Timestamp},
		{Name: "ServiceName", Data: b.ServiceName},
		{Name: "PagePath", Data: b.PagePath},
		{Name: "Category", Data: b.Category},
		{Name: "Grade", Data: b.Grade},
		{Name: "ErrorUrl", Data: b.ErrorUrl},
		{Name: "Line", Data: b.Line},
		{Name: "Col", Data: b.Col},
		{Name: "Message", Data: b.Message},
		{Name: "Stack", Data: b.Stack},
		{Name: "UserId", Data: b.UserId},
		{Name: "RawData", Data: b.RawData},
	}
	err := b.exec(ch.Query{Body: input.Into("err_log_data"), Input: input})
	if err != nil {
		klog.Errorln(err)
	}
	for _, col := range input {
		if resettable, ok := col.Data.(chproto.Resettable); ok {
			resettable.Reset()
		}
	}
}

func (c *Collector) ErrLog(w http.ResponseWriter, r *http.Request) {
	project, err := c.getProject(r.Header.Get(ApiKeyHeader))
	if err != nil {
		klog.Errorln(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	_, err = c.getClickhouseClient(project)
	if err != nil {
		klog.Errorln(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "unsupported content type", http.StatusBadRequest)
		return
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		klog.Errorln(err)
		http.Error(w, "failed to read request", http.StatusInternalServerError)
		return
	}
	var payload ErrLogPayload
	err = json.Unmarshal(data, &payload)
	if err != nil {
		klog.Errorln(err)
		http.Error(w, "failed to decode JSON", http.StatusBadRequest)
		return
	}
	dp := DataPointErr{
		UniqueId:    payload.UniqueId,
		Timestamp:   time.Unix(0, payload.Timestamp*int64(time.Millisecond)),
		ServiceName: payload.Service,
		PagePath:    payload.PagePath,
		Category:    payload.Category,
		Grade:       payload.Grade,
		ErrorUrl:    payload.ErrorUrl,
		Line:        payload.Line,
		Col:         payload.Col,
		Message:     payload.Message,
		Stack:       payload.Stack,
		UserId:      payload.UserId,
	}
	c.getErrLogBatch(project).Add(dp, string(data))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}
