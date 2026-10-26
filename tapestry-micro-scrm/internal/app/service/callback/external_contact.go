package callbackSvc

import (
	"errors"
	"fmt"
	"micro-scrm/internal/app"
	activityQRCodeDao "micro-scrm/internal/app/dao/activity_qrcode"
	dashboardDao "micro-scrm/internal/app/dao/dashboard"
	drainageContactDao "micro-scrm/internal/app/dao/drainage"
	externalSearchDao "micro-scrm/internal/app/dao/external_search"
	membersDao "micro-scrm/internal/app/dao/members"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	activityQRCodeModel "micro-scrm/internal/app/models/activity_qrcode"
	excitationModels "micro-scrm/internal/app/models/excitation"
	externalSearchModel "micro-scrm/internal/app/models/external_search"
	membersModel "micro-scrm/internal/app/models/members"
	microBossModel "micro-scrm/internal/app/models/micro_boss"
	"micro-scrm/internal/app/service"
	excitationSvc "micro-scrm/internal/app/service/excitation"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/boss"
	"micro-scrm/internal/pkg/sdk/cgema"
	"micro-scrm/internal/pkg/sdk/microwork"
	"micro-scrm/internal/pkg/sdk/work"
	"strings"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
)

// 企微回调
var (
	// AddExternalContact 添加企业客户事件
	AddExternalContact = "add_external_contact"
	// EditExternalContact 编辑企业客户事件
	EditExternalContact = "edit_external_contact"
	// AddHalfExternalContact 外部联系人免验证添加成员事件
	AddHalfExternalContact = "add_half_external_contact"
	// DelExternalContact 成员删除外部联系人
	DelExternalContact = "del_external_contact"
	// DelFollowUser 外部联系人员删除成员
	DelFollowUser = "del_follow_user"
	// TransferFail 客户接替失败事件
	TransferFail = "transfer_fail"
	// MigrateExternalContact 迁移企微历史数据
	MigrateExternalContact = "migrate_external_contact"
)

// 三方推送
var (
	PushThirdCdpEventAdd = "push_third_cdp_add"
	PushThirdCdpEventDel = "push_third_cdp_del"
)

type ExternalContact struct {
	service.Base

	enterprise   microBossModel.EnterpriseModel
	member       membersModel.MemberModel
	externalUser *work.ExternalContactResponse
	workImpl     *work.QyWorkImpl
	callbackData *work.MemberExternalCallback
}

func NewExternalContact() *ExternalContact {
	return &ExternalContact{}
}

func (s *ExternalContact) Handel(c ctx.Context, callbackData *work.MemberExternalCallback, workImpl *work.QyWorkImpl) error {
	s.workImpl = workImpl
	s.callbackData = callbackData

	var err error
	// 查询企业信息
	s.enterprise, err = microBossDao.NewEnterprise().GetOne(c, microBossModel.EnterpriseModel{
		CorpID: callbackData.ToUserName,
	})
	if err != nil {
		s.Log().Error(c).Msgf("查询企业信息异常 ：%v", err)
		return err
	}

	s.Log().Info(c).Msgf("外部联系人回调,事件类型：%s,事件数据：%v,企业数据：%v", callbackData.ChangeType, callbackData, s.enterprise)

	s.member, err = membersDao.NewMembers().GetOne(c,
		map[string]interface{}{
			"enterprise_id": s.enterprise.ID,
			"userid":        callbackData.UserID,
		})
	if err != nil {
		s.Log().Error(c).Msgf("查询员工信息异常 ：%v", err)
		return err
	}

	s.externalUser, err = s.workImpl.ExternalDetail(c, work.ExternalContactRequest{
		ExternalUserID: callbackData.ExternalUserID,
		QyGwRequest:    work.QyGwRequest{},
	})
	if err != nil {
		if !s.isDeleteExternalContactEvent(callbackData.ChangeType) || !s.isWorkExternalRelationMissing(err) {
			s.Log().Error(c).Msgf("查询外部联系人信息异常 ：%v", err)
			return err
		}
		if err = s.fillExternalUserByLocal(c, callbackData.ExternalUserID); err != nil {
			s.Log().Error(c).Msgf("企微外部联系人关系不存在,查询本地外部联系人信息异常 ：%v", err)
			return err
		}
	}

	switch callbackData.ChangeType {
	case AddExternalContact, MigrateExternalContact, AddHalfExternalContact:
		threading.GoSafe(func() error {
			if err = s.asyncAddExternalContact(c, callbackData); err != nil {
				s.Log().Error(c).Msgf("异步处理添加外部联系人,异常 ：%v", err)
			}
			return err
		})
		err = s.syncAddExternalContact(c, callbackData)
	case EditExternalContact:
		threading.GoSafe(func() error {
			if err = s.asyncEditExternalContact(c, callbackData); err != nil {
				s.Log().Error(c).Msgf("异步处理编辑外部联系人,异常 ：%v", err)
			}
			return err
		})
		err = s.syncEditExternalContact(c, callbackData)
	case DelFollowUser:
		threading.GoSafe(func() error {
			if err = s.asyncDelFollowUser(c, callbackData); err != nil {
				s.Log().Error(c).Msgf("异步处理外部联系人,异常 ：%v", err)
			}
			return err
		})
		err = s.syncDelFollowUser(c, callbackData)
	case DelExternalContact:
		err = s.syncDelExternalContact(c, callbackData)
	}

	return err
}

