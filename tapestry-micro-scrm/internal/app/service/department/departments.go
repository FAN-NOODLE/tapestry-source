package department

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	"micro-scrm/internal/app/http/request/department"
	membersResp "micro-scrm/internal/app/http/response/members"
	departmentsModels "micro-scrm/internal/app/models/departments"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	"micro-scrm/internal/pkg/helper"
	"strconv"
	"strings"
)

type DepartmentsService struct {
	departmentsDao *departmentsDao.Departments
}

func NewDepartmentsService() DepartmentsService {
	return DepartmentsService{
		departmentsDao: departmentsDao.NewDepartments(),
	}
}

func (s DepartmentsService) GetAllDepartments(c app.Ctx, enterpriseId int64) (error, []*departmentsModels.DepartmentsModel) {
	return s.departmentsDao.GetAllDepartments(c, enterpriseId)
}

func (s DepartmentsService) GetSaStoreLists(c app.Ctx, request department.SAStoreListReq, isCoach bool) (int64, department.DepartmentsLists, error) {
	var (
		total  int64
		result []department.NewMomentDepartments
		err    error
	)
	var res department.DepartmentsLists
	var thirdIds []int64
	storeMap := make(map[int64]int)
	if isCoach {
		total, result, err = s.departmentsDao.GetSaStoreLists(c, request.EnterpriseID, request.StoreType)
	} else {
		switch request.StoreType {
		case 1: // 奥莱
			outletsIds, err := microWorkDao.NewStoreDepartment().GetDpIdsByType(c, request.EnterpriseID, microWorkDao.STORE_TYPE_OUTLETS)
			if err != nil {
				return total, res, err
			}
			for _, v := range outletsIds {
				item := v
				parseInt, err := strconv.ParseInt(item, 10, 64)
				if err != nil {
					return total, res, err
				}
				thirdIds = append(thirdIds, parseInt)
			}
		case 2: // 零售
			retailIds, err := microWorkDao.NewStoreDepartment().GetDpIdsByType(c, request.EnterpriseID, microWorkDao.STORE_TYPE_RETAIL)
			if err != nil {
				return total, res, err
			}
			// 零售
			for _, v := range retailIds {
				item := v
				parseInt, err := strconv.ParseInt(item, 10, 64)
				if err != nil {
					return total, res, err
				}
				thirdIds = append(thirdIds, parseInt)
			}
		case 0:
			many, err := microWorkDao.NewStoreDepartment().GetMany(c, map[string]interface{}{
				"enterprise_id": request.EnterpriseID,
				"deleted_at":    0,
			})
			if err != nil {
				return 0, department.DepartmentsLists{}, err
			}
			for _, model := range many {
				if !helper.InArray(int64(model.WechatDepartmentID), thirdIds) {
					storeMap[int64(model.WechatDepartmentID)] = int(model.StoreType)
					thirdIds = append(thirdIds, int64(model.WechatDepartmentID))
				}
			}
		}
		total, result, err = s.departmentsDao.GetSaStoreListsByThirdIds(c, request.EnterpriseID, request.StoreType, thirdIds)
	}

	var OCList []department.NewMomentDepartments
	var OCFList []department.NewMomentDepartments

	switch request.StoreType {
	case 0:
		for _, departments := range result {
			info := departments
			// var storeDepartmentName = departments.DpName
			var storeDepartmentName = departments.Name
			if isCoach {
				if strings.Contains(storeDepartmentName, "OCF") == true { // 奥莱
					info.StoreType = 1
					OCFList = append(OCFList, info)
				} else if strings.Contains(storeDepartmentName, "OC") == true { // 零售
					info.StoreType = 2
					OCList = append(OCList, info)
				}
			} else {
				if item, ok := storeMap[info.ThirdID]; ok {
					info.StoreType = item
					if item == 1 {
						OCFList = append(OCFList, info)
					}
					if item == 2 {
						OCList = append(OCList, info)
					}
				}
			}

		}

	case 1: // 奥莱
		OCFList = result
	case 2: // 零售
		OCList = result
	}
	res.OCList = OCList
	res.OCFList = OCFList

	return total, res, err
}

func (s DepartmentsService) GetStoreCodesByThirdIds(c app.Ctx, enterpriseId int64, thirdIds []int) ([]string, error) {
	// 查询门店编码
	storeList, err := microWorkDao.NewStoreDepartment().GetAllByWechatDepartmentIds(c, enterpriseId, thirdIds)
	if err != nil {
		return nil, err
	}
	var storeCodeSlice []string
	for _, store := range storeList {
		storeCodeSlice = append(storeCodeSlice, store.ThirdStoreCode)
	}
	return storeCodeSlice, nil
}

