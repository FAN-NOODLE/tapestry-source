package callback

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	activityQRCodeDao "micro-scrm/internal/app/dao/activity_qrcode"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	drainageContactDao "micro-scrm/internal/app/dao/drainage"
	membersDao "micro-scrm/internal/app/dao/members"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	tagsDao "micro-scrm/internal/app/dao/tags"
	thirdDao "micro-scrm/internal/app/dao/third"
	"micro-scrm/internal/app/models"
	activityQRCodeModel "micro-scrm/internal/app/models/activity_qrcode"
	drainageContactModel "micro-scrm/internal/app/models/drainage"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	thirdModel "micro-scrm/internal/app/models/third"
	cacheSvc "micro-scrm/internal/app/service/cache"
	membersSvc "micro-scrm/internal/app/service/member"
	"micro-scrm/internal/config"
	drainage2 "micro-scrm/internal/orm/model/micro_scrm/drainage"
	welcome2 "micro-scrm/internal/orm/model/micro_scrm/welcome"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/redis"
	"micro-scrm/internal/pkg/sdk"
	"micro-scrm/internal/pkg/sdk/boss"
	"micro-scrm/internal/pkg/sdk/microwork"
	"strconv"
	"strings"
	"time"

	"micro-scrm/internal/orm/model/micro_scrm/work_external_user"

	"micro-scrm/internal/events"
	"micro-scrm/internal/events/entity"
	"micro-scrm/internal/logic/fission/welcome"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"

	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/drainage"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/modules/fission"
	"micro-scrm/internal/modules/members"
	"micro-scrm/internal/modules/user"
	workexternal "micro-scrm/internal/modules/work-external"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/sdk/work"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/pkg/uid"
)

// Service  应用服务
type Service struct {
	MemberService      *members.Service
	UserService        *user.Service
	CustomerService    *customer.Service
	MessageService     *customer.MessageService
	ExternalTagService *customer.ExternalTagService
	DrainageService    *drainage.DrainageService
	ContactService     *drainage.ContactService
	ExternalService    ExternalService
	FissionService     *fission.Service
}

const (
	SystemWelcomeMessage        = 1 // 系统原有欢迎语
	AdminActivityWelcomeMessage = 2 // 管理员创建裂变活动的欢迎语
	ActivityWelcomeMessage      = 3 // 用户裂变活动裂变的欢迎语
)

// ChangeExternalContact 外部联系人通讯录变更事件处理
func (s Service) ChangeExternalContact(ctx ctx.Context, callbackData *work.MemberExternalCallback,
	enterpriseData *enterprise.Corp, workImpl *work.QyWorkImpl, addContactsChannelId int64) error {
	// 查询成员是否存在
	mReq := members.Member{
		Userid:       callbackData.UserID,
		EnterpriseID: enterpriseData.EnterpriseID,
	}
	member := s.MemberService.FindMemberByUserid(ctx, mReq)
	if member == nil {
		log.Info(ctx).Msgf("新增联系人验证回调,未查询到企业成员%s, 查询信息：%+v", mReq.Userid, mReq)
		var tmpErr error
		member, tmpErr = s.MemberService.CreateMemberOnLock(ctx, mReq, workImpl)
		if tmpErr != nil {
			return tmpErr
		}
		//如果新增用户职位符合要求，添加【会话存档_正式】标签
		if helper.InArray(member.Position, membersDao.ConversationArchiveTag) {
			storeThirdIds := s.getStoreThirdIds(ctx, enterpriseData.EnterpriseID)
			if len(storeThirdIds) > 0 && len(member.Department) > 0 && len(helper.ArrayIntersect(storeThirdIds, member.Department)) > 0 {
				conf := config.CfgData.CoachSaTag
				_ = workImpl.AddTagUsers(ctx, work.AddTagUsersRequest{
					QyGwRequest: work.QyGwRequest{ReqPL: "communalContact"},
					TagID:       conf.TagId,
					UserList:    []string{member.Userid},
				})
			}
		}
	}
	// 调用合并接口，处理合并用户逻辑
	source := &user.SourceUser{
		EnterpriseID: enterpriseData.EnterpriseID,
		SourceUserBase: customer.SourceUserBase{
			AppID:        callbackData.ToUserName,
			OpenID:       callbackData.ExternalUserID,
			CustomerType: customer.CustomerTypeExternal,
		},
		CreateAt: callbackData.CreateTime,
	}

	var dealErr error
	var relationId int64
	var externalDetail *work.ExternalContactResponse
	switch callbackData.ChangeType {
	case AddExternalContact, MigrateExternalContact, AddHalfExternalContact, EditExternalContact:
		// 从企业微信查询外部联系人详情
		externalDetailTmp, err := workImpl.ExternalDetail(ctx, work.ExternalContactRequest{
			ExternalUserID: callbackData.ExternalUserID,
			QyGwRequest:    work.QyGwRequest{},
		})
		if err != nil {
			log.Info(ctx).Msgf("新增联系人验证回调,企微查询的外部联系人信息失败：%v", err)
			return err
		}
		externalDetail = externalDetailTmp
		log.Info(ctx).Msgf("新增联系人验证回调,查询到的外部联系人信息：%+v %v", externalDetail, err)
		source.UnionID = externalDetail.Base.UnionID
		source.Nick = externalDetail.Base.Name
		//source.Name = externalDetail.Base.Name	//name用三方CRM的姓名填充
		source.Avatar = externalDetail.Base.Avatar
		source.Sex = externalDetail.Base.Gender
		source.Position = externalDetail.Base.Position
		source.CorpName = externalDetail.Base.CorpName
		source.CorpFullName = externalDetail.Base.CorpFullName
		source.Type = externalDetail.Base.Type
		source.ExternalProfile = string(externalDetail.Base.ExternalProfile)

		if _, err := s.UserService.MergeUser(ctx, source); err != nil {
			log.Info(ctx).Msgf("新增联系人验证回调,SynWxShopUser合并用户数据MergeUser失败,err:%#v,data:%+v", err.Error(), source)
			return err
		}

		// 获取回调渠道id(联系我添加好友数量)
		state, _ := strconv.ParseInt(callbackData.State, 10, 64)
		s.handleChannelAddUser(ctx, source.ID, state, enterpriseData.EnterpriseID)

		log.Info(ctx).Msgf("新增联系人验证回调,查询到的外部联系人信息source：%+v %v", source, err)
		if callbackData.ChangeType == AddExternalContact || callbackData.ChangeType == MigrateExternalContact || callbackData.ChangeType == EditExternalContact {
			// 外部联系人关系入库处理
			relationModel, err := s.dealBatchInsertRelations(ctx, externalDetail, source, callbackData)
			if err != nil {
				log.Error(ctx).Msgf("新增联系人验证回调,处理外部联系人关系数据失败，错误信息：%+v", err)
			}
			relationId = relationModel.ID
			// 添加足迹
			if externalDetail.Base.UnionID != "" && (callbackData.ChangeType != EditExternalContact) {
				eventDesc := "员工：" + member.Name + "（" + member.Userid + "）"
				_ = s.UserService.AddUserEvent(ctx, enterpriseData.EnterpriseID, source.EuID, string(common.EventEnum.AddExternalUser), eventDesc, "", nil, callbackData.CreateTime)
			}
			// 处理标签
			_ = s.UpdateExternalUserTag(ctx, source.ID, callbackData, externalDetail, enterpriseData)
		}
	default:
		if _, err := s.UserService.MergeUser(ctx, source); err != nil {
			log.Info(ctx).Msgf("新增联系人验证回调,SynWxShopUser合并用户数据MergeUser失败,err:%#v,data:%+v", err.Error(), source)
			return err
		}
	}

	switch callbackData.ChangeType {
	case AddExternalContact, MigrateExternalContact:
		// 查询该关系的channelID
		contactsChannelId, err := s.RelationChannelId(ctx, callbackData, enterpriseData.EnterpriseID)
		if contactsChannelId != 0 && err == nil {
			// 操作外部联系人状态表
			s.AddExternalTraceHandle(ctx, callbackData, enterpriseData.EnterpriseID)
			// 操作外部联系人渠道状态表
			s.AddExternalChannelTraceHandle(ctx, callbackData, enterpriseData.EnterpriseID, contactsChannelId)
			// 恢复外部联系人删除员工
			s.RecoveryDelRelationChannel(ctx, callbackData, enterpriseData.EnterpriseID, contactsChannelId)
			// 更新企业微信外部联系人回调记录
			callbackData.ChannelID = contactsChannelId
			_, _ = s.CustomerService.UpdateCustomerCallbackRecord(ctx, *callbackData)
		}
		// 穿透case向下执行一次
		fallthrough
	case AddHalfExternalContact:
		if callbackData.WelcomeCode != "" {
			//先处理活动欢迎语和投放渠道欢迎语
			if !s.HandleActivityQrWelcome(ctx, member, callbackData, enterpriseData, workImpl, externalDetail) {
				// 再处理联系我二维码活动欢迎语、渠道欢迎语
				if !s.HandleDrainageWelcomesWelcome(ctx, member, callbackData, enterpriseData, workImpl, externalDetail) {
					// 先处理三方欢迎语
					if !s.HandleThirdWelcome(ctx, member, callbackData, enterpriseData.Domain, workImpl, externalDetail) {
						//发送系统默认欢迎语
						sendWelcomeRes, err := s.SendDefaultWelcome(ctx, member, callbackData, enterpriseData.Domain, workImpl, externalDetail)
						if !sendWelcomeRes || err != nil {
							sdk.SendMessage(ctx, "发送欢迎语失败", "", member, callbackData, dealErr)
							dealErr = err
						}
					}
				}
			}

		}
		// 企微用户打标签
		err := s.UserTagHandler(ctx, callbackData, workImpl, relationId, enterpriseData.EnterpriseID)
		if err != nil {
			log.Error(ctx).Msgf("新增联系人验证回调,ChangeExternalContact UserTagHandler %+v", err)
		}
		//处理会员信息
		s.setCrmDataInfo(ctx, enterpriseData.EnterpriseID, externalDetail.Base.UnionID)
		break
	case DelFollowUser:
		dealErr = s.delFollowUser(ctx, enterpriseData, callbackData.CreateTime, source, member)
		if addContactsChannelId != 0 {
			// 操作外部联系人状态表
			s.DelExternalTraceHandle(ctx, source.ID, callbackData, enterpriseData.EnterpriseID)
			// 操作外部联系人渠道状态表
			s.DelExternalChannelTraceHandle(ctx, source.ID, callbackData, enterpriseData.EnterpriseID, addContactsChannelId)
			// 操作好友删除员工表
			s.DelRelationChannelTrace(ctx, callbackData, enterpriseData.EnterpriseID, addContactsChannelId)
			// 更新企业微信外部联系人回调记录
			callbackData.ChannelID = addContactsChannelId
			_, _ = s.CustomerService.UpdateCustomerCallbackRecord(ctx, *callbackData)
		}
		break
	case DelExternalContact:
		dealErr = s.delExternalContact(ctx, enterpriseData, callbackData.CreateTime, source, member)
		if addContactsChannelId != 0 {
			// 操作外部联系人状态表
			s.DelExternalTraceHandle(ctx, source.ID, callbackData, enterpriseData.EnterpriseID)
			// 操作外部联系人渠道状态表
			s.DelExternalChannelTraceHandle(ctx, source.ID, callbackData, enterpriseData.EnterpriseID, addContactsChannelId)
			// 更新状态ID
			callbackData.ChannelID = addContactsChannelId
			_, _ = s.CustomerService.UpdateCustomerCallbackRecord(ctx, *callbackData)
		}
		break
	}
	return dealErr
}

