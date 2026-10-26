import commonRequest from '../utils/common-request'

/** 导购--邀请函模板列表 */
const getInviteConfigList = commonRequest({
  url: 'invite/config-list',
  method: 'POST'
})

/** 导购--邀请函模板详情 */
const getInvitationInfo = commonRequest({
  url: 'invite/config-detail',
  method: 'POST'
})

/** 发送邀请函 */
const sendInvite = commonRequest({
  url: 'invite/create-record',
  method: 'POST'
})

/** 查看邀请函 */
const recordInfo = commonRequest({
  url: 'invite/record-info',
  method: 'POST'
})

export { getInviteConfigList, sendInvite, getInvitationInfo
  , recordInfo }
