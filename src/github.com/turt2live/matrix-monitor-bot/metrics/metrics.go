package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/turt2live/matrix-monitor-bot/config"
	"time"
)

const namespace = "monbot"

// Metrics:
// [ bot ] --A-> [ matrix.org ] --B-> [ t2bot.io ] --C-> [ bot ]
// A: Ping remote send delay
// B: Ping federation delay
// C: Ping sync delay

var pingSendDelay *prometheus.HistogramVec    // Metric A
var pingReceiveDelay *prometheus.HistogramVec // Metric BC
var pingTime *prometheus.HistogramVec         // Metric ABC

// TODO: Calculate and export time between pings
// TODO: Detect and export missed pings (by threshold)
// TODO: Detect and export missed pongs (by threshold)

func initMetrics() {
	logrus.Info("Creating metrics...")

	pingSendDelay = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "ping_send_delay_seconds",
		Help:      "Number of seconds for the origin to send a ping to their homeserver",
		Namespace: namespace,
	}, []string{"sourceDomain"})
	prometheus.MustRegister(pingSendDelay)

	pingReceiveDelay = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "ping_receive_delay_seconds",
		Help:      "Number of seconds for a bot to receive a ping",
		Namespace: namespace,
	}, []string{"sourceDomain", "receivingDomain"})
	prometheus.MustRegister(pingReceiveDelay)

	pingTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "ping_time_seconds",
		Help:      "Total number of seconds a ping lasts",
		Namespace: namespace,
	}, []string{"sourceDomain", "receivingDomain"})
	prometheus.MustRegister(pingTime)
}

func RecordPingSendDelay(domain string, duration time.Duration) {
	if !config.Get().Metrics.Enabled {
		return
	}

	pingSendDelay.With(prometheus.Labels{
		"sourceDomain": domain,
	}).Observe(duration.Seconds())
}

func RecordPingReceiveDelay(sourceDomain string, receivingDomain string, duration time.Duration) {
	if !config.Get().Metrics.Enabled {
		return
	}

	pingReceiveDelay.With(prometheus.Labels{
		"sourceDomain":    sourceDomain,
		"receivingDomain": receivingDomain,
	}).Observe(duration.Seconds())
}

func RecordPingTime(sourceDomain string, receivingDomain string, duration time.Duration) {
	if !config.Get().Metrics.Enabled {
		return
	}

	pingTime.With(prometheus.Labels{
		"sourceDomain":    sourceDomain,
		"receivingDomain": receivingDomain,
	}).Observe(duration.Seconds())
}
