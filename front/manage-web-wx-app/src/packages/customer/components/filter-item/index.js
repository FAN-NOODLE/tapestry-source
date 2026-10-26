Component({
  options: {
    multipleSlots: true
  },
  properties: {
    item: {
      type: Object,
      value: {}
    },
    value: {
      type: Object | Number,
      value: [] | 0,
      observer(v) {
        const { type } = this.data.item
        let innerValue = []
        // 多选
        if (type === 1) {
        // 单选
          innerValue = v ? [v] : []
        } else if (type === 2) {
          innerValue = v || []
        } else if (type === 3) {
          if (v) {
            innerValue[0] = v.start_time ? v.start_time * 1000 : ''
            innerValue[1] = v.start_time ? v.end_time * 1000 : ''
          } else {
            innerValue = []
          }
        } else {
          innerValue = v || {}
        }
        this.setData({ innerValue })
      }
    },
    options: {
      type: Array,
      value: [],
      observer(v) {
        if (v) {
          this.setData({ showExpand: v.length > 8 })
        }
      }
    }
  },
  data: {
    innerValue: [],
    /** 是否显示展开按钮 */
    showExpand: false,
    /** 展开收起 */
    expand: false,
    /** 显示提示 */
    showTip: false
  },

  /*    系统事件开始   */
  created () {
  },
  attached () {},
  /*    系统事件结束   */

  methods: {
    /**
     * 单选/多选
     * @param {*} e
     */
    handleSelect (e) {
      const { item: { type } } = this.data
      let { innerValue } = this.data

      const { key } = e.currentTarget.dataset
      const index = innerValue.findIndex(v => v === key)

      // 取消选择
      if (index > -1) {
        innerValue.splice(index, 1)
      } else {
        // 多选
        if (type === 1) {
          // 单选
          innerValue = [key]
        } else {
          innerValue.push(key)
        }
      }

      this.setData({ innerValue })
      this.triggerEvent('change', {
        value: type === 1 ? innerValue[0] || '' : innerValue,
        key: this.data.item.key,
        type: this.data.item.type
      })
    },
    /**
     * 展开折叠
     */
    handleExpand () {
      this.setData({ expand: !this.data.expand })
    },
    /**
     * 选择时间
     */
    handleSelectDate() {
      // this.setData({
      //   showPickerDate: true
      // })
      this.triggerEvent('select-date', {
        item: this.data.item,
        value: this.data.innerValue
      })
    },
    /**
     * 显示提示
     */
    handleTips() {
      // this.setData({
      //   showTip: !this.data.showTip
      // })
      this.triggerEvent('show-tips', this.data.item)
    },
    handleInput(e) {
      const { key } = e.currentTarget.dataset
      const { value } = e.detail
      // if (key === 'end_number') {
      //   if (!Number(this.data.innerValue.start_number) && !Number(value)) {
      //     console.log('请输入大于0的数字')
      //     wx.showToast('请输入大于0的数字')
      //     return
      //   }
      // }
      this.triggerEvent('change', {
        value: {
          ...this.data.innerValue,
          [key]: value ? Number(value) : ''
        },
        key: this.data.item.key,
        type: this.data.item.type
      })
    }
  }
})
