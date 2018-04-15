package events

type DisplayHints struct {
	Hints [][]string `json:"display_hints"`
}

type TextBody struct {
	Body string `json:"body"`
}

type RelatesTo struct {
	InReplyTo ReplyTo `json:"m.in_reply_to"`
}

type ReplyTo struct {
	EventId string `json:"event_id"`
}
