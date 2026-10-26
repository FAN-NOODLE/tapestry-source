

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
  onLoad: function(e) {
    const eventChannel = this.getOpenerEventChannel()
    const that = this
    this.setData({
      options: e
    })
    if (eventChannel.hasOwnProperty('once')) {
      eventChannel.once('getExternalId', data => {
        that.setData({
          'externalIds': data
        })
      })
    }
    // 神策数据上报 -- 进入素材库
    wx.track.trackAssetLibrary({
      entry_point: 3,
      library_type: 3
    })
  },
  onReachBottom() {
    const { current } = this.data
    const selector = ['#com-send-material', '#com-goods-material'][current]
    this.selectComponent(selector).loadMore()
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
  }
})

