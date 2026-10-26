import commonRequest from '../utils/common-request'

/** 获取预约列表 */
export const reserveListApi = commonRequest({
  url: 'activity/reserve-check/lists',
  method: 'POST'
})

/** 获取已经核销的列表 */
export const reserveCheckListApi = commonRequest({
  url: 'activity/reserve-check/check-lists',
  method: 'POST'
})

/** 核销 */
export const writeOffCouponApi = commonRequest({
  url: 'activity/reserve-check/write-off-coupon',
  method: 'POST'
})

/** 获取活动预约时刻表 */
export const timeListsAPi = commonRequest({
  url: 'activity/reserve-check/time-lists',
  method: 'POST'
})

/** 编辑活动预约 */
export const timeConfigApi = commonRequest({
  url: 'activity/reserve-check/time-config',
  method: 'POST'
})

/** 白名单权限 */
export const whiteListApi = commonRequest({
  url: 'activity/reserve-check/white-list',
  method: 'POST'
})
