
import commonRequest from '../utils/common-request'

// 获取筛选条件
export const getChoiceOptionApi = commonRequest({
  url: 'customer/external-search/get-choice-option',
  method: 'POST'
})
// 获取顾客列表
export const getCustomerListApi = commonRequest({
  url: 'customer/external-search/list',
  method: 'POST'
})
// 获取顾客详情
export const getCustomerDetailApi = commonRequest({
  url: 'customer/external-search/detail',
  method: 'POST'
})
// 获取顾客详情 crm
export const getCustomerDetailByCrmApi = commonRequest({
  url: 'customer/external-search/detail-by-crm',
  method: 'POST'
})
// 获取顾客详情 推荐商品
export const getCustomerInterestedProductsApi = commonRequest({
  url: 'customer/external-search/customer-interested-products',
  method: 'POST'
})

// 获取顾客详情 编辑用户描述
export const updateCustomerDescApi = commonRequest({
  url: 'customer/update-desc',
  method: 'POST'
})

// 修改顾客备注
export const updateCustomerRemarkApi = commonRequest({
  url: 'customer/update-remark',
  method: 'POST'
})

// 添加人群包
export const addCrowdApi = commonRequest({
  url: 'customer/external-crowd/add',
  method: 'POST'
})
// 获取人群包
export const gerCrowdListApi = commonRequest({
  url: 'customer/external-crowd/list',
  method: 'POST'
})

// 获取人群包详情
export const gerCrowdInfoApi = commonRequest({
  url: 'customer/external-crowd/info',
  method: 'POST'
})

// 更新排序
export const updateCrowdSortApi = commonRequest({
  url: 'customer/external-crowd/update-sort',
  method: 'POST'
})

// 更新名称
export const updateCrowdApi = commonRequest({
  url: 'customer/external-crowd/edit',
  method: 'POST'
})

// 删除
export const deleteCrowdApi = commonRequest({
  url: 'customer/external-crowd/del',
  method: 'POST'
})

// 更新联系状态
export const updateCrowdContactApi = commonRequest({
  url: 'customer/external-crowd/update-contact',
  method: 'POST'
})

// 获取所有可以用的openid
export const getOpenIdsApi = commonRequest({
  url: 'customer/external-crowd/get-open-ids',
  method: 'POST'
})
// 获取所有可以用的openid
export const getCustomerAllCountApi = commonRequest({
  url: 'customer/external-search/get-open-ids',
  method: 'POST'
})

// 获取待办任务
export const getCustomerTaskListApi = commonRequest({
  url: 'marketing-plan/customer-task-list',
  method: 'POST'
})
