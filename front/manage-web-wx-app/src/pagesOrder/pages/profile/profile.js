import { getMyCenter, getMySA } from '../../api/profile'

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    showLoading: true, // 页面loading
    myInfo: {},

    /** SA信息 */
    mySaInfo: {},

    /** 未添加SA弹窗是否显示 */
    showNoSa: false
  },
  // 生命周期函数--监听页面加载
  async onLoad () {
    await this.getMySA()
  },
  onShow () {
    wx.hideHomeButton()
    // TODO 获取用户信息
    this.getMyCenter()
  },
  // 点击跳转预约列表
  goRecordList () {
    wx.navigateTo({
      url: '/pagesOrder/pages/store-invite/record-list/record-list'
    })
    // 神策数据上报 -- 查看邀约记录
    wx.track.trackViewAppointmentHistory({
      entry_point: 5
    })
  },
  // 获取用户信息
  getMyCenter () {
    getMyCenter()
      .then((res) => {
        this.setData({
          showLoading: false
        })
        if (res.code !== 200) return
        this.setData({
          myInfo: res.data
        })
      })
      .catch((err) => {
        this.setData({
          showLoading: false
        })
        console.log(err, '错误')
      })
  },

  /** 获取用户SA信息 */
  async getMySA () {
    try {
      const res = await getMySA()
      if (res.code === 200) {
        this.setData({
          mySaInfo: res.data,
          showNoSa: Object.keys(res.data).length <= 0
        })
      } else {
        wx.showToast(res.message)
      }
    } catch (error) {
      console.log(
        '%c [ error ]-57',
        'font-size:13px; background:pink; color:#bf2c9f;',
        error
      )
    }
  },

  /** 跳转 */
  handleOperate (e) {
    const { type } = e.currentTarget.dataset

    switch (type) {
      case '2':
        wx.navigateTo({
          url: '/pagesOrder/pages/order-list/order-list'
        })
        break
      case '3':
        wx.navigateTo({
          url: `/pagesOrder/pages/append-sa/append-sa?qr_code=${encodeURIComponent(this.data.mySaInfo.qr_code)}`
        })
        break
    }
  }
})
