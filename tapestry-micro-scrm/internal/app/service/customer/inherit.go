package customerSvc

import (
	"errors"
	"micro-scrm/internal/app"
	customerDao "micro-scrm/internal/app/dao/customer"
	membersDao "micro-scrm/internal/app/dao/members"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	workDao "micro-scrm/internal/app/dao/work"
	workGroupDao "micro-scrm/internal/app/dao/work_group"
	customerReq "micro-scrm/internal/app/http/request/customer"
	memberReq "micro-scrm/internal/app/http/request/member"
	"micro-scrm/internal/app/models"
	customerModels "micro-scrm/internal/app/models/customer"
	membersModel "micro-scrm/internal/app/models/members"
	workModels "micro-scrm/internal/app/models/work"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/work"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gorm.io/gorm"
	"micro-scrm/internal/pkg/uid"
)

type Inherit struct {
	service.Base
	inheritDao     *customerDao.Inherit
	chatInheritDao *customerDao.ChatInherit
	memberDao      *membersDao.Members
	workGroupDao   *workGroupDao.WorkGroupChat
	qyWorkImpl     *work.QyWorkImpl
}

func NewInherit() *Inherit {
	return &Inherit{
		inheritDao:     customerDao.NewInherit(),
		chatInheritDao: customerDao.NewChatInherit(),
		memberDao:      membersDao.NewMembers(),
		workGroupDao:   workGroupDao.NewWorkGroupChat(),
		qyWorkImpl:     work.NewQyWorkImpl(),
	}
}

// GetInheritList 获取离职员工列表
func (s Inherit) GetInheritList(c app.Ctx, request customerReq.CustomerListReq) ([]customerModels.ResponseCustomerInheritList, int64, error) {
	// 列表
	result, total, err := s.inheritDao.GetInheritList(c, request)
	if err != nil {
		return nil, 0, nil
	}
	var userIds []string
	for _, v := range result {
		userIds = append(userIds, v.HandoverUserid)
	}

	// 未分配社群
	resultGroup, err := s.workGroupDao.GetUserGroupChatNumForUserId(c, request.EnterpriseID, userIds)
	if err != nil {
		return nil, 0, err
	}
	unAssignedGroupMap := make(map[string]int, 0)
	for _, user := range resultGroup {
		unAssignedGroupMap[user.Userid] = user.Num
	}

	// 已分配社群
	var assignedGroupNumList []customerReq.InsertChatNumList
	assignedGroupMap := make(map[string]int, 0)
	assignedGroupNumList, err = customerDao.NewChatInherit().GetChatCountNumForUserid(c, request.EnterpriseID, userIds, 1)
	if err != nil {
		return nil, 0, err
	}
	for _, user := range assignedGroupNumList {
		assignedGroupMap[user.HandoverUserid] = user.Num
	}

	// 查询群数量
	for k, v := range result {
		if _, ok := unAssignedGroupMap[v.HandoverUserid]; ok {
			result[k].UnassignedGroupNum = unAssignedGroupMap[v.HandoverUserid]
		}
		if _, ok := assignedGroupMap[v.HandoverUserid]; ok {
			result[k].AssignedGroupNum = assignedGroupMap[v.HandoverUserid]
		}
		result[k].HandoverMobile = helper.HandleMobile(v.HandoverMobile)
	}

	return result, total, err
}
func (s Inherit) GetExternalAllotList(c app.Ctx, request customerReq.ExternalListReq) ([]customerReq.ExternalListRes, int64, error) {
	data, total, err := s.inheritDao.GetExternalAllotList(c, request)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, nil
	}
	return data, total, err
}

func (s Inherit) GetExternalAllotCount(c app.Ctx, request customerReq.ExternalListReq) (int64, error) {
	return s.inheritDao.GetExternalAllotCount(c, request)
}

