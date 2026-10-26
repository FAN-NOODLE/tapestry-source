import {
  orderMemberApi
} from '../../api/order-api'
import { debounce } from '../../utils/tools'

const { globalData, showToast } = getApp()

let searchKeywords = ''

Page({
  // 页面的初始数据
  data: {
    currentEmptyText: '暂未搜索到相关导购～', // 当前列表缺省文案
    ready: false,
    searchResult: [],
    searching: false,
    searchKeywords: '',
    hasMore: false,
    meta: {
      page: 0,
      size: 10
    },
    orderType: 1, // 订单类型：1推广订单 2服务订单
    serviceOrderType: 2, // 服务订单类型：1待接订单 2:我的订单 3：可抢单,
    memberId: 0, // 员工id,
    check_all: 0, // 0:未选中 1:选中
    isShowMore: false // 是否展示更多
  },
  // 生命周期函数--监听页面加载
  onLoad (options) {
    this.setData({
      serviceOrderType: 2,
      memberId: options.memberId,
      memberName: options.memberName,
      check_all: options.check_all
    })
    searchKeywords = ''
    this.getSearchResult()
  },
  /**
   * 搜索搜索接口
   * keyword 关键词
   * */
  async getSearchResult () {
    try {
      let params = {
        keywords: searchKeywords,
        page: ++this.data.meta.page,
        size: this.data.meta.size,
        enterprise_id: globalData.enterprise[0].enterprise_id || 0,
        store_id: globalData.store.store_id || 0
      }
      let list = []
      const { code, message, data } = await orderMemberApi(params)
      if (parseInt(code) === 200) {
        list = data.list.map((d) => {
          return {
            myself: d.id === globalData.memberInfo.memberId,
            ...d
          }
        })
        this.setData({
          searchResult: this.data.isShowMore ? [...this.data.searchResult, ...list] : list,
          ready: true,
          searching: false
        }, () => {
          this.setData({
            hasMore: data.total > data.page * data.size
          })
        })
      } else {
        this.setData({
          ready: true
        })
        showToast(message)
      }
    } catch (error) {
      this.setData({
        ready: true
      })
      console.error('guideSelect-getSearchResult', error)
    }
  },
  clearSearchResult () {
    return new Promise(resolve => {
      this.setData({
        searchResult: []
      }, resolve)
    })
  },
  // 输入框改变
  onChange: debounce(async function (e) {
    console.log(e, '朱二龙')
    const value = e.detail.value
    console.log(value, '朱二龙')
    this.setData({
      value,
      ready: false,
      'meta.page': 0
    })

    if (!value) {
      this.setData({
        searchResult: [],
        searching: false
      })
      return
    }
    // 请求接口
    searchKeywords = value
    await this.clearSearchResult()
    await this.getSearchResult()
  }, 500),
  // 清除输入框值
  async clearValue () {
    this.setData({
      value: '',
      'meta.page': 0
    })
    searchKeywords = ''
    await this.clearSearchResult()
    await this.getSearchResult()
  },
  // 选择导购
  selectGuide (e) {
    let memberId = e.currentTarget.dataset.id
    let memberName = e.currentTarget.dataset.name
    let pages = getCurrentPages()// 获取页面栈
    let prevPage = pages[pages.length - 2] // 上一个页面
    // 直接调用上一个页面的setData()方法，把数据存到上一个页面中去
    const isSelectSelf = memberId === globalData.memberInfo.memberId
    prevPage.setData({
      memberId,
      memberName: isSelectSelf ? `${memberName}(自己)` : memberName,
      check_all: 0,
      isFilterSelectSelf: isSelectSelf
    })
    wx.navigateBack({
      delta: 1
    })
  },
  // 全选
  selectAll () {
    let pages = getCurrentPages()// 获取页面栈
    let prevPage = pages[pages.length - 2] // 上一个页面
    // 直接调用上一个页面的setData()方法，把数据存到上一个页面中去
    prevPage.setData({
      memberId: 0,
      memberName: '全部导购',
      check_all: 1,
      isFilterSelectSelf: false
    })
    wx.navigateBack({
      delta: 1
    })
  },
  // 触底翻页
  async onReachBottom () {
    if (!this.data.hasMore) return
    this.setData({
      isShowMore: true
    })
    this.getSearchResult()
  }
})
