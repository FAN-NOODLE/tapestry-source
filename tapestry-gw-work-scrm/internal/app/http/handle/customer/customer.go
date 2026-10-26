package customer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/json"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"gw-work-scrm/internal/app/http/common"
	"gw-work-scrm/internal/app/http/handle"
	"gw-work-scrm/internal/config"
	"gw-work-scrm/internal/pkg/redis"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
	microScrmCustomer "gw-work-scrm/internal/pkg/sdk/microscrm/customer"
	"gw-work-scrm/internal/pkg/sdk/work"
	"gw-work-scrm/internal/pkg/validate_params"
	"net/http"
	"time"
)

type CustomerHandle struct {
	handle.BaseHandle
	customerSdk     *microScrm.CustomerApi
	siderbarSdk     *microScrm.SidebarGroupChatApi
	sdkCustomerInfo microScrmCustomer.CustomerInfoApi
	customerApiSdk  microScrmCustomer.CustomerApi
}

func NewCustomerHandle() CustomerHandle {
	return CustomerHandle{
		customerSdk:     microScrm.NewCustomer(config.CfgData.Service.MicroScrm.Domain),
		siderbarSdk:     microScrm.NewSidebarGroupChat(config.CfgData.Service.MicroScrm.Domain),
		sdkCustomerInfo: microScrmCustomer.NewCustomerInfoApi(),
		customerApiSdk:  microScrmCustomer.NewCustomerApi(),
	}
}

// Detail 客户画像详情
func (h CustomerHandle) Detail(c *gin.Context) {
	req := &microScrm.CustomerDetailRequest{}

	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberID = memberInfo.ID
	req.MemberUserID = memberInfo.Userid

	// cacheKey := fmt.Sprintf("CustomerDetail:EnterpriseID:%d:MemberID:%d:Userid:%s", req.EnterpriseID, req.MemberID, req.MemberUserID)
	//
	// var result string
	// if err := redis.Client().Get(c, cacheKey).Scan(&result); err == nil {
	//	var newResult microScrm.PortraitDetailResponse
	//	_ = json.Unmarshal([]byte(result), &newResult)
	//	h.ResSuccess(c, newResult)
	//	return
	// }

	detail, err := h.customerSdk.UserDetail(ctx.Wrap(c), req)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	reqChat := &microScrm.GroupUsersAddEmployeesRequest{
		EnterpriseId:   req.EnterpriseID,
		ExternalUserId: req.ExternalUserID,
		Page:           1,
		Size:           3,
	}
	members, _ := h.siderbarSdk.GroupUsersAddEmployeesRequest(ctx.Wrap(c), *reqChat)
	detail.CustomerMember.CustomerMemberList = members.List
	detail.CustomerMember.MemberCount = members.Total

	if detail.UnionID != "" {
		// crm信息对接
		userReq := work.UserRequest{
			EnterpriseID: req.EnterpriseID,
			UnionId:      detail.UnionID,
			EuId:         detail.UserID,
		}
		userTagsReq := microScrm.UserTagsRequest{
			EnterpriseID:       req.EnterpriseID,
			WorkExternalUserID: detail.WorkExternalUserID,
			MemberID:           req.MemberID,
		}
		userImpressionTagsReq := microScrm.UserTagsRequest{
			EnterpriseID:       req.EnterpriseID,
			WorkExternalUserID: detail.UserID, // enterprise_user——id
			MemberID:           req.MemberID,
			Userid:             req.MemberUserID,
		}

		group := threading.NewRoutineGroup()

		// 用户会员信息
		group.RunSafe(func() {
			userName, errs := work.GetCrmUserName(ctx.Wrap(c), &userReq)
			if errs != nil {
				ginplus.ResError(c, errs)
				return
			}
			detail.CustomerCode = userName.CustomerCode
			detail.Membership = userName.Membership
			detail.CustomerType = userName.CustomerType
			detail.Name = userName.Name
			detail.BirthdayStr = userName.Birthday
			detail.Mobile = userName.Mobile
			detail.ExpiryDate = userName.ExpiryDate
			detail.UpgradeDateTime = userName.UpgradeDateTime
			detail.UpgradeNeedPoints = userName.UpgradeNeedPoints
			detail.Last12mPurchaseAmount = userName.Last12mPurchaseAmount
			detail.Last36mPurchaseAmount = userName.Last36mPurchaseAmount
			detail.ClientType = userName.ClientType
			detail.TotalValidPoints = userName.TotalValidPoints
			detail.Address = userName.Address

			if err != nil {
				log.Info(ctx.Wrap(c)).Msgf("vipInfo_Mobile缓存信息失败, UnionID:%s, err:%s", detail.UnionID, err.Error())
			}

		})

		// 客户标签
		group.RunSafe(func() {
			userTags, errs1 := h.customerSdk.GetUserTags(ctx.Wrap(c), &userTagsReq)
			if errs1 != nil {
				ginplus.ResError(c, errs1)
				return
			}
			detail.UserTags = userTags
		})

		// 客户印象标签
		group.RunSafe(func() {
			userImpressionTags, errs2 := h.customerSdk.GetUserImpressionTags(ctx.Wrap(c), &userImpressionTagsReq)
			if errs2 != nil {
				ginplus.ResError(c, errs2)
				return
			}
			detail.UserImpressionTags = userImpressionTags
		})

		// 客户地址、鞋码、生日
		group.RunSafe(func() {
			infoAdd, errs2 := h.sdkCustomerInfo.GetInfoAdd(ctx.Wrap(c), &userReq)
			if errs2 != nil {
				ginplus.ResError(c, errs2)
				return
			}

			detail.UserInfoAdd = infoAdd
		})
		group.Wait()
	}

	// detailStr := json.MarshalToString(detail)

	// 缓存5分钟
	// redis.Client().Set(c, cacheKey, detailStr, 5*60*time.Second)

	h.ResSuccess(c, detail)
	return
}