func (s *ExternalContact) isDeleteExternalContactEvent(changeType string) bool {
	return changeType == DelFollowUser || changeType == DelExternalContact
}

func (s *ExternalContact) isWorkExternalRelationMissing(err error) bool {
	var errx *errorx.Error
	return errors.As(err, &errx) && errx.Code == 84061
}

func (s *ExternalContact) fillExternalUserByLocal(c ctx.Context, externalUserID string) error {
	externalUser, err := microScrmDao.NewWxUserCorpExternal().FindOne(c, s.enterprise.ID, externalUserID)
	if err != nil {
		return err
	}
	if externalUser.Unionid == "" {
		s.externalUser = &work.ExternalContactResponse{}
		return nil
	}
	s.externalUser = &work.ExternalContactResponse{
		Base: work.ExternalContactBase{
			ExternalUserID: externalUserID,
			UnionID:        externalUser.Unionid,
			Name:           externalUser.Nick,
			Avatar:         externalUser.Avatar,
		},
	}
	s.Log().Info(c).Msgf("企微外部联系人关系不存在,使用本地外部联系人信息继续处理删除回调, externalUserID:%s, unionid:%s", externalUserID, externalUser.Unionid)
	return nil
}

// 同步添加处理
func (s *ExternalContact) syncAddExternalContact(ctx ctx.Context, callbackData *work.MemberExternalCallback) error {
	return nil
}

// 异步添加处理
func (s *ExternalContact) asyncAddExternalContact(ctx ctx.Context, callbackData *work.MemberExternalCallback) error {
	// 三方推送添加员工时间
	err := s.PushThirdCdp(ctx, PushThirdCdpEventAdd, callbackData.State)
	// 新增好友加积分
	s.producePoints(ctx)
	return err
}

// 同步编辑处理
func (s *ExternalContact) syncEditExternalContact(ctx ctx.Context, callbackData *work.MemberExternalCallback) error {
	var err error
	return err
}

// 异步编辑处理
func (s *ExternalContact) asyncEditExternalContact(ctx ctx.Context, callbackData *work.MemberExternalCallback) error {
	var err error
	return err
}

// 同步删除处理
func (s *ExternalContact) syncDelFollowUser(ctx ctx.Context, callbackData *work.MemberExternalCallback) error {
	var err error
	return err
}

// 异步删除处理
func (s *ExternalContact) asyncDelFollowUser(ctx ctx.Context, callbackData *work.MemberExternalCallback) error {
	var err error
	err = s.PushThirdCdp(ctx, PushThirdCdpEventDel, callbackData.State)
	return err
}

// 同步处理成员删除外部联系人
func (s *ExternalContact) syncDelExternalContact(ctx ctx.Context, callbackData *work.MemberExternalCallback) error {
	var err error
	err = s.PushThirdCdp(ctx, PushThirdCdpEventDel, callbackData.State)
	return err
}

