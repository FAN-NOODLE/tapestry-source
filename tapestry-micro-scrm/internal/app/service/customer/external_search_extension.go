package customerSvc

import (
	"micro-scrm/internal/app"
	externalSearchDao "micro-scrm/internal/app/dao/external_search"
	openReq "micro-scrm/internal/app/http/request/open"
	externalSearchModel "micro-scrm/internal/app/models/external_search"
	"sort"
	"strings"
	"time"
)

const (
	couponStatusMask = 1 << iota
	couponBy30daysStatusMask
	birthdayCouponStatusMask
	renewMembershipCouponStatusMask
	welcomeCouponStatusMask
	upgradeCouponStatusMask
)

type customerCouponExtensionGroup struct {
	columns []string
	data    []externalSearchModel.ExternalCustomerExtensionModel
}

func (s *ExternalSearch) saveReserveExtension(c app.Ctx, enterpriseId int64, dataMap map[string]externalSearchDao.ExternalExtraData) {
	now := time.Now().Unix()
	data := make([]externalSearchModel.ExternalCustomerExtensionModel, 0, len(dataMap))
	for unionid, item := range dataMap {
		if unionid == "" {
			continue
		}
		data = append(data, externalSearchModel.ExternalCustomerExtensionModel{
			Unionid:          unionid,
			ValidReserveAt:   int64(item.ValidReserveAt),
			SuccessReserveAt: int64(item.SuccessReserveAt),
			CreatedAt:        now,
			UpdatedAt:        now,
		})
	}
	if err := s.externalCustomerExtensionDao.UpsertReserveData(c, enterpriseId, data); err != nil {
		s.Log().Error(c).Msgf("客户扩展表邀约数据影子写入失败,enterprise_id:%d,error:%v", enterpriseId, err)
	}
}

func (s *ExternalSearch) saveTaskExtension(c app.Ctx, enterpriseId int64, dataMap map[string]externalSearchDao.ExternalExtraData) {
	now := time.Now().Unix()
	data := make([]externalSearchModel.ExternalCustomerExtensionModel, 0, len(dataMap))
	for unionid, item := range dataMap {
		if unionid == "" {
			continue
		}
		data = append(data, externalSearchModel.ExternalCustomerExtensionModel{
			Unionid:       unionid,
			TaskFinishAt:  int64(item.TaskFinishAt),
			TaskCreatedAt: item.TaskCreatedAt,
			CreatedAt:     now,
			UpdatedAt:     now,
		})
	}
	if err := s.externalCustomerExtensionDao.UpsertTaskData(c, enterpriseId, data); err != nil {
		s.Log().Error(c).Msgf("客户扩展表任务数据影子写入失败,enterprise_id:%d,error:%v", enterpriseId, err)
	}
}

func (s *ExternalSearch) saveCouponExtension(c app.Ctx, enterpriseId int64, dataMap map[string]openReq.UpdateCouponParams) {
	groups := buildCouponExtensionGroups(dataMap, time.Now().Unix())
	for _, group := range groups {
		if err := s.externalCustomerExtensionDao.UpsertCouponData(c, enterpriseId, group.data, group.columns); err != nil {
			s.Log().Error(c).Msgf("客户扩展表优惠券数据影子写入失败,enterprise_id:%d,columns:%s,error:%v", enterpriseId, strings.Join(group.columns, ","), err)
		}
	}
}

func (s *ExternalSearch) saveRelationshipExtension(c app.Ctx, req openReq.BatchUpdateActiveReq, relationships []externalSearchModel.ExternalRelationshipModel) {
	data := buildRelationshipExtensionData(req.DataList, relationships, time.Now().Unix())
	if err := s.externalRelationshipExtensionDao.UpsertData(c, req.EnterpriseId, data); err != nil {
		s.Log().Error(c).Msgf("关系扩展表互动数据影子写入失败,enterprise_id:%d,error:%v", req.EnterpriseId, err)
	}
}

