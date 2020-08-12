package utils

import (
	"time"
)

const (
	TIME_UTC_FORMAT = "2006-01-02T15:04:05Z"
)

// GetTodayStartUnix 获取今日的开头
func GetTodayStartUnix() int64 {
	dayStart := GetUnixByFormat("2006-01-02")
	return dayStart
}

// GetMonthStartUnix 获取当前月份的开头
func GetMonthStartUnix() int64 {
	monthStart := GetUnixByFormat("2006-01")
	return monthStart
}

// GetYearStartUnix 获取当前年份的开头
func GetYearStartUnix() int64 {
	monthStart := GetUnixByFormat("2006")
	return monthStart
}

// GetUnixByFormat 获取当前format的时间
func GetUnixByFormat(format string) int64 {
	date := time.Now().Format(format)
	return GetUnixByTimeFormat(date, format)
}

// GetUnixByUnixFormat 获取具体时间format的时间
func GetUnixByUnixFormat(unix int64, format string) int64 {
	date := time.Unix(unix, 0).Format(format)
	return GetUnixByTimeFormat(date, format)
}

// GetUnixByTimeFormat 获取time的format的时间
func GetUnixByTimeFormat(timeShow, format string) int64 {
	location, _ := time.ParseInLocation(format, timeShow, time.Local)
	return location.Unix()
}

// GetCSTDateByUnix 获取unix的正确日期(注意，只能用于获取日期)
func GetCSTDateByUnix(unix int64) time.Time {
	// 先获取cst和utc的差值(固定8小时)
	var between int64
	between = 8 * 60 * 60
	// 再转换
	dateTime := time.Unix(unix+between, 0)
	return dateTime
}

// GetUnixByCSTDate 获取正确日期的unix(注意，只能从日期转换)
func GetUnixByCSTDate(date time.Time) int64 {
	unix := date.Unix()
	unix -= 8 * 60 * 60
	return unix
}

// IsSameDay 是否是同一天
func IsSameDay(t1, t2 int64) bool {
	day1 := time.Unix(t1, 0).YearDay()
	year1 := time.Unix(t1, 0).Year()
	day2 := time.Unix(t2, 0).YearDay()
	year2 := time.Unix(t2, 0).Year()
	if day1 == day2 && year1 == year2 {
		return true
	}
	return false
}
