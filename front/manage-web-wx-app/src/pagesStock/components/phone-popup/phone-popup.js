// Component Object
Component({
  externalClasses: ['custom-class'],
  properties: {
    show: {
      type: Boolean
    },
    phone: {
      type: String
    },
    isHorizontal: {
      type: Boolean
    }
  },
  data: {
    show: false
  },
  methods: {
    handleCancel() {
      this.setData({
        show: false
      })
    },
    handleCall() {
      this.triggerEvent('call', this.properties.phone)
    },
    handleCopy() {
      this.triggerEvent('copy', this.properties.phone)
    }
  }
})

