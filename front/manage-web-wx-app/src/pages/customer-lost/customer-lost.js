import { getCusomerLostOption, getCusomerLostList } from '@/api/customer-api'
import { toDate } from '../../utils/tools'
// 获取全局应用程序实例对象
// const app = getApp();

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    time: 0, // 传入的时间
    showTime: '',
    title: '客户流失详情',
    active: 2,
    lost_total: 0, // 流失总人数
    selectStoreId: '0',
    selectStore: '全部',
    selectBAId: '0',
    selectBA: '全部',
    options1: [],
    options2: [],
    lost_count: 0,
    customerCardShow: true, // 判断暂无顾客是否展示
    ready: false,
    hasMore: false,
    isOnReachBottom: false,
    page: 1,
    customerList: [],
    member: {}
  },

  // 生命周期函数--监听页面加载
  onLoad (options) {
    // TODO: onLoad
    if (options) {
      this.data.time = Number(options.time)
      this.setData({
        showTime: toDate(this.data.time, 'yyyy年MM月dd日').nowTime
      })
    }
    this.getCusomerLostOption()
    this.getCusomerLostList()
  },

  changeActive (e) {
    const { index } = e.currentTarget.dataset
    if (index === this.data.active) return
    this.setData({
      active: index,
      page: 1,
      lost_total: 0,
      selectStoreId: '0',
      selectStore: '全部',
      selectBAId: '0',
      selectBA: '全部',
      options1: [],
      options2: []
    })
    this.getCusomerLostOption()
    this.getCusomerLostList()
  },

  changeStore (e) {
    const index = Number(e.detail.value)
    this.data.options2 = []
    const member = JSON.parse(JSON.stringify(this.data.member))
    if (this.data.options1[index].dp_id === '0' || 0) {
      for (const item in member) {
        member[item].forEach(childitem => {
          this.data.options2.push(childitem)
        })
      }
    } else {
      for (const key in member) {
        if (Number(key) === Number(this.data.options1[index].dp_id)) {
          this.data.options2 = member[Number(key)]
        }
      }
    }
    this.data.options2.unshift({user_id: '0', username: '全部'})
    this.setData({
      page: 1,
      selectBAId: '0',
      selectBA: '全部',
      options2: this.data.options2,
      selectStoreId: this.data.options1[index].dp_id,
      selectStore: this.data.options1[index].dp_name
    })
    this.getCusomerLostList()
  },

  changeBA (e) {
    const index = Number(e.detail.value)
    this.setData({
      page: 1,
      selectBAId: this.data.options2[index].user_id,
      selectBA: this.data.options2[index].username
    })
    this.getCusomerLostList()
  },

  // 获取流失客户统计和筛选数据
  getCusomerLostOption () {
    const params = {
      time: this.data.time,
      delete_type: this.data.active
    }
    getCusomerLostOption(params).then(res => {
      const { code, data } = res
      if (code === 200) {
        this.data.member = data.member
        this.data.options2 = []
        data.departments.unshift({dp_id: '0', dp_name: '全部'})
        this.data.options2.push({user_id: '0', username: '全部'})
        for (const item in this.data.member) {
          this.data.member[item].forEach(childitem => {
            this.data.options2.push(childitem)
          })
        }
        // 处理接口返回数据
        this.setData({
          options1: data.departments,
          lost_total: data.lost_total || 0,
          options2: this.data.options2,
          member: this.data.member
        })
      } else {
        wx.showToast(res.message)
      }
    }).catch(err => {
      console.log(err)
    })
  },

  // 获取客户流失列表
  getCusomerLostList () {
    const params = {
      time: this.data.time,
      delete_type: this.data.active,
      page: this.data.page,
      size: 10
    }
    if (this.data.selectStoreId !== '0') {
      params.dp_id = this.data.selectStoreId
    }
    if (this.data.selectBAId !== '0') {
      params.user_id = this.data.selectBAId
    }
    getCusomerLostList(params).then(res => {
      const { code, data } = res
      if (code === 200) {
        let customerList = []
        if (this.data.isOnReachBottom) {
          customerList = this.data.customerList.concat(data.lost_users)
          this.data.isOnReachBottom = false
        } else {
          customerList = data.lost_users
        }
        // 处理接口返回数据
        this.setData({
          lost_count: data.lost_count || 0,
          customerList: customerList,
          customerCardShow: !(customerList.length > 0),
          hasMore: customerList.length < data.lost_count,
          ready: true
        })
      } else {
        this.setData({
          isLoading: false,
          ready: true
        })
        wx.showToast(res.message)
      }
    }).catch(err => {
      console.log(err)
    })
  },
  // 上拉加载
  onReachBottom () {
    if (this.data.hasMore) {
      this.data.page += 1
      this.data.isOnReachBottom = true
      this.getCusomerLostList()
    }
  },
  // 页面分享设置
  onShareAppMessage () {
    return {
      title: '客户流失详情',
      path: 'pages/customer-lost/customer-lost'
    }
  }
})
