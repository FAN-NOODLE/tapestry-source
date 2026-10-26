package openSvc

import (
	"errors"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	drainageContactDao "micro-scrm/internal/app/dao/drainage"
	membersDao "micro-scrm/internal/app/dao/members"
	workExternalDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	openDao "micro-scrm/internal/app/dao/open"
	thirdDao "micro-scrm/internal/app/dao/third"
	reserveReq "micro-scrm/internal/app/http/request/open"
	reserveRequest "micro-scrm/internal/app/http/request/reserve"
	reserveResp "micro-scrm/internal/app/http/response/open"
	"micro-scrm/internal/app/models"
	departmentsModels "micro-scrm/internal/app/models/departments"
	drainageContactModel "micro-scrm/internal/app/models/drainage"
	membersModel "micro-scrm/internal/app/models/members"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	openModel "micro-scrm/internal/app/models/open"
	thirdModel "micro-scrm/internal/app/models/third"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/work"
	"sort"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gorm.io/gorm"
)

type ReserveService struct {
	service.Base

	membersDao     *membersDao.Members
	departmentsDao *departmentsDao.Departments
}

func NewReserve() *ReserveService {
	return &ReserveService{
		membersDao:     membersDao.NewMembers(),
		departmentsDao: departmentsDao.NewDepartments(),
	}
}

const (
	AssignGuideTypeStoreWorker  = 0
	AssignGuideTypeReserveGuide = 1
	assignGuideDefaultState     = "0"
)

