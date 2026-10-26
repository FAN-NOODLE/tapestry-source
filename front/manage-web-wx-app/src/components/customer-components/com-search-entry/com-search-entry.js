import { debounce } from '@/utils/tools'

// Component Object
Component({
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
      wx.navigateBack({
        delta: 1
      })
    }
  }
})

