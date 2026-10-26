import { userCrmCouponsListApi } from '@/api/order-api'
// import { toDate } from '../../utils/tools'

// 创建页面实例对象
Page({
  data: {
    userInfo: {}, // 用户信息
    couponList: [], // 订单列表
    meta: {
      page: 1,
      size: 10,
      isMore: false
    },
    isLoading: false
  },
  onLoad (options) {
    const { info } = options

    this.setData({ userInfo: JSON.parse(info) })
  },
  onShow () {
    this.getCouponList()
  },
  /**
   * 获取卡券详情
   */
  async getCouponList (page = false) {
    const { meta, couponList, userInfo } = this.data

    if (page) {
      this.setData({ 'meta.page': meta.page + 1 })
    }

    wx.showNavigationBarLoading()

    this.setData({ isLoading: true })

    const params = {
      union_id: userInfo.union_id,
      customer_code: userInfo.customer_code,
      page: meta.page,
      size: meta.size
    }

    try {
      const { code, data, message } = await userCrmCouponsListApi(params)
      this.setData({ isLoading: false })
      wx.hideNavigationBarLoading()

      if (code !== 200) {
        wx.showToast(message)
        return
      }

      const { total, records } = data

      const newCouponList = page ? couponList.concat(records) : records

      this.setData({
        'meta.total': total,
        'meta.isMore': newCouponList.length < total,
        couponList: newCouponList
      })
    } catch (error) {
      console.log(error)
      wx.hideNavigationBarLoading()
      this.setData({ isLoading: false, 'meta.page': Math.max(1, meta.page - 1) })
    }
  },
  /**
   * 展示更多
   */
  showRule (e) {
    const { index } = e.currentTarget.dataset
    const { couponList } = this.data
    couponList[index].isShow = !couponList[index].isShow

    this.setData({ couponList })
  },
  onReachBottom () {
    const { meta, isLoading } = this.data

    if (!meta.isMore || isLoading) return

    this.getOrderList(true)
  }
})
