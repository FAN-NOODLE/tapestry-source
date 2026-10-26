Component({
  properties: {
    // 缺省一级文案
    slogan: String,
    // 缺省二级文案
    content: String,
    // 缺省图
    icon: String,
    // 操作按钮
    btnValue: String,
    // 缺省图距离顶部的高度
    marginTop: {
      type: Number,
      value: 30
    },
    center: {
      type: Boolean,
      value: false
    }
  },


  methods: {
    /*    行为事件开始   */
    goTap () {
      this.triggerEvent('goclick')
    }
    /*    行为事件结束   */
  }
})
