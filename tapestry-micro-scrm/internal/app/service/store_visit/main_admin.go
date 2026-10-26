package storeVisitSvc

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	dashboardDao "micro-scrm/internal/app/dao/dashboard"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	storeVisitDao "micro-scrm/internal/app/dao/store_visit"
	"micro-scrm/internal/app/http/request"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	storeVisitReq "micro-scrm/internal/app/http/request/store_visit"
	storeVisitResp "micro-scrm/internal/app/http/response/store_visit"
	dashboardModel "micro-scrm/internal/app/models/dashboard"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	storeVisitModel "micro-scrm/internal/app/models/store_visit"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	"sort"
	"strconv"
	"sync/atomic"
	"time"
)

// MainAdminSvc 巡店管理后台的业务服务
type MainAdminSvc struct {
	service.Base
	visitMainDao   *storeVisitDao.MainDao
	visitMemberDao *storeVisitDao.MemberRelationDao
}

func NewMainAdminSvc() *MainAdminSvc {
	return &MainAdminSvc{
		visitMainDao:   storeVisitDao.NewMainDao(),
		visitMemberDao: storeVisitDao.NewMemberRelationDao(),
	}
}

func (s *MainAdminSvc) AdminStoreVisitList(c ctx.Context, req storeVisitReq.AdminStoreVisitListReq) ([]storeVisitResp.AdminStoreVisitInfo, int64, error) {
	var result = make([]storeVisitResp.AdminStoreVisitInfo, 0)
	thirdIds := make([]int, 0)
	if req.StoreName != "" {
		storeList, err := microWorkDao.NewStoreDepartment().GetMany(c, "enterprise_id = ? AND (store_name like ? or third_store_code like ?)", req.EnterpriseId, "%"+req.StoreName+"%", "%"+req.StoreName+"%")
		if err != nil {
			return nil, 0, err
		}
		if len(storeList) == 0 {
			return result, 0, nil
		}
		for _, store := range storeList {
			thirdIds = append(thirdIds, int(store.WechatDepartmentID))
		}
	}
	visitIds := make([]int64, 0)
	if req.DeptName != "" {

		memberIds, err := storeVisitDao.NewDeptMemberDao().GetMemberIdsByDeptName(c, req.EnterpriseId, req.DeptName)
		if err != nil {
			return nil, 0, err
		}
		if len(memberIds) > 0 {
			res, err := s.visitMemberDao.GetVisitIdsByMemberIds(c, req.EnterpriseId, memberIds)
			if err != nil || len(res) == 0 {
				return result, 0, err
			}
			visitIds = res
		}

	}
	visitList, total, err := s.visitMainDao.GetAdminStoreVisitList(c, req, thirdIds, visitIds)
	if err != nil || len(visitList) == 0 {
		return result, 0, err
	}
	storeThirdIds := make([]int, 0)
	for _, visit := range visitList {
		storeThirdIds = append(storeThirdIds, visit.ThirdId)
	}
	storeList, err := microWorkDao.NewStoreDepartment().GetStoreName(c, req.EnterpriseId, storeThirdIds)
	if err != nil {
		return nil, 0, err
	}
	var storeMap = make(map[int]microWorkModel.StoreDepartmentModel)
	for _, store := range storeList {
		storeMap[int(store.WechatDepartmentID)] = store
	}
	for _, row := range visitList {
		status := row.Status
		if row.DeletedAt > 0 {
			status = storeVisitModel.VisitClose
		}
		tmp := storeVisitResp.AdminStoreVisitInfo{
			ID:        row.ID,
			Date:      row.Date,
			TimeSpan:  row.TimeSpan,
			Status:    status,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}
		//门店信息
		if store, ok := storeMap[row.ThirdId]; ok {
			tmp.StoreCode = store.ThirdStoreCode
			tmp.StoreName = store.StoreName
		}
		result = append(result, tmp)
	}

	return result, total, nil
}

