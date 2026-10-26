package http

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/ba_kanban"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type baKanBan struct {
}

func NewBaKanBan() *baKanBan {
	return &baKanBan{}
}

func (ba *baKanBan) GetBaKanBanData(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId

	res, rErr := microscrm.ApiBaKanBan{}.GetBaKanBanData(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

func (ba *baKanBan) GetSalesDataDetail(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId

	res, rErr := microscrm.ApiBaKanBan{}.GetSalesDataDetail(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

func (ba *baKanBan) ExportSalesDataDetail(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId

	res, rErr := microscrm.ApiBaKanBan{}.GetSalesDataDetail(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}

	bytes, err := json.Marshal(res)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}

	var data ba_kanban.OperationDataDetailResponseV1
	err = json.Unmarshal(bytes, &data)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}

	excelFile := excelize.NewFile()
	sheetName := "BA销售数据明细"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	//设置宽度
	err = excelFile.SetColWidth(sheetName, "A", "P", 40)
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
	_ = excelFile.SetCellStr(sheetName, "A1", "导购")
	_ = excelFile.SetCellStr(sheetName, "B1", "销售额")
	_ = excelFile.SetCellStr(sheetName, "C1", "目标金额")
	_ = excelFile.SetCellStr(sheetName, "D1", "完成率")
	_ = excelFile.SetCellStr(sheetName, "E1", "差额")
	_ = excelFile.SetCellStr(sheetName, "F1", "订单数")
	_ = excelFile.SetCellStr(sheetName, "G1", "客人数")
	_ = excelFile.SetCellStr(sheetName, "H1", "新客人数")
	_ = excelFile.SetCellStr(sheetName, "I1", "新客销售额")
	_ = excelFile.SetCellStr(sheetName, "J1", "新客销售占比")
	_ = excelFile.SetCellStr(sheetName, "K1", "新二购人数")
	_ = excelFile.SetCellStr(sheetName, "L1", "新二购销售额")
	_ = excelFile.SetCellStr(sheetName, "M1", "新二购销售占比")
	_ = excelFile.SetCellStr(sheetName, "N1", "老客人数")
	_ = excelFile.SetCellStr(sheetName, "O1", "老客销售额")
	_ = excelFile.SetCellStr(sheetName, "P1", "老客销售占比")
	startRow := 2
	//填充数据
	style, err = excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"}}`)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	for _, item := range data.List {
		rowStartStr := strconv.Itoa(startRow)
		_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.BaName)
		_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, fmt.Sprintf("%.0f", item.Sales))
		_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, fmt.Sprintf("%.0f", item.MonthlySales))
		_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, fmt.Sprintf("%d", item.FinishRate)+"%")
		_ = excelFile.SetCellStr(sheetName, "E"+rowStartStr, fmt.Sprintf("%.0f", item.DifferenceSales))
		_ = excelFile.SetCellStr(sheetName, "F"+rowStartStr, fmt.Sprintf("%d", item.OrderNum))
		_ = excelFile.SetCellStr(sheetName, "G"+rowStartStr, fmt.Sprintf("%d", item.CustomerNum))
		_ = excelFile.SetCellStr(sheetName, "H"+rowStartStr, fmt.Sprintf("%d", item.NewCustomerNum))
		_ = excelFile.SetCellStr(sheetName, "I"+rowStartStr, fmt.Sprintf("%f", item.NewCustomerSales))
		_ = excelFile.SetCellStr(sheetName, "J"+rowStartStr, fmt.Sprintf("%d", item.NewCustomerRate)+"%")
		_ = excelFile.SetCellStr(sheetName, "K"+rowStartStr, fmt.Sprintf("%d", item.NewSecondCustomerNum))
		_ = excelFile.SetCellStr(sheetName, "L"+rowStartStr, fmt.Sprintf("%f", item.NewSecondCustomerSales))
		_ = excelFile.SetCellStr(sheetName, "M"+rowStartStr, fmt.Sprintf("%d", item.NewSecondCustomerSalesRate)+"%")
		_ = excelFile.SetCellStr(sheetName, "N"+rowStartStr, fmt.Sprintf("%d", item.OldCustomerNum))
		_ = excelFile.SetCellStr(sheetName, "O"+rowStartStr, fmt.Sprintf("%.0f", item.OldCustomerSales))
		_ = excelFile.SetCellStr(sheetName, "P"+rowStartStr, fmt.Sprintf("%d", item.OldCustomerSalesRate)+"%")

		startRow = startRow + 1
		//设置样式
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
	_, _ = c.Writer.Write(des.Bytes())

	return
}

