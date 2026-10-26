import Calendar from '@/utils/calendar'
import { monthListApi } from '@/api/plan-task-api'

/** 周 */
const WEEK_NAME_MAP = ['SUN', 'MON', 'TUE', 'WED', 'THU', 'FRI', 'SAT']

/** 生成日期选择data */
const generateTime = (time) => {
  const currentYear = time.getFullYear()
  const currentMonth = time.getMonth() + 1

  const thisYear = new Date().getFullYear()
  const years = [thisYear - 2, thisYear - 1, thisYear]

  let months = []

  for (let i = 1; i <= 12; i++) {
    months.push(i)
  }

  const pickerIndex = [years.indexOf(currentYear) !== -1 ? years.indexOf(currentYear) : '', months.indexOf(currentMonth) !== -1 ? months.indexOf(currentMonth) : '']
  return {years, months, currentYear, currentMonth, pickerIndex}
}

Component({
  /** 自定义样式 */
  externalClasses: ['custom-calendar'],
  properties: {
    /** 完成状态 1 已完成 2 待完成 3 不做展示 其他数字展示全部 */
    status: { type: Number, value: 3 },
    /** 页面状态 1计划任务，2播种，3到店邀约 */
    pageType: { type: Number, value: 1 },
    /** 默认选中时间 */
    calendarTime: { type: Number },
    /** 其他参数 */
    otherParams: {type: Object, value: {}}
  },
  /** 页面的初始数据 */
  data: {
    /** 周 */
    weekData: WEEK_NAME_MAP,
    /** 日历 */
    calendarData: [],
    /** 判断日期选择是否打开 */
    pickerShow: false,
    /** 日期选择器选中picker index */
    pickerIndex: [],
    /** 本月详细信息 */
    monthInfo: {},
    /** 是否折叠日历面板 */
    collapse: true,
    /** 判断日期组件是否开始滚动 */
    scrollEnd: false,
    /** 颜色值，右上角状态 */
    colors: [],
    /** 是否点击了返回今天 */
    isClickToday: false,
    /** 是否是刷新 */
    isRefresh: false,
    /** 当前折叠的月份 */
    collapseWeek: 1
  },
  /** 生命周期函数--监听页面加载 */
  ready () {
    this.init()
  },
  methods: {
    /** 获取日历状态 */
    getMonthList (start_time, end_time) {
      return new Promise((resolve) => {
        const {pageType, otherParams} = this.properties
        const params = {
          start_time,
          end_time: end_time + 86399,
          // 1计划任务，2播种，3到店邀约
          type: Number(pageType)
        }

        // 如果是任务计划并且id为真，则是子任务详情日历接口
        if (Number(pageType) === 1 && otherParams.id) {
          params.plan_id = Number(otherParams.id)
        }

        monthListApi(params)
          .then((res) => {
            resolve(res.data.List || [])
          })
          .catch(e => console.error('获取日历状态报错', e))
      })
    },

    /** 判断数组中是否存在连续相同的属性 */
    continuedData (data, attributes = []) {
      let fillet = new Array(attributes.length).fill(0)
      data.forEach((item, index) => {
        attributes.forEach((attribute, key) => {
          if (item[attribute] !== undefined && data[index - 1] && data[index - 1][attribute] === item[attribute] && data[index + 1] && data[index + 1][attribute] === item[attribute]) {
            item[`log-${attribute}`] = item[attribute]
            data[index - 1][`log-${attribute}`] = item[attribute]
            data[index + 1][`log-${attribute}`] = item[attribute]

            // 第一个元素
            if (fillet[key] === 0) {
              data[index - 1][`log-${attribute}-start`] = item[attribute]
              fillet[key] = 1
            }

            // 最后一个元素
            if (fillet[key] === 1 && ((data[index + 2] && data[index + 2][attribute] !== item[attribute]) || !data[index + 2])) {
              data[index + 1][`log-${attribute}-end`] = item[attribute]
              fillet[key] = 0
            }
          }
        })
      })
      return data
    },

    /** 生成日历 */
    generateCalendar (year, month) {
      // 只触发一次默认点击
      let once = true
      const calendarClass = new Calendar(
        async (list, data) => {
          const monthList = await this.getMonthList(data.firstDayNum, data.lastDayNum)
        
          // 这里可以添加业务处理函数
          let newList = list.map((item) => {
            const matchingData = monthList.find(matchingItem => item.dayToNum === matchingItem.time)
            return {
              ...item,
              ...matchingData
            }
          })
          newList = this.continuedData(newList, ['finished', 'afoot'])
          // 默认触发点击事件target.dataset.item
          if (once) {
            let currentDay = data
            // 指定选中
            if (this.properties.calendarTime && !this.data.isClickToday) {
              currentDay = newList.filter(item => item.dayToNum === this.properties.calendarTime)[0]
            }
            once = false
            this.clickDay({target: {dataset: {item: currentDay}}})
          }
         
          const monthInfo = {...data, currentWeek: data.currentWeek || 1, ...this.data.monthInfo}
          this.setData({
            isClickToday: false,
            calendarData: newList,
            monthInfo,
            collapse: true,
            collapseWeek: monthInfo.month === data.month ? monthInfo.currentWeek || 1 : 1
          })
          return list
        },
        year, month
      )

      this.data.calendarClass = calendarClass
    },

    /** 初始化数据 */
    init (isRefresh = false) {
      // 设置刷新情况
      this.setData({
        isRefresh
      })
      let timeStamp = this.properties.calendarTime
      if (this.data.isClickToday) {
        timeStamp = new Date()
      } else {
        timeStamp = timeStamp ? new Date(timeStamp * 1000) : new Date()
      }
     
      const {currentYear, currentMonth, pickerIndex, ...rest} = generateTime(timeStamp)
      this.generateCalendar(currentYear, currentMonth)

      const colors = this.properties.pageType === 3
        ? [
          {name: '完成/成功', color: '#BBD1A7'},
          {name: '进行中', color: '#D1BFA7'},
          {name: '连续完成/成功', color: '#BBD1A7'},
          {name: '连续进行中', color: '#D1BFA7'}
        ]
        : [
          {name: '已完成', color: '#BBD1A7'},
          {name: '待完成', color: '#E6B8B8'},
          {name: '连续完成', color: '#BBD1A7'},
          {name: '连续待完成', color: '#E6B8B8'}
        ]
      // 生成颜色
      this.setData({
        currentYear,
        currentMonth,
        colors,
        ...rest
      })

      // 防止 首次选中默认值不生效
      this.setData({
        pickerIndex
      })
    },

    /** 点击某一天 */
    clickDay (e) {
      const dayInfo = e.target.dataset.item || (e.currentTarget && e.currentTarget.dataset.item)
      // 置灰 或者day 未空 的不可点击
      if ((this.data.pageType === 1 && dayInfo.disabled) || !dayInfo.day) return
      // 防止重复触发
      if (this.data.monthInfo.day === dayInfo.day && !this.data.isRefresh) return

      // 时间戳
      let dayToNum = dayInfo.dayToNum
      let day = dayInfo.day || 1

      // 如果没有dayToNum，则是默认触发的点击事件（默认查询1号数据）
      if (!dayInfo.dayToNum) {
        const dayToString = `${dayInfo.year}/${dayInfo.month}/${day}`
        dayToNum = ~~(new Date(dayToString).getTime() / 1000)
      }

      // 设置选中
      this.setData({
        monthInfo: dayInfo,
        isRefresh: false,
        collapseWeek: dayInfo.currentWeek
      })
      // 配合接口如果没有选中的日期默认查询所有 0，
      const timeNum = dayInfo.day === 0 ? 0 : dayToNum
      // 通知父组件 选中日历变化
      this.triggerEvent('changeDay', timeNum)
    },

    /** 回到今天 */
    handleToday () {
      this.setData({
        isClickToday: true
      })
      this.init()
    },

    /** 选择确认 */
    confirm () {
      const { pickerIndex, years, months } = this.data

      const currentYear = years[pickerIndex[0]]
      const currentMonth = months[pickerIndex[1]]
      this.setData({
        currentYear,
        currentMonth,
        pickerShow: false
      })
      // 指定月份查询
      this.data.calendarClass.appointYearOrMonth(currentYear, currentMonth)
    },

    /** 收起折叠日历 */
    pullDown () {
      this.setData({
        collapse: !this.data.collapse
      })
    },

    /** 打开日期选择 */
    openPicker () {
      this.setData({
        pickerShow: true
      })
    },

    /** 关闭日期选择 */
    closePicker () {
      this.setData({
        pickerShow: false
      })
    },

    /** 日期组件滚动停止 */
    pickerEnd () {
      this.setData({
        scrollEnd: false
      })
    },

    /** 日期组件滚动开始 */
    pickerStart () {
      this.setData({
        scrollEnd: true
      })
    },

    /** 日期选择change */
    pickerChange (e) {
      this.setData({
        pickerIndex: e.detail.value
      })
    }
  }
})