func buildRelationshipExtensionData(requestData []openReq.UpdateActiveParams, relationships []externalSearchModel.ExternalRelationshipModel, now int64) []externalSearchModel.ExternalRelationshipExtensionModel {
	requestMap := make(map[string]openReq.UpdateActiveParams, len(requestData))
	for _, item := range requestData {
		requestMap[relationshipIdentity(item.UnionId, item.UserId)] = item
	}
	data := make([]externalSearchModel.ExternalRelationshipExtensionModel, 0, len(relationships))
	for _, relationship := range relationships {
		if relationship.ID == 0 {
			continue
		}
		item, ok := requestMap[relationshipIdentity(relationship.Unionid, relationship.UserID)]
		if !ok || item.ActiveTime <= 0 && item.InteractTime <= 0 {
			continue
		}
		data = append(data, externalSearchModel.ExternalRelationshipExtensionModel{
			RelationID:   relationship.ID,
			Unionid:      relationship.Unionid,
			UserID:       relationship.UserID,
			InteractTime: item.InteractTime,
			ActiveTime:   item.ActiveTime,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
	}
	return data
}

func buildCouponExtensionGroups(dataMap map[string]openReq.UpdateCouponParams, now int64) map[int]customerCouponExtensionGroup {
	groups := make(map[int]customerCouponExtensionGroup)
	for unionid, item := range dataMap {
		if unionid == "" {
			continue
		}
		data := externalSearchModel.ExternalCustomerExtensionModel{
			Unionid:         unionid,
			CouponUpdatedAt: int64(item.UpdatedAt),
			CreatedAt:       now,
			UpdatedAt:       now,
		}
		mask := 0
		if item.HasCoupon > 0 {
			mask |= couponStatusMask
			data.CouponStatus = couponStatus(item.HasCoupon)
		}
		if item.HasCouponBy30days > 0 {
			mask |= couponBy30daysStatusMask
			data.CouponBy30daysStatus = couponStatus(item.HasCouponBy30days)
		}
		if item.HasBirthdayCoupon > 0 {
			mask |= birthdayCouponStatusMask
			data.BirthdayCouponStatus = couponStatus(item.HasBirthdayCoupon)
		}
		if item.HasRenewMembershipCoupon > 0 {
			mask |= renewMembershipCouponStatusMask
			data.RenewMembershipCouponStatus = couponStatus(item.HasRenewMembershipCoupon)
		}
		if item.HasWelcomeCoupon > 0 {
			mask |= welcomeCouponStatusMask
			data.WelcomeCouponStatus = couponStatus(item.HasWelcomeCoupon)
		}
		if item.HasUpgradeCoupon > 0 {
			mask |= upgradeCouponStatusMask
			data.UpgradeCouponStatus = couponStatus(item.HasUpgradeCoupon)
		}
		if mask == 0 {
			continue
		}
		group := groups[mask]
		if len(group.columns) == 0 {
			group.columns = couponColumns(mask)
		}
		group.data = append(group.data, data)
		groups[mask] = group
	}
	return groups
}

func couponColumns(mask int) []string {
	columns := make([]string, 0, 6)
	columnMasks := map[string]int{
		"coupon_status":                  couponStatusMask,
		"coupon_by_30days_status":        couponBy30daysStatusMask,
		"birthday_coupon_status":         birthdayCouponStatusMask,
		"renew_membership_coupon_status": renewMembershipCouponStatusMask,
		"welcome_coupon_status":          welcomeCouponStatusMask,
		"upgrade_coupon_status":          upgradeCouponStatusMask,
	}
	for column, columnMask := range columnMasks {
		if mask&columnMask != 0 {
			columns = append(columns, column)
		}
	}
	sort.Strings(columns)
	return columns
}

func couponStatus(value int) int {
	if value == openReq.HasCoupon {
		return externalSearchModel.HasCoupon
	}
	return externalSearchModel.NoCoupon
}

func relationshipIdentity(unionid, userId string) string {
	return unionid + "\x00" + userId
}
