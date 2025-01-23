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

type PerfPayload struct {
	RedirectTime    int64  `json:"redirectTime"`
	DnsTime         int64  `json:"dnsTime"`
	TtfbTime        int64  `json:"ttfbTime"`
	TcpTime         int64  `json:"tcpTime"`
	TransTime       int64  `json:"transTime"`
	DomAnalysisTime int64  `json:"domAnalysisTime"`
	FptTime         int64  `json:"fptTime"`
	DomReadyTime    int64  `json:"domReadyTime"`
	LoadPageTime    int64  `json:"loadPageTime"`
	ResTime         int64  `json:"resTime"`
	TtlTime         int64  `json:"ttlTime"`
	FirstPackTime   int64  `json:"firstPackTime"`
	FmpTime         int64  `json:"fmpTime"`
	PagePath        string `json:"pagePath"`
	Domain          string `json:"domain"`
	ServiceVersion  string `json:"serviceVersion"`
	Service         string `json:"service"`
	Os              string `json:"os"`
	Device          string `json:"device"`
	Browser         string `json:"browser"`
	CountryCode     string `json:"countryCode"`
	SyntheticUser   bool   `json:"syntheticUser"`
}

type DataPoint struct {
	TimestampUnixNano uint64
	ServiceName       string
	PageName          string
	DeviceId          string
	UserId            string
	TransTime         int64
	LoadPageTime      int64
	ResTime           int64
	Browser           string
}

type PerfRequestType struct {
	DataPoints []DataPoint `json:"dataPoints"`
}

type PerfBatch struct {
	limit int
	exec  func(query ch.Query) error

	lock sync.Mutex
	done chan struct{}

	Timestamp    *chproto.ColDateTime64
	ServiceName  *chproto.ColLowCardinality[string]
	PageName     *chproto.ColLowCardinality[string]
	DeviceId     *chproto.ColStr
	UserId       *chproto.ColStr
	TransTime    *chproto.ColInt64
	LoadPageTime *chproto.ColInt64
	ResTime      *chproto.ColInt64
	Browser      *chproto.ColStr
	RawData      *chproto.ColStr
}

func NewPerfBatch(limit int, timeout time.Duration, exec func(query ch.Query) error) *PerfBatch {
	b := &PerfBatch{
		limit:        limit,
		exec:         exec,
		done:         make(chan struct{}),
		Timestamp:    new(chproto.ColDateTime64).WithPrecision(chproto.PrecisionNano),
		ServiceName:  new(chproto.ColStr).LowCardinality(),
		PageName:     new(chproto.ColStr).LowCardinality(),
		DeviceId:     new(chproto.ColStr),
		UserId:       new(chproto.ColStr),
		TransTime:    new(chproto.ColInt64),
		LoadPageTime: new(chproto.ColInt64),
		ResTime:      new(chproto.ColInt64),
		Browser:      new(chproto.ColStr),
		RawData:      new(chproto.ColStr),
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

func (b *PerfBatch) Close() {
	b.done <- struct{}{}
	b.lock.Lock()
	defer b.lock.Unlock()
	b.save()
}

func (b *PerfBatch) Add(perfData *PerfRequestType, raw string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	for _, dataPoint := range perfData.DataPoints {
		b.Timestamp.Append(time.Unix(0, int64(dataPoint.TimestampUnixNano)))
		b.ServiceName.Append(dataPoint.ServiceName)
		b.PageName.Append(dataPoint.PageName)
		b.DeviceId.Append(dataPoint.DeviceId)
		b.UserId.Append(dataPoint.UserId)
		b.TransTime.Append(dataPoint.TransTime)
		b.LoadPageTime.Append(dataPoint.LoadPageTime)
		b.ResTime.Append(dataPoint.ResTime)
		b.Browser.Append(dataPoint.Browser)
		b.RawData.Append(raw)
	}
	if b.Timestamp.Rows() >= b.limit {
		b.save()
	}
}

func (b *PerfBatch) save() {
	if b.Timestamp.Rows() == 0 {
		return
	}
	input := chproto.Input{
		{Name: "Timestamp", Data: b.Timestamp},
		{Name: "ServiceName", Data: b.ServiceName},
		{Name: "PageName", Data: b.PageName},
		{Name: "DeviceId", Data: b.DeviceId},
		{Name: "UserId", Data: b.UserId},
		{Name: "TransTime", Data: b.TransTime},
		{Name: "LoadPageTime", Data: b.LoadPageTime},
		{Name: "ResTime", Data: b.ResTime},
		{Name: "Browser", Data: b.Browser},
		{Name: "RawData", Data: b.RawData},
	}
	err := b.exec(ch.Query{Body: input.Into("perf_data"), Input: input})
	if err != nil {
		klog.Errorln(err)
	}
	for _, col := range input {
		if resettable, ok := col.Data.(chproto.Resettable); ok {
			resettable.Reset()
		}
	}
}

func (c *Collector) Perf(w http.ResponseWriter, r *http.Request) {
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
	var payload PerfPayload
	err = json.Unmarshal(data, &payload)
	if err != nil {
		klog.Errorln(err)
		http.Error(w, "failed to decode JSON", http.StatusBadRequest)
		return
	}
	// Create DataPoint from payload
	dp := DataPoint{
		TimestampUnixNano: uint64(time.Now().UnixNano()),
		ServiceName:       payload.Service,
		PageName:          payload.PagePath,
		DeviceId:          payload.Device,
		UserId:            "", // Not provided in payload
		TransTime:         payload.TransTime,
		LoadPageTime:      payload.LoadPageTime,
		ResTime:           payload.ResTime,
		Browser:           payload.Browser,
	}
	perfReq := &PerfRequestType{DataPoints: []DataPoint{dp}}
	c.getPerfBatch(project).Add(perfReq, string(data))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}
