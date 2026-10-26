package customer

import (
	"errors"
	"github.com/gin-gonic/gin"
	membersDao "micro-scrm/internal/app/dao/members"
	"micro-scrm/internal/app/http/handle"
	customerReq "micro-scrm/internal/app/http/request/customer"
	tagsReq "micro-scrm/internal/app/http/request/tags"
	membersModel "micro-scrm/internal/app/models/members"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	cacheSvc "micro-scrm/internal/app/service/cache"
	customerSvc "micro-scrm/internal/app/service/customer"
	membersSvc "micro-scrm/internal/app/service/member"
	workTagSvc "micro-scrm/internal/app/service/tags"
	"micro-scrm/internal/pkg/helper"
)

type ExternalUserHandle struct {
	handle.BaseHandle
}

func NewExternalUserHandle() ExternalUserHandle {
	return ExternalUserHandle{}
}

// List 根据条件筛选导购和外部联系人列表
func (h ExternalUserHandle) List(c *gin.Context) {

	var request customerReq.ExternalUserListsReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	if request.MemberType == 4 && len(request.ThirdIds) == 0 {
		h.ResError(c, errors.New("third_ids不能为空"))
		return
	}

	// 获取member信息
	memberInfo, err1 := membersSvc.NewMembersService().GetInfoByMobile(h.WarpContext(c), request.EnterpriseID, request.Mobile)
	if err1 != nil {
		h.ResError(c, err1)
		return
	}
	request.MemberId = memberInfo.ID
	request.Userid = memberInfo.Userid
	resConfig, err := cacheSvc.NewRedisData().GetEnterpriseConfig(h.WarpContext(c), request.EnterpriseID, "store_name")
	if err != nil {
		h.ResError(c, err)
	}

	// 先查ck如果ck异常就查mysql
	resultCk, totalCk, errCk := customerSvc.NewCustomer().GetExternalUserListCk(h.WarpContext(c), request, resConfig.Name)
	if errCk != nil {

		// MySQL查询业务
		result, total, err := customerSvc.NewCustomer().GetExternalUserList(h.WarpContext(c), request, resConfig.Name)
		if err != nil {
			h.ResError(c, err)
			return
		}
		h.ResList(c, total, request.Page, request.Size, result)

	} else {

		// ck返回结果
		h.ResList(c, totalCk, request.Page, request.Size, resultCk)
	}
	return
}

// ListCk 根据条件筛选导购和外部联系人列表
func (h ExternalUserHandle) ListCk(c *gin.Context) {

	var request customerReq.ExternalUserListsReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	if request.MemberType == 4 && len(request.ThirdIds) == 0 {
		h.ResError(c, errors.New("third_ids不能为空"))
		return
	}

	// 获取member信息
	memberInfo, err1 := membersSvc.NewMembersService().GetInfoByMobile(h.WarpContext(c), request.EnterpriseID, request.Mobile)
	if err1 != nil {
		h.ResError(c, err1)
		return
	}
	request.MemberId = memberInfo.ID
	request.Userid = memberInfo.Userid
	resConfig, err := cacheSvc.NewRedisData().GetEnterpriseConfig(h.WarpContext(c), request.EnterpriseID, "store_name")
	if err != nil {
		h.ResError(c, err)
	}

	result, total, err := customerSvc.NewCustomer().GetExternalUserListCk(h.WarpContext(c), request, resConfig.Name)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResList(c, total, request.Page, request.Size, result)

	return
}

