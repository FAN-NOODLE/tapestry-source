package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gw-scrm/internal/config"
	"io/ioutil"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func CostTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		p := c.Request.URL.Path
		method := c.Request.Method
		reqBody := ""

		// 如果是POST/PUT请求，并且内容类型为JSON，则读取内容体
		if method == http.MethodPost || method == http.MethodPut {
			mediaType, _, _ := mime.ParseMediaType(c.GetHeader("Content-Type"))
			if mediaType == "application/json" {
				body, err := ioutil.ReadAll(c.Request.Body)
				c.Request.Body.Close()
				if err == nil {
					buf := bytes.NewBuffer(body)
					c.Request.Body = ioutil.NopCloser(buf)
					reqBody = string(body)
				}
			}
		}
		c.Next()

		timeConsuming := time.Since(start).Seconds()
		str := fmt.Sprintf("gw-scrm项目请求异常。\n"+
			">请求时间：<font color='comment'>%vs</font>\n"+
			">请求地址：<font color='comment'>%v</font>\n"+
			">请求方式：<font color='comment'>%v</font>\n"+
			">请求IP：<font color='comment'>%v</font>\n"+
			">返回状态：<font color='comment'>%v</font>\n"+
			">请求参数：<font color='comment'>%v</font>\n", fmt.Sprintf("%.2f", timeConsuming), p, c.Request.Method, c.ClientIP(), c.Writer.Status(), reqBody)
		timeLimit := config.CfgData.Chatrobot.Time
		statusPre := string(strconv.Itoa(c.Writer.Status())[0])
		if (statusPre != "2" && statusPre != "4") || timeConsuming > float64(timeLimit) {
			RequestRobot(ctx.Wrap(c), str)
		}
	}
}

func RequestRobot(c ctx.Context, postDataStr string) {
	PostUrl := config.CfgData.Chatrobot.Url
	time.Sleep(time.Second)
	client := &http.Client{Timeout: 30 * time.Second}
	postData := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"content": postDataStr,
		},
	}
	data, _ := json.Marshal(postData)
	req, _ := http.NewRequest("POST", PostUrl, strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Info(c).Err(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Info(c).Msg(string(body))
}
