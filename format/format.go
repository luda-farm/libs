package format

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func Cents[I constraints.Integer](i I) string {
	if i < 0 {
		return fmt.Sprintf("%d.%02d", i/100, -i%100)
	} else {
		return fmt.Sprintf("%d.%02d", i/100, i%100)
	}
}
