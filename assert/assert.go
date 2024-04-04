package assert

import "fmt"

func Nil(a any) {
	if a != nil {
		panic(fmt.Errorf("assert.Nil(%v)", a))
	}
}

func NotNil(a any) {
	if a == nil {
		panic(fmt.Errorf("assert.NotNil(%v)", a))
	}
}