// PushThirdCdp 三方推送员工事件
func (s *ExternalContact) PushThirdCdp(ctx ctx.Context, event string, state string) error {
	if !s.getBossConfigSwitch(ctx, boss.EnterpriseConfigThirdAssignmentPushKey) {
		s.Log().Debug(ctx).Msgf("获取推送配置信息失败")
		return nil
	}

	if s.externalUser.Base.UnionID == "" {
		return nil
	}

	var err error
	switch event {
	case PushThirdCdpEventAdd:
		err = s.PushThirdCdpEventAddHandle(ctx, state)
	case PushThirdCdpEventDel:
		err = s.PushThirdCdpEventDelHandle(ctx)
	}

	return err
}

func (s *ExternalContact) PushThirdCdpEventAddHandle(ctx ctx.Context, state string) error {
	// 获取员工和外部联系人是否存在关系
	userRelation, err := microScrmDao.NewUserRelation().GetOne(ctx, map[string]interface{}{
		"enterprise_id": s.enterprise.ID,
		"unionid":       s.externalUser.Base.UnionID,
		"member_id":     s.member.ID,
	})
	if err != nil {
		s.Log().Error(ctx).Msgf("查询外部联系人关系异常 ：%v", err)
	}

	// 首次绑定时间/最后绑定时间
	firstBindDate := helper.GetFormatDateByParams(0)
	latestBindDate := firstBindDate
	if userRelation.ID > 0 {
		firstBindDate = helper.GetFormatDateByParams(int64(userRelation.Createtime))
	}
	bindTime := helper.GetFormatDateByParams(0)
	// 获取员工详情信息
	membersInfo, err := membersDao.NewMembers().GetMemberInfo(ctx,
		map[string]interface{}{
			"enterprise_id": s.enterprise.ID,
			"userid":        s.callbackData.UserID,
		})
	if err != nil {
		s.Log().Error(ctx).Msgf("获取员工详情信息异常 ：%v", err)
		return err
	}

	if len(membersInfo.MemberDepartments) == 0 {
		s.Log().Error(ctx).Msgf("员工详情未存在部门信息")
		return err
	}

	source := ""
	if state != "" {
		drainageContactInfo, err := drainageContactDao.NewDrainageContactsDao().GetOne(ctx, map[string]interface{}{
			"id": state,
		})
		if err != nil {
			s.Log().Error(ctx).Msgf("查询联系我渠道信息错误:%+v", err)
		}
		source = drainageContactInfo.Name
	}
	for _, department := range membersInfo.MemberDepartments {
		if department.Department.Name == "ESA" {
			s.Log().Debug(ctx).Msgf("%d::开始给ec_hub推送绑定数据", s.enterprise.ID)
			_, err = microwork.NewEcHub().SaAssignment(ctx, microwork.SaAssignmentEcHubReq{
				EnterpriseID: s.enterprise.ID,
				Params: microwork.SaAssignmentParamsEcHubReq{
					SaAccount:      s.callbackData.UserID,
					FirstBindDate:  firstBindDate,
					LatestBindDate: latestBindDate,
					UnionId:        s.externalUser.Base.UnionID,
					Nickname:       s.externalUser.Base.Name,
					AvatarUrl:      s.externalUser.Base.Avatar,
					Channel:        source,
				},
			})
			break
		}
	}
	// 获取当前员工的门店编码
	storeDepartmentInfo, err := microWorkDao.NewStoreDepartment().GetInfo(ctx, map[string]any{
		"enterprise_id":        s.enterprise.ID,
		"wechat_department_id": membersInfo.MemberDepartments[0].Department.ThirdID,
	})
	if err == nil {
		s.pushCgeRelationBind(ctx, storeDepartmentInfo.ThirdStoreCode, firstBindDate, latestBindDate, bindTime)
	}
	if err != nil {
		s.Log().Error(ctx).Msgf("获取当前员工的门店编码异常:%+v", err)
		return err
	}

	s.Log().Debug(ctx).Msgf("%d::开始给acxiom推送绑定数据", s.enterprise.ID)
	_, err = microwork.NewCdp().SaAssignment(ctx, microwork.SaAssignmentReq{
		EnterpriseID: s.enterprise.ID,
		Params: microwork.SaAssignmentParamsReq{
			SaAccount:      s.callbackData.UserID,
			FirstBindDate:  firstBindDate,
			LatestBindDate: latestBindDate,
			UnionId:        s.externalUser.Base.UnionID,
			StoreCode:      storeDepartmentInfo.ThirdStoreCode,
			Source:         source,
			BindTime:       bindTime,
		},
	})

	// 查询活动二维码点位id 推送神策
	pointInfo, _ := activityQRCodeDao.NewPointDao().GetOne(ctx, s.enterprise.ID, "id = ?", state)
	if pointInfo.ID > 0 {
		categoryModel, _ := activityQRCodeDao.NewCategoryDao().GetOne(ctx, s.enterprise.ID, "id = ?", pointInfo.CategoryID)
		activityInfo, _ := activityQRCodeDao.NewMainDao().GetOne(ctx, s.enterprise.ID, "id = ?", pointInfo.ActivityID)
		if activityInfo.ID > 0 {
			var scopeStrSlice []string
			memberScopes, _ := activityQRCodeDao.NewMemberScopeDao().GetMany(ctx, s.enterprise.ID, "activity_id = ?", pointInfo.ActivityID)
			// 转换成员范围
			var areaIds []int
			var memberIds []int64
			for _, ms := range memberScopes {
				if ms.AreaId > 0 {
					areaIds = append(areaIds, ms.AreaId)
				}
				if ms.MemberId > 0 {
					memberIds = append(memberIds, ms.MemberId)
				}
			}
			memberList, _ := membersDao.NewMembers().GetUnscopedMemberInfoByIds(ctx, s.enterprise.ID, memberIds)
			for _, model := range memberList {
				scopeStrSlice = append(scopeStrSlice, model.Name)
			}
			areaList, _ := dashboardDao.NewArea().GetMany(ctx, s.enterprise.ID, "id in ?", areaIds)
			for _, model := range areaList {
				scopeStrSlice = append(scopeStrSlice, model.Name)
			}
			//获取导购管理柜台
			extraModel, _ := externalSearchDao.NewCdpExternalTagCk().GetOne(ctx, s.enterprise.ID,
				"unionid = ? AND tag_type = ?", s.externalUser.Base.UnionID, externalSearchModel.CdpTagManagerStoreCode)
			var addedStatusStr = "未添加过投放范围内的导购，且没有绑定对应的管理柜台"
			if extraModel.TagValue != "" && extraModel.TagValue == storeDepartmentInfo.ThirdStoreCode {
				addedStatusStr = "未添加过投放范围内的导购，有绑定过投放范围内的绑定柜台"
			}
			// 查询联系我二维码
			wayModel, _ := activityQRCodeDao.NewContactWayDao().GetOne(ctx, s.enterprise.ID, "state = ? AND third_id = ?", state, storeDepartmentInfo.WechatDepartmentID)
			var profiles = map[string]string{
				"button_type":         "添加",
				"app_name":            "Wecom_QRCode",
				"union_id":            s.externalUser.Base.UnionID,
				"campaign_name":       activityInfo.Name,
				"campaign_id":         helper.Int64ToString(activityInfo.ID),
				"valid_period_type":   helper.RuleThree(activityInfo.TimeType == activityQRCodeModel.TimeTypeLongTerm, "长期有效", "指定时间"),
				"campaign_channel":    categoryModel.Name + helper.RuleThree(pointInfo.Name == "", "", "_"+pointInfo.Name),
				"campaign_scope":      strings.Join(scopeStrSlice, ";"),
				"qr_poster_code_id":   state,
				"store_name":          storeDepartmentInfo.StoreName,
				"store_code":          storeDepartmentInfo.ThirdStoreCode,
				"SA_name":             membersInfo.Name,
				"SA_id":               s.callbackData.UserID,
				"QR_code_id":          wayModel.QrCode,
				"SA_added_status":     addedStatusStr,
				"$latest_utm_source":  helper.ConvertChineseToPinyin(strings.ReplaceAll(categoryModel.Name, " ", "_")),
				"$latest_utm_medium":  "qrcode",
				"$latest_utm_content": helper.ConvertChineseToPinyin(strings.ReplaceAll(activityInfo.Name, " ", "_")),
			}
			if activityInfo.TimeType != activityQRCodeModel.TimeTypeLongTerm {
				profiles["valid_period"] = helper.GetFormatDate(int64(activityInfo.StartTime), 1) + "至" + helper.GetFormatDate(int64(activityInfo.EndTime), 1)
			}
			if activityInfo.AliasName != "" {
				profiles["display_campaign_name"] = activityInfo.AliasName
				profiles["$latest_utm_campaign"] = helper.ConvertChineseToPinyin(strings.ReplaceAll(activityInfo.AliasName, " ", "_"))
			}
			if pointInfo.Name != "" {
				profiles["$latest_utm_term"] = helper.ConvertChineseToPinyin(strings.ReplaceAll(pointInfo.Name, " ", "_"))
			}
			_, err = microwork.NewShenCe().AddWecomSA(ctx, s.enterprise.ID, s.externalUser.Base.UnionID, profiles)
			if err != nil {
				return err
			}
		}
	}

	return err
}

