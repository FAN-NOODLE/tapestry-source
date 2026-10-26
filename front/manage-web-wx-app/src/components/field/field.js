Component({
  options: {
    multipleSlots: true // 启用多slot支持
  },
  properties: {
    bottom: {
      type: Boolean,
      value: true
    },
    arrow: {
      type: Boolean,
      value: false
    },
    height: {
      type: String,
      value: 'auto'
    }
  },

  /*    系统事件开始   */
  created () {},
  attached () {},
  ready () {},
  /*    系统事件结束   */

  methods: {
    /*    普通事件开始  */
    /*    普通事件结束  */

    /*    cdp分割事件   */
    cdpSplitEvent () {}
    /*    cdp分割事件   */

    /*    行为事件开始   */

    /*    行为事件结束   */
  }
})
