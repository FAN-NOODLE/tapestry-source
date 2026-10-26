import { getCustomerListApi, getCusomerAllCountApi, getExternalIdsApi } from '@/api/customer-api'
import {
  memberSearchListApi,
  getMemberAllCountApi,
  getMemberExternalIdApi
} from '@/api/member-search-api'
import {
  checkFirstOpenApi
} from '@/api/common-api'
import Toast from '@/components/common/mp-toast/toast'

const { globalData } = getApp()
// Page Object
Page({
  data: {
    /** 静态资源域名 */
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
    /** 仅看好友 */
    isRelation: false,
    /** 好友数量 */
    relationCount: 0,
    /** 选中的数量 */
    checkedCount: 0,
    /** 是否选中全部 */
    isAllChecked: false,
    /** 选中的列表 */
    checkedIds: [],
    /** 排除的选项 */
    excludeIds: [],
    /** 选择群发素材 */
    showSelectMaterialPopup: false
  },
  // options(Object)
  /**
   * onLoad
   * @param {Object} options
   * @param {Number} options.from  1-顾客列表；2-会员查询结果列表
   * @param {String} options.customParams  筛选参数
   */
  onLoad: function(options) {
    const { customerName, filterCount, filterResult } = options.customParams ? JSON.parse(options.customParams) : {}
    this.setData({
      options,
      customerName: customerName || '',
      'filter.count': filterCount || 0,
      'filter.value': filterResult || {}
    })
    this.getData()
    this.getCusomerAllCount()
    this.checkFirstOpen()
  },
  onReachBottom() {
    // Do something when page reach bottom.
    this.loadMore()
  },
  onPageScroll() {
    this.selectComponent('#com-create-entry').scroll()
  },
  /**
   * @description 打开选择群发素材弹框
   */
  checkFirstOpen() {
    checkFirstOpenApi({ type: 3 }).then(res => {
      console.log(res)
      if (res.code === 200) {
        if (res.data.is_first_open === 1) {
          Toast('请选择您要群发的好友')
        }
      }
    }).catch(err => {
      console.log(err)
    })
  },
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
      customerName,
      isRelation
    } = this.data
    const params = {
      name: customerName,
      ...filterResult
    }
    if (page) {
      params.page = page
      params.size = size
      params.is_self = isRelation ? 1 : 0
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
      options,
      excludeIds,
      checkedIds
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
      const { data, code, message } = Number(options.from) === 2 ? await memberSearchListApi(params) : await getCustomerListApi(params)
      wx.hideNavigationBarLoading()
      wx.hideLoading()
      if (code !== 200) {
        this.setData({ isLoading: false })
        wx.showToast(message)
        return
      }

      this.data.list = (page === 1 ? [] : this.data.list).concat((data.list || []).map(item => {
        item.checked = false
        // 只能选择好友
        if (item.is_relation === 1) {
          // 全选并且没有排出的选项 或者 非全选并且选中的选项
          if ((isAllChecked && !excludeIds.includes(item.openid)) || (!isAllChecked && checkedIds.includes(item.openid))) {
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
      // if (this.isSearch && customerName) {
      //   // 神策数据上报 -- 顾客列表搜索
      //   wx.track.trackSearch({
      //     current_page: 3,
      //     search_type: 1,
      //     search_key_word: customerName,
      //   })
      //   this.isSearch = false
      // }
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
  /**
   * @description 获取数量
   */
  async getCusomerAllCount() {
    try {
      const { options, isAllChecked, checkedCount } = this.data
      const params = this.getParams()
      const { data, code, message } = Number(options.from) === 2 ? await getMemberAllCountApi(params) : await getCusomerAllCountApi(params)
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
    this.getCusomerAllCount()
    this.handleFilterPopupHide()
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
    // 数据上报用的
    // this.isSearch = true
    this.getData()
    this.getCusomerAllCount()
  },
  /**
   * @description 单个选择/取消选择
   */
  hanldeCheckOne(e) {
    const { checked, openid, work_external_user_id } = e.detail
    const { index } = e.currentTarget.dataset
    let { checkedCount, relationCount, checkedIds, isAllChecked, excludeIds } = this.data
    if (checked) {
      // 选择
      checkedCount++
      checkedIds.push(openid)
      if (isAllChecked) {
        excludeIds = excludeIds.filter(item => item !== work_external_user_id)
      }
      // 神策数据上报 -- 创建群发页操作，勾选好友
      wx.track.trackBatchMsgOperation({
        operation_type: 4
      })
    } else {
      // 取消选择
      checkedCount--
      checkedIds = checkedIds.filter(item => item !== openid)
      // 全选状态下记录排除的选项
      if (isAllChecked) {
        excludeIds.push(work_external_user_id)
      }
      // 神策数据上报 -- 创建群发页操作，勾除好友
      wx.track.trackBatchMsgOperation({
        operation_type: 5
      })
    }
    this.setData({
      checkedCount,
      checkedIds,
      excludeIds,
      [`list[${index}].checked`]: checked
    })

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
   * @description 全部选择/取消选择
   */
  handleCheckAll() {
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
      // 神策数据上报 -- 创建群发页操作，取消全选
      wx.track.trackBatchMsgOperation({
        operation_type: 2
      })
    } else {
      // 如果是非全选,
      this.setData({
        isAllChecked: true,
        checkedCount: relationCount
      })
      // 神策数据上报 -- 创建群发页操作全选
      wx.track.trackBatchMsgOperation({
        operation_type: 1
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
      if (item.is_relation === 1) {
        prev.checked[`list[${index}].checked`] = isAllChecked
        if (isAllChecked) {
          prev.checkedIds.push(item.openid)
        }
      }
      return prev
    }, { checked: {}, checkedIds: [] })

    console.log(checked, checkedIds)
    this.setData({
      ...checked,
      checkedIds
    })
  },
  /**
   * @description 仅查看好友
   */
  handleOnlyViewRelation() {
    const { isRelation } = this.data
    this.setData({
      isRelation: !isRelation
    })
    wx.showLoading({
      title: '加载中'
    })
    this.getData()
    // 神策数据上报 -- 创建群发页操作
    wx.track.trackBatchMsgOperation({
      operation_type: 0
    })
  },
  /**
   * @description 选择群发素材
   */
  handleSelectMaterial() {
    if (this.data.checkedCount === 0) {
      wx.showToast('请选择您要群发的好友')
      return
    }
    this.setData({
      showSelectMaterialPopup: true
    })
    // 神策数据上报 -- 创建群发页操作,选择群发素材
    wx.track.trackBatchMsgOperation({
      operation_type: 6
    })
  },
  /**
   * @description 选择群发素材
   */
  async handleSelectMaterialConfirm(e) {
    const { type } = e.currentTarget.dataset
    try {
      // 获取外部联系人ID
      const ids = await this.getExternalId()
      if (type === 1) {
        wx.navigateTo({
          url: '/pagesMaterial/pages/send-list/send-list',
          success: (result) => {
            result.eventChannel.emit('getExternalId', ids)
          }
        })
        // 神策数据上报 -- 创建群发页操作,选择群发素材,素材库选择
        wx.track.trackBatchMsgAssetSelect({
          operation_type: 0
        })
      } else {
        // 自定义素材
        this.shareToExternalContact(ids)
        // 神策数据上报 -- 创建群发页操作,选择群发素材,自定义素材
        wx.track.trackBatchMsgAssetSelect({
          operation_type: 1
        })
        // 神策数据上报 -- 群发素材
        wx.track.trackBatchMsgConfirm({
          asset_coverage: ids.length,
          asset_list: ['自定义素材无法获取'],
          customer_list: ids.join(','),
          operation_type: 1
        })
      }
      // 关闭弹框
      this.setData({
        showSelectMaterialPopup: false
      })
    } catch (error) {
      console.log(error)
    }
  },
  /**
   * @description 取消选择群发素材
   */
  handleSelectMaterialCancel() {
    this.setData({
      showSelectMaterialPopup: false
    })
  },
  /**
   * @description 获取外部联系人ID
   * @returns
   */
  async getExternalId() {
    try {
      const { isAllChecked, checkedIds, excludeIds, options } = this.data
      let ids = []
      // 如果点击了全选
      if (isAllChecked) {
        const params = this.getParams()

        params.filter_external_user_id = excludeIds

        const { data, code, message } = Number(options.from) === 2 ? await getMemberExternalIdApi(params) : await getExternalIdsApi(params)
        if (code === 200) {
          ids = data.external_ids
        } else {
          wx.showToast(message)
        }
      } else {
        ids = checkedIds
      }
      return ids
    } catch (error) {
      console.log(error)
    }
  },
  /**
   * @description 创建群发
   */
  shareToExternalContact(ids) {
    console.log(ids)
    wx.showLoading({
      title: '加载中...',
      mask: true
    })
    try {
      wx.qy.shareToExternalContact({
        externalUserIds: ids,
        text: {
          content: ' '
        },
        success: res => {
          wx.showToast('发送成功', true, () => {
            wx.navigateBack({
              delta: 1
            })
          })
        },
        fail: (err) => {
          console.log(err)
          wx.showToast('发送失败')
        },
        complete: () => {
          wx.hideLoading()
        }
      })
    } catch (error) {
      console.log(error)
    }
  }

})

