Component({
  properties: {
    closeIcon: {
      type: Boolean,
      value: false
    },
    titleBtn: {
      type: Boolean,
      value: true
    },
    leftBtn: {
      type: String,
      value: '取消'
    },
    rightBtn: {
      type: String,
      value: '确定'
    },
    title: {
      type: String,
      value: 'Dialog'
    }
  },
  data: {
    title: 1
  },
  /*    系统事件开始   */
  created () { },
  attached () { },
  ready () {
    wx.event.$on('modal-slot-close', this.cancelTap.bind(this, 'close'))
    wx.event.$on('modal-slot-confirm', this.cancelTap.bind(this, 'confirm'))
    this.animationStart()
  },
  /*    系统事件结束   */

  methods: {
    animation () {
      const animation = wx.createAnimation({
        duration: 200,
        timingFunction: 'linear',
        delay: 0
      })
      animation.translateY(0).step()
      this.setData({
        animationData: animation.export()
      })

      animation.translateY(574).step()
      this.setData({
        animationData: animation.export()
      })
    },
    // 弹窗动画
    animationStart () {
      const animation = wx.createAnimation({
        duration: 200,
        timingFunction: 'linear',
        delay: 0
      })
      animation.translateY(547).step()
      this.setData({
        animationData: animation.export()
      })

      animation.translateY(0).step()
      this.setData({
        animationData: animation.export()
      })
    },

    // 关闭弹层
    cancelTap (event) {
      this.animation()
      setTimeout(() => {
        this.triggerEvent(event.type ? 'close' : event)
      }, 150)
    },
    confirmTap () {
      this.cancelTap('confirm')
    }
  }
})
