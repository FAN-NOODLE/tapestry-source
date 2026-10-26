// Component Object
Component({
  properties: {
    show: {
      type: Boolean,
      value: false
    },
    img: {
      type: String,
      value: ''
    }
  },
  data: {
    list: []
  },
  methods: {
    show() {
      this.setData({
        show: true
      })
    },
    hide() {
      this.setData({
        show: false
      })
    }
  }
})

