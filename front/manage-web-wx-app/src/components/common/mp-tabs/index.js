/* eslint-disable promise/catch-or-return */
// Component Object
/**
 * components/common/mp-tabs/index
 * @name 组件名称
 * @author author <emailAddress>
 * @description 组件描述
 * @property {string} style '' 定义 style 属性可以拿到 style 属性上设置的值
 * @property {boolean} virtualHost false 虚拟化组件节点
 * @property {string} styleIsolation apply-shared 表示页面 wxss 样式将影响到自定义组件，但自定义组件 wxss 中指定的样式不会影响页面；
 * @property {boolean} multipleSlots true 支持多插槽
 */

import touch from '../common/touch'
import { useChildren } from '../mixins/relation'
import { getRect, getAllRect } from '../common/utils'

const { relations, mixin } = useChildren('mp-tab', function () {
  this.updateTabs()
})
Component({
  relations,
  externalClasses: [
    'custom-class',
    'nav-class',
    'tab-class',
    'tab-active-class',
    'line-class',
    'wrap-class'
  ],
  options: {
    virtualHost: false,
    multipleSlots: true,
    styleIsolation: 'apply-shared'
  },
  behaviors: [mixin, touch],
  properties: {
    sticky: Boolean,
    border: Boolean,
    swipeable: Boolean,
    titleActiveColor: String,
    titleInactiveColor: String,
    color: String,
    active: {
      optionalTypes: [String, Number],
      value: 0,
      observer: function (name) {
        if (name !== this.getCurrentName()) {
          this.setCurrentIndexByName(name)
        }
      }
    },
    animated: {
      type: Boolean,
      value: false
    },
    lineWidth: {
      type: Number,
      value: 40,
      observer: 'resize'
    },
    lineHeight: {
      type: null,
      value: -1
    },
    ellipsis: {
      type: Boolean,
      value: false
    },
    duration: {
      type: Number,
      value: 0.3
    },
    zIndex: {
      type: Number,
      value: 1
    },
    scrollable: {
      type: Boolean,
      value: false
    },
    swipeThreshold: {
      type: Number,
      value: 5,
      observer: function (value) {
        this.setData({
          scrollable: this.children.length > value || !this.data.ellipsis
        })
      }
    },
    useBeforeChange: {
      type: Boolean,
      value: false
    },
    type: {
      type: String,
      value: 'line'
    }
  },
  data: {
    tabs: [],
    currentIndex: 0,
    lineOffsetLeft: 0,
    scrollLeft: 0,
    scrollable: false,
    skipTransition: true,
    scrollWithAnimation: false,
    inited: false
  },
  lifetimes: {
    ready() {
      setTimeout(() => {
        this.resize()
        this.scrollIntoView()
      }, 1000 / 30)
    }
  },
  methods: {
    /**
     * @description 获取tabs
     */
    updateTabs() {
      const { active, swipeThreshold, ellipsis } = this.data
      const tabs = this.children.map(item => item.data)
      const scrollable = this.children.length > swipeThreshold || !ellipsis
      this.setData({
        tabs,
        scrollable
      })
      this.setCurrentIndexByName(active || this.getCurrentName())
    },
    /**
     * @description 设置通过name当前下标
     */
    setCurrentIndexByName(name) {
      // 获取当前选中项
      const matched = this.children.filter(
        child => child.getComputedName() === name,
      )
      if (matched && matched.length) {
        this.setCurrentIndex(matched[0].index)
      }
    },
    /**
     * @description 设置通过name当前下标
     * @param {Number} currentIndex
     */
    setCurrentIndex(currentIndex) {
      if (
        currentIndex === undefined ||
        currentIndex === null ||
        currentIndex >= this.children.length ||
        currentIndex < 0
      ) {
        return
      }
      this.children.forEach((item, index) => {
        const active = index === currentIndex
        if (active !== item.data.active || !item.inited) {
          item.updateRender(active, this)
        }
      })

      if (currentIndex === this.data.currentIndex) {
        if (!this.data.inited) {
          this.resize()
        }
        return
      }

      const shouldEmitChange = currentIndex !== null
      this.setData({ currentIndex })
      setTimeout(() => {
        this.resize()
        this.scrollIntoView()
      }, 1000 / 30)

      wx.nextTick(() => {
        this.trigger('input')
        if (shouldEmitChange) {
          this.trigger('change')
        }
      })
    },
    /**
     * @description 获取当前默认选中值
     */
    getCurrentName() {
      const activeTab = this.children[this.data.currentIndex]
      if (activeTab) {
        return activeTab.getComputedName()
      }
    },
    /**
     * @description 重新计算底线位置
     */
    resize() {
      const { type, currentIndex, ellipsis, skipTransition, duration } =
        this.data
      if (type !== 'line') return

      Promise.all([
        getAllRect(this, '.mp-tab'),
        getRect(this, '.mp-tabs__line')
      ]).then(([rects, lineRect]) => {
        const rect = rects[currentIndex]
        if (!rect) return

        let lineOffsetLeft = rects
          .slice(0, currentIndex)
          .reduce(function (prev, curr) {
            return prev + curr.width
          }, 0)

        lineOffsetLeft +=
          (rect.width - lineRect.width) / 2 + (ellipsis ? 0 : 8)

        this.setData({
          lineOffsetLeft,
          inited: true
        })

        if (skipTransition) {
          setTimeout(() => {
            this.setData({
              skipTransition: false
            })
          }, duration)
        }
      })
    },
    /**
     * @description 重新计算滚动位置
     */
    scrollIntoView() {
      const { scrollable, currentIndex, scrollWithAnimation } = this.data
      if (!scrollable) return

      Promise.all([
        getAllRect(this, '.mp-tab'),
        getRect(this, '.mp-tabs__nav')
      ]).then(([tabRects, navRect]) => {
        const tabRect = tabRects[currentIndex]
        const offsetLeft = tabRects
          .slice(0, currentIndex)
          .reduce(function (prev, curr) {
            return prev + curr.width
          }, 0)
        this.setData({
          scrollLeft: offsetLeft - (navRect.width - tabRect.width) / 2
        })
        if (!scrollWithAnimation) {
          wx.nextTick(() => {
            this.setData({
              scrollWithAnimation: true
            })
          })
        }
      })
    },
    /**
     * 抛出事件
     * @param {string} eventName
     * @param {Object} child
     * @returns
     */
    trigger(eventName, child) {
      const { currentIndex } = this.data
      const data = this.getChildData(currentIndex, child)
      if (!data) return

      this.triggerEvent(eventName, data)
    },
    /**
     * 获取子元素数据
     * @param {Number} index
     * @param {Object} child
     * @returns
     */
    getChildData(index, child) {
      const currentChild = child || this.children[index]

      if (!currentChild) {
        return {}
      }

      return {
        index: currentChild.index,
        name: currentChild.getComputedName(),
        title: currentChild.data.title
      }
    },
    /**
     * 点击切换
     */
    onTap(e) {
      const { index } = e.currentTarget.dataset
      const child = this.children[index]
      if (child.data.disabled) {
        this.trigger('disabled', child)
        return
      }

      this.onBeforeChange(index).then(() => {
        this.setCurrentIndex(index)
        wx.nextTick(() => {
          this.trigger('click')
        })
      })
    },
    /**
     * @description 切换前拦截
     * @param {number} index
     * @returns
     */
    onBeforeChange(index) {
      const { useBeforeChange } = this.data
      if (!useBeforeChange) return Promise.resolve()

      return new Promise((resolve, reject) => {
        this.triggerEvent('before-change', {
          ...this.getChildData(index),
          callback(status) {
            status ? resolve() : reject(new Error('before-change rejected'))
          }
        })
      })
    },
    /**
     * @description 触点开始
     */
    onTouchStart(e) {
      const { swipeable } = this.data
      if (!swipeable) return
      this.swiping = true
      this.touchStart(e)
    },
    /**
     * @description 触点移动中
     */
    onTouchMove: function (e) {
      if (!this.data.swipeable || !this.swiping) return
      this.touchMove(e)
    },
    /**
     * @description 触点结束
     */
    onTouchEnd: function () {
      if (!this.data.swipeable || !this.swiping) return
      const { direction, deltaX, offsetX } = this
      const minSwipeDistance = 50
      if (direction === 'horizontal' && offsetX >= minSwipeDistance) {
        var index_1 = this.getAvaiableTab(deltaX)
        if (index_1 !== -1) {
          this.onBeforeChange(index_1).then(() => {
            this.setCurrentIndex(index_1)
          })
        }
      }
      this.swiping = false
    },

    getAvaiableTab(direction) {
      const { tabs, currentIndex } = this.data
      var step = direction > 0 ? -1 : 1
      for (
        var i = step;
        currentIndex + i < tabs.length && currentIndex + i >= 0;
        i += step
      ) {
        var index = currentIndex + i
        if (
          index >= 0 &&
          index < tabs.length &&
          tabs[index] &&
          !tabs[index].disabled
        ) {
          return index
        }
      }
      return -1
    }
  }
})

