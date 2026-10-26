// 获取全局应用程序实例对象
// const app = getApp();
import { getInventoryListApi } from '@/api/stock-api'
import { getStoreListApi } from '@/api/store-api'
import { otherHeader } from '../../utils/data'
import { groupArr } from '../../utils/index'
import { getRect } from '@/pagesMaterial/utils/util'

const { getBrandMaterialDetailApi } = require('@/api/material-api')

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    active: 0,
    skuInfo: {},
    cityNum: 0,
    storeNum: 0,
    tableHeader: otherHeader,
    oldTableData: [],
    tableData: [],
    storeList: [],
    storeCityList: [],
    routeParams: {},
    showLoading: false,
    checkActive: true,
    showStoreSelect: false,
    checkedList: [],
    phone: '',
    showPhonePopup: false,
    isFirstOpen: true,
    isHorizontal: false
  },

  // 生命周期函数--监听页面加载
  async onLoad(options) {
    // TODO: onLoad
    this.setData({
      options: options
    })
    if (options && options.searchParams) {
      const params = decodeURIComponent(options.searchParams)
      const routeParams = JSON.parse(params)
      this.setData({
        from: options.from,
        skuInfo: routeParams,
        active: Number(routeParams.active),
        isHorizontal: options.isHorizontal === 'true'
      })
      this.getBrandMaterialDetail()
      await this.getInventoryList()

      if (this.data.options.isHorizontal === 'true') {
        this.exitScreen()
      }
    }
  },
  getBrandMaterialDetail() {
    getBrandMaterialDetailApi({
      upc_code: this.data.skuInfo.upc_code
    }).then(res => {
      if (res.data.url) {
        this.setData({
          'skuInfo.image_url': res.data.url,
          'skuInfo.sku_item': res.data.sku_item,
          'skuInfo.color': res.data.color,
          'skuInfo.style': res.data.style || ''
        })
      }
      // 神策数据上报 -- 跨店库存查询
      wx.track.trackOtherStockInquiry({
        product_id: res.data.resource_key,
        product_name: res.data.product_name
      })
    }).catch(err => {
      console.log(err)
    })
  },
  /**
   * @description:获取库存列表
   */
  async getInventoryList() {
    this.setData({
      showLoading: true
    })
    const params = {
      'upc_code': this.data.skuInfo.upc_code
    }
    try {
      const res = await getInventoryListApi(params)
      // 判断data是否是数组
      let data = res.data.list
      if (data instanceof Array) {
        this.getStoreList(data)
      } else {
        this.getStoreList([])
      }
    } catch (error) {
      wx.showToast('获取库存列表失败')
      this.setData({
        showLoading: false
      })
    }
  },
  /**
   * @des.getStoreListApi
   * @return {*}
   */
  getStoreList(data) {
    getStoreListApi().then(res => {
      const storeList = res.data.list
      // todo: 将所有的门店数据添加到data中,并且将门店数据添加到data中
      storeList.forEach(item => {
        const index = data.findIndex(item1 => item1.loc_code === item.third_store_code)
        if (index === -1) {
          data.push({
            ...item,
            loc_code: item.third_store_code,
            on_hand: 0,
            in_transit: 0
          })
        } else {
          data[index] = {
            ...data[index],
            ...item
          }
        }
      })
      // 创建树形结构,将门店按照城市分组
      const treeData = groupArr(data, 'city')
      this.setData({
        storeList: data,
        storeCityList: treeData,
        oldTableData: data
      })
      this.dealStockList(data).then(res => {
        this.setData({
          tableData: res,
          showLoading: false
        })
      }).catch(err => {
        console.log(err)
      })
    }).catch(err => {
      console.log(err)
      this.setData({
        showLoading: false
      })
    })
  },
  /**
   * @description: 筛选逻辑
   */
  dealStockList(data) {
    return new Promise((resolve, reject) => {
      const { checkedList, checkActive, isFirstOpen } = this.data
      let tableData = data || this.data.tableData
      if (checkedList.length > 0) {
        checkedList.forEach(item => {
          const index = tableData.findIndex(item1 => item1.loc_code === item.third_store_code)
          item.loc_code = item.third_store_code
          if (index > -1) {
            item.on_hand = tableData[index].on_hand
            item.in_transit = tableData[index].in_transit
          } else {
            item.on_hand = 0
            item.in_transit = 0
          }
        })
        tableData = checkedList
      }
      // 判断是否点击了排序
      if (this.data.tableHeader.some(item => item.sort_type)) {
        const sortItem = this.data.tableHeader.find(item => item.sort_type)
        tableData = tableData.sort((a, b) => {
          return sortItem.sort_type === 'desc' ? a[sortItem.prop] - b[sortItem.prop] : b[sortItem.prop] - a[sortItem.prop]
        })
      }
      if (checkActive) {
        tableData = tableData.filter(item => item.on_hand + item.in_transit > 0)
        if (this.data.checkedList.length === 0) {
          this.setData({
            checkedList: tableData
          })
        }
      }
      // 获取有库存的城市数量和门店数量
      if (isFirstOpen) {
        const cityNum = new Set(tableData.map(item => item.city)).size
        const storeNum = tableData.length
        this.setData({
          cityNum,
          storeNum
        })
      }
      resolve(tableData)
    })
  },
  /**
   * @description: 点击查看有库存
   */
  checkStock() {
    // 将tableData的数据过滤
    this.setData({
      checkActive: !this.data.checkActive,
      showLoading: true
    })
    this.dealStockList(this.data.oldTableData).then(res => {
      this.setData({
        tableData: res
      })
      setTimeout(() => {
        this.setData({
          showLoading: false
        })
      }, 500)
    }).catch(err => {
      console.log(err)
    })
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
        delta: 2
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
    const tableData = this.data.tableData.sort((a, b) => {
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
      tableData
    })
  },
  /**
   * @description: 显示门店选择
   */
  showStore() {
    this.setData({
      showStoreSelect: true
    })
    if (!this.data.isFirstOpen) {
      this.selectComponent('#storeSelect').checkStock()
    }
  },
  /**
   * @description: 点击门店选择
   */
  confirmStoreSelect(e) {
    const { checkedList, checkedNum } = e.detail
    //  获取选中城市数量
    const cityNum = new Set(checkedList.map(item => item.city)).size
    // 更新storeCityList
    const storeCityList = this.data.storeCityList.map(item => {
      const children = item.children.map(item1 => {
        const index = checkedList.findIndex(item2 => item2.third_store_code === item1.third_store_code)
        if (index > -1) {
          item1.checked = true
        } else {
          item1.checked = false
        }
        return item1
      })
      return {
        ...item,
        checkedNum: children.filter(item => item.checked).length,
        allChecked: children.filter(item => item.checked).length === children.length,
        children
      }
    })
    this.setData({
      cityNum,
      storeNum: checkedNum,
      storeCityList,
      showStoreSelect: false,
      checkedList,
      showLoading: true,
      isFirstOpen: false
    })
    this.dealStockList(this.data.oldTableData).then(res => {
      this.setData({
        tableData: res
      })
      setTimeout(() => {
        this.setData({
          showLoading: false
        })
      }, 500)
    }).catch(err => {
      console.log(err)
    })
  },
  /**
   * @description: 关闭门店选择
   */
  closeStoreSelect() {
    this.setData({
      showStoreSelect: false
    })
  },
  /**
   * @description: 点击电话图标
   * @return {*}
   */
  callPhone(e) {
    const phone = e.detail
    this.setData({
      showPhonePopup: true,
      phone
    })
  },
  /**
   * @description: 拨打电话
   * @param {*} e
   */
  handleCall(e) {
    const phone = e.detail
    if (phone) {
      wx.makePhoneCall({
        phoneNumber: phone
      })
    }
  },
  /**
   * @description: 复制电话
   * @param {*} e
   */
  handleCopy(e) {
    const phone = e.detail
    wx.setClipboardData({
      data: phone
    })
  },
  /**
     * @description: 点击横屏
     */
  exitScreen() {
    this.setData({
      isHorizontal: true
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
  },
  /**
     * @description: 点击竖屏
     */
  enterScreen() {
    this.setData({
      isHorizontal: false
    })
  }
})
