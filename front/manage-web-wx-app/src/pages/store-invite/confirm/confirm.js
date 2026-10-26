// 获取全局应用程序实例对象

import { guideTimeApi, orderInfoApi, serviceListApi, cancelOrderApi, getConfigMessageApi, commitOrderApi, messageSuccessApi } from '@/api/store-invite'

const { globalData } = getApp()

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    // wxwork: globalData.system.environment === 'wxwork',
    wxwork: false,
    dateModel: false,
    serverIndex: -1, // 服务选择的下标
    serverList: [],
    timeObj: {},
    showTime: false,
    id: 0,
    enterpriseId: 0,
    reserveNumber: 0,
    orderInfo: {},
    selectTime: {},
    selectServer: {},
    userRemark: '',
    messageSuccess: false,
    isCheck: true,
    showModel: false,
    isSingle: false,
    modelTitle: '提醒',
    description: '该时段已有顾客预约服务，',
    descriptionTwo: '建议重新选择时',
    statusList: [
      {
        id: 0,
        status_name: '进行中'
      },
      {
        id: 1,
        status_name: '已完成'
      },
      {
        id: 2,
        status_name: '待完成'
      },
      {
        id: 3,
        status_name: '已取消'
      },
      {
        id: 4,
        status_name: '已过期'
      }
    ]

  },

  // 生命周期函数--监听页面加载
  async onLoad (options) {
    this.setData({
      wxwork: globalData.system.environment === 'wxwork'
    })
    this.setData({
      id: options.orderId - 0 || 0,
      enterpriseId: options.id || 0,
      reserveNumber: options.num || 0
    })
    // 导购分享直接进入详情页面
    if (!(options.corp_id && this.data.wxwork)) {
      // 获取服务列表
      await this.getServerList()
      this.getConfigMessage()
      // 获取预约单详情
      this.getDetail()
      wx.hideShareMenu({
        withShareTicket: true,
        menus: ['shareAppMessage', 'shareTimeline']
      })
      wx.hideHomeButton()
    }
  },
  onShow () {
    let options = globalData.appOnShowoption
    if (options.query.corp_id && this.data.wxwork) {
      wx.reLaunch({
        url: `/pages/store-invite/invite-detail/invite-detail?id=${options.query.orderId}&num=${options.query.num}&corp_id=${options.query.corp_id}`
      })
    }
  },

  // 获取订阅消息配置
  getConfigMessage () {
    getConfigMessageApi().then(res => {
      console.log('获取订阅消息配置res=====》', res)
      this.data.tmplIds = res.data['confirm-subscribe']
    }).catch(err => {
      console.error(err, '====')
    })
  },
  // 是否有改变
  hasCheck () {
    const { orderInfo, selectTime, userRemark, selectServer } = this.data
    let isCheck = false
    if (
      orderInfo.reserve_id !== selectServer.id ||
      selectTime.startTime !== orderInfo.reserve_time ||
      userRemark
    ) {
      isCheck = true
    }
    this.setData({
      isCheck
    })
  },
  // 获取服务列表
  async getServerList () {
    const params = {
      page: 1,
      size: 9999,
      enterprise_id: this.data.enterpriseId - 0
    }
    serviceListApi(params).then(
      res => {
        console.log('获取服务列表res===>', res)
        this.setData({
          serverList: res.data.list
        })
      }
    ).catch(err => {
      console.error(err, '错误❎')
    })
  },
  // 获取用户详情
  getDetail () {
    console.log('confirm页面开始获取页面详情====getDetail===>')
    const { id, reserveNumber } = this.data
    const params = { id, reserve_number: reserveNumber }
    orderInfoApi(params).then(
      res => {
        console.log('用户详情获取，返回数据=====>', res)
        const orderInfo = res.data
        this.setData({
          userRemark: orderInfo.user_remark
        })
        this.getTime(orderInfo.member_id)
        const serverIndex = this.data.serverList.findIndex(item => {
          return item.id === orderInfo.reserve_id
        })
        let selectServer = {}
        // 兼容详情服务id被删之后的展示与编辑
        if (serverIndex === -1) {
          selectServer = {
            id: orderInfo.reserve_id,
            reserve_img: orderInfo.reserve_img || '',
            reserve_name: orderInfo.serve_name || orderInfo.reserve_name || '',
            reserve_time: orderInfo.serve_time
          }
        } else {
          selectServer = this.data.serverList[serverIndex]
        }
        orderInfo.serverIndex = serverIndex
        const date = new Date(orderInfo.reserve_time * 1000)
        const select_name = `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`
        const selectTime = {
          select_name,
          inviteNum: 0,
          isFinish: false,
          name: `${date.getHours()}:00 - ${date.getHours() + 1}:00`,
          startTime: orderInfo.reserve_time
        }
        console.log(date, '错误✅')
        this.setData({
          orderInfo,
          serverIndex,
          selectServer,
          selectTime,
          hasCheck: false
        })
      }
    ).catch(err => {
      console.error(err, '错误❎')
    })
  },
  cancelOrder () {
    wx.showLoading({
      title: '取消中'
    })
    const { id, reserveNumber } = this.data
    const params = { id, reserve_number: reserveNumber }
    cancelOrderApi(params).then(
      res => {
        wx.hideLoading()
        if (res.code !== 200) {
          wx.showToast(res.message)
          return
        }
        this.setData({
          'orderInfo.status': 3
        })
        wx.showToast('取消成功')
      }
    ).catch(err => {
      wx.hideLoading()
      console.error(err, '错误❎')
    })
  },
  confirmData () {
    // if (!this.data.isCheck) return
    try {
      const { tmplIds } = this.data
      const self = this
      wx.requestSubscribeMessage({
        tmplIds,
        success (res) {
          self.data.messageSuccess = true
          self.confirmOrder()
        },
        fail (err) {
          console.log(err, '----')
          self.data.messageSuccess = false
          self.confirmOrder()
        }
      })
    } catch (error) {
      console.log(error, '999')
    }
  },
  confirmOrder () {
    wx.showLoading({
      title: '确认中'
    })
    const { id, enterpriseId, reserveNumber, selectTime, selectServer, userRemark, tmplIds } = this.data
    const params = {
      id,
      user_remark: userRemark,
      reserve_number: reserveNumber,
      reserve_id: selectServer.id,
      reserve_time: selectTime.startTime,
      status: 1
    }
    commitOrderApi(params).then(
      res => {
        wx.hideLoading()
        if (res.code !== 200) {
          wx.showToast(res.message)
          return
        }
        this.getDetail()
        wx.showToast('确认成功')
        if (this.data.messageSuccess) {
          messageSuccessApi({
            rel_id: id,
            tmpl_list: tmplIds.map(item => {
              return {
                pri_tmpl_id: item,
                url: `pages/store-invite/confirm/confirm?id=${enterpriseId}&orderId=${id}&num=${reserveNumber}`
              }
            })
          })
        }
      }
    ).catch(err => {
      wx.hideLoading()
      console.error(err, '错误❎')
    })
    console.log(params)
  },
  changeInput (e) {
    const key = e.currentTarget.dataset.key
    const value = e.detail.value
    this.setData({
      [key]: value
    })
    this.hasCheck()
  },
  selectGuide () {
    this.data.goGuide = 1
    const a = [] // 当前预约单门店的id
    const b = [] // 当前预约导购
    wx.navigateTo({
      url: `/pages/date-guide-store-select/date-guide-store-select?type=2&dpIds=${JSON.stringify(a)}&hasChooseGuideIds=${JSON.stringify(b)}&role=3&isBa=${false}&mautiple=${false}`
    })
  },
  getTime (member_id) {
    guideTimeApi({ member_id }).then(res => {
      this.setData({
        timeObj: res.data,
        showTime: true
      })
    }).catch(err => {
      console.error(err, '错误')
    })
  },
  // chooseTime () {
  // const { is_commit, status } = this.data.orderInfo
  // if (is_commit === 0 && status !== 3) {
  //   this.setData({
  //     dateModel: true
  //   })
  // }
  // },
  closeModel () {
    this.setData({
      dateModel: false
    })
  },
  confirmModel () {
    this.setData({
      showModel: false
    })
  },
  sumbit (e) {
    if (e.detail.inviteNum !== 0) {
      this.setData({
        selectTime: e.detail,
        showModel: true,
        isSingle: true,
        dateModel: false,
        modelTitle: '提醒',
        description: '该时段已有顾客预约服务，',
        descriptionTwo: '建议重新选择时'
      })
    }
    this.setData({
      dateModel: false,
      selectTime: e.detail
    })
    this.hasCheck()
  },
  // changeServer (e) {
  //   this.setData({
  //     serverIndex: Number(e.detail.value),
  //     selectServer: this.data.serverList[Number(e.detail.value)]
  //   })
  //   this.hasCheck()
  // },
  // 页面分享设置
  onShareAppMessage () {
    return {
      title: 'share title',
      path: 'pages/store-invite/create/create'
    }
  }
})
