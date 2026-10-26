const { globalData } = getApp()
Component({
  // 使用全局 & 页面样式
  options: {
    addGlobalClass: true
  },
  properties: {
    // 弹窗标题
    title: {
      type: String,
      value: ''
    },
    /** 回显结果集合 */
    stepList: {
      type: Object,
      value: {}
    }
  },
  data: {
    /** 是否展示、关闭预约信息弹窗 */
    show: false,
    /** 判断是否适配 */
    isIPhoneX: globalData.isIPhoneX,
    /** 图片高度 */
    imageHeightList: []
  },
  methods: {
    /** 关闭弹窗-内部方法 */
    close() {
      // this.resetData({}, true)
      this.updatePopup(false)
    },

    /** 打开、关闭弹窗, 供外部调用 */
    updatePopup(show = true) {
      // 如果是打开状态, 获取数据
      this.setData({
        show
      })
      // 通知父组件子组件打开了，关闭父组件滚动
      this.triggerEvent('onScroll', show)
    },
    /** 重置数据 */
    resetData(e, isClose = false) {
      this.updatePopup(false)
    },

    /** 确定提交筛选数据 */
    confirm() {
      this.updatePopup(false)
      wx.exitMiniProgram()
    }
  }
})
