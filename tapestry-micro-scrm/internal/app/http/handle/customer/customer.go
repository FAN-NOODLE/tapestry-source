package customer

import (
	"errors"
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	customerReq "micro-scrm/internal/app/http/request/customer"
	workExternalModel "micro-scrm/internal/app/models/micro_scrm"
	cacheSvc "micro-scrm/internal/app/service/cache"
	customerSvc "micro-scrm/internal/app/service/customer"
	memberSvc "micro-scrm/internal/app/service/member"
	"micro-scrm/internal/pkg/common"
)

type CustomerHandle struct {
	handle.BaseHandle
}

func NewCustomerHandle() CustomerHandle {
	return CustomerHandle{}
}

func (h CustomerHandle) UserList(c *gin.Context) {
	var request customerReq.UserList
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	request.EnterpriseID = common.EnterpriseID(cc)

	result, total, err := customerSvc.NewCustomer().UserList(cc, request)
	if err != nil {
		return
	}

	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResList(c, total, request.Page, request.Size, result)

	return
}

func (h CustomerHandle) UserDetail(c *gin.Context) {
	var request struct {
		UserID       int64 `json:"user_id" binding:"required" required:"客户ID必填"`
		EnterpriseID int64 `json:"enterprise_id"`
	}
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	request.EnterpriseID = common.EnterpriseID(h.WarpContext(c))

	cc := h.WarpContext(c)
	result, err := customerSvc.NewCustomer().UserDetail(cc, request.EnterpriseID, request.UserID)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)

	return
}

func (h CustomerHandle) UpdateRemark(c *gin.Context) {
	var request customerReq.UpdateRemarkRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	if request.MemberID == 0 {
		memberInfo, err := memberSvc.NewMembersService().GetInfoByMobile(h.WarpContext(c), request.EnterpriseID, request.Mobile)
		if err != nil {
			h.ResError(c, err)
			return
		}
		request.MemberID = memberInfo.ID
	}

	if err := customerSvc.NewCustomer().UpdateRemark(h.WarpContext(c), request); err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)

	return
}

// UserSimpleInfo 获取简单的用户信息
func (h CustomerHandle) UserSimpleInfo(c *gin.Context) {
	var request customerReq.SimpleInfoRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	UserInfo, err := customerSvc.NewCustomer().SimpleInfo(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	var result customerReq.SimpleInfoRes
	result.Avatar = UserInfo.Avatar
	result.Name = UserInfo.Name
	result.UnionId = UserInfo.Unionid
	result.ExternalUserID = UserInfo.Openid
	result.Nick = UserInfo.Nick
	h.ResSuccess(c, result)

	return
}

// SetInfo 同步unionid的微信openid
func (h CustomerHandle) SetInfo(c *gin.Context) {
	var request customerReq.OpenIdInfoReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if err := customerSvc.NewCustomer().SetInfo(h.WarpContext(c), request); err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)

	return
}

// MemberLists 获取外部联系人有关的导购手机号
func (h CustomerHandle) MemberLists(c *gin.Context) {
	var request customerReq.MemberListsReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	result, err := customerSvc.NewCustomer().GetMemberListsByUnionId(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
	return
}

// MyCenter 顾客端--我的
func (h CustomerHandle) MyCenter(c *gin.Context) {
	var request customerReq.UserInfoRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	userInfo, err := customerSvc.NewCustomer().GetInfoByUnionId(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	var result customerReq.SimpleInfoRes
	result.Avatar = userInfo.Avatar
	result.Name = userInfo.Name
	result.UnionId = userInfo.Unionid
	result.ExternalUserID = userInfo.Openid
	result.Nick = userInfo.Nick
	h.ResSuccess(c, result)

	return
}

// MySa 我的-添加专属顾问
func (h CustomerHandle) MySa(c *gin.Context) {
	var request customerReq.UserSaRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	userInfo, err := customerSvc.NewCustomer().GetInfoByUnionId(h.WarpContext(c), request.UserInfoRequest)
	if err != nil {
		h.ResError(c, err)
		return
	}
	request.WorkExternalUserID = userInfo.ID
	info, err := customerSvc.NewCustomer().GetUserSAByExternalUserIds(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	var userSa = workExternalModel.UserSaData{}
	if len(info) == 0 {
		var data []int64
		h.ResSuccess(c, data)
		return
	}
	userSa = workExternalModel.UserSaData{
		MemberId:   info[0].MemberId,
		Unionid:    info[0].Unionid,
		Userid:     info[0].Userid,
		QrCode:     info[0].QrCode,
		IsKeySa:    info[0].IsKeySa,
		SaUserId:   info[0].SaUserId,
		SaUserName: info[0].SaUserName,
	}

	h.ResSuccess(c, userSa)

	return
}

// UserDetailTipsNum 导购或店长获取顾客邀约和播种数量
func (h CustomerHandle) UserDetailTipsNum(c *gin.Context) {
	var request customerReq.UserDetailTipsRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	if request.MemberType == 4 && len(request.ThirdIds) == 0 {
		h.ResError(c, errors.New("third_ids不能为空"))
		return
	}
	memberInfo, err := memberSvc.NewMembersService().GetInfoByMobile(h.WarpContext(c), request.EnterpriseID, request.Mobile)
	if err != nil {
		h.ResError(c, err)
		return
	}
	request.MemberID = memberInfo.ID
	if request.MemberType == 3 {
		request.ThirdIds = []int{}
	}

	// 查询邀约和播种数量
	result, err := customerSvc.NewCustomer().GetUserDetailTips(h.WarpContext(c), request, memberInfo.ID, memberInfo.Userid)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)

	return
}

// UserBindRelation 顾客和导购的绑定关系
func (h CustomerHandle) UserBindRelation(c *gin.Context) {
	var request customerReq.UserBindRelationRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	memberInfo, err := memberSvc.NewMembersService().GetInfoByMobile(h.WarpContext(c), request.EnterpriseID, request.Mobile)
	if err != nil {
		h.ResError(c, err)
		return
	}
	request.MemberId = memberInfo.ID

	resConfig, err := cacheSvc.NewRedisData().GetEnterpriseConfig(h.WarpContext(c), request.EnterpriseID, "store_name")
	if err != nil {
		h.ResError(c, err)
	}
	result, err := customerSvc.NewCustomer().GetUserSAListByExternalUserIds(h.WarpContext(c), request, resConfig.Name)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)

	return
}

func (h CustomerHandle) UpdateDescription(c *gin.Context) {
	var request customerReq.UpdateDescriptionRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if err := customerSvc.NewCustomer().UpdateDescription(h.WarpContext(c), request); err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)

	return
}