// GenerateReserve 生成预约记录
func (s ReserveService) GenerateReserve(c app.Ctx, data reserveReq.GenerateReserveRequest) (reserveResp.Success, error) {
	var resp reserveResp.Success

	reserveRecordDao := thirdDao.NewReserveRecordDao()
	// 查询预约记录
	recordInfo, err := reserveRecordDao.GetInfo(c, thirdModel.ReserveRecordModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: data.EnterpriseId,
		},
		// StoreCode:       data.StoreCode,//查询条件不需要带这个
		ReserveNumberID: data.ReserveNumberId,
		// Unionid:         data.UnionId, //查询条件不需要带这个
	})
	if err != nil {
		return resp, err
	}

	// 由RSVP小程序向WeClient pro推送本次预约的“门店Code”，“预约单ID”，“预约渠道ID”，“预约人uninoid”
	// 是否首次预约，向RSVP返回顾客已预约的“SA Code”，“门店Code”和“预约单ID”
	// 非首次预约，向顾客推送门店中的好友导购。若导购存在多个，推送首次预约的导购好友。
	if recordInfo.ID > 0 {
		// 获取员工信息
		membersInfo, _ := s.membersDao.GetOne(c, membersModel.MemberModel{
			MicroBaseModel: models.MicroBaseModel{
				BaseModel:    models.BaseModel{ID: recordInfo.MemberID},
				EnterpriseID: recordInfo.EnterpriseID,
			},
		})
		resp = reserveResp.Success{
			StoreCode:       recordInfo.StoreCode,
			ReserveNumberId: recordInfo.ReserveNumberID,
			MemberId:        recordInfo.MemberID,
			MemberName:      membersInfo.Name,
			WechatUserId:    membersInfo.Userid,
			ContactQrcode:   recordInfo.ContactQrcode,
		}
		return resp, nil
	}
	// 查询该用户在该门店上次预约记录
	lastRecord, err := reserveRecordDao.GetInfo(c, thirdModel.ReserveRecordModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: data.EnterpriseId,
		},
		StoreCode: data.StoreCode,
		Unionid:   data.UnionId,
	})
	if err != nil {
		return resp, err
	}
	var lastMemberId int64
	if lastRecord.ID > 0 {
		lastMemberId = lastRecord.MemberID
	}

	// -------------------处理不存在预约记录--------------------
	// 查询部门信息
	departmentInfo, err := s.getDepartmentByStoreCode(c, data.EnterpriseId, data.StoreCode)
	if err != nil {
		return resp, err
	}

	// 获取部门下的所有员工Id集合
	memberIds, err := s.getDepartmentMemberIds(c, data.EnterpriseId, departmentInfo.ID)
	if err != nil {
		return resp, err
	}

	// 查询客户和员工之前的绑定关系
	workUserRelations, err := workExternalDao.NewUserRelation().GetUserByUnionIdAndMemberIds(c, data.EnterpriseId, data.UnionId, memberIds)
	if err != nil {
		return resp, err
	}

	// 获取可以分配的员工id集合
	assignMembersIds, err := s.getAssignMemberIds(c, data.EnterpriseId, departmentInfo.ID, memberIds)
	if err != nil {
		return resp, err
	}

	var assignMemberId int64
	// 首次预约-无好友关系
	if len(workUserRelations) == 0 {
		assignMemberId = assignMembersIds[0]
	}
	// 首次预约-有好友关系-唯一好友关系
	if len(workUserRelations) == 1 {
		assignMemberId = workUserRelations[0].MemberID
	}

	// 首次预约-有好友关系-多好友关系
	if len(workUserRelations) > 1 {
		var ksaMemberId int64
		// 获取专属导购员工id
		for _, relation := range workUserRelations {
			if relation.IsKeySa == 1 {
				ksaMemberId = relation.MemberID
			}
		}
		// 存在专属sa导购
		if ksaMemberId > 0 {
			assignMemberId = ksaMemberId
		} else {
			assignMemberId = workUserRelations[0].MemberID
		}
	}
	for _, memberId := range memberIds {
		if memberId == lastMemberId {
			assignMemberId = lastMemberId
		}
	}

	reserveRecordAssignDao := thirdDao.NewReserveRecordAssignDao()

	// 查询员工绑定的部门分配记录
	reserveRecordAssignData, err := reserveRecordAssignDao.GetInfo(c, thirdModel.ReserveRecordAssignModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: data.EnterpriseId},
		DepartmentID:   departmentInfo.ID,
		MemberID:       assignMemberId,
	})

	reserveRecordId := app.UIDContext(c)
	// 生成联系我二维码
	contactQrcode, err := s.genContactQrcode(c, data.EnterpriseId, assignMemberId, reserveRecordId, "")
	if err != nil || contactQrcode.Errcode != 0 {
		if contactQrcode != nil && contactQrcode.Errcode != 0 {
			err = errors.New(contactQrcode.Errmsg)
		}
		s.Log().Error(c).Msgf("生成联系我二维码,异常：%v", err)
		return resp, err
	}

	// 开启事务
	err = app.DB().WithContext(c).Transaction(func(tx *gorm.DB) error {
		if reserveRecordAssignData.ID > 0 {
			reserveRecordAssignData.AssignNumber++
			err = reserveRecordAssignDao.WithTx(tx).Save(c, reserveRecordAssignData)
			if err != nil {
				s.Log().Error(c).Msgf("更新分配记录,异常：%v", err)
				tx.Rollback()
				return err
			}
		} else {
			// 添加分配记录
			reserveRecordAssignSaveData := thirdModel.ReserveRecordAssignModel{
				MicroBaseModel: models.CreateMicroBase(data.EnterpriseId, app.UIDContext(c)),
				DepartmentID:   departmentInfo.ID,
				MemberID:       assignMemberId,
				AssignNumber:   1,
			}
			err = reserveRecordAssignDao.WithTx(tx).Create(c, reserveRecordAssignSaveData)
			if err != nil {
				s.Log().Error(c).Msgf("添加分配记录,异常：%v", err)
				tx.Rollback()
				return err
			}
		}

		// 添加预约记录
		reserveRecordSaveData := thirdModel.ReserveRecordModel{
			MicroBaseModel:  models.CreateMicroBase(data.EnterpriseId, reserveRecordId),
			StoreCode:       data.StoreCode,
			ReserveNumberID: data.ReserveNumberId,
			ServiceName:     data.ServiceName,
			Unionid:         data.UnionId,
			DepartmentID:    departmentInfo.ID,
			MemberID:        assignMemberId,
			ConfigID:        contactQrcode.ConfigId,
			ContactQrcode:   contactQrcode.QrCode,
		}
		err = reserveRecordDao.WithTx(tx).Create(c, reserveRecordSaveData)
		if err != nil {
			s.Log().Error(c).Msgf("添加预约记录,异常：%v", err)
			tx.Rollback()
			return err
		}
		return err
	})

	info, err := s.membersDao.GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: data.EnterpriseId,
			BaseModel: models.BaseModel{ID: assignMemberId},
		},
	})

	if err != nil {
		return reserveResp.Success{}, err
	}
	resp = reserveResp.Success{
		StoreCode:       data.StoreCode,
		ReserveNumberId: data.ReserveNumberId,
		MemberId:        assignMemberId,
		WechatUserId:    info.Userid,
		MemberName:      info.Name,
		ContactQrcode:   contactQrcode.QrCode,
	}

	return resp, nil
}

