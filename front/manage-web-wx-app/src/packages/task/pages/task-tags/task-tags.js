import {
  getCustomerApi,
  leaveGroupApi,
  getMarketingTaskDetail
} from '@/api/plan-task-api'

// Page Object
Page({
  data: {
    /** tabs */
    tabs: [
      {
        value: 0,
        title: '全部'
      },
      {
        value: 1,
        title: '好友'
      },
      {
        value: 2,
        title: '非好友'
      }
    ],
    /** current tab index */
    active: 0,
    /** 顾客列表 */
    list: [],
    /** 分页数据 */
    meta: {
      page: 1,
      size: 10,
      total: 0,
      status: ''
    },
    /** 搜索参数 */
    searchParams: {
      /** 任务ID */
      id: '',
      /** 顾客的筛选类型，0全部，1已联系，2待联系 */
      filterType: 0,
      /** 联系说明的关键字 */
      keywords: '',
      /** 0为不展示tag,1展示tag */
      is_tag: 0,
      /** 任务类型 */
      task_category: ''
    },
    /** 任务详情 */
    detail: null
  },
  /**
   * onLoad
   * @param {Object} options
   * @param {number} options.id 任务ID
   * @param {number} options.dimension
   * @param {number} options.task_category
   */
  onLoad(options) {
    this.setData({
      options,
      'searchParams.id': Number(options.id),
      'searchParams.is_tag': Number(options.dimension) === 2 ? 1 : 0,
      'searchParams.task_category': Number(options.task_category)
    })
    this.getData()
    this.getDetail()
  },
  onReachBottom() {
    this.loadMore()
  },
  /**
   * @description 获取接口请求参数
   * @param {number|null} page
   * @returns
   */
  getParams(page) {
    const {
      meta: { size },
      searchParams,
      tabs,
      active,
      options
    } = this.data
    if (Number(options.dimension) === 1) {
      return {
        id: Number(options.id),
        page: page || 1,
        size: size,
        status: tabs[active].value
      }
    }
    return {
      ...searchParams,
      page: page || 1,
      size: size,
      is_union: tabs[active].value
    }
  },
  /**
   * 获取数据
   * @param {Number} page
   */
  async getData(page = 1) {
    try {
      const {
        meta: { status },
        options
      } = this.data
      if (status === 'loading') return
      this.setData({
        'meta.status': 'loading'
      })
      const params = this.getParams(page)
      const { data, code, message } = options.dimension === 1 ? await leaveGroupApi(params) : await getCustomerApi(params)
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
      }, () => {
        // this.loadMore()
      })
    } catch (error) {
      this.setData({
        'meta.status': 'loadError'
      })
    } finally {
      this.setData({
        showLoading: false
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
  /**
   * 获取任务详情
   * @returns
   */
  async getDetail() {
    try {
      const params = {
        id: Number(this.data.options.id)
      }
      const { data, code, message } = await getMarketingTaskDetail(params)
      if (code !== 200) {
        wx.showToast(message)
        return
      }
      this.setData({
        detail: data
      })
    } catch (error) {
      console.log(error)
    }
  },
  /**
   * tab切换
   * @param {Object} e
   */
  handleTabChange(e) {
    const { index } = e.detail
    const { tabs } = this.data
    this.setData({
      active: index,
      'searchParams.customer_status': tabs[index].value
    }, () => {
      this.getData()
    })
  },
  /**
   * 筛选
   * @param {Object} e
   */
  handleFilterChange(e) {
    const { value } = e.detail
    this.setData({
      'searchParams.filterType': value
    }, () => {
      this.getData()
    })
  },
  /**
   * 输入
   * @param {*} e
   */
  handleInputChange(e) {
    const { value } = e.detail
    this.setData({
      'searchParams.keywords': value.trim()
    })
  },
  /**
   * 搜索
   */
  handleSearch() {
    this.getData()
  },
  /**
   * 搜索
   */
  handleClear() {
    this.setData({
      'searchParams.keywords': ''
    })
    this.getData()
  }
})

