// Component Object
Component({
  properties: {
    show: {
      type: Boolean,
      value: false
    },
    current: {
      type: Number,
      value: 0
    },
    list: {
      type: Array,
      value: []
    }

  },
  data: {

  },
  methods: {
    /**
     * 切换
     */
    handleChange(e) {
      const { index } = e.currentTarget.dataset
      this.setData({
        current: index
      })
    },

    /**
     * 关闭
     */
    handleClose() {
      this.triggerEvent('close')
    },
    /**
     * @description: 确定
     * @return {*}
     */
    handleConfirm() {
      this.triggerEvent('confirm', {
        value: this.data.current
      })
    }
  }
})
