import {choiceOptionsApi} from '@/api/plan-task-api'
import dayjs from 'dayjs'

const { globalData } = getApp()

/** 生成筛选条件 */
const generateData = (keys = [], list) => {
  // 默认全部类型
  let data = [
    {show: true, title: '顾问筛选', key: 'manage_type', multiple: false, list: []},
    {show: true, title: '员工筛选', key: 'member_list', multiple: true, list: []},
    {show: true, title: '会员卡等', key: 'card_type', multiple: true, list: []},
    // {show: true, title: 'CLT类型', key: 'clt_type', multiple: true, list: []},
    {show: false, title: '客户标签', key: 'tag_list', multiple: true, list: []},
    {show: false, title: '印象标签', key: 'impression_tag_list', multiple: true, list: []}
  ]

  // 根据key 过滤数据,并且添加数据
  data = data.filter(item => keys.includes(item.key)).map(item => ({...item, list: list[item.key] || []}))
  return data
}

Component({
  // 使用全局 & 页面样式
  options: {
    addGlobalClass: true
  },
  properties: {
    /** 展示过滤类型，对应接口choiceOptionsApi 的key，默认展示全部 */
    keys: {
      type: Array,
      value: ['manage_type', 'member_list', 'card_type', 'clt_type', 'tag_list', 'impression_tag_list']
    },
    /** 回显结果集合 */
    result: {
      type: Object,
      value: {}
    },
    currentPage: {
      type: Number,
      value: 2
    },
    filterType: {
      type: Number,
      value: 1
    }
  },
  data: {
    /** 是否展示、关闭预约信息弹窗 */
    show: false,

    /** 渲染数据list */
    selectList: [],

    /** 判断是否适配 */
    isIPhoneX: globalData.isIPhoneX
  },
  methods: {
    /** 关闭弹窗-内部方法 */
    close () {
      // this.resetData({}, true)
      this.updatePopup(false)
    },

    /** 打开、关闭弹窗, 供外部调用 */
    updatePopup (show = true) {
      // 如果是打开状态, 获取数据
      if (show) {
        wx.hideTabBar()
        this.initData()
      } else {
        wx.showTabBar()
      }
      this.setData({
        show
      })

      // 通知父组件子组件打开了，关闭父组件滚动
      this.triggerEvent('onScroll', !show)
    },

    /** 初始化数据 */
    initData () {
      choiceOptionsApi()
        .then((res) => {
          if (res.code !== 200) return wx.showToast(res.message)
          // 判断页面类型，区分显示
          const list = generateData(this.properties.keys, res.data)
          this.setData({
            selectList: list
          })
        })
        .catch(e => console.error('获取筛选数据出错', e))
    },

    /** 重置数据 */
    resetData (e, isClose = false) {
      const { selectList } = this.data
      const result = {}

      selectList.forEach(v => {
        this.selectComponent(`#${v.key}`).setData({ result: [] })
      })
      this.updatePopup(false)

      !isClose && this.triggerEvent('confirm', result)
    },

    /** 确定提交筛选数据 */
    confirm () {
      const { selectList } = this.data
      const result = {}

      selectList.forEach(v => {
        const itemResult = this.selectComponent(`#${v.key}`).data.result

        if (itemResult.length > 0) {
          result[v.key] = v.multiple ? itemResult : itemResult[0]
        }
      })

      this.track('确定')

      this.updatePopup(false)
      console.log('confirm', result)
      this.triggerEvent('confirm', result)
    },
    track(button_name) {
      const { currentPage, filterType, selectList } = this.data

      const result = []

      selectList.forEach(v => {
        const itemResult = this.selectComponent(`#${v.key}`).data.result

        if (itemResult.length > 0) {
          if (v.multiple) {
            result.push(...itemResult.map(item => {
              return {
                tag_level_1: '',
                tag_level_2: v.title,
                tag_value: v.list.find(val => val.key === item).value
              }
            }))
          } else {
            result.push({
              tag_level_1: '',
              tag_level_2: v.title,
              tag_value: v.list.find(val => val.key === itemResult[0]).value
            })
          }
        }
      })


      const result_option = result.reduce((prev, cur) => {
        if (prev[cur.tag_level_2]) {
          prev[cur.tag_level_2].push(cur.tag_value)
        } else {
          const item = selectList.find(item => item.title === cur.tag_level_2)
          prev[cur.tag_level_2] = item.multiple ? [cur.tag_value] : cur.tag_value
        }
        return prev
      }, {})

      // 神策数据上报 -- 筛选
      wx.track.trackFilter({
        current_page: currentPage,
        button_name: button_name,
        filter_type_sv: filterType,
        filter_option: result_option
      })

      result.forEach(item => {
        // 神策数据上报 -- 筛选详情
        wx.track.trackFilterDetails({
          current_page: currentPage,
          button_name: button_name,
          filter_type_sv: filterType,
          ...item
        })
      })
    }
  }
})
