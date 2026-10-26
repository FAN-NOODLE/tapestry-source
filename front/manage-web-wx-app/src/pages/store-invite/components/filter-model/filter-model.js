import { StatusList, sourceList } from '../../constants'

/** 生成数据来源 */
function generateOptions () {
  return [
    {title: '状态', key: 'status', list: [{id: -1, status_name: '全部'}].concat(StatusList)},
    {title: '来源', key: 'order_origin', list: [{id: 0, status_name: '全部'}].concat(sourceList)}
  ]
}

Component({
  // 使用全局 & 页面样式
  options: {
    addGlobalClass: true
  },
  data: {
    /** 是否展示、关闭预约信息弹窗 */
    show: false,
    /** 渲染options */
    options: generateOptions(),
    /** 选中列表 */
    checkedObject: {}
  },
  methods: {
    /** 点击（选中）tag */
    clickTag (e) {
      const {key, id} = e.target.dataset

      this.data.checkedObject[key] = id
      this.setData({
        checkedObject: this.data.checkedObject
      })
      console.log('选中列表', key, id)
    },

    /** 打开、关闭弹窗, 供外部调用, 注意：子组件是根据id 选中， 需要父组件传递{对应的key: 对应的id} */
    updatePopup (show, result = {}) {
      this.setData({
        show,
        checkedObject: result
      })
    },

    /** 重置数据 */
    resetData () {
      this.setData({
        show: false
      })
      this.triggerEvent('onReset')
    },

    /** 确定提交筛选数据 */
    confirm () {
      this.setData({
        show: false
      })
      this.triggerEvent('onSubmit', this.data.checkedObject)
    }
  }
})
