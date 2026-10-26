package display

import (
	"archive/zip"
	"bytes"
	"errors"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	display2 "micro-scrm/internal/app/dao/display"
	membersDao "micro-scrm/internal/app/dao/members"
	"micro-scrm/internal/app/http/request/display"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/tealeg/xlsx"
	"github.com/xuri/excelize/v2"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

type TaskSvc struct {
	service.Base
	taskDao *display2.TaskDao
}

func NewTaskSvc() *TaskSvc {
	return &TaskSvc{
		taskDao: display2.NewTaskDao(),
	}
}

func (s TaskSvc) CreateTask(c ctx.Context, request display.AddTaskReq) (int64, error) {
	if request.StartTime >= request.EndTime {
		return 0, errors.New("任务的开始时间不能大于结束时间")
	}
	id := app.UID()
	nowTime := time.Now().Unix()
	data := microScrmModel.DisplayTaskModel{
		BaseModel: models.BaseModel{
			ID:        id,
			CreatedAt: nowTime,
			UpdatedAt: nowTime,
		},
		DisplayName:  request.DisplayName,
		StoreType:    request.StoreType,
		Content:      request.Content,
		ClassURL:     helper.MarshalToString(request.ClassURL),
		TemplateInfo: helper.MarshalToString(request.TemplateInfo),
		ShowType:     request.ShowType,
		Status:       request.Status,
		StartTime:    request.StartTime,
		EndTime:      request.EndTime,
		ReviewType:   request.ReviewType,
		ReviewFile:   request.ReviewFile,
	}
	err := s.taskDao.CreateData(c, request.EnterpriseId, data)
	if err != nil {
		return id, err
	}
	if request.StoreType == microScrmModel.VMStoreTypeLimit { // 插入门店关联数据
		err = s.taskDao.CreateDisplayTaskRelation(c, request.EnterpriseId, id, request.StoreList)
	}
	return id, err

}

func (s TaskSvc) List(c ctx.Context, request display.SearchReq) (int64, []display.TaskInfoRes, error) {
	result := make([]display.TaskInfoRes, 0)
	total, res, err := s.taskDao.GetList(c, request)
	if err != nil {
		return 0, result, err
	}
	for _, item := range res {
		info := display.TaskInfoRes{
			ID:          item.ID,
			DisplayName: item.DisplayName,
			StoreType:   item.StoreType,
			Content:     item.Content,
			ShowType:    item.ShowType,
			Status:      item.Status,
			StartTime:   item.StartTime,
			EndTime:     item.EndTime,
			CreatedAt:   item.CreatedAt,
		}
		result = append(result, info)
	}
	return total, result, err
}

// Detail VM任务详情
func (s TaskSvc) Detail(c ctx.Context, request display.TaskBaseReq) (display.TaskInfoRes, error) {
	var res display.TaskInfoRes
	info, err := s.taskDao.GetOne(c, request.EnterpriseId, map[string]interface{}{
		"id":         request.ID,
		"deleted_at": 0,
	})
	if err != nil {
		return res, err
	}
	res = display.TaskInfoRes{
		ID:          info.ID,
		DisplayName: info.DisplayName,
		StoreType:   info.StoreType,
		Content:     info.Content,
		ShowType:    info.ShowType,
		Status:      info.Status,
		StartTime:   info.StartTime,
		EndTime:     info.EndTime,
		CreatedAt:   info.CreatedAt,
		ReviewType:  info.ReviewType,
		ReviewFile:  info.ReviewFile,
	}
	helper.Unmarshal(info.ClassURL, &res.ClassURL)
	helper.Unmarshal(info.TemplateInfo, &res.TemplateInfo)
	res.StoreList = []display.StoreDepartmentInfo{}
	if info.StoreType == microScrmModel.VMStoreTypeLimit {
		storeList, err := s.taskDao.GetDisplayStore(c, request.EnterpriseId, request.ID)
		if err != nil {
			return res, err
		}
		for _, storeInfo := range storeList {
			item := display.StoreDepartmentInfo{
				ThirdID:     storeInfo.StoreID,
				StoreType:   int(storeInfo.StoreType),
				DpStoreName: storeInfo.StoreName,
				DpThirdID:   storeInfo.StoreID,
			}
			res.StoreList = append(res.StoreList, item)
		}
	}
	return res, err
}

// Copy VM任务复制
func (s TaskSvc) Copy(c ctx.Context, request display.TaskBaseReq) (int64, error) {
	copyId := app.UID()
	info, err := s.taskDao.GetOne(c, request.EnterpriseId, map[string]interface{}{
		"id":         request.ID,
		"deleted_at": 0,
	})
	if err != nil {
		return copyId, err
	}
	info.ID = copyId
	info.StartTime = 0
	info.EndTime = 0
	info.CreatedAt = time.Now().Unix()
	info.UpdatedAt = time.Now().Unix()
	info.Status = microScrmModel.DisplayStatus
	info.ShowType = microScrmModel.ShowStatusDown
	info.DisplayName = info.DisplayName + "-副本"
	err = s.taskDao.CreateData(c, request.EnterpriseId, info)
	if err != nil {
		return copyId, err
	}

	if info.StoreType == microScrmModel.VMStoreTypeLimit {
		var createStoreList []display.StoreDepartmentInfo
		storeList, err := s.taskDao.GetDisplayStore(c, request.EnterpriseId, request.ID)
		if err != nil {
			return copyId, err
		}
		for _, storeInfo := range storeList {
			item := display.StoreDepartmentInfo{
				DepartmentID: int(storeInfo.DepartmentId),
				ThirdID:      storeInfo.StoreID,
				StoreType:    int(storeInfo.StoreType),
				DpThirdID:    storeInfo.StoreID,
				DpStoreName:  storeInfo.StoreName,
			}

			createStoreList = append(createStoreList, item)
		}
		err = s.taskDao.CreateDisplayTaskRelation(c, request.EnterpriseId, info.ID, createStoreList)
	}
	return copyId, err
}

// UpdateShow VM任务上下架
func (s TaskSvc) UpdateShow(c ctx.Context, request display.TaskBaseReq) error {
	one, err := s.taskDao.GetOne(c, request.EnterpriseId, map[string]interface{}{
		"id":         request.ID,
		"deleted_at": 0,
	})
	if err != nil {
		return err
	}
	showType := helper.RuleThree(one.ShowType == microScrmModel.ShowStatusDown, microScrmModel.ShowStatusUp, microScrmModel.ShowStatusDown)
	err = s.taskDao.UpdateData(c, request.EnterpriseId, map[string]interface{}{
		"id":         request.ID,
		"deleted_at": 0,
	}, map[string]interface{}{
		"show_type": showType,
	})
	return err

}

// UpdateTask 编辑任务
func (s TaskSvc) UpdateTask(c ctx.Context, request display.AddTaskReq) error {
	taskInfo, err := s.taskDao.GetOne(c, request.EnterpriseId, map[string]interface{}{
		"id": request.ID,
	})
	if err != nil {
		return err
	}
	if request.StartTime <= int(time.Now().Unix()) {
		return errors.New("开始时间不能小于当前时间")
	}
	if taskInfo.Status >= microScrmModel.DisplayStatusIng {
		return errors.New("任务已开始不能编辑")
	}

	updateMap := map[string]interface{}{
		"display_name":  request.DisplayName,
		"store_type":    request.StoreType,
		"content":       request.Content,
		"class_url":     helper.MarshalToString(request.ClassURL),
		"template_info": helper.MarshalToString(request.TemplateInfo),
		"show_type":     request.ShowType,
		"start_time":    request.StartTime,
		"end_time":      request.EndTime,
		"review_type":   request.ReviewType,
		"review_file":   request.ReviewFile,
		"updated_at":    time.Now().Unix(),
	}
	err = s.taskDao.UpdateData(c, request.EnterpriseId, map[string]interface{}{
		"id": request.ID,
	}, updateMap)
	if err != nil {
		return err
	}
	if request.StoreType == microScrmModel.VMStoreTypeLimit {
		// 整删，整增
		err = s.taskDao.DeleteDisplayStore(c, request.EnterpriseId, map[string]interface{}{"display_task_id": request.ID})
		err = s.taskDao.CreateDisplayTaskRelation(c, request.EnterpriseId, request.ID, request.StoreList)
	}
	err = display2.NewAllotDao().DeleteData(c, request.EnterpriseId, map[string]interface{}{"display_task_id": request.ID})
	return err
}

func (s TaskSvc) GetAllAllotList(c ctx.Context, enterpriseId int64) ([]display.SimpleTaskInfo, error) {
	result := make([]display.SimpleTaskInfo, 0)
	many, err := s.taskDao.GetMany(c, enterpriseId, map[string]interface{}{
		"show_type":  microScrmModel.ShowStatusUp,
		"deleted_at": 0,
	})
	if err != nil {
		return result, err
	}
	sort.Slice(many, func(i, j int) bool {
		return many[i].ID > many[j].ID
	})
	for _, model := range many {
		info := display.SimpleTaskInfo{
			ID:          model.ID,
			DisplayName: model.DisplayName,
		}
		result = append(result, info)
	}
	return result, err

}

func (s TaskSvc) GetStoreTaskList(c ctx.Context, enterpriseId int64, storeId, finishStatus int, taskIds []int) ([]display.StoreListRes, error) {
	res := make([]display.StoreListRes, 0)
	list, err := s.taskDao.GetStoreTaskList(c, enterpriseId, storeId, finishStatus, taskIds)
	if err != nil {
		return res, err
	}
	for _, value := range list {
		if len(value.TaskAllotInfo) > 0 {
			allotInfo := value.TaskAllotInfo[0]
			info := display.StoreListRes{
				ID:           value.ID,
				DisplayName:  value.DisplayName,
				StoreType:    value.StoreType,
				Content:      value.Content,
				ShowType:     value.ShowType,
				Status:       value.Status,
				StartTime:    value.StartTime,
				EndTime:      value.EndTime,
				AllotId:      allotInfo.ID,
				StoreId:      allotInfo.StoreID,
				FinishStatus: allotInfo.FinishStatus,
			}
			res = append(res, info)
		}

	}
	return res, err
}

func (s TaskSvc) CheckReviewFile(c ctx.Context, enterpriseId int64, fileUrl string) (display.ReviewImportResp, error) {
	fileData, _, result, err := s.ReadReviewFile(c, enterpriseId, fileUrl)
	if err != nil {
		return result, err
	}
	if result.FailedNum > 0 { // 如果有错误的销售单号，则全都不更新转货单，展示错误信息
		fileName := "定向分配名单" + time.Now().Format("20060102150405") + ".xlsx"
		excelFile := excelize.NewFile()
		sheetName := "定向分配名单错误列表"
		// 创建sheet
		excelFile.SetSheetName("Sheet1", sheetName)
		err = excelFile.SaveAs(fileName)
		defer excelFile.Close()
		if err != nil {
			return result, err
		}
		styleID, err := excelFile.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Color: "#FF0000", // 红色
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
		})
		if err != nil {
			return result, err
		}
		stream, err := excelFile.NewStreamWriter(sheetName)
		if err != nil {
			return result, err
		}
		row := make([]any, 0)
		row = []any{
			"员工姓名", "导购工号", "门店code", "失败数据原因",
		}
		err = stream.SetRow("A1", row)
		startRow := int64(2)
		for _, item := range fileData {
			rowStartStr := strconv.Itoa(int(atomic.LoadInt64(&startRow)))
			row = []any{item.Name, item.UserId, item.StoreCode, excelize.Cell{
				StyleID: styleID,
				Value:   item.FailReason,
			}}
			err = stream.SetRow("A"+rowStartStr, row)
			atomic.AddInt64(&startRow, 1)
		}
		if err != nil {
			return result, err
		}
		err = stream.Flush()
		if err != nil {
			return result, err
		}
		err = excelFile.SaveAs(fileName)
		if err != nil {
			return result, err
		}
		failUrl, err := helper.UploadFileToCos(c, fileName, false)
		if err != nil {
			return result, err
		}
		result.FailUrl = failUrl
		return result, err
	}
	return result, err
}

