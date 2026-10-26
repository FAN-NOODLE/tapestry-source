// import { request } from '../utils/customer-request'
import commonRequest from '../utils/common-request'
// 核销二维码
const verifyCodeApi = commonRequest({
  url: 'activity/check/write-off-coupon',
  method: 'POST'
})
// 已核销的列表
const HasVerifyListApi = commonRequest({
  url: 'activity/check/coupon-list',
  method: 'POST'
})
export {
  verifyCodeApi, HasVerifyListApi
}