func (s ReserveService) AssignGuide(c app.Ctx, data reserveReq.AssignGuideRequest) (reserveResp.AssignGuideSuccessResp, error) {
	var resp reserveResp.AssignGuideSuccessResp
	state, err := s.getAssignGuideState(c, data.EnterpriseId, data.Type)
	if err != nil {
		return resp, err
	}
	// 查询部门信息
	departmentInfo, err := s.getDepartmentByStoreCode(c, data.EnterpriseId, data.StoreCode)
	if err != nil || departmentInfo.ID == 0 {
		return resp, errors.New("部门不存在")
	}

	// 获取部门下的所有员工Id集合
	departmentMembers, err := s.departmentsDao.GetDepartmentsMembersByDepartmentId(c, data.EnterpriseId, departmentInfo.ID)
	if err != nil || len(departmentMembers) == 0 {
		return resp, errors.New("该部门下没有人员")
	}
	var memberIds []int64
	var memberMap = make(map[int64]departmentsModels.DepartmentsMembers)
	for _, member := range departmentMembers {
		memberIds = append(memberIds, member.MemberId)
		memberMap[member.MemberId] = member
	}
	// 查询客户和员工之前的绑定关系
	workUserRelations, err := workExternalDao.NewUserRelation().GetUserByUnionIdAndMemberIds(c, data.EnterpriseId, data.UnionId, memberIds)
	if err != nil {
		return resp, err
	}

	var assignMemberId int64
	// 首次预约-无好友关系
	if len(workUserRelations) == 0 {
		var saleMemberIds []int64
		for _, member := range departmentMembers {
			if helper.InArray(strings.ToLower(member.Position), helper.GetSalesPosition()) {
				saleMemberIds = append(saleMemberIds, member.MemberId)
			}
		}
		if len(saleMemberIds) == 0 {
			return resp, errors.New("该部门下没有sale人员")
		}
		assignMemberId = helper.ArrayRand(saleMemberIds)
	} else if len(workUserRelations) == 1 {
		// 首次预约-有好友关系-唯一好友关系
		assignMemberId = workUserRelations[0].MemberID
	} else if len(workUserRelations) > 1 {
		// 首次预约-有好友关系-多好友关系
		var ksaMemberId int64
		// 获取专属导购员工id
		for _, relation := range workUserRelations {
			if relation.IsKeySa == 1 {
				ksaMemberId = relation.MemberID
			}
		}
		// 存在专属sa导购
		if ksaMemberId > 0 {
			assignMemberId = ksaMemberId
		} else {
			assignMemberId = workUserRelations[0].MemberID
		}
	}
	memberInfo := memberMap[assignMemberId]
	//获取联系我二维码缓存
	wayModel, _ := openDao.NewOpenQrcodeContactWayDao().GetOne(c, data.EnterpriseId, "user_id = ? AND state = ?", memberInfo.Userid, state)
	if wayModel.QrCode != "" {
		return reserveResp.AssignGuideSuccessResp{
			StoreCode:     data.StoreCode,
			Type:          data.Type,
			MemberId:      assignMemberId,
			WechatUserId:  memberInfo.Userid,
			MemberName:    memberInfo.MemberName,
			ContactQrcode: wayModel.QrCode,
		}, nil
	}
	//生成联系我二维码
	contactQrcode, err := s.genContactQrcode(c, data.EnterpriseId, assignMemberId, 0, state)
	if err != nil || contactQrcode.Errcode != 0 {
		if contactQrcode != nil && contactQrcode.Errcode != 0 {
			err = errors.New(contactQrcode.Errmsg)
		}
		s.Log().Error(c).Msgf("生成联系我二维码,异常：%v", err)
		return resp, err
	}
	//入库
	_ = openDao.NewOpenQrcodeContactWayDao().CreateData(c, data.EnterpriseId, openModel.OpenQrcodeContactWayModel{
		BaseModel: models.CreateBase(app.UID()),
		State:     state,
		UserId:    memberInfo.Userid,
		ConfigId:  contactQrcode.ConfigId,
		QrCode:    contactQrcode.QrCode,
	})

	resp = reserveResp.AssignGuideSuccessResp{
		StoreCode:     data.StoreCode,
		Type:          data.Type,
		MemberId:      assignMemberId,
		WechatUserId:  memberInfo.Userid,
		MemberName:    memberInfo.MemberName,
		ContactQrcode: contactQrcode.QrCode,
	}

	return resp, nil
}

