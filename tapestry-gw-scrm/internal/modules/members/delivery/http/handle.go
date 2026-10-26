package members_http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/enterprise"
	"gw-scrm/internal/modules/members"
	"gw-scrm/internal/modules/members/api"
	"gw-scrm/internal/pkg/common"
	scrm_errs "gw-scrm/internal/pkg/errs"
	"gw-scrm/internal/pkg/sdk/microboss"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

// membersDetail 员工详情
func membersDetail(c *gin.Context) {
	var requestInfo members.MemberDetailRequest
	err := c.BindJSON(&requestInfo)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	requestInfo.EnterpriseID = tokenInfo.EnterpriseId
	act := &api.MemberApi{}
	res := act.MemberDetail(ctx.Wrap(c), &requestInfo)
	ginplus.ResSuccess(c, res)
}

// membersList 员工列表
func membersList(c *gin.Context) {
	var requestInfo members.MemberListRequest
	err := c.BindJSON(&requestInfo)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	requestInfo.EnterpriseID = tokenInfo.EnterpriseId
	act := &api.MemberApi{}
	res := act.MemberList(ctx.Wrap(c), &requestInfo)
	ginplus.ResSuccess(c, res)
}

// membersSync 同步通讯录
func membersSync(c *gin.Context) {
	type syncReq struct {
		Domain       string `json:"domain"`
		EnterpriseID int64  `json:"enterprise_id"`
		AdminID      int64  `json:"admin_id"`
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	// 查询企业信息
	enterpriseData := enterprise.Info(ctx.Wrap(c), tokenInfo.EnterpriseId)
	enterpriseDataRes := enterpriseData.(microboss.Enterprise)
	req := &syncReq{
		Domain:       enterpriseDataRes.Data.Domain,
		EnterpriseID: tokenInfo.EnterpriseId,
		AdminID:      tokenInfo.AdminId,
	}
	response, err := microscrm.SyncMemberContact(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// momentList 朋友圈创建请求员工
func momentList(c *gin.Context) {
	var requestInfo members.MomentListRequest
	err := c.BindJSON(&requestInfo)

	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	requestInfo.EnterpriseId = tokenInfo.EnterpriseId
	act := &api.MemberApi{}
	res := act.MomentList(ctx.Wrap(c), &requestInfo)
	ginplus.ResSuccess(c, res)
}

// 获取组织架构BA数据
func momentBaList(c *gin.Context) {
	var requestInfo members.MomentListRequest
	err := c.BindJSON(&requestInfo)

	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	requestInfo.EnterpriseId = tokenInfo.EnterpriseId
	act := &api.MemberApi{}
	res := act.MomentBaList(ctx.Wrap(c), &requestInfo)
	ginplus.ResSuccess(c, res)
}

// getSimpleMembersByDepartmentIds 获取多组织架构下员工信息
func getSimpleMembersByDepartmentIds(c *gin.Context) {
	res, err := microscrm.GetMembersByDepartmentIds(ctx.Wrap(c), common.GetBodyParams(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

func MemberListByThirdIds(c *gin.Context) {
	var requestInfo members.MemberListByThirdIdsReq
	err := c.BindJSON(&requestInfo)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	requestInfo.EnterpriseId = tokenInfo.EnterpriseId
	act := &api.MemberApi{}
	res := act.MemberListByThirdIds(ctx.Wrap(c), &requestInfo)
	ginplus.ResSuccess(c, res)
}
