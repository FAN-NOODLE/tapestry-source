package customerhttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"gw-work-scrm/internal/pkg/redis"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
	"gw-work-scrm/internal/pkg/sdk/work"
	"gw-work-scrm/internal/pkg/validate_params"
	"math"
	"strconv"
	"sync"
	"time"
)

var customerSdk *microScrm.CustomerApi
var siderbarSdk *microScrm.SidebarGroupChatApi

// footPrint 足迹
func footPrint(c *gin.Context) {
	req := &microScrm.CustomerFootPrintRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	req.MemberID = c.Value("memberId").(int64)
	list, err := customerSdk.FootPrintList(ctx.Wrap(c), *req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
}

// detail 详情
func detail(c *gin.Context) {
	req := &microScrm.CustomerDetailRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	req.MemberID = c.Value("memberId").(int64)
	detail, err := customerSdk.Detail(ctx.Wrap(c), req)
	fmt.Println(detail)
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
	members, _ := siderbarSdk.GroupUsersAddEmployeesRequest(ctx.Wrap(c), *reqChat)
	detail.CustomerMember.CustomerMemberList = members.List
	detail.CustomerMember.MemberCount = members.Total
	if detail.UnionID != "" {
		// crm信息对接
		userReq := work.UserRequest{
			EnterpriseID: req.EnterpriseID,
			//UnionId:      detail.UnionID,
			UnionId: "o_GI00elT3HiZ72K4N-k8tLshed4",
		}
		eID := strconv.FormatInt(req.EnterpriseID, 10)
		//用户级别
		//group.RunSafe(func() {
		//userGrade, err := work.GetUserGrade(ctx.Wrap(c), &userReq)
		//if err == nil {
		//	detail.Grade = userGrade.Grade
		//	detail.GradeCode = userGrade.GradeCode
		//	//detail.EndTime = userGrade.GradeEndDate
		//	//detail.UPTime = userGrade.GradeStartDate
		//	err = redis.Client().HSet(c, eID+":"+detail.UnionID, "Grade", userGrade.Grade).Err()
		//	err = redis.Client().HSet(c, eID+":"+detail.UnionID, "GradeCode", userGrade.GradeCode).Err()
		//	err = redis.Client().HSet(c, eID+":"+detail.UnionID, "EndTime", userGrade.GradeEndDate).Err()
		//	err = redis.Client().HSet(c, eID+":"+detail.UnionID, "UPTime", userGrade.GradeStartDate).Err()
		//	if err != nil {
		//		log.Info(ctx.Wrap(c)).Msgf("Grade缓存信息失败, UnionID:%s, err:%s", detail.UnionID, err.Error())
		//	}
		//}
		//})
		//所有会员等级
		//group.RunSafe(func() {
		//dlData, err := work.GetAllGrade(ctx.Wrap(c), &userReq)
		//if err == nil {
		//	currentGradeCode := detail.GradeCode
		//	for k, v := range dlData {
		//		if k != 0 && v.TierCode == currentGradeCode {
		//			detail.NextGrade = dlData[k-1].TierName
		//			detail.NextGradeCode = dlData[k-1].TierCode
		//			err = redis.Client().HSet(c, eID+":"+detail.UnionID, "NextGrade", dlData[k-1].TierName).Err()
		//			err = redis.Client().HSet(c, eID+":"+detail.UnionID, "NextGradeCode", dlData[k-1].TierCode).Err()
		//			break
		//		}
		//	}
		//	if err != nil {
		//		log.Info(ctx.Wrap(c)).Msgf("NextGrade缓存信息失败, UnionID:%s, err:%s", detail.UnionID, err.Error())
		//	}
		//}
		//})
		//用户升级信息
		//group.RunSafe(func() {
		//gcReq := work.UserRequest{
		//	EnterpriseID: req.EnterpriseID,
		//	//Type:         "union_code",
		//	//Value:        "158S0002K9L",
		//	Type:    "unionId",
		//	Value:   detail.UnionID,
		//	Grade:   detail.GradeCode,
		//	UpGrade: detail.NextGradeCode,
		//}
		//dlData1, err := work.GetUserGradeCalculator(ctx.Wrap(c), &gcReq)
		//if err == nil {
		//	upTex := "再消费￥" + fmt.Sprintf("%.2f", dlData1.RequiredAmount) + "升级至" + detail.NextGrade
		//	if detail.NextGrade == "" {
		//		upTex = ""
		//	}
		//	requiredAmount := dlData1.RequiredAmount
		//	currentAmount := dlData1.CurrentAmount
		//	var progress int64
		//	if requiredAmount == currentAmount+requiredAmount {
		//		progress = 0
		//	} else {
		//		total := currentAmount + requiredAmount
		//		bili := Wrap(requiredAmount/total, 2)
		//		progress = bili
		//	}
		//	detail.UpGradeInfo.UpGradeTxt = upTex
		//	detail.UpGradeInfo.Percentage = progress
		//	err = redis.Client().HSet(c, eID+":"+detail.UnionID, "UpGradeTxt", upTex).Err()
		//	err = redis.Client().HSet(c, eID+":"+detail.UnionID, "Percentage", progress).Err()
		//	if err != nil {
		//		log.Info(ctx.Wrap(c)).Msgf("Grade缓存信息失败, UnionID:%s, err:%s", detail.UnionID, err.Error())
		//	}
		//}
		//})
		group := threading.NewRoutineGroup()
		////用户积分
		//group.RunSafe(func() {
		//	userPoint, err := work.GetUserPoint(ctx.Wrap(c), &userReq)
		//	if err == nil {
		//		//会员积分
		//		if len(userPoint) > 0 {
		//			detail.EvalidPoints = userPoint[0].EvalidPoints
		//			err = redis.Client().HSet(c, eID+":"+detail.UnionID, "EvalidPoints", userPoint[0].EvalidPoints).Err()
		//			if err != nil {
		//				log.Info(ctx.Wrap(c)).Msgf("EvalidPoints缓存信息失败, UnionID:%s, err:%s", detail.UnionID, err.Error())
		//			}
		//		}
		//	}
		//})

		//用户过期积分
		//group.RunSafe(func() {
		//	epReq := work.UserRequest{
		//		EnterpriseID: req.EnterpriseID,
		//		//Type:         "union_code",
		//		//Value:        "158S0002K9L",
		//		Type:           "unionId",
		//		Value:          detail.UnionID,
		//		PointTypeGroup: "BP",
		//		Startdate:      time.Unix(now.BeginningOfMonth().Unix(), 0).Format("20060102"),
		//		Enddate:        time.Unix(now.EndOfMonth().Unix(), 0).Format("20060102"),
		//		//Enddate: "20220930",
		//	}
		//	data, err := work.GetUserPointsExpire(ctx.Wrap(c), &epReq)
		//	if err == nil {
		//		detail.ExpiredPoints = data.ExpireNum
		//		err = redis.Client().HSet(c, eID+":"+detail.UnionID, "ExpiredPoints", data.ExpireNum).Err()
		//		if err != nil {
		//			log.Info(ctx.Wrap(c)).Msgf("Grade缓存信息失败, UnionID:%s, err:%s", detail.UnionID, err.Error())
		//		}
		//	}
		//})

		//用户会员信息
		group.RunSafe(func() {
			userName, _ := work.GetUserName(ctx.Wrap(c), &userReq)
			fmt.Println(userName)
			if err == nil {

				detail.CustomerCode = userName.CustomerCode
				detail.CustomerType = userName.CustomerType
				detail.Name = userName.Name
				detail.BirthdayStr = userName.Birthday
				detail.Mobile = userName.Mobile
				detail.ExpiryDate = userName.ExpiryDate
				detail.UpgradeDateTime = userName.UpgradeDateTime
				detail.Last36mPurchaseAmount = userName.Last36mPurchaseAmount

				err = redis.Client().HSet(c, eID+":"+detail.UnionID+":vipInfo", "CustomerCode", userName.CustomerCode).Err()
				err = redis.Client().HSet(c, eID+":"+detail.UnionID+":vipInfo", "CustomerType", userName.CustomerType).Err()
				err = redis.Client().HSet(c, eID+":"+detail.UnionID+":vipInfo", "Name", userName.Name).Err()
				err = redis.Client().HSet(c, eID+":"+detail.UnionID+":vipInfo", "BirthdayStr", userName.Birthday).Err()
				err = redis.Client().HSet(c, eID+":"+detail.UnionID+":vipInfo", "Mobile", userName.Mobile).Err()
				err = redis.Client().HSet(c, eID+":"+detail.UnionID+":vipInfo", "ExpiryDate", userName.ExpiryDate).Err()
				err = redis.Client().HSet(c, eID+":"+detail.UnionID+":vipInfo", "UpgradeDateTime", userName.UpgradeDateTime).Err()
				err = redis.Client().HSet(c, eID+":"+detail.UnionID+":vipInfo", "Last36mPurchaseAmount", userName.Last36mPurchaseAmount).Err()
				if err != nil {
					log.Info(ctx.Wrap(c)).Msgf("vipInfo_Mobile缓存信息失败, UnionID:%s, err:%s", detail.UnionID, err.Error())
				}
			}
		})

		isExists := redis.Client().Exists(c, eID+":"+detail.UnionID).Val()
		isVipExists := redis.Client().Exists(c, eID+":"+detail.UnionID+":vipInfo").Val()
		if isExists == 1 && isVipExists == 1 {

			detail.EvalidPoints, _ = redis.Client().HGet(c, eID+":"+detail.UnionID, "EvalidPoints").Int64()
			detail.ExpiredPoints, _ = redis.Client().HGet(c, eID+":"+detail.UnionID, "ExpiredPoints").Int64()
			detail.Grade = redis.Client().HGet(c, eID+":"+detail.UnionID, "Grade").Val()
			detail.GradeCode = redis.Client().HGet(c, eID+":"+detail.UnionID, "GradeCode").Val()
			detail.CustomerCode = redis.Client().HGet(c, eID+":"+detail.UnionID+":vipInfo", "CustomerCode").Val()
			detail.CustomerType = redis.Client().HGet(c, eID+":"+detail.UnionID+":vipInfo", "CustomerType").Val()
			detail.Name = redis.Client().HGet(c, eID+":"+detail.UnionID+":vipInfo", "Name").Val()
			detail.BirthdayStr = redis.Client().HGet(c, eID+":"+detail.UnionID+":vipInfo", "BirthdayStr").Val()
			detail.Mobile = redis.Client().HGet(c, eID+":"+detail.UnionID+":vipInfo", "Mobile").Val()
			detail.ExpiryDate = redis.Client().HGet(c, eID+":"+detail.UnionID, "ExpiryDate").Val()
			detail.UpgradeDateTime = redis.Client().HGet(c, eID+":"+detail.UnionID, "UpgradeDateTime").Val()
			detail.Last36mPurchaseAmount, _ = redis.Client().HGet(c, eID+":"+detail.UnionID, "Last36mPurchaseAmount").Float32()
			detail.Status = redis.Client().HGet(c, eID+":"+detail.UnionID+":vipInfo", "Status").Val()
			detail.Location = redis.Client().HGet(c, eID+":"+detail.UnionID+":vipInfo", "Location").Val()
			//detail.UpGradeInfo.UpGradeTxt = redis.Client().HGet(c, eID+":"+detail.UnionID, "UpGradeTxt").Val()
			//detail.UpGradeInfo.Percentage, _ = redis.Client().HGet(c, eID+":"+detail.UnionID, "Percentage").Int64()
		} else {
			group.Wait()
		}
	}
	ginplus.ResSuccess(c, detail)
}

//DL会员信息订单列表
func dlOrderList(c *gin.Context) {
	req := &microScrm.CustomerDLOrderGoodsRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	//DL信息对接
	userReq := work.UserRequest{
		EnterpriseID: c.Value("enterpriseId").(int64),
		Type:         "unionId",
		Value:        req.UnionID,
		Tradeno:      req.TradeNo,
		//Type:  "union_code",
		//Value: "158S0002K9L",
	}
	userOrderList, err := work.GetUserOrderList(ctx.Wrap(c), &userReq)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("dl订单列表接口数据错误：%+v", err)
	}
	ginplus.ResSuccess(c, userOrderList)
}

