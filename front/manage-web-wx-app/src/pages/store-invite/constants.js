const StatusList = [
  {
    id: 0,
    status_name: '进行中',
    desc: '用户未打开邀约链接或用户已打开链接，但距邀约单发出时间（不含当天）未超过30天',
    status_bg: 'https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/inprogress_bg.png',
    detail_status_bg: 'https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/d_inprogress_bg.png'
  },
  {
    id: 1,
    status_name: '邀约成功',
    // desc: '用户到店且有关联订单，若用户在邀约完成阶段未评价，SA在该阶段仍可发送邀请评价链接',
    desc: '邀约成功：邀约单发出（不含当天）30天内有关联成功订单',
    status_bg: 'https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/success_bg.png',
    detail_status_bg: 'https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/d_success_bg.png'
  },
  {
    id: 2,
    status_name: '邀约完成',
    desc: '用户到店，且导购手动确认用户到店',
    status_bg: 'https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/complete_bg.png',
    detail_status_bg: 'https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/d_complete_bg.png'
  },
  {
    id: 3,
    status_name: '已取消',
    desc: 'SA/SM手动取消预约单',
    status_bg: 'https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/cancel_overdue_bg.png',
    detail_status_bg: 'https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/d_cancel_overdue_bg.png'
  },
  {
    id: 4,
    status_name: '已过期',
    desc: '邀约发出时间超过30天，且未匹配到订单',
    status_bg: 'https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/cancel_overdue_bg.png',
    detail_status_bg: 'https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/d_cancel_overdue_bg.png'
  }
]
/** 来源 */
const sourceList = [
  {id: 1, status_name: '到店邀约'},
  {id: 2, status_name: '邀请函'}
]
export { StatusList, sourceList }
