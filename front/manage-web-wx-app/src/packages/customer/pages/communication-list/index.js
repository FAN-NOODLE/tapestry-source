import { gerCrowdListApi, gerCrowdInfoApi, updateCrowdContactApi, getOpenIdsApi } from '@/api/customer-beta-api'
import { formatDate } from '@/utils/tools'

// Component Object
Page({
  data: {
    /** 人群包列表 */
    crowdList: [],
    /** 当前人群包下标 */
    current: 0,
    /** 加载中 */
    showLoading: true,
    /** 当前人群包下的用户列表 */
    list: [],
    /** 用户列表 分页数据 */
    meta: {
      page: 1,
      size: 10,
      total: 0,
      status: ''
    },
    /** 用户列表 筛选条件 */
    searchParams: {
      name: '',
      contacted_status: 0,
      relationship_type: 0
    },
    /** 好友状态选项 */
    pickerOptions: [
      {
        label: '全部',
        value: 0
      },
      {
        label: '好友',
        value: 1
      },
      {
        label: '非好友',
        value: 2
      }
    ],
    /** 好友状态选择弹框 */
    showPicker: false,
    /** 好友状态下标 */
    value: 0,
    /** 人群包条件弹框 */
    showTips: false,
    /** 选择 */
    selection: {
      'crowd_id': {
        value: [],
        isSelectedAll: false,
        allIds: []
      }
    },
    /** 选择群发素材 */
    showSelectMaterialPopup: false
  },
  onLoad() {
    this.init()
    // 神策数据上报 -- 沟通名单页浏览
    wx.track.trackViewCommunicate()
  },
  async init() {
    try {
      await this.gerCrowdList()
      this.setData({
        showLoading: false
      })

      this.gerCrowdInfo()
      this.getOpenIds()
    } catch (error) {
      this.setData({
        showLoading: false
      })
    }
  },
  /**
   * @description 获取人群包列表
   */
  async gerCrowdList() {
    const { selection } = this.data
    const res = await gerCrowdListApi()
    const crowdList = res.data.map(item => {
      item.search.map = JSON.parse(item.search.map).map(val => {
        if (val.type === 1 || val.type === 2) {
          val.result_text = val.result.map(s => s.value).join('/')
        } else if (val.type === 3) {
          const { start_time, end_time } = val.result

          val.result_text = [formatDate(start_time * 1000, 'yyyy/MM/dd'), formatDate(end_time * 1000, 'yyyy/MM/dd')].join(' - ')
        } else {
          const { start_number, end_number } = val.result
          val.result_text = [start_number, end_number].join(' - ')
        }
        return val
      })

      if (!selection[item.crowd_id]) {
        selection[item.crowd_id] = {
          value: [],
          allIds: [],
          isSelectedAll: false
        }
      }

      return item
    })
    this.setData({
      crowdList,
      selection
    })
  },
  /**
   * @description 获取人群包明细
   * @param {*} page
   * @returns
   */
  async gerCrowdInfo(page = 1) {
    const { crowdList, current, searchParams, meta } = this.data
    const item = crowdList[current]
    if (!item) return
    try {
      if (meta.status === 'loading') return
      this.setData({
        'meta.status': 'loading'
      })
      const res = await gerCrowdInfoApi({
        page,
        size: meta.size,
        crowd_id: item.crowd_id,
        ...searchParams
      })
      if (res.code !== 200) {
        this.setData({
          'meta.status': 'loadError'
        })
        wx.showToast(res.message)
        return
      }
      const { list, total } = res.data
      this.data.list = page === 1 ? list : this.data.list.concat(list)

      this.setData({
        list: this.data.list,
        'meta.total': total,
        'meta.page': page,
        'meta.status': this.data.list.length < total ? 'loadMore' : 'loaded'
      })
    } catch (error) {
      this.setData({
        'meta.status': 'loadError'
      })
    }
  },
  /**
   * 获取所有可选的opened
   */
  async getOpenIds() {
    const { crowdList, current, selection } = this.data
    const { crowd_id } = crowdList[current]
    if (selection[crowd_id].allIds.length > 0) return
    try {
      const res = await getOpenIdsApi({
        crowd_id,
        contacted_status: 0,
        relationship_type: 1
      })
      if (res.code === 200) {
        this.setData({
          [`selection.${crowd_id}.allIds`]: res.data
        })
      }
    } catch (error) {
      console.log(error)
    }
  },
  /**
   * @description 加载更多
   */
  loadMore() {
    const { page, status } = this.data.meta
    if (status === 'loadMore') {
      this.gerCrowdInfo(page + 1)
    }
  },

  /**
   * @description 更新联系状态
   */
  async updateCrowdContact(e) {
    const item = e.detail
    if (item.is_contact === 1) return
    const { crowdList, current, list } = this.data
    const params = {
      crowd_id: crowdList[current].crowd_id,
      work_external_user_id: item.id
    }

    const res = await updateCrowdContactApi(params)

    if (res.code === 200) {
      const index = list.findIndex(val => val.id === item.id)
      this.setData({
        [`list[${index}].is_contact`]: 1
      })
    }
  },
  /**
   * @description 联系状态切换
   * @param {*} e
   */
  handleTabChange(e) {
    console.log(e.detail)
    this.setData({
      'searchParams.contacted_status': e.detail.name
    })
    this.gerCrowdInfo()

    // 神策数据上报 -- 沟通名单页浏览
    wx.track.trackCommunicateOperation({
      current_page: '沟通名单页',
      operation_type: e.detail.title
    })
  },
  /**
   * @description 搜索输入
   * @param {*} e
   */
  handleInput(e) {
    this.setData({
      'searchParams.name': e.detail.value
    })
  },
  /**
   * @description 确定搜索
   */
  handleSearch() {
    this.gerCrowdInfo()
  },
  /**
   * @description 切换侧边人群包
   * @param {*} e
   */
  handleMenuChange(e) {
    const { index } = e.currentTarget.dataset
    this.setData({
      current: index
    })
    this.gerCrowdInfo()
    this.getOpenIds()
  },
  /**
   * @description 去编辑人群包
   */
  handleEdit() {
    wx.navigateTo({
      url: '/packages/customer/pages/communication-edit/index',
      events: {
        update: () => {
          this.init()
        }
      }
    })

    // 神策数据上报 -- 沟通名单页浏览
    wx.track.trackCommunicateOperation({
      current_page: '沟通名单页',
      operation_type: '沟通名单编辑'
    })
  },
  /**
   * @description 选择好友状态
   */
  handlePicker() {
    this.setData({
      showPicker: true
    })
  },
  /**
   * @description 选择好友状态 取消
   */
  handlePickerCancel() {
    this.setData({
      showPicker: false
    })
  },
  /**
   * @description 选择好友状态 确定
   */
  handlePickerConfirm(e) {
    const { index, value } = e.detail

    this.setData({
      showPicker: false,
      value: index,
      'searchParams.relationship_type': value.value
    })
    this.gerCrowdInfo()
  },

  /**
   * @description 查看筛选条件
   */
  handleShow() {
    this.setData({
      showTips: true
    })
    // 神策数据上报 -- 沟通名单页
    wx.track.trackCommunicateOperation({
      current_page: '名单配置详情弹框',
      operation_type: '弹框显示'
    })
    // 神策数据上报 -- 名单配置详情弹窗
    wx.track.trackViewConfigureDetail()
  },
  /**
   * @description 关闭弹框
   */
  handleClose() {
    this.setData({
      showTips: false
    })
  },
  /**
   * @description 单个选择
   */
  handleSelect(e) {
    const { item, selected } = e.detail

    const { selection, crowdList, current } = this.data
    const { crowd_id } = crowdList[current]
    if (!selection[crowd_id]) {
      selection[crowd_id] = {
        value: [],
        allIds: [],
        isSelectedAll: false
      }
    }
    const ids = selection[crowd_id].value
    const index = ids.findIndex(id => item.openid === id)

    if (selected) {
      selection[crowd_id].value.splice(index, 1)
    } else {
      selection[crowd_id].value.push(item.openid)
    }

    selection[crowd_id].isSelectedAll = selection[crowd_id].value.length === selection[crowd_id].allIds.length

    this.setData({
      [`selection.${crowd_id}`]: selection[crowd_id]
    })
  },
  /**
   * @description 全选
   */
  handleSelectAll() {
    const { selection, crowdList, current } = this.data
    const { crowd_id } = crowdList[current]
    if (!selection[crowd_id]) {
      selection[crowd_id] = {
        value: [],
        allIds: [],
        isSelectedAll: false
      }
    }
    const value = !selection[crowd_id].isSelectedAll

    this.setData({
      [`selection.${crowd_id}.value`]: value ? [...selection[crowd_id].allIds] : [],
      [`selection.${crowd_id}.isSelectedAll`]: value
    })
  },
  /**
   * @description 选择群发素材
   */
  handleSelectMaterial() {
    this.setData({
      showSelectMaterialPopup: true
    })

    // 神策数据上报 -- 沟通名单页
    wx.track.trackCommunicateOperation({
      current_page: '沟通名单页',
      operation_type: '选择群发素材'
    })
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
   * @description 选择群发素材
   */
  async handleSelectMaterialConfirm(e) {
    const { type } = e.currentTarget.dataset
    const { selection, crowdList, current } = this.data
    const { crowd_id } = crowdList[current]
    const { value, allIds, isSelectedAll } = selection[crowd_id]
    try {
      // 获取外部联系人ID
      if (type === 1) {
        wx.navigateTo({
          url: '/pagesMaterial/pages/send-list/send-list?back=1',
          success: (result) => {
            result.eventChannel.emit('getExternalId', isSelectedAll ? allIds : value)
          }
        })
        // 神策数据上报 -- 创建群发页操作,选择群发素材,素材库选择
        // wx.track.trackBatchMsgAssetSelect({
        //   operation_type: 0
        // })
      } else {
        // 自定义素材
        this.shareToExternalContact(isSelectedAll ? allIds : value)
        // // 神策数据上报 -- 创建群发页操作,选择群发素材,自定义素材
        // wx.track.trackBatchMsgAssetSelect({
        //   operation_type: 1
        // })
        // // 神策数据上报 -- 群发素材
        // wx.track.trackBatchMsgConfirm({
        //   asset_coverage: ids.length,
        //   asset_list: '自定义素材无法获取',
        //   customer_list: ids.join(','),
        //   operation_type: 1
        // })
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
          wx.showToast('发送成功')
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
