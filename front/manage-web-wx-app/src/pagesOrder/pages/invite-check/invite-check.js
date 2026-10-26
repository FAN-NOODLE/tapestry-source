import { createOrder, inviteDetail, inviteInfo } from '../../api/invite'
import { getMyCenter } from '../../api/profile'
import { getConfigMessageApi, orderNotifyApi } from '../../api/guide-invite'
import { toDate } from '@/utils/tools'
import { sensorsLoginInit } from '@/lib/sensors'

const { globalData } = getApp()

/** 时间拼接0 */
let tf = function (i) {
  return (i < 10 ? '0' : '') + i
}

Page({
  data: {
    /** 已选择时间 */
    timeObj: {},
    dateModel: false,
    showTime: true,

    /** 邀请函详情 */
    detail: {},

    /** 是否正在加载 */
    isLoading: true,

    /** 邀请函ID */
    invite_record_id: '',

    /** 选择的预约时间 */
    selectTime: null,

    /** 用户信息 */
    myInfo: {},
    /** 页面来源 code为后台扫码预览 */
    from: ''
  },

  onLoad ({ id, scene }) {
    let from = ''
    if (scene) {
      const query = decodeURIComponent(scene).split('&')
      console.log(query)
      id = query[0].split('=')[1]
      console.log('id', id)
      from = query[1].split('=')[1]
      console.log('from', from)
    }
    // 如果是企微端打开邀约详情，需要重定向到企微端的邀约详情
    if (globalData.system.environment === 'wxwork' && from !== 'code') {
      return wx.reLaunch({
        url: `/pages/invite-check/invite-check?corp_id=${globalData.currentEnterprise.corp_id || ''}&invite_record_id=${id || ''}`
      })
    }

    this.setData({
      invite_record_id: id || '',
      from
    })
    wx.hideHomeButton()
  },

  async onShow () {
    await this.getMyCenter()
    this.data.from === 'code' ? this.getInviteInfo() : this.getInviteDetail()
  },

  /** 查看邀请函详情 */
  async getInviteDetail () {
    wx.showLoading({
      title: '加载中'
    })
    try {
      const { invite_record_id } = this.data
      let params = {
        invite_record_id: Number(invite_record_id)
      }
      const res = await inviteDetail(params)
      wx.hideLoading()
      if (res.code === 200) {
        // todo: 到店开始结束时间
        const date = new Date(res.data.reserve_time * 1000)
        let select_hours = ''
        if (res.data.reserve_end_time === 0) {
          select_hours = `${date.getHours()}:00 - ${date.getHours() + 1}:00`
        } else {
          const endDate = new Date(res.data.reserve_end_time * 1000)
          select_hours = `${date.getHours()}:00 - ${endDate.getHours()}:00`
        }
        let detail = {
          ...res.data,
          reserve_time: res.data.reserve_time ? this.formatData(res.data.reserve_time, select_hours) : ''
        }
        this.setData({
          detail
        })
      } else {
        wx.showToast({
          title: res.message
        })
      }
    } catch (error) {
      wx.hideLoading()
      console.log(error)
    } finally {
      this.setData({
        isLoading: false
      })
    }
  },

  /** 查看邀请函预览详情 */
  async getInviteInfo () {
    wx.showLoading({
      title: '加载中'
    })
    try {
      const { invite_record_id } = this.data
      let params = {
        id: Number(invite_record_id)
      }
      const res = await inviteInfo(params)
      wx.hideLoading()
      if (res.code === 200) {
        let detail = res.data
        this.setData({
          detail
        })
      } else {
        wx.showToast({
          title: res.message
        })
      }
    } catch (error) {
      wx.hideLoading()
      console.log(error)
    } finally {
      this.setData({
        isLoading: false
      })
    }
  },

  /** 获取用户信息 */
  async getMyCenter () {
    try {
      let res = await getMyCenter()
      if (res.code === 200) {
        this.setData({
          myInfo: res.data
        })
      } else {
        wx.showToast({
          title: res.message
        })
      }
    } catch (error) {
      console.log(error)
    }
  },

  /** 选择到店时间 */
  handleChooseTime () {
    if (this.data.from === 'code') return
    this.setData({
      dateModel: true
    })
  },

  /** 确定选择到店时间 */
  submit (e) {
    this.setData(
      {
        selectTime: e.detail
      },
      () => {
        this.closeModel()
      }
    )
  },

  /** 关闭选择到店时间 */
  closeModel () {
    this.setData({
      dateModel: false
    })
  },

  formatData (time, select_hours) {
    const {nowYear, nowMonth, nowDate} = toDate(
      time,
      'yyyy-MM-dd hh:mm:ss'
    )
    const time_str = `${tf(nowYear)}-${tf(nowMonth)}-${tf(nowDate)} ${select_hours}`

    return time_str
  },

  /** 提交预约 */
  async handleConfirm () {
    if (this.data.from === 'code') return

    try {
      const { invite_record_id, selectTime, detail } = this.data

      if (!selectTime) {
        wx.showToast({
          title: '请选择到店时间',
          icon: 'none'
        })
        return
      }

      wx.showLoading({
        title: '加载中'
      })

      const params = {
        invite_record_id: Number(invite_record_id),
        reserve_time: selectTime.startTime,
        reserve_end_time: selectTime.endTime
      }

      const res = await createOrder(params)
      wx.hideLoading()
      if (res.code === 200) {
        wx.showToast({
          title: '预约成功'
        })
        this.getInviteDetail()
        // 获取订阅消息配置
        this.getConfigMessage(res.data)

        // 神策数据上报 -- 开始
        sensorsLoginInit(detail.member_user_id)
        // 神策数据上报 -- 邀请函被接受
        wx.track.trackInvitationAccept({
          invite_record_id: invite_record_id,
          invite_template_id: detail.id,
          user_external_id: detail.user_openid,
          reserve_name: detail.reserve_service_name,
          reserve_time: selectTime.startTime,
          time_period: selectTime.startTime,
          store_code: detail.store_code || detail.store_name,
          store_name: detail.store_name
        })
      } else {
        wx.showToast({
          title: res.message
        })
      }
    } catch (error) {
      wx.hideLoading()
      console.log(error)
    }
  },

  /** 获取订阅消息配置 */
  getConfigMessage (data) {
    getConfigMessageApi()
      .then((res) => {
        if (res.code === 200) {
          this.data.templateIds = res.data['confirm-subscribe']
          this.subscribe(data)
        } else {
          wx.showToast({
            title: res.message
          })
        }
      })
      .catch((err) => {
        console.error(err, '====')
      })
  },

  /** 发起订阅 */
  subscribe (data) {
    const { templateIds } = this.data
    wx.requestSubscribeMessage({
      tmplIds: [...templateIds],
      success (res) {
        // 过滤出接受的template_id
        const tempRes = Object.keys(res).filter(
          (e) => e !== 'errMsg' && res[e] === 'accept'
        )
        if (res.errMsg === 'requestSubscribeMessage:ok' && tempRes.length > 0) {
          const params = {
            rel_id: data.id,
            reserve_number: data.reserve_number,
            tmpl_list: templateIds.map((item) => {
              return {
                pri_tmpl_id: item,
                url: `?id=${data.id}&reserve_number=${data.reserve_number}`
              }
            })
          }
          orderNotifyApi(params)
            .then((res) => {
              if (res.code === 200) {
              } else {
                wx.codeException(res)
              }
            })
            .catch((err) => {
              console.error('订阅错误信息-->', err)
            })
        } else {
          wx.showToast({
            title:
              '您已关闭/拒收此类订阅消息，可在小程序设置-接受订阅消息进行设置开启',
            icon: 'none',
            duration: 2000,
            mask: false
          })
        }
      },
      fail (err) {
        console.error('订阅错误信息', err)
        wx.showToast({
          title:
            '您已关闭/拒收此类订阅消息，可在小程序设置-接受订阅消息进行设置开启',
          icon: 'none',
          duration: 2000,
          mask: false
        })
      },
      complete: () => {
        // 确认订单
        // this.sureOperation();
      }
    })
  }
})
