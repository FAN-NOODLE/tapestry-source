package invite

import (
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/app/http/handle"
	inviteReq "gw-scrm/internal/app/http/request/invite"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
	"math"
	"strconv"
	"time"
)

type InviteHandle struct {
	handle.BaseHandle
	inviteSdk microscrm.InviteApi
}

func NewInviteHandle() InviteHandle {
	return InviteHandle{
		inviteSdk: microscrm.NewInviteApi(),
	}
}

// TemplateSave 创建或修改模板
func (s InviteHandle) TemplateSave(c *gin.Context) {

	var req inviteReq.TemplateSaveReq
	if err := c.ShouldBind(&req); err != nil {
		s.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseID

	res, err := s.inviteSdk.TemplateSave(ctx.Wrap(c), req)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

// TemplateList 邀请函模板列表
func (s InviteHandle) TemplateList(c *gin.Context) {

	var reqList inviteReq.TemplateListReq
	if err := c.ShouldBind(&reqList); err != nil {
		s.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	reqList.EnterpriseId = enterpriseID

	res, err := s.inviteSdk.TemplateList(ctx.Wrap(c), reqList)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

// TemplateInfo 邀请函模板信息
func (s InviteHandle) TemplateInfo(c *gin.Context) {

	var reqInfo inviteReq.TemplateInfoReq
	if err := c.ShouldBind(&reqInfo); err != nil {
		s.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	reqInfo.EnterpriseId = enterpriseID

	res, err := s.inviteSdk.TemplateInfo(ctx.Wrap(c), reqInfo)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

// TemplateDelete 删除邀请函模板
func (s InviteHandle) TemplateDelete(c *gin.Context) {

	var req inviteReq.TemplateDeleteReq
	if err := c.ShouldBind(&req); err != nil {
		s.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseID

	res, err := s.inviteSdk.TemplateDelete(ctx.Wrap(c), req)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

// RecordOrderList 邀请数据列表
func (s InviteHandle) RecordOrderList(c *gin.Context) {

	var req inviteReq.RecordOrderListReq
	if err := c.ShouldBind(&req); err != nil {
		s.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseID

	res, err := s.inviteSdk.RecordOrderList(ctx.Wrap(c), req)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

// RecordOrderListExport 邀请数据列表导出
func (s InviteHandle) RecordOrderListExport(c *gin.Context) {
	var request inviteReq.RecordOrderListReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	request.Size = 1000
	request.Page = 1
	request.Total = int64(0)

	list, err := s.inviteSdk.RecordOrderList(ctx.Wrap(c), request)

	if err != nil {
		s.ResError(c, err)
		return
	}

	jsonStr, _ := json.Marshal(list)
	var data inviteReq.RecordOrderList
	errJson := json.Unmarshal(jsonStr, &data)
	if errJson != nil {
		s.ResError(c, err)
		return
	}
	var datalist []inviteReq.RecordOrder
	for _, v := range data.List {
		datalist = append(datalist, v)
	}

	limitEnd := request.Page * request.Size
	if int(data.Total) > limitEnd {
		request.Total = data.Total
		maxPage := int(math.Ceil(float64(data.Total / int64(request.Size))))
		for i := 2; i <= maxPage; i++ {
			request.Page = i
			newList, err1 := s.inviteSdk.RecordOrderList(ctx.Wrap(c), request)
			if err1 != nil {
				s.ResError(c, err1)
				return
			}
			jsonStr1, _ := json.Marshal(newList)
			var data1 inviteReq.RecordOrderList
			errJson1 := json.Unmarshal(jsonStr1, &data1)
			if errJson1 != nil {
				s.ResError(c, err)
				return
			}
			for _, v := range data1.List {
				datalist = append(datalist, v)
			}
		}
	}

	excelFile := excelize.NewFile()
	sheetName := "邀请数据列表"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	// 设置宽度
	err = excelFile.SetColWidth(sheetName, "A", "H", 15)
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
	err = excelFile.SetCellStyle(sheetName, "A1", "J1", style)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	// 写入表头
	_ = excelFile.SetCellStr(sheetName, "A1", "模板ID")
	_ = excelFile.SetCellStr(sheetName, "B1", "模板名称")
	_ = excelFile.SetCellStr(sheetName, "C1", "关联预约服务")
	_ = excelFile.SetCellStr(sheetName, "D1", "邀约导购")
	_ = excelFile.SetCellStr(sheetName, "E1", "所在部门")
	_ = excelFile.SetCellStr(sheetName, "F1", "顾客姓名")
	_ = excelFile.SetCellStr(sheetName, "G1", "到店时间")
	_ = excelFile.SetCellStr(sheetName, "H1", "创建时间")

	startRow := 2
	// 填充数据
	style, err = excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"}}`)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	for _, item := range datalist {
		rowStartStr := strconv.Itoa(startRow)
		_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, strconv.FormatInt(item.Id, 10))
		_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, item.TitleName)
		_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, item.ReserveServiceName)
		_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, item.MemberName)
		_ = excelFile.SetCellStr(sheetName, "E"+rowStartStr, item.StoreName)
		_ = excelFile.SetCellStr(sheetName, "F"+rowStartStr, item.UserName)
		_ = excelFile.SetCellStr(sheetName, "G"+rowStartStr, item.ReserveTimeStr)
		_ = excelFile.SetCellStr(sheetName, "H"+rowStartStr, item.CreateTimeStr)

		startRow = startRow + 1
		err = excelFile.SetCellStyle(sheetName, "A"+rowStartStr, "G"+rowStartStr, style)
		if nil != err {
			ginplus.ResError(c, err)
			return
		}
	}
	excelFile.SetActiveSheet(index)
	fileName := time.Now().Format("200601021504") + ".xlsx"

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	des, _ := excelFile.WriteToBuffer()
	_, _ = c.Writer.Write(des.Bytes())
	return

}

func (s InviteHandle) UpdateStatus(c *gin.Context) {

	var req inviteReq.TemplateUpdateStatusReq
	if err := c.ShouldBind(&req); err != nil {
		s.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseID

	res, err := s.inviteSdk.UpdateStatus(ctx.Wrap(c), req)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}
