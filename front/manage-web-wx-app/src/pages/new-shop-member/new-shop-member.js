const { globalData } = getApp()

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    isIpx: globalData.isIpx
  },

  // 生命周期函数--监听页面加载
  onLoad () {
    this.setData({
      poster: globalData.store.worker_qr_code
    })
  },
  // 打开相册权限设置
  openConfirm () {
    let self = this
    wx.showModal({
      content: '检测到您没打开保存图片到相册的权限，是否去设置打开？',
      confirmText: '确认',
      cancelText: '取消',
      success: function (res) {
      // 点击“确认”时打开设置页面
        if (res.confirm) {
          wx.openSetting({
            success (res) {
              self.save()
            }
          })
        }
      }
    })
  },
  save () {
    let self = this
    wx.getImageInfo({
      src: globalData.store.worker_qr_code,
      success: function (img) {
        // 将图片保存到相册
        wx.saveImageToPhotosAlbum({
          filePath: img.path,
          success () {
            wx.showToast('已保存至手机相册')
          },
          fail () {
            // 判断是否获得了用户保存相册授权
            wx.getSetting({
              success: res => {
                if (!res.authSetting['scope.writePhotosAlbum']) {
                  self.openConfirm()
                }
              }
            })
          }
        })
      },
      fail: function (res) {
        wx.showToast('获取图片信息失败！', false)
      }
    })
  }
})
