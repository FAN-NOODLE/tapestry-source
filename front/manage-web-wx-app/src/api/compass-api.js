import commonRequest from '../utils/common-request'

// 顾客排行榜
const customerRankApi = commonRequest({
  url: 'customer/rank-data',
  method: 'POST'
})
// 订单排行榜
const orderRankApi = commonRequest({
  url: 'order/ranking-list',
  method: 'POST'
})
// 分销排行榜
const distributionRankApi = commonRequest({
  url: 'distribution-statistics/statistical-ranking-list',
  method: 'POST'
})
// 排行榜邀请
const inviteRankApi = commonRequest({
  url: 'reserve/rank-list',
  method: 'POST'
})
// 顾客数据概览
const customerStatisticApi = commonRequest({
  url: 'customer/statistic',
  method: 'POST'
})
// 分销统计数据概览
const distributionStatisticApi = commonRequest({
  url: 'distribution-statistics/statistical',
  method: 'POST'
})
// 订单统计数据概览
const orderStatisticApi = commonRequest({
  url: 'order/order-statistical',
  method: 'POST'
})
// 导购身份和权限类型
const memberPositionApi = commonRequest({
  url: 'member/guide-position',
  method: 'POST'
})
// 导购身份和权限类型
const storeOrDepartmentListApi = commonRequest({
  url: 'store/store-or-department-list',
  method: 'POST'
})
// ba 导购身份和权限类型
const baStoreOrDepartmentListApi = commonRequest({
  url: 'ba-kanban/ba_source_depte',
  method: 'POST'
})
// 根据门店或部门获取导购列表
const guideListApi = commonRequest({
  url: 'member/guide-list',
  method: 'POST'
})

const setViewTabApi = commonRequest({
  url: 'customer/set-view-tab',
  method: 'POST'
})

// 社群数据
const communityDataAPI = commonRequest({
  url: 'community-statistics/count',
  method: 'POST',
  apiVersion: 'brand'
})

// 未参加活动人数
const communityListAPI = commonRequest({
  url: 'community-statistics/list',
  method: 'POST',
  apiVersion: 'brand'
})

const baDataApi = commonRequest({
  url: 'ba-kanban',
  method: 'POST'
})

const allDataApi = commonRequest({
  url: 'customer/get-compass-summary',
  method: 'POST'
})
// 明星榜
const performanceApi = commonRequest({
  url: 'performance-list/star',
  method: 'POST'
})
// 门店数据
const storeListsApi = commonRequest({
  url: 'member/store-lists',
  method: 'POST'
})
// 财年财月时间维度数据
const fiscalMonthApi = commonRequest({
  url: 'customer/fiscal-month',
  method: 'POST'
})
// 邀约榜
const invitationApi = commonRequest({
  url: 'performance-list/invitation',
  method: 'POST'
})
// 活跃榜/顾客榜
const activeCustomersApi = commonRequest({
  url: 'performance-list/active-customers',
  method: 'POST'
})

export {
  customerRankApi,
  orderRankApi,
  distributionRankApi,
  guideListApi,
  customerStatisticApi,
  distributionStatisticApi,
  orderStatisticApi,
  memberPositionApi,
  storeOrDepartmentListApi,
  baStoreOrDepartmentListApi,
  setViewTabApi,
  communityDataAPI,
  communityListAPI,
  baDataApi,
  allDataApi,
  inviteRankApi,
  performanceApi,
  storeListsApi,
  activeCustomersApi,
  invitationApi,
  fiscalMonthApi
}
