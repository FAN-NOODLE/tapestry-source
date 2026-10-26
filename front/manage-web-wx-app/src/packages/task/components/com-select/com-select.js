// Component Object
Component({
  properties: {
    value: {
      type: Number,
      observer: function(value) {
        this.setData({
          current: this.data.list.findIndex(item => item.value === value)
        })
      }
    }
  },
  data: {
    list: [
      {
        label: '全部顾客',
        value: 0
      },
      {
        label: '已联系顾客',
        value: 1
      },
      {
        label: '待联系顾客',
        value: 2
      }
    ],
    current: 0,
    show: false
  },
  methods: {
    handleToggle() {
      this.setData({
        show: !this.data.show
      })
    },
    handleChange(e) {
      const { list } = this.data
      const { index } = e.currentTarget.dataset
      if (index === this.data.current) return

      this.setData({
        current: index
      })
      this.triggerEvent('change', {
        value: list[index].value
      })
      this.handleToggle()
    }
  }
})
