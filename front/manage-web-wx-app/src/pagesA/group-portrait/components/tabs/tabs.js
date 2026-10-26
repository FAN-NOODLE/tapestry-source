
Component({
  /** 自定义样式 */
  externalClasses: ['custom-tab-header', 'custom-tab-item'],
  properties: {
    /** table list */
    list: { type: Array },
    /** 选中tabs index */
    currentIndex: { type: Number }
  },
  methods: {
    /** 切换tab */
    clickTab (e) {
      // 判断重复点击
      if (this.properties.currentIndex === e.target.dataset.current) return false
      // 通知父组件tab 选中变化
      this.triggerEvent('changeTab', e.target.dataset.current)
    }
  }
})

