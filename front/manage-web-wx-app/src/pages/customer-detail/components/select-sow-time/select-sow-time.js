
const dayjs = require('dayjs')

// Component Object
Component({
  properties: {
    value: {
      type: Number,
      value: 0,
      observer(newVal, oldVal) {
        if (newVal !== oldVal) {
          this.setData({
            dateTimeText: newVal ? dayjs(newVal * 1000).format('YYYY-MM-DD HH:mm') : ''
          })
        }
      }
    }

  },
  data: {
    showPopup: false,
    minDate: Date.now()
  },
  lifetimes: {
    ready() {
      this.setData({
        minDate: Date.now() + 5 * 60 * 1000
      })
    }
  },
  methods: {
    handleClick() {
      this.setData({
        showPopup: !this.data.showPopup
      })
    },
    handleCancel() {
      this.setData({
        showPopup: false
      })
    },
    handleConfirm(e) {
      const value = parseInt(e.detail.value / 1000)
      console.log(value)
      this.setData({ value: value, dateTimeText: dayjs(value * 1000).format('YYYY-MM-DD HH:mm') })
      this.triggerEvent('change', { value: value })
      this.setData({
        showPopup: false
      })
    }
  }
})
