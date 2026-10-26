package drainage_http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/drainage"
	"gw-scrm/internal/pkg/common"
	scrm_errs "gw-scrm/internal/pkg/errs"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

// channelCreate 渠道添加
func channelsCreate(c *gin.Context) {

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	var request map[string]interface{}
	c.BindJSON(&request)
	request["enterprise_id"] = tokenInfo.EnterpriseId
	request["admin_id"] = tokenInfo.AdminId

	response, err := microscrm.ChannelsCreate(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// channelsList 渠道列表
func channelsList(c *gin.Context) {

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	var request map[string]interface{}
	c.BindJSON(&request)
	request["enterprise_id"] = tokenInfo.EnterpriseId
	request["admin_id"] = tokenInfo.AdminId

	response, err := microscrm.ChannelsList(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// channelsDetail 渠道详情
func channelsDetail(c *gin.Context) {

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	var request map[string]interface{}
	c.BindJSON(&request)
	request["enterprise_id"] = tokenInfo.EnterpriseId
	request["admin_id"] = tokenInfo.AdminId

	response, err := microscrm.ChannelsDetail(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// channelCreate 渠道编辑
func channelsEdit(c *gin.Context) {

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	var request map[string]interface{}
	c.BindJSON(&request)
	request["enterprise_id"] = tokenInfo.EnterpriseId
	request["admin_id"] = tokenInfo.AdminId

	response, err := microscrm.ChannelsEdit(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// welcomeCreate 欢迎语新建
func welcomeCreate(c *gin.Context) {
	params := &drainage.WelcomeCreateRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseId = tokenInfo.EnterpriseId
	response, err := microscrm.WelcomeCreate(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

// welcomeDetail 欢迎语详情
func welcomeDetail(c *gin.Context) {
	params := &drainage.WelcomeDetailResponse{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseId = tokenInfo.EnterpriseId
	response, err := microscrm.WelcomeDetail(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

// welcomeDefault 欢迎语默认详情
func welcomeDefault(c *gin.Context) {
	params := &drainage.WelcomeDetailResponse{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseId = tokenInfo.EnterpriseId
	response, err := microscrm.WelcomeDefault(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

// welcomeEdit 欢迎语编辑
func welcomeEdit(c *gin.Context) {
	params := &drainage.WelcomeEditRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseId = tokenInfo.EnterpriseId
	response, err := microscrm.WelcomeEdit(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

// welcomeDelete 欢迎语删除
func welcomeDelete(c *gin.Context) {
	params := &drainage.WelcomeDeleteRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	response, err := microscrm.WelcomeDelete(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

// welcomeList 欢迎语列表
func welcomeList(c *gin.Context) {
	params := &drainage.WelcomeListRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseId = tokenInfo.EnterpriseId
	response, err := microscrm.WelcomeList(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

// welcomeValidate 欢迎语校验
func welcomeChannelList(c *gin.Context) {
	params := &drainage.WelcomeValidateRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseId = tokenInfo.EnterpriseId
	response, err := microscrm.WelcomeChannelList(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

// channelsData 渠道数据统计
func channelsData(c *gin.Context) {
	params := &drainage.ChannelsDataRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseID = tokenInfo.EnterpriseId

	response, err := microscrm.StatisticsChannelFacades(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

func KdTypeAll(c *gin.Context) {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	var param drainage.KdTypaAll
	param.EnterpriseID = tokenInfo.EnterpriseId
	response, err := microscrm.GetAllKd(ctx.Wrap(c), param)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

func WelcomeDetailNew(c *gin.Context) {
	params := &drainage.WelcomeBase{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseID = tokenInfo.EnterpriseId

	response, err := microscrm.WelcomeDetailNew(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

func WelcomeNewEdit(c *gin.Context) {
	params := &drainage.WelcomeNewEditRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseId = tokenInfo.EnterpriseId

	response, err := microscrm.WelcomeNewEdit(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

func WelcomeDel(c *gin.Context) {
	params := &drainage.WelcomeBase{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseID = tokenInfo.EnterpriseId

	response, err := microscrm.DelWelcome(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

func WelcomeAll(c *gin.Context) {
	params := &drainage.WelcomeAllRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseID = tokenInfo.EnterpriseId

	response, err := microscrm.GetAllWelcome(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

// channelsDataSum 渠道数据统计
func channelsDataSum(c *gin.Context) {
	params := &drainage.ChannelsDataRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseID = tokenInfo.EnterpriseId

	ctx := ctx.Wrap(c)
	if params.ChannelsId == 0 {
		// 查询所有渠道
		response, err := microscrm.StatisticsAllBySum(ctx, params)
		if err != nil {
			ginplus.ResError(c, err, http.StatusOK)
			return
		}
		ginplus.ResSuccess(c, response)

		return
	}

	// 查询某个渠道的总览
	params.Type = 2
	response, err := microscrm.StatisticsChannelFacades(ctx, params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	data, ok := response.(*drainage.ChannelsCount)
	if !ok || len(data.List) != 1 {
		ginplus.ResError(c, errorx.NewErrorMsg(5000, "响应结果出错"), http.StatusOK)
		return
	}

	ginplus.ResSuccess(c, data.List[0])
}

func channelsCountExport(c *gin.Context) {
	params := &drainage.ChannelsDataRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseID = tokenInfo.EnterpriseId
	params.Page = 1
	params.Size = 9999
	response, err := microscrm.StatisticsChannelFacades(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	// 类型断言
	data, ok := response.(*drainage.ChannelsCount)
	if !ok {
		ginplus.ResError(c, errorx.NewErrorMsg(5000, "响应结果出错"), http.StatusOK)
		return
	}

	excelFile := excelize.NewFile()
	sheetName := "渠道统计"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	if params.Type == 1 {
		_ = excelFile.SetCellStr(sheetName, "A1", "日期")
	} else {
		_ = excelFile.SetCellStr(sheetName, "A1", "渠道")
	}

	// 写入表头
	_ = excelFile.SetCellStr(sheetName, "B1", "留存好友数")
	_ = excelFile.SetCellStr(sheetName, "C1", "新增好友数")
	_ = excelFile.SetCellStr(sheetName, "D1", "净增好友数")
	_ = excelFile.SetCellStr(sheetName, "E1", "流失好友数")
	_ = excelFile.SetCellStr(sheetName, "F1", "好友删除员工数")
	_ = excelFile.SetCellStr(sheetName, "G1", "员工删除好友数")
	rowStart := 2
	if data.List != nil {
		for _, channelData := range data.List {
			rowStartStr := strconv.Itoa(rowStart)
			excelFile.SetCellValue(sheetName, "A"+rowStartStr, channelData.TimeOrName)
			excelFile.SetCellValue(sheetName, "B"+rowStartStr, channelData.ExistRelationCustomer)
			excelFile.SetCellValue(sheetName, "C"+rowStartStr, channelData.IncreaseCustomerCount)
			excelFile.SetCellValue(sheetName, "D"+rowStartStr, channelData.NewlyAddedCustomerCount)
			excelFile.SetCellValue(sheetName, "E"+rowStartStr, channelData.LossCustomerCount)
			excelFile.SetCellValue(sheetName, "F"+rowStartStr, channelData.CustomerDelMember)
			excelFile.SetCellValue(sheetName, "G"+rowStartStr, channelData.MemberDelCustomer)
			rowStart = rowStart + 1
		}
	}
	excelFile.SetActiveSheet(index)
	fileName := "渠道统计" + time.Now().Format("200601021504") + ".xlsx"
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	des, _ := excelFile.WriteToBuffer()
	_, _ = c.Writer.Write(des.Bytes())
	return

}

// cacheMedia 上传图片素材
func cacheMedia(c *gin.Context) {
	params := &drainage.CacheMediaRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseId = tokenInfo.EnterpriseId
	response, err := microscrm.CacheMedia(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

// channelSimpleList 渠道列表
func channelSimpleList(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	tk := c.Value(common.ContextKeyEnum.TokenInfo)
	if tk != nil {
		tokenInfo = tk.(*common.TokenInfo)
		params["enterprise_id"] = tokenInfo.EnterpriseId
	}
	res, err := microscrm.GetSimpleChannelList(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

func externalWelcomeTypeList(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	tk := c.Value(common.ContextKeyEnum.TokenInfo)
	if tk != nil {
		tokenInfo = tk.(*common.TokenInfo)
		params["enterprise_id"] = tokenInfo.EnterpriseId
	}
	res, err := microscrm.ExternalWelcomeTypeList(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}
