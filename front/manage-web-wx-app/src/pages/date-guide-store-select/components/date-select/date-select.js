import { toDate } from '@/utils/tools'
import dayjs from 'dayjs'

const { globalData } = getApp()

const getLastYear = () => {
  let date = new Date()
  date.setHours(0, 0, 0, 0)
  return date.setFullYear(date.getFullYear() - 1) / 1000
}

Component({
  properties: {
    startTime: {
      type: Number,
      value: new Date(new Date().setHours(0, 0, 0, 0)) / 1000
    },
    endTime: {
      type: Number,
      value: 0
    },
    max: {
      type: Number,
      value: new Date().setHours(23, 59, 59, 0) / 1000
    },
    min: {
      type: Number,
      value: getLastYear()
    },
    showDate: {
      // 选项是否具体到日
      type: Boolean,
      value: true
    },
    selectDateText: {
      type: String,
      value: 'selectDate'
    }
  },
  data: {
    /*  私有数据，可用于模版渲染 */
    startTime: '',
    endTime: '',
    currentClickTime: 1,
    isIpx: globalData.isIpx,
    displayMonths: [0, 11],
    displayDays: [0, 30]
  },
  /*    系统事件开始   */
  created () {},
  attached () {},
  ready () {
    if (this.data.showDate) {
      this.setData({
        startTimeText: this.properties.startTime
          ? toDate(this.properties.startTime, 'yyyy-MM-dd').nowTime
          : toDate(new Date().getTime(), 'yyyy-MM-dd').nowTime,
        endTimeText: this.properties.endTime
          ? toDate(this.properties.endTime, 'yyyy-MM-dd').nowTime
          : ''
      })
    } else {
      this.setData({
        startTimeText: this.properties.startTime
          ? toDate(this.properties.startTime, 'yyyy-MM').nowTime
          : toDate(new Date().getTime(), 'yyyy-MM').nowTime,
        endTimeText: ''
      })
    }
    this.initDateOption(this.properties.startTime)
    this.initPick(this.properties.startTime)
  },
  /*    系统事件结束   */

  methods: {
    /*    普通事件开始  */
    // 根据选择时间获取当前月的最大天数
    getDayMaxNumWithMonth (year, month) {
      const dayAll = [31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31]
      if ((year % 4 === 0 && year % 100 !== 0) || year % 400 === 0) {
        dayAll[1] = 29
      }
      return dayAll[month - 1]
    },
    // 计算年份的可选项
    getYearOptions (min, max) {
      let years = []
      for (let i = min.getFullYear(); i <= max.getFullYear(); i++) {
        years.push(i)
      }
      return years
    },
    // 计算月份的可选项
    getMonthOptions () {
      let start = 1
      let end = 12
      let months = []
      for (let i = start; i <= end; i++) {
        months.push(i)
      }
      return months
    },
    // 计算日期的可选项
    getDayOptions () {
      let start = 1
      let end = 31
      let days = []
      for (let i = start; i <= end; i++) {
        days.push(i)
      }
      return days
    },
    // 补零
    formatDate (value) {
      return value < 10 ? `0${value}` : value
    },
    // 初始化时间筛选框
    initDateOption (time) {
      const current = new Date(time * 1000) // 传入的时间
      const min = new Date(this.properties.min * 1000) // 最小时间
      const max = new Date(this.properties.max * 1000) // 最大时间
      const years = this.getYearOptions(min, max)
      const months = this.getMonthOptions()
      const days = this.getDayOptions()
      this.setData({
        years,
        months,
        days,
        showPicker: true
      })
      setTimeout(() => {
        this.setData({
          displayDays: this.getDispalyDays(min, max, current),
          displayMonths: this.getDispalyMonths(min, max, current)
        })
      }, 30)
    },
    initPick (time) {
      const current = new Date(time * 1000) // 传入的时
      const { years, months, days } = this.data
      let yearIndex = years.indexOf(current.getFullYear())
      let monthIndex = months.indexOf(current.getMonth() + 1)
      let dayIndex = days.indexOf(current.getDate())
      this.setData({
        pickerValue: [yearIndex, monthIndex, dayIndex]
      })
    },
    // 获取要显示的月份索引范围
    getDispalyMonths (min, max, current) {
      let displayMonths = [0, 11]
      if (min.getFullYear() === current.getFullYear()) {
        displayMonths[0] = min.getMonth()
      }
      if (max.getFullYear() === current.getFullYear()) {
        displayMonths[1] = max.getMonth()
      }
      return displayMonths
    },
    // 获取要显示的月份索引范围
    getDispalyDays (min, max, current) {
      let displayDays = [
        0,
        this.getDayMaxNumWithMonth(
          current.getFullYear(),
          current.getMonth() + 1
        ) - 1
      ]
      if (
        min.getFullYear() === current.getFullYear() &&
        min.getMonth() === current.getMonth()
      ) {
        displayDays[0] = min.getDate() - 1
      }
      if (
        max.getFullYear() === current.getFullYear() &&
        max.getMonth() === current.getMonth()
      ) {
        displayDays[1] = max.getDate() - 1
      }
      return displayDays
    },
    // 日期筛选框修改赋值
    changeValue (e) {
      let [yearIndex, monthIndex, dayIndex] = e.detail.value
      const { years, months, days, currentClickTime } = this.data
      console.log(months, monthIndex)
      const min = new Date(this.properties.min * 1000) // 最小时间
      const max = new Date(this.properties.max * 1000) // 最大时间
      let year = years[yearIndex]
      let month = months[monthIndex]
      let day = days[dayIndex]
      const maxDay = this.getDayMaxNumWithMonth(year, month)
      if (day > maxDay) {
        day = maxDay
      }
      // 赋新值给当前要修改的时间 先回滚月份
      let currentTimeText = `${year}-${this.formatDate(
        month
      )}-${this.formatDate(day)}`
      let currentTimeStamp = parseInt(
        new Date(currentTimeText).setHours(0, 0, 0, 0) / 1000
      )
      if (currentTimeStamp < this.properties.min) {
        month = min.getMonth() + 1
      }
      if (currentTimeStamp > this.properties.max) {
        month = max.getMonth() + 1
      }
      // 然后回滚日期
      currentTimeText = `${year}-${this.formatDate(month)}-${this.formatDate(
        day
      )}`
      currentTimeStamp = parseInt(
        new Date(currentTimeText).setHours(0, 0, 0, 0) / 1000
      )
      if (currentTimeStamp < this.properties.min) {
        day = min.getDate()
      }
      if (currentTimeStamp > this.properties.max) {
        day = max.getDate()
      }
      monthIndex = months.indexOf(month)
      dayIndex = days.indexOf(day)
      e.detail.value[1] = monthIndex
      e.detail.value[2] = dayIndex
      this.setData({
        pickerValue: [yearIndex, monthIndex, dayIndex]
      })

      // 选择到月份时 日期的显示
      const endDay =
        month === new Date().getMonth() + 1 ? new Date().getDate() : maxDay
      const newTimeText =
        !this.data.showDate && currentClickTime === 1
          ? `${year}-${this.formatDate(month)}`
          : !this.data.showDate && currentClickTime === 2
            ? `${year}-${this.formatDate(month)}-${endDay}`
            : `${year}-${this.formatDate(month)}-${this.formatDate(day)}`

      const newTime = dayjs(newTimeText).toDate()
      const newTimeStamp = parseInt(newTime.getTime() / 1000)
      console.log('newTimeText', newTimeText)
      console.log('newTime', newTime)
      console.log('newTimeStamp', newTimeStamp)

      if (currentClickTime === 1) {
        this.setData({
          startTimeText: newTimeText,
          startTime: newTimeStamp
        })
      } else {
        this.setData({
          endTimeText: newTimeText,
          endTime: newTimeStamp + 86399
        })
      }
      setTimeout(() => {
        this.setData({
          displayMonths: this.getDispalyMonths(min, max, newTime),
          displayDays: this.getDispalyDays(min, max, newTime)
        })
      }, 30)
    },
    // 切换日期
    changeTimeClick (e) {
      this.setData({
        currentClickTime: parseInt(e.currentTarget.dataset.type),
        showPicker: false
      })
      const time =
        this.data.currentClickTime === 1
          ? this.data.startTime
          : this.data.endTime
            ? this.data.endTime
            : new Date(new Date().setHours(23, 59, 59, 0)) / 1000
      this.initDateOption(time)
      this.initPick(time)
      if (this.data.currentClickTime === 2) {
        this.setData({
          endTimeText: toDate(time, 'yyyy-MM-dd').nowTime,
          endTime: time
        })
      }
    },
    // 重置按钮
    reset () {
      this.setData({
        currentClickTime: 1,
        showPicker: false,
        startTimeText: '',
        endTimeText: '',
        startTime: new Date(new Date().setHours(0, 0, 0, 0)) / 1000,
        endTime: new Date(new Date().setHours(23, 59, 59, 0)) / 1000
      })
      this.initDateOption(this.data.startTime)
      this.initPick(this.data.startTime)
    },
    // 确认按钮
    sure () {
      if (
        this.data.startTimeText &&
        this.data.endTimeText &&
        this.data.startTime > this.data.endTime
      ) {
        wx.showToast('结束时间不可早于开始时间')
        return
      }
      if (this.data.startTimeText && this.data.endTimeText) {
        globalData[this.data.selectDateText] = [
          {
            time: toDate(this.data.startTime, 'yyyy.MM.dd').nowTime,
            s_timestamp: this.data.startTime
          },
          {
            time: toDate(this.data.endTime, 'yyyy.MM.dd').nowTime,
            s_timestamp: this.data.endTime
          }
        ]
      } else if (this.data.startTimeText && !this.data.endTimeText) {
        // wx.showToast('请选择截至时间')
        // return
        // 如果仅“结束时间”未填入，则认为截止到当前日期
        if (this.data.showDate === false) {
          globalData[this.data.selectDateText] = [
            {
              time: toDate(this.data.startTime, 'yyyy-MM').nowTime,
              s_timestamp: this.data.startTime
            },
            {
              time: toDate(this.data.startTime, 'yyyy-MM').nowTime,
              s_timestamp: this.data.startTime
            }
          ]
        } else {
          wx.showToast('请选择截至时间')
          return
          // globalData.selectDate = [
          //   {
          //     time: toDate(this.data.startTime, 'yyyy.MM.dd').nowTime,
          //     s_timestamp: this.data.startTime
          //   },
          //   {
          //     time: toDate(new Date(new Date().setHours(23, 59, 59, 0)) / 1000, 'yyyy.MM.dd').nowTime,
          //     s_timestamp: new Date(new Date().setHours(23, 59, 59, 0)) / 1000
          //   }
          // ]
        }
      } else if (!this.data.startTimeText && this.data.endTimeText) {
        // 如果仅“开始时间”未填入，则认为开始时间与截止时间相同
        globalData[this.data.selectDateText] = [
          {
            time: toDate(
              this.data.endTime - ((this.data.endTime + 8 * 3600) % 86400),
              'yyyy.MM.dd'
            ).nowTime,
            s_timestamp:
              this.data.endTime - ((this.data.endTime + 8 * 3600) % 86400)
          },
          {
            time: toDate(this.data.endTime, 'yyyy.MM.dd').nowTime,
            s_timestamp: this.data.endTime
          }
        ]
      } else {
        // 如果“开始时间”、“结束时间”未填入，则点击跳页后页面不修改，与进入页面之前保持一致
        globalData[this.data.selectDateText] = []
      }
      console.log(
        '==============date',
        this.data.showDate,
        globalData[this.data.selectDateText]
      )
      wx.navigateBack({
        delta: 1
      })
    },
    /*    普通事件结束  */

    /*    cdp分割事件   */
    cdpSplitEvent () {}
    /*    cdp分割事件   */

    /*    行为事件开始   */

    /*    行为事件结束   */
  }
})
