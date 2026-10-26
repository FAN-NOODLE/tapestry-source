package cacheSvc

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/json"
	"micro-scrm/internal/app"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	microBossModel "micro-scrm/internal/app/models/micro_boss"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/redis"
	"time"
)

type RedisData struct {
	service.Base
}

func NewRedisData() *RedisData {
	return &RedisData{}
}

// GetEnterpriseConfig 获取企业配置信息
func (s *RedisData) GetEnterpriseConfig(c app.Ctx, enterpriseID int64, enterpriseConfigKey string) (microBossModel.EnterpriseConfigModel, error) {

	var (
		resConfig microBossModel.EnterpriseConfigModel
		err       error
	)
	configKey := fmt.Sprintf("%d", enterpriseID) + ":ENTERPRISE_CONFIG_KEY:" + fmt.Sprintf("%s", enterpriseConfigKey)
	resConfigCache := redis.Client().Get(c, configKey).Val()
	if resConfigCache != "" {
		jErr := json.Unmarshal([]byte(resConfigCache), &resConfig)
		if jErr != nil {
			return resConfig, jErr
		}

		return resConfig, nil
	} else {
		// 获取企业配置信息
		resConfig, err = microBossDao.NewEnterpriseConfig().GetInfo(c,
			map[string]interface{}{"enterprise_id": enterpriseID, "key": enterpriseConfigKey})
		if err != nil {
			return resConfig, err
		}
		if resConfig.ID == 0 {
			return resConfig, errors.New("未获取到企业配置信息")
		}

		// 缓存180分钟
		resConfigStr := json.MarshalToString(resConfig)
		redis.Client().Set(c, configKey, resConfigStr, 3*60*60*time.Second)
	}

	return resConfig, nil

}

// GetEnterpriseWorkPlugin 获取企业配置信息
func (s *RedisData) GetEnterpriseWorkPlugin(c app.Ctx, enterpriseID int64, identifier string) (microBossModel.EnterpriseWorkPluginModel, error) {
	configKey := fmt.Sprintf("%d", enterpriseID) + ":ENTERPRISE_WORK_PLUGIN_KEY:" + fmt.Sprintf("%s", identifier)
	resConfigCache := redis.Client().Get(c, configKey).Val()
	if resConfigCache != "" {
		var pluginInfo microBossModel.EnterpriseWorkPluginModel
		jErr := json.Unmarshal([]byte(resConfigCache), &pluginInfo)
		if jErr != nil {
			return pluginInfo, jErr
		}
		return pluginInfo, nil
	}

	// 获取企业插件表 小程序appid
	pluginInfo, err := microBossDao.NewEnterpriseWorkPlugin().GetOne(c, map[string]interface{}{
		"enterprise_id": enterpriseID,
		"identifier":    identifier,
	})

	if err != nil || pluginInfo.Secret == "" {
		return pluginInfo, errors.New("获取企业微信自建应用失败")
	}

	// 缓存180分钟
	resConfigStr := json.MarshalToString(pluginInfo)
	redis.Client().Set(c, configKey, resConfigStr, time.Minute*180)

	return pluginInfo, nil
}
