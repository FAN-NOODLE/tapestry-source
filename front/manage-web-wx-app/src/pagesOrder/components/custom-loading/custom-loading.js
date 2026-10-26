Component({
  properties: {
    isLoading: {
      type: Boolean,
      value: true
    },
    hasMore: {
      type: Boolean,
      value: false
    },
    customText: {
      type: String,
      value: '没有更多了~'
    }
  },

  /*    系统事件开始   */
  created () {},
  attached () {},
  ready () {
  },
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
