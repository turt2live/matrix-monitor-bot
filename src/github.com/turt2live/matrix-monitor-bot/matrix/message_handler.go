package matrix

import (
	"github.com/matrix-org/gomatrix"
	"github.com/sirupsen/logrus"
)

// Shared structs for message building
// =====================================================
type displayHints struct {
	Hints [][]string `json:"display_hints"`
}

type textBody struct {
	Body string `json:"body"`
}

type relatesTo struct {
	InReplyTo replyTo `json:"m.in_reply_to"`
}

type replyTo struct {
	EventId string `json:"event_id"`
}

// =====================================================

func (c *Client) handleMessage(ev *gomatrix.Event) {
	log := logrus.WithFields(logrus.Fields{
		"sender":  ev.Sender,
		"eventId": ev.ID,
		"roomId":  ev.RoomID,
	})

	if ev.Content == nil {
		log.Warn("Event has no content (redacted?)")
		return
	}

	if ev.Content["io.t2bot.monitor.ping"] != nil {
		if ev.Sender == c.UserId {
			return // Don't pong ourselves
		}
		c.handlePing(log, ev)
		return
	}

	if ev.Content["io.t2bot.monitor.pong"] != nil {
		c.handlePong(log, ev)
		return
	}

	log.Warn("Unexpected event - is someone talking in the monitor room?")
}
