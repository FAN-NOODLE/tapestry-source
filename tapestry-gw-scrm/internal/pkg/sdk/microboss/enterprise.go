package microboss

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

// 企业积分规则type
const ScoreType = 10

// 企业积分规则key
const ScoreKey = "score_show_name"

// 企业积分规则name
const ScoreName = "企业展示积分名称"

// Enterprise 企业信息构造体
type Enterprise struct {
	NoCheckResponse
	Data struct {
		ID          int64  `json:"id"`
		Hash        string `json:"hash"`
		Name        string `json:"name"`
		CorpId      string `json:"corp_id"`
		Domain      string `json:"domain"`
		Version     int64  `json:"version"`
		IsSplitFlow int64  `json:"is_split_flow"`
	}
}

// 保存企业配置信息请求参数
type ConfigInfoRequest struct {
	EnterpriseID int64    `json:"enterprise_id"`
	AdminId      int64    `json:"admin_id"`
	Name         string   `json:"name"`
	Data         []Config `json:"data"`
}

// 企业配置信息
type Config struct {
	Type  int64  `json:"type"`
	Key   string `json:"key"`
	Name  string `json:"name"`
	Value int64  `json:"value"`
	Desc  string `json:"desc"`
}

// ConfigInfo 企业配置信息构造体
type ConfigInfo struct {
	NoCheckResponse
	Data struct {
		Type  int64  `json:"type"`
		Key   string `json:"key"`
		Value int64  `json:"value"`
		Desc  string `json:"desc"`
		Name  string `json:"name"`
	}
}

type Response interface{}

// EnterpriseInfo 获取企业信息
func EnterpriseInfo(c ctx.Context, request interface{}) (Enterprise, error) {
	var response Enterprise
	httpClient.Post(c, GetUrl(urlEnum.EnterpriseInfoUrl), request, &response)
	return response, nil
}

// GetScoreName 获取企业配置-积分名称
func GetScoreName(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.EnterpriseInfoGetUrl), request)
	return response, err
}

// SaveScoreName 设置企业配置-积分名称
func SaveScoreName(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.EnterpriseInfoSetUrl), request)
	return response, err
}

// GetEnterpriseName 获取企业配置-积分名称
func GetEnterpriseName(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.EnterpriseInfoGetUrl), request)

	return response, err
}