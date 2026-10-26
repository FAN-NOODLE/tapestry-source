import { getStoreListApi } from '@/api/store-api'
import Enums from '../../utils/enums'

const { showToast, globalData } = getApp()
// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    loading: true,
    list: [],
    pageSet: {
      page: 0,
      size: 10,
      hasMore: true,
      loading: true
    },
    currentPage: ''
  },
  onLoad (option) {
    if (option && option.currentPage) {
      this.setData({
        currentPage: option.currentPage
      })
    }
    this.getStore()
  },
  onShow () {
    if (globalData.store && this.data.list.length > 0) {
      const list = this.data.list
      list.forEach((item) => {
        item.select = false
        item.id === globalData.store.id && (item.select = true)
      })
      this.setData({ list })
    }
  },
  async getStore () {
    try {
      const res = await getStoreListApi({
        page: ++this.data.pageSet.page,
        size: this.data.pageSet.size
      })
      this.setData({ loading: false })
      if (res.code === 200) {
        this.setData({
          list: [...this.data.list, ...res.data.list],
          'pageSet.loading': false,
          // 'pageSet.hasMore': res.data.total > res.data.page * res.data.size,
          'pageSet.hasMore':
            res.data.total > this.data.pageSet.page * this.data.pageSet.size
        })
      } else {
        showToast(res.message)
        this.setData({
          'pageSet.loading': false,
          loading: false
        })
      }
    } catch (error) {
      console.error(error)
      this.setData({
        'pageSet.loading': false,
        loading: false
      })
    }
  },
  select (e) {
    const { list } = this.data
    list.forEach((item) => {
      item.select = false
      item.id === e.detail.id && (item.select = true)
    })
    this.setData({ list })
    // globalData.store = e.detail
    // globalData.store.store_id = e.detail.id
    let canSwithPage = ['pages/task/task', 'pages/compass-data/compass-data', 'pages/customer-list/customer-list', 'pages/my/my']
    if (this.data.currentPage !== '') {
      if (canSwithPage.some((item) => item === this.data.currentPage)) {
        wx.redirectTo({
          url: '/pages/empty/empty?currentPage=' + this.data.currentPage
        })
      } else {
        wx.redirectTo({
          url: '/' + this.data.currentPage
        })
      }
    } else if (globalData.identityType === Enums.identityType.partner) {
      console.log('跳转至品牌合伙人')
    } else {
      wx.switchTab({
        url: '/pages/compass-data/compass-data'
      })
    }
  },
  computed: {
    total () {
      return this.data.pageSet.size
    }
  },
  watch: {
    loading: {
      handler (k, value, newVal) {
      }
    }
  },
  onReachBottom () {
    if (!this.data.pageSet.loading && this.data.pageSet.hasMore) {
      this.setData({
        'pageSet.loading': true
      })
      this.getStore()
    }
  }
})
