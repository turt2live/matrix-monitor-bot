package tracker

import (
	"sync"
	"github.com/patrickmn/go-cache"
	"github.com/turt2live/matrix-monitor-bot/config"
)

type RemoteDomain struct {
	domain *Domain
	name   string
	rooms  *sync.Map
}

func (r *RemoteDomain) GetRoom(roomId string) (*Room) {
	i, ok := r.rooms.Load(roomId)
	if !ok {
		i = &Room{
			remoteDomain: r,
			timings:      cache.New(config.WebAverageInterval, config.WebAverageInterval),
		}
		r.rooms.Store(roomId, i)
	}

	return i.(*Room)
}

func (r *RemoteDomain) GetRooms() ([]string) {
	m := make(map[string]interface{})
	r.rooms.Range(func(k interface{}, v interface{}) bool {
		m[k.(string)] = 1
		return true // keep going if we can
	})

	rooms := make([]string, 0, len(m))
	for k := range m {
		rooms = append(rooms, k)
	}

	return rooms
}
