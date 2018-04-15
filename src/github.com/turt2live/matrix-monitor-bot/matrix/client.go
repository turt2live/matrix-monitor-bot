package matrix

import (
	"github.com/matrix-org/gomatrix"

	"encoding/json"
	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
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
