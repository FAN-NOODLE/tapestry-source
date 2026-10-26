
import {storeListsApi, orderListApi, changeStatusApi} from '@/api/plan-task-api'
import { memberPositionApi } from '@/api/compass-api'
import { formatWithArray } from '@/utils/tools'

const { globalData } = getApp()

/** 是否购买 */
const SHOP_LIST = [
  {value: 1, name: '是'},
  {value: 2, name: '否'}
]

/** 顾客类型 */
const PEOPLE_LIST = [
  {value: 1, name: '新客'},
  {value: 2, name: '老客'}
]

/** 筛选状态 */
const STATUS_LIST = [
  { key: -1, value: '全部' },
  { key: 0, value: '进行中' },
  { key: 1, value: '邀约成功' },
  { key: 2, value: '邀约完成' },
  { key: 3, value: '已取消' },
  { key: 4, value: '已过期' }
]

/** 设置分页数据 */
const generatePageInfo = (info = {}) => {
  return {
    page: 1,
    size: 10,
    hasMore: true,
    loading: true,
    ...info
  }
}

Component({
  options: {
    addGlobalClass: true
  },
  properties: {
  },
  /** 页面的初始数据 */
  data: {
    /** 是否购买 */
    shopOptions: SHOP_LIST,
    /** 是否新/老客户 */
    newPeoples: PEOPLE_LIST,
    /** 客户名称搜索value */
    searchValue: '',
    /** 分页数据 */
    pageInfo: generatePageInfo(),
    /** 开始时间-日历中选中日期 */
    start_time: 0,
    /** 日历中选中展示状态  1 已完成 2 待完成（进行中） */
    calendarStatus: 2,
    /** 预约列表 */
    list: [],
    /** 状态picker */
    statusList: STATUS_LIST,
    /** 门店picker */
    storeList: [],
    /** 选中状态 index */
    checkedStatus: 1,
    /** 选中的门店 */
    checkedStore: -1,
    /** 筛选key */
    filterKeys: [],
    /** 判断是否企业环境 */
    wxwork: globalData.system.environment === 'wxwork',
    /** 搜索数据 */
    searchParams: {
      status: 0
    },
    /** 筛选框回显结果集合 */
    filterResult: {}
  },

  /** 生命周期函数--监听页面加载 */
  ready () {
    // 获取组件实例
    this.calendarCom = this.selectComponent('#calendar')
    // 获取实例
    this.filterSelection = this.selectComponent('#filterSelection')
    this.getMember()
  },
  methods: {
    /** 刷新日历 */
    refreshCalendar () {
      this.calendarCom.init(true)
    },

    /** 发送邀约评价给顾客 */
    sendApplets (e) {
      const appid = wx.getAccountInfoSync().miniProgram.appId
      const { id, reserve_number } = e.currentTarget.dataset.item
      let title = '邀约评价'
      let url = `pagesOrder/pages/store-invite/comment/comment.html?id=${id}&reserve_number=${reserve_number}&corp_id=${globalData.currentEnterprise.corp_id || ''}`

      console.log()
      wx.qy.shareToExternalContact({
        text: {
          content: title
        },
        attachments: [
          {
            msgtype: 'miniprogram', // 消息类型，必填
            miniprogram: {
              appid, //  三合一小程序的appid
              title: title, // 小程序消息的title
              imgUrl: 'https://cdn-cos.coach.com.cn/uploads/2022/5/23/60f86724d69e71c10979b4e5bb17fa23.png', // 小程序消息的封面图。
              page: url // 小程序消息打开后的路径，注意要以.html作为后缀，否则在微信端打开会提示找不到页面
            }
          }
        ],
        success: (res) => {
          console.log('成功', res)
          // TODO
        },
        error: (e) => {
          console.log('失败', e)
        }
      })
    },

    /** 跳转查看预约 */
    goInviteDetail (e) {
      let item = e.currentTarget.dataset.item
      wx.navigateTo({ url: `/pages/store-invite/invite-detail/invite-detail?id=${item.id}&reserveNumber=${item.reserve_number}&role=${this.data.role}` })
    },

    /** 立即沟通 */
    communicate (e) {
      let id = e.currentTarget.dataset.id
      if (!this.data.wxwork) {
        wx.showToast('联系功能仅支持企业微信登录时使用，请从企业微信打开WeClient Pro后再试')
        return
      }
      wx.qy.openEnterpriseChat({
        userIds: '',
        externalUserIds: id,
        groupName: '',
        success: function (res) {
          console.log(res, '成功')
        },
        fail: function (err) {
          console.log(err, '失败')
        },
        complete: function (complete) {
          console.log(complete, '完成')
        }
      })
    },
    /** 获取身份信息 */
    getMember () {
      memberPositionApi({is_new: 0}).then(async res => {
        if (res.code === 200) {
          // role 1 店员 2 店长 3 企业导购 4 社群导购 5 企业管理（目前只有3和5）
          // range 0 初始化值 1 无权限 2 门店权限 3 部门权限
          const {role} = res.data
          // 'card_type', 'clt_type', 暂时注释
          let filterKeys = []
          let storeList = []
          if (role === 5) {
            filterKeys = ['manage_type', 'member_list', 'card_type', 'clt_type', 'tag_list', 'impression_tag_list']
            // 获取门店列表
            storeList = await this.getStoreList()
          } else {
            filterKeys = ['manage_type', 'tag_list', 'card_type', 'clt_type', 'impression_tag_list']
          }
          this.setData({
            filterKeys,
            role,
            storeList
          })
        }
      }).catch(err => {
        console.log(err)
      })
    },

    /** 打开筛选 */
    handleFilter () {
      this.filterSelection.updatePopup()
    },

    /** 筛选确认 */
    confirm (e) {
      const { manage_type, impression_tag_list, tag_list, member_list, card_type, clt_type } = e.detail
      const searchParams = {
        ...this.data.searchParams,
        member_ids: member_list
      }

      const search = {
        // 顾问筛选
        manage_type: manage_type,
        // 客户标签
        tag_list: tag_list,
        // 印象标签
        impression_tag_list: impression_tag_list,
        // clt类型
        clt_type,
        // 会员卡等
        card_type,
        // 是否有search 0为否,1为是
        is_search: 0
      }

      // 响应后端接口，如果没有值就不传
      if (manage_type || impression_tag_list || tag_list || card_type || clt_type) {
        searchParams.search = {...search, is_search: 1}
      } else {
        searchParams.search = {is_search: 0}
      }

      this.setData({
        searchParams: {...searchParams},
        filterResult: e.detail
      })

      this.getList(true)
    },

    /** 切换类型、顾客类型/是否到店预约 */
    changeInput (e) {
      wx.showLoading({
        title: '加载中',
        mask: true
      })
      const {name, item} = e.currentTarget.dataset
      const value = Number(e.detail.value)
      // 是否到店
      let is_store = name === 'is_store' ? value : item.is_store
      // 是否新客户
      let customer_status = name === 'customer_status' ? value : item.customer_status

      // 组合参数
      const params = {
        id: item.id,
        reserve_number: item.reserve_number,
        is_store: Number(is_store),
        customer_status: Number(customer_status)
      }

      changeStatusApi(params)
        .then((res) => {
          wx.hideLoading()
          const {list} = this.data
          // 获取对应index
          const index = list.findIndex(listItem => listItem.id === item.id)

          if (res.code !== 200) {
            // 如果出错恢复原有值
            if (name === 'is_store') {
              is_store = item.is_store
            } else {
              customer_status = item.customer_status
            }
            wx.showToast(res.message)
          }

          // 不管请求成功还是失败都会重新setData
          list[index]['is_store'] = is_store
          list[index]['customer_status'] = customer_status

          this.setData({
            list
          })
        })
        .catch(e => {
          wx.hideLoading()
          console.error('切换类型出错', e)
        })
    },

    /** 重置状态 */
    reset () {
      let { checkedStatus, checkedStore, searchParams } = this.data
      // 状态默认进行中
      checkedStatus = 1
      checkedStore = -1
      searchParams.status = 0
      delete searchParams.store_ids
      this.setData({
        checkedStatus,
        checkedStore,
        searchParams,
        calendarStatus: 2
      })
      this.getList(true)
    },

    /** 选中元素变化了 */
    bindChange (e) {
      const index = Number(e.detail.value)
      const pickerId = e.target.dataset.id
      let attribute = ''

      // 获取日历中展示的状态
      let calendarStatus = this.data.checkedStatus
      // 组合搜索条件
      const searchParams = {}
      // 当提选择的picker 是状态
      if (pickerId === 0) {
        attribute = 'checkedStatus'
        const key = STATUS_LIST[index].key
        searchParams.status = key
        if ([1, 2].includes(key)) {
        // 预约计划状态1（邀约完成）,2（邀约成功） 对应日历中的1（任务完成）
          calendarStatus = 1
        } else if (key === 0) {
        // 预约计划状态0（进行中 对应日历中的2（进行中）
          calendarStatus = 2
        } else if ([3, 4].includes(key)) {
        // 预约计划状态3(已取消)、4（已过期 对应日历中的3（不做展示）
          calendarStatus = 3
        } else {
          calendarStatus = 0
        }
        // 神策数据上报 -- 顾客列表筛选
        wx.track.trackFilter({
          current_page: 3,
          filter_type_sv: 2,
          filter_option: {
            '邀约状态': STATUS_LIST[index].value
          }
        })
        wx.track.trackFilterDetails({
          current_page: 3,
          button_name: '确定',
          filter_type_sv: 2,
          tag_level_1: '',
          tag_level_2: '邀约状态',
          tag_value: STATUS_LIST[index].value
        })
      } else {
        attribute = 'checkedStore'
        const {storeList} = this.data
        searchParams.store_ids = [storeList[index].key]
        // 神策数据上报 -- 顾客列表筛选
        wx.track.trackFilter({
          current_page: 3,
          filter_type_sv: 5,
          filter_option: {
            '邀约门店': storeList[index].value
          }
        })
        wx.track.trackFilterDetails({
          current_page: 3,
          button_name: '确定',
          filter_type_sv: 5,
          tag_level_1: '',
          tag_level_2: '邀约门店',
          tag_value: storeList[index].value
        })
      }

      this.setData({
        [attribute]: index,
        calendarStatus,
        searchParams: {...this.data.searchParams, ...searchParams}
      })

      this.getList(true)
    },

    /** 获取门店列表-下拉 */
    getStoreList () {
      return new Promise((resolve) => {
        storeListsApi()
          .then(res => {
            if (res.code !== 200) {
              resolve([])
              return
            }

            const newData = res.data.map(item => ({
              value: item.store_name,
              key: item.store_id
            }))
            resolve(newData)
          }).catch(e => console.error('获取门店列表出错', e))
      })
    },

    /** 上拉加载 */
    bindDownLoad () {
      const { pageInfo } = this.data
      if (!pageInfo.hasMore) return
      pageInfo.page += 1
      this.getList()
    },

    /** 搜索 */
    handleSearch (e) {
      this.setData({
        searchValue: e.detail.value
      })
      this.isSearch = true
      this.getList(true)
    },

    /** 日期变化 */
    changeDay (e) {
      console.log('日期点击触发了')
      const timeNum = e.detail
      this.setData({
        start_time: timeNum
      })
      this.getList(true)
    },

    /**
     * 获取预约列表
     * @param { boolean } isClear 是否需要情空数据
     */
    getList (isClear = false) {
      let { pageInfo, start_time, searchParams, searchValue } = this.data
      console.log('接口调用了', searchParams)
      if (isClear) {
        // 重置分页数据
        pageInfo = generatePageInfo()
      } else {
        pageInfo.loading = true
      }

      // 请求前改变加载状态
      this.setData({
        pageInfo
      })

      const params = {
        page: pageInfo.page,
        size: pageInfo.size,
        start_time,
        // 后端是更距时间段来查询，86399 = 59:59
        end_time: start_time + 86399,
        name: searchValue,
        ...searchParams
        // is_search: searchValue ? 1 : searchParams.is_search || 0
      }

      wx.showLoading({
        title: '加载中'
      })
      orderListApi(params)
        .then(res => {
          wx.hideLoading()
          // 判断请求是否成功
          if (res.code !== 200) {
            this.setData({
              'pageInfo.loading': false
            })
            return wx.showToast(res.message)
          }
          const newList = res.data.list ? res.data.list.map(item => ({...item, reserve_order_price: item.reserve_order_price ? formatWithArray(item.reserve_order_price / 100) : ''})) : []

          this.setData({
            'pageInfo.loading': false,
            'pageInfo.hasMore': res.data.total > pageInfo.page * pageInfo.size,
            list: pageInfo.page === 1 ? (newList || []) : [...this.data.list, ...newList]
          })
          if (this.isSearch && searchValue) {
            // 神策数据上报 -- 顾客列表搜索
            wx.track.trackSearch({
              current_page: 2,
              search_type: 1,
              search_key_word: searchValue,
              result_number: res.data.total
            })
            this.isSearch = false
          }
        })
        .catch(e => {
          this.setData({
            'pageInfo.loading': false
          })
          console.error('播种计划列表获取出错', e)
        })
    }
  }
})

