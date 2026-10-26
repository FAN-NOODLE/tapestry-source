package callbackSvc

import (
	"errors"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	membersDao "micro-scrm/internal/app/dao/members"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	"micro-scrm/internal/app/models"
	departmentsModels "micro-scrm/internal/app/models/departments"
	membersModel "micro-scrm/internal/app/models/members"
	microBossModel "micro-scrm/internal/app/models/micro_boss"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/boss"
	"micro-scrm/internal/pkg/sdk/microwork"
	"micro-scrm/internal/pkg/sdk/uc"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/workwechat"
	"strings"
	"time"
)

var (
	// CreateUser 新增成员事件
	CreateUser = "create_user"
	// UpdateUser 更新成员事件
	UpdateUser = "update_user"
	// DeleteUser 删除成员事件
	DeleteUser = "delete_user"
)

type User struct {
	service.Base
	bossDetail   microBossModel.EnterpriseModel
	memberDetail membersModel.MemberModel

	storeWorkerDao     *microWorkDao.StoreWorker
	enterpriseAdminDao *microBossDao.EnterpriseAdmin

	membersDao          *membersDao.Members
	memberDepartmentDao *membersDao.MemberDepartments
	departmentsDao      *departmentsDao.Departments

	userRelation *microScrmDao.UserRelation
}

func NewUser() *User {
	return &User{
		storeWorkerDao:      microWorkDao.NewStoreWorker(),
		enterpriseAdminDao:  microBossDao.NewEnterpriseAdmin(),
		membersDao:          membersDao.NewMembers(),
		memberDepartmentDao: membersDao.NewMemberDepartments(),
		departmentsDao:      departmentsDao.NewDepartments(),
		userRelation:        microScrmDao.NewUserRelation(),
	}
}

func (s *User) Handel(c ctx.Context, data work.ContactCallbackUser, bossDetail microBossModel.EnterpriseModel) error {
	s.Log().Info(c).Msgf("通讯录回调,事件类型：%s,成员事件数据：%+v,企业数据：%+v", data.ChangeType, data, bossDetail)
	s.bossDetail = bossDetail

	// 查询员工数据
	memberInfo, err := s.membersDao.GetOneByUnscoped(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: bossDetail.ID},
		Userid:         data.UserID,
	})
	if err != nil {
		s.Log().Error(c).Msgf("查询员工数据,异常：%v", err)
	}
	s.memberDetail = memberInfo

	switch data.ChangeType {
	case CreateUser:
		err = s.createEvent(c, data)
		break
	case UpdateUser:
		err = s.updateEvent(c, data)
		break
	case DeleteUser:
		err = s.deleteEvent(c, data)
		break
	}

	return err
}

func (s *User) createEvent(c ctx.Context, data work.ContactCallbackUser) (err error) {
	newClient, err := workwechat.NewWorkWeChat(c, s.bossDetail.ID)
	if err != nil {
		s.Log().Error(c).Msgf("创建成员事件,初始化企微客户端异常：%v", err)
		return err
	}
	userInfo, err := newClient.BADetail(data.UserID)
	if err != nil {
		s.Log().Error(c).Msgf("创建成员事件,查询企微成员详情异常：%v", err)
		return err
	}

	// 同步员工和部门关系
	if err = s.addMember(c, data, userInfo); err != nil {
		s.Log().Error(c).Msgf("同步员工和部门关系,异常：%v", err)
	}

	// 添加员工为导购
	if err = s.addMemberToGuide(c, data); err != nil {
		s.Log().Error(c).Msgf("添加员工为导购,异常：%v", err)
	}

	return err
}

