package coupon

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/modules/marketingPromotion"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/sdk/boss"
	"micro-scrm/internal/pkg/uid"
)

// Service  应用服务
type Service struct {
	repo           Repository
	customerRepo   customer.Repository
	enterpriseRepo enterprise.Repository
}

// SyncCoupon 同步优惠券
func (s *Service) SyncCoupon(c ctx.Context, request *BasicInfoEntity, event common.DdcEvent) (bool, error) {
	log.Info(c).Msgf("请求同步优惠券参数，request:%+v,事件:%s", request, event)
	nowTime := time.Now().Unix()
	// 3.1 删除事件
	if event == common.DdcEventEnum.WxCouponDelete {
		request.DeletedAt = nowTime
		_, err := s.repo.DeleteCoupon(c, request)
		if err != nil {
			return true, err
		}
		return true, nil
	}

	// 3.2 新增事件
	if event == common.DdcEventEnum.WxCouponAdd {
		request.ID = uid.ID(c)
		request.CreatedAt = nowTime
		request.UpdatedAt = nowTime
		_, err := s.repo.AddCoupon(c, request)
		if err != nil {
			return true, err
		}
	}

	// 3.3 更新事件
	if event == common.DdcEventEnum.WxCouponUpdate {
		total := s.repo.CountCoupon(c, request)
		request.UpdatedAt = nowTime
		if total > 0 {
			_, err := s.repo.UpdateCoupon(c, request)
			if err != nil {
				return true, err
			}
		} else {
			request.ID = uid.ID(c)
			request.CreatedAt = nowTime
			_, err := s.repo.AddCoupon(c, request)
			if err != nil {
				return true, err
			}
		}
	}
	return true, nil
}

// SyncCouponUser 同步用户优惠券
func (s *Service) SyncCouponUser(c ctx.Context, request *[]BasicCouponUserInfoEntity, event common.DdcEvent) (bool, error) {
	log.Info(c).Msgf("请求同步用户优惠券参数，request:%+v,事件:%s", request, event)

	nowTime := time.Now().Unix()
	// 3.1 新增事件
	if event == common.DdcEventEnum.CouponUserAdd {
		for _, value := range *request {
			appid, _ := s.enterpriseRepo.QueryOneApp(c, value.EnterpriseId, customer.CustomerTypeShop)

			user, err := s.customerRepo.QueryEnterpriseUser(c, value.EnterpriseId, &customer.SourceUserBase{
				UnionID:      value.UnionId,
				CustomerType: customer.CustomerTypeShop,
				AppID:        appid,
				OpenID:       value.OpenId,
			})
			if err != nil {
				log.Err(c, err).Msgf("用户不存在: %d %s", value.EnterpriseId, value.UnionId)
			} else {
				value.UserId = user.ID
			}
			value.ID = uid.ID(c)
			value.CreatedAt = nowTime
			value.UpdatedAt = nowTime
			_, err = s.repo.AddCouponUser(c, &value)

			if err != nil {
				return true, err
			}
		}
	}

	// 3.2 更新事件走
	if event == common.DdcEventEnum.CouponUserUpdate {

		where := &BatchCouponUserEntity{}

		for _, value := range *request {
			where.CouponUserIds = append(where.CouponUserIds, value.CouponUserId)
			where.EnterpriseId = value.EnterpriseId
		}

		list, _ := s.repo.SelectCouponUser(c, where)
		var CouponUserIds []int64

		if list != nil {
			for _, val := range *list {
				CouponUserIds = append(CouponUserIds, val.CouponUserId)
			}
		}

		for _, reqValue := range *request {
			user, err := s.customerRepo.QueryEnterpriseUser(c, reqValue.EnterpriseId, &customer.SourceUserBase{UnionID: reqValue.UnionId})
			if err != nil {
				log.Err(c, err).Msgf("用户不存在: %d %s", reqValue.EnterpriseId, reqValue.UnionId)
			} else {
				reqValue.UserId = user.ID
			}
			reqValue.UpdatedAt = nowTime
			if IsContain(CouponUserIds, reqValue.CouponUserId) {
				if reqValue.Status == 3 {
					reqValue.UseTime = nowTime
				}
				_, err := s.repo.UpdateCouponUser(c, &reqValue)
				if err != nil {
					return true, err
				}
			} else {
				reqValue.ID = uid.ID(c)
				reqValue.CreatedAt = nowTime
				_, err := s.repo.AddCouponUser(c, &reqValue)
				if err != nil {
					return true, err
				}
			}
		}
	}

	return true, nil
}

// CouponUserList 用户优惠券列表
func (s *Service) CouponUserList(c ctx.Context, req *UserCouponListRequest) (int64, *[]*UserCouponListResponse, error) {
	total, list, err := s.repo.CouponUserList(c, req)
	if total > 0 {
		for _, item := range *list {
			item.RangeGoods = make([]int64, 0)
			type spu struct {
				SpuIds []int64
			}
			spuData := &spu{}
			if item.IsAll != WxCouponIsAllYes {
				json.Unmarshal([]byte(item.Range), spuData)

				item.RangeGoods = spuData.SpuIds
			}
			item.CouponName = fmt.Sprintf("满%d减%s", item.CouponMininumn/100, strconv.FormatFloat(float64(item.CouponValue)/100, 'f', 2, 32))
		}
	}
	return total, list, err
}