// ChoiceOption 筛选外部联系人列表的条件
func (h ExternalUserHandle) ChoiceOption(c *gin.Context) {

	var request customerReq.ReqChoice
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if request.MemberType == 4 && len(request.ThirdIds) == 0 {
		h.ResError(c, errors.New("third_ids不能为空"))
		return
	}
	cc := h.WarpContext(c)
	resConfig, err := cacheSvc.NewRedisData().GetEnterpriseConfig(cc, request.EnterpriseID, "store_name")
	if err != nil {
		h.ResError(c, err)
		return
	}

	// SM管理的门店导购及自己
	var choiceOptionList customerReq.ChoiceOptionList

	if request.MemberType == 4 {
		// 员工筛选
		var memberList []customerReq.ChoiceList
		// 获取member信息
		memberInfo, err1 := membersSvc.NewMembersService().GetInfoByMobile(cc, request.EnterpriseID, request.Mobile)
		if err1 != nil {
			h.ResError(c, err1)
			return
		}
		// 查询店长管理的门店下店员 GetMemberByThirdIds
		memberLists, err1 := membersDao.NewMembers().GetMemberByThirdIds(cc, request.EnterpriseID, request.ThirdIds, []int64{})
		if err1 != nil {
			h.ResError(c, err1)
			return
		}

		isHave := 0
		for _, v := range memberLists {
			temp := customerReq.ChoiceList{Key: v.MembersId, Value: v.MembersName}
			if v.MembersId == memberInfo.ID {
				isHave = 1
			}
			memberList = append(memberList, temp)
		}
		if isHave == 0 {
			memberTemp := customerReq.ChoiceList{Key: memberInfo.ID, Value: memberInfo.Name}
			memberList = append(memberList, memberTemp)
		}
		choiceOptionList.MemberList = memberList
	}

	switch resConfig.Name {
	case "coach":
		// 邀约筛选文案
		var reserveType []customerReq.ReserveOption
		reserveType = append(reserveType,
			customerReq.ReserveOption{
				Key:   1,
				Value: "邀约中",
			},
			customerReq.ReserveOption{
				Key:   2,
				Value: "无邀约",
			},
			customerReq.ReserveOption{
				Key:   3,
				Value: "邀约成功",
			})
		// 任务筛选
		var taskType []customerReq.TaskOption
		taskType = append(taskType,
			customerReq.TaskOption{
				Key:   1,
				Value: "未触达",
			},
			customerReq.TaskOption{
				Key:   2,
				Value: "已触达",
			})
		// 近5天是否有任务
		var taskHave []customerReq.TaskOption
		taskHave = append(taskHave,
			customerReq.TaskOption{
				Key:   1,
				Value: "无任务",
			},
			customerReq.TaskOption{
				Key:   2,
				Value: "有任务",
			})
		// 任务顾客列表筛选 好友状态
		var isUnion []customerReq.TaskOption
		isUnion = append(isUnion,
			customerReq.TaskOption{
				Key:   1,
				Value: "好友顾客",
			},
			customerReq.TaskOption{
				Key:   2,
				Value: "非好友顾客",
			})
		// 优惠券筛选
		var couponType []customerReq.CouponOption
		couponType = append(couponType,
			customerReq.CouponOption{
				Key:   1,
				Value: "30天内临期券",
			},
			customerReq.CouponOption{
				Key:   2,
				Value: "待使用券",
			})
		// 近一年消费金额
		var consumeType []customerReq.ConsumeOption
		consumeType = append(consumeType,
			customerReq.ConsumeOption{
				Key:   7,
				Value: "30000及以上",
			},
			customerReq.ConsumeOption{
				Key:   6,
				Value: "20000-29999",
			},
			customerReq.ConsumeOption{
				Key:   5,
				Value: "10000-19999",
			},
			customerReq.ConsumeOption{
				Key:   4,
				Value: "8000-9999",
			},
			customerReq.ConsumeOption{
				Key:   3,
				Value: "5000-7999",
			},
			customerReq.ConsumeOption{
				Key:   2,
				Value: "3000-4999",
			},
			customerReq.ConsumeOption{
				Key:   1,
				Value: "2999及以下",
			})
		// 顾问筛选
		var manageType []customerReq.ChoiceList
		manageType = append(manageType,
			customerReq.ChoiceList{
				Key:   1,
				Value: "专属",
			},
			customerReq.ChoiceList{
				Key:   2,
				Value: "非专属",
			},
		)
		// 会员卡等,卡等信息是coach的会员卡等信息
		var cardType []customerReq.ChoiceOption
		cardType = append(cardType,
			//customerReq.ChoiceOption{
			//	Key:   customerReq.POTENTIAL,
			//	Value: "潜在客",
			//},
			customerReq.ChoiceOption{
				Key:   "O",
				Value: "非会员",
			},
			customerReq.ChoiceOption{
				Key:   "PROSPECT",
				Value: helper.GetCoachInfoByTime("品牌粉丝", "注册会员"),
			},
			customerReq.ChoiceOption{
				Key:   "C",
				Value: helper.GetCoachInfoByTime("新晋会员", "已购会员"),
			},
			//customerReq.ChoiceOption{
			//	Key:   "D",
			//	Value: "商场会员",
			//},
			customerReq.ChoiceOption{
				Key:   "SILVER",
				Value: "银卡会员",
			},
			customerReq.ChoiceOption{
				Key:   "GOLD",
				Value: "金卡会员",
			},
			customerReq.ChoiceOption{
				Key:   "DIAMOND",
				Value: "钻卡会员",
			},
		)
		// CLT类型
		var cltType []customerReq.ChoiceOption
		cltType = append(cltType,
			customerReq.ChoiceOption{
				Key:   "Top Client",
				Value: "Top Client",
			},
			customerReq.ChoiceOption{
				Key:   "New High Potential",
				Value: "New High Potential",
			},
		)
		var ageTag []customerReq.ChoiceList
		ageTag = append(ageTag,
			customerReq.ChoiceList{
				Key:   microScrmModel.AgeTagOne,
				Value: "Gen Z",
			}, customerReq.ChoiceList{
				Key:   microScrmModel.AgeTagTwo,
				Value: "YM",
			})
		choiceOptionList.TaskType = taskType
		choiceOptionList.ReserveType = reserveType
		choiceOptionList.CouponType = couponType
		choiceOptionList.ConsumeType = consumeType
		choiceOptionList.ManageType = manageType
		choiceOptionList.CardType = cardType
		choiceOptionList.CltType = cltType
		choiceOptionList.AgeTag = ageTag
		choiceOptionList.TaskHave = taskHave
		choiceOptionList.IsUnion = isUnion
		chatTagList, err := workTagSvc.NewWorkTags().GetChatTagList(cc, request.EnterpriseID)
		if err != nil {
			h.ResError(c, err)
			return
		}
		chatTag := make([]customerReq.ChoiceList, 0)
		for _, v := range chatTagList {
			temp := customerReq.ChoiceList{Key: v.ID, Value: v.Name}
			chatTag = append(chatTag, temp)
		}
		choiceOptionList.ChatTag = chatTag
		break
	case "sw":
		// 顾客类型
		var customerType []customerReq.ChoiceList
		customerType = append(customerType,
			customerReq.ChoiceList{
				Key:   1,
				Value: "新客",
			},
			customerReq.ChoiceList{
				Key:   2,
				Value: "老客",
			},
			customerReq.ChoiceList{
				Key:   3,
				Value: "VIP客",
			},
			customerReq.ChoiceList{
				Key:   4,
				Value: "沉睡客",
			},
		)
		// 会员卡等,卡等信息是sw的会员卡等信息
		var cardType []customerReq.ChoiceOption
		cardType = append(cardType,
			// customerReq.ChoiceOption{
			//	Key:   "Z",
			//	Value: "钻石「履」行家Z",
			// }, 默认选B带上Z
			customerReq.ChoiceOption{
				Key:   "B",
				Value: "钻石「履」行家",
			},
			customerReq.ChoiceOption{
				Key:   "C",
				Value: "金牌「履」行家",
			},
			customerReq.ChoiceOption{
				Key:   "W",
				Value: "新星「履」行家",
			},
		)

		choiceOptionList.CustomerType = customerType
		choiceOptionList.CardType = cardType
		break
	case "ks":
		// 会员卡等,卡等信息是ks的会员卡等信息
		var cardType []customerReq.ChoiceOption
		cardType = append(cardType,
			customerReq.ChoiceOption{
				Key:   "0",
				Value: "注册粉丝",
			},
			customerReq.ChoiceOption{
				Key:   "1",
				Value: "粉蜜会员",
			},
			customerReq.ChoiceOption{
				Key:   "2",
				Value: "银卡",
			},
			customerReq.ChoiceOption{
				Key:   "3",
				Value: "金卡",
			},
			customerReq.ChoiceOption{
				Key:   "4",
				Value: "钻卡",
			},
		)
		choiceOptionList.CardType = cardType
		break
	}

	// 客户企微标签
	var tagList []customerReq.ChoiceList
	tagLists, err := workTagSvc.NewWorkTags().GetChoiceTagList(cc, request.EnterpriseID)
	if err != nil {
		h.ResError(c, err)
		return
	}
	for _, v := range tagLists {
		temp := customerReq.ChoiceList{Key: v.ID, Value: v.Name}
		tagList = append(tagList, temp)
	}

	// 印象标签
	var impressionTagList []customerReq.ChoiceList
	var impressionTagReq tagsReq.ReqList
	impressionTagReq.EnterpriseID = request.EnterpriseID
	impressionTagLists, err := workTagSvc.NewImpression().GetList(cc, impressionTagReq)
	if err != nil {
		h.ResError(c, err)
		return
	}
	for _, v := range impressionTagLists {
		temp := customerReq.ChoiceList{Key: v.ID, Value: v.Name}
		impressionTagList = append(impressionTagList, temp)
	}

	choiceOptionList.TagList = tagList
	choiceOptionList.ImpressionTagList = impressionTagList

	h.ResSuccess(c, choiceOptionList)

	return
}

