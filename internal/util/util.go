package util

func Filter[T any](user string, ss []T, test func(user string, album T) bool) (ret []T) {
	for _, s := range ss {
		if test(user, s) {
			ret = append(ret, s)
		}
	}
	return
}
