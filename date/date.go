package date

import "time"

func UnixTimeStampToDate(timeStamp int64) time.Time {
	return time.Unix(timeStamp, 0)
}

func Now() int64 {
	return time.Now().Unix()
}

func GoToDay(day int) int64 {
	return time.Now().Add(time.Duration(day*24) * time.Hour).Unix()
}

func GoToHour(hour int) int64 {
	return time.Now().Add(time.Duration(hour) * time.Hour).Unix()
}

func GoToMinute(min int) int64 {
	return time.Now().Add(time.Duration(min) * time.Minute).Unix()
}