func (s *User) updateEvent(c ctx.Context, data work.ContactCallbackUser) error {
	// 调用Uc查询员工详情
	newClient, err := workwechat.NewWorkWeChat(c, s.bossDetail.ID)
	if err != nil {
		s.Log().Error(c).Msgf("更新成员事件,初始化企微客户端异常：%v", err)
		return err
	}
	userInfo, err := newClient.BADetail(data.UserID)
	if err != nil {
		s.Log().Error(c).Msgf("调用Uc查询员工详情,异常：%v", err)
		return err
	}
	// 判断当前员工是否在数据库中
	// 会存在三方使用api批量拉员工，不会走通讯录回调的创建事件，只有当员工首次进入企微激活才会触发员工信息变更->激活事件
	if s.memberDetail.ID == 0 {
		data.Department = helper.Int64ToString(int64(userInfo.MainDepartment))
		data.Status = helper.Int64ToString(int64(userInfo.Status))
		if err := s.addMember(c, data, userInfo); err != nil {
			s.Log().Error(c).Msgf("更新成员事件,同步员工和部门关系,异常：%v", err)
		}
		// 查询员工数据
		memberInfo, err := s.membersDao.GetOne(c, membersModel.MemberModel{
			MicroBaseModel: models.MicroBaseModel{EnterpriseID: s.bossDetail.ID},
			Userid:         data.UserID,
		})
		if err != nil {
			s.Log().Error(c).Msgf("更新成员事件,查询员工数据,异常：%v", err)
		}
		s.memberDetail = memberInfo
	}

	// 判断是否为更改手机号
	if data.Mobile != "" {
		if err := s.updateMemberMobile(c, data.Mobile, data.UserID); err != nil {
			s.Log().Error(c).Msgf("员工主动退出企业,异常：%v", err)
			return err
		}
	}

	// 员工主动退出企业
	if data.Status == helper.Int64ToString(membersModel.MembersStatusLeave) {
		err := s.membersDao.RemoveByUserId(c, data.UserID, s.bossDetail.ID, int64(membersModel.MembersStatusLeave))
		if err != nil {
			s.Log().Error(c).Msgf("员工主动退出企业,异常：%v", err)
			return err
		}
	}

	// 更换部门
	if data.Department != "" {
		departmentIds := strings.Split(data.Department, ",")
		if err := s.changeDepartment(c, departmentIds); err != nil {
			s.Log().Error(c).Msgf("员工更换部门,异常：%v", err)
			return err
		}
	}

	// 员工禁用
	if data.Status == helper.Int64ToString(membersModel.MembersStatusDisable) {
		err := s.membersDao.RemoveByUserId(c, data.UserID, s.bossDetail.ID, int64(membersModel.MembersStatusDisable))
		if err != nil {
			s.Log().Error(c).Msgf("员工禁用企业,异常：%v", err)
			return err
		}
	}

	// 更换userId
	if data.NewUserID != "" {
		if err := s.changeUserId(c, data.UserID, data.NewUserID); err != nil {
			s.Log().Error(c).Msgf("员工更换userId,异常：%v", err)
			return err
		}
	}

	membersData := map[string]interface{}{
		"name":         userInfo.Name,
		"mobile":       userInfo.Mobile,
		"gender":       helper.StringToInt(userInfo.Gender),
		"avatar":       userInfo.Avatar,
		"thumb_avatar": userInfo.ThumbAvatar,
		"qr_code":      userInfo.QrCode,
		"position":     userInfo.Position,
		"deleted_at":   0,
	}
	if userInfo.Extattr.Attrs != nil {
		for _, attr := range userInfo.Extattr.Attrs {
			if attr.Name == "sid" && attr.Value != "" {
				membersData["open_userid"] = attr.Value
			}
		}
	}
	//  员工激活
	if data.Status == helper.Int64ToString(membersModel.MembersStatusEnable) || helper.StringToInt(data.Status) == membersModel.MembersStatusEnable {
		membersData["status"] = membersModel.MembersStatusEnable
		// 添加员工为导购
		if err = s.addMemberToGuide(c, data); err != nil {
			s.Log().Error(c).Msgf("添加员工为导购,异常：%v", err)
		}
	}

	// 更新用户基本信息
	err = s.membersDao.UpdateData(c, map[string]any{"id": s.memberDetail.ID, "enterprise_id": s.bossDetail.ID}, membersData)
	if err != nil {
		s.Log().Error(c).Msgf("更新member表员工数据,异常：%v", err)
		return err
	}

	return err
}

func (s *User) deleteEvent(c ctx.Context, data work.ContactCallbackUser) (err error) {
	// 移除员工
	if err = s.removeMember(c, data.UserID); err != nil {
		return err
	}
	return err
}

