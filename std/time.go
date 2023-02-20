package std

import (
	"fmt"
	"time"
)

func FirstOfNextMonth() time.Time {
	now := time.Now()
	year := now.Year() + int(now.Month())/12
	month := now.Month()%12 + 1
	date, err := time.Parse("20061", fmt.Sprintf("%d%d", year, month))
	if err != nil {
		panic(err)
	}

	return date
}
