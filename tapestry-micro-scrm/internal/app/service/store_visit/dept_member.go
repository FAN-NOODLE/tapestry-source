package storeVisitSvc

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	storeVisitDao "micro-scrm/internal/app/dao/store_visit"
	storeVisitReq "micro-scrm/internal/app/http/request/store_visit"
	storeVisitResp "micro-scrm/internal/app/http/response/store_visit"
	"micro-scrm/internal/app/models"
	storeVisitModel "micro-scrm/internal/app/models/store_visit"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	"time"
)

type DeptMemberSvc struct {
	service.Base
	visitDeptMemberDao *storeVisitDao.DeptMemberDao
	emailContactDao    *storeVisitDao.EmailContactDao
}

func NewDeptMemberSvc() *DeptMemberSvc {
	return &DeptMemberSvc{
		visitDeptMemberDao: storeVisitDao.NewDeptMemberDao(),
		emailContactDao:    storeVisitDao.NewEmailContactDao(),
	}
}

// DeptMemberList 预置的部门的人员列表
func (s *DeptMemberSvc) DeptMemberList(c ctx.Context, enterpriseId int64) ([]storeVisitResp.DeptInfo, error) {
	list, err := s.visitDeptMemberDao.GetAllInfo(c, enterpriseId, "dept_name != ''")
	if err != nil {
		return []storeVisitResp.DeptInfo{}, err
	}
	var deptList = make([]storeVisitResp.DeptInfo, 0)
	deptIndex := make(map[string]int) // Track department positions
	for _, contact := range list {
		member := storeVisitResp.MemberInfo{
			MemberId: contact.MemberID,
			Name:     contact.Name,
			UserId:   contact.Userid,
			Email:    contact.Email,
		}
		if contact.Members.Position != "" {
			if tmpP, ok := membersDao.StoreManagerPositionMap[contact.Members.Position]; ok {
				member.PositionAbbr = tmpP
			}
		}
		if len(contact.Members.MemberDepartments) > 0 {
			for _, department := range contact.Members.MemberDepartments {
				if department.DeletedAt == 0 && department.Department.DeletedAt == 0 {
					member.ThirdId = department.Department.ThirdID
					break
				}
			}
		}
		// Check if we've seen this department before
		index, exists := deptIndex[contact.DeptName]
		if !exists {
			// New department - add to the list and record its position
			deptInfo := storeVisitResp.DeptInfo{
				DeptName:   contact.DeptName,
				MemberList: []storeVisitResp.MemberInfo{member},
			}
			deptList = append(deptList, deptInfo)
			deptIndex[contact.DeptName] = len(deptList) - 1
		} else {
			// Existing department - append the member
			deptList[index].MemberList = append(deptList[index].MemberList, member)
		}
	}
	return deptList, nil

}

// GetEmailContact   个人通讯录人员列表
func (s *DeptMemberSvc) GetEmailContact(c ctx.Context, req storeVisitReq.BaseReq) ([]storeVisitResp.EmailContact, error) {
	list, err := s.emailContactDao.GetMany(c, req.EnterpriseId, "member_id = ?", req.MemberId)
	if err != nil {
		return []storeVisitResp.EmailContact{}, err
	}
	var emailList = make([]storeVisitResp.EmailContact, 0)
	for _, model := range list {
		emailList = append(emailList, storeVisitResp.EmailContact{
			ID:    model.ID,
			Email: model.Email,
		})
	}
	return emailList, nil
}

// SaveEmail    个人通讯录保存
func (s *DeptMemberSvc) SaveEmail(c ctx.Context, req storeVisitReq.EmailContactReq) (int64, error) {
	if req.ID > 0 {
		err := s.emailContactDao.UpdateData(c, req.EnterpriseId,
			map[string]interface{}{"id": req.ID, "member_id": req.MemberId},
			map[string]interface{}{"email": req.Email, "updated_at": time.Now().Unix()})
		return req.ID, err
	}
	info, _ := s.emailContactDao.GetOne(c, req.EnterpriseId, "member_id = ? AND email = ?", req.MemberId, req.Email)
	if info.ID > 0 {
		return info.ID, nil
	}
	var newId = app.UID()
	err := s.emailContactDao.CreateData(c, req.EnterpriseId, storeVisitModel.StoreVisitEmailContactModel{
		BaseModel: models.CreateBase(newId),
		MemberID:  req.MemberId,
		Email:     req.Email,
	})
	return newId, err
}

// DelEmail   个人通讯录删除
func (s *DeptMemberSvc) DelEmail(c ctx.Context, req storeVisitReq.EmailContactReq) error {
	if req.ID == 0 {
		return nil
	}
	return s.emailContactDao.DeleteData(c, req.EnterpriseId, "member_id = ? AND id = ?", req.MemberId, req.ID)
}