func (ba *baKanBan) GetOperationDataDetail(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId

	res, rErr := microscrm.ApiBaKanBan{}.GetOperationDataDetail(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

func (ba *baKanBan) ExportOperationDataDetail(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId

	res, rErr := microscrm.ApiBaKanBan{}.GetOperationDataDetail(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}

	bytes, err := json.Marshal(res)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}

	var data ba_kanban.OperationDataDetailResponse
	err = json.Unmarshal(bytes, &data)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}

	excelFile := excelize.NewFile()
	sheetName := "BA销售数据明细"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	//设置宽度
	err = excelFile.SetColWidth(sheetName, "A", "P", 40)
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
	_ = excelFile.SetCellStr(sheetName, "A1", "导购")
	_ = excelFile.SetCellStr(sheetName, "B1", "新二购且一码三关注指标人数")
	_ = excelFile.SetCellStr(sheetName, "C1", "新二购且一码三关注完成人数")
	_ = excelFile.SetCellStr(sheetName, "D1", "新二购一码三关注率")
	_ = excelFile.SetCellStr(sheetName, "E1", "新二购工资完成进度")
	_ = excelFile.SetCellStr(sheetName, "F1", "老客且一码三关注指标人数")
	_ = excelFile.SetCellStr(sheetName, "G1", "老客且一码三关注完成人数")
	_ = excelFile.SetCellStr(sheetName, "H1", "老客一码三关注率")
	_ = excelFile.SetCellStr(sheetName, "I1", "老客工资完成进度")
	_ = excelFile.SetCellStr(sheetName, "J1", "新客且一码三关注指标人数")
	_ = excelFile.SetCellStr(sheetName, "K1", "新客且一码三关注完成人数")
	_ = excelFile.SetCellStr(sheetName, "L1", "新客一码三关注率")
	_ = excelFile.SetCellStr(sheetName, "M1", "新客工资完成进度")
	startRow := 2
	//填充数据
	style, err = excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"}}`)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	for _, item := range data.List {
		rowStartStr := strconv.Itoa(startRow)
		_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.Name)
		_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, fmt.Sprintf("%d", item.MonthlyNewCustomerAndThreeScanNum))
		_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, fmt.Sprintf("%d", item.NewSecondAndThreeScanNum))
		_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, fmt.Sprintf("%d", item.NewSecondAndThreeScanRate))
		_ = excelFile.SetCellStr(sheetName, "E"+rowStartStr, fmt.Sprintf("%d", item.NewSecondAndThreeScanRate))
		_ = excelFile.SetCellStr(sheetName, "F"+rowStartStr, fmt.Sprintf("%d", item.MonthlyOldCustomerAndThreeScanNum))
		_ = excelFile.SetCellStr(sheetName, "G"+rowStartStr, fmt.Sprintf("%d", item.OldAndThreeScanNum))
		_ = excelFile.SetCellStr(sheetName, "H"+rowStartStr, fmt.Sprintf("%d", item.OldAndThreeScanRate))
		_ = excelFile.SetCellStr(sheetName, "I"+rowStartStr, fmt.Sprintf("%d", item.OldAndThreeScanRate))
		_ = excelFile.SetCellStr(sheetName, "J"+rowStartStr, fmt.Sprintf("%d", item.MonthlyNewCustomerAndThreeScanNum))
		_ = excelFile.SetCellStr(sheetName, "K"+rowStartStr, fmt.Sprintf("%d", item.NewCustomerAndThreeScanNum))
		_ = excelFile.SetCellStr(sheetName, "L"+rowStartStr, fmt.Sprintf("%d", item.NewCustomerAndThreeScanRate))
		_ = excelFile.SetCellStr(sheetName, "M"+rowStartStr, fmt.Sprintf("%d", item.NewCustomerAndThreeScanRate))

		startRow = startRow + 1
		//设置样式
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
	_, _ = c.Writer.Write(des.Bytes())

	return
}

