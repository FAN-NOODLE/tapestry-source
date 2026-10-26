/**
 * components/common/mp-dragable/index
 * @name 组件名称
 * @author author <emailAddress>
 * @description 组件描述
 * @property {string} style '' 定义 style 属性可以拿到 style 属性上设置的值
 * @property {boolean} virtualHost false 虚拟化组件节点
 * @property {string} styleIsolation apply-shared 表示页面 wxss 样式将影响到自定义组件，但自定义组件 wxss 中指定的样式不会影响页面；
 * @property {boolean} multipleSlots true 支持多插槽
 */

import { getAllRect, getRect, getCurrentPage } from '../common/utils'

Component({
  externalClasses: ['custom-class'],
  options: {
    virtualHost: true,
    multipleSlots: true,
    styleIsolation: 'apply-shared'
  },
  behaviors: [],
  properties: {
    customStyle: {
      type: String
    },
    selector: {
      type: String,
      value: '.drag-item'
    },
    list: {
      type: Array,
      value: []
    },
    disabled: {
      type: Boolean,
      value: false
    }
  },
  data: {
    title: 'mp-dragable',
    /** 拖拽区域 X */
    areaX: 0,
    /** 拖拽区域 Y */
    areaY: 0,
    /** 拖拽元素位置 X */
    x: 0,
    /** 拖拽元素位置 Y */
    y: 0,
    /** 触点位置 X */
    sX: 0,
    /** 触点位置 Y */
    sY: 0,
    /** 所有可以拖拽的元素 */
    rects: [],
    /** 当前拖动元素的下标 */
    current: -1,
    /** 是否可以拖 */
    isDrag: false
  },
  methods: {
    /**
     * @description 初始化
     */
    async init() {
      if (this.data.disabled) return
      const page = getCurrentPage()
      if (!page) return
      const [area, rects] = await Promise.all([
        getRect(this, '#movable-area'),
        getAllRect(page, this.data.selector)
      ])
      this.data.areaX = area.left
      this.data.areaY = area.top
      this.data.rects = rects.map(item => {
        return {
          id: item.id,
          x: item.left - area.left,
          y: item.top - area.top,
          width: item.width,
          height: item.height
        }
      })
    },
    /**
     * @description 长按启动拖拽
     */
    onLongPress() {
      this.setData({
        isDrag: true
      })
      this.triggerEvent('dragBefore', { current: this.data.current })
    },
    /**
     * @description 触屏开始
     * @param {*} e
     */
    async onTouchStart(e) {
      if (this.data.disabled) return
      // 获取数据
      await this.init()
      const { rects, areaX, areaY } = this.data
      const { clientX, clientY } = e.changedTouches[0]
      const index = rects.findIndex(item =>
        this.pointInRect({ x: clientX - areaX, y: clientY - areaY }, item),
      )
      if (index !== -1) {
        const { x, y, width, height } = rects[index]
        this.setData({
          x,
          y,
          width,
          height,
          sX: clientX,
          sY: clientY,
          current: index
        })
      }
    },
    /**
     * @description 触屏移动
     * @param {*} e
     */
    onTouchMove(e) {
      if (this.data.disabled) return
      const { sX, sY, x, y } = this.data
      const { clientX, clientY } = e.changedTouches[0]
      this.setData({
        x: x + clientX - sX,
        y: y + clientY - sY,
        sX: clientX,
        sY: clientY
      })
    },
    /**
     * @description 触屏结束
     * @param {*} e
     */
    onTouchEnd(e) {
      if (this.data.disabled) return
      const { clientX, clientY } = e.changedTouches[0]
      const { rects, list, current, areaX, areaY } = this.data

      const index = rects.findIndex(item =>
        this.pointInRect({ x: clientX - areaX, y: clientY - areaY }, item),
      )
      if (index !== -1) {
        const value = list.splice(current, 1)[0]
        list.splice(index, 0, value)
      }
      this.setData({
        sX: 0,
        sY: 0,
        isDrag: false
      })
      this.triggerEvent('change', {
        list,
        current: -1
      })
    },
    /**
     * 判断点是否在矩形内
     * @param {Object} point 要判断的点，形如 { x, y }
     * @param {Object} rect 矩形，形如 { x, y, width, height }
     * @returns {boolean} 点在矩形内返回 true，否则返回 false
     */
    pointInRect(point, rect) {
      return (
        point.x >= rect.x &&
        point.x <= rect.x + rect.width &&
        point.y >= rect.y &&
        point.y <= rect.y + rect.height
      )
    }
  }
})
