import {
  recordInfo
} from '@/api/invitation-letter'
import { toDate } from '@/utils/tools'

/** 时间拼接0 */
let tf = function (i) {
  return (i < 10 ? '0' : '') + i
}

Page({
  data: {
    /** 邀请函详情 */
    detail: {},

    /** 是否正在加载 */
    isLoading: true,

    /** 邀请函ID */
    invite_record_id: ''

  },

  onLoad ({ invite_record_id }) {
    this.setData({
      invite_record_id: invite_record_id || ''
    })
    wx.hideHomeButton()
  },

  async onShow () {
    this.getInviteDetail()
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
      const res = await recordInfo(params)
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

  formatData (time, select_hours) {
    const {nowYear, nowMonth, nowDate} = toDate(
      time,
      'yyyy-MM-dd hh:mm:ss'
    )
    const time_str = `${tf(nowYear)}-${tf(nowMonth)}-${tf(nowDate)} ${select_hours}`

    return time_str
  }
})
