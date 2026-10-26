import { split } from '../../utils/tools'
import { getMemberDetailApi } from '@/api/store-api'

const { globalData } = getApp()
// 创建页面实例对象
Page({
  data: {
    show: false // 控制页面显隐
  },
  onLoad () {
    this.getDetail()
  },
  onShow () { },
  // 获取员工信息接口
  getDetail () {
    const { memberId, login_cache_key } = globalData.memberInfo
    wx.showNavigationBarLoading()
    getMemberDetailApi({
      id: memberId,
      login_cache_key
    })
      .then((res) => {
        if (res.code === 200) {
          this.setData({
            ...res.data,
            show: true,
            switchChecked: res.data.is_send_sms === 1,
            mobile: split(res.data.mobile, 'space'),
            sms_mobile: split(res.data.sms_mobile, 'space')
          })
        } else {
          wx.showToast(res.message)
        }
        wx.hideNavigationBarLoading()
      })
      .catch((err) => {
        wx.hideNavigationBarLoading()
        wx.showToast(err.message)
      })
  },
  // 查看名片操作
  preview () {
    wx.navigateTo({
      url: '/pages/business-card/business-card'
    })
  }
})
