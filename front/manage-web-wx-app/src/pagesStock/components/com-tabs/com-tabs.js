/* eslint-disable promise/catch-or-return */
// Component Object
function getRect(context, selector) {
  return new Promise(function (resolve) {
    wx.createSelectorQuery()
      .in(context)
      .select(selector)
      .boundingClientRect()
      .exec(function (rect) {
        if (rect === void 0) { rect = [] }
        return resolve(rect[0])
      })
  })
}
function getAllRect(context, selector) {
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
Component({
  externalClasses: ['custom-class'],
  options: {
    virtualHost: false,
    multipleSlots: true,
    styleIsolation: 'apply-shared'
  },
  behaviors: [],
  properties: {
    customStyle: {
      type: String
    },
    active: {
      optionalTypes: [String, Number],
      observer: function (index) {
        this.setCurrentIndex(index)
      }
    },
    tabs: {
      type: Array,
      value: []
    },
    lineWidth: {
      type: Number,
      value: 60,
      observer: 'resize'
    }
  },
  data: {
    title: 'mp-tabs',
    scrollable: true,
    currentIndex: 0,
    lineOffsetLeft: 0,
    scrollLeft: 0
  },
  methods: {
    /**
     * @description: 切换tab
     */
    onTap(e) {
      const { index } = e.currentTarget.dataset
      this.triggerEvent('change', { index })
    },
    setCurrentIndex(currentIndex) {
      this.setData({ currentIndex })
      this.resize()
      this.scrollIntoView()
    },
    resize() {
      const { currentIndex } = this.data
      // eslint-disable-next-line promise/catch-or-return
      Promise.all([
        getAllRect(this, '.com-tab'),
        getRect(this, '.com-tabs__line')
      ]).then(([rects, lineRect]) => {
        const rect = rects[currentIndex]

        let lineOffsetLeft = rects
          .slice(0, currentIndex)
          .reduce(function (prev, curr) { return prev + curr.width }, 0)
        lineOffsetLeft += (rect.width - lineRect.width) / 2
        this.setData({
          lineOffsetLeft: lineOffsetLeft
        })
      })
    },
    scrollIntoView() {
      const { currentIndex } = this.data
      Promise.all([
        getAllRect(this, '.com-tab'),
        getRect(this, '.com-tabs__nav')
      ]).then(([rects, navRect]) => {
        var rect = rects[currentIndex]
        var offsetLeft = rects
          .slice(0, currentIndex)
          .reduce(function (prev, curr) { return prev + curr.width }, 0)

        this.setData({
          scrollLeft: offsetLeft - (navRect.width - rect.width) / 2
        })
      })
    }
  }
})

