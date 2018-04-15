package matrix

import (
	"github.com/sirupsen/logrus"
	"github.com/matrix-org/gomatrix"
	"encoding/json"
	"time"
	"github.com/turt2live/matrix-monitor-bot/config"
	"github.com/turt2live/matrix-monitor-bot/util"
)

func (c *Client) handlePong(log *logrus.Entry, ev *gomatrix.Event) {
	pong := pongInfo{}
	pongAsStr, _ := json.Marshal(ev.Content["io.t2bot.monitor.pong"])
	_ = json.Unmarshal(pongAsStr, &pong)

	if pong.Version > 1 {
		log.Warn("Pong is of a higher version: ", pong.Version)
	}
	if pong.Version < 1 {
		log.Warn("Pong is too old for processing: ", pong.Version)
		return
	}

	domain, err := ExtractUserHomeserver(ev.Sender)
	if err != nil {
		log.Error("Error parsing domain from which we received a pong: ", err)
		return
	}
	log.Info("Pong received for ", pong.InReplyTo, " from ", domain)

	// TODO: Verify that the original ping wasn't tampered with:
	// * Event ID exists
	// * Ping event has the same io.t2bot.monitor.ping content
	// * Ping event is in the same room

	remoteSendDelay := time.Duration(ev.Timestamp-pong.GeneratedMs) * time.Millisecond
	if remoteSendDelay >= config.RemoteSendDelayThreshold || remoteSendDelay <= 0 {
		log.Warn(domain, " has a ", remoteSendDelay, " delay in sending events (origin_server_ts vs generated_ms) on pong")
	}
	if remoteSendDelay < 0 {
		remoteSendDelay = 0 // For sanity, even though it's not supposed to be possible
	}

	// TODO: Export remoteSendDelay (metric DE)

	receiveDelay := (time.Duration(util.NowMillis()-pong.GeneratedMs) * time.Millisecond) - remoteSendDelay
	if receiveDelay >= config.ReceiveDelayThreshold || receiveDelay <= 0 {
		log.Warn("Pong received from ", domain, " has a receive delay of ", receiveDelay)
	}

	// TODO: Export receiveDelay (metric F)

	processingDelay := time.Duration(pong.GeneratedMs-pong.ReceivedMs) * time.Millisecond
	if processingDelay >= config.ProcessingDelayThreshold || processingDelay < 0 {
		log.Warn(domain, " has a processing delay of ", processingDelay)
	}

	// TODO: Export processingDelay (metric G)

	remoteDomain := pong.OriginalPing.Domain

	pingDelay := time.Duration(util.NowMillis()-pong.OriginalPing.GeneratedMs) * time.Millisecond
	pongDelay := time.Duration(util.NowMillis()-pong.GeneratedMs) * time.Millisecond
	rtt := pingDelay + pongDelay

	// TODO: Export pingDelay (metric ABC)
	// TODO: Export pongDelay (metric DEF)
	// TODO: Export rtt (metric ABCDEF - not G)

	log.Info("Ping delay (", remoteDomain, " -> ", domain, ") is ", pingDelay)
	log.Info("Pong delay (", domain, " -> ", remoteDomain, ") is ", pongDelay)
	log.Info("Round trip delay (", remoteDomain, " -> ", domain, " -> ", remoteDomain, ") is ", rtt)

	// TODO: Detect out of order pongs
	// TODO: Disregard obviously old pongs to prevent throwing off metrics from bots that are recovering
}
