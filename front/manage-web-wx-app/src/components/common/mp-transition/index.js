/**
 * components/common/mp-transition/index
 * @name 组件名称
 * @author author <emailAddress>
 * @description 组件描述
 * @property {string} style '' 定义 style 属性可以拿到 style 属性上设置的值
 * @property {boolean} virtualHost false 虚拟化组件节点
 * @property {string} styleIsolation apply-shared 表示页面 wxss 样式将影响到自定义组件，但自定义组件 wxss 中指定的样式不会影响页面；
 * @property {boolean} multipleSlots true 支持多插槽
 */
import { transition } from '../mixins/transition';
Component({
  externalClasses: [
    'custom-class',
    'enter-class',
    'enter-active-class',
    'enter-to-class',
    'leave-class',
    'leave-active-class',
    'leave-to-class',
  ],
  options: {
    virtualHost: false,
    multipleSlots: true,
    styleIsolation: 'apply-shared',
  },
  behaviors: [transition(false)],
  properties: {
    customStyle: String,
  },
});
