/**
 * 半屏弹窗
 * show         是否显示          false
 * */
Component({
  options: {
    /** 启用多slot支持 */
    multipleSlots: true
  },
  /** 可编辑样式 */
  externalClasses: ['custom-dialog', 'custom-mask'],
  properties: {
    /** 判断是否显示弹窗 */
    show: {
      type: Boolean,
      value: false
    }
  }
})

