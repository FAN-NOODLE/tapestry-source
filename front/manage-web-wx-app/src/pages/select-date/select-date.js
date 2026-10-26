import { getTaskCalendarApi } from '@/api/date-api'

const { globalData } = getApp()


// =====================
// type==1 转发顾客 type==2 转发朋友圈
// =====================
// mediaType==1 单图  mediaType==2 多图  mediaType==3 视频

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    week: ['日', '一', '二', '三', '四', '五', '六'],
    dateList: [],
    currentLocaldate: '', // 当前日期
    currentSelectdate: '', // 当前选中时间戳
    shouldFiexdSelect: false, // 浮动选框
    scrollTop: '', // 卷动值
    isFlash: false // 当前日期是否需要flash闪动
  },

  onReady () {
  },
  // 生命周期函数--监听页面加载
  onLoad (options) {
    this.setData({
      currentSelectdate: Number(options.curDate)
    })
  },
  onShow () {
    this.getTaskCalendartFn()
  },
  // 页面滚动事件
  onPageScroll (e) {
    this.setData({
      scrollTop: e.scrollTop
    })
    let query = wx.createSelectorQuery()

    query.select('#day-word').boundingClientRect((rect) => {
      let top = rect.bottom
      if (top > 620) { // 临界值，根据自己的需求来调整
        this.setData({
          shouldFiexdSelect: true // 是否显示返回今日
        })
      } else {
        this.setData({
          shouldFiexdSelect: false
        })
      }
    }).exec()
  },
  // 定位选中
  goBackToday (type) {
    var query = wx.createSelectorQuery().in(this)
    query.selectViewport().scrollOffset()
    if (type && type === 'scroll') {
      query.select('#day-active').boundingClientRect()
    } else {
      query.select('#day-word').boundingClientRect()
    }

    query.exec(function (res) {
      var miss = res[0].scrollTop + res[1].top - 300
      wx.pageScrollTo({
        scrollTop: miss,
        duration: 300
      })
    })
    if (type && type === 'scroll') {
      console.log('scroll')
    } else {
      this.setData({
        isFlash: true
      })
      setTimeout(() => {
        this.setData({
          isFlash: false
        })
      }, 1400)
    }
  },
  // 获取营销计划日历
  getTaskCalendartFn () {
    getTaskCalendarApi({
      current_date: this.data.currentSelectdate
    })
      .then((res) => {
        if (res.code === 200) {
          this.setData({
            dateList: res.data
          })
          this.getFormatDate()// 获取当前日期
          // 定位到选中的日期
          this.goBackToday('scroll')
        } else {
        }
      })
      .catch((err) => {
        console.log(err)
      })
  },
  // 选择日期
  selectDateFn (e) {
    let {item} = e.currentTarget.dataset
    var selectTimestamp = new Date(`${item.year}/${item.month}/${item.day}`).getTime() / 1000
    globalData.selectTimeStamp = selectTimestamp
    globalData.defaultDate = selectTimestamp
    globalData.isRefreshDate = true
    wx.switchTab({
      url: '/pages/compass-data/compass-data'
    })
  },
  // 获取当前年月日
  getFormatDate () {
    var date = new Date()
    var year = date.getFullYear().toString()
    var month = (date.getMonth() + 1).toString()
    var strDate = date.getDate().toString()
    this.setData({
      currentLocaldate: year + month + strDate
    })
  }
})
