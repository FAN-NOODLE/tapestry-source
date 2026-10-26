
// const { globalData } = getApp()

Component({
  properties: {
    show: {
      type: Boolean,
      value: false
    },
    // 淡入效果
    fadein: {
      type: Boolean,
      value: false
    },
    title: {
      type: String,
      value: '提醒'
    },
    description: {
      type: String,
      value: ''
    },
    descriptionTwo: {
      type: String,
      value: ''
    },
    cancelText: {
      type: String,
      value: '取消'
    },
    sureText: {
      type: String,
      value: '确定'
    },
    reverse: {
      type: Boolean,
      value: false
    },
    // 是否要展示单个操作按钮
    isSingle: {
      type: Boolean,
      value: false
    },
    isAuthent: {
      type: Boolean,
      value: true
    },
    isShare: {
      type: Boolean,
      value: false
    },
    contentClass: {
      type: String,
      value: ''
    }
  },
  data: {
    showModel: true
  },

  /*    系统事件开始   */
  created () { },
  attached () { },
  ready () {
  },
  /*    系统事件结束   */

  methods: {
    /*    普通事件开始  */
    cancel () {
      this.triggerEvent('cancel')
    },
    sure () {
      this.triggerEvent('sure')
    },
    close () {
      this.triggerEvent('close')
    },
    /*    普通事件结束  */

    /*    cdp分割事件   */
    cdpSplitEvent () { }
    /*    cdp分割事件   */

    /*    行为事件开始   */

    /*    行为事件结束   */
  }
})
