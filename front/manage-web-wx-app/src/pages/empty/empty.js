// 创建页面实例对象
Page({
  onLoad (option) {
    if (option.currentPage) {
      wx.switchTab({
        url: '/' + option.currentPage
      })
    }
  }
})