func (s *ExternalContact) PushThirdCdpEventDelHandle(ctx ctx.Context) error {
	unbindTime := helper.GetFormatDateByParams(0)
	s.Log().Debug(ctx).Msgf("%d::开始给acxiom推送取消绑定数据", s.enterprise.ID)
	_, err := microwork.NewCdp().SaUnAssignment(ctx, microwork.SaUnAssignmentReq{
		EnterpriseID: s.enterprise.ID,
		Params: microwork.SaUnAssignmentParamsReq{
			SaAccount:  s.callbackData.UserID,
			UnionId:    s.externalUser.Base.UnionID,
			UnbindTime: unbindTime,
		},
	})
	membersInfo, err := membersDao.NewMembers().GetMemberInfo(ctx,
		map[string]interface{}{
			"enterprise_id": s.enterprise.ID,
			"userid":        s.callbackData.UserID,
		})
	if err != nil {
		s.Log().Error(ctx).Msgf("获取员工详情信息异常 ：%v", err)
		return err
	}

	if len(membersInfo.MemberDepartments) == 0 {
		s.Log().Error(ctx).Msgf("员工详情未存在部门信息")
		return err
	}

	storeCode := ""
	storeDepartmentInfo, storeErr := microWorkDao.NewStoreDepartment().GetInfo(ctx, map[string]any{
		"enterprise_id":        s.enterprise.ID,
		"wechat_department_id": membersInfo.MemberDepartments[0].Department.ThirdID,
	})
	if storeErr != nil {
		s.Log().Error(ctx).Msgf("获取当前员工的门店编码异常:%+v", storeErr)
	} else {
		storeCode = storeDepartmentInfo.ThirdStoreCode
	}
	s.pushCgeRelationUnbind(ctx, storeCode, unbindTime)

	for _, department := range membersInfo.MemberDepartments {
		if department.Department.Name == "ESA" {
			s.Log().Debug(ctx).Msgf("%d::开始给EcHub推送取消绑定数据", s.enterprise.ID)
			_, err = microwork.NewEcHub().SaUnAssignment(ctx, microwork.SaUnAssignmentEcHubReq{
				EnterpriseID: s.enterprise.ID,
				Params: microwork.SaUnAssignmentParamsEcHubReq{
					SaAccount:  s.callbackData.UserID,
					UnionId:    s.externalUser.Base.UnionID,
					UnbindDate: helper.GetFormatDateByParams(0),
				},
			})
			break
		}
	}
	userRelation, err := microScrmDao.NewUserRelation().GetOne(ctx, map[string]interface{}{
		"enterprise_id": s.enterprise.ID,
		"unionid":       s.externalUser.Base.UnionID,
		"user_id":       s.callbackData.UserID,
	})
	if userRelation.State != "" {
		// 查询活动二维码点位id 推送神策
		pointInfo, _ := activityQRCodeDao.NewPointDao().GetOne(ctx, s.enterprise.ID, "id = ?", userRelation.State)
		if pointInfo.ID > 0 {
			// 获取当前员工的门店编码
			storeDepartmentInfo, err := microWorkDao.NewStoreDepartment().GetInfo(ctx, map[string]any{
				"enterprise_id":        s.enterprise.ID,
				"wechat_department_id": membersInfo.MemberDepartments[0].Department.ThirdID,
			})
			categoryModel, _ := activityQRCodeDao.NewCategoryDao().GetOne(ctx, s.enterprise.ID, "id = ?", pointInfo.CategoryID)
			activityInfo, _ := activityQRCodeDao.NewMainDao().GetOne(ctx, s.enterprise.ID, "id = ?", pointInfo.ActivityID)
			if activityInfo.ID > 0 {
				var scopeStrSlice []string
				memberScopes, _ := activityQRCodeDao.NewMemberScopeDao().GetMany(ctx, s.enterprise.ID, "activity_id = ?", pointInfo.ActivityID)
				// 转换成员范围
				var areaIds []int
				var memberIds []int64
				for _, ms := range memberScopes {
					if ms.AreaId > 0 {
						areaIds = append(areaIds, ms.AreaId)
					}
					if ms.MemberId > 0 {
						memberIds = append(memberIds, ms.MemberId)
					}
				}
				memberList, _ := membersDao.NewMembers().GetUnscopedMemberInfoByIds(ctx, s.enterprise.ID, memberIds)
				for _, model := range memberList {
					scopeStrSlice = append(scopeStrSlice, model.Name)
				}
				areaList, _ := dashboardDao.NewArea().GetMany(ctx, s.enterprise.ID, "id in ?", areaIds)
				for _, model := range areaList {
					scopeStrSlice = append(scopeStrSlice, model.Name)
				}
				//获取导购管理柜台
				extraModel, _ := externalSearchDao.NewCdpExternalTagCk().GetOne(ctx, s.enterprise.ID,
					"unionid = ? AND tag_type = ?", s.externalUser.Base.UnionID, externalSearchModel.CdpTagManagerStoreCode)
				var addedStatusStr = "未添加过投放范围内的导购，且没有绑定对应的管理柜台"
				if extraModel.TagValue != "" && extraModel.TagValue == storeDepartmentInfo.ThirdStoreCode {
					addedStatusStr = "未添加过投放范围内的导购，有绑定过投放范围内的绑定柜台"
				}
				// 查询联系我二维码
				wayModel, _ := activityQRCodeDao.NewContactWayDao().GetOne(ctx, s.enterprise.ID, "state = ? AND third_id = ?", userRelation.State, storeDepartmentInfo.WechatDepartmentID)
				var profiles = map[string]string{
					"button_type":         "删除",
					"app_name":            "Wecom_QRCode",
					"union_id":            s.externalUser.Base.UnionID,
					"campaign_name":       activityInfo.Name,
					"campaign_id":         helper.Int64ToString(activityInfo.ID),
					"valid_period_type":   helper.RuleThree(activityInfo.TimeType == activityQRCodeModel.TimeTypeLongTerm, "长期有效", "指定时间"),
					"campaign_channel":    categoryModel.Name + helper.RuleThree(pointInfo.Name == "", "", "_"+pointInfo.Name),
					"campaign_scope":      strings.Join(scopeStrSlice, ";"),
					"qr_poster_code_id":   userRelation.State,
					"store_name":          storeDepartmentInfo.StoreName,
					"store_code":          storeDepartmentInfo.ThirdStoreCode,
					"SA_name":             membersInfo.Name,
					"SA_id":               s.callbackData.UserID,
					"QR_code_id":          wayModel.QrCode,
					"SA_added_status":     addedStatusStr,
					"$latest_utm_source":  helper.ConvertChineseToPinyin(strings.ReplaceAll(categoryModel.Name, " ", "_")),
					"$latest_utm_medium":  "qrcode",
					"$latest_utm_content": helper.ConvertChineseToPinyin(strings.ReplaceAll(activityInfo.Name, " ", "_")),
				}
				if activityInfo.TimeType != activityQRCodeModel.TimeTypeLongTerm {
					profiles["valid_period"] = helper.GetFormatDate(int64(activityInfo.StartTime), 1) + "至" + helper.GetFormatDate(int64(activityInfo.EndTime), 1)
				}
				if activityInfo.AliasName != "" {
					profiles["display_campaign_name"] = activityInfo.AliasName
					profiles["$latest_utm_campaign"] = helper.ConvertChineseToPinyin(strings.ReplaceAll(activityInfo.AliasName, " ", "_"))
				}
				if pointInfo.Name != "" {
					profiles["$latest_utm_term"] = helper.ConvertChineseToPinyin(strings.ReplaceAll(pointInfo.Name, " ", "_"))
				}
				_, err = microwork.NewShenCe().AddWecomSA(ctx, s.enterprise.ID, s.externalUser.Base.UnionID, profiles)
				if err != nil {
					return err
				}
			}
		}
	}

	return err
}

