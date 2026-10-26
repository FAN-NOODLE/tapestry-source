const { toggleTab } = require('../../utils')

// Component Object
Component({
  properties: {
    show: {
      type: Boolean,
      value: false
    },
    value: {
      type: Object,
      value: {}
    }

  },
  data: {
    showConfirmPopup: false
  },
  methods: {
    handleClose() {
      this.setData({
        show: false,
        showConfirmPopup: false
      })
      this.triggerEvent('close')
      toggleTab()
    },
    handleCloseConfirm() {
      this.setData({
        showConfirmPopup: false
      })
      this.triggerEvent('close-confirm')
    },
    handleShow() {
      this.setData({
        show: true
      })
      toggleTab()
    },
    handleShowConfirm() {
      this.setData({
        showConfirmPopup: true
      })
    },
    handleChange(e) {
      const { key } = e.currentTarget.dataset
      this.setData({
        [key]: e.detail.value
      })
    },
    handleConfirm() {
      this.triggerEvent('confirm', this.data.value)
    }
  }
})
