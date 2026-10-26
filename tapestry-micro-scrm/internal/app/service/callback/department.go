package callbackSvc

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	membersDao "micro-scrm/internal/app/dao/members"
	"micro-scrm/internal/app/models"
	departmentsModels "micro-scrm/internal/app/models/departments"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/work"
	"strconv"
)

var (
	// CreateParty 新增部门事件
	CreateParty = "create_party"
	// UpdateParty 更新成员事件
	UpdateParty = "update_party"
	// DeleteParty 删除成员事件
	DeleteParty = "delete_party"
)

type Department struct {
	service.Base
	bossDetail     boss.EnterpriseModel
	departmentsDao *departmentsDao.Departments
}

func NewDepartment() *Department {
	return &Department{
		departmentsDao: departmentsDao.NewDepartments(),
	}
}

func (s *Department) Handel(c ctx.Context, data work.ContactCallbackDepartment, bossDetail boss.EnterpriseModel) error {
	s.Log().Info(c).Msgf("通讯录回调,事件类型：%s,部门事件数据：%v,企业数据：%v", data.ChangeType, data, bossDetail)
	s.bossDetail = bossDetail

	var err error
	switch data.ChangeType {
	case CreateParty:
		err = s.createEvent(c, data)
		break
	case UpdateParty:
		err = s.updateEvent(c, data)
		break
	case DeleteParty:
		err = s.deleteEvent(c, data)
		break
	}
	return err
}

func (s *Department) createEvent(c ctx.Context, data work.ContactCallbackDepartment) (err error) {
	thirdParentID := helper.StringToInt(data.ParentId)
	departmentsData := departmentsModels.DepartmentsModel{
		ThirdID:        data.Id,
		ThirdParentID:  thirdParentID,
		Name:           data.Name,
		Order:          int(data.Order),
		MicroBaseModel: models.CreateMicroBase(s.bossDetail.ID, app.UID()),
	}

	s.Log().Info(c).Msgf("添加部门数据：%v", helper.MarshalToString(departmentsData))
	err = s.departmentsDao.Create(c, departmentsData)
	if err != nil {
		s.Log().Error(c).Msgf("添加部门,异常：%v", err)
		return err
	}

	err = s.getDepartmentsAndChildDepartments(c)
	if err != nil {
		s.Log().Error(c).Msgf("获取部门及子部门及子子部门数据,异常：%v", err)
		return err
	}

	return err
}

func (s *Department) updateEvent(c ctx.Context, data work.ContactCallbackDepartment) (err error) {
	departmentsInfo, err := s.departmentsDao.GetInfo(c, departmentsModels.DepartmentsModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: s.bossDetail.ID,
		},
		ThirdID: data.Id,
	})
	if err != nil || departmentsInfo.ID == 0 {
		s.Log().Error(c).Msgf("查询部门,异常：%v", err)
		return err
	}

	// 判断是否更换部门
	if data.ParentId != "" {
		thirdParentID, _ := strconv.Atoi(data.ParentId)
		departmentsInfo.ThirdParentID = thirdParentID
	}

	// 判断是否更换部门名称
	if data.Name != "" {
		departmentsInfo.Name = data.Name
	}

	s.Log().Info(c).Msgf("更新部门数据：%v", helper.MarshalToString(departmentsInfo))
	err = s.departmentsDao.Update(c, departmentsInfo, map[string]interface{}{
		"enterprise_id": s.bossDetail.ID,
		"third_id":      data.Id,
	})
	if err != nil {
		s.Log().Error(c).Msgf("更新部门,异常：%v", err)
		return err
	}

	err = s.getDepartmentsAndChildDepartments(c)
	if err != nil {
		s.Log().Error(c).Msgf("获取部门及子部门及子子部门数据,异常：%v", err)
		return err
	}

	return err
}

func (s *Department) deleteEvent(c ctx.Context, data work.ContactCallbackDepartment) (err error) {
	departmentsInfo, err := s.departmentsDao.GetInfo(c, departmentsModels.DepartmentsModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: s.bossDetail.ID,
		},
		ThirdID: data.Id,
	})
	if err != nil {
		s.Log().Error(c).Msgf("查询部门,异常：%v", err)
		return err
	}

	// 删除部门
	err = s.departmentsDao.RemoveById(c, departmentsInfo.ID, s.bossDetail.ID)
	if err != nil {
		s.Log().Error(c).Msgf("删除部门,异常：%v", err)
		return err
	}

	// 删除员工部门绑定关系
	err = membersDao.NewMemberDepartments().RemoveByDepartmentId(c, departmentsInfo.ID, s.bossDetail.ID)
	if err != nil {
		s.Log().Error(c).Msgf("删除员工部门绑定关系,异常：%v", err)
		return err
	}

	return err
}

func (s *Department) getChildren(fatherMap map[int]int, children []int, belongToFather map[int][]int) {
	for _, third := range children {
		fatherMap[third] = third
		s.getChildren(fatherMap, belongToFather[third], belongToFather)
	}
}

// 获取部门及子部门及子子部门数据
func (s *Department) getDepartmentsAndChildDepartments(c ctx.Context) error {
	var (
		thirdToDpMap   = make(map[int]int64) // thirdPid => departmentId 一个thirdId对应的一个部门
		belongToFather = make(map[int][]int) // father => []children  一个父thirdId对应的多个子thirdId
	)
	dps, err := s.departmentsDao.GetAll(c, departmentsModels.DepartmentsModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: s.bossDetail.ID,
		},
	})
	if err != nil {
		s.Log().Error(c).Msgf("获取所有部门,异常：%v", err)
		return err
	}

	// 映射有子部门
	for _, v := range dps {
		thirdToDpMap[v.ThirdID] = v.ID
		if v.ThirdParentID > 0 {
			belongToFather[v.ThirdParentID] = append(belongToFather[v.ThirdParentID], v.ThirdID)
		}
	}

	var fatherMap = make(map[int]map[int]int)
	for father, children := range belongToFather {
		var childMap = make(map[int]int)
		childMap[father] = father
		s.getChildren(childMap, children, belongToFather)

		fatherMap[father] = childMap // 所有子部门，子子部门
	}

	var dpMap = make(map[int64][]int64)
	for father, children := range fatherMap {
		for _, child := range children {
			dpMap[thirdToDpMap[father]] = append(dpMap[thirdToDpMap[father]], thirdToDpMap[child])
		}
	}

	// 没有子部门的也放进来
	for third, dpId := range thirdToDpMap {
		if _, ok := dpMap[thirdToDpMap[third]]; !ok {
			dpMap[thirdToDpMap[third]] = append(dpMap[thirdToDpMap[third]], dpId)
		}
	}

	for father, children := range dpMap {
		childDepartments := helper.MarshalToString(children)
		updateData := departmentsModels.DepartmentsModel{
			ChildDepartments: childDepartments,
		}
		err = s.departmentsDao.Update(c, updateData, map[string]interface{}{
			"enterprise_id": s.bossDetail.ID,
			"id":            father,
		})
		if err != nil {
			s.Log().Error(c).Msgf("更新部门-子部门,异常：%v", err)
		}
	}

	return err
}
