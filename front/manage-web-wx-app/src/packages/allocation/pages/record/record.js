import { getInheritRecordListApi } from '@/api/inherit-api'
import { toDate } from '@/utils/tools'

// 创建页面实例对象
Page({
  data: {
    tabs: [
      {
        status: -1,
        title: '全部'
      },
      {
        status: 0,
        title: '分配中'
      },
      {
        status: 1,
        title: '分配成功'
      },
      {
        status: 2,
        title: '分配失败'
      }
    ],
    active: 0,
    total: {
      on: 0,
      err: 0,
      success: 0
    }, // 页面信息
    recordList: [], // 记录列表
    dateFilter: [], // 日期筛选
    statusFilter: -1, // 状态筛选
    meta: {
      page: 1,
      size: 10,
      total: 0,
      isMore: false
    },
    isLoading: false,
    isShowDate: false, // 是否展示时间筛选
    /** 时间选择配置 最大时间 */
    maxDate: new Date().getTime(),
    /** 时间选择配置 最小时间 */
    minDate: new Date().setFullYear(new Date().getFullYear() - 3),
    create_date_text: '' // 日期筛选文本
  },
  onLoad () {
    this.getRecordList()
  },
  onShow () {
  },
  /**
   * @description 切换tab
   */
  handleTabChange(e) {
    const { index } = e.detail
    this.setData({
      active: index,
      statusFilter: this.data.tabs[index].status
    }, () => {
      this.getRecordList()
    })
  },
  /**
   * 获取分配列表
   */
  async getRecordList (page = false) {
    const { meta, recordList, statusFilter, dateFilter } = this.data

    wx.showNavigationBarLoading()

    this.setData({ isLoading: true })

    if (page) {
      this.setData({ 'meta.page': meta.page + 1 })
    } else {
      this.setData({ 'meta.page': 1, recordList: [] })
    }

    const params = {
      page: meta.page,
      size: meta.size,
      created_at: dateFilter.map(i => parseInt(i / 1000))
    }

    if (statusFilter > -1) {
      params.status = statusFilter
    }

    try {
      const { code, data, message } = await getInheritRecordListApi(params)
      wx.hideNavigationBarLoading()
      this.setData({ isLoading: false })

      if (code !== 200) {
        wx.showToast(message)
        return
      }

      const { list, top, total } = data

      const newList = list.map(i => ({
        ...i,
        created_at: toDate(i.created_at, 'yyyy/MM/dd hh:mm:ss').nowTime,
        updated_at: toDate(i.updated_at, 'yyyy/MM/dd hh:mm:ss').nowTime
      }))

      const newRecordList = page ? recordList.concat(newList) : newList

      this.setData({
        recordList: newRecordList,
        total: top,
        'meta.total': total,
        'meta.isMore': newRecordList.length < total
      })
    } catch (error) {
      console.error()
      wx.hideNavigationBarLoading()
      this.setData({ isLoading: false, 'meta.page': Math.max(1, meta.page - 1) })
    }
  },
  /**
   * 触底翻页
   */
  async onReachBottom () {
    let { meta, isLoading } = this.data
    if (!meta.isMore || isLoading) return

    this.getRecordList(true)
  },
  /**
   * 打开时间筛选
   */
  openDateSelect () {
    this.setData({
      isShowDate: true,
      dateFilter: this.data.dateFilter
    })
  },
  /**
   * 时间选择
   */
  dateConfirm (e) {
    console.log('时间选择', e)
    const {value} = e.detail
    // 获取时间
    const create_date = value.map(i => {
      return toDate(i, 'yyyy/MM/dd').nowTime
    })
    const create_date_text = create_date.join('-')
    this.setData({
      dateFilter: value,
      // 将.替换成/
      create_date_text: create_date_text.replace(/\./g, '/')
    })
    this.getRecordList()
    // 关闭时间弹窗
    this.dateCancel()
  },
  /**
   * @description: 关闭时间提示tips
   */
  closeTips () {
    this.selectComponent('#mp-tooltip').onClickOverlay()
  },
  /**
   * 时间选择重置
   */
  dateReset () {
    this.setData({
      dateFilter: []
    })
    this.getRecordList()
    // 关闭时间弹窗
    this.dateCancel()
  },
  /**
   * 时间选择关闭
   */
  dateCancel() {
    this.setData({
      isShowDate: false
    })
  },
  /**
   * 状态选择
   */
  statusConfirm (e) {
    console.log('状态选择', e)
    this.setData({
      statusFilter: e.detail
    })
    this.getRecordList()
  },
  /**
   * 跳转详情页
   */
  goDetail (e) {
    const { item } = e.currentTarget.dataset
    const params = encodeURIComponent(JSON.stringify(item))
    let url = `/packages/allocation/pages/record-detail/record-detail?routerParams=${params}`
    wx.navigateTo({
      url: url
    })
  }
})
