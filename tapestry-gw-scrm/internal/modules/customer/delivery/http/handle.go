package customer_http

import (
	"strconv"
	"strings"
	"time"

	scrm_errs "gw-scrm/internal/pkg/errs"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gw-scrm/internal/modules/customer"
	"gw-scrm/internal/modules/customer/api"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/errorIgnore"
	response "gw-scrm/internal/pkg/resoponse"
	"gw-scrm/internal/pkg/validate_params"
	"gw-scrm/pkg/micro-scrm-sdk/microScrm"
)

func updateDelMsgDetail(c *gin.Context) {
	params := &customer.UpdateMessageRequest{}

	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseID = tokenInfo.EnterpriseId
	act := &api.CustomerApi{}
	res := act.UpdateDelMsgDetail(ctx.Wrap(c), params)
	ginplus.ResSuccess(c, res)
}

func customerDelMsgList(c *gin.Context) {
	params := &customer.MessageListRequest{}

	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseID = tokenInfo.EnterpriseId
	act := &api.CustomerApi{}
	res := act.MessageList(ctx.Wrap(c), params)
	ginplus.ResSuccess(c, res)
}

// 重复客户
func repeatCustomer(c *gin.Context) {
	params := &customer.RepeatListRequest{}

	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseID = tokenInfo.EnterpriseId
	act := &api.CustomerApi{}
	res := act.RepeatCustomer(ctx.Wrap(c), params)
	ginplus.ResSuccess(c, res)
}

// RepeatFindMember
func repeatFindMember(c *gin.Context) {
	params := &customer.ExternalRelationByUserIdRequest{}

	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseID = tokenInfo.EnterpriseId
	act := &api.CustomerApi{}
	res := act.RepeatFindMember(ctx.Wrap(c), params)
	ginplus.ResSuccess(c, res)
}

// lossCustomersList 流失客户列表
func lossCustomersList(c *gin.Context) {
	params := &customer.LossCustomersListRequest{}

	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseID = tokenInfo.EnterpriseId
	act := &api.CustomerApi{}
	res := act.LossCustomersList(ctx.Wrap(c), params)
	ginplus.ResSuccess(c, res)
}

const (
	// 后台管理查询信息时，不需要传递member_id. 用改常量判断是否来自哪个端。RequestCustomer中添加ChannelType来判断
	ChannelSuperAdmin = 1
)

// footprintList 客户足迹列表
func footprintList(c *gin.Context) {
	params := &customer.RequestCustomer{}
	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	reqParams := &microScrm.RequestCustomer{
		UserID:       params.UserID,
		Page:         params.Page,
		Size:         params.Size,
		EnterpriseID: getEnterpriseId(c),
		ChannelType:  ChannelSuperAdmin,
	}
	total, res, err := customerSDK.FootprintList(ctx.Wrap(c), reqParams)

	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, scrm_errs.FootprintListErr, err)
		return
	}

	response.SuccessList(c, total, reqParams.Page, reqParams.Size, res)
	return
}

// orderList 客户订单列表
func orderList(c *gin.Context) {
	params := &customer.RequestCustomer{}
	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	reqParams := &microScrm.RequestCustomer{
		UserID:       params.UserID,
		Page:         params.Page,
		Size:         params.Size,
		SortType:     params.SortType,
		SortKey:      params.SortKey,
		EnterpriseID: getEnterpriseId(c),
	}
	total, res, err := customerSDK.OrderList(ctx.Wrap(c), reqParams)

	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, scrm_errs.OrderListErr, err)
		return
	}
	response.SuccessList(c, total, reqParams.Page, reqParams.Size, res)
	return
}

