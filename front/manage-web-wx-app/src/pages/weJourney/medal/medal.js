import { medalWallApi, medalWallDetailApi } from '@/api/we-journey-api'
// 获取全局应用程序实例对象
const app = getApp()
const imagePath = 'https://cdn-cos.coach.com.cn/uploads/'

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    navigationBarHeight: app.globalData.system.statusBarHeight + 44 + 'px', // 状态栏高度
    windowWidth: wx.getSystemInfoSync().windowWidth, // 屏幕宽度
    imagePath: imagePath, // 图片公共地址
    medalInfoList: [], // 某一类奖章列表
    detailInfo: {}, // 某一类奖章详情
    showPopup: false, // 是否展示奖章弹窗
    position: 'center', // page-container position属性
    duration: 1000, // share-element动画时长
    show: true, // 是否展示某一类奖章详情
    overlay: false, // 是否显示page-container遮罩层
    current: -1, // 当前滑动index
    showShare: true, // 是否展示share-element
    transform: true, // 是否进行动画过度
    total: 0, // 总奖章数量
    gained: 0, // 已获得奖章数量
    animateInfo: {}, // 打开奖章动画信息
    ordinaryInfo: {
      list: [],
      total: 0,
      value: 0
    }, // 普通奖章
    rareInfo: {
      list: []
    }, // 稀有奖章
    isRare: false, // 是否是稀有奖章
    isLoadingDetail: false, // 是否正在获取奖章详情
    isScroll: false, // 页面是否滚动
    needSway: false, // 是否需要摇摆
    yearList: [], // 财年列表
    currentYear: 0, // 当前财年
    showMore: {
      ordinary: 9, // 普通奖章是否展示更多
      rare: 9 // 稀有奖章是否展示更多
    }
  },

  // 生命周期函数--监听页面加载
  onLoad () {
    this.setData({
      show: false
    })
    this.getMedalWall(true)
  },
  onPageScroll (e) {
    const { scrollTop } = e
    const { isScroll } = this.data

    if (scrollTop > 33 && !isScroll) {
      this.setData({ isScroll: true })
    } else if (scrollTop < 33) {
      this.setData({ isScroll: false })
    }
  },


  /* 打开弹窗 */
  openPopup (e) {
    const { isLoadingDetail, show, showPopup } = this.data
    // 如果正在加载再次点击无效
    console.log(isLoadingDetail, show, showPopup)
    if (isLoadingDetail || show || showPopup) return

    const { item: { is_show, rule_id, rule_gained, time_index }, type } = e.currentTarget.dataset
    this.setData({ isRare: type === 2 })

    this.medalWallDetail(is_show, rule_gained, rule_id, time_index)
  },
  /* 弹窗打开动画 */
  popupAnimation () {
    this.animationOne()
  },
  /* 背景色动画 */
  animationOne () {
    this.animate('.popup .background', [
      { opacity: 0 },
      { opacity: 0.8 }
    ], 500, () => {
      console.log('背景色')
      this.animationTwo()
    })
  },
  /* 底部动画展开 */
  animationTwo () {
    this.animate('.popup .footer', [
      { translateX: '-50%', scale: [1, 1] },
      { translateX: '-50%', scale: [7, 1] }
    ], 300, () => {
      console.log('放大')
    })
    // 动画中途执行另一个动画
    setTimeout(() => {
      this.animationThree()
    }, 200)
  },
  /* 奖章展示动画 */
  animationThree () {
    this.animate('.popup-medal', [
      { opacity: 0.2, scale: [0.5, 0.5], rotateY: 0 },
      { opacity: 1, scale: [1, 1], rotateY: 360 }
    ], 800, () => {
      this.clearAnimation('.popup-medal', { rotateY: true }, () => {
        console.log('奖章展示动画清除完毕')
        this.animate('.popup-medal', [
          { rotateY: 0 },
          { rotateY: -50 },
          { rotateY: 30 },
          { rotateY: 0 }
        ], 500, () => {
          console.log('抖动')
          setTimeout(() => {
            this.setData({ show: true })
          }, 100)
        })
      })
    })
    this.animate('.grow-box', [
      { height: '200px' },
      { height: '40%' }
    ], 800)
  },

  /* 获取奖章墙 */
  getMedalWall (isFirst = false) {
    medalWallApi().then(res => {
      console.log(res)
      const { code, message, data } = res

      if (code !== 200) {
        wx.showToast(message)
        return
      }

      const { total = 0, gained = 0, now_year, list } = data
      const yearList = (list || []).map((item) => {
        const ordinary = item.ordinary || []
        // 排序
        ordinary.sort((a, b) => {
          return b.get_time - a.get_time
        })
        // console.log('ordinary', ordinary.map(item => item.get_time))
        // 普通奖章数量
        item.total = ordinary.reduce((pre, nex) => pre + nex.rule_total, 0)
        item.value = ordinary.reduce((pre, nex) => pre + nex.rule_gained, 0)

        // 稀有奖章数量
        item.rare = item.rare || []
        if (item.rare.length) {
          const rare_total = item.rare.reduce((pre, nex) => pre + nex.rule_total, 0)
          const rare_gained = item.rare.reduce((pre, nex) => pre + nex.rule_gained, 0)
          item.total += rare_total
          item.value += rare_gained
        }
        // 获取名称
        if (item.year === now_year) {
          item.name = '当前财年'
        } else {
          item.name = `FY${(item.year + '').slice(-2)}`
        }
        return item
      })
      // 财年排序
      yearList.sort((a, b) => b.year - a.year)
      // 当前财年
      if (isFirst) {
        this.data.currentYear = yearList.findIndex(item => item.year === now_year)
      }


      this.setData({ total, gained, yearList, currentYear: this.data.currentYear })
    }).catch(err => {
      console.log(err)
    })
  },
  /* 关闭page-container */
  closePage () {
    console.log('closePage')
    this.setData({ show: false }, () => {
      console.log('closePage111', this.data.show)
    })
    this.getMedalWall()
  },
  /* 奖章详情 */
  medalWallDetail (isShow, ruleGained, id, time_index = 0) {
    const params = {
      rule_id: id,
      year: this.data.yearList[this.data.currentYear].year
    }
    this.setData({ isLoadingDetail: true })

    medalWallDetailApi(params).then(res => {
      const { code, data, message } = res
      this.setData({ isLoadingDetail: false })

      if (code !== 200) {
        wx.showToast(message)
        return
      }
      const { medal_style, ...detailInfo } = data

      this.setData({ medalInfoList: medal_style, detailInfo })

      // 已经全部查看或者已获得奖章数量为0 不展示查看动画
      if (isShow === 1 || ruleGained === 0) {
      // if (!isShow) {
        console.log('不展示', isShow)
        this.setData({ show: true, duration: 0, current: time_index, needSway: false })
      } else {
        // 排序
        const newMedalStyle = JSON.parse(JSON.stringify(medal_style))
        newMedalStyle.sort(function (a, b) {
          return a.create_time > b.create_time ? -1 : 1
        })

        const newIndex = newMedalStyle.length ? newMedalStyle[0].style_id : 0
        const medalIndex = medal_style.findIndex(item => item.style_id === newIndex)

        /* 最后一个已获取奖章 */
        const index = Math.max(0, medalIndex)
        // 模拟数据使用
        // const index = Math.min(parseInt(Math.random() * 5), medal_style.length - 1)
        console.log('index index', index)
        this.setData({duration: 1000, showPopup: true, needSway: true}, () => {
          this.setData({ current: index, animateInfo: medal_style[index] }, () => {
            this.popupAnimation()
          })
        })
      }
    }).catch(err => {
      console.log(err)
    })
  },
  /* current改变 */
  changeCurrent ({ value }) {
    this.setData({ current: value })
  },
  /* page-container打开进入之后触发 */
  onAfterEnter (res) {
    console.log(9999999999, res)
    this.setData({ showPopup: false, duration: 0 })
    this.data.needSway && this.afterLeave()
  },

  /* 动画结束后左右摇摆 */
  afterLeave () {
    this.animate('.swiper-item-checked', [
      { rotateY: 0 },
      { rotateY: -30 },
      { rotateY: 20 },
      { rotateY: 0 }
    ], 500)
  },

  /** 内嵌页面手动右滑返回 */
  onAfterLeave () {
    this.setData({
      show: false
    })
  },
  /**
   * @description 财年change
   * @param {*} e
   */
  changeYear (e) {
    const { index } = e.currentTarget.dataset
    this.setData({ currentYear: index })
  },
  /** 展示更多 */
  handleShowMore (e) {
    const { key } = e.currentTarget.dataset
    const { showMore, yearList, currentYear } = this.data
    const { ordinary, rare } = yearList[currentYear]
    const list = key === 'ordinary' ? ordinary : rare
    this.setData({
      [`showMore.${key}`]: showMore[key] === 9 ? list.length : 9
    })
  }
})
