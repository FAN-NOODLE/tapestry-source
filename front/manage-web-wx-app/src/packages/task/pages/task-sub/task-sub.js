
import { toDate } from '@/utils/tools'
import { getMarketingPlanChildren, planDetailApi, marketingPlanDeleteChild } from '@/api/plan-task-api'

Page({
  data: {
    /** 主任务数据 */
    mainTaskInfo: {},
    /** 子任务数据 */
    childTaskInfo: {},
    /** 日历默认选中时间 */
    calendarTime: 0,
    /** 路由参数 */
    routerParams: {},
    /** 标注是否从返回跳转的 */
    isBack: false,
    showTips: false,
    showLoading: false
  },

  /**  生命周期函数--监听页面加载 */
  onLoad(option) {
    // todo：判断是否从返回跳转的
    let start_time = wx.getStorageSync('start_time')
    if (start_time) {
      start_time = Number(start_time)
    } else {
      start_time = Number(option.start_time)
      // 存储当前选中日期
      wx.setStorageSync('start_time', start_time)
    }
    this.setData({
      option,
      routerParams: {
        ...option,
        id: Number(option.id),
        start_time: start_time
      }
    })
  },

  /** 生命周期函数--监听页面显示 */
  onShow() {
    // 获取组件实例
    this.calendarCom = this.selectComponent('#calendar')
    // 如果是返会，刷新数据
    // if (this.data.isBack && this.data.continued) {
    //   this.refreshCalendar()
    // }

    this.getMainDetail(this.data.routerParams.start_time)
  },

  /** 刷新日历 */
  refreshCalendar() {
    this.calendarCom.init(true)
  },

  /** 日期变化 */
  changeDay(e) {
    this.setData({
      changeTime: e.detail
    })
    // 需要存储当前选中日期
    wx.setStorageSync('start_time', e.detail)
    this.getPlanTaskList(e.detail)
  },

  /** 获取主任务详情 */
  getMainDetail(start_time) {
    //  增加loading
    this.setData({
      showLoading: true
    })
    const { routerParams } = this.data
    const params = {
      id: routerParams.id,
      start_time,
      // 后端是更距时间段来查询，86399 = 59:59
      end_time: start_time + 86399
    }
    planDetailApi(params)
      .then(res => {
        // 判断请求是否成功
        if (res.code !== 200) return wx.showToast(res.message)
        const mainTaskInfo = {
          ...res.data,
          start: toDate(res.data.start_time, 'yyyy-MM-dd hh:mm:ss').nowTime,
          end: toDate(res.data.end_time, 'yyyy-MM-dd hh:mm:ss').nowTime
        }
        this.getPlanTaskList(routerParams.start_time)
        this.setData({
          mainTaskInfo,
          calendarTime: routerParams.start_time
        })
      })
      .catch(e => console.error('获取任务详情失败', e))
  },

  /** 获取子任务列表 */
  getPlanTaskList(time) {
    const { routerParams } = this.data
    // 组合查询参数，无分页
    const params = {
      start_time: time,
      page: 1,
      size: 1000,
      id: routerParams.id
    }

    getMarketingPlanChildren(params)
      .then(res => {
        // 判断请求是否成功
        if (res.code !== 200) return wx.showToast(res.message)
        res.data.list = res.data.list.map(item => {
          item.task_start_time_text = toDate(item.task_start_time, 'yyyy-MM-dd hh:mm:ss').nowTime
          item.task_end_time_text = toDate(item.task_end_time, 'yyyy-MM-dd hh:mm:ss').nowTime
          return item
        })
        this.setData({
          childTaskInfo: res.data,
          isBack: false,
          showLoading: false
        })
      }).catch(err => {
        console.log(err)
        this.setData({
          showLoading: false
        })
      })
  },

  //  显示提示
  showPopover() {
    this.setData({
      showTips: true
    })
  }
})
