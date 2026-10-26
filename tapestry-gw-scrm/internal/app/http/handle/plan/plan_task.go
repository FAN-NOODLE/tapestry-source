package plan

import (
	"encoding/json"
	"gw-scrm/internal/app/http/handle"
	"gw-scrm/internal/app/http/request/plan"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/helper"
	"gw-scrm/internal/pkg/sdk"
	"gw-scrm/internal/pkg/sdk/downloadcenter"
	"gw-scrm/internal/pkg/sdk/microscrm"
	"os"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
)

type TaskHandle struct {
	handle.BaseHandle
}

func NewTaskHandle() TaskHandle {
	return TaskHandle{}
}

func (h TaskHandle) TaskDetailExport(c *gin.Context) {
	var request plan.TaskDetailExportReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	result, err := sdk.SendPostNoLog(h.WarpContext(c), microscrm.GetUrl("%s/micro-scrm/v1/plan/task-detail-export-base"), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	marshal, err := json.Marshal(result)
	if err != nil {
		h.ResError(c, err)
		return
	}
	var taskAllotBase plan.TaskAllotDetailExportBase
	err = json.Unmarshal(marshal, &taskAllotBase)
	if err != nil {
		h.ResError(c, err)
		return
	}
	req := plan.TaskDetailExportListReq{
		BaseReq: plan.BaseReq{
			EnterpriseId: enterpriseId,
		},
		PlanId:         taskAllotBase.PlanId,
		TaskAllotLogId: taskAllotBase.TaskAllotLogIds,
		Page:           0,
		Size:           1000,
		Total:          0,
	}
	fileName := time.Now().Format("20060102150405") + ".xlsx"
	excelFile := excelize.NewFile()
	sheetName := "任务数据明细"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	excelFile.SetActiveSheet(index)
	err = excelFile.SaveAs(fileName)
	defer os.Remove(fileName)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	stream, err := excelFile.NewStreamWriter(sheetName)
	if err != nil {
		return
	}
	groupSendName := "群发结果状态"
	if taskAllotBase.TaskCategory == 1 {
		groupSendName = "朋友圈群发结果状态"
	}

	row := make([]any, 0)

	row = []any{
		"任务计划名称", "子任务名称", "任务标签", "子任务类型", "触达导购", "导购ID", "所属门店", "门店属性", "子任务是否完成", "完成时间", "Union ID", "Customer ID", "是否匹配成功",
		groupSendName, "是否点击转发给顾客", "是否点击立即沟通", "是否点击拨打电话", "联系说明", "是否点击添加好友", "顾客状态",
	}
	err = stream.SetRow("A1", row)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	memberListInfo := make(map[string]plan.TaskAllotMemberLists)
	for _, v := range taskAllotBase.MemberLists {
		info := v
		if _, ok := memberListInfo[info.SaCode]; !ok {
			memberListInfo[info.SaCode] = info
		}
	}
	startRow := 2
	for {
		req.Page++
		list, err := sdk.SendPostNoLog(h.WarpContext(c), microscrm.GetUrl("%s/micro-scrm/v1/plan/task-detail-export-list"), req)
		marshal, err = json.Marshal(list)
		if err != nil {
			log.Error(ctx.Wrap(c)).Msgf("解析数据失败,err:%v", err)
			continue
		}
		var dataLists []plan.TaskDetailExportRes
		err = json.Unmarshal(marshal, &dataLists)
		if err != nil {
			log.Error(ctx.Wrap(c)).Msgf("解析数据失败,err:%v", err)
			continue
		}
		if len(dataLists) == 0 {
			break
		}
		for _, item := range dataLists {
			rowStartStr := strconv.Itoa(startRow)
			finishedTimeStr := "--"
			if taskAllotBase.TaskFinish[item.Id].EndTime > 0 {
				finishedTimeStr = helper.GetFormatDate(int64(taskAllotBase.TaskFinish[item.Id].EndTime), 4)
			}
			categoryName := "--"
			switch taskAllotBase.Category {
			case 1:
				categoryName = "手动创建"
			case 2:
				categoryName = "任务模板创建"
			}
			taskCategoryName := "--"
			switch taskAllotBase.TaskCategory {
			case 0:
				taskCategoryName = "沟通任务"
			case 1:
				taskCategoryName = "朋友圈任务"
			}
			touchName := "--"
			if taskAllotBase.GroupSendType == 1 {
				if taskAllotBase.TaskCategory == 0 {
					switch item.IsTouch {
					case 0:
						touchName = "未发送"
					case 1:
						touchName = "成功"
					case 2:
						touchName = "好友已删除 "
					case 3:
						touchName = "接收次数已达上限 "
					}
				} else {
					touchName = "未发送"
					if item.IsTouch == 1 {
						touchName = "可展示"
					}
					if taskAllotBase.TaskFinish[item.Id].IsFinished == 1 && item.IsTouch == 0 {
						touchName = "不可展示"
					}
				}
			}
			row = []any{
				taskAllotBase.PlanName, taskAllotBase.TaskName, categoryName, taskCategoryName, memberListInfo[item.UserId].SaName, memberListInfo[item.UserId].SaCode,
				memberListInfo[item.UserId].StoreName, memberListInfo[item.UserId].StoreType, threeBoolValue(taskAllotBase.TaskFinish[item.Id].IsFinished), finishedTimeStr, item.UnionID, item.CustomerID,
				threeBoolValue(item.IsMatch), touchName, threeBoolValue(taskAllotBase.TaskForwards[item.Id]), threeBoolValue(item.IsContact), threeBoolValue(item.IsCall), item.Remark, threeBoolValue(item.IsAdd), threeContactValue(item.IsFinished),
			}
			err = stream.SetRow("A"+rowStartStr, row)
			startRow++
		}
	}
	err = stream.Flush()
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	err = excelFile.SaveAs(fileName)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	des, err := os.ReadFile(fileName)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	_, _ = c.Writer.Write(des)
}

// NewStat 主任务的任务数据新接口
func (h TaskHandle) NewStat(c *gin.Context) {
	var request plan.NewStatReq
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	stat, err := microscrm.TaskNewDataStat(ctx.Wrap(c), &request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, &stat)
}

// NewStatDetail 主任务的数据明细接口
func (h TaskHandle) NewStatDetail(c *gin.Context) {
	var request plan.NewStatDetailReq
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	stat, err := microscrm.TaskNewDataStatDetail(ctx.Wrap(c), &request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, &stat)
}

func (h TaskHandle) NewStatDetailExport(c *gin.Context) {
	var request plan.NewStatDetailReq
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	stat, err := microscrm.TaskNewDataStatDetail(ctx.Wrap(c), &request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	marshal, err := json.Marshal(stat)
	if err != nil {
		h.ResError(c, err)
		return
	}
	var result plan.NewStatDetailRes
	err = json.Unmarshal(marshal, &result)
	if err != nil {
		h.ResError(c, err)
		return
	}
	excelFile := excelize.NewFile()
	cstSh, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		ginplus.ResError(c, err)
		cstSh = time.Local
	}
	startRow := 2
	if request.Dimension == 1 {
		sheetName := "社群明细数据"
		// 创建sheet
		excelFile.SetSheetName("Sheet1", sheetName)
		// 写入表头
		_ = excelFile.SetCellStr(sheetName, "A1", "子任务名称")
		_ = excelFile.SetCellStr(sheetName, "B1", "子任务类型")
		_ = excelFile.SetCellStr(sheetName, "C1", "子任务群发类型")
		_ = excelFile.SetCellStr(sheetName, "D1", "条件维度")
		_ = excelFile.SetCellStr(sheetName, "E1", "触达群主")
		_ = excelFile.SetCellStr(sheetName, "F1", "匹配社群数")
		_ = excelFile.SetCellStr(sheetName, "G1", "任务总数（件）")
		_ = excelFile.SetCellStr(sheetName, "H1", "已完成任务数（件）")
		_ = excelFile.SetCellStr(sheetName, "I1", "任务完成率（%）")
		for _, item := range result.Data {
			rowStartStr := strconv.Itoa(startRow)
			ReachLordNumber := strconv.FormatInt(item.ReachLordNumber, 10)
			MatchGroupNumber := strconv.FormatInt(item.MatchGroupNumber, 10)
			TaskNumber := strconv.FormatInt(item.TaskNumber, 10)
			var TaskFinishNumber string
			if item.TaskFinishNumber != 0 {
				TaskFinishNumber = strconv.FormatInt(item.TaskFinishNumber, 10)
			} else {
				TaskFinishNumber = "-"
			}
			taskCategoryName := "沟通任务"
			groupSendTypeName := "标签任务群发"
			if item.TaskCategory == 1 {
				taskCategoryName = "朋友圈任务"
				if item.GroupSendType == 1 {
					groupSendTypeName = "企业任务群发"
				}
			} else {
				if item.GroupSendType == 1 {
					groupSendTypeName = "企微一键群发"
				}
				if item.GroupSendType == 2 {
					groupSendTypeName = "灵活沟通转发"
				}
			}

			_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.TaskName)
			_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, taskCategoryName)
			_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, groupSendTypeName)
			_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, "社群任务")
			_ = excelFile.SetCellStr(sheetName, "E"+rowStartStr, ReachLordNumber)
			_ = excelFile.SetCellStr(sheetName, "F"+rowStartStr, MatchGroupNumber)
			_ = excelFile.SetCellStr(sheetName, "G"+rowStartStr, TaskNumber)
			_ = excelFile.SetCellStr(sheetName, "H"+rowStartStr, TaskFinishNumber)
			_ = excelFile.SetCellStr(sheetName, "I"+rowStartStr, item.TaskFinishPercentage)
			startRow = startRow + 1
		}
	} else {
		sheetName := "客户画像明细数据"
		// 创建sheet
		excelFile.SetSheetName("Sheet1", sheetName)
		// 写入表头
		_ = excelFile.SetCellStr(sheetName, "A1", "子任务名称")
		_ = excelFile.SetCellStr(sheetName, "B1", "子任务类型")
		_ = excelFile.SetCellStr(sheetName, "C1", "任务群发类型")
		_ = excelFile.SetCellStr(sheetName, "D1", "条件维度")
		_ = excelFile.SetCellStr(sheetName, "E1", "累计触达导购人数")
		_ = excelFile.SetCellStr(sheetName, "F1", "累计匹配客户数")
		_ = excelFile.SetCellStr(sheetName, "G1", "累计匹配客户数（导购去重）")
		_ = excelFile.SetCellStr(sheetName, "H1", "任务总数（件）")
		_ = excelFile.SetCellStr(sheetName, "I1", "已完成任务数（件）")
		for _, item := range result.Data {
			rowStartStr := strconv.Itoa(startRow)
			ReachGuideNumber := strconv.FormatInt(item.ReachGuideNumber, 10)
			MatchCustomersNumber := strconv.FormatInt(item.MatchCustomersNumber, 10)
			MatchCustomersNumber2 := strconv.FormatInt(item.MatchCustomersNumber2, 10)
			TaskNumber := strconv.FormatInt(item.TaskNumber, 10)
			var TaskFinishNumber string
			if item.TaskFinishNumber != 0 {
				TaskFinishNumber = strconv.FormatInt(item.TaskFinishNumber, 10)
			} else {
				TaskFinishNumber = "-"
			}
			taskCategoryName := "沟通任务"
			groupSendTypeName := "标签任务群发"
			if item.TaskCategory == 1 {
				taskCategoryName = "朋友圈任务"
				if item.GroupSendType == 1 {
					groupSendTypeName = "企业任务群发"
				}
			} else {
				if item.GroupSendType == 1 {
					groupSendTypeName = "企微一键群发"
				}
				if item.GroupSendType == 2 {
					groupSendTypeName = "灵活沟通转发"
				}
			}
			_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.TaskName)
			_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, taskCategoryName)
			_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, groupSendTypeName)
			if item.Dimension == 2 {
				_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, "用户画像")
			}
			if item.Dimension == 3 {
				_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, "导入顾客名单")
			}
			if item.Dimension == 4 {
				_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, "导入顾客&导购匹配关系")
			}
			if item.Dimension == 1 {
				_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, "社群任务")
			}
			_ = excelFile.SetCellStr(sheetName, "E"+rowStartStr, ReachGuideNumber)
			_ = excelFile.SetCellStr(sheetName, "F"+rowStartStr, MatchCustomersNumber)
			_ = excelFile.SetCellStr(sheetName, "G"+rowStartStr, MatchCustomersNumber2)
			_ = excelFile.SetCellStr(sheetName, "H"+rowStartStr, TaskNumber)
			_ = excelFile.SetCellStr(sheetName, "I"+rowStartStr, TaskFinishNumber)
			startRow = startRow + 1
		}
	}
	fileName := "任务明细数据导出-" + time.Now().In(cstSh).Format("200601021504") + ".xlsx"
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	des, _ := excelFile.WriteToBuffer()
	_, _ = c.Writer.Write(des.Bytes())
	return

}

