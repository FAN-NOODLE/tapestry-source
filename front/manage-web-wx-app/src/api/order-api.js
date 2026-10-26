import commonRequest from '../utils/common-request'

// 顾客详情线上订单 顾客详情线下订单
const memberLineListApi = commonRequest({
  url: 'customer/customer-expense-record',
  method: 'POST'
})

// 顾客详情 顾客优惠卷
const customerCouponListApi = commonRequest({
  url: 'customer/customer-coupon-list',
  method: 'POST'
})

const customerCouponGoodsListApi = commonRequest({
  url: 'customer/customer-coupon-goods-list',
  method: 'POST'
})

// 店员列表
const orderMemberApi = commonRequest({
  url: 'member/list',
  method: 'GET'
})

// 购物订单列表
const userCrmOrderListApi = commonRequest({
  url: 'customer/user-crm-order-list',
  method: 'POST'
})

// 购物订单列表
const userCrmCouponsListApi = commonRequest({
  url: 'customer/user-crm-coupons-list',
  method: 'POST'
})

export {
  memberLineListApi,
  customerCouponListApi,
  customerCouponGoodsListApi,
  orderMemberApi,
  userCrmOrderListApi,
  userCrmCouponsListApi
}
