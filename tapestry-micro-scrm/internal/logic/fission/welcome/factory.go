package welcome

import (
	"encoding/json"
	"fmt"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/wx/wxoffiaccount/cache"
	"micro-scrm/internal/config"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm/fission"
	micro_redis "micro-scrm/internal/pkg/redis"
)

const (
	WelcomeType_Fission_AddEmp          = 1 + iota // 加BA
	WelcomeType_Fission_AddEmpAndGroup             // 加BA且入群
	WelcomeType_Fission_AddEmpAndMember            // 加BA且成为品牌会员
)

type welcomeGetter interface {
	Get() ([]map[string]interface{}, error)
}

type welcomeFactory struct {
	ctx ctx.Context
	Param
	data  fission.FissionActivityModel
	cache *cache.RedisCache
}

// 工厂参数
type Param struct {
	EnterpriseID int64 // 企业ID
	ActivityID   int64 // 欢迎语ID
}

func NewWelcomeFactory(ctx ctx.Context, param Param) (wf *welcomeFactory, err error) {
	wf = &welcomeFactory{
		ctx:   ctx,
		Param: param,
		cache: cache.NewRedisCache(micro_redis.Client()),
	}

	if wf.EnterpriseID == 0 || wf.ActivityID == 0 {
		return wf, welcomeParamError{}
	}

	// 获取活动信息
	key := fmt.Sprintf("%s::%d::fission_activity_detail::%d", config.CfgData.App, wf.EnterpriseID, wf.ActivityID)
	var data string
	if err = wf.cache.Get(wf.ctx, key, &data); err != nil {
		// 获取活动信息
		wf.data, err = orm.Service.Fission.GetFissionDetail(wf.EnterpriseID, wf.ActivityID)
		if err != nil {
			return
		}

		jsonData, jsonErr := json.Marshal(wf.data)
		if jsonErr != nil {
			return
		}
		expireIn := wf.data.EndTime - wf.data.StartTime
		wf.cache.Set(wf.ctx, key, jsonData, int64(expireIn))
		data = string(jsonData)
	}

	if err = json.Unmarshal([]byte(data), &wf.data); err != nil {
		return
	}

	return
}

// 获取对应类型活动的欢迎语
func (w welcomeFactory) Get() ([]map[string]interface{}, error) {
	// 这个类型，后期可能不会遵循fission_activity.fission_type。 所以，重新构建了const
	var iWelcome welcomeGetter

	switch w.data.FissionType {
	case WelcomeType_Fission_AddEmp:
		iWelcome = NewWelcomeFissionAndEmp(w.ctx, w.data)
	case WelcomeType_Fission_AddEmpAndGroup:
		iWelcome = NewWelcomeFissionAndEmpGroup(w.ctx, w.data)
	case WelcomeType_Fission_AddEmpAndMember:
		// 裂变活动 - 客户添加BA且成为品牌会员
		iWelcome = NewWelcomeFissionAndEmpMember(w.ctx, w.data)
	default:
		return nil, welcomeObjError{}
	}
	return iWelcome.Get()
}

// drainage_welcomes.content 转结构体
func DecodeWelcomeContent(contentJson string) (msg []map[string]interface{}, err error) {
	err = json.Unmarshal([]byte(contentJson), &msg)

	return
}
