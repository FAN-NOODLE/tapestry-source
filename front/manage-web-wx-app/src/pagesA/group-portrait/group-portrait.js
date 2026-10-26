import { chatGroupMembers, chatIndexList } from '@/api/group-portrait-api'
import { toDate } from '../../utils/tools'

const TAB_LIST = [
  {label: '群成员列表', value: 0}
]

Page({
  data: {
    /** 群人员tab 切换 */
    tabList: TAB_LIST,
    /** 选中index */
    tabIndex: 0,
    /** 客户概述 */
    c_summary_list: [
      {
        name: '群成员总数',
        actions: '当前查询日期，该社群存在的群成员总数',
        prop: 'valid_uv'
      },
      {
        name: '当月新入群人数',
        actions: '当前查询日期内，该社群新入群人员数（人数去重）',
        prop: 'add_uv'
      },
      {
        name: '累计退群人数',
        actions: '群创建日起至查询日期该社群每日退群人数之和',
        prop: 'all_loss_uv'
      },
      {
        name: '退群率',
        actions: '当前查询日期，新增退群人数与退群前累计群留存成员成员数的占比：新增退群人数/（新增退群人数+群留存成员数）*100%',
        prop: 'loss_uv_fee'
      }
    ],
    /** 群统计 */
    g_summary_list: [
      {
        name: '群会员数',
        actions: '当前查询日期，群内留存用户中属于会员的总人数',
        prop: 'member_uv'
      },
      {
        name: '非会员数',
        actions: '当前查询日期，群内留存用户中属于非会员的总人数',
        prop: 'general_uv'
      },
      {
        name: '会员率',
        actions: '当前查询日期，群内会员人数在群留存成员中的占比：（会员总人数/群留存成员数（去除员工））*100%',
        prop: 'member_card_uv'
      }
    ],
    /** 群成员列表 */
    list: [],
    /** 活动参与排行榜 */
    activityList: [],
    /** 社群id  wr9hx_EQAAKjtYaKg0IS7WkbtXU26JLQ */
    chatId: ''
  },
  onLoad (option) {
    // 取消顶部返回按钮
    wx.hideHomeButton()
    this.setData({
      chatId: option.chatId
    })
    this.getMemberList()
  },
  onShow () {
    this.getChatInfo()
  },

  /** tab 切换 */
  tabChange (e) {
    this.setData({
      tabIndex: e.detail
    })

    this.getMemberList()
  },

  /** 群成员列表 */
  async getMemberList () {
    wx.showLoading({
      title: '加载中'
    })
    const params = {
      chat_id: this.data.chatId,
      page: 1,
      size: 5,
      friend: 0,
      list_type: this.data.tabIndex
    }
    const {
      code,
      message = '',
      data
    } = await chatGroupMembers(params)

    wx.hideLoading()
    if (code === 200) {
      if (data.items && data.items.length) {
        const attr = this.data.tabIndex ? 'activityList' : 'list'
        this.setData({
          [attr]: data.items
        })
      }
    } else {
      wx.showToast(message)
    }
  },

  /** 社群画像首页接口 */
  async getChatInfo () {
    wx.showLoading({
      title: '加载中'
    })
    const params = {
      chat_id: this.data.chatId,
      sts: Math.floor((new Date().getTime() - 86400000) / 1000),
      ets: Math.floor((new Date().getTime() + 1) / 1000)
    }
    const {
      code,
      message = '',
      data
    } = await chatIndexList(params)

    wx.hideLoading()
    if (code === 200) {
      data.chat_info.create_time = toDate(data.chat_info.create_time, 'yyyy/MM/dd hh:mm:ss').nowTime

      this.setData({
        info: data
      })
      this.initCharts(data['valid'] || [])
    } else {
      wx.showToast(message)
      this.initCharts([])
    }
  },

  /** 初始化统计图 */
  initCharts (data) {
    const bodyW = wx.getSystemInfoSync().windowWidth
    const options = {
      silent: true,
      legend: [
        {
          orient: 'horizontal',
          bottom: 24,
          icon: 'rect',
          itemGap: 15,
          x: 'center',
          itemWidth: 10,
          itemHeight: 10,
          itemStyle: {
            opacity: 1,
            fontSize: 12
          },
          data: ['黑钻卡', '钻卡', '白金卡', '金卡']
        },
        {
          orient: 'horizontal',
          bottom: 0,
          icon: 'rect',
          itemGap: 15,
          itemWidth: 10,
          itemHeight: 10,
          itemStyle: {
            opacity: 1,
            fontSize: 12
          },
          data: ['名媛', 'GC', '会员', '非会员']
        }
      ],
      color: [
        '#EBD95D',
        '#E19E57',
        '#D67B5E',
        '#CDB7CB',
        '#D7E9E9',
        '#6A99A7',
        '#116F6A'
      ],
      series: [
        {
          type: 'pie',
          radius: '45%',
          center: ['50%', '40%'],
          data,
          itemStyle: {
            opacity: 0.2
          },
          labelLine: {
            smooth: 0.2,
            length: bodyW / 20,
            length2: bodyW / 25,
            lineStyle: {
              type: 'dashed'
            }
          },
          label: {
            opacity: 1,
            show: true,
            align: 'left',
            fontSize: 14,
            color: '#33363C',
            lineHeight: 19,
            formatter: function (params) {
              return `${params.name}\n{x|${params.value}}人\n`
            },
            rich: {
              x: {
                color: '#000000',
                fontSize: 16
              }
            }
          }
        },
        {
          type: 'pie',
          radius: '40%',
          center: ['50%', '40%'],
          data,
          labelLine: false,
          label: {
            opacity: 1,
            show: true,
            position: 'inner',
            align: 'left',
            fontSize: 10,
            color: '#ffffff',
            lineHeight: 21,
            formatter: function (params) {
              const num = params.percent
              return `${Math.round(num) === num ? num : num.toFixed(1)}%`
            }
          }
        }
      ]
    }
    this.setData({
      options
    })
  },

  /** 跳转到群成员列表 */
  goPersons () {
    wx.navigateTo({
      url: `/pagesA/group-portrait/group-persons/group-persons?chatId=${this.data.chatId}`
    })
  }
})
