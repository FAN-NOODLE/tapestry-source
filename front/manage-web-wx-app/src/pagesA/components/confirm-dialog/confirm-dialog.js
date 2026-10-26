Component({
  properties: {
    popShow: {
      type: Boolean,
      value: false
    },
    content: {
      type: String,
      value: '是否完成任务？'
    }
  },
  methods: {
    /**
     * 取消
     */
    cancelBtn () {
      this.triggerEvent('confirm', { type: 'cancel' })
    },
    /**
     * 完成
     */
    confirmBtn () {
      this.triggerEvent('confirm', { type: 'ok' })
    }
  }
})
