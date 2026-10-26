Component({
  options: {
    multipleSlots: true // 启用多slot支持
  },
  properties: {
    src: {
      type: String,
      value: 'https://oss-crs.vchangyi.com/standard/lottery_default_avatar.png'
    },
    preview: {
      type: Boolean,
      value: false
    },
    bubble: {
      type: Boolean,
      value: false
    }
  },

  /*    系统事件开始   */
  created () {},
  attached () {},
  ready () {},
  /*    系统事件结束   */

  methods: {
    /*    普通事件开始  */
    imgLoad () {
      this.triggerEvent('imgLoad')
    },
    // 预览图片
    previewImage () {
      const { src: current, preview } = this.data
      if (preview) {
        wx.previewImage(current)
      }
    },
    /*    普通事件结束  */

    /*    cdp分割事件   */
    cdpSplitEvent () {}
    /*    cdp分割事件   */

    /*    行为事件开始   */

    /*    行为事件结束   */
  }
})
