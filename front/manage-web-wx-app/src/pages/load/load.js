import { wxworkLoginApi } from '@/api/login-api'

const { globalData } = getApp()

Page({
  data: {
    title: 'load',
    shadeShow: true
  },
  onLoad () {
    // 顾客端进入
    if (globalData.system.environment !== 'wxwork') {
      return wx.redirectTo({
        url: '/pagesOrder/pages/profile/profile'
      })
    }

    // 导购端进入, 重置登录信息
    wx.removeStorageSync('token')
    globalData.checkSessionOK = false
    globalData.token = ''
    this.login()
  },

  /** 登录，空接口（wxworkLoginApi） */
  async login () {
    const res = await wxworkLoginApi()
    this.enterprise(res)
  },

  // 选择企业
  enterprise (res) {
    console.log(res, '登录后参数')
    if (Array.isArray(res.data.enterprise) && res.data.enterprise.length > 0) {
      // 存储用户可选企业列表
      globalData.enterprise = res.data.enterprise
      // 存储缓存key
      globalData.login_cache_key = res.data.login_cache_key
      if (res.data.enterprise.length === 1) {
        // TODO 调用选择企业
        const currentEnterprise = res.data.enterprise[0]
        globalData.logo = currentEnterprise.logo
        // 存储当前所在企业
        globalData.currentEnterprise = currentEnterprise
        // 存储当前角色(PS：店长、店员角色会在选择门店时调用member/detail接口重新赋值)
        globalData.roleMark = currentEnterprise && currentEnterprise.guide_roles && currentEnterprise.guide_roles[0].role_mark
        this.setData({
          shadeShow: false
        })
        // 跳转数据罗盘首页
        wx.switchTab({
          url: '/pages/compass-data/compass-data'
        })
      } else {
        wx.redirectTo({
          url: '/pages/select-enterprise/select-enterprise'
        })
      }
    } else {
      wx.redirectTo({
        url: '/pages/login/login'
      })
    }
  }
})
