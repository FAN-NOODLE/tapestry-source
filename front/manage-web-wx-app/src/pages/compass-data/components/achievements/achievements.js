import {
  storeListsApi,
  performanceApi,
  fiscalMonthApi,
  activeCustomersApi,
  invitationApi
} from '@/api/compass-api'
import { toDate } from '@/utils/tools'
import { getCurrentDate } from '@/utils/get-time-period'

const { globalData } = getApp()
const closeIcon = require('../../../../images/invite-store/close-icon.png')
/** 快速排序 - 时间 */
function sortTime (arr, attribute) {
  return arr.sort((a, b) => { return a[attribute] - b[attribute] })
}


Component({
  properties: {
    /** 自定义时间 */
    dateVals: {
      type: Array,
      value: []
    },
    /** 用户身份 */
    role: {
      type: Number
    }
  },
  observers: {
    dateVals () {
      if (this.data.dateVals.length) {
        console.log(666666666666, this.data.dateVals)
        this.setData({
          performanceLists: [],
          hasMore: true,
          'meta.page': 1
        })
        this.getData(this.data.tabValue)
      }
    }
  },
  data: {
    closeIcon: closeIcon,
    /** tab渲染数据 */
    tabList: [
      {
        label: '明星榜',
        value: 1
      },
      {
        label: '邀约榜',
        value: 2
      },
      {
        label: '顾客榜',
        value: 3
      },
      {
        label: '活跃榜',
        value: 4
      }
    ],
    /** 当前选中的tab项  */
    tabValue: 1,
    /** 明星榜数据 */
    performanceList: [],
    /** 店铺下拉数据 */
    selectPointList: [],
    // 选中的门店
    selectPicker: 0,
    /** 时间维度状态  */
    dateType: 'month',
    /** 邀约榜单更新时间 */
    updateDate: '',
    /** 除了邀约榜单更新时间 */
    updateDate2: '',
    /** 排序提示 */
    sortTitle: '',
    /** 排序提示弹窗状态 */
    isSort: false,
    /** 分页数据 */
    meta: {
      page: 1, // 页数
      size: 10 // 每页的数据
      // total: 1, // 总数
    },
    /** 是否有更多数据 */
    hasMore: true
  },
  ready () {
    /** 页面展示统计截止时间  */
    const updateDate = toDate(
      getCurrentDate()[0].s_timestamp - 86400 * 2,
      'yyyy年MM月dd日'
    ).nowTime
    const updateDate2 = toDate(
      getCurrentDate()[0].s_timestamp - 86400,
      'yyyy年MM月dd日'
    ).nowTime
    this.setData({
      updateDate,
      updateDate2
    })
    this.getStoreLists()
  },
  methods: {
    /** 切换时间维度类型
     * @param {Number} e.currentTarget.dataset.type 当前时间维度唯一值  month:月| year:年
     */
    switchRule (e) {
      let dateType = e.currentTarget.dataset.type
      // 和当前时间维度一致 return出去
      if (dateType === this.data.dateType) {
        return
      }
      if (dateType === 'month') {
        this.setData({
          dateType,
          dateIndex: this.data.dateMonthIndex
        })
      } else if (dateType === 'year') {
        const current = this.data.yearList.findIndex(item => item.end_time >= parseInt(Date.now() / 1000))
        this.setData({
          dateType,
          dateIndex: current === -1 ? this.data.yearList.length - 1 : current
        })
      }
      /** 切换的时候若自定义时间存在 清除掉后掉借口 */
      if (this.data.dateVals.length) {
        this.deleteTime()
      } else {
        this.setData({
          performanceLists: []
        })
        this.data.meta.page = 1
        this.getData(this.data.tabValue)
      }
    },
    /** 切换时间维度
     *  @param {Number} e.currentTarget.dataset.type 当前切换时间左右的唯一值  left:左边| right:右边
     */
    switchDate (e) {
      let type = e.currentTarget.dataset.type
      if (type === 'left') {
        if (this.data.dateIndex === 0) {
          return
        }
        this.setData({
          dateIndex: this.data.dateIndex - 1
        })
      } else if (type === 'right') {
        if (this.data.dateType === 'year') {
          if (this.data.dateIndex + 1 === this.data.yearList.length) {
            return
          }
          this.setData({
            dateIndex: this.data.dateIndex + 1
          })
        }
        if (this.data.dateType === 'month') {
          if (this.data.dateIndex + 1 === this.data.monthList.length) {
            return
          }
          this.setData({
            dateIndex: this.data.dateIndex + 1
          })
        }
      }
      this.setData({
        performanceLists: []
      })
      this.data.meta.page = 1
      this.getData(this.data.tabValue)
    },
    /** 删除自定义时间 */
    deleteTime () {
      this.setData({
        dateVals: [],
        performanceLists: []
      })
      this.data.meta.page = 1
      globalData.selectDate1 = null
      this.getData(this.data.tabValue)
    },
    /** 自定义时间 */
    getDateSelect () {
      let showDate = true
      const { monthList } = this.data
      let min = 0
      let max = 0

      if (monthList.length) {
        min = monthList[0].start_time
        max = monthList[monthList.length - 1].end_time
      }
      /** 有选择自定义时间数据 携带用来回显 */
      if (this.data.dateVals.length) {
        wx.navigateTo({
          url: `/pages/date-guide-store-select/date-guide-store-select?type=3&startTime=${this.data.dateVals[0]}&endTime=${this.data.dateVals[1]}&showDate=${showDate}&selectDateText=selectDate1&min=${min}&max=${max}`
        })
      } else {
        console.log('走了这个几口')
        wx.navigateTo({
          url: `/pages/date-guide-store-select/date-guide-store-select?type=3&showDate=${showDate}&selectDateText=selectDate1&min=${min}&max=${max}`
        })
      }
    },
    /** tab切换
     * @param {Number} e.currentTarget.dataset.item.value 当前选中tab值
     */
    switchTab (e) {
      let v = e.currentTarget.dataset.item.value
      if (this.data.tabValue === v) {
        return
      }
      this.setData({
        tabValue: e.currentTarget.dataset.item.value,
        performanceList: {},
        performanceLists: []
      })
      this.data.hasMore = true
      this.data.meta.page = 1
      this.getData(this.data.tabValue)
    },
    /** 根据tab获取列表数据
     * @param {Number} v 当前选中tab值
     */
    getData (v) {
      switch (v) {
        case 1:
          this.getPerformance()
          break
        case 2:
          this.getInvitation()
          break
        case 3:
          this.getActiveCustomers(2)
          break
        case 4:
          this.getActiveCustomers(1)
          break

        default:
          break
      }

      // 神策数据上报 -- 查看绩效榜单
      if (this.data.meta.page === 1) {
        wx.track.trackViewRanking({
          ranking_cycle: this.data.dateType,
          ranking_type: v,
          time_range: this.getQueryTime()
        })
      }
    },
    /** 点击门店数据
     * @param {object} e
     * @param {string} e.detail.value 门店选择项的值
     */
    bindPickerChange (e) {
      this.setData({
        selectPicker: e.detail.value - 0,
        performanceLists: [],
        'meta.page': 1
      })
      console.log(e, '111')
      this.getData(this.data.tabValue)
    },
    /** 展开关闭后五名
     * @param {object} data
     * @param {string} e.currentTarget.dataset.key 当前列表唯一key
     */
    setEnable (e) {
      let key = e.currentTarget.dataset.key
      let k = 'performanceList.' + key
      let ks = 'performanceList.' + key + 'Show'
      if (this.data.performanceList[key + 'Show']) {
        this.setData({
          [k]: this.data.performanceList[key + 'All'].slice(0),
          [ks]: false
        })
      } else {
        let v = this.data.performanceList[key + 'All'].slice(0, 5)
        this.setData({
          [k]: v,
          [ks]: true
        })
      }
    },
    /** 查看更多 */
    handleToViewMore () {
      if (this.data.hasMore) {
        console.log(123)
        this.data.meta.page++
        this.getData(this.data.tabValue)
      }
    },
    /**
     * 请求时间参数整合
     * @returns arr 时间参数
     */
    getQueryTime () {
      let arr = []
      if (this.data.dateVals && this.data.dateVals.length) {
        arr = [this.data.dateVals[0], this.data.dateVals[1]]
      } else if (this.data.dateType === 'month') {
        arr = [
          this.data.monthList[this.data.dateIndex].start_time,
          this.data.monthList[this.data.dateIndex].end_time
        ]
      } else if (this.data.dateType === 'year') {
        console.log(1111111111, this.data.yearList, this.data.dateIndex)
        arr = [
          this.data.yearList[this.data.dateIndex].start_time,
          this.data.yearList[this.data.dateIndex].end_time
        ]
      }
      return arr
    },
    /** 获取明星列表 */
    getPerformance () {
      let query = {
        financial_time: this.getQueryTime()
        // third_id:
        //   this.data.selectPicker || this.data.role === 3
        //     ? this.data.selectPointList[this.data.selectPicker].store_id
        //     : ''
      }
      performanceApi(query)
        .then((res) => {
          if (res.code === 200 && Object.keys(res.data).length) {
            let { active, customer, invitation } = res.data
            this.setData({
              performanceList: {
                // 列表使用的真实数据
                active: active.length > 5 ? active.slice(0, 5) : active,
                customer: customer.length > 5 ? customer.slice(0, 5) : customer,
                invitation:
                  invitation.length > 5 ? invitation.slice(0, 5) : invitation,
                invitationShow: true,
                customerShow: true,
                activeShow: true,
                // 保存三个列表的全部数据
                activeAll: active,
                customerAll: customer,
                invitationAll: invitation
              }
            })
          } else {
            this.setData({
              performanceList: {
                invitationShow: true,
                customerShow: true,
                activeShow: true,
                activeAll: [],
                customerAll: [],
                invitationAll: [],
                active: [],
                customer: [],
                invitation: []
              }
            })
          }
        })
        .catch((err) => {
          console.log(err)
        })
      // }
    },
    /** 活跃榜/顾客榜
     * @param {Number} active_type 1:活跃榜 2:顾客榜
     */
    getActiveCustomers (active_type) {
      let query = {
        financial_time: this.getQueryTime(),
        active_type, // 1.活跃榜 2.顾客榜
        ...this.data.meta,
        third_id:
          this.data.selectPicker || this.data.role === 3
            ? this.data.selectPointList[this.data.selectPicker].store_id
            : ''
      }
      wx.showLoading({
        title: '加载中'
      })
      activeCustomersApi(query)
        .then((res) => {
          if (res.code === 200) {
            wx.hideLoading()
            if (res.data.list.length && Array.isArray(res.data.list)) {
              const list = res.data.list.map((item) => {
                return item
              })

              // 判断更多
              this.data.hasMore =
        Math.ceil(res.data.total / res.data.size) > res.data.page

              this.setData({
                performanceLists: this.data.performanceLists.concat(list),
                hasMore: this.data.hasMore
              })
            } else {
              this.setData({
                performanceLists: []
              })
            }
          }
        })
        .catch((err) => {
          console.log(err)
        })
    },
    /** 邀约榜 */
    getInvitation () {
      let query = {
        financial_time: this.getQueryTime(),
        ...this.data.meta,
        third_id:
          this.data.selectPicker || this.data.role === 3
            ? this.data.selectPointList[this.data.selectPicker].store_id
            : ''
      }
      wx.showLoading({
        title: '加载中'
      })
      invitationApi(query)
        .then((res) => {
          if (res.code === 200) {
            wx.hideLoading()
            if (res.data.list.length && Array.isArray(res.data.list)) {
              const list = res.data.list.map((item) => {
                return item
              })
              // 判断更多
              this.data.hasMore =
        Math.ceil(res.data.total / res.data.size) > res.data.page

              this.setData({
                performanceLists: this.data.performanceLists.concat(list),
                hasMore: this.data.hasMore,
                isLoading: false
              })
            } else {
              this.setData({
                performanceLists: []
              })
            }
          }
        })
        .catch((err) => {
          console.log(err)
        })
        // .finally(() => {
        //   this.setData({
        //     isLoading: false
        //   })
        // })
    },
    /** 获取店铺 */
    getStoreLists () {
      storeListsApi()
        .then((res) => {
          if (res.code === 200) {
            if (res.data.length && Array.isArray(res.data)) {
              /** 门店数据添加一项 全部选项 */
              if (this.data.role !== 3) {
                res.data.unshift({ store_id: 0, store_name: '全部门店' })
              }
              this.setData({
                selectPointList: res.data
              })
            }
          }
          this.getFiscalMonthApi()
        })
        .catch((err) => {
          console.log(err)
        })
    },

    /** 获取财年&财月数据  */
    getFiscalMonthApi () {
      fiscalMonthApi()
        .then((res) => {
          if (res.code === 200) {
            // 排序-按照年份
            res.data = res.data.map(item => ({...item, yearNum: item.year.split('-')[0]}))
            res.data = sortTime(res.data, 'yearNum')
            /** 维度分为几年 */
            let round = res.data.length / 12
            let arr = new Array(round).fill('')
            /** 根据维度分为几年 来整理财年数据 */
            let yearList = arr.map((item, index) => {
              let dif = index * 12
              return {
                start_time: res.data[0 + dif].start_time,
                end_time: res.data[11 + dif].end_time,
                label: res.data[0 + dif].year,
                now: 0
              }
            })
            /** 根据维度整理财月数据 */
            let monthList = res.data.map((e, index) => {
              /** e.now == 1 则当前维度是默认维度 */
              if (e.now === 1) {
                /** 给年维度数据 加上默认维度 */
                yearList.forEach((ev) => {
                  if (ev.label === e.year) {
                    ev.now = 1
                  }
                })
                this.setData({
                  /** 当前选中的时间维度索引 */
                  dateIndex: index,
                  /** 缓存当前选中月的维度索引 */
                  dateMonthIndex: index
                })
              }
              return {
                /** 月维度展示文字 */
                label:
                  toDate(e.start_time, 'MM月dd日').nowTime +
                  '至' +
                  toDate(e.end_time, 'MM月dd日').nowTime,
                ...e
              }
            })
            this.setData({
              /** 财月数据 */
              monthList,
              /** 财年数据 */
              yearList
            })
          }
          this.getData(this.data.tabValue)
        })
        .catch((err) => {
          this.getData(this.data.tabValue)
          console.log(err)
        })
    },
    /** 跳转邀约详情页面 */
    goRecordPage (e) {
      const {member_id, reserve_numbers} = e.currentTarget.dataset.item
      wx.navigateTo({
        url: `/pages/store-invite/record-list/record-list?member_id=${member_id}&reserve_numbers=1`,
        success: function(res) {
          // 通过eventChannel向被打开页面传送数据
          res.eventChannel.emit('accept', reserve_numbers)
        }
      })
      // 神策数据上报 -- 查看邀约记录
      wx.track.trackViewAppointmentHistory({
        entry_point: 3
      })
    },
    /** 排序提示 */
    switchSort () {
      let title = ''
      switch (this.data.tabValue) {
        case 2:
          title = '按成功金额从高到低排序'
          break
        case 3:
          title = '按新增顾客数从高到低排序'
          break
        case 4:
          title = '按聊天总数从高到低排序'
          break
        default:
          break
      }
      this.setData({
        /** 当前提示弹窗的文字 */
        sortTitle: title,
        isSort: true
      })
    },
    /** 关闭提示弹窗  */
    close () {
      this.setData({
        isSort: false
      })
    }
  }
})
