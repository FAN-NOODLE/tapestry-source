package member

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm"
)

//// GetBAByStoreID_Scrm 根据Scrm的部门ID，获取该部门下的所有导购
//func (logic memberLogic) GetBAByStoreID_Scrm(ctx ctx.Context, departmentID int64) (data []micro_scrm.DepartmentsModel, err error) {
//	departmentMemberList, memberErr := logic.GetMemberListByDepartmentID_Scrm(ctx, departmentID)
//	if memberErr != nil {
//		err = memberErr
//		return
//	}
//
//	// 获取userid
//	userIDList := []string{}
//	for _, item := range departmentMemberList {
//		userIDList = append(userIDList, item.Member.Userid)
//	}
//
//	// 在store_worker表中查询导购信息
//
//
//	return
//}

// GetMemberListByDepartmentIDScrm 获取制定部门下的所有员工
func (logic memberLogic) GetMemberListByDepartmentIDScrm(ctx ctx.Context, departmentID int64) (data []micro_scrm.MemberDepartmentsModel, err error) {
	departmentList, departmentErr := logic.GetDepartmentListByDepartmentIDScrm(ctx, departmentID)
	if departmentErr != nil {
		err = departmentErr
		return
	}

	idList := []int64{}
	for _, item := range departmentList {
		idList = append(idList, item.ID)
	}

	// 获取员工的uid
	data, err = orm.Service.Structure.GetMemberInfoByDepartmentIDs(logic.enterpriseID, idList)
	return
}

// GetDepartmentListByDepartmentIDScrm 获取部门下的所有部门的信息
func (logic memberLogic) GetDepartmentListByDepartmentIDScrm(ctx ctx.Context, departmentID int64) (data []micro_scrm.DepartmentsModel, err error) {
	// 获取企业的组织架构
	departmentList, departmentErr := orm.Service.Structure.GetDepartmentAll(logic.enterpriseID)
	if departmentErr != nil {
		err = departmentErr
		return
	}

	return orm.Service.Structure.GetDepartmentInfoByThirdIDList(logic.enterpriseID, logic.recursiveDepartment(departmentList, departmentID))
}

func (logic memberLogic) recursiveDepartment(departmentList []micro_scrm.DepartmentsModel, departmentID int64) (data []int) {
	// 获取当前节点的ID
	rootID := 0
	for _, item := range departmentList {
		if item.ID == departmentID {
			rootID = item.ThirdID
		}
	}
	return queryNode(departmentList, []int{rootID})
}

func queryNode(departmentList []micro_scrm.DepartmentsModel, nodeIDList []int) (result []int) {
	for _, item := range departmentList {
		for _, nodeID := range nodeIDList {
			if item.ThirdParentID == nodeID {
				result = append(result, item.ThirdID)
			}
		}
	}

	if len(result) == 0 {
		return nodeIDList
	}
	return append(nodeIDList, queryNode(departmentList, result)...)
}
