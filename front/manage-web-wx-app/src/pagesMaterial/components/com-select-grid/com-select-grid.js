// Component Object
Component({
  properties: {
    current: {
      type: Number,
      value: 0
    },
    showFirstTips: {
      type: Boolean,
      value: false
    }
  },
  data: {},
  methods: {
    /**
     * @description 隐藏首次提示
     */
    handleHideTips() {
      this.setData({
        showFirstTips: false
      })
      this.triggerEvent('hide')
    },
    /**
     * @description 切换宫格
     */
    handleGideChange(e) {
      const { value } = e.currentTarget.dataset
      this.triggerEvent('change', {
        value
      })
      this.handleHideTips()
    }
  }
})
