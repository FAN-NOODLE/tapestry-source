// Component Object
// Component Object

const { debounce } = require('@/utils/tools')

// Component Object
Component({
  properties: {
    show: {
      type: Boolean,
      value: false,
      observer(newVal) {
        if (newVal) {
          this.setData({
            content: this.data.value || '',
            length: Math.min(this.data.value.length, this.data.maxlength)
          })
        }
      }
    },
    title: {
      type: String,
      value: '客户描述'
    },
    value: {
      type: String,
      value: ''
    },
    placeholder: {
      type: String,
      value: '添加更多描述信息'
    },
    maxlength: {
      type: Number,
      value: 150
    },
    type: {
      type: Number | String,
      value: 1
    },
    confirmText: {
      type: String,
      value: '确认打包'
    },
    cancelText: {
      type: String,
      value: '取消'
    }
  },
  data: {
    content: ''
  },
  methods: {
    handleConfirm: debounce(function() {
      const { content } = this.data
      this.triggerEvent('confirm', {
        value: content
      })
    }),
    handleCancel() {
      this.triggerEvent('cancel')
    },
    handleInput(e) {
      const { value } = e.detail
      const len = Math.min(value.length, this.data.maxlength)
      this.data.content = value.slice(0, len)
      this.setData({
        length: Math.min(value.length, this.data.maxlength)
      })
    },
    handleBlur(e) {
      const { value } = e.detail
      this.data.content = value
    }
  }
})
