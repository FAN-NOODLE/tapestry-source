package common

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gw-work-scrm/internal/config"
	"gw-work-scrm/internal/pkg/redis"
	"gw-work-scrm/internal/pkg/sdk/boss"
	"time"
)

// EnterpriseExpire 企业信息缓存时间
const EnterpriseExpiration = 2 * time.Hour

// EnterpriseID 从上下文中获取企业id
func EnterpriseID(c ctx.Context) int64 {
	var enterpriseID int64
	if c.Value(ContextKeyEnum.EnterpriseID) != nil {
		switch c.Value(ContextKeyEnum.EnterpriseID).(type) {
		case int64:
			enterpriseID = c.Value(ContextKeyEnum.EnterpriseID).(int64)
		default:
			enterpriseID = 0
		}
	}
	return enterpriseID
}

// GetEnterpriseInfo 获取企业信息
func GetEnterpriseInfo(c ctx.Context, hash string) (*EnterpriseInfo, error) {
	var enterprise EnterpriseInfo
	key := fmt.Sprintf("%s:enterprise:info:%s", config.CfgData.App, hash)
	gErr := redis.Client().Get(c, key).Scan(&enterprise)
	if gErr == redis.Nil {
		log.Info(c).Msgf("redis中不存在企业信息， hash:%s", hash)
	} else if gErr != nil {
		log.Info(c).Msgf("redis查询企业信息失败：hash:%s, err:%s", hash, gErr.Error())
	} else {
		return &enterprise, nil
	}

	var request = boss.EnterpriseDetailRequest{
		Hash: hash,
	}
	response, err := boss.EnterpriseDetail(c, &request)
	if err != nil {
		return nil, err
	}
	log.Info(c).Msgf("请求企业信息结果：%+v", response)
	enterprise = EnterpriseInfo{
		ID:             response.ID,
		EnterpriseID:   response.ID,
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
		BrandCode:      response.BrandCode,
	}
	err = redis.Client().Set(c, key, &enterprise, EnterpriseExpiration).Err()
	if err != nil {
		log.Info(c).Msgf("redis设置缓存信息失败, hash:%s, err:%s", hash, err.Error())
	}
	return &enterprise, nil
}
