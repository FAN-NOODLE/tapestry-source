package http

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"gw-scrm/internal/logic"
	"gw-scrm/internal/modules/file/storage"
	"gw-scrm/internal/modules/file/storage/cos"
	"gw-scrm/internal/modules/tag_csv_upload/delivery"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx"

	planTask "gw-scrm/internal/modules/plan_task"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
)

/**
*子任务详情
 */
func planTaskInfo(c *gin.Context) {
	params := planTask.PlanTasksInfoRequst{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	versionInfo, err := microscrm.PlanTaskInfo(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, versionInfo)
}

/**
*新增任务
 */
func planTasksAdd(c *gin.Context) {
	params := planTask.PlanTasksRequst{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if params.Dimension == 3 || params.Dimension == 4 {
		params.ConditionType = 1
		params.TaskType = 1
	}
	versionInfo, err := microscrm.PlanTasksAdd(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, versionInfo)
}

/**
*修改任务
 */
func planTasksUpdate(c *gin.Context) {
	params := planTask.PlanTasksRequst{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	versionInfo, err := microscrm.PlanTasksUpdate(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, versionInfo)
}

/**
*删除任务
 */
func planTasksDelete(c *gin.Context) {
	params := planTask.PlanTasksInfoRequst{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	versionInfo, err := microscrm.PlanTasksDelete(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, versionInfo)
}

/**
*计划列表
 */
func getPlanList(c *gin.Context) {
	params := planTask.PlanListRequest{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	versionInfo, err := microscrm.PlanList(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, versionInfo)
}

/**
*任务列表
 */
func getTaskList(c *gin.Context) {
	params := planTask.TaskListRequest{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	versionInfo, err := microscrm.TaskList(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, versionInfo)
}

/**
//@author:LiKun
//计划任务相关
// planAdd:新增计划任务
//planEdit:计划任务编辑
//planDelete:计划任务删除
//planCopy:计划任务复制
*/

func planAdd(c *gin.Context) {
	var receive planTask.PlanRequest
	if err := c.ShouldBindJSON(&receive); err != nil {
		log.Println(err)
		ginplus.ResError(c, err)
		return
	}
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	if tokenInfo.AdminId != 0 {
		receive.OperationUserId = tokenInfo.AdminId
	}
	if res, err := microscrm.PlanAdd(ctx.Wrap(c), &receive); err != nil {
		ginplus.ResError(c, err)
		return
	} else {
		ginplus.ResSuccess(c, res)
	}
}

func planEdit(c *gin.Context) {
	var receive planTask.ReceivePlanModel
	if err := c.ShouldBindJSON(&receive); err != nil {
		ginplus.ResError(c, err)
		return
	}
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	if tokenInfo.AdminId != 0 {
		receive.OperationUserId = tokenInfo.AdminId
	}

	if res, err := microscrm.PlanEdit(ctx.Wrap(c), &receive); err != nil {
		ginplus.ResError(c, err)
		return
	} else {
		ginplus.ResSuccess(c, res)
	}
}

func planDelete(c *gin.Context) {
	var info = planTask.PlanId{}
	if err := c.ShouldBindJSON(&info); err != nil {
		ginplus.ResError(c, errors.New("Data parsing error"))
		return
	} else {
		if res, err := microscrm.PlanRemove(ctx.Wrap(c), info); err != nil {
			ginplus.ResError(c, err)
			return
		} else {
			ginplus.ResSuccess(c, res)
		}

	}

}

func planCopy(c *gin.Context) {
	var receive planTask.ReceivePlanModel
	if err := c.ShouldBindJSON(&receive); err != nil {
		ginplus.ResError(c, errors.New("Data parsing error"))
		return
	} else {
		tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
		if tokenInfo.AdminId != 0 {
			receive.OperationUserId = tokenInfo.AdminId
		}
		if res, err := microscrm.PlanCopy(ctx.Wrap(c), &receive); err != nil {
			ginplus.ResError(c, err)
			return
		} else {
			ginplus.ResSuccess(c, res)
		}
	}

}

func planStatus(c *gin.Context) {
	var receive planTask.PlanStatus
	if err := c.ShouldBindJSON(&receive); err != nil {
		ginplus.ResError(c, errors.New("Data parsing error"))
		return
	}

	if res, err := microscrm.PlanStatus(ctx.Wrap(c), &receive); err != nil {
		ginplus.ResError(c, err)
		return
	} else {
		ginplus.ResSuccess(c, res)
	}

}
func planTasksConditionsCompletion(c *gin.Context) {
	params := planTask.PlanTasksConditionsRequst{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	versionInfo, err := microscrm.PlanTasksConditionsCompletion(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, versionInfo)
}
func planTasksConditionsUserCompletion(c *gin.Context) {
	params := planTask.PlanTasksConditionsBaRequst{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	versionInfo, err := microscrm.PlanTasksConditionsUserCompletion(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, versionInfo)
}

func customerData(c *gin.Context) {
	params := planTask.CustomerDataReq{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	versionInfo, err := microscrm.CustomerData(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, versionInfo)
}

func customerlistDataExport(c *gin.Context) {
	params := planTask.ExportTasksRequst{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	customerList, err := microscrm.CustomerListDataExport(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	jsonStr, _ := json.Marshal(customerList)
	var data []*planTask.ExportTaskDataResp
	errJson := json.Unmarshal(jsonStr, &data)
	if errJson != nil {
		ginplus.ResError(c, err)
		return
	}
	excelFile := excelize.NewFile()
	sheetName := "子任务用户画像列表数据"
	// 创建sheet
	excelFile.SetSheetName("Sheet1", sheetName)
	// 写入表头
	_ = excelFile.SetCellStr(sheetName, "A1", "任务计划名称")
	_ = excelFile.SetCellStr(sheetName, "B1", "子任务名称")
	_ = excelFile.SetCellStr(sheetName, "C1", "任务计划类型")
	_ = excelFile.SetCellStr(sheetName, "D1", "任务时间")
	_ = excelFile.SetCellStr(sheetName, "E1", "完成时间")
	_ = excelFile.SetCellStr(sheetName, "F1", "导购姓名")
	_ = excelFile.SetCellStr(sheetName, "G1", "区域")
	_ = excelFile.SetCellStr(sheetName, "H1", "所属部门")
	_ = excelFile.SetCellStr(sheetName, "I1", "任务状态")
	_ = excelFile.SetCellStr(sheetName, "J1", "符合条件客户数（人）")
	_ = excelFile.SetCellStr(sheetName, "K1", "触达客户数（去重）")
	_ = excelFile.SetCellStr(sheetName, "L1", "是否在职")
	_ = excelFile.SetCellStr(sheetName, "M1", "任务完成率")
	startRow := 2
	cstSh, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		ginplus.ResError(c, err)
		cstSh = time.Local
	}
	for _, item := range data {
		// 0进行中，1完成，2未完成
		finishStr := ""
		CompletedTime := "         -"
		if item.IsFinished == 1 {
			finishStr = "完成"
			if item.CompletedTime > 0 {
				CompletedTime = time.Unix(item.CompletedTime, 0).Format("2006-01-02 15:04:05")
			}
		} else {
			finishStr = "未完成"
		}

		startTime := "         -"
		if item.StartTime > 0 {
			startTime = time.Unix(item.StartTime, 0).Format("2006-01-02 15:04:05")
		}
		var status string
		if item.WhetherToResign == 0 {
			status = "是"
		} else {
			status = "否"
		}
		categoryName := "--"
		switch item.Category {
		case 1:
			categoryName = "手动创建"
		case 2:
		case 3:
			categoryName = "任务模板创建"

		}
		rowStartStr := strconv.Itoa(startRow)
		_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.PlanName)
		_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, item.TaskName)
		_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, categoryName)
		_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, startTime)
		_ = excelFile.SetCellStr(sheetName, "E"+rowStartStr, CompletedTime)
		_ = excelFile.SetCellStr(sheetName, "F"+rowStartStr, item.UserName)
		_ = excelFile.SetCellStr(sheetName, "G"+rowStartStr, item.BelongingArea)
		_ = excelFile.SetCellStr(sheetName, "H"+rowStartStr, item.DepartmentName)
		_ = excelFile.SetCellStr(sheetName, "I"+rowStartStr, finishStr)
		_ = excelFile.SetCellStr(sheetName, "J"+rowStartStr, fmt.Sprintf("%v", item.CustomerNum))
		_ = excelFile.SetCellValue(sheetName, "K"+rowStartStr, item.TotalReachCustomerNumber)
		_ = excelFile.SetCellStr(sheetName, "L"+rowStartStr, status)
		_ = excelFile.SetCellStr(sheetName, "M"+rowStartStr, item.FinishRate)
		startRow = startRow + 1
	}
	fileName := "子任务用户画像列表数据导出-" + time.Now().In(cstSh).Format("200601021504") + ".xlsx"
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	des, _ := excelFile.WriteToBuffer()
	_, _ = c.Writer.Write(des.Bytes())
	return
}

func customerDetailDataExport(c *gin.Context) {
	params := planTask.ExportTasksRequst{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	DetailList, err := microscrm.CustomerDetailDataExport(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	jsonStr, _ := json.Marshal(DetailList)
	var data []*planTask.ExportDetailDataResp
	errJson := json.Unmarshal(jsonStr, &data)
	if errJson != nil {
		ginplus.ResError(c, err)
		return
	}
	excelFile := excelize.NewFile()
	sheetName := "子任务用户画像明细数据"
	// 创建sheet
	excelFile.SetSheetName("Sheet1", sheetName)
	// 写入表头
	_ = excelFile.SetCellStr(sheetName, "A1", "任务计划名称")
	_ = excelFile.SetCellStr(sheetName, "B1", "子任务名称")
	if params.TaskType == 2 {
		_ = excelFile.SetCellStr(sheetName, "C1", "顾客标签")
		_ = excelFile.SetCellStr(sheetName, "D1", "触达导购")
		_ = excelFile.SetCellStr(sheetName, "E1", "是否完成")
		_ = excelFile.SetCellStr(sheetName, "F1", "完成时间")
		_ = excelFile.SetCellStr(sheetName, "G1", "unionid")
		_ = excelFile.SetCellStr(sheetName, "H1", "会员卡号")
		_ = excelFile.SetCellStr(sheetName, "I1", "是否在职")
	} else {
		_ = excelFile.SetCellStr(sheetName, "C1", "触达导购")
		_ = excelFile.SetCellStr(sheetName, "D1", "是否完成")
		_ = excelFile.SetCellStr(sheetName, "E1", "完成时间")
		_ = excelFile.SetCellStr(sheetName, "F1", "unionid")
	}
	startRow := 2
	cstSh, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		ginplus.ResError(c, err)
		cstSh = time.Local
	}
	for _, item := range data {
		// 0进行中，1完成，2未完成
		finistStr := ""
		completedStr := "         -"
		if item.IsFinished == 1 {
			finistStr = "完成"
			if item.CompletedTime > 0 {
				completedStr = time.Unix(item.CompletedTime, 0).Format("2006-01-02 15:04:05")
			}
		} else {
			finistStr = "未完成"
		}
		var status string
		if item.WhetherToResign == 0 {
			status = "是"
		} else {
			status = "否"
		}

		rowStartStr := strconv.Itoa(startRow)
		_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.PlanName)
		_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, item.TaskName)
		if params.TaskType == 2 {
			_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, item.TagName)
			_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, item.BaName)
			_ = excelFile.SetCellStr(sheetName, "E"+rowStartStr, finistStr)
			_ = excelFile.SetCellStr(sheetName, "F"+rowStartStr, completedStr)
			_ = excelFile.SetCellStr(sheetName, "G"+rowStartStr, item.UnionID)
			_ = excelFile.SetCellStr(sheetName, "H"+rowStartStr, item.VipCode)
			_ = excelFile.SetCellStr(sheetName, "I"+rowStartStr, status)
		} else {
			_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, item.BaName)
			_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, finistStr)
			_ = excelFile.SetCellStr(sheetName, "E"+rowStartStr, completedStr)
			_ = excelFile.SetCellStr(sheetName, "F"+rowStartStr, item.UnionID)
		}

		startRow = startRow + 1
	}
	fileName := "子任务用户画像明细数据-" + time.Now().In(cstSh).Format("200601021504") + ".xlsx"
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	des, _ := excelFile.WriteToBuffer()
	_, _ = c.Writer.Write(des.Bytes())
	return
}

func groupChatData(c *gin.Context) {
	params := planTask.CustomerDataReq{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	versionInfo, err := microscrm.GroupChatTaskData(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, versionInfo)
}

func groupTaskListData(c *gin.Context) {
	params := planTask.TaskListRequst{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	versionInfo, err := microscrm.GroupChatTaskListData(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, versionInfo)
}

func customerTaskList(c *gin.Context) {
	params := planTask.TaskListRequst{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	versionInfo, err := microscrm.CustomerTaskListData(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, versionInfo)
}

func groupChatTaskDataExport(c *gin.Context) {
	params := planTask.ExportTasksRequst{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	groupList, err := microscrm.GroupChatTaskDataExport(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	jsonStr, _ := json.Marshal(groupList)
	var data []planTask.ExportGroupChatDataResp
	errJson := json.Unmarshal(jsonStr, &data)
	if errJson != nil {
		ginplus.ResError(c, err)
		return
	}
	excelFile := excelize.NewFile()
	sheetName := "子任务社群任务列表数据"
	// 创建sheet
	excelFile.SetSheetName("Sheet1", sheetName)
	// 写入表头
	_ = excelFile.SetCellStr(sheetName, "A1", "任务计划名称")
	_ = excelFile.SetCellStr(sheetName, "B1", "子任务名称")
	_ = excelFile.SetCellStr(sheetName, "C1", "任务时间")
	_ = excelFile.SetCellStr(sheetName, "D1", "完成时间")
	_ = excelFile.SetCellStr(sheetName, "E1", "社群名称")
	_ = excelFile.SetCellStr(sheetName, "F1", "社群群主")
	_ = excelFile.SetCellStr(sheetName, "G1", "区域")
	_ = excelFile.SetCellStr(sheetName, "H1", "所属部门")
	_ = excelFile.SetCellStr(sheetName, "I1", "任务状态")
	_ = excelFile.SetCellStr(sheetName, "J1", "触达客户数（去重）")
	startRow := 2
	cstSh, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		ginplus.ResError(c, err)
		cstSh = time.Local
	}
	for _, item := range data {
		// 0进行中，1完成，2未完成
		finishStr := ""
		CompletedTime := "         -"
		if item.IsFinished == 1 {
			finishStr = "完成"
			if item.CompletedTime > 0 {
				CompletedTime = time.Unix(item.CompletedTime, 0).Format("2006-01-02 15:04:05")
			}
		} else {
			finishStr = "未完成"
		}

		startTime := "         -"
		if item.StartTime > 0 {
			startTime = time.Unix(item.StartTime, 0).Format("2006-01-02 15:04:05")
		}

		rowStartStr := strconv.Itoa(startRow)

		_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.PlanName)
		_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, item.TaskName)
		_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, startTime)
		_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, CompletedTime)
		_ = excelFile.SetCellStr(sheetName, "E"+rowStartStr, item.GroupChatName)
		_ = excelFile.SetCellStr(sheetName, "F"+rowStartStr, item.UserName)
		_ = excelFile.SetCellStr(sheetName, "G"+rowStartStr, item.BelongingArea)
		_ = excelFile.SetCellStr(sheetName, "H"+rowStartStr, item.DepartmentName)
		_ = excelFile.SetCellStr(sheetName, "I"+rowStartStr, finishStr)

		_ = excelFile.SetCellValue(sheetName, "J"+rowStartStr, item.TotalReachCustomerNumber)
		startRow = startRow + 1
	}
	fileName := "子任务社群任务列表数据-" + time.Now().In(cstSh).Format("200601021504") + ".xlsx"
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	des, _ := excelFile.WriteToBuffer()
	_, _ = c.Writer.Write(des.Bytes())
	return
}

// 任务统计数据
func taskDataStat(c *gin.Context) {
	var request planTask.PlanStatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	if tokenInfo.EnterpriseId == 0 {
		ginplus.ResError(c, errors.New("企业ID为空！"))
		return
	}
	request.EnterpriseID = tokenInfo.EnterpriseId
	stat, err := microscrm.TaskDataStat(ctx.Wrap(c), &request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, &stat)
}

// 任务明细
func taskDataStatDetail(c *gin.Context) {
	var request planTask.StatDetailRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	if tokenInfo.EnterpriseId == 0 {
		ginplus.ResError(c, errors.New("企业ID为空！"))
		return
	}
	request.EnterpriseID = tokenInfo.EnterpriseId
	detail, err := microscrm.TaskDataStatDetail(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, &detail)
}

// 任务明细导出
func taskDetailStatExport(c *gin.Context) {
	var request planTask.StatDetailRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	if tokenInfo.EnterpriseId == 0 {
		ginplus.ResError(c, errors.New("企业ID为空！"))
		return
	}
	request.EnterpriseID = tokenInfo.EnterpriseId
	request.Size = 0
	request.Page = 0
	detail, err := microscrm.TaskDataStatDetail(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	jsonStr, _ := json.Marshal(detail)
	// 社群数据excel
	if request.TaskType == 1 {
		var data planTask.ResultCommunityDataDetail
		errJson := json.Unmarshal(jsonStr, &data)
		if errJson != nil {
			ginplus.ResError(c, err)
			return
		}
		excelFile := excelize.NewFile()
		sheetName := "社群明细数据"
		// 创建sheet
		excelFile.SetSheetName("Sheet1", sheetName)
		// 写入表头
		_ = excelFile.SetCellStr(sheetName, "A1", "子任务名称")
		_ = excelFile.SetCellStr(sheetName, "B1", "任务类型")
		_ = excelFile.SetCellStr(sheetName, "C1", "触达群主")
		_ = excelFile.SetCellStr(sheetName, "D1", "匹配社群数")
		_ = excelFile.SetCellStr(sheetName, "E1", "任务总数（件）")
		_ = excelFile.SetCellStr(sheetName, "F1", "已完成任务数（件）")
		_ = excelFile.SetCellStr(sheetName, "G1", "任务完成率（%）")

		startRow := 2
		cstSh, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			ginplus.ResError(c, err)
			cstSh = time.Local
		}
		for _, item := range data.Data {
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

			_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.SubtasksNumber)
			_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, "社群任务")
			_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, ReachLordNumber)
			_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, MatchGroupNumber)
			_ = excelFile.SetCellStr(sheetName, "E"+rowStartStr, TaskNumber)
			_ = excelFile.SetCellStr(sheetName, "F"+rowStartStr, TaskFinishNumber)
			_ = excelFile.SetCellStr(sheetName, "G"+rowStartStr, item.TaskFinishPercentage)

			startRow = startRow + 1
		}
		fileName := "任务明细数据导出-" + time.Now().In(cstSh).Format("200601021504") + ".xlsx"
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
		des, _ := excelFile.WriteToBuffer()
		_, _ = c.Writer.Write(des.Bytes())
		return
	}

	var data planTask.ResultCustomerPortraitDetail
	errJson := json.Unmarshal(jsonStr, &data)
	if errJson != nil {
		ginplus.ResError(c, err)
		return
	}
	excelFile := excelize.NewFile()
	sheetName := "客户画像明细数据"
	// 创建sheet
	excelFile.SetSheetName("Sheet1", sheetName)
	// 写入表头
	_ = excelFile.SetCellStr(sheetName, "A1", "子任务名称")
	_ = excelFile.SetCellStr(sheetName, "B1", "任务类型")
	_ = excelFile.SetCellStr(sheetName, "C1", "累计触达导购人数")
	_ = excelFile.SetCellStr(sheetName, "D1", "累计匹配客户数")
	_ = excelFile.SetCellStr(sheetName, "E1", "累计匹配客户数（导购去重）")
	_ = excelFile.SetCellStr(sheetName, "F1", "任务总数（件）")
	_ = excelFile.SetCellStr(sheetName, "G1", "已完成任务数（件）")

	startRow := 2
	cstSh, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		ginplus.ResError(c, err)
		cstSh = time.Local
	}
	for _, item := range data.Data {
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

		_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.SubtasksNumber)
		if request.TaskType == 2 {
			_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, "用户画像")
		}
		if request.TaskType == 3 {
			_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, "导入顾客名单")
		}
		if request.TaskType == 4 {
			_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, "导入顾客&导购匹配关系")
		}

		_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, ReachGuideNumber)
		_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, MatchCustomersNumber)
		_ = excelFile.SetCellStr(sheetName, "E"+rowStartStr, MatchCustomersNumber2)
		_ = excelFile.SetCellStr(sheetName, "F"+rowStartStr, TaskNumber)
		_ = excelFile.SetCellStr(sheetName, "G"+rowStartStr, TaskFinishNumber)

		startRow = startRow + 1
	}
	fileName := "任务明细数据导出-" + time.Now().In(cstSh).Format("200601021504") + ".xlsx"
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	des, _ := excelFile.WriteToBuffer()
	_, _ = c.Writer.Write(des.Bytes())
	return

}

