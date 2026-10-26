const { getBrandMaterialListApi } = require('@/api/material-api')
const { getRect } = require('@/components/common/common/utils')

/* eslint-disable standard/object-curly-even-spacing */
const app = getApp()
Component({
  /** 组件的属性列表 */
  properties: {
    data: {
      type: Array,
      value: []
    },
    headers: {
      type: Array,
      value: []
    },
    /**  table的高度, 溢出可滚动 */
    height: {
      type: String,
      value: 'auto'
    },
    width: {
      type: Number || String,
      value: '100%'
    },
    minWidth: {
      type: Number
    },
    hasMore: {
      type: Boolean,
      value: true
    },
    storeCode: {
      type: String,
      value: ''
    },
    maskWidth: {
      type: String,
      value: '392rpx'
    },
    isHorizontal: {
      type: Boolean,
      value: false
    },
    emptyText: {
      type: String,
      value: '未找到相关商品信息'
    }
  },

  /** 组件的初始数据 */
  data: {
    scrollWidth: '100%',
    isDrag: true,
    showFixed: false,
    page: 1,
    hideMask: wx.getStorageSync('hideMask')
  },
  /** 组件的监听属性 */
  observers: {
    // 在 numberA 或者 numberB 被设置时，执行这个函数
    'headers': function (headers) {
      this.init()
    },
    'minWidth': function (minWidth) {
      this.init()
    }
  },
  /** 组件的方法列表 */
  methods: {
    init() {
      const reducer = (accumulator, currentValue) => {
        return accumulator + Number(currentValue.width)
      }
      const scrollWidth = this.data.headers.reduce(reducer, 0)
      const { screenWidth } = app.globalData.system
      const pixelRatio = 750 / screenWidth
      const showFixed = (scrollWidth / pixelRatio) > (this.data.minWidth || screenWidth)
      this.setData({
        scrollWidth: this.data.minWidth ? this.data.minWidth * pixelRatio : scrollWidth - 10,
        showFixed,
        columnWidth: !showFixed ? ((this.data.minWidth * pixelRatio - scrollWidth) / this.data.headers.length) : 0
      })
      console.log(scrollWidth / pixelRatio)

      wx.nextTick(() => {
        getRect(this, '.table').then(res => {
          console.log(2222, res, scrollWidth)
          if (res) {
            this.setData({
              tableHeight: res.height - 4
            })
          }
        }).catch(err => {
          console.log(err)
        })
      })
    },
    bindscrolltolower(e) {
      const {direction} = e.detail
      if (direction === 'bottom') {
        let { page, hasMore} = this.data
        if (!hasMore) return
        page++
        this.setData({
          page: page
        })
        console.log(`==========第${page}页`)
        this.triggerEvent('loadMoreData', page)
      }
    },
    showImageInfo(e) {
      this.triggerEvent('showImageInfo', e, e.currentTarget.dataset.item)
    },
    showSearchMore(e) {
      this.triggerEvent('showSearchMore', e, e.currentTarget.dataset.item)
    },
    copyCode(e) {
      const { item } = e.currentTarget.dataset
      wx.setClipboardData({
        data: item.upc_code,
        success () {
          wx.showToast('已复制')
        }
      })
    },
    imgError(e) {
      this.triggerEvent('imgError', e.currentTarget.dataset.index)
    },
    sortChange(e) {
      console.log(1111, e.currentTarget.dataset.item.sort)
      if (e.currentTarget.dataset.item.sort) {
        this.triggerEvent('sortChange', e)
      }
    },
    callPhone(e) {
      const { phone } = e.currentTarget.dataset.item
      this.triggerEvent('callPhone', phone)
    },
    showToast(e) {
      const { header, index, headerItem, content } = e.currentTarget.dataset
      if (headerItem.showTooltip && content && !this.data.isHorizontal) {
        const selector = `#${header}_${index}`
        console.log(selector)
        this.selectComponent('#mp-tooltip').show(this, selector)
      }
    },
    /** 加载图片 */
    handleLoadImage(e) {
      const { callback } = e.detail
      const { item, index } = e.currentTarget.dataset
      if (item._loaded) {
        callback()
        return
      }
      console.log(item)
      getBrandMaterialListApi({
        upc_code: item.upc_code,
        store_code: this.data.storeCode
      }).then(res => {
        if (res.data && res.data.length) {
          const list = res.data[0].brand_source_list
          const currentImage = list.find(image => image.is_main === 1) || list[0]
          this.triggerEvent('loaded', {
            index,
            url: currentImage ? currentImage.url : '',
            source_num: list.length,
            info: res.data[0]
          })
          setTimeout(() => {
            callback()
          }, 10)
        } else {
          this.triggerEvent('loaded', {
            index,
            info: null
          })
          callback()
        }
      }).catch(err => {
        console.log(err)
        callback()
      })
    },
    scroll() {
      if (!this.data.hideMask) {
        this.setData({
          hideMask: true
        })
        wx.setStorageSync('hideMask', true)
      }
    }
  }
})
