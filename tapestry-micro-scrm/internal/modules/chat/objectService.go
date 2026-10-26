package chat

import (
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/customer"
	userEs "micro-scrm/internal/modules/customer/respository/elastic"
	customerRepository "micro-scrm/internal/modules/customer/respository/mysql"
	"micro-scrm/internal/modules/members"
	memberRespository "micro-scrm/internal/modules/members/respository/mysql"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/uid"
)

// ObjectService 聊天对象
type ObjectService struct {
	repo ObjectRepository
}

// MemberToMemberList 员工聊天列表
func (s *ObjectService) MemberToMemberList(c ctx.Context, cond MemberToMemberPageListRequest) (MemberToMemberPageListResponse, error) {
	var cmp MemberToMemberPageListResponse

	offset := (cond.Page - 1) * cond.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(cond.Size)}
	option := &repos.FindPageOption{
		TplPrefix: "t_chat_object/member-to-member-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]MemberToMemberInfo{}
		},
	}
	total, list, err := s.repo.FindPage(c, pageVo, cond, option)
	if errorIgnore.Ignore(err) != nil {
		return cmp, err
	}
	resList := list.(*[]MemberToMemberInfo)

	cmp.Page = cond.Page
	cmp.Size = cond.Size
	cmp.List = *resList
	cmp.Total = total

	return cmp, nil
}

// MemberToExternalList 员工与外部联系人聊天列表
func (s *ObjectService) MemberToExternalList(c ctx.Context, cond MemberToExternalPageListRequest) (MemberToExternalPageListResponse, error) {
	var cmp MemberToExternalPageListResponse

	offset := (cond.Page - 1) * cond.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(cond.Size)}
	option := &repos.FindPageOption{
		TplPrefix: "t_chat_object/member-to-external-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]MemberToExternalInfo{}
		},
	}
	total, list, err := s.repo.FindPage(c, pageVo, cond, option)
	if errorIgnore.Ignore(err) != nil {
		return cmp, err
	}
	resList := list.(*[]MemberToExternalInfo)

	cmp.Page = cond.Page
	cmp.Size = cond.Size
	cmp.List = *resList
	cmp.Total = total

	return cmp, nil
}

// MemberToGroupList 员工与群聊聊天列表
func (s *ObjectService) MemberToGroupList(c ctx.Context, cond MemberToGroupPageListRequest) (MemberToGroupPageListResponse, error) {
	var cmp MemberToGroupPageListResponse

	offset := (cond.Page - 1) * cond.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(cond.Size)}
	option := &repos.FindPageOption{
		TplPrefix: "t_chat_object/member-to-group-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]MemberToGroupInfo{}
		},
	}
	total, list, err := s.repo.FindPage(c, pageVo, cond, option)
	if errorIgnore.Ignore(err) != nil {
		return cmp, err
	}
	resList := list.(*[]MemberToGroupInfo)

	for key, value := range *resList {
		if value.OwnerMemberID == cond.MemberID {
			(*resList)[key].IsOwner = 1
		}
	}

	cmp.Page = cond.Page
	cmp.Size = cond.Size
	cmp.List = *resList
	cmp.Total = total

	return cmp, nil
}

