import { objToUrlParams } from '@/utils/tools'

const { wejourney_companycode, wejourney_token, wejourney_redirect } = require('@/lib/config')

const imagePath = 'https://cdn-cos.coach.com.cn/uploads/'

const defaultTypeIcon = new Map([
  [1, `${imagePath}2023/v3.5.0/task/icon_1.png`], // 邀约类
  [2, `${imagePath}2023/v3.5.0/task/icon_2.png`], // 任务类
  [3, `${imagePath}2023/v3.5.0/task/icon_3.png`], // 培训类
  [4, `${imagePath}2023/v3.5.0/task/icon_2.png`] // 好友类
])

const linkList = [
  '/pages/task-no-tab/task-no-tab?type=0',
  '/pages/task-no-tab/task-no-tab?type=1',
  '/pages/task-no-tab/task-no-tab?type=2',
  'webview'
]
export const defaultPointsMap = new Map([
  [1, {
    title: '发送到店邀约',
    unit: '次',
    link: linkList[2],
    points: 1,
    icon: defaultTypeIcon.get(1),
    tip: '每天最多可得5积分哦~'
  }],
  [7, {
    title: '任务完成',
    unit: '次',
    link: linkList[0],
    points: 1,
    icon: defaultTypeIcon.get(2)
  }],
  [15, {
    title: '完成welearning培训可得积分',
    unit: '',
    rule_value: '',
    link: linkList[3],
    hideNeedPoint: true,
    icon: defaultTypeIcon.get(3)
  }],
  [14, {
    title: '完成播种',
    unit: '次',
    link: linkList[1],
    points: 1,
    icon: defaultTypeIcon.get(4),
    tip: '每天最多可得5积分哦~'
  }],
  [6, {
    title: '邀约到店确认',
    unit: '次',
    link: linkList[2],
    points: 1,
    icon: defaultTypeIcon.get(1),
    tip: '每天最多可得5积分哦~'
  }],
  [9, {
    title: '任务触达',
    unit: '人',
    link: linkList[0],
    points: 1,
    icon: defaultTypeIcon.get(2)
  }],
  [4, {
    title: '邀约成功',
    unit: '次',
    link: linkList[2],
    points: 5,
    icon: defaultTypeIcon.get(1)
  }],
  [10, {
    title: '新增好友',
    unit: '人',
    link: false,
    points: 5,
    icon: defaultTypeIcon.get(4)
  }],
  [11, {
    title: '每月客户流失高于',
    unit: '人',
    link: false,
    points: 2,
    icon: defaultTypeIcon.get(4)
  }]
])

/**
 * 获取任务信息
 * @param {object} item
 */
export function getTaskInfo (item) {
  const value = defaultPointsMap.get(item.rule_category_id)
  const frequency = ['', '天', '周', '月', '年'][item.astrict_type]
  // 获取提示信息
  if (frequency && item.max_integral) {
    value.tip = `每${frequency}最多可得${item.max_integral}积分哦~`
  }

  if (item.max_integral !== 0) {
    item.title_text = `${item.now_integral}/${item.max_integral}`
  }

  console.log(value)
  return {
    ...item,
    ...value
  }
}

/**
 * 生成loginToken
 * @param {string} token_salt
 * @returns
 */
export function getRedirectPrarms (moduleType, time, isWelearning) {
  console.log('getRedirectPrarms', moduleType, time, isWelearning)
  const app = getApp()
  const {mobile} = app.globalData.currentEnterprise
  const {userid} = app.globalData
  let loginToken = ''
  if (!mobile && isWelearning) {
    wx.showToast('请先绑定手机号')
    throw new Error('请先绑定手机号')
  }
  if (!userid && !isWelearning) {
    wx.showToast('请先绑定工号')
    throw new Error('请先绑定工号')
  }
  if (isWelearning) {
    loginToken = {
      moduleType: moduleType,
      moduleId: '0',
      uniqueKey: 'mobile',
      uniqueValue: '',
      companyCode: 'coach',
      time: time,
      token: 'REDACTED_FOR_SOURCE_SHARING'
    }
  } else {
    loginToken = {
      moduleType: moduleType,
      moduleId: '0',
      uniqueKey: 'workNo',
      uniqueValue: '',
      companyCode: wejourney_companycode,
      time: time,
      token: wejourney_token
    }
  }
  let params = ''
  let redirectUrl = ''
  if (isWelearning) {
    loginToken.uniqueValue = mobile
    params = objToUrlParams(loginToken)
    redirectUrl = 'https://cust.strong365.com/index.php?g=Wap&m=SingleLogin&a=index&' + params
  } else {
    loginToken.uniqueValue = userid
    params = objToUrlParams(loginToken)
    redirectUrl = wejourney_redirect + '/index.php?g=Wap&m=SingleLogin&a=index&' + params
  }
  return {
    params: params,
    redirectUrl
  }
}
