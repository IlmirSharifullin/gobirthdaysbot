package dates

import (
	"fmt"
	"time"
)

func FindNextBirthday(date time.Time) time.Time {
	today := time.Now().Truncate(time.Hour * 24)
	thisYearDate := time.Date(today.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	if thisYearDate.Before(today) {
		nextYearDate := time.Date(thisYearDate.Year()+1, date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
		return nextYearDate
	} else {
		return thisYearDate
	}
}

func CalculateDate(date time.Time) (int, string) {
	today := time.Now().Truncate(24 * time.Hour).UTC()
	d := FindNextBirthday(date)
	years := int(d.Sub(date).Truncate(365*24*time.Hour).Hours()) / (24 * 365)
	days := int(d.Sub(today).Truncate(24*time.Hour).Hours()) / 24
	if days == 0 {
		return years, "today"
	} else if days == 1 {
		return years, "tomorrow"
	} else if days == 7 {
		return years, "in a week"
	} else {
		return years, fmt.Sprintf("in %d days", days)
	}
}
