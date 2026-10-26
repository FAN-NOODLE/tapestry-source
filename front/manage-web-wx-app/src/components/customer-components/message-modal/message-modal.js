// Component Object
Component({
  externalClasses: ['custom-class'],
  properties: {
    show: {
      type: Boolean
    }
  },
  data: {
    show: false
  },
  methods: {
    handleCancel() {
      this.triggerEvent('cancel')
    }
  }
})

