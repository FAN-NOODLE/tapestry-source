
import { toDate, throttle } from '@/utils/tools'
// 创建页面实例对象
import { planListApi, marketingPlanDelete } from '@/api/plan-task-api'

/**
 * tabs
 * type 神策数据上报
 */
const TAB_LIST = [
  {label: '待完成', value: 0, type: 1.1},
  {label: '已完成', value: 1, type: 1.2},
  {label: '已过期', value: 2, type: 1.3},
  {label: '全部', value: -1, type: 1.4}
]

Component({
  properties: {
  },
  /** 页面的初始数据 */
  data: {
    /** tabs 切换数据 */
    tabList: TAB_LIST,
    /** tabs 选中TAB_LIST 中value */
    tabIndex: 0,
    /** 对应日历状态 */
    calendarStatus: 2,
    /** 任务名称搜索 */
    searchValue: '',
    /** 开始时间-日历中选中日期 */
    start_time: 0,
    /** 任务列表 */
    planList: []
  },
  /** 生命周期函数--监听组件显示 */
  pageLifetimes: {
    show() {
      // todo:读取当前选中日期
      const start_time = wx.getStorageSync('start_time')
      console.log('start_time', start_time)
      this.setData({
        start_time: start_time || 0
      })
    },
    hide() {
      // 删除选中日期
      wx.removeStorageSync('start_time')
    }
  },
  /** 生命周期函数--监听组件显示 */
  ready() {
    // 获取组件实例
    this.calendarCom = this.selectComponent('#calendar')
  },
  methods: {
    /** 刷新日历 */
    refreshCalendar() {
      this.calendarCom.init(true)
    },

    /** tab 切换 */
    tabChange(e) {
      // 默认全部
      let calendarStatus = 0
      // 日历中状态 1 已完成 2 待完成 3 不做展示 其他数字展示全部
      // 主任务状态 0为待完成,1为已完成,2为未完成
      switch (e.detail) {
        case 0:
          // 任务中心状态1（待完成）对应日历中的2
          calendarStatus = 2
          break
        case 1:
          // 任务中心状态1（已完成）对应日历中的1
          calendarStatus = 1
          break
        case 2:
          // 任务中心状态3（未完成）对应日历中的其他
          calendarStatus = 3
          break
        default:
          calendarStatus = 0
      }
      this.setData({
        tabIndex: e.detail,
        calendarStatus
      })
      this.getPlanTaskList(true)
    },

    /** 搜索 */
    handleSearch(e) {
      this.setData({
        searchValue: e.detail.value
      })
      this.isSearch = true
      this.getPlanTaskList()
    },

    /** 日期变化 */
    changeDay(e) {
      console.log('changeDay', e)
      const timeNum = e.detail
      this.setData({
        start_time: timeNum
      })
      this.getPlanTaskList()
    },
    // todo:隐藏planlist中某一个提醒
    hideRemind: throttle(function(e) {
      console.log('hideRemind', e)
      const { index, id } = e.currentTarget.dataset
      marketingPlanDelete({ id: id }).then(res => {
        if (res.code === 200) {
          this.setData({
            [`planList[${index}].close_tips`]: 1
          })
        } else {
          wx.showToast(res.message)
        }
      }).catch(err => {
        console.error('hideRemind 异常', err)
      })
    }, 2000),
    /** 查看任务 */
    viewTask(e) {
      // todo:跳转到任务详情
      console.log('viewTask', e)
      const { id, status } = e.currentTarget.dataset.item
      // 日历选中时间
      const { start_time } = this.data
      wx.navigateTo({ url: `/packages/task/pages/task-sub/task-sub?id=${id}&status=${status}&start_time=${start_time}` })
      // 神策数据上报 -- 查看任务
      // wx.track.trackTaskView({
      //   task_id: id,
      //   task_status: status
      // })
    },

    /**
     * 获取计划列表
     */
    getPlanTaskList() {
      wx.showLoading({
        title: '加载中'
      })
      const { searchValue, start_time, tabIndex } = this.data
      const params = {
        status: tabIndex,
        name: searchValue,
        start_time,
        // 后端是更距时间段来查询，86399 = 59:59
        end_time: start_time + 86399
      }
      planListApi(params).then(res => {
        wx.hideLoading()
        // 判断请求是否成功
        if (res.code !== 200) return wx.showToast(res.message)

        // 重新组合数据
        const taskList = res.data.Lists ? res.data.Lists.map(item => {
          return {
            ...item,
            start: toDate(item.start_time, 'yyyy-MM-dd hh:mm:ss').nowTime,
            end: toDate(item.end_time, 'yyyy-MM-dd hh:mm:ss').nowTime
          }
        }) : []

        this.setData({
          planList: taskList
        })
        if (this.isSearch && searchValue) {
          // 神策数据上报 -- 顾客列表搜索
          wx.track.trackSearch({
            current_page: 0,
            search_type: 0,
            search_key_word: searchValue,
            result_number: taskList.length
          })
          this.isSearch = false
        }
      }).catch(err => {
        wx.hideLoading()
        console.error('getPlanTaskList 异常', err)
      })
    }
  }
})

