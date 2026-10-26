Component({
  /** 可编辑样式 */
  externalClasses: ['custom-content'],
  behaviors: ['wx://component-export'],
  /** 指定导出 */
  export () {
    return {
      open: this.open.bind(this),
      close: this.close.bind(this)
    }
  },
  options: {
    addGlobalClass: true,
    multipleSlots: true // 复数插槽: 是
  },
  properties: {
    /**
       * top          顶部中间位置
       * top-start    顶部左侧位置
       * top-end      顶部右侧位置
       * bottom       底部中间位置
       * bottom-start 底部左侧位置
       * bottom-end   底部右侧位置
       */
    placement: {
      type: String,
      value: 'top'
    },

    /** 主题, 黑色主题：dark 白色主题：light */
    theme: {
      type: String,
      value: 'dark'
    },
    offsetY: {
      type: Number,
      value: 0
    },
    offsetX: {
      type: Number,
      value: 0
    }

  },
  data: {
    left: -1000,
    top: -1000,
    show: false
  },
  methods: {
    /** 初始化内容 */
    init () {
      // 获取弹窗dom  #popover
      // 获取点击元素dom  #targetPopover
      // eslint-disable-next-line promise/catch-or-return
      Promise.all([this.getScrollTop(), this.getRect('#targetPopover'), this.getRect('#popover')]).then(res => {
        const [scrollTop, targetDom, popoverDom] = res
        console.log(111111, scrollTop, targetDom, popoverDom)
        this.setRect(scrollTop, targetDom, popoverDom)
      })
    },
    /** 打开弹窗 */
    open () {
      this.init()
    },

    /** 关闭弹窗 */
    close () {
      this.setData({
        show: false,
        left: -1000,
        right: -1000
      })
    },

    /** 计算left 距离 */
    calcList (targetDom, popoverDom) {
      const { placement, offsetX } = this.properties
      let left = this.data.left
      if (placement === 'top-end') {
        // 居右对齐
        left = targetDom.right - popoverDom.width
      } else if (placement === 'top-start') {
        // 居左对齐
        left = targetDom.left
      } else {
        // 居中对齐
        left = targetDom.left - popoverDom.width / 2 + targetDom.width / 2
      }

      return left + offsetX
    },

    /** 根据 placement 设置页面position。 */
    setRect (scrollTop, targetDom, popoverDom) {
      const { placement } = this.properties
      let { left, top, offsetY } = this.data

      left = this.calcList(targetDom, popoverDom)
      if (placement === 'top' || placement === 'top-end' || placement === 'top-start') {
        // 顶部
        top = targetDom.top + scrollTop - popoverDom.height - 6
      } else if (placement === 'bottom' || placement === 'bottom-end' || placement === 'bottom-start') {
        // 底部
        top = targetDom.bottom + 6 + scrollTop + offsetY
      } else {
        // 无效值
        wx.showToast('请传入正确的placement')
      }

      console.log('left', left)
      console.log('top', top)
      this.setData({
        left,
        top
      })
    },

    /** 获取元素到顶部的距离 */
    getScrollTop () {
      return new Promise(resolve => {
        const query = wx.createSelectorQuery()
        query.selectViewport().scrollOffset()
        query.exec(function (res) {
          resolve(res[0].scrollTop)
        })
      })
    },

    /** 获取dom 位置 */
    getRect (id) {
      this.setData({
        show: true
      })
      return new Promise(resolve => {
        this.createSelectorQuery().select(id).boundingClientRect(function (rect) {
          resolve(rect)
        }).exec()
      })
    }
  }
})