// subTaskUploadUsers 子任务上传员工数据
func subTaskUploadUsers(c *gin.Context) {
	var params = make(map[string]string)
	params["enterprise_id"] = c.PostForm("enterprise_id")
	params["file_type"] = c.PostForm("file_type")
	fh, err := c.FormFile("file")
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	tk := c.Value(common.ContextKeyEnum.TokenInfo)
	if tk != nil {
		tokenInfo := tk.(*common.TokenInfo)
		params["enterprise_id"] = fmt.Sprintf("%+v", tokenInfo.EnterpriseId)
	}

	resp, err := microscrm.UploadSubTaskUsers(ctx.Wrap(c), params, "file", fh)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, resp)
	return
}

func zeroOmit2(info string) string {
	if info != "" {
		split := strings.Split(info, ".")
		if len(split) == 2 {
			if split[1] == "0" {
				return split[0]
			} else {
				return info
			}
		}
	}
	return "0"
}

// taskCSVFileUpload vsc标签文件上传
func taskCSVFileUpload(c *gin.Context) {
	// 上传文件流
	formFile, _ := c.FormFile("file")
	openFile, err := formFile.Open()
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var fileUrl string
	prefix := logic.GetBolbPrefix()
	prefix += "/private/plan-task"
	// 获取其他参数
	var dataMap = make(map[string]string, 0)
	d, _ := c.MultipartForm()
	for k, v := range d.Value {
		dataMap[k] = v[0]
	}
	for _, v2 := range d.File {
		dataMap["file_name"] = v2[0].Filename
	}
	if _, ok := dataMap["task_type"]; !ok {
		dataMap["task_type"] = "3"
	}
	taskType, _ := strconv.ParseInt(dataMap["task_type"], 10, 64)
	ext := filepath.Ext(formFile.Filename)
	if ext != ".xlsx" {
		ginplus.ResError(c, errors.New("只允许上传以xlsx的文件"))
		return
	}
	zr, err := zip.NewReader(openFile, formFile.Size)
	if err != nil {
		ginplus.ResError(c, errors.New("文件格式错误"))
		return
	}
	fileData, err := xlsx.ReadZipReader(zr)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	err = checkXlsxData(fileData.Sheets, int(taskType))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	// 上传数据到blob
	var fileName = fmt.Sprintf("%d", time.Now().Unix()) + "_taskConsumerFile" + ext
	// 获取enterpriseID
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	var filedetection delivery.FileDetection
	filedetection.FileName = dataMap["file_name"]
	filedetection.EnterpriseID = tokenInfo.EnterpriseId
	init := storage.Init(cos.COS)
	newFile, err := os.Create(fileName)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	// 保证utf-8不会中文数据乱码
	bytes, err := ioutil.ReadAll(openFile)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if _, err0 := newFile.Write(bytes); err0 != nil {
		ginplus.ResError(c, err0)
		return
	}
	if err := newFile.Close(); err != nil {
		ginplus.ResError(c, err)
		return
	}

	upload, err := init.Upload(ctx.Wrap(c), fileName, fileName, prefix)
	if err != nil {
		os.Remove(fileName)
		ginplus.ResError(c, err)
		return
	}
	fileUrl = init.GetURL(upload)
	if err := os.Remove(fileName); err != nil {
		ginplus.ResError(c, err)
		return
	}
	// 写入文件路径信息到数据库，获取id
	var tagDownloadInfo delivery.DownLoadFileInfo
	tagDownloadInfo.EnterpriseID = tokenInfo.EnterpriseId
	tagDownloadInfo.FileUrl = fileUrl
	tagDownloadInfo.FileName = dataMap["file_name"]
	tagDownloadInfo.TaskType = int(taskType)
	IDInfo, err := microscrm.PlanImportDownloadFileInfo(ctx.Wrap(c), &tagDownloadInfo)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var newID delivery.IDInfo
	marshal2, err := json.Marshal(IDInfo)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if err := json.Unmarshal(marshal2, &newID); err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, newID)
	return
}

