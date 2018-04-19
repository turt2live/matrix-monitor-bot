package events

import (
	"github.com/turt2live/matrix-monitor-bot/tracker"
)

type PingInfo struct {
	Version      int                `json:"version"`
	GeneratedMs  int64              `json:"generated_ms"`
	SenderDomain string             `json:"domain"` // Legacy
	Tree         tracker.RemoteTree `json:"tree"`
}

type PingContent struct {
	Msgtype      string       `json:"msgtype"`
	Body         string       `json:"body"`
	DisplayHints DisplayHints `json:"m.display_hints"`
	TextBody     TextBody     `json:"m.text"`

	// This is the actual object we end up parsing ourselves. The rest of the stuff is so the event
	// doesn't look too atrocious in Riot/clients.
	PingInfo PingInfo `json:"io.t2bot.monitor.ping"`
}
