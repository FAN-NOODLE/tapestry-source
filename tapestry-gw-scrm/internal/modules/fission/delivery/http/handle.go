package fissionhttp

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/fission"
	"gw-scrm/internal/modules/fission/api"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/errorIgnore"
	"gw-scrm/internal/pkg/http"
	"gw-scrm/internal/pkg/sdk/microscrm"
	"gw-scrm/internal/pkg/validate_params"
	"strconv"
	"time"
)

func exportExchangeRecord(c *gin.Context) {

	var (
		req       fission.ExchangeDetailRequest
		err       error
		bytes     []byte
		excelData fission.ExchangeDetailResponse
		style     int
	)
	ok := validate_params.ValidateParams(c, &req)
	if !ok {
		return
	}
	tokenInfo := http.GetTokenInfo(c)
	req.EnterpriseID = tokenInfo.EnterpriseId
	//获取数据
	response, err := microscrm.ExportExchangeRecord(ctx.Wrap(c), &req)
	if errorIgnore.Ignore(err) != nil {
		ginplus.ResError(c, err)
		return
	}
	bytes, err = json.Marshal(response)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	err = json.Unmarshal(bytes, &excelData)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	excelFile := excelize.NewFile()
	sheetName := "裂变活动"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	//设置宽度
	err = excelFile.SetColWidth(sheetName, "A", "G", 40)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	//设置样式
	style, err = excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"},"font":{"bold":true}}`)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	err = excelFile.SetCellStyle(sheetName, "A1", "G1", style)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	// 写入表头
	_ = excelFile.SetCellStr(sheetName, "A1", "活动名称")
	_ = excelFile.SetCellStr(sheetName, "B1", "客户昵称")
	_ = excelFile.SetCellStr(sheetName, "C1", "兑换奖品")
	_ = excelFile.SetCellStr(sheetName, "D1", "兑换时间")
	_ = excelFile.SetCellStr(sheetName, "E1", "收货人")
	_ = excelFile.SetCellStr(sheetName, "F1", "收货人号码")
	_ = excelFile.SetCellStr(sheetName, "G1", "收货地址")
	startRow := 2
	//填充数据
	style, err = excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"}}`)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	for _, item := range excelData.List {
		rowStartStr := strconv.Itoa(startRow)
		_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.ActivityName)
		_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, item.Nickname)
		_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, item.PrizeName)
		//时间需要转为标准格式
		tm := time.Unix(item.ExchangeTime, 0)
		_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, tm.Format("2006-01-02 15:04:05"))
		_ = excelFile.SetCellStr(sheetName, "E"+rowStartStr, item.DeliverName)
		_ = excelFile.SetCellValue(sheetName, "F"+rowStartStr, item.DeliverMobile)
		address := item.DeliverArea + item.DeliverAddress
		_ = excelFile.SetCellStr(sheetName, "G"+rowStartStr, address)
		startRow = startRow + 1
		//设置样式
		err = excelFile.SetCellStyle(sheetName, "A"+rowStartStr, "G"+rowStartStr, style)
		if nil != err {
			ginplus.ResError(c, err)
			return
		}
	}
	excelFile.SetActiveSheet(index)
	fileName := "裂变活动" + time.Now().Format("200601021504") + ".xlsx"

	//本地生成 暂时保留
	/*	if err := excelFile.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}*/
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	des, _ := excelFile.WriteToBuffer()
	_, _ = c.Writer.Write(des.Bytes())
	return
}

