import {
  getOrderList,
  getRefundOrderList,
  getOrderLogisticsTrack
} from '@/api/order'
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
      },
      {
        label: '已退款',
        value: 4
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

    /** 搜索框--销售单号 */
    sales_order_no: '',

    /** 查看物流弹窗 */
    showLogistics: false,

    /** 订单状态背景 */
    orderStatusBg: ['wait', 'finish', 'success', 'refound']
  },

  onLoad ({ statusTabs, searchTabs }) {
    this.setData({
      statusTabs: statusTabs * 1 || 0,
      searchTabs: searchTabs * 1 || 0
    })
    this.init()
  },
  init () {
    this.setData({
      list: []
    })
    this.data.meta.page = 1
    this.getOrderList()
  },

  /** 订单列表 */
  async getOrderList () {
    try {
      const params = {
        order_state: this.data.statusTabs,
        sales_order_no: this.data.sales_order_no,
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
          item.pay_time = formatDate(
            item.pay_time * 1000,
            'yyyy-MM-dd hh:mm:ss'
          )
          return item
        })

        // 判断更多
        this.data.hasMore =
          Math.ceil(res.data.total / res.data.size) > res.data.page

        this.setData({
          list: this.data.list.concat(list),
          hasMore: this.data.hasMore
        })
      } else {
        wx.hideLoading()
        wx.showToast(res.message)
      }
    } catch (error) {
      wx.hideLoading()
    } finally {
      this.setData({
        isLoading: false
      })
    }
  },

  /** 搜索框输入 */
  handleInput (e) {
    this.setData({
      sales_order_no: e.detail.value
    })
  },

  /** 确认搜索 */
  handleConfirm () {
    this.init()
  },

  /** 切换搜索栏下TABS */
  handleChangeHeaderTabs (e) {
    const { value } = e.currentTarget.dataset

    this.setData({
      searchTabs: value
    })
    this.init()
  },

  /** 切换订单状态TABS */
  handleChangeStatus (e) {
    const { value } = e.currentTarget.dataset

    this.setData({
      statusTabs: value
    })
    this.init()
  },

  /** 查看更多 */
  handleToViewMore () {
    if (this.data.hasMore && !this.data.isLoading) {
      this.data.meta.page++
      this.getOrderList()
    }
  },

  /** 查看订单详情 */
  handleToView (e) {
    const { list, searchTabs } = this.data
    const { index } = e.currentTarget.dataset
    const { id } = list[index]

    wx.navigateTo({
      url: `/pages/order-detail/order-detail?id=${id}&search_tab=${searchTabs}`,
      events: {
        /** 更新列表 */
        update: (data) => {
          this.updateData(e, data)
        }
      }
    })
  },

  /** 发货 */
  handleToSend (e) {
    const { list } = this.data
    const { index } = e.currentTarget.dataset
    const { id, delivery_type } = list[index]

    wx.navigateTo({
      url: `/pages/send-goods/send-goods?id=${id}&delivery_type=${delivery_type}`,
      events: {
        /** 更新列表 */
        update: (data) => {
          this.updateData(e, data)
        }
      }
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
    }).catch(err => {
      console.log(err)
    })
  },

  /**
   * 更新列表
   * @param {object} e
   * @param {object} data
   */
  updateData () {
    this.init()
  }
})
