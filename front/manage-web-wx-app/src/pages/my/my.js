import { split } from '../../utils/tools'
import { getMemberDetailApi, changeSmsStatusApi } from '@/api/store-api'
import { memberPositionApi } from '@/api/compass-api'
import { whiteListApi } from '@/api/booking-api'
import { getOnlineOrderPage, customWhiteListApi } from '@/api/order'

const { globalData } = getApp()
// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    /** 查询是否展示在线订单入口 */
    onlineOrder: false,
    switchChecked: false, // 控制接受短信提醒按钮开关的开or关
    showLogout: false, // 控制退出登录二次确认弹框显隐
    showTelModal: false, // 控制接受提醒手机号码修改弹框显隐
    wxwork: globalData.system.environment === 'wxwork', // 判断当前是否在企业微信环境打开
    /** 是否白名单用户 */
    whitePeople: false,
    role: 3, // role 3 企业导购 5 企业管理（目前只有3和5）
    /** 对外收款订单图标数据 */
    externalOrderData: [
      {
        type: 1,
        iconName: 'wait',
        title: '待发货'
      },
      {
        type: 2,
        iconName: 'finish',
        title: '已发货'
      },
      {
        type: 3,
        iconName: 'success',
        title: '已完成'
      },
      {
        type: 4,
        iconName: 'refound',
        title: '退款'
      }
    ],
    showMember: false
  },

  // 生命周期函数--监听页面加载
  onLoad () {},
  onShow () {
    this.getUserIdentity()
    this.setTabBar()
    this.getDetail()
    this.getOnlineOrder()
    this.getMember()
    this.getShowMember()
  },
  /** 查询导购所在门店是否配置入口 */
  getOnlineOrder () {
    let params = {
      type: 1 // 1 在线订单(个人)
    }
    getOnlineOrderPage(params)
      .then((res) => {
        // 错误判断
        if (res.code !== 200) return wx.showToast(res.message)
        this.setData({
          onlineOrder: !!res.data.url
        })
      })
      .catch((err) => {
        console.error(err, 'getOnlineOrderPage出错')
      })
  },

  /** 是否展示会员查询 */
  getShowMember () {
    customWhiteListApi().then((res) => {
      if (res.code !== 200) return
      this.setData({
        showMember: res.data.show === 1
      })
    }).catch(e => console.error(e, '查询是否展示报错'))
  },

  /** 获取用户身份 */
  getUserIdentity () {
    whiteListApi()
      .then((res) => {
        // 错误处理
        if (res.code !== 200) return wx.showToast(res.message)
        // 0为不在白名单，1为在白名单
        const { sa, ma } = res.data
        const whitePeople = sa === 1 || ma === 1
        this.setData({
          whitePeople
        })
      })
      .catch((e) => console.error(e, '用户身份获取出错'))
  },

  setTabBar () {
    if (typeof this.getTabBar === 'function' && this.getTabBar()) {
      this.getTabBar().setData({
        selected: 2
      })
    }
  },
  // 打开对外名片
  goCard () {
    wx.navigateTo({
      url: '/pages/business-card/business-card'
    })
  },
  goOrderList () {
    wx.navigateTo({
      url: '/pages/store-invite/record-list/record-list'
    })
    // 神策数据上报 -- 查看邀约记录
    wx.track.trackViewAppointmentHistory({
      entry_point: 4
    })
  },
  goWriteOff () {
    wx.navigateTo({
      url: '/pages/write-off/write-off'
    })
  },

  /** 跳转至客户预约页面  */
  goBooking () {
    wx.navigateTo({
      url: '/pages/booking-activities/booking-activities'
    })
  },
  /**
   * 跳转顾客分配页面
   */
  goAllocation () {
    wx.navigateTo({
      url: '/packages/allocation/pages/customer-list/customer-list'
    })
  },
  /** 跳转到线上订单页面 */
  goOnlineOrder () {
    wx.navigateTo({
      url: '/pages/online-order/online-order?scene=list'
    })
  },
  /**
   * @description: 会员查询
   */
  goMemberSearch () {
    wx.navigateTo({
      url: '/packages/customer/pages/member-search/member-search?from=1'
    })
  },
  /**
   * @description: 库存查询
   */
  goStockSearch () {
    // 清空门店信息,重新读取全局门店
    wx.removeStorageSync('storeInfo')
    wx.navigateTo({
      url: '/pagesStock/pages/index/index'
    })
  },
  // 去设置页面
  goSetting () {
    wx.navigateTo({
      url: '/pages/setting/setting'
    })
  },
  // 获取员工信息接口
  getDetail () {
    const { memberId, login_cache_key } = globalData.memberInfo
    wx.showNavigationBarLoading()
    getMemberDetailApi({
      id: memberId,
      login_cache_key
    })
      .then((res) => {
        if (res.code === 200) {
          this.setData({
            ...res.data,
            show: true,
            switchChecked: res.data.is_send_sms === 1,
            mobile: split(res.data.mobile, 'space'),
            sms_mobile: split(res.data.sms_mobile, 'space')
          })
        } else {
          wx.showToast(res.message)
        }
        wx.hideNavigationBarLoading()
      })
      .catch((err) => {
        wx.hideNavigationBarLoading()
        wx.showToast(err.message)
      })
  },
  // 打开/关闭短信提醒按钮
  switchChange () {
    const { switchChecked } = this.data
    wx.showNavigationBarLoading()
    changeSmsStatusApi({
      is_send_sms: !switchChecked ? 1 : 2,
      worker_id: this.data.id
    })
      .then((res) => {
        if (res.code === 200) {
          this.setData({
            switchChecked: !switchChecked
          })
        } else {
          wx.showToast(res.message)
        }
        wx.hideNavigationBarLoading()
      })
      .catch((err) => {
        wx.hideNavigationBarLoading()
        wx.showToast(err.message)
      })
  },
  // 点击退出登录
  logout () {
    this.setData({
      showLogout: true
    })
  },
  // 退出登录二次确认弹框取消按钮回调
  close () {
    this.setData({
      showLogout: false
    })
  },
  // 退出登录二次确认弹框确认按钮回调
  sure () {
    globalData.token = ''
    wx.removeStorageSync('token')
    this.close()
    wx.redirectTo({
      url: '/pages/login/login'
    })
  },
  // 打开接受提醒手机号码修改弹框
  showTel () {
    this.setData({
      showTelModal: true
    })
  },
  // 接受提醒手机号码修改弹框的关闭操作回调
  closeTel () {
    this.setData({
      showTelModal: false
    })
  },
  // 接受提醒手机号码修改弹框的确认操作回调
  confirm (e) {
    this.setData({
      sms_mobile: split(e.detail.sms_mobile, 'space')
    })
    this.closeTel()
  },

  /** 跳转订单列表 */
  handleToOrderList (e) {
    const { type } = e.currentTarget.dataset
    if (type === 4) {
      // 已退款跳转退款信息列表
      wx.navigateTo({
        url: '/pages/order-list-foreign/order-list-foreign?searchTabs=1'
      })
    } else {
      wx.navigateTo({
        url: `/pages/order-list-foreign/order-list-foreign?statusTabs=${type}`
      })
    }
  },
  // 获取身份信息
  async getMember () {
    await memberPositionApi({ is_new: 0 }).then(res => {
      if (res.code === 200) {
        // role 3 企业导购 5 企业管理（目前只有3和5）
        this.setData({
          role: res.data.role
        })
      }
    }).catch(err => {
      console.log(err)
    })
  }
})
