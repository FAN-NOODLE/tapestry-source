// 获取全局应用程序实例对象

// import {
//   memberPositionApi
// } from '../../../api/compass-api'
import {
  guideTimeApi,
  serviceListApi,
  userInfoApi,
  createOrderApi,
  storeInfoApi,
  orderShareNumApi
} from '../../../api/guide-invite'

import { throttle } from '../../../utils/tools'

const { globalData } = getApp()

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    wxwork: globalData.system.environment === 'wxwork',
    dateModel: false,
    serverIndex: -1, // 服务选择的下标
    timeObj: {},
    showTime: false,
    canClick: false, // 创建按钮是否可以点击，防止多次点击
    userInfo: {}, // 用户的信息
    selectTime: {}, // 选择时间
    serverList: [],
    storeName: '',
    storeAddress: '',
    reserveRemark: '',
    userName: '',
    storeInfo: {},
    showModel: false,
    isSingle: false,
    modelTitle: '提醒',
    description: '该时段已有顾客预约服务，',
    descriptionTwo: '建议重新选择时段',
    showLoading: true,
    orderInfo: {},
    appid: require('../../../lib/config').appid,
    isError: false,
    errorTxt: '',
    allDay: 30
  },

  // 生命周期函数--监听页面加载
  onLoad () {
    if (!this.data.wxwork) {
      this.setData({
        isError: true,
        showLoading: false,
        errorTxt: '请通过私聊窗口发送邀请'
      })
      return
    }
    wx.hideHomeButton()
  },

  onShow () {
    // 初始化数据
    this.setData({
      serverIndex: -1,
      userInfo: {},
      selectTime: {},
      serverList: [],
      storeName: '',
      storeAddress: '',
      reserveRemark: '',
      storeInfo: {},
      userName: ''
    })
    this.init()
  },
  onHide () {
    // 清空用户身份。防止从安卓进入不刷新
    this.setData({
      userName: '',
      userInfo: {}
    })
  },
  getTime () {
    return new Promise((resolve, reject) => {
      guideTimeApi()
        .then((res) => {
          if (res.code === 200) {
            this.setData({
              timeObj: res.data,
              showTime: true
            })
            resolve()
          }
          reject(res.message)
        })
        .catch((err) => {
          reject(err)
        })
    })
  },
  init () {
    let p = [
      this.getServerList(), // 获取服务列表
      this.getStoreInfo(), // 获取门店
      this.getTime() // 获取一个月之内时间段预约的人数
    ]
    Promise.all(p)
      .then(() => {
        this.setData({
          showLoading: false
        })
        // 获取用户userId 然后获得当前用户的基本信息
        this.getUserId()
      })
      .catch((err) => {
        this.setData({
          showLoading: false
        })
        console.log(err, '出错')
      })
  },
  // 获取门店信息
  getStoreInfo () {
    return new Promise((resolve, reject) => {
      storeInfoApi()
        .then((res) => {
          if (res.code === 200) {
            this.setData({
              storeInfo: res.data,
              storeName: res.data.store_name,
              storeAddress: res.data.address
            })
            resolve()
          }
          reject(res.message)
        })
        .catch((err) => {
          reject(err)
        })
    })
  },
  // 获取用户ID
  getUserId () {
    const _this = this
    wx.qy.getContext({
      success: (res) => {
        var entry = res.entry // 返回进入小程序的入口类型
        if (entry === 'single_chat_tools') {
          // 从单聊会话的工具栏进入
          wx.qy.getCurExternalContact({
            success: (res) => {
              console.log('外部联系人返回res,=======>', res)
              _this.getUserInfo(res.userId) // 返回当前外部联系人userId
              // 神策数据上报 -- 点击侧边栏
              wx.track.trackClickItem({
                item_name: _this.route,
                user_external_id: res.userId
              })
            },
            fail: (err) => {
              console.error(err, 'getCurExternalContact===')
            }
          })
        } else {
          _this.setData({
            isError: true,
            showLoading: false,
            errorTxt: '请通过私聊窗口发送邀请'
          })
          console.log('wx.qy.getCurExternalContact==')
        }
      },
      fail: function (err) {
        console.error('wx.qy.getContext 异常', err)
        console.log(_this, 'xxxxxx')
        // _this.getUserInfo('wmvjp0DgAA-vjnccwO0FzDLkE2WZZvdA')
        _this.setData({
          isError: true,
          showLoading: false,
          errorTxt: '获取用户ID失败'
        })
      }
    })
  },
  // 获取用户信息
  async getUserInfo (external_userid) {
    try {
      const { code, data, message } = await userInfoApi({ external_userid })
      if (code === 200) {
        this.setData({
          userInfo: data,
          userName: data.name || data.nick
        })
      } else {
        wx.showToast(message)
      }
    } catch (error) {
      console.error(error, '获取用户userInfo错误')
    }
  },
  // 获取服务列表
  getServerList () {
    const params = {
      page: 1,
      size: 9999
    }
    return new Promise((resolve, reject) => {
      serviceListApi(params)
        .then((res) => {
          if (res.code === 200) {
            const serverList = res.data.list.filter(i => {
              if (i.service_type === 0) {
                const nowDate = this.getTodayTime()
                const startTime = this.getTodayTime(i.start_time * 1000)

                if (startTime > nowDate) {
                  return false
                }
              }
              return true
            })
            this.setData({
              serverList
            })
            resolve()
          }
          reject(res.message)
        })
        .catch((err) => {
          reject(err)
        })
    })
  },
  // 发送邀约单
  sendMessage (param) {
    const {
      serverIndex,
      serverList,
      orderInfo
    } = this.data
    let self = this
    wx.qy.sendChatMessage({
      // 向当前会话发送消息
      msgtype: 'miniprogram', // 消息类型，必填
      enterChat: true,
      miniprogram: {
        appid: wx.getAccountInfoSync().miniProgram.appId,
        title: '到店邀约', // 小程序消息的title
        imgUrl:
          'https://cdn-cos.coach.com.cn/uploads/2022/5/23/60f86724d69e71c10979b4e5bb17fa23.png', // 小程序消息的封面图
        page: `pagesOrder/pages/store-invite/share-detail/share-detail.html?id=${orderInfo.id}&reserve_number=${orderInfo.reserve_number}` // 小程序消息打开后的路径，注意要以.html作为后缀，否则在微信端打开会提示找不到页面
      },
      success: function (res) {
        console.log(res, '++++')
        self.setData({
          canClick: false
        })
        // 发送成功后清空信息
        self.init()
        self.getUserId()
        setTimeout(() => {
          self.setData({
            serverIndex: -1,
            selectTime: {},
            reserveRemark: '', // 备注
            storeAddress: ''
          })
        }, 100)
        const params = {
          id: orderInfo.id,
          reserve_number: orderInfo.reserve_number
        }
        orderShareNumApi(params)
        // 神策数据上报 -- 发送到店邀约
        wx.track.trackAppointmentSend({
          reserve_number: orderInfo.reserve_number,
          user_external_id: param.user_external_id,
          reserve_time: param.reserve_time,
          time_period: param.reserve_time,
          reserve_name: serverList[serverIndex].reserve_name
        })
      },
      fail: function (err) {
        console.log(err, '-----error')
        self.setData({
          canClick: false
        })
      }
    })
  },
  // 创建预约单
  create: throttle(function () {
    const params = this.getParams()
    if (typeof params === 'string') {
      wx.showToast(params)
      return
    }
    if (this.data.canClick) return
    this.setData({
      canClick: true
    })
    wx.showLoading({
      title: '加载中'
    })
    createOrderApi(params)
      .then((res) => {
        wx.hideLoading()
        if (res.code === 200) {
          this.setData({
            orderInfo: res.data
          })
          this.sendMessage(params)
          return
        }
        this.setData({
          canClick: false
        })
        wx.showToast(res.message)
      })
      .catch((err) => {
        this.setData({
          canClick: false
        })
        console.error(err, '=+++')
        wx.hideLoading()
      })
  }, 1500),
  // 跳转邀约记录页面
  goRecordList () {
    wx.navigateTo({
      url: `/pages/store-invite/record-list/record-list?union_id=${this.data.userInfo.union_id}`
    })
    // 神策数据上报 -- 查看邀约记录
    wx.track.trackViewAppointmentHistory({
      entry_point: 6
    })
  },
  /**
   * 获取当日00:00:00时间戳
   */
  getTodayTime (target) {
    let time = 0

    if (target) {
      time = new Date(new Date(new Date(new Date(target).setHours(0)).setMinutes(0)).setSeconds(0)).setMilliseconds(0)
    } else {
      time = new Date(new Date(new Date(new Date().setHours(0)).setMinutes(0)).setSeconds(0)).setMilliseconds(0)
    }

    return time
  },
  getParams () {
    const {
      serverIndex,
      userInfo, // 用户的信息
      selectTime, // 选择时间
      serverList,
      storeName,
      storeAddress,
      reserveRemark,
      storeInfo,
      userName
    } = this.data
    console.log(this.data, '__________')
    if (!userName) {
      return '请填写用户名称'
    }
    if (serverIndex === -1) {
      return '请选择服务类型'
    }
    if (!selectTime.startTime) {
      return '请选择到店时间'
    }
    if (!storeName) {
      return '请填写门店名称'
    }
    // if (!storeAddress) {
    //   return '请填写门店地址'
    // }
    const params = {
      reserve_id: serverList[serverIndex].id,
      reserve_time: selectTime.startTime,
      store_id: storeInfo.wechat_department_id,
      store_name: storeName,
      store_address: storeAddress.slice(0, 45),
      user_name: userName,
      user_unionid: userInfo.union_id,
      user_external_id: userInfo.external_userid,
      reserve_remark: reserveRemark
    }
    return params
  },
  changeInput (e) {
    const key = e.currentTarget.dataset.key
    const value = e.detail.value
    this.setData({
      [key]: value
    })
  },
  chooseTime () {
    this.setData({
      dateModel: true
    })
  },
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
    this.closeModel()
    // if (e.detail.inviteNum !== 0) {
    //   this.setData({
    //     selectTime: e.detail,
    //     showModel: true,
    //     isSingle: true,
    //     modelTitle: '提醒',
    //     description: '该时段已有顾客预约服务，',
    //     descriptionTwo: '建议重新选择时段'
    //   })
    // }
    this.setData({
      selectTime: e.detail
    })
  },
  changeServer (e) {
    const { serverList } = this.data
    const index = e.detail.value

    let allDay = 30
    if (index !== -1) {
      const item = serverList[index]
      if (item.service_type === 0) {
        const nowDate = this.getTodayTime()
        const endTime = this.getTodayTime(item.end_time * 1000)

        allDay = Math.min(Math.floor((endTime - nowDate) / (24 * 60 * 60 * 1000)) - 1, 30)
      }
    }

    this.setData({
      serverIndex: Number(index),
      allDay,
      selectTime: {}
    })
  }
})
