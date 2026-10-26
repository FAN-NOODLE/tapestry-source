
import { debounce } from '@/utils/tools'
import { getCustomerPhoneApi } from '@/api/customer-api'
import { reportAddFriendApi, reportCallApi, reportSendMsgApi, tagCustomerApi } from '@/api/plan-task-api'
import Retry from '@/lib/retry'
import { openUserProfile, openEnterpriseChat } from '@/utils/wx.qy'

const {globalData} = getApp()
Component({
  externalClasses: ['custom-class'],
  properties: {
    item: {
      type: Object
    },
    options: {
      type: Object
    },
    showSelect: {
      type: Boolean,
      value: false
    },
    type: {
      type: Number,
      value: 1
    },
    index: {
      type: Number
    }
  },
  data: {
    imageHost: 'https://cdn-cos.coach.com.cn/uploads',
    wxwork: globalData.system.environment === 'wxwork',
    defaultAvatar: 'https://cdn-cos.coach.com.cn/uploads/2023/v3.2.0.1/default.png',
    showPopup: false,
    remarkPopup: false
  },
  methods: {
    async reportUpdate(key) {
      const api = {
        is_contact: reportSendMsgApi,
        is_call: reportCallApi,
        is_friend: reportAddFriendApi
      }[key]

      const retryRequest = new Retry(api)

      const { task_allot_log_id, union_id, customer_id } = this.data.item
      // 标记完成
      const params = {
        task_allot_log_id: task_allot_log_id,
        union_id: union_id,
        customer_id: customer_id,
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
        const value = {
          [`item.${key}`]: 1
        }
        // 子任务状态 关联 标记联系 0 待完成,1 已完成,2 已过期
        if (this.data.options.expire_status !== 2) {
          value['item.is_finished'] = 1// coach品牌对顾客进行打电话或者发消息，顾客状态自动标记为【已联系】
        }
        this.setData(value)
        // 维护状态
        this.triggerEvent('update', { value: this.data.item, key })
      } catch (error) {
        console.log(error)
      }
    },
    /**
     * 电话联系
     */
    handleCall: debounce(async function () {
      wx.showLoading({ title: '', mask: true })
      // 上报打电话行为
      this.reportUpdate('is_call')
      this.handleTrackTaskDetailOperation(3)
      try {
        const { item } = this.data
        const params = {
          unionid: item.union_id,
          customer_code: item.customer_id
        }
        const { code, data, message } = await getCustomerPhoneApi(params)
        wx.hideLoading()
        this.data.isLoadingPhone = false
        if (code !== 200) {
          wx.showToast(message)
          return
        }
        if (data.mobile) {
          // 对某个人点击了打电话，则允许对该用户标记
          wx.makePhoneCall({
            phoneNumber: data.mobile.replace(/\s*/g, ''),
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
        console.log(error)
      }
    }),
    /**
     * 发送消息
     */
    handleSendMessage: debounce(async function () {
      const { work_external_user_id } = this.data.item
      // 打开用户详情页
      openEnterpriseChat({
        externalUserIds: work_external_user_id
      }).then(() => {
        // 上报发消息行为
        if (this.data.options.dimension !== 1) {
          this.reportUpdate('is_contact')
        }
        this.handleTrackTaskDetailOperation(4)
      }).catch(err => {
        console.log(err)
      })
    }),
    /**
     * 添加好友
     */
    handleAddUnion: debounce(async function (e) {
      const { work_external_user_id, is_call } = this.data.item
      // 对当前用户未执行打电话，直接点击添加好友，需阻止(退群名单例外)
      if (is_call === 0 && this.data.options.dimension !== 1) {
        this.setData({
          showPopup: true
        })
        return
      }
      // 打开用户详情页
      openUserProfile({
        type: 2, // 1表示该userid是企业成员，2表示该userid是外部联系人
        userid: work_external_user_id // 可以是企业成员，也可以是外部联系人
      }).then(() => {
        // 上报发消息行为
        this.reportUpdate('is_friend')
        this.handleTrackTaskDetailOperation(6)
      }).catch(err => {
        console.log(err)
      })
    }),
    /**
     * 添加备注
     */
    handleRemark() {
      this.setData({
        remarkPopup: true
      })
      this.handleTrackTaskDetailOperation(5)
    },
    /**
     * 添加备注确认
     * @param {*} e
     */
    handleRemarkConfirm(e) {
      const { detail } = e
      if (detail.type === 'ok') {
        const { remark } = detail.data
        this.setData({
          'item.remark': remark,
          remarkPopup: false
        })
      } else {
        this.setData({
          remarkPopup: false
        })
      }
    },
    /**
     * 关闭弹窗
     */
    handlePopupClose() {
      this.setData({
        showPopup: false
      })
    },
    /**
     * 点击事件
     */
    handleClick(e) {
      const { work_external_user_id, is_union, customer_id } = this.data.item

      if (!customer_id && !work_external_user_id) {
        return
      }

      let url = `/pages/customer-detail/customer-detail?openid=${work_external_user_id}&isRelation=${is_union}&isTask=1&index=${0}&customer_id=${customer_id}`
      // beta
      if (getApp().globalData.is_beta_user === 1) {
        url = `/packages/customer/pages/detail/detail?openid=${work_external_user_id}&isTask=1`
      }
      wx.navigateTo({
        url: url,
        events: {
          // 为指定事件添加一个监听器，获取被打开页面传送到当前页面的数据
          acceptDataFromOpenedPage: async (data) => {
            // 在顾客详情操作了某条数据的拨打电话，此时，该条数据的添加好友可以操作了，编辑完成也可以操作了。
            console.log('acceptDataFromOpenedPage', data)
            const { key } = data
            const upStatus = { }
            if (key === 'is_contact') {
              upStatus['item.is_contact'] = 1
            } else if (key === 'is_call') {
              upStatus['item.is_call'] = 1
            }
            // 子任务状态 关联 标记联系 0待完成,1已完成,2已过期
            if (this.data.options.expire_status !== 2) {
              upStatus['item.is_finished'] = 1 // coach品牌对顾客进行打电话或者发消息，顾客状态自动标记为【已联系】
            }
            this.setData(upStatus)
            // 对该顾客进行标记已联系
            await this.reportUpdate(key)
          }
        }
      })
    },
    /**
     * 对顾客进行标记完成
     */
    handleTag: debounce(async function(e) {
      this.handleTrackTaskDetailOperation(10)
      const { is_contact, is_call, task_allot_log_id, union_id, customer_id, is_finished } = this.data.item
      const { expire_status: status } = this.data.options
      if (Number(status) === 2 || is_finished === 1) return
      // 没弹过loading && 未联系过客户
      if (is_contact === 0 && is_call === 0) {
        wx.showToast('建议联系客户后再标记完成哦~')
        return
      }
      wx.showLoading({ mask: true, title: '' })
      const params = {
        task_allot_log_id: task_allot_log_id,
        union_id: union_id,
        customer_id: customer_id,
        status: 1
      }
      try {
        const { code, message } = await tagCustomerApi(params)
        wx.hideLoading()
        if (code !== 200) return wx.showToast(message)
        // 标记成功后，维护tag状态
        this.setData({
          'item.is_finished': 1
        })
      } catch (error) {
        wx.hideLoading()
      }
    }),
    /**
     * 选择
     */
    handleSelect() {
      this.setData({
        'item.checked': !this.data.item.checked
      })
      this.triggerEvent('select', {
        value: this.data.item
      })
      this.handleTrackTaskDetailOperation(this.data.item.checked ? 61 : 62)
    },
    /** 数据上报 */
    handleTrackTaskDetailOperation(type) {
      this.triggerEvent('track', {
        type
      })
    },
    /**
     * @description 显示tooltip
     */
    handleShowTooltip (e) {
      const { selector, content } = e.currentTarget.dataset
      if (content) {
        this.selectComponent('#mp-tooltip').show(this, selector, content)
      }
    }
  }
})
