import { request } from '../utils/customer-request'

export const getWxlogin = data => request({
  url: "common/wx-login",
  method: 'POST'
}, data)


export const getScratchActivity = data => request({
  url: "activity/scratch/info",
  method: 'POST'
}, data)

// 刮卡
export const scratchPrize = data => request({
  url: "activity/scratch/scratch-prize",
  method: 'POST'
}, data)
// 券列表
export const couponList = data => request({
  url: "activity/scratch/coupon-list",
  method: 'POST'
}, data)
export const checkCoupon = data => request({
  url: "activity/scratch/check-coupon",
  method: 'POST'
}, data)
//分享上报
export const shareInfo = data => request({
  url: "activity/scratch/share-info",
  method: 'POST'
}, data)
