Component({
  properties: {
    filterItem: {
      type: Object,
      value: {}
    },
    itemResult: {
      type: Object | Number,
      value: [] | 0
    }

  },
  data: {
    result: [],
    isShowExpand: false,
    expand: true
  },
  observers: {
    'filterItem.list': function (v) {
      if (v.length > 0) {
        this.computedDoubleLine()
      }
    },
    itemResult: function(v) {
      if (v) {
        const { filterItem: { multiple } } = this.properties
        let result = []
        if (v) {
          if (multiple) {
            result = v
          } else {
            result = [v]
          }
        }
        this.setData({ result })
      }
    }
  },

  /*    系统事件开始   */
  created () {
  },
  attached () {},
  ready () {
    const { itemResult, filterItem: { multiple } } = this.properties
    let result = []
    if (itemResult) {
      if (multiple) {
        result = itemResult
      } else {
        result = [itemResult]
      }
    }
    this.setData({ result })
    this.triggerEvent('change', {
      value: result,
      item: this.properties.filterItem
    })
  },
  /*    系统事件结束   */

  methods: {
    /*    普通事件开始  */
    tapSelectItem (e) {
      const { filterItem: { multiple } } = this.properties
      let { result } = this.data
      const { key } = e.currentTarget.dataset
      const index = result.findIndex(v => v === key)

      // 取消选择
      if (index > -1) {
        result.splice(index, 1)
      } else {
        // 多选
        if (multiple) {
          result.push(key)
        } else {
          // 单选
          result = [key]
        }
      }

      this.setData({ result })
      this.triggerEvent('change', {
        value: result,
        item: this.properties.filterItem
      })
    },
    /**
     * 计算是否超过两行
     */
    computedDoubleLine () {
      var query = wx.createSelectorQuery().in(this)
      query.selectAll('.filter-list').boundingClientRect(labelDoms => {
        console.log(labelDoms)
        if (labelDoms.length) {
          const height = labelDoms[0].bottom - labelDoms[0].top

          this.setData({ isShowExpand: height > (39 * 1.5) })
        }
      }).exec()
    },
    /**
     * 展开折叠
     */
    expand () {
      this.setData({ expand: !this.data.expand })
    },
    /**
     * 显示提示
     */
    showTips () {
      this.triggerEvent('showTips', this.properties.filterItem)
    },
    handleSelectDate() {
      this.triggerEvent('select-date', this.properties.filterItem)
    }
  }
})
