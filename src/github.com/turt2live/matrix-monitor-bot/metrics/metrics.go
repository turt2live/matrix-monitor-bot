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

var pingSendDelay *prometheus.SummaryVec    // Metric A
var pingReceiveDelay *prometheus.SummaryVec // Metric BC
var pingProcessDelay *prometheus.SummaryVec // Metric G
var pongSendDelay *prometheus.SummaryVec    // Metric D
var pongReceiveDelay *prometheus.SummaryVec // Metric EF
var pingTime *prometheus.SummaryVec         // Metric ABC
var pongTime *prometheus.SummaryVec         // Metric DEF
var rtt *prometheus.SummaryVec              // Metric ABCDEF (no G)

// TODO: Calculate and export time between pings
// TODO: Detect and export missed pings (by threshold)
// TODO: Detect and export missed pongs (by threshold)

func initMetrics() {
	logrus.Info("Creating metrics...")

	pingSendDelay = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      "ping_send_delay_seconds",
		Help:      "Number of seconds for the origin to send a ping to their homeserver",
		Namespace: namespace,
	}, []string{"sourceDomain"})
	prometheus.MustRegister(pingSendDelay)

	pingReceiveDelay = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      "ping_receive_delay_seconds",
		Help:      "Number of seconds for a bot to receive a ping",
		Namespace: namespace,
	}, []string{"sourceDomain", "receivingDomain"})
	prometheus.MustRegister(pingReceiveDelay)

	pingProcessDelay = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      "ping_process_delay_seconds",
		Help:      "Number of seconds for a bot to process a ping event",
		Namespace: namespace,
	}, []string{"sourceDomain", "receivingDomain"})
	prometheus.MustRegister(pingProcessDelay)

	pongSendDelay = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      "pong_send_delay_seconds",
		Help:      "Number of seconds for the origin to send a pong in response to a ping to their homeserver",
		Namespace: namespace,
	}, []string{"sourceDomain", "receivingDomain"})
	prometheus.MustRegister(pongSendDelay)

	pongReceiveDelay = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      "pong_receive_delay_seconds",
		Help:      "Number of seconds for a bot to receive a pong",
		Namespace: namespace,
	}, []string{"sourceDomain", "receivingDomain"})
	prometheus.MustRegister(pongReceiveDelay)

	pingTime = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      "ping_time_seconds",
		Help:      "Total number of seconds a ping lasts",
		Namespace: namespace,
	}, []string{"sourceDomain", "receivingDomain"})
	prometheus.MustRegister(pingTime)

	pongTime = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      "pong_time_seconds",
		Help:      "Total number of seconds a pong lasts",
		Namespace: namespace,
	}, []string{"sourceDomain", "receivingDomain"})
	prometheus.MustRegister(pongTime)

	rtt = prometheus.NewSummaryVec(prometheus.SummaryOpts{
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