// 添加员工为导购
func (s *User) addMemberToGuide(c ctx.Context, data work.ContactCallbackUser) error {
	// 判断员工是否为导购
	memberInfo, err := s.storeWorkerDao.GetInfo(c, microWorkModel.StoreWorkerModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: s.bossDetail.ID},
		MemUserid:      data.UserID,
	})
	if err != nil {
		s.Log().Error(c).Msgf("查询员工是否为导购,错误:%v", err)
		return err
	}

	s.Log().Info(c).Msgf("员工信息", memberInfo)

	if memberInfo.ID == 0 || memberInfo.DeletedAt > 0 {
		// 获取企业配置要同步的部门
		configData, err := boss.EnterpriseConfig(c, map[string]interface{}{
			"enterprise_id": s.bossDetail.ID,
			"key":           "store_root_dept_id",
		})

		s.Log().Info(c).Msgf("获取企业配置要同步的部门", configData)

		if err != nil {
			s.Log().Error(c).Msgf("获取企业自定义配置,错误:%v", err)
			return err
		}

		if len(configData.Data) == 0 {
			return errors.New("获取企业自定义配置为空")
		}

		var thirdIds []int64
		thirdId := configData.Data[0].Value
		thirdIds = append(thirdIds, thirdId)

		// 获取要同步部门的Uc部门详情
		departmentDetails, err := uc.NewDepartment().List(c,
			map[string]interface{}{"thirds": thirdIds},
			uc.CommunalContact,
			s.bossDetail.Domain,
		)

		s.Log().Info(c).Msgf("获取要同步部门的Uc部门详情", departmentDetails)

		if err != nil {
			s.Log().Error(c).Msgf("获取要同步部门的Uc部门详情,错误:%v", err)
			return err
		}

		// 获取要同步部门下的所有部门
		departmentList, err := uc.NewDepartment().List(c,
			map[string]interface{}{"dpId": departmentDetails[0].DpId, "departmentChildrenFlag": 1, "pageSize": 1000},
			uc.CommunalContact,
			s.bossDetail.Domain,
		)

		s.Log().Info(c).Msgf("获取要同步部门下的所有部门", departmentList)

		if err != nil {
			s.Log().Error(c).Msgf("获取要同步部门下的所有部门Id集合,错误:%v", err)
			return err
		}

		thirdIdStr := helper.Int64ToString(thirdId)
		dpThirdIds := []string{thirdIdStr}
		for _, depart := range departmentList {
			dpThirdIds = append(dpThirdIds, depart.DpThirdid)
		}

		s.Log().Info(c).Msgf("判断员工是否在自动创建导购部门的范围下", thirdIdStr, dpThirdIds)

		// 判断员工是否在自动创建导购部门的范围下
		if helper.InArray[string](thirdIdStr, dpThirdIds) {
			member, err := uc.MemberByUserId(c, uc.MemberByUserIdRequest{MemUserID: data.UserID}, uc.CommunalContact, s.bossDetail.Domain)
			if err != nil {
				s.Log().Error(c).Msgf("成员列表,错误:%v", err)
				return err
			}

			s.Log().Info(c).Msgf("成员列表", member)

			// 组装添加导购数据
			var addMembersData []microwork.AddGuideMemberRequest
			addMembersData = append(addMembersData, microwork.AddGuideMemberRequest{
				MemUid:   member.MemUid,
				Avatar:   member.MemFace,
				UserId:   member.MemUserID,
				UserName: member.MemUsername,
				Mobile:   member.MemMobile,
			})

			requestData := microwork.AddGuideRequest{
				Type:         3,
				EnterpriseID: s.bossDetail.ID,
				Member:       addMembersData,
				OperatorId:   0,
				OperatorName: "system",
			}

			s.Log().Info(c).Msgf("组装添加导购数据", helper.MarshalToString(requestData))

			_, err = microwork.NewStoreWorker().AddGuide(c, requestData)
			if err != nil {
				s.Log().Error(c).Msgf("添加导购(企业导购,社群导购),错误:%v", err)
				return err
			}
		}
	}

	return err
}

// 更新员工手机号
func (s *User) updateMemberMobile(c ctx.Context, mobile, userId string) error {
	// 1.1 通过用户名查员工的旧手机号
	memberInfo, err := s.membersDao.GetOneByUserId(c, s.bossDetail.ID, userId)
	if err != nil {
		s.Log().Error(c).Msgf("查询员工手机号,异常：%v", err)
		return err
	}
	// 1.2 更新员工表信息
	err = s.storeWorkerDao.UpdateDataByMobile(c,
		microWorkModel.StoreWorkerModel{Mobile: mobile, SmsMobile: mobile},
		memberInfo.Mobile,
		s.bossDetail.ID)
	if err != nil {
		s.Log().Error(c).Msgf("更新员工表手机号为新手机号,异常：%v", err)
		return err
	}

	// 1.3 查询管理员表,对应的旧手机号
	enterpriseAdminInfo, err := s.enterpriseAdminDao.GetInfo(c, microBossModel.EnterpriseAdminModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: s.bossDetail.ID},
		Mobile:         memberInfo.Mobile})
	if err != nil {
		s.Log().Error(c).Msgf("查询管理员表,对应的旧手机号,异常：%v", err)
	}
	if enterpriseAdminInfo.ID > 0 {
		// 1.4 更新管理员表手机号为新手机号
		err = s.enterpriseAdminDao.UpdateMobileByOldMobile(c, memberInfo.Mobile, mobile, s.bossDetail.ID)
		if err != nil {
			s.Log().Error(c).Msgf("更新管理员表手机号为新手机号,异常：%v", err)
			return err
		}
	}

	// 更新员工表用户手机号
	membersData := membersModel.MemberModel{Mobile: mobile}
	err = s.membersDao.UpdateDataByUserId(c, membersData, userId, s.bossDetail.ID)
	return err
}

