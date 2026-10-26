Page({
  data: {
    /** SA信息 */
    qr_code: ''
  },
  onLoad (options) {
    this.setData({
      qr_code: decodeURIComponent(options.qr_code)
    })
  }
})
