Component({
  externalClasses: ['custom-class'],
  properties: {
    // 缺省一级文案
    slogan: String,
    // 缺省二级文案
    content: String,
    // 缺省图
    icon: {
      type: String,
      value: 'no-data'
    },
    // 操作按钮
    btnValue: String,
    // 缺省图距离顶部的高度
    marginTop: {
      type: Number,
      value: 30
    },
    center: {
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

    /*    普通事件结束  */

    /*    cdp分割事件   */
    cdpSplitEvent () {},
    /*    cdp分割事件   */

    /*    行为事件开始   */
    goTap () {
      this.triggerEvent('goclick')
    }
    /*    行为事件结束   */
  }
})
