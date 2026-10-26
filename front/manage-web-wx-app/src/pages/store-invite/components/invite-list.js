import { changeStatus, commentShareApi, verifyShareApi } from '@/api/guide-invite'
import { StatusList } from '../constants'

Component({
  properties: {
    item: Object,
    role: Number,
    index: Number// 当前点击的列表index
  },
  data: {
    is_store: 1,
    customer_status: '',
    statusList: StatusList,
    isOpenTipsModal: false,
    // appid: require('../../../lib/config').appid,
    isShowPriceTips: false,
    priceTips: '订单金额',
    priceTipsTwo: '此处的订单金额为含税金额',
    isComment: false,
    /** 是否从侧边栏直接打开 */
    singleChat: false
  },

  attached() { },
  ready() {
    this.judgeEntrance()
    const { item, role } = this.properties

    this.setData({
      role: Number(role) || 0,
      is_store: item.is_store || 0,
      customer_status: item.customer_status || 0,
      isComment: item.status === 1 || item.status === 2
    })
  },
  methods: {
    /** 获取用户是否从侧边栏打开 */
    judgeEntrance() {
      wx.qy.getContext({
        success: (res) => {
          var entry = res.entry // 返回进入小程序的入口类型
          if (entry === 'single_chat_tools') {
            this.setData({
              singleChat: true
            })
          }
        }
      })
    },
    // 状态信息弹窗
    onTips() {
      this.setData({
        isShowPriceTips: false,
        isOpenTipsModal: true
      })
    },
    // 关闭状态信息弹窗
    closeTipsModel() {
      this.setData({
        isOpenTipsModal: false
      })
    },
    // 显示订单金额弹窗
    orderPriceTips() {
      this.setData({
        isOpenTipsModal: true,
        isShowPriceTips: true
      })
    },
    goDetail (e) {
      const {item, index} = e.currentTarget.dataset
      const _this = this
      wx.navigateTo({
        url: `/pages/store-invite/invite-detail/invite-detail?id=${item.id}&reserveNumber=${item.reserve_number}&role=${this.data.role}&index=${this.properties.index}`,
        events: {
          // 为指定事件添加一个监听器，获取被打开页面传送到当前页面的数据
          acceptDataFromOpenedPage: function (res) {
            console.log('acceptDataFromOpenedPage', res)
            // 状态与当前不一致
            if (_this.properties.item.status !== res.status) {
              _this.triggerEvent('statusChange', {index, item: res})
            }
          }
        }
      })
    },
    // 操作 - 是否到店购买 & 选择顾客类型
    onChangeStatus(e) {
      if (e.detail.value) {
        // 更新接口信息
        if (e.currentTarget.dataset.item && e.currentTarget.dataset.type) {
          const { item, type } = e.currentTarget.dataset
          const params = {
            id: item.id,
            reserve_number: item.reserve_number,
            [type]: e.detail.value && Number(e.detail.value)
          }
          changeStatus(params).then(res => {
            if (res.code === 200) {
              this.setData({
                [type]: e.detail.value
              })
            } else {
              if (type === 'is_store') {
                this.setData({
                  is_store: this.data.is_store
                })
              } else {
                this.setData({
                  customer_status: this.data.customer_status
                })
              }
              // todo: 提示错误信息
              wx.showToast(res.message)
            }
          }).catch(err => {
            console.log(err)
          })
        }
      }
      return false
    },

    /** 侧边栏打开返送邀请评价 */
    sidebarSend(id, reserve_number, title, url) {
      wx.qy.sendChatMessage({
        msgtype: 'miniprogram', // 消息类型，必填
        enterChat: true,
        miniprogram: {
          appid: wx.getAccountInfoSync().miniProgram.appId,
          title: title, // 小程序消息的title
          imgUrl: 'https://cdn-cos.coach.com.cn/uploads/2022/5/23/60f86724d69e71c10979b4e5bb17fa23.png', // 小程序消息的封面图。
          page: url // 小程序消息打开后的路径，注意要以.html作为后缀，否则在微信端打开会提示找不到页面
        },
        success: (res) => {
          console.log('inviteComments', res, id, reserve_number)
          // 记录邀请评价日志
          const params = {
            id: id, // 主键id
            reserve_number: reserve_number// 预约单号
          }
          commentShareApi(params)
        },
        fail: (err) => {
          console.log('inviteCommentsFail', err)
          // 记录邀请评价日志
        }
      })
    },

    /** 判断是否导购自己的邀请单 */
    verifyShareSend(id, reserve_number) {
      return new Promise(resolve => {
        const params = {
          id,
          reserve_number
        }
        verifyShareApi(params).then(res => {
          if (res.code !== 200) {
            resolve(false)
            return wx.showToast(res.message)
          }
          resolve(true)
        }).catch(e => {
          resolve(false)
          console.error(e, 'verifyShareSend 出错')
        })
      })
    },

    // 邀请顾客评价
    async inviteComments(item) {
      const { id, reserve_number, is_comment } = item.currentTarget && item.currentTarget.dataset.item
      const isOwn = await this.verifyShareSend(id, reserve_number)
      // 判断是否导购自己的单子
      if (!isOwn) return
      if (is_comment === 1) return
      let title = '邀约评价'
      let url = `pagesOrder/pages/store-invite/comment/comment.html?id=${id}&reserve_number=${reserve_number}`
      this.sidebarSend(id, reserve_number, title, url)
    }
  }
})
