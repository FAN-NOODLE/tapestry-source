import { gradePointsApi, pointsListApi, getPointsRulesApi, getLevelRecordApi, getCoinsListApi, getCoinsTokenApi } from '@/api/we-journey-api'
// eslint-disable-next-line standard/object-curly-even-spacing
import { toDate} from '@/utils/tools'
import { getRedirectPrarms } from '../utils'

// 获取全局应用程序实例对象
const app = getApp()
const imagePath = 'https://cdn-cos.coach.com.cn/uploads/'

const defaultDimension = new Map([
  [1, '邀约能力'],
  [2, '任务能力'],
  [3, '互动能力'],
  [4, '培训能力'],
  [5, '线上能力']
])

const defaultCondition = new Map([
  [2, '高于'],
  [3, '达到'],
  [4, '低于']
])

const defaultTypeIcon = new Map([
  [1, `${imagePath}2023/v3.5.0/task/icon_1.png`], // 邀约类
  [2, `${imagePath}2023/v3.5.0/task/icon_2.png`], // 任务类
  [3, `${imagePath}2023/v3.5.0/task/icon_3.png`], // 培训类
  [4, `${imagePath}2023/v3.5.0/task/icon_2.png`] // 好友类
])

const linkList = [
  '/pages/task-no-tab/task-no-tab?type=0',
  '/pages/task-no-tab/task-no-tab?type=1',
  '/pages/task-no-tab/task-no-tab?type=2',
  'webview'
]

const defaultPointsMap = new Map([
  [1, {
    title: '发送到店邀约',
    unit: '次',
    link: linkList[2],
    points: 1,
    icon: defaultTypeIcon.get(1),
    tip: '每天最多可得5积分哦~'
  }],
  [7, {
    title: '任务完成',
    unit: '次',
    link: linkList[0],
    points: 1,
    icon: defaultTypeIcon.get(2)
  }],
  [15, {
    title: '完成welearning培训可得积分',
    unit: '',
    rule_value: '',
    link: linkList[3],
    hideNeedPoint: true,
    icon: defaultTypeIcon.get(3)
  }],
  [14, {
    title: '完成播种',
    unit: '次',
    link: linkList[1],
    points: 1,
    icon: defaultTypeIcon.get(4),
    tip: '每天最多可得5积分哦~'
  }],
  [6, {
    title: '邀约到店确认',
    unit: '次',
    link: linkList[2],
    points: 1,
    icon: defaultTypeIcon.get(1),
    tip: '每天最多可得5积分哦~'
  }],
  [9, {
    title: '任务触达',
    unit: '人',
    link: linkList[0],
    points: 1,
    icon: defaultTypeIcon.get(2)
  }],
  [4, {
    title: '邀约成功',
    unit: '次',
    link: linkList[2],
    points: 5,
    icon: defaultTypeIcon.get(1)
  }],
  [10, {
    title: '新增好友',
    unit: '人',
    link: false,
    points: 5,
    icon: defaultTypeIcon.get(4)
  }],
  [11, {
    title: '每月客户流失高于',
    unit: '人',
    link: false,
    points: 2,
    icon: defaultTypeIcon.get(4)
  }]
])

/**
 * 获取任务信息
 * @param {object} item
 */
