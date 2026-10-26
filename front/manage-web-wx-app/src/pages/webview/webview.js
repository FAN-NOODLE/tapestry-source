Page({
  /** 页面的初始数据 */
  data: {
    /** 页面中需要的数据 */
    url: ''
  },

  /** 生命周期函数--监听页面加载 */
  onLoad: function (options) {
    // 从跳转页面中传过来的url在options中可以拿到
    this.setData({
      url: decodeURIComponent(options.url)
    })
    if (options.nav) {
      wx.setNavigationBarTitle({
        title: options.nav
      })
    }
  }
})