// UpdateExternalUserTag 更新用户标签
func (s Service) UpdateExternalUserTag(c ctx.Context, workExternalUserID int64, callbackData *work.MemberExternalCallback, externalDetail *work.ExternalContactResponse, enterpriseData *enterprise.Corp) error {
	ev := log.Info(c).Interface("更新用户标签,回调数据:", callbackData).Interface("外部联系人信息", externalDetail)
	defer func() {
		ev.Msg("更新用户标签,结束")
	}()

	if callbackData.ChangeType != EditExternalContact &&
		callbackData.ChangeType != AddExternalContact &&
		callbackData.ChangeType != MigrateExternalContact {
		return nil
	}

	var editFollowUser *work.ExternalFollowUser
	for _, v := range externalDetail.FollowUser {
		if v.UserID == callbackData.UserID {
			editFollowUser = &v
			break
		}
	}
	if editFollowUser == nil {
		ev.Interface("外部联系人关系数据", editFollowUser)
		return nil
	}
	// 获取外部联系人关系数据
	relation, err := s.ExternalTagService.GetExternalUserRelation(c, enterpriseData.EnterpriseID, workExternalUserID, callbackData.UserID)
	if err != nil {
		ev.AnErr("获取外部联系人关系数据", err)
		return err
	}
	if relation == nil {
		ev.Interface("未找到外部联系人关系数据", relation)
	}

	// 获取外部联系人关系对应的标签数据
	dbExternalTags, err := s.ExternalTagService.GetExternalUserTags(c, relation.ID)
	if err != nil {
		ev.Interface("获取外部联系人关系对应的标签数据", err)
		return err
	}
	var externalTagIDList []int64
	for _, v := range dbExternalTags {
		externalTagIDList = append(externalTagIDList, v.ExternalTagId)
	}
	// 通过标签id获取标签信息
	dbTagList, err := s.ExternalTagService.GetExternalTagInfo(c, externalTagIDList)
	if err != nil {
		ev.Interface("通过标签id获取标签信息", err)
		return err
	}
	var dbCompanyTagMap = make(map[string]*customer.WorkExternalTag)
	var dbCustomerTagMap = make(map[string]*customer.WorkExternalTag)
	var dbTagMapOfID = make(map[int64]*customer.WorkExternalTag)
	for _, v := range dbTagList {
		switch v.Type {
		case customer.ExternalTagTypeCompany:
			dbCompanyTagMap[v.TagID] = v
		case customer.ExternalTagTypeCustomer:
			dbCustomerTagMap[v.Name] = v
		}
		dbTagMapOfID[v.ID] = v
	}
	for _, v := range dbExternalTags {
		if _, ok := dbTagMapOfID[v.ExternalTagId]; !ok {
			continue
		}
		v.TagID = dbTagMapOfID[v.ExternalTagId].TagID
		v.TagName = dbTagMapOfID[v.ExternalTagId].Name
		v.Type = dbTagMapOfID[v.ExternalTagId].Type
	}
	// 筛选出数据库中后果不存在的标签关系
	var newCompanyTagList []*work.ExternalTag
	var newCustomerTagList []*work.ExternalTag
	for i, v := range editFollowUser.Tags {
		if v.Type == customer.ExternalTagTypeCompany {
			_, ok := dbCompanyTagMap[v.TagID]
			if !ok {
				newCompanyTagList = append(newCompanyTagList, &editFollowUser.Tags[i])
				continue
			} else {
				delete(dbCompanyTagMap, v.TagID)
				continue
			}
		} else {
			if _, ok := dbCustomerTagMap[v.TagName]; !ok {
				newCustomerTagList = append(newCustomerTagList, &editFollowUser.Tags[i])
				continue
			} else {
				delete(dbCustomerTagMap, v.TagName)
				continue
			}
		}
	}
	err = s.AddExternalUserCompanyTag(c, relation, newCompanyTagList, enterpriseData)
	if err != nil {
		ev.Interface("添加企业标签", err)
	}
	err = s.AddExternalUserCustomerTag(c, relation, newCustomerTagList)
	if err != nil {
		log.Info(c).Msgf("添加自定义标签", err)
	}
	err = s.DeleteExternalUserTag(c, dbExternalTags, dbCompanyTagMap, dbCustomerTagMap)
	if err != nil {
		log.Info(c).Msgf("删除外部联系人标签", err)
	}

	// 不同客户的标签和客户分组
	events.ScrmEvent.ChangeUserTag.SyncDispatch(entity.Data{
		Ctx:          c,
		EnterpriseID: relation.EnterpriseID,
		RelationID:   relation.ID,
	})

	return err
}

// delFollowUser 成员被外部联系人删除
func (s Service) delFollowUser(ctx ctx.Context, enterpriseData *enterprise.Corp, msgTime int64, source *user.SourceUser, memberData *members.Member) error {
	log.Info(ctx).Msgf("成员被外部联系人删除事件处理开始")
	// 查询外部联系人关系数据
	req := customer.WorkUserMemberRelation{
		MemberID:           memberData.ID,
		EnterpriseID:       enterpriseData.EnterpriseID,
		WorkExternalUserID: source.ID,
	}
	relation, err := s.CustomerService.FindExternalRelation(ctx, req)
	if err != nil {
		log.Error(ctx).Msgf("查询外部联系人关系失败，错误信息：%+v", err)
		return err
	}
	// 判断在外部联系人在删除成员之前，成员有无删除外部联系人
	if relation.RelationStatus == 0 {
		relation.RelationStatus = customer.WorkUserRelationStatusDeletedByExternal
	} else {
		relation.RelationStatus = customer.WorkUserRelationStatusDeleteEachOther
	}
	relation.ExternalDeletedTime = msgTime
	relation.UpdatedAt = time.Now().Unix()
	rErr := s.CustomerService.UpdateExternalRelation(ctx, *relation)
	if rErr != nil {
		log.Info(ctx).Msgf("更新外部联系人关系失败，错误信息：%+v", rErr)
		return rErr
	}
	// 外部联系人删除成员记录日志
	relation.ID = uid.ID(ctx)
	err = s.CustomerService.AddExternalRelationLog(ctx, []customer.WorkUserMemberRelation{*relation})

	// 查询用户部门信息
	memberDpName := orm.Service.Member.GetMemberDepartment(memberData.Userid, enterpriseData.EnterpriseID)
	var data micro_scrm.DeleteFriendRelationRecordModel
	// 记录成员被外部联系人删除
	data.ID = int64(uint64(uid.ID(ctx)))
	data.ExternalUserID = uint64(source.ID)
	data.DpID = uint64(memberDpName.DpId)
	data.DpName = memberDpName.DpName
	data.UserID = memberData.Userid
	data.EnterpriseID = int64(uint64(enterpriseData.EnterpriseID))
	data.DeleteType = micro_scrm.DeleteFriendDeleteTypeExternalUser
	if err := gormx.GetInstance().Create(&data).Error; err != nil {
		log.Info(ctx).Msgf("记录员工被外部联系人删除，错误信息：%+v", err.Error())
		return err
	}

	eventDesc := "员工：" + memberData.Name + "（" + memberData.Userid + "）"
	_ = s.UserService.AddUserEvent(ctx, enterpriseData.EnterpriseID, source.EuID, string(common.EventEnum.DeleteExternalUser), eventDesc, "", nil, msgTime)

	// 发送删除消息
	// go func() {
	//	defer wg.Done()
	//	message := &customer.SendMessageRequest{
	//		EnterpriseID:        enterpriseData.EnterpriseID,
	//		Type:                customer.SendTypeFromCustomer,
	//		UserName:            source.Name,
	//		MemberName:          memberData.Name,
	//		DeletedMemberUserID: memberData.Userid,
	//	}
	//	s.MessageService.SendMessage(ctx, message, enterpriseData)
	//	select {
	//	case <-ticker.C:
	//		return
	//	}
	// }()

	return nil
}