function getTaskInfo (item) {
  const value = defaultPointsMap.get(item.rule_category_id)
  const frequency = ['', '天', '周', '月', '年'][item.astrict_type]
  // 获取提示信息
  if (frequency && item.max_integral) {
    value.tip = `每${frequency}最多可得${item.max_integral}积分哦~`
  }

  if (item.max_integral !== 0) {
    item.title_text = `${item.now_integral}/${item.max_integral}`
  }

  console.log(value)
  return {
    ...item,
    ...value
  }
}

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    navigationBarHeight: app.globalData.system.statusBarHeight + 44 + 'px',
    imagePath: imagePath,
    levelArr: [], // 等级列表
    settingList: [], // 获取积分渠道列表
    pointsList: [], // 积分领取信息列表
    foregroundWidth: 0, // 进度条已进行长度
    progressWidth: 0, // 进度条总长度
    listIndex: 0, // tap栏id
    gradePoints: {}, // 积分页面总信息
    pageInfo: {
      page: 1,
      size: 10,
      total: 0
    }, // 分页信息
    isLoading: false, // 是否正在加载分页
    fLen: -1, // 等级数量
    // loginToken
    loginToken: {
      moduleType: 'portal',
      moduleId: '0',
      uniqueKey: 'mobile',
      uniqueValue: '',
      companyCode: 'coach',
      time: new Date().getTime(),
      token: 'REDACTED_FOR_SOURCE_SHARING'
    },
    showPopupFootprint: false, // 是否显示足迹弹窗
    levelRecord: [], // 等级足迹
    showMaxLimit: 5, // 最多显示多少条
    showMore: {
      task: true,
      point: true,
      coins: true
    },
    scrollTop: {
      task: 0,
      point: 0,
      gold: 0
    },
    offset: 0,
    offsetY: 0,
    offsetX: 0,
    // 金币列表
    coinsList: [],
    // 金币分页
    coinsPageInfo: {
      page: 1,
      size: 10,
      total: 0
    }
  },

  // 生命周期函数--监听页面加载
  onLoad (options) {
    this.initTab(options)
    this.getLevelRecord()
    this.gradePoints()
  },
  /** todo: 初始化tab */
  initTab(options) {
    const { tab } = options
    if (tab) {
      const numIndex = Number(tab)
      if (numIndex === 0) {
        this.getPointsRules()
      } else if (numIndex === 1) {
        this.setData({
          'pageInfo.page': 1,
          listIndex: 1
        })
        this.pointsList()
      } else if (numIndex === 2) {
        this.setData({
          'coinsPageInfo.page': 1,
          listIndex: 2
        })
        this.coinsList()
      }
    } else {
      // 如果没有tab参数，默认获取积分渠道列表
      this.getPointsRules()
    }
  },
  /** todo:导航改变 */
  navChange (e) {
    const { index } = e.currentTarget.dataset
    const numIndex = Number(index)
    if (numIndex === 0) {
      this.getPointsRules()
    } else if (numIndex === 1) {
      this.setData({ 'pageInfo.page': 1 })
      this.pointsList()
    } else if (numIndex === 2) {
      this.setData({ 'coinsPageInfo.page': 1 })
      this.coinsList()
    }
    this.setData({
      listIndex: numIndex,
      offsetY: this.data.offset
    })
  },
  /** 跳转规则页面 */
  toRules () {
    wx.navigateTo({
      url: '/pages/weJourney/rule/rule'
    })
  },
  /* 获取积分总信息 */
  gradePoints () {
    gradePointsApi().then(res => {
      const { data, message, code } = res

      if (code !== 200) {
        wx.showToast(message)
        return
      }

      const { level_rule: levelRule } = data

      const levelArr = levelRule.map((i, idx) => {
        const min = `${i.min_value}≤`
        const max = idx === levelRule.length - 1 ? '' : `＜${i.max_value}`
        let word = `${min}总积分${max}`

        /* 有规则计算规则 */
        if (i.rule) {
          i.rule.forEach(i => {
            const dimension = defaultDimension.get(i.incentive_config_type)
            const condition = defaultCondition.get(i.rule_condition)
            let w = ''

            switch (i.rule_condition) {
              case 1:
                w += `,${dimension}TOP${i.rule_value}`
                break
              case 2:
              case 3:
              case 4:
                w += `,${dimension}${condition}${i.rule_give}`
                break
              default:
                break
            }

            word += w
          })
        }

        return { ...i, word }
      })

      // 当前处于什么等级
      const findIndex = levelArr.findIndex(i => i.selected)
      const filterArr = findIndex === -1 ? [] : levelArr.slice(0, findIndex + 1)
      const fLen = filterArr.length
      const lLen = levelArr.length
      const rpx = wx.getSystemInfoSync().windowWidth / 750
      // 计算进度条长度 = 两端的宽度 + 中间的宽度 + 最后一个等级的宽度
      // 两端的宽度 = (边距 + 边宽) * 2
      // 中间的宽度 = (等级数量 - 1) * (等级宽度 + 等级间距)
      const progressWidth = ((48 + 100) * 2 + (lLen - 1) * (34 + 135) + 42) * 2

      this.setData({
        gradePoints: {
          ...data,
          level_updated_at: toDate(data.level_updated_at, 'yyyy.MM.dd').nowTime
        },
        levelArr,
        progressWidth,
        fLen
      }, () => {
        if (fLen === 0) return

        const query = wx.createSelectorQuery()


        query.select('.progress').fields({
          computedStyle: ['padding-left', 'width', 'left']
        })

        query.selectAll('.level').boundingClientRect()
        query.select('.scroll-box').fields({
          computedStyle: ['left']
        })

        query.exec(res => {
          console.log(11111, res)
          if (fLen === lLen) {
            this.setData({ foregroundWidth: parseInt(res[0]['width']) })
          } else {
            // 距离左侧的距离+自身宽度一半-进度条左侧内边距
            const width = res[1][fLen - 1].left - parseInt(res[2].left) + res[1][fLen - 1].width / 2 - parseInt(res[0]['padding-left'])
            console.log(res[1][fLen - 1].left - parseInt(res[2].left))
            this.setData({ foregroundWidth: width / rpx })
          }
        })
      })
    }).catch(err => {
      console.log(err)
    })
  },
  /* 导购积分列表 */
  pointsList () {
    this.setData({ isLoading: true })
    pointsListApi().then(res => {
      this.setData({ isLoading: false })
      const { data, message, code } = res

      if (code !== 200) {
        wx.showToast(message)
        return
      }
      const { list = [] } = data

      const newList = list.map(i => ({
        ...i,
        created_at: toDate(i.created_at, 'yyyy/MM/dd hh:mm:ss').nowTime,
        icon: defaultTypeIcon.get(i.incentive_type) || defaultTypeIcon.get(2)
      }))

      this.setData({ pointsList: newList, 'pageInfo.total': data.total })
    }).catch(err => {
      this.setData({ isLoading: false })
      console.log(err)
    })
  },
  /* 获取积分信息 */
  getPointsRules () {
    const params = { rule_category_id: [...defaultPointsMap.keys()] }
    getPointsRulesApi(params).then(res => {
      const { data = [], message, code } = res

      if (code !== 200) {
        wx.showToast(message)
        return
      }

      if (!(data instanceof Array)) return

      const list = data.map(item => getTaskInfo(item))

      list.sort((a, b) => Number(Boolean(b.link)) - Number(Boolean(a.link)))

      this.setData({ settingList: list })
    }).catch(err => {
      console.log(err)
    })
  },
  /* 点击列表去完成 */
  goPage (e) {
    const { item } = e.currentTarget.dataset
    console.log(item)
    // 企掌道
    if (item.title === '完成welearning培训可得积分') {
      const timeStamps = new Date().getTime()
      getCoinsTokenApi({
        login_str: getRedirectPrarms('portal', timeStamps, true).params
      }).then(res => {
        if (res.code === 200) {
          const redirectUrl = `${getRedirectPrarms('portal', timeStamps, true).redirectUrl}&loginToken=${res.data.login_token}`
          console.log(redirectUrl)
          wx.navigateTo({ url: `/pages/webview/webview?url=${encodeURIComponent(redirectUrl)}` })
        }
      }).catch(err => {
        console.log(err)
      })
    // 发送邀约去顾客列表
    } else if (item.title === '发送到店邀约') {
      wx.navigateTo({ url: '/pagesA/customer-list/customer-list' })
    } else {
      wx.navigateTo({ url: item.link })
    }
  },
  /* 滑动到底部 */
  scrollToLower () {
    console.log('滑动到底部')
    // todo:增加金币列表
    const { isLoading, pageInfo, pointsList, coinsPageInfo, coinsList, listIndex } = this.data
    // 获取积分列表
    if (listIndex === 1 && !isLoading && pageInfo.total > pointsList.length) {
      this.setData({ 'pageInfo.page': pageInfo.page + 1 })
      this.pointsList()
    }
    // 获取金币列表
    if (listIndex === 2 && !isLoading && coinsPageInfo.total > coinsList.length) {
      this.setData({ 'coinsPageInfo.page': coinsPageInfo.page + 1 })
      this.pointsList()
    }
  },
  /* 显示足迹弹窗 */
  showPopupFootprint () {
    this.setData({ showPopupFootprint: true })
  },
  /** 获取等级足迹 */
  getLevelRecord () {
    getLevelRecordApi().then(res => {
      const { data, message, code } = res

      if (code !== 200) {
        wx.showToast(message)
        return
      }
      console.log(data)
      // 排序
      data.sort((a, b) => b.created_at - a.created_at)

      const list = data.map(i => ({
        ...i,
        created_at: toDate(i.created_at, 'yyyy年MM月dd日').nowTime
      }))

      this.setData({ levelRecord: [...list] })
    }).catch(err => {
      console.log(err)
    })
  },
  /** 展示更多 */
  loadMore (e) {
    const {name} = e.currentTarget.dataset
    this.setData({ [`scrollTop.${name}`]: 400 })
  },
  /** 滚动 */
  handleScroll (e) {
    const {scrollTop} = e.detail
    const {name} = e.currentTarget.dataset
    this.setData({
      [`showMore.${name}`]: false
    })

    if (name === 'task') {
      this.setData({
        offsetY: scrollTop + this.data.offset
      })
    }
  },
  /** todo:导购金币列表 */
  coinsList () {
    this.setData({ isLoading: true })
    getCoinsListApi().then(res => {
      this.setData({ isLoading: false })
      const { data, message, code } = res

      if (code !== 200) {
        wx.showToast(message)
        return
      }
      const { list = [] } = data

      const newList = list.map(i => ({
        ...i,
        created_at: toDate(i.created_at, 'yyyy/MM/dd hh:mm:ss').nowTime
      }))

      this.setData({ coinsList: newList, 'coinsPageInfo.total': data.total })
    }).catch(err => {
      this.setData({ isLoading: false })
      console.log(err)
    })
  }
})
