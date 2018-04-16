package metrics

import (
	"time"
	"sync"
	"github.com/patrickmn/go-cache"
)

var sendTimesMap = &sync.Map{} // domain_key: Cache

func CalculateSendTime(fromDomain string, toDomain string) (time.Duration) {
	domainKey := fromDomain + " TO " + toDomain
	i, ok := sendTimesMap.Load(domainKey)
	if !ok {
		return -1
	}

	c := i.(*cache.Cache)
	n := int64(0)
	s := time.Duration(0)
	for _, v := range c.Items() {
		n++
		s += (v.Object).(time.Duration)
	}

	return time.Duration(s.Nanoseconds()/n) * time.Nanosecond
}

func RecordSendTime(fromDomain string, toDomain string, duration time.Duration, id string) {
	domainKey := fromDomain + " TO " + toDomain
	i, ok := sendTimesMap.Load(domainKey)
	if !ok {
		i = cache.New(5*time.Minute, 5*time.Minute)
		sendTimesMap.Store(domainKey, i)
	}

	c := i.(*cache.Cache)
	c.Set(id, duration, cache.DefaultExpiration)
}