func (s *ExternalContact) pushCgeRelationBind(ctx ctx.Context, storeCode, firstBindTime, latestBindTime, bindTime string) {
	_, err := cgema.NewClient().PushRelationEvent(ctx, cgema.RelationEventRequest{
		EventType: "bind",
		EventTime: bindTime,
		Customer: cgema.RelationCustomer{
			UnionID:        s.externalUser.Base.UnionID,
			ExternalUserID: s.callbackData.ExternalUserID,
		},
		Sa: cgema.RelationSa{
			UserID:    s.callbackData.UserID,
			StoreCode: storeCode,
		},
		Relation: cgema.RelationInfo{
			FirstBindTime:  firstBindTime,
			LatestBindTime: latestBindTime,
		},
	})
	if err != nil {
		s.Log().Error(ctx).Msgf("push cge relation bind failed: %+v", err)
	}
}

func (s *ExternalContact) pushCgeRelationUnbind(ctx ctx.Context, storeCode, unbindTime string) {
	_, err := cgema.NewClient().PushRelationEvent(ctx, cgema.RelationEventRequest{
		EventType:    "unbind",
		EventTime:    unbindTime,
		UnbindReason: s.cgeUnbindReason(),
		Customer: cgema.RelationCustomer{
			UnionID:        s.externalUser.Base.UnionID,
			ExternalUserID: s.callbackData.ExternalUserID,
		},
		Sa: cgema.RelationSa{
			UserID:    s.callbackData.UserID,
			StoreCode: storeCode,
		},
		Relation: cgema.RelationInfo{
			UnbindTime: unbindTime,
		},
	})
	if err != nil {
		s.Log().Error(ctx).Msgf("push cge relation unbind failed: %+v", err)
	}
}

