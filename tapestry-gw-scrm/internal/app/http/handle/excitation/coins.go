package excitation

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/app/http/handle"
	excitationReq "gw-scrm/internal/app/http/request/excitation"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/helper"
	"gw-scrm/internal/pkg/sdk/downloadcenter"
	"gw-scrm/internal/pkg/sdk/microscrm"
	"strconv"
)

type CoinsHandle struct {
	handle.BaseHandle
	api microscrm.CoinsApi
}

func NewCoinsHandle() CoinsHandle {
	return CoinsHandle{
		api: microscrm.NewCoinsApi(),
	}
}

// RecordList 导购获取积分历史记录列表
func (h CoinsHandle) RecordList(c *gin.Context) {
	var request excitationReq.CoinsRecordReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if request.Page == 0 {
		request.Page = 1
	}

	if request.Size == 0 {
		request.Size = 10
	}
	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	res, err := h.api.RecordList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

// CoinsSettingsInfo  金币配置
func (h CoinsHandle) CoinsSettingsInfo(c *gin.Context) {
	var request excitationReq.GetCoinSettingInfoReq
	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	request.HasRecord = true
	res, err := h.api.CoinsSettingsInfo(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

// SaveCoinsSettings  金币配置设置
func (h CoinsHandle) SaveCoinsSettings(c *gin.Context) {
	var request excitationReq.SaveCoinsSettingsReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	request.AdminId = adminId
	res, err := h.api.SaveCoinsSettings(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

func (h CoinsHandle) DownloadRecord(c *gin.Context) {
	var request excitationReq.CoinsRecordReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	params := request
	sourceTitle := make([]string, 0)
	sourceTitleValue := make([]string, 0)
	sourceTitleAlias := make(map[string]string, 0)

	sourceTitle = []string{
		"member_user_id", "member_name", "department_name", "change_type_str", "coins", "source_type_str", "coins_remark", "order_num", "created_time_str",
	}
	sourceTitleAlias = map[string]string{
		"member_user_id": "导购id", "member_name": "导购姓名", "department_name": "所在部门", "change_type_str": "类型",
		"coins": "金币数", "source_type_str": "来源", "coins_remark": "变动原因", "order_num": "订单号", "created_time_str": "变更时间",
	}
	sourceTitleValue = []string{
		"coins",
	}

	_, err := downloadcenter.DC.CreateExport(ctx.Wrap(c), downloadcenter.CreateExportRequest{
		FileName:           "Wejourney--金币明细数据",
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
			SourceRestfulPath:             microscrm.GetUrl("%s/micro-scrm/v1/excitation/coins/record-list"),
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
