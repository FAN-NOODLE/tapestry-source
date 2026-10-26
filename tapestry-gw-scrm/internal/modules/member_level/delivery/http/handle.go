package member_level_http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/member_level/api"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/gomicrouser"
)

// List 会员等级列表
func List(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.MemberLevelApi{}
	res, rErr := act.LevelList(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

// Detail 会员等级详情
func Detail(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.MemberLevelApi{}
	res, rErr := act.LevelDetail(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

// Set 会员等级设置
func Set(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.MemberLevelApi{}
	res, rErr := act.LevelSet(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

// UpdateStatus 会员等级状态修改
func UpdateStatus(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.MemberLevelApi{}
	res, rErr := act.UpdateStatus(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

// ContentSort 等级内容排序
func ContentSort(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.MemberLevelApi{}
	res, rErr := act.ContentSort(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

// LevelUpList 等级跃迁列表
func LevelUpList(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)
	response, err := gomicrouser.LevelUpList(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// AddLevelUp 新增等级跃迁
func AddLevelUp(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)
	response, err := gomicrouser.AddLevelUp(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// EditLevelUp 编辑等级跃迁
func EditLevelUp(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)
	response, err := gomicrouser.EditLevelUp(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// DeleteLevelUp 删除等级跃迁
func DeleteLevelUp(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)
	response, err := gomicrouser.DeleteLevelUp(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// LevelUpInfo 删除等级跃迁
func LevelUpInfo(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)
	response, err := gomicrouser.LevelUpInfo(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// LevelList 会员等级列表
func LevelList(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.MemberLevelApi{}
	res, rErr := act.SelectLevelList(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}