// 添加员工
func (s *User) addMember(c ctx.Context, data work.ContactCallbackUser, userInfo *work.MemberDetail) error {
	enterpriseId := s.bossDetail.ID
	userId := data.UserID
	memberInfo := s.memberDetail
	thirdId := helper.StringToInt(data.Department)
	if memberInfo.ID == 0 {
		thirdId = int(userInfo.MainDepartment)
	}
	// 查询部门信息
	departmentInfo, err := s.departmentsDao.GetInfo(c, departmentsModels.DepartmentsModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: enterpriseId},
		ThirdID:        thirdId,
	})
	if err != nil {
		s.Log().Error(c).Msgf("查询部门信息,异常：%v", err)
		return err
	}

	// 判断员工是是否在数据库存在
	// 不存在则添加员工数据，添加员工和部门关系数据
	// 存在则更新员工表和员工部门关系表数据
	if memberInfo.ID == 0 {
		// 添加员工数据
		membersData := membersModel.MemberModel{
			MicroBaseModel: models.CreateMicroBase(enterpriseId, app.UID()),
			Userid:         userInfo.Userid,
			Name:           userInfo.Name,
			Mobile:         userInfo.Mobile,
			Position:       userInfo.Position,
			Gender:         int32(helper.StringToInt(userInfo.Gender)),
			Email:          userInfo.Email,
			Avatar:         userInfo.Avatar,
			ThumbAvatar:    userInfo.ThumbAvatar,
			Telephone:      userInfo.Telephone,
			Alias:          userInfo.Alias,
			Status:         int32(userInfo.Status),
			QrCode:         userInfo.QrCode,
			Address:        userInfo.Address,
			Type:           membersModel.MembersTypeSA,
		}
		for _, attr := range userInfo.Extattr.Attrs {
			if attr.Name == "sid" && attr.Value != "" {
				membersData.OpenUserid = attr.Value
				break
			}
		}

		err = s.membersDao.Create(c, membersData)
		if err != nil {
			s.Log().Error(c).Msgf("添加员工数据,异常：%v", err)
			return err
		}
		s.addConversationArchiveTag(c, membersData, userInfo.Department)
		// 添加员工和部门关系数据
		err = s.memberDepartmentDao.Create(c, membersModel.MemberDepartmentsModel{
			MicroBaseModel: models.CreateMicroBase(enterpriseId, app.UID()),
			MemberID:       membersData.ID,
			DepartmentID:   departmentInfo.ID,
		})
		if err != nil {
			s.Log().Error(c).Msgf("添加员工和部门关系数据,异常：%v", err)
		}
		return err
	}

	// 1.更新员工表数据
	updateData := map[string]interface{}{
		"name":       data.Name,
		"alias":      data.Alias,
		"mobile":     data.Mobile,
		"position":   data.Position,
		"gender":     data.Gender,
		"status":     data.Status,
		"type":       membersModel.MembersTypeSA,
		"deleted_at": 0,
	}
	if userInfo.Extattr.Attrs != nil {
		for _, attr := range userInfo.Extattr.Attrs {
			if attr.Name == "sid" && attr.Value != "" {
				updateData["open_userid"] = attr.Value
			}
		}
	}
	err = s.membersDao.UpdateDataByUserId(c, updateData, userId, enterpriseId)
	if err != nil {
		s.Log().Error(c).Msgf("更新员工表数据,异常：%v", err)
		return err
	}

	// 2.查询员工和部门关系数据关系是否存在
	memberDepartmentInfo, err := s.memberDepartmentDao.GetInfo(c, membersModel.MemberDepartmentsModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: enterpriseId},
		MemberID:       memberInfo.ID,
		DepartmentID:   departmentInfo.ID,
	})
	if err != nil {
		s.Log().Error(c).Msgf("查询员工和部门关系数据,异常：%v", err)
		return err
	}

	if memberDepartmentInfo.ID > 0 {
		// 3.更新员工和部门关系数据
		err = s.memberDepartmentDao.EnableByMemberIdAndDepartmentId(c, enterpriseId, memberInfo.ID, departmentInfo.ID)
		if err != nil {
			s.Log().Error(c).Msgf("更新员工和部门关系数据,异常：%v", err)
			return err
		}
		return err
	}

	// 4.添加员工和部门关系数据
	err = s.memberDepartmentDao.Create(c, membersModel.MemberDepartmentsModel{
		MicroBaseModel: models.CreateMicroBase(enterpriseId, app.UID()),
		MemberID:       memberInfo.ID,
		DepartmentID:   departmentInfo.ID,
	})
	if err != nil {
		s.Log().Error(c).Msgf("添加员工和部门关系数据,异常：%v", err)
		return err
	}

	return err
}

