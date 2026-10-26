package right_rule_http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/right_rule"
	"gw-scrm/internal/modules/right_rule/api"
	"gw-scrm/internal/pkg/common"
	scrm_errs "gw-scrm/internal/pkg/errs"
	response "gw-scrm/internal/pkg/resoponse"
	"gw-scrm/internal/pkg/sdk/microboss"
	"strings"
)

// GetScoreName 获取企业积分名称
func GetScoreName(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	params["enterprise_id"] = tokenInfo.EnterpriseId
	params["admin_id"] = tokenInfo.AdminId
	params["type"] = microboss.ScoreType
	params["key"] = microboss.ScoreKey

	response, err := microboss.GetScoreName(ctx.Wrap(c), params)
	if err != nil {
		response = microboss.Config{
			Name: "",
		}
	}
	ginplus.ResSuccess(c, response)
}

// SaveScoreName 设置企业展示积分名称
func SaveScoreName(c *gin.Context) {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	var request microboss.ConfigInfoRequest
	request.EnterpriseID = tokenInfo.EnterpriseId
	request.AdminId = tokenInfo.AdminId
	c.BindJSON(&request)

	// 过滤空格及验证为空
	name := strings.Replace(request.Name, " ", "", -1)
	if name == "" {
		response.Notice(c, scrm_errs.InputEmpty, nil)
		return
	}

	var config = microboss.Config{
		Type:  microboss.ScoreType,
		Key:   microboss.ScoreKey,
		Name:  name,
		Value: 0,
		Desc:  microboss.ScoreName,
	}

	request.Data = append(request.Data, config)

	response, err := microboss.SaveScoreName(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// scoreRuleSave 积分权益新增保存
func scoreRuleSave(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	// 获取管理员姓名
	var AdminDetailRequest = right_rule.AdminDetailRequest{
		AdminId: tokenInfo.AdminId,
	}

	AdminDetailRes, AdminDetailErr := microboss.GetAdminDetail(ctx.New(), AdminDetailRequest)
	if AdminDetailErr != nil {
		ginplus.ResError(c, AdminDetailErr)
		return
	}

	params["enterprise_id"] = tokenInfo.EnterpriseId
	params["operation_id"] = tokenInfo.AdminId
	params["operation_name"] = AdminDetailRes.Data.Username
	act := &api.RightRuleApi{}
	res, rErr := act.ScoreRuleSave(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

// scoreDetail 积分权益规则详情
func scoreDetail(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	params["operation_id"] = tokenInfo.AdminId
	act := &api.RightRuleApi{}
	res, rErr := act.ScoreDetail(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

// setUse 权益规则-启用/禁用
func setUse(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	params["operation_id"] = tokenInfo.AdminId
	act := &api.RightRuleApi{}
	res, rErr := act.SetUse(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

// ruleDelete 权益规则-删除
func ruleDelete(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	params["operation_id"] = tokenInfo.AdminId
	act := &api.RightRuleApi{}
	res, rErr := act.RuleDelete(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

// saveCouponRule 新增编辑优惠券权益规则
func saveCouponRule(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	// 获取管理员姓名
	var AdminDetailRequest = right_rule.AdminDetailRequest{
		AdminId: tokenInfo.AdminId,
	}

	AdminDetailRes, AdminDetailErr := microboss.GetAdminDetail(ctx.New(), AdminDetailRequest)
	if AdminDetailErr != nil {
		ginplus.ResError(c, AdminDetailErr)
		return
	}

	params["enterprise_id"] = tokenInfo.EnterpriseId
	params["operation_id"] = tokenInfo.AdminId
	params["operation_name"] = AdminDetailRes.Data.Username
	act := &api.RightRuleApi{}

	res, rErr := act.SaveCouponRule(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

// wxCouponList 优惠券列表
func wxCouponList(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.RightRuleApi{}
	res, rErr := act.WxCouponList(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

// infoCouponRule 优惠券权益规则详情
func infoCouponRule(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.RightRuleApi{}
	res, rErr := act.InfoCouponRule(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

// infoCouponRule 优惠券权益规则详情
func couponPackageList(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.RightRuleApi{}
	res, rErr := act.CouponPackageList(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

// list 权益规则列表
func list(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.RightRuleApi{}
	res, rErr := act.RuleList(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}
