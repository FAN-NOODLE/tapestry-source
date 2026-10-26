/**
 * components/common/mp-image/index
 * @name 组件名称
 * @author author <emailAddress>
 * @description 组件描述
 * @property {string} style '' 定义 style 属性可以拿到 style 属性上设置的值
 * @property {boolean} virtualHost false 虚拟化组件节点
 * @property {string} styleIsolation apply-shared 表示页面 wxss 样式将影响到自定义组件，但自定义组件 wxss 中指定的样式不会影响页面；
 * @property {boolean} multipleSlots true 支持多插槽
 */
Component({
  externalClasses: [
    'custom-class',
    'image-class',
    'loading-class',
    'error-class',
    'icon-loading'
  ],
  options: {
    virtualHost: false,
    multipleSlots: false,
    styleIsolation: 'apply-shared'
  },
  behaviors: [],
  properties: {
    src: {
      type: String,
      value: ''
    },
    mode: {
      type: String,
      value: 'scaleToFill'
    },
    webp: {
      type: Boolean,
      value: false
    },
    lazyLoad: {
      type: Boolean,
      value: true
    },
    showMenuByLongpress: {
      type: Boolean,
      value: false
    },
    preview: {
      type: Boolean,
      value: false
    },
    previewList: {
      type: Array,
      value: []
    },
    viewport: {
      type: Object,
      value: { bottom: 0 }
    },
    customStyle: {
      type: String,
      value: ''
    },
    useErrorSlot: {
      type: Boolean,
      value: false
    },
    useLoadingSlot: {
      type: Boolean,
      value: false
    },
    showLoading: {
      type: Boolean,
      value: true
    },
    showError: {
      type: Boolean,
      value: true
    },
    round: {
      type: Boolean,
      value: false
    },
    radius: {
      type: String,
      value: 0
    },
    customLoad: {
      type: Boolean,
      value: false
    }
  },
  data: {
    title: 'mp-image',
    laoded: false,
    laodError: false,
    loading: true,
    show: false
  },
  // 组件生命周期函数
  lifetimes: {
    ready() {
      this.addObserver()
    },
    detached() {
      this.clean()
    }
  },
  methods: {
    /**
     * @description 添加监听
     */
    addObserver() {
      if (!wx.createIntersectionObserver || !this.data.lazyLoad) {
        this.setData({
          show: true
        })
        return
      }
      if (this.observer) return
      try {
        const observer = this.createIntersectionObserver()
        observer
          .relativeToViewport(this.data.viewport)
          .observe('.mp-image', () => {
            if (this.data.customLoad) {
              this.triggerEvent('custom-load', {
                callback: () => {
                  this.setData({
                    show: true
                  })
                  this.clean()
                }
              })
            } else {
              this.setData({
                show: true
              })
              this.clean()
            }
          })
        this.observer = observer
      } catch (e) {
        this.setData({
          show: true
        })
      }
    },
    /**
     * @description 移除监听
     */
    clean() {
      if (this.observer) {
        this.observer.disconnect()
      }
      this.observer = null
    },
    /**
     * @description 点击事件
     */
    onClick(e) {
      const { preview } = this.data
      // 是否可以预览
      if (preview) {
        this.previewImage()
      }
      this.triggerEvent('click', e.detail)
    },
    /**
     * @description 图片加载完成事件
     * @param {object} e 事件对象
     */
    onImageLoaded(e) {
      this.setData({
        loaded: true,
        loading: false
      })
      this.triggerEvent('load', e.detail)
    },
    /**
     * @description 图片加载失败事件
     * @param {object} e 事件对象
     */
    onImageError(e) {
      this.setData({
        laoded: false,
        loading: false
      })
      this.triggerEvent('error', e.detail)
    },
    /**
     * @description 预览图片
     */
    previewImage() {
      const { src, previewList } = this.data

      if (previewList.length === 0) {
        previewList.push(src)
      }

      wx.previewImage({
        current: src,
        urls: previewList
      })
    }
  }
})
