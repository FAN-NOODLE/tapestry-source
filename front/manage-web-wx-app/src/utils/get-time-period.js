
/***
    * 获得当前时间
    *
    */
// eslint-disable-next-line no-unused-vars
const getCurrentDate = function (type = 'yyyy-MM-dd') {
  let date = new Date()
  const param = [{
    ms_timestamp: new Date().setHours(0, 0, 0, 0),
    s_timestamp: parseInt(new Date().setHours(0, 0, 0, 0) / 1000),
    time: date.format(type)
  }, {
    ms_timestamp: new Date().getTime(),
    s_timestamp: parseInt(new Date().getTime() / 1000),
    time: date.format(type)
  }]
  return param
}
/***
    * 获得本周起止时间
    */
// eslint-disable-next-line no-unused-vars
const getCurrentWeek = function (type = 'yyyy-MM-dd') {
  // 起止日期数组
  // eslint-disable-next-line no-array-constructor
  var startStop = new Array()
  // 获取当前时间
  var currentDate = new Date()
  // 返回date是一周中的某一天
  var week = currentDate.getDay()
  // 一天的毫秒数
  var millisecond = 1000 * 60 * 60 * 24
  // 减去的天数
  // eslint-disable-next-line eqeqeq
  var minusDay = week != 0 ? week - 1 : 6
  // alert(minusDay);
  // 本周 周一
  var monday = new Date(currentDate.getTime() - (minusDay * millisecond))
  // 本周 周日
  var sunday = new Date(monday.getTime() + (6 * millisecond))

  // 添加本周时间
  startStop.push({
    ms_timestamp: monday.setHours(0, 0, 0, 0),
    s_timestamp: parseInt(monday.setHours(0, 0, 0, 0) / 1000),
    time: monday.format(type)
  }) // 本周起始时间
  // 添加本周最后一天时间
  startStop.push({
    ms_timestamp: sunday.getTime(),
    s_timestamp: parseInt(sunday.getTime() / 1000),
    time: sunday.format(type)
  }) // 本周终止时间
  // 返回
  return startStop
}
/***
    * 获得本月的起止时间
    */
// eslint-disable-next-line no-unused-vars
const getCurrentMonth = function (type = 'yyyy-MM-dd') {
  // 起止日期数组
  // eslint-disable-next-line no-array-constructor
  var startStop = new Array()
  // 获取当前时间
  var currentDate = new Date()
  // 获得当前月份0-11
  var currentMonth = currentDate.getMonth()
  // 获得当前年份4位年
  var currentYear = currentDate.getFullYear()
  // 求出本月第一天
  var firstDay = new Date(currentYear, currentMonth, 1)


  // 当为12月的时候年份需要加1
  // 月份需要更新为0 也就是下一年的第一个月
  // eslint-disable-next-line eqeqeq
  if (currentMonth == 11) {
    currentYear++
    currentMonth = 0 // 就为
  } else {
    // 否则只是月份增加,以便求的下一月的第一天
    currentMonth++
  }


  // 一天的毫秒数
  var millisecond = 1000 * 60 * 60 * 24
  // 下月的第一天
  var nextMonthDayOne = new Date(currentYear, currentMonth, 1)
  // 求出上月的最后一天
  var lastDay = new Date(nextMonthDayOne.getTime() - millisecond)

  // 添加至数组中返回
  startStop.push({
    ms_timestamp: firstDay.getTime(),
    s_timestamp: parseInt(firstDay.getTime() / 1000),
    time: firstDay.format(type)
  })
  startStop.push({
    ms_timestamp: lastDay.getTime(),
    s_timestamp: parseInt(lastDay.getTime() / 1000),
    time: lastDay.format(type)
  })
  // 返回
  return startStop
}
/***
    * 获得下月的起止时间
    */
