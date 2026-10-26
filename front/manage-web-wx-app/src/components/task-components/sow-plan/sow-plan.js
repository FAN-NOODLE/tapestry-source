
import { sowListApi, sowFinishedApi, userMemberRelationApi } from '@/api/plan-task-api'
import { toDate } from '@/utils/tools'

const { globalData } = getApp()

/** tabs */
const TAB_LIST = [
  {label: '待完成', value: 1, type: 2.1},
  {label: '已完成', value: 2, type: 2.2},
  {label: '全部', value: 0, type: 2.3}
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

Component({
  properties: {
    /** 消息通知进入默认选中天 */
    sowTime: {
      type: Number,
      value: 0
    }
  },
  /** 页面的初始数据 */
  data: {
    /** tabs 切换数据 */
    tabList: TAB_LIST,
    /** tabs 选中index */
    tabIndex: 1,
    /** 客户名称搜索value */
    searchValue: '',
    /** 分页数据 */
    pageInfo: generatePageInfo(),
    /** 开始时间-日历中选中日期 */
    start_time: 0,
    /** 播种列表 */
    list: [],
    /** 完成id */
    finishId: '',
    /** 是否打开二次确认弹窗 */
    showModel: false,
    /** 判断是否企业环境 */
    wxwork: globalData.system.environment === 'wxwork',
    /** 对应日历状态 */
    calendarStatus: 2,
    /** 搜索数据 */
    searchParams: {},
    /** 筛选框回显结果集合 */
    filterResult: {}
  },

  observers: {
    'sowTime': function (value) {
      if (value) {
        this.setData({
          start_time: value
        })
        this.refreshCalendar()
      }
    }
  },
  /** 生命周期函数--监听页面加载 */
  ready () {
    // 获取组件实例
    this.calendarCom = this.selectComponent('#calendar')
    // 获取实例
    this.filterSelection = this.selectComponent('#filterSelection')
  },
  methods: {
    /** 刷新日历 */
    refreshCalendar () {
      this.calendarCom && this.calendarCom.init(true)
    },

    /** 点击立即沟通 */
    handleCommunicate (e) {
      if (!this.data.wxwork) {
        wx.showToast('联系功能仅支持企业微信登录时使用，请从企业微信打开WeClient Pro后再试')
        return
      }
      let { open_id, member_id, work_external_user_id } = e.currentTarget.dataset.item

      const params = {
        member_ids: [member_id],
        work_external_user_ids: [work_external_user_id]
      }
      userMemberRelationApi(params)
        .then((res) => {
          if (res.code !== 200) return
          // relation_status 0-正常，1-员工删除联系人，2-联系人删除员工，3-双方互相删除，4-关系转移中'
          // 判断是否好友
          if (res.data.length && res.data[0].relation_status) {
            // 加好友
            this.addPerson(open_id)
          } else {
          // 打开聊天
            this.communicate(open_id)
          }
        })
        .catch(e => console.error('接口失败', e))
    },

    /** 打开立即沟通聊天 */
    communicate (id) {
      wx.qy.openEnterpriseChat({
        userIds: '',
        externalUserIds: id,
        groupName: '',
        success: function (res) {
          console.log(res, '成功')
        },
        fail: function (err) {
          console.log(err, '错误')
        }
      })
    },

    /** 添加好友 */
    addPerson (id) {
      wx.qy.openUserProfile({
        type: 2,
        userid: id,
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

    /** 筛选 */
    handleFilter () {
      this.filterSelection.updatePopup()
    },

    /** 筛选确认 */
    confirm (e) {
      // 避免为空查询
      if (JSON.stringify({}) === e.detail) return
      const { manage_type, impression_tag_list, tag_list, card_type, clt_type } = e.detail
      const other = {
        // 顾问筛选
        self_manager: manage_type,
        // 客户标签
        impression_tag_ids: impression_tag_list,
        // 印象标签
        work_external_tag_ids: tag_list,
        // clt类型
        clt_type,
        // 会员卡等
        card_type
      }

      this.setData({
        searchParams: other,
        filterResult: e.detail
      })
      this.getList(true)
    },

    /** 下拉加载 */
    bindDownLoad () {
      const { pageInfo } = this.data
      if (!pageInfo.hasMore) return
      pageInfo.page += 1
      this.getList()
    },

    /** tab 切换 */
    tabChange (e) {
      // 对应日历中的状态对应表
      const calendarStatusMap = [0, 2, 1]
      this.setData({
        tabIndex: e.detail,
        calendarStatus: calendarStatusMap[e.detail]
      })
      this.getList(true)
    },

    /** 搜索 */
    handleSearch (e) {
      this.setData({
        searchValue: e.detail.value
      })
      this.isSearch = true
      this.getList(true, {username: e.detail.value})
    },

    /** 日期变化 */
    changeDay (e) {
      const timeNum = e.detail
      this.setData({
        start_time: timeNum
      })
      this.getList(true)
    },

    /** 点击完成 */
    onFinish (e) {
      sowFinishedApi({id: this.data.finishId})
        .then((res) => {
          if (res.code !== 200) return wx.showToast(res.message)

          // 操作成功更改状态
          this.refreshCalendar()
          // 通知父组件更新播种计划右上角提示数量
          this.triggerEvent('sowRefresh')
          // 关闭弹窗
          this.closeModel()
          wx.showToast('操作成功')
          // 神策数据上报 -- 播种完成
          wx.track.trackReminderComplete({
            sow_task_id: this.data.finishId,
            sow_task_content: this.data.currentSow.content,
            user_external_id: this.data.currentSow.open_id
          })
        })
        .catch(e => console.error('播种计划点击完成', e))
    },

    /** 打开完成二次确认弹窗 */
    openModel (e) {
      const {id, item} = e.currentTarget.dataset
      this.setData({
        showModel: true,
        finishId: id,
        currentSow: item
      })
    },

    /** 点击关闭完成二次弹窗 */
    closeModel () {
      this.setData({
        showModel: false
      })
    },

    /**
     * 获取计划列表
     * @param { boolean } isClear 是否需要情空数据
     */
    getList (isClear = false, other = {}) {
      let { pageInfo, tabIndex, start_time, searchValue, searchParams } = this.data

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

      wx.showLoading({
        title: '加载中'
      })
      const params = {
        page: pageInfo.page,
        size: pageInfo.size,
        // 0 全部 1 未完成 2 已完成
        status: tabIndex,
        time: start_time,
        username: searchValue,
        ...searchParams
      }
      sowListApi(params)
        .then(res => {
          wx.hideLoading()
          // 判断请求是否成功
          if (res.code !== 200) {
            this.setData({
              'pageInfo.loading': false
            })
            return wx.showToast(res.message)
          }

          const newList = res.data.list ? res.data.list.map((item) => ({
            ...item,
            remind_time_str: toDate(item.remind_time, 'yyyy.MM.dd hh:mm:ss').nowTime
          })) : []

          this.setData({
            'pageInfo.loading': false,
            'pageInfo.hasMore': res.data.total > pageInfo.page * pageInfo.size,
            list: pageInfo.page === 1 ? newList : [...this.data.list, ...newList]
          })
          if (this.isSearch && searchValue) {
            // 神策数据上报 -- 顾客列表搜索
            wx.track.trackSearch({
              current_page: 1,
              search_type: 2,
              search_key_word: searchValue,
              result_number: res.data.total
            })
            this.isSearch = false
          }
        })
        .catch(e => console.error('播种计划列表获取出错', e))
    }
  }
})