// CrmOrderList 客户订单
func (h CustomerHandle) CrmOrderList(c *gin.Context) {
	req := &microScrm.CustomerCrmOrderRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}

	now := time.Now()
	// last12MonthFirstDay := now.AddDate(0, -12, -now.Day()+1)
	// last12MonthStart := time.Date(last12MonthFirstDay.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).UnixNano() / 1e6
	timeNow := now.UnixNano() / 1e6 // 当前毫秒时间戳
	last12MonthStart := 0

	// CRM信息对接
	userReq := work.CrmUserRequest{
		EnterpriseID:  h.EnterpriseID(c),
		StartDateTime: last12MonthStart,
		EndDateTime:   int(timeNow),
		UnionId:       req.UnionID,
	}

	userOrderList, err := work.CrmGetUserOrderList(ctx.Wrap(c), &userReq)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("CRM订单列表接口数据错误：%+v", err)
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, userOrderList)

	return
}

// FootPrint 客户足迹
func (h CustomerHandle) FootPrint(c *gin.Context) {
	req := microScrm.CustomerFootPrintRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}

	req.EnterpriseID = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberID = memberInfo.ID

	list, err := h.customerSdk.FootPrintList(ctx.Wrap(c), req)

	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

// CrmPointsList 客户积分
func (h CustomerHandle) CrmPointsList(c *gin.Context) {
	req := &microScrm.CustomerCrmPointsRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	// Crm信息对接
	userReq := work.CrmPointsListRequest{
		EnterpriseID: h.EnterpriseID(c),
		UnionId:      req.UnionID,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		PageNum:      req.Page,
		PageSize:     req.Size,
	}

	userPointList, err := work.CrmGetUserPointList(ctx.Wrap(c), &userReq)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("CRM客户积分明细列表接口数据错误：%+v", err)
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, userPointList)

	return

}

// CrmCouponsList 卡券列表
func (h CustomerHandle) CrmCouponsList(c *gin.Context) {
	req := &microScrm.CustomerCrmCouponsRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}

	// Crm信息对接
	userReq := work.CustomerCrmCouponsRequest{
		EnterpriseID: h.EnterpriseID(c),
		UnionId:      req.UnionID,
		ChannelType:  req.ChannelType,
		Status:       1, // 优惠券状态 1未使用 2已使用 3过期 0 all
		Source:       1,
		PageNum:      req.Page,
		PageSize:     req.Size,
	}
	req.EnterpriseID = h.EnterpriseID(c)
	userCouponList, err := work.CrmGetUserCouponList(ctx.Wrap(c), &userReq)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("CRM客户优惠卷列表接口数据错误：%+v", err)
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, userCouponList)

	return
}

