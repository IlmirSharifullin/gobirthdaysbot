package dates

import (
	"fmt"
	"time"
)

func FindNextBirthday(date time.Time) time.Time {
	today := time.Now().Truncate(time.Hour * 24)
	if c := date.Compare(today); c >= 0 {
		return date.Truncate(24 * time.Hour)
	}
	thisYearDate := time.Date(today.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	if thisYearDate.Before(today) {
		nextYearDate := time.Date(thisYearDate.Year()+1, date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
		return nextYearDate
	} else {
		return thisYearDate
	}
}

func CalculateDate(date time.Time) string {
	today := time.Now().Truncate(24 * time.Hour)
	d := FindNextBirthday(date)
	days := int(d.Sub(today).Truncate(24*time.Hour).Hours()) / 24
	if days == 0 {
		return "today"
	} else if days == 1 {
		return "tomorrow"
	} else if days == 7 {
		return "in a week"
	} else {
		return fmt.Sprintf("in %d days", days)
	}
}