// 检查文件内容
func checkXlsxData(sheets []*xlsx.Sheet, taskType int) error {
	riMin := 0
	if taskType != 3 {
		riMin = 1
	}
	for index, sheet := range sheets {
		if index > 0 {
			return errors.New("文件格式错误")
		}
		for ri, row := range sheet.Rows {

			if ri == riMin {
				if taskType == 3 {
					if row.Cells[0].String() != "unionid" || row.Cells[1].String() != "customerid" {
						return errors.New("上传的文件错误")
					}
				}
				if taskType == 4 {
					if row.Cells[0].String() != "Unionid ID" || row.Cells[1].String() != "Customer ID" || row.Cells[2].String() != "*SA Code" {
						return errors.New("上传的文件错误")
					}
				}
			}
			if len(row.Cells) == 0 {
				continue
			}
		}
	}
	return nil
}

func taskConsumerExport(c *gin.Context) {
	var request planTask.ConsumerList
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}

	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	if tokenInfo.EnterpriseId == 0 {
		ginplus.ResError(c, errors.New("企业ID为空！"))
		return
	}
	request.EnterpriseID = tokenInfo.EnterpriseId
	request.Size = 1000
	request.Page = 1
	request.Total = int64(0)
	list, err := microscrm.TaskConsumeList(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	jsonStr, _ := json.Marshal(list)
	var data planTask.ResultConsumerList
	errJson := json.Unmarshal(jsonStr, &data)
	if errJson != nil {
		ginplus.ResError(c, err)
		return
	}
	var datalist []planTask.ResultConsumerData
	for _, v := range data.List {
		datalist = append(datalist, v)
	}
	limitEnd := request.Page * request.Size
	if int(data.Total) > limitEnd {
		request.Total = data.Total
		maxPage := int(math.Ceil(float64(data.Total / int64(request.Size))))
		maxPage += 1
		for i := 2; i <= maxPage; i++ {
			request.Page = i
			newList, errs := GetConsumerDataExport(c, request)
			if errs != nil {
				ginplus.ResError(c, err)
				return
			}
			datalist = append(datalist, newList...)
		}
	}
	excelFile := excelize.NewFile()
	sheetName := "导入结果"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	// 设置宽度
	err = excelFile.SetColWidth(sheetName, "A", "D", 15)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	// 设置样式
	style, styleErr := excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"},"font":{"bold":true}}`)
	if nil != styleErr {
		ginplus.ResError(c, styleErr)
		return
	}
	err = excelFile.SetCellStyle(sheetName, "A1", "D1", style)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}

	// 写入表头
	if request.TaskType == 3 {
		_ = excelFile.SetCellStr(sheetName, "A1", "unionID")
		_ = excelFile.SetCellStr(sheetName, "B1", "导入结果")
		_ = excelFile.SetCellStr(sheetName, "C1", "专属SAID")
		_ = excelFile.SetCellStr(sheetName, "D1", "匹配SA姓名")
	}

	if request.TaskType == 4 {
		_ = excelFile.SetCellStr(sheetName, "A1", "unionID")
		_ = excelFile.SetCellStr(sheetName, "B1", "SaCode")
		_ = excelFile.SetCellStr(sheetName, "C1", "匹配结果")
	}

	startRow := 2
	// 填充数据
	style, err = excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"}}`)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	for _, item := range datalist {
		rowStartStr := strconv.Itoa(startRow)
		if request.TaskType == 3 {
			_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.UnionId)
			_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, item.ImportRes)
			_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, item.Said)
			_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, item.Name)
		}
		if request.TaskType == 4 {
			_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.UnionId)
			_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, item.Said)
			_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, item.ImportRes)
		}
		startRow = startRow + 1
		err = excelFile.SetCellStyle(sheetName, "A"+rowStartStr, "P"+rowStartStr, style)
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

func GetConsumerDataExport(c *gin.Context, request planTask.ConsumerList) (data []planTask.ResultConsumerData, err error) {
	list, err := microscrm.TaskConsumeList(ctx.Wrap(c), request)
	if err != nil {
		return nil, err
	}
	var newList planTask.ResultConsumerList
	byteinfo, errJson := json.Marshal(list)
	if errJson != nil {
		return nil, errJson
	}
	if errs := json.Unmarshal(byteinfo, &newList); err != nil {
		return nil, errs
	}
	return newList.List, nil
}
