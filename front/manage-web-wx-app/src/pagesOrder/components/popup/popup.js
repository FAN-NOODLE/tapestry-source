/**
 * Created By Fanyj@2019.11.13
 * 弹出层
 * properties
 * show         是否显示          false
 * zIndex       层级              100
 * mask         是否遮罩          true
 * position     弹出位置          center/top/left/right/bottom
 * round        是否圆角          false
 * event
 * maskClick    点击遮罩
 * */
Component({
  options: {
    multipleSlots: true // 启用多slot支持
  },
  // 属性
  properties: {
    show: {
      type: Boolean,
      value: false
    },
    zIndex: {
      type: Number,
      value: 100
    },
    mask: {
      type: Boolean,
      value: true
    },
    position: {
      type: String,
      value: 'center'
    },
    round: {
      type: Boolean,
      value: false
    }
  },
  data: {
    maskClass: '',
    showClass: '',
    roundClass: ''
  },
  observers: {
    show: function (show) {
      const p = this.data.position
      const _self = this
      this._timerId && clearTimeout(this._timerId)
      this.setData({
        showClass: show ? `position-${p}` : `position-${p}-hide`,
        maskClass: show ? 'popup-mask' : 'popup-mask-hide',
        maskContainerClass: 'popup-content-container'
      })
      if (!show) {
        this._timerId = setTimeout(() => {
          _self.setData({
            showClass:
              p === 'center'
                ? 'position-center-hide-all'
                : _self.data.showClass,
            maskClass: 'popup-mask-hide-all',
            maskContainerClass: 'popup-mask-hide-all'
          })
        }, 500)
      }
    }
  },
  methods: {
    maskClick: function () {
      this.triggerEvent('maskClick')
    },
    catchTap () {},
    touchMove () {}
  },
  lifetimes: {
    ready () {
      const { position, round, show } = this.data
      this.setData({
        roundClass: round ? `popup-round-${position}` : '',
        showClass: show
          ? `position-${position}`
          : `position-${position}-hide${position === 'center' ? '-all' : ''}`,
        maskContainerClass: 'popup-mask-hide-all'
      })
    }
  }
})
