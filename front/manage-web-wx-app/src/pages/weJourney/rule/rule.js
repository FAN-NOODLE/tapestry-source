import { integralRuleDescriptionApi } from '@/api/we-journey-api'
// 获取全局应用程序实例对象
// const app = getApp();

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    content: '' // 积分规则
  },

  // 生命周期函数--监听页面加载
  onLoad () {
    // TODO: onLoad
    this.integralRuleDescription()
  },

  /* 获取积分规则 */
  integralRuleDescription () {
    integralRuleDescriptionApi().then(res => {
      const { data, code, message } = res

      if (code !== 200) {
        wx.showToast(message)
        return
      }

      const { content } = data

      if (content) {
        this.setData({ content })
      }
    }).catch(err => {
      console.log(err)
    })
  }
})
