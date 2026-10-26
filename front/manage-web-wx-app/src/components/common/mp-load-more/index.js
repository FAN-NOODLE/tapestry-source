/**
 * components/common/mp-load-more/index
 * @name 组件名称
 * @author author <emailAddress>
 * @description 组件描述
 * @property {string} style '' 定义 style 属性可以拿到 style 属性上设置的值
 * @property {boolen} virtualHost false 虚拟化组件节点
 * @property {string} styleIsolation apply-shared 表示页面 wxss 样式将影响到自定义组件，但自定义组件 wxss 中指定的样式不会影响页面；
 * @property {boolen} multipleSlots true 支持多插槽
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
    csutomStyle: {
      type: String
    },
    status: {
      type: String,
      value: 'loading'
    },
    loadingText: {
      type: String,
      value: '加载中...'
    },
    loadedText: {
      type: String,
      value: '～～我是有底线的～～'
    },
    loadMoreText: {
      type: String,
      value: '加载更多'
    },
    loadErrorText: {
      type: String,
      value: '加载失败'
    }
  },
  data: {
    title: 'mp-load-more'
  },
  // 组件生命周期函数
  lifetimes: {
    attached() {
      // 在组件实例进入页面节点树时执行
    },
    detached() {
      // 在组件实例被从页面节点树移除时执行
    }
  },
  pageLifetimes: {
    // 组件所在页面的生命周期函数
    show() {},
    hide() {},
    resize() {}
  },
  methods: {
  }
})
