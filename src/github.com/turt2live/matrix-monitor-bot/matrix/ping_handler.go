package matrix

import (
	"github.com/matrix-org/gomatrix"
	"github.com/sirupsen/logrus"
	"encoding/json"
	"github.com/turt2live/matrix-monitor-bot/util"
	"github.com/turt2live/matrix-monitor-bot/events"
	"github.com/turt2live/matrix-monitor-bot/tracker"
)

func (c *Client) handlePing(log *logrus.Entry, ev *gomatrix.Event) {
	ping := events.PingInfo{}
	pingAsStr, _ := json.Marshal(ev.Content["io.t2bot.monitor.ping"])
	_ = json.Unmarshal(pingAsStr, &ping)

	if ping.Version > 2 {
		log.Warn("Ping is of a higher version (", ping.Version, "). Will attempt to parse")
	}
	if ping.Version < 1 {
		log.Warn("Ping version is too old for processing (", ping.Version, "). Ignoring ping")
		return
	}
	if ping.Version < 2 {
		log.Warn("Ping version is old, but compatible (", ping.Version, ")")
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

	receivedMs := util.NowMillis()
	tracker.RecordPing(domain, c.Domain, ev.RoomID, ev.ID, ping.GeneratedMs, ev.Timestamp, receivedMs, log)

	// Parse the remote tree
	if ping.Tree != nil {
		for k, v := range ping.Tree {
			for eventId, record := range v {
				tracker.RecordPing(k, domain, ev.RoomID, eventId, record.GeneratedTs, record.OriginTs, record.ReceivedTs, log)
			}
		}
	}
}
