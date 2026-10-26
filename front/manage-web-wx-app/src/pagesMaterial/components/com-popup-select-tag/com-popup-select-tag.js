// Component Object
import { getMaterialTagListApi } from '@/api/material-api'

Component({
  externalClasses: ['custom-class'],
  properties: {
    value: {
      type: Array,
      value: [],
      observer: 'init'
    },
    show: {
      type: Boolean,
      value: false
    }
  },
  data: {
    list: [],
    name: ''
  },
  lifetimes: {
    attached() {
      this.getMaterialTagList()
    }
  },
  methods: {
    /**
     * @description 获取标签列表
     */
    getMaterialTagList() {
      getMaterialTagListApi({
        page: 1,
        size: 999
      }).then(res => {
        this.setData({
          list: res.data.list
        }, () => {
          this.init()
        })
      }).catch(err => {
        console.log(err)
      })
    },
    init() {
      const { value, list } = this.data

      list.forEach(item => {
        if (value.includes(item.id)) {
          item.checked = true
        } else {
          item.checked = false
        }
      })
      this.setData({
        list
      })
    },
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
      })
    },
    /**
     * @description 选择
     */
    handleClick(e) {
      const { index } = e.currentTarget.dataset
      const { list } = this.data
      this.setData({
        [`list[${index}].checked`]: !list[index].checked
      })
    },
    /**
     * @description 确认
     */
    handleConfirm() {
      const { list } = this.data
      const value = list.filter(item => item.checked)
      if (value.length) {
        this.hide()
        this.setData({
          name: value.map(item => item.name).join()
        })
        this.triggerEvent('change', { value })
      } else {
        wx.showToast({
          title: '请选择素材标签',
          icon: 'none'
        })
      }
    },
    /**
     * @description 取消
     */
    handleCancel() {
      this.hide()
      this.init()
    }
  }
})

