package helper

import (
	"fmt"
	"strings"
	"time"
)

// 时间、日期格式
const (
	YYYY_MM_DD          = "2006-01-02"
	YYYY_MM_DD_HH_MM_SS = "2006-01-02 15:04:05"
	YYYYMMDDHHMMSS      = "20060102150405"
	YYYYMMDD            = "20060102"
	YYYYMM              = "200601"
)

// PeriodOfTimestamp 今日凌晨到23点时间戳
func PeriodOfTimestamp() (int64, int64) {
	nowDate := time.Now().Format(YYYY_MM_DD)
	startNowDateTime := nowDate + " 00:00:00"
	endNowDateTime := nowDate + " 23:59:59"
	startTimestamp, err := time.ParseInLocation(YYYY_MM_DD_HH_MM_SS, startNowDateTime, time.Local)
	if err != nil {
		return 0, 0
	}
	endTimestamp, err := time.ParseInLocation(YYYY_MM_DD_HH_MM_SS, endNowDateTime, time.Local)
	if err != nil {
		return 0, 0
	}
	return startTimestamp.Unix(), endTimestamp.Unix()
}

// YesterdayPeriodOfTimestamp 昨天凌晨到23点时间戳
func YesterdayPeriodOfTimestamp(day int64) (int64, int64) {
	start, end := PeriodOfTimestamp()
	startTimestamp := start - (day * 86400)
	endTimestamp := end - (day * 86400)
	return startTimestamp, endTimestamp
}

// GetNowDate 获取系统当前日期
func GetNowDate() string {
	return time.Now().Format(YYYY_MM_DD)
}

// GetNowDateTimestamp 获取当前时间戳
func GetNowDateTimestamp() int64 {
	return time.Now().Unix()
}

// GetDateMorningTimestamp 根据日期获取凌晨的时间戳
func GetDateMorningTimestamp(date string) int64 {

	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(YYYY_MM_DD_HH_MM_SS, date+" 00:00:00", loc)
	return theTime.Unix()
}

// StrToTime 字符串转time
func StrToTime(str string) time.Time {
	// 定义时间格式
	buf := RuleThree(strings.Contains(str, "-"), "-", "/")
	layouts := []string{
		"2006" + buf + "01" + buf + "02",
		"2006" + buf + "01" + buf + "02 15",
		"2006" + buf + "01" + buf + "02 15:04",
		"2006" + buf + "01" + buf + "02 15:04:05",
	}
	var trueLayout string
	for _, layout := range layouts {
		if len(str) == len(layout) {
			trueLayout = layout
			break
		}
	}
	// 解析时间字符串
	t, err := time.Parse(trueLayout, str)
	if err != nil {
		fmt.Println(err)
		return time.Now()
	}
	return t
}

// GetDateNightTimestamp 根据日期获取23点的时间戳
func GetDateNightTimestamp(date string) int64 {
	timeLayout := YYYY_MM_DD_HH_MM_SS
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, date+" 23:59:59", loc)
	return theTime.Unix()
}

func GetFormatDate(params int64, formatType int) string {
	if params <= 0 {
		params = time.Now().Unix()
	}
	date := time.Unix(params, 0).Format(YYYY_MM_DD)
	if formatType == 0 {
		return date
	}
	switch formatType {
	case 1:
		date = time.Unix(params, 0).Format("2006/01/02")
	case 2:
		date = time.Unix(params, 0).Format(YYYYMMDD)
	case 3:
		date = time.Unix(params, 0).Format("01/02")
	}
	return date
}
func GetFormatDateByParams(params int64) string {
	if params == 0 {
		return time.Now().Format("2006-01-02 15:04:05")
	}
	return time.Unix(params, 0).Format(YYYY_MM_DD_HH_MM_SS)
}

func GetFormatDateDay(params int64) string {
	if params == 0 {
		return time.Now().Format("2006-01-02")
	}
	return time.Unix(params, 0).Format(YYYY_MM_DD)
}

// GetDateToTimes 时间格式转时间戳
func GetDateToTimes(date string) int64 {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, date, loc)
	if theTime.Unix() < 0 {
		return 0
	}
	return theTime.Unix()
}

func Second(t time.Duration) time.Duration {
	return t * time.Second
}

func GetZeroDateByParams(params int64) int64 {
	var startNowDateTime string
	if params == 0 {
		nowDate := time.Now().Format(YYYY_MM_DD)
		startNowDateTime = nowDate + " 00:00:00"
	} else {
		Date := time.Unix(params, 0).Format(YYYY_MM_DD)
		startNowDateTime = Date + " 00:00:00"
	}
	startTimestamp, err := time.ParseInLocation(YYYY_MM_DD_HH_MM_SS, startNowDateTime, time.Local)
	if err != nil {
		return 0
	}
	return startTimestamp.Unix()

}

