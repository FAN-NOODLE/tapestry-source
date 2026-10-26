/**
 * components/common/mp-overlay/index
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
    show: Boolean,
    customStyle: String,
    duration: {
      type: null,
      value: 300
    },
    zIndex: {
      type: Number,
      value: 1
    },
    lockScroll: {
      type: Boolean,
      value: true
    },
    rootPortal: {
      type: Boolean,
      value: false
    }
  },
  data: {
    title: 'mp-overlay'
  },
  methods: {
    onClick: function () {
      this.triggerEvent('click')
    },
    // for prevent touchmove
    noop: function () {}
  }
})
