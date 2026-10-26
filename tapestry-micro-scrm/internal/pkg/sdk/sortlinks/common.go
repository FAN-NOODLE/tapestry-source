package sortlinks

import (
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/http"
)

type SortLinkInfo struct {
	ID       int64     `json:"id"`
	OriginID int64     `json:"origin_id"`
	PID      int64     `json:"pid"`
	Level    int64     `json:"level"`
	Args     []ArgInfo `json:"args"`
}

type ArgInfo struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Passed int8   `json:"passed"`
}

// sendGet 发送get请求
func sendGet(c ctx.Context, urlStr string, response *SortLinkInfo) error {

	log.Info(c).Msgf("request sortLinks interface, url:%s", urlStr)
	var options []clients.Option

	res, err := http.CdpClientObject.Get(c, urlStr, options...)
	if err != nil {
		log.Info(c).Msgf("sortLinks接口请求失败：%s, url:%s", err.Error(), urlStr)
		return errorx.NewErrorMsg(errs.ErrorStatus, fmt.Sprintf("sortLinks接口请求失败：%s", err.Error()))
	}
	resStr, err := res.Read(clients.Text)
	if err != nil {
		log.Info(c).Msgf("sortLinks接口返回数据获取失败：%s, url:%s", err.Error(), urlStr)
		return errorx.NewErrorMsg(errs.ErrorStatus, fmt.Sprintf("sortLinks接口返回数据获取失败：%s, url:%s", err.Error(), urlStr))
	}
	log.Info(c).Msgf("request sortLinks interface, url:%s, response:%s", urlStr, resStr.(string))
	err = json.Unmarshal([]byte(resStr.(string)), &response)
	if err != nil {
		log.Info(c).Msgf("sortLinks接口返回数据解析失败：%s, url:%s", err.Error(), urlStr)
		return errorx.NewErrorMsg(errs.ErrorStatus, fmt.Sprintf("sortLinks接口返回数据解析失败：%s, url:%s", err.Error(), urlStr))
	}
	return nil
}
