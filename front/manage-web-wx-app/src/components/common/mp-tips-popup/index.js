const { toggleTab } = require('@/packages/customer/utils')

// Component Object
Component({
  properties: {
    title: {
      type: String,
      value: '提示标题'
    },
    content: {
      type: [String, Array],
      value: '提示内容',
      observer(newVal) {
        console.log(newVal)
        if (typeof newVal === 'string') {
          this.setData({
            textList: [newVal]
          })
        } else {
          this.setData({
            textList: newVal
          })
        }
      }
    },
    useSlot: {
      type: Boolean,
      value: false
    }
  },
  data: {
    show: false,
    textList: []
  },
  methods: {
    open() {
      this.setData({
        show: true
      })
      toggleTab()
    },
    close() {
      this.setData({
        show: false
      })
      toggleTab()
      this.triggerEvent('cancel')
    }
  }
})
