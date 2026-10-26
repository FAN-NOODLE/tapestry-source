import { memberPositionApi } from '@/api/compass-api'
import { userDetailTipsNumApi } from '@/api/customer-api'
import { getCustomerDetailApi, updateCustomerDescApi, updateCustomerRemarkApi } from '@/api/customer-beta-api'
import { update_time } from '@/utils/const'
import { openEnterpriseChat } from '@/utils/wx.qy'
import dayjs from 'dayjs'

const { talk_link } = require('@/lib/config')

const { globalData } = getApp()
// Page Object
Page({
  data: {
    tabs: [
      {
        value: 1,
        label: '客户信息',
        key: 3
      },
      {
        value: 2,
        label: '邀约记录',
        info: 0,
        key: 4
      },
      {
        value: 3,
        label: '待办任务',
        key: 9
      },
      {
        value: 4,
        label: '绑定关系',
        key: 6
      },
      {
        value: 5,
        label: '播种',
        info: 0,
        key: 5
      }
    ],
    /** 当前tab页 */
    current: 0,
    /** 显示备注弹框 */
    showRemarkPopup: false,
    /** 描述弹框 */
    showDescPopup: false,
    /** 客户详情 */
    detail: null,
    crmDetail: null,
    /** 是否是关联客户 */
    isRelation: false,
    /** 判断是否从企微侧边栏打开 */
    sidebar: false,
    /** 是否显示错误提示 */
    isError: false,
    /** 判断是否从任务计划顾客名单 */
    isTask: false,
    /** 用户所在的索引，反向查询用户数据 */
    customerIndex: 0,
    /** 客户id */
    openid: '',
    // 顾客id
    customer_id: '',
    /** role 3 企业导购 5 企业管理（目前只有3和5） */
    role: 3,
    /** 是否是企微环境 */
    wxwork: globalData.system.environment === 'wxwork',
    showTab: true,
    /** 6.2.7 迭代更新日期 */
    update: new Date() >= new Date(update_time)

  },
  // options(Object)
  /**
   * @description 获取页面参数
   * @param {Object} options
   * @param {Number} options.nav 导航栏索引
   * @param {String} options.memberUserid memberUserid
   * @param {Number} options.memberId memberId
   * @param {Number} options.isRelation 是否是关联客户，1是，0否
   * @param {Number} options.isTask 是否来自任务计划，1是，0否
   * @param {Number} options.index 客户索引
   * @param {String} options.customer_id 客户id
   *
   * @param {String} options.id 客户id
   * @param {String} options.openid 客户openid
   * @param {String} options.unionid unionid
   */
  onLoad(options) {
    this.setData({
      options,
      current: Number(options.nav) || 0,
      isRelation: Number(options.isRelation) === 1,
      isTask: options.isTask === '1',
      customerIndex: Number(options.index),
      customer_id: options.customer_id,
      openid: options.openid,
      memberUserid: options.memberUserid,
      memberId: Number(options.memberId)
    })
    this.getMember()
  },
  onShow () {
    this.getNum()
  },
  onReachBottom() {
    this.loadMore()
  },
  /**
   * 获取身份信息
   */
  async getMember () {
    await memberPositionApi({ is_new: 0 }).then(async res => {
      if (res.code === 200) {
        // role 3 企业导购 5 企业管理（目前只有3和5）
        this.setData({
          role: res.data.role
        })
        // todo:确保首个接口不是调用企微接口
        await this.getCustomerDetail()
        this.getNum()
      }
    }).catch(err => {
      console.log(err)
    })
  },
  /**
   * 获取顾客详情
   */
  async getCustomerDetail () {
    const { openid, options } = this.data
    let params = {}
    // 先判断work_external_user_id 这边是openid
    if (openid) {
      // 判断是否好友关系
      params = {
        openid: openid
      }
    } else {
      try {
        params.openid = await this.getUserId()
        console.log('侧边栏ID', params.openid)
        this.setData({
          openid: params.openid
        })
      } catch (e) {
        console.log(e)
        return
      }
    }
    const pages = getCurrentPages()

    // 神策数据上报 -- 查看标签
    wx.track.trackViewLabel({
      customer_id: params.openid,
      label_entrance: this.data.sidebar ? (Number(options.nav) ? 'sidebar-1' : 'sidebar') : pages[pages.length - 2] ? pages[pages.length - 2].route : '-1'
    })

    // params.is_work = sidebar
    // console.log('侧边栏参数', params)
    // wx.showNavigationBarLoading()
    try {
      const { code, data, message } = await getCustomerDetailApi(params)
      // wx.hideNavigationBarLoading()

      if (code !== 200) {
        wx.showToast(message)
        return
      }
      const now = dayjs().unix()
      data.is_relation = !!data.chat_tag
      data.history_amount_total_text = Number(data.history_amount_total).toLocaleString()
      data.amount_total_text = Number(data.amount_total).toLocaleString()
      data.active_time_text = data.active_time ? dayjs(data.active_time * 1000).format('YYYY-MM-DD') : '--'
      data.interact_time_text = data.interact_time ? dayjs(data.interact_time * 1000).format('YYYY-MM-DD') : '--'
      if (!data.last_purchase_date) {
        data.last_purchase_date_text = '--'
      } else if (now - data.last_purchase_date > 180 * 24 * 60 * 60) {
        data.last_purchase_date_text = '>180D'
      } else if (now - data.last_purchase_date > 90 * 24 * 60 * 60) {
        data.last_purchase_date_text = '>90D'
      } else if (now - data.last_purchase_date > 30 * 24 * 60 * 60) {
        data.last_purchase_date_text = '>30D'
      } else {
        data.last_purchase_date_text = '<30D'
      }

      if (now - data.valid_reserve_at < 30 * 24 * 60 * 60) {
        data.valid_reserve_at_text = '邀约中'
      } else {
        data.valid_reserve_at_text = '无邀约'
      }
      if (now - data.task_finish_at < 30 * 24 * 60 * 60) {
        data.task_finish_at_text = '已触达'
      } else {
        data.task_finish_at_text = '未触达'
      }

      // 标签
      data.tags = []
      if (data.price_sensitivity) {
        data.tags.push('价格敏感度: ' + {H: '高', L: '低'}[data.price_sensitivity])
      }
      if (data.high_willing === 'Y') {
        data.tags.push('高意愿用户')
      }
      if (data.category_preference && data.category_preference.length) {
        data.tags.push('商品偏好: ' + data.category_preference.join('/'))
      }

      const ageTag = {
        'Z世代': 'Gen Z',
        '千禧一代': 'YM'
      }
      if (data.age_tag && ageTag[data.age_tag]) {
        data.tags.push('顾客类型AGE: ' + (ageTag[data.age_tag] || data.age_tag))
      }

      this.setData({ detail: data })
    } catch (error) {
      console.log(error)
      wx.hideNavigationBarLoading()
    }
  },
  /**
   * 获取用户ID
   */
  getUserId () {
    const _this = this
    return new Promise((resolve, reject) => {
      wx.qy.getContext({
        success: res => {
          var entry = res.entry // 返回进入小程序的入口类型

          // 如果从企微侧边栏播种入口打开,
          if (entry === 'group_chat_tools' && _this.data.navId === 2) {
            _this.setData({ isError: true })
            reject(new Error('从群聊进入'))
          }
          if (entry === 'single_chat_tools' || entry === 'contact_profile') {
            wx.qy.getCurExternalContact({
              success: res => {
                console.log('外部联系人返回res,=======>', res)
                _this.setData({ openid: res.userId, isRelation: true, sidebar: true })// 返回当前外部联系人userId
                resolve(res.userId)
                // 取消顶部返回按钮
                wx.hideHomeButton()
                // 神策数据上报 -- 点击侧边栏
                wx.track.trackClickItem({
                  item_name: _this.route,
                  user_external_id: res.userId
                })
              },
              fail: err => {
                console.error(err, 'getCurExternalContact===')
                reject(err)
              }
            })
          } else {
            console.log('wx.qy.getCurExternalContact==')
            reject(Error)
          }
        },
        fail: function (err) {
          console.error('wx.qy.getContext 异常', err)
          reject(err)
        }
      })
    })
  },
  /**
   * 获取顾客邀约和播种数量
   */
  async getNum () {
    try {
      if (!this.data.openid) return
      const { data, code, message } = await userDetailTipsNumApi({ external_userid: this.data.openid })

      if (code !== 200) {
        wx.showToast(message)
        return
      }

      const { reserve_num = 0, sow_task_num = 0, plan_task_num = 0 } = data
      const { tabs } = this.data

      tabs[1].info = reserve_num
      tabs[2].info = plan_task_num
      tabs[4].info = sow_task_num

      this.setData({ tabs })
    } catch (error) {
      console.log(error)
    }
  },
  loadMore() {
    const { current } = this.data
    console.log(current)

    let id = ''

    // 预约列表
    if (current === 1) {
      id = '#invite-list'
    } else if (current === 2) {
      id = '#task-list'
    } else if (current === 4) {
      id = '#sow-list'
    }

    if (id) {
      this.selectComponent(id).loadMore()
    }
  },

  handleUpdateCrmDetail(e) {
    this.setData({
      crmDetail: e.detail
    })
  },

  handleTabChange(e) {
    console.log(e.detail.index)

    this.setData({
      current: e.detail.index
    })

    if (e.detail.index === 0) {
      this.selectComponent('#base-info').reload()
    } else {
      this.selectComponent('#base-info').hide()
    }

    const { key } = this.data.tabs[e.detail.index]
    // 神策数据上报 -- 顾客详情页操作
    this.trackClickDetails(key)
  },

  /**
   * 基础操作 - 拨打电话
   */
  handleMakePhone() {
    if (!this.data.crmDetail.mobile) {
      wx.showToast('该用户无电话')
      return
    }
    wx.makePhoneCall({
      phoneNumber: this.data.crmDetail.mobile,
      success: () => {
        // 神策数据上报 -- 拨打顾客电话
        wx.track.trackCallCustomer({
          user_external_id: this.data.detail.openid
        })
        // 神策数据上报 -- 顾客详情页操作
        this.trackClickDetails(12)
      }
    })

    // 从任务场景查看的客户详情，有拨打电话操作 联动任务计划
    if (this.data.isTask) {
      const eventChannel = this.getOpenerEventChannel()
      eventChannel.emit('acceptDataFromOpenedPage', {
        key: 'is_call'
      })
    }
  },

  /**
   * 基础操作 - 打开聊天框
   */
  handleOpenChat() {
    if (!this.data.detail.is_relation) {
      wx.showToast('该用户非您好友，无法打开对话框')
      return
    }
    openEnterpriseChat({
      externalUserIds: this.data.detail.openid
    })
    // 从任务场景查看的客户详情，有拨打电话操作 联动任务计划
    if (this.data.isTask) {
      const eventChannel = this.getOpenerEventChannel()
      eventChannel.emit('acceptDataFromOpenedPage', {
        key: 'is_contact'
      })
    }
    // 神策数据上报 -- 顾客详情页操作
    this.trackClickDetails(0)
  },
  /**
   * 基础操作 - 去购物清单页面
   */
  handleViewShopList () {
    const url = `/packages/customer/pages/consume-record/consume-record?unionId=${this.data.detail.unionid}&external_userid=${this.data.openid}`
    wx.navigateTo({ url })
    // 神策数据上报 -- 顾客详情页操作
    this.trackClickDetails(1)
  },
  /**
   * 基础操作 - 去卡券详情页面
   */
  handleViewCouponDetail () {
    const {unionid, membership} = this.data.detail
    const params = {
      union_id: unionid || '',
      customer_code: membership || ''
    }
    const url = `/pages/coupon-list/coupon-list?info=${JSON.stringify(params)}`
    wx.navigateTo({ url })
    // 神策数据上报 -- 顾客详情页操作
    this.trackClickDetails(2)
  },
  /**
   * 基础操作 - 显示备注弹框
   */
  handleShowRemarkPopup() {
    if (!this.data.detail.is_relation) {
      wx.showToast('该用户非您好友，无法编辑')
      return
    }
    this.setData({
      showRemarkPopup: true
    })
  },
  /**
   * 基础操作 - 显示备注弹框 确定
   * @param {*} e
   */
  async handleEditRemarkConfirm(e) {
    const { value } = e.detail
    if (value === this.data.detail.remark) {
      this.setData({
        showRemarkPopup: false
      })
      wx.showToast('编辑成功')
      return
    }
    wx.showLoading()

    const res = await updateCustomerRemarkApi({
      external_userid: this.data.openid,
      remark: value
    })
    if (res.code === 200) {
      this.setData({
        'detail.remark': value,
        showRemarkPopup: false
      })
      wx.showToast('编辑成功')
      this.trackClickDetails(7)
    } else {
      wx.showToast(res.message)
    }
  },
  /**
   * 基础操作 - 显示描述弹框
   */
  handleShowDescPopup() {
    if (!this.data.detail.is_relation) {
      wx.showToast('该用户非您好友，无法编辑')
      return
    }
    this.setData({
      showDescPopup: true
    })
  },
  /**
   * 基础操作 - 显示描述弹框 确定
   * @param {*} e
   */
  async handleEditDescConfirm(e) {
    const { value } = e.detail
    if (value === this.data.detail.description) {
      this.setData({
        showDescPopup: false
      })
      wx.showToast('编辑成功')
      return
    }
    wx.showLoading()
    const res = await updateCustomerDescApi({
      external_userid: this.data.openid,
      description: value
    })
    if (res.code === 200) {
      this.setData({
        'detail.description': value,
        showDescPopup: false
      })
      wx.showToast('编辑成功')
      this.trackClickDetails(8)
    } else {
      wx.showToast(res.message)
    }
  },
  /**
   * 公用 - 关闭弹框
   */
  handleCancel() {
    this.setData({
      showRemarkPopup: false,
      showDescPopup: false
    })
  },
  toggleTab() {
    if (!this.data.showTab) {
      setTimeout(() => {
        this.setData({
          showTab: true
        })
      }, 500)
    } else {
      this.setData({
        showTab: false
      })
    }
  },

  // 神策数据上报 -- 顾客详情页操作
  trackClickDetails(type) {
    const { openid } = this.data.detail
    wx.track.trackClickDetails({
      user_external_id: openid,
      operation_type: type
    })
  },
  handleOpenTalk() {
    wx.navigateTo({
      url: `/pages/webview/webview?url=${encodeURIComponent(talk_link)}`
    })
  }

})
