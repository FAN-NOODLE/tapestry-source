import { memberSearchListApi } from '@/api/member-search-api'

const { globalData } = getApp()
// 创建页面实例对象
Page({
  data: {
    imageHost: 'https://cdn-cos.coach.com.cn/uploads/2022/v3.0.2/weclient/',
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
    }
  },

  /*    系统事件开始   */
  /**
   *
   * @param {Object} options
   * @param {Number} options.from  1-个人中心；2-顾客列表
   * @param {String} options.searchParams  筛选条件
   */
  onLoad (options) {
    if (options.searchParams) {
      const params = JSON.parse(options.searchParams)
      this.setData({
        'filter.value': {
          ...params
        }
      })
    }
    this.setData({
      options,
      showMemberLoading: true
    })
    this.getData()
  },
  onShow () {
  },
  /**
   * 触底翻页
   */
  onReachBottom() {
    // Do something when page reach bottom.
    this.loadMore()
  },
  onPageScroll() {
    if (this.selectComponent('#com-create-entry')) {
      this.selectComponent('#com-create-entry').scroll()
    }
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
      }
    } = this.data
    const params = {
      ...filterResult
    }
    if (page) {
      params.page = page
      params.size = size
    }
    return params
  },
  /**
   * 获取顾客列表数据
   */
  async getData (page = 1) {
    const {
      meta: { status },
      options
    } = this.data
    if (status === 'loading') return
    this.setData({
      'meta.status': 'loading'
    })
    wx.showNavigationBarLoading()

    const params = this.getParams(page)

    if (page === 1) {
      // 神策数据上报 -- 会员查询
      wx.track.trackInquireMember({
        entry_point: Number(options.from),
        inquiry_condition: JSON.stringify(params)
      })
    }

    try {
      const { data, code, message } = await memberSearchListApi(params)
      wx.hideNavigationBarLoading()
      this.setData({
        showMemberLoading: false
      })
      if (code !== 200) {
        this.setData({
          'meta.status': 'loadError'
        })
        wx.showToast(message)
        return
      }


      this.data.list = (page === 1 ? [] : this.data.list).concat((data.list || []))
      this.setData({
        'meta.page': page,
        'meta.total': data.total,
        'meta.status': this.data.list.length < data.total ? 'loadMore' : 'loaded',
        list: this.data.list
      })
      if (page === 1) {
        // 神策数据上报 -- 会员查询结果
        wx.track.trackInquireMemberResult({
          entry_point: Number(options.from),
          total: data.total
        })
      }
    } catch (error) {
      console.log(error)
      this.setData({
        showMemberLoading: false,
        'meta.status': 'loadError'
      })
      wx.hideNavigationBarLoading()
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
   * @description 创建群发
   */
  handleCreate() {
    const { value } = this.data.filter
    const customParams = JSON.stringify({
      filterResult: value,
      customerName: value.name
    })
    // 神策数据上报 -- 点击创建群发
    wx.track.trackBatchMsgCreate({
      is_filter: !!Object.keys(value).length,
      entry_point: '会员查询',
      filter_option: value
    })
    wx.navigateTo({
      url: '/packages/customer/pages/create-group-post/create-group-post?from=2&customParams=' + customParams
    })
  },
  /**
   * @description 回顾客列表
   */
  handleViewCustomerList() {
    wx.navigateBack({
      delta: 2
    })
  },
  /**
   * @description 重新查询
   */
  handleRequery() {
    wx.navigateBack({
      delta: 1
    })
  }
  /*    行为事件结束   */
})
