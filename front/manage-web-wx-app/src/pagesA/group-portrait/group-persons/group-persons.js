import { chatGroupMembers } from '@/api/group-portrait-api'

const TAB_LIST = [
  {label: '全部', value: 1},
  {label: '好友', value: 2},
  {label: '非好友', value: 3}
]

/** 设置分页数据 */
const generatePageInfo = (info = {}) => {
  return {
    page: 1,
    size: 10,
    hasMore: true,
    ...info
  }
}

Page({
  data: {
    /** 群人员tab 切换 */
    tabList: TAB_LIST,
    /** 选中index */
    tabIndex: 1,
    /** 分页数据 */
    pageInfo: generatePageInfo(),
    /** 好友列表 */
    list: [],
    /** 社群id  wr9hx_EQAAKjtYaKg0IS7WkbtXU26JLQ */
    chatId: ''
  },
  onLoad (option) {
    this.setData({
      chatId: option.chatId
    })
    this.getMemberList()
  },
  /** 页面触底上拉加载更多 */
  onReachBottom () {
    const { pageInfo } = this.data
    if (!pageInfo.hasMore) return
    pageInfo.page += 1
    this.getMemberList()
  },

  /** tab 切换 */
  tabChange (e) {
    this.setData({
      tabIndex: e.detail,
      list: [],
      pageInfo: generatePageInfo()
    })

    this.getMemberList()
  },


  /** 群成员列表 */
  async getMemberList () {
    wx.showLoading({
      title: '加载中'
    })

    const {pageInfo, tabIndex} = this.data
    const params = {
      chat_id: this.data.chatId,
      page: pageInfo.page,
      size: pageInfo.size,
      friend: tabIndex - 1,
      list_type: 0
    }
    const {
      code,
      message = '',
      data
    } = await chatGroupMembers(params)

    wx.hideLoading()
    if (code === 200) {
      if (data.items && data.items.length) {
        this.setData({
          'pageInfo.hasMore': data.total_count > pageInfo.page * pageInfo.size,
          list: pageInfo.page === 1 ? data.items : [...this.data.list, ...data.items]
        })
      }
    } else {
      wx.showToast(message)
    }
  }
})
