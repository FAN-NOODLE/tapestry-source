import { getInheritCustomerListApi, getSelectAllCountApi } from '@/api/inherit-api'

const { globalData } = getApp()

Page({
  // properties: {
  //   list: {
  //     type: Array,
  //     value: []
  //   }, // 数据反显
  //   memberIds: {
  //     type: Object,
  //     value: []
  //   }
  // },
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
    /** 可分配人数 */
    relationCount: 0,
    /** 选中的数量 */
    checkedCount: 0,
    /** 是否选中全部 */
    isAllChecked: false,
    /** 选中的列表 */
    checkedIds: [],
    /** 排除的选项 */
    excludeIds: []
  },
  /*    系统事件开始   */
  onLoad () {
    wx.showLoading({
      title: '加载中'
    })
    this.getData()
    this.getSelectAllCount()
  },
  onReachBottom () {
    this.loadMore()
  },
  /*    系统事件结束   */

  /*    普通事件开始  */
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
      is_group: false,
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
      isAllChecked,
      checkedIds,
      excludeIds
    } = this.data

    if (status === 'loading') return
    this.setData({
      'meta.status': 'loading'
    })

    wx.showNavigationBarLoading()

    const params = this.getParams(page)
    try {
      const { data, code, message } = await getInheritCustomerListApi(params)
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
        .concat((data.list || []).map(item => {
          item.checked = false
          if (item.transfer_num > 0) {
            // 全选并且没有排出的选项 或者 非全选并且选中的选项
            if ((isAllChecked && !excludeIds.includes(item.id)) || (!isAllChecked && checkedIds.includes(item.id))) {
              item.checked = true
            }
          }
          return item
        }))

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
   * @description 获取可分配人数
   */
  async getSelectAllCount() {
    try {
      const {isAllChecked, checkedCount} = this.data
      const params = this.getParams()
      const { data, code, message } = await getSelectAllCountApi(params)
      if (code === 200) {
        this.setData({
          relationCount: data.num,
          checkedCount: isAllChecked ? data.num : checkedCount
        })
      } else {
        wx.showToast(message)
      }
    } catch (error) {
      console.log(error)
    }
  },
  /**
   * @description  清空选项
   */
  clearChecked() {
    this.setData({
      isAllChecked: false,
      checkedCount: 0,
      checkedIds: [],
      excludeIds: []
    })
  },

  /*    普通事件结束  */

  /*    行为事件开始   */
  /**
   * @description 列表头部筛选
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
    this.clearChecked()
    wx.showLoading({
      title: '加载中'
    })
    this.getData()
    this.getSelectAllCount()
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
    this.clearChecked()
    wx.showLoading({
      title: '加载中'
    })
    this.getData()
    this.handleFilterPopupHide()
    this.getSelectAllCount()
  },
  /**
   * @description 搜索确认
   */
  handleSearchConfirm(e) {
    this.setData({
      customerName: e.detail.value
    })
    this.clearChecked()
    wx.showLoading({
      title: '加载中'
    })
    this.isSearch = true
    this.getData()
    this.getSelectAllCount()
  },
  /**
   * @description 单个选择/取消选择
   */
  hanldeCheckOne(e) {
    const { checked, id } = e.detail
    let { isAllChecked, relationCount, checkedCount, checkedIds, excludeIds } = this.data
    const { index } = e.currentTarget.dataset

    if (checked) {
      checkedCount++
      checkedIds.push(id)
      // 全选状态下记录排除的选项
      if (isAllChecked) {
        excludeIds = excludeIds.filter(item => item !== id)
      }
    } else {
      checkedCount--
      checkedIds = checkedIds.filter(item => item !== id)
      // 全选状态下记录排除的选项
      if (isAllChecked) {
        excludeIds.push(id)
      }
    }

    this.setData({
      checkedCount,
      checkedIds,
      excludeIds,
      [`list[${index}].checked`]: checked
    })

    // 如果选中全部则设置isAllChecked为true
    if (checkedCount > 0 && checkedCount === relationCount) {
      this.setData({
        isAllChecked: true
      })
    } else if (checkedCount === 0) {
      this.setData({
        isAllChecked: false
      })
    }
  },
  /**
   * @description 单个选择/取消选择
   */
  hanldeCheckAll() {
    const { isAllChecked, checkedCount, relationCount } = this.data

    if (relationCount === 0) {
      wx.showToast('您没有可以创建群发的好友哦')
      return
    }

    if (isAllChecked && relationCount === checkedCount) {
      // 如果已经全选了，则取消全选
      this.setData({
        isAllChecked: false,
        checkedCount: 0
      })
    } else {
      // 如果是非全选,
      this.setData({
        isAllChecked: true,
        checkedCount: relationCount
      })
    }

    this.updateListChecked()
  },
  /**
   * @description 更新列表选中状态
   */
  updateListChecked() {
    const { list, isAllChecked } = this.data

    // 列表没数据
    if (list.length === 0) return

    const {checked, checkedIds} = list.reduce((prev, item, index) => {
      if (item.transfer_num > 0) {
        prev.checked[`list[${index}].checked`] = isAllChecked
        if (isAllChecked) {
          prev.checkedIds.push(item.id)
        }
      }
      return prev
    }, { checked: {}, checkedIds: [] })

    this.setData({
      ...checked,
      checkedIds
    })
  },
  /**
   * @description 取消分配
   */
  handleBack() {
    wx.navigateBack({
      delta: 1
    })
  },
  /**
   * @description 下一步
   */
  handleNext() {
    const { isAllChecked, checkedIds, excludeIds, checkedCount } = this.data
    if (!isAllChecked && checkedIds.length === 0) {
      wx.showToast('请选择要分配的客户')
      return
    }
    if (checkedCount > 10000) {
      wx.showToast('每次最多只能分配10000个客户')
      return
    }
    const params = this.getParams()
    params.is_all_checked = isAllChecked ? 1 : 0
    /** 需要排除的ID */
    if (isAllChecked && excludeIds.length > 0) {
      params.filter_relation_id = excludeIds
    }
    /** 需要转移的ID */
    if (!isAllChecked && checkedIds.length > 0) {
      params.relation_id = checkedIds
    }
    wx.navigateTo({
      url: `/packages/allocation/pages/member-select-list/member-select-list?checkedCount=${checkedCount}&customParams=${JSON.stringify(params)}`,
      events: {
        success: () => {
          const eventChannel = this.getOpenerEventChannel()
          eventChannel.emit('success')
        }
      }
    })
  }
  /*    行为事件结束   */

})