//DL会员积分明细列表
func dlPointsDetailList(c *gin.Context) {
	req := &microScrm.CustomerDLPointsDetailRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	dateResult := GetMonthStartAndEnd()
	//DL信息对接
	userReq := work.PointsListRequest{
		EnterpriseID: c.Value("enterpriseId").(int64),
		Type:         "unionId",
		Value:        req.UnionID,
		//Type:  "union_code",
		//Value: "158S0002K9L",
		//Startdate: time.Unix(now.BeginningOfMonth().Unix(), 0).Format("20060102"),
		//Enddate:   time.Unix(now.EndOfMonth().Unix(), 0).Format("20060102"),
		Startdate: dateResult["start"],
		Enddate:   dateResult["end"],
		PageNo:    req.PageNo,
		Limit:     req.Limit,
	}

	var list []work.PointList
	//pointTypeGroup := []string{"BP", "RP", "ON", "OP"} //, "RP", "ON", "OP"
	pointTypeGroup := []string{"BP"} //, "RP", "ON", "OP"
	var wg sync.WaitGroup
	lock := &sync.Mutex{}
	wg.Add(len(pointTypeGroup))
	//不同积分分组类型的各自明细汇总到一个切片里
	for _, v := range pointTypeGroup {
		go func(wg *sync.WaitGroup, lock *sync.Mutex, pointType string) {
			lock.Lock()
			userReq.PointTypeGroup = pointType
			dlData, err := work.GetUserPointsDetailList(ctx.Wrap(c), userReq)
			if err != nil {
				dlData = &[]work.PointList{}
			}
			for _, v1 := range *dlData {
				tm, _ := time.Parse("2006-01-02 15:04:05", v1.CreateDateTime)
				v1.SortDateTime = tm.Unix() - 28800
				list = append(list, v1)
			}
			lock.Unlock()
			wg.Done()
		}(&wg, lock, v)
	}

	wg.Wait()

	//按时间倒序
	for i := 0; i < len(list)-1; i++ {
		for j := i + 1; j < len(list); j++ {
			if list[i].SortDateTime < list[j].SortDateTime {
				list[i], list[j] = list[j], list[i]
			}
		}
	}
	ginplus.ResSuccess(c, list)
}

