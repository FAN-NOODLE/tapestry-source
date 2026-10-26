Component({
  properties: {
    /** 预约客户 */
    info: { type: Object }
  },
  methods: {
    close () {
      this.triggerEvent('close')
    },
    /** 核销 */
    writeOff () {
      // is_check 0 核销
      const { is_check, activity_code } = this.data.info
      if (is_check !== 0) return
      this.triggerEvent('writeOff', activity_code)
    }
  }
})
