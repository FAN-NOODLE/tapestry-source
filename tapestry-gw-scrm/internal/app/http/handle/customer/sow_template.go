package customer

import (
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/app/http/handle"
	customerReq "gw-scrm/internal/app/http/request/customer"
	"gw-scrm/internal/pkg/common"
	scrm_errs "gw-scrm/internal/pkg/errs"
	"gw-scrm/internal/pkg/helper"
	"gw-scrm/internal/pkg/sdk/downloadcenter"
	"gw-scrm/internal/pkg/sdk/microscrm"
	"strconv"
	"time"
)

type SowTaskHandle struct {
	handle.BaseHandle
}

func NewSowTaskHandle() SowTaskHandle {
	return SowTaskHandle{}
}

// List  播种模版列表
func (h SowTaskHandle) List(c *gin.Context) {
	var request customerReq.SowTemplateListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId

	result, err := microscrm.SowTemplateList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)
	return
}

func (h SowTaskHandle) Save(c *gin.Context) {
	var request customerReq.SowTemplateSaveReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseId, _, mobile := common.GetEnterpriseAndMobile(c)
	request.EnterpriseID = enterpriseId
	request.Mobile = mobile
	result, err := microscrm.SowTemplateSave(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)
	return
}

func (h SowTaskHandle) Detail(c *gin.Context) {
	var request customerReq.SowTemplateDetailReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId
	result, err := microscrm.SowTemplateDetail(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)

	return
}

func (h SowTaskHandle) AllTaskList(c *gin.Context) {
	var request customerReq.SowAllTaskListReq
	err := c.ShouldBind(&request)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId

	res, err := microscrm.SowAllTaskList(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

func (h SowTaskHandle) UpdateStatus(c *gin.Context) {
	var request customerReq.SowTemplateDetailReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId
	result, err := microscrm.SowTemplateUpdateStatus(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)

	return
}

func (h SowTaskHandle) Delete(c *gin.Context) {
	var request customerReq.SowTemplateDetailReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId
	result, err := microscrm.SowTemplateDelete(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)

	return
}

func (h SowTaskHandle) DownloadSowTaskRecord(c *gin.Context) {
	var request customerReq.SowAllTaskListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId
	params := request
	sourceTitle := make([]string, 0)
	sourceTitleValue := make([]string, 0)
	sourceTitleAlias := make(map[string]string, 0)

	sourceTitle = []string{
		"template_name", "template_type", "template_id_str", "sa_name", "sa_dp_name", "username", "union_code", "status_str", "created_at_str", "remind_time_str", "finished_at_str", "content",
	}
	sourceTitleAlias = map[string]string{
		"template_name": "播种模板名称", "template_type": "播种类型", "template_id_str": "模板ID", "sa_name": "导购姓名", "sa_dp_name": "所在部门", "username": "顾客姓名",
		"union_code": "会员卡号", "status_str": "播种状态", "created_at_str": "创建时间", "remind_time_str": "提醒时间", "finished_at_str": "完成时间", "content": "播种内容",
	}
	sourceTitleValue = []string{}

	_, err := downloadcenter.DC.CreateExport(ctx.Wrap(c), downloadcenter.CreateExportRequest{
		FileName:           "播种记录列表",
		Type:               downloadcenter.FileTypeExcel,
		StoreId:            1,
		EnterpriseId:       enterpriseId,
		OperationId:        adminId,
		OperationName:      "adminId",
		SourceType:         downloadcenter.SourceTypeRestful,
		SourceTitleAlias:   sourceTitleAlias,
		SourceTitle:        sourceTitle,
		SourceTitleValue:   sourceTitleValue,
		SourceHeaderPrefix: nil,
		Label:              downloadcenter.ReserveOrderExport,
		CreateData:         downloadcenter.CreateData{},
		CreateSql:          downloadcenter.CreateSql{},
		CreateRestful: downloadcenter.CreateRestful{
			SourceRestfulPath:             microscrm.GetUrl("%s/micro-scrm/v1/customer/sow-template/all-task-list"),
			SourceRestfulMethod:           "POST",
			SourceRestfulHeaders:          map[string]interface{}{"enterprise-id": strconv.FormatInt(enterpriseId, 10)},
			SourceRestfulParams:           helper.StructToMap(params),
			SourceSuffix:                  nil,
			SourceRestfulCountField:       "data.total",
			SourceRestfulDataField:        "data.list",
			SourceRestfulRequestPageField: "page",
			SourceRestfulRequestSizeField: "size",
		},
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, "下载成功")
	return
}

func (h SowTaskHandle) TemplateListExport(c *gin.Context) {
	var request customerReq.SowTemplateListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID

	request.Size = 9999
	request.Page = 1

	result, err := microscrm.SowTemplateList(h.WarpContext(c), request)

	if err != nil {
		h.ResError(c, err)
		return
	}

	jsonStr, _ := json.Marshal(result)
	var data customerReq.SowTemplateList
	errJson := json.Unmarshal(jsonStr, &data)
	if errJson != nil {
		h.ResError(c, err)
		return
	}
	excelFile := excelize.NewFile()
	sheetName := "播种模版数据"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	// 设置宽度
	err = excelFile.SetColWidth(sheetName, "A", "H", 15)
	if nil != err {
		h.ResError(c, err)
		return
	}
	// 设置样式
	style, styleErr := excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"},"font":{"bold":true}}`)
	if nil != styleErr {
		h.ResError(c, styleErr)
		return
	}
	err = excelFile.SetCellStyle(sheetName, "A1", "I1", style)
	if nil != err {
		h.ResError(c, err)
		return
	}
	_ = excelFile.SetCellStr(sheetName, "A1", "模板ID")
	_ = excelFile.SetCellStr(sheetName, "B1", "播种模板名称")
	_ = excelFile.SetCellStr(sheetName, "C1", "播种次数")
	_ = excelFile.SetCellStr(sheetName, "D1", "使用模版次数")
	_ = excelFile.SetCellStr(sheetName, "E1", "引用导购人数")
	_ = excelFile.SetCellStr(sheetName, "F1", "创建者")
	_ = excelFile.SetCellStr(sheetName, "G1", "创建时间")
	_ = excelFile.SetCellStr(sheetName, "H1", "最近更新时间")

	startRow := 2
	// 填充数据
	style, err = excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"}}`)
	if nil != err {
		h.ResError(c, err)
		return
	}
	for _, item := range data.List {
		rowStartStr := strconv.Itoa(startRow)
		if item.CreatedAt > 0 {
			item.CreatedAtStr = time.Unix(item.CreatedAt, 0).Format("2006-01-02 15:04:05")
		}
		if item.UpdatedAt > 0 {
			item.UpdatedAtStr = time.Unix(item.UpdatedAt, 0).Format("2006-01-02 15:04:05")
		}
		_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, strconv.FormatInt(item.Id, 10))
		_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, item.Name)
		_ = excelFile.SetCellInt(sheetName, "C"+rowStartStr, item.TaskNum)
		_ = excelFile.SetCellInt(sheetName, "D"+rowStartStr, item.UseNum)
		_ = excelFile.SetCellInt(sheetName, "E"+rowStartStr, item.UseSaNum)
		_ = excelFile.SetCellStr(sheetName, "F"+rowStartStr, item.AdminName)
		_ = excelFile.SetCellStr(sheetName, "G"+rowStartStr, item.CreatedAtStr)
		_ = excelFile.SetCellStr(sheetName, "H"+rowStartStr, item.UpdatedAtStr)
		startRow = startRow + 1
		err = excelFile.SetCellStyle(sheetName, "A"+rowStartStr, "H"+rowStartStr, style)
		if nil != err {
			h.ResError(c, err)
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