//dlCouponsList 卡券列表
func dlCouponsList(c *gin.Context) {
	req := &microScrm.CustomerCouponsRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	if req.ConsumerID == "" {
		req.ConsumerID = "a"
	}
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	dlData, err := customerSdk.CouponsList(ctx.Wrap(c), req)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("dl卡券列表接口数据错误：%+v", err)
	}
	ginplus.ResSuccess(c, dlData)
}

//DL会员信息订单列表包含商品
func dlOrderGoods(c *gin.Context) {
	req := &microScrm.CustomerDLOrderGoodsRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	//DL信息对接
	userReq := work.UserRequest{
		EnterpriseID: c.Value("enterpriseId").(int64),
		Type:         "unionId",
		Value:        req.UnionID,
		Tradeno:      req.TradeNo,
		TradeNo:      req.TradeNo,
		//Type:    "union_code",
		//Value:   "158S0002K9L",
	}
	var totalAmount int64
	group := threading.NewRoutineGroup()
	group.RunSafe(func() {
		//订单金额
		dlData, err := work.GetUserOrderList(ctx.Wrap(c), &userReq)
		if err == nil {
			totalAmount = dlData[0].TradePrice
		}
	})
	group.Wait()
	userOrderList, err := work.GetUserOrderListGoods(ctx.Wrap(c), &userReq)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	repData := make(map[string]interface{})
	repData["tradePrice"] = totalAmount
	repData["list"] = userOrderList
	ginplus.ResSuccess(c, repData)
}

