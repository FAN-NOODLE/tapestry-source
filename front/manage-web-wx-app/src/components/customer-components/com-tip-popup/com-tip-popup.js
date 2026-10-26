// Component Object
Component({
  properties: {
    title: {
      type: String
    },
    list: {
      type: Array
    },
    show: {
      type: Boolean,
      value: false,
      observer() {
        console.log(this.data.list)
      }
    }
  },
  data: {
    showTip: false
  },
  methods: {
    /**
     * @description 关闭弹框
     */
    onClose(e) {
      this.triggerEvent('close', e.detail)
    }
  }
})

