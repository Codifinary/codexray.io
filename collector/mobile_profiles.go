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

type MobileUserRegistrationPayload struct {
	UserId           string `json:"userId"`
	OS               string `json:"os"`
	Platform         int32  `json:"platform"`
	ServiceVersion   string `json:"serviceVersion"`
	Device           string `json:"device"`
	Service          string `json:"service"`
	Country          string `json:"country"`
	RegistrationTime int64  `json:"registrationTime"`
	IpAddress        string `json:"ipAddress"`
	TimeBucket       int32  `json:"timeBucket"`
}

type MobileUserRegistrationBatch struct {
	limit int
	exec  func(query ch.Query) error

	lock sync.Mutex
	done chan struct{}

	UserId           *chproto.ColStr
	OS               *chproto.ColLowCardinality[string]
	Platform         *chproto.ColInt32
	ServiceVersion   *chproto.ColLowCardinality[string]
	Device           *chproto.ColStr
	Service          *chproto.ColStr
	Country          *chproto.ColLowCardinality[string]
	RegistrationTime *chproto.ColDateTime64
	IpAddress        *chproto.ColStr
	TimeBucket       *chproto.ColInt32
	RawData          *chproto.ColStr
}

func NewMobileUserRegistrationBatch(limit int, timeout time.Duration, exec func(query ch.Query) error) *MobileUserRegistrationBatch {
	b := &MobileUserRegistrationBatch{
		limit: limit,
		exec:  exec,
		done:  make(chan struct{}),

		UserId:           new(chproto.ColStr),
		OS:               new(chproto.ColStr).LowCardinality(),
		Platform:         new(chproto.ColInt32),
		ServiceVersion:   new(chproto.ColStr).LowCardinality(),
		Device:           new(chproto.ColStr),
		Service:          new(chproto.ColStr),
		Country:          new(chproto.ColStr).LowCardinality(),
		RegistrationTime: new(chproto.ColDateTime64).WithPrecision(chproto.PrecisionNano),
		IpAddress:        new(chproto.ColStr),
		TimeBucket:       new(chproto.ColInt32),
		RawData:          new(chproto.ColStr),
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

func (b *MobileUserRegistrationBatch) Close() {
	b.done <- struct{}{}
	b.lock.Lock()
	defer b.lock.Unlock()
	b.save()
}

func (b *MobileUserRegistrationBatch) Add(registration *MobileUserRegistrationPayload, raw string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	registrationTime := time.UnixMilli(registration.RegistrationTime)
	b.UserId.Append(registration.UserId)
	b.OS.Append(registration.OS)
	b.Platform.Append(registration.Platform)
	b.ServiceVersion.Append(registration.ServiceVersion)
	b.Device.Append(registration.Device)
	b.Service.Append(registration.Service)
	b.Country.Append(registration.Country)
	b.RegistrationTime.Append(registrationTime)
	b.IpAddress.Append(registration.IpAddress)
	b.TimeBucket.Append(registration.TimeBucket)
	b.RawData.Append(raw)

	if b.UserId.Rows() >= b.limit {
		b.save()
	}
}

func (b *MobileUserRegistrationBatch) save() {
	if b.UserId.Rows() == 0 {
		return
	}

	input := chproto.Input{
		{Name: "UserId", Data: b.UserId},
		{Name: "OS", Data: b.OS},
		{Name: "Platform", Data: b.Platform},
		{Name: "ServiceVersion", Data: b.ServiceVersion},
		{Name: "Device", Data: b.Device},
		{Name: "Service", Data: b.Service},
		{Name: "Country", Data: b.Country},
		{Name: "RegistrationTime", Data: b.RegistrationTime},
		{Name: "IpAddress", Data: b.IpAddress},
		{Name: "TimeBucket", Data: b.TimeBucket},
		{Name: "RawData", Data: b.RawData},
	}

	err := b.exec(ch.Query{Body: input.Into("mobile_user_registration"), Input: input})
	if err != nil {
		klog.Errorln(err)
	}

	for _, col := range input {
		if resettable, ok := col.Data.(chproto.Resettable); ok {
			resettable.Reset()
		}
	}
}

func (c *Collector) MobileUserRegistration(w http.ResponseWriter, r *http.Request) {
	project, err := c.getProject(r.Header.Get(ApiKeyHeader))
	if err != nil {
		klog.Errorln(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	_, err = c.getClickhouseClient(project)
	if err != nil {
		klog.Errorln(err)
		http.Error(w, err.Error(),
			http.StatusNotFound)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "unsupported content type: "+contentType, http.StatusBadRequest)
		return
	}

	decoder, err := getDecoder(r.Header.Get("Content-Encoding"), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(decoder)
	if err != nil {
		klog.Errorln(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var registration MobileUserRegistrationPayload
	if err := json.Unmarshal(data, &registration); err != nil {
		klog.Errorln(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	c.getMobileUserRegistrationBatch(project).Add(&registration, string(data))

	w.WriteHeader(http.StatusOK)

}