// dealBatchInsertRelations 处理外部联系人关系
func (s Service) dealBatchInsertRelations(ctx ctx.Context, externalDetail *work.ExternalContactResponse, source *user.SourceUser, callbackData *work.MemberExternalCallback) (*customer.WorkUserMemberRelation, error) {

	var memberUserIDs []string
	var relationsFromQw []*customer.WorkUserMemberRelation
	var requestRelation *customer.WorkUserMemberRelation
	var batchInsertRelation []customer.WorkUserMemberRelation
	var batchInsertRelationLog []customer.WorkUserMemberRelation
	for _, relationV := range externalDetail.FollowUser {
		memberUserIDs = append(memberUserIDs, relationV.UserID)
		tempRelation := &customer.WorkUserMemberRelation{}
		customer.TransferExternalRelation(relationV, tempRelation)
		tempRelation.ExternalUserID = externalDetail.Base.ExternalUserID
		tempRelation.UserID = relationV.UserID
		relationsFromQw = append(relationsFromQw, tempRelation)
	}
	// 从成员表里查询所有的成员
	memberListFromDb := s.SearchChatMemberInner(ctx, memberUserIDs, source.EnterpriseID)
	// 从数据库查询外部联系人关系
	relationReq := customer.BatchWorkUserMemberRelationRequest{
		WorkExternalUserID: source.ID,
		UserIds:            memberUserIDs,
		UserIDCnt:          len(memberUserIDs),
		EnterpriseID:       source.EnterpriseID,
	}
	customerRelations, err := s.CustomerService.BatchExternalRelations(ctx, relationReq)
	if err != nil {
		log.Error(ctx).Msgf("外部联系人关系列表查询失败，错误信息：%+v, 请求参数：%+v", err, relationReq)
		return requestRelation, err
	}
	var customerRelationMap = make(map[string]customer.BatchWorkUserMemberRelationResponse)
	for _, vItem := range customerRelations {
		customerRelationMap[vItem.MemberUserID+source.OpenID] = vItem
	}
	for _, relationV := range relationsFromQw {
		index := relationV.UserID + relationV.ExternalUserID
		// 更新
		if dbRelationItem, ok := customerRelationMap[index]; ok {
			relationV.ID = dbRelationItem.RelationID
			relationV.UpdatedAt = time.Now().Unix()
			relationV.EnterpriseID = source.EnterpriseID
			relationV.UserId = dbRelationItem.MemberUserID
			relationV.UnionId = externalDetail.Base.UnionID
			relationV.MemberDeletedTime = dbRelationItem.MemberDeletedTime
			relationV.ExternalDeletedTime = dbRelationItem.ExternalDeletedTime
			// 编辑回调不更新关系状态
			if relationV.UserID == callbackData.UserID && relationV.ExternalUserID == callbackData.ExternalUserID {
				if callbackData.ChangeType != EditExternalContact {
					relationV.RelationStatus = customer.WorkUserRelationStatusNormal // 回调新增关系状态设置为0
				}
				requestRelation = relationV
			}
			_ = s.CustomerService.UpdateExternalRelation(ctx, *relationV)
		} else {
			if vMember, ok := memberListFromDb[relationV.UserID]; ok {
				relationV.ID = uid.ID(ctx)
				relationV.WorkExternalUserID = source.ID
				relationV.MemberID = vMember.ID
				relationV.CreatedAt = time.Now().Unix()
				relationV.UpdatedAt = time.Now().Unix()
				relationV.EnterpriseID = source.EnterpriseID
				relationV.UserId = callbackData.UserID
				relationV.UnionId = externalDetail.Base.UnionID
				relationV.CreateTime = time.Now().Unix()
				batchInsertRelation = append(batchInsertRelation, *relationV)
				batchInsertRelationLog = append(batchInsertRelationLog, *relationV)
				if relationV.UserID == callbackData.UserID && relationV.ExternalUserID == callbackData.ExternalUserID {
					requestRelation = relationV
				}
			}
		}
	}
	if len(batchInsertRelation) > 0 {
		err = s.CustomerService.AddExternalRelation(ctx, batchInsertRelation)
		if err != nil {
			log.Error(ctx).Msgf("外部联系人关系列表批量入库失败，错误信息：%+v", err)
			return requestRelation, err
		}
		// 外部联系人添加日志
		err = s.CustomerService.AddExternalRelationLog(ctx, batchInsertRelationLog)
	}

	return requestRelation, nil
}

// delExternalContact 成员删除外部联系人
func (s Service) delExternalContact(ctx ctx.Context, enterpriseData *enterprise.Corp, msgTime int64, source *user.SourceUser, memberData *members.Member) error {
	// 查询外部联系人关系数据
	req := customer.WorkUserMemberRelation{
		MemberID:           memberData.ID,
		EnterpriseID:       enterpriseData.EnterpriseID,
		WorkExternalUserID: source.ID,
	}
	relation, err := s.CustomerService.FindExternalRelation(ctx, req)
	if err != nil {
		log.Error(ctx).Msgf("查询外部联系人关系失败，错误信息：%+v", err)
		return err
	}
	// 判断在成员删除外部联系人之前，外部联系人有无删除成员
	if relation.RelationStatus == 0 {
		relation.RelationStatus = customer.WorkUserRelationStatusDeletedByMember
	} else {
		relation.RelationStatus = customer.WorkUserRelationStatusDeleteEachOther
	}
	relation.MemberDeletedTime = msgTime
	relation.UpdatedAt = time.Now().Unix()
	_ = s.CustomerService.UpdateExternalRelation(ctx, *relation)

	// 成员删除外部联系人记录日志
	relation.ID = uid.ID(ctx)
	err = s.CustomerService.AddExternalRelationLog(ctx, []customer.WorkUserMemberRelation{*relation})

	// 查询用户部门信息
	memberDpName := orm.Service.Member.GetMemberDepartment(memberData.Userid, enterpriseData.EnterpriseID)
	var data micro_scrm.DeleteFriendRelationRecordModel
	// 记录成员被外部联系人删除
	data.ID = int64(uint64(uid.ID(ctx)))
	data.ExternalUserID = uint64(source.ID)
	data.DpID = uint64(memberDpName.DpId)
	data.DpName = memberDpName.DpName
	data.UserID = memberData.Userid
	data.EnterpriseID = int64(uint64(enterpriseData.EnterpriseID))
	data.DeleteType = micro_scrm.DeleteFriendDeleteTypeUser
	if err := gormx.GetInstance().Create(&data).Error; err != nil {
		log.Info(ctx).Msgf("记录成员删除外部联系人删除，错误信息：%+v", err.Error())
		return err
	}

	// 发送删除消息
	// go func() {
	//	message := &customer.SendMessageRequest{
	//		EnterpriseID:        enterpriseData.EnterpriseID,
	//		Type:                customer.SendTypeFromMember,
	//		UserName:            source.Name,
	//		MemberName:          memberData.Name,
	//		DeletedMemberUserID: memberData.Userid,
	//	}
	//	s.MessageService.SendMessage(ctx, message, enterpriseData)
	//	return
	// }()

	// 增加员工删除好友事件
	eventDesc := "员工：" + memberData.Name + "（" + memberData.Userid + "）"
	_ = s.UserService.AddUserEvent(ctx, enterpriseData.EnterpriseID, source.EuID, string(common.EventEnum.DeleteExternalUserByMember), eventDesc, "", nil, msgTime)

	return nil
}

// AddExternalUserCustomerTag 添加自定义标签
func (s Service) AddExternalUserCustomerTag(c ctx.Context, relation *customer.WorkUserMemberRelation, tags []*work.ExternalTag) error {
	var tagNameList []string
	for _, v := range tags {
		tagNameList = append(tagNameList, v.TagName)
	}
	// 从数据库中查询已经存在的自定义标签
	dbTags, err := s.ExternalTagService.GetCustomerTags(c, relation.EnterpriseID, relation.MemberID, tagNameList)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	var dbTagMap = make(map[string]*customer.WorkExternalTag)
	for _, v := range dbTags {
		dbTagMap[v.Name] = v
	}
	var newTags []*customer.WorkExternalTag
	for _, v := range tags {
		if _, ok := dbTagMap[v.TagName]; ok {
			continue
		}
		newTags = append(newTags, &customer.WorkExternalTag{
			ID:           uid.ID(c),
			EnterpriseID: relation.EnterpriseID,
			Type:         customer.ExternalTagTypeCustomer,
			Name:         v.TagName,
			MemberID:     relation.MemberID,
			ExpiredTime:  time.Now().Unix(),
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		})
	}
	_, err = s.ExternalTagService.AddExternalTag(c, newTags)
	if err != nil {
		log.Info(c).Msgf("添加自定义标签失败：err:%s, newTags:%+v", err.Error(), newTags)
		return err
	}
	var addTagIDList []int64
	for _, v := range dbTags {
		addTagIDList = append(addTagIDList, v.ID)
	}
	for _, v := range newTags {
		addTagIDList = append(addTagIDList, v.ID)
	}
	_, err = s.ExternalTagService.AddExternalUserTag(c, relation.EnterpriseID, relation.ID, addTagIDList)
	if err != nil {
		log.Info(c).Msgf("添加外部联系人标签关系失败：err:%s, relation:%+v, tagIDList:%+v", err.Error(), relation, addTagIDList)
		return err
	}
	return nil
}

// AddExternalUserCompanyTag 处理用户标签
func (s Service) AddExternalUserCompanyTag(c ctx.Context, relation *customer.WorkUserMemberRelation, tags []*work.ExternalTag, enterpriseData *enterprise.Corp) error {
	// 循环处理每个标签
	var workTagIDList []string
	var customerTagNameList []string
	for _, v := range tags {
		if v.Type == work.ExternalTagTypeOfficial {
			workTagIDList = append(workTagIDList, v.TagID)
		} else {
			customerTagNameList = append(customerTagNameList, v.TagName)
		}
	}
	dbTagList, err := s.ExternalTagService.GetExternalTagOfficial(c, relation.EnterpriseID, workTagIDList)
	if err != nil {
		log.Info(c).Msgf("获取企业标签数据失败，err:%v", err)
		return err
	}
	var dbTagMap = make(map[string]*customer.WorkExternalTag)
	for _, v := range dbTagList {
		dbTagMap[v.TagID] = v
	}
	var newTagIDList []string
	for _, v := range workTagIDList {
		if _, ok := dbTagMap[v]; !ok {
			newTagIDList = append(newTagIDList, v)
		}
	}
	newTagList, err := s.ExternalTagService.GetAndCreateTagListFromWork(c, relation.EnterpriseID, enterpriseData.Domain, newTagIDList)
	newTagList = append(newTagList, dbTagList...)
	var userTagIDs []int64
	for _, v := range newTagList {
		userTagIDs = append(userTagIDs, v.ID)
	}
	_, err = s.ExternalTagService.AddExternalUserTag(c, enterpriseData.EnterpriseID, relation.ID, userTagIDs)
	if err != nil {
		log.Info(c).Msgf("添加用户标签数据失败, err:%s, relationInfo:%+v", err.Error(), relation)
	}
	return nil
}

// DeleteExternalUserTag 删除外部联系人标签
func (s Service) DeleteExternalUserTag(c ctx.Context, userTags []*customer.WorkExternalUserTag, deleteCompanyTagMap map[string]*customer.WorkExternalTag, deleteCustomerTagMap map[string]*customer.WorkExternalTag) error {
	var deleteUserTagIDList []int64
	for _, v := range userTags {
		switch v.Type {
		case customer.ExternalTagTypeCompany:
			if _, ok := deleteCompanyTagMap[v.TagID]; ok {
				deleteUserTagIDList = append(deleteUserTagIDList, v.ID)
			}
		case customer.ExternalTagTypeCustomer:
			if _, ok := deleteCustomerTagMap[v.TagName]; ok {
				deleteUserTagIDList = append(deleteUserTagIDList, v.ID)
			}
		}

	}
	if len(deleteUserTagIDList) == 0 {
		return nil
	}
	_, err := s.ExternalTagService.DeleteExternalUserTags(c, deleteUserTagIDList)
	return err
}

