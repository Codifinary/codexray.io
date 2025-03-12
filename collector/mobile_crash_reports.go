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

type MobileCrashReportPayload struct {
	UniqueId        string `json:"uniqueId"`
	SessionId       string `json:"sessionId"`
	CrashTime       int64  `json:"crashTime"`
	CrashReason     string `json:"crashReason"`
	FileName        string `json:"fileName"`
	LineNo          string `json:"lineNo"`
	CrashStackTrace string `json:"crashStackTrace"`
	MemoryUsage     int64  `json:"memoryUsage"`
	Os              string `json:"os"`
	Platform        string `json:"platform"`
	ServiceVersion  string `json:"serviceVersion"`
	DeviceInfo      string `json:"deviceInfo"`
	Service         string `json:"service"`
	Country         string `json:"country"`
}

type MobileCrashReportDataPoint struct {
	TimestampUnixNano uint64
	UniqueId          string
	SessionId         string
	CrashTime         int64
	CrashReason       string
	FileName          string
	LineNo            string
	CrashStackTrace   string
	MemoryUsage       int64
	Os                string
	Platform          string
	ServiceVersion    string
	DeviceInfo        string
	Service           string
	Country           string
}

type MobileCrashReportRequestType struct {
	DataPoints []MobileCrashReportDataPoint `json:"dataPoints"`
}

type MobileCrashReportBatch struct {
	limit int
	exec  func(query ch.Query) error

	lock sync.Mutex
	done chan struct{}

	Timestamp       *chproto.ColDateTime64
	UniqueId        *chproto.ColStr
	SessionId       *chproto.ColStr
	CrashTime       *chproto.ColInt64
	CrashReason     *chproto.ColStr
	FileName        *chproto.ColStr
	LineNo          *chproto.ColStr
	CrashStackTrace *chproto.ColStr
	MemoryUsage     *chproto.ColInt64
	Os              *chproto.ColStr
	Platform        *chproto.ColStr
	ServiceVersion  *chproto.ColStr
	DeviceInfo      *chproto.ColStr
	Service         *chproto.ColStr
	Country         *chproto.ColStr
	RawData         *chproto.ColStr
}

func NewMobileCrashReportBatch(limit int, timeout time.Duration, exec func(query ch.Query) error) *MobileCrashReportBatch {
	b := &MobileCrashReportBatch{
		limit: limit,
		exec:  exec,
		done:  make(chan struct{}),

		Timestamp:       new(chproto.ColDateTime64).WithPrecision(chproto.PrecisionNano),
		UniqueId:        new(chproto.ColStr),
		SessionId:       new(chproto.ColStr),
		CrashTime:       new(chproto.ColInt64),
		CrashReason:     new(chproto.ColStr),
		FileName:        new(chproto.ColStr),
		LineNo:          new(chproto.ColStr),
		CrashStackTrace: new(chproto.ColStr),
		MemoryUsage:     new(chproto.ColInt64),
		Os:              new(chproto.ColStr),
		Platform:        new(chproto.ColStr),
		ServiceVersion:  new(chproto.ColStr),
		DeviceInfo:      new(chproto.ColStr),
		Service:         new(chproto.ColStr),
		Country:         new(chproto.ColStr),
		RawData:         new(chproto.ColStr),
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

func (b *MobileCrashReportBatch) Close() {
	b.done <- struct{}{}
	b.lock.Lock()
	defer b.lock.Unlock()
	b.save()
}

func (b *MobileCrashReportBatch) Add(crashData *MobileCrashReportRequestType, raw string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	for _, dataPoint := range crashData.DataPoints {
		b.Timestamp.Append(time.Unix(0, int64(dataPoint.TimestampUnixNano)))
		b.UniqueId.Append(dataPoint.UniqueId)
		b.SessionId.Append(dataPoint.SessionId)
		b.CrashTime.Append(dataPoint.CrashTime)
		b.CrashReason.Append(dataPoint.CrashReason)
		b.FileName.Append(dataPoint.FileName)
		b.LineNo.Append(dataPoint.LineNo)
		b.CrashStackTrace.Append(dataPoint.CrashStackTrace)
		b.MemoryUsage.Append(dataPoint.MemoryUsage)
		b.Os.Append(dataPoint.Os)
		b.Platform.Append(dataPoint.Platform)
		b.ServiceVersion.Append(dataPoint.ServiceVersion)
		b.DeviceInfo.Append(dataPoint.DeviceInfo)
		b.Service.Append(dataPoint.Service)
		b.Country.Append(dataPoint.Country)
		b.RawData.Append(raw)
	}

	if b.Timestamp.Rows() >= b.limit {
		b.save()
	}
}

func (b *MobileCrashReportBatch) save() {
	if b.Timestamp.Rows() == 0 {
		return
	}

	input := chproto.Input{
		{Name: "Timestamp", Data: b.Timestamp},
		{Name: "UniqueId", Data: b.UniqueId},
		{Name: "SessionId", Data: b.SessionId},
		{Name: "CrashTime", Data: b.CrashTime},
		{Name: "CrashReason", Data: b.CrashReason},
		{Name: "FileName", Data: b.FileName},
		{Name: "LineNo", Data: b.LineNo},
		{Name: "CrashStackTrace", Data: b.CrashStackTrace},
		{Name: "MemoryUsage", Data: b.MemoryUsage},
		{Name: "Os", Data: b.Os},
		{Name: "Platform", Data: b.Platform},
		{Name: "ServiceVersion", Data: b.ServiceVersion},
		{Name: "DeviceInfo", Data: b.DeviceInfo},
		{Name: "Service", Data: b.Service},
		{Name: "Country", Data: b.Country},
		{Name: "RawData", Data: b.RawData},
	}

	err := b.exec(ch.Query{Body: input.Into("mobile_crash_reports"), Input: input})
	if err != nil {
		klog.Errorln(err)
	}

	for _, col := range input {
		if resettable, ok := col.Data.(chproto.Resettable); ok {
			resettable.Reset()
		}
	}
}

func (c *Collector) MobileCrashReports(w http.ResponseWriter, r *http.Request) {
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

	var payload MobileCrashReportPayload
	err = json.Unmarshal(data, &payload)
	if err != nil {
		klog.Errorln(err)
		http.Error(w, "failed to decode JSON", http.StatusBadRequest)
		return
	}

	dp := MobileCrashReportDataPoint{
		TimestampUnixNano: uint64(time.Now().UnixNano()),
		UniqueId:          payload.UniqueId,
		SessionId:         payload.SessionId,
		CrashTime:         payload.CrashTime,
		CrashReason:       payload.CrashReason,
		FileName:          payload.FileName,
		LineNo:            payload.LineNo,
		CrashStackTrace:   payload.CrashStackTrace,
		MemoryUsage:       payload.MemoryUsage,
		Os:                payload.Os,
		Platform:          payload.Platform,
		ServiceVersion:    payload.ServiceVersion,
		DeviceInfo:        payload.DeviceInfo,
		Service:           payload.Service,
		Country:           payload.Country,
	}

	crashReq := &MobileCrashReportRequestType{DataPoints: []MobileCrashReportDataPoint{dp}}
	c.getMobileCrashReportBatch(project).Add(crashReq, string(data))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}