// GetRelation 获取预约关系
func (s ReserveService) GetRelation(c app.Ctx, data reserveReq.GetRelationRequest) (reserveResp.Success, error) {
	var resp reserveResp.Success

	reserveRecordDao := thirdDao.NewReserveRecordDao()
	// 查询预约记录
	recordInfo, err := reserveRecordDao.GetInfo(c, thirdModel.ReserveRecordModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: data.EnterpriseId,
		},
		ReserveNumberID: data.ReserveNumberId,
	})
	if err != nil || recordInfo.ID == 0 {
		return resp, errors.New("获取预约信息异常")
	}

	// 获取员工信息
	membersInfo, err := s.membersDao.GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{
			BaseModel: models.BaseModel{
				ID: recordInfo.MemberID,
			},
			EnterpriseID: recordInfo.EnterpriseID,
		},
	})
	// 获取员工和部门的关系
	memberDepartmentInfo, err := membersDao.NewMemberDepartments().GetInfo(c, membersModel.MemberDepartmentsModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: recordInfo.EnterpriseID,
		},
		MemberID:     recordInfo.MemberID,
		DepartmentID: recordInfo.DepartmentID,
	})

	// 判断导购离职和转店/转岗
	if membersInfo.ID > 0 || membersInfo.Status == membersModel.MembersStatusEnable || memberDepartmentInfo.ID >= 0 {
		resp = reserveResp.Success{
			StoreCode:       recordInfo.StoreCode,
			ReserveNumberId: recordInfo.ReserveNumberID,
			MemberId:        recordInfo.MemberID,
			ContactQrcode:   recordInfo.ContactQrcode,
		}
		return resp, err
	}

	// 获取部门下的所有员工Id集合
	memberIds, err := s.getDepartmentMemberIds(c, data.EnterpriseId, recordInfo.DepartmentID)
	if err != nil {
		return resp, err
	}
	// 获取分配的员工id集合
	assignMembersIds, err := s.getAssignMemberIds(c, recordInfo.EnterpriseID, recordInfo.DepartmentID, memberIds)
	if err != nil {
		return resp, err
	}

	assignMemberId := assignMembersIds[0]

	reserveRecordAssignDao := thirdDao.NewReserveRecordAssignDao()
	// 开启事务
	tx := app.DB().WithContext(c).Begin()

	// 删除员工离职或转店的预约关系
	err = reserveRecordAssignDao.WithTx(tx).Delete(c, thirdModel.ReserveRecordAssignModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: recordInfo.EnterpriseID,
		},
		DepartmentID: recordInfo.DepartmentID,
		MemberID:     recordInfo.MemberID,
	})
	if err != nil {
		s.Log().Error(c).Msgf("删除员工离职或转店的预约关系错误,异常：%v", err)
		tx.Rollback()
		return resp, err
	}

	// 添加分配记录
	reserveRecordAssignSaveData := thirdModel.ReserveRecordAssignModel{
		MicroBaseModel: models.CreateMicroBase(data.EnterpriseId, app.UIDContext(c)),
		DepartmentID:   recordInfo.DepartmentID,
		MemberID:       assignMemberId,
		AssignNumber:   1,
	}
	err = reserveRecordAssignDao.WithTx(tx).Create(c, reserveRecordAssignSaveData)
	if err != nil {
		s.Log().Error(c).Msgf("添加分配记录,异常：%v", err)
		tx.Rollback()
		return resp, err
	}

	// 生成联系我二维码
	contactQrcode, err := s.genContactQrcode(c, data.EnterpriseId, assignMemberId, recordInfo.ID, "")
	if err != nil || contactQrcode.Errcode != 0 {
		if contactQrcode != nil && contactQrcode.Errcode != 0 {
			err = errors.New(contactQrcode.Errmsg)
		}
		s.Log().Error(c).Msgf("生成联系我二维码,异常：%v", err)
		tx.Rollback()
		return resp, err
	}

	// 更新预约记录
	reserveRecordUpdateData := thirdModel.ReserveRecordModel{
		MicroBaseModel: models.MicroBaseModel{
			BaseModel: models.BaseModel{
				ID: recordInfo.ID,
			},
			EnterpriseID: recordInfo.EnterpriseID,
		},
		MemberID:      assignMemberId,
		ConfigID:      contactQrcode.ConfigId,
		ContactQrcode: contactQrcode.QrCode,
	}
	err = reserveRecordDao.WithTx(tx).Save(c, reserveRecordUpdateData)
	if err != nil {
		s.Log().Error(c).Msgf("更新预约记录,异常：%v", err)
		tx.Rollback()
		return resp, err
	}

	tx.Commit()

	info, err := s.membersDao.GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: data.EnterpriseId,
			BaseModel: models.BaseModel{
				ID: assignMemberId,
			},
		},
	})
	if err != nil {
		s.Log().Error(c).Msgf("查询导购信息,异常：%v", err)
		return resp, err
	}
	resp = reserveResp.Success{
		StoreCode:       recordInfo.StoreCode,
		ReserveNumberId: recordInfo.ReserveNumberID,
		MemberId:        reserveRecordUpdateData.MemberID,
		MemberName:      info.Name,
		ContactQrcode:   reserveRecordUpdateData.ContactQrcode,
	}

	return resp, nil
}

