const { getUploadAuthApi } = require('@/api/material-api')


// Page Object
Page({
  data: {
    current: 0,
    tabs: [
      {
        title: 'KOS素材',
        value: 1
      },
      {
        title: '官网商品素材',
        value: 2
      }
    ]
  },
  // options(Object)
  onLoad(options) {
    wx.hideHomeButton()
    this.getUploadAuth()
    const { from, current, sku_item } = options
    if (from === 'TALK') {
      this.setData({
        current: Number(current),
        entryPoint: 5
      })
      this.selectComponent('#com-goods-material').search(sku_item)
    }

    // 神策数据上报 -- 进入素材库
    wx.track.trackAssetLibrary({
      entry_point: from === 'TALK' ? 5 : 1,
      library_type: 1
    })
  },
  onReachBottom() {
    const { current } = this.data
    const selector = ['#com-normal-material', '#com-goods-material'][current]
    this.selectComponent(selector).loadMore()
  },
  onPageScroll({ scrollTop }) {
    const { current } = this.data
    const selector = ['#com-normal-material', '#com-goods-material'][current]
    if (this.selectComponent(selector).scroll) {
      this.selectComponent(selector).scroll()
    }
    if (this.selectComponent('#com-create-entry')) {
      this.selectComponent('#com-create-entry').scroll()
    }
    this.selectComponent('#com-create-entry-2').scroll()
  },

  /**
   * @description 获取上传素材权限
   */
  async getUploadAuth() {
    try {
      const res = await getUploadAuthApi()
      if (res.code === 200) {
        this.setData({
          canUpload: res.data.auth === 1
        })
      }
    } catch (error) {
      console.log(error)
    }
  },
  handleShow() {
    if (this.selectComponent('#com-create-entry')) {
      this.selectComponent('#com-create-entry').scroll()
    }
  },
  /**
   * 类型切换
   * @param {*} e
   */
  handleTypeChange(e) {
    const { index } = e.detail
    this.setData({
      current: index
    })
  },
  handleViewCollect() {
    wx.navigateTo({
      url: '/pagesMaterial/pages/collect/collect?entryPoint=' + this.data.entryPoint,
      events: {
        update: ({item}) => {
          const { current } = this.data
          const selector = ['#com-normal-material', '#com-goods-material'][current]
          if (this.selectComponent(selector).setUpdateList) {
            this.selectComponent(selector).setUpdateList()
          }
          this.selectComponent(selector).updateCurrent(item)
        }
      }
    })
    // 神策数据上报 -- 点击收藏入口
    wx.track.trackMyFavoriteEntranceClick({
      operation_type: '我的收藏',
      entry_point: '工作台',
      library_type: '品牌素材库'
    })
  },
  handleCreate() {
    this.selectComponent('#com-normal-material').handleCreate()
  }
})

