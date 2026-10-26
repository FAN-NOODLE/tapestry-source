// 获取全局应用程序实例对象
// const app = getApp();
import {
  communityListAPI
} from '../../api/compass-api'
import {toDate} from '../../utils/tools'

const {showToast, globalData} = getApp()
// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    wxwork: globalData.system.environment === 'wxwork',
    confirmModelShow: false,
    amount: 0,
    percentage: '0%',
    list: [],
    pageSet: {
      page: 0,
      size: 10,
      hasMore: true,
      loading: true
    },
    start_time: '',
    end_time: '',
    join_time: 1624527283
  },

  // 生命周期函数--监听页面加载
  onLoad (options) {
    this.data.start_time = options.start_time
    this.data.end_time = options.end_time
    this.setData({
      amount: options.amount,
      percentage: options.percentage
    })
    this.getCommunityList()
    // TODO: onLoad
  },

  // 关闭二次确认弹框
  closeConfirmModel () {
    this.setData({
      confirmModelShow: false
    })
  },

  scrolltolower () {
    if (!this.data.pageSet.loading && this.data.pageSet.hasMore) {
      this.setData({
        'pageSet.loading': true
      })
      this.getCommunityList()
    }
  },

  // 获取未参加活动列表
  async getCommunityList () {
    try {
      const res = await communityListAPI({
        page: ++this.data.pageSet.page,
        size: this.data.pageSet.size,
        start_time: this.data.start_time,
        end_time: this.data.end_time
      })
      if (res.code === 200) {
        res.data.list.forEach((item) => {
          item.join_time = toDate(item.join_time, 'yyyy年MM月dd日').nowTime
        })
        this.setData({
          list: [...this.data.list, ...res.data.list],
          'pageSet.loading': false,
          'pageSet.hasMore':
            res.data.total > this.data.pageSet.page * this.data.pageSet.size
        })
      } else {
        showToast(res.message)
        this.data.pageSet.loading = false
      }
      this.setData({ loading: false })
    } catch (error) {
      this.setData({
        'pageSet.loading': false,
        loading: false
      })
    }
  },

  // 立即沟通
  contact (e) {
    const { extid } = e.currentTarget.dataset
    if (this.data.wxwork) {
      wx.qy.openEnterpriseChat({
        // 注意：userIds和externalUserIds至少选填一个，且userIds+externalUserIds总数不能超过2000。
        userIds: '',
        externalUserIds: extid, // 参与会话的外部联系人列表，格式为userId1;userId2;…，用分号隔开。
        groupName: '', // 必填，会话名称。单聊时该参数传入空字符串""即可。
        success: function (res) {
          console.log('success', res)
          // 回调
        },
        fail: function (res) {
          console.log('fail', res)
          // 失败处理
        }
      })
    } else {
      this.setData({
        confirmModelShow: true,
        dialogText: '无法直接使用聊天功能',
        dialogTwoText: '请在企业微信中打开'
      })
    }
  }
})
