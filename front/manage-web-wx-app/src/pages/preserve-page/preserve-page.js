// 获取全局应用程序实例对象
// const app = getApp();

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    title: 'preserve-page'
  },

  // 生命周期函数--监听页面加载
  onLoad() {
    // TODO: onLoad
  },

  // 点击重试
  goclick() {
    wx.navigateBack({
      delta: 1,
      success: function () {
        console.log('返回成功')
      },
      fail: function () {
        console.log('返回失败')
      }
    })
  }
})
