import commonRequest from '../utils/common-request'

// 在职继承顾客列表
const getInheritCustomerListApi = commonRequest({
  url: 'customer-inherit/customer-list',
  method: 'POST'
})
// 在职继承顾客列表 - beta
const getInheritCustomerBetaListApi = commonRequest({
  url: 'customer-inherit/search-list',
  method: 'POST'
})

// 在职继承顾客列表
const getInheritMemberListApi = commonRequest({
  url: 'customer-inherit/member-list',
  method: 'POST'
})

// 分配记录列表
const getInheritRecordListApi = commonRequest({
  url: 'customer-inherit/list',
  method: 'POST'
})

// 分配详情
const getInheritRecordInfoApi = commonRequest({
  url: 'customer-inherit/info',
  method: 'POST'
})

// 提交转移信息
const getInheritTransferInfoApi = commonRequest({
  url: 'customer-inherit/transfer',
  method: 'POST'
})

// 获取全选数量
const getSelectAllCountApi = commonRequest({
  url: 'customer-inherit/select-all-num',
  method: 'POST'
})

export {
  getInheritCustomerListApi,
  getInheritMemberListApi,
  getInheritRecordListApi,
  getInheritRecordInfoApi,
  getInheritTransferInfoApi,
  getSelectAllCountApi,
  getInheritCustomerBetaListApi
}
