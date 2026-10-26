package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/socialDraw/api"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

// DefaultApi 默认初始化社群画家接口
var DefaultApi = &api.SocialDrawApi{}

// getBodyParams 获取Body参数
func getBodyParams(c *gin.Context) interface{} {
	params := common.GetRawMap(c)
	tk := c.Value(common.ContextKeyEnum.TokenInfo)
	if tk != nil {
		tokenInfo := tk.(*common.TokenInfo)
		params["enterprise_id"] = tokenInfo.EnterpriseId
	}
	return params
}

// getWorkGroupChatList 获取社群列表
func getWorkGroupChatList(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	tk := c.Value(common.ContextKeyEnum.TokenInfo)
	if tk != nil {
		tokenInfo = tk.(*common.TokenInfo)
		params["enterprise_id"] = tokenInfo.EnterpriseId
	}

	res, err := DefaultApi.GetWorkGroupChatList(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

// getWorkGroupChatInfo 获取社群详情
func getWorkGroupChatInfo(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	tk := c.Value(common.ContextKeyEnum.TokenInfo)
	if tk != nil {
		tokenInfo = tk.(*common.TokenInfo)
		params["enterprise_id"] = tokenInfo.EnterpriseId
	}

	res, err := DefaultApi.GetWorkGroupChatInfo(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

// setWorkGroupChat 社群设置
func setWorkGroupChat(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	tk := c.Value(common.ContextKeyEnum.TokenInfo)
	if tk != nil {
		tokenInfo = tk.(*common.TokenInfo)
		params["enterprise_id"] = tokenInfo.EnterpriseId
	}

	res, err := DefaultApi.SetWorkGroupChat(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

// getWorkChatGroupTypeList 社群类型列表
func getWorkGroupChatTypeList(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	tk := c.Value(common.ContextKeyEnum.TokenInfo)
	if tk != nil {
		tokenInfo = tk.(*common.TokenInfo)
		params["enterprise_id"] = tokenInfo.EnterpriseId
	}

	res, err := DefaultApi.GetWorkGroupChatTypeList(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

// manageWorkGroupChatType 社群类型管理
func manageWorkGroupChatType(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	tk := c.Value(common.ContextKeyEnum.TokenInfo)
	if tk != nil {
		tokenInfo = tk.(*common.TokenInfo)
		params["enterprise_id"] = tokenInfo.EnterpriseId
	}

	res, err := DefaultApi.ManageWorkGroupChatType(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

// getWorkGroupChatFilters 社群类型管理
func getWorkGroupChatFilters(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	tk := c.Value(common.ContextKeyEnum.TokenInfo)
	if tk != nil {
		tokenInfo = tk.(*common.TokenInfo)
		params["enterprise_id"] = tokenInfo.EnterpriseId
	}

	res, err := DefaultApi.GetWorkGroupChatFilters(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

// getWorkGroupChatMembers 社群成员明细
func getWorkGroupChatMembers(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	tk := c.Value(common.ContextKeyEnum.TokenInfo)
	if tk != nil {
		tokenInfo = tk.(*common.TokenInfo)
		params["enterprise_id"] = tokenInfo.EnterpriseId
	}

	res, err := DefaultApi.GetWorkGroupChatMembers(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

// getWorkGroupChatDataGrid 获取社群详情统计数据
func getWorkGroupChatDataGrid(c *gin.Context) {
	res, err := microscrm.WorkGroupChatDataGrid(ctx.Wrap(c), getBodyParams(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

// WorkGroupChatUploadSetting 上传社群批量设置数据
func WorkGroupChatUploadSetting(c *gin.Context) {
	var params = make(map[string]string)
	params["enterprise_id"] = c.PostForm("enterprise_id")

	tk := c.Value(common.ContextKeyEnum.TokenInfo)
	if tk != nil {
		tokenInfo := tk.(*common.TokenInfo)
		params["enterprise_id"] = fmt.Sprintf("%+v", tokenInfo.EnterpriseId)
	}

	fh, err := c.FormFile("file")
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	res, err := microscrm.WorkGroupUploadSetting(ctx.Wrap(c), params, "file", fh)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

// WorkGroupChatBatchSetting 批量设置社群数据
func WorkGroupChatBatchSetting(c *gin.Context) {
	res, err := microscrm.WorkGroupBatchSetting(ctx.Wrap(c), getBodyParams(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

// WorkGroupChatDataUpdateTime 数据最新更新时间
func WorkGroupChatDataUpdateTime(c *gin.Context) {
	res, err := microscrm.WorkGroupDataUpdateTime(ctx.Wrap(c), getBodyParams(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

// WorkGroupChatBatchSync 批量同步群信息
func WorkGroupChatBatchSync(c *gin.Context) {
	res, err := microscrm.WorkGroupBatchSync(ctx.Wrap(c), getBodyParams(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}