// CustomerList 三方接口的会员卡号查询信息
func (h ExternalUserHandle) CustomerList(c *gin.Context) {

	var request customerReq.CustomerListsReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if request.MemberType == 4 && len(request.ThirdIds) == 0 {
		h.ResError(c, errors.New("third_ids不能为空"))
		return
	}
	// 获取member信息
	memberInfo, err1 := membersSvc.NewMembersService().GetInfoByMobile(h.WarpContext(c), request.EnterpriseID, request.Mobile)
	if err1 != nil {
		h.ResError(c, err1)
		return
	}
	request.MemberId = memberInfo.ID
	request.Userid = memberInfo.Userid
	result, err := customerSvc.NewCustomer().GetCustomerList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)

	return
}

// Detail 外部联系人详情
func (h ExternalUserHandle) Detail(c *gin.Context) {

	// v3.2.0.1 查看外部联系人详情，不验证好友关系
	var request customerReq.ExternalUserDetailReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	// 当前导购member信息
	memberInfo, err1 := membersSvc.NewMembersService().GetInfoByMobile(h.WarpContext(c), request.EnterpriseId, request.Mobile)
	if err1 != nil {
		h.ResError(c, err1)
		return
	}

	if request.CustomerCode != "" && request.ExternalUserid == "" {
		// 查询customer_code
		var whereConf = map[string]interface{}{
			"customer_code": request.CustomerCode,
		}

		resultExternalUser, err := customerSvc.NewCustomer().GetCustomerDetail(h.WarpContext(c), request.EnterpriseId, whereConf)
		if err != nil {
			h.ResError(c, err)
			return
		}
		if resultExternalUser.Openid == "" {
			h.ResSuccess(c, err)
			return
		}
		// 不为空赋值
		request.ExternalUserid = resultExternalUser.Openid
	}

	// 根据open_id 获取外部联系人详情
	var whereUser = map[string]interface{}{
		"openid": request.ExternalUserid,
	}
	wxUserCorpExternalInfo, err := customerSvc.NewCustomer().GetCustomerDetail(h.WarpContext(c), request.EnterpriseId, whereUser)
	if err != nil {
		h.ResError(c, err)
		return
	}
	// 不验证关系
	userSARelation, err := customerSvc.NewCustomer().GetUserSARelation(h.WarpContext(c), request.EnterpriseId, memberInfo.ID, wxUserCorpExternalInfo.ID)
	if err != nil {
		h.ResError(c, err)
		return
	}
	resConfig, err := cacheSvc.NewRedisData().GetEnterpriseConfig(h.WarpContext(c), request.EnterpriseId, "store_name")
	if err != nil {
		h.ResError(c, err)
		return
	}
	// 对话框侧边栏打开 查看导购自己和客户的标签关系
	var memberIds []int64
	var memberUserIds []string
	memberIds = append(memberIds, memberInfo.ID)
	memberUserIds = append(memberUserIds, memberInfo.Userid)
	if request.IsWork == false {
		// 判断当前员工类型 店员或店长
		if request.MemberType == membersModel.MembersTypeSM {
			// 查询店长管理的门店下店员
			memberList, errs := membersDao.NewMembers().GetMemberByThirdIds(h.WarpContext(c), request.EnterpriseId, request.ThirdIds, []int64{})
			if errs != nil {
				h.ResError(c, errs)
				return
			}
			for _, v := range memberList {
				memberIds = append(memberIds, v.MembersId)
				memberUserIds = append(memberUserIds, v.MembersUserid)
			}
		}
	}

	// 各企业客户详情
	if resConfig.Name == "coach" {

		userDetail, errs := customerSvc.NewCustomer().CustomerDetailCoach(h.WarpContext(c), request, wxUserCorpExternalInfo, userSARelation, memberIds, memberUserIds)

		if errs != nil {
			h.ResError(c, errs)
			return
		}
		userDetail.ConfigStoreName = resConfig.Name

		h.ResSuccess(c, userDetail)

	}

	if resConfig.Name == "sw" {
		userDetail, errs := customerSvc.NewCustomer().CustomerDetailSW(h.WarpContext(c), request, wxUserCorpExternalInfo, userSARelation, memberIds, memberUserIds)

		if errs != nil {
			h.ResError(c, errs)
			return
		}
		userDetail.ConfigStoreName = resConfig.Name

		h.ResSuccess(c, userDetail)

	}

	if resConfig.Name == "ks" {
		userDetail, errs := customerSvc.NewCustomer().CustomerDetailKS(h.WarpContext(c), request, wxUserCorpExternalInfo, userSARelation, memberIds, memberUserIds)

		if errs != nil {
			h.ResError(c, errs)
			return
		}
		userDetail.ConfigStoreName = resConfig.Name

		h.ResSuccess(c, userDetail)

	}

	return

}

