// 获取全局应用程序实例对象
import { makeScoreApi, scoreInfoApi } from '../../../api/store-invite'
import { toDate } from '@/utils/tools'
import { StatusList } from '../constants'

const { globalData } = getApp()

/** 时间拼接0 */
let tf = function (i) {
  return (i < 10 ? '0' : '') + i
}
// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    errPage: false, // 无权限查看
    title: 'comment',
    showLoading: true,
    storeIndex: -1,
    serverIndex: -1,
    remark: '',
    id: 0,
    enterpriseId: 0,
    reserveNumber: 0,
    orderInfo: {},
    storeDate: '05月28日',
    is_comment: 0, // 1-已评价 2-未评价
    isDisabled: false,
    notShare: false, // 是否通过分享进入评价页面
    statusList: StatusList,
    showConfirmModal: false, // 二次确认弹窗
    confirmDescTwo: '确定提交评价',
    wxwork: globalData.system.environment === 'wxwork'
  },

  // 生命周期函数--监听页面加载
  onLoad (options) {
    // wx.setStorageSync('commentOption', options)
    const { id, reserve_number, notShare } = options
    this.setData({
      id: Number(id),
      // enterpriseId: Number(options.id),
      reserveNumber: reserve_number,
      notShare
    })
    // 获取预约单详情
    this.getDetail()
    wx.hideShareMenu({
      withShareTicket: true,
      menus: ['shareAppMessage', 'shareTimeline']
    })
    wx.hideHomeButton()
  },
  onShow () {},
  // 获取评价及订单详情
  getDetail () {
    this.setData({
      errMsg: '',
      errPage: false,
      showLoading: true
    })
    const { id, reserveNumber } = this.data
    const params = { id, reserve_number: reserveNumber }
    scoreInfoApi(params)
      .then((res) => {
        if (res.code === 200) {
          const {nowYear, nowMonth, nowDate, nowHour} = toDate(
            res.data.reserve_time,
            'yyyy-MM-dd hh:mm:ss'
          )
          res.data.reserve_time_str = `${tf(nowYear)}-${tf(nowMonth)}-${tf(nowDate)} ${nowHour}:00 ~ ${nowHour + 1}:00`
          const orderInfo = res.data
          this.setData({
            orderInfo,
            showLoading: false,
            storeIndex: orderInfo.store_score - 1,
            serverIndex: orderInfo.service_score - 1,
            remark: orderInfo.remark,
            is_comment: orderInfo.is_comment,
            isDisabled: orderInfo.is_comment === 1,
            storeDate: toDate(res.data.reserve_time, 'MM月dd号').nowTime
          })
        } else if (res.code === 10001) {
          wx.showToast({
            icon: 'none',
            title: res.message
          })
          this.setData({
            errMsg: '您还没有创建到店邀约服务',
            errPage: true,
            showLoading: false
          })
        } else {
          wx.showToast({
            icon: 'none',
            title: res.message
          })
          this.setData({
            errPage: true,
            showLoading: false
          })
        }
      })
      .catch((err) => {
        this.setData({
          errPage: true,
          showLoading: false
        })
        console.error(err, '错误')
      })
  },
  // 点击星星进行评价
  clickScore (e) {
    // if (this.data.orderInfo.is_score === 1) return
    if (this.data.isDisabled) return
    const { key, index } = e.currentTarget.dataset
    let newIndex = index
    const oldIndex = this.data[key]
    if (newIndex === oldIndex) {
      newIndex -= 1
    }
    this.setData({
      [key]: newIndex
    })
  },
  // 预约单确认 -备注信息
  changeInput (e) {
    const value = e.detail.value
    this.setData({
      remark: value
    })
  },
  // 点击修改评价
  modifyComment () {
    this.setData({
      isDisabled: false
    })
  },
  // 确认预约单信息
  confirmOrder () {
    this.setData({
      showConfirmModal: true
    })
  },
  sureData () {
    wx.showLoading({
      title: '加载中'
    })
    const { id, reserveNumber, storeIndex, serverIndex, remark } = this.data
    const params = {
      id,
      remark,
      reserve_number: reserveNumber,
      service_score: serverIndex + 1,
      store_score: storeIndex + 1
    }
    this.cancelModel()
    makeScoreApi(params)
      .then((res) => {
        wx.hideLoading()
        if (res.code === 200) {
          wx.showToast({
            icon: 'none',
            title: '评价成功'
          })
          this.getDetail()
          return
        }
        wx.showToast({
          icon: 'none',
          title: res.message
        })
      })
      .catch((err) => {
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
  },
  // 关闭确认弹窗
  cancelModel () {
    this.setData({
      showConfirmModal: false
    })
  }
})
