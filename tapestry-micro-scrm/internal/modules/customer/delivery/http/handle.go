// Package callback_http
package customerHttp

import (
	"fmt"
	"micro-scrm/internal/pkg/helper"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/json"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/customer/delivery"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/modules/members"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/errs"
	response "micro-scrm/internal/pkg/resoponse"
	"micro-scrm/internal/pkg/uid"
)

var s *customer.Service
var m *customer.MessageService
var enterpriseService *enterprise.Service
var memberService *members.Service

// footprintList 客户足迹列表
func footprintList(c *gin.Context) {
	params := &customer.RequestCustomer{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	total, list, err := s.FootprintList(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errs.FootprintListErr, err)
		return
	}

	response.SuccessList(c, total, params.Page, params.Size, list)
	return
}

// orderDetail 订单详情
func orderDetail(c *gin.Context) {
	params := &customer.RequestOrderDetail{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	response.Success(c, customer.ResponseOrderDetail{})

	return
}

// orderList 客户订单列表
func orderList(c *gin.Context) {
	params := &customer.RequestOrderList{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	if params.SortKey != "" {
		if !customer.OrderSortKeyList.Contains(params.SortKey) {
			params.SortKey = ""
		}
	}
	params.SortType = strings.ToUpper(params.SortType)
	if params.SortType != "ASC" && params.SortType != "DESC" {
		params.SortType = "DESC"
	}

	response.SuccessList(c, 0, params.Page, params.Size, customer.ResponseOrderList{})

	return
}

// orderGoods 客户订单商品列表
func orderGoods(c *gin.Context) {
	params := &customer.OrderGoodsRequest{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	response.SuccessList(c, 0, params.Page, params.Size, customer.GoodsList{})
	return
}

// customerDetail 客户详情
func customerDetail(c *gin.Context) {
	params := &customer.RequestCustomer{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	detail, err := s.CustomerDetail(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errs.CustomerDetailErr, err)
		return
	}

	if detail.UserID == 0 {
		response.Success(c, nil)
		return
	}

	response.Success(c, detail)
	return
}

// belongToMembers 所属员工列表
func belongToMembers(c *gin.Context) {
	params := &customer.RequestCustomer{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	total, list, err := s.BelongToMembers(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errs.BelongToMembersErr, err)
		return
	}

	response.SuccessList(c, total, params.Page, params.Size, list)
	return
}

// belongToChats 所属群列表
func belongToChats(c *gin.Context) {
	params := &customer.RequestCustomerGroup{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	total, list, err := s.BelongToChats(ctx.Wrap(c), params)

	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errs.BelongToChatErr, err)
		return
	}

	response.SuccessList(c, total, params.Page, params.Size, list)
	return
}

// consigneeAddressList 收货地址列表
func consigneeAddressList(c *gin.Context) {
	params := &customer.RequestCustomer{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	response.SuccessList(c, 0, params.Page, params.Size, &customer.ResponseAddressList{})
}

// lossCustomersList 流失客户列表
func lossCustomersList(c *gin.Context) {
	params := &customer.LossCustomersListRequest{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	if params.MemberID > 0 {
		member, err := memberService.Detail(ctx.Wrap(c), params.MemberID)
		if err != nil {
			return
		}
		params.UserId = append(params.UserId, member.Userid)
	} else {
		if params.DepartmentID > 0 {
			userIdList, err := memberService.GetUserIdsByDepId(ctx.Wrap(c), params.DepartmentID, params.EnterpriseID)
			if err != nil {
				return
			}
			if len(userIdList) > 0 {
				params.UserId = append(params.UserId, userIdList...)
			}
		}
	}

	total, list, err := s.LossCustomersList(ctx.Wrap(c), *params)
	if err != nil {
		return
	}
	response.SuccessList(c, total, params.Page, params.Size, list)

	return
}

func addCustomerDelMessage(c *gin.Context) {
	params := &customer.AddMessageRequest{}

	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	cc := ctx.Wrap(c)
	params.ID = uid.ID(cc)
	id, errAdd := m.AddMessage(cc, customer.TranAddParam(params))

	if errAdd != nil {
		ginplus.ResError(c, errAdd)

		return
	}

	responseData := customer.AddMessageResponse{ID: id}
	ginplus.ResSuccess(c, responseData)
}

func customerDelMsgDetail(c *gin.Context) {
	params := &customer.MessageDetailRequest{}

	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	res, _ := m.MessageDetail(ctx.Wrap(c), params.ID)
	ginplus.ResSuccess(c, res)
}

func updateDelMsgDetail(c *gin.Context) {
	params := &customer.UpdateMessageRequest{}

	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	if params.ID == 0 {
		cc := ctx.Wrap(c)
		params.ID = uid.ID(cc)

		_, errAdd := m.AddMessage(cc, customer.TranUpdateParam(params))
		if errAdd != nil {
			ginplus.ResError(c, errAdd)
			return
		}
		ginplus.ResSuccess(c, []struct{}{})
		return
	} else {
		_, UpErr := m.UpdateMessage(ctx.Wrap(c), params)
		if UpErr != nil {
			ginplus.ResError(c, UpErr)
			return
		}

		ginplus.ResSuccess(c, []struct{}{})
		return
	}
}

func sendMessage(c *gin.Context) {
	params := &customer.SendMessageRequest{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	cc := ctx.Wrap(c)

	corpReq := enterprise.Corp{
		EnterpriseID: params.EnterpriseID,
	}
	// 查询企业信息
	enterpriseData, enterpriseErr := enterpriseService.Detail(ctx.Wrap(c), corpReq)
	if enterpriseErr != nil {
		ginplus.ResError(c, enterpriseErr)
		return
	}
	_, errRes := m.SendMessage(cc, params, enterpriseData)
	if errRes != nil {
		ginplus.ResError(c, errRes)
		return
	}

	ginplus.ResSuccess(c, []struct{}{})
}

func customerDelMsgList(c *gin.Context) {
	params := &customer.MessageListRequest{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	cc := ctx.Wrap(c)

	list, errRes := m.MessageList(cc, params)
	if errRes != nil {
		ginplus.ResError(c, errRes)
		return
	}
	statObj := &customer.MessageMember{}
	statObjNew := &customer.MessageMemberArr{}
	for _, sendData := range *list {
		_ = json.Unmarshal([]byte(sendData.QyStatTxt), statObj)
		memCount := len(statObj.Member)
		if memCount > 0 {
			reqMembers := members.MemberListRequest{
				UserIDs:      statObj.Member,
				UserIDCount:  memCount,
				Page:         1,
				Size:         memCount,
				EnterpriseID: params.EnterpriseID,
			}
			total, memberData := memberService.Page(cc, reqMembers)
			if total > 0 {
				for _, v := range *memberData {
					value := customer.StatResponse{
						v.Name,
						v.Userid,
					}
					statObjNew.Member = append(statObjNew.Member, value)
				}
				sendData.StatTxtObj = statObjNew
			}
		}
		if sendData.StatTxtObj == nil {
			sendData.StatTxtObj = &customer.MessageMemberArr{
				Member:     []customer.StatResponse{},
				Department: []customer.StatResponse{},
			}
		}
		ginplus.ResSuccess(c, sendData)
		return
	}

	ginplus.ResSuccess(c, struct{}{})
}

func repeatCustomer(c *gin.Context) {
	params := &customer.RepeatListRequest{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	cc := ctx.Wrap(c)

	total, list, errRes := s.RepeatCustomerList(cc, delivery.TransListParams(params))
	if errRes != nil {
		ginplus.ResError(c, errRes)
		return
	}

	if total == 0 {
		response.SuccessList(c, total, params.Page, params.Size, []struct{}{})
		return
	}
	response.SuccessList(c, total, params.Page, params.Size, list)
}

// FindExternalRelationByUserId 外部联系人添加的重复成员列表
func FindExternalRelationByUserId(c *gin.Context) {
	params := &customer.ExternalRelationByUserIdRequest{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	cc := ctx.Wrap(c)

	total, list, errRes := s.FindExternalRelationByUserId(cc, params)
	if errRes != nil {
		ginplus.ResError(c, errRes)
		return
	}
	response.SuccessList(c, total, params.Page, params.Size, list)
}

// portraitDetail 客户画像详情
func portraitDetail(c *gin.Context) {
	req := customer.PortraitRequest{}
	if err := c.ShouldBind(&req); err != nil {
		ginplus.ResError(c, err)
		return
	}
	response, err := s.PortraitDetail(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// remarkSave 客户备注信息保存
func remarkSave(c *gin.Context) {
	params := &customer.RemarkInfo{}

	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	cc := ctx.Wrap(c)

	userID, err := s.RemarkSave(cc, params)
	if err != nil {
		return
	}

	response.Success(c, userID)
	return
}

// SetExternalTags 设置外部联系人企业标签并同步到企微
func SetExternalTags(c *gin.Context) {
	params := &customer.SetExternalTags{}

	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	cc := ctx.Wrap(c)

	_, err := s.SetExternalTags(cc, params)
	if err != nil {
		return
	}

	response.Success(c, struct{}{})
	return
}

func SyncEterpriseUserToES(c *gin.Context) {
	params := &customer.SyncEterpriseUserToEsRequest{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	s.SyncEterpriseUserToES(ctx.Wrap(c), params)
	response.Success(c, struct{}{})
	return
}

// lostOption 丢失客户记录，额外数据
func lostOption(c *gin.Context) {
	var request struct {
		EnterpriseId int64  `json:"enterprise_id"  binding:"required" comment:"企业Id"`
		Time         int64  `json:"time"  binding:"required" comment:"时间戳"`
		DeleteType   uint   `json:"delete_type" binding:"min=1,max=2" comment:"类型"`
		Mobile       string `json:"mobile" binding:"required" comment:"员工手机号"`
	}
	err := c.Bind(&request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	// 通过导购员工手机号查询员工信息
	var member micro_scrm.MemberModel
	gormx.GetInstance().Where("mobile", request.Mobile).First(&member)

	var customersDelUserMessage []micro_scrm.CustomersDelUserMessageModel
	// 查询用户对应的组织架构配置信息
	gormx.GetInstance().Select("dp_id,send_type").
		Where("enterprise_id", request.EnterpriseId).
		Where("user_id", member.Userid).
		Find(&customersDelUserMessage)

	var sendTypeDpIDs1, sendTypeDpIDs2 []int64
	for _, cus := range customersDelUserMessage {
		if cus.SendType == 0 || cus.SendType == 1 {
			sendTypeDpIDs1 = append(sendTypeDpIDs1, cus.DpID)
		}
		if cus.SendType == 0 || cus.SendType == 2 {
			sendTypeDpIDs2 = append(sendTypeDpIDs2, cus.DpID)
		}
	}

	// 查询导购被动流失好友统计
	var deleteFriendRelationRecords []struct {
		ExternalUserID uint64 `json:"external_user_id"`
		UserId         string `json:"user_id"`
		MemberName     string `json:"member_name"`
		DpId           int64  `json:"dp_id"`
		DpName         string `json:"dp_name"`
		DeleteType     uint32 `json:"delete_type"`
		Mobile         string `json:"mobile"`
	}

	gormx.GetInstance().Table("delete_friend_relation_record as dfrr").
		Select("dfrr.external_user_id,delete_type,dfrr.user_id,dfrr.dp_id,members.name as member_name,departments.name as dp_name,members.mobile").
		Joins("left join members on members.userid = dfrr.user_id").
		Joins("left join departments on departments.id = dfrr.dp_id").
		Where("dfrr.dp_id IN (?)", func() (dp_id []string) {
			dpErr := gormx.GetInstance().Model(&micro_scrm.CustomersDelUserMessageModel{}).Where("user_id = ?", member.Userid).Where("(send_type = ? or send_type = 0)", request.DeleteType).Pluck("dp_id", &dp_id).Error
			if dpErr != nil {
				return []string{}
			}
			return
		}()).
		Where("dfrr.enterprise_id=? and dfrr.created_at >= ? and dfrr.created_at < ?", request.EnterpriseId, request.Time, request.Time+86400-1).
		Where("(dfrr.delete_type = ? and dfrr.dp_id in (?)) or (dfrr.delete_type = ? and dfrr.dp_id in (?))", request.DeleteType, sendTypeDpIDs1, request.DeleteType, sendTypeDpIDs2).
		Find(&deleteFriendRelationRecords)

	var departmentInfo []map[string]interface{}

	var departmentIds []int64
	var userIds []string

	var externalUserIds []uint64
	var lostTotal int64

	// 部门和员工数据联动
	departmentMemberListMap := make(map[int64][]interface{})

	for _, deleteFriend := range deleteFriendRelationRecords {
		if !helper.InArray[string](deleteFriend.UserId, userIds) {
			userIds = append(userIds, deleteFriend.UserId)
			if _, ok := departmentMemberListMap[deleteFriend.DpId]; !ok {
				departmentMemberListMap[deleteFriend.DpId] = []interface{}{}
			}
			departmentMemberListMap[deleteFriend.DpId] = append(departmentMemberListMap[deleteFriend.DpId], gin.H{"user_id": deleteFriend.UserId, "username": deleteFriend.MemberName, "mobile": deleteFriend.Mobile})
		}
		if !helper.InArray[int64](deleteFriend.DpId, departmentIds) {
			departmentIds = append(departmentIds, deleteFriend.DpId)
			departmentInfo = append(departmentInfo, gin.H{"dp_id": deleteFriend.DpId, "dp_name": deleteFriend.DpName})
		}
		if !helper.InArray[uint64](deleteFriend.ExternalUserID, externalUserIds) && deleteFriend.DeleteType == uint32(request.DeleteType) {
			externalUserIds = append(externalUserIds, deleteFriend.ExternalUserID)
			lostTotal++
		}
	}

	ginplus.ResSuccess(c, gin.H{
		"lost_total":  lostTotal,
		"departments": departmentInfo,
		"member":      departmentMemberListMap,
	})

	return
}

// lostRecord 丢失客户记录
func lostRecord(c *gin.Context) {
	var request customer.LostRecordRequest
	err := c.Bind(&request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	// 通过导购员工手机号查询员工信息
	var member micro_scrm.MemberModel
	gormx.GetInstance().Where("mobile", request.Mobile).First(&member)

	// 部门人数
	var lostUsers []*micro_scrm.DeleteFriendRelationRecordData
	query := gormx.GetInstance().Table("delete_friend_relation_record as dfrr").
		Select("dfrr.external_user_id,any_value(dfrr.user_id) as user_id,any_value(dfrr.dp_name) as dp_name,any_value(dfrr.dp_id) as dp_id,any_value(dfrr.created_at) as created_at,any_value(eu.name) as name,any_value(eu.nick) as nick,any_value(eu.avatar) as avatar,any_value(eu.birthday) as birthday,any_value(eu.mobile) as mobile,any_value(eu.sex) as sex").
		Joins(fmt.Sprintf("LEFT JOIN micro_scrm_enterprise_%d.wx_user_corp_external as wuce on wuce.id=dfrr.external_user_id", request.EnterpriseId)).
		Joins(fmt.Sprintf("LEFT JOIN micro_scrm_enterprise_%d.enterprise_user as eu on eu.id=wuce.euid", request.EnterpriseId)).
		Where("dfrr.enterprise_id=? and dfrr.created_at >= ? and dfrr.created_at < ?", request.EnterpriseId, request.Time, request.Time+86400-1).
		Where("dfrr.delete_type", request.DeleteType).
		Where("dfrr.dp_id IN (?)", func() (dp_id []string) {
			dpErr := gormx.GetInstance().Model(&micro_scrm.CustomersDelUserMessageModel{}).Where("user_id = ?", member.Userid).Where("(send_type = ? or send_type = 0)", request.DeleteType).Pluck("dp_id", &dp_id).Error
			if dpErr != nil {
				return []string{}
			}
			return
		}()).
		Group("dfrr.external_user_id")

	if request.UserId != "" {
		query.Where("dfrr.user_id", request.UserId)
	}

	if request.DpId != 0 {
		query.Where("dfrr.dp_id", request.DpId)
	}

	var lostCount int64
	query.Count(&lostCount)

	page := request.Page
	size := request.Size
	query.Offset((page - 1) * size).Limit(size).Find(&lostUsers)
	for _, user := range lostUsers {
		if user.Mobile != "" {
			user.Mobile = user.Mobile[:3] + "****" + user.Mobile[7:]
		}
	}

	ginplus.ResSuccess(c, gin.H{
		"lost_total": 0,
		"lost_users": lostUsers,
		"lost_count": lostCount,
	})

	return
}
