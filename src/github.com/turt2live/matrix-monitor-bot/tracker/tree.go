package tracker

type RemoteTree map[string]RemoteTimings // domain : times
type RemoteTimings map[string]Record     // eventId : record

func CalculateRemoteTree(domainName string, roomId string) (RemoteTree) {
	tree := RemoteTree{}
	domain := GetDomain(domainName)

	for _, remoteName := range domain.GetRemotes() {
		remote := domain.GetRemote(remoteName)
		timings := RemoteTimings{}

		room := remote.GetRoom(roomId)
		for _, r := range room.GetPings() {
			timings[r.EventId] = r.Record
		}

		tree[remoteName] = timings
	}

	return tree
}
