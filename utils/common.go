package utils

func IfThen[V any](ok bool, v1, v2 V) V {
	if ok {
		return v1
	}
	return v2
}