// SendDefaultWelcome  发送默认欢迎语
func (s Service) SendDefaultWelcome(ctx ctx.Context, member *members.Member, callbackData *work.MemberExternalCallback, domain string, workImpl *work.QyWorkImpl, externalDetail *work.ExternalContactResponse) (bool, error) {
	ev := log.Info(ctx).Interface("SendDefaultWelcome-callback", callbackData)
	defer func() {
		ev.Msg("SendDefaultWelcomeOver")
	}()
	var welcomeContentJson work.SendWelcomeRequest
	var attachments = make([]map[string]interface{}, 0)
	var sendAttachments = make([]map[string]interface{}, 0)
	// 查询配置的默认欢迎语
	welcomeDefault, err := s.DrainageService.WelcomeDefault(ctx, drainage.WelcomeDetailRequest{EnterpriseId: member.EnterpriseID})
	if err != nil {
		ev = ev.Err(err)
		return false, nil
	}

	if err = json.Unmarshal([]byte(welcomeDefault.Content), &attachments); err != nil {
		ev = ev.Err(err)
		return false, nil
	}

	welcomeContentJson.Attachments = sendAttachments
	if len(attachments) == 0 {
		return false, nil
	}
	senWelcomeRequest, err := s.genSendAttachments(ctx, attachments, member.EnterpriseID, domain)
	//sw小程序路径添加参数
	resConfig, _ := cacheSvc.NewRedisData().GetEnterpriseConfig(ctx, member.EnterpriseID, "store_name")
	if resConfig.Name == "sw" {
		var pathMapArr = make(map[int64]string) //欢迎语内容中的原始小程序path
		for i, attachmentsArr := range senWelcomeRequest.Attachments {
			for key, val := range attachmentsArr {
				if key == "miniprogram" {
					tmpVal := val.(work.MiniProgram)
					pathMapArr[int64(i)] = tmpVal.Page //将多个小程序path收集,下面批量处理
				}
			}
		}
		if len(pathMapArr) > 0 {
			//批量处理小程序path得到新的path
			newPathMap, _ := membersSvc.NewMembersService().SetMiniProgramLinkParams(ctx, member.EnterpriseID, member.ID, pathMapArr)
			for i, attachmentsArr := range senWelcomeRequest.Attachments {
				for key, val := range attachmentsArr {
					if key == "miniprogram" {
						tmpVal := val.(work.MiniProgram)
						if _, ok := newPathMap[int64(i)]; ok {
							tmpVal.Page = newPathMap[int64(i)] //反写回欢迎语内容中
							attachmentsArr[key] = tmpVal
						}
					}
				}
			}
		}
	}

	return s.FormatDataSendWelcome(ctx,
		member,
		callbackData,
		workImpl,
		externalDetail,
		senWelcomeRequest,
		microwork.UserInfoResponse{}), nil
}

// 生成欢迎语组装信息
func (s Service) genSendAttachments(ctx ctx.Context, attachments []map[string]interface{}, enterpriseId int64, domain string) (*work.SendWelcomeRequest, error) {
	var (
		err                error
		welcomeContentJson work.SendWelcomeRequest
		sendAttachments    []map[string]interface{}
	)

	for index, vMap := range attachments {
		for k, v := range vMap {
			vj, _ := json.Marshal(v)
			if k == "text" {
				var vdata work.Text
				err = json.Unmarshal(vj, &vdata)
				if err != nil {
					log.Error(ctx).Msgf("send-welcome-attachments-%s %+v", k, err)
					return nil, err
				}
				welcomeContentJson.Text = vdata
				attachments[index][k] = vdata
			}

			if k == "link" {
				attachments[index]["msgtype"] = k
				var vdata work.Link
				err = json.Unmarshal(vj, &vdata)
				if err != nil {
					log.Error(ctx).Msgf("send-welcome-attachments-%s %+v", k, err)
					return nil, err
				}
				vdata.WxPicUrl = vdata.PicUrl
				attachments[index][k] = vdata
				sendAttachments = append(sendAttachments, attachments[index])
			}
			if k == "image" {
				attachments[index]["msgtype"] = k
				var vdata work.Image
				err = json.Unmarshal(vj, &vdata)
				if err != nil {
					log.Error(ctx).Msgf("send-welcome-attachments-%s %+v", k, err)
					return nil, err
				}
				if vdata.PicUrl != "" {
					cacheMediaRequest := drainage.CacheMediaRequest{
						EnterpriseId: enterpriseId,
						ImgUrl:       vdata.PicUrl,
						FileType:     "1",
					}
					cacheMediaResponse, err := s.DrainageService.CacheMedia(ctx, cacheMediaRequest, domain)
					if err != nil || cacheMediaResponse.Result == false {
						return nil, errs.WeChatImgUrlCacheError
					}
					vdata.MediaId = cacheMediaResponse.MediaId

				}
				attachments[index][k] = vdata
				sendAttachments = append(sendAttachments, attachments[index])
			}
			if k == "miniprogram" {
				attachments[index]["msgtype"] = k
				var vdata work.MiniProgram
				err = json.Unmarshal(vj, &vdata)
				if err != nil {
					log.Error(ctx).Msgf("send-welcome-attachments-%s %+v", k, err)
					return nil, err
				}
				if vdata.CoverUrl != "" {
					ImgMediaId := redis.Client().Get(ctx, vdata.CoverUrl).Val()
					if ImgMediaId != "" {
						vdata.PicMediaId = ImgMediaId
					} else {
						cacheMediaRequest := drainage.CacheMediaRequest{
							EnterpriseId: enterpriseId,
							ImgUrl:       vdata.CoverUrl,
							FileType:     "1",
						}
						cacheMediaResponse, err := s.DrainageService.CacheMedia(ctx, cacheMediaRequest, domain)
						if err != nil || cacheMediaResponse.Result == false {
							return nil, errs.WeChatImgUrlCacheError
						}
						vdata.PicMediaId = cacheMediaResponse.MediaId
					}
				}

				attachments[index][k] = vdata
				if vdata.Appid != "" {
					sendAttachments = append(sendAttachments, attachments[index])
				}
			}

			if k == "video" {
				attachments[index]["msgtype"] = k
				var vdata work.Video
				err = json.Unmarshal(vj, &vdata)
				if err != nil {
					log.Error(ctx).Msgf("send-welcome-attachments-%s %+v", k, err)
					return nil, err
				}
				if vdata.CoverUrl != "" {
					ImgMediaId := redis.Client().Get(ctx, vdata.CoverUrl).Val()
					if ImgMediaId != "" {
						vdata.MediaId = ImgMediaId
					} else {
						cacheMediaRequest := drainage.CacheMediaRequest{
							EnterpriseId: enterpriseId,
							ImgUrl:       vdata.CoverUrl,
							FileType:     "3",
						}
						cacheMediaResponse, err := s.DrainageService.CacheMedia(ctx, cacheMediaRequest, domain)
						if err != nil || cacheMediaResponse.Result == false {
							return nil, errs.WeChatImgUrlCacheError
						}
						vdata.MediaId = cacheMediaResponse.MediaId
					}
				}
				attachments[index][k] = vdata
				sendAttachments = append(sendAttachments, attachments[index])
			}

			if k == "file" {
				attachments[index]["msgtype"] = k
				var vdata work.File
				err = json.Unmarshal(vj, &vdata)
				if err != nil {
					log.Error(ctx).Msgf("send-welcome-attachments-%s %+v", k, err)
					return nil, err
				}
				if vdata.CoverUrl != "" {
					ImgMediaId := redis.Client().Get(ctx, vdata.CoverUrl).Val()
					if ImgMediaId != "" {
						vdata.MediaId = ImgMediaId
					} else {
						cacheMediaRequest := drainage.CacheMediaRequest{
							EnterpriseId: enterpriseId,
							ImgUrl:       vdata.CoverUrl,
							FileType:     "99",
						}
						cacheMediaResponse, err := s.DrainageService.CacheMedia(ctx, cacheMediaRequest, domain)
						if err != nil || cacheMediaResponse.Result == false {
							return nil, errs.WeChatImgUrlCacheError
						}
						vdata.MediaId = cacheMediaResponse.MediaId
					}
				}
				attachments[index][k] = vdata
				sendAttachments = append(sendAttachments, attachments[index])
			}
		}
	}

	welcomeContentJson.Attachments = sendAttachments
	return &welcomeContentJson, nil
}

// GetFissionUserActivity 根据用户裂变欢迎语信息查询
func (s Service) GetFissionUserActivity(c ctx.Context, contactID int64, enterpriseId int64) ([]map[string]interface{}, error) {
	// 通过查询用户发起的裂变活动ID
	var ContentJson = make([]map[string]interface{}, 0)
	fissionUserInfo, err := s.FissionService.GetActivityUserByContactID(c, contactID)
	log.Info(c).Msgf("====进入发送欢迎语模块-查询用户发起的裂变活动：%+v", fissionUserInfo)
	if err != nil {
		return ContentJson, errs.FissionActivityErr
	}

	activityInfo, fissionUserInfoErr := s.FissionService.FissionActivityDetail(c, fissionUserInfo.ActivityID, enterpriseId)
	log.Info(c).Msgf("====进入发送欢迎语模块-裂变活动详情：%+v", activityInfo)
	if fissionUserInfoErr != nil {
		return ContentJson, errs.ActivityInfoErr
	}

	// 获取发送的欢迎语
	wf, wfErr := welcome.NewWelcomeFactory(c, welcome.Param{
		EnterpriseID: enterpriseId,
		ActivityID:   fissionUserInfo.ActivityID,
	})
	if wfErr != nil {
		return ContentJson, wfErr
	}

	return wf.Get()
}

func (s Service) GetFissionActivity(c ctx.Context, contactID int64, enterpriseId int64) ([]map[string]interface{}, error) {
	var ContentJson = make([]map[string]interface{}, 0)
	log.Info(c).Msgf("====进入发送欢迎语模块-获取活动信息请求参数：%+v，%+v", contactID, enterpriseId)
	activityInfo, err := s.FissionService.ByContactIdActivityDetail(c, contactID, enterpriseId)
	log.Info(c).Msgf("====进入发送欢迎语模块-获取活动信息：%+v", activityInfo)
	if err != nil {
		log.Info(c).Msgf("====进入发送欢迎语模块-获取活动信息出错：%+v", err)
		return ContentJson, errs.ActivityInfoErr
	}
	if activityInfo == nil {
		return ContentJson, errs.ActivityInfoErr
	}

	// 获取发送的欢迎语
	wf, wfErr := welcome.NewWelcomeFactory(c, welcome.Param{
		EnterpriseID: enterpriseId,
		ActivityID:   activityInfo.ID,
	})
	if wfErr != nil {
		return ContentJson, wfErr
	}

	return wf.Get()
}

