/**
 * components/common/mp-movable/index
 * @name 组件名称
 * @author author <emailAddress>
 * @description 组件描述
 * @property {string} style '' 定义 style 属性可以拿到 style 属性上设置的值
 * @property {boolen} virtualHost false 虚拟化组件节点
 * @property {string} styleIsolation apply-shared 表示页面 wxss 样式将影响到自定义组件，但自定义组件 wxss 中指定的样式不会影响页面；
 * @property {boolen} multipleSlots true 支持多插槽
 */

// import { systemInfo } from '@utils/getSystemInfo';
const { globalData } = getApp()
Component({
  externalClasses: ['custom-class'],
  options: {
    virtualHost: false,
    multipleSlots: true,
    styleIsolation: 'apply-shared'
  },
  properties: {
    direction: {
      type: String,
      value: 'all'
    },
    position: {
      type: Object,
      value: {
        left: 0,
        right: 0,
        bottom: 0,
        top: 0
      }
    },
    left: {
      type: Number,
      value: 0
    },
    right: {
      type: Number,
      value: 0
    },
    top: {
      type: Number,
      value: 0
    },
    bottom: {
      type: Number,
      value: 0
    },
    hMargin: {
      // 拖动后的水平方向最小边距（左右边距）
      type: Number,
      value: 10
    },
    vMargin: {
      // 拖动后的垂直方向最小边距(上下边距)
      type: Number,
      value: 10
    }
  },
  data: {
    x: 999,
    y: 999,
    windowWidth: globalData.system.windowWidth,
    windowHeight: globalData.system.windowHeight,
    elementWidth: 0,
    elementHeight: 0,
    animation: false
  },
  ready() {
    this.init()
  },
  methods: {
    /**
     * 初始化
     */
    init() {
      wx.createSelectorQuery()
        .in(this)
        .select('#movable-view')
        .boundingClientRect()
        .exec(res => {
          console.log(res[0])
          if (!res[0]) return
          this.data.elementWidth = res[0].width
          this.data.elementHeight = res[0].height
          if (this.data.position.left || this.data.left) {
            this.data.x = this.data.position.left || this.data.left
          }
          if (this.data.position.right || this.data.right) {
            this.data.x =
              this.data.windowWidth -
              this.data.elementWidth -
              (this.data.position.right
                ? this.data.position.right
                : this.data.right)
          }
          if (this.data.position.top || this.data.top) {
            this.data.y = this.data.position.top || this.data.top
          }
          if (this.data.position.bottom || this.data.bottom) {
            this.data.y =
              this.data.windowHeight -
              this.data.elementHeight -
              (this.data.position.bottom
                ? this.data.position.bottom
                : this.data.bottom)
          }
          console.log(1111, this.data)
          this.setData(this.data)
        })
    },
    /**
     * 触摸结束
     */
    onTouchEnd(e) {
      if (
        Math.abs(this.data.x - e.changedTouches[0].clientX) <
          this.data.elementWidth &&
        Math.abs(this.data.y - e.changedTouches[0].clientY) <
          this.data.elementHeight
      ) {
        return
      }
      const currentX = e.changedTouches[0].clientX
      let currentY = e.changedTouches[0].clientY
      if (currentY <= this.data.vMargin) {
        currentY = this.data.vMargin + this.data.elementHeight / 2
      }
      if (currentY >= this.data.windowHeight - this.data.vMargin) {
        currentY =
          this.data.windowHeight -
          this.data.vMargin -
          this.data.elementHeight / 2
      }
      if (currentX + this.data.elementWidth / 2 > this.data.windowWidth / 2) {
        this.data.x =
          this.data.windowWidth - this.data.hMargin - this.data.elementWidth
        this.data.y = currentY - this.data.elementHeight / 2
      }
      if (currentX + this.data.elementWidth / 2 <= this.data.windowWidth / 2) {
        this.data.x = this.data.hMargin
        this.data.y = currentY - this.data.elementHeight / 2
      }
    },
    /**
     * 触摸开始
     */
    onTouchStart() {
      if (!this.data.animation) {
        this.data.animation = true
      }
    }
  }
})
