package utils

import "time"

func GetTimeDaysAgo(days int) time.Time {
	return time.Now().AddDate(0, 0, -days)
}
func SubtractDaysFromUnix(timestamp int64, days int) time.Time {
	return time.Unix(timestamp, 0).AddDate(0, 0, -days)
}
