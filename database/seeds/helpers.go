package seeds

import (
	"fmt"
	"time"
)

func ParseDate(date string) time.Time {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(fmt.Sprintf("invalid date: %s", date))
	}

	return t
}

func ParseDatePtr(date string) *time.Time {
	t := ParseDate(date)
	return &t
}