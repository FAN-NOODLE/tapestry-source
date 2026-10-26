// Component Object
Component({
  properties: {
    list: {
      type: Array,
      value: [
        {
          label: '全部',
          key: 0,
          type: 1
        },
        {
          label: '专属顾客',
          key: 1,
          type: 2
        },
        {
          label: '非专属顾客',
          key: 2,
          type: 3
        }
      ]
    },
    current: {
      type: Number,
      value: 0
    },
    showIcon: {
      type: Boolean,
      value: true
    }
  },
  data: {
    imageHost: 'https://cdn-cos.coach.com.cn/uploads/',
    show: false
  },
  methods: {
    /**
     * @description 点击切换
     * @param {*} e
     */
    onClick(e) {
      const { key } = e.currentTarget.dataset
      this.setData({
        current: key
      })
      this.triggerEvent('change', { value: key })
    },
    /**
     * @description 显示弹框
     */
    handleShowPopup() {
      wx.hideTabBar()
      this.setData({ show: true })
    },
    /**
     * @description 隐藏弹框
     */
    handleClosePopup() {
      this.setData({ show: false })
      setTimeout(() => {
        wx.showTabBar()
      }, 300)
    }
  }
})

