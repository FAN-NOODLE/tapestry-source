
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

function getCurrentDate(...value) {
  return new Date().setHours(...value)
}

function defaultFormatter(type, value) {
  return value
}
/**
 *
 * 时间选择器，(时间多列选择器)
 *
 * 获取每一列的选择项
 * 通过条件限制选项内容
 * 切换修改选项内容
 *
 */

Component({
  properties: {
    show: {
      type: Boolean,
      value: false,
      observer: 'updateValue'
    },
    title: {
      type: String,
      value: '请选择时间'
    },
    value: {
      type: Number,
      value: 0
    },
    // datetime | date | year-month
    type: {
      type: String,
      value: 'datetime'
    },
    defaultTime: {
      type: String,
      value: '00:00:00'
    },
    minDate: {
      type: Number,
      value: new Date(currentYear - 10, 0, 1).getTime()
    },
    maxDate: {
      type: Number,
      value: new Date(currentYear + 10, 11, 31, 23, 59, 59).getTime()
    },
    format: {
      type: String,
      value: 'yyyy-MM-dd hh:mm'
    },
    formatter: {
      type: null,
      value: defaultFormatter
    },
    useSlot: {
      type: Boolean,
      value: false
    }
  },
  data: {
    /** 内部时间: [start, end] */
    innerValue: 0,
    /** 默认先选择开始时间 index */
    current: 0,
    /** 选择区间 [year[], month[], day[], hour[], minute[]] */
    columns: [],
    /** 选中下标 [yearIndex, monthIndex, dayIndex, hourIndex, minuteIndex] */
    columnIndex: [],
    types: [],
    isScrolling: false
  },
  lifetimes: {
    attached() {}
  },
  methods: {
    handleShow() {
      this.setData({
        show: true
      })
    },
    /**
     * 异步设置
     * @param {*} data
     * @returns
     */
    set(data) {
      this.setData(data)
      return new Promise(function (resolve) {
        return wx.nextTick(resolve)
      })
    },
    updateValue(show) {
      if (!show) return
      const { value, type, defaultTime } = this.data
      // 获取有效的index
      // 获取当前下标对应的值
      var currentValue = value || getCurrentDate(...defaultTime.split(':'))
      // 格式化值的范围
      currentValue = this.correctValue(currentValue)
      console.log(currentValue)
      // 获取lable
      const types = ['年', '月', '日', '时', '分']
      if (type === 'date') types.splice(3, 2)
      if (type === 'year-month') types.splice(2, 3)
      // 更新数据
      this.set({
        types,
        innerValue: currentValue,
        isScrolling: false
      }).then(() => {
        this.updateColumnValue(currentValue)
      }).catch((err) => {
        console.log(err)
      })
    },
    /**
     *
     * @param {Number} value 时间戳
     * @returns
     */
    updateColumnValue(value) {
      console.log('updateColumnValue', value)
      let values = []
      const { type } = this.data
      const formatter = this.data.formatter || defaultFormatter
      const date = new Date(value)
      values = [
        formatter('year', ''.concat(date.getFullYear())),
        formatter('month', padZero(date.getMonth() + 1))
      ]
      if (type === 'date') {
        values.push(formatter('day', padZero(date.getDate())))
      }
      if (type === 'datetime') {
        values.push(
          formatter('day', padZero(date.getDate())),
          formatter('hour', padZero(date.getHours())),
          formatter('minute', padZero(date.getMinutes())),
        )
      }
      return this.set({ innerValue: value })
        .then(async () => {
          await this.updateColumns()
          console.log('updateColumnValue', values)
          const columnIndex = this.data.columns.map((item, index) =>
            item.values.findIndex(val => values[index] === val),
          )
          console.log('updateColumnValue', columnIndex)

          this.set({
            columnIndex
          })
        })
    },
    /**
     * @description 获取最终选想
     */
    updateColumns() {
      const formatter = this.data.formatter || defaultFormatter
      const result = this.getOriginColumns().map(function (column) {
        return {
          values: column.values.map(function (value) {
            return formatter(column.type, value)
          })
        }
      })
      console.log('updateColumns', result)
      return this.set({
        columns: result
      })
    },
    /**
     * 获取原始选项
     * @returns
     */
    getOriginColumns() {
      return this.getRanges().map(({ type, range }) => {
        const values = times(range[1] - range[0] + 1, function (index) {
          const value = range[0] + index
          return type === 'year' ? ''.concat(value) : padZero(value)
        })
        return { type: type, values: values }
      })
    },
    /**
     * 获取时间范围
     * @returns
     */
    getRanges() {
      const { innerValue, type } = this.data
      const currentValue = innerValue
      const { maxYear, maxMonth, maxDate, maxHour, maxMinute } =
        this.getBoundary('max', currentValue)
      const { minYear, minMonth, minDate, minHour, minMinute } =
        this.getBoundary('min', currentValue)

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
        }
      ]
      if (type === 'date') result.splice(3, 2)
      if (type === 'year-month') result.splice(2, 3)
      return result
    },
    /**
     * 获取边界值
     * @param {String} type max/min
     * @param {Number} currentValue 时间戳
     * @returns
     */
    getBoundary(type, currentValue) {
      var result = {}
      var value = new Date(currentValue)
      var boundary = new Date(this.data[''.concat(type, 'Date')])
      var year = boundary.getFullYear()
      var month = 1
      var date = 1
      var hour = 0
      var minute = 0
      if (type === 'max') {
        month = 12
        date = getMonthEndDay(value.getFullYear(), value.getMonth() + 1)
        hour = 23
        minute = 59
      }
      if (value.getFullYear() === year) {
        month = boundary.getMonth() + 1
        if (value.getMonth() + 1 === month) {
          date = boundary.getDate()
          if (value.getDate() === date) {
            hour = boundary.getHours()
            if (value.getHours() === hour) {
              minute = boundary.getMinutes()
            }
          }
        }
      }

      result[''.concat(type, 'Year')] = year
      result[''.concat(type, 'Month')] = month
      result[''.concat(type, 'Date')] = date
      result[''.concat(type, 'Hour')] = hour
      result[''.concat(type, 'Minute')] = minute
      return result
    },
    /**
     * 获取时间
     * @param {*} value
     * @returns
     */
    correctValue(value) {
      const { minDate, maxDate } = this.data
      // date type
      value = Math.max(value, minDate)
      value = Math.min(value, maxDate)
      return value
    },
    /**
     * todo:picker滚动开始，判断是否是手动更新
     */
    handleStart(e) {
      console.log('handleStart', e)
      this.data.isScrolling = true
    },
    /**
     * picker变化
     * @param {*} e
     */
    handleChange(e) {
      const { value } = e.detail
      const { type, isScrolling, minDate } = this.data
      console.log('isScrolling', isScrolling)
      if (!isScrolling) return
      this.data.isScrolling = false
      var result
      var originColumns = this.getOriginColumns()
      var values = value.map(function (value, index) {
        return originColumns[index].values[value]
      })

      var year = getTrueValue(values[0])
      var month = getTrueValue(values[1])
      // var maxDate = getMonthEndDay(year, month)
      var date = getTrueValue(values[2])
      if (type === 'year-month') {
        date = 1
      }
      // date = date > maxDate ? maxDate : date
      var hour = 0
      var minute = 0
      if (type === 'datetime') {
        hour = getTrueValue(values[3])
        minute = getTrueValue(values[4])
      }
      result = new Date(year, month - 1, date, hour, minute)

      result = result < minDate ? minDate : result

      result = this.correctValue(result)

      console.log('handleChange', result, e)
      this.updateColumnValue(result)
    },
    handleEnd(e) {
      console.log('handleEnd', e)
      // this.data.isScrolling = false
    },
    /**
     * @description 取消
     */
    handleCancel() {
      this.triggerEvent('cancel')
    },
    /**
     * @description 确定
     */
    handleConfirm() {
      const { innerValue } = this.data
      this.triggerEvent('confirm', {
        value: innerValue
      })
    },
    /**
     * @description
     */
    handleReset() {
      this.triggerEvent('cancel')
    }
  }
})
