import { systemInfo, getNavInfo } from '@/utils/getSystemInfo'
import { getRect } from '../common/utils'

Component({
  options: {
    virtualHost: false,
    multipleSlots: true,
    styleIsolation: 'shared'
  },
  externalClasses: ['custom-class', 'title-class'],
  properties: {
    title: String,
    fixed: {
      type: Boolean,
      value: true,
      observer: 'setHeight'
    },
    placeholder: {
      type: Boolean,
      value: true,
      observer: 'setHeight'
    },
    leftText: String,
    rightText: String,
    customStyle: String,
    leftArrow: {
      type: Boolean,
      value: true
    },
    useBeforeClick: Boolean,
    autoBack: {
      type: Boolean,
      value: true
    },
    zIndex: {
      type: Number,
      value: 1
    },
    opacity: {
      type: Number,
      value: 1
    },
    backgroundColor: {
      type: String,
      value: '#fff'
    },
    safeAreaInsetTop: {
      type: Boolean,
      value: true
    }
  },
  data: {
    /** 导航栏的高度 */
    height: 0,
    /** 状态栏的高度 */
    statusBarHeight: 0,
    /** 导航栏整体的高度 */
    navbarHeight: 0,
    /** 导航栏整体的高度 */
    totalHeight: 0,
    /** 页面只有一个路由 */
    onlyOnePage: false
  },
  lifetimes: {
    attached() {
      this.init()
    },
    detached() {}
  },
  methods: {
    /** 初始化 */
    init() {
      if (systemInfo.brand !== 'devtools') {
        wx.setNavigationBarColor({
          frontColor: this.data.color,
          backgroundColor: this.data.backgroundColor
        })
      }
      const { width, height, statusBarHeight, totalHeight } = getNavInfo()
      this.setData({
        statusBarHeight,
        height,
        width,
        totalHeight,
        onlyOnePage: getCurrentPages().length === 1
      })
      this.setHeight()
    },
    /**
     * 获取状态栏高度
     * @returns
     */
    setHeight() {
      if (!this.data.fixed || !this.data.placeholder) {
        return
      }
      wx.nextTick(() => {
        getRect(this, '.mp-navigation').then(res => {
          if (res && 'height' in res) {
            this.setData({ navbarHeight: res.height })
          }
        })
      })
    },
    /**  */
    onClick() {
      if (this.data.onlyOnePage) {
        wx.reLaunch({
          url: '/pages/index/index'
        })
      } else {
        wx.navigateBack({
          delta: 1
        })
      }
    },
    /**
     * @description 左侧点击
     */
    onClickLeft() {
      this.onBeforeClick().then(() => {
        if (this.data.autoBack) this.onClick()
        else this.triggerEvent('click-left')
      })
    },
    /**
     * @description 右侧点击
     */
    onClickRight() {
      this.triggerEvent('click-right')
    },
    /**
     * @description 切换前拦截
     * @returns
     */
    onBeforeClick() {
      const { useBeforeClick } = this.data
      if (!useBeforeClick) return Promise.resolve()

      return new Promise((resolve, reject) => {
        this.triggerEvent('before-click', {
          callback(status) {
            status ? resolve() : reject()
          }
        })
      })
    },
    /**
     * @description 页面滚动
     */
    scroll(e) {
      const max = 100
      if (e.scrollTop < max) {
        const opacity = e.scrollTop / max
        this.setData({
          opacity: opacity > 0 ? opacity : 0
        })
      }
    }
  }
})
