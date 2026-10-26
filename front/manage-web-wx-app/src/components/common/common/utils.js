export function requestAnimationFrame(cb) {
  setTimeout(() => {
    cb && cb()
  }, 1000 / 30)
}
export function getRect(context, selector) {
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

export function getAllRect(context, selector) {
  return new Promise(function (resolve) {
    wx.createSelectorQuery()
      .in(context)
      .selectAll(selector)
      .boundingClientRect()
      .exec(function (rect) {
        if (rect === void 0) { rect = [] }
        return resolve(rect[0])
      })
  })
}

export function getCurrentPage() {
  let currentPages = getCurrentPages()
  let len = currentPages.length - 1
  return currentPages[len]
}
