import { reserveListApi, reserveCheckListApi, whiteListApi, writeOffCouponApi } from '@/api/booking-api'
import { formatDate } from '@/utils/tools'

const { globalData } = getApp()

/** tabs */
const TABS = [
  { name: '全部预约客户', id: 1 },
  { name: '我核销的客户', id: 2 }
]

/** 设置分页数据 */
const generatePageInfo = (info = {}) => {
  return {
    page: 1,
    size: 10,
    hasMore: true,
    loading: false,
    ...info
  }
}

Page({
  data: {
    /** 切换tab */
    tabs: TABS,
    /** 当前选中tab, 默认选中第一个 */
    tabCheckedIndex: TABS[0].id,
    /** 搜索内容 */
    searchValue: '',
    /** 客户列表 */
    customerList: [],
    /** 分页数据，用于下拉加载 */
    pageInfo: generatePageInfo(),
    /** 总人数 */
    count: 0,
    /** 是否管理导购 */
    isAdmin: false,
    /** 判断是否适配 */
    isIPhoneX: globalData.isIPhoneX,

    activity_id: globalData.activity_id
  },

  onLoad () {
    // 获取编辑预约组件实例
    this.toastPopup = this.selectComponent('#toast')

    this.getUserIdentity()
    this.getCustomerList()
  },

  /** 清空搜索条件 */
  clearSearch () {
    console.log('触发删除')
    this.setData({
      searchValue: ''
    })
    this.handleSearch()
  },

  /** 获取用户身份 */
  getUserIdentity () {
    whiteListApi()
      .then(res => {
        // 错误处理
        if (res.code !== 200) return wx.showToast(res.message)
        // 0为不在白名单，1为在白名单
        const { ma } = res.data
        this.setData({
          isAdmin: ma === 1
        })
      })
      .catch(e => console.error(e, '用户身份获取出错'))
  },

  /** 上拉加载 */
  handleReachBottom () {
    console.log('上拉加载触发了')
    const { pageInfo } = this.data
    if (!pageInfo.hasMore) return
    pageInfo.page += 1
    this.getCustomerList()
  },

  /** 扫码核销 */
  scanCodeEvent () {
    const self = this
    wx.scanCode({
      // 只允许从相机扫码
      onlyFromCamera: false,
      success (res) {
        console.log('扫码成功：' + JSON.stringify(res))

        // 扫码成功后
        if (res.errMsg === 'scanCode:ok') {
          self.checkCode(res.result)
        }
      }
    })
  },

  /** 子组件核销 */
  writeOff (e) {
    wx.showModal({
      title: '提示',
      content: '确认核销？',
      success: (res) => {
        if (res.confirm) {
          this.checkCode(e.detail)
        } else {
          console.log('取消')
        }
      }
    })
  },

  /** 调用接口核销 */
  checkCode (info) {
    writeOffCouponApi({code: info, activity_id: this.data.activity_id})
      .then(res => {
        // 组合参数
        const message = res.code === 200 ? '核销成功' : res.message
        const type = res.code === 200 ? 'success' : 'error'
        // 设置弹窗参数
        this.toastPopup.openPopup(type, message, 3000)

        res.code === 200 && this.handleSearch()
      })
      .catch(e => console.error(e, '核销出错'))
  },

  /** 页面search change 变化 触发，判断是否展示删除按钮 */
  changeInput (e) {
    this.setData({
      searchValue: e.detail.value
    })
  },

  /** 页面search blur触发搜索 */
  handleSearch () {
    console.log('触发blur')
    // 初始化内容
    this.setData({
      pageInfo: generatePageInfo(),
      customerList: []
    })
    this.getCustomerList()
  },

  /** tab 切换, 重新赋值tabCheckedIndex */
  switchTab (e) {
    const { id } = e.currentTarget.dataset

    // 防止重复触发
    if (this.data.tabCheckedIndex === id) return

    // 初始化内容
    this.setData({
      tabCheckedIndex: id,
      pageInfo: generatePageInfo(),
      customerList: []
    })

    this.getCustomerList()
  },

  /** 获取预约客户列表 */
  getCustomerList () {
    const { pageInfo, customerList, searchValue, tabCheckedIndex } = this.data
    console.log('分页数据为', pageInfo)

    // 请求前改变加载状态
    this.setData({
      'pageInfo.loading': true
    })

    // 获取不同接口
    const api = tabCheckedIndex === 1 ? reserveListApi : reserveCheckListApi

    // 组合接口参数
    const params = {
      activity_id: this.data.activity_id,
      page: pageInfo.page,
      size: pageInfo.size,
      search: searchValue
    }

    // 接口请求
    api(params)
      .then((res) => {
        // 错误处理
        if (res.code !== 200) {
          this.setData({
            'pageInfo.loading': false
          })
          return wx.showToast(res.message)
        }

        // 数据处理
        const newCustomerList = res.data.lists.map(item => {
          return {
            ...item,
            reserveTime: formatDate(item.reserve_time * 1000, 'yyyy/MM/dd hh:mm:ss')
          }
        })
        // 请求成功处理
        this.setData({
          count: res.data.count,
          'pageInfo.loading': false,
          'pageInfo.hasMore': res.data.count > pageInfo.page * pageInfo.size,
          customerList: [...customerList, ...newCustomerList]

        })
      })
      .catch(e => console.error(e, '获取核销列表出错'))
  },

  /** 跳转至配置页面 */
  goSetupPage () {
    wx.navigateTo({
      url: '/pages/booking-setting/booking-setting'
    })
  }
})
