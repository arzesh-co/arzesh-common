package creators

import "time"

func ConvertUnixTimestampToDate(unix int64) time.Time {
	return time.Unix(unix, 0)
}
func NowTime() int64 {
	return time.Now().Unix()
}
func AddDaysToNow(days int) int64 {
	return time.Now().Add(time.Duration(days) * 24 * time.Hour).Unix()
}
func AddHourToNow(hour int) int64 {
	return time.Now().Add(time.Duration(hour) * time.Hour).Unix()
}
func AddMinuteToNow(minute int) int64 {
	return time.Now().Add(time.Duration(minute) * time.Minute).Unix()
}
