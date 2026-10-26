import { login } from '@/utils/customer-login'
import { getActivityInfoApi } from '@/api/booking-customer-api'

const { globalData } = getApp()

/** 说明文案 */
const TIPS = [
  '1、预约提交后预约时间可在当天未核销前更改。',
  '2、每个用户在活动期间，只能预约一次，且只能预约当日活动。',
  '3、预约成功后，请于预约时间段内到达指定地点，凭预约二维码或您填写的预约手机号或称谓核实身份。',
  '4、如因不可抗力或者其它超出我们合理控制的事情导致您无法享受预约的皮牌压印服务，我们将及时通知您并不对此承担任何责任。'
]

Page({
  data: {
    /** 说明 */
    tips: TIPS,
    /** 活动信息 */
    detail: {
      btnText: ''
    },
    /** 活动id */
    activity_id: globalData.activity_id,
    /** 判断是否需要适配 */
    isIPhoneX: globalData.isIPhoneX
  },

  onLoad () {
    // 静默登录
    login(() => {
      // 登录成功后回调
      console.log('登录成功了')
      this.getActivityInfo()
    })

    // 获取编辑预约组件实例
    this.appointmentPopup = this.selectComponent('#appointmentPopup')
    // 获取查看预约信息组件实例
    this.lookAppointmentPopup = this.selectComponent('#lookAppointmentPopup')
  },

  onShow () {
    // 取消顶部返回按钮
    wx.hideHomeButton()
  },

  /** 获取活动信息 */
  getActivityInfo () {
    getActivityInfoApi({activity_id: this.data.activity_id})
      .then((res) => {
        // 错误判断
        if (res.code !== 200) return wx.showToast(res.message)
        const {is_reserve, qr_code, activity_status} = res.data

        // 处理按钮文案
        // activity_status 0未开始，1进行中，2已结束
        let btnText = ''
        const textList = ['预约活动未开始', '预约活动进行中', '预约活动已结束']

        if (activity_status === 1) {
          btnText = is_reserve ? '查看预约信息' : '立即预约'
          is_reserve && this.lookAppointmentPopup.updatePopup()
        } else {
          btnText = textList[activity_status]
        }
        this.setData({
          detail: {
            ...res.data,
            btnText,
            qr_code: 'data:image/jpg;base64,' + qr_code
          }
        })
      }).catch((err) => {
        console.log(err, '获取活动信息失败')
      })
  },

  /** 跳转第三方小程序 */
  goMiniProgram () {
    wx.navigateToMiniProgram({
      appId: 'wx449e93ebf040f98c',
      path: 'pages/index/index',
      success (res) {
        // 打开成功
        console.log('成功')
      }
    })
  },

  /** appointment-popup 子组件操作事件 */
  handleSuccess () {
    this.getActivityInfo()
    this.lookAppointmentPopup.updatePopup()
  },

  /** look-appointment-popup 子组件操作事件 */
  operationSuccess (e) {
    if (e.detail) {
      // 打开编辑弹窗
      this.appointmentPopup.updatePopup(true, true)
    } else {
      this.getActivityInfo()
    }
  },

  /** 打开预约弹窗 */
  openPopup () {
    // is_reserve 1为已预约，0为未预约
    const { is_reserve, activity_status } = this.data.detail

    if (activity_status !== 1) return
    is_reserve
      ? this.lookAppointmentPopup.updatePopup()
      : this.appointmentPopup.updatePopup()
    console.log('打开预约弹窗', this.data.showAboutPopup)
  },

  /** 跳转公众号备选方案 */
  handleToWebSite () {
    const url = this.data.detail.mp_url // 跳转的外链
    // const navTitle = '测试' // 这个标题是你自己可以设置的 &nav=${navTitle}
    wx.navigateTo({
      // 跳转到webview页面
      url: `/pages/webview/webview?url=${url}`
    })
  }
})

