import {
  getOrderDetail, getOrderLogisticsTrack, getRefoundOrderDetail
} from '@/api/order'
import { formatDate } from '@/utils/tools'

Page({
  data: {
    /** 订单ID */
    id: 0,

    /** 物流信息弹窗是否显示 */
    showLogistics: false,

    /** 物流信息数据 */
    logisticsTrack: {},

    /** 是否退款单 */
    isRefundOrder: false,

    /** 订单状态ICON */
    orderStatusIcon: ['wait', 'finish', 'success', 'wait'],

    /** 是否正在加载 */
    isLoading: true
  },

  onLoad ({ id, search_tab }) {
    this.setData({
      id,
      isRefundOrder: search_tab === 1 || search_tab === '1'
    })
  },

  onShow () {
    this.getOrderDetail()
  },

  /** 订单详情 */
  async getOrderDetail () {
    try {
      let params = {
        id: this.data.id
      }
      wx.showLoading({
        title: '加载中'
      })

      let res = !this.data.isRefundOrder
        ? await getOrderDetail(params)
        : await getRefoundOrderDetail(params)
      if (res.code === 200) {
        wx.hideLoading()
        let detail = {
          ...res.data,
          pay_time: formatDate(res.data.pay_time * 1000, 'yyyy-MM-dd hh:mm:ss'),
          refund_reqtime: formatDate(
            res.data.refund_reqtime * 1000,
            'yyyy-MM-dd hh:mm:ss'
          ),
          delivery_time: formatDate(
            res.data.delivery_time * 1000,
            'yyyy-MM-dd'
          )
        }
        this.setData({
          detail,
          isLoading: false
        })
      } else {
        wx.hideLoading()

        wx.showToast(res.message)
      }
    } catch (error) {}
  },

  /** 查看小票 */
  handleToCheckReceipt () {
    wx.previewImage(this.data.detail.tickets, [this.data.detail.tickets])
  },

  /** 复制 */
  handleToCopy (e) {
    wx.setClipboardData({
      data: e.currentTarget.dataset.logistics_no
    })
  },

  /** 发货、修改发货信息 */
  handleToSend (e) {
    const { id, delivery_type, is_edit } = e.currentTarget.dataset

    wx.navigateTo({
      url: `/pages/send-goods/send-goods?id=${id}&delivery_type=${delivery_type}&isEdit=${is_edit}`
    })
  },

  /** 查看物流 */
  handleToCheck (e) {
    const { id } = e.currentTarget.dataset

    let params = {
      id
    }

    wx.showLoading({
      title: '加载中'
    })

    getOrderLogisticsTrack(params).then((res) => {
      if (res.code === 200) {
        wx.hideLoading()
        let logisticsTrack = res.data
        this.setData({
          showLogistics: true,
          logisticsTrack
        })
      } else {
        wx.hideLoading()
        wx.showToast(res.message)
      }
    }).catch((error) => {
      console.log(error)
    })
  }
})