func IsContain(items []int64, item int64) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// WxCouponList 优惠券列表
func (s *Service) WxCouponList(c ctx.Context, cond WxCouponRequest) (WxCouponListResponse, error) {
	var cmp WxCouponListResponse

	offset := (cond.Page - 1) * cond.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(cond.Size)}
	option := &repos.FindPageOption{
		TplPrefix: "t_coupon/wx-coupon-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]WxCoupon{}
		},
	}
	total, list, err := s.repo.FindPage(c, pageVo, cond, option)
	if errorIgnore.Ignore(err) != nil {
		return cmp, err
	}
	resList := list.(*[]WxCoupon)

	// 获取小程序的logo+name
	appReq := &boss.AuthDetailRequest{
		EnterpriseID: cond.EnterpriseID,
		AuthType:     1,
	}
	authDetailResponse, authDetailErr := boss.AuthDetail(c, appReq)
	if authDetailErr != nil {
		return cmp, authDetailErr
	}

	var responseWxCoupon []WxCoupon
	for _, value := range *resList {
		value.CouponName = fmt.Sprintf("满%s减%s", strconv.FormatFloat(float64(value.CouponMininumn)/100, 'f', 2, 32), strconv.FormatFloat(float64(value.CouponValue)/100, 'f', 2, 32))
		value.ShareAppId = authDetailResponse.Data.AppID
		value.SharePath = getParams(c, cond, value)
		value.DefaultPicture = "https://img-crs.vchangyi.com/goods16099009171810.png"
		responseWxCoupon = append(responseWxCoupon, value)
	}

	cmp.Page = cond.Page
	cmp.Size = cond.Size
	cmp.List = responseWxCoupon
	if responseWxCoupon == nil {
		cmp.List = []WxCoupon{}
	}
	cmp.Total = total

	return cmp, nil
}

// 获取太阳码携带参数
func getParams(c ctx.Context, req WxCouponRequest, value WxCoupon) string {
	params := "/" + marketingPromotion.MarketingPath[marketingPromotion.MarketingCoupon] + ".html"
	params += "?id=" + strconv.FormatInt(value.CouponID, 10) + "&coupon_id=" + strconv.FormatInt(value.CouponID, 10) + "&store_id=0&mem_uid=" + req.MemberUserId + "&share_uid=" + req.MemberUserId + "&source=scrm"
	return params
}

// couponInfo 查询详情
func (s Service) CouponInfo(c ctx.Context, request *WithGoodsRequest) (*BasicInfoEntity, error) {
	return s.repo.DetailByCoupon(c, request)
}

func (s *Service) GetAllWxCouponAvailableList(c ctx.Context, cond WxCouponAvailableRequest) (*WxCouponAvailableListResponse, error) {
	resp := &WxCouponAvailableListResponse{}
	offset := (cond.Page - 1) * cond.Size
	end := offset + cond.Size
	cond.NowTime = time.Now().Unix()

	ResStrategyList, err := s.repo.GetAllCouponAvailableList(c, cond)

	if len(ResStrategyList) == 0 || err != nil {
		return resp, nil
	}
	//查询有哪些券正在被其他策略占用
	var wxCouponIds []int64
	for _, value := range ResStrategyList {
		wxCouponIds = append(wxCouponIds, value.CouponID)
	}
	couponIDs := make(map[int64]int64)
	res, _ := s.repo.GetMarketingStrategyIsUsedBatch(c, wxCouponIds)
	for _, v := range res {
		couponIDs[v] = 1
	}

	//查询是否被占用，进行重新排序
	var responseWxCoupon []*WxCouponAvailable
	var responseWxCouponIsUsed []*WxCouponAvailable
	for _, value := range ResStrategyList {
		value.Condition = fmt.Sprintf("满%s减%s", strconv.FormatFloat(float64(value.CouponMininumn)/100, 'f', 2, 32), strconv.FormatFloat(float64(value.CouponValue)/100, 'f', 2, 32))
		if _, ok := couponIDs[value.CouponID]; ok {
			value.IsUsed = 1
			responseWxCouponIsUsed = append(responseWxCouponIsUsed, value)
		} else {
			responseWxCoupon = append(responseWxCoupon, value)
		}
	}

	for _, vv := range responseWxCouponIsUsed {
		responseWxCoupon = append(responseWxCoupon, vv)
	}

	resp.Total = int64(len(ResStrategyList))
	resp.Page = cond.Page
	resp.Size = cond.Size
	if end >= resp.Total {
		end = resp.Total
	}
	for _, vv := range responseWxCoupon[offset:end] {
		if vv != nil {
			resp.List = append(resp.List, vv)
		}
	}
	return resp, nil
}