// orderExport 客户订单列表
func orderExport(c *gin.Context) {
	params := &customer.RequestCustomer{}
	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}
	reqParams := &microScrm.RequestCustomer{
		UserID:       params.UserID,
		Size:         1000,
		Page:         1,
		EnterpriseID: getEnterpriseId(c),
	}
	_, res, err := customerSDK.OrderList(ctx.Wrap(c), reqParams)

	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, scrm_errs.OrderListErr, err)
		return
	}
	excelFile := excelize.NewFile()
	sheetName := "客户订单"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	// 写入表头
	_ = excelFile.SetCellStr(sheetName, "A1", "订单编号")
	_ = excelFile.SetCellStr(sheetName, "B1", "下单账户")
	_ = excelFile.SetCellStr(sheetName, "C1", "下单商品")
	_ = excelFile.SetCellStr(sheetName, "D1", "优惠券")
	_ = excelFile.SetCellStr(sheetName, "E1", "含券支付金额")
	_ = excelFile.SetCellStr(sheetName, "F1", "实付金额")
	_ = excelFile.SetCellStr(sheetName, "G1", "支付时间")
	_ = excelFile.SetCellStr(sheetName, "H1", "退款金额")
	_ = excelFile.SetCellStr(sheetName, "I1", "退款时间")
	startRow := 2
	cstSh, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("LoadLocation:%s", err.Error())
		cstSh = time.Local
	}
	for _, item := range *res {
		rowStartStr := strconv.Itoa(startRow)
		var goodsName []string
		for _, itemGoods := range item.GoodsList {
			goodsName = append(goodsName, itemGoods.SpuName)
		}
		goodsNameStr := strings.Join(goodsName, ";")
		var couponName []string
		for _, couponItem := range item.CouponList {
			couponName = append(couponName, couponItem.Name+"（"+couponItem.StockId+"）")
		}
		couponNameStr := strings.Join(couponName, ";")
		payAt := time.Unix(item.PayAt, 0).In(cstSh).Format("2006-01-02 15:04")
		var refundAt string
		if item.RefundAt > 0 {
			refundAt = time.Unix(item.RefundAt, 0).In(cstSh).Format("2006-01-02 15:04")
		}

		_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.OrderNo)
		_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, item.UserInfo.Name)
		_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, goodsNameStr)
		_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, couponNameStr)
		_ = excelFile.SetCellValue(sheetName, "E"+rowStartStr, strconv.FormatFloat(float64(item.FinalAmount)/100, 'f', 2, 32))
		_ = excelFile.SetCellValue(sheetName, "F"+rowStartStr, strconv.FormatFloat(float64(item.CashFee)/100, 'f', 2, 32))
		_ = excelFile.SetCellStr(sheetName, "G"+rowStartStr, payAt)
		_ = excelFile.SetCellValue(sheetName, "H"+rowStartStr, strconv.FormatFloat(float64(item.TotalRefundFee)/100, 'f', 2, 32))
		_ = excelFile.SetCellStr(sheetName, "I"+rowStartStr, refundAt)
		startRow = startRow + 1
	}
	excelFile.SetActiveSheet(index)
	fileName := "客户订单" + time.Now().In(cstSh).Format("200601021504") + ".xlsx"
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	des, _ := excelFile.WriteToBuffer()
	_, _ = c.Writer.Write(des.Bytes())
	return
}

// belongToMembers 所属员工列表
func belongToMembers(c *gin.Context) {
	params := &customer.RequestCustomer{}
	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	reqParams := &microScrm.RequestCustomer{
		UserID:       params.UserID,
		Page:         params.Page,
		Size:         params.Size,
		EnterpriseID: getEnterpriseId(c),
	}
	total, res, err := customerSDK.BelongToMembers(ctx.Wrap(c), reqParams)

	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, scrm_errs.BelongToMembersErr, err)
		return
	}

	response.SuccessList(c, total, reqParams.Page, reqParams.Size, res)
	return
}

// belongToChats 所属群列表
func belongToChats(c *gin.Context) {
	params := &customer.RequestCustomerGroupChat{}
	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	params.Status = 1
	params.EnterpriseID = getEnterpriseId(c)
	act := &api.CustomerApi{}
	res, err := act.BelongToChats(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

// consigneeAddressList 收货地址列表
func consigneeAddressList(c *gin.Context) {
	params := &customer.RequestCustomer{}
	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	reqParams := &microScrm.RequestCustomer{
		UserID: params.UserID,
		Page:   params.Page,
		Size:   params.Size,
	}
	total, res, err := customerSDK.ConsigneeAddressList(ctx.Wrap(c), reqParams)

	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, scrm_errs.ConsigneeAddressListErr, err)
		return
	}

	response.SuccessList(c, total, reqParams.Page, reqParams.Size, res)
	return
}

