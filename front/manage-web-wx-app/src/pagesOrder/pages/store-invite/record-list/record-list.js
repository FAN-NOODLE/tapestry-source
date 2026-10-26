import { orderListApi } from '../../../api/guide-invite'
import { StatusList } from '../constants'

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    title: 'record-list',
    list: [],
    selectDate: '',
    statusIndex: -1,
    statusList: StatusList,
    pageLoading: true, // 页面首屏渲染
    isLoading: false, // 分页loading
    page: 1,
    hasMore: true,
    dateModel: false,
    showModel: false,
    showLoading: true
  },

  // 生命周期函数--监听页面加载
  onLoad () {
    this.getList()
    wx.hideHomeButton()
  },

  /** 重置 */
  searchReset () {
    this.setData({
      selectDate: '',
      statusIndex: -1
    })
    this.reset()
  },

  // 点击跳转详情
  goDetail (e) {
    const item = e.currentTarget.dataset.item
    const self = this
    wx.navigateTo({
      url: `/pagesOrder/pages/store-invite/invite-detail/invite-detail?id=${item.id}&reserve_number=${item.reserve_number}`,
      events: {
        // 为指定事件添加一个监听器，获取被打开页面传送到当前页面的数据
        acceptDataFromOpenedPage: function (res) {
          console.log('acceptDataFromOpenedPage', res, item.status)
          // 状态与当前不一致
          if (item.status !== res.status) {
            self.reset()
          }
        }
      }
    })
  },
  // 重置列表数据
  reset (pullDown = false) {
    this.setData(
      {
        pageLoading: !pullDown,
        list: [],
        page: 1
      },
      () => {
        this.getList(pullDown)
      }
    )
  },
  /**
   * 下拉刷新
   */
  onPullDownRefresh () {
    console.log('我下拉刷新了')
    this.reset(true) // 请求当前页面数据
  },
  // 获取预约记录列表
  getList (pullDown = false) {
    this.setData({
      isLoading: true
    })
    const { statusIndex, page, selectDate } = this.data
    const params = {
      page,
      size: 10
    }
    if (selectDate) {
      const date = selectDate.replace(/-/g, '/')
      params.start_time = new Date(`${date} 00:00:00`).getTime() / 1000
    }
    if (statusIndex > -1) {
      params.status = this.data.statusList[Number(statusIndex)].id
    }

    orderListApi(params)
      .then((res) => {
        pullDown && wx.stopPullDownRefresh()
        if (res.code === 200) {
          const { total, page, size } = res.data
          let list = res.data.list || []
          this.setData({
            list: [...this.data.list, ...list],
            hasMore: page * size < total,
            pageLoading: false,
            isLoading: false
          })
        } else {
          wx.showToast({
            icon: 'none',
            title: res.message
          })
          this.setData({
            pageLoading: false,
            isLoading: false
          })
          if (this.data.page > 1) this.data.page -= 1
        }
      })
      .catch((err) => {
        pullDown && wx.stopPullDownRefresh()
        wx.showToast({
          icon: 'none',
          title: err.message
        })
        this.setData({
          pageLoading: false,
          isLoading: false
        })
        if (this.data.page > 1) this.data.page -= 1
      })
  },
  // 筛选条件变化
  changeKey (e) {
    const key = e.currentTarget.dataset.key
    const value = e.detail.value
    console.log(key, value)
    this.setData({
      [key]: value
    })
    this.reset()
  },
  // 列表滑动到底部 -滚动加载
  onReachBottom () {
    if (this.data.hasMore) {
      this.data.page += 1
      this.getList()
    }
  },
  // 点击跳转评价页面
  onComments (e) {
    const { id, reserve_number } =
      e.currentTarget && e.currentTarget.dataset.item
    wx.redirectTo({
      url: `/pagesOrder/pages/store-invite/comment/comment?id=${id}&reserve_number=${reserve_number}&notShare=true`
    })
  }
})
