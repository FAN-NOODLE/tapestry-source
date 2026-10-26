/**
 * components/common/mp-tab/index
 * @name 组件名称
 * @author author <emailAddress>
 * @description 组件描述
 * @property {string} style '' 定义 style 属性可以拿到 style 属性上设置的值
 * @property {boolean} virtualHost false 虚拟化组件节点
 * @property {string} styleIsolation apply-shared 表示页面 wxss 样式将影响到自定义组件，但自定义组件 wxss 中指定的样式不会影响页面；
 * @property {boolean} multipleSlots true 支持多插槽
 */

import { useParent } from '../mixins/relation'

const { relations, mixin } = useParent('mp-tabs')
Component({
  relations,
  externalClasses: ['custom-class'],
  options: {
    virtualHost: false,
    multipleSlots: true,
    styleIsolation: 'apply-shared'
  },
  behaviors: [mixin],
  properties: {
    dot: {
      type: Boolean,
      observer: 'update'
    },
    info: {
      type: null,
      observer: 'update'
    },
    title: {
      type: String,
      observer: 'update'
    },
    disabled: {
      type: Boolean,
      observer: 'update'
    },
    titleStyle: {
      type: String,
      observer: 'update'
    },
    name: {
      optionalTypes: [String, Number],
      value: ''
    }
  },
  data: {
    active: false
  },
  methods: {
    getComputedName() {
      if (this.data.name !== '') {
        return this.data.name
      }
      return this.index
    },
    updateRender(active, parent) {
      var parentData = parent.data
      this.inited = this.inited || active
      this.setData({
        active: active,
        shouldRender: this.inited || !parentData.lazyRender,
        shouldShow: active || parentData.animated
      })
    },
    update() {
      if (this.parent) {
        this.parent.updateTabs()
      }
    }
  }
})