func (h ExternalUserHandle) GetMemberUserIds(c *gin.Context) {

	var request customerReq.ExternalUseridReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	if request.MemberType == 4 && len(request.ThirdIds) == 0 {
		h.ResError(c, errors.New("third_ids不能为空"))
		return
	}

	// 获取member信息
	memberInfo, err1 := membersSvc.NewMembersService().GetInfoByMobile(h.WarpContext(c), request.EnterpriseID, request.Mobile)
	if err1 != nil {
		h.ResError(c, err1)
		return
	}

	// 判断当前员工类型 店员或店长
	if request.MemberType == membersModel.MembersTypeSM { // SM
		// 查询店长管理的门店下店员 userid
		memberIds, err1 := membersDao.NewMemberDepartments().GetMemberIdsByThirdIds(h.WarpContext(c), request.EnterpriseID, request.ThirdIds, []int64{})
		if err1 != nil {
			h.ResError(c, err1)
		}
		// memberIds = helper.ArrayUniq(memberIds)

		// 店长和管理的店员
		memberIds = append(memberIds, memberInfo.ID)
		userIds, err1 := membersDao.NewMembers().GetUserIdsByMemberIDs(h.WarpContext(c), request.EnterpriseID, memberIds)
		userIds = append(userIds, memberInfo.Userid)

		// SA 店员
		h.ResSuccess(c, userIds)
		return

	} else {
		var userIds []string
		userIds = append(userIds, memberInfo.Userid)
		// SA 店员
		h.ResSuccess(c, userIds)
		return
	}

}

func (h ExternalUserHandle) SelectAllNum(c *gin.Context) {
	var request customerReq.ExternalUserListBaseReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	total, err := customerSvc.NewCustomer().GetExternalUserTotalCk(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, total)
}

func (h ExternalUserHandle) SelectAllList(c *gin.Context) {
	var request customerReq.ExternalUserListBaseReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	externalIds, err := customerSvc.NewCustomer().GetExternalUserIdsCk(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, externalIds)
}
