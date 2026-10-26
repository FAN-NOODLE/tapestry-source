

import { getInheritRecordInfoApi } from '@/api/inherit-api'
import { toDate } from '@/utils/tools'

// 创建页面实例对象
Page({
  data: {
    headerData: {},
    customerList: [], // 顾客列表
    meta: {
      page: 1,
      size: 20,
      total: 0,
      isMore: false
    },
    inheritId: null, // 记录id
    isLoading: false
  },
  onLoad (option) {
    const { routerParams } = option
    if (routerParams) {
      const params = JSON.parse(decodeURIComponent(routerParams))
      params.created_at = toDate(params.created_at, 'yyyy/MM/dd hh:mm:ss').nowTime
      params.updated_at = toDate(params.updated_at, 'yyyy/MM/dd hh:mm:ss').nowTime
      console.log(params)
      this.setData({
        inheritId: params.id,
        headerData: params
      })
    }
  },
  onShow () {
    this.getRecordDetail()
  },
  /**
   * 获取记录详情
   */
  async getRecordDetail (page = false) {
    wx.showNavigationBarLoading()
    this.setData({ isLoading: true })

    const { inheritId, meta, customerList } = this.data

    // 分页
    if (page) {
      this.setData({ 'meta.page': meta.page + 1 })
    }

    const params = {
      inherit_id: inheritId,
      page: meta.page,
      size: meta.size
    }

    try {
      const { code, data, message } = await getInheritRecordInfoApi(params)
      wx.hideNavigationBarLoading()
      this.setData({ isLoading: false })

      if (code !== 200) {
        wx.showToast(message)
        return
      }

      const { list = [], total } = data
      const newList = list.map(i => ({...i, takeover_time: i.takeover_time === 0 ? '' : toDate(i.takeover_time, 'yyyy/MM/dd hh:mm:ss').nowTime}))
      const newCustomerList = page ? customerList.concat(newList) : newList
      this.setData({
        customerList: newCustomerList,
        'meta.total': total,
        'meta.isMore': newCustomerList.length < total
      })
    } catch (error) {
      console.log(error)
      wx.hideNavigationBarLoading()
      this.setData({ isLoading: false, 'meta.page': Math.max(1, meta.page - 1) })
    }
  },
  // 触底翻页
  async onReachBottom () {
    let { meta, isLoading } = this.data

    if (!meta.isMore || isLoading) return

    this.getRecordDetail(true)
  }
})
