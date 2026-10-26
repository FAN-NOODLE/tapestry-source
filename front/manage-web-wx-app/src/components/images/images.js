Component({
  properties: {
    src: {
      type: String
    },
    preview: {
      type: Boolean,
      value: false
    },
    list: {
      type: Array,
      value: []
    }
  },
  data: {
    name: 'images',
    src: ''
  },

  /*    系统事件开始   */
  created () {},
  attached () {},
  ready () {},
  methods: {
    /*    普通事件开始  */
    // 预览图片
    previewImage () {
      const { src: current, preview, list } = this.data
      if (preview) {
        wx.previewImage(current, list)
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
