package matrix

import (
	"github.com/matrix-org/gomatrix"
	"github.com/sirupsen/logrus"
	"encoding/json"
	"time"
	"github.com/turt2live/matrix-monitor-bot/config"
	"github.com/turt2live/matrix-monitor-bot/util"
	"github.com/turt2live/matrix-monitor-bot/events"
	"github.com/turt2live/matrix-monitor-bot/tracker"
)

func (c *Client) handlePing(log *logrus.Entry, ev *gomatrix.Event) {
	ping := events.PingInfo{}
	pingAsStr, _ := json.Marshal(ev.Content["io.t2bot.monitor.ping"])
	_ = json.Unmarshal(pingAsStr, &ping)

	if ping.Version > 1 {
		log.Warn("Ping is of a higher version (", ping.Version, "). Replying with an older pong")
	}
	if ping.Version < 1 {
		log.Warn("Ping version is too old for processing (", ping.Version, "). Ignoring ping")
		return
	}

	domain, err := ExtractUserHomeserver(ev.Sender)
	if err != nil {
		log.Error("Error parsing domain from which we received a ping: ", err)
		return
	}
	log.Info("Ping received from ", domain)

	if domain != ping.SenderDomain {
		log.Warn("Ping domain (", ping.SenderDomain, ") does not match sender domain (", domain, "). Ignoring ping.")
		return
	}

	err = tracker.GetPingTracker().StorePing(ev.ID, ev.RoomID, ping)
	if err != nil {
		logrus.Error("Error storing ping: ", err)
		return
	}

	// Analyze the ping to see if the target server is having sending issues
	remoteSendDelay := time.Duration(ev.Timestamp-ping.GeneratedMs) * time.Millisecond
	if remoteSendDelay >= config.RemoteSendDelayThreshold || remoteSendDelay <= 0 {
		log.Warn(domain, " has a ", remoteSendDelay, " delay in sending events (origin_server_ts vs generated_ms) on ping")
	}
	if remoteSendDelay < 0 {
		remoteSendDelay = 0 // For sanity, even though it's not supposed to be possible
	}

	// TODO: Export remoteSendDelay (metric A)

	receiveDelay := (time.Duration(util.NowMillis()-ping.GeneratedMs) * time.Millisecond) - remoteSendDelay
	if receiveDelay >= config.ReceiveDelayThreshold || receiveDelay <= 0 {
		log.Warn("Ping received from ", domain, " has a receive delay of ", receiveDelay)
	}

	// TODO: Export receiveDelay (metric BC)

	response := &events.PongContent{
		Msgtype:      "m.text",
		Body:         "Pong for " + ev.ID,
		DisplayHints: events.DisplayHints{Hints: [][]string{{"io.t2bot.monitor.pong"}, {"m.text"}}},
		RelatesTo:    events.RelatesTo{InReplyTo: events.ReplyTo{EventId: ev.ID}},
		TextBody:     events.TextBody{Body: "Pong for " + ev.ID},
		PongInfo: events.PongInfo{
			Version:          1,
			InReplyTo:        ev.ID,
			ReceivedMs:       util.NowMillis(),
			ReceivedNano:     util.NowNano(),
			GeneratedMs:      util.NowMillis(),
			GeneratedNano:    util.NowNano(),
			ReceiveDelayMs:   receiveDelay.Nanoseconds() / 1000000,
			ReceiveDelayNano: receiveDelay.Nanoseconds(),
			OriginalPing:     ping,
		},
	}

	go func() {
		util.Retry(20, 500*time.Millisecond, func() error {
			log.Info("Sending pong...")
			r, err := c.mxClient.SendMessageEvent(ev.RoomID, "m.room.message", response)
			if err != nil {
				log.Error(err)
				return err
			}

			log.Info("Pong sent as ", r.EventID)
			return nil
		})
	}()
}