func (s Inherit) GetGroupChatAllotList(c app.Ctx, request customerReq.GroupChatList) (data []customerReq.GroupChatInfo, err error) {
	workImpl := s.qyWorkImpl.SetEnterprise(request.EnterpriseID)
	var groupRequest = work.ExternalGroupChatListRequest{}
	groupRequest.OwnerFilter.UserIDList = []string{request.HandoverUserId}
	groupRequest.Limit = 1000
	groupRequest.Offset = 0
	groupRequest.StatusFilter = 1
	groupList, groupErr := workImpl.ExternalGroupChatList(c, groupRequest)
	if groupErr != nil {
		return nil, groupErr
	}
	var groups []string
	for _, item := range groupList {
		groups = append(groups, item.ChatID)
	}
	if len(groups) == 0 {
		return []customerReq.GroupChatInfo{}, nil
	}
	data, err = workDao.NewGroupChat().GetGroupChatAllotList(c, request.EnterpriseID, groups)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return data, nil
}

func (s Inherit) GetTakeOverList(c app.Ctx, request customerReq.GroupChatList) (data []customerReq.TakeOverInfo, err error) {
	data, err = s.inheritDao.GetTakeOverList(c, request)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return data, nil
}

func (s Inherit) GetUniqueStr(c app.Ctx, enterpriseID int64, newList [][]string) (data []string, err error) {
	return s.inheritDao.GetUniqueStr(c, enterpriseID, newList)
}

func (s Inherit) GetListByUserIds(c app.Ctx, enterpriseID int64, userIds []string) (list []customerModels.CustomerInheritModel, err error) {
	return s.inheritDao.GetListByUserIds(c, enterpriseID, userIds)
}

func (s Inherit) InsertNewData(c app.Ctx, data []customerModels.CustomerInheritModel) error {
	return s.inheritDao.InsertNewData(c, data)
}
func (s Inherit) UpdateDataByUnique(c app.Ctx, data [][]string) error {
	return s.inheritDao.UpdateDataByUnique(c, data)
}

func (s Inherit) UpdateByIds(c app.Ctx, enterpriseID int64, ids []int64) error {
	return s.inheritDao.UpdateByIds(c, enterpriseID, ids)
}

func (s Inherit) DelInvalidData(c app.Ctx, enterpriseId int64) error {
	return s.inheritDao.DelInvalidData(c, enterpriseId)
}

func (s Inherit) UpdateIncompleteData(c app.Ctx, enterpriseID int64) error {
	err := s.inheritDao.UpdateHandoverData(c, enterpriseID)
	corpExternal := microScrmDao.NewWxUserCorpExternal()
	if err != nil {
		return err
	}
	var lastID int64
	for {
		list, err := s.inheritDao.GetCustomerLimit(c, enterpriseID, lastID, 100)
		if err != nil {
			return err
		}
		if len(list) == 0 {
			break
		}
		var openIds []string
		for _, model := range list {
			openIds = append(openIds, model.ExternalUserid)
		}

		externalList, err := corpExternal.GetMany(c, enterpriseID, "openid in ?", openIds)
		if err != nil {
			return err
		}
		for _, model := range externalList {
			_ = s.inheritDao.UpdateData(c, enterpriseID, map[string]any{
				"external_userid":   model.Openid,
				"external_nickname": "",
			}, map[string]any{
				"external_avatar":   model.Avatar,
				"external_nickname": model.Nick,
			})
		}

		lastID = list[len(list)-1].ID
	}
	return err
}

func (s Inherit) AllocationCustomer(ctx app.Ctx, data memberReq.ResignedTransferCustomerReq) ([]string, []string, error) {
	requestQyData := work.ResignedTransferCustomerReq{
		HandoverUserid: data.HandoverUserid,
		TakeoverUserid: data.TakeoverUserid,
		ExternalUserid: data.ExternalUserid,
	}

	resp, err := s.qyWorkImpl.SetEnterprise(data.EnterpriseID).ResignedTransferCustomer(ctx, requestQyData)
	if err != nil {
		return nil, nil, err
	}

	// 成功，失败分配的用户id集合
	var successUserIds, errorUserIds []string
	for _, user := range resp.Customer {
		if user.Errcode == 0 {
			successUserIds = append(successUserIds, user.ExternalUserid)
		} else {
			errorUserIds = append(errorUserIds, user.ExternalUserid)
		}
	}

	// 查询接替人的部门
	memberDepartment, err := s.memberDao.GetMemberInfoAndDepartmentsByUserId(ctx, data.EnterpriseID, []string{data.TakeoverUserid})
	if err != nil {
		return nil, nil, err
	}

	updateData := customerModels.CustomerInheritModel{
		TakeoverUserid:         data.TakeoverUserid,
		TakeoverName:           memberDepartment[0].MembersName,
		TakeoverDepartmentName: memberDepartment[0].DepartmentName,
		Status:                 customerModels.InheritStatusAssigned,
	}
	// 更新接替人和状态
	err = customerDao.NewInherit().ChangeTakeoverUserByHandoverIdsExternalUserIds(ctx, data.EnterpriseID, data.HandoverUserid, successUserIds, updateData)
	if err != nil {
		return nil, nil, err
	}

	return successUserIds, errorUserIds, err
}

