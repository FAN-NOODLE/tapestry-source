// Component Object
// Component Object
Component({
  properties: {
    show: {
      type: Boolean,
      value: false,
      observer(newVal) {
        if (newVal) {
          this.setData({
            innerCurrent: this.data.current || 0
          })
        }
      }
    },
    keyName: {
      type: String,
      value: 'label'
    },
    current: {
      type: Number,
      value: 0
    },
    options: {
      type: Array,
      value: []
    }
  },
  data: {
    innerCurrent: 0
  },
  methods: {
    handleSelect(e) {
      const { index } = e.currentTarget.dataset
      this.setData({
        innerCurrent: index
      })
      this.triggerEvent('change', {
        index
      })
    },
    handleConfirm() {
      const { innerCurrent, options } = this.data
      this.triggerEvent('confirm', {
        index: innerCurrent,
        value: options[innerCurrent]
      })
    },
    handleCancel() {
      this.triggerEvent('cancel')
    }
  }
})
