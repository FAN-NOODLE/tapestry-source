// Component Object
Component({
  properties: {
    value: {
      type: String,
      value: '',
      observer: function() {}
    },
    maxlength: {
      type: Number,
      value: 200
    },
    disabled: {
      type: Boolean,
      value: false
    }
  },
  data: {

  },
  methods: {
    handleInput(e) {
      const value = e.detail.value.slice(0, this.data.maxlength)

      this.setData({
        value
      })
      this.triggerEvent('change', {value})
      return {value}
    },
    tapDisableEdit () {
      if (this.properties.disabled) {
        wx.showToast('该用户非您好友，暂无法设置播种')
      }
    }
  }
})