// ReadReviewFile 读取定向分配的名单文件
func (s TaskSvc) ReadReviewFile(c ctx.Context, enterpriseId int64, fileUrl string) ([]display.ReviewData, []display.ReviewImportData, display.ReviewImportResp, error) {

	var data []display.ReviewData
	var importData []display.ReviewImportData
	var res display.ReviewImportResp
	specialUserIds, err := s.taskDao.GetSpecialList(c, enterpriseId, map[string]interface{}{
		"is_special": 1,
		"deleted_at": 0,
	})
	if err != nil {
		return data, importData, res, err
	}
	list, err := membersDao.NewMembers().GetMemberInfoAndDepartmentsByUserId(c, enterpriseId, specialUserIds)
	if err != nil {
		return data, importData, res, err
	}
	specialMap := make(map[string]string, 0)
	for _, v := range list {
		specialMap[v.MembersUserid] = v.DepartmentName
	}
	departMents, err := departmentsDao.NewDepartments().GetMany(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"deleted_at":    0,
	})
	if err != nil {
		return data, importData, res, nil
	}
	departmentCode := make([]string, 0)
	for _, ment := range departMents {
		departmentCode = append(departmentCode, ment.Name)
	}
	if fileUrl == "" { // 如果是空地址则直接返回
		return data, importData, res, nil
	}
	b, err := helper.ReadRemoteFile(c, fileUrl)
	if err != nil {
		return data, nil, res, err
	}

	zr, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return data, nil, res, err
	}
	fileData, err := xlsx.ReadZipReader(zr)
	if err != nil {
		return data, nil, res, err
	}
	riMin := 1
	storeCodes := make([]string, 0)
	for index, sheet := range fileData.Sheets {
		if index > 0 {
			return nil, nil, res, errors.New("文件格式错误")
		}
		for ri, row := range sheet.Rows {
			if ri < riMin || len(row.Cells) == 0 {
				continue
			}
			if len(row.Cells) != 3 {
				return nil, nil, res, errors.New("模板文件错误")
			}
			res.UploadNum++
			failReason := ""
			name := strings.TrimSpace(row.Cells[0].String())
			userId := strings.TrimSpace(row.Cells[1].String())
			storeCode := strings.TrimSpace(row.Cells[2].String())
			switch {
			case userId == "" && storeCode == "":
				res.FailedNum++
				failReason = "无效的数据"
			case userId == "":
				res.FailedNum++
				failReason = "未匹配到员工"
			case storeCode == "":
				res.FailedNum++
				failReason = "未匹配到门店信息"
			case !helper.InArray(storeCode, departmentCode):
				res.FailedNum++
				failReason = "不存在的门店code"
			case helper.InArray(storeCode, storeCodes):
				res.DuplicateNum++
				failReason = "该门店被重复分配"
			case !helper.InArray(userId, specialUserIds):
				res.FailedNum++
				failReason = "该员工不是special人员"
			case specialMap[userId] == storeCode:
				res.FailedNum++
				failReason = "special人员不能审核自己门店"
			default:
				storeCodes = append(storeCodes, storeCode)
				res.Success++
				failReason = ""
			}
			data = append(data, display.ReviewData{
				ReviewImportData: display.ReviewImportData{
					Name:      name,
					UserId:    userId,
					StoreCode: storeCode,
				},
				FailReason: failReason,
			})
			if failReason == "" {
				importData = append(importData, display.ReviewImportData{
					Name:      name,
					UserId:    userId,
					StoreCode: storeCode,
				})
			}
		}
	}
	if len(data) == 0 {
		return nil, nil, res, errors.New("模板内容为空")
	}
	return data, importData, res, nil
}
