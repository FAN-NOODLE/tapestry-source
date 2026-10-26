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
    imagePath,
    // 等级足迹
    levelList: [],
    // 最高等级信息
    maxLevelInfo: {
      level: 1,
      time: 0,
      name: '新手上路'
    },
    showTips: true
  },
  methods: {
    // 关闭弹窗
    closePopup () {
      this.setData({
        show: false
      })
    },
    init () {
      const { list } = this.data
      // 获取最高等级信息
      const maxLevelInfo = list.find(item => item.is_max === 1)

      const length = list.length
      // 设置位置
      if (length >= 5) {
        const position = [
          ['left', 40, 90],
          ['right', 74, 218],
          ['left', 80, -28],
          ['right', 104, 202]
        ]
        list.forEach((item, index) => {
          const i = index % 4
          const offsetY = index > 0 && i === 0 ? 120 : 0
          console.log(offsetY)
          const [direction, top, left] = position[i]
          item.marginTop = top + offsetY
          item.marginLeft = left
          item.direction = direction
          item.styles = `margin-top: ${item.marginTop}rpx; margin-left: ${item.marginLeft}rpx;`
        })
      }
      this.setData({
        maxLevelInfo,
        levelList: list,
        showTips: list.length > 5
      })
    },
    onScroll () {
      this.setData({
        showTips: false
      })
    }
  },
  created: function () {

  },
  attached: function () {
    this.init()
  },
  ready: function () {

  },
  moved: function () {

  },
  detached: function () {

  }
})

