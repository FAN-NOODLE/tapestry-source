/* eslint-disable standard/object-curly-even-spacing */
// 获取全局应用程序实例对象
// const app = getApp();
import { getInventoryListApi, getFirstOpenApi } from '@/api/stock-api'
import { skuItemHeader, upcCodeHeader } from '../../utils/data'
import { groupArr } from '@/pagesStock/utils'
import { getStoreListApi } from '@/api/store-api'
import { getRect } from '@/pagesMaterial/utils/util'
const {
  normalizeInventoryFailure,
  normalizeInventoryPayload
} = require('../../utils/inventory-status')
// 获取全局应用程序实例对象
const app = getApp()

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    store: null,
    active: 0,
    searchValue: '',
    colorValue: 'All Colors',
    sizeValue: 'All sizes',
    checkActive: false,
    tableHeader: skuItemHeader,
    tableData: [],
    copyTableData: [],
    showLoading: false,
    colorList: [],
    sizeList: [],
    showSelect: false,
    colorIndex: 0,
    sizeIndex: 0,
    type: '',
    showListTips: false,
    showPreview: false,
    previewImg: '',
    tableArea: 0,
    isSelect: false,
    hasMore: true,
    pageSize: 30,
    // 选择门店
    storeCityList: [],
    storeList: [],
    showWarnPopup: false,
    showStorePopup: false,
    isHorizontal: false,
    inventoryWarnings: [],
    failedUPCs: [],
    emptyMessage: ''
  },

  // 生命周期函数--监听页面加载
  async onLoad(options) {
    // TODO: onLoad
    console.log(options)
    if (options && options.searchParams) {
      const params = decodeURIComponent(options.searchParams)
      const routeParams = JSON.parse(params)
      this.setData({
        ...routeParams,
        active: Number(routeParams.active),
        from: options.from
      })
      // 设置不同的tableHeader
      if (this.data.active === 0) {
        this.setData({
          tableHeader: skuItemHeader
        })
      } else if (this.data.active === 1) {
        this.setData({
          tableHeader: upcCodeHeader
        })
      }
      // 如果没有门店信息，获取门店信息
      if (!routeParams.store) {
        const storeInfo = wx.getStorageSync('storeInfo')
        if (storeInfo) {
          this.setData({
            store: storeInfo
          })
        } else {
          await this.getStoreList()
        }
      }
    }
    if (this.data.store) {
      this.getInventoryList()
    }
    // getFirstOpenApi({
    //   type: 10
    // })
  },

  onUnload() {
    // 清空门店信息,重新读取全局门店
    wx.removeStorageSync('storeInfo')
  },
  /**
   * @description: 判断是否是首次打开
   */
  getFirstOpen() {
    getFirstOpenApi({
      type: 1
    }).then(res => {
      const data = res.data
      if (data.is_first_open) {
        // 获取tableDataArea的高度
        const query = wx.createSelectorQuery()
        query.select('#tableDataArea').boundingClientRect()
        query.exec(res => {
          console.log(res)
          const tableArea = res[0]
          this.setData({
            tableArea,
            showListTips: true
          })
        })
      }
    }).catch(err => {
      console.log(err)
    })
  },
  /**
   * @description:获取库存列表
   */
  getInventoryList() {
    console.log(1111, this.data.store)
    this.setData({
      showLoading: true
    })
    // 测试用参数
    // let params = {
    //   'sku_item': 'CE772',
    //   'store_code_arr': [
    //     'OC245'
    //   ]
    // }
    let params = {
      'store_code_arr': [
        this.data.store.store_code
      ],
      scan: !!this.data.isScan
    }
    // todo: 根据active判断是skuitem查询还是upc查询
    if (this.data.active === 0) {
      params.sku_item = this.data.searchValue
    } else if (this.data.active === 1) {
      params.upc_code = this.data.searchValue
    }
    getInventoryListApi(params).then(res => {
	  if (!res || Number(res.code) !== 200) {
	    throw new Error((res && res.message) || 'inventory request failed')
	  }
	  const normalized = normalizeInventoryPayload(res.data)
	  let data = normalized.list
      // todo: 处理数据,获取所有的颜色和尺寸去重
      const colorList = [...new Set(data.map(item => item.color))]
      const sizeList = [...new Set(data.map(item => item.style))]
      // todo: 将all colors和all sizes添加到数组的第一位
      colorList.unshift('All Colors')
      sizeList.unshift('All sizes')
      if (data && data.length > 0) {
        this.setData({
          isSelect: true
        })
      }
      this.setData({
        copyTableData: data,
        colorList,
        sizeList,
        inventoryWarnings: normalized.warnings,
        failedUPCs: normalized.failedUPCs,
        emptyMessage: normalized.emptyMessage
      })
      // 如果是筛选条件需要重新筛选
      this.loadMore(1)
      if (data && data.length > 0) {
        // 判断是否是首次打开
        this.getFirstOpen()
      }
    }).catch(err => {
      console.log(err)
      const normalized = normalizeInventoryFailure()
      this.setData({
        showLoading: false,
        copyTableData: normalized.list,
        tableData: normalized.list,
        inventoryWarnings: normalized.warnings,
        failedUPCs: normalized.failedUPCs,
        emptyMessage: normalized.emptyMessage
      })
    })
  },
  /**
   * @description: 筛选方法
   */
  handleSelect() {
    return new Promise((resolve) => {
      let data = this.data.copyTableData
      // todo: 点击只看有库存
      if (this.data.checkActive) {
        data = data.filter(item => item.on_hand + item.in_transit > 0)
      }
      // todo: 判断是否选择了颜色
      if (this.data.colorValue !== 'All Colors') {
        data = data.filter(item => item.color === this.data.colorValue)
      }
      // todo: 判断是否选择了尺寸
      if (this.data.sizeValue !== 'All sizes') {
        data = data.filter(item => item.style === this.data.sizeValue)
      }
      // todo: 判断是否点击了排序
      const sortItem = this.data.tableHeader.find(item => item.sort_type)
      if (sortItem) {
        const sortType = sortItem.sort_type
        data = data.sort((a, b) => {
          return sortType === 'desc' ? a[sortItem.prop] - b[sortItem.prop] : b[sortItem.prop] - a[sortItem.prop]
        })
      }
      resolve(data)
    })
  },
  /**
   * @description: 对tableData进行分页
   */
  loadMore(page) {
    this.handleSelect().then((data) => {
      if (page === 1) {
        this.setData({
          showLoading: true,
          tableData: [],
          hasMore: true
        })
        // this.data.hasMore = true
        // this.data.tableData = []
        // this.data.showLoading = true
      }
      const { tableData, pageSize } = this.data
      const start = (page - 1) * pageSize
      const end = page * pageSize
      const sliceData = data.slice(start, end)
      // 判断是否还有数据
      if (sliceData && sliceData.length > 0) {
        const newData = [...tableData, ...sliceData]
        this.setData({
          tableData: newData
        }, () => {
          setTimeout(() => {
            this.setData({
              showLoading: false
            })
          }, 500)
        })
      } else {
        setTimeout(() => {
          this.setData({
            showLoading: false,
            hasMore: false
          })
        }, 500)
      }
    }).catch((err) => {
      console.log(err)
    })
  },
  loadMoreData(e) {
    const page = e.detail
    this.loadMore(page)
  },
  /**
   * @description:查询赋值
   */
  changeInput (e) {
    const value = e.detail.value
    const key = e.currentTarget.dataset.key
    this.setData({
      [key]: value
    })
  },
  /**
   * @description: 点击查看有库存
   */
  checkStock() {
    // 将tableData的数据过滤
    this.setData({
      checkActive: !this.data.checkActive
    })
    // 如果是筛选条件需要重新筛选
    this.loadMore(1)
  },
  /**
   * @description: 点击刷新查询
   */
  refreshStock() {
    this.getInventoryList()
  },
  /**
   * @description: 点击重新查询
   */
  handleSubmit() {
    // 清空表头的排序,默认净库存为降序，其他清空
    const tableHeader = this.data.tableHeader.map(item => {
      if (item.prop === 'on_hand') {
        item.sort_type = 'asc'
      } else {
        item.sort_type = ''
      }
      return item
    })
    this.setData({
      tableHeader
    })
    if (this.data.from === '3') {
      wx.exitMiniProgram()
    } else {
      wx.navigateBack({
        delta: 1
      })
    }
  },
  /**
   * @description: 图片加载失败替换图片
   */
  imgError(e) {
    const index = e.detail
    const image_url = `https://cdn-cos.coach.com.cn/uploads/2023/v3.2.0.1/default.png?time=${new Date().getTime()}`
    this.setData({
      [`tableData[${index}].image_url`]: image_url
    })
  },
  imgLoaded(e) {
    const {index, url, source_num, info} = e.detail
    const image_url = `https://cdn-cos.coach.com.cn/uploads/2023/v3.2.0.1/default.png?time=${new Date().getTime()}`
    if (info) {
      console.log(111, info)
      this.setData({
        [`tableData[${index}]._loaded`]: true,
        [`tableData[${index}].image_url`]: url || this.data.tableData[index].image_url,
        [`tableData[${index}].source_num`]: source_num,
        [`tableData[${index}].sku_tag_price`]: info.sku_tag_price ? `¥${info.sku_tag_price}` : '--',
        [`tableData[${index}].store_price`]: info.store_price ? `¥${info.store_price}` : '--'
      })
    } else {
      this.setData({
        [`tableData[${index}]._loaded`]: true,
        [`tableData[${index}].image_url`]: image_url,
        [`tableData[${index}].sku_tag_price`]: '--',
        [`tableData[${index}].store_price`]: '--'
      })
    }
  },
  /**
   * @description: 点击排序
   * @param {Object} e
   */
  sortChange(e) {
    const {index, item} = e.detail.currentTarget.dataset
    // todo: 将tableData的数据排序
    const sortType = item.sort_type === 'desc' ? 'asc' : 'desc'
    const copyTableData = this.data.copyTableData.sort((a, b) => {
      return sortType === 'desc' ? a[item.prop] - b[item.prop] : b[item.prop] - a[item.prop]
    })
    // todo: 将除去当前点击之外的sort_type改为''
    this.data.tableHeader.forEach((item, otherindex) => {
      if (item.sort_type && index !== otherindex) {
        this.setData({
          [`tableHeader[${otherindex}].sort_type`]: ''
        })
      }
    })
    this.setData({
      // todo: 将tableHeader的sortType改变
      [`tableHeader[${index}].sort_type`]: item.sort_type === 'desc' ? 'asc' : 'desc',
      copyTableData
    })
    this.loadMore(1)
  },
  /**
   * @description: 显示颜色筛选弹窗
   */
  showColorSelect() {
    // 判断列表是否有数据
    const { tableData, isSelect} = this.data
    if (tableData.length > 0 || isSelect) {
      this.setData({
        showSelect: true,
        type: 'color'
      })
    }
  },
  /**
   * @description: 显示尺寸筛选弹窗
   */
  showSizeSelect() {
    const { tableData, isSelect} = this.data
    if (tableData.length > 0 || isSelect) {
      this.setData({
        showSelect: true,
        type: 'size'
      })
    }
  },
  /**
   * @description: 点击颜色，尺寸筛选弹窗选择
   * @param {Object} e
   * @return {*}
   */
  confirmSelect(e) {
    const {index, type} = e.detail
    if (type === 'color') {
      this.setData({
        colorIndex: index,
        colorValue: this.data.colorList[index],
        showSelect: false,
        isSelect: true
      })
    } else {
      this.setData({
        sizeIndex: index,
        sizeValue: this.data.sizeList[index],
        showSelect: false,
        isSelect: true
      })
    }
    // 如果是筛选条件需要重新筛选
    this.loadMore(1)
  },
  /**
   * @description: 点击颜色，尺寸筛选弹窗取消
   */
  cancelSelect() {
    this.setData({
      showSelect: false
    })
  },
  /**
   * @description: 点击跳转到其他门店库存页面
   * @param {Object} e
   */
  showSearchMore(e) {
    let {item} = e.detail.currentTarget.dataset
    item.active = this.data.active
    const params = encodeURIComponent(JSON.stringify(item))
    wx.navigateTo({
      url: `/pagesStock/pages/other-list/other-list?from=${this.data.from}&searchParams=${params}&isHorizontal=${this.data.isHorizontal}`
    })
  },
  /**
   * @description: 关闭提示弹窗
   */
  tipsCancel() {
    this.setData({
      showListTips: false
    })
  },
  /**
   * @description: 显示图片预览
   */
  showImageInfo(e) {
    const {item} = e.detail.currentTarget.dataset
    const { store_code } = this.data.store
    // console.log(item.image_url)
    // this.setData({
    //   showPreview: true,
    //   previewImg: item.image_url
    // })
    wx.navigateTo({
      url: `/pagesMaterial/pages/select/select?from=2&store_code=${store_code}&upc_code=${item.upc_code}`
    })
    // 神策数据上报 -- 进入素材库
    wx.track.trackAssetLibrary({
      entry_point: 4,
      library_type: 4
    })
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
    this.getInventoryList()
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
   * @de.getStoreList
   * @return {*}
   */
  async getStoreList() {
    const res = await getStoreListApi()
    const storeList = res.data.list
    const {storeInfo} = app.globalData

    if (storeInfo.store_code || (this.data.store && this.data.store.store_code)) {
      this.setData({
        store: storeInfo
      })
    } else {
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
  },
  /**
   * @description: 点击横屏
   */
  exitScreen() {
    this.setData({
      isHorizontal: !this.data.isHorizontal
    })

    wx.nextTick(() => {
      getRect(this, '#tableDataArea2').then(res => {
        console.log(2222, res)
        if (res) {
          this.setData({
            minWidth: Math.max(res.width, res.height)
          })
        }
      }).catch(err => {
        console.log(err)
      })
    })
  }
})