// UserTagHandler 用户标签处理
func (s Service) UserTagHandler(c ctx.Context, callbackData *work.MemberExternalCallback, workImpl *work.QyWorkImpl, relationID, enterpriseId int64) error {
	log.Info(c).Msg("用户标签处理")
	if relationID == 0 {
		log.Info(c).Msg("关系relationID没获取到")
		return nil
	}
	req := work.UserTagHandlerRequest{
		Userid:         callbackData.UserID,
		ExternalUserid: callbackData.ExternalUserID,
	}
	//联系我二维码渠道tag
	drainageContactsAddTag, err := s.drainageContactsTagHandler(c, callbackData, relationID)
	if err != nil {
		return err
	}
	//活动二维码tag
	activityQrCodeTags, _ := s.activityQrCodeTagHandler(c, callbackData, enterpriseId, relationID)
	//合并要打的标签
	req.AddTag = helper.ArrayMerge(drainageContactsAddTag, activityQrCodeTags)
	if len(req.AddTag) == 0 {
		log.Info(c).Msg("无需打企业标签")
		return nil
	}
	isOK, err := workImpl.UserTagHandler(c, req)
	if err != nil || isOK == false {
		return err
	}
	return nil
}

func (s Service) drainageContactsTagHandler(c ctx.Context, callbackData *work.MemberExternalCallback, relationID int64) ([]string, error) {
	// 渠道ID
	State, _ := strconv.ParseInt(callbackData.State, 10, 64)

	// 查询渠道id
	contactReq := drainage.ContactDetailsRequest{
		ID: State,
	}
	detail, err := s.DrainageService.ContactSelect(c, contactReq)
	if err != nil {
		log.Error(c).Interface("req", contactReq).Err(err).Send()
		return []string{}, nil
	}
	cReq := drainage.ContactChannelLabelsListRequest{
		ChannelID: detail.ChannelId,
	}
	tags, err := s.DrainageService.TagList(c, &cReq)
	if err != nil {
		log.Error(c).Interface("req", cReq).Err(err).Send()
		return []string{}, nil
	}

	if len(tags) > 0 {
		var dataList []work_external_user.WorkExternalUserTagModel
		for _, tag := range tags {
			dataList = append(dataList, work_external_user.WorkExternalUserTagModel{
				MicroScrmBaseModel: model.MicroScrmBaseModel{
					BaseModel: model.BaseModel{
						ID:        uid.ID(c),
						CreatedAt: time.Now().Unix(),
						UpdatedAt: time.Now().Unix(),
					},
					EnterpriseID: detail.EnterpriseId,
				},
				RelationID:    relationID,
				ExternalTagID: tag.LabelID,
			})
		}
		tagErr := orm.Service.WorkExternalUser.InsertUserTag(dataList)
		if tagErr != nil {
			log.Error(c).Msgf("插入标签异常: %v", tagErr)
			return []string{}, nil
		}
		events.ScrmEvent.ChangeUserTag.Dispatch(entity.Data{
			Ctx:          c,
			EnterpriseID: detail.EnterpriseId,
			RelationID:   relationID,
		})
		var resTag []string
		// 查询所属渠道TagID列表
		for _, v := range tags {
			if v.WxTagID == "" {
				continue
			}
			resTag = append(resTag, v.WxTagID)
		}
		return resTag, nil
	}
	return []string{}, nil
}

// ActivityQrCodeTagHandler  活动用户标签处理
func (s Service) activityQrCodeTagHandler(c ctx.Context, callbackData *work.MemberExternalCallback, enterpriseId int64, relationId int64) ([]string, error) {
	log.Info(c).Msg("用户活动标签处理")
	// 渠道ID
	State, _ := strconv.ParseInt(callbackData.State, 10, 64)
	// 查询互动点位id
	pointInfo, err := activityQRCodeDao.NewPointDao().GetOne(c, enterpriseId, "id = ?", State)
	if pointInfo.ID == 0 {
		return []string{}, nil
	}
	activityInfo, err := activityQRCodeDao.NewMainDao().GetOne(c, enterpriseId, "id = ?", pointInfo.ActivityID)
	if activityInfo.ID == 0 {
		return []string{}, nil
	}
	var tagIds []string
	if activityInfo.TagIds != "" {
		tagIds = strings.Split(activityInfo.TagIds, ",")
	}
	if activityInfo.AddCategoryTag != 0 {
		categoryInfo, _ := activityQRCodeDao.NewCategoryDao().GetOne(c, enterpriseId, "id = ?", pointInfo.CategoryID)
		if categoryInfo.TagId > 0 {
			tagIds = append(tagIds, helper.Int64ToString(categoryInfo.TagId))
		}
	}
	if activityInfo.AddNameTag != 0 {
		tagIds = append(tagIds, helper.Int64ToString(activityInfo.NameTagId))
	}
	if err != nil {
		return []string{}, nil
	}
	if len(tagIds) == 0 {
		return []string{}, nil
	}
	// 是否发送
	var dataList []work_external_user.WorkExternalUserTagModel
	for _, tagId := range tagIds {
		dataList = append(dataList, work_external_user.WorkExternalUserTagModel{
			MicroScrmBaseModel: model.MicroScrmBaseModel{
				BaseModel: model.BaseModel{
					ID:        uid.ID(c),
					CreatedAt: time.Now().Unix(),
					UpdatedAt: time.Now().Unix(),
				},
				EnterpriseID: enterpriseId,
			},
			RelationID:    relationId,
			ExternalTagID: int64(helper.StringToInt(tagId)),
		})
	}
	tagErr := orm.Service.WorkExternalUser.InsertUserTag(dataList)
	if tagErr != nil {
		return []string{}, nil
	}

	tagList, _ := tagsDao.NewWorkExternalTag().GetMany(c, "enterprise_id=? AND id in ?", enterpriseId, tagIds)
	// 查询所属渠道TagID列表
	var qyTagIds []string
	for _, v := range tagList {
		if v.TagID == "" {
			continue
		}
		qyTagIds = append(qyTagIds, v.TagID)
	}
	return qyTagIds, nil
}

// AddExternalTraceHandle 外部联系人加入系统
func (s Service) AddExternalTraceHandle(c ctx.Context, callbackData *work.MemberExternalCallback, enterpriseId int64) {
	if callbackData.State == "" {
		return
	}

	isExists, err := s.ExternalService.ExternalExists(c, callbackData.ExternalUserID, enterpriseId)
	if err != nil {
		log.Error(c).Msgf("查询外部联系人存在出错,%v", err)
		return
	}
	// 存在则直接返回
	if isExists {
		return
	}

	err = s.ExternalService.JoinExternal(c, callbackData, enterpriseId)
	if err != nil {
		log.Error(c).Msgf("查询外部联系人存在出错,%v", err)
		return
	}
	return
}

// AddExternalChannelTraceHandle 添加渠道外部联系人
func (s Service) AddExternalChannelTraceHandle(c ctx.Context, callbackData *work.MemberExternalCallback, enterpriseId, channelId int64) {

	// 查询该外部联系人是否已经添加过该渠道
	isExists, err := s.ExternalService.ExternalExistsInChannel(c, callbackData.ExternalUserID, enterpriseId, channelId)
	if err != nil {
		log.Error(c).Msgf("查询渠道外部联系人存在出错,%v", err)
		return
	}
	// 存在则直接返回
	if isExists {
		return
	}

	err = s.ExternalService.JoinExternalInChannel(c, callbackData, enterpriseId, channelId)
	if err != nil {
		log.Error(c).Msgf("查询渠道外部联系人存在出错,%v", err)
		return
	}

	return
}

// DelExternalTraceHandle 删除外部联系人关系
func (s Service) DelExternalTraceHandle(c ctx.Context, workExternalUserID int64, data *work.MemberExternalCallback, enterpriseId int64) {

	// 查询外部联系人是否存在正常的关系
	isExists, err := s.CustomerService.HasNormalRelation(c, workExternalUserID, enterpriseId, 0)
	if err != nil {
		log.Error(c).Msgf("查询外部联系人出错,%v", err)
		return
	}

	if isExists {
		return
	}

	// 删除该外部联系人
	err = s.ExternalService.LeaveExternal(c, data, enterpriseId)
	if err != nil {
		log.Error(c).Msgf("删除外部联系人出错,%v", err)
		return
	}
	return
}

// DelExternalChannelTraceHandle 删除渠道外部联系人关系
func (s Service) DelExternalChannelTraceHandle(c ctx.Context, workExternalUserID int64, data *work.MemberExternalCallback, enterpriseId int64, channelId int64) {

	// 查询渠道外部联系人是否存在正常的关系
	isExists, err := s.CustomerService.HasNormalRelation(c, workExternalUserID, enterpriseId, channelId)
	if err != nil {
		log.Error(c).Msgf("查询渠道外部联系人出错,%v", err)
		return
	}

	if isExists {
		return
	}

	// 删除该外部联系人
	err = s.ExternalService.LeaveExternalInChannel(c, data, enterpriseId, channelId)
	if err != nil {
		log.Error(c).Msgf("删除渠道外部联系人出错,%v", err)
		return
	}
	return
}

// DelRelationChannelTrace 外部联系人删除员工
func (s Service) DelRelationChannelTrace(c ctx.Context, data *work.MemberExternalCallback, enterpriseId, channelId int64) {

	// 插入好友删除员工
	err := s.ExternalService.AddDelExternalContact(c, data, enterpriseId, channelId)
	if err != nil {
		log.Error(c).Msgf("外部联系人删除员工出错,%v", err)
	}

	return
}

// RecoveryDelRelationChannel 恢复外部联系人删除员工
func (s Service) RecoveryDelRelationChannel(c ctx.Context, data *work.MemberExternalCallback, enterpriseId, channelId int64) {
	if channelId == 0 {
		return
	}

	existDelRelation, err := s.ExternalService.ExistDelRelation(c, data, enterpriseId)
	if err != nil {
		log.Error(c).Msgf("查询恢复外部联系人删除员工出错,%v", err)
	}
	if existDelRelation == nil {
		return
	}

	_ = s.ExternalService.RecoveryDelRelation(c, data, existDelRelation, channelId)
	if err != nil {
		log.Error(c).Msgf("恢复外部联系人删除员工出错,%v", err)
	}

	return
}

