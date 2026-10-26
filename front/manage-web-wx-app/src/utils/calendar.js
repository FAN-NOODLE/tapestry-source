/** 日期 */
export const WEEK_NAME_MAP = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']

class Calendar {
 data = {
   year: 0,
   month: 0,
   day: 0,
   currentWeek: 0,
   // 本月第一天时间戳
   firstDayNum: 0,
   // 本月最后一天时间戳
   lastDayNum: 0,
   // 当前时间戳
   dayToNum: 0
 };

  /** 计算后的日期 */
  dayList = [];

  /** 给时间属性添加特殊处理，添加参数。。。。。。 */
  businessFn = null

  /** 构造函数 */
  constructor (businessFn, year, month) {
    this.businessFn = businessFn
    this.data.year = year
    this.data.month = month
    this.handleDateChange({ year, month })
  }

  /** 内部函数-计算当前日期 */
  calcDayList () {
    const { year, month } = this.data

    const totalDay = this.getThisMonthDays(year, month)
    const firstDay = this.getFirstDayOfWeek(year, month)

    const dayList = []
    // 初始化数据
    for (let i = 0; i < totalDay; i++) {
      const dayToString = `${year}/${month}/${i + 1}`
      const dayToNum = ~~(new Date(dayToString).getTime() / 1000)
      const todayNum = new Date().setHours(0, 0, 0, 0) / 1000
      dayList.push({
        day: i + 1,
        dayToString: dayToString,
        // 时间戳
        dayToNum: dayToNum,
        // 本周name
        weekName: WEEK_NAME_MAP[(firstDay + i) % 7],
        // 是否今天
        isToday: false,
        // 月
        month,
        // 是否小于今天
        disabled: todayNum < dayToNum
      })
    }

    // 补充当月第一天和最后一天的时间戳
    this.data.firstDayNum = dayList[0].dayToNum
    this.data.lastDayNum = dayList[dayList.length - 1].dayToNum

    // 向前添加
    for (let i = 0; i < firstDay; i++) {
      dayList.unshift({ day: '', isToday: false })
    }

    // 初始化默认值
    this.data.day = 0
    // 当天
    const now = new Date()
    if (now.getFullYear() === year && now.getMonth() + 1 === month) {
      // 寻找今天
      const currentDay = dayList.find((d) => d.day === now.getDate())
      if (currentDay) {
        currentDay.isToday = true
        this.data.day = currentDay.day
        this.data.dayToNum = currentDay.dayToNum
      }
    }

    // 初始化默认值
    this.data.currentWeek = 0

    // 以7为周切分日历，判断是否一周，currentWeek 为当前元素添加是属于那一周
    dayList.forEach((item, index) => {
      item.currentWeek = Math.ceil((index + 1) / 7)
      // 保存今天视第几周
      if (item.isToday) {
        this.data.currentWeek = item.currentWeek
      }
    })

    // 赋值
    this.dayList = this.businessFn(dayList, this.data)
  }

  /** 内部函数-处理日期变化 */
  handleDateChange (time) {
    const date = { ...this.data, ...time }
    if (date) {
      // 当date 有值时, 保存当前值
      this.data = date
      // 如果年份、月份变化，重新计算日期列表
      if (time.year || time.month) this.calcDayList()
    }
  }

  /** 内部函数-获取当月共多少天 */
  getThisMonthDays (year, month) {
    return new Date(year, month, 0).getDate()
  }

  /** 内部函数-获取当月第一天星期几   0 星期日 */
  getFirstDayOfWeek (year, month) {
    return new Date(Date.UTC(year, month - 1, 1)).getDay()
  }

  /** 点击上个月 */
  handlePrevMonthClick () {
    console.log('点击了上个月', this)
    let year = this.data.year
    let month = this.data.month
    if (this.data.month === 1) {
      year = this.data.year - 1
      month = 12
    } else {
      month = this.data.month - 1
    }
    this.handleDateChange({ year, month })
  }

  /** 点击下个月 */
  handleNextMonthClick () {
    console.log('点击了下个月')
    let year = this.data.year
    let month = this.data.month
    if (this.data.month === 12) {
      year = this.data.year + 1
      month = 1
    } else {
      month = this.data.month + 1
    }
    this.handleDateChange({ year, month })
  }

  /** 指定月份查询 */
  appointYearOrMonth (year, month) {
    this.handleDateChange({ year, month })
  }
}

export default Calendar
