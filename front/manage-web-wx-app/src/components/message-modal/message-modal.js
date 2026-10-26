import Enums from '../../utils/enums'

const { globalData } = getApp()

Component({
  properties: {
    show: {
      type: Boolean,
      value: false
    },
    // 淡入效果
    fadein: {
      type: Boolean,
      value: false
    },
    description: {
      type: String,
      value: ''
    },
    descriptionTwo: {
      type: String,
      value: ''
    },
    ispartner: {
      type: String,
      value: ''
    },
    cancelText: {
      type: String,
      value: '取消'
    },
    sureText: {
      type: String,
      value: '确定'
    },
    reverse: {
      type: Boolean,
      value: false
    },
    // 是否要展示单个操作按钮
    isSingle: {
      type: Boolean,
      value: false
    },
    // 下方是否要展示按钮
    noBtn: {
      type: Boolean,
      value: false
    },
    isAuthent: {
      type: Boolean,
      value: true
    },
    contentClass: {
      type: String,
      value: ''
    },
    modalW: {
      type: String,
      value: ''
    },
    descBox: {
      type: String,
      value: ''
    },
    descStyle: {
      type: String,
      value: ''
    },
    descStyleTwo: {
      type: String,
      value: ''
    },
    isCatchFun: {
      type: Boolean,
      value: false
    }
  },
  data: {
    showModel: true,
    isPartner: globalData.identityType === Enums.identityType.partner
  },

  /*    系统事件开始   */
  created () {},
  attached () {},
  ready () {
    this.setData({
      isPartner: globalData.identityType === Enums.identityType.partner
    })
  },
  /*    系统事件结束   */

  methods: {
    /*    普通事件开始  */
    cancel () {
      this.triggerEvent('cancel')
    },
    sure () {
      this.triggerEvent('sure')
    },
    close () {
      this.triggerEvent('close')
    },
    catchFun () {
      if (this.properties.isCatchFun) {
        return false
      }
    },
    /*    普通事件结束  */

    /*    cdp分割事件   */
    cdpSplitEvent () {}
    /*    cdp分割事件   */
  }
})
