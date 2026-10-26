// import { request } from '../utils/customer-request'
import commonRequest from '../utils/common-request'
// 首页
const indexApi = commonRequest({
  url: 'excitation/index',
  method: 'POST'
})
// 导购等级信息
const gradePointsApi = commonRequest({
  url: 'excitation/grade-points',
  method: 'POST'
})
// 积分规则说明
const integralRuleDescriptionApi = commonRequest({
  url: 'excitation/integral-rule-description',
  method: 'POST'
})
// 导购积分列表
const pointsListApi = commonRequest({
  url: 'excitation/points-list',
  method: 'POST'
})
// 获取积分信息
const getPointsRulesApi = commonRequest({
  url: 'excitation/get-points-rules',
  method: 'POST'
})

// 奖章墙
const medalWallApi = commonRequest({
  url: 'excitation/medal/medal-wall',
  method: 'POST'
})

// 奖章墙详情
const medalWallDetailApi = commonRequest({
  url: 'excitation/medal/medal-wall-detail',
  method: 'POST'
})

// 领取奖章数量和弹框
const recordAlreadyGotApi = commonRequest({
  url: 'excitation/medal/record-already-got',
  method: 'POST'
})

// 领取奖章
const recordMemberReceivedApi = commonRequest({
  url: 'excitation/medal/record-member-received',
  method: 'POST'
})
// 等级足迹
const getLevelRecordApi = commonRequest({
  url: 'excitation/level-record',
  method: 'POST'
})
// todo:导购金币明细
const getCoinsListApi = commonRequest({
  url: 'excitation/coins-list',
  method: 'POST'
})
// todo:获取导购金币数量第三方token
const getCoinsTokenApi = commonRequest({
  url: 'excitation/coins-settings-info',
  method: 'POST'
})
export {
  indexApi,
  gradePointsApi,
  integralRuleDescriptionApi,
  pointsListApi,
  getPointsRulesApi,
  medalWallApi,
  medalWallDetailApi,
  recordAlreadyGotApi,
  recordMemberReceivedApi,
  getLevelRecordApi,
  getCoinsListApi,
  getCoinsTokenApi
}
