package matrix

import (
	"github.com/matrix-org/gomatrix"
	"github.com/sirupsen/logrus"
	"encoding/json"
	"time"
	"github.com/turt2live/matrix-monitor-bot/config"
	"github.com/turt2live/matrix-monitor-bot/util"
)

type pingInfo struct {
	Version int `json:"version"`

	// The timestamp we generated the ping at
	GeneratedMs   int64 `json:"generated_ms"`
	GeneratedNano int64 `json:"generated_nano"`

	// The domain is provided for ease of troubleshooting pongs
	SenderDomain string `json:"domain"`
}

type pongInfo struct {
	Version int `json:"version"`

	// The event we're responding to
	InReplyTo string `json:"in_reply_to"`

	// The timestamp when we received the ping
	ReceivedMs   int64 `json:"received_ms"`
	ReceivedNano int64 `json:"received_nano"`

	// The timestamp we generated the pong at
	// This is specified for clarity, despite the received timestamp usually being the same
	GeneratedMs   int64 `json:"generated_ms"`
	GeneratedNano int64 `json:"generated_nano"`

	// The time it took to receive the event over federation
	ReceiveDelayMs   int64 `json:"receive_delay_ms"`
	ReceiveDelayNano int64 `json:"receive_delay_nano"`

	OriginalPing pingInfo `json:"original_ping"`
}

type PingContent struct {
	Msgtype      string       `json:"msgtype"`
	Body         string       `json:"body"`
	DisplayHints displayHints `json:"m.display_hints"`
	TextBody     textBody     `json:"m.text"`

	// This is the actual object we end up parsing ourselves. The rest of the stuff is so the event
	// doesn't look too atrocious in Riot/clients.
	PingInfo pingInfo `json:"io.t2bot.monitor.ping"`
}

type PongContent struct {
	Msgtype      string       `json:"msgtype"`
	Body         string       `json:"body"`
	DisplayHints displayHints `json:"m.display_hints"`
	TextBody     textBody     `json:"m.text"`
	RelatesTo    relatesTo    `json:"m.relates_to"`

	// This is the actual object we end up parsing ourselves. The rest of the stuff is so the event
	// doesn't look too atrocious in Riot/clients.
	PongInfo pongInfo `json:"io.t2bot.monitor.pong"`
}

func (c *Client) handlePing(log *logrus.Entry, ev *gomatrix.Event) {
	ping := pingInfo{}
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

	response := &PongContent{
		Msgtype:      "m.text",
		Body:         "Pong for " + ev.ID,
		DisplayHints: displayHints{[][]string{{"io.t2bot.monitor.pong"}, {"m.text"}}},
		RelatesTo:    relatesTo{replyTo{EventId: ev.ID}},
		TextBody:     textBody{"Pong for " + ev.ID},
		PongInfo: pongInfo{
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
