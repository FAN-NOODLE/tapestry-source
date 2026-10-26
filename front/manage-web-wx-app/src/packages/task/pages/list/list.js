import { getCustomerTaskListApi } from '@/api/customer-beta-api'


// Page Object
Page({
  data: {
    /** 顾客列表 */
    list: [],
    /** 分页数据 */
    meta: {
      page: 1,
      size: 10,
      total: 0,
      status: ''
    },
    searchParams: {
      union_id: '',
      task_category: -1,
      status: -1
    },
    showTypePicker: false,
    typePickerOptions: [
      { label: '全部', value: -1 },
      { label: '企微沟通任务', value: 0 },
      { label: '朋友圈任务', value: 1 }
    ],
    showStatusPicker: false,
    statusPickerOptions: [
      { label: '全部', value: -1 },
      { label: '待完成', value: 0 },
      { label: '已完成', value: 1 },
      { label: '已过期', value: 2 }
    ]
  },
  // options(Object)
  onLoad: function(options) {
    this.setData({
      options,
      'searchParams.union_id': options.unionid
    })
    this.getData()

    // 神策数据上报 -- 任务列表页浏览
    wx.track.trackViewTaskList()
  },
  onReachBottom: function() {
    this.loadMore()
  },
  async getData(page = 1) {
    const {
      meta: { size, status },
      searchParams
    } = this.data

    if (status === 'loading') return
    this.setData({
      'meta.status': 'loading'
    })

    wx.showNavigationBarLoading()

    const params = {
      page,
      size,
      ...searchParams
    }
    try {
      const { data, code, message } = await getCustomerTaskListApi(params)
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
        .concat(data.list || [])

      this.setData({
        'meta.page': page,
        'meta.total': data.total,
        'meta.status': this.data.list.length < data.total ? 'loadMore' : 'loaded',
        list: this.data.list
      })
    } catch (error) {
      wx.hideLoading()
      wx.hideNavigationBarLoading()
      this.setData({
        'meta.status': 'loadError'
      })
    }
  },
  loadMore() {
    const { page, status } = this.data.meta
    if (status === 'loadMore') {
      this.getData(page + 1)
    }
  },
  handleSelectType() {
    this.setData({
      showTypePicker: true
    })
  },
  handleSelectTypeConfirm(e) {
    console.log(e)

    this.setData({
      showTypePicker: false,
      'searchParams.task_category': e.detail.value.value
    })
    this.getData()

    // 神策数据上报 -- 任务列表页筛选
    wx.track.trackFilter({
      current_page: 10,
      filter_type_sv: 3,
      filter_option: {
        '任务类型': e.detail.value.label
      }
    })
    // 神策数据上报 -- 筛选详情
    wx.track.trackFilterDetails({
      current_page: 10,
      button_name: '确定',
      filter_type_sv: 3,
      tag_level_1: '',
      tag_level_2: '任务类型',
      tag_value: e.detail.value.label
    })
  },
  handleSelectStatus() {
    this.setData({
      showStatusPicker: true
    })
  },
  handleSelectStatusConfirm(e) {
    this.setData({
      showStatusPicker: false,
      'searchParams.status': e.detail.value.value
    })
    this.getData()
    // 神策数据上报 -- 任务列表页筛选
    wx.track.trackFilter({
      current_page: 10,
      filter_type_sv: 4,
      filter_option: {
        '任务状态': e.detail.value.label
      }
    })

    // 神策数据上报 -- 筛选详情
    wx.track.trackFilterDetails({
      current_page: 10,
      button_name: '确定',
      filter_type_sv: 4,
      tag_level_1: '',
      tag_level_2: '任务状态',
      tag_value: e.detail.value.label
    })
  },
  handleCancel() {
    this.setData({
      showTypePicker: false,
      showStatusPicker: false
    })
  }
})
