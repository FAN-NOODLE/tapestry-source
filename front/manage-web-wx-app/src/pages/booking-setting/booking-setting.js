import { timeListsAPi, timeConfigApi } from '@/api/booking-api'

const { globalData } = getApp()

const tableHeader = [
  {
    prop: 'renderTime',
    label: '开放时间段'
  },
  {
    prop: 'reserve_num',
    label: '剩余预约数'
  },
  {
    prop: 'reserved_num',
    label: '已预约人数'
  }
]

/** 生成默认时间戳 */
const generateDefaultTime = (date) => {
  const myDate = new Date(date)
  // 获取年月日
  const year = myDate.getFullYear()
  const month = myDate.getMonth() + 1
  const day = myDate.getDate()

  /** 默认时间段 */
  const defaultTimes = [
    {time: '9:00', reserve_num: 0, reserved_num: 0},
    {time: '10:00', reserve_num: 0, reserved_num: 0},
    {time: '11:00', reserve_num: 0, reserved_num: 0},
    {time: '13:00', reserve_num: 0, reserved_num: 0},
    {time: '14:00', reserve_num: 0, reserved_num: 0},
    {time: '15:00', reserve_num: 0, reserved_num: 0},
    {time: '16:00', reserve_num: 0, reserved_num: 0}
  ]

  // 生成上午时间戳
  const times = defaultTimes.map((item) => {
    let time = year + '/' + month + '/' + day + ' ' + item.time
    console.log(444, time)
    // time = time.replace(/-/g, '/')
    const reserve_time = new Date(time).getTime()
    return {
      ...item,
      reserve_time: reserve_time / 1000
    }
  })

  return times
}

Page({
  data: {
    /** 表格配置 */
    tableHeader: tableHeader,
    /** 上午时间段 */
    tableDataTop: [],
    /** 下午时间段 */
    tableDataBottom: [],
    /** 时间选择 */
    rangeList: [],
    /** 下拉时间段 */
    searchTime: '',
    /** 选中时间第0个 */
    searchTimeIndex: 0,
    /** 是否开启了编辑 */
    isEdit: false,
    /** 判断是否需要适配 */
    isIPhoneX: globalData.isIPhoneX
  },

  onLoad () {
    this.getTimeList(true)
  },

  /** 保存时间段 */
  saveTimeSetting () {
    // 组合参数
    const { tableDataTop, tableDataBottom } = this.data
    let data = [...tableDataTop, ...tableDataBottom]

    // 重组数据格式
    data = data.map((item) => ({
      reserve_time: item.reserve_time,
      reserve_num: item.reserve_num
    }))
    timeConfigApi({configs: data, activity_id: globalData.activity_id})
      .then(res => {
      // 错误处理
        if (res.code !== 200) {
          this.cancel()
          return wx.showToast(res.message)
        }

        // 设置编辑状态
        this.setEdit()
      })
      .catch(e => console.error(e, '获取活动预约时间段列表报错'))
  },

  /** 修改编辑状态 */
  setEdit () {
    this.setData({
      isEdit: !this.data.isEdit
    })
  },

  /** 取消 */
  cancel () {
    this.getTimeList()
    this.setEdit()
  },

  /** 区分上午、下午分别存储 */
  differenceTime (data) {
    const tableDataTop = [] // 上午时间段
    const tableDataBottom = [] // 下午时间段
    const currentTime = new Date().getTime() // 当前时间（时间戳）
    data.forEach((item) => {
      // 当前小时
      let hour = new Date(item.reserve_time * 1000).getHours()
      const disabled = currentTime >= item.reserve_time * 1000 + 3600000

      let renderTime = ''
      // 区分上下午
      if (hour >= 12) {
        // 下午

        hour = hour - 12
        renderTime = `${hour}:00~${hour + 1}:00`
        tableDataBottom.push({...item, renderTime, disabled})
      } else {
        // 上午
        renderTime = `${hour}:00~${hour + 1}:00`
        tableDataTop.push({...item, renderTime, disabled})
      }
    })
    // 返回上午、下午时间段，统一setData
    return {tableDataTop, tableDataBottom}
  },

  /** 获取活动预约时间段列表 */
  getTimeList (once = false) {
    const { searchTime } = this.data
    const time = searchTime.replace(/-/g, '/')
    let reserve_time = 0
    if (time) {
      reserve_time = ~~(new Date(time).getTime() / 1000)
    } else {
      reserve_time = ~~(new Date().getTime() / 1000)
    }


    timeListsAPi({ reserve_time })
      .then(res => {
        // 错误处理
        if (res.code !== 200) return wx.showToast(res.message)

        const { days, today, lists } = res.data

        let times = []
        // 判断时间戳list是否为空
        if (!lists.length) {
          // 使用默认时间戳
          times = generateDefaultTime(reserve_time * 1000)
          console.log(333333, times)
        } else {
          // 采用接口数据
          // const cache = [{
          //   reserve_num: 60,
          //   reserve_time: 1665651600,
          //   reserved_num: 0
          // },
          // {
          //   reserve_num: 60,
          //   reserve_time: 1665655200,
          //   reserved_num: 0
          // }]
          times = [...lists]
        }
        const tableData = this.differenceTime(times)
        const search = {}
        // 首次进入为搜索赋值
        if (once) {
          const index = days.findIndex(item => item === today)
          search.searchTime = today
          search.searchTimeIndex = index
        }
        this.setData({
          rangeList: days,
          ...search,
          ...tableData
        })
      })
      .catch(e => console.error(e, '获取活动预约时间段列表报错'))
  },

  /** 预约日期变化 */
  bindChange (e) {
    // 下标
    const index = e.detail.value
    this.setData({
      searchTime: this.data.rangeList[index],
      searchTimeIndex: index
    })
    this.getTimeList()
  },

  /** input 只能输入数字 */
  inputNumber (e) {
    let num = e.detail.value
    return num.replace(/[^0-9]/gi, '').replace(/^0+/g, '')
  },

  /** 剩余预约人数变化 */
  peopleChange (e) {
    // 修改下标e.target.dataset.index, 值为e.detail.value
    const index = e.target.dataset.index
    const name = e.target.dataset.name
    const value = e.detail.value

    const { tableDataTop, tableDataBottom } = this.data
    if (name === 'top') {
      tableDataTop[index]['reserve_num'] = value || 0
    } else {
      tableDataBottom[index]['reserve_num'] = value || 0
    }

    this.setData({
      tableDataTop,
      tableDataBottom
    })
  }
})