func (s *MainAdminSvc) AdminStoreVisitClose(c ctx.Context, req storeVisitReq.VisitInfoReq) error {
	return s.visitMainDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
		"id": req.ID,
	}, map[string]interface{}{
		"deleted_at": time.Now().Unix(),
		"updated_at": time.Now().Unix(),
		"status":     storeVisitModel.VisitClose,
	})
}

func (s *MainAdminSvc) AdminStoreVisitExport(c ctx.Context, req storeVisitReq.AdminStoreVisitListReq) (string, error) {
	thirdIds := make([]int, 0)
	visitIds := make([]int64, 0)
	if req.StoreName != "" {
		storeList, err := microWorkDao.NewStoreDepartment().GetMany(c, "enterprise_id = ? AND (store_name like ? or third_store_code like ?)", req.EnterpriseId, "%"+req.StoreName+"%", "%"+req.StoreName+"%")
		if err != nil {
			return "", err
		}
		if len(storeList) == 0 {
			return "", err
		}
		for _, store := range storeList {
			thirdIds = append(thirdIds, int(store.WechatDepartmentID))
		}
	}
	if req.DeptName != "" {
		memberIds, err := storeVisitDao.NewDeptMemberDao().GetMemberIdsByDeptName(c, req.EnterpriseId, req.DeptName)
		if err != nil {
			return "", err
		}
		if len(memberIds) > 0 {
			res, err := s.visitMemberDao.GetVisitIdsByMemberIds(c, req.EnterpriseId, memberIds)
			if err != nil || len(res) == 0 {
				return "", err
			}
			visitIds = res
		}
	}
	storeList, err := microWorkDao.NewStoreDepartment().GetMany(c, "enterprise_id=?", req.EnterpriseId)
	if err != nil {
		return "", err
	}
	var storeMap = make(map[int]microWorkModel.StoreDepartmentModel)
	for _, store := range storeList {
		storeMap[int(store.WechatDepartmentID)] = store
	}

	fileName := "巡店列表明细" + time.Now().Format("20060102150405") + ".xlsx"
	excelFile := excelize.NewFile()
	sheetName := "巡店列表明细"
	// 创建sheet
	excelFile.SetSheetName("Sheet1", sheetName)
	err = excelFile.SaveAs(fileName)
	defer excelFile.Close()
	if err != nil {
		return fileName, err
	}
	stream, err := excelFile.NewStreamWriter(sheetName)
	if err != nil {
		return fileName, err
	}
	row := make([]any, 0)
	row = []any{
		"巡店记录ID", "巡店门店", "门店编号", "巡店时间", "创建时间", "状态", "完成时间",
	}
	err = stream.SetRow("A1", row)
	if err != nil {
		return fileName, err
	}
	startRow := int64(2)
	wg := threading.RoutineGroup{}
	StatusMap := map[int]string{
		storeVisitModel.VisitStatusPlan: "巡店预约",
		storeVisitModel.VisitVisiting:   "巡店中",
		storeVisitModel.VisitComplete:   "巡店完成",
		storeVisitModel.VisitClose:      "关闭",
	}
	req.Page = 0
	req.Size = 10000
	for {
		req.Page++
		// 获取数据
		list, _, err := s.visitMainDao.GetAdminStoreVisitList(c, req, thirdIds, visitIds)
		if err != nil || len(list) == 0 {
			break
		}
		wg.RunSafe(func() {
			for _, item := range list {
				rowStartStr := strconv.Itoa(int(atomic.LoadInt64(&startRow)))
				var storeCode string
				var storeName string
				if store, ok := storeMap[item.ThirdId]; ok {
					storeCode = store.ThirdStoreCode
					storeName = store.StoreName
				}
				createTime := helper.GetFormatDateByParams(item.CreatedAt)
				updateTime := helper.GetFormatDateByParams(item.UpdatedAt)
				finishTIme := "--"
				status := item.Status
				if item.DeletedAt > 0 {
					status = storeVisitModel.VisitClose
				}
				if status >= storeVisitModel.VisitComplete {
					finishTIme = updateTime
				}
				statusStr := StatusMap[status]
				row = []any{
					strconv.Itoa(int(item.ID)) + "\t", storeName, storeCode, fmt.Sprintf("%s  %s", item.Date, item.TimeSpan), createTime, statusStr, finishTIme,
				}
				err = stream.SetRow("A"+rowStartStr, row)
				atomic.AddInt64(&startRow, 1)
			}
		})

	}
	wg.Wait()
	err = stream.Flush()
	if err != nil {
		return fileName, err
	}
	err = excelFile.SaveAs(fileName)
	if err != nil {
		return fileName, err
	}
	return fileName, err
}

