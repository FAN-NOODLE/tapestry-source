Component({
  options: {
    addGlobalClass: true
  },
  data: {
    name: 'loading',
    imgIndex: 0
  },

  /*    系统事件开始   */
  detached () {
    this.timer && clearInterval(this.timer)
  },
  ready () {
    this.init()
  },
  /*    系统事件结束   */

  methods: {
    /*    普通事件开始  */

    /*    普通事件结束  */
    init () {
      let self = this
      this.data.isTimer = true
      this.timer = setInterval(() => {
        let imgIndex = ++self.data.imgIndex
        if (imgIndex === 3) {
          imgIndex = 0
        }
        self.setData({
          imgIndex
        })
      }, 1000)
    }
  }
})