// RelationChannelId 查询该关系的channelID
func (s Service) RelationChannelId(c ctx.Context, data *work.MemberExternalCallback, enterpriseId int64) (int64, error) {
	switch data.ChangeType {
	case AddExternalContact:
		if data.State != "" {
			contactsId, _ := strconv.ParseInt(data.State, 10, 64)
			contactDetailsRequest := drainage.ContactDetailsRequest{
				EnterpriseId: enterpriseId,
				ID:           contactsId,
			}
			drainageContacts, _ := s.ContactService.GetContact(c, contactDetailsRequest)
			if drainageContacts == nil {
				return 0, nil
			}

			return drainageContacts.ChannelId, nil
		}
		return 0, nil
	case DelFollowUser, DelExternalContact:
		// 查询上次添加的时候的渠道值
		req := &work.MemberExternalCallback{
			ToUserName:     data.ToUserName,
			ChangeType:     AddExternalContact,
			UserID:         data.UserID,
			ExternalUserID: data.ExternalUserID,
		}
		detail, err := s.CustomerService.CallbackRecordDetail(c, req)
		if err != nil || detail == nil {
			return 0, err
		}

		return detail.ChannelID, nil
	default:
		return 0, errors.New(fmt.Sprintf("不支持类型%s,统计渠道值", data.ChangeType))
	}
}

func (s Service) IsRunFission(c ctx.Context, data *work.MemberExternalCallback, enterpriseId int64) (bool, error) {
	switch data.ChangeType {
	case DelFollowUser, DelExternalContact:
		// 查询上次添加的时候记录
		req := &work.MemberExternalCallback{
			ToUserName:     data.ToUserName,
			ChangeType:     AddExternalContact,
			UserID:         data.UserID,
			ExternalUserID: data.ExternalUserID,
		}
		detail, err := s.CustomerService.CallbackRecordDetail(c, req)
		if err != nil && err != sql.ErrNoRows {
			return false, err
		}
		if detail == nil {
			return false, err
		}
		data = detail
		fallthrough
	case AddExternalContact:
		if data.State != "" {
			contactsId, _ := strconv.ParseInt(data.State, 10, 64)
			contactDetailsRequest := drainage.ContactDetailsRequest{
				EnterpriseId: enterpriseId,
				ID:           contactsId,
			}
			drainageContacts, contactErr := s.ContactService.ContactDetails(c, contactDetailsRequest)
			if contactErr != nil && contactErr != sql.ErrNoRows {
				return false, contactErr
			}
			if drainageContacts == nil {
				return false, contactErr
			}
			return drainageContacts.Type == ActivityWelcomeMessage, nil
		}
		return false, nil
	default:
		return false, nil
	}
}

// AddRelationRecord 添加记录
func (s Service) AddRelationRecord(c ctx.Context, data *work.MemberExternalCallback) (*workexternal.WorkRelationAddRecord, error) {
	// 查詢渠道ID
	addChannelID, err := s.RelationChannelId(c, data, 0)
	if err != nil {
		return nil, err
	}
	return s.ExternalService.AddRelationRecord(c, data, addChannelID)
}

// FindLastAddRelationRecord 查找上一条记录
func (s Service) FindLastAddRelationRecord(c ctx.Context, data *work.MemberExternalCallback) (*workexternal.WorkRelationAddRecord, error) {
	return s.ExternalService.FindLastAddRelationRecord(c, data)
}

// FinishAddRelation 完成状态
func (s Service) FinishAddRelation(c ctx.Context, data *workexternal.WorkRelationAddRecord) error {
	return s.ExternalService.FinishAddRelation(c, data)
}

// 渠道添加成员
func (s Service) handleChannelAddUser(ctx ctx.Context, workExternalUserId, channelId, enterpriseID int64) {
	if channelId > 0 {
		// 累加渠道添加成员数量
		relation := customer.ExternalUserRelation{
			WorkExternalUserId: workExternalUserId,
			State:              channelId,
			EnterpriseId:       enterpriseID,
		}
		// 查询用户渠道关联关系是否存在
		relationResult, relationErr := s.DrainageService.GetWorkExternalUserRelation(ctx, relation)
		log.Info(ctx).Msgf("新增联系人验证回调,外部联系人关系，查询条件:%v 查询结果:%v", relation, relationResult)
		if relationErr != nil {
			log.Info(ctx).Msgf("新增联系人验证回调,查询外部联系人关系,失败:%v", relationErr)
		}
		// 不存渠道关联关系,则更新联系我新增好友数+1
		if relationResult.ID == 0 {
			err := s.DrainageService.SumAddMemberNumber(ctx, channelId)
			log.Info(ctx).Msgf("新增联系人验证回调,更新联系我新增好友数,失败:%v", err)
		}
	}
	return
}

// HandleActivityQrWelcome  活动欢迎语和投放渠道欢迎语
func (s Service) HandleActivityQrWelcome(
	c ctx.Context,
	member *members.Member,
	callbackData *work.MemberExternalCallback,
	enterpriseData *enterprise.Corp,
	workImpl *work.QyWorkImpl,
	externalDetail *work.ExternalContactResponse,
) bool {
	log.Info(c).Msg("活动、投放渠道欢迎语")
	// 渠道ID
	state, _ := strconv.ParseInt(callbackData.State, 10, 64)
	// 查询互动点位id
	pointInfo, _ := activityQRCodeDao.NewPointDao().GetOne(c, enterpriseData.EnterpriseID, "id = ?", state)
	if pointInfo.ID == 0 {
		log.Info(c).Msg("没有找到活动")
		return false
	}
	welcomeList, _ := activityQRCodeDao.NewWelcomeDao().GetMany(c, enterpriseData.EnterpriseID, "rel_type = ? AND rel_id = ?", activityQRCodeModel.WelcomeTypeActivity, pointInfo.ActivityID)
	if len(welcomeList) == 0 {
		welcomeList, _ = activityQRCodeDao.NewWelcomeDao().GetMany(c, enterpriseData.EnterpriseID, "rel_type = ? AND rel_id = ?", activityQRCodeModel.WelcomeTypeCategory, pointInfo.CategoryID)
		if len(welcomeList) == 0 {
			log.Info(c).Msg("没查到活动、投放渠道欢迎语")
			return false
		}
	}
	var crmUserInfo microwork.UserInfoResponse
	var attachments = make([]map[string]interface{}, 0)
	if len(welcomeList) == 1 && welcomeList[0].Type == 0 {
		_ = json.Unmarshal([]byte(welcomeList[0].Content), &attachments)
	} else {
		var hasCardLevel bool
		// 查询work请求crm获取用户信息
		crmUserInfo, _ = microwork.NewCrm().GetUserInfoByUnionId(c, microwork.UserInfoRequest{
			EnterpriseID: member.EnterpriseID,
			UnionId:      externalDetail.Base.UnionID,
		})
		// 判断是否查询到用户信息(根据status返回判断) true已查到 false未查到
		if crmUserInfo.Data.Status {
			// 判断是否为会员
			if crmUserInfo.Data.CustomerLevel != "" {
				hasCardLevel = true
			}
		}

		for _, welcomeModel := range welcomeList {
			if hasCardLevel && welcomeModel.Type == activityQRCodeModel.WelcomeTypeHasCardLevel {
				_ = json.Unmarshal([]byte(welcomeModel.Content), &attachments)
			} else if !hasCardLevel && welcomeModel.Type == activityQRCodeModel.WelcomeTypeNoCardLevel {
				_ = json.Unmarshal([]byte(welcomeModel.Content), &attachments)
			}
		}
	}
	if len(attachments) == 0 {
		log.Info(c).Msg("没找到符合条件的活动、投放渠道欢迎语")
		return false
	}
	welcomeContentJson, err := s.genSendAttachments(c, attachments, enterpriseData.EnterpriseID, enterpriseData.Domain)
	if err != nil {
		log.Info(c).Msg("生成欢迎语组装信息异常")
		return false
	}
	return s.FormatDataSendWelcome(c,
		member,
		callbackData,
		workImpl,
		externalDetail,
		welcomeContentJson,
		crmUserInfo)
}

// HandleDrainageWelcomesWelcome  联系我渠道欢迎语
func (s Service) HandleDrainageWelcomesWelcome(
	c ctx.Context,
	member *members.Member,
	callbackData *work.MemberExternalCallback,
	enterpriseData *enterprise.Corp,
	workImpl *work.QyWorkImpl,
	externalDetail *work.ExternalContactResponse,
) bool {
	log.Info(c).Msg("联系我渠道欢迎语")
	var attachments = make([]map[string]interface{}, 0)
	// 渠道ID
	state, _ := strconv.ParseInt(callbackData.State, 10, 64)
	// 查询联系我二维码
	drainageContactsInfo, _ := drainageContactDao.NewDrainageContactsDao().GetOne(c, enterpriseData.EnterpriseID, "id = ?", state)
	if drainageContactsInfo.ID == 0 {
		log.Info(c).Msg("没有找到联系我渠道")
		return false
	}
	if drainageContactsInfo.ContactGroupId > 0 {
		//联系我活动
		contactGroupInfo, _ := drainageContactDao.NewDrainageContactGroupDao().GetOnePreloadWelcomes(c, enterpriseData.EnterpriseID, drainageContactsInfo.ContactGroupId)
		if contactGroupInfo.Welcomes != nil && len(contactGroupInfo.Welcomes) > 0 {
			//联系我活动配置了欢迎语
			return s.handleDrainageContactsGroupWelcome(c, contactGroupInfo.Welcomes, member, callbackData, workImpl, externalDetail, enterpriseData)
		}
	}
	// 查询到了渠道，查看渠道是否配置了欢迎语
	welcomeList, err := s.DrainageService.WelcomeList(c, drainage.WelcomeListRequest{
		Page:         1,
		Size:         1,
		EnterpriseId: member.EnterpriseID,
		ChannelId:    drainageContactsInfo.ChannelID,
	})
	if err != nil {
		return false
	}
	if welcomeList.Total == 0 {
		log.Info(c).Msg("没有找到联系我渠道相关欢迎语")
		return false
	}
	// 有联系我渠道相关欢迎语
	attachments = welcomeList.List[0].Content
	welcomeContentJson, err := s.genSendAttachments(c, attachments, enterpriseData.EnterpriseID, enterpriseData.Domain)
	if err != nil {
		log.Info(c).Msg("生成欢迎语组装信息异常")
		return false
	}
	return s.FormatDataSendWelcome(c,
		member,
		callbackData,
		workImpl,
		externalDetail,
		welcomeContentJson,
		microwork.UserInfoResponse{})
}

