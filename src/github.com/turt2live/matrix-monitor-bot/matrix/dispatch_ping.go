package matrix

import (
	"github.com/turt2live/matrix-monitor-bot/util"
	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
	"encoding/json"
	"github.com/turt2live/matrix-monitor-bot/events"
	"github.com/turt2live/matrix-monitor-bot/tracker"
)

type RoomPing struct {
	Ping             *events.PingInfo
	EventId          string
	ExpectingServers []string
}

type DispatchedPing struct {
	Rooms map[string]RoomPing
}

func (c *Client) DispatchPing() (*DispatchedPing, error) {
	domain, err := ExtractUserHomeserver(c.UserId)
	if err != nil {
		return nil, err
	}

	rooms, err := c.mxClient.JoinedRooms()
	if err != nil {
		return nil, err
	}

	dispatchResults := &DispatchedPing{Rooms: make(map[string]RoomPing)}

	var aggregateErr error
	for _, roomId := range rooms.JoinedRooms {
		expectingReplyFrom, err := c.GetMonitoredDomainsInRoom(roomId)
		if err != nil {
			aggregateErr = multierror.Append(aggregateErr, err)
			continue
		}

		if len(expectingReplyFrom) <= 0 {
			logrus.Warn("Empty room: ", roomId)
			continue
		}

		ping := events.PingContent{
			Msgtype:      "m.text",
			Body:         "Ping from " + domain,
			DisplayHints: events.DisplayHints{Hints: [][]string{{"io.t2bot.monitor.ping"}, {"m.text"}}},
			TextBody:     events.TextBody{Body: "Ping from " + domain},
			PingInfo: events.PingInfo{
				Version:       1,
				GeneratedMs:   util.NowMillis(),
				SenderDomain:  domain,
			},
		}

		logrus.Info("Expecting a reply from ", len(expectingReplyFrom), " servers in ", roomId, ": ", expectingReplyFrom)
		evt, err := c.mxClient.SendMessageEvent(roomId, "m.room.message", ping)
		if err != nil {
			aggregateErr = multierror.Append(aggregateErr, err)
			continue
		}
		logrus.Info("Ping in ", roomId, " is event ", evt.EventID)

		err = tracker.GetPingTracker().StorePing(evt.EventID, roomId, &ping.PingInfo)
		if err != nil {
			logrus.Error("Error storing ping: ", err)
		}
		dispatchResults.Rooms[roomId] = RoomPing{
			Ping:             &ping.PingInfo,
			EventId:          evt.EventID,
			ExpectingServers: expectingReplyFrom,
		}
	}

	return dispatchResults, aggregateErr
}

func (c *Client) GetMonitoredDomainsInRoom(roomId string) ([]string, error) {
	members, err := c.mxClient.JoinedMembers(roomId)
	if err != nil {
		return nil, err
	}

	expectingReplyFrom := make([]string, 0)
	for userId, profile := range members.Joined {
		if userId == c.UserId {
			continue // Skip ourselves
		}

		if profile.DisplayName == nil {
			logrus.Warn("User ", userId, " in ", roomId, " is not a bot (no display name)")
			continue
		}

		info := &BotInfo{}
		err := json.Unmarshal([]byte(*profile.DisplayName), info)
		if err != nil {
			logrus.Warn("User ", userId, " in ", roomId, " does not look like a bot. Display name is '", *profile.DisplayName, "'. Error parsing display name: ", err)
			continue
		}

		if info.FormatVersion > 1 {
			logrus.Warn("User ", userId, " in ", roomId, " has a newer format for bot info. We're expecting a reply, however it may be incompatible.")
		}
		if info.FormatVersion < 1 {
			logrus.Warn("Not considering ", userId, " in ", roomId, " to be a bot because the format version is too old")
			continue
		}

		// Just in case someone does something weird
		if !info.IsBot {
			logrus.Warn("User ", userId, " in ", roomId, " has a display name that is JSON and matches the format, but claims it is not a bot. Ignoring user.")
			continue
		}

		domain, err := ExtractUserHomeserver(userId)
		if err != nil {
			logrus.Warn("Error determining domain for user ", userId, " in ", roomId, ": ", err)
			continue
		}

		if domain != info.Domain {
			logrus.Warn("User ", userId, " has a mismatch between the advertised domain (", info.Domain, ") and their user ID domain. Ignoring user")
			continue
		}

		expectingReplyFrom = append(expectingReplyFrom, domain)
	}

	return expectingReplyFrom, nil
}