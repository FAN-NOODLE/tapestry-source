package display

import (
	"fmt"
	"gw-scrm/internal/app/http/handle"
	request2 "gw-scrm/internal/app/http/request"
	"gw-scrm/internal/app/http/request/display"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/helper"
	"gw-scrm/internal/pkg/sdk/downloadcenter"
	"gw-scrm/internal/pkg/sdk/microscrm"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
)

type TaskHandle struct {
	handle.BaseHandle
	displaySdk microscrm.DisplayApi
}

func NewTaskHandle() TaskHandle {
	return TaskHandle{
		displaySdk: microscrm.NewDisplayApi(),
	}
}

func (h TaskHandle) AddTask(c *gin.Context) {
	var request display.AddTaskReq
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	create, err := h.displaySdk.Create(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, create)
}

func (h TaskHandle) List(c *gin.Context) {
	var request display.SearchReq
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	list, err := h.displaySdk.List(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h TaskHandle) Detail(c *gin.Context) {
	var request display.TaskBaseReq
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	info, err := h.displaySdk.Detail(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h TaskHandle) Copy(c *gin.Context) {
	var request display.TaskBaseReq
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	res, err := h.displaySdk.Copy(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h TaskHandle) UpdateShow(c *gin.Context) {
	var request display.TaskBaseReq
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	res, err := h.displaySdk.Status(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h TaskHandle) UpdateTask(c *gin.Context) {
	var request display.AddTaskReq
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	create, err := h.displaySdk.Update(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, create)
}

func (h TaskHandle) AllotList(c *gin.Context) {
	var request display.SearchAllotReq
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	create, err := h.displaySdk.AllotList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, create)
}

func (h TaskHandle) AllotNum(c *gin.Context) {
	var request display.SearchAllotBaseReq
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	create, err := h.displaySdk.AllotNum(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, create)
}

func (h TaskHandle) AllotListExport(c *gin.Context) {
	var request display.SearchAllotBaseReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	enterpriseId := tokenInfo.EnterpriseId
	adminId := tokenInfo.AdminId
	request.EnterpriseId = enterpriseId
	sourceTitle := make([]string, 0)
	sourceTitleAlias := make(map[string]string, 0)
	sourceTitleValue := make([]string, 0)
	sourceTitle = []string{
		"display_name", "store_code", "store_name", "store_type_str", "store_city", "finish_status_str", "grade_str", "submit_name", "owner_member_name",
		"submit_time_str", "review_time_str", "review_name",
	}
	sourceTitleAlias = map[string]string{
		"display_name": "标题", "store_code": "门店ID", "store_name": "门店名称", "store_type_str": "门店类型", "store_city": "所属城市", "finish_status_str": "完成情况", "grade_str": "评分",
		"submit_name": "最近提交导购", "owner_member_name": "所属VM", "submit_time_str": "最近提交时间", "review_time_str": "最近审批时间", "review_name": "最近审批人",
	}

	// 邀约单导出
	_, err := downloadcenter.DC.CreateExport(ctx.Wrap(c), downloadcenter.CreateExportRequest{
		FileName:           "VM陈列--陈列数据导出",
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
		Label:              downloadcenter.VmListExportSave,
		CreateData:         downloadcenter.CreateData{},
		CreateSql:          downloadcenter.CreateSql{},
		CreateRestful: downloadcenter.CreateRestful{
			SourceRestfulPath:             microscrm.GetUrl("%s/micro-scrm/v1/display/task/allot-list-export"),
			SourceRestfulMethod:           "POST",
			SourceRestfulHeaders:          map[string]interface{}{"enterprise-id": strconv.FormatInt(enterpriseId, 10)},
			SourceRestfulParams:           helper.StructToMap(request),
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
}

func (h TaskHandle) NoticeSetting(c *gin.Context) {
	var request display.NoticeSettingReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	info, err := h.displaySdk.NoticeSetting(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h TaskHandle) AllotInfo(c *gin.Context) {
	var request display.AllotBaseReq
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	request.IsAdmin = true
	info, err := h.displaySdk.AllotInfo(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h TaskHandle) AllotCheck(c *gin.Context) {
	var request display.StoreTaskReviewReq
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	request.AdminId = adminId

	info, err := h.displaySdk.AllotCheck(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h TaskHandle) NoticeSettingInfo(c *gin.Context) {
	var request request2.EnterpriseBaseReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	info, err := h.displaySdk.NoticeSettingInfo(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h TaskHandle) AllotSendNotice(c *gin.Context) {
	var request display.SearchAllotBaseReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	info, err := h.displaySdk.AllotSendNotice(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h TaskHandle) SaveImage(c *gin.Context) {
	var request display.SaveImageReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	enterpriseId := tokenInfo.EnterpriseId
	adminId := tokenInfo.AdminId
	request.EnterpriseId = enterpriseId
	info, err := h.displaySdk.GetImageList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	list := info.([]interface{})
	sourceData := make([][]interface{}, 0)
	sourceData = append(sourceData, list)
	sourceTitle := make([]string, 0)
	sourceTitleAlias := make(map[string]string, 0)
	sourceTitleValue := make([]string, 0)
	fileName := fmt.Sprintf("%s-%s-%s", request.DisplayName, request.StoreCode, request.StoreName)
	_, err = downloadcenter.DC.CreateExport(ctx.Wrap(c), downloadcenter.CreateExportRequest{
		FileName:           fileName,
		Type:               downloadcenter.FileTypeZip,
		StoreId:            1,
		EnterpriseId:       enterpriseId,
		OperationId:        adminId,
		OperationName:      "adminId",
		SourceType:         downloadcenter.SourceTypeData,
		SourceTitleAlias:   sourceTitleAlias,
		SourceTitle:        sourceTitle,
		SourceTitleValue:   sourceTitleValue,
		SourceHeaderPrefix: nil,
		Label:              downloadcenter.VmImageSave,
		CreateData: downloadcenter.CreateData{
			SourceData: sourceData,
		},
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, "下载成功")
}

func (h TaskHandle) AllotNoticeNum(c *gin.Context) {
	var request display.SearchAllotBaseReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	info, err := h.displaySdk.AllotNoticeNum(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h TaskHandle) VmList(c *gin.Context) {
	var request request2.EnterpriseBaseReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	info, err := h.displaySdk.VmList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h TaskHandle) CosAvatar(c *gin.Context) {
	var request display.CosAvatar
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.displaySdk.CosAvatar(ctx.Wrap(c), &request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h TaskHandle) ReviewFileCheck(c *gin.Context) {
	var request display.ReviewFileReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	res, err := h.displaySdk.ReviewFileCheck(ctx.Wrap(c), &request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h TaskHandle) ReviewExport(c *gin.Context) {
	var request display.TaskBaseReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	res, err := h.displaySdk.ReviewExport(ctx.Wrap(c), &request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}