func (s Inherit) AllocationGroup(ctx app.Ctx, data memberReq.ResignedTransferCustomerGroupReq) ([]string, []string, error) {
	requestQyData := work.ResignedTransferCustomerGroupReq{
		ChatIdList: data.ChatIdList,
		NewOwner:   data.TakeoverUserid,
	}

	resp, err := s.qyWorkImpl.SetEnterprise(data.EnterpriseID).ResignedTransferCustomerGroup(ctx, requestQyData)
	if err != nil {
		return nil, nil, err
	}

	// 查询员工的部门
	memberDepartment, err := s.memberDao.GetMemberInfoAndDepartmentsByUserId(ctx, data.EnterpriseID, []string{data.HandoverUserid, data.TakeoverUserid})
	if err != nil {
		return nil, nil, err
	}
	memberDepartmentIndex := make(map[string]membersModel.MemberDepartment)
	for _, department := range memberDepartment {
		memberDepartmentIndex[department.MembersUserid] = department
	}

	// 查询群组名字
	groupChat, err := workDao.NewGroupChat().GetListByChatIds(ctx, data.EnterpriseID, data.ChatIdList)
	if err != nil {
		return nil, nil, err
	}
	groupChatIndex := make(map[string]workModels.GroupChatModel)
	for _, chat := range groupChat {
		groupChatIndex[chat.ChatID] = chat
	}

	var (
		createData           []customerModels.CustomerChatInheritModel
		successIds, errorIds []string
	)

	// 判断企微返回是否存在分配失败群组
	if len(resp.FailedChatList) == 0 {
		successIds = data.ChatIdList
		for _, chat := range data.ChatIdList {
			tmpData := customerModels.CustomerChatInheritModel{
				MicroBaseModel:         models.CreateMicroBase(data.EnterpriseID, uid.ID(ctx)),
				ChatID:                 chat,
				ChatName:               groupChatIndex[chat].Name,
				TakeoverUserid:         data.TakeoverUserid,
				TakeoverName:           memberDepartmentIndex[data.TakeoverUserid].MembersName,
				TakeoverDepartmentName: memberDepartmentIndex[data.TakeoverUserid].DepartmentName,
				HandoverUserid:         data.HandoverUserid,
				HandoverName:           memberDepartmentIndex[data.HandoverUserid].MembersName,
				HandoverDepartmentName: memberDepartmentIndex[data.HandoverUserid].DepartmentName,
				Status:                 customerModels.InheritChatStatusSuccess,
			}
			createData = append(createData, tmpData)
		}
	} else {
		for _, chat := range resp.FailedChatList {
			tmpData := customerModels.CustomerChatInheritModel{
				MicroBaseModel:         models.CreateMicroBase(data.EnterpriseID, uid.ID(ctx)),
				ChatID:                 chat.ChatId,
				ChatName:               groupChatIndex[chat.ChatId].Name,
				TakeoverUserid:         data.TakeoverUserid,
				TakeoverName:           memberDepartmentIndex[data.TakeoverUserid].MembersName,
				TakeoverDepartmentName: memberDepartmentIndex[data.TakeoverUserid].DepartmentName,
				HandoverUserid:         data.HandoverUserid,
				HandoverName:           memberDepartmentIndex[data.HandoverUserid].MembersName,
				HandoverDepartmentName: memberDepartmentIndex[data.HandoverUserid].DepartmentName,
				Status:                 customerModels.InheritChatStatusError,
			}
			errorIds = append(errorIds, chat.ChatId)
			createData = append(createData, tmpData)
		}
	}

	// 添加群组分配记录
	if err = customerDao.NewChatInherit().CreateInBatches(ctx, createData); err != nil {
		return nil, nil, err
	}

	return successIds, errorIds, err
}