// 通过门店编码获取部门信息
func (s ReserveService) getDepartmentByStoreCode(c app.Ctx, enterpriseId int64, StoreCode string) (departmentsModels.DepartmentsModel, error) {
	var department departmentsModels.DepartmentsModel
	// 通过store_code查询部门id
	storeDepartmentInfo, err := microWorkDao.NewStoreDepartment().GetInfo(c, microWorkModel.StoreDepartmentModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: enterpriseId,
		},
		ThirdStoreCode: StoreCode,
	})

	if err != nil {
		s.Log().Error(c).Msgf("通过store_code查询部门,异常：%v", err)
		return department, errors.New("门店异常")
	}

	if storeDepartmentInfo.ID == 0 {
		return department, errors.New("门店不存在")
	}

	// 查询部门主键id
	department, err = s.departmentsDao.GetInfo(c, departmentsModels.DepartmentsModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: enterpriseId,
		},
		ThirdID: int(storeDepartmentInfo.WechatDepartmentID),
	})

	if department.ID == 0 {
		return department, errors.New("部门未找到")
	}

	return department, err
}

// 获取指定部门下的员工id集合
func (s ReserveService) getDepartmentMemberIds(c app.Ctx, enterpriseId int64, departmentId int64) ([]int64, error) {
	// 获取部门下的员工
	departmentMembers, err := s.departmentsDao.GetDepartmentsMembersByDepartmentId(c, enterpriseId, departmentId)
	if err != nil {
		s.Log().Error(c).Msgf("获取部门下的员工,异常：%v", err)
		return nil, err
	}

	if len(departmentMembers) == 0 {
		return nil, errors.New("部门下没有员工")
	}

	var memberIds []int64
	for _, member := range departmentMembers {
		memberIds = append(memberIds, member.MemberId)
	}
	return memberIds, nil
}

// 获取分配的员工id集合
func (s ReserveService) getAssignMemberIds(c app.Ctx, enterpriseId int64, departmentId int64, departmentMemberIds []int64) ([]int64, error) {
	// 查询该部门下已绑定的员工
	reserveRecordAssignDao := thirdDao.NewReserveRecordAssignDao()

	reserveRecordAssignResults, err := reserveRecordAssignDao.GetAll(c, thirdModel.ReserveRecordAssignModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: enterpriseId,
		},
		DepartmentID: departmentId,
	})
	if err != nil {
		return nil, err
	}

	// 根据分配次数升序排列
	sort.Slice(reserveRecordAssignResults, func(i, j int) bool {
		return reserveRecordAssignResults[i].AssignNumber < reserveRecordAssignResults[j].AssignNumber
	})

	// 已绑定的员工id集合
	var bindMemberIds []int64
	for _, reserveRecordAssign := range reserveRecordAssignResults {
		// 过滤离职员工
		if !helper.InArray[int64](reserveRecordAssign.MemberID, departmentMemberIds) {
			continue
		}
		bindMemberIds = append(bindMemberIds, reserveRecordAssign.MemberID)
	}
	// 取出部门下的所有员工和已绑定的员工和的差集
	assignMembers := helper.ArrayDiff(departmentMemberIds, bindMemberIds)

	// 如果没有差集 则返回已经绑定过的员工集合
	if len(assignMembers) == 0 {
		assignMembers = bindMemberIds
	}

	return assignMembers, nil
}

