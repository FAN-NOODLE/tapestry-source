package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/data_compass"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

func dataCompassDataSummary(c *gin.Context) {
	var request data_compass.CommunityDataDetailExportRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	if tokenInfo.EnterpriseId == 0 {
		ginplus.ResError(c, errors.New("无效的企业编号"))
	}
	request.EnterpriseID = tokenInfo.EnterpriseId
	data, err := microscrm.GetDatacompassCommunityDataSummary(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, data)
}

func datacompassCommunityData(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	if tokenInfo.EnterpriseId == 0 {
		ginplus.ResError(c, errors.New("无效的企业编号"))
		return
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	result, err := microscrm.DatacompassCommunityData(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, result)
}

// 社群数据明细导出
func datacompassCommunityDataExport(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	if tokenInfo.EnterpriseId == 0 {
		ginplus.ResError(c, errors.New("无效的企业编号"))
		return
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	res, err := microscrm.DatacompassCommunityData(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	bytes, err := json.Marshal(res)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}

	data := data_compass.CommunityDate{}
	err = json.Unmarshal(bytes, &data)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	excelFile := excelize.NewFile()
	sheetName := "数据罗盘-社群数据明细"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	//设置宽度
	err = excelFile.SetColWidth(sheetName, "A", "P", 15)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	//设置样式
	style, styleErr := excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"},"font":{"bold":true}}`)
	if nil != styleErr {
		ginplus.ResError(c, styleErr)
		return
	}
	err = excelFile.SetCellStyle(sheetName, "A1", "P1", style)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	// 写入表头

	excelFile.SetCellStr(sheetName, "A1", "社群名称")
	excelFile.SetCellStr(sheetName, "B1", "累计群成员数")
	excelFile.SetCellStr(sheetName, "C1", "新入群人数")
	excelFile.SetCellStr(sheetName, "D1", "退群人数")
	excelFile.SetCellStr(sheetName, "E1", "退群率")
	excelFile.SetCellStr(sheetName, "F1", "BA人数")

	//填充数据
	style, err = excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"}}`)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	startRow := 2
	for _, item := range data.List {
		rowStartStr := strconv.Itoa(startRow)
		excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.GroupName)
		excelFile.SetCellStr(sheetName, "B"+rowStartStr, fmt.Sprintf("%d", item.CumulativeGroupNumber))
		excelFile.SetCellStr(sheetName, "C"+rowStartStr, fmt.Sprintf("%d", item.NumberOfNewcomers))

		excelFile.SetCellStr(sheetName, "D"+rowStartStr, fmt.Sprintf("%d", item.LeaveGroupNumber))
		excelFile.SetCellStr(sheetName, "E"+rowStartStr, item.LeaveGroupPercentage)

		excelFile.SetCellStr(sheetName, "F"+rowStartStr, fmt.Sprintf("%d", item.BANumber))

		startRow++
		err = excelFile.SetCellStyle(sheetName, "A"+rowStartStr, "P"+rowStartStr, style)
		if nil != err {
			ginplus.ResError(c, err)
			return
		}
	}
	excelFile.SetActiveSheet(index)
	fileName := sheetName + time.Now().Format("200601021504") + ".xlsx"
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	des, _ := excelFile.WriteToBuffer()
	c.Writer.Write(des.Bytes())
	return
}
func datacompassFriendData(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	if tokenInfo.EnterpriseId == 0 {
		ginplus.ResError(c, errors.New("无效的企业编号"))
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	result, err := microscrm.DatacompassFriendData(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, result)
}

// 社群数据明细导出
func datacompassFriendDataExport(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	if tokenInfo.EnterpriseId == 0 {
		ginplus.ResError(c, errors.New("无效的企业编号"))
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId

	res, err := microscrm.DatacompassFriendData(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, errors.New("无效的企业编号"))
	}
	bytes, err := json.Marshal(res)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	data := data_compass.BuddyData{}
	err = json.Unmarshal(bytes, &data)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	excelFile := excelize.NewFile()
	sheetName := "数据罗盘-好友数据明细"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	//设置宽度
	err = excelFile.SetColWidth(sheetName, "A", "F", 15)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	//设置样式
	style, styleErr := excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"},"font":{"bold":true}}`)
	if nil != styleErr {
		ginplus.ResError(c, styleErr)
		return
	}
	err = excelFile.SetCellStyle(sheetName, "A1", "F1", style)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	// 写入表头

	excelFile.SetCellStr(sheetName, "A1", "导购ID")
	excelFile.SetCellStr(sheetName, "B1", "导购")
	excelFile.SetCellStr(sheetName, "C1", "所属门店")
	excelFile.SetCellStr(sheetName, "D1", "新增好友数")
	excelFile.SetCellStr(sheetName, "E1", "好友流失率")
	excelFile.SetCellStr(sheetName, "F1", "累计好友数")

	//填充数据
	style, err = excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"}}`)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	startRow := 2
	for _, item := range data.List {
		rowStartStr := strconv.Itoa(startRow)
		excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.ShoppingGuideID)
		excelFile.SetCellStr(sheetName, "B"+rowStartStr, item.ShoppingGuideName)
		excelFile.SetCellStr(sheetName, "C"+rowStartStr, item.BelongingStore)
		excelFile.SetCellStr(sheetName, "D"+rowStartStr, fmt.Sprintf("%d", item.NewAddBuddyNumber))
		excelFile.SetCellStr(sheetName, "E"+rowStartStr, item.BuddyChurnPercentage)
		excelFile.SetCellStr(sheetName, "F"+rowStartStr, fmt.Sprintf("%d", item.GrandTotalNumber))
		startRow = startRow + 1
		err = excelFile.SetCellStyle(sheetName, "A"+rowStartStr, "F"+rowStartStr, style)
		if nil != err {
			ginplus.ResError(c, err)
			return
		}
	}
	excelFile.SetActiveSheet(index)
	fileName := sheetName + time.Now().Format("200601021504") + ".xlsx"
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	des, _ := excelFile.WriteToBuffer()
	c.Writer.Write(des.Bytes())
	return
}
