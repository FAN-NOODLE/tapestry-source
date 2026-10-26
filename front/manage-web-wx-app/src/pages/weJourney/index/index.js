import { indexApi, recordAlreadyGotApi, recordMemberReceivedApi, getPointsRulesApi, getCoinsTokenApi } from '@/api/we-journey-api'
import { toDate } from '@/utils/tools'
import { defaultPointsMap, getRedirectPrarms, getTaskInfo } from '../utils'

const imagePath = 'https://cdn-cos.coach.com.cn/uploads/'

const defaultDimension = [
  { subTitle: '邀请目标达成率', value: 'invite_value', ruleValue: 'invite_rule_value', time: 'invite_updated_at', name: '邀约', key: 1, ruleText: 'invite_rule_str' },
  { subTitle: '任务完成率', value: 'task_value', ruleValue: 'task_rule_value', time: 'task_updated_at', name: '任务', key: 2, ruleText: 'task_rule_str' },
  { subTitle: '总客户互动积分排行榜', value: 'interact_value', ruleValue: 'interact_rule_value', time: 'interact_updated_at', name: '互动', key: 3, ruleText: 'interact_rule_str' },
  { subTitle: '总培训积分排行榜', value: 'train_value', ruleValue: 'train_rule_value', time: 'train_updated_at', name: '培训', key: 4, ruleText: 'train_rule_str' },
  { subTitle: '目标达成率', value: 'online_value', ruleValue: 'online_rule_value', time: 'online_updated_at', name: '线上销售', key: 5, ruleText: 'online_rule_str' }
]
const defaultLevel = new Map([['SS', 0], ['S', 1], ['A+', 2], ['A', 3], ['A-', 4]])

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    isAnimation: false, // 是否进行动画
    imagePath: imagePath,
    pageInfo: {}, // 页面信息
    levelArr: [], // 能力图谱维度集合
    isLoading: true, // 是否正在加载
    medalNum: 0, // 已经领取奖章数量
    listLength: 0, // 未领取奖章数量
    newMedalInfo: {}, // 未领取奖章
    showMedalPopup: false, // 是否展示领取奖章弹窗
    isClosing: false, // 是否正在关闭领取奖章弹窗
    isShowList: false, // 是否展示下方列表
    isShowRedDot: false, // 是否展示红点
    /** 奖章列表 */
    medalList: []
  },

  // 生命周期函数--监听页面加载
  onLoad () {
  },

  onShow () {
    this.getIndexData()
    this.recordAlreadyGot()
    this.getPointsRules()
    wx.hideHomeButton()
  },
  onUnload () {
    clearTimeout(this.timer1)
    clearTimeout(this.timer2)
    clearTimeout(this.timer3)
    clearTimeout(this.timer4)
  },

  /** 奖章信息动画 */
  medalAnimation () {
    console.log('执行了执行了')
    // 已经进行动画
    this.setData({ isAnimation: true })
    // 模块放大
    this.animate('.medal', [
      { scale: [0, 0] },
      { scale: [1, 1] }
    ], 500, () => {
      // 文字展示
      this.animate('.word-box', [
        { opacity: 0.2 },
        { opacity: 1 }
      ], 500, () => {
        console.log('动画完成')
      })
    })

    this.timer4 = setTimeout(() => {
      // 背景展示
      this.animate('.bgc', [
        { opacity: 0.2 },
        { opacity: 1 }
      ], 500, () => {

      })
      // 进度条背景展示
      this.animate('.progress-background', [
        { width: '10%' },
        { width: '100%' }
      ], 500)
    }, 200)

    this.timer1 = setTimeout(() => {
      // 积分模块展示
      this.animate('.points', [
        { opacity: 0 },
        { opacity: 1 }
      ], 500, () => {

      })
    }, 350)

    this.timer2 = setTimeout(() => {
      // 进度条展示
      const { pageInfo: { max_level_value, point_value } } = this.data

      let progress = '0%'
      if (Number(point_value) === 0 || Number(max_level_value) === 0) {
        progress = '0%'
      } else {
        progress = Math.min(100, (Number(point_value) / Number(max_level_value) * 100).toFixed(2)) + '%'
      }

      this.animate('.progress-foreground', [
        { width: 0 },
        { width: progress }
      ], 500, () => {
        // 进度条头部icon展示
        this.animate('.progress-foreground > .icon', [
          { opacity: 0.2 },
          { opacity: 1 }
        ], 250, () => {
          // 底部展示
          this.animate('.footer', [
            { opacity: 0 },
            { opacity: 1 }
          ], 400, () => {
            console.log('全部完成')
          })
        })
      })

      // 执行canvas组建内部动画
      if (this.data.levelArr.length > 0) {
        this.timer3 = setTimeout(() => {
          this.selectComponent('#radar-chart').draw()
        }, 10)
      }
    }, 500)
  },
  /** 动画结束 */
  animateEnd () {
    if (this.data.isShowList) {
      this.animate('.list', [
        { opacity: 0 },
        { opacity: 1 }
      ], 1000)
    }
  },
  /* 跳转积分页面 */
  toPointsPage () {
    wx.navigateTo({ url: '/pages/weJourney/points/points' })
  },
  /* 跳转奖章页面 */
  goMedalPage () {
    wx.navigateTo({ url: '/pages/weJourney/medal/medal' })
  },
  /* 亮光特效 */
  lightAnimation () {
    this.animate('.light', [
      { top: '40%', opacity: 0.8, translateX: '10%', rotateZ: 30 },
      { top: '60%', opacity: 0.3, translateX: '-10%', rotateZ: 30 }
    ], 1000, () => {
      this.clearAnimation('.light', { opacity: true }, () => {
        this.animate('.red-dot', [
          // { opacity: 0 },
          { opacity: 1 },
          { opacity: 1, scale: [1.3, 1.3] },
          { opacity: 1, scale: [1, 1] },
          { opacity: 1, scale: [1.3, 1.3] }
        ], 2000)
      })
    })
  },
  /* 获取首页数据 */
  getIndexData () {
    indexApi().then(res => {
      this.setData({ isLoading: false })
      const { code, message, data } = res

      if (code !== 200) {
        wx.showToast(message)
        return
      }
      console.log('pageInfo', data)
      this.setData({ pageInfo: data })
      this.initData(data)
    }).catch(err => {
      this.setData({ isLoading: false })
      console.log(err)
    })
  },
  /* 初始化数据 */
  initData (data) {
    let filterDimension = []
    if (data.invite_rule && data.invite_rule.length > 0) {
      const list = [...new Set(data.invite_rule.map(i => i.incentive_config_type))]
      filterDimension = defaultDimension.filter(i => list.includes(i.key))
    }

    let levelArr = []

    // 是否有能力图谱维度
    if (filterDimension.length > 0) {
      const dimensionMap = (i, idx) => {
        const value = data[i.value] || '-'
        const time = data[i.time]

        const getInviteRule = rule => {
          const newRule = rule.sort(p => defaultLevel.get(value))

          const giveList = ['', '%', '%', '名', '名', '%']
          return newRule.map(j => {
            let condition = ''
            let value = ''

            switch (j.rule_condition) {
              case 1:
                condition = '达到'
                value = `TOP ${j.rule_value}名`
                break
              case 2:
                condition = '高于'
                value = `${j.rule_value}${giveList[j.incentive_config_type]}`
                break
              case 3:
                condition = '达到'
                value = `${j.rule_value}${giveList[j.incentive_config_type]}`
                break
              case 4:
                condition = '低于'
                value = `${j.rule_value}${giveList[j.incentive_config_type]}`
                break
              case 5:
                condition = 'Top榜占比'
                value = `${j.rule_value}%`
                break
              default:
                break
            }

            return { ...j, condition, value }
          })
        }

        return {
          id: idx,
          title: i.name,
          subTitle: i.subTitle,
          value,
          time: time ? toDate(time, 'yyyy.MM.dd').nowTime : '-',
          word: data[i.ruleText] || '-',
          inviteRule: getInviteRule(data['invite_rule'].filter(k => k['incentive_config_type'] === i.key)),
          index: defaultLevel.get(value)
        }
      }
      levelArr = filterDimension.map(dimensionMap)
    }

    console.log('level arr', levelArr)
    let isShowList = true

    if (levelArr.some(i => i.index === undefined)) {
      isShowList = false
    }

    this.setData({ levelArr, isShowList }, () => {
      !this.data.isAnimation && this.medalAnimation()
    })
  },
  /* 领取奖章数量和弹框 */
  recordAlreadyGot () {
    recordAlreadyGotApi().then(res => {
      const { code, message, data } = res

      if (code !== 200) {
        wx.showToast(message)
        return
      }

      const { medal_list, number, new_medal_sign } = data

      let info = {}

      const listLength = medal_list ? medal_list.length : 0

      if (listLength > 0) {
        info = medal_list[0]
        info.received_time = toDate(info.received_time, 'yyyy.MM.dd').nowTime
        this.setData({ showMedalPopup: true })
      }

      this.setData({ medalNum: number || 0, newMedalInfo: info, listLength, medalList: medal_list, isShowRedDot: new_medal_sign > 0 })
    }).catch(err => {
      console.log(err)
    })
  },
  /* 领取奖章 */
  drawMedal () {
    // 正在进行领取动画时不能继续领取
    if (this.data.isClosing) return

    const params = { record_id: [this.data.newMedalInfo.record_id] }
    recordMemberReceivedApi(params).then(res => {
      if (this.data.listLength === 1) {
        this.animateMedal()
      } else {
        // 重新定义info
        let info = {}
        // 删除第一项
        this.data.medalList.shift()
        info = this.data.medalList[0]
        info.received_time = toDate(info.received_time, 'yyyy.MM.dd').nowTime
        console.log(111111111, info)
        this.setData({ newMedalInfo: info, medalList: this.data.medalList, listLength: this.data.medalList.length })
      }
    }).catch(err => {
      console.log(err)
    })
  },
  /* 领取奖章动画 */
  animateMedal () {
    this.setData({ isClosing: true })
    let left = 0
    let top = 0

    const query = wx.createSelectorQuery()

    query.select('.medal-image').boundingClientRect()
    query.select('.medal').boundingClientRect()

    query.exec(res => {
      console.log('medal-image', res)
      left = (res[1].left - res[0].left) - 10
      top = (res[1].top - res[0].top) - 10

      this.animate('.popup', [
        { opacity: 1 },
        { opacity: 0 }
      ], 1000, () => {
        this.setData({ showMedalPopup: false, isClosing: false })
        this.lightAnimation()
      })

      this.animate('.medal-image', [
        { left: 0, top: 0, scale: [1, 1] },
        { left: left + 'px', top: top + 'px', scale: [0, 0] }
      ], 1000, () => {
        this.recordAlreadyGot()
      })
    })
  },
  /* 关闭领取奖章 */
  closeMedal () {
    this.setData({ showMedalPopup: false })
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
  }
})
