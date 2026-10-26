Component({
  properties: {
    /** 弹窗是否打开 */
    show: {
      type: Boolean,
      default: false
    }
  },

  data: {},

  created () {},

  methods: {
    /** 关闭弹窗 */
    _handleClose () {
      this.setData({
        show: false
      })
    }
  }
})
