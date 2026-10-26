package common

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strconv"
	"time"
)

// 获取当前日期
func GetDate(sep string) string {
	t := time.Now()

	y := strconv.Itoa(t.Year())
	m := strconv.Itoa(int(t.Month()))
	d := strconv.Itoa(t.Day())

	return y + sep + m + sep + d
}

// 获取当前日期时间
func GetDateTime() string {
	t := time.Now()

	y := strconv.Itoa(t.Year())
	m := strconv.Itoa(int(t.Month()))
	d := strconv.Itoa(t.Day())

	h := strconv.Itoa(t.Hour())
	mm := strconv.Itoa(t.Minute())
	s := strconv.Itoa(t.Second())

	return y + "-" + m + "-" + d + " " + h + ":" + mm + ":" + s
}
func JSONToMap(str string) map[string]interface{} {

	var tempMap map[string]interface{}

	err := json.Unmarshal([]byte(str), &tempMap)

	if err != nil {
		panic(err)
	}

	return tempMap
}

// 获取原始body转map
func GetRawMap(c *gin.Context) (input map[string]interface{}) {
	buf, _ := c.GetRawData()
	input = make(map[string]interface{})

	if string(buf) == "" {
		return
	}

	_ = json.Unmarshal(buf, &input)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

	return
}

// GetBodyParams 获取Body参数
func GetBodyParams(c *gin.Context) map[string]interface{} {
	params := GetRawMap(c)
	tk := c.Value(ContextKeyEnum.TokenInfo)
	if tk != nil {
		tokenInfo := tk.(*TokenInfo)
		params["enterprise_id"] = tokenInfo.EnterpriseId
	}
	return params
}

func GetEnterprise(c *gin.Context) (int64, int64) {
	var enterpriseId, adminId int64
	if c.Value(ContextKeyEnum.TokenInfo) != nil {
		tokenInfo := c.Value(ContextKeyEnum.TokenInfo).(*TokenInfo)
		enterpriseId = tokenInfo.EnterpriseId
		adminId = tokenInfo.AdminId
	}
	return enterpriseId, adminId
}

func GetEnterpriseAndMobile(c *gin.Context) (int64, int64, string) {
	var enterpriseId, adminId int64
	var mobile string
	if c.Value(ContextKeyEnum.TokenInfo) != nil {
		tokenInfo := c.Value(ContextKeyEnum.TokenInfo).(*TokenInfo)
		enterpriseId = tokenInfo.EnterpriseId
		adminId = tokenInfo.AdminId
		mobile = tokenInfo.Mobile
	}
	return enterpriseId, adminId, mobile
}
