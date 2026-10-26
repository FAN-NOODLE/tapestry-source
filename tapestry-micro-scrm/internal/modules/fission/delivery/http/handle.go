package http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gorm.io/gorm"
	logic_common "micro-scrm/internal/logic/common"
	"micro-scrm/internal/modules/fission"
	"micro-scrm/internal/orm"
	micro_scrm_fission "micro-scrm/internal/orm/model/micro_scrm/fission"
	fission2 "micro-scrm/internal/orm/service/fission"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/errs"
	response "micro-scrm/internal/pkg/resoponse"
	"sync"
	"time"
)

var fissionService *fission.Service

type RequestFissionUser struct {
	UserID     int64 `json:"user_id"`
	ActivityID int64 `json:"activity_id"`
}

type RequestActivityID struct {
	ActivityID int64 `json:"activity_id"`
}

// exchangePrizeList 裂变奖品兑换列表
func exchangePrizeList(c *gin.Context) {
	var request struct {
		common.RequestPage
		RequestFissionUser
		EnterpriseID int64 `json:"enterprise_id"`
	}
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	var response = fission.ExchangeList{
		PageType: common.PageType{
			Page: request.Page,
			Size: request.Size,
		},
		List: []*fission.ExchangeInfo{},
	}
	if request.UserID == 0 || request.ActivityID == 0 {
		ginplus.ResSuccess(c, response)
		return
	}
	response.Total, response.List, err = fissionService.ExchangeList(ctx.Wrap(c), request.EnterpriseID, request.ActivityID, request.UserID, request.Page, request.Size)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// activityList 获取裂变活动列表
func activityList(c *gin.Context) {
	var request fission.ActivityListRequest
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	total, list, err := fissionService.ActivityList(ctx.Wrap(c), &request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	resp := fission.ActivityListResponse{}
	if list != nil {
		resp.List = list
	} else {
		resp.List = []*fission.ActivityEntity{}
	}
	resp.Total = total
	resp.Page = request.Page
	resp.Size = request.Size

	ginplus.ResSuccess(c, resp)
	return
}

// exchangePrizeDetail 裂变奖品兑换详情
func exchangePrizeDetail(c *gin.Context) {
	var request common.RequestID
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	detail, err := fissionService.ExchangeDetail(ctx.Wrap(c), request.ID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, detail)
}

func prizeList(c *gin.Context) {
	var request struct {
		RequestActivityID
		common.RequestPage
		EnterpriseID int64 `json:"enterprise_id"`
	}
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	list, err := fissionService.PrizeList(ctx.Wrap(c), request.EnterpriseID, request.ActivityID, request.Page, request.Size)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
	return
}

func prizeInfo(c *gin.Context) {
	var request common.RequestID
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	prizeInfo, err := fissionService.PrizeInfoWithDeleted(ctx.Wrap(c), request.ID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	activityInfo, err := fissionService.ActivityDetail(ctx.Wrap(c), prizeInfo.EnterpriseID, prizeInfo.ActivityID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if activityInfo == nil {
		ginplus.ResError(c, errs.FissionActivityNotExistErr)
		return
	}

	prizeDetail := fission.PrizeDetail{
		PrizeEntity:  *prizeInfo,
		ActivityDesc: activityInfo.Desc,
	}
	now := time.Now().Unix()
	if prizeInfo.DeletedAt != 0 {
		prizeDetail.Status = fission.PrizeStatusDeleted
	} else if activityInfo.Status == fission.StatusNotStarted {
		prizeDetail.Status = fission.PrizeStatusNotStart
	} else if prizeInfo.Stocks <= prizeInfo.ExchangeNum {
		prizeDetail.Status = fission.PrizeStatusEmpty
	} else if now > activityInfo.ExchangeEndTime {
		prizeDetail.Status = fission.PrizeStatusExchangeEnd
	}
	ginplus.ResSuccess(c, prizeDetail)
	return
}

func assistantList(c *gin.Context) {
	var request struct {
		RequestFissionUser
		common.RequestPage
		EnterpriseID int64 `json:"enterprise_id"`
	}
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	list, err := fissionService.AssistantList(ctx.Wrap(c), request.EnterpriseID, request.ActivityID, request.UserID, request.Page, request.Size)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
}

func exchangePrize(c *gin.Context) {
	var request fission.ExchangePrizeRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	_, err = fissionService.ExchangePrize(ctx.Wrap(c), request.EnterpriseID, request.ActivityID, request.UserID, request.PrizeID, &request.AddressInfo)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

func exchangePrizeCheck(c *gin.Context) {
	var request fission.UserPrizeRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	_, err = fissionService.ExchangePrizeCheck(ctx.Wrap(c), request.EnterpriseID, request.ActivityID, request.UserID, request.PrizeID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	record, err := fissionService.QueryUserLatestExchangeRecord(ctx.Wrap(c), request.UserID)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("获取用户最近兑换记录失败,err:%s", err.Error())
	}
	var userAddress fission.AddressInfo

	if record != nil {
		userAddress.DeliverName, _ = common.AesDecrypt(record.DeliverName)
		userAddress.DeliverMobile, _ = common.AesDecrypt(record.DeliverMobile)
		userAddress.Address, _ = common.AesDecrypt(record.DeliverAddress)
		userAddress.Area, _ = common.AesDecrypt(record.DeliverArea)
		userAddress.AreaCode, _ = common.AesDecrypt(record.DeliverAreaCode)
	}
	ginplus.ResSuccess(c, struct {
		UserAddress fission.AddressInfo `json:"user_address"`
	}{
		UserAddress: userAddress,
	})
}

func scrollList(c *gin.Context) {
	var request fission.ScrollListRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	total, res, err := fissionService.GetScrollList(ctx.Wrap(c), request.EnterpriseID, request.ActivityID, request.Page, request.Size)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	ginplus.ResSuccess(c, response.List{
		Total: total,
		Page:  int(request.Page),
		Size:  int(request.Size),
		List:  res,
	})
}

func joinActivity(c *gin.Context) {
	var request struct {
		RequestFissionUser
		EnterpriseID int64 `json:"enterprise_id"`
	}
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	fissionUser, err := fissionService.JoinActivity(ctx.Wrap(c), request.EnterpriseID, request.ActivityID, request.UserID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var response struct {
		QrCode string `json:"qr_code"`
	}
	response.QrCode, _, _ = fissionService.GetContactQrCode(ctx.Wrap(c), fissionUser.ContactID, request.EnterpriseID)
	ginplus.ResSuccess(c, response)
}

func mpHome(c *gin.Context) {
	var request struct {
		EnterpriseID int64 `json:"enterprise_id"`
		RequestFissionUser
	}
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	activityInfo, err := fissionService.ActivityDetailWithDeleted(ctx.Wrap(c), request.EnterpriseID, request.ActivityID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if activityInfo == nil {
		ginplus.ResError(c, errs.FissionActivityNotExistErr)
		return
	}
	var homeInfo fission.MpHomeInfo
	homeInfo.ActivityInfo.Name = activityInfo.Name
	homeInfo.ActivityInfo.FissionType = activityInfo.FissionType
	homeInfo.ActivityInfo.StartTime = activityInfo.StartTime
	homeInfo.ActivityInfo.EndTime = activityInfo.EndTime
	homeInfo.ActivityInfo.ExchangeEndTime = activityInfo.ExchangeEndTime
	homeInfo.ActivityInfo.Desc = activityInfo.Desc
	homeInfo.ActivityInfo.ActivityStatus = activityInfo.Status
	homeInfo.ActivityInfo.SeverTime = time.Now().Unix()
	homeInfo.ActivityInfo.MemberQrCode, _, _ = fissionService.GetContactQrCode(ctx.Wrap(c), activityInfo.ContactId, activityInfo.EnterpriseID)
	homeInfo.ActivityInfo.ChatQrCode, _ = fissionService.ByActivityContactQrCode(ctx.Wrap(c), request.ActivityID, request.EnterpriseID)
	style, _ := fissionService.GetStyleConfig(request.EnterpriseID)
	homeInfo.ActivityInfo.Styles = style.Styles
	homeInfo.ActivityInfo.RgbColors = style.RgbColors
	bg, _ := fissionService.GetBgConfig(request.EnterpriseID)
	homeInfo.ActivityInfo.BgImg = bg.BgImg
	homeInfo.ActivityInfo.ShareImg = bg.ShareImg
	// 查询用户昵称和头像
	mpUserInfo, err := fissionService.GetMpUser(ctx.Wrap(c), request.UserID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if mpUserInfo != nil {
		homeInfo.UserInfo.ImgUrl = mpUserInfo.ImgUrl
		homeInfo.UserInfo.UnionID = mpUserInfo.UnionId
	}

	// 如果裂变活动是：加BA且VIP的，则同步裂变客户的vip状态，在裂变信息表中 fission_assistant_records.is_vip = 1
	if activityInfo.FissionType == micro_scrm_fission.FA_FissionType_VIP {
		// 获取该客户在该裂变活动中的裂变粉丝信息
		dataList, err := orm.Service.Fission.GetFissionAssistantUserOneFieldList(activityInfo.EnterpriseID, activityInfo.ID, request.UserID)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		if len(dataList) != 0 {
			recordIDList := []int64{}
			var waitGroup sync.WaitGroup
			for _, data := range dataList {
				waitGroup.Add(1)
				go func(data fission2.UpdateVipFlag) {
					defer waitGroup.Done()
					if isVip, _ := logic_common.NewVipUser(ctx.Wrap(c)).SetUser(logic_common.VipUserParam{
						EnterpriseID: activityInfo.EnterpriseID,
						UnionID:      data.UnionID,
					}).IsVip(); isVip {
						recordIDList = append(recordIDList, data.ID)
					}
				}(data)
			}
			waitGroup.Wait()

			// 更新会员标志 fission_assistant_records.is_vip
			if len(recordIDList) != 0 {
				orm.Service.Fission.Transaction(func(tx *gorm.DB) error {
					if err := orm.Service.Fission.UpdateFissionIsVipField(activityInfo.EnterpriseID, recordIDList); err != nil {
						return err
					}
					if err := orm.Service.Fission.UpdateFissionUserInfo(activityInfo.EnterpriseID, activityInfo.ID, request.UserID, len(recordIDList)); err != nil {
						return err
					}

					return nil
				})
			}
		}
	}

	// 目前的场景没有使用商城的东西，故注释
	//authInfo, err := fissionService.GetEnterpriseAuthInfo(ctx.Wrap(c), request.EnterpriseID)
	//if err != nil {
	//	ginplus.ResError(c, err)
	//	return
	//}
	//if authInfo != nil {
	//	homeInfo.ShopInfo = *authInfo
	//}
	if activityInfo.Status == fission.StatusDeleted {
		ginplus.ResSuccess(c, homeInfo)
		return
	}

	activityUser, err := fissionService.GetActivityUser(ctx.Wrap(c), request.EnterpriseID, request.ActivityID, request.UserID)
	if err != nil && err != errs.FissionActivityUserNotExist {
		ginplus.ResError(c, err)
		return
	}
	// 检查用户是否参加其他互斥活动
	if activityInfo.MutualType == fission.IsMutualType {
		homeInfo.UserInfo.JoinOtherActivity, err = fissionService.IsUserInMutualActivity(ctx.Wrap(c), request.EnterpriseID, request.ActivityID, request.UserID)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
	}
	if activityUser != nil {
		homeInfo.UserInfo.CurAssistantNum = activityUser.CurAssistantNum
		homeInfo.UserInfo.TotalAssistantNum = activityUser.TotalAssistantNum
		homeInfo.UserInfo.ContactQrCode, _, _ = fissionService.GetContactQrCode(ctx.Wrap(c), activityUser.ContactID, activityUser.EnterpriseID)
	}

	ginplus.ResSuccess(c, homeInfo)
}

//activityDetail 活动详情
func activityDetail(c *gin.Context) {
	var (
		request  fission.ActivityDetailRequest
		activity *fission.ActivityEntity
		err      error
	)
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	activity, err = fissionService.ActivityDetail(ctx.Wrap(c), request.EnterpriseID, request.ActivityID)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, activity)
	return
}

//activityDetailList活动详情列表
func activityDetailList(c *gin.Context) {
	var (
		err      error
		total    int64
		request  fission.ActivityDetailRequest
		list     *[]fission.ActivityUserEntity
		response fission.ActivityDetailListResponse
	)
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	total, list, err = fissionService.ActivityDetailList(ctx.Wrap(c), request)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	response.List = list
	response.Total = total
	response.Page = request.Page
	response.Size = request.Size
	ginplus.ResSuccess(c, response)
	return
}

//inviteDetail发起人邀请明细
func inviteDetail(c *gin.Context) {
	var (
		request  fission.InviteDetailRequest
		response fission.InviteDetailResponse
		list     []*fission.AssistantRecordInfo
		err      error
		total    int64
	)
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	total, list, err = fissionService.InviteDetail(ctx.Wrap(c), request)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	response.List = list
	response.Total = total
	response.Page = request.Page
	response.Size = request.Size
	ginplus.ResSuccess(c, response)
	return
}

//exchangeRecord发起人兑换记录
func exchangeRecord(c *gin.Context) {
	var (
		request  fission.ExchangeDetailRequest
		response fission.ExchangeDetailResponse
		list     *[]fission.ExchangeRecordInfo
		err      error
		total    int64
	)
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	total, list, err = fissionService.ExchangeRecord(ctx.Wrap(c), request)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	response.List = list
	response.Total = total
	response.Page = request.Page
	response.Size = request.Size
	ginplus.ResSuccess(c, response)
	return

}

//prizeExchangeRecord 兑奖记录
func prizeExchangeRecord(c *gin.Context) {

	var (
		request  fission.ExchangeDetailRequest
		response fission.ExchangeDetailResponse
		list     *[]fission.ExchangeRecordInfo
		err      error
		total    int64
	)
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	total, list, err = fissionService.ExchangeRecord(ctx.Wrap(c), request)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	response.List = list
	response.Total = total
	response.Page = request.Page
	response.Size = request.Size
	ginplus.ResSuccess(c, response)
	return
}

//openExchangeWarn 兑换提醒设置
func openExchangeWarn(c *gin.Context) {
	var (
		request fission.ExchangeWarnRequest
		err     error
	)
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	err = fissionService.OpenExchangeWarn(ctx.Wrap(c), request)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nil)
	return
}

//ExchangeCount 奖品兑换记录总数
func ExchangeTotal(c *gin.Context) {
	var (
		request fission.ExchangeDetailRequest
		err     error
		total   int64
		result  = make(map[string]interface{})
	)
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	total, err = fissionService.ExchangeTotal(ctx.Wrap(c), request)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	result["total"] = total
	ginplus.ResSuccess(c, result)
	return
}

//activityListSimple 活动列表 全部活动
func activityListSimple(c *gin.Context) {
	var (
		request fission.ActivityListRequest
		list    = make([]*fission.ActivityEntity, 0)
		err     error
	)
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	list, err = fissionService.ActivityListSimple(ctx.Wrap(c), request)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
	return
}

//activitySet 裂变活动-设置
func activitySet(c *gin.Context) {
	var request fission.ActivitySetRequest
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if request.Activity.ExchangeEndTime < request.Activity.EndTime {
		ginplus.ResError(c, errs.ExchangeEndTimeErr)
		return
	}
	if request.Activity.ID == 0 && request.Activity.StartTime < time.Now().Unix() {
		ginplus.ResError(c, errs.ExchangeStartTimeErr)
		return
	}
	activityID, err := fissionService.ActivitySet(ctx.Wrap(c), request.EnterpriseID, request.Activity, request.MemberIDs, request.ChatGroups, request.WelcomeContent, request.Prize, request.MutualActivityIDs)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, activityID)
	return
}

//getExchangeConfigMember 获取兑换提醒的额外员工
func getExchangeConfigMember(c *gin.Context) {
	var (
		request  fission.ExchangeWarnRequest
		err      error
		response *fission.ExchangeWarnReponse
	)
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	response, err = fissionService.GetExchangeConfigMember(ctx.Wrap(c), request.EnterpriseID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
	return
}

//styleConfigList 获取配色列表
func styleConfigList(c *gin.Context) {
	var request fission.StyleConfigRequest
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	list, err := fissionService.GetStyleConfigList(request.EnterpriseID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	bgList, err := fissionService.GetBgConfigList(request.EnterpriseID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	data := make(map[string]interface{})
	data["style-list"] = list
	data["bg-list"] = bgList
	ginplus.ResSuccess(c, data)
	return
}

//addStyleConfig 获取配色列表
func addStyleConfig(c *gin.Context) {
	var request fission.StyleConfigRequest
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	model := micro_scrm_fission.FissionStyleConfigModel{
		EnterpriseID: request.EnterpriseID,
		RgbColors:    request.RgbColors,
		Styles:       request.Styles,
	}
	err := fissionService.AddStyleConfig(model)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nil)
	return
}

//delStyleConfig 获取配色列表
func delStyleConfig(c *gin.Context) {
	var request fission.DelStyleConfigRequest
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	err := fissionService.DelStyleConfig(request.EnterpriseID, request.ID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nil)
	return
}

//setStyleConfig 设置选中配色
func setStyleConfig(c *gin.Context) {
	var request fission.StyleConfigRequest
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	//设置配色
	err := fissionService.SetStyleConfig(request.EnterpriseID, request.ID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	//设置背景图
	err = fissionService.SetBgConfig(request.EnterpriseID, request.BgImg, request.ShareImg)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nil)
	return
}

//getCheckStyle 获取当前配色
func getCheckStyle(c *gin.Context) {
	var request fission.StyleConfigRequest
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	list, err := fissionService.GetStyleConfig(request.EnterpriseID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
	return
}

// mutualActivityList 选择互斥活动
func mutualActivityList(c *gin.Context) {
	var request fission.MutualActivityListRequest
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	request.Action = "list"
	total, list, err := fissionService.MutualActivityList(ctx.Wrap(c), &request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	resp := fission.MutualActivityListResponse{}
	if list != nil {
		resp.List = list
	} else {
		resp.List = []*fission.MutualActivityEntity{}
	}
	resp.Total = total
	resp.Page = request.Page
	resp.Size = request.Size
	ginplus.ResSuccess(c, resp)
	return
}

// chatGroupList 选择社群列表
func chatGroupList(c *gin.Context) {
	var request fission.ChatGroupListRequest
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	total, list, err := fissionService.ChatGroupList(ctx.Wrap(c), &request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	resp := fission.ChatGroupListListResponse{}
	if list != nil {
		resp.List = list
	} else {
		resp.List = []*fission.WorkChatGroupEntity{}
	}
	resp.Total = total
	resp.Page = request.Page
	resp.Size = request.Size
	ginplus.ResSuccess(c, resp)
	return
}

// activityDelete 删除裂变活动
func activityDelete(c *gin.Context) {
	var request fission.ActivityInfoRequest
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	err := fissionService.ActivityDelete(ctx.Wrap(c), request.EnterpriseID, request.ActivityID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
	return
}

// activityInfo 裂变活动详情
func activityInfo(c *gin.Context) {
	var request fission.ActivityInfoRequest
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	info, err := fissionService.ActivityInfo(ctx.Wrap(c), request.ActivityID, request.EnterpriseID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, info)
	return
}

// updateStatus 更新活动状态
func updateStatus(c *gin.Context) {
	var request fission.UpdateStatusRequest
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	err := fissionService.UpdateStatus(ctx.Wrap(c), request.EnterpriseID, request.ActivityID, request.Status)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
	return
}

// fissionActivityList 企微裂变活动列表
func fissionActivityList(c *gin.Context) {
	var request fission.QwActivityListRequest
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}

	request.NowTime = time.Now().Unix()
	total, list, listErr := fissionService.QwActivityList(ctx.Wrap(c), &request)
	if listErr != nil {
		ginplus.ResError(c, listErr)
		return
	}

	resp := fission.QwActivityListResponse{}
	if len(*list) > 0 {
		resp.List = list
	} else {
		resp.List = &[]fission.QwActivityList{}
	}

	info, _ := fissionService.StatisticsInfo(ctx.Wrap(c), &request)

	resp.Total = total
	resp.Page = request.Page
	resp.Size = request.Size

	if info != nil {
		resp.Statistics.UserNum = info.UserNum
		resp.Statistics.FissionNum = info.FissionNum
		resp.Statistics.RunAwayNum = info.RunAwayNum
		resp.Statistics.JoinGroupNum = info.JoinGroupNum
		resp.Statistics.LeaveGroupNum = info.LeaveGroupNum
	}
	ginplus.ResSuccess(c, resp)
	return
}
func sendMessage(c *gin.Context) {
	var (
		request fission.ExchangeWarnRequest
		err     error
	)
	err = c.BindJSON(&request)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	err = fissionService.SendMessage(ctx.Wrap(c), request.EnterpriseID, 132822132281728, "裂变活动0113", "张三")
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
}

// mpUserInfo 用户信息
func mpUserInfo(c *gin.Context) {
	req := new(common.RequestID)
	err := c.ShouldBindJSON(req)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	user, err := fissionService.GetMpUser(ctx.Wrap(c), req.ID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, user)
}
