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

type MobilePerfPayload struct {
	ProjectId           string `json:"projectId"`
	Platform            string `json:"platform"`
	RequestPayloadSize  int    `json:"requestPayloadSize"`
	EndpointName        string `json:"endpointName"`
	RequestTime         int64  `json:"requestTime"`
	Service             string `json:"serviceName"`
	Status              bool   `json:"status"`
	ResponseTime        int    `json:"responseTime"`
	ResponsePayloadSize int    `json:"responsePayloadSize"`
	UserID              string `json:"userId"`
	SessionId           string `json:"sessionId"`
	Host                string `json:"host"`
	Device              string `json:"device"`
	StatusCode          int    `json:"statusCode"`
	ServiceVersion      string `json:"serviceVersion"`
	Country             string `json:"country"`
	OS                  string `json:"os"`
}

type MobilePerfDataPoint struct {
	ProjectId           string
	TimestampUnixNano   uint64
	Platform            string
	RequestPayloadSize  int
	EndpointName        string
	RequestTime         int64
	Service             string
	Status              bool
	ResponseTime        int
	ResponsePayloadSize int
	UserID              string
	SessionId           string
	Host                string
	Device              string
	StatusCode          int
	ServiceVersion      string
	Country             string
	OS                  string
	AppType             string
}

type MobilePerfRequestType struct {
	DataPoints []MobilePerfDataPoint `json:"dataPoints"`
}

type MobilePerfBatch struct {
	limit int
	exec  func(query ch.Query) error

	lock sync.Mutex
	done chan struct{}

	Timestamp           *chproto.ColDateTime64
	ProjectId           *chproto.ColStr
	Platform            *chproto.ColStr
	RequestPayloadSize  *chproto.ColInt64
	EndpointName        *chproto.ColStr
	RequestTime         *chproto.ColInt64
	Service             *chproto.ColLowCardinality[string]
	Status              *chproto.ColBool
	ResponseTime        *chproto.ColInt64
	ResponsePayloadSize *chproto.ColInt64
	UserID              *chproto.ColStr
	SessionId           *chproto.ColStr
	Host                *chproto.ColStr
	Device              *chproto.ColStr
	StatusCode          *chproto.ColInt64
	ServiceVersion      *chproto.ColStr
	Country             *chproto.ColStr
	OS                  *chproto.ColStr
	AppType             *chproto.ColStr
	RawData             *chproto.ColStr
}

