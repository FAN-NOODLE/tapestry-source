import { debounce } from '@/utils/tools'

// Component Object
Component({
  properties: {
    isCancel: Boolean
  },
  data: {
    show: true
  },
  methods: {
    hide: debounce(function() {
      this.setData({
        show: true
      })
    }),
    scroll() {
      if (this.data.show) {
        this.setData({
          show: false
        })
      }
      this.hide()
    },
    onClick() {
      this.triggerEvent('click')
    },
    onCancel() {
      // 神策数据上报 -- 创建群发页操作，取消群发
      wx.track.trackBatchMsgOperation({
        operation_type: 3
      })
      wx.navigateBack({
        delta: 1
      })
    }
  }
})

