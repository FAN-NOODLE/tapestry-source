package common

import (
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/http"
	"micro-scrm/internal/pkg/redis"
	"micro-scrm/internal/pkg/sdk/boss"
	"time"
)

// EnterpriseExpire 企业信息缓存时间
const EnterpriseExpiration = 2 * time.Hour

// EnterpriseID 从上下文中获取企业id
func EnterpriseID(c ctx.Context) int64 {
	var enterpriseID int64
	if c.Value(ContextKeyEnum.EnterpriseID) != nil {
		switch v := c.Value(ContextKeyEnum.EnterpriseID).(type) {
		case int64:
			enterpriseID = v
		default:
			enterpriseID = 0
		}
	}
	return enterpriseID
}

// 生成短参数
func CreateShortParamId(c ctx.Context, params Params) (int64, error) {

	var inner Inner
	var innerArr []Inner
	for key, value := range params {
		inner.Name = key
		inner.Value = value
		innerArr = append(innerArr, inner)
	}
	requestParams := requestParams{
		Args: innerArr,
	}

	url := fmt.Sprintf("%s0", config.CfgData.Service.SortLinks.Domain)
	res, errP := http.ClientObject.Post(c, url, clients.JSONBody(requestParams))
	if errP != nil {
		return 0, errP
	}
	rt, errR := res.Read(clients.Text)
	if errR != nil {
		return 0, errR
	}
	var response = &ResponseParams{}
	errJ := json.Unmarshal([]byte(rt.(string)), response)
	if errJ != nil {
		return 0, errJ
	}

	return response.ID, nil
}

// EnterpriseID 从上下文中获取企业id
func EnterpriseHash(c ctx.Context) string {
	var enterpriseHash string
	if c.Value(ContextKeyEnum.EnterpriseHash) != nil {
		switch v := c.Value(ContextKeyEnum.EnterpriseHash).(type) {
		case string:
			enterpriseHash = v
		default:
			enterpriseHash = ""
		}
	}
	return enterpriseHash
}

// GetEnterpriseInfo 获取企业信息
func GetEnterpriseInfo(c ctx.Context, enterpriseID int64) (*EnterpriseInfo, error) {
	var enterprise EnterpriseInfo
	key := fmt.Sprintf("%s:enterprise:info:%d", config.CfgData.App, enterpriseID)
	gErr := redis.Client().Get(c, key).Scan(&enterprise)
	if gErr == redis.Nil {
		log.Info(c).Msgf("redis中不存在企业信息， enterpriseID:%d", enterpriseID)
	} else if gErr != nil {
		log.Info(c).Msgf("redis查询企业信息失败：enterpriseID:%d, err:%s", enterpriseID, gErr.Error())
	} else {
		return &enterprise, nil
	}

	var request = boss.EnterpriseDetailRequest{
		EnterpriseID: enterpriseID,
	}
	response, err := boss.EnterpriseDetail(c, &request)
	if err != nil {
		return nil, err
	}
	log.Info(c).Msgf("请求企业信息结果：%+v", response)
	enterprise = EnterpriseInfo{
		ID:             response.ID,
		EnterpriseID:   enterpriseID,
		Hash:           response.Hash,
		Name:           response.Name,
		Logo:           response.Logo,
		IndustryID:     response.IndustryID,
		CorpID:         response.CorpID,
		Domain:         response.Domain,
		Status:         response.Status,
		BusinessSource: response.BusinessSource,
		Version:        response.Version,
		VersionTag:     response.VersionTag,
		IsSplitFlow:    response.IsSplitFlow,
		Identifier:     response.Identifier,
	}
	err = redis.Client().Set(c, key, &enterprise, EnterpriseExpiration).Err()
	if err != nil {
		log.Info(c).Msgf("redis设置缓存信息失败, enterpriseID:%d, err:%s", enterpriseID, err.Error())
	}
	return &enterprise, nil
}
