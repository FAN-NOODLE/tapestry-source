/**
 * components/common/mp-observer/index
 * @name 组件名称
 * @author author <emailAddress>
 * @description 组件描述
 * @property {string} style '' 定义 style 属性可以拿到 style 属性上设置的值
 * @property {boolean} virtualHost false 虚拟化组件节点
 * @property {string} styleIsolation apply-shared 表示页面 wxss 样式将影响到自定义组件，但自定义组件 wxss 中指定的样式不会影响页面；
 * @property {boolean} multipleSlots true 支持多插槽
 */

Component({
  externalClasses: ['custom-class'],
  options: {
    virtualHost: false,
    multipleSlots: true,
    styleIsolation: 'apply-shared'
  },
  behaviors: [],
  properties: {
    thresholds: {
      type: Array,
      value: [0, 1]
    },
    viewport: {
      type: Object,
      value: { bottom: 200 }
    }
  },
  data: {
    title: 'mp-observer'
  },
  lifetimes: {
    ready() {
      this.init()
    },
    detached() {
      this.clean()
    }
  },
  methods: {
    /**
     * @description 初始化监听
     */
    init() {
      this.observer = wx.createIntersectionObserver(this, {
        thresholds: this.data.thresholds
      })
      this.observer
        .relativeToViewport(this.data.viewport)
        .observe('.mp-observer', res => {
          if (res.intersectionRatio === 1) {
            this.triggerEvent('change', { value: true })
          } else {
            this.triggerEvent('change', { value: false })
          }
        })
    },
    /**
     * @description 移除监听
     */
    clean() {
      if (this.observer) {
        this.observer.disconnect()
      }
      this.observer = null
    }
  }
})
