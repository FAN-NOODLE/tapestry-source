
// Component Object
const currentYear = new Date().getFullYear()
function padZero(val) {
  return '00'.concat(val).slice(-2)
}
function getMonthEndDay(year, month) {
  return 32 - new Date(year, month - 1, 32).getDate()
}
function times(n, iteratee) {
  var index = -1
  var result = Array(n < 0 ? 0 : n)
  while (++index < n) {
    result[index] = iteratee(index)
  }
  return result
}
function getTrueValue(formattedValue) {
  if (formattedValue === undefined) {
    formattedValue = '1'
  }
  while (isNaN(parseInt(formattedValue, 10))) {
    formattedValue = formattedValue.slice(1)
  }
  return parseInt(formattedValue, 10)
}
Component({
  properties: {
    value: {
      type: Array,
      value: [],
      observer: 'updateValue'
    },
    minDate: {
      type: Number,
      value: new Date(currentYear - 10, 0, 1).getTime(),
      observer: 'updateValue'
    },
    maxDate: {
      type: Number,
      value: new Date(currentYear + 10, 11, 31, 23, 59, 59).getTime(),
      observer: 'updateValue'
    },
    formatter: {
      type: String,
      value: 'yyyy-MM-dd hh:mm:ss'
    },
    useSlot: {
      type: Boolean,
      value: false
    }
  },
  data: {
    show: false,
    /** 内部时间 */
    innerValue: [],
    /** 默认先选择开始时间 */
    current: 0,
    /** 当前数据 */
    currentValue: [],
    /** 选择区间 */
    columns: [],
    /** 选中下标 */
    columnIndex: []
  },
  methods: {
    /**
     * @description 显示弹窗
     */
    show() {
      this.setData({
        show: true,
        current: 0
      })
    },
    /**
     * @description 隐藏弹窗
     */
    hide() {
      this.setData({
        show: false
      })
    },
    /**
     * @description 更新数据
     */
    updateValue() {
      const { value, minDate, maxDate, innerValue } = this.data
      // 排序
      if (value.length) {
        value.sort((a, b) => a - b)
      } else {
        // 设置默认开始时间为当前时间+5分钟
        const defaultStartTime = new Date().getTime() + 5 * 60 * 1000
        value.push(defaultStartTime)
        value.push(defaultStartTime)
      }

      // 格式化收藏
      if (value[1] < minDate) {
        value[0] = minDate
        value[1] = minDate
      } else if (value[0] > maxDate) {
        value[0] = maxDate
        value[1] = maxDate
      }

      // 获取开始时间 - 保留原有时分秒
      const startTime = new Date(Math.max(value[0], minDate))
      innerValue[0] = startTime.getTime()

      // 获取结束时间 - 保留原有时分秒
      const endTime = new Date(Math.min(value[1], maxDate))
      innerValue[1] = endTime.getTime()
      console.log('innerValue2', innerValue)

      this.setData({
        innerValue
      })
      // 先处理第一个
      this.updateColumnValue(innerValue[0])
    },
    /**
     * @description 更新数据
     */
    updateColumnValue(value) {
      const date = new Date(value)
      const values = []
      // 获取年份
      values.push(''.concat(date.getFullYear()))
      // 获取月份
      values.push(padZero(date.getMonth() + 1))
      // 获取日期
      values.push(padZero(date.getDate()))
      // 获取小时
      values.push(padZero(date.getHours()))
      // 获取分钟
      values.push(padZero(date.getMinutes()))
      // 获取秒
      values.push(padZero(date.getSeconds()))
      this.setData({
        currentValue: values
      }, () => {
        this.updateColumn()
      })
    },
    /**
     * @description 更新数据下标
     */
    updateColumn() {
      const result = this.getOriginColumns().map(function (column) {
        return {
          values: column.values.map(function (value) { return value })
        }
      })
      this.setData({
        columns: result
      }, () => {
        this.updateColumnIndex()
      })
    },
    updateColumnIndex() {
      const { currentValue, columns } = this.data
      const columnIndex = currentValue.map((item, index) => {
        return columns[index].values.findIndex(column => column === item)
      })

      this.setData({
        columnIndex
      })
    },
    getOriginColumns() {
      return this.getRanges().map(({type, range}) => {
        const values = times(range[1] - range[0] + 1, function (index) {
          const value = range[0] + index
          return type === 'year' ? ''.concat(value) : padZero(value)
        })
        return { type: type, values: values }
      })
    },
    /** 获取时间区间 */
    getRanges () {
      const { currentValue } = this.data
      const { maxYear, maxDate, maxMonth, maxHour, maxMinute, maxSecond } = this.getBoundary('max', currentValue)
      const { minYear, minDate, minMonth, minHour, minMinute, minSecond } = this.getBoundary('min', currentValue)

      const result = [
        {
          type: 'year',
          range: [minYear, maxYear]
        },
        {
          type: 'month',
          range: [minMonth, maxMonth]
        },
        {
          type: 'day',
          range: [minDate, maxDate]
        },
        {
          type: 'hour',
          range: [minHour, maxHour]
        },
        {
          type: 'minute',
          range: [minMinute, maxMinute]
        },
        {
          type: 'second',
          range: [minSecond, maxSecond]
        }
      ]
      return result
    },
    getBoundary(type, currentValue) {
      var result = {}
      // 直接从currentValue数组中获取年、月、日等信息，并转换为数字
      var year = getTrueValue(currentValue[0])
      var month = getTrueValue(currentValue[1])
      var day = getTrueValue(currentValue[2])
      var hour = getTrueValue(currentValue[3])
      var minute = getTrueValue(currentValue[4])

      var boundary = new Date(this.data[''.concat(type, 'Date')])
      var boundaryYear = boundary.getFullYear()
      var boundaryMonth = boundary.getMonth() + 1
      var boundaryDay = boundary.getDate()
      var boundaryHour = boundary.getHours()
      var boundaryMinute = boundary.getMinutes()
      var boundarySecond = boundary.getSeconds()

      var resultYear = boundaryYear
      var resultMonth = 1
      var resultDate = 1
      var resultHour = 0
      var resultMinute = 0
      var resultSecond = 0

      if (type === 'max') {
        resultMonth = 12
        // 使用当前选择的年月来计算月末天数
        resultDate = getMonthEndDay(year, month)
        resultHour = 23
        resultMinute = 59
        resultSecond = 59
      }

      if (year === boundaryYear) {
        resultMonth = boundaryMonth
        if (month === boundaryMonth) {
          resultDate = boundaryDay
          if (day === boundaryDay) {
            resultHour = boundaryHour
            if (hour === boundaryHour) {
              resultMinute = boundaryMinute
              if (minute === boundaryMinute) {
                resultSecond = boundarySecond
              }
            }
          }
        }
      }

      result[''.concat(type, 'Year')] = resultYear
      result[''.concat(type, 'Month')] = resultMonth
      result[''.concat(type, 'Date')] = resultDate
      result[''.concat(type, 'Hour')] = resultHour
      result[''.concat(type, 'Minute')] = resultMinute
      result[''.concat(type, 'Second')] = resultSecond
      return result
    },
    handleInputChange(e) {
      const { innerValue } = this.data
      const { index } = e.currentTarget.dataset
      if (index === 1 && innerValue[1] < innerValue[0]) {
        innerValue[1] = innerValue[0]
      }
      this.setData({
        current: index,
        innerValue
      }, () => {
        this.updateColumnValue(innerValue[index])
      })
    },
    handleDateChange(e) {
      const { value } = e.detail
      const { columns, current } = this.data

      const currentValue = value.map((item, index) => {
        return columns[index].values[item]
      })

      this.setData({
        columnIndex: value,
        currentValue
      }, () => {
        const year = getTrueValue(currentValue[0])
        const month = getTrueValue(currentValue[1])
        const date = getTrueValue(currentValue[2])
        const hour = getTrueValue(currentValue[3])
        const minute = getTrueValue(currentValue[4])
        const second = getTrueValue(currentValue[5])

        // 检查日期是否超过当月天数，如果超过则调整为当月最后一天
        const maxDate = getMonthEndDay(year, month)
        const adjustedDate = Math.min(date, maxDate)

        const newTimestamp = new Date(year, month - 1, adjustedDate, hour, minute, second).getTime()

        this.setData({
          [`innerValue[${current}]`]: newTimestamp
        }, () => {
          // 更新columnValue以重新生成columns数组，确保天数范围根据新月份更新
          this.updateColumnValue(newTimestamp)
        })
      })
    },
    handleCancel() {
      this.hide()
    },
    /**
     * @description 确认
     */
    handleConfirm() {
      const { innerValue } = this.data
      console.log('innerValue1', innerValue)

      this.triggerEvent('change', { value: innerValue })
      this.hide()
    }
  }
})

