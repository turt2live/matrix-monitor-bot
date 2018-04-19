package tracker

import (
	"github.com/patrickmn/go-cache"
	"github.com/turt2live/matrix-monitor-bot/config"
	"time"
	"github.com/pkg/errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/turt2live/matrix-monitor-bot/metrics"
)

type Room struct {
	remoteDomain *RemoteDomain
	timings      *cache.Cache
}

type Record struct {
	GeneratedTs int64 `json:"generated_ts"`
	OriginTs    int64 `json:"origin_ts"`
	ReceivedTs  int64 `json:"received_ts"`
}

type PingRecord struct {
	EventId string
	Record  Record
}

func (r *Room) RecordPing(eventId string, generatedTs int64, originTs int64, receivedTs int64, log *logrus.Entry) (error) {
	age := time.Duration(receivedTs-generatedTs) * time.Millisecond
	if age >= config.WebAverageInterval {
		return errors.New("Event too old to store: " + fmt.Sprint(age))
	}
	_, exists := r.timings.Get(eventId)
	if exists {
		log.Warn("Event ID ", eventId, " already stored")
		return nil // Not technically an error, but we don't want to trigger stats
	}

	record := Record{
		GeneratedTs: generatedTs,
		OriginTs:    originTs,
		ReceivedTs:  receivedTs,
	}
	r.timings.Set(eventId, record, cache.DefaultExpiration)
	r.recordMetrics(record, log)

	return nil
}

func (r *Room) GetPings() ([]PingRecord) {
	records := make([]PingRecord, 0, r.timings.ItemCount())

	for k, v := range r.timings.Items() {
		records = append(records, PingRecord{EventId: k, Record: v.Object.(Record)})
	}

	return records
}

func (r *Room) recordMetrics(record Record, log *logrus.Entry) {
	remoteSendDelay := time.Duration(record.OriginTs-record.GeneratedTs) * time.Millisecond
	receiveDelay := time.Duration(record.ReceivedTs-record.OriginTs) * time.Millisecond
	pingTime := time.Duration(record.ReceivedTs-record.GeneratedTs) * time.Millisecond

	sourceDomain := r.remoteDomain.name
	receivingDomain := r.remoteDomain.domain.name

	metrics.RecordPingSendDelay(sourceDomain, remoteSendDelay)
	if remoteSendDelay >= config.RemoteSendDelayWarnThreshold || remoteSendDelay <= 0 {
		log.Warn(sourceDomain, " has a ", remoteSendDelay, " delay in sending events to their homeserver")
	}

	metrics.RecordPingReceiveDelay(sourceDomain, receivingDomain, receiveDelay)
	if receiveDelay >= config.ReceiveDelayWarnThreshold || receiveDelay <= 0 {
		log.Warn(receivingDomain, " has a ", receiveDelay, " delay in receiving events from ", sourceDomain)
	}

	metrics.RecordPingTime(sourceDomain, receivingDomain, pingTime)
	if pingTime >= config.PingTimeWarnThreshold || pingTime <= 0 {
		log.Warn("Ping time for ", sourceDomain, " -> ", receivingDomain, " is ", pingTime)
	}
}
