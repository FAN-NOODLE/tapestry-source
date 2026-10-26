package departments_http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/departments"
	"gw-scrm/internal/modules/departments/api"
	"gw-scrm/internal/pkg/common"
)

// departmentList 获取组织列表
func departmentList(c *gin.Context) {
	var requestInfo departments.DepartmentListRequest
	err := c.BindJSON(&requestInfo)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	requestInfo.EnterpriseId = tokenInfo.EnterpriseId
	act := &api.DepartmentApi{}
	res := act.DepartmentList(ctx.Wrap(c), &requestInfo)
	if res == nil {
		res = struct{}{}
	}
	ginplus.ResSuccess(c, res)
}

// departmentDetail 获取组织详情
func departmentDetail(c *gin.Context) {
	var requestInfo departments.DepartmentListRequest
	err := c.BindJSON(&requestInfo)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	act := &api.DepartmentApi{}
	res := act.DepartmentList(ctx.Wrap(c), &requestInfo)
	ginplus.ResSuccess(c, res)
}

// updateRecord 获取更新记录
func updateRecord(c *gin.Context) {
	var requestInfo departments.UpdateRecordRequest
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	requestInfo.EnterpriseId = tokenInfo.EnterpriseId
	act := &api.DepartmentApi{}
	res := act.UpdateRecord(ctx.Wrap(c), &requestInfo)
	ginplus.ResSuccess(c, res)
}

// departmentMomentList 选择部门
func departmentMomentList(c *gin.Context) {
	var requestInfo departments.MomentListRequest
	err := c.BindJSON(&requestInfo)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	requestInfo.EnterpriseId = tokenInfo.EnterpriseId
	act := &api.DepartmentApi{}
	res := act.DepartmentMomentList(ctx.Wrap(c), &requestInfo)
	ginplus.ResSuccess(c, res)
}

func userLostConfig(c *gin.Context) {
	var request departments.LostConfigRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	request.EnterpriseId = tokenInfo.EnterpriseId
	act := &api.DepartmentApi{}
	res := act.UserLostConfig(ctx.Wrap(c), &request)
	ginplus.ResSuccess(c, res)
}

func userLostInfo(c *gin.Context) {
	var request departments.LostInfoRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	request.EnterpriseId = tokenInfo.EnterpriseId
	act := &api.DepartmentApi{}
	res := act.UserLostInfo(ctx.Wrap(c), &request)
	ginplus.ResSuccess(c, res)
}