func (s *MainAdminSvc) AdminDashboard(c ctx.Context, req storeVisitReq.AdminStoreDashboardReq, isExport bool) ([]storeVisitResp.AdminDashboardInfo, int64, int64, error) {
	thirdIds := make([]int, 0)
	memberIds := make([]int64, 0)
	result := make([]storeVisitResp.AdminDashboardInfo, 0)
	var total int64
	var sumVisitNum int64
	if req.AreaId > 0 {
		res, err := dashboardDao.NewArea().GetThirdIdsByAreaIds(c, dashboardReq.ManagerAuthAreaIds{
			EnterpriseBaseReq: request.EnterpriseBaseReq{
				EnterpriseId: req.EnterpriseId,
			},
			AreaIds: []int64{req.AreaId},
		})
		if err != nil || len(res) == 0 {
			return result, 0, 0, err
		}
		thirdIds = append(thirdIds, res...)
	}
	if req.DeptName != "" {
		res, err := storeVisitDao.NewDeptMemberDao().GetMemberIdsByDeptName(c, req.EnterpriseId, req.DeptName)
		if err != nil {
			return nil, 0, 0, err
		}
		memberIds = append(memberIds, res...)
	}

	if req.Type == storeVisitReq.DashboardStore {
		many, err := storeVisitDao.NewDeptMemberDao().GetMany(c, req.EnterpriseId, map[string]interface{}{
			"deleted_at": 0,
		})
		if err != nil {
			return result, 0, 0, err
		}
		var memberDeptMap = make(map[int64]string)
		for _, item := range many {
			memberDeptMap[item.MemberID] = item.DeptName
		}
		storeList, err := microWorkDao.NewStoreDepartment().GetMany(c, "enterprise_id=?", req.EnterpriseId)
		if err != nil {
			return result, 0, 0, err
		}
		var storeMap = make(map[int]microWorkModel.StoreDepartmentModel)
		for _, store := range storeList {
			storeMap[int(store.WechatDepartmentID)] = store
		}
		list, count, err := s.visitMainDao.GetStoreVisitDashboardByStore(c, req, thirdIds, memberIds, false)
		if err != nil {
			return result, 0, 0, err
		}
		for _, item := range list {
			var storeCode string
			var storeName string
			if store, ok := storeMap[item.ThirdId]; ok {
				storeCode = store.ThirdStoreCode
				storeName = store.StoreName
			}
			deptName := memberDeptMap[item.MemberId]
			rate := 0.00
			if item.FinishNum > 0 && item.VisitNum > 0 {
				rate = float64(item.FinishNum) / float64(item.VisitNum) * 100
			}
			finishRate := fmt.Sprintf("%.2f%%", rate)
			result = append(result, storeVisitResp.AdminDashboardInfo{
				Name:      storeName,
				StoreCode: storeCode,
				VisitNum:  item.VisitNum,
				DeptName:  deptName,
				VisitRate: finishRate,
			})
		}
		total = count
		if !isExport && req.Page == 1 {
			num, err := s.visitMainDao.GetStoreVisitDashboardByStoreVisitNum(c, req, thirdIds, memberIds)
			if err != nil {
				return result, 0, 0, err
			}
			sumVisitNum = num
		}
	} else {
		if req.Page != 1 {
			return result, 0, 0, nil
		}
		areaType := req.Type - 1
		reginMap := make(map[int]string)
		if req.Type == storeVisitReq.DashboardArea {
			regionList, err := dashboardDao.NewArea().GetMany(c, req.EnterpriseId, map[string]interface{}{
				"type":       dashboardModel.AreaRegion,
				"is_virtual": 0,
				"closed_at":  0,
			})
			if err != nil {
				return result, 0, 0, err
			}
			for _, region := range regionList {
				reginMap[int(region.ID)] = region.Name
			}
			if req.AreaId > 0 { // 如果是小区汇总，并且筛选条件有了分区，则只查询当前分区数据
				res, _, err := s.visitMainDao.GetStoreVisitDashboardByStore(c, req, thirdIds, memberIds, true)
				if err != nil {
					return result, 0, 0, err
				}
				var visitNum int64
				var finishNum int64
				rate := 0.00
				for _, v := range res {
					sumVisitNum += v.VisitNum
					visitNum += v.VisitNum
					finishNum += v.FinishNum
				}
				if finishNum > 0 && visitNum > 0 {
					rate = float64(finishNum) / float64(visitNum) * 100
				}
				finishRate := fmt.Sprintf("%.2f%%", rate)
				areaInfo, err := dashboardDao.NewArea().GetOne(c, req.EnterpriseId, map[string]interface{}{
					"id": req.AreaId,
				})
				if err != nil {
					return result, 0, 0, err
				}
				result = append(result, storeVisitResp.AdminDashboardInfo{
					Name:       areaInfo.Name,
					ParentName: reginMap[areaInfo.ParentID],
					VisitNum:   visitNum,
					VisitRate:  finishRate,
				})
				total = int64(len(result))
				return result, total, sumVisitNum, nil
			}
		}
		many, err := dashboardDao.NewArea().GetMany(c, req.EnterpriseId, map[string]interface{}{
			"type":       areaType,
			"is_virtual": 0,
			"closed_at":  0,
		})
		if err != nil {
			return result, 0, 0, err
		}
		resultThirdIdMap := make(map[int64][]int)
		resultMap := make(map[int64]struct {
			VisitNum   int64
			FinishNum  int64
			Name       string
			ParentName string
		})
		allThirdIds := make([]int, 0)
		for _, item := range many {
			resAreaIds, err := dashboardDao.NewArea().GetThirdIdsByAreaIds(c, dashboardReq.ManagerAuthAreaIds{
				EnterpriseBaseReq: request.EnterpriseBaseReq{
					EnterpriseId: req.EnterpriseId,
				},
				AreaIds: []int64{item.ID},
			})
			if err != nil {
				return result, 0, 0, err
			}
			if len(thirdIds) > 0 {
				isContinue := false
				for _, id := range resAreaIds { // 如果这个区域下的所有门店都不在查询的小区门店下，则不展示这个区域的数据
					if helper.InArray(id, thirdIds) {
						isContinue = true
					}
				}
				if !isContinue {
					continue
				}
			}
			resultThirdIdMap[item.ID] = resAreaIds
			allThirdIds = append(allThirdIds, resAreaIds...)
			reginName := ""
			if req.Type == storeVisitReq.DashboardArea {
				reginName = reginMap[item.ParentID]
			}
			resultMap[item.ID] = struct {
				VisitNum   int64
				FinishNum  int64
				Name       string
				ParentName string
			}{VisitNum: 0, FinishNum: 0, Name: item.Name, ParentName: reginName}

		}
		interSectIds := allThirdIds
		if len(thirdIds) > 0 {
			interSectIds = helper.ArrayIntersect(thirdIds, allThirdIds)
		}
		res, _, err := s.visitMainDao.GetStoreVisitDashboardByStore(c, req, interSectIds, memberIds, true)
		if err != nil {
			return result, 0, 0, err
		}
		for _, item := range res {
			for k, v := range resultThirdIdMap {
				if helper.InArray(item.ThirdId, v) {
					visitNum := resultMap[k].VisitNum
					finishNum := resultMap[k].FinishNum
					sumVisitNum += item.VisitNum
					visitNum += item.VisitNum
					finishNum += item.FinishNum
					resultMap[k] = struct {
						VisitNum   int64
						FinishNum  int64
						Name       string
						ParentName string
					}{VisitNum: visitNum, FinishNum: finishNum, Name: resultMap[k].Name, ParentName: resultMap[k].ParentName}
				}
			}
		}
		for _, item := range resultMap {
			rate := 0.00
			if item.FinishNum > 0 && item.VisitNum > 0 {
				rate = float64(item.FinishNum) / float64(item.VisitNum) * 100
			}
			finishRate := fmt.Sprintf("%.2f%%", rate)

			result = append(result, storeVisitResp.AdminDashboardInfo{
				Name:       item.Name,
				ParentName: item.ParentName,
				VisitNum:   item.VisitNum,
				VisitRate:  finishRate,
			})

		}
		sort.Slice(result, func(i, j int) bool {
			return result[i].VisitNum > result[j].VisitNum
		})
		total = int64(len(result))
	}

	return result, total, sumVisitNum, nil
}

