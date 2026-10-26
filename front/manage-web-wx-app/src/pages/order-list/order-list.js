import { userCrmOrderListApi } from '@/api/order-api'
import { toDate } from '../../utils/tools'

// 创建页面实例对象
Page({
  data: {
    userInfo: {}, // 用户信息
    orderList: [], // 订单列表
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
    this.getOrderList()
  },
  /**
   * 获取顾客详情
   */
  async getOrderList (page = false) {
    const { meta, orderList, userInfo } = this.data
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
      const { code, data, message } = await userCrmOrderListApi(params)
      this.setData({ isLoading: false })
      wx.hideNavigationBarLoading()

      if (code !== 200) {
        wx.showToast(message)
        return
      }

      const { total, records } = data

      const newRecords = records.map((i) => ({
        ...i,
        purchase_date: toDate(i.purchase_date, 'yyyy-MM-dd hh:mm:ss').nowTime
      }))
      const newOrderList = page ? orderList.concat(newRecords) : newRecords

      this.setData({
        'meta.total': total,
        'meta.isMore': newOrderList.length < total,
        orderList: newOrderList
      })
    } catch (error) {
      console.log(error)
      wx.hideNavigationBarLoading()
      this.setData({
        isLoading: false,
        'meta.page': Math.max(1, meta.page - 1)
      })
    }
  },
  onReachBottom () {
    const { meta, isLoading } = this.data

    if (!meta.isMore || isLoading) return

    this.getOrderList(true)
  }
})
