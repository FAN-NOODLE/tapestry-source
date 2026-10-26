import { getCustomerListApi } from '@/api/customer-api'

const { globalData } = getApp()
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
    customerName: ''
  },
  /*    系统事件开始   */
  onLoad (options) {
    wx.showLoading({
      title: '加载中'
    })
    this.getData()
  },
  onPageScroll() {
    this.selectComponent('#com-create-entry').scroll()
    // todo: 会员查询
    this.selectComponent('#com-search-entry').scroll()
  },
  /**
   * 触底翻页
   */
  onReachBottom () {
    this.loadMore()
  },
  /*    普通事件开始   */
  /**
   * 获取顾客列表数据
   */
  async getData (page = 1) {
    const {
      meta: { size, status },
      filter: {
        value: filterResult
      },
      customerName
    } = this.data

    if (status === 'loading') return
    this.setData({
      'meta.status': 'loading'
    })
    wx.showNavigationBarLoading()

    const params = {
      page: page,
      size: size,
      ...filterResult,
      name: customerName
    }

    console.log(params)

    try {
      const { data, code, message } = await getCustomerListApi(params)
      wx.hideNavigationBarLoading()
      wx.hideLoading()
      if (code !== 200) {
        this.setData({ isLoading: false })
        wx.showToast(message)
        return
      }
      this.data.list = (page === 1 ? [] : this.data.list).concat(data.list || [])

      this.setData({
        'meta.page': page,
        'meta.total': data.total,
        'meta.status': this.data.list.length < data.total ? 'loadMore' : 'loaded',
        list: this.data.list,
        isLoading: false
      })
      if (this.isSearch && customerName) {
        // 神策数据上报 -- 顾客列表搜索
        wx.track.trackSearch({
          current_page: 3,
          search_type: 1,
          search_key_word: customerName
        })
        this.isSearch = false
      }
    } catch (error) {
      wx.hideLoading()
      console.log(error)
      wx.hideNavigationBarLoading()
    }
  },
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
   * 客户名称搜索
   */
  inputBlur (e) {
    this.setData({
      customerName: e.detail.value,
      customerList: []
    })
    wx.showLoading({
      title: '加载中'
    })
    this.isSearch = true
    this.getData()
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
  /*    统事件结束   */

  /*    行为件开始   */
  /**
   * @description 显示筛选弹框
   */
  handleFilterPopupShow() {
    wx.hideTabBar()
    this.setData({ 'filter.show': true })
  },
  /**
   * @description 关闭筛选弹框
   */
  handleFilterPopupHide() {
    this.setData({ 'filter.show': false })
    wx.showTabBar()
  },
  /**
   * 筛选条件提交
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
   * @description 会员查询
   */
  handleSearchMember() {
    wx.navigateTo({
      url: '/packages/customer/pages/member-search/member-search?from=2'
    })
  },
  /**
   * @description 创建群发
   */
  handleCreate() {
    const { filter: {value, count}, customerName } = this.data
    const customParams = JSON.stringify({
      filterResult: value, filterCount: count, customerName
    })
    wx.navigateTo({
      url: '/packages/customer/pages/create-group-post/create-group-post?from=1&customParams=' + customParams
    })
  }
  /*    行为事件结束   */
})
