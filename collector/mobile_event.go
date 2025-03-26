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

type MobileEventPayload struct {
	UserId         string `json:"userId"`
	ProjectId      string `json:"projectId"`
	Name           string `json:"name"`
	StartTime      int64  `json:"startTime"`
	SessionId      string `json:"sessionId"`
	Os             string `json:"os"`
	Platform       string `json:"platform"`
	ServiceVersion string `json:"serviceVersion"`
	Device         string `json:"device"`
	Service        string `json:"service"`
	Country        string `json:"country"`
}

type MobileEventDataPoint struct {
	TimestampUnixNano uint64
	UserId            string
	ProjectId         string
	Name              string
	StartTime         int64
	SessionId         string
	Os                string
	Platform          string
	ServiceVersion    string
	Device            string
	Service           string
	Country           string
}

type MobileEventRequestType struct {
	DataPoints []MobileEventDataPoint `json:"dataPoints"`
}

type MobileEventBatch struct {
	limit int
	exec  func(query ch.Query) error

	lock sync.Mutex
	done chan struct{}

	Timestamp      *chproto.ColDateTime64
	UserId         *chproto.ColStr
	ProjectId      *chproto.ColStr
	Name           *chproto.ColStr
	StartTime      *chproto.ColInt64
	SessionId      *chproto.ColStr
	Os             *chproto.ColStr
	Platform       *chproto.ColStr
	ServiceVersion *chproto.ColStr
	Device         *chproto.ColStr
	Service        *chproto.ColStr
	Country        *chproto.ColStr
	RawData        *chproto.ColStr
}

func NewMobileEventBatch(limit int, timeout time.Duration, exec func(query ch.Query) error) *MobileEventBatch {
	b := &MobileEventBatch{
		limit: limit,
		exec:  exec,
		done:  make(chan struct{}),

		Timestamp:      new(chproto.ColDateTime64).WithPrecision(chproto.PrecisionNano),
		ProjectId:      new(chproto.ColStr),
		UserId:         new(chproto.ColStr),
		Name:           new(chproto.ColStr),
		StartTime:      new(chproto.ColInt64),
		SessionId:      new(chproto.ColStr),
		Os:             new(chproto.ColStr),
		Platform:       new(chproto.ColStr),
		ServiceVersion: new(chproto.ColStr),
		Device:         new(chproto.ColStr),
		Service:        new(chproto.ColStr),
		Country:        new(chproto.ColStr),
		RawData:        new(chproto.ColStr),
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

func (b *MobileEventBatch) Close() {
	b.done <- struct{}{}
	b.lock.Lock()
	b.save()
	b.lock.Unlock()
}

func (b *MobileEventBatch) Add(eventData *MobileEventRequestType, raw string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	for _, dataPoint := range eventData.DataPoints {
		b.Timestamp.Append(time.Unix(0, int64(dataPoint.TimestampUnixNano)))
		b.ProjectId.Append(dataPoint.ProjectId)
		b.UserId.Append(dataPoint.UserId)
		b.Name.Append(dataPoint.Name)
		b.StartTime.Append(dataPoint.StartTime)
		b.SessionId.Append(dataPoint.SessionId)
		b.Os.Append(dataPoint.Os)
		b.Platform.Append(dataPoint.Platform)
		b.ServiceVersion.Append(dataPoint.ServiceVersion)
		b.Device.Append(dataPoint.Device)
		b.Service.Append(dataPoint.Service)
		b.Country.Append(dataPoint.Country)
		b.RawData.Append(raw)
	}

	if b.Timestamp.Rows() >= b.limit {
		b.save()
	}
}

func (b *MobileEventBatch) save() {
	if b.Timestamp.Rows() == 0 {
		return
	}

	input := chproto.Input{
		{Name: "Timestamp", Data: b.Timestamp},
		{Name: "ProjectId", Data: b.ProjectId},
		{Name: "UserId", Data: b.UserId},
		{Name: "Name", Data: b.Name},
		{Name: "StartTime", Data: b.StartTime},
		{Name: "SessionId", Data: b.SessionId},
		{Name: "Os", Data: b.Os},
		{Name: "Platform", Data: b.Platform},
		{Name: "ServiceVersion", Data: b.ServiceVersion},
		{Name: "Device", Data: b.Device},
		{Name: "Service", Data: b.Service},
		{Name: "Country", Data: b.Country},
		{Name: "RawData", Data: b.RawData},
	}

	err := b.exec(ch.Query{Body: input.Into("mobile_event_data"), Input: input})
	if err != nil {
		klog.Errorln(err)
	}

	for _, col := range input {
		if resettable, ok := col.Data.(chproto.Resettable); ok {
			resettable.Reset()
		}
	}
}

func (c *Collector) MobileEvent(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var payload MobileEventPayload
	err = json.Unmarshal(data, &payload)
	if err != nil {
		klog.Errorln(err)
		http.Error(w, "failed to decode JSON", http.StatusBadRequest)
		return
	}

	dp := MobileEventDataPoint{
		TimestampUnixNano: uint64(time.Now().UnixNano()),
		ProjectId:         payload.ProjectId,
		UserId:            payload.UserId,
		Name:              payload.Name,
		StartTime:         payload.StartTime,
		SessionId:         payload.SessionId,
		Os:                payload.Os,
		Platform:          payload.Platform,
		ServiceVersion:    payload.ServiceVersion,
		Device:            payload.Device,
		Service:           payload.Service,
		Country:           payload.Country,
	}

	eventReq := &MobileEventRequestType{DataPoints: []MobileEventDataPoint{dp}}
	c.getMobileEventBatch(project).Add(eventReq, string(data))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}
