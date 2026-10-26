import commonRequest from '../utils/common-request'

// 顾客管理列表
const getCustomerListApi = commonRequest({
  url: 'customer/external-user-list',
  method: 'POST'
})
// 顾客详情
const getCustomerDetailApi = commonRequest({
  url: 'customer/external-user-detail',
  method: 'POST'
})
// 获取顾客手机号
const getCustomerPhoneApi = commonRequest({
  url: 'customer/external-user-mobile',
  method: 'POST'
})
// 获取顾客邀约和播种数量
const userDetailTipsNumApi = commonRequest({
  url: 'customer/user-detail-tips-num',
  method: 'POST'
})
// 顾客和导购的绑定关系
const userBindRelationApi = commonRequest({
  url: 'customer/user-bind-relation',
  method: 'POST'
})
// 操作顾客的印象标签
const impressionTagsOperateApi = commonRequest({
  url: 'customer/impression-tags-operate',
  method: 'POST'
})
// 操作顾客的标签
const workTagsOperateApi = commonRequest({
  url: 'customer/work-tags-operate',
  method: 'POST'
})
// 顾客列表筛选条件
const getChoiceOptionApi = commonRequest({
  url: 'customer/get-choice-option',
  method: 'POST'
})
// 修改顾客备注
const updateRemarkApi = commonRequest({
  url: 'customer/update-remark',
  method: 'POST'
})
// 添加自定义标签
const addCustomerUsertagApi = commonRequest({
  url: 'customer/usertag-add-tag',
  method: 'POST'
})
// 添加自定义标签
const deleteCustomerUsertagApi = commonRequest({
  url: 'customer/usertag-delete-tag',
  method: 'POST'
})
// 绑定自定义标签
const bindCustomerUsertagApi = commonRequest({
  url: 'customer/usertag-add-relation',
  method: 'POST'
})
// 获取绑定的自定义标签
const getUsertagApi = commonRequest({
  url: 'customer/usertag-get-relation',
  method: 'POST'
})
// 解除用户自定义标签
const deleteUsertagApi = commonRequest({
  url: 'customer/usertag-delete-relation',
  method: 'POST'
})
// 员工自定义标签列表
const getCustomerUsertagApi = commonRequest({
  url: 'customer/usertag-get-list',
  method: 'POST'
})
// 员工自定义标签（未绑定）
const getCustomerUserUnbindtagApi = commonRequest({
  url: 'customer/usertag-get-tags',
  method: 'POST'
})
// 获取流失客户统计和筛选数据
const getCusomerLostOption = commonRequest({
  url: 'customer/lost-option',
  method: 'POST'
})
// 获取客户流失列表
const getCusomerLostList = commonRequest({
  url: 'customer/lost-record',
  method: 'POST'
})

// 顾客列表全选的人员数据
const getCusomerAllCountApi = commonRequest({
  url: 'customer/select-all-num',
  method: 'POST'
})
// 获取选中的外部联系人ID
const getExternalIdsApi = commonRequest({
  url: 'customer/get-external-ids',
  method: 'POST'
})

export {
  getCustomerListApi,
  getCustomerDetailApi,
  getChoiceOptionApi,
  addCustomerUsertagApi,
  deleteCustomerUsertagApi,
  bindCustomerUsertagApi,
  getUsertagApi,
  deleteUsertagApi,
  getCustomerUsertagApi,
  getCustomerUserUnbindtagApi,
  getCusomerLostOption,
  getCusomerLostList,
  userDetailTipsNumApi,
  userBindRelationApi,
  updateRemarkApi,
  impressionTagsOperateApi,
  getCustomerPhoneApi,
  getCusomerAllCountApi,
  getExternalIdsApi,
  workTagsOperateApi
}
