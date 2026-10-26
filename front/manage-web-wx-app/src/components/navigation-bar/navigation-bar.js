const app = getApp()
Component({
  properties: {
    text: {
      type: String,
      value: '我的'
    },
    back: {
      type: Boolean,
      value: false
    },
    icon: {
      type: Boolean,
      value: false
    },
    backActive: {
      type: String,
      value: ''
    },
    throughTop: {
      // 页面通顶展示
      type: Boolean,
      value: false
    },
    // 颜色
    color: {
      type: String,
      value: '#000'
    },
    // 背景色
    bgColor: {
      type: String,
      value: ''
    },
    showNaviTitle: {
      type: Boolean,
      value: true
    },
    fontSize: {
      type: Number,
      value: 16
    },
    backText: {
      type: String,
      value: ''
    },
    backBoxWidth: {
      type: Number,
      value: 50
    }
  },
  data: {
    statusBarHeight: app.globalData.system.statusBarHeight + 'px',
    navigationBarHeight: app.globalData.system.statusBarHeight + 44 + 'px',

    vipLeave: false // 离开确认
  },

  /*    系统事件开始   */
  created () {},
  attached () {},
  ready () {
    const route = app.getHistoryPage(1)
    this.setData({
      show: Boolean(route)
    })
  },
  /*    系统事件结束   */

  methods: {
    /*    普通事件开始  */

    /*    普通事件结束  */

    /*    cdp分割事件   */
    cdpSplitEvent () {},
    /*    cdp分割事件   */

    /*    行为事件开始   */
    backPage: function () {
      this.triggerEvent('back')
      if (this.data.backActive === 'goback') {
        wx.navigateBack({
          delta: 1
        })
      }
    },
    leave: function () {
      this.backPage()
      wx.navigateBack()
      this.setData({
        isLeave: false
      })
    },
    // 关闭
    closeVipModel (e) {
      this.setData({
        vipLeave: false
      })
      if (e.currentTarget.dataset.type === 'close') return
      wx.navigateBack({
        delta: 1
      })
    },
    // 继续授权
    authVip () {
      this.setData({
        vipLeave: false
      })
    },
    cancel: function () {
      this.setData({
        isLeave: false
      })
    }
    /*    行为事件结束   */
  }
})
