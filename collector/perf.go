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
	CountryCode     string `json:"countryCode"`
	SyntheticUser   bool   `json:"syntheticUser"`
	SslTime         int64  `json:"sslTime"`
	Browser         string `json:"browser"`
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
	AppType           string
	DnsTime           int64
	TcpTime           int64
	SslTime           int64
	DomAnalysisTime   int64
	DomReadyTime      int64
	FirstPackTime     int64
	FmpTime           int64
	FptTime           int64
	RedirectTime      int64
	TtfbTime          int64
	TtlTime           int64
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

	Timestamp       *chproto.ColDateTime64
	ServiceName     *chproto.ColLowCardinality[string]
	PageName        *chproto.ColLowCardinality[string]
	DeviceId        *chproto.ColStr
	UserId          *chproto.ColStr
	TransTime       *chproto.ColInt64
	LoadPageTime    *chproto.ColInt64
	ResTime         *chproto.ColInt64
	AppType         *chproto.ColStr
	DnsTime         *chproto.ColInt64
	TcpTime         *chproto.ColInt64
	SslTime         *chproto.ColInt64
	DomAnalysisTime *chproto.ColInt64
	DomReadyTime    *chproto.ColInt64
	FirstPackTime   *chproto.ColInt64
	FmpTime         *chproto.ColInt64
	FptTime         *chproto.ColInt64
	RedirectTime    *chproto.ColInt64
	TtfbTime        *chproto.ColInt64
	TtlTime         *chproto.ColInt64
	RawData         *chproto.ColStr
	Browser         *chproto.ColStr
}

func NewPerfBatch(limit int, timeout time.Duration, exec func(query ch.Query) error) *PerfBatch {
	b := &PerfBatch{
		limit:           limit,
		exec:            exec,
		done:            make(chan struct{}),
		Timestamp:       new(chproto.ColDateTime64).WithPrecision(chproto.PrecisionNano),
		ServiceName:     new(chproto.ColStr).LowCardinality(),
		PageName:        new(chproto.ColStr).LowCardinality(),
		DeviceId:        new(chproto.ColStr),
		UserId:          new(chproto.ColStr),
		TransTime:       new(chproto.ColInt64),
		LoadPageTime:    new(chproto.ColInt64),
		ResTime:         new(chproto.ColInt64),
		AppType:         new(chproto.ColStr),
		DnsTime:         new(chproto.ColInt64),
		TcpTime:         new(chproto.ColInt64),
		SslTime:         new(chproto.ColInt64),
		DomAnalysisTime: new(chproto.ColInt64),
		DomReadyTime:    new(chproto.ColInt64),
		FirstPackTime:   new(chproto.ColInt64),
		FmpTime:         new(chproto.ColInt64),
		FptTime:         new(chproto.ColInt64),
		RedirectTime:    new(chproto.ColInt64),
		TtfbTime:        new(chproto.ColInt64),
		TtlTime:         new(chproto.ColInt64),
		RawData:         new(chproto.ColStr),
		Browser:         new(chproto.ColStr),
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
		b.AppType.Append("Browser")
		b.DnsTime.Append(dataPoint.DnsTime)
		b.TcpTime.Append(dataPoint.TcpTime)
		b.SslTime.Append(dataPoint.SslTime)
		b.DomAnalysisTime.Append(dataPoint.DomAnalysisTime)
		b.DomReadyTime.Append(dataPoint.DomReadyTime)
		b.FirstPackTime.Append(dataPoint.FirstPackTime)
		b.FmpTime.Append(dataPoint.FmpTime)
		b.FptTime.Append(dataPoint.FptTime)
		b.RedirectTime.Append(dataPoint.RedirectTime)
		b.TtfbTime.Append(dataPoint.TtfbTime)
		b.TtlTime.Append(dataPoint.TtlTime)
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
		{Name: "AppType", Data: b.AppType},
		{Name: "DnsTime", Data: b.DnsTime},
		{Name: "TcpTime", Data: b.TcpTime},
		{Name: "SslTime", Data: b.SslTime},
		{Name: "DomAnalysisTime", Data: b.DomAnalysisTime},
		{Name: "DomReadyTime", Data: b.DomReadyTime},
		{Name: "FirstPackTime", Data: b.FirstPackTime},
		{Name: "FmpTime", Data: b.FmpTime},
		{Name: "FptTime", Data: b.FptTime},
		{Name: "RedirectTime", Data: b.RedirectTime},
		{Name: "TtfbTime", Data: b.TtfbTime},
		{Name: "TtlTime", Data: b.TtlTime},
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
		DnsTime:           payload.DnsTime,
		TcpTime:           payload.TcpTime,
		SslTime:           payload.SslTime,
		DomAnalysisTime:   payload.DomAnalysisTime,
		DomReadyTime:      payload.DomReadyTime,
		FirstPackTime:     payload.FirstPackTime,
		FmpTime:           payload.FmpTime,
		FptTime:           payload.FptTime,
		RedirectTime:      payload.RedirectTime,
		TtfbTime:          payload.TtfbTime,
		TtlTime:           payload.TtlTime,
		AppType:           "Browser",
		Browser:           payload.Browser,
	}
	perfReq := &PerfRequestType{DataPoints: []DataPoint{dp}}
	c.getPerfBatch(project).Add(perfReq, string(data))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}
