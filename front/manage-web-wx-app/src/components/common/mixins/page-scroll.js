function getCurrentPage() {
  let currentPages = getCurrentPages()
  let len = currentPages.length - 1
  return currentPages[len]
}

function onPageScroll(event) {
  const pageScroller = getCurrentPage().pageScroller || []
  pageScroller.forEach(function (scroller) {
    if (typeof scroller === 'function') {
      // @ts-ignore
      scroller(event)
    }
  })
}
export function pageScrollMixin(scroller) {
  return Behavior({
    attached: function () {
      var page = getCurrentPage()
      if (!page) {
        return
      }
      var _scroller = scroller.bind(this)
      const pageScroller = page.pageScroller || []
      if (
        typeof page.onPageScroll === 'function' &&
        page.onPageScroll !== onPageScroll
      ) {
        pageScroller.push(page.onPageScroll.bind(page))
      }
      pageScroller.push(_scroller)
      page.pageScroller = pageScroller
      page.onPageScroll = onPageScroll
      this._scroller = _scroller
    },
    detached: function () {
      var _this = this
      var page = getCurrentPage()
      if (!page || !page.pageScroller) {
        return
      }
      var pageScroller = page.pageScroller
      var index = pageScroller.findIndex(function (v) {
        return v === _this._scroller
      })
      if (index > -1) {
        page.pageScroller.splice(index, 1)
      }
      this._scroller = undefined
    }
  })
}
