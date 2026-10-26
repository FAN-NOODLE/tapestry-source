
// Component Object
Component({
  options: {
    multipleSlots: true
  },
  properties: {
    show: {
      type: Boolean,
      value: false
    },
    showFirstPreviewTips: {
      type: Boolean,
      value: false
    },
    item: {
      type: Object
    },
    selection: {
      type: Array,
      value: []
    },
    /** 群发还是保存  1保存 2群发 */
    enterType: {
      type: Number,
      value: 1
    },
    zIndex: {
      type: Number,
      value: 88
    }
  },
  data: {
    previewIndex: 0
  },
  methods: {
    /**
     * @description 显示弹窗
     */
    show() {
      this.setData({
        show: true
      })
    },
    /**
     * @description 隐藏弹窗
     */
    hide() {
      this.setData({
        show: false
      }, () => {
        setTimeout(() => {
          this.triggerEvent('close')
        }, 500)
      })
    },
    /**
     * @description 选择
     */
    handleSelect(e) {
      const { status, item, isMergeImage } = e.currentTarget.dataset
      // 返回 当前列表 下标 是否选中
      this.triggerEvent('select', {
        item: {
          type: this.data.item.type,
          ...item
        },
        isMergeImage,
        status
      })
    },
    /**
     * @description: 预览图片
     */
    handlePreviewImage(e) {
      const { index } = e.currentTarget.dataset
      this.triggerEvent('preview', {
        previewIndex: index
      })
    },
    /**
     * @description 保存图片
     */
    handleSave(e) {
      const { item } = this.data
      const { index } = e.currentTarget.dataset
      this.triggerEvent('save', {
        index,
        item: item
      })
    },
    handleSaveAll() {
      const { item, selection } = this.data
      const ids = item.merge_image.map(item => item.id)
      const imageList = selection.filter(x => ids.includes(x.id))
      console.log(imageList)
      if (imageList.length === 0) return

      this.triggerEvent('saveAll', {
        item: item,
        list: imageList
      })
    }
  }
})