func NewMobilePerfBatch(limit int, timeout time.Duration, exec func(query ch.Query) error) *MobilePerfBatch {
	b := &MobilePerfBatch{
		limit: limit,
		exec:  exec,
		done:  make(chan struct{}),

		Timestamp:           new(chproto.ColDateTime64).WithPrecision(chproto.PrecisionNano),
		ProjectId:           new(chproto.ColStr),
		Platform:            new(chproto.ColStr),
		RequestPayloadSize:  new(chproto.ColInt64),
		EndpointName:        new(chproto.ColStr),
		RequestTime:         new(chproto.ColInt64),
		Service:             new(chproto.ColStr).LowCardinality(),
		Status:              new(chproto.ColBool),
		ResponseTime:        new(chproto.ColInt64),
		ResponsePayloadSize: new(chproto.ColInt64),
		UserID:              new(chproto.ColStr),
		SessionId:           new(chproto.ColStr),
		Host:                new(chproto.ColStr),
		Device:              new(chproto.ColStr),
		StatusCode:          new(chproto.ColInt64),
		ServiceVersion:      new(chproto.ColStr),
		Country:             new(chproto.ColStr),
		OS:                  new(chproto.ColStr),
		AppType:             new(chproto.ColStr),
		RawData:             new(chproto.ColStr),
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

func (b *MobilePerfBatch) Close() {
	b.done <- struct{}{}
	b.lock.Lock()
	defer b.lock.Unlock()
	b.save()
}

func (b *MobilePerfBatch) Add(perfData *MobilePerfRequestType, raw string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	for _, dataPoint := range perfData.DataPoints {
		b.Timestamp.Append(time.Unix(0, int64(dataPoint.TimestampUnixNano)))
		b.ProjectId.Append(dataPoint.ProjectId)
		b.Platform.Append(dataPoint.Platform)
		b.RequestPayloadSize.Append(int64(dataPoint.RequestPayloadSize))
		b.EndpointName.Append(dataPoint.EndpointName)
		b.RequestTime.Append(int64(dataPoint.RequestTime))
		b.Service.Append(dataPoint.Service)
		b.Status.Append(dataPoint.Status)
		b.ResponseTime.Append(int64(dataPoint.ResponseTime))
		b.ResponsePayloadSize.Append(int64(dataPoint.ResponsePayloadSize))
		b.UserID.Append(dataPoint.UserID)
		b.SessionId.Append(dataPoint.SessionId)
		b.Host.Append(dataPoint.Host)
		b.Device.Append(dataPoint.Device)
		b.StatusCode.Append(int64(dataPoint.StatusCode))
		b.ServiceVersion.Append(dataPoint.ServiceVersion)
		b.Country.Append(dataPoint.Country)
		b.OS.Append(dataPoint.OS)
		b.AppType.Append(dataPoint.AppType)
		b.RawData.Append(raw)
	}

	if b.Timestamp.Rows() >= b.limit {
		b.save()
	}
}

func (b *MobilePerfBatch) save() {
	if b.Timestamp.Rows() == 0 {
		return
	}

	input := chproto.Input{
		{Name: "Timestamp", Data: b.Timestamp},
		{Name: "ProjectId", Data: b.ProjectId},
		{Name: "Platform", Data: b.Platform},
		{Name: "RequestPayloadSize", Data: b.RequestPayloadSize},
		{Name: "EndpointName", Data: b.EndpointName},
		{Name: "RequestTime", Data: b.RequestTime},
		{Name: "Service", Data: b.Service},
		{Name: "Status", Data: b.Status},
		{Name: "ResponseTime", Data: b.ResponseTime},
		{Name: "ResponsePayloadSize", Data: b.ResponsePayloadSize},
		{Name: "UserID", Data: b.UserID},
		{Name: "SessionId", Data: b.SessionId},
		{Name: "Host", Data: b.Host},
		{Name: "Device", Data: b.Device},
		{Name: "StatusCode", Data: b.StatusCode},
		{Name: "ServiceVersion", Data: b.ServiceVersion},
		{Name: "Country", Data: b.Country},
		{Name: "OS", Data: b.OS},
		{Name: "AppType", Data: b.AppType},
		{Name: "RawData", Data: b.RawData},
	}

	err := b.exec(ch.Query{Body: input.Into("mobile_perf_data"), Input: input})
	if err != nil {
		klog.Errorln(err)
	}

	for _, col := range input {
		if resettable, ok := col.Data.(chproto.Resettable); ok {
			resettable.Reset()
		}
	}
}

func (c *Collector) MobilePerf(w http.ResponseWriter, r *http.Request) {
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

	var payload MobilePerfPayload
	err = json.Unmarshal(data, &payload)
	if err != nil {
		klog.Errorln(err)
		http.Error(w, "failed to decode JSON", http.StatusBadRequest)
		return
	}

	dp := MobilePerfDataPoint{
		TimestampUnixNano:   uint64(time.Now().UnixNano()),
		ProjectId:           payload.ProjectId,
		Platform:            payload.Platform,
		RequestPayloadSize:  payload.RequestPayloadSize,
		EndpointName:        payload.EndpointName,
		RequestTime:         payload.RequestTime,
		Service:             payload.Service,
		Status:              payload.Status,
		ResponseTime:        payload.ResponseTime,
		ResponsePayloadSize: payload.ResponsePayloadSize,
		UserID:              payload.UserID,
		SessionId:           payload.SessionId,
		Host:                payload.Host,
		Device:              payload.Device,
		StatusCode:          payload.StatusCode,
		ServiceVersion:      payload.ServiceVersion,
		Country:             payload.Country,
		OS:                  payload.OS,
		AppType:             "mobile",
	}

	perfReq := &MobilePerfRequestType{DataPoints: []MobilePerfDataPoint{dp}}
	c.getMobilePerfBatch(project).Add(perfReq, string(data))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}
