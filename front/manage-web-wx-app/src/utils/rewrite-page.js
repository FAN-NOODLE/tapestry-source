import watch from './watch'

const originPage = Page // 保存原对象
Page = function (page) {
  const methods = ['onShow', 'onUnload']

  if (page.onLoad) {
    const tempMethod = page.onLoad
    page.onLoad = function (options) {
      const pages = getCurrentPages()
      const currentPage = pages[pages.length - 1]
      // 支持watch
      if (page.watch) {
        for (let key in page.watch) {
          const { handler, deep = false } = page.watch[key]
          watch(currentPage.data, (_key, value, newVal) => {
            if (key === _key) handler.call(currentPage, _key, value, newVal)
          }, deep)
        }
      }
      // 支持computed
      if (page.computed) {
        for (let key in page.computed) {
          currentPage.setData({
            [key]: page.computed[key].call(currentPage)
          })
        }
      }
      tempMethod.call(this, options)
    }
  }
  // 防止预览图片触发onshow
  methods.forEach(name => {
    const tempMethod = page[name]
    if (typeof tempMethod === 'function') {
      page[name] = function (options) {
        const { globalData } = getApp()
        // 兼容垃圾企业微信-ios机型，在图片放大完后，又不触发onshow生命周期
        if (globalData.system && globalData.system.environment === 'wxwork' && globalData.system.model.search('iPhone') !== -1) {
          globalData.isReviewImg = false
        }
        if (globalData.isReviewImg) return (globalData.isReviewImg = false)
        tempMethod.call(this, options)
      }
    }
  })
  originPage.call(this, page)
}