func customerDetail(c *gin.Context) {
	params := &customer.RequestCustomer{}
	ok := validate_params.ValidateParams(c, params)
	if !ok {
		return
	}

	reqParams := &microScrm.RequestCustomer{
		UserID:       params.UserID,
		Page:         params.Page,
		Size:         params.Size,
		EnterpriseID: getEnterpriseId(c),
	}
	res, err := customerSDK.CustomerDetail(ctx.Wrap(c), reqParams)

	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, scrm_errs.CustomerDetailErr, err)
		return
	}
	var result *customer.ResponseCustomerDetail
	if res != nil {
		result = &customer.ResponseCustomerDetail{
			UserID:       res.UserID,
			UserName:     res.UserName,
			Nickname:     res.Nickname,
			Avatar:       res.Avatar,
			Gender:       res.Gender,
			Country:      res.Country,
			Province:     res.Province,
			City:         res.City,
			RemarkList:   res.RemarkList,
			TagList:      res.TagList,
			GroupList:    res.GroupList,
			IsWorkUser:   res.IsWorkUser,
			OfficialList: translateOfficialList(res.OfficialList),
			ShowOfficial: res.ShowOfficial,
			LevelName:    res.LevelName,
			Birthday:     res.Birthday,
			UnionCode:    res.UnionCode,
			Point:        res.Point,
			Grade:        res.Grade,
			StoreCode:    res.StoreCode,
			BACode:       res.BACode,
			Mobile:       res.Mobile,
		}
		reqAddressParams := &microScrm.RequestCustomer{
			UserID:       params.UserID,
			Page:         0,
			Size:         1000,
			EnterpriseID: getEnterpriseId(c),
		}

		_, list, err := customerSDK.ConsigneeAddressList(ctx.Wrap(c), reqAddressParams)

		if err == nil {
			result.AddressList = list
		}
	}

	response.Success(c, result)
	return
}

func translateOfficialList(accountList []*microScrm.DetailOfficialAccount) []*customer.DetailOfficialAccount {
	var list = make([]*customer.DetailOfficialAccount, 0)
	for _, v := range accountList {
		list = append(list, &customer.DetailOfficialAccount{
			UserID:        v.UserID,
			Name:          v.Name,
			SubscribeTime: v.SubscribeTime,
		})
	}
	return list
}

func customerList(c *gin.Context) {
	var request interface{}
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}

	act := &api.CustomerApi{}
	res, err := act.CustomerList(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

func editCustomerRemark(c *gin.Context) {
	var request interface{}
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
	}

	act := &api.CustomerApi{}
	res, err := act.EditCustomerRemark(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

func tagList(c *gin.Context) {
	var request interface{}
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
	}

	act := &api.CustomerApi{}
	res, err := act.TagList(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

func listConfig(c *gin.Context) {
	var request interface{}
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
	}

	act := &api.CustomerApi{}
	res, err := act.ListConfig(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

func orderGoods(c *gin.Context) {
	var request interface{}
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
	}

	act := &api.CustomerApi{}
	res, err := act.OrderGoods(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

// getEnterpriseId 获取企业ID
func getEnterpriseId(c *gin.Context) int64 {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	if tokenInfo == nil {
		return 0
	}

	return tokenInfo.EnterpriseId
}

// 会员等级列表
func memberLevelList(c *gin.Context) {
	ginplus.ResSuccess(c, gin.H{})
	//request := common.GetRawMap(c)
	//var tokenInfo *common.TokenInfo
	//if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
	//	tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	//}
	//request["enterprise_id"] = tokenInfo.EnterpriseId
	//act := &api.CustomerApi{}
	//res, err := act.MemberLevelList(ctx.Wrap(c), request)
	//if err != nil {
	//	ginplus.ResError(c, err)
	//	return
	//}
	//ginplus.ResSuccess(c, res)
}

func SyncEterpriseUserToES(c *gin.Context) {
	var request interface{}
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
	}

	act := &api.CustomerApi{}
	res, err := act.SyncEterpriseUserToES(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}
