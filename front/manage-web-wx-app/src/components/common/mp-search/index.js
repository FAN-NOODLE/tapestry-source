// Component Object
Component({
  externalClasses: [
    'custom-class'
  ],
  properties: {
    value: String,
    placeholder: {
      type: String,
      value: '请搜索顾客昵称和备注名'
    }
  },
  data: {

  },
  methods: {
    onBlur(e) {
      this.triggerEvent('blur', e.detail)
    }
  }
})

