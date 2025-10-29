package std

import (
	"time"
)

func FirstOfCurrentMonthUtc() time.Time {
	now := time.Now().In(time.UTC)
	year, month, _ := now.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}

func FirstOfNextMonthUtc() time.Time {
	firstOfCurrentMonth := FirstOfCurrentMonthUtc()
	return firstOfCurrentMonth.AddDate(0, 1, 0)
}
