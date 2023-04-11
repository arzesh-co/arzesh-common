package date

import (
	"github.com/hablullah/go-hijri"
	ptime "github.com/yaa110/go-persian-calendar"
	"strconv"
	"time"
)

func convertorToInt64(n any) int64 {
	switch n := n.(type) {
	case int:
		return int64(n)
	case int8:
		return int64(n)
	case int16:
		return int64(n)
	case int32:
		return int64(n)
	case int64:
		return int64(n)
	case string:
		num, err := strconv.ParseInt(n, 10, 64)
		if err == nil {
			return int64(0)
		}
		return num
	}
	return int64(0)
}

func UnixTimeStampToDate(timeStamp int64) time.Time {
	return time.Unix(timeStamp, 0)
}

func UnixTimeStampToStringDateByLocation(timeStamp int64, Location string) string {
	date := UnixTimeStampToDate(timeStamp)
	switch Location {
	case "Asia/Tehran":
		pt := ptime.New(date)
		return pt.Format("yyyy-MM-dd")
	case "Asia/Dubai":
		ht, err := hijri.CreateHijriDate(date, hijri.Default)
		if err != nil {
			return ""
		}
		sYear := strconv.Itoa(int(ht.Year))
		sMonth := strconv.Itoa(int(ht.Month))
		sDay := strconv.Itoa(int(ht.Day))
		return sYear + "-" + sMonth + "-" + sDay
	default:
		return date.Format("2006-02-01")
	}
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

func StartToday(days any) any {
	year, month, day := time.Now().Date()
	theTime := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return theTime.Unix()
}
func NDayBefore(days any) any {
	intDay := convertorToInt64(days)
	return time.Now().Add((time.Duration(intDay) * -24) * time.Hour).Unix()
}
func NDayAfter(days any) any {
	intDay := convertorToInt64(days)
	return time.Now().Add((time.Duration(intDay) * -24) * time.Hour).Unix()
}
func BeginOfThisYear(days any) any {
	year, _, _ := time.Now().Date()
	theTime := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	return theTime.Unix()
}
func BeginOfThisMonth(days any) any {
	year, month, _ := time.Now().Date()
	theTime := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	return theTime.Unix()
}