// HandleDrainageContactsGroupWelcome   联系我活动发送欢迎语
func (s Service) handleDrainageContactsGroupWelcome(
	c ctx.Context,
	welcomes []drainageContactModel.DrainageContactsWelcomesModel,
	member *members.Member,
	callbackData *work.MemberExternalCallback,
	workImpl *work.QyWorkImpl,
	externalDetail *work.ExternalContactResponse,
	enterpriseData *enterprise.Corp,
) bool {
	//查询是否有联系我活动的欢迎语
	var crmUserInfo microwork.UserInfoResponse
	var attachments = make([]map[string]interface{}, 0)
	if len(welcomes) == 1 && welcomes[0].Type == 0 {
		_ = json.Unmarshal([]byte(welcomes[0].Content), &attachments)
	} else {
		var hasCardLevel bool
		// 查询work请求crm获取用户信息
		crmUserInfo, _ = microwork.NewCrm().GetUserInfoByUnionId(c, microwork.UserInfoRequest{
			EnterpriseID: member.EnterpriseID,
			UnionId:      externalDetail.Base.UnionID,
		})
		// 判断是否查询到用户信息(根据status返回判断) true已查到 false未查到
		if crmUserInfo.Data.Status {
			// 判断是否为会员
			if crmUserInfo.Data.CustomerLevel != "" {
				hasCardLevel = true
			}
		}

		for _, welcomeModel := range welcomes {
			if hasCardLevel && welcomeModel.Type == activityQRCodeModel.WelcomeTypeHasCardLevel {
				_ = json.Unmarshal([]byte(welcomeModel.Content), &attachments)
			} else if !hasCardLevel && welcomeModel.Type == activityQRCodeModel.WelcomeTypeNoCardLevel {
				_ = json.Unmarshal([]byte(welcomeModel.Content), &attachments)
			}
		}
	}
	if len(attachments) == 0 {
		log.Info(c).Msg("没找到符合条件的联系我活动二维码")
		return false
	}
	welcomeContentJson, err := s.genSendAttachments(c, attachments, enterpriseData.EnterpriseID, enterpriseData.Domain)
	if err != nil {
		log.Info(c).Msg("生成欢迎语组装信息异常")
		return false
	}
	return s.FormatDataSendWelcome(c,
		member,
		callbackData,
		workImpl,
		externalDetail,
		welcomeContentJson,
		crmUserInfo,
	)
}

// FormatDataSendWelcome  整理欢迎语数据发送欢迎语
func (s Service) FormatDataSendWelcome(
	c ctx.Context,
	member *members.Member,
	callbackData *work.MemberExternalCallback,
	workImpl *work.QyWorkImpl,
	externalDetail *work.ExternalContactResponse,
	welcomeContentJson *work.SendWelcomeRequest,
	crmUserInfo microwork.UserInfoResponse,
) bool {
	log.Info(c).Msg("欢迎语发送操作开始")
	welcomeContent := welcomeContentJson.Text.Content
	welcomeContent = strings.Replace(welcomeContent, "{顾客昵称}", externalDetail.Base.Name, -1)
	welcomeContent = strings.Replace(welcomeContent, "{导购姓名}", member.Name, -1)
	if strings.Contains(welcomeContent, "{会员姓名}") {
		if crmUserInfo.Data.Id == "" {
			// 查询work请求crm获取用户信息
			crmUserInfo, _ = microwork.NewCrm().GetUserInfoByUnionId(c, microwork.UserInfoRequest{
				EnterpriseID: member.EnterpriseID,
				UnionId:      externalDetail.Base.UnionID,
			})
		}

		// 默认会员姓名为昵称
		customerRealName := externalDetail.Base.Name
		if crmUserInfo.Data.Name != "" {
			customerRealName = crmUserInfo.Data.Name
		}
		welcomeContent = strings.Replace(welcomeContent, "{会员姓名}", customerRealName, -1)
	}
	if strings.Contains(welcomeContent, "{门店名称}") {
		// 查询员工的部门
		memberInfo, _ := membersDao.NewMembers().GetMemberInfoAndDepartmentsByUserId(c, member.EnterpriseID, []string{member.Userid})
		// 查询员工部门信息
		if len(memberInfo) > 0 {
			storeDepartmentInfo, _ := microWorkDao.NewStoreDepartment().GetInfo(c, microWorkModel.StoreDepartmentModel{
				MicroBaseModel:     models.MicroBaseModel{EnterpriseID: member.EnterpriseID},
				WechatDepartmentID: uint(memberInfo[0].DepartmentThirdId),
			})
			welcomeContent = strings.Replace(welcomeContent, "{门店名称}", storeDepartmentInfo.StoreName, -1)
		}
	}
	welcomeContentJson.Text.Content = welcomeContent

	welcomeContentJson.WelcomeCode = callbackData.WelcomeCode

	_ = workImpl.SendWelcome(c, *welcomeContentJson)
	log.Info(c).Msg("发送欢迎语操作结束")
	return true
}

