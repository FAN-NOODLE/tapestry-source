const { getCustomerListApi } = require('@/api/customer-beta-api')

const app = getApp()
// Component Object
Component({
  properties: {
    myProperty: {
      type: String,
      value: '',
      observer: function() {}
    }

  },
  data: {
    imageHost: 'https://cdn-cos.coach.com.cn/uploads/',
    filterParams: {
      rel_create_time_type: 1,
      rel_interact_time_type: 1,
      rel_active_time_type: 2
    },
    /** 筛选 */
    filter: {
      /** 筛选弹框控制 */
      show: false,
      /** 筛选条件 */
      value: {},
      /** 筛选条件数量 */
      count: 0
    },
    /** 顾客列表 */
    list: [],
    /** 分页数据 */
    meta: {
      page: 1,
      size: 10,
      total: 0,
      status: ''
    },
    /** 搜索条件 */
    customerName: '',
    exclude: []
  },

  lifetimes: {
    attached: function () {
      wx.showLoading()
      this.getData()
      this.setData({
        exclude: app.globalData.is_manager === 1 ? [] : ['member_ids']
      })
    }
  },
  methods: {
    /**
     * 获取顾客列表数据
     */
    async getData (page = 1) {
      const {
        meta: { size, status },
        filter: {
          value: filterResult,
          count
        },
        customerName
      } = this.data

      if (status === 'loading') return
      this.setData({
        'meta.status': 'loading'
      })


      const params = {
        page: page,
        size: size,
        ...filterResult,
        name: customerName,
        is_search: (customerName || count > 0) ? 1 : 0
      }

      console.log(params)
      console.log(1111, app.globalData.is_manager)

      try {
        const { data, code, message } = await getCustomerListApi(params)
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
        list: []
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
    /**
     * @description 筛选
     */
    handleFilterToggle() {
      if (this.data.filter.show) {
        this.handleFilterClose()
      } else {
        this.handleFilter()
      }
    },
    /**
     * @description 筛选
     */
    handleFilter() {
      this.setData({
        'filter.show': !this.data.filter.show
      })
      if (this.data.filter.show) {
        wx.hideTabBar()
      } else {
        wx.showTabBar()
      }
    },
    /**
     * @description 筛选 取消
     */
    handleFilterClose() {
      wx.showTabBar()
      this.setData({
        'filter.show': false
      })
    },
    /**
     * @description 筛选 确认
     */
    handleFilterConfirm(e) {
      this.setData({
        'filter.show': false,
        'filter.value': e.detail.value,
        'filter.count': e.detail.count,
        list: []
      })
      setTimeout(() => {
        wx.showTabBar()
      }, 500)
      wx.showLoading()
      this.getData()
    },

    /**
     * @description 滚动中
     */
    handleScroll() {
      this.selectComponent('#com-create-entry').scroll()
      // todo: 会员查询
      this.selectComponent('#com-search-entry').scroll()
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
      // 神策数据上报 -- 点击创建群发
      wx.track.trackBatchMsgCreate({
        is_filter: !!Object.keys(value).length,
        entry_point: '顾客列表',
        filter_option: value
      })
      wx.navigateTo({
        url: '/packages/customer/pages/create-group-post-beta/create-group-post-beta?from=1&customParams=' + customParams
      })
    }
  }
})
