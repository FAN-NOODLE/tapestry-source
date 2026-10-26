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

type PointsHandle struct {
	handle.BaseHandle
	api microscrm.PointsApi
}

func NewPointsHandle() PointsHandle {
	return PointsHandle{
		api: microscrm.NewPointsApi(),
	}
}

// RecordList 导购获取积分历史记录列表
func (h PointsHandle) RecordList(c *gin.Context) {
	var request excitationReq.PointsRecordReq
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

func (h PointsHandle) DownloadRecord(c *gin.Context) {
	var request excitationReq.PointsRecordReq
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
	if helper.GetEnterpriseConfig(c, enterpriseId) == "coach" {
		sourceTitle = []string{
			"member_user_id", "member_name", "department_name", "integral_rule_name", "incentive_type_str", "rule_category_name", "points_remark", "points", "coins", "source_type_str", "created_time_str", "fy_time_str",
		}
		sourceTitleAlias = map[string]string{
			"member_user_id": "导购id", "member_name": "导购姓名", "department_name": "所在部门", "integral_rule_name": "关联规则名称", "incentive_type_str": "积分维度",
			"rule_category_name": "规则类目", "points_remark": "积分获取条件", "points": "积分", "coins": "金币", "source_type_str": "积分来源", "created_time_str": "获取时间", "fy_time_str": "财年财月",
		}
		sourceTitleValue = []string{
			"points", "coins",
		}
	} else {
		sourceTitle = []string{
			"member_user_id", "member_name", "department_name", "integral_rule_name", "incentive_type_str", "rule_category_name", "points_remark", "points", "source_type_str", "created_time_str", "fy_time_str",
		}
		sourceTitleAlias = map[string]string{
			"member_user_id": "导购id", "member_name": "导购姓名", "department_name": "所在部门", "integral_rule_name": "关联规则名称", "incentive_type_str": "积分维度",
			"rule_category_name": "规则类目", "points_remark": "积分获取条件", "points": "积分", "source_type_str": "积分来源", "created_time_str": "获取时间", "fy_time_str": "财年财月",
		}
	}

	_, err := downloadcenter.DC.CreateExport(ctx.Wrap(c), downloadcenter.CreateExportRequest{
		FileName:           "Wejourney--积分兑领数据",
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
			SourceRestfulPath:             microscrm.GetUrl("%s/micro-scrm/v1/excitation/points/record-list"),
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
