import { getRect } from '../common/utils'
import { transition } from '../mixins/transition'

const { globalData } = getApp()

Component({
  externalClasses: [
    'overlay-custom-class',
    'triangle-custom-class',
    'content-custom-class',
    'custom-class',
    'enter-class',
    'enter-active-class',
    'enter-to-class',
    'leave-class',
    'leave-active-class',
    'leave-to-class',
    'close-icon-class'
  ],
  options: {
    virtualHost: false,
    multipleSlots: true,
    styleIsolation: 'apply-shared'
  },
  behaviors: [transition(false)],
  properties: {
    content: {
      type: String,
      value: ''
    },
    zIndex: {
      type: Number,
      value: 521
    },
    /**
     * 'top' | 'top-start' | 'top-end' | 'bottom' | 'bottom-start' | 'bottom-end' | 'left' | 'left-start' | 'left-end' | 'right' | 'right-start' | 'right-end'
     */
    placement: {
      type: String,
      value: 'bottom',
      observer: 'observeClass'
    },
    overlay: {
      type: Boolean,
      value: true
    },
    customStyle: {
      type: String,
      value: ''
    },
    duration: {
      type: Number,
      value: 300
    },
    overlayStyle: {
      type: String,
      value: ''
    },
    closeOnClickOverlay: {
      type: Boolean,
      value: true
    },
    closeable: {
      type: Boolean,
      value: false
    },
    closeIcon: {
      type: String,
      value: 'cross'
    },
    closeIconPosition: {
      type: String,
      value: 'top-right'
    },
    lockScroll: {
      type: Boolean,
      value: true
    },
    rootPortal: {
      type: Boolean,
      value: true
    },
    transition: {
      type: String,
      observer: 'observeClass'
    }
  },
  data: {
    styles: {},
    contentStyle: {}
  },
  lifetimes: {
    attached() {
      this.observeClass()
    }
  },
  methods: {
    show(context, selector, content) {
      let placement = this.data.placement
      const { screenHeight, screenWidth } = globalData.system
      const maxHeight = screenHeight / 2
      const maxWidth = 155 / 2
      /** 获取元素位置信息 */
      getRect(context, selector).then(res => {
        const { top, left, right, width, height, dataset } = res
        let contentStyle = {}
        // 判断左右边界
        if (screenWidth - right + width / 2 < maxWidth) {
          // 右边界
          contentStyle = {
            right: 0,
            transform: 'translateX(-90%)'
          }
        } else if (width / 2 + left < maxWidth) {
          // 左边界
          contentStyle = {
            left: 0,
            transform: 'none'
          }
        } else {
          // 中间
          contentStyle = {
            left: '50%',
            transform: 'translateX(-50%)'
          }
        }
        // 判断上下
        if (placement === 'auto') {
          if (top + height / 2 > maxHeight) {
            placement = 'top'
          } else {
            placement = 'bottom'
          }
        }
        // 默认展示在底部
        const styles = {}
        switch (placement) {
          case 'bottom':

            styles.top = top + height + 'px'
            styles.left = left + 'px'
            styles.width = width + 'px'
            break
          case 'top':
            styles.top = top + 'px'
            styles.left = left + 'px'
            styles.width = width + 'px'
            break

          default:
            break
        }
        console.log(placement)
        this.setData({
          styles,
          content: content || dataset.content,
          placement,
          contentStyle: JSON.stringify(contentStyle),
          show: true
        })
      }).catch(err => {
        console.log(err)
      })
    },
    onClickCloseIcon() {
      this.triggerEvent('close')
    },
    onClickOverlay() {
      this.triggerEvent('click-overlay')
      if (this.data.closeOnClickOverlay) {
        this.setData({
          show: false,
          inited: false
        })
        this.triggerEvent('close')
      }
    },
    observeClass() {
      const { transition, placement, duration } = this.data
      var updateData = {
        name: transition || placement
      }
      if (transition === 'none') {
        updateData.duration = 0
        this.originDuration = duration
      } else if (
        this.originDuration !== null &&
        this.originDuration !== undefined
      ) {
        updateData.duration = this.originDuration
      }
      this.setData(updateData)
    },
    noop() {}
  }
})
