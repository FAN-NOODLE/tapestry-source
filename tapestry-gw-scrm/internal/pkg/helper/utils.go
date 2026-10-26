package helper

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microboss"
	"math"
	"strconv"
	"strings"
	"time"
)

const enterpriseNameKey = "store_name"

// PeriodOfTimestamp 今日凌晨到23点时间戳
func PeriodOfTimestamp() (int64, int64) {
	nowDate := time.Now().Format("2006-01-02")
	startNowDateTime := nowDate + " 00:00:00"
	endNowDateTime := nowDate + " 23:59:59"
	startTimestamp, err := time.ParseInLocation("2006-01-02 15:04:05", startNowDateTime, time.Local)
	if err != nil {
		return 0, 0
	}
	endTimestamp, err := time.ParseInLocation("2006-01-02 15:04:05", endNowDateTime, time.Local)
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
func StringTofloat64(param string) float64 {
	sales, _ := strconv.ParseFloat(param, 64)
	return sales
}
func StringToInt(param string) int {
	paramint, _ := strconv.Atoi(param)
	return paramint
}
func GetRatefloat64(param string) float64 {
	if param == "" {
		return 0
	}
	if params := strings.Split(param, "%"); len(params) > 0 {
		ratefloat := StringTofloat64(params[0]) * 100
		n10 := math.Pow10(2)
		return math.Trunc((ratefloat+0.5/n10)*n10) / n10
	}
	return 0
}
func GetRate(param string) int {
	if param == "" {
		return 0
	}
	if params := strings.Split(param, "%"); len(params) > 0 {
		ratefloat := StringTofloat64(params[0]) * 100
		return int(math.Floor(ratefloat + 0.5))
	}
	return 0
}

/**
*获取系统当前日期
 */
func GetNowDate() string {
	return time.Now().Format("2006-01-02")
}

/**
*获取当前时间戳
 */
func GetNowDatetimestamp() int64 {
	return time.Now().Unix()
}

/**
*根据日期获取凌晨的时间戳
 */
func GetDateMorningtimestamp(date string) int64 {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, date+" 00:00:00", loc)
	return theTime.Unix()
}

/**
*根据日期获取23点的时间戳
 */
func GetDateNighttimestamp(date string) int64 {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, date+" 23:59:59", loc)
	return theTime.Unix()
}
func GetFormatDate(params int64, formatType int) string {
	date := time.Unix(params, 0).Format("2006-01-02")
	if formatType == 0 {
		return date
	}
	switch formatType {
	case 1:
		date = time.Unix(params, 0).Format("2006/01/02")
	case 2:
		date = time.Unix(params, 0).Format("20060102")
	case 3:
		date = time.Unix(params, 0).Format("01/02")
	case 4:
		date = time.Unix(params, 0).Format("2006-01-02 15:04:05")
	}
	return date
}

// GetEnterpriseConfig 获取企业名称
func GetEnterpriseConfig(c *gin.Context, enterpriseID int64) string {
	// 获取企业配置
	params := common.GetRawMap(c)
	params["enterprise_id"] = enterpriseID
	params["key"] = enterpriseNameKey
	enterpriseConfigs, err := microboss.GetEnterpriseName(ctx.Wrap(c), params)
	var enterpriseConfigArr []microboss.Config
	var enterpriseConfig microboss.Config
	if err != nil {
		enterpriseConfig = microboss.Config{
			Name: "",
		}
	}

	res, _ := json.Marshal(enterpriseConfigs)
	err = json.Unmarshal(res, &enterpriseConfigArr)
	if err != nil {
		enterpriseConfig.Name = ""
	}
	for _, v := range enterpriseConfigArr {
		enterpriseConfig.Name = v.Name
	}

	return enterpriseConfig.Name

}
