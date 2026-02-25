package utils

import (
	"errors"
	"time"
)

const (
	dateLayout            = "2006-01-02"
	timeLayout            = "2006/01/02 15:04:05"
	timeAndTimezoneLayout = "2006/01/02 15:04:05-07:00"
	layoutUS              = "January 2, 2006"
	safaricomTimeLayout   = "2006-01-02 15:04:05"
)

func ParseUSFormat(date time.Time) string {
	return date.Format(layoutUS)
}

func ParseTime(timeString string) (time.Time, error) {
	parsedTime, err := time.ParseInLocation(timeLayout, timeString, time.UTC)
	if err == nil {
		return parsedTime, nil
	}

	parsedTime, err = time.ParseInLocation(timeAndTimezoneLayout, timeString, time.UTC)
	if err == nil {
		return parsedTime, nil
	}

	return parsedTime, errors.New("invalid time format")
}

func FormatDate(timeToFormat time.Time) string {
	return timeToFormat.In(time.UTC).Format(dateLayout)
}

func FormatTime(timeToFormat time.Time) string {
	return timeToFormat.In(time.UTC).Format(timeLayout)
}

func FormatDateTime(timeToFormat time.Time) string {
	return timeToFormat.In(time.UTC).Format(safaricomTimeLayout)
}
