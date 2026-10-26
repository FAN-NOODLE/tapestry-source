import { toDate } from '../../utils/tools'
import {
  memberPositionApi,
  allDataApi,
  guideListApi
} from '../../api/compass-api'
import {
  getCurrentDate,
  getYesterDays,
  getSevenDays,
  getThirtyDays
} from '@/utils/get-time-period'
import { choiceEnterpriseApi } from '../../api/login-api'
import { customWhiteListApi } from '@/api/order'
// eslint-disable-next-line no-undef
// const { echarts } = requirePlugin('echarts')
const { globalData } = getApp()
// customerStatisticApi
// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    role: 0,
    range: 0,
    updateDate: '',
    tabbars: [
      {
        text: '绩效榜单',
        type: 3
      },
      {
        text: '好友数据',
        type: 2
      },
      {
        text: '社群数据',
        type: 1
      }
    ],
    tabIndex: 3,
    dateIndex: 2,
    selectDate: [
      getCurrentDate()[0].s_timestamp,
      getCurrentDate()[1].s_timestamp
    ],
    tipsIndex: '',

    worker_ids: [], // 导购id
    worker_texts: [], // 导购名称
    selectDepartment: [], // 已选部门、负责人
    department_ids: [], // 已选门店id集合
    selectStore: [], // 已选门店、负责人
    store_ids: [],

    storeUserIds: [], // mem_userid
    store_worker_ids: [],
    baworker_ids: [],

    community_data: {
      cumulative_group_number: 0, // 累计群成员数
      leave_group_percentage: '0%', // 退群率
      number_of_newcomers: 0, // 新入群人数
      wx_member_binding_percentage: '0%', // 微信会员绑定率
      group_number: 0, // 群数量
      wx_member_binding_number: 0, // 微信会员绑定人数
      add_friend_percentage: '0%', // 加好友率
      add_friend_number: 0, // 加好友人数
      leave_group_number: 0 // 退群人数
    }, // 社群数据
    community_activity_data: {
      activity_group_number: 0, // 活动群数
      participate_number: 0, // 群活动参与人数
      participate_percentage: '0%', // 群活动参与率
      into_group_number: 0, // 活动进群人数
      into_group_percentage: '0%' // 活动进群率
    }, // 社群活动数据
    buddy_data: {
      new_add_buddy_number: 0, // 新增好友数
      grand_total_number: 0, // 累计好友数
      buddy_into_number: 0, // 好友入群数
      buddy_into_percentage: '0%', // 好友入群率
      buddy_churn_number: 0, // 好友流失人数
      buddy_churn_percentage: '0%', // 好友流失率
      new_buddy_member_binding_number: '0%', // 新增好友会员绑定率
      new_buddy_into_number: '0%', // 新增好友入群率
      ba_number: 0, // ba 导购人数
      wx_member_binding_number: 0, // 微信会员绑定人数
      wx_member_binding_percentage: '0%' // 微信会员绑定率
    }, // 好友数据

    storeName: '',
    baName: '',

    dateVal: {
      value: '2',
      data: [getCurrentDate()[0].s_timestamp, getCurrentDate()[1].s_timestamp],
      label: '今天'
    },
    dateVals: [],
    dateList: [
      {
        value: '1',
        label: '昨天'
      },
      {
        value: '2',
        label: '今天',
        checked: true
      },
      {
        value: '3',
        label: '近7天'
      },
      {
        value: '4',
        label: '近30天'
      },
      {
        value: '5',
        label: '自定义时间'
      }
    ],
    shadeShow: true,
    selDate: false,
    shadeShowNew: false,
    firstShow: true
  },
  // 生命周期函数--监听页面加载
  onLoad (options) {
    const updateDate = toDate(
      getCurrentDate()[0].s_timestamp,
      'yyyy.MM.dd'
    ).nowTime
    this.setData({
      updateDate: updateDate
    })
    this.getCustomWhiteList()
  },
  async onShow () {
    !this.data.role && this.getMember()
    if (!globalData.choiceEnterprise) {
      await this.getEnterInfo()
    }
    !this.data.storeName &&
      this.setData({
        baworker_ids: [globalData.choiceEnterprise.mem_uid],
        storeName: globalData.choiceEnterprise.department_name,
        baName: globalData.choiceEnterprise.name
      })
    this.getData()
    wx.showTabBar()
  },
  onUnload () {
    globalData.selectDate = ''
    globalData.selectDate1 = ''
    globalData.selectGuide = ''
    globalData.selectDepartment = ''
    globalData.selectStore = ''
  },
  getCustomWhiteList() {
    customWhiteListApi({
      type: 5
    }).then(res => {
      if (res.code === 200) {
        if (res.data.show === 1) {
          this.setData({
            tabIndex: 3,
            tabbars: [
              {
                text: '绩效榜单',
                type: 3
              },
              {
                text: '好友数据',
                type: 2
              },
              {
                text: '社群数据',
                type: 1
              }
            ]
          })
        } else {
          this.setData({
            tabIndex: 2,
            tabbars: [
              {
                text: '好友数据',
                type: 2
              },
              {
                text: '社群数据',
                type: 1
              }
            ]
          })
        }
      }
    }).catch(err => {
      console.log(err)
    })
  },
  // 排行
  goRanking () {
    wx.redirectTo({
      url: `/pages/ranking-list/ranking-list?type=${this.data.tabIndex}`
    })
  },
  getData () {
    if (globalData.selectDate && globalData.selectDate.length > 0) {
      this.setData({
        dateVal: {
          value: '5',
          label: '自定义',
          showData:
            globalData.selectDate[0].time + '-' + globalData.selectDate[1].time,
          data: [
            globalData.selectDate[0].s_timestamp,
            globalData.selectDate[1].s_timestamp
          ]
        }
      })
    }
    if (globalData.selectDate1 && globalData.selectDate1.length > 0) {
      this.setData({
        dateVals: [
          globalData.selectDate1[0].s_timestamp,
          globalData.selectDate1[1].s_timestamp,
          globalData.selectDate1[0].time + '-' + globalData.selectDate1[1].time
        ]
      })
    }

    const selectGuide = globalData.selectGuide || []
    let worker_ids = selectGuide.map((item) => {
      return item.mem_uid
    })
    this.setData({
      worker_ids: worker_ids
    })
    if (this.data.range === 3) {
      const selectDepartment = globalData.selectDepartment || []
      this.setData({
        selectDepartment: globalData.selectDepartment || [],
        department_ids: selectDepartment.map((item) => {
          return item.third_id
        })
      })
      const a = this.data.department_ids.slice(0) //
      this.getList(a)
    } else {
      const selectStore = globalData.selectStore || []
      this.setData({
        selectStore: globalData.selectStore || [],
        store_ids: selectStore.map((item) => {
          return item.third_id
        })
      })
      let a = this.data.store_ids.slice(0)
      this.getList(a)
    }
  },
  async getEnterInfo () {
    let params = {
      login_cache_key: globalData.login_cache_key,
      enterprise_id:
        globalData.enterprise[0] && globalData.enterprise[0].enterprise_id
    }
    const data = await choiceEnterpriseApi(params)
    globalData.choiceEnterprise = data.data
  },
  // 打开日期选项
  openDateSel () {
    this.setData({
      selDate: !this.data.selDate
    })
  },
  // 获取身份信息
  getMember () {
    memberPositionApi({ is_new: 0 })
      .then((res) => {
        if (res.code === 200) {
          // role 1 店员 2 店长 3 企业导购 4 社群导购 5 企业管理（目前只有3和5）
          // range 0 初始化值 1 无权限 2 门店权限 3 部门权限
          // if (roleIndex === 1 || roleIndex === 3) {
          //   role = 1
          // } else if (roleIndex === 2) {
          //   role = 2
          // } else if (roleIndex === 5 || roleIndex === 4) {
          //   role = 3
          // }
          this.setData({
            role: res.data.role, // 改回来
            range: res.data.range
          })
        }
      })
      .catch((err) => {
        console.log(err)
      })
  },
  // 获取导购列表
  getList (type) {
    const params = {
      dp_ids: type || []
    }
    guideListApi(params)
      .then((res) => {
        if (res.code === 200) {
          if (res.data.list.length > 0) {
            let store_worker_ids = res.data.list.map((item) => {
              return item.mem_uid
            })
            let storeUserIds = res.data.list.map((item) => {
              return item.mem_userid
            })
            this.setData({
              store_worker_ids: store_worker_ids,
              storeUserIds: storeUserIds
            })
          } else {
            this.setData({
              store_worker_ids: [''],
              storeUserIds: ['']
            })
          }
          this.getAllData()
        } else {
          wx.showToast(res.data.message)
        }
      }).catch(err => {
        console.log(err)
      })
  },
  getGuide () {
    const { role, worker_ids, store_worker_ids } = this.data
    // role 1 店员 2 店长 3 企业导购 5 企业管理 (1，3只获取自己的数据)
    if (role === 1 || role === 3) {
      return [globalData.choiceEnterprise.mem_uid]
    }
    // 选择员工搜索
    if (globalData.selectGuide) {
      return worker_ids
    }
    // 选择部门
    if (globalData.selectDepartment || globalData.selectStore) {
      return store_worker_ids
    }
    return worker_ids.length === 0 ? store_worker_ids : worker_ids
  },

  /** 获取userIds */
  getUserIds () {
    const { role, storeUserIds } = this.data
    // role 1 店员 2 店长 3 企业导购 5 企业管理 (1，3只获取自己的数据)
    let worker_ids = []

    if (role === 1 || role === 3) {
      return [globalData.choiceEnterprise.mem_userid]
    }
    // 选择员工搜索
    if (globalData.selectGuide) {
      const selectGuide = globalData.selectGuide || []
      let worker_ids = selectGuide.map((item) => {
        return item.mem_userid
      })
      return worker_ids
    }
    // 选择部门
    if (globalData.selectDepartment || globalData.selectStore) {
      return storeUserIds
    }
    return worker_ids.length === 0 ? storeUserIds : worker_ids
  },
  // 获取数据
  getAllData (type) {
    let params = {
      shopping_guide: this.getGuide(),
      userids: this.getUserIds(),
      start_time: this.data.dateVal.data[0],
      end_time: this.data.dateVal.data[1],
      search_all: 0
    }
    const storeList = this.data.range === 3
      ? this.data.department_ids
      : this.data.store_ids
    if (
      this.data.role === 5 &&
      storeList.length === 0 &&
      this.data.worker_ids.length === 0
    ) {
      params.shopping_guide = ['']
      params.search_all = 1
      delete params.userids
    }
    console.log(params, '这是id')
    // eslint-disable-next-line no-mixed-operators
    // if (globalData.selectGuide) {
    //   params.shopping_guide = this.data.worker_ids
    // } else if ((globalData.selectDepartment || globalData.selectStore) && !globalData.selectGuide) {
    //   params.shopping_guide = this.data.store_worker_ids
    // } else if (this.data.role === 1 || this.data.role === 3) {
    //   params.shopping_guide = [globalData.choiceEnterprise.mem_uid]
    // } else {
    //   params.shopping_guide = this.data.worker_ids
    // }
    allDataApi(params)
      .then((res) => {
        if (res.code === 200) {
          this.setData({
            shadeShow: false,
            community_data: res.data.community_data,
            community_activity_data: res.data.community_activity_data,
            buddy_data: res.data.buddy_data
          })
          console.log(this.data.shadeShow)
        }
      })
      .catch((err) => {
        this.setData({
          shadeShow: false
        })
        console.log(err)
      })
  },
  // tab切换
  switchTab (e) {
    // console.log(e.target.dataset)
    const { index } = e.target.dataset
    this.setData({
      tabIndex: index,
      selDate: false
    })
    if (globalData.selectDate1 == null) {
      this.setData({
        dateVals: []
      })
    }
  },
  // 时间切换
  switchDate (e) {
    // console.log(e.target.dataset)
    const { type } = e.target.dataset
    this.setData({
      dateIndex: type
    })
    switch (type) {
      case 1:
        this.setData({
          shadeShow: true,
          selectDate: [
            getYesterDays()[0].s_timestamp,
            getYesterDays()[1].s_timestamp
          ]
        })
        globalData.selectDate = ''
        this.getAllData()
        break
      case 2:
        this.setData({
          shadeShow: true,
          selectDate: [
            getCurrentDate()[0].s_timestamp,
            getCurrentDate()[1].s_timestamp
          ]
        })
        globalData.selectDate = ''
        this.getAllData()
        break
      case 3:
        this.setData({
          shadeShow: true,
          selectDate: [
            getSevenDays()[0].s_timestamp,
            getSevenDays()[1].s_timestamp
          ]
        })
        globalData.selectDate = ''
        this.getAllData()
        break
      case 4:
        this.setData({
          shadeShow: true,
          selectDate: [
            getThirtyDays()[0].s_timestamp,
            getThirtyDays()[1].s_timestamp
          ]
        })
        globalData.selectDate = ''
        this.getAllData()
        break
      case 5:
        const showDate = true // 时间选择器是否具体到日
        if (globalData.selectDate && globalData.selectDate.length) {
          wx.navigateTo({
            url: `/pages/date-guide-store-select/date-guide-store-select?type=3&startTime=${globalData.selectDate[0].s_timestamp}&endTime=${globalData.selectDate[1].s_timestamp}&showDate=${showDate}`
          })
        } else {
          wx.navigateTo({
            url: `/pages/date-guide-store-select/date-guide-store-select?type=3&showDate=${showDate}`
          })
        }
        break
      default:
        break
    }
  },
  // 时间切换
  dateChecked (e) {
    const index = parseInt(e.currentTarget.dataset.index)

    this.data.dateList.forEach((item) => {
      item.checked = false
      item.data = []
    })

    if (index === 5) {
      // 兼容BA数据 时间选项只有两项
      this.data.dateList[index - 4]['checked'] = true
    } else {
      this.data.dateList[index - 1]['checked'] = true
    }

    const list = this.data.dateList.slice(0)
    if (index === 1) {
      this.setData({
        selDate: false,
        dateList: list,
        shadeShow: true,
        dateVal: {
          ...list[index - 1],
          data: [
            getYesterDays()[0].s_timestamp,
            getYesterDays()[1].s_timestamp
          ]
        }
      })
      globalData.selectDate = ''
      this.getAllData()
    }
    if (index === 2) {
      this.setData({
        selDate: false,
        dateList: list,
        shadeShow: true,
        dateVal: {
          ...list[index - 1],
          data: [
            getCurrentDate()[0].s_timestamp,
            getCurrentDate()[1].s_timestamp
          ]
        }
      })
      globalData.selectDate = ''
      this.getAllData()
    }
    if (index === 3) {
      this.setData({
        selDate: false,
        dateList: list,
        shadeShow: true,
        dateVal: {
          ...list[index - 1],
          data: [getSevenDays()[0].s_timestamp, getSevenDays()[1].s_timestamp]
        }
      })
      globalData.selectDate = ''
      this.getAllData()
    }
    if (index === 4) {
      this.setData({
        selDate: false,
        dateList: list,
        shadeShow: true,
        dateVal: {
          ...list[index - 1],
          data: [
            getThirtyDays()[0].s_timestamp,
            getThirtyDays()[1].s_timestamp
          ]
        }
      })
      globalData.selectDate = ''
      this.getAllData()
    }
    if (index === 5) {
      this.setData({
        selDate: false
      })
      const showDate = true // 时间选择器是否具体到日
      if (globalData.selectDate && globalData.selectDate.length) {
        wx.navigateTo({
          url: `/pages/date-guide-store-select/date-guide-store-select?type=3&startTime=${globalData.selectDate[0].s_timestamp}&endTime=${globalData.selectDate[1].s_timestamp}&showDate=${showDate}`
        })
      } else {
        wx.navigateTo({
          url: `/pages/date-guide-store-select/date-guide-store-select?type=3&showDate=${showDate}`
        })
      }
    }
    const currentPage = 2 - this.data.tabIndex
    // 神策数据上报 -- 顾客列表筛选
    wx.track.trackFilter({
      current_page: currentPage,
      filter_type_sv: 0,
      filter_option: {
        '时间': this.data.dateList[index - 1].label
      }
    })
    // 神策数据上报 -- 筛选详情
    wx.track.trackFilterDetails({
      current_page: currentPage,
      button_name: '确定',
      filter_type_sv: 0,
      tag_level_1: '',
      tag_level_2: '时间',
      tag_value: this.data.dateList[index - 1].label
    })
  },
  // 点击显示字段说明文案
  handleTipExpress (e) {
    console.log(e.currentTarget.dataset.index)
    this.setData({
      tipsIndex: Number(e.currentTarget.dataset.index),
      shadeShowNew: true
    })
  },
  // 隐藏字段说明文案
  handleHideTipExpress () {
    this.setData({
      tipsIndex: '',
      shadeShowNew: false
    })
  },
  // 选择门店
  openOrgaSel () {
    this.setData({
      selDate: false
    })
    const b =
      this.data.range === 3
        ? this.data.department_ids.slice(0)
        : this.data.store_ids.slice(0)

    const mautiple = false // 门店选择是否单选
    wx.navigateTo({
      url: `/pages/date-guide-store-select/date-guide-store-select?type=5&hasIds=${JSON.stringify(
        b
      )}&mautiple=${mautiple}`
    })
  },
  // 选择导购
  openGuideSel () {
    this.setData({
      selDate: false
    })
    const isBa = false
    if (this.data.range === 3) {
      const a = this.data.department_ids.slice(0) //
      const b = this.data.worker_ids.slice(0) //
      wx.navigateTo({
        url: `/pages/date-guide-store-select/date-guide-store-select?type=2&dpIds=${JSON.stringify(
          a
        )}&hasChooseGuideIds=${JSON.stringify(b)}&role=3&isBa=${isBa}`
      })
    } else {
      let a = this.data.store_ids.slice(0)
      const b = this.data.worker_ids.slice(0)
      wx.navigateTo({
        url: `/pages/date-guide-store-select/date-guide-store-select?type=2&storeIds=${JSON.stringify(
          a
        )}&hasChooseGuideIds=${JSON.stringify(b)}&role=2&isBa=${isBa}`
      })
    }
  }
})
