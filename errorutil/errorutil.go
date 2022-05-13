package errorutil

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func Must[V any](v V, e error) V {
	Check(e)
	return v
}
