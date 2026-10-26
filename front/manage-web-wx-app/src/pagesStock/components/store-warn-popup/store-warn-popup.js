// Component Object
Component({
  externalClasses: ['custom-class'],
  properties: {
    show: {
      type: Boolean
    }
  },
  data: {
    show: false
  },
  methods: {
    handleCancel() {
      this.setData({
        show: false
      })
      // 清空缓存
      wx.removeStorageSync('storeInfo')
      // 页面返回
      wx.navigateBack({
        delta: 1
      })
      // 关闭小程序
      const pages = getCurrentPages()
      if (pages.length === 1) {
        wx.exitMiniProgram()
      }
    },
    handleConfirm() {
      this.triggerEvent('confirm')
    }
  }
})

