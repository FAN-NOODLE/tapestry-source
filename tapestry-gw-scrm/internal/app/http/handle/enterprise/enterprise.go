package enterprise

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/app/http/handle"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microboss"
)

const enterpriseNameCoach = "coach"
const enterpriseNameKey = "store_name"

type EnterpriseHandle struct {
	handle.BaseHandle
}

func NewEnterpriseHandle() EnterpriseHandle {
	return EnterpriseHandle{
	}
}

func (s EnterpriseHandle) EnterpriseConfig(c *gin.Context) {

	enterpriseID, _ := common.GetEnterprise(c)

	// 获取企业配置
	params := common.GetRawMap(c)
	params["enterprise_id"] = enterpriseID
	params["key"] = enterpriseNameKey
	enterpriseConfigs, err := microboss.GetEnterpriseName(ctx.Wrap(c), params)
	var enterpriseConfigArr []microboss.Config
	var enterpriseConfig microboss.Config
	if err != nil {
		enterpriseConfig = microboss.Config{
		}
	}

	result , _ :=json.Marshal(enterpriseConfigs)
	err = json.Unmarshal(result, &enterpriseConfigArr)
	if err != nil {
		enterpriseConfig.Name = ""
	}
	for _,v := range enterpriseConfigArr {
		enterpriseConfig.Name = v.Name
		enterpriseConfig.Key = v.Key
	}

	s.ResSuccess(c, enterpriseConfig)
	return
}