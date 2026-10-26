import commonRequest from '../utils/common-request'

// 店员列表
const getMemberListApi = commonRequest({
  url: 'member/list',
  method: 'GET'
})

// 个人详情
const getMemberDetailApi = commonRequest({
  url: 'member/detail',
  method: 'GET'
})

// 店员详情
const getMemberWorkerDetailApi = commonRequest({
  url: 'member/worker-detail',
  method: 'GET'
})

// 离岗复岗
const editStatus = commonRequest({
  url: 'member/edit-status',
  method: 'GET'
})

// 店员订单业绩统计
const getMemberCountApi = commonRequest({
  url: 'member/member-statistics',
  method: 'POST'
})

// 合伙人我的信息
const getPartnerInfoApi = commonRequest({
  url: 'distribution-partner/my',
  method: 'POST'
})

// 合伙人我的页面联系他发送企业微信消息
const contactInviterApi = commonRequest({
  url: 'partner-member/contact',
  method: 'POST'
})

// 手机号保护
const mobileChangeApi = commonRequest({
  url: 'distribution-partner/mobile-switch',
  method: 'POST'
})

// 接收提醒
const receiveRemindApi = commonRequest({
  url: 'distribution-partner/sms-switch',
  method: 'POST'
})

export {
  getMemberListApi,
  getMemberDetailApi,
  editStatus,
  getMemberCountApi,
  getMemberWorkerDetailApi,
  getPartnerInfoApi,
  contactInviterApi,
  mobileChangeApi,
  receiveRemindApi
}