// GetZeroTimeByTime 获取某个时间的0点时间
func GetZeroTimeByTime(t time.Time) time.Time {
	startNowDateTime := t.Format(YYYY_MM_DD) + " 00:00:00"
	zeroTime, err := time.ParseInLocation(YYYY_MM_DD_HH_MM_SS, startNowDateTime, time.Local)
	if err != nil {
		return time.Now()
	}
	return zeroTime

}

// GetBetweenDates 根据开始日期和结束日期计算出时间段内所有日期
func GetBetweenDates(s, e string, acs bool) []string {
	var d []string
	timeFormatTpl := YYYY_MM_DD_HH_MM_SS
	if len(timeFormatTpl) != len(s) {
		timeFormatTpl = timeFormatTpl[0:len(s)]
	}
	date, err := time.Parse(timeFormatTpl, s)
	if err != nil {
		// 时间解析，异常
		return d
	}
	date2, err := time.Parse(timeFormatTpl, e)
	if err != nil {
		// 时间解析，异常
		return d
	}
	if date2.Before(date) {
		// 如果结束时间小于开始时间，异常
		return d
	}
	// 输出日期格式固定

	date2Str := date2.Format(YYYY_MM_DD)
	d = append(d, date.Format(YYYY_MM_DD))

	// 只有当天的数据
	if date2 == date {
		return d
	}

	for {
		date = date.AddDate(0, 0, 1)
		dateStr := date.Format(timeFormatTpl)
		d = append(d, dateStr)
		if dateStr == date2Str {
			break
		}
	}

	if !acs {
		for i, j := 0, len(d)-1; i < j; i, j = i+1, j-1 {
			d[i], d[j] = d[j], d[i]
		}
	}
	return d
}

// GetTimestamp 获取当天的零时和24时(强制传入类型为2006-01-02)
func GetTimestamp(timeStr string) (beginTimeNum, endTimeNum int64) {
	t, _ := time.ParseInLocation(YYYY_MM_DD, timeStr, time.Local)
	beginTimeNum = t.Unix()
	endTimeNum = beginTimeNum + 86400
	return beginTimeNum, endTimeNum
}

// StartAndEndByDayTime 指定时间的当天凌晨和结束
func StartAndEndByDayTime(dayTime int64) (startTime, endTime int64) {
	timeStr := time.Unix(dayTime, 0).Format(YYYY_MM_DD)
	tm, _ := time.Parse(YYYY_MM_DD, timeStr)
	startTime = tm.Unix() - 28800
	endTime = startTime + 86399
	return
}

// GetDayByUnix 根据时间戳获取日
func GetDayByUnix(dayTime int64) (day int) {
	timeStr := time.Unix(dayTime, 0).Format(YYYY_MM_DD)
	tm, _ := time.Parse(YYYY_MM_DD, timeStr)
	return tm.Day()
}

// MonthIntervalTime 2.获取某月的开始和结束时间mon为0本月,-1上月，1下月以此类推
func MonthIntervalTime(mon int) (startTime, endTime string) {
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	startTime = thisMonth.AddDate(0, mon, 0).Format("2006-01-02") + " 00:00:00"
	endTime = thisMonth.AddDate(0, mon+1, -1).Format("2006-01-02") + " 23:59:59"
	return startTime, endTime
}

// GetFirstDateOfWeek 获取本周周一日期
func GetFirstDateOfWeek(t time.Time) time.Time {
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).
		AddDate(0, 0, offset)
}

// GetNextFirstDateOfWeek 获取下周周一
func GetNextFirstDateOfWeek(t time.Time) time.Time {
	return GetFirstDateOfWeek(t).AddDate(0, 0, 7)
}

// GetLastDateOfWeek 获取本周周日
func GetLastDateOfWeek(t time.Time) time.Time {
	return GetFirstDateOfWeek(t).AddDate(0, 0, 6)
}

// GetLastWeekFirstDate 获取上周的周一日期
func GetLastWeekFirstDate(t time.Time) time.Time {
	return GetFirstDateOfWeek(t).AddDate(0, 0, -7)
}

// GetLastWeekLastDate 获取下周周日
func GetLastWeekLastDate(t time.Time) time.Time {
	return GetLastDateOfWeek(t).AddDate(0, 0, -7)

}

// GetBetweenByTime 获取时间范围
func GetBetweenByTime(day int64) (start, end int64) {
	now := time.Now()
	todayZero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	switch day {
	case 0:
		// 今日
		end = now.Unix()
		start = todayZero
	case 1:
		// 昨日
		end = todayZero - 1
		start = todayZero - 24*60*60
	default:
		// 近几天
		end = now.Unix()
		start = todayZero - (day-1)*24*60*60
	}

	return
}

// GetCoachInfoByTime   根据时间获取
func GetCoachInfoByTime(s1, s2 string) string {
	return RuleThree(CheckCoachTime(), s2, s1)
}

// CheckCoachTime   根据时间判断
func CheckCoachTime() bool {
	signTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.Local)
	return time.Now().After(signTime)
}
