// Component Object
Component({
  properties: {
    current: {
      type: Number,
      value: 0
    },
    list: {
      type: Array,
      value: [
        {
          id: 0,
          name: '自定义播种'
        }
      ]
    }
  },
  data: {
    showPopup: false
  },
  methods: {
    handleClick() {
      this.setData({
        showPopup: !this.data.showPopup
      })
    },
    handleSelect(e) {
      const { index } = e.currentTarget.dataset
      this.setData({
        current: index,
        showPopup: false
      })
      this.triggerEvent('change', {
        value: index
      })
    }
  }
})
