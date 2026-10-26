Component({
  data: {
    /** 图片地址 */
    url: '',
    /** 判断是否显示 */
    show: false,
    /** 弹窗类型 */
    type: 'success',
    /** 提示文案 */
    message: '',
    /** 自动关闭时长 */
    timing: 1000
  },
  methods: {
    openPopup (type, message, timing = 1000) {
      const url = `https://cdn-cos.coach.com.cn/uploads/2022/v2.8.0/${type}.png`
      this.setData({
        url,
        type,
        message,
        timing,
        show: true
      })
      this.autoClose()
    },
    /** 自动关闭 */
    autoClose () {
      setTimeout(() => {
        this.setData({
          show: false
        })
      }, this.data.timing)
    }
  }
})

