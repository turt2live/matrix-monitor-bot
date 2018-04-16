package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/turt2live/matrix-monitor-bot/config"
	"time"
)

const namespace = "monbot"

// Metrics:
// [ bot ] --A-> [ matrix.org ] --B-> [ t2bot.io ] --C-> [ bot (G) ]
//                                                          |
// [ bot ] <-F-- [ matrix.org ] <-E-- [ t2bot.io ] <-D------+
// A: Ping remote send delay
// B: Ping federation delay
// C: Ping sync delay
// D: Pong send delay
// E: Pong federation delay
// F: Pong sync delay
// G: The processing delay for a ping

var pingSendDelay *prometheus.HistogramVec    // Metric A
var pingReceiveDelay *prometheus.HistogramVec // Metric BC
var pingProcessDelay *prometheus.HistogramVec // Metric G
var pongSendDelay *prometheus.HistogramVec    // Metric D
var pongReceiveDelay *prometheus.HistogramVec // Metric EF
var pingTime *prometheus.HistogramVec         // Metric ABC
var pongTime *prometheus.HistogramVec         // Metric DEF
var rtt *prometheus.HistogramVec              // Metric ABCDEF (no G)

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

	pingProcessDelay = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "ping_process_delay_seconds",
		Help:      "Number of seconds for a bot to process a ping event",
		Namespace: namespace,
	}, []string{"sourceDomain", "receivingDomain"})
	prometheus.MustRegister(pingProcessDelay)

	pongSendDelay = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "pong_send_delay_seconds",
		Help:      "Number of seconds for the origin to send a pong in response to a ping to their homeserver",
		Namespace: namespace,
	}, []string{"sourceDomain", "receivingDomain"})
	prometheus.MustRegister(pongSendDelay)

	pongReceiveDelay = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "pong_receive_delay_seconds",
		Help:      "Number of seconds for a bot to receive a pong",
		Namespace: namespace,
	}, []string{"sourceDomain", "receivingDomain"})
	prometheus.MustRegister(pongReceiveDelay)

	pingTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "ping_time_seconds",
		Help:      "Total number of seconds a ping lasts",
		Namespace: namespace,
	}, []string{"sourceDomain", "receivingDomain"})
	prometheus.MustRegister(pingTime)

	pongTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "pong_time_seconds",
		Help:      "Total number of seconds a pong lasts",
		Namespace: namespace,
	}, []string{"sourceDomain", "receivingDomain"})
	prometheus.MustRegister(pongTime)

	rtt = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "rtt_seconds",
		Help:      "Total number of seconds for a given ping/pong sequence",
		Namespace: namespace,
	}, []string{"sourceDomain", "receivingDomain"})
	prometheus.MustRegister(rtt)
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

func RecordPingProcessDelay(sourceDomain string, receivingDomain string, duration time.Duration) {
	if !config.Get().Metrics.Enabled {
		return
	}

	pingProcessDelay.With(prometheus.Labels{
		"sourceDomain":    sourceDomain,
		"receivingDomain": receivingDomain,
	}).Observe(duration.Seconds())
}

func RecordPongSendDelay(sourceDomain string, receivingDomain string, duration time.Duration) {
	if !config.Get().Metrics.Enabled {
		return
	}

	pongSendDelay.With(prometheus.Labels{
		"sourceDomain":    sourceDomain,
		"receivingDomain": receivingDomain,
	}).Observe(duration.Seconds())
}

func RecordPongReceiveDelay(sourceDomain string, receivingDomain string, duration time.Duration) {
	if !config.Get().Metrics.Enabled {
		return
	}

	pongReceiveDelay.With(prometheus.Labels{
		"sourceDomain":    sourceDomain,
		"receivingDomain": receivingDomain,
	}).Observe(duration.Seconds())
}

func RecordPongTime(sourceDomain string, receivingDomain string, duration time.Duration) {
	if !config.Get().Metrics.Enabled {
		return
	}

	pongTime.With(prometheus.Labels{
		"sourceDomain":    sourceDomain,
		"receivingDomain": receivingDomain,
	}).Observe(duration.Seconds())
}

func RecordRtt(sourceDomain string, receivingDomain string, duration time.Duration) {
	if !config.Get().Metrics.Enabled {
		return
	}

	rtt.With(prometheus.Labels{
		"sourceDomain":    sourceDomain,
		"receivingDomain": receivingDomain,
	}).Observe(duration.Seconds())
}
