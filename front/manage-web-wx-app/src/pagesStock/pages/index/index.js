import { getStoreListApi, getSearchHistoryApi, deleteSearchHistoryApi, searchListGroupSkuApi } from '@/api/store-api'
import { groupArr } from '../../utils/index'
import { debounce } from '@/utils/tools'
// 获取全局应用程序实例对象
const app = getApp()

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    title: 'index',
    store: null,
    /** tab列表 */
    tabs: [
      {
        title: '商品款号查询'
      },
      {
        title: '条形码查询'
      }
    ],
    /** tab下标 */
    active: 0,
    searchValue: '',
    storeCityList: [],
    storeList: [],
    showWarnPopup: false,
    showStorePopup: false,
    historyList: [],
    historyEdit: false,
    showModal: false,
    showLoading: true,
    showSearchLoading: false,
    skuList: []
  },

  // 生命周.getStoreList
  async onShow() {
    // TODO: 判断localstorage中是否有storeInfo 有则赋值给store 没有则获取门店
    const storeInfo = wx.getStorageSync('storeInfo')
    if (storeInfo) {
      this.setData({
        store: storeInfo,
        showLoading: false
      })
    } else {
      await this.getStoreList()
    }
    this.getHistoryList()
  },
  onUnload() {
    // 清空门店信息,重新读取全局门店
    wx.removeStorageSync('storeInfo')
  },
  /**
   * @de.getStoreList
   * @return {*}
   */
  getStoreList() {
    const { store } = this.data
    getStoreListApi().then(res => {
      const storeList = res.data.list
      const {storeInfo} = app.globalData
      console.log(storeInfo)
      if (storeInfo.store_code) {
        this.setData({
          store: storeInfo
        })
      } else if (!store || !store.store_code) {
        this.setData({
          showWarnPopup: true
        })
      }
      // 创建树形结构,将门店按照城市分组
      const treeData = groupArr(storeList, 'city')
      this.setData({
        storeCityList: treeData,
        storeList,
        showLoading: false
      })
    }).catch(err => {
      console.log(err)
      this.setData({
        showLoading: false
      })
    })
  },
  /**
   * @description 切换tab
   */
  handleTabChange(e) {
    const { index } = e.detail
    this.setData({
      active: index,
      searchValue: ''
    })
  },
  /**
   * @description:skuitem查询赋值
   */
  changeInput (e) {
    const value = e.detail.value
    this.setData({
      searchValue: value.trim()
    })

    if (this.data.active === 1) return
    if (value.trim().length < 4) {
      this.setData({
        showResult: false,
        skuList: []
      })
      return
    }
    this.searchListGroupSku()
  },
  searchListGroupSku: debounce(async function () {
    try {
      this.setData({
        showSearchLoading: true
      })
      const res = await searchListGroupSkuApi({
        page: 1,
        size: 50,
        sku_item: this.data.searchValue,
        include_none_image_product: true
      })
      console.log(res)
      if (res.code === 200) {
      // const search = res.data.find(item => item === this.data.searchValue)
        this.setData({
          skuList: res.data,
          showResult: res.data.length > 0,
          showSearchLoading: false
        // searchValue: search || this.data.searchValue
        })
      } else {
        this.setData({
          showSearchLoading: false
        })
      }
    } catch (error) {
      this.setData({
        showSearchLoading: false
      })
    }
  }, 500, false),
  handleCloseResult() {
    this.setData({
      showResult: false
    })
  },
  handleSelectSkuItem(e) {
    const { item } = e.currentTarget.dataset
    this.setData({
      searchValue: item,
      skuList: []
    })
    this.handleSubmit(e, false)
  },
  /**
   * @description: 获取相机
   */
  scanCode(e) {
    const that = this
    wx.scanCode({
      success (res) {
        console.log(res)
        const code = res.result
        that.setData({
          searchValue: code
        })
        that.handleSubmit(e, true)
      }
    })
  },
  /**
   * @description: 立即查询
   */
  handleSubmit(e, isScan = false) {
    const {store, searchValue, active} = this.data
    if (!searchValue) {
      return
    }
    // 判断searchValue是否大于3个字符
    // if (active === 0 && searchValue.length < 3) {
    //   wx.showToast('至少输入3位字符')
    //   return
    // }
    const params = {
      store,
      searchValue,
      active,
      isScan
    }
    // 转成字符串
    const paramsStr = encodeURIComponent(JSON.stringify(params))
    const that = this
    wx.navigateTo({
      url: `/pagesStock/pages/list/list?from=1&searchParams=${paramsStr}`,
      success: function() {
        // 清空搜索框
        that.setData({
          searchValue: ''
        })
      }
    })

    // 神策数据上报 -- 库存查询
    wx.track.trackStockInquiry({
      entry_point: 0,
      inquiry_type: active,
      search_type: '自然搜索',
      is_have_historical_search_word: !!this.data.historyList[active].list.find(item => item.keywords === searchValue)
    })
  },
  /**
   * @description: 页面返回
   */
  handleBack() {
    // 清空门店信息,重新读取全局门店
    wx.removeStorageSync('storeInfo')
  },
  /**
   * @description: 显示选择门店弹窗
   */
  handleShowStorePopup() {
    this.setData({
      showStorePopup: true,
      showWarnPopup: false
    })
  },
  /**
   * @description: 选择门店
   */
  handleConfirmStore(e) {
    const item = e.detail
    this.setData({
      store: {
        store_inventory_name: item.inventory_name,
        store_code: item.third_store_code,
        store_name: item.store_name,
        store_city: item.city,
        store_phone: item.phone
      },
      showStorePopup: false
    })
    // 将门店信息存入localstorage
    wx.setStorageSync('storeInfo', this.data.store)
  },
  /**
   * @description: 关闭门店弹窗
   */
  handleCloseStore() {
    this.setData({
      showWarnPopup: true
    })
  },

  /**
   * @description: 获取历史查询记录
   */
  async getHistoryList() {
    const res = await getSearchHistoryApi()
    const list = res.data.length ? res.data : []
    const itemHistoryList = list.filter(item => item.keyword_type === 'sku_item')
    const upcHistoryList = list.filter(item => item.keyword_type === 'upc_code')
    this.setData({
      historyList: [
        {
          edit: false,
          list: itemHistoryList
        },
        {
          edit: false,
          list: upcHistoryList
        }
      ]
    })
  },
  /**
   * @description: 管理历史查询记录
   */
  async handleEdit() {
    const { active, historyList } = this.data
    this.setData({
      [`historyList[${active}].edit`]: !historyList[active].edit
    })
  },
  /**
   * @description: 管理历史查询记录
   */
  async handleDeleteHistory(e) {
    const { historyList, active } = this.data
    const list = historyList[active].list
    const { index = -1 } = e.currentTarget.dataset
    const params = {
      ids: [],
      is_all: 0
    }
    if (index === -1) {
      params.is_all = active + 1
    } else {
      params.ids = [list[index].id]
    }
    await deleteSearchHistoryApi(params)
    if (index === -1) {
      this.getHistoryList()
    } else {
      historyList[active].list.splice(index, 1)
      this.setData({
        [`historyList[${active}].edit`]: historyList[active].list.length > 0,
        [`historyList[${active}].list`]: historyList[active].list
      })
    }
  },

  handleDelete() {
    this.setData({
      showModal: true
    })
  },
  handleDeleteCancel() {
    this.setData({
      showModal: false
    })
  },
  handleDeleteConfirm(e) {
    this.setData({
      showModal: false
    })
    this.handleDeleteHistory(e)
  },

  handleHistoryClick(e) {
    const { store, historyList, active } = this.data
    const { index } = e.currentTarget.dataset
    const item = historyList[active].list[index]

    const params = {
      store,
      searchValue: item.keywords,
      active: item.keyword_type === 'sku_item' ? 0 : 1
    }
    // 转成字符串
    const paramsStr = encodeURIComponent(JSON.stringify(params))
    wx.navigateTo({
      url: `/pagesStock/pages/list/list?from=1&searchParams=${paramsStr}`
    })

    // 神策数据上报 -- 库存查询
    wx.track.trackStockInquiry({
      entry_point: 0,
      inquiry_type: params.active,
      search_type: '历史搜索',
      is_have_historical_search_word: true
    })
  }
})
