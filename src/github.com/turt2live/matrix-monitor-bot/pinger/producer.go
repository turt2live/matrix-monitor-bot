package pinger

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
		// Because the ticker doesn't send anything until $interval, we'll trigger a ping manually
		doPing(time.Now(), p)

		ticker := time.NewTicker(p.interval)
		for now := range ticker.C {
			logrus.Info("Scheduling a ping at ", now)

			// We add a little bit of jitter so we don't obviously look like a bot
			// It also gives us the opportunity to ensure that other bots aren't just echoing back at regular intervals
			jitter := time.Duration(rand.Int63n(int64(p.interval)))
			time.Sleep(jitter / 4)

			doPing(now, p)
		}
	}()
}

func doPing(now time.Time, p *Producer) {
	logrus.Info("Dispatching the ping for ", now)
	err := p.client.DispatchPing()
	if err != nil {
		logrus.Error("Error producing ping: ", err)
	}
}