// belongToMembers 添加的员工
func belongToMembers(c *gin.Context) {
	req := &microScrm.CustomerRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	list, err := customerSdk.BelongToMembers(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
}

// groupChatList 添加的员工
func groupChatList(c *gin.Context) {
	req := &microScrm.CustomerRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	response, err := customerSdk.BelongToChat(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	loginMemberID := c.Value("memberId").(int64)
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
	ginplus.ResSuccess(c, response)
}

// RemarkSave 客户备注信息保存
func RemarkSave(c *gin.Context) {
	req := &microScrm.RemarkInfo{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	detail, err := customerSdk.RemarkSave(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, detail)
}

// orderList 客户订单列表
func orderList(c *gin.Context) {
	req := &microScrm.RequestOrderList{}

	// 展示未付款的订单
	req.ShowCancel = 1

	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	list, err := customerSdk.OrderList(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
}

// orderDetail 客户订单详情
func orderDetail(c *gin.Context) {
	req := &microScrm.RequestOrderDetail{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	detail, err := customerSdk.OrderDetail(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, detail)
}

// SetExternalTags 设置外部联系人企业标签并同步到企微
func SetExternalTags(c *gin.Context) {
	req := &microScrm.SetExternalTags{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	detail, err := customerSdk.SetExternalTags(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, detail)
}

// couponList 用户优惠券列表
func couponList(c *gin.Context) {
	req := &microScrm.CustomerRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	detail, err := customerSdk.CustomerCouponList(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, detail)
}

// addressList 地址列表
func addressList(c *gin.Context) {
	req := &microScrm.CustomerRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	detail, err := customerSdk.CustomerAddressList(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, detail)
}

func Wrap(num float64, retain int) int64 {
	return int64(num * math.Pow10(retain))
}

func GetMonthStartAndEnd() map[string]string {
	myYear := strconv.Itoa(time.Now().Year())
	myMonth := strconv.Itoa(int(time.Now().Month()))
	day := time.Now().Day()
	// 数字月份必须前置补零
	if len(myMonth) == 1 {
		myMonth = "0" + myMonth
	}
	//当前时间-年
	yInt, _ := strconv.Atoi(myYear)

	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, myYear+"-"+myMonth+"-01 00:00:00", loc)
	//当前月
	newMonth := theTime.Month()
	//半年前
	month := newMonth - 6
	year := yInt
	if newMonth <= 6 {
		year = yInt - 1
		month = 12 - (6 - newMonth)
	}

	t1 := time.Date(year, month, day, 0, 0, 0, 0, time.Local).Format("20060102")
	t2 := time.Date(yInt, newMonth, day, 0, 0, 0, 0, time.Local).Format("20060102")
	result := map[string]string{"start": t1, "end": t2}
	return result
}
