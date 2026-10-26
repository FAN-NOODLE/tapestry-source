/**
 * 两种方式
 * 一：主动获取状态；（用于入口页面，主动弹出）
 * wx.getPrivacySetting 获取授权状态；如果需要授权，返回true；弹框授权；
 * 同意后，可使用相关api；（微信能记住授权状态，下次不会弹）
 *
 * 二：被动触发状态; （用户业务页面，监听api事件）
 * wx.onNeedPrivacyAuthorization 监听事件，在需要授权时候，会触发监听事件；
 * wx.onNeedPrivacyAuthorization 监听隐私接口需要用户授权事件，当需要用户进行隐私授权时会触发；
 * 用户同意后，调用 wx.onNeedPrivacyAuthorization 的回调事件 resolve
 * 状态会同步给微信记录下来；并且会触发 wx.requirePrivacyAuthorize 的事件；同意了走success；不同意走fail事件；
 *
 *
 * 采用方案：引入组件，默认会监听，如果传了 type == get,会主动弹
 *
 * components/business/com-model-privicy/index
 * 结果：
 * 全局根据 pricacyStatus 字段状态，做适配；用户有拒绝的权利，所以需要产品思考拒绝后的逻辑
 * globalData.pricacyStatus: 0;   //0: 没授权，1: 同意，2:拒绝
 */

let privacyHandler
let privacyResolves = new Set()
let closeOtherPagePopUpHooks = new Set()

/**
 * @description 监听
 */
function onNeedPrivacyAuthorization() {
  if (wx.canIUse('onNeedPrivacyAuthorization')) {
    console.log('监听', wx.canIUse('onNeedPrivacyAuthorization'))
    wx.onNeedPrivacyAuthorization(resolve => {
      if (typeof privacyHandler === 'function') {
        privacyHandler(resolve)
      }
    })
  }
}
// 启动监听
onNeedPrivacyAuthorization()

const closeOtherPagePopUp = closePopUp => {
  closeOtherPagePopUpHooks.forEach(hook => {
    if (closePopUp !== hook) {
      hook()
    }
  })
}

Component({
  externalClasses: ['custom-class'],
  properties: {
    /** 方案 get: 方案一主动获取；listening: 监听页面api */
    openType: {
      type: String,
      value: 'listening'
    },
    /**
     * 隐私协议名字
     */
    privacyContractName: {
      type: String,
      value: ''
    }
  },
  data: {
    /** 是否展示授权弹框 */
    showPrivacy: false
  },
  pageLifetime: {
    show() {
      this.init()
    },
    hide() {
      privacyResolves.clear()
      closeOtherPagePopUpHooks.clear()
    }
  },
  lifetimes: {
    attached: function () {
      console.log('组件启动')
      this.init()
      privacyResolves.clear()
      closeOtherPagePopUpHooks.clear()
    },
    detached: function () {
      closeOtherPagePopUpHooks.delete(this.closePopUp)
    }
  },
  methods: {
    /** 初始化 */
    init() {
      // 默认会监听
      this.listeningForPrivacy()
      // 如果传了 get ，会主动拉取；
      this.getPrivacyStatus()
    },
    /**
     * 监听隐私接口需要用户授权事件
     */
    listeningForPrivacy() {
      const closePopUp = () => {
        this.disPopUp()
      }
      privacyHandler = resolve => {
        console.log('触发弹框')
        privacyResolves.add(resolve)
        this.popUp()
        // 额外逻辑：当前页面的隐私弹窗弹起的时候，关掉其他页面的隐私弹窗
        closeOtherPagePopUp(closePopUp)
      }

      this.closePopUp = closePopUp
      closeOtherPagePopUpHooks.add(this.closePopUp)
    },
    /**
     * 方案一：主动获取授权状态；
     */
    getPrivacyStatus() {
      if (wx.canIUse('getPrivacySetting')) {
        wx.getPrivacySetting({
          success: res => {
            console.log('getPrivacySetting', res)
            // 返回结果为: res = { needAuthorization: true/false, privacyContractName: '《xxx隐私保护指引》' }
            if (res.needAuthorization && this.data.openType === 'get') {
              // 需要弹出隐私协议
              this.setData({
                privacyContractName: res.privacyContractName,
                showPrivacy: true
              })
            } else {
              // 用户已经同意过隐私协议，所以不需要再弹出隐私协议，也能调用隐私接口
              this.setData({
                privacyContractName: res.privacyContractName
              })
            }
          },
          fail: () => {},
          complete: res => {
            // 用于调试
            console.log('res', res)
          }
        })
      }
    },
    /**
     * 用户同意
     */
    handleAgree() {
      this.disPopUp()
      // 这里演示了同时调用多个wx隐私接口时要如何处理：让隐私弹窗保持单例，点击一次同意按钮即可让所有pending中的wx隐私接口继续执行 （看page/index/index中的 wx.getClipboardData 和 wx.startCompass）
      privacyResolves.forEach(resolve => {
        resolve({
          event: 'agree',
          buttonId: 'agree-btn'
        })
      })
      privacyResolves.clear()
      this.triggerEvent('result', { value: true })
    },
    /**
     * 用户不同意
     */
    handleDisagree() {
      this.disPopUp()
      privacyResolves.forEach(resolve => {
        resolve({
          event: 'disagree'
        })
      })
      privacyResolves.clear()
      this.triggerEvent('result', { value: false })
    },
    /**
     * 打开
     */
    popUp() {
      if (this.data.showPrivacy === false) {
        this.setData({
          showPrivacy: true
        })
      }
    },
    /** 关闭 */
    disPopUp() {
      if (this.data.showPrivacy === true) {
        this.setData({
          showPrivacy: false
        })
      }
    },
    /**
     * tabbar 页面，来回切换
     */
    tabBarPageShow() {
      console.log('tabbar show')
      // console.log("privacyHandler", privacyHandler)
      // console.log("privacyResolves", privacyResolves)
      // console.log("closeOtherPagePopUpHooks", closeOtherPagePopUpHooks)
      if (this.closePopUp) {
        privacyHandler = resolve => {
          privacyResolves.add(resolve)
          this.popUp()
          // 额外逻辑：当前页面的隐私弹窗弹起的时候，关掉其他页面的隐私弹窗
          closeOtherPagePopUp(this.closePopUp)
        }
      }
    },
    /** 打开隐私弹窗 */
    handleOpenPrivacyContract() {
      // 打开隐私协议页面
      if (wx.canIUse('openPrivacyContract')) {
        wx.openPrivacyContract({
          success: res => {
            console.log('openPrivacyContract - success', res)
          }, // 打开成功
          fail: () => {}, // 打开失败
          complete: res => {
            console.log('openPrivacyContract - complete', res)
          }
        })
      } else {
        wx.showToast({
          title: '微信版本太低，请升级',
          icon: 'none',
          duration: 2000,
          mask: true
        })
      }
    }
  }
})
