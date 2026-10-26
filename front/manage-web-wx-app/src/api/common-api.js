import commonRequest from '../utils/common-request'

// 是否是第一次，啥都行
export const checkFirstOpenApi = commonRequest({
  url: 'common/check-first-open',
  method: 'POST'
})

