// Component Object
Component({
  externalClasses: ['custom-class'],
  properties: {
    show: {
      type: Boolean
    },
    list: {
      type: Array
    },
    type: {
      type: String
    },
    tableArea: {
      type: Object
    },
    chooseIndex: {
      type: Number,
      observer(newVal) {
        this.setData({
          chooseIndex: newVal
        })
      }
    }
  },
  data: {
    show: false,
    chooseIndex: 0
  },
  methods: {
    handleSelectColor(e) {
      console.log(e)
      const index = e.currentTarget.dataset.index
      this.setData({
        chooseIndex: index
      })
    },
    handleConfirm() {
      this.triggerEvent('confirm', {
        index: this.data.chooseIndex,
        type: this.properties.type
      })
    },
    handleCancel() {
      this.triggerEvent('cancel')
    }
  }
})

