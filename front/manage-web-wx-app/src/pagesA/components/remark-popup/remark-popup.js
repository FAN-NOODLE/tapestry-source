// pagesA/task-tags/components/remark-popup/remark-popup.js
import { updateRemarkApi } from '@/api/plan-task-api'

Component({
  /**
   * 组件的属性列表
   */
  properties: {
    /** 显示隐藏控制 */
    show: {
      type: 'boolean',
      value: true
    },
    content: {
      type: 'object',
      value: () => {
        return {
          item: {}, // 子任务
          index: -1 // 数组中的索引
        }
      }
    },
    /** 子任务状态 0待完成,1已完成,2已过期 */
    status: {
      type: 'number',
      value: 0
    }
  },
  /**
   * 组件的初始数据
   */
  data: {
    /** 说明 */
    remark: '',
    /** 最多50个字 */
    maxLength: 50,
    /** 当前统计字数 */
    nowCount: 0,
    /** 提交中 */
    loading: false,
    tools: ['电话无人接听', '手机号码空号', '顾客正忙，下次联系'],
    /** 光标的位置 */
    cursorIndex: 0,
    /** 是否失去焦点 */
    isBlur: true
  },

  /**
   * 组件的方法列表
   */
  methods: {
    /** 文本框 */
    changeInput (e) {
      let { value } = e.detail
      value = value.trim()
      let len = value.length // 获取textarea的内容长度
      this.setData({
        nowCount: len,
        remark: value
      })
    },
    /** 说明弹窗反选 */
    init () {
      const { content } = this.properties
      let len = content.item.remark.length // 获取textarea的内容长度
      this.setData({
        cursorIndex: len, // 掺入文案的位置
        nowCount: len,
        remark: content.item.remark
      })
    },
    /** 保存弹窗 */
    async confirmBtn () {
      if (this.properties.status === 2) {
        // 子任务已过期
        return
      }

      if (this.data.loading) {
        return
      }

      this.data.loading = true
      const { content } = this.properties
      const { item, index } = content
      const { remark } = this.data
      const params = {
        remark,
        task_allot_log_id: item.task_allot_log_id,
        union_id: item.union_id,
        customer_id: item.customer_id
      }
      try {
        const { code, message } = await updateRemarkApi(params)
        this.data.loading = false
        if (code !== 200) {
          console.log(message)
          wx.showToast(message)
          return
        }
        // 同步成功后，更新页面并关闭弹窗
        this.triggerEvent('confirm', {
          type: 'ok',
          data: {
            remark,
            index
          }
        })
      } catch (error) {
        this.data.loading = false
      }
    },
    /** 取消 */
    cancelBtn () {
      this.triggerEvent('confirm', { type: 'cancel' })
    },
    /** 插入文本到输入框 */
    insertBefore (e) {
      const { value } = e.currentTarget.dataset
      const { maxLength, nowCount } = this.data
      if (value.length + nowCount > maxLength) {
        return wx.showToast('超出长度限制')
      }
      // 插入
      this.insertOperate(value)
    },
    /** 把文本插入到内容中 */
    insertOperate (insertContent) {
      if (!this.data.isBlur) {
        setTimeout(() => {
          this.insertOperate(insertContent)
        }, 100)
      } else {
        const { cursorIndex, remark } = this.data
        console.log('插入位置', cursorIndex)
        const newEditValue = `${remark.slice(
          0,
          cursorIndex
        )}${insertContent}${remark.slice(cursorIndex)}`
        let len = newEditValue.length // 获取textarea的内容长度
        this.setData({ nowCount: len, remark: newEditValue })
      }
    },
    /**
     * 编辑器失焦
     */
    editorBlur (e) {
      console.log('编辑器失焦', e)
      this.data.cursorIndex = e.detail.cursor
      this.data.isBlur = true
    },
    /**
     * 编辑器聚焦
     */
    editorFocus () {
      this.data.isBlur = false
    }
  },
  lifetimes: {
    attached: function () {
      // 在组件实例进入页面节点树时执行
      this.init()
    },
    detached: function () {
      // 在组件实例被从页面节点树移除时执行
    }
  }
})
