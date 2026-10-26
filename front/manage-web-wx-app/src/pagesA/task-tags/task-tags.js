import {
  getCustomerApi,
  leaveGroupApi,
  reportSendMsgApi,
  reportCallApi,
  reportAddFriendApi,
  tagCustomerApi
} from '@/api/plan-task-api'
import { getCustomerPhoneApi } from '@/api/customer-api'
import Retry from '@/lib/retry'
import { throttle } from '@/utils/tools'
// 创建页面实例对象
const { globalData } = getApp()
Page({
  // 页面的初始数据
  data: {
    /** 首次弹出 建议联系客户后再标记完成哦~ */
    // firstLoading: false,
    /** 请先拨打电话 */
    showComPopup: false,
    /** 是否正在查询手机号 */
    isLoadingPhone: false,
    /** 是否正在标记 */
    isTagLoading: false,
    /** 填写说明弹窗 */
    remarkPopup: false,
    /** 说明弹窗内容  */
    remarkContent: {
      remark: '',
      id: '',
      index: -1
    },
    /** 子任务id */
    id: '',
    wxwork: globalData.system.environment === 'wxwork',
    tabList: [
      {
        type: 0,
        text: '全部'
      },
      {
        type: 1,
        text: '好友'
      },
      {
        type: 2,
        text: '非好友'
      }
    ],
    /** 用户名单 */
    userList: [
    ],
    /** 0待完成,1已完成,2已过期 */
    status: 0,
    isLoading: false,
    /** 是否还有分页数据 */
    noMore: false,
    /** 分页信息 */
    meta: {
      pageNo: 1,
      pageSize: 50,
      totalPage: 0
    },

    /** 输入框关键字 */
    keywords: '',
    /** 查询顾客类型  0全部 1已联系 2待联系 */
    filterType: 0,
    filterList: [
      {
        label: '全部顾客',
        value: 0
      },
      {
        label: '已联系顾客',
        value: 1
      },
      {
        label: '待联系顾客',
        value: 2
      }
    ],
    /** 是否展示 下拉选项 */
    showType: false,
    /** todo:增加参数 */
    category: 0,
    group_send_type: 0,
    dimension: 0
  },
  /** 打电话: 异步查询手机号码  */
  handleCall: throttle(async function (e) {
    wx.showLoading({ title: '', mask: true })

    if (this.data.isLoadingPhone) {
      return
    }
    this.data.isLoadingPhone = true

    const {
      item,
      index
    } = e.currentTarget.dataset
    const that = this
    // 上报打电话行为
    that.tagCommunicate(index, 'is_call')
    // TODO: 换个接口
    const params = {
      unionid: item.union_id,
      customer_code: item.customer_id
    }
    try {
      const { code, data, message } = await getCustomerPhoneApi(params)
      wx.hideLoading()
      this.data.isLoadingPhone = false
      if (code !== 200) {
        wx.showToast(message)
        return
      }
      let { mobile } = data

      if (mobile) {
        // 对某个人点击了打电话，则允许对该用户标记

        mobile = mobile.replace(/\s*/g, '')
        wx.makePhoneCall({
          phoneNumber: mobile,
          success: function (res) {
            console.log(res, '成功')
            // 神策数据上报 -- 拨打顾客电话
            wx.track.trackCallCustomer({
              user_external_id: item.openid
            })
          }
        })
      } else {
        wx.showToast('该用户暂无手机号')
      }
    } catch (error) {
      wx.hideLoading()
      this.data.isLoadingPhone = false
      console.log(error)
    }
  }, 2000),
  /** 关闭“电话联系客户后才能添加好友哦”我知道了 */
  closeComPopup() {
    this.setData({
      showComPopup: false
    })
  },
  /** 添加好友 */
  addUnion: throttle(async function (e) {
    let {
      item: { work_external_user_id, is_call },
      index
    } = e.currentTarget.dataset

    if (!this.data.wxwork) {
      wx.showToast(
        '联系功能仅支持企业微信登录时使用，请从企业微信打开WeClient Pro后再试'
      )
      return
    }
    // 对当前用户未执行打电话，直接点击添加好友，需阻止(退群名单例外)
    if (is_call === 0 && this.data.pageType !== 1) {
      this.setData({
        showComPopup: true
      })
      return
    }
    // 上报添加好友行为
    this.tagCommunicate(index, 'is_friend')
    wx.qy.openUserProfile({
      type: 2, // 1表示该userid是企业成员，2表示该userid是外部联系人
      userid: work_external_user_id, // 可以是企业成员，也可以是外部联系人
      success: function (res) {
        // 回调
        console.log(res, '成功')
      }
    })
  }, 2000),
  /** 发消息 */
  communicate: throttle(async function (e) {
    let {
      item: { work_external_user_id },
      index
    } = e.currentTarget.dataset
    const that = this
    if (!this.data.wxwork) {
      wx.showToast(
        '联系功能仅支持企业微信登录时使用，请从企业微信打开WeClient Pro后再试'
      )
      return
    }
    // 上报发消息行为
    if (that.data.pageType !== 1) {
      that.tagCommunicate(index, 'is_contact')
    }
    // 对某个人点击了发消息
    wx.qy.openEnterpriseChat({
      userIds: '',
      externalUserIds: work_external_user_id,
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
  }, 2000),
  /** 打电话或者发消息后需要更新顾客的 is_contact /is_call 状态 ，coach需要自动标记已联系 */
  async tagCommunicate(customerIndex, key) {
    let reportCommunicationApi = ''
    if (key === 'is_contact') {
      reportCommunicationApi = reportSendMsgApi
    } else if (key === 'is_call') {
      reportCommunicationApi = reportCallApi
    } else {
      reportCommunicationApi = reportAddFriendApi
    }
    const retryRequest = new Retry(reportCommunicationApi)

    let userInfo = this.data.userList[customerIndex]
    // 标记完成
    let params = {
      task_allot_log_id: userInfo.task_allot_log_id,
      union_id: userInfo.union_id,
      customer_id: userInfo.customer_id,
      status: 1// 导购执行了拨打，或者联系顾客
    }

    try {
      const { code, message } = await retryRequest.request(params)

      if (code !== 200) {
        wx.showToast(message)
        return
      }
      if (key === 'is_friend') {
        // 添加好友，无需执行编辑完成和状态操作
        return
      }
      // 维护communication状态
      let oneUserList = `userList[${customerIndex}]`
      let upStatus = {
        [key]: 1
      }
      // 子任务状态 关联 标记联系 0 待完成,1 已完成,2 已过期
      if (this.data.status !== 2) {
        upStatus.is_finished = 1// coach品牌对顾客进行打电话或者发消息，顾客状态自动标记为【已联系】
      }
      this.setData({
        [oneUserList]: {
          ...userInfo,
          ...upStatus
        }
      })
    } catch (error) {
      console.log(error)
    }
  },
  /**
   * 对顾客进行标记完成
   */
  async handleTag(e) {
    const {
      item: { is_contact, is_call }
    } = e.currentTarget.dataset
    console.log(is_contact, is_call)
    // 没弹过loading && 未联系过客户
    if (is_contact === 0 && is_call === 0) {
      wx.showToast('建议联系客户后再标记完成哦~')
      return
    }

    if (this.data.isTagLoading) return

    this.data.isTagLoading = true

    let {
      index
    } = e.currentTarget.dataset
    await this.handleTagFinish(index)
  },
  /**
   * @description: 标记完成接口
   * @param {index} 顾客索引
   */
  async handleTagFinish(index) {
    wx.showLoading({ mask: true, title: '' })
    let userInfo = this.data.userList[index]
    const params = {
      task_allot_log_id: userInfo.task_allot_log_id,
      union_id: userInfo.union_id,
      customer_id: userInfo.customer_id,
      status: 1
    }

    try {
      const { code, message } = await tagCustomerApi(params)
      this.data.isTagLoading = false
      wx.hideLoading()
      if (code !== 200) return wx.showToast(message)
      // 标记成功后，维护tag状态
      let oneUserList = `userList[${index}]`
      this.setData({
        [oneUserList]: {
          ...userInfo,
          is_finished: 1
        }
      })
    } catch (error) {
      this.data.isTagLoading = false
      wx.hideLoading()
    }
  },
  /** 编辑说明 */
  handleRemark(e) {
    const { item, index } = e.currentTarget.dataset
    console.log('handleRemark', item)
    this.setData({
      remarkPopup: true,
      remarkContent: {
        item,
        index
      }
    })
  },
  /** 保存 取消说明 */
  confirmFn(e) {
    const { detail } = e
    if (detail.type === 'ok') {
      const { remark, index } = detail.data
      console.log('confirmFn', index, remark)
      let userInfo = this.data.userList[index]
      let oneUserList = `userList[${index}]`
      this.setData({
        [oneUserList]: {
          ...userInfo,
          remark
        },
        remarkPopup: false
      })
    } else {
      this.setData({
        remarkPopup: false
      })
    }
  },
  // 切换tab
  tabStairList(e) {
    let id = e.currentTarget.dataset.id
    let list = this.data.tabList
    list.forEach((item) => {
      if (id === item.type) {
        item.checked = true
      } else {
        item.checked = false
      }
    })
    this.setData({
      'meta.pageNo': 1,
      userList: [],
      tabList: list
    })
    this.getList()
  },
  /** 获取标签用户裂表 */
  getCustomer(pendRenderData = []) {
    if (pendRenderData.length === 0) {
      this.setData({
        isLoading: true
      })
      wx.showNavigationBarLoading()
    }
    const is_union = this.data.tabList.find((tab) => {
      return tab.checked
    }).type
    const { keywords, filterType } = this.data
    getCustomerApi({
      is_tag: this.data.pageType === 2 ? 1 : 0, // 0标签顾客  2导入顾客
      id: this.data.id,
      is_union: Number(is_union),
      page: this.data.meta.pageNo,
      size: this.data.meta.pageSize,
      task_category: this.data.category,
      keywords,
      filterType
    })
      .then((res) => {
        if (res.code === 200) {
          // 分页数据处理
          let { list = [], total, size } = res.data
          let totalPage = Math.ceil(total / size)
          let pageNo = this.data.meta.pageNo + 1 // 下一次请求页码

          // 判断自动请求下一页
          let { next, noMore, value } = this.justifyNext(
            pendRenderData,
            list,
            pageNo,
            totalPage
          )
          // 自动加载
          if (next) {
            this.setData({
              noMore: noMore
            })
            this.getList(value)
          } else {
            // 加载完毕，执行渲染
            wx.hideNavigationBarLoading()
            const userList = this.data.userList.concat(value)
            this.setData({
              noMore: noMore,
              isLoading: false,
              userList: userList // 更新视图
            })
          }
        } else {
          // 请求或者自动请求失败了
          this.dealError(pendRenderData)
        }
      })
      .catch((err) => {
        console.log(err)
        // 请求或者自动请求失败了
        this.dealError(pendRenderData)
      })
  },

  /** 获取退群名单 */
  getLeaveGroupList(pendRenderData = []) {
    const { id, tabList } = this.data
    const is_union = tabList.find((tab) => {
      return tab.checked
    }).type
    // 退群名单对应表 0全部，1好友，2非好友, 对应页面的状态
    const params = {
      id,
      status: Number(is_union), // 筛选状态，0全部，1好友，2非好友
      page: this.data.meta.pageNo,
      size: this.data.meta.pageSize
    }
    if (pendRenderData.length === 0) {
      this.setData({
        isLoading: true
      })
      wx.showNavigationBarLoading()
    }
    leaveGroupApi(params)
      .then((res) => {
        if (res.code === 200) {
          // 分页数据处理
          let { list = [], total, size } = res.data
          let totalPage = Math.ceil(total / size)
          let pageNo = this.data.meta.pageNo + 1 // 下一次请求页码

          // 判断自动请求下一页
          let { noMore, next, value } = this.justifyNext(
            pendRenderData,
            list,
            pageNo,
            totalPage
          )
          if (next) {
            this.setData({
              noMore: noMore
            })
            this.getList(value)
          } else {
            wx.hideNavigationBarLoading()
            const userList = this.data.userList.concat(value)
            this.setData({
              noMore: noMore,
              isLoading: false,
              userList: userList // 更新视图
            })
          }
        } else {
          this.dealError(pendRenderData)
        }
      })
      .catch((err) => {
        console.log(err)
        this.dealError(pendRenderData)
      })
  },

  /** 区分条件判断调用哪一个接口 */
  getList(pendRenderData = [], isSearch = false) {
    // pageType: 1.群/社群数据 2.用户标签/用户画像 3.导入顾客名单 4.顾客与导购&导购匹配关系
    // 查看退群顾客名单
    // 如果是筛选需要重置页码
    if (isSearch) {
      this.setData({
        'meta.pageNo': 1,
        userList: []
      })
    }
    if (this.data.pageType === 1) {
      this.getLeaveGroupList(pendRenderData)
    } else {
      // 查看标签顾客名单/导入顾客名单
      this.getCustomer(pendRenderData)
    }
  },

  /** 分页查询失败，页码回退 */
  dealError(pendRenderData) {
    wx.hideNavigationBarLoading()
    if (pendRenderData.length) {
      this.setData({
        isLoading: false,
        userList: pendRenderData
      })
    } else {
      this.setData({
        isLoading: false
      })
    }
  },
  /**
   * 1、判断是否需要继续加载
   * 2、更新页码
   * @returns object {next:true自动, false无需自动,value：数据}
   */
  justifyNext(pendRenderData, curData, pageNo, totalPage) {
    this.setData({
      'meta.pageNo': pageNo,
      'meta.totalPage': totalPage
    })
    // 当前请求页 无数据，自动请求
    if (curData.length === 0) {
      return {
        noMore: pageNo > totalPage,
        next: pageNo <= totalPage,
        value: pendRenderData
      }
    }
    // 页面数据+本次请求数据 小于 一页 时
    if (pendRenderData.length + curData.length < this.data.meta.pageSize) {
      return {
        noMore: pageNo > totalPage,
        next: pageNo <= totalPage,
        value: pendRenderData.concat(curData)
      }
    }
    // 其他场景无需自动加载,需要当前渲染数据
    return {
      noMore: pageNo > totalPage,
      next: false,
      value: pendRenderData.concat(curData)
    }
  },
  /**
   * 触底翻页
   */
  loadmore: throttle(async function () {
    console.log('触底翻页')
    let { meta, isLoading } = this.data
    // 分页完毕或者正在请求时，阻止下一次请求
    if (meta.pageNo > meta.totalPage || isLoading) return
    this.getList()
  }, 1000),
  /** 查看顾客详情 */
  viewDetail(e) {
    const self = this
    let {
      item: { work_external_user_id, is_union, customer_id },
      index
    } = e.currentTarget.dataset
    if (!customer_id && !work_external_user_id) {
      return
    }

    console.log('viewDetail', work_external_user_id, customer_id)
    let url = `/pages/customer-detail/customer-detail?openid=${work_external_user_id}&isRelation=${is_union}&isTask=1&index=${index}&customer_id=${customer_id}`
    wx.navigateTo({
      url: url,
      events: {
        // 为指定事件添加一个监听器，获取被打开页面传送到当前页面的数据
        acceptDataFromOpenedPage: async function (data) {
          // 在顾客详情操作了某条数据的拨打电话，此时，该条数据的添加好友可以操作了，编辑完成也可以操作了。
          console.log('acceptDataFromOpenedPage', data)
          let { customerIndex, key } = data
          let userInfo = self.data.userList[customerIndex]
          let oneUserList = `userList[${customerIndex}]`
          let upStatus = { }
          if (key === 'is_contact') {
            upStatus.is_contact = 1
          } else if (key === 'is_call') {
            upStatus.is_call = 1
          }
          // 子任务状态 关联 标记联系 0待完成,1已完成,2已过期
          if (self.data.status !== 2) {
            upStatus.is_finished = 1 // coach品牌对顾客进行打电话或者发消息，顾客状态自动标记为【已联系】
          }
          self.setData({
            [oneUserList]: {
              ...userInfo,
              ...upStatus
            }
          })
          // todo: 对该顾客进行标记已联系
          await self.tagCommunicate(customerIndex, key)
        }
      }
    })
  },
  /** 根据说明关键字查找 */
  handleSearch() {
    this.setData({
      'meta.pageNo': 1,
      userList: []
    })
    this.getList()
  },
  /** 清空输入框文案 */
  clearFilter() {
    this.setData({
      keywords: ''
    })
  },
  /** 切换过滤类型 */
  chooseType(e) {
    const { type } = e.currentTarget.dataset
    this.setData({
      filterType: type,
      showType: false
    })
    this.getList([], true)
  },
  /** 打开选择弹窗 */
  openSelect() {
    this.setData({
      showType: true
    })
  },
  /** 关闭选择弹窗 */
  cancelSelect() {
    this.setData({
      showType: false
    })
  },
  /** 输入框失去焦点 */
  inputChange(e) {
    console.log('inputBlur', e)
    const keywords = e.detail.value
    this.setData({
      keywords: keywords.trim()
    })
  },
  /** 生命周期 */
  onLoad(option) {
    console.log('onLoad', option)
    let list = this.data.tabList
    list.forEach((item) => {
      item.checked = false
    })
    // 默认勾选全部
    list[0].checked = true
    if (option.routerParams) {
      const routerParams = JSON.parse(decodeURIComponent(option.routerParams))
      console.log('routerParams', routerParams)
      this.setData({
        tabList: list,
        ...routerParams
      })
    }
    this.getList()
  }
})
