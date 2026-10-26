Component({
  properties: {
    /** 弹窗是否打开 */
    show: {
      type: Boolean,
      default: false
    },
    /** 物流信息数据 */
    list: {
      type: Object,
      default: {}
    }
  },
  data: {},

  created () {},

  methods: {
    /** 复制快递单号 */
    _handleToCopy (e) {
      wx.setClipboardData({
        data: e.currentTarget.dataset.logistics_no
      })
    },

    /** 关闭弹窗 */
    _handleClose () {
      this.setData({
        show: false
      })
    }
  }
})
