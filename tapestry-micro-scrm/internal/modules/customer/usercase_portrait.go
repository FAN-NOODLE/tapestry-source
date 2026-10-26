package customer

import (
	"fmt"
	"strconv"
	"time"

	"micro-scrm/internal/events"
	"micro-scrm/internal/events/entity"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/pkg/assist"
)

// RemarkSave 客户备注信息保存
func (s Service) RemarkSave(c ctx.Context, params *RemarkInfo) (int64, error) {

	// 查询是否已存在用户备注信息
	remarkData, _ := s.repo.GetUserRemark(c, ReturnUID{UserID: params.UserID})

	if remarkData == nil || remarkData.UserID == 0 {
		params.CreatedAt = assist.Timestamp()
	} else {
		params.UpdatedAt = assist.Timestamp()
	}

	// 生日格式转换
	if params.Birthday != 0 {
		timeLayout := "20060102"
		var birthdayDate string
		birthdayDate = time.Unix(params.Birthday, 0).Format(timeLayout)
		params.BirthdayDate, _ = strconv.Atoi(birthdayDate)
	}

	userID, error := s.repo.RemarkSave(c, *params)

	if error != nil || userID == 0 {
		return 0, error
	}

	return userID, nil
}

// SetExternalTags 设置外部联系人企业标签
func (s Service) SetExternalTags(c ctx.Context, params *SetExternalTags) (int64, error) {

	_, err := s.repo.SetExternalTags(c, *params)
	if err != nil {
		return 0, err
	}

	// 同步enterprise_user.groups|tags
	go events.ScrmEvent.ChangeUserTag.Dispatch(entity.Data{
		EnterpriseID: params.EnterpriseID,
		RelationID:   params.RelationID,
		Ctx:          c,
	})

	return 0, nil
}

// PortraitDetail 客户画像详情
func (s Service) PortraitDetail(c ctx.Context, request PortraitRequest) (PortraitResponse, error) {
	var portraitResponse PortraitResponse
	appid, err := s.enterpriseRepo.QueryOneApp(c, request.EnterpriseID, CustomerTypeExternal)
	if err != nil {
		return portraitResponse, err
	}

	user, err := s.repo.QueryEnterpriseUser(c, request.EnterpriseID, &SourceUserBase{
		AppID:        appid,
		OpenID:       request.ExternalUserID,
		CustomerType: CustomerTypeExternal,
	})

	if err != nil {
		return portraitResponse, err
	}

	appUser, ok := user.WxExternal[appid]
	if !ok {
		return portraitResponse, fmt.Errorf("用户数据异常:%s", request.ExternalUserID)
	}

	// 不查询好友关系
	//if request.MemberID == 0 {
	//	var tmpMemberId int64
	//	err = app.DB().Table(membersModel.MemberModel{}.TableName()).
	//		Where("userid", request.MemberUserID).
	//		Where("enterprise_id", request.EnterpriseID).
	//		Pluck("id", &tmpMemberId).Error
	//	if err != nil {
	//		return nil, err
	//	}
	//	request.MemberID = tmpMemberId
	//}

	// 先判断当前的外部联系和员工是否有关系
	relation, fErr := s.repo.FindCustomerRelation(c, request.MemberID, appUser.ID, request.EnterpriseID)
	//// 查询外部联系人关系出错
	if fErr != nil {
		return portraitResponse, fErr
	}

	// 未查询到外部联系人关系数据
	//if relation == nil {
	//	return nil, errs.ExternalRelationNotExistErr
	//}
	nickName := user.Nick
	if nickName == "" {
		nickName = user.Name
	}
	group := threading.NewRoutineGroup()
	portraitResponse = PortraitResponse{
		UserID:             user.ID,
		WorkExternalUserID: appUser.ID,
		ExternalUserID:     appUser.ThirdID,
		NickName:           nickName,
		Avatar:             user.Avatar,

		// 有好友关系
		RelationID: relation.ID,
		Remark:     relation.Remark,

		//Sex:                user.Sex,
		//CreateTime:         relation.CreateTime,
		//TagNames:           make([]string, 0),
		//GroupName:          make([]string, 0),
		//GenderIndex: int(user.Sex),
		//UserOrder: UserOrderStatistics{
		//	NumberOfOrders:   user.PayOrder,
		//	NumberOfRefunded: user.RefundOrder,
		//	OrdersAmount:     user.TotalFee,
		//	RefundAmount:     user.RefundFee,
		//},
	}

	for _, unionID := range user.WxUnion {
		portraitResponse.UnionID = unionID
		break
	}

	//if len(user.Tags) > 0 {
	//	// 根据relation_id获取客户的标签，
	//	if tagList, tagError := s.tagRepo.QueryWorkExternalTagList(c, user.Tags, relation.ID); tagError == nil {
	//		for _, item := range tagList {
	//			portraitResponse.TagNames = append(portraitResponse.TagNames, item.TagName)
	//		}
	//	}
	//}
	//if len(user.Groups) > 0 {
	//	if groups, err := s.tagRepo.QueryUserTagGroups(c, user.Groups); err == nil {
	//		for _, group := range groups {
	//			portraitResponse.GroupName = append(portraitResponse.GroupName, group.Name)
	//		}
	//	}
	//}
	group.RunSafe(func() {
		// 再查询客户备注信息
		remarkInfo, err := s.repo.GetUserRemark(c, ReturnUID{
			UserID: user.ID,
		})
		if err == nil && remarkInfo != nil {
			//portraitResponse.Birthday = remarkInfo.Birthday
			portraitResponse.Mobile = remarkInfo.Mobile
			portraitResponse.Gender = remarkInfo.Gender

			//portraitResponse.Age = remarkInfo.Age
			//portraitResponse.QQ = remarkInfo.Qq
			//portraitResponse.Email = remarkInfo.Email
			//portraitResponse.Describe = remarkInfo.Describe
		}
	})
	// 查询客户群个数
	//group.RunSafe(func() {
	//	total, cErr := s.repo.BelongToChatsCount(c, relation.EnterpriseID, user.WxExternal, nil, 1)
	//	if cErr == nil {
	//		portraitResponse.GroupNumber = total
	//	}
	//})

	// 查询渠道信息
	//group.RunSafe(func() {
	//	qrCodeID, _ := strconv.ParseInt(relation.State, 10, 64)
	//	_, channels := DrainageContact.ListChannelByQrCodeIds(c, []int64{qrCodeID}, request.EnterpriseID, 1, 1)
	//	if channels != nil {
	//		portraitResponse.ChannelName = channels[qrCodeID].ChannelName
	//	}
	//})
	//
	//if user.ID > 0 {
	//	// 查询绑定导购信息
	//	group.RunSafe(func() {
	//		guide, err := GuideService.FindGuideByUserIDForCustomer(c, request.EnterpriseID, user.ID)
	//		if err == nil {
	//			portraitResponse.BindGuideName = guide.BindMemberName
	//		}
	//	})
	//}

	group.Wait()
	return portraitResponse, nil
}

func isAjB(a, b []int64) bool {
	mapping := make(map[int64]interface{})
	for _, v := range a {
		mapping[v] = 0
	}

	for _, v := range b {
		if _, ok := mapping[v]; ok {
			return true
		}
	}

	return false
}
