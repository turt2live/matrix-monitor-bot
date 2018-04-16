package metrics

import (
	"time"
	"sync"
	"github.com/patrickmn/go-cache"
	"strings"
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

	if n == 0 {
		n = 1
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

func ListDomainsWithSendTimes(exceptDomain string) []string {
	domains := make(map[string]struct{})

	sendTimesMap.Range(func(key interface{}, value interface{}) bool {
		domainKey := key.(string)
		i := strings.Index(domainKey, " TO ")
		domainA := domainKey[:i]
		domainB := domainKey[i+len(" TO "):]

		if _, ok := domains[domainA]; !ok && domainA != exceptDomain {
			domains[domainA] = struct{}{}
		}
		if _, ok := domains[domainB]; !ok && domainB != exceptDomain {
			domains[domainB] = struct{}{}
		}

		return true // more elements please
	})

	domainList := make([]string, 0, len(domains))
	for k := range domains {
		domainList = append(domainList, k)
	}

	return domainList
}
