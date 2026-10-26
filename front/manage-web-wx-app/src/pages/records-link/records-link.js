// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    type: 0
  },
  // 生命周期函数--监听页面加载
  onLoad (options) {
    this.setData({
      type: Number(options.type)
    })
  }
})
