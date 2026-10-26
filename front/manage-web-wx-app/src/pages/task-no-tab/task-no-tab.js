import { planTipsApi } from '@/api/plan-task-api'

/** 生成tabs */
const generateTabs = (taskNum = 0, sowNum = 0, invitationNum = 0) => {
  return [
    {id: 0, label: '任务计划', num: taskNum},
    {id: 1, label: '播种计划', num: sowNum},
    {id: 2, label: '邀约计划', num: invitationNum}
  ]
}
const app = getApp()

Page({
  // 页面的初始数据
  data: {
    navigationBarHeight: app.globalData.system.statusBarHeight + 44 + 'px',
    /** tab数据 */
    tabs: generateTabs(),
    /** 默认选中tab */
    currentTab: 0,
    /** 标注是否从返回跳转的 */
    isBack: false,
    /** 是否离开页面 */
    isLevel: false,
    /** 消息通知进入选中播种的天 */
    sowTime: 0
  },
  onLoad (options) {
    const {type, time} = options
    // 从消息通知进入跳到播种计划的
    if (type) {
      this.setData({
        currentTab: Number(type),
        sowTime: Number(time)
      })
    }
  },
  onShow () {
    // 获取组件实例
    this.taskPlan = this.selectComponent('#taskPlan')
    this.invitationPlan = this.selectComponent('#invitationPlan')
    this.sowPlan = this.selectComponent('#sowPlan')
    if (this.data.isLevel) {
      this.refresh()
    }
    this.getPlanTips()
  },
  onHide () {
    this.setData({
      isLevel: true
    })
  },

  /** 播种计划通知需要刷新提示数量 */
  sowRefresh () {
    this.getPlanTips()
  },

  /** 当页面切换时，需要刷新接口 */
  refresh () {
    console.log('页面返回，接口刷新了。。。')
    switch (this.data.currentTab) {
      case 0:
        this.taskPlan.refreshCalendar()
        break
      case 1:
        this.sowPlan.refreshCalendar()
        break
      case 2:
        this.invitationPlan.refreshCalendar()
        break
    }
  },

  /** 页面触底上拉加载更多 */
  onReachBottom () {
    if (this.data.currentTab === 0) return
    // 组件实例
    const componentId = ['taskPlan', 'sowPlan', 'invitationPlan']
    const example = this.selectComponent(`#${componentId[this.data.currentTab]}`)
    example && example.bindDownLoad()
  },

  /** 点击切换 */
  clickTab (e) {
    // 判断重复点击
    if (this.data.currentTab === e.target.dataset.current) return false

    this.setData({
      currentTab: e.target.dataset.current
    })
  },

  /** 获取提示信息 */
  getPlanTips () {
    const start_time = new Date().getTime() / 1000
    const params = {
      start_time,
      // 后端是更距时间段来查询，86399 = 59:59
      end_time: start_time + 86399
    }
    planTipsApi(params)
      .then((res) => {
        if (res.code !== 200) return
        const {plan_nums, sow_task_nums, reserve_nums} = res.data
        const tabs = generateTabs(plan_nums, sow_task_nums, reserve_nums)

        this.setData({
          tabs,
          isBack: false,
          isLevel: false
        })
      })
      .catch(e => console.error('获取提示信息出错', e))
  }
})
