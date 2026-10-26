/**
 * components/common/mp-row/index
 * @name 组件名称
 * @author author <emailAddress>
 * @description 组件描述
 * @property {string} style '' 定义 style 属性可以拿到 style 属性上设置的值
 * @property {boolean} virtualHost false 虚拟化组件节点
 * @property {string} styleIsolation apply-shared 表示页面 wxss 样式将影响到自定义组件，但自定义组件 wxss 中指定的样式不会影响页面；
 * @property {boolean} multipleSlots true 支持多插槽
 */
import { useChildren } from '../mixins/relation'

const { relations, mixin } = useChildren('mp-col', function (target) {
  this.init(target)
})
Component({
  relations: relations,
  externalClasses: ['custom-class'],
  options: {
    virtualHost: false,
    multipleSlots: true,
    styleIsolation: 'apply-shared'
  },
  behaviors: [mixin],
  properties: {
    gutter: {
      type: Number,
      value: 0,
      observer: 'setGutter'
    }
  },
  data: {
    title: 'mp-row'
  },
  methods: {
    init(target) {
      const { gutter } = this.data
      if (gutter) {
        target.setData({ gutter: gutter })
      }
    },
    setGutter() {
      if (this.children) {
        this.children.forEach(col => {
          col.setData(this.data)
        })
      }
    }
  }
})
