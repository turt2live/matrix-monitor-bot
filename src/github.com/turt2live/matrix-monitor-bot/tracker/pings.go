package tracker

import (
	"sync"
	"github.com/pkg/errors"
	"time"
	"github.com/turt2live/matrix-monitor-bot/util"
	"github.com/turt2live/matrix-monitor-bot/config"
	"github.com/sirupsen/logrus"
	"fmt"
	"github.com/turt2live/matrix-monitor-bot/events"
)

type PingTracker struct {
	pingsByDomain *sync.Map // domain: {roomId: []TrackedPing}
}

type trackedPing struct {
	EventId string
	Info    *events.PingInfo
}

var ptInstance *PingTracker

func init() {
	ptInstance = &PingTracker{
		pingsByDomain: &sync.Map{},
	}
}

func GetPingTracker() (*PingTracker) {
	return ptInstance
}

func (p *PingTracker) TryGetPing(eventId string, roomId string, domain string) (*events.PingInfo, int, error) {
	i, ok := p.pingsByDomain.Load(domain)
	if !ok {
		return nil, -1, errors.New("Domain not found")
	}

	byRoom := i.(*sync.Map)
	i, ok = byRoom.Load(roomId)
	if !ok {
		return nil, -1, errors.New("Room not found")
	}

	pings := i.([]trackedPing)
	for idx, ping := range pings {
		if ping.EventId != eventId {
			continue
		}

		age := time.Duration(util.NowMillis()-ping.Info.GeneratedMs) * time.Millisecond
		if age >= config.PingTtl {
			logrus.Warn("Ping lookup requested but it has expired. Removing from the collection. The ping is: ", eventId, " in ", roomId, " sent by ", domain)
			fmt.Println(pings)
			byRoom.Store(roomId, append(pings[:idx], pings[idx+1:]...))
			fmt.Println(byRoom.Load(roomId))
			return nil, -1, errors.New("Ping expired")
		}

		return ping.Info, idx, nil
	}

	return nil, -1, errors.New("Ping not found")
}

func (p *PingTracker) StorePing(eventId string, roomId string, ping *events.PingInfo) (error) {
	age := time.Duration(util.NowMillis()-ping.GeneratedMs) * time.Millisecond
	if age >= config.PingTtl {
		return errors.New("Ping is too old to be stored")
	}

	i, ok := p.pingsByDomain.Load(ping.SenderDomain)
	if !ok {
		i = &sync.Map{}
		p.pingsByDomain.Store(ping.SenderDomain, i)
	}

	byRoom := i.(*sync.Map)
	i, ok = byRoom.Load(roomId)
	if !ok {
		i = make([]trackedPing, 0)
		byRoom.Store(roomId, i)
	}

	pings := i.([]trackedPing)
	pings = append(pings, trackedPing{Info: ping, EventId: eventId})
	byRoom.Store(roomId, pings)

	return nil
}
