import { debounce } from '@/utils/tools'

// Component Object
Component({
  properties: {
    show: {
      type: Boolean,
      value: false
    },
    tip: {
      type: String,
      value: '点击确定后，所选的{count}位顾客将分配给顾问{name}'
    },
    content: {
      type: String,
      value: '您好，由于个人工作变更，后续将由我的同事【接替者名字】接替我的工作，继续为您服务。'
    }

  },
  data: {
    /** 是否编辑 */
    isEdit: false,
    /** 提示 */
    innerTip: '',
    /** 顾问名称 */
    innerName: '',
    /** 编辑的内容 */
    innerContent: '',
    /** 是否获取焦点 */
    isBlur: true,
    /** 光标位置 */
    cursorIndex: 0,
    /** 富文本内容 */
    innerHtmlContent: ''
  },
  methods: {
    /**
     * show
     * @param {string} name 顾问姓名
     * @param {number} count 需要转移到顾客数量
     */
    show(name, count) {
      const { content, tip } = this.data
      const innerName = `【${name}】`
      this.setData({
        show: true,
        innerName,
        innerContent: content,
        innerHtmlContent: content.replace('【接替者名字】', '<b>【接替者名字】</b>'),
        innerTip: tip.replace('{count}', count).replace('{name}', innerName)
      })
    },
    hide() {
      this.setData({
        show: false
      })
    },
    /**
     * @description 编辑
     */
    hanleEdit() {
      this.setData({
        isEdit: true
      })
    },

    /**
     * 取消编辑
     */
    handleEditCancel () {
      const { content } = this.data
      this.setData({
        isEdit: false,
        innerContent: content
      })
    },
    /**
     * 确认编辑
     */
    handleEditConfirm () {
      const { innerContent } = this.data
      const innerHtmlContent = innerContent.replace(/【接替者名字】/g, '<b>【接替者名字】</b>')
      this.setData({
        isEdit: false,
        content: innerContent,
        innerHtmlContent: innerHtmlContent
      })
    },
    /**
     * 编辑器输入
     */
    handleInput (e) {
      console.log('编辑器输入', e)
      this.data.innerContent = e.detail.value
    },
    /**
     * 编辑器聚焦
     */
    handleFocus (e) {
      console.log('聚焦位置', e)
      this.data.isBlur = false
    },
    /**
     * 编辑器失焦
     */
    handleBlur (e) {
      console.log('编辑器失焦', e)
      this.data.cursorIndex = e.detail.cursor
      this.data.isBlur = true
    },
    /**
     * 插入接替者昵称
     */
    handleInsertName () {
      if (!this.data.isBlur) {
        setTimeout(() => { this.handleInsertName() }, 100)
      } else {
        const { innerContent, cursorIndex } = this.data
        console.log('插入位置', cursorIndex)
        const content = `${innerContent.slice(0, cursorIndex)}【接替者名字】${innerContent.slice(cursorIndex)}`
        this.setData({ innerContent: content })
      }
    },
    /** 取消 */
    handleCancel() {
      this.setData({ show: false })
    },
    /** 确定 */
    handleConfirm: debounce(function () {
      const { content } = this.data
      this.triggerEvent('confirm', {
        content
      })
    })
  }
})

