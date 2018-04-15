package events

type PingInfo struct {
	Version int `json:"version"`

	// The timestamp we generated the ping at
	GeneratedMs   int64 `json:"generated_ms"`

	// The domain is provided for ease of troubleshooting pongs
	SenderDomain string `json:"domain"`
}

type PongInfo struct {
	Version int `json:"version"`

	// The event we're responding to
	InReplyTo string `json:"in_reply_to"`

	// The timestamp when we received the ping
	ReceivedMs   int64 `json:"received_ms"`

	// The timestamp we generated the pong at
	// This is specified for clarity, despite the received timestamp usually being the same
	GeneratedMs   int64 `json:"generated_ms"`


	ReceiveDelayMs   int64 `json:"receive_delay_ms"`

	OriginalPing PingInfo `json:"original_ping"`
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

type PongContent struct {
	Msgtype      string       `json:"msgtype"`
	Body         string       `json:"body"`
	DisplayHints DisplayHints `json:"m.display_hints"`
	TextBody     TextBody     `json:"m.text"`
	RelatesTo    RelatesTo    `json:"m.relates_to"`

	// This is the actual object we end up parsing ourselves. The rest of the stuff is so the event
	// doesn't look too atrocious in Riot/clients.
	PongInfo PongInfo `json:"io.t2bot.monitor.pong"`
}
