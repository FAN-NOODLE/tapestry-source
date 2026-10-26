import commonRequest from '../utils/common-request'
// 获取会员查询筛选条件picker
const memberSearchPickerApi = commonRequest({
  url: 'customer/customer-option',
  method: 'POST'
})
// 获取会员查询列表
const memberSearchListApi = commonRequest({
  url: 'customer/customer-list',
  method: 'POST'
})
// 获取消费记录详情
const userCrmOrderApi = commonRequest({
  url: 'customer/user-crm-order-list',
  method: 'POST'
})
// 会员聚合全选的人员数据
const getMemberAllCountApi = commonRequest({
  url: 'customer/customer-all-num',
  method: 'POST'
})
// 会员聚合全选的人员群发的external_ids
const getMemberExternalIdApi = commonRequest({
  url: 'customer/customer-external-ids',
  method: 'POST'
})
export {
  memberSearchPickerApi,
  memberSearchListApi,
  userCrmOrderApi,
  getMemberAllCountApi,
  getMemberExternalIdApi
}
