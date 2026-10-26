// Component Object
Component({
  properties: {
    show: {
      type: Boolean,
      value: false
    }
  },
  data: {
    list: [1, 2, 3, 4]
  },
  methods: {
    show(value) {
      this.setData({
        show: true,
        list: value
      })
    },
    hide() {
      this.setData({
        show: false,
        list: []
      })
    }
  }
})

