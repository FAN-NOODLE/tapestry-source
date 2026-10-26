// 事件总线
import './utils/event-manage'
import { setInit, setSensorsConfig } from './lib/sensors/index'
// todo:异常上报
import { report } from './lib/robot/index'
// 挂载到全局
wx.report = report
// 新Page
require('./utils/rewrite-page')
// 扩展
require('./utils/extend-wx')

// 神策数据上报 -- 初始化
setSensorsConfig()
setInit()

App({
  globalData: {
    activity_id: 4, // 活动id
    token: '', // token，传给后端
    customer_enterprise_id: '',
    customerOption: {},
    mobileList: [],
    customer_token: '',
    checkSessionOK: false, // 阻止多次checksession
    isCanLogin: false, // 是否可以调用登陆接口
    sessionKeyLogin: false, // 是否有效期内重新登陆中
    refreshSubscribers: [], // 并发接口需要登陆时缓存的接口promies
    refreshEnterpriseSubscribers: [], // 并发接口需要登录并且入参需要企业id缓存的接口promise
    /* ———————————————————————— 业务参数  —————————————————————— */
    store: {
      store_id: ''
    },
    logo: '',
    enterprise: [], // 企业列表
    isFirstLogin: false,
    appOptions: false, // onshow options
    sysScene: false, // 系统场景
    apiSceneObj: false, // 登陆接口返回的场景值 crsScene / scene =》自购 / 分享使用
    memberInfo: {
      isBoss: false,
      memberId: 0
    }, // 接口返回来的用户信息
    isIpx: false, // 是否是iphonex
    system: {}, // 系统信息
    isReviewImg: false, // 是否预览图片控制
    canUseCustomerNav: false, // 是否可以使用自定义nav
    logisticsList: [], // 物流信息
    showBtn: true, // 是否显示订单控制按钮
    identityType: '', // 判断是导购还是合伙人
    imagePrefix: process.env.NODE_ENV === 'dev' ? 'http://localhost:8080/crs-manage-wxapp/' : 'https://img-crs.vchangyi.com/store-manage-wxapp/',
    globalPartnerStatus: 0,
    selectTimeStamp: '', // 日历选择时间戳
    defaultDate: '', // 日历默认值
    isRefreshDate: false, // 是否刷新
    currentEnterprise: {
      enterprise_id: ''
    },
    isPopLeaveBox: false, // 转发页面返回是否弹出离开框
    version: `V${wx.getAccountInfoSync().miniProgram.version || '1.0.0'}`,
    currentEnterPage: '', // TODO 企微登录
    isIos: wx.getSystemInfoSync().system.indexOf('iOS') !== -1,
    customer_login_token: '', // 顾客端登录token
    storeInfo: {}, // 门店信息
    userid: ''
  },
  // 弹框
  showMessage (option) {
    wx.showMessage(option)
  },
  // toast
  showToast (title, is, cb) {
    wx.showToast(title, is, cb)
  },
  // 版本检测 更新
  systemUpdateFn () {
    if (wx.canIUse('getUpdateManager')) {
      const updateManager = wx.getUpdateManager()

      updateManager.onCheckForUpdate(res => {
        console.log(res, '查看小程序是否更新')
        if (res.hasUpdate) {
          updateManager.onUpdateReady(() => {
            wx.showModal({
              title: '更新提示',
              content: '新版本已经准备好，是否重启应用？',
              success: res => {
                if (res.confirm) {
                  updateManager.applyUpdate()
                }
              }
            })
          })
        }
      })
    }
  },
  // 获取系统环境
  getSystemInfo () {
    let self = this
    try {
      const res = wx.getSystemInfoSync()
      console.log(res, '获取设备情况信息等')
      const { model, version, SDKVersion } = res
      self.globalData.system = res
      // ipx
      if (model.search('iPhone X') !== -1 || model.search('iPhone 11') !== -1 || model.search('iPhone 12') !== -1 || model.search('iPhone13') !== -1) {
        self.globalData.isIpx = true
      } else {
        self.globalData.isIpx = false
      }
      // customerNav
      let _version = parseInt(version.split('.').join(''))
      let _SDKVersion = parseInt(SDKVersion.split('.').join(''))
      if (_version > 700 && _SDKVersion > 243) {
        self.globalData.canUseCustomerNav = true
      }
      console.log('process.env', process.env)
    } catch (e) {
      console.log('e', e)
    }
  },
  // 获取历史第n个页面栈
  getHistoryPage (n = 1) {
    const pages = getCurrentPages() || []
    const historyPage = pages.length - n - 1 >= 0 ? pages[pages.length - (n + 1)].route : ''
    return String(historyPage)
  },

  // 第一次进入时 获取基本参数
  onLaunch (option) {
    //  版本更新
    this.systemUpdateFn()
    // 运行环境
    this.getSystemInfo()
    if (option.query.id) {
      console.log(option, '-----')
      this.globalData.customer_enterprise_id = option.query.id - 0
      this.globalData.customerOption = option
    }
  },
  // 每次打开小程序时 登录
  onShow (option) {
    let route = option.path
    console.log(route, option, '路径')

    // ApponShow的参数存储
    this.globalData.appOnShowoption = option
    // TODO 企微登录
    if (option.path !== 'pages/load/load') {
      this.globalData.currentEnterPage = option.path
      this.globalData.corp_id = option.query.corp_id
    }
  },
  // 比较特殊不要删除
  onHide () { },
  onUnLoad () { },
  onPageNotFound (res) {
    console.error(res)
    wx.redirectTo({
      url: '/pages/load/load'
    })
  },
  onError (err) {
    console.log(err)
    // todo: 错误上报
    wx.report({
      title: 'coach小程序异常',
      time: new Date().getTime(),
      errMsg: err
    })
  }
})
