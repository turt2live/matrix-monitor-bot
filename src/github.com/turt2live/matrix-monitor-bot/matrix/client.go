package matrix

import (
	"github.com/matrix-org/gomatrix"

	"encoding/json"
	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
	"github.com/turt2live/matrix-monitor-bot/util"
)

type Client struct {
	mxClient      *gomatrix.Client
	info          *BotInfo
	infoStr       string
	joinedRoomIds []string

	UserId            string // readonly
	AutoAcceptInvites bool
}

type BotInfo struct {
	FormatVersion int    `json:"formatVersion"`
	IsBot         bool   `json:"isBot"`
	Domain        string `json:"domain"`

	// Other properties we care to read from bots would go here
}

func NewClient(csUrl string, accessToken string) (*Client, error) {
	client := &Client{}
	mxClient, err := gomatrix.NewClient(csUrl, "", accessToken)
	if err != nil {
		return nil, err
	}

	client.mxClient = mxClient

	resp := &WhoAmIResponse{}
	url := mxClient.BuildURL("/account/whoami")
	_, err = mxClient.MakeRequest("GET", url, nil, resp)
	if err != nil {
		return nil, err
	}

	client.UserId = resp.UserId
	mxClient.UserID = resp.UserId

	server, err := ExtractUserHomeserver(client.UserId)
	if err != nil {
		return nil, err
	}

	client.info = &BotInfo{
		FormatVersion: 1,
		IsBot:         true, // obviously
		Domain:        server,
	}

	// TODO: use extensible profiles instead of the display name when that is a thing
	b, _ := json.Marshal(client.info)
	botInfoStr := string(b)
	client.infoStr = botInfoStr
	err = mxClient.SetDisplayName(botInfoStr)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) JoinRoomByAliases(aliases []string) (error) {
	var joinError error

	for _, alias := range aliases {
		logrus.Info("Trying to join room ", alias)
		resp, err := c.mxClient.JoinRoom(alias, "", nil)
		if err != nil {
			logrus.Warn(err)
			joinError = multierror.Append(joinError, err)
			continue
		}

		logrus.Info("Joined ", resp.RoomID, " through ", alias)
		c.joinedRoomIds = append(c.joinedRoomIds, resp.RoomID)
		return nil
	}

	return joinError
}

func (c *Client) StartSync() (error) {
	syncer := c.mxClient.Syncer.(*gomatrix.DefaultSyncer)
	syncer.OnEventType("m.room.member", c.handleMembership)
	syncer.OnEventType("m.room.message", c.handleMessage)
	return c.mxClient.Sync()
}

func (c *Client) DispatchPing() (error) {
	domain, err := ExtractUserHomeserver(c.UserId)
	if err != nil {
		return err
	}

	pingContent := &PingContent{
		Msgtype:      "m.text",
		Body:         "Ping from " + domain,
		DisplayHints: displayHints{[][]string{{"io.t2bot.monitor.ping"}, {"m.text"}}},
		TextBody:     textBody{"Ping from " + domain},
		PingInfo: pingInfo{
			Version:       1,
			GeneratedMs:   util.NowMillis(),
			GeneratedNano: util.NowNano(),
			Domain:        domain,
		},
	}

	rooms, err := c.mxClient.JoinedRooms()
	if err != nil {
		return err
	}

	var aggregateErr error
	for _, roomId := range rooms.JoinedRooms {
		members, err := c.mxClient.JoinedMembers(roomId)
		if err != nil {
			aggregateErr = multierror.Append(aggregateErr, err)
			continue
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

		if len(expectingReplyFrom) <= 0 {
			logrus.Warn("Empty room: ", roomId)
			continue
		}

		logrus.Info("Expecting a reply from ", len(expectingReplyFrom), " servers in ", roomId, ": ", expectingReplyFrom)
		evt, err := c.mxClient.SendMessageEvent(roomId, "m.room.message", pingContent)
		if err != nil {
			aggregateErr = multierror.Append(aggregateErr, err)
			continue
		}
		logrus.Info("Ping in ", roomId, " is event ", evt.EventID)
	}

	return aggregateErr
}
