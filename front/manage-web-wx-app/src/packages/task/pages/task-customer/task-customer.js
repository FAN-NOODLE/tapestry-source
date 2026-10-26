import { searchCustomerOpenidApi, searchCustomerListApi, searchCustomerBetaListApi, getMarketingTaskDetail, marketingPlanGetVideo, sendFinishApi, moveCustomerApi } from '@/api/plan-task-api'
import { getCosCompressUrl } from '@/utils/upload/utils'
import { shareToExternalContact } from '@/utils/wx.qy'

// const { globalData } = getApp()
// Page Object
Page({

  data: {
    imageHost: 'https://cdn-cos.coach.com.cn/uploads/',
    tabs: [
      {
        title: '群发顾客',
        value: 0,
        unionIds: [],
        selectedList: [],
        selectedAll: false,
        selectCount: 0,
        /** 顾客列表 */
        list: [],
        openids: [],
        /** 分页数据 */
        meta: {
          page: 1,
          size: 10,
          total: 0,
          status: ''
        }
      },
      {
        title: '单独沟通',
        value: 1,
        unionIds: [],
        selectedList: [],
        selectedAll: false,
        selectCount: 0,
        /** 顾客列表 */
        list: [],
        openids: [],
        /** 分页数据 */
        meta: {
          page: 1,
          size: 10,
          total: 0,
          status: ''
        }
      }
    ],
    active: 0,
    /** 顾客列表 */
    list: [],
    /** 分页数据 */
    meta: {
      page: 1,
      size: 1000,
      total: 0,
      status: ''
    },
    /** 筛选 */
    filter: {
      /** 筛选弹框控制 */
      show: false,
      /** 筛选弹框顶部便宜 */
      offsetTop: 0,
      /** 筛选条件 */
      value: {},
      /** 筛选条件数量 */
      count: 0
    },
    /** 搜索条件 */
    searchParams: {
      /** 任务ID */
      id: '',
      /** 顾客的筛选类型，0全部，1已联系，2待联系 */
      filterType: 0,
      /** 联系说明的关键字 */
      keywords: '',
      /** 0为群发顾客，1为沟通顾客 */
      customer_status: 0,
      /** 0为不展示tag,1展示tag */
      is_tag: 0
    },
    /** loading */
    showLoading: true,
    /** 子任务详情 */
    detail: null,
    /** 确认弹框 */
    confirmPopup: false
  },
  /**
   * onLoad
   * @param {Object} options
   * @param {number} options.id 任务ID
   * @param {number} options.dimension
   */
  onLoad(options) {
    this.setData({
      options,
      isBeat: getApp().globalData.is_beta_user === 1,
      'searchParams.id': Number(options.id),
      'searchParams.is_tag': Number(options.dimension) === 2 ? 1 : 0
    })
    this.init()
  },
  onShow() {
  },
  onUnload() {
    const eventChannel = this.getOpenerEventChannel()
    eventChannel.emit('update')
  },
  onReachBottom() {
    this.loadMore()
  },
  async init() {
    this.getDetail()
    await this.getData()
    await this.searchCustomerOpenid()
    // 设置可选择的总数
    this.setSelectCount()
  },
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
   * 获取数据
   * @param {Number} page
   */
  async getData(page = 1) {
    const {
      meta: { status },
      tabs,
      active,
      isBeat
    } = this.data
    try {
      const selectedList = tabs[active].selectedList
      if (status === 'loading') return
      this.setData({
        [`tabs[${active}].meta.status`]: 'loading'
      })
      const params = this.getParams(page)
      const { data, code, message } = isBeat ? await searchCustomerBetaListApi(params) : await searchCustomerListApi(params)
      if (code !== 200) {
        this.setData({
          [`tabs[${active}].meta.status`]: 'loadError'
        })
        wx.showToast(message)
        return
      }
      data.list = data.list || []
      // 反选
      const list = (page === 1 ? [] : this.data.tabs[active].list).concat((data.list || []).map(item => {
        item.checked = selectedList.includes(item.work_external_user_id)
        return item
      }))
      this.setData({
        [`tabs[${active}].meta.page`]: page,
        [`tabs[${active}].meta.total`]: data.total,
        [`tabs[${active}].meta.status`]: data.list.length > 0 && list.length < data.total ? 'loadMore' : 'loaded',
        [`tabs[${active}].list`]: list
      })
    } catch (error) {
      console.log(error)

      this.setData({
        [`tabs[${active}].meta.status`]: 'loadError'
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
    const { tabs, active } = this.data
    const { page, status } = tabs[active].meta
    if (status === 'loadMore') {
      this.getData(page + 1)
    }
  },
  async searchCustomerOpenid() {
    try {
      const params = this.getParams(1)
      const { data, code } = await searchCustomerOpenidApi({
        ...params,
        filterType: 2
      })
      if (code !== 200) {
        return
      }
      const { active } = this.data
      this.data.tabs[active].openids = data.length ? data : []
      this.update = true
    } catch (error) {
      console.log(error)
    }
  },
  /**
   * @description 设置可选择的总数
   */
  setSelectCount() {
    const { active, tabs } = this.data
    const { openids, selectedList } = tabs[active]
    const count = openids.length

    const selectCount = count
    this.updated = false
    this.setData({
      [`tabs[${active}].selectCount`]: selectCount,
      [`tabs[${active}].selectedAll`]: selectCount > 0 && selectedList.length === selectCount
    })
  },
  setSelect() {
    const { tabs, active } = this.data
    const { list } = tabs[active]
    const selectedList = []
    const unionIds = []
    list.forEach(item => {
      if (item.is_finished === 0) {
        item.checked = true
        unionIds.push(item.union_id)
        selectedList.push(item.work_external_user_id)
      }
    })
    this.setData({
      [`tabs[${active}].list`]: list,
      [`tabs[${active}].unionIds`]: unionIds,
      [`tabs[${active}].selectedList`]: selectedList,
      [`tabs[${active}].selectedAll`]: true
    })
  },
  /**
   * @description 获取接口请求参数
   * @param {number|null} page
   * @returns
   */
  getParams(page) {
    const {
      filter: { value, count },
      searchParams,
      tabs,
      active
    } = this.data
    const { size } = tabs[active].meta
    const params = {
      ...searchParams,
      search: {
        is_search: 0
      },
      customer_status: tabs[active].value
    }
    if (count > 0) {
      params.search = {
        is_search: 1,
        ...value
      }
    }
    if (searchParams.keywords) {
      params.search.name = searchParams.keywords
      params.search.is_search = 1
    }
    if (page) {
      params.page = page
      params.size = size
    }
    return params
  },
  /**
   * tab切换
   * @param {Object} e
   */
  async handleTabChange(e) {
    const { index } = e.detail
    // const { tabs } = this.data
    const data = {
      'searchParams.filterType': 0,
      /** 联系说明的关键字 */
      'searchParams.keywords': '',
      filter: {
        /** 筛选弹框控制 */
        show: false,
        /** 筛选弹框顶部便宜 */
        offsetTop: 0,
        /** 筛选条件 */
        value: {},
        /** 筛选条件数量 */
        count: 0
      },
      active: index,
      [`tabs[${index}].list`]: []
    }
    this.setData(data)
    await this.getData()
    await this.searchCustomerOpenid()
    // 设置可选择的总数
    this.setSelectCount()
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
  async handleSearch() {
    await this.getData()
    await this.searchCustomerOpenid()
    // 设置可选择的总数
    this.setSelectCount()
    // 神策数据上报 -- 顾客列表搜索
    wx.track.trackSearch({
      current_page: [4, 5][this.data.active],
      search_type: 1,
      search_key_word: this.data.searchParams.keywords
    })
  },
  /**
   * 搜索
   */
  async handleClear() {
    this.setData({
      'searchParams.keywords': ''
    })
    await this.getData()
    await this.searchCustomerOpenid()
    // 设置可选择的总数
    this.setSelectCount()
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
  async handleFilterConfirm (e) {
    this.setData({
      'filter.value': e.detail.value,
      'filter.count': e.detail.count
    })
    this.handleFilterPopupHide()
    await this.getData()
    await this.searchCustomerOpenid()
    // 设置可选择的总数
    this.setSelectCount()
  },
  handleNavigate() {
    const that = this
    const { options, filter, searchParams } = this.data
    wx.navigateTo({
      url: '/packages/task/pages/task-customer/task-customer',
      events: {
        update() {
          that.updated = true
          that.init()
        }
      },
      success(result) {
        result.eventChannel.emit('select', JSON.parse(JSON.stringify({
          options,
          filter,
          searchParams
        })))
      }
    })
  },
  /**
   * 更新数据
   * @param {*} e
   */
  async handleUpdate(e) {
    const { value } = e.detail
    const { index } = e.currentTarget.dataset

    const { active } = this.data
    this.setData({
      [`tabs[${active}].list[${index}]`]: value
    })
    await this.searchCustomerOpenid()
    // 设置可选择的总数
    this.setSelectCount()
    // 更新选中项
    if (value.is_finished === 1 && e.detail.value.checked) {
      e.detail.value.checked = false
      this.handleSelect(e)
    }
  },
  /**
   * 选中/取消选中
   * @param {*} e
   */
  handleSelect(e) {
    const { value } = e.detail
    const { tabs, active } = this.data
    const { selectedList, selectCount } = tabs[active]

    if (value.checked) {
      // 选中
      selectedList.push(value.work_external_user_id)
    } else {
      // 取消选中
      const index = selectedList.indexOf(value.work_external_user_id)
      console.log(111, index)

      if (index !== -1) {
        selectedList.splice(index, 1)
      }
    }

    this.setData({
      [`tabs[${active}].selectedList`]: selectedList,
      [`tabs[${active}].selectedAll`]: selectCount > 0 && selectedList.length === selectCount
    })
  },
  /**
   * 全选/取消全选
   */
  handleSelectAll() {
    const { tabs, active } = this.data
    const { list, selectedAll, selectCount, openids } = tabs[active]
    if (selectCount === 0) return
    if (selectedAll) {
      this.setData({
        [`tabs[${active}].selectedList`]: [],
        [`tabs[${active}].selectedAll`]: false
      })
    } else {
      this.setData({
        [`tabs[${active}].selectedList`]: [...openids],
        [`tabs[${active}].selectedAll`]: true
      })
    }
    // 更新列表
    this.setData({
      [`tabs[${active}].list`]: list.map(item => {
        if (item.is_finished === 0) {
          item.checked = !selectedAll
        }
        return item
      })
    })
    this.handleTrackTaskDetailOperation(!selectedAll ? 71 : 72)
  },
  /**
   * 获取视频素材id
   */
  getVideoMediaId(id) {
    return new Promise((resolve) => {
      marketingPlanGetVideo({ id: Number(id) })
        .then((res) => {
          resolve(res.data.media_id || '')
        })
        .catch((e) => console.error('获取media_id 失败'))
    })
  },
  /** 创建消息内容 */
  async createAttachments() {
    const { content_type, content_info } = this.data.detail
    const { image, video, link, applets, text } = content_info
    const message = {
      text: {
        content: text.content
      },
      attachments: []
    }
    let assetIds = []
    let types = {}
    switch (content_type) {
      case 2:
        message.attachments = image.imgList.map((item) => {
          return {
            msgtype: 'image',
            image: {
              imgUrl: getCosCompressUrl({
                url: item.pic_url,
                maxWidth: 1080,
                maxHeight: 1440
              })
            }
          }
        })
        //  取出素材id
        assetIds = image.imgList.map((item) => {
          types[item.media_id] = 2
          return item.media_id
        })
        break
      case 4:
        const media_id = await this.getVideoMediaId(video.media_id)
        message.attachments = [
          {
            msgtype: 'video',
            video: {
              mediaid: media_id
            }
          }
        ]
        //  取出素材id
        assetIds = [video.media_id]
        types[video.media_id] = 5
        break
      case 6:
        message.attachments = [
          {
            msgtype: 'link',
            link: {
              url: link.link_url,
              title: link.title,
              imgUrl: link.cover_url,
              desc: link.description
            }
          }
        ]
        //  取出素材id
        assetIds = [link.id]
        types[link.id] = 3
        break
      case 7:
        message.attachments = [
          {
            msgtype: 'miniprogram',
            miniprogram: {
              appid: applets.app_id,
              title: applets.title,
              imgUrl: applets.cover_url,
              page: applets.path
            }
          }
        ]
        //  取出素材id
        assetIds = [applets.id]
        types[applets.id] = 4
        break
      default:
        break
    }
    return { message, assetIds, types }
  },
  /**
   * @description 分享到外部联系人
   */
  async shareToExternalContact() {
    try {
      wx.showLoading({
        title: '加载中',
        mask: true
      })
      const { tabs, active } = this.data
      const { selectedList } = tabs[active]
      const { assetIds, message, types } = await this.createAttachments()

      // 获取顾客列表
      await shareToExternalContact({
        ...message,
        externalUserIds: selectedList
      })
      wx.showToast('发送成功')
      this.handleFinish(assetIds, types)
    } catch (error) {
      console.log(error)
    } finally {
      wx.hideLoading()
    }
  },
  handleFinish(assetIds, types) {
    const { detail, tabs, active } = this.data
    const { selectedList } = tabs[active]
    if (selectedList.length === 0) return

    this.sendFinish()

    this.handleTrackTaskDetailOperation(9)
    // 神策数据上报 -- 群发素材
    wx.track.trackCommunicationSend({
      task_id: detail.task_id,
      customer_number: selectedList.length
    })
    // 神策数据上报 -- 发送素材
    selectedList.forEach(item => {
      assetIds.forEach((child) => {
        wx.track.trackAssetSend({
          asset_id: child,
          task_id: detail.task_id,
          customer_id: item,
          operation_type: 3,
          is_official: 1,
          asset_type_modify: types[child]
        })
      })
    })
  },
  /**
   * @description 完成任务
   */
  async sendFinish() {
    const { options, tabs, active } = this.data
    const { selectedList } = tabs[active]
    await sendFinishApi({
      task_allot_log_id: Number(options.id),
      open_ids: selectedList
    })
    await this.getData()
    await this.searchCustomerOpenid()
    // 设置可选择的总数
    this.setSelectCount()
    // 更新选中列表
    this.setData({
      [`tabs[${active}].unionIds`]: [],
      [`tabs[${active}].selectedList`]: []
    })
  },
  /**
   * @description 移动顾客
   */
  async moveCustomer(e) {
    const { status } = e.currentTarget.dataset
    const { options, tabs, active } = this.data
    const { selectCount, selectedList } = tabs[active]
    const { code, message } = await moveCustomerApi({
      id: Number(options.id),
      open_ids: selectedList,
      customer_status: Number(status)
    })
    if (code !== 200) {
      wx.showToast(message)
      return
    }
    wx.showToast(Number(status) === 1 ? '移入成功' : '操作成功')
    this.setData({
      confirmPopup: false,
      [`tabs[${active}].unionIds`]: [],
      [`tabs[${active}].selectedList`]: [],
      [`tabs[${active}].selectedAll`]: false,
      [`tabs[${active}].selectCount`]: selectCount - selectedList.length
    })

    await this.getData()
    await this.searchCustomerOpenid()
    // 设置可选择的总数
    this.setSelectCount()
    this.handleTrackTaskDetailOperation([11, 8][status])
  },
  /**
   * @description 显示弹窗
   */
  handleShowPopup() {
    this.setData({
      confirmPopup: true
    })
  },
  /**
   * @description 关闭弹窗
   */
  handleCancel() {
    this.setData({
      confirmPopup: false
    })
  },
  handleTrack(e) {
    const { type } = e.detail
    this.handleTrackTaskDetailOperation(type)
  },
  /** 数据上报 */
  handleTrackTaskDetailOperation(type) {
    // 神策数据上报 -- 沟通子任务详情页操作
    wx.track.trackTaskDetailOperation({
      task_id: this.data.detail.task_id,
      operation_type: type
    })
  }
})
