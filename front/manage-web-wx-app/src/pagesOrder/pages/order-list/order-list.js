import {
  getOrderList,
  getRefundOrderList,
  getOrderLogisticsTrack
} from '../../api/order'
import { formatDate } from '@/utils/tools'

Page({
  data: {
    /** 订单信息、退款信息TABS */
    headerTabs: [
      {
        label: '订单信息',
        value: 0
      },
      {
        label: '退款信息',
        value: 1
      }
    ],

    /** 订单状态TABS */
    tabs: [
      {
        label: '全部',
        value: 0
      },
      {
        label: '待发货',
        value: 1
      },
      {
        label: '已发货',
        value: 2
      },
      {
        label: '已完成',
        value: 3
      }
    ],

    /** 分页数据 */
    meta: {
      page: 1, // 页数
      page_size: 10 // 每页的数据
      // total: 1, // 总数
    },
    /** 是否有更多数据 */
    hasMore: true,
    /** 是否正在加载 */
    isLoading: true,

    /** 列表数据 */
    list: [],

    /** 订单状态TABS */
    statusTabs: 0,

    /** 订单信息、退款信息TABS */
    searchTabs: 0,

    /** 查看物流弹窗 */
    showLogistics: false,

    /** 订单状态背景 */
    orderStatusBg: ['wait', 'finish', 'success']
  },

  onLoad () {},

  onShow () {
    this.setData({
      list: [],
      'meta.page': 1
    })
    this.getOrderList()
  },

  /** 订单列表 */
  async getOrderList () {
    try {
      const params = {
        order_state: this.data.statusTabs,
        ...this.data.meta
      }

      wx.showLoading({
        title: '加载中'
      })

      let res =
        this.data.searchTabs === 0
          ? await getOrderList(params)
          : await getRefundOrderList(params)
      if (res.code === 200) {
        wx.hideLoading()
        const list = res.data.list.map((item) => {
          if (this.data.searchTabs === 0) {
            item.pay_time = formatDate(
              item.pay_time * 1000,
              'yyyy-MM-dd hh:mm:ss'
            )
          } else {
            item.refund_time = formatDate(
              item.refund_time * 1000,
              'yyyy-MM-dd hh:mm:ss'
            )
          }

          return item
        })

        // 判断更多
        this.data.hasMore =
          Math.ceil(res.data.total / res.data.size) > res.data.page

        this.setData({
          list: this.data.list.concat(list),
          hasMore: this.data.hasMore,
          'meta.page': res.data.page
        })
      } else {
        wx.hideLoading()
        wx.showToast(res.message)
      }
    } catch (error) {
      console.log(
        '%c [ error ]-118',
        'font-size:13px; background:pink; color:#bf2c9f;',
        error
      )
      wx.hideLoading()
    } finally {
      this.setData({
        isLoading: false
      })
    }
  },

  /** 切换搜索栏下TABS */
  handleChangeHeaderTabs (e) {
    const { value } = e.currentTarget.dataset

    this.setData({
      searchTabs: value,
      list: [],
      'meta.page': 1
    })
    this.getOrderList()
  },

  /** 切换订单状态TABS */
  handleChangeStatus (e) {
    const { value } = e.currentTarget.dataset

    this.setData({
      statusTabs: value,
      list: [],
      'meta.page': 1
    })

    this.getOrderList()
  },

  /** 查看更多 */
  handleToViewMore () {
    if (this.data.hasMore && !this.data.isLoading) {
      this.data.meta.page++
      this.getOrderList()
    }
  },

  /** 查看小票 */
  handleToCheckReceipt (e) {
    const { index } = e.currentTarget.dataset
    wx.previewImage(this.data.list[index].tickets, [this.data.list[index].tickets])
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
    }).catch(e => console.error('接口请求报错'))
  }
})
