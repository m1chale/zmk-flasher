package slices

func GetAddedElements[T any](left []T, right []T, compare func(T, T) bool) (ret []T) {
	for _, r := range right {
		found := false
		for _, l := range left {
			if compare(l, r) {
				found = true
				break
			}
		}
		if !found {
			ret = append(ret, r)
		}
	}
	return
}

func GetRemovedElements[T any](left []T, right []T, compare func(T, T) bool) (ret []T) {
	for _, r := range left {
		found := false
		for _, l := range right {
			if compare(l, r) {
				found = true
				break
			}
		}
		if !found {
			ret = append(ret, r)
		}
	}
	return
}
