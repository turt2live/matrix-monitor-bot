package ping_producer

import (
	"time"
	"github.com/turt2live/matrix-monitor-bot/matrix"
	"math/rand"
	"github.com/sirupsen/logrus"
)

type Producer struct {
	interval time.Duration
	client   *matrix.Client
}

func NewProducer(interval time.Duration, client *matrix.Client) (*Producer) {
	producer := &Producer{
		interval: interval,
		client:   client,
	}

	return producer
}

func (p *Producer) Start() {
	go func() {
		ticker := time.NewTicker(p.interval)
		for now := range ticker.C {
			logrus.Info("Scheduling a ping at ", now)
			// We add a little bit of jitter so we don't obviously look like a bot
			// It also gives us the opportunity to ensure that other bots aren't just echoing back at regular intervals
			jitter := time.Duration(rand.Int63n(int64(p.interval)))
			time.Sleep(jitter / 4)

			logrus.Info("Dispatching the ping for ", now)
			p.client.DispatchPing()
		}
	}()
}
