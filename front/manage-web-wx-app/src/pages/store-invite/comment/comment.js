// 获取全局应用程序实例对象
import { makeScoreApi, scoreInfoApi } from '@/api/store-invite'
import { sensorsLoginInit } from '@/lib/sensors'

const { globalData } = getApp()
// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    title: 'comment',
    storeIndex: -1,
    serverIndex: -1,
    remark: '',
    starImg: require('../../../images/invite-store/star.png'),
    activeImg: require('../../../images/invite-store/star-gray.png'),
    id: 0,
    enterpriseId: 0,
    reserveNumber: 0,
    orderInfo: {}
  },

  // 生命周期函数--监听页面加载
  onLoad (options) {
    wx.setStorageSync('commentOption', options)
    // TODO: onLoad
    this.setData({
      id: options.orderId - 0 || 0,
      enterpriseId: options.id || 0,
      reserveNumber: options.num || 0
    })
    // 获取预约单详情
    this.getDetail()
    wx.hideShareMenu({
      withShareTicket: true,
      menus: ['shareAppMessage', 'shareTimeline']
    })
    wx.hideHomeButton()
  },
  onShow () {
    let options = wx.getStorageSync('commentOption')
    let iswxwork = globalData.system.environment === 'wxwork'
    console.log('commentOption-----commentOption----', options)
    if (iswxwork) {
      wx.reLaunch({
        url: `/pages/store-invite/invite-detail/invite-detail?id=${options.orderId}&num=${options.num}&corp_id=${options.corp_id}&role=${options.role}`
      })
    }
  },

  getDetail () {
    const { id, reserveNumber } = this.data
    const params = { id, reserve_number: reserveNumber }
    scoreInfoApi(params).then(
      res => {
        if (res.code === 200) {
          const orderInfo = res.data
          this.setData({
            orderInfo,
            storeIndex: orderInfo.store_score - 1,
            serverIndex: orderInfo.service_score - 1,
            remark: orderInfo.remark
          })
          return
        }
        wx.showToast(res.message)
      }
    ).catch(err => {
      console.error(err, '错误❎')
    })
  },
  clickScore (e) {
    if (this.data.orderInfo.is_score === 1) return
    const { key, index } = e.currentTarget.dataset
    let newIndex = index
    const oldIndex = this.data[key]
    if (newIndex === oldIndex) {
      newIndex -= 1
    }
    this.setData({
      [`${key}`]: newIndex
    })
  },
  sureData () {
    wx.showLoading({
      title: '加载中'
    })
    const { id, reserveNumber, storeIndex, serverIndex, remark, orderInfo } = this.data
    const params = {
      id,
      remark,
      reserve_number: reserveNumber,
      service_score: serverIndex + 1,
      store_score: storeIndex + 1
    }
    makeScoreApi(params).then(
      res => {
        wx.hideLoading()
        if (res.code === 200) {
          wx.showToast('评价成功')
          this.getDetail()

          // 神策数据上报 -- 开始
          sensorsLoginInit(orderInfo.member_user_id)
          // 神策数据上报 -- 邀约被评价
          wx.track.trackAppointmentReview({
            reserve_number: reserveNumber,
            user_external_id: orderInfo.user_external_id,
            evaluation_result: remark,
            activity_star: params.service_score,
            service_star: params.store_score
          })
          return
        }
        wx.showToast(res.message)
      }
    ).catch(err => {
      wx.hideLoading()
      console.error(err, '错误❎')
    })
  },
  // 页面分享设置
  onShareAppMessage () {
    return {
      title: 'share title',
      path: 'pages/store-invite/comment/comment'
    }
  }
})