// GroupChatList 加入的社群
func (h CustomerHandle) GroupChatList(c *gin.Context) {
	req := &microScrm.CustomerRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = h.EnterpriseID(c)
	response, err := h.customerSdk.BelongToChat(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	memberInfo := h.MemberInfo(c)
	loginMemberID := memberInfo.ID
	// 判断当前的登录人是否在群里
	if response.Total > 0 {
		for k, item := range response.ListData {
			for _, subItem := range item.GroupMember {
				if subItem.MemberID == loginMemberID {
					response.ListData[k].MemberIsJoin = 1
					break
				}
			}
		}
	}
	h.ResSuccess(c, response)

	return
}

// UpdateRemark 修改备注
func (h CustomerHandle) UpdateRemark(c *gin.Context) {
	req := &microScrm.UpdateRemarkRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}

	req.EnterpriseID = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberID = memberInfo.ID
	req.Mobile = memberInfo.Mobile

	// 删除客户画像缓存
	cacheKey := fmt.Sprintf("USER:CRM:%d:info:%s", req.EnterpriseID, req.ExternalUserID)
	redis.Client().Del(c, cacheKey)

	list, err := h.customerSdk.UpdateRemarkForRelation(ctx.Wrap(c), req)

	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

// CrmUserDetail h5-客户画像
func (h CustomerHandle) CrmUserDetail(c *gin.Context) {
	req := &microScrm.CustomerDetailRequest{}

	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.Mobile = memberInfo.Mobile
	req.MemberType = memberInfo.Type
	req.MemberID = memberInfo.ID
	req.MemberUserID = memberInfo.Userid

	if req.MemberType == 4 {
		thirdIds, err := common.NewBaseHandle().GetMemberThirdsByType(ctx.Wrap(c), req.EnterpriseID, req.Mobile)
		if err != nil {
			h.ResError(c, err)
			return
		}
		req.ThirdIds = thirdIds
	}

	cacheKey := fmt.Sprintf("USER:CRM:%d:info:%s", req.EnterpriseID, req.ExternalUserID)
	var result string
	if err := redis.Client().Get(c, cacheKey).Scan(&result); err == nil {
		var newResult microScrm.CrmUserDetail
		_ = json.Unmarshal([]byte(result), &newResult)
		h.ResSuccess(c, newResult)
		return
	}

	data, err := h.customerApiSdk.CrmUserDetail(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	// 断言
	detail := data.(microScrm.CrmUserDetail)
	if detail.UnionID != "" {
		// crm信息对接
		userReq := work.UserRequest{
			EnterpriseID: req.EnterpriseID,
			UnionId:      detail.UnionID,
			EuId:         detail.UserID,
		}

		userTagsReq := microScrm.UserTagsRequest{
			EnterpriseID:       req.EnterpriseID,
			WorkExternalUserID: detail.WorkExternalUserID,
			MemberID:           req.MemberID,
		}
		userImpressionTagsReq := microScrm.UserTagsRequest{
			EnterpriseID:       req.EnterpriseID,
			WorkExternalUserID: detail.WorkExternalUserID,
			MemberID:           req.MemberID,
			Userid:             req.MemberUserID,
		}

		group := threading.NewRoutineGroup()
		// 用户会员信息
		group.RunSafe(func() {
			userName, errs := work.GetCrmUserName(ctx.Wrap(c), &userReq)
			if errs != nil {
				ginplus.ResError(c, errs)
				return
			}
			detail.CustomerCode = userName.CustomerCode
			detail.CustomerCard = userName.CustomerCard
			detail.CustomerCardStr = userName.CustomerCardStr
			detail.Membership = userName.Membership
			detail.Name = userName.Name
			detail.BirthdayStr = userName.Birthday
			detail.Mobile = userName.Mobile
			detail.ExpiryDate = userName.ExpiryDate
			detail.UpgradeDateTime = userName.UpgradeDateTime
			detail.UpgradeNeedPoints = userName.UpgradeNeedPoints
			detail.Last12mPurchaseAmount = userName.Last12mPurchaseAmount
			detail.Last36mPurchaseAmount = userName.Last36mPurchaseAmount
			detail.TotalValidPoints = userName.TotalValidPoints
			detail.ClientType = userName.ClientType
			detail.Address = userName.Address

			detail.VipStatus = microScrm.CrmVipStatus{
				IssueDate:    userName.VipStatus.IssueDate,
				IssueShop:    userName.VipStatus.IssueShop,
				PtsToUpgrade: userName.VipStatus.PtsToUpgrade,
				ExpiryDate:   userName.VipStatus.ExpiryDate,
			}

			if err != nil {
				log.Info(ctx.Wrap(c)).Msgf("vipInfo_Mobile缓存信息失败, UnionID:%s, err:%s", detail.UnionID, err.Error())
			}

		})

		// 客户标签
		group.RunSafe(func() {
			userTags, errs1 := h.customerSdk.GetUserTags(ctx.Wrap(c), &userTagsReq)
			if errs1 != nil {
				ginplus.ResError(c, errs1)
				return
			}
			detail.UserTags = userTags
		})

		// 客户印象标签
		group.RunSafe(func() {
			userImpressionTags, errs2 := h.customerSdk.GetUserImpressionTags(ctx.Wrap(c), &userImpressionTagsReq)
			if errs2 != nil {
				ginplus.ResError(c, errs2)
				return
			}
			detail.UserImpressionTags = userImpressionTags
		})

		// 客户地址、鞋码、生日
		group.RunSafe(func() {
			infoAdd, errs2 := h.sdkCustomerInfo.GetInfoAdd(ctx.Wrap(c), &userReq)
			if errs2 != nil {
				ginplus.ResError(c, errs2)
				return
			}

			detail.UserInfoAdd = infoAdd
		})

		group.Wait()
	}

	// 缓存30秒
	detailStr := json.MarshalToString(detail)
	redis.Client().Set(c, cacheKey, detailStr, 30*time.Second)

	h.ResSuccess(c, detail)
	return
}

// UserDetailTipsNum h5-获取顾客邀约和播种数量
func (h CustomerHandle) UserDetailTipsNum(c *gin.Context) {
	req := &microScrm.CustomerDetailRequest{}

	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberID = memberInfo.ID
	req.MemberUserID = memberInfo.Userid
	req.Mobile = memberInfo.Mobile
	req.MemberType = memberInfo.Type

	if req.MemberType == 4 {
		thirdIds, err := common.NewBaseHandle().GetMemberThirdsByType(ctx.Wrap(c), req.EnterpriseID, req.Mobile)
		if err != nil {
			h.ResError(c, err)
			return
		}
		req.ThirdIds = thirdIds
	}

	data, err := h.customerApiSdk.GetUserDetailTipsNum(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, data)
	return
}

// UserBindRelation 顾客和导购的绑定关系
func (h CustomerHandle) UserBindRelation(c *gin.Context) {
	req := &microScrm.CustomerDetailRequest{}

	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberID = memberInfo.ID
	req.MemberUserID = memberInfo.Userid
	req.Mobile = memberInfo.Mobile
	req.MemberType = memberInfo.Type

	if req.MemberType == 4 {
		thirdIds, err := common.NewBaseHandle().GetMemberThirdsByType(ctx.Wrap(c), req.EnterpriseID, req.Mobile)
		if err != nil {
			h.ResError(c, err)
			return
		}
		req.ThirdIds = thirdIds
	}

	data, err := h.customerApiSdk.GetUserBindRelation(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, data)
	return
}

// SwUserCrmCouponsList  SW卡券列表
func (h CustomerHandle) SwUserCrmCouponsList(c *gin.Context) {
	req := &work.SwCrmCouponsRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	// Crm信息对接
	userReq := work.SwCrmCouponsRequest{
		EnterpriseID: h.EnterpriseID(c),
		UnionId:      req.UnionId,
	}
	req.EnterpriseID = h.EnterpriseID(c)
	userCouponList, err := work.SwUserCrmGetUserCouponList(ctx.Wrap(c), &userReq)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResJson(c, http.StatusOK, userCouponList)

	return
}