func threeBoolValue(value int) string {
	str := "否"
	if value == 1 {
		str = "是"
	}
	return str
}
func threeContactValue(value int) string {
	str := "待联系"
	if value == 1 {
		str = "已联系"
	}
	return str
}

// TaskDetailExportNew 子任务数据明细导出---下载中心
func (h TaskHandle) TaskDetailExportNew(c *gin.Context) {
	var request plan.TaskDetailExportReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	enterpriseId := tokenInfo.EnterpriseId
	adminId := tokenInfo.AdminId
	request.EnterpriseId = enterpriseId
	result, err := sdk.SendPostNoLog(h.WarpContext(c), microscrm.GetUrl("%s/micro-scrm/v1/plan/task-detail-export-base"), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	marshal, err := json.Marshal(result)
	if err != nil {
		h.ResError(c, err)
		return
	}
	var taskAllotBase plan.TaskAllotDetailExportBase
	err = json.Unmarshal(marshal, &taskAllotBase)
	if err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseConfigName := helper.GetEnterpriseConfig(c, enterpriseId)
	groupSendName := "群发结果状态"
	if taskAllotBase.TaskCategory == 1 {
		groupSendName = "朋友圈群发结果状态"
	}
	sourceTitle := make([]string, 0)
	sourceTitleAlias := make(map[string]string, 0)
	if enterpriseConfigName != "coach" || taskAllotBase.GroupSendType < 2 { // 非coach的企业，和群发类型不是企业沟通群发类型，下载模板为之前模板
		sourceTitle = []string{
			"plan_name", "task_name", "category_name", "task_category_name", "sa_name", "sa_code", "store_name", "store_type", "is_task_finished", "finished_time_str", "union_id", "customer_id",
			"is_match", "touch_name", "is_forward", "is_contact", "is_call", "remark", "is_add", "customer_status", "customer_finished_time_str",
		}
		sourceTitleAlias = map[string]string{
			"plan_name": "任务计划名称", "task_name": "子任务名称", "category_name": "任务标签", "task_category_name": "子任务类型", "sa_name": "触达导购", "sa_code": "导购ID", "store_name": "所属门店",
			"store_type": "门店属性", "is_task_finished": "子任务是否完成", "finished_time_str": "完成时间", "union_id": "Union ID", "customer_id": "Customer ID", "is_match": "是否匹配成功",
			"touch_name": groupSendName, "is_forward": "是否点击转发给顾客", "is_contact": "是否点击立即沟通", "is_call": "是否点击拨打电话", "remark": "联系说明", "is_add": "是否点击添加好友", "customer_status": "顾客状态", "customer_finished_time_str": "联系时间",
		}
	} else {
		sourceTitle = []string{
			"plan_name", "task_name", "category_name", "task_category_name", "sa_name", "sa_code", "store_name", "store_type", "is_task_finished", "finished_time_str", "union_id", "customer_id",
			"is_match", "is_forward", "is_contact", "is_call", "remark", "is_add", "customer_status", "customer_finished_time_str", "touch_type",
		}
		sourceTitleAlias = map[string]string{
			"plan_name": "任务计划名称", "task_name": "子任务名称", "category_name": "任务标签", "task_category_name": "子任务类型", "sa_name": "触达导购", "sa_code": "导购ID", "store_name": "所属门店",
			"store_type": "门店属性", "is_task_finished": "子任务是否完成", "finished_time_str": "完成时间", "union_id": "Union ID", "customer_id": "Customer ID", "is_match": "是否匹配成功",
			"is_forward": "是否点击转发给顾客", "is_contact": "是否点击立即沟通", "is_call": "是否点击拨打电话", "remark": "联系说明", "is_add": "是否点击添加好友", "customer_status": "顾客状态", "customer_finished_time_str": "联系时间", "touch_type": "触达方式",
		}
	}

	params := plan.TaskDetailExportListNewReq{
		BaseReq: plan.BaseReq{
			EnterpriseId: enterpriseId,
		},
		TaskAllotDetailExportBase: taskAllotBase,
		TaskAllotLogId:            taskAllotBase.TaskAllotLogIds,
	}
	_, err = downloadcenter.DC.CreateExport(ctx.Wrap(c), downloadcenter.CreateExportRequest{
		FileName:           "任务计划--子任务详情数据明细导出",
		Type:               downloadcenter.FileTypeExcel,
		StoreId:            1,
		EnterpriseId:       enterpriseId,
		OperationId:        adminId,
		OperationName:      "adminId",
		SourceType:         downloadcenter.SourceTypeRestful,
		SourceTitleAlias:   sourceTitleAlias,
		SourceTitle:        sourceTitle,
		SourceHeaderPrefix: nil,
		Label:              downloadcenter.ReserveOrderExport,
		CreateData:         downloadcenter.CreateData{},
		CreateSql:          downloadcenter.CreateSql{},
		CreateRestful: downloadcenter.CreateRestful{
			SourceRestfulPath:             microscrm.GetUrl("%s/micro-scrm/v1/plan/task-detail-export-list-new"),
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
}

func (h TaskHandle) IsShow(c *gin.Context) {
	var request plan.NewStatReq
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseId
	result, err := sdk.SendPost(h.WarpContext(c), microscrm.GetUrl("%s/micro-scrm/v1/plan/is-show"), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, &result)
}

func (h TaskHandle) TaskPlanNameList(c *gin.Context) {
	var request plan.TaskPlanNameListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	enterpriseId := tokenInfo.EnterpriseId
	request.EnterpriseId = enterpriseId
	result, err := sdk.SendPostNoLog(h.WarpContext(c), microscrm.GetUrl("%s/micro-scrm/v1/plan/task/plan-name-list"), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
}