func (ba *baKanBan) GetBADeptInfo(c *gin.Context) {
	var receive ba_kanban.BaSourceDeptRequst

	if err := c.ShouldBindJSON(&receive); err != nil {
		ginplus.ResError(c, err)
		return
	}
	res, err := microscrm.ApiBaKanBan{}.GetBAInfo(ctx.Wrap(c), receive)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

func (ba *baKanBan) GetSalesDataDetailV1(c *gin.Context) {
	var receive ba_kanban.BaKanBanDataRequest

	if err := c.ShouldBindJSON(&receive); err != nil {
		ginplus.ResError(c, err)
		return
	}
	res, err := microscrm.ApiBaKanBan{}.GetKanBanInfo(ctx.Wrap(c), receive)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}
func (ba *baKanBan) BaSaleListOperation(c *gin.Context) {
	var receive ba_kanban.BaKanBanpOerationDataRequest

	if err := c.ShouldBindJSON(&receive); err != nil {
		ginplus.ResError(c, err)
		return
	}
	res, err := microscrm.ApiBaKanBan{}.BaSaleList(ctx.Wrap(c), receive)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}
func (ba *baKanBan) GrowthRate(c *gin.Context) {
	var receive ba_kanban.BaKanBanpOerationDataRequest

	if err := c.ShouldBindJSON(&receive); err != nil {
		ginplus.ResError(c, err)
		return
	}
	res, err := microscrm.ApiBaKanBan{}.GrowthRate(ctx.Wrap(c), receive)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}
func (ba *baKanBan) ExportbaSale(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId

	res, rErr := microscrm.ApiBaKanBan{}.GetSalesDataDetailV1(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}

	bytes, err := json.Marshal(res)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}

	var data ba_kanban.OperationDataDetailResponseV1
	err = json.Unmarshal(bytes, &data)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}

	excelFile := excelize.NewFile()
	sheetName := "BA销售数据明细"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	//设置宽度
	err = excelFile.SetColWidth(sheetName, "A", "P", 40)
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
	_ = excelFile.SetCellStr(sheetName, "A1", "导购")
	_ = excelFile.SetCellStr(sheetName, "B1", "新二购且一码三关注指标人数")
	_ = excelFile.SetCellStr(sheetName, "C1", "新二购且一码三关注完成人数")
	_ = excelFile.SetCellStr(sheetName, "D1", "新二购一码三关注率")
	_ = excelFile.SetCellStr(sheetName, "E1", "新二购工资完成进度")
	_ = excelFile.SetCellStr(sheetName, "F1", "老客且一码三关注指标人数")
	_ = excelFile.SetCellStr(sheetName, "G1", "老客且一码三关注完成人数")
	_ = excelFile.SetCellStr(sheetName, "H1", "老客一码三关注率")
	_ = excelFile.SetCellStr(sheetName, "I1", "老客工资完成进度")
	_ = excelFile.SetCellStr(sheetName, "J1", "新客且一码三关注指标人数")
	_ = excelFile.SetCellStr(sheetName, "K1", "新客且一码三关注完成人数")
	_ = excelFile.SetCellStr(sheetName, "L1", "新客一码三关注率")
	_ = excelFile.SetCellStr(sheetName, "M1", "新客工资完成进度")
	startRow := 2
	//填充数据
	style, err = excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"}}`)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	for _, item := range data.List {
		rowStartStr := strconv.Itoa(startRow)
		_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.BaName)
		_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, fmt.Sprintf("%d", item.MonthlySecondBuyAndThreeScan))
		_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, fmt.Sprintf("%d", item.NewSecondAndThreeScan))
		_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, fmt.Sprintf("%d", item.NewSecondAndThreeScanRate)+"%")
		_ = excelFile.SetCellStr(sheetName, "E"+rowStartStr, fmt.Sprintf("%.0f", item.NewSecondAndThreeScanWageProgress)+"%")
		_ = excelFile.SetCellStr(sheetName, "F"+rowStartStr, fmt.Sprintf("%d", item.MonthlyOldCustomerAndThreeScan))
		_ = excelFile.SetCellStr(sheetName, "G"+rowStartStr, fmt.Sprintf("%d", item.OldAndThreeScan))
		_ = excelFile.SetCellStr(sheetName, "H"+rowStartStr, fmt.Sprintf("%d", item.MonthlyOldCustomerAndThreeScanRate)+"%")
		_ = excelFile.SetCellStr(sheetName, "I"+rowStartStr, fmt.Sprintf("%.0f", item.MonthlyOldCustomerAndThreeWageProgress)+"%")
		_ = excelFile.SetCellStr(sheetName, "J"+rowStartStr, fmt.Sprintf("%d", item.MonthlyNewCustomerAndThreeScan))
		_ = excelFile.SetCellStr(sheetName, "K"+rowStartStr, fmt.Sprintf("%d", item.NewCustomerAndThreeScan))
		_ = excelFile.SetCellStr(sheetName, "L"+rowStartStr, fmt.Sprintf("%d", item.NewCustomerAndThreeScanRate)+"%")
		_ = excelFile.SetCellStr(sheetName, "M"+rowStartStr, fmt.Sprintf("%d", 0))

		startRow = startRow + 1
		//设置样式
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
	_, _ = c.Writer.Write(des.Bytes())

	return
}