// GetDpListByStoreCode 根据门店编码查询门店信息
func (s DepartmentsService) GetDpListByStoreCode(c app.Ctx, enterpriseId int64, storeCodeSlice []string) (map[string]microWorkModel.StoreDepartmentModel, error) {
	// 查询门店信息
	list, err := microWorkDao.NewStoreDepartment().GetStoreInfoByStoreCode(c, enterpriseId, storeCodeSlice)
	if err != nil {
		return nil, err
	}
	var storeMapByCode = make(map[string]microWorkModel.StoreDepartmentModel)
	for _, info := range list {
		storeMapByCode[info.ThirdStoreCode] = info
	}
	return storeMapByCode, nil
}

func (s DepartmentsService) GetAllStoreCode(c app.Ctx, enterpriseId int64) ([]string, error) {
	// 查询门店编码
	storeList, err := microWorkDao.NewStoreDepartment().GetAll(c, enterpriseId)
	if err != nil {
		return nil, err
	}
	var storeCodeSlice []string
	for _, store := range storeList {
		storeCodeSlice = append(storeCodeSlice, store.ThirdStoreCode)
	}
	return storeCodeSlice, nil
}

func (s DepartmentsService) GetAllStore(c app.Ctx, enterpriseId int64, storeParentThirdId int) ([]microWorkModel.StoreDepartmentModel, error) {
	departments, err := departmentsDao.NewDepartments().GetListByThirdParentId(c, enterpriseId, storeParentThirdId)
	if err != nil {
		return nil, err
	}
	var dpIds []int
	for _, val := range departments {
		dpIds = append(dpIds, val.ThirdID)
	}
	if len(dpIds) == 0 {
		return []microWorkModel.StoreDepartmentModel{}, nil
	}
	// 查询门店
	return microWorkDao.NewStoreDepartment().GetAllByWechatDepartmentIds(c, enterpriseId, dpIds)

}

func (s DepartmentsService) GetStoreNameMap(c app.Ctx, enterpriseId int64) (map[int]string, error) {
	return microWorkDao.NewStoreDepartment().GetStoreNameMap(c, enterpriseId)

}

func (s DepartmentsService) GetDepList(c ctx.Context, enterpriseId int64, ids []int) ([]membersResp.DepListResp, error) {
	res := make([]membersResp.DepListResp, 0)
	result, err := s.departmentsDao.GetDepartmentsInfoThirdIds(c, enterpriseId, ids)
	if err != nil {
		return res, err
	}
	storeDepartments, err := microWorkDao.NewStoreDepartment().GetMany(c, map[string]interface{}{
		"enterprise_id":        enterpriseId,
		"wechat_department_id": ids,
	})
	if err != nil {
		return res, err
	}
	storeNameMap := make(map[int]microWorkModel.StoreDepartmentModel)
	for _, v := range storeDepartments {
		storeNameMap[int(v.WechatDepartmentID)] = v
	}

	for _, v := range result {
		storeName := ""
		storeCode := v.Name
		if item, ok := storeNameMap[v.ThirdID]; ok {
			storeName = item.StoreName
			storeCode = item.ThirdStoreCode
		}
		res = append(res, membersResp.DepListResp{
			DepartmentsModel: v,
			StoreName:        storeName,
			StoreCode:        storeCode,
		})
	}
	return res, err
}

func (s DepartmentsService) GetSaList(c ctx.Context, enterpriseId int64, ids []int) ([]membersResp.SaListResp, error) {
	res := make([]membersResp.SaListResp, 0)
	result, err := s.departmentsDao.GetMemberDataByThirdIds(c, enterpriseId, ids)
	if err != nil {
		return res, err
	}
	storeDepartments, err := microWorkDao.NewStoreDepartment().GetMany(c, map[string]interface{}{
		"enterprise_id":        enterpriseId,
		"wechat_department_id": ids,
	})
	if err != nil {
		return res, err
	}
	storeNameMap := make(map[int]microWorkModel.StoreDepartmentModel)
	for _, v := range storeDepartments {
		storeNameMap[int(v.WechatDepartmentID)] = v
	}

	for _, v := range result {
		storeName := ""
		storeCode := v.ThirdName
		if item, ok := storeNameMap[v.ThirdId]; ok {
			storeName = item.StoreName
			storeCode = item.ThirdStoreCode
		}
		res = append(res, membersResp.SaListResp{
			DepartmentsMembersData: v,
			StoreName:              storeName,
			StoreCode:              storeCode,
		})
	}
	return res, err

}
