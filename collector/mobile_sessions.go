package collector

import (
	"context"
	"fmt"
	"sync"
	"time"

	"codexray/db"

	"github.com/ClickHouse/ch-go"
	chproto "github.com/ClickHouse/ch-go/proto"
	"k8s.io/klog"
)

type MobileSessionDataPoint struct {
	Timestamp uint64
	SessionId string
	UserId    string
	StartTime uint64
	EndTime   uint64
	Country   string
	Device    string
	OS        string
}

type MobileSessionRequestType struct {
	DataPoints []MobileSessionDataPoint `json:"dataPoints"`
}

type MobileSessionBatch struct {
	limit int
	exec  func(query ch.Query) error

	lock sync.Mutex
	done chan struct{}

	Timestamp      *chproto.ColDateTime64
	SessionId      *chproto.ColStr
	UserId         *chproto.ColStr
	StartTime      *chproto.ColDateTime64
	EndTime        *chproto.ColNullable[time.Time]
	Country        *chproto.ColStr
	Device         *chproto.ColStr
	OS             *chproto.ColStr
	existingCombos map[string]bool
}

func NewMobileSessionBatch(limit int, timeout time.Duration, exec func(query ch.Query) error) *MobileSessionBatch {
	endTimeCol := new(chproto.ColDateTime64).WithPrecision(chproto.PrecisionNano)
	b := &MobileSessionBatch{
		limit: limit,
		exec:  exec,
		done:  make(chan struct{}),

		Timestamp:      new(chproto.ColDateTime64).WithPrecision(chproto.PrecisionNano),
		SessionId:      new(chproto.ColStr),
		UserId:         new(chproto.ColStr),
		StartTime:      new(chproto.ColDateTime64).WithPrecision(chproto.PrecisionNano),
		EndTime:        chproto.NewColNullable[time.Time](endTimeCol),
		Country:        new(chproto.ColStr),
		Device:         new(chproto.ColStr),
		OS:             new(chproto.ColStr),
		existingCombos: make(map[string]bool),
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

func (b *MobileSessionBatch) Close() {
	b.done <- struct{}{}
	b.lock.Lock()
	b.save()
	b.lock.Unlock()
	b.existingCombos = make(map[string]bool)
}

func (b *MobileSessionBatch) Add(sessionData *MobileSessionRequestType) {
	b.lock.Lock()
	defer b.lock.Unlock()
	for _, dataPoint := range sessionData.DataPoints {
		key := dataPoint.SessionId + "_" + dataPoint.UserId

		if _, exists := b.existingCombos[key]; exists {
			continue
		}
		b.existingCombos[key] = true

		b.Timestamp.Append(time.Unix(0, int64(dataPoint.Timestamp)))
		b.SessionId.Append(dataPoint.SessionId)
		b.UserId.Append(dataPoint.UserId)
		b.StartTime.Append(time.Unix(0, int64(dataPoint.StartTime)))
		b.EndTime.Append(chproto.Null[time.Time]())
		b.Country.Append(dataPoint.Country)
		b.Device.Append(dataPoint.Device)
		b.OS.Append(dataPoint.OS)
	}

	if b.SessionId.Rows() >= b.limit {
		b.save()
	}
}

func (b *MobileSessionBatch) save() {
	if b.SessionId.Rows() == 0 {
		return
	}

	input := chproto.Input{
		{Name: "Timestamp", Data: b.Timestamp},
		{Name: "SessionId", Data: b.SessionId},
		{Name: "UserId", Data: b.UserId},
		{Name: "StartTime", Data: b.StartTime},
		{Name: "EndTime", Data: b.EndTime},
		{Name: "Country", Data: b.Country},
		{Name: "Device", Data: b.Device},
		{Name: "OS", Data: b.OS},
	}

	query := ch.Query{Body: input.Into("mobile_session_data"), Input: input}

	err := b.exec(query)
	if err != nil {
		klog.Errorf("Error saving to mobile_session_data: %v", err)
	}

	for _, col := range input {
		if resettable, ok := col.Data.(chproto.Resettable); ok {
			resettable.Reset()
		}
	}

	b.existingCombos = make(map[string]bool)
}

func (c *Collector) UpdateSessionEndTime(project *db.Project, sessionId string, endTime time.Time) error {
	chClient, err := c.getClickhouseClient(project)
	if err != nil {
		return err
	}

	endTimeStr := endTime.UTC().Format("2006-01-02 15:04:05.000000000")

	query := fmt.Sprintf(
		"ALTER TABLE mobile_session_data UPDATE EndTime = toDateTime64('%s', 9) WHERE SessionId = '%s' AND EndTime IS NULL",
		endTimeStr,
		sessionId,
	)

	ctx := context.Background()
	err = chClient.pool.Do(ctx, ch.Query{Body: query})
	if err != nil {
		klog.Errorf("Failed to update EndTime for session %s: %v", sessionId, err)
		return err
	}

	return nil
}
