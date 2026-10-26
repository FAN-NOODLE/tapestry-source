package mysql

import (
	"database/sql"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/coupon"
)

type CouponRepo struct {
	*repos.SQLRepo
}

// AddCoupon 添加优惠券
func (u *CouponRepo) AddCoupon(c ctx.Context, data *coupon.BasicInfoEntity) (int64, error) {
	stmt, err := u.DB.Template("t_coupon/insert-coupon")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// 更新优惠券
func (u *CouponRepo) UpdateCoupon(c ctx.Context, data *coupon.BasicInfoEntity) (int64, error) {
	stmt, err := u.DB.Template("t_coupon/update-coupon")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// 统计优惠券
func (u *CouponRepo) CountCoupon(c ctx.Context, request *coupon.BasicInfoEntity) int64 {
	stmt, err := u.DB.Template("t_coupon/count-coupon", u.Cond2Args(request))
	if err != nil {
		return 0
	}
	countResponseInfo := &coupon.CountResponse{}
	getErr := stmt.NamedGet(c, countResponseInfo, request)
	if getErr != nil {
		return 0
	}
	return countResponseInfo.Count
}

// 删除优惠券
func (u *CouponRepo) DeleteCoupon(c ctx.Context, data *coupon.BasicInfoEntity) (int64, error) {
	stmt, err := u.DB.Template("t_coupon/delete-coupon")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

//  添加用户优惠券
func (u *CouponRepo) AddCouponUser(c ctx.Context, data *coupon.BasicCouponUserInfoEntity) (int64, error) {
	stmt, err := u.DB.Template("t_coupon_user/insert-coupon-user")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// 更新用户优惠券
func (u *CouponRepo) UpdateCouponUser(c ctx.Context, data *coupon.BasicCouponUserInfoEntity) (int64, error) {
	stmt, err := u.DB.Template("t_coupon_user/update-coupon-user")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// 更新用户优惠券用户ID
func (u *CouponRepo) UpdateCouponUserUserId(c ctx.Context, data *coupon.BasicCouponUserInfoEntity) (int64, error) {
	stmt, err := u.DB.Template("t_coupon_user/update-coupon-user-user-id")
	if err != nil {
		return 0, err
	}
	res, errs := stmt.NamedExec(c, data)
	if errs != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// 统计用户优惠券
func (u *CouponRepo) CountCouponUser(c ctx.Context, request *coupon.BasicCouponUserInfoEntity) int64 {
	stmt, err := u.DB.Template("t_coupon_user/count-coupon-user", u.Cond2Args(request))
	if err != nil {
		return 0
	}
	countResponseInfo := &coupon.CountResponse{}
	getErr := stmt.NamedGet(c, countResponseInfo, request)
	if getErr != nil {
		return 0
	}
	return countResponseInfo.Count
}

// CouponUserList 用户优惠券列表
func (u *CouponRepo) CouponUserList(c ctx.Context, request *coupon.UserCouponListRequest) (int64, *[]*coupon.UserCouponListResponse, error) {
	option := &repos.FindPageOption{
		TplPrefix: "t_coupon_user/select-user-coupon-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]*coupon.UserCouponListResponse{}
		},
	}
	offset := (request.Page - 1) * request.Size
	total, list, err := u.FindPage(c, repos.PageVO{
		Offset: offset,
		Limit:  request.Size,
	}, *request, option)

	if err != nil && err != sql.ErrNoRows {
		return 0, nil, err
	}

	return total, (list).(*[]*coupon.UserCouponListResponse), nil
}

func (u *CouponRepo) SelectCouponUser(c ctx.Context, batch *coupon.BatchCouponUserEntity) (*[]coupon.BasicCouponUserInfoEntity, error) {
	mapping := make(map[string]interface{})
	mapping["CouponUserIds"] = batch.CouponUserIds
	mapping["enterprise_id"] = batch.EnterpriseId

	stmt, err := u.DB.Template("t_coupon_user/select-coupon-user", &mapping)
	if err != nil {
		return nil, err
	}

	list := &[]coupon.BasicCouponUserInfoEntity{}
	stmtErr := stmt.NamedSelect(c, list, mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	return list, nil
}

// DetailByCoupon 详情返回接口
func (u CouponRepo) DetailByCoupon(c ctx.Context, request *coupon.WithGoodsRequest) (*coupon.BasicInfoEntity, error) {
	cop := new(coupon.BasicInfoEntity)
	stmt, err := u.DB.Template("t_coupon/detail-coupon")
	if err != nil {
		return nil, err
	}
	errS := stmt.Get(c, cop, request.EnterpriseID, request.CouponId)
	if errS != nil {
		return nil, errS
	}

	return cop, nil
}

func (u *CouponRepo) SelectCouponUserUserIdList(c ctx.Context) (*[]coupon.BasicCouponUserInfoEntity, error) {
	mapping := make(map[string]interface{})
	mapping["UserId"] = 0
	stmt, err := u.DB.Template("t_coupon_user/select-user-old-coupon-list", &mapping)
	if err != nil {
		return nil, err
	}

	list := &[]coupon.BasicCouponUserInfoEntity{}
	stmtErr := stmt.NamedSelect(c, list, mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	return list, nil
}

// 更新老用户数据
func (u *CouponRepo) UpdateUnionIdCouponUserUserId(c ctx.Context, couponInfo *coupon.BasicCouponUserInfoEntity) (int64, error) {
	stmt, err := u.DB.Template("t_coupon_user/update-coupon-user-union-id")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, couponInfo)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (m *CouponRepo) GetMarketingStrategyIsUsedBatch(cc ctx.Context, wxCouponIds []int64) (UsedCouponIds []int64, err error) {
	mapping := make(map[string]interface{})
	mapping["wx_coupon_ids"] = wxCouponIds

	stmt, err := m.DB.Template("t_coupon/select-coupon-is-used-by-coupon_ids", &mapping)
	if err != nil {
		return nil, err
	}

	stmtErr := stmt.Select(cc, &UsedCouponIds)
	if stmtErr != nil {
		return nil, stmtErr
	}
	return UsedCouponIds, err
}

func (m *CouponRepo) GetAllCouponAvailableList(cc ctx.Context, request coupon.WxCouponAvailableRequest) ([]*coupon.WxCouponAvailable, error) {
	mapping := make(map[string]interface{})
	mapping["cond"] = request
	stmt, err := m.DB.Template("t_coupon/wx-coupon-available-list", &mapping)
	if err != nil {
		return nil, err
	}

	var list = make([]*coupon.WxCouponAvailable, 0)
	stmtErr := stmt.NamedSelect(cc, &list, mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	return list, nil
}
