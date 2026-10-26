import { getRect } from '../../utils/util'
import { getMaterialTagListApi } from '@/api/material-api'

// Component Object
Component({
  properties: {
    show: {
      type: Boolean,
      value: false
    },
    isCollect: {
      type: Number,
      value: 0
    }
  },
  data: {
    list: [],
    top: 0,
    checkedList: []
  },
  lifetimes: {
    attached() {
      this.getMaterialTagList()
    },
    ready() {
      getRect(this, '.btn').then(res => {
        this.setData({
          top: res.bottom + 20
        })
      }).catch(err => { console.log(err) })
    }
  },
  methods: {
    /**
     * @description 获取标签列表
     */
    getMaterialTagList() {
      getMaterialTagListApi({
        page: 1,
        size: 999,
        is_collect: this.data.isCollect
      }).then(res => {
        this.setData({
          list: res.data.list
        })
      }).catch(err => {
        console.log(err)
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
    handleSelect(e) {
      const { index } = e.currentTarget.dataset
      const { list } = this.data
      list[index].checked = !list[index].checked
      // todo: 获取选中数量
      const checkedList = list.filter(item => item.checked)
      this.setData({
        list,
        checkedList
      })
    },
    /**
     * @description 确认选择
     */
    handleConfirm() {
      const { list, checkedList } = this.data
      if (checkedList.length === 0) {
        return
      }
      this.triggerEvent('confirm', {
        value: list.filter(item => item.checked)
      })
      this.hide()
    },
    /**
     * @description 取消选择
     */
    handleReset() {
      const { list } = this.data
      list.forEach(item => {
        item.checked = false
      })
      this.setData({
        list,
        checkedList: []
      })
      this.triggerEvent('confirm', {
        value: []
      })
      this.hide()
    }
  }
})

