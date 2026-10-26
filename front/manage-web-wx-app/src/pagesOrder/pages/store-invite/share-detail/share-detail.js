// 获取全局应用程序实例对象
import { sensorsLoginInit } from '@/lib/sensors'
import {
  orderDetailApi,
  orderEditApi,
  getConfigMessageApi,
  orderNotifyApi
} from '../../../api/guide-invite'
import { StatusList } from '../constants'

const { globalData } = getApp()

/** 时间补0 */
let tf = function (i) {
  return (i < 10 ? '0' : '') + i
}

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    id: 0,
    serverIndex: -1, // 服务选择的下标
    serverList: [], // 服务列表
    reserveNumber: '', // 备注
    statusList: StatusList,
    showLoading: true, // 页面loading
    errPage: false, // 无权限查看
    errMsg: '',
    dateModel: false, // 是否可以展示选择时间
    detail: {}, // 预约单详情
    showModel: false, // 店员完成服务弹框
    // canShare: true, // 是否向用户分享评论页面
    goShare: false, // 去分享了
    tipModel: false, // 选择时间是是否占用
    writeOffModel: false, // 展示核销弹框
    orderNumber: '',
    showConfirmModal: false, // 二次确认弹窗
    confirmDesc: '确认',
    confirmDescTwo: '确定提交该预约单',
    user_remark: '', // 用户备注信息
    tmplIds: [] // 模版ID
  },
  // 生命周期函数--监听页面加载
  onLoad (options) {
    this.data.id = Number(options.id)
    this.data.reserveNumber = options.reserve_number

    // 如果是企微端打开邀约详情，需要重定向到企微端的邀约详情
    if (globalData.system.environment === 'wxwork') {
      return wx.reLaunch({
        url: `/pages/store-invite/invite-detail/invite-detail?id=${this.data.id}&reserveNumber=${this.data.reserveNumber}`
      })
    }
    wx.hideShareMenu({
      withShareTicket: true,
      menus: ['shareAppMessage', 'shareTimeline']
    })
    wx.hideHomeButton()
  },
  onShow () {
    this.getDetail()
    this.getConfigMessage()
  },
  // 预约单确认 -备注信息
  changeInput (e) {
    const key = e.currentTarget.dataset.key
    const value = e.detail.value
    this.setData({
      [key]: value
    })
  },
  // 获取预约详情信息
  getDetail () {
    this.setData({
      errMsg: '',
      errPage: false
    })
    const { id, reserveNumber } = this.data
    const params = {
      id,
      reserve_number: reserveNumber
    }
    orderDetailApi(params)
      .then((res) => {
        if (res.code === 200) {
          const orderInfo = res.data

          if (res.data.status === 0) {
            wx.showShareMenu({
              withShareTicket: true,
              menus: ['shareAppMessage', 'shareTimeline']
            })
          }

          const date = new Date(orderInfo.reserve_time * 1000)
          // todo: 到店开始结束时间
          let select_hours = ''
          if (orderInfo.reserve_end_time === 0) {
            select_hours = `${date.getHours()}:00 - ${date.getHours() + 1}:00`
          } else {
            const endDate = new Date(orderInfo.reserve_end_time * 1000)
            select_hours = `${date.getHours()}:00 - ${endDate.getHours()}:00`
          }
          const select_name = `${date.getFullYear()}-${
            tf(date.getMonth() + 1)
          }-${tf(date.getDate())}`
          const selectTime = {
            select_name,
            inviteNum: 0,
            isFinish: false,
            name: select_hours,
            startTime: orderInfo.reserve_time
          }

          this.setData({
            selectTime,
            showLoading: false,
            detail: res.data,
            user_remark: orderInfo.user_remark
          })
        } else if (res.code === 10001) {
          this.setData({
            showLoading: false,
            errPage: true,
            errMsg: '抱歉，您非邀约客户，暂无法查看预约内容'
          })
          wx.showToast({
            icon: 'none',
            title: res.message
          })
        } else {
          this.setData({
            errPage: true,
            showLoading: false
          })

          wx.showToast({
            icon: 'none',
            title: res.message
          })
        }
      })
      .catch((err) => {
        console.log(err)
        this.setData({
          errPage: true,
          showLoading: false
        })
      })
  },
  // 确认预约单信息
  confirmOrder () {
    if (this.data.tmplIds) {
      this.handleSubscribe()
    }
  },
  // 关闭确认弹窗
  cancelModel () {
    this.setData({
      showConfirmModal: false
    })
  },
  // 获取订阅消息配置
  getConfigMessage () {
    getConfigMessageApi()
      .then((res) => {
        console.log('获取订阅消息配置res=====》', res)
        if (res.code !== 200) return
        this.data.tmplIds = res.data['confirm-subscribe']
      })
      .catch((err) => {
        console.error(err, '====')
      })
  },
  /**
   * 执行订阅
   */
  handleSubscribe () {
    const { tmplIds, id, reserveNumber } = this.data

    console.log('tmplIds: ', tmplIds)
    wx.requestSubscribeMessage({
      tmplIds: [...tmplIds],
      success: (res) => {
        // 过滤出接受的template_id
        const tempRes = Object.keys(res).filter(
          (e) => e !== 'errMsg' && res[e] === 'accept'
        )
        console.log('handleSubscribe', tempRes)
        if (res.errMsg === 'requestSubscribeMessage:ok' && tempRes.length > 0) {
          const params = {
            rel_id: id, // 订单ID
            tmpl_list: tmplIds.map((item) => {
              return {
                pri_tmpl_id: item,
                url: `?id=${id}&reserve_number=${reserveNumber}`
              }
            })
          }
          orderNotifyApi(params)
            .then(() => {
              wx.showToast('如有其他需要，请联系您的导购')
            })
            .catch((err) => {
              console.error('订阅错误信息-->', err)
            })
        } else {
          wx.showToast({
            title:
              '您已关闭/拒收此类订阅消息，可在小程序设置-接受订阅消息进行设置开启',
            icon: 'none',
            duration: 2000,
            mask: false
          })
        }
      },
      fail (err) {
        console.error('订阅错误信息', err)
        wx.showToast({
          title:
            '您已关闭/拒收此类订阅消息，可在小程序设置-接受订阅消息进行设置开启',
          icon: 'none',
          duration: 2000,
          mask: false
        })
      },
      complete: () => {
        // 确认订单
        this.sureOperation()
      }
    })
  },
  // 点击确认预约单 -拉起订阅弹窗
  sureOperation () {
    // 1. 确认预约单
    this.submitOrder()
  },
  // 提交预约单
  submitOrder () {
    const { id, reserveNumber, user_remark, detail } = this.data
    const params = {
      id,
      user_remark,
      reserve_number: reserveNumber,
      reserve_time: detail.reserve_time,
      reserve_end_time: detail.reserve_end_time,
      reserve_id: detail.reserve_id,
      status: 1
    }
    wx.showLoading({
      title: '加载中'
    })
    this.cancelModel()
    orderEditApi(params)
      .then((res) => {
        wx.hideLoading()
        if (res.code === 200) {
          wx.showToast({
            title: '确认成功',
            duration: 2000,
            mask: false
          })
          this.getDetail()
          // 神策数据上报 -- 开始
          sensorsLoginInit(detail.member_user_id)
          // 神策数据上报 -- 邀约被接受
          wx.track.trackAppointmentAccept({
            reserve_number: detail.reserve_number,
            user_external_id: detail.user_external_id,
            reserve_name: detail.serve_name,
            reserve_time: detail.reserve_time,
            time_period: detail.reserve_time,
            store_name: detail.store_name,
            store_code: detail.store_name
          })
          return
        }
        wx.showToast({
          icon: 'none',
          title: res.message
        })
      })
      .catch((err) => {
        wx.hideLoading()
        console.log(err)
      })
  }
})
