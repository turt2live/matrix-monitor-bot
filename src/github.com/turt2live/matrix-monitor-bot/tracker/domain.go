package tracker

import (
	"sync"
	"time"
)

type Domain struct {
	name    string
	remotes *sync.Map
}

type DomainTimings struct {
	Send       time.Duration
	Receive    time.Duration
	HasSend    bool
	HasReceive bool
}

var domainCache = &sync.Map{}

func GetDomain(name string) (*Domain) {
	i, ok := domainCache.Load(name)
	if !ok {
		i = ResetDomain(name)
	}

	return i.(*Domain)
}

func ResetDomain(name string) (*Domain) {
	i := &Domain{remotes: &sync.Map{}, name: name}
	domainCache.Store(name, i)
	return i
}

func (d *Domain) GetRemote(name string) (*RemoteDomain) {
	i, ok := d.remotes.Load(name)
	if !ok {
		i = d.ResetRemote(name)
	}

	return i.(*RemoteDomain)
}

func (d *Domain) ResetRemote(name string) (*RemoteDomain) {
	i := &RemoteDomain{rooms: &sync.Map{}, domain: d, name: name}
	d.remotes.Store(name, i)
	return i
}

func (d *Domain) GetRemotes() ([]string) {
	r := make(map[string]interface{})
	d.remotes.Range(func(k interface{}, v interface{}) bool {
		r[k.(string)] = 1
		return true // keep going if we can
	})

	remotes := make([]string, 0, len(r))
	for k := range r {
		remotes = append(remotes, k)
	}

	return remotes
}

func (d *Domain) CompareTo(other string) DomainTimings {
	otherDomain := GetDomain(other)
	remote := d.GetRemote(other)
	usRemote := otherDomain.GetRemote(d.name)

	// Calculate the send times first (us -> them)
	sendTimeTotal := int64(0)
	sendTimeCount := 0
	for _, roomId := range remote.GetRooms() {
		room := remote.GetRoom(roomId)

		for _, ping := range room.GetPings() {
			sendTimeCount++
			sendTimeTotal += ping.Record.ReceivedTs - ping.Record.GeneratedTs
		}
	}

	// Calculate the receive times (them -> us)
	receiveTimeTotal := int64(0)
	receiveTimeCount := 0
	for _, roomId := range usRemote.GetRooms() {
		room := usRemote.GetRoom(roomId)

		for _, ping := range room.GetPings() {
			receiveTimeCount++
			receiveTimeTotal += ping.Record.ReceivedTs - ping.Record.GeneratedTs
		}
	}

	times := DomainTimings{}

	if sendTimeCount > 0 {
		times.HasSend = true
		avg := float64(sendTimeTotal) / float64(sendTimeCount)
		times.Send = time.Duration(avg) * time.Millisecond
	} else {
		times.HasSend = false
	}

	if receiveTimeCount > 0 {
		times.HasReceive = true
		avg := float64(receiveTimeTotal) / float64(receiveTimeCount)
		times.Receive = time.Duration(avg) * time.Millisecond
	} else {
		times.HasReceive = false
	}

	return times
}

func GetDomains() []string {
	return GetDomainsExcept() // No exceptions
}

func GetDomainsExcept(excludeDomains ...string) []string {
	d := make(map[string]interface{})
	domainCache.Range(func(k interface{}, v interface{}) bool {
		d[k.(string)] = 1
		return true // keep going if we can
	})

	domains := make([]string, 0, len(d))
	for k := range d {
		excluded := false
		for _, e := range excludeDomains {
			if e == k {
				excluded = true
				break
			}
		}

		if !excluded {
			domains = append(domains, k)
		}
	}

	return domains
}