func (s Inherit) GetTakeOverExternalList(ctx app.Ctx, request customerReq.ExternalListReq) (int64, []customerReq.TakeOverExternalInfo, error) {
	requestQyData := work.GetResignedTransferResultRequest{
		HandoverUserid: request.HandoverUserId,
		TakeoverUserid: request.TakeoverUserId,
		Cursor:         "",
	}
	resp, err := s.qyWorkImpl.SetEnterprise(request.EnterpriseID).GetResignedTransferResult(ctx, requestQyData)
	if respErr, ok := err.(*errorx.Error); ok {
		if respErr.Code == 40003 || respErr.Code == 60111 { //微信提示无效的userid,有可能是因为接替人离职，调用企微接口报错,所以再去本地数据库查下分配的记录
			total, lists, errs := s.GetTakeOverExternalListByLocal(ctx, request)
			if errs != nil {
				return 0, nil, errs
			}
			if total > 0 {
				var result []customerReq.TakeOverExternalInfo
				for _, v := range lists {
					Info := customerReq.TakeOverExternalInfo{
						TakeOverTime:       0,
						ExternalStatusName: "接替人已离职",
						ExternalListRes:    v,
					}
					result = append(result, Info)
				}
				return total, result, nil
			}
			return 0, []customerReq.TakeOverExternalInfo{}, nil
		}
	}
	if err != nil {
		return 0, nil, err
	}
	var ExternalUserIds []string
	ExternalStatusMap := make(map[string]int, 0)
	TakeOverTimeMap := make(map[string]int, 0)
	for _, item := range resp.Customer {
		ExternalUserIds = append(ExternalUserIds, item.ExternalUserid)
		ExternalStatusMap[item.ExternalUserid] = item.Status
		TakeOverTimeMap[item.ExternalUserid] = item.TakeoverTime
	}
	total := len(ExternalUserIds)
	if len(ExternalUserIds) == 0 {
		return 0, []customerReq.TakeOverExternalInfo{}, nil
	}
	start := (request.Page - 1) * request.Size
	end := request.Page * request.Size
	var Ids []string
	if end >= total {
		Ids = ExternalUserIds[start:total]
	} else {
		Ids = ExternalUserIds[start:end]
	}
	list, errs := s.inheritDao.GetTakeOverExternalList(ctx, request, Ids)
	if errs != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil, errs
	}
	for k, v := range list {
		switch ExternalStatusMap[v.ExternalUserid] {
		case 1:
			list[k].ExternalStatusName = "接替完毕"
		case 2:
			list[k].ExternalStatusName = "等待接替"
		case 3:
			list[k].ExternalStatusName = "客户拒绝"
		case 4:
			list[k].ExternalStatusName = "接替成员客户达到上限"
		}
		list[k].TakeOverTime = TakeOverTimeMap[v.ExternalUserid]
	}
	return int64(total), list, nil
}
func (s Inherit) GetTakeOverExternalListByLocal(c app.Ctx, request customerReq.ExternalListReq) (int64, []customerReq.ExternalListRes, error) {
	return s.inheritDao.GetTakeOverExternalListByLocal(c, request)
}
func (s Inherit) GetTakeOverGroupList(c app.Ctx, request customerReq.GroupChatList) (int64, []customerReq.TakeOverGroupInfo, error) {
	total, err := s.chatInheritDao.GetTakeOverGroupCount(c, request)
	if err != nil {
		return 0, nil, err
	}
	if total <= 0 {
		return 0, nil, nil
	}
	list, errs := s.chatInheritDao.GetTakeOverGroupList(c, request)
	if errs != nil && !errors.Is(errs, gorm.ErrRecordNotFound) {
		return total, nil, errs
	}
	return total, list, nil
}

func (s Inherit) GetTakeOverCount(c app.Ctx, enterpriseID int64, UserID string) (int64, error) {
	return s.inheritDao.GetTakeOverCount(c, enterpriseID, UserID)
}
