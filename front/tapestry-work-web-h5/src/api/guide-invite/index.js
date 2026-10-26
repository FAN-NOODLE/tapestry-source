import http from "@/utils/axios";

// 获取导购当前导购时间
export const guideTimeApi = (data) => http('reserve/time-list', {
  method: "POST",
  data
})

// 创建预约单reserve/create-order
export const createOrderApi = (data) => http('reserve/create-order', {
  method: "POST",
  data
})
// 获取用户信息
export const userInfoApi = (data) => http('customer/user-simple-info', {
  method: "POST",
  data
})
// 获取门店信息
export const storeInfoApi = (data) => http('reserve/member-info', {
  method: "POST",
  data
})

// 获取服务列表
export const serviceListApi = (data) => http('reserve/service-list', {
  method: "POST",
  data
})
// 分享次数reserve/
export const orderShareNumApi = (data) => http('reserve/order-share', {
  method: "POST",
  data
})
// 删除邀约记录
export const deleteOrderApi = (data) => http('reserve/order-delete', {
  method: "POST",
  data
})
// todo:显示邀约记录
export const showOrderApi = (data) => http('reserve/order-show', {
  method: "POST",
  data
})
// 是否到店购买 & 选择顾客类型
export const changeStatus = (data) => http('reserve/change-status', {
  method: "POST",
  data
})
// 发送评价记录日志
export const commentShareApi = (data) => http('reserve/comment-share', {
  method: "POST",
  data
})
// 发送评价记录日志
export const verifyShareApi = (data) => http('reserve/verify-share', {
  method: "POST",
  data
})
// 邀约记录列表
export const orderListApi = (data)=> http('reserve/customer-reserve-lists',{
  method: "POST",
  data
})
// 验证到店时间
export const verifyTimeApi = (data) => http('reserve/check-reserve-time', {
  method: "POST",
  data
})
