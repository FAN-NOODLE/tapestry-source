import { PolygonShape, CircleShape, CircleOneShape, CircleTwoShape, CircleOneLineShape } from '@/utils/chart'

Component({
  properties: {
    chartArr: {
      type: Array,
      value: []
    }
  },
  data: {
    tapIdx: -1, // 说明弹窗index
    isTapPopover: false, // 说明弹窗是否打开
    isShowChart: false, // 是否展示第二个canvas
    size: 236, // 5: 178 4: 162 3: 118
    style: 1
  },

  lifetimes: {
    ready: function () {
      // 通过 SelectorQuery 获取 Canvas 节点
      console.log('ready')
      const query = this.createSelectorQuery().in(this)

      query.select('#chart')
        .fields({
          size: true,
          node: true
        })
        .exec(this.init.bind(this))

      let size = 0

      switch (this.properties.chartArr.length) {
        case 5:
          size = 356
          break
        case 4:
          size = 324
          break
        case 3:
        case 1:
          size = 236
          break
        case 2:
          size = 204
          break
        default:
          break
      }

      this.setData({ size })
    }
  },
  methods: {
    init (res) {
      console.log('init', res)
      const canvas = res[0].node
      const width = res[0].width
      const height = res[0].height
      const { chartArr } = this.properties
      let chart = null

      console.log('chartArr', chartArr)

      if (chartArr.length === 5 || chartArr.length === 4) {
        this.setData({ style: 1 })
        chart = new PolygonShape(canvas, chartArr.length, width, height)
      } else if (chartArr.length === 3) {
        this.setData({ style: 2 })
        chart = new CircleShape(canvas, chartArr.length, width, height)
      } else if (chartArr.length === 2) {
        if (chartArr.some(i => i.inviteRule.length < 2)) {
          // 展示两个canvas
          this.setData({ isShowChart: true, style: 3 }, () => {
            const query = this.createSelectorQuery().in(this)

            query.select('#chart2')
              .fields({
                size: true,
                node: true
              })
              .exec(res => {
                this.chart2 = new CircleOneShape(res[0].node, 0, width, height)
              })
          })
          chart = new CircleOneShape(canvas, 0, width, height)
        } else {
          this.setData({ style: 4 })
          chart = new CircleTwoShape(canvas, 0, width, height)
        }
      } else if (chartArr.length === 1) {
        this.setData({ style: 5 })
        chart = new CircleOneLineShape(canvas, 0, width, height)
      }

      this.chart = chart
    },
    draw () {
      const levelArr = this.properties.chartArr.map(i => i.index)

      const { chartArr } = this.properties
      if (chartArr.length === 2 && chartArr.some(i => i.inviteRule.length < 2)) {
        this.chart.draw(chartArr[0].value)
        this.chart2.draw(chartArr[1].value)
      } else {
        this.chart.draw(levelArr)
      }

      const classStr = [...new Array(levelArr.length).keys()].map(i => `.index-${i + 1}`).join(',')

      this.animate(classStr, [
        { opacity: 0 },
        { opacity: 1 }
      ], 500, () => {
        this.triggerEvent('animateEnd')
      })
    },
    showPopover (e) {
      const { tapIdx } = this.data
      const { index } = e.currentTarget.dataset

      this.setData({
        tapIdx: e.currentTarget.dataset.index,
        isTapPopover: index !== tapIdx
      })
    },
    tapChart () {
      this.setData({ tapIdx: -1 })
      this.setData({ isTapPopover: false })
    }
  }
})
