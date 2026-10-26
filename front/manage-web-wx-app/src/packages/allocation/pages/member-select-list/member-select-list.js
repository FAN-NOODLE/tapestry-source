import { getInheritMemberListApi, getInheritTransferInfoApi } from '@/api/inherit-api'


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
    /** 搜索条件 */
    search: '',
    /** 选中的顾问 */
    checkedInfo: null,
    /** 顾客的筛选条件 */
    customParams: {},
    /** 选中的顾客数量 */
    checkedCount: 0
  },
  /*    系统事件开始   */
  /**
   * onLoad
   * @param {Object} options
   * @param {String} options.customParams 自定义筛选参数
   */
  onLoad (options) {
    if (options.customParams) {
      this.setData({
        checkedCount: options.checkedCount,
        customParams: JSON.parse(options.customParams)
      })
    }
    wx.showLoading({
      title: '加载中'
    })
    this.getData()
  },
  onReachBottom () {
    this.loadMore()
  },
  /*    系统事件结束   */

  /*    普通事件开始  */
  /**
   * @description 获取顾客列表
   */
  async getData (page = 1) {
    const {
      meta: { size, status },
      customParams,
      search
    } = this.data

    if (status === 'loading') return
    this.setData({
      'meta.status': 'loading'
    })

    wx.showNavigationBarLoading()

    const params = {
      page,
      size,
      search,
      ...customParams
    }
    try {
      const { data, code, message } = await getInheritMemberListApi(params)
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
   * 神策数据上报 -- 顾客接替 转移到后端
   * @param {Number} inherit_id 分配记录id
   */
  trackEventData (inherit_id) {
    const { selectCustomerList, selectedMember } = this.data
    if (!wx.track) return
    selectCustomerList.forEach(item => {
      // 神策数据上报 -- 顾客分配
      wx.track.trackAssignCustomer({
        inherit_id: inherit_id,
        user_external_id: item.openid,
        handover_userid: item.sa_user_id,
        takeover_userid: selectedMember.member_id
      })
      // 神策数据上报 -- 顾客被分配（SA）
      wx.track.trackCustomerAssigned({
        inherit_id: inherit_id,
        user_external_id: item.openid,
        takeover_userid: selectedMember.member_id
      })
    })
  },
  /*    普通事件结束  */

  /*    行为事件开始   */
  /**
   * @description 搜索确认
   */
  handleSearchConfirm(e) {
    this.setData({
      search: e.detail.value
    })
    wx.showLoading({
      title: '加载中'
    })
    this.getData()
  },
  /**
   * @description 选择切换
   */
  handleCheckOne(e) {
    const { item } = e.currentTarget.dataset
    // todo: 选中的顾问不能是自己
    if (item.check) {
      this.setData({
        checkedInfo: item
      })
    }
  },
  /**
   * @description 返回
   */
  handleBack() {
    wx.navigateBack({
      delta: 1
    })
  },
  /**
   * @description 交转移
   * @returns
   */
  handleSend() {
    const { checkedInfo, checkedCount } = this.data
    if (!checkedInfo) {
      wx.showToast('请选择顾问')
      return
    }
    this.selectComponent('#com-send-popup').show(checkedInfo.name, checkedCount)
  },
  /**
   * @description 确认转移
   * @param {*} e
   */
  async handleSendConfirm(e) {
    const { checkedInfo, customParams } = this.data
    const { content } = e.detail
    const params = {
      ...customParams,
      takeover_member_id: checkedInfo.member_id,
      transfer_success_msg: content
    }
    console.log(params)
    wx.showLoading({title: ''})
    try {
      const { code, data, message } = await getInheritTransferInfoApi(params)
      wx.hideLoading()
      if (code !== 200) {
        wx.showToast(message)
      }
      wx.navigateBack({
        delta: 2
      })
      const eventChannel = this.getOpenerEventChannel()
      eventChannel.emit('success')
      this.selectComponent('#com-send-popup').hide()
      // 神策数据上报 -- 顾客接替 转移到后端
      // this.trackEventData(data.id)
      console.log(data)
    } catch (error) {
      wx.hideLoading()
      console.log(error)
    }
  }
  /*    行为事件结束   */
})