func (s *ExternalContact) cgeUnbindReason() string {
	if s.callbackData.ChangeType == DelFollowUser {
		return "customer_delete_member"
	}
	return "member_delete_customer"
}

// 获取企业自定义配置
func (s *ExternalContact) getBossConfigSwitch(c app.Ctx, key string) bool {
	req := &boss.EnterpriseConfigRequest{
		EnterpriseId: s.enterprise.ID,
		Key:          key,
	}
	configData, configErr := boss.EnterpriseConfig(c, req)
	s.Log().Debug(c).Msgf("新增联系人验证回调,获取企业自定义配置", configData, configErr)
	if configErr != nil {
		s.Log().Error(c).Msgf("新增联系人验证回调,获取企业自定义配置:错误 %+v", configErr)
		return false
	}

	if len(configData.Data) > 0 && configData.Data[0].Value == 1 {
		return true
	}

	return false
}

// 新增好友加积分
func (s *ExternalContact) producePoints(ctx ctx.Context) {
	// 积分队列埋点
	excitationSvc.NewPointsProduceService().ProduceQueue(ctx, excitationReq.PointsProduceSave{
		MemberId:     s.member.ID,
		Type:         excitationModels.RuleCategoryId10,
		Value:        1,
		UniqueHash:   s.externalUser.Base.ExternalUserID + fmt.Sprintf("-%d-%d", s.member.ID, excitationModels.RuleCategoryId10),
		EnterpriseId: s.enterprise.ID,
	})
}
