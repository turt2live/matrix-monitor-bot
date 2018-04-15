package matrix

import (
	"github.com/matrix-org/gomatrix"
	"github.com/sirupsen/logrus"
	"github.com/turt2live/matrix-monitor-bot/util"
	"time"
)

func (c *Client) handleMembership(ev *gomatrix.Event) {
	if ev.StateKey == nil || *ev.StateKey != c.UserId || ev.Content == nil || ev.Content["membership"] != "invite" {
		return // Not an invite for us
	}

	go func() {
		util.Retry(10, time.Second, func() error {
			reply := &RoomMemberEventContent{
				DisplayName: c.infoStr,
			}

			if c.AutoAcceptInvites {
				logrus.Info("Accepting invite to ", ev.RoomID)
				reply.Membership = "join"
			} else {
				logrus.Info("Declining invite to ", ev.RoomID)
				reply.Membership = "leave"
			}

			_, err := c.mxClient.SendStateEvent(ev.RoomID, "m.room.member", c.UserId, reply)
			if err != nil {
				logrus.Error("Error replying to invite in ", ev.RoomID, ": ", err)
				return err
			}

			return nil
		})
	}()
}