// HandleThirdWelcome 处理三方欢迎语
func (s Service) HandleThirdWelcome(
	ctx ctx.Context,
	member *members.Member,
	callbackData *work.MemberExternalCallback,
	domain string,
	workImpl *work.QyWorkImpl,
	externalDetail *work.ExternalContactResponse,
) bool {
	ev := log.Info(ctx).Interface("处理三方欢迎语开始", callbackData)
	defer func() {
		ev.Msg("处理三方欢迎语结束")
	}()
	// 外部联系人是不就不存在
	if externalDetail.Base.ExternalUserID == "" {
		ev.AnErr("外部联系人是不就不存在", errors.New("外部联系人是不就不存在"))
		return false
	}

	enterpriseId := member.EnterpriseID

	ev = ev.Interface("开始callbackData数据", callbackData)

	// 查询企业配置的所有个性欢迎语
	allWelcome, total, err := orm.Service.Welcome.GetAllWelcome(enterpriseId, 9999, 1, "")
	if err != nil || total == 0 {
		ev = ev.AnErr("查询个性欢迎语", err)
		return false
	}

	// 获取企业配置的规则条件数据
	welcomesExternalTypeData, _ := orm.Service.Welcome.GetAllType(enterpriseId)
	if len(welcomesExternalTypeData) == 0 {
		ev.AnErr("获取企业配置的规则条件数据", errors.New("企业未配置规则"))
		return false
	}
	var welcomesExternalTypeCode []string
	for _, item := range welcomesExternalTypeData {
		welcomesExternalTypeCode = append(welcomesExternalTypeCode, item.Code)
	}

	// 查询三方预约记录
	var reserveRecordInfo thirdModel.ReserveRecordModel
	// 预约状态开关
	var reserveStatusSwitch bool
	contactsId := helper.StringToInt(callbackData.State)
	reserveStatus := drainage2.ReserveStatusNo
	// 判断是否存在预约规则
	if helper.InArray[string](welcome2.ReserveStatus, welcomesExternalTypeCode) {
		reserveStatusSwitch = true
		if contactsId > 0 {
			reserveRecordInfo, err = thirdDao.NewReserveRecordDao().GetInfo(ctx, thirdModel.ReserveRecordModel{
				MicroBaseModel: models.MicroBaseModel{
					BaseModel:    models.BaseModel{ID: int64(contactsId)},
					EnterpriseID: enterpriseId,
				},
			})
			if err != nil || reserveRecordInfo.ID == 0 {
				ev = ev.AnErr("查询三方预约记录", err)
			}
			ev = ev.Interface("预约信息", reserveRecordInfo)
			reserveStatus = drainage2.ReserveStatusYes
		}
	}

	// 查询员工的部门
	memberInfo, err := membersDao.NewMembers().GetMemberInfoAndDepartmentsByUserId(ctx, enterpriseId, []string{member.Userid})
	if err != nil || len(memberInfo) == 0 {
		ev = ev.AnErr("查询员工的部门异常或不存在数据", err)
		return false
	}

	ev = ev.Interface("员工的部门", memberInfo)

	// 查询员工部门信息
	storeDepartmentInfo, err := microWorkDao.NewStoreDepartment().GetInfo(ctx, microWorkModel.StoreDepartmentModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: enterpriseId,
		},
		WechatDepartmentID: uint(memberInfo[0].DepartmentThirdId),
	})

	ev = ev.Interface("查询员工部门信息", storeDepartmentInfo)

	if err != nil {
		ev = ev.AnErr("查询员工部门信息", err)
	}
	// 门店类型转换
	var storeTypeName string
	if storeDepartmentInfo.StoreType == microWorkModel.StoreTypeOle {
		storeTypeName = "ole"
	}
	if storeDepartmentInfo.StoreType == microWorkModel.StoreTypeRetail {
		storeTypeName = "retail"
	}
	if storeDepartmentInfo.ThirdStoreCode == microWorkModel.StoreTypeESA {
		storeTypeName = microWorkModel.StoreTypeESA
	}
	// 三方外部联系人信息
	var thirdUserInfo struct {
		Level               string `json:"level" comment:"会员等级"`
		CustomType          string `json:"customType" comment:"客户类型"`
		ReceiveCouponStatus string `json:"receiveCouponStatus" comment:"领券状态"`
		RealName            string `json:"real_name" comment:"真实姓名"`
	}
	if helper.InArray(welcome2.MemberLevel, welcomesExternalTypeCode) || helper.InArray(welcome2.MemberLevelFY26, welcomesExternalTypeCode) || helper.InArray(welcome2.CustomType, welcomesExternalTypeCode) || helper.InArray(welcome2.ReceiveCouponStatus, welcomesExternalTypeCode) {
		// 默认会员姓名为昵称
		thirdUserInfo.RealName = externalDetail.Base.Name

		// 查询work请求crm获取用户信息
		crmUserInfo, err := microwork.NewCrm().GetUserInfoByUnionId(ctx, microwork.UserInfoRequest{
			EnterpriseID: enterpriseId,
			UnionId:      externalDetail.Base.UnionID,
		})
		if err != nil {
			ev = ev.AnErr("查询work请求crm获取用户信息", err)
		}

		ev = ev.Interface("crm获取用户信息", crmUserInfo)

		// 判断是否查询到用户信息(根据status返回判断) true已查到 false未查到
		if crmUserInfo.Data.Status {
			// 判断是否为会员
			if crmUserInfo.Data.CustomerLevel != "" {
				thirdUserInfo.Level = crmUserInfo.Data.CustomerLevel
			}
		}
		// 会员姓名
		if crmUserInfo.Data.Name != "" {
			thirdUserInfo.RealName = crmUserInfo.Data.Name
		}
	}

	// 客户类型开关
	var customTypeStatusSwitch bool
	// 小程序配置
	var miniProgram struct {
		AppId string `json:"app_id"`
		Path  string `json:"path"`
	}
	// 判断是否存在用户类型规则（新老用户，目前只有ks）
	if helper.InArray[string](welcome2.CustomType, welcomesExternalTypeCode) {
		customTypeStatusSwitch = true
		// 查询work请求crm获取用户是否为新老用户
		crmCheckMemberRes, err := microwork.NewCrm().CheckMemberByUnionId(ctx, microwork.UserInfoRequest{
			EnterpriseID: enterpriseId,
			UnionId:      externalDetail.Base.UnionID,
		})
		if err != nil {
			ev = ev.AnErr("查询work请求crm获取用户是否为新老用户", err)
		}

		ev = ev.Interface("crm获取用户是否为新老用户", crmCheckMemberRes)

		miniProgram.AppId = crmCheckMemberRes.Data.Config.AppId
		// 默认为新客户(根据status返回判断) true已查到 false未查到
		if crmCheckMemberRes.Data.Status {
			thirdUserInfo.CustomType = welcome2.CodeCustomTypeOldCustomer
			miniProgram.Path = crmCheckMemberRes.Data.Config.Path
		} else {
			thirdUserInfo.CustomType = welcome2.CodeCustomTypeNewCustomer
			miniProgramPath := crmCheckMemberRes.Data.Config.Path
			miniProgramPath = strings.Replace(miniProgramPath, "{shopCode}", storeDepartmentInfo.ThirdStoreCode, 1)
			miniProgramPath = strings.Replace(miniProgramPath, "{clerkCode}", memberInfo[0].MembersUserid, 1)
			miniProgram.Path = miniProgramPath
		}
	}

	// 领券开关
	var receiveCouponStatusSwitch bool
	// 判断是否存在领券规则（新老用户，目前只有ks）
	if helper.InArray[string](welcome2.ReceiveCouponStatus, welcomesExternalTypeCode) {
		receiveCouponStatusSwitch = true
		thirdUserInfo.ReceiveCouponStatus = welcome2.ReceiveCouponStatusNo
		havedResult, err := microwork.NewCrm().WelfareHaved(ctx, microwork.WelfareHavedRequest{
			EnterpriseID: enterpriseId,
			UnionId:      externalDetail.Base.UnionID,
			ExtUserid:    callbackData.ExternalUserID,
			ShopCode:     storeDepartmentInfo.ThirdStoreCode,
			ClerkCode:    memberInfo[0].MembersUserid,
			StoreType:    storeDepartmentInfo.StoreType,
		})
		if err != nil {
			ev = ev.AnErr("查询work请求crm判断是否存在领券规则", err)
		}
		if havedResult.Data.Haved {
			thirdUserInfo.ReceiveCouponStatus = welcome2.ReceiveCouponStatusYes
		}

		ev = ev.Interface("crm用户是否领券", havedResult)
	}

	// 渠道来源开关
	var sourceSwitch bool
	for _, code := range welcomesExternalTypeCode {
		if strings.Contains(code, "source_") {
			sourceSwitch = true
		}
	}

	var attachments = make([]map[string]interface{}, 0)
	// 循环判断是否满足配置的条件不满足其中一条则退出
	for _, wc := range allWelcome {
		ev = ev.Interface("个性化欢迎语规则", wc)

		var conditionContent map[string][]string
		helper.Unmarshal(wc.ConditionContent, &conditionContent)
		// 取出所有code值
		var codeValues []string
		for _, codes := range conditionContent {
			if len(codes) > 0 {
				codeValues = append(codeValues, codes...)
			}
		}

		// 渠道来源
		if sourceSwitch {
			if !helper.InArray(callbackData.State, codeValues) {
				ev = ev.Interface("渠道来源不满足", []any{wc.ID, callbackData.State, codeValues})
				continue
			}
		}

		// 预约回复(预约记录存在,且欢迎语规则不是预约则跳出)
		if reserveStatusSwitch {
			val, ok := conditionContent["reserve_status"]
			if ok && !helper.InArray(reserveStatus, val) {
				ev = ev.Interface("预约回复不满足", []any{wc.ID, reserveStatus, val})
				continue
			}
		}

		// 门店类型（零售，奥莱）
		if val, ok := conditionContent["store_type"]; ok {
			if len(val) != 0 && !helper.InArray(storeTypeName, val) {
				ev = ev.Interface("门店类型不满足", []any{wc.ID, storeTypeName, val})
				continue
			}
		}

		if helper.CheckCoachTime() {
			// 会员等级
			if val, ok := conditionContent["member_level_26"]; ok {
				// 兼容非会员
				if thirdUserInfo.Level == "" {
					thirdUserInfo.Level = "No"
				}
				if !helper.InArray(thirdUserInfo.Level, val) {
					//如果会员等级不满足设置判断是否是26年之后的逻辑
					if !(helper.InArray("C", val) && helper.InArray(thirdUserInfo.Level, []string{"SILVER", "GOLD", "DIAMOND"})) {
						//FY26之后Coach的已购会员【C】包含银【SILVER】金【GOLD】、钻【DIAMOND】
						ev = ev.Interface("会员等级不满足", []any{wc.ID, thirdUserInfo.Level, val})
						continue
					}
				}
			}
		} else {
			// 会员等级
			if val, ok := conditionContent["member_level"]; ok {
				// 兼容非会员
				if thirdUserInfo.Level == "" {
					thirdUserInfo.Level = "No"
				}
				if !helper.InArray(thirdUserInfo.Level, val) {
					ev = ev.Interface("会员等级不满足", []any{wc.ID, thirdUserInfo.Level, val})
					continue
				}
			}
		}

		// 客户类型（新客户，老客户）
		if customTypeStatusSwitch {
			val, ok := conditionContent["custom_type"]
			if ok && !helper.InArray(thirdUserInfo.CustomType, val) {
				ev = ev.Interface("会员等级不满足", []any{wc.ID, thirdUserInfo.CustomType, val})
				continue
			}
		}

		// 是否领券
		if receiveCouponStatusSwitch {
			val, ok := conditionContent["receive_coupon_status"]
			if ok && !helper.InArray(thirdUserInfo.ReceiveCouponStatus, val) {
				ev = ev.Interface("是否领券不满足", []any{wc.ID, thirdUserInfo.ReceiveCouponStatus, val})
				continue
			}
		}

		err = json.Unmarshal([]byte(wc.Content), &attachments)
		if err != nil {
			ev.AnErr("满足发送条件,解析Content", err)
			continue
		}

		if len(attachments) > 0 {
			break
		}
	}

	if len(attachments) == 0 {
		ev = ev.Interface("未找到满足条件的欢迎语", attachments)
		return false
	}

	welcomeContentJson, err := s.genSendAttachments(ctx, attachments, enterpriseId, domain)
	if err != nil {
		ev = ev.AnErr("生成欢迎语组装信息", err)
		return false
	}

	welcomeContent := welcomeContentJson.Text.Content
	welcomeContent = strings.Replace(welcomeContent, "{顾客昵称}", externalDetail.Base.Name, -1)
	welcomeContent = strings.Replace(welcomeContent, "{会员姓名}", thirdUserInfo.RealName, -1)
	welcomeContent = strings.Replace(welcomeContent, "{门店名称}", storeDepartmentInfo.StoreName, -1)
	welcomeContent = strings.Replace(welcomeContent, "{导购姓名}", memberInfo[0].MembersName, -1)
	if reserveStatusSwitch && reserveRecordInfo.ID > 0 {
		welcomeContent = strings.Replace(welcomeContent, "{预约服务名称}", reserveRecordInfo.ServiceName, -1)
		welcomeContent = strings.Replace(welcomeContent, "{预约时间}", helper.GetFormatDateByParams(reserveRecordInfo.CreatedAt), -1)
	}
	welcomeContentJson.Text.Content = welcomeContent

	// 只针对ks企业动态替换小程序路径
	if customTypeStatusSwitch && miniProgram.AppId != "" {
		for _, attachmentsArr := range welcomeContentJson.Attachments {
			for key, val := range attachmentsArr {
				if key == "miniprogram" {
					tmpVal := val.(work.MiniProgram)
					tmpVal.Appid = miniProgram.AppId
					tmpVal.Page = miniProgram.Path
					attachmentsArr[key] = tmpVal
				}
			}
		}
	}
	welcomeContentJson.WelcomeCode = callbackData.WelcomeCode

	err = workImpl.SendWelcome(ctx, *welcomeContentJson)
	ev = ev.Interface("发送欢迎语", welcomeContentJson).Err(err)
	return true
}

func (s Service) setCrmDataInfo(c ctx.Context, enterpriseId int64, unionId string) {
	//sw小程序路径添加参数
	resConfig, _ := cacheSvc.NewRedisData().GetEnterpriseConfig(c, enterpriseId, "store_name")
	if resConfig.Name == "sw" {
		log.Info(c).Msgf("外部联系人回调获取客户会员信息:%s", unionId)
		// 查询work请求crm获取用户信息
		crmUserInfo, _ := microwork.NewCrm().GetUserInfoByUnionId(c, microwork.UserInfoRequest{
			EnterpriseID: enterpriseId,
			UnionId:      unionId,
		})
		log.Info(c).Msgf("回调获取客户会员信息结果：%+v", crmUserInfo)
		if crmUserInfo.Data.CustomerCode != "" {
			err := microScrmDao.NewWxUserCorpExternal().UpdateData(c, enterpriseId, map[string]any{"unionid": unionId}, map[string]interface{}{
				"customer_code":            crmUserInfo.Data.CustomerCode,
				"card_level":               crmUserInfo.Data.CustomerLevel,
				"birthday":                 helper.EncodeToString(crmUserInfo.Data.Birthday),
				"birth_year":               helper.EncodeToString(helper.Int64ToString(int64(crmUserInfo.Data.BirthYear))),
				"birth_month":              helper.EncodeToString(helper.Int64ToString(int64(crmUserInfo.Data.BirthMonth))),
				"birth_day":                helper.EncodeToString(helper.Int64ToString(int64(crmUserInfo.Data.BirthDay))),
				"customer_type_updated_at": time.Now().Unix(),
				"updated_at":               time.Now().Unix(),
			})
			log.Info(c).Msgf("回调获取客户会员信息更新结果：%+v", err)
		}
	}
}

func (s Service) getStoreThirdIds(c ctx.Context, enterpriseId int64) []int64 {
	var storeThirdIds []int64
	// 获取企业配置的门店父级部门id信息
	configData, _ := boss.EnterpriseConfig(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"key":           "store_root_dept_id",
	})
	storeParentThirdId := int(configData.Data[0].Value) //所有门店的父级企微部门id
	departments, _ := departmentsDao.NewDepartments().GetListByThirdParentId(c, enterpriseId, storeParentThirdId)
	if len(departments) > 0 {
		for _, storeDept := range departments {
			storeThirdIds = append(storeThirdIds, int64(storeDept.ThirdID))
		}
	}
	return storeThirdIds
}
