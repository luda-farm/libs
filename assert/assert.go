package assert

func Nil(a any) {
	if a != nil {
		panic(a)
	}
}