//activityDetail裂变活动明细-统计信息
func activityDetail(c *gin.Context) {

	var (
		tokenInfo *common.TokenInfo
		act       = &api.FissionApi{}
		params    map[string]interface{}
		res       interface{}
		err       error
	)
	params = common.GetRawMap(c)
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	res, err = act.ActivityDetail(ctx.Wrap(c), params)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

//activityDetailList 裂变活动详情列表
func activityDetailList(c *gin.Context) {
	var (
		tokenInfo *common.TokenInfo
		act       = &api.FissionApi{}
		params    map[string]interface{}
		res       interface{}
		err       error
	)
	params = common.GetRawMap(c)
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	res, err = act.ActivityDetailList(ctx.Wrap(c), params)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return

}

//inviteDetail 发起人邀请明细
func inviteDetail(c *gin.Context) {
	var (
		tokenInfo *common.TokenInfo
		act       = &api.FissionApi{}
		params    map[string]interface{}
		res       interface{}
		err       error
	)
	params = common.GetRawMap(c)
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	res, err = act.InviteDetail(ctx.Wrap(c), params)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

//exchangeRecord 发起人兑换记录
func exchangeRecord(c *gin.Context) {
	var (
		tokenInfo *common.TokenInfo
		act       = &api.FissionApi{}
		params    map[string]interface{}
		res       interface{}
		err       error
	)
	params = common.GetRawMap(c)
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	res, err = act.ExchangeRecord(ctx.Wrap(c), params)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

//prizeExchangeRecord兑换奖品记录
func prizeExchangeRecord(c *gin.Context) {
	var (
		tokenInfo *common.TokenInfo
		act       = &api.FissionApi{}
		params    map[string]interface{}
		res       interface{}
		err       error
	)
	params = common.GetRawMap(c)
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	res, err = act.PrizeExchangeRecord(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		ginplus.ResError(c, err)
		return
	}
	bytes, err := json.Marshal(res)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	var excelData fission.ExchangeDetailResponse
	err = json.Unmarshal(bytes, &excelData)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}

	for index, item := range excelData.List {
		excelData.List[index].DeliverMobile = fmt.Sprintf("%s%s%s", item.DeliverMobile[0:3], "****", item.DeliverMobile[7:])
	}

	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, excelData)
	return
}

//OpenExchangeWarn 开启奖品兑换提醒
func openExchangeWarn(c *gin.Context) {
	var (
		tokenInfo *common.TokenInfo
		act       = &api.FissionApi{}
		params    map[string]interface{}
		res       interface{}
		err       error
	)
	params = common.GetRawMap(c)
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	res, err = act.OpenExchangeWarn(ctx.Wrap(c), params)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

//exchangeTotal 兑换记录总数
func exchangeTotal(c *gin.Context) {
	var (
		tokenInfo *common.TokenInfo
		act       = &api.FissionApi{}
		params    map[string]interface{}
		res       interface{}
		err       error
	)
	params = common.GetRawMap(c)
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	res, err = act.ExchangeTotal(ctx.Wrap(c), params)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

//activityListSimple 活动列表 全部有效活动
func activityListSimple(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId

	act := &api.FissionApi{}
	res, rErr := act.ActivityListSimple(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)

	return
}

//getExchangeConfigMember获取兑换提醒的额外员工
func getExchangeConfigMember(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.FissionApi{}
	res, rErr := act.GetExchangeConfigMember(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)

	return
}

// activityList 裂变活动列表
func activityList(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId

	act := &api.FissionApi{}
	res, rErr := act.ActivityList(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)

	return
}

// activityAdd 裂变活动新增
func activityAdd(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId

	act := &api.FissionApi{}
	res, rErr := act.ActivityAdd(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)

	return
}

// activityUpdate 裂变活动编辑
func activityUpdate(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId

	act := &api.FissionApi{}
	res, rErr := act.ActivityUpdate(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
	return
}

// activityDelete 裂变活动删除
func activityDelete(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId

	act := &api.FissionApi{}
	res, rErr := act.ActivityDelete(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
	return
}

// updateStatus 裂变活动更新状态
func updateStatus(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId

	act := &api.FissionApi{}
	res, rErr := act.UpdateStatus(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
	return
}

// activityInfo 裂变活动详情
func activityInfo(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId

	act := &api.FissionApi{}
	res, rErr := act.ActivityInfo(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
	return
}

// chatGroupList 选择社群列表
func chatGroupList(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId

	act := &api.FissionApi{}
	res, rErr := act.ChatGroupList(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
	return
}

// mutualActivityList 互斥活动选择列表
func mutualActivityList(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId

	act := &api.FissionApi{}
	res, rErr := act.MutualActivityList(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
	return
}

//getStyleConfig
func getStyleConfig(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.FissionApi{}
	res, rErr := act.GetStyleConfig(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)

	return
}

//setStyleConfig
func setStyleConfig(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.FissionApi{}
	res, rErr := act.SetStyleConfig(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)

	return
}

//addStyleConfig
func addStyleConfig(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.FissionApi{}
	res, rErr := act.AddStyleConfig(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)

	return
}

//delStyleConfig
func delStyleConfig(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.FissionApi{}
	res, rErr := act.DelStyleConfig(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)

	return
}
