// Component Object
Component({
  properties: {
    width: {
      type: String,
      value: ''
    },
    height: {
      type: String,
      value: ''
    },
    title: {
      type: String,
      value: ''
    },
    legendRight: {
      type: String,
      value: '0'
    },
    legendWidth: {
      type: String,
      value: '210rpx'
    }

  },
  data: {
    options: [],
    showTips: false
  },

  lifetimes: {
    attached: function () {
    },
    detached: function () {

    }
  },
  methods: {
    /** 初始化统计图 */
    init (data) {
      const options = {
        silent: true,
        legend: [
          {
            orient: 'horizontal',
            bottom: 24,
            icon: 'rect',
            itemGap: 15,
            x: 'center',
            itemWidth: 10,
            itemHeight: 10,
            itemStyle: {
              opacity: 1,
              fontSize: 12
            },
            data: ['黑钻卡', '钻卡', '白金卡', '金卡']
          },
          {
            orient: 'horizontal',
            bottom: 0,
            icon: 'rect',
            itemGap: 15,
            itemWidth: 10,
            itemHeight: 10,
            itemStyle: {
              opacity: 1,
              fontSize: 12
            },
            data: ['名媛', 'GC', '会员', '非会员']
          }
        ],
        color: [
          '#E8CE9F',
          '#F2B183',
          '#F28383',
          '#E89FC2',
          '#8FDBBF',
          '#75E0FF',
          '#98A3D9'
        ],
        series: [
          {
            type: 'pie',
            radius: ['35%', '65%'],
            center: ['50%', '50%'],
            data,
            itemStyle: {
              opacity: 1,
              borderWidth: 1,
              borderColor: '#fff'
            },
            labelLine: {
              smooth: 0.2,
              length: 2,
              length2: 4,
              lineStyle: {
                // type: ''
              }
            },
            label: {
              opacity: 1,
              show: true,
              align: 'left',
              fontSize: 9,
              color: '#333333',
              formatter: function (params) {
                return `${params.name}`
              },
              rich: {
                x: {
                  color: '#333333',
                  fontSize: 9
                }
              }
            }
          },
          {
            type: 'pie',
            radius: '35%',
            center: ['50%', '50%'],
            name: this.data.title,
            emptyCircleStyle: {
              opacity: 0

            }
          }
        ]
      }
      this.setData({
        options,
        data
      })
    },
    handleShowTips() {
      this.setData({
        showTips: true
      })
    },
    handleHideTips() {
      this.setData({
        showTips: false
      })
    }
  }
})