func (s ReserveService) getAssignGuideState(c app.Ctx, enterpriseId int64, guideType int) (string, error) {
	var sourceType int
	switch guideType {
	case AssignGuideTypeStoreWorker:
		sourceType = drainageContactModel.DrainageContactGroupSourceStoreWorker
	case AssignGuideTypeReserveGuide:
		sourceType = drainageContactModel.DrainageContactGroupSourceReserveGuide
	default:
		return "", errors.New("无效的二维码场景类型")
	}

	groupInfo, err := drainageContactDao.NewDrainageContactGroupDao().GetOneBySourceType(c, enterpriseId, sourceType)
	if err != nil {
		return "", err
	}
	if groupInfo.ID == 0 {
		return assignGuideDefaultState, nil
	}
	contactInfo, err := drainageContactDao.NewDrainageContactsDao().GetFirstByGroupId(c, enterpriseId, groupInfo.ID)
	if err != nil {
		return "", err
	}
	if contactInfo.ID == 0 {
		return assignGuideDefaultState, nil
	}
	return helper.Int64ToString(contactInfo.ID), nil
}

// 生成联系我的二维码
func (s ReserveService) genContactQrcode(c app.Ctx, enterpriseId int64, memberId, channelId int64, state string) (*work.GetWxCodeResult, error) {
	// 获取员工信息
	membersInfo, err := s.membersDao.GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{
			BaseModel: models.BaseModel{
				ID: memberId,
			},
			EnterpriseID: enterpriseId,
		},
	})
	if err != nil || membersInfo.ID == 0 {
		return nil, errorx.NewErrorMsg(1001, "获取员工信息失败或员工信息不存在")
	}

	request := work.ContactCreateRequest{
		ID:           channelId,
		CodeCame:     "三方预约二维码",
		Type:         1, // 1-单人, 2-多人
		AutoFriend:   1,
		UserIds:      []string{membersInfo.Userid},
		EnterpriseId: enterpriseId,
		State:        state,
	}

	return work.NewQyWorkImpl().SetEnterprise(enterpriseId).GetWxEexternalContact(ctx.New(), request)
}

func (s ReserveService) UserReserveList(c ctx.Context, request reserveReq.UserReserveListReq) ([]reserveReq.UserReserveListRes, int64, error) {
	var res []reserveReq.UserReserveListRes
	result, total, err := workExternalDao.NewReservesOrder().GetUserReserveOrderLists(c, request)
	if err != nil {
		return res, 0, err
	}
	if total == 0 {
		return res, 0, err
	}
	// 查询服务信息
	var listReq reserveRequest.ReserveServiceListReq
	serviceList, _, errs := workExternalDao.NewReserveServiceDao().GetListAll(c, request.EnterpriseId, listReq)
	if errs != nil {
		return nil, total, errs
	}
	serviceListMap := make(map[int64]microScrmModel.ReserveServiceModel, 0)
	for _, item := range serviceList {
		serviceListMap[item.ID] = item
	}
	for _, v := range result {
		info := v
		reserveStartTime := time.Unix(int64(info.ReserveTime), 10).Format("2006-01-02 15:04")
		reserveEndTimeUnix := info.ReserveTime + 3600
		reserveEndTime := time.Unix(int64(reserveEndTimeUnix), 10).Format("15:04")
		item := reserveReq.UserReserveListRes{
			ID:              info.ID,
			CreatedAt:       info.CreatedAt,
			ReserveNumber:   info.ReserveNumber,
			ReserveTimeStr:  reserveStartTime + "-" + reserveEndTime,
			StoreName:       info.StoreName,
			StoreAddress:    info.StoreAddress,
			Status:          info.Status,
			UserName:        info.UserName,
			UserUnionid:     info.UserUnionid,
			MemberID:        info.MemberID,
			MemberName:      info.MemberName,
			ReserveRemark:   info.ReserveRemark,
			UserRemark:      info.UserRemark,
			Remark:          info.ReserveRemark,
			OrderOrigin:     info.OrderOrigin,
			ServiceDuration: int(serviceListMap[info.ReserveID].ReserveTime),
		}
		res = append(res, item)
	}
	return res, total, err
}
