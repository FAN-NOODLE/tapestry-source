import commonRequest from '../utils/common-request'

// 获取验证码
const getCodeApi = commonRequest({
  url: 'common/send-sms-code',
  method: 'POST'
})
// 手机号登陆
const loginApi = commonRequest({
  url: 'common/login',
  method: 'POST'
})
// 获取用户信息
const checkTokenApi = commonRequest({
  url: 'common/check-login',
  method: 'POST'
})

const wxworkLoginApi = commonRequest({
  url: '',
  method: 'POST'
})

// 获取用户信息
const choiceEnterpriseApi = commonRequest({
  url: 'common/choice-enterprise',
  method: 'POST'
})
// 获取注册状态
const getRegisterApi = commonRequest({
  url: 'partner-member/wx-login',
  method: 'POST'
})
// 保存注册信息
const saveRegisterApi = commonRequest({
  url: 'partner-member/register',
  method: 'POST'
})
// 获取手机号
const getMobileApi = commonRequest({
  url: 'partner-member/mobile',
  method: 'POST'
})
// 合伙人联系导购
const getPartnerContactApi = commonRequest({
  url: 'partner-member/contact',
  method: 'POST'
})
// 合伙人注册前获取hash
const getRegisterHashApi = commonRequest({
  url: 'common/hash',
  method: 'POST'
})


export {
  getCodeApi,
  getPartnerContactApi,
  getMobileApi,
  getRegisterApi,
  loginApi,
  checkTokenApi,
  wxworkLoginApi,
  saveRegisterApi,
  getRegisterHashApi,
  choiceEnterpriseApi
}
