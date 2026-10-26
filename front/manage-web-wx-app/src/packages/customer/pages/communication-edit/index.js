// Component Object

const { gerCrowdListApi, updateCrowdSortApi, deleteCrowdApi, updateCrowdApi } = require('@/api/customer-beta-api')

// Page Object
Page({
  data: {
    list: [],
    current: -1,
    showEditPopup: false,
    showDeletePopup: false,
    currentInfo: null
  },
  // options(Object)
  onLoad: function(options) {
    this.gerCrowdList()
    // 神策数据上报 -- 沟通名单页页面浏览
    wx.track.trackViewPackagePage()
  },
  onUnload() {
    if (this.data.update) {
      const eventChannel = this.getOpenerEventChannel()
      eventChannel.emit('update')
    }
  },
  async gerCrowdList() {
    const res = await gerCrowdListApi()
    console.log(res)
    this.setData({
      list: res.data
    })
  },
  handleDragBefore(e) {
    const { current } = e.detail
    console.log(current)

    this.setData({
      current
    })
    // 神策数据上报 -- 沟通名单页
    wx.track.trackCommunicateOperation({
      current_page: '打包名单',
      operation_type: '排序'
    })
  },
  handleDragChange(e) {
    const { list, current } = e.detail
    console.log(list, current)

    this.setData({
      list,
      current
    })
    const sort_list = list.map((item, index) => {
      return {
        sort: index,
        crowd_id: item.crowd_id
      }
    })

    updateCrowdSortApi({
      sort_list
    })
    this.data.update = true
  },
  handleEdit(e) {
    const { item } = e.currentTarget.dataset
    this.setData({
      currentInfo: item,
      showEditPopup: true
    })
    // 神策数据上报 -- 沟通名单页
    wx.track.trackCommunicateOperation({
      current_page: '打包名单',
      operation_type: '编辑'
    })
  },
  handleEditCancel() {
    this.setData({
      currentInfo: null,
      showEditPopup: false
    })
    // 神策数据上报 -- 沟通名单页
    wx.track.trackCommunicateOperation({
      current_page: '打包名单编辑弹窗',
      operation_type: '取消'
    })
  },
  async handleEditConfirm(e) {
    // 神策数据上报 -- 沟通名单页
    wx.track.trackCommunicateOperation({
      current_page: '打包名单编辑弹窗',
      operation_type: '确定修改'
    })
    const { value } = e.detail
    const { currentInfo } = this.data
    if (!value.trim()) {
      wx.showToast('沟通名单名称不能为空')
      return
    }

    const res = await updateCrowdApi({
      crowd_id: currentInfo.crowd_id,
      name: value
    })
    if (res.code === 200) {
      this.setData({
        currentInfo: null,
        showEditPopup: false
      })
      this.gerCrowdList()
      this.data.update = true
      wx.showToast('更新成功')
    } else {
      wx.showToast(res.message)
    }
  },
  handleDelete(e) {
    const { item } = e.currentTarget.dataset
    this.setData({
      currentInfo: item,
      showDeletePopup: true
    })
    // 神策数据上报 -- 沟通名单页
    wx.track.trackCommunicateOperation({
      current_page: '打包名单',
      operation_type: '删除'
    })
  },
  handleDeleteCancel(e) {
    this.setData({
      showDeletePopup: false
    })
    // 神策数据上报 -- 沟通名单页
    wx.track.trackCommunicateOperation({
      current_page: '打包名单删除弹窗',
      operation_type: '取消'
    })
  },
  async handleDeleteConfirm() {
    // 神策数据上报 -- 沟通名单页
    wx.track.trackCommunicateOperation({
      current_page: '打包名单删除弹窗',
      operation_type: '删除'
    })
    const { currentInfo } = this.data
    const res = await deleteCrowdApi({
      crowd_id: currentInfo.crowd_id
    })
    if (res.code === 200) {
      this.setData({
        showDeletePopup: false
      })
      this.data.update = true
      this.gerCrowdList()
      wx.showToast('删除成功')
    }
  }
})
