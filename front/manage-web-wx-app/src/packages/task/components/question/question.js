// Component Object
const imagePath = 'https://cdn-cos.coach.com.cn/uploads/'

Component({
  properties: {
    show: {
      type: Boolean,
      value: false
    },
    list: {
      type: Array,
      value: [],
      observer: 'init'
    }
  },
  data: {
    imagePath
  },
  methods: {
    // 关闭弹窗
    closePopup () {
      this.setData({
        show: false
      })
      this.triggerEvent('close')
    }
  },
  created: function () {

  },
  attached: function () {
  },
  ready: function () {

  },
  moved: function () {

  },
  detached: function () {

  }
})