func (s *User) addConversationArchiveTag(c ctx.Context, member membersModel.MemberModel, departmentThirdIDs []int64) {
	if !helper.InArray(member.Position, membersDao.ConversationArchiveTag) {
		return
	}

	configData, err := boss.EnterpriseConfig(c, map[string]interface{}{
		"enterprise_id": member.EnterpriseID,
		"key":           "store_root_dept_id",
	})
	if err != nil {
		s.Log().Error(c).Msgf("添加会话存档标签,获取门店根部门配置异常：%v", err)
		return
	}
	if len(configData.Data) == 0 {
		s.Log().Error(c).Msg("添加会话存档标签,未配置门店根部门")
		return
	}

	storeDepartments, err := s.departmentsDao.GetListByThirdParentId(c, member.EnterpriseID, int(configData.Data[0].Value))
	if err != nil {
		s.Log().Error(c).Msgf("添加会话存档标签,查询门店部门异常：%v", err)
		return
	}

	storeThirdIDs := make([]int64, 0, len(storeDepartments))
	for _, department := range storeDepartments {
		storeThirdIDs = append(storeThirdIDs, int64(department.ThirdID))
	}
	if len(helper.ArrayIntersect(storeThirdIDs, departmentThirdIDs)) == 0 {
		return
	}

	workImpl := new(work.QyWorkImpl).SetDomain(s.bossDetail.Domain)
	err = workImpl.AddTagUsers(c, work.AddTagUsersRequest{
		QyGwRequest: work.QyGwRequest{ReqPL: work.QWCommunalContact},
		TagID:       config.CfgData.CoachSaTag.TagId,
		UserList:    []string{member.Userid},
	})
	if err != nil {
		s.Log().Error(c).Msgf("添加会话存档标签异常, userid:%s, err:%v", member.Userid, err)
	}
}

// 移除员工
func (s *User) removeMember(c ctx.Context, userId string) error {
	enterpriseId := s.bossDetail.ID

	memberInfo := s.memberDetail

	// 禁用导购
	if err := s.storeWorkerDao.DisableByMemUserId(c, userId, enterpriseId); err != nil {
		s.Log().Error(c).Msgf("禁用导购,异常：%v", err)
	}

	// 更新member表员工状态
	err := s.membersDao.UpdateData(c, map[string]any{"id": s.memberDetail.ID, "enterprise_id": s.bossDetail.ID}, map[string]any{
		"status":     membersModel.MembersStatusLeave,
		"deleted_at": time.Now().Unix(),
	})

	if err != nil {
		s.Log().Error(c).Msgf("更新member表员工状态,异常：%v", err)
	}

	// 移除员工和部门关系
	err = s.memberDepartmentDao.RemoveByMemberId(c, memberInfo.ID, enterpriseId)
	if err != nil {
		s.Log().Error(c).Msgf("移除员工和部门关系,异常：%v", err)
	}
	// 禁用企业管理员
	if err = s.enterpriseAdminDao.RemoveByMobile(c, memberInfo.Mobile, enterpriseId); err != nil {
		s.Log().Error(c).Msgf("禁用企业管理员,异常：%v", err)
	}
	// 解绑所有好友关系
	if err = microScrmDao.NewUserRelation().BatchDeleteExternalRelationByMemberIds(c, enterpriseId, time.Now().Unix()-10, []int64{memberInfo.ID}); err != nil {
		s.Log().Error(c).Msgf("解绑所有好友关系,异常：%v", err)
	}

	return err
}