func (s *MainAdminSvc) AdminDashboardExport(c ctx.Context, req storeVisitReq.AdminStoreDashboardReq) (string, error) {
	var err error
	fileName := "巡店统计明细" + time.Now().Format("20060102150405") + ".xlsx"
	excelFile := excelize.NewFile()
	sheetName := "巡店统计明细"
	// 创建sheet
	excelFile.SetSheetName("Sheet1", sheetName)
	err = excelFile.SaveAs(fileName)
	defer excelFile.Close()
	if err != nil {
		return fileName, err
	}
	stream, err := excelFile.NewStreamWriter(sheetName)
	if err != nil {
		return fileName, err
	}
	row := make([]any, 0)
	switch req.Type {
	case storeVisitReq.DashboardStore:
		row = []any{
			"门店名称", "门店编码", "部门", "巡店次数", "巡店完成率",
		}
	case storeVisitReq.DashboardSource:
		row = []any{
			"渠道", "巡店次数", "巡店完成率",
		}
	case storeVisitReq.DashboardRegion:
		row = []any{
			"大区名称", "巡店次数", "巡店完成率",
		}
	case storeVisitReq.DashboardArea:
		row = []any{
			"小区名称", "大区名称", "巡店次数", "巡店完成率",
		}
	}
	err = stream.SetRow("A1", row)
	if err != nil {
		return fileName, err
	}
	startRow := int64(2)
	wg := threading.RoutineGroup{}
	req.Page = 0
	req.Size = 10000
	for {
		req.Page++
		// 获取数据
		list, _, _, err := s.AdminDashboard(c, req, true)
		if err != nil || len(list) == 0 {
			break
		}
		wg.RunSafe(func() {
			for _, item := range list {
				rowStartStr := strconv.Itoa(int(atomic.LoadInt64(&startRow)))
				switch req.Type {
				case storeVisitReq.DashboardStore:
					row = []any{
						item.Name, item.StoreCode, item.DeptName, item.VisitNum, item.VisitRate,
					}
				case storeVisitReq.DashboardSource, storeVisitReq.DashboardRegion:
					row = []any{
						item.Name, item.VisitNum, item.VisitRate,
					}
				case storeVisitReq.DashboardArea:
					row = []any{
						item.Name, item.ParentName, item.VisitNum, item.VisitRate,
					}
				}
				err = stream.SetRow("A"+rowStartStr, row)
				atomic.AddInt64(&startRow, 1)
			}
		})
	}
	wg.Wait()
	err = stream.Flush()
	if err != nil {
		return fileName, err
	}
	err = excelFile.SaveAs(fileName)
	if err != nil {
		return fileName, err
	}
	return fileName, err
}
