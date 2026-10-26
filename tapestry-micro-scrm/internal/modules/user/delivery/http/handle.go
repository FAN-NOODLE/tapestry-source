package userhttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/modules/marketing"
	"micro-scrm/internal/modules/user"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/errs"
)

var userService *user.Service
var orderService *user.OrderService
var enterpriseService *enterprise.Service
var marketService *marketing.Service
var customerService *customer.Service

// AddEvent 添加用户事件
func AddEvent(c *gin.Context) {
	var request user.SyncUserEventRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	var userEvent string
	switch request.Method {
	case common.DdcEventEnum.UserJoinActivity.String():
		userEvent = string(common.EventEnum.UserJoinActivity)
	default:
		userEvent = string(request.Method)
	}
	event := request.Event
	userEntity, err := customerService.QueryEnterpriseUser(ctx.Wrap(c), event.EnterpriseID, &customer.SourceUserBase{UnionID: event.UnionID})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	err = userService.AddUserEvent(ctx.Wrap(c), event.EnterpriseID, userEntity.ID, userEvent, event.EventDesc, event.Object, event.Extra, int64(event.CreatedAt))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

func SyncOrder(c *gin.Context) {
	var request user.OrderSyncRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	err = orderService.SyncOrder(ctx.Wrap(c), &request.OrderData, userService, request.Method, marketService)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// customerList 客户列表
func customerList(c *gin.Context) {
	request := customer.ListRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	list, err := customerService.QueryList(ctx.Wrap(c), &request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
}

// tagGroupList 客户分组下预计触达人数统计
// 该接口暂时下线，因SQL对CPU营销较大。
func tagGroupUserCount(c *gin.Context) {
	var (
		request marketing.UserCategory
		err     error
		res     *customer.TagUserCount
	)
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	cond := &customer.ListCond{EnterpriseID: request.EnterpriseID}
	if request.TagOrGroupID > 0 {
		cond.CustomerGroupCond.Include = true
		cond.CustomerGroupCond.GroupIDList = []int64{request.TagOrGroupID}
	}

	for _, channel := range request.ChannelTypeList {
		switch channel {
		case marketing.StrategyChannelOfficial:
			cond.CustomerSource = append(cond.CustomerSource, customer.CustomerTypeOpen)
		case marketing.StrategyChannelQW:
			cond.CustomerSource = append(cond.CustomerSource, customer.CustomerTypeExternal)
		case marketing.StrategyChannelShop:
			cond.CustomerSource = append(cond.CustomerSource, customer.CustomerTypeShop)
		case marketing.StrategyChannelMobile:
			cond.CustomerSource = append(cond.CustomerSource, customer.CustomerTypeMobile)
		}
	}

	request.Type = "group"
	res, err = userService.QueryUserCount(ctx.Wrap(c), cond)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

func editRemarks(c *gin.Context) {
	var request user.EditRemarksRequest
	err := c.BindJSON(&request)
	request.EnterpriseID = common.EnterpriseID(ctx.Wrap(c))
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	res, err := userService.EditRemarks(ctx.Wrap(c), request, enterpriseService)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

// listConfig 客户列表展示配置
func listConfig(c *gin.Context) {
	res, err := userService.ListConfig(ctx.Wrap(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}
