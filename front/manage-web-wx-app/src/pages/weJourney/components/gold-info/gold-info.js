import { getCoinsTokenApi } from '@/api/we-journey-api'
import {getRedirectPrarms} from '../../utils'

const imagePath = 'https://cdn-cos.coach.com.cn/uploads/'

Component({
  properties: {
    pageInfo: {
      type: Object,
      value: {}
    }
  },
  data: {
    imagePath: imagePath
  },
  lifetimes: {
    attached () {
      this.getCoinsToken(false)
    }
  },
  methods: {
    /* 前往金币详细信息 */
    goCoinsInfo () {
      wx.navigateTo({
        url: '/pages/weJourney/points/points?tab=2'
      })
    },
    /* 获取兑换中心入口配置 */
    getCoinsToken (isGetToken = true) {
      const timeStamps = new Date().getTime()
      getCoinsTokenApi({
        login_str: getRedirectPrarms('mall', timeStamps).params
      }).then(res => {
        if (res.code === 200) {
          const { coins_settings_info } = res.data
          this.setData({
            coins_settings_info: coins_settings_info
          })
          if (coins_settings_info.mall_is_open === 0 && isGetToken) {
            wx.showToast('兑换中心暂未开放')
          } else if (isGetToken) {
            const redirectUrl = `${getRedirectPrarms('mall', timeStamps).redirectUrl}&loginToken=${res.data.login_token}`
            console.log(redirectUrl)
            wx.navigateTo({ url: `/pages/webview/webview?url=${encodeURIComponent(redirectUrl)}` })
          }
        }
      }).catch(err => {
        console.log(err)
      })
    }
  }
})
