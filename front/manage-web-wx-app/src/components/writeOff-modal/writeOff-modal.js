
Component({
  properties: {
    writeOffModal: Boolean,
    couponObj:Object
  },
  data: {
    // writeOffModal: {}
  },

  /*    系统事件开始   */
  created () { },
  attached () { },
  ready () {
   
  },
  /*    系统事件结束   */

  methods: {
    closeModal(){
      this.setData({
        writeOffModal:false
      })
    }
  }
})
