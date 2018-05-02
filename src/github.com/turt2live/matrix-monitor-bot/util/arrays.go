package util

func StrArrayContains(a []string, e string) bool {
	for _, i := range a {
		if i == e {
			return true
		}
	}

	return false
}