// 更换部门
func (s *User) changeDepartment(c ctx.Context, departmentIds []string) error {
	enterpriseId := s.bossDetail.ID

	// 通过企微部门id查询本地部门详情
	departmentInfoList, _ := s.departmentsDao.GetDepartmentsThirdIds(c, enterpriseId, departmentIds)
	var addMemberDepartmentIds []int64
	for _, model := range departmentInfoList {
		addMemberDepartmentIds = append(addMemberDepartmentIds, model.ID)
	}

	// 删除的部门和添加的部门及更新部门id集合
	var addDepartmentIds, delDepartmentIds, updateDepartmentIds []int64

	// 查询员工已绑定的部门
	existedMemberDepartmentInfo, _ := s.memberDepartmentDao.GetAllByMemberId(c, enterpriseId, s.memberDetail.ID)
	if len(existedMemberDepartmentInfo) > 0 {
		var existedMemberDepartmentIds []int64
		for _, model := range existedMemberDepartmentInfo {
			existedMemberDepartmentIds = append(existedMemberDepartmentIds, model.DepartmentID)
			if helper.InArray[int64](model.DepartmentID, addMemberDepartmentIds) && model.DeletedAt > 0 {
				updateDepartmentIds = append(updateDepartmentIds, model.DepartmentID)
			}
			if !helper.InArray[int64](model.DepartmentID, addMemberDepartmentIds) {
				delDepartmentIds = append(delDepartmentIds, model.DepartmentID)
			}
		}

		addDepartmentIds = helper.ArrayDiff(addMemberDepartmentIds, existedMemberDepartmentIds)

		if len(updateDepartmentIds) > 0 {
			err := s.memberDepartmentDao.EnableByMemberIdAndDepIds(c, s.memberDetail.ID, enterpriseId, updateDepartmentIds)
			if err != nil {
				s.Log().Error(c).Msgf("启动已删除员工和部门关系数据,异常：%v", err)
				return err
			}
		}

		if len(delDepartmentIds) > 0 {
			// 删除员工和部门关系数据
			err := s.memberDepartmentDao.RemoveByMemberIdAndDepIds(c, s.memberDetail.ID, enterpriseId, delDepartmentIds)
			if err != nil {
				s.Log().Error(c).Msgf("删除员工和部门关系数据,异常：%v", err)
				return err
			}
		}
	}

	// 添加员工和部门关系数据
	if len(addDepartmentIds) > 0 {
		var createMemberDepartmentData []membersModel.MemberDepartmentsModel
		for _, departmentId := range addDepartmentIds {
			memberDepartmentData := membersModel.MemberDepartmentsModel{
				MicroBaseModel: models.CreateMicroBase(enterpriseId, app.UIDContext(c)),
				MemberID:       s.memberDetail.ID,
				DepartmentID:   departmentId,
			}
			createMemberDepartmentData = append(createMemberDepartmentData, memberDepartmentData)
		}
		err := s.memberDepartmentDao.BatchesCreate(c, createMemberDepartmentData)
		if err != nil {
			s.Log().Error(c).Msgf("添加员工和部门关系数据,异常：%v", err)
		}
	}

	return nil
}

// 更换userId
func (s *User) changeUserId(c ctx.Context, userId, newUserId string) error {
	enterpriseId := s.bossDetail.ID

	// 查询员工手机号
	memberInfo, err := s.storeWorkerDao.GetInfo(c, microWorkModel.StoreWorkerModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: enterpriseId},
		MemUserid:      userId,
	})
	if err != nil {
		s.Log().Error(c).Msgf("查询员工手机号,异常：%v", err)
		return err
	}

	// 更新员工表员工企微id
	err = s.membersDao.UpdateDataByUserId(c, map[string]interface{}{"userid": newUserId}, userId, enterpriseId)
	if err != nil {
		s.Log().Error(c).Msgf("更新员工表员工企微id,异常：%v", err)
		return err
	}

	// 更新外部联系人关系表员工企微id
	err = s.userRelation.UpdateData(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"member_id":     memberInfo.ID,
		"user_id":       newUserId,
	}, map[string]interface{}{
		"userid": newUserId,
	})
	if err != nil {
		s.Log().Error(c).Msgf("更新外部联系人关系表员工企微id,异常：%v", err)
		return err
	}

	return nil
}
