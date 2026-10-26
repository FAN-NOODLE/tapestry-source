import { getInheritCustomerListApi, getInheritCustomerBetaListApi } from '@/api/inherit-api'
import { toDate } from '@/utils/tools'
import Toast from '@/components/common/mp-toast/toast'

const { globalData } = getApp()
// 创建页面实例对象
Page({
  data: {
    imageHost: 'https://cdn-cos.coach.com.cn/uploads/',
    /** 顾客列表 */
    list: [],
    /** 分页数据 */
    meta: {
      page: 1,
      size: 10,
      total: 0,
      status: ''
    },
    /** 筛选 */
    filter: {
      /** 筛选弹框控制 */
      show: false,
      /** 筛选弹框顶部便宜 */
      offsetTop: globalData.system.statusBarHeight + 44 + 58 + 'px',
      /** 筛选条件 */
      value: {},
      /** 筛选条件数量 */
      count: 0
    },
    /** 搜索条件 */
    customerName: '',
    isBeat: false
  },
  /*    系统事件开始   */
  onLoad (options) {
    // this.setData({
    //   isBeat: getApp().globalData.is_beta_user === 1
    // })
    this.getData()
  },
  onShow () {},
  /**
   * 触底翻页
   */
  onReachBottom () {
    this.loadMore()
  },
  /*    统事件结束   */
  /*    普通事件开始   */
  /**
   * @description 获取接口请求参数
   * @param {number|null} page
   * @returns
  */
  getParams(page) {
    const {
      meta: { size },
      filter: {
        value: filterResult
      },
      customerName
    } = this.data
    const params = {
      name: customerName,
      is_group: true,
      ...filterResult
    }
    if (page) {
      params.page = page
      params.size = size
    }
    return params
  },
  /**
   * @description 获取顾客列表
   */
  async getData (page = 1) {
    const {
      meta: { status },
      isBeat
    } = this.data

    if (status === 'loading') return
    this.setData({
      'meta.status': 'loading'
    })
    wx.showLoading({
      title: '加载中'
    })
    wx.showNavigationBarLoading()

    const params = this.getParams(page)
    try {
      const { data, code, message } = isBeat ? await getInheritCustomerBetaListApi(params) : await getInheritCustomerListApi(params)
      wx.hideNavigationBarLoading()
      wx.hideLoading()
      if (code !== 200) {
        this.setData({
          'meta.status': 'loadError'
        })
        wx.showToast(message)
        return
      }

      this.data.list = (page === 1 ? [] : this.data.list)
        .concat((data.list || []).map(i => {
          const { more_member, ...v } = i
          // v.created_at = toDate(v.created_at, 'yyyy/MM/dd').nowTime
          if (!more_member || more_member.length === 0) {
            return { id: v.id, list: [v] }
          } else {
            // more_member.forEach(i => {
            //   i.created_at = toDate(i.created_at, 'yyyy/MM/dd').nowTime
            // })
            return { id: v.id, list: [v, ...more_member] }
          }
        }))

      this.setData({
        'meta.page': page,
        'meta.total': data.total,
        'meta.status': this.data.list.length < data.total ? 'loadMore' : 'loaded',
        list: this.data.list
      })
    } catch (error) {
      console.log(error)

      wx.hideLoading()
      wx.hideNavigationBarLoading()
      this.setData({
        'meta.status': 'loadError'
      })
    }
  },
  /**
   * @description 加载更多
   */
  loadMore() {
    const { page, status } = this.data.meta
    if (status === 'loadMore') {
      this.getData(page + 1)
    }
  },
  /*    普通事件结束   */

  /*    行为件开始   */
  /**
   * 列表头部筛选
   */
  handleSectionChange (e) {
    const { value } = e.detail
    let { value: originValue, count } = this.data.filter
    // 如果顾客筛选选择的是全部则需要减掉这个条件
    if (value === 0 && originValue.manage_type) {
      count--
      delete originValue.manage_type
    } else if (value !== 0 && !originValue.manage_type) {
      count++
      originValue.manage_type = value
    } else {
      originValue.manage_type = value
    }
    this.setData({
      'filter.value': originValue,
      'filter.count': count
    })
    wx.showLoading({
      title: '加载中'
    })
    this.getData()
  },
  /**
   * @description 显示筛选弹框
   */
  handleFilterPopupShow() {
    this.setData({ 'filter.show': true })
  },
  /**
   * @description 关闭筛选弹框
   */
  handleFilterPopupHide() {
    this.setData({ 'filter.show': false })
  },
  /**
   * @description 筛选确认
   */
  handleFilterConfirm (e) {
    this.setData({
      'filter.value': e.detail.value,
      'filter.count': e.detail.count
    })
    wx.showLoading({
      title: '加载中'
    })
    this.getData()
    this.handleFilterPopupHide()
  },
  /**
   * @description 搜索确认
   */
  handleSearchConfirm(e) {
    this.setData({
      customerName: e.detail.value
    })
    wx.showLoading({
      title: '加载中'
    })
    this.isSearch = true
    this.getData()
  },
  /**
   * @description 分配记录
   */
  handleViewRecord () {
    let url = '/packages/allocation/pages/record/record'
    wx.navigateTo({
      url: url
    })
  },
  /**
   * @description 分配顾客
   */
  handleAllocationCustomer () {
    wx.navigateTo({
      url: '/packages/allocation/pages/customer-select-list/customer-select-list',
      events: {
        success: () => {
          setTimeout(() => {
            Toast('可在分配记录中查看')
          }, 600)
          this.getData()
        }
      }
    })
  }
  /*    行为事件结束   */
})
