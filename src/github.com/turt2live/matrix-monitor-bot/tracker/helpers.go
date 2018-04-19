package tracker

import (
	"github.com/sirupsen/logrus"
)

func RecordPing(fromDomain string, toDomain string, roomId string, eventId string, generatedTs int64, originTs int64, receivedTs int64, log *logrus.Entry) (error) {
	// TODO: Verify the event was in the specified room
	d := GetDomain(toDomain)
	r := d.GetRemote(fromDomain)
	c := r.GetRoom(roomId)
	return c.RecordPing(eventId, generatedTs, originTs, receivedTs, log)
}