// Add 新增聊天对象记录
func (s *ObjectService) Add(c ctx.Context, chatService Service, memberRelationService MemberService, row Object) error {
	t := time.Now().Unix()

	row.AgreeStatus = AgreeStatusDefaultAgree
	row.IsReply = 2
	row.LastChatAt = t
	row.CreatedAt = t
	row.UpdatedAt = t
	row.DeletedAt = 0

	// 找出员工ID,外部联系人ID
	var memberIDArr []int64
	memberService := members.New(memberRespository.New())
	customerService := customer.New(customerRepository.NewCustomer(), customerRepository.NewWorkUser(), userEs.NewEsRepo())

	// 员工与员工发送
	if row.FromChannel == UserTypeMember && row.ToChannel == UserTypeMember {
		fromMemberInfo := memberService.FindMemberByUserid(c, members.Member{
			Userid:       row.FromUserid,
			EnterpriseID: row.EnterpriseID,
			DeletedAt:    0,
		})
		if fromMemberInfo != nil {
			row.FromObjectID = fromMemberInfo.ID
			memberIDArr = append(memberIDArr, fromMemberInfo.ID)
		}

		toMemberInfo := memberService.FindMemberByUserid(c, members.Member{
			Userid:       row.ToUserid,
			EnterpriseID: row.EnterpriseID,
			DeletedAt:    0,
		})
		if toMemberInfo != nil {
			row.ToObjectID = toMemberInfo.ID
			memberIDArr = append(memberIDArr, toMemberInfo.ID)
		}
	}

	var agreeInfoReq AgreeInfoRequest
	agreeInfoReq.EnterpriseID = row.EnterpriseID
	agreeInfoListReq := AgreeInfoRequestListInfo{}

	// 员工发送给外部联系人
	if row.FromChannel == UserTypeMember && (row.ToChannel == UserTypeWorkExternalUser || row.ToChannel == UserTypeExternalUser) {

		fromMemberInfo := memberService.FindMemberByUseridStatus(c, members.Member{
			Userid:       row.FromUserid,
			EnterpriseID: row.EnterpriseID,
			Status:       MemberStatusNormal,
			DeletedAt:    0,
		})
		if fromMemberInfo.ID != 0 {
			row.FromObjectID = fromMemberInfo.ID
			memberIDArr = append(memberIDArr, fromMemberInfo.ID)
		}

		// 找外部联系人ID
		externalUser, findExternalUserErr := customerService.FindExternal(c, customer.WorkExternalUser{
			EnterpriseID:   row.EnterpriseID,
			ExternalUserID: row.ToUserid,
		})
		if findExternalUserErr != nil {
			log.Error(ctx.Wrap(c)).Msgf("处理聊天对象,查找外部联系人ID错误:" + findExternalUserErr.Error())
		} else {
			row.ToObjectID = externalUser.ID
		}

		agreeInfoListReq.Userid = row.FromUserid
		agreeInfoListReq.ExteranalUserid = row.ToUserid
	} else if (row.FromChannel == UserTypeWorkExternalUser || row.FromChannel == UserTypeExternalUser) && row.ToChannel == UserTypeMember {
		// 如果是外部联系人发给员工
		toMemberInfo := memberService.FindMemberByUseridStatus(c, members.Member{
			Userid:       row.ToUserid,
			EnterpriseID: row.EnterpriseID,
			Status:       MemberStatusNormal,
			DeletedAt:    0,
		})
		if toMemberInfo.ID != 0 {
			row.ToObjectID = toMemberInfo.ID
			memberIDArr = append(memberIDArr, toMemberInfo.ID)
		}

		// 找外部联系人ID
		externalUser, findExternalUserErr := customerService.FindExternal(c, customer.WorkExternalUser{
			EnterpriseID:   row.EnterpriseID,
			ExternalUserID: row.FromUserid,
		})
		if findExternalUserErr != nil {
			log.Error(ctx.Wrap(c)).Msgf("处理聊天对象,查找外部联系人ID错误:" + findExternalUserErr.Error())
		} else {
			row.FromObjectID = externalUser.ID
		}

		agreeInfoListReq.Userid = row.ToUserid
		agreeInfoListReq.ExteranalUserid = row.FromUserid
	} else if row.ToChannel == UserTypeGroup {
		groupRow := row

		// 获取该群聊下的员工
		fromMemberInfo, fromErr := s.repo.GetMemberByGroup(c, groupRow.ToUserid, groupRow.EnterpriseID)
		if fromErr != nil {
			log.Error(ctx.Wrap(c)).Msgf("群聊下无员工:%v", fromErr)
			return nil
		}

		if fromMemberInfo == nil || len(*fromMemberInfo) == 0 {
			return nil
		}

		groupRow.ToObjectID = (*fromMemberInfo)[0].GroupChatID
		for _, value := range *fromMemberInfo {
			groupRow.FromUserid = value.UserID
			groupRow.FromObjectID = value.MemberID

			// 判断记录是否存在,存在则更新row
			chatObject := s.Find(c, groupRow)
			if chatObject.ID != 0 {
				groupRow.ID = chatObject.ID
				_, updateErr := s.repo.Update(c, groupRow)
				if updateErr != nil {
					log.Error(ctx.Wrap(c)).Msgf("处理聊天对象,更新聊天对象错误:" + updateErr.Error())

					return updateErr
				}
			} else {
				groupRow.ID = uid.ID(c)
				_, addErr := s.repo.Add(c, groupRow)
				if addErr != nil {
					log.Error(ctx.Wrap(c)).Msgf("新增聊天对象记录错误:" + addErr.Error())

					return addErr
				}
			}

			// 记录员工聊天关系
			memberRelationService.Add(c, MemberRelation{
				ID:           uid.ID(c),
				MemberID:     value.MemberID,
				LastChatAt:   t,
				EnterpriseID: groupRow.EnterpriseID,
				CreatedAt:    t,
				UpdatedAt:    t,
				DeletedAt:    0,
			})
		}

	}

	if row.ToChannel != UserTypeGroup {

		// 查外部联系人是否开启了同意
		agreeInfoReq.List = append(agreeInfoReq.List, agreeInfoListReq)
		agreeInfoListRes, _ := chatService.GetCheckSingleAgree(c, agreeInfoReq)

		if len(agreeInfoListRes) != 0 {
			row.AgreeStatus = chatService.ConversionAgreeStatus(agreeInfoListRes[0].AgreeStatus)
			row.AgreeStatusAt = agreeInfoListRes[0].StatusChangeTime
		}

		// 判断记录是否存在,存在则更新row
		chatObject := s.Find(c, row)
		if chatObject.ID != 0 {
			row.ID = chatObject.ID
			row.UpdatedAt = t
			_, updateErr := s.repo.Update(c, row)
			if updateErr != nil {
				log.Error(ctx.Wrap(c)).Msgf("处理聊天对象,更新聊天对象错误:" + updateErr.Error())

				return updateErr
			}
			// 倒置用户再次更新同意状态
			convertRow := Object{
				FromUserid:   row.ToUserid,
				ToUserid:     row.FromUserid,
				EnterpriseID: row.EnterpriseID,
			}
			convertObject := s.Find(c, convertRow)
			if chatObject.ID != 0 {
				convertObject.AgreeStatus = row.AgreeStatus
				convertObject.AgreeStatusAt = row.AgreeStatusAt
				convertObject.UpdatedAt = t
			}
			_, updateConvertErr := s.repo.Update(c, convertObject)
			if updateConvertErr != nil {
				log.Error(ctx.Wrap(c)).Msgf("处理聊天对象,更新聊天倒置对象错误:" + updateConvertErr.Error())

				return updateConvertErr
			}
		} else {
			if row.ID == 0 {
				row.ID = uid.ID(c)
			}
			_, addErr := s.repo.Add(c, row)
			if addErr != nil {
				log.Error(ctx.Wrap(c)).Msgf("新增聊天对象记录错误:" + addErr.Error())

				return addErr
			}

			return nil
		}

		// 发送人与接收人倒置找记录是否存在,存在则更新最后聊天时间 是否回复
		newT := time.Now().Unix()
		pourChatObject := s.Find(c, Object{
			FromUserid:   row.ToUserid,
			ToUserid:     row.FromUserid,
			EnterpriseID: row.EnterpriseID,
		})
		if pourChatObject.ID != 0 {
			pourChatObject.IsReply = 1
			pourChatObject.LastChatAt = newT
			pourChatObject.UpdatedAt = newT
			_, pourUpdateErr := s.repo.Update(c, pourChatObject)
			if pourUpdateErr != nil {
				log.Error(ctx.Wrap(c)).Msgf("处理聊天对象,更新倒置聊天对象错误:" + pourUpdateErr.Error())
			}
			// 再更新发送人的记录
			row.IsReply = 1
			row.LastChatAt = newT
			pourChatObject.UpdatedAt = newT
			_, originalUpdateErr := s.repo.Update(c, row)
			if originalUpdateErr != nil {
				log.Error(ctx.Wrap(c)).Msgf("处理聊天对象,更新原始聊天对象错误:" + originalUpdateErr.Error())
			}
		}

		// 记录员工聊天关系
		for _, memberID := range memberIDArr {
			memberRelationService.Add(c, MemberRelation{
				ID:           uid.ID(c),
				MemberID:     memberID,
				LastChatAt:   t,
				EnterpriseID: row.EnterpriseID,
				CreatedAt:    t,
				UpdatedAt:    t,
				DeletedAt:    0,
			})
		}
	}

	return nil
}

// Find 查一条记录
func (s *ObjectService) Find(c ctx.Context, cond Object) Object {
	var obj Object
	_ = s.repo.Find(c, cond, &obj)

	return obj
}
