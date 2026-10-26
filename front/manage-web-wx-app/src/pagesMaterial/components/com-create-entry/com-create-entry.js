import debounce from '@/utils/debounce'

// Component Object
Component({
  properties: {
    type: {
      type: Number,
      value: 1
    },
    bottom: {
      type: Number,
      value: 100
    }
  },
  data: {
    show: true
  },
  methods: {
    hide: debounce(function() {
      console.log(11111)

      this.setData({
        show: true
      })
      this.triggerEvent('show')
    }, 100, {
      immediate: false
    }),
    scroll() {
      if (this.data.show) {
        this.setData({
          show: false
        })
      }
      this.hide()
    },
    show() {
      this.setData({
        show: true
      })
    },
    onClick() {
      this.triggerEvent('click', this.data.type)
    }
  }
})

