const imagePath = 'https://cdn-cos.coach.com.cn/uploads/'

Component({
  properties: {
    info: {
      type: Object,
      value: {}
    }
  },
  data: {
    imagePath: imagePath
  },
  methods: {
    /* 点击奖章 */
    tapMedal () {
      this.triggerEvent('openPopup')
    }
  }
})