const getNextMonth = function (type = 'yyyy-MM-dd') {
  // 起止日期数组
  // eslint-disable-next-line no-array-constructor
  var startStop = new Array()
  // 获取当前时间
  var currentDate = new Date()
  // 获得当前月份0-11
  var currentMonth = currentDate.getMonth()
  // 获得当前年份4位年
  var currentYear = currentDate.getFullYear()
  // 获取下个月所在的月份、年份
  if (currentMonth < 11) {
    var nextMonth = currentMonth + 1
    var nextYear = currentYear
  } else {
    // eslint-disable-next-line no-redeclare
    var nextMonth = 0
    // eslint-disable-next-line no-redeclare
    var nextYear = currentYear + 1
  }
  // 求出下个月第一天
  var firstDay = new Date(nextYear, nextMonth, 1)
  var lastDay = new Date(nextYear, nextMonth + 1, 0)
  // 添加至数组中返回
  startStop.push({
    ms_timestamp: firstDay.getTime(),
    s_timestamp: parseInt(firstDay.getTime() / 1000),
    time: firstDay.format(type)
  })
  startStop.push({
    ms_timestamp: lastDay.getTime(),
    s_timestamp: parseInt(lastDay.getTime() / 1000),
    time: lastDay.format(type)
  })
  // 返回
  return startStop
}
/**
    * 获得本季度的起止日期
    */
// eslint-disable-next-line no-unused-vars
const getCurrentSeason = function (type = 'yyyy-MM-dd') {
  // 起止日期数组
  // eslint-disable-next-line no-array-constructor
  var startStop = new Array()
  // 获取当前时间
  var currentDate = new Date()
  // 获得当前月份0-11
  var currentMonth = currentDate.getMonth()
  // 获得当前年份4位年
  var currentYear = currentDate.getFullYear()
  // 获得本季度开始月份
  var quarterSeasonStartMonth = this.getQuarterSeasonStartMonth(currentMonth)
  // 获得本季度结束月份
  var quarterSeasonEndMonth = quarterSeasonStartMonth + 2

  // 获得本季度开始的日期
  var quarterSeasonStartDate = new Date(currentYear, quarterSeasonStartMonth, 1)
  // 获得本季度结束的日期
  var quarterSeasonEndDate = new Date(currentYear, quarterSeasonEndMonth, this.getMonthDays(currentYear, quarterSeasonEndMonth))
  // 加入数组返回
  startStop.push({
    ms_timestamp: quarterSeasonStartDate.getTime(),
    s_timestamp: parseInt(quarterSeasonStartDate.getTime() / 1000),
    time: quarterSeasonStartDate.format(type)
  })
  startStop.push({
    ms_timestamp: quarterSeasonEndDate.getTime(),
    s_timestamp: parseInt(quarterSeasonEndDate.getTime() / 1000),
    time: quarterSeasonEndDate.format(type)
  })
  // 返回
  return startStop
}

/***
    * 得到本年的起止日期
    *
    */
// eslint-disable-next-line no-unused-vars
const getCurrentYear = function (type = 'yyyy-MM-dd') {
  // 起止日期数组
  // eslint-disable-next-line no-array-constructor
  var startStop = new Array()
  // 获取当前时间
  var currentDate = new Date()
  // 获得当前年份4位年
  var currentYear = currentDate.getFullYear()

  // 本年第一天
  var currentYearFirstDate = new Date(currentYear, 0, 1)
  // 本年最后一天
  var currentYearLastDate = new Date(currentYear, 11, 31)
  // 添加至数组
  startStop.push({
    ms_timestamp: currentYearFirstDate.getTime(),
    s_timestamp: parseInt(currentYearFirstDate.getTime() / 1000),
    time: currentYearFirstDate.format(type)
  })
  startStop.push({
    ms_timestamp: currentYearLastDate.getTime(),
    s_timestamp: parseInt(currentYearLastDate.getTime() / 1000),
    time: currentYearLastDate.format(type)
  })
  // 返回
  return startStop
}
/**
    * 得到本季度开始的月份
    * @param month 需要计算的月份
    ***/
// eslint-disable-next-line no-unused-vars
const getQuarterSeasonStartMonth = function (month) {
  var spring = 0 // 春
  var summer = 3 // 夏
  var fall = 6 // 秋
  var winter = 9 // 冬
  // 月份从0-11
  if (month < 3) {
    return spring
  }

  if (month < 6) {
    return summer
  }

  if (month < 9) {
    return fall
  }

  return winter
}

/**
    * 获得该月的天数
    * @param year年份
    * @param month月份
    * */
