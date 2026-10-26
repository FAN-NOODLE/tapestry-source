import { debounce } from '@/utils/tools'
import { getCurrentDate, getYesterDays, getSevenDays, getThirtyDays } from '@/utils/get-time-period'
import { customerRankApi, orderRankApi, distributionRankApi, inviteRankApi } from '../../api/compass-api'

const { globalData } = getApp()
Page({
  // 页面的初始数据
  data: {
    isIpx: globalData.isIpx,
    currentTab: 1, // 当前tab
    tabs: [
      { text: '顾客数据', type: 1 },
      { text: '分销数据', type: 2 },
      { text: '订单数据', type: 3 }
    ],
    currentCustomerTab: 1, // 当前筛选的顾客数据tab
    customerTabs: [
      {
        id: 1,
        text: '顾客总数'
      },
      {
        id: 2,
        text: '新增顾客数'
      },
      {
        id: 8,
        text: '邀约榜'
      }
    ],
    currentDistributionTab: 1, // 当前筛选的分销数据tab
    distributionTabs: [
      {
        id: 1,
        text: '分销结算总金额'
      },
      {
        id: 2,
        text: '分销结算笔数'
      }
    ],
    currentOrderTab: 1, // 当前筛选的订单数据tab
    orderTabs: [
      {
        id: 1,
        text: '推广订单数'
      },
      {
        id: 2,
        text: '推广订单实收金额'
      },
      {
        id: 3,
        text: '服务订单数'
      },
      {
        id: 4,
        text: '服务订单实收金额'
      }
    ],
    currentDate: 5, // 当前选择时间tab, 1: 近七天; 2: 近30天; 3: 自定义时间
    startTime: getCurrentDate()[0].s_timestamp, // 开始时间
    endTime: getCurrentDate()[1].s_timestamp, // 结束时间
    customerList: [],
    distributionList: [],
    orderList: [],
    rankList: []
  },

  // 生命周期函数--监听页面加载
  onLoad (options) {
    this.setData({
      currentTab: parseInt(options.type) || 1
    })
    this.requestSwitch()
  },
  onShow () {
    // 选择时间回调
    if (globalData.selectDate && globalData.selectDate.length > 0) {
      this.setData({
        startTime: globalData.selectDate[0].s_timestamp,
        endTime: globalData.selectDate[1].s_timestamp,
        startTimeText: globalData.selectDate[0].time,
        endTimeText: globalData.selectDate[1].time
      })
      this.requestSwitch()
    } else {
      this.setData({
        currentDate: 5,
        startTime: getCurrentDate()[0].s_timestamp,
        endTime: getCurrentDate()[1].s_timestamp,
        startTimeText: '',
        endTimeText: ''
      })
      this.requestSwitch()
    }
  },
  onUnload () {
    globalData.selectDate = ''
  },
  getRankList () {
    let param = {
      start_time: this.data.currentCustomerTab !== 1 ? this.data.startTime : 0,
      end_time: this.data.currentCustomerTab !== 1 ? this.data.endTime : getCurrentDate()[1].s_timestamp,
      size: 50,
      page: 1
    }
    inviteRankApi(param).then(res => {
      if (res.code === 200) {
        this.setData({
          ready: true,
          rankList: res.data
        })
      } else {
        this.setData({
          ready: true
        })
        wx.showToast(res.message)
      }
    }).catch(err => {
      this.setData({
        ready: true
      })
      wx.showToast(err.message)
    })
  },
  // 获取顾客排行榜
  customerRankList () {
    if (this.data.currentCustomerTab === 8) {
      this.getRankList()
      return
    }
    let param = {
      start_time: this.data.currentCustomerTab !== 1 ? this.data.startTime : 0,
      end_time: this.data.currentCustomerTab !== 1 ? this.data.endTime : getCurrentDate()[1].s_timestamp,
      type: 1
    }
    customerRankApi(param).then(res => {
      if (res.code === 200) {
        this.setData({
          ready: true,
          customerList: res.data.list
        })
      } else {
        this.setData({
          ready: true
        })
        wx.showToast(res.message)
      }
    }).catch(err => {
      this.setData({
        ready: true
      })
      wx.showToast(err.message)
    })
  },
  // 获取分销排行榜
  distributionRankList () {
    let param = {
      start_time: this.data.startTime,
      end_time: this.data.endTime,
      type: this.data.currentDistributionTab
    }
    distributionRankApi(param).then(res => {
      if (res.code === 200) {
        this.setData({
          ready: true,
          distributionList: res.data.list
        })
      } else {
        this.setData({
          ready: true
        })
        wx.showToast(res.message)
      }
    }).catch(err => {
      this.setData({
        ready: true
      })
      wx.showToast(err.message)
    })
  },
  // 获取订单排行榜
  orderRankList () {
    let param = {
      start_time: this.data.startTime,
      end_time: this.data.endTime,
      type: this.data.currentOrderTab
    }
    orderRankApi(param).then(res => {
      if (res.code === 200) {
        this.setData({
          ready: true,
          orderList: res.data.list
        })
      } else {
        this.setData({
          ready: true
        })
        wx.showToast(res.message)
      }
    }).catch(err => {
      this.setData({
        ready: true
      })
      wx.showToast(err.message)
    })
  },
  // 请求接口的逻辑列表
  requestSwitch () {
    switch (this.data.currentTab) {
      case 1:
        this.customerRankList()
        break
      case 2:
        this.distributionRankList()
        break
      case 3:
        this.orderRankList()
        break
    }
  },
  inviteList () {
    console.log('邀约排行榜')
  },
  // 顶部tab切换
  changeTab: debounce(function (e) {
    const { type } = e.detail
    this.setData({
      ready: false,
      currentDate: 5,
      startTime: getCurrentDate()[0].s_timestamp,
      endTime: getCurrentDate()[1].s_timestamp,
      startTimeText: '',
      endTimeText: '',
      customerList: [],
      orderList: [],
      distributionList: [],
      currentCustomerTab: 1,
      currentDistributionTab: 1,
      currentOrderTab: 1,
      currentTab: type
    })
    this.requestSwitch()
  }, 300),
  // 顾客tab切换
  changeCustomerTab: debounce(function (e) {
    const tab = e.currentTarget.dataset.item
    if (tab.id === this.data.currentCustomerTab) return
    this.setData({
      ready: false,
      currentCustomerTab: tab.id
    })
    this.customerRankList()
  }, 300),
  // 分销tab切换
  changeTaskTab: debounce(function (e) {
    const tab = e.currentTarget.dataset.item
    if (tab.id === this.data.currentDistributionTab) return
    this.setData({
      ready: false,
      currentDistributionTab: tab.id
    })
    this.distributionRankList()
  }, 300),
  // 订单tab切换
  changeOrderTab: debounce(function (e) {
    const tab = e.currentTarget.dataset.item
    if (tab.id === this.data.currentOrderTab) return
    this.setData({
      ready: false,
      currentOrderTab: tab.id
    })
    this.orderRankList()
  }, 300),
  // 选择时间
  selectDate (e) {
    const type = parseInt(e.currentTarget.dataset.type)
    this.setData({
      currentDate: type
    })
    switch (type) {
      case 1:
        // 选择近七天
        this.setData({
          startTimeText: '',
          endTimeText: '',
          startTime: getSevenDays()[0].s_timestamp,
          endTime: getSevenDays()[1].s_timestamp
        })
        break
      case 2:
        // 选择近30天
        this.setData({
          startTimeText: '',
          endTimeText: '',
          startTime: getThirtyDays()[0].s_timestamp,
          endTime: getThirtyDays()[1].s_timestamp
        })
        break
      case 3:
        // 选择自定义时间
        let url = '/pages/date-guide-store-select/date-guide-store-select?type=3'
        if (this.data.startTimeText) {
          url += `&startTime=${this.data.startTime}&endTime=${this.data.endTime}`
        }
        wx.navigateTo({
          url: url
        })
        break
      case 4:
        // 昨天
        this.setData({
          startTimeText: '',
          endTimeText: '',
          startTime: getYesterDays()[0].s_timestamp,
          endTime: getYesterDays()[1].s_timestamp
        })
        break
      case 5:
        // 今天
        this.setData({
          startTimeText: '',
          endTimeText: '',
          startTime: getCurrentDate()[0].s_timestamp,
          endTime: getCurrentDate()[1].s_timestamp
        })
        break
    }
    if (type !== 3) {
      this.requestSwitch()
    }
  },
  // 去数据概览页面
  goCompass () {
    wx.switchTab({
      url: `/pages/compass-data/compass-data?type=${this.data.currentTab}`
    })
  }
})
