import commonRequest from '../utils/common-request'

// 退群人员顾客信息列表接口
const leaveGroupApi = commonRequest({
  url: 'marketing-plan/leave-group-chat-lists',
  method: 'POST'
})

// 首页-任务计划列表
const planDetailApi = commonRequest({
  url: 'marketing-plan/plan-detail',
  method: 'POST'
})

// 首页-任务计划列表
const planListApi = commonRequest({
  url: 'marketing-plan/plan-list',
  method: 'POST'
})

// 首页-任务计划列表
const getMarketingPlan = commonRequest({
  url: 'marketing-plan/list',
  method: 'POST'
})

// 任务计划-子任务列表
const getMarketingPlanChildren = commonRequest({
  url: 'marketing-plan/task-list',
  method: 'POST'
})

// 任务计划-子任务详情
const getMarketingTaskDetail = commonRequest({
  url: 'marketing-plan/task-detail',
  method: 'POST'
})
// 任务 顾客列表
const getCustomerApi = commonRequest({
  url: 'marketing-plan/sa-customer-list',
  method: 'POST'
})

// 任务计划-子任务-完成状态
const marketingPlanSetStatus = commonRequest({
  url: 'marketing-plan/set-status',
  method: 'POST'
})

// 获取视频id
const marketingPlanGetVideo = commonRequest({
  url: 'marketing-plan/getVideoMedia',
  method: 'POST'
})

// 新增人员任务详情-群任务完成状态
const marketinhPlanUpdateStatus = commonRequest({
  url: 'marketing-plan/update-status',
  method: 'POST'
})

/** 获取导购列表 */
const staffListsApi = commonRequest({
  url: 'member/staff-lists',
  method: 'POST'
})

/** 获取门店列表 */
const storeListsApi = commonRequest({
  url: 'member/store-lists',
  method: 'POST'
})

/** 播种计划列表 */
const sowListApi = commonRequest({
  url: 'customer/sow-task/list',
  method: 'POST'
})

/** 播种计划完成 */
const sowFinishedApi = commonRequest({
  url: 'customer/sow-task/finished',
  method: 'POST'
})

/** 播种计划完成 */
const orderListApi = commonRequest({
  url: 'reserve/order-list',
  method: 'POST'
})

/** 确认用户到店购买，设置新客老客 */
const changeStatusApi = commonRequest({
  url: 'reserve/change-status',
  method: 'POST'
})

/** 日历接口（计划任务，播种，邀约） */
const monthListApi = commonRequest({
  url: 'marketing-plan/month-list',
  method: 'POST'
})

/** 筛选条件 */
const choiceOptionsApi = commonRequest({
  url: 'customer/get-choice-option',
  method: 'POST'
})

/** 筛选条件 */
const planTipsApi = commonRequest({
  url: 'marketing-plan/plan-tips',
  method: 'POST'
})

/** 判断是否外部联系人 */
const userMemberRelationApi = commonRequest({
  url: 'customer/get-user-member-relation',
  method: 'POST'
})
/** 标记客户沟通状态-发消息 */
const reportSendMsgApi = commonRequest({
  url: 'marketing-plan/confirm-contact',
  method: 'POST'
})

/** 标记客户沟通状态-打电话 */
const reportCallApi = commonRequest({
  url: 'marketing-plan/confirm-call',
  method: 'POST'
})
/** 标记客户沟通状态-加好友 */
const reportAddFriendApi = commonRequest({
  url: 'marketing-plan/confirm-friend',
  method: 'POST'
})

/** 标记客户 */
const tagCustomerApi = commonRequest({
  url: 'marketing-plan/confirm-finish',
  method: 'POST'
})

/** 标记客户说明 */
const updateRemarkApi = commonRequest({
  url: 'marketing-plan/set-remark',
  method: 'POST'
})
/** 任务计划-转发给顾客操作 上报 */
const reportShareApi = commonRequest({
  url: 'marketing-plan/set-forward',
  method: 'POST'
})
// todo: 任务未触达提醒删除
const marketingPlanDelete = commonRequest({
  url: 'marketing-plan/close_touch_tips',
  method: 'POST'
})
// todo: 子任务未触达提醒删除
const marketingPlanDeleteChild = commonRequest({
  url: 'marketing-plan/close_task_touch_tips',
  method: 'POST'
})
// 获取沟通顾客列表
const searchCustomerListApi = commonRequest({
  url: 'marketing-plan/search-customer-list',
  method: 'POST'
})
// 获取顾客列表
const searchCustomerBetaListApi = commonRequest({
  url: 'marketing-plan/search-external-list',
  method: 'POST'
})
// 获取顾客openid
const searchCustomerOpenidApi = commonRequest({
  url: 'marketing-plan/search-external-all-openid',
  method: 'POST'
})
// 移动顾客
const moveCustomerApi = commonRequest({
  url: 'marketing-plan/customer-move',
  method: 'POST'
})
// 完成
const sendFinishApi = commonRequest({
  url: 'marketing-plan/send-finish',
  method: 'POST'
})

export {
  marketingPlanDeleteChild,
  marketingPlanDelete,
  getMarketingPlan,
  getMarketingPlanChildren,
  getMarketingTaskDetail,
  getCustomerApi,
  marketingPlanSetStatus,
  marketingPlanGetVideo,
  marketinhPlanUpdateStatus,
  staffListsApi,
  storeListsApi,
  sowListApi,
  sowFinishedApi,
  orderListApi,
  changeStatusApi,
  monthListApi,
  choiceOptionsApi,
  planDetailApi,
  planListApi,
  planTipsApi,
  leaveGroupApi,
  userMemberRelationApi,
  reportSendMsgApi,
  reportCallApi,
  reportAddFriendApi,
  updateRemarkApi,
  tagCustomerApi,
  reportShareApi,
  searchCustomerListApi,
  moveCustomerApi,
  sendFinishApi,
  searchCustomerBetaListApi,
  searchCustomerOpenidApi
}