// eslint-disable-next-line no-unused-vars
const getMonthDays = function (year, month) {
  // 本月第一天 1-31
  var relativeDate = new Date(year, month, 1)
  // 获得当前月份0-11
  var relativeMonth = relativeDate.getMonth()
  // 获得当前年份4位年
  var relativeYear = relativeDate.getFullYear()

  // 当为12月的时候年份需要加1
  // 月份需要更新为0 也就是下一年的第一个月
  // eslint-disable-next-line eqeqeq
  if (relativeMonth == 11) {
    relativeYear++
    relativeMonth = 0
  } else {
    // 否则只是月份增加,以便求的下一月的第一天
    relativeMonth++
  }
  // 一天的毫秒数
  var millisecond = 1000 * 60 * 60 * 24
  // 下月的第一天
  var nextMonthDayOne = new Date(relativeYear, relativeMonth, 1)
  // 返回得到上月的最后一天,也就是本月总天数
  return new Date(nextMonthDayOne.getTime() - millisecond).getDate()
}

/***
    * 获得昨天时间
    */
// eslint-disable-next-line no-unused-vars
const getYesterDays = function (type = 'yyyy-MM-dd') {
  // 起止日期数组
  let date = new Date()

  let millisecond = 1000 * 60 * 60 * 24
  // 近七天的开始时间
  let yesterday = new Date(date.getTime() - (1 * millisecond))
  const param = [{
    ms_timestamp: yesterday.setHours(0, 0, 0, 0),
    s_timestamp: parseInt(yesterday.setHours(0, 0, 0, 0) / 1000),
    time: yesterday.format(type)
  }, {
    ms_timestamp: yesterday.setHours(23, 59, 59, 999),
    s_timestamp: parseInt(yesterday.setHours(23, 59, 59, 999) / 1000),
    time: yesterday.format(type)
  }]
  return param
}
/***
    * 获得近七天起止时间
    */
// eslint-disable-next-line no-unused-vars
const getSevenDays = function (type = 'yyyy-MM-dd') {
  // 起止日期数组
  // eslint-disable-next-line no-array-constructor
  let startStop = new Array()
  // 获取当前时间
  let currentDate = new Date()
  // 一天的毫秒数
  let millisecond = 1000 * 60 * 60 * 24
  // 近七天的开始时间
  let week = new Date(currentDate.getTime() - (6 * millisecond))

  // 添加近七天的开始时间
  startStop.push({
    ms_timestamp: week.setHours(0, 0, 0, 0),
    s_timestamp: parseInt(week.setHours(0, 0, 0, 0) / 1000),
    time: week.format(type)
  })
  // 添加当前时间
  startStop.push({
    ms_timestamp: currentDate.getTime(),
    s_timestamp: parseInt(currentDate.getTime() / 1000),
    time: currentDate.format(type)
  })
  // 返回
  return startStop
}
/***
    * 获得近30天起止时间
    */
// eslint-disable-next-line no-unused-vars
const getThirtyDays = function (type = 'yyyy-MM-dd') {
  // 起止日期数组
  // eslint-disable-next-line no-array-constructor
  let startStop = new Array()
  // 获取当前时间
  let currentDate = new Date()
  // 一天的毫秒数
  let millisecond = 1000 * 60 * 60 * 24
  // 近30天的开始时间
  let week = new Date(currentDate.getTime() - (29 * millisecond))

  // 添加近30天的开始时间
  startStop.push({
    ms_timestamp: week.setHours(0, 0, 0, 0),
    s_timestamp: parseInt(week.setHours(0, 0, 0, 0) / 1000),
    time: week.format(type)
  })
  // 添加当前时间
  startStop.push({
    ms_timestamp: currentDate.getTime(),
    s_timestamp: parseInt(currentDate.getTime() / 1000),
    time: currentDate.format(type)
  })
  // 返回
  return startStop
}

// eslint-disable-next-line no-extend-native
Date.prototype.format = function (fmt) {
  var o = {
    'M+': this.getMonth() + 1, // 月份
    'd+': this.getDate(), // 日
    'h+': this.getHours(), // 小时
    'm+': this.getMinutes(), // 分
    's+': this.getSeconds(), // 秒
    'q+': Math.floor((this.getMonth() + 3) / 3), // 季度
    'S': this.getMilliseconds() // 毫秒
  }
  if (/(y+)/.test(fmt)) {
    fmt = fmt.replace(RegExp.$1, (this.getFullYear() + '').substr(4 - RegExp.$1.length))
  }
  for (var k in o) {
    if (new RegExp('(' + k + ')').test(fmt)) {
      // eslint-disable-next-line eqeqeq
      fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (('00' + o[k]).substr(('' + o[k]).length)))
    }
  }
  return fmt
}

export { getCurrentDate, getYesterDays, getCurrentWeek, getCurrentMonth, getNextMonth, getCurrentSeason, getCurrentYear, getSevenDays, getThirtyDays }
