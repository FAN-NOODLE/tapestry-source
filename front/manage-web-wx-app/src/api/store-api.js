import commonRequest from '../utils/common-request'

// 门店列表
const getStoreListApi = commonRequest({
  url: 'store/list',
  method: 'POST'
})
// 门店详情
const getStoreDetailApi = commonRequest({
  url: 'store/detail',
  method: 'GET'
})
// 门店TAB业绩统计(推广订单、服务订单红点数字)
const getOrderCountApi = commonRequest({
  url: 'store/home-order-tab',
  method: 'POST'
})

// 店员首页订单业绩统计(推广订单、服务订单业绩统计接口)
const getCountApi = commonRequest({
  url: 'store/home-statistics',
  method: 'POST'
})
// 门店信息修改
const setStoreApi = commonRequest({
  url: 'store/edit',
  method: 'POST'
})
// 门店经纬度转省市区
const gisToAddressApi = commonRequest({
  url: 'store/gis-to-address',
  method: 'GET'
})
// 店员详情
const getMemberDetailApi = commonRequest({
  url: 'member/detail',
  method: 'GET'
})
// 修改接收短信状态
const changeSmsStatusApi = commonRequest({
  url: 'member/update-sms-status',
  method: 'POST'
})
// 修改电话
const changeSmsApi = commonRequest({
  url: 'member/update-sms-mobile',
  method: 'POST'
})
// 导购生成海报
const generateWorkerPosterApi = commonRequest({
  url: 'member/generate-worker-poster',
  method: 'POST'
})
// 获取搜索历史记录
const getSearchHistoryApi = commonRequest({
  url: 'inventory/search-keywords',
  method: 'POST'
})
// 删除搜索历史记录
const deleteSearchHistoryApi = commonRequest({
  url: 'inventory/del-keywords',
  method: 'POST'
})
// 查询商品列表
const searchListGroupSkuApi = commonRequest({
  url: 'source/brand/search-list-group-sku',
  method: 'POST'
})


export {
  getStoreListApi,
  getStoreDetailApi,
  getOrderCountApi,
  getCountApi,
  setStoreApi,
  gisToAddressApi,
  getMemberDetailApi,
  changeSmsStatusApi,
  changeSmsApi,
  generateWorkerPosterApi,
  getSearchHistoryApi,
  deleteSearchHistoryApi,
  searchListGroupSkuApi
}