// SearchMembers    通过姓名或邮箱搜索用户
func (s *DeptMemberSvc) SearchMembers(c ctx.Context, req storeVisitReq.SearchMemberReq) ([]storeVisitResp.MemberInfo, error) {
	var memberList = make([]storeVisitResp.MemberInfo, 0)
	//先查自定义部门
	deptMemberList, err := s.visitDeptMemberDao.GetAllInfo(c, req.EnterpriseId, "name like ? AND dept_name != ''", "%"+req.Str+"%")
	if err != nil {
		return nil, err
	}
	var memberIds []int64
	for _, member := range deptMemberList {
		memberIds = append(memberIds, member.MemberID)
		memberInfo := storeVisitResp.MemberInfo{
			MemberId: member.MemberID,
			Name:     member.Name,
			Email:    member.Email,
			UserId:   member.Userid,
			DeptName: member.DeptName,
		}
		if member.Members.Position != "" {
			if tmpP, ok := membersDao.StoreManagerPositionMap[member.Members.Position]; ok {
				memberInfo.PositionAbbr = tmpP
			}
		}
		if len(member.Members.MemberDepartments) > 0 {
			for _, department := range member.Members.MemberDepartments {
				if department.DeletedAt == 0 && department.Department.DeletedAt == 0 {
					memberInfo.ThirdId = department.Department.ThirdID
					break
				}
			}
		}
		memberList = append(memberList, memberInfo)
	}
	//在查普通用户
	list, err := membersDao.NewMembers().GetMemberInfoByNameOrEmail(c, req.EnterpriseId, req.Str)
	if err != nil {
		return []storeVisitResp.MemberInfo{}, err
	}
	for _, member := range list {
		if helper.InArray(member.ID, memberIds) {
			//如果自定义部门已经有这个人了，则过滤
			continue
		}
		memberInfo := storeVisitResp.MemberInfo{
			MemberId: member.ID,
			Name:     member.Name,
			Email:    member.Email,
			UserId:   member.Userid,
		}
		if member.Position != "" {
			if tmpP, ok := membersDao.StoreManagerPositionMap[member.Position]; ok {
				memberInfo.PositionAbbr = tmpP
			}
		}
		if len(member.MemberDepartments) > 0 {
			for _, department := range member.MemberDepartments {
				if department.DeletedAt == 0 && department.Department.DeletedAt == 0 {
					memberInfo.ThirdId = department.Department.ThirdID
					memberInfo.DeptName = department.Department.Name
					break
				}
			}
		}
		memberList = append(memberList, memberInfo)
	}

	return memberList, nil
}

// LoginUserInfo     当前用户信息
func (s *DeptMemberSvc) LoginUserInfo(c ctx.Context, req storeVisitReq.BaseReq) (storeVisitResp.LoginUserInfo, error) {
	memberDepartmentList, err := membersDao.NewMembers().GetMemberInfoAndDepartmentByMemberIds(c, req.EnterpriseId, []int64{req.MemberId})
	if err != nil || len(memberDepartmentList) == 0 {
		return storeVisitResp.LoginUserInfo{}, err
	}
	// 整理导购id和导购信息的key=>value
	var loginUserInfo = storeVisitResp.LoginUserInfo{}
	for _, info := range memberDepartmentList {
		if info.DeletedAt > 0 {
			continue
		}
		loginUserInfo.MemberId = info.ID
		loginUserInfo.Name = info.Name
		loginUserInfo.UserId = info.Userid
		loginUserInfo.Avatar = info.Avatar
		loginUserInfo.Position = info.Position
		for _, department := range info.MemberDepartments {
			if department.DeletedAt > 0 {
				continue
			}
			if department.Department.ID > 0 && department.Department.DeletedAt == 0 {
				loginUserInfo.StoreCode = department.Department.Name
				loginUserInfo.DeptName = department.Department.Name
				loginUserInfo.ThirdId = department.Department.ThirdID
				break
			}
		}
	}
	managerInfo, _ := storeVisitDao.NewDeptMemberDao().GetOne(c, req.EnterpriseId, "member_id = ?", req.MemberId)
	loginUserInfo.IsStoreVisitManager = managerInfo.IsStoreVisitManager > 0
	loginUserInfo.IsStoreManager = helper.InArray(loginUserInfo.Position, membersDao.StoreVisitManagerPosition)
	if managerInfo.DeptName != "" {
		loginUserInfo.DeptName = managerInfo.DeptName
	}
	//增加门店中文名
	storeInfo, _ := microWorkDao.NewStoreDepartment().GetInfo(c, map[string]any{
		"enterprise_id":        req.EnterpriseId,
		"wechat_department_id": loginUserInfo.ThirdId,
	})
	loginUserInfo.StoreName = helper.RuleThree(req.Language == "en" && storeInfo.InventoryName != "", storeInfo.InventoryName, storeInfo.StoreName)
	return loginUserInfo, nil
}

func (s *DeptMemberSvc) DeptNameList(c ctx.Context, enterpriseId int64) ([]string, error) {
	return s.visitDeptMemberDao.GetDeptNameList(c, enterpriseId)
}
