function getRect(context, selector) {
  return new Promise(function (resolve) {
    wx.createSelectorQuery()
      .in(context)
      .select(selector)
      .boundingClientRect()
      .exec(function (rect) {
        if (rect === void 0) {
          rect = []
        }
        return resolve(rect[0])
      })
  })
}
const systemInfo = wx.getSystemInfoSync()
/**
 * 页面底部安全区适配组件
 * @author gaoshifei <2270029397@qq.com>
 * @description 页面底部安全区适配组件
 * @typedef {object} OptionsParams
 * @property {boolean} virtualHost false 虚拟化组件节点
 * @property {string} styleIsolation apply-shared 表示页面 wxss 样式将影响到自定义组件，但自定义组件 wxss 中指定的样式不会影响页面；
 * @property {boolean} multipleSlots true 支持多插槽
 */
Component({
  externalClasses: ['custom-class'],
  /** @type {OptionsParams} */
  options: {
    virtualHost: false,
    multipleSlots: true,
    styleIsolation: 'apply-shared'
  },
  properties: {
    fixed: {
      type: Boolean,
      value: true
    },
    backgroundColor: {
      type: String,
      value: '#ffffff'
    },
    zIndex: {
      type: Number,
      value: 100
    }
  },
  data: {
    title: 'mp-footer',
    paddingBottom: systemInfo.screenHeight - systemInfo.safeArea.bottom
  },
  lifetimes: {
    ready() {
      wx.nextTick(() => {
        this.init()
      })
    }
  },
  methods: {
    init() {
      getRect(this, '#mp-footer').then(res => {
        if (res) {
          this.setData({ height: res.height })
        }
      }).catch(err => {
        console.log(err)
      })
    }
  }
})
