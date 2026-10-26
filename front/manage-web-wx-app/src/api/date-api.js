import commonRequest from '../utils/common-request'

// 营销计划日历
const getTaskCalendarApi = commonRequest({
  url: 'marketing-plan/task-calendar',
  method: 'POST'
})

// 获取5天的营销计划列表
const getTaskListApi = commonRequest({
  url: 'marketing-plan/task-list',
  method: 'POST'
})

// 营销计划详情
const getTaskDetailApi = commonRequest({
  url: 'marketing-plan/task-detail',
  method: 'POST'
})

//
const setTaskFinishApi = commonRequest({
  url: 'marketing-plan/task-finish',
  method: 'POST'
})

export {
  getTaskCalendarApi,
  getTaskListApi,
  getTaskDetailApi,
  setTaskFinishApi
}
