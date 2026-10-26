// eslint-disable-next-line standard/object-curly-even-spacing
import { verifyCodeApi, HasVerifyListApi } from '@/api/verify'

Page({
  // 页面的初始数据
  data: {
    page: 1,
    size: 10,
    verifyList: [], // 已核销列表
    hasMore: true
  },
  // 生命周期函数--监听页面加载
  onLoad() {
    this.getVerifyList()
  },
  onShow() {
  },
  // 扫码核销
  scanQRCode() {
    let that = this
    wx.scanCode({
      success(res) {
        console.log(res, res.result)
        if (res.errMsg === 'scanCode:ok') {
          try {
            that.check(res.result)
          } catch (e) {
            console.log('e', e)
            wx.showToast('请求失败')
          }
        }
      }
    })
  },
  check(result) {
    verifyCodeApi({ code: result }).then(res => {
      if (res.code === 200) {
        wx.showToast('核销成功')
        this.setData({
          verifyList: []
        })
        this.getVerifyList()
      } else {
        wx.showToast(res.message)
      }
    }).catch(err => {
      console.log(err, '错误')
      wx.showToast('核销异常，再次核销')
    })
  },
  // 已核销数据列表
  async getVerifyList() {
    let query = {
      page: this.data.page,
      size: this.data.size
    }
    let res = await HasVerifyListApi(query)
    console.log(res, '数据')
    if (res.code === 200) {
      let { page, total, size } = res.data.meta
      this.setData({
        verifyList: [...this.data.verifyList, ...res.data.data],
        hasMore: page * size < total
      })
    }
  },
  // 上滑加载
  onReachBottom() {
    if (this.data.hasMore) {
      this.data.page += 1
      this.getVerifyList()
    }
  }
})
