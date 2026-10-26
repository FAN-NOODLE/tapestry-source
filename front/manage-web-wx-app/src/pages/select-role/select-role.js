// 获取全局应用程序实例对象
import Enums from '@/utils/enums'

const { globalData, choosoeEnterpri, showToast } = getApp()
// const app = getApp();

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    title: 'select-role',
    role: '',
    res: {}
  },

  // 生命周期函数--监听页面加载
  onLoad () {
    this.data.res = globalData.tokenObj
  },
  // 选择角色
  selectRole (e) {
    const tokenIndex = Number(e.currentTarget.dataset.type)
    if (tokenIndex === 1) {
      globalData.token = this.data.res.data.token
      wx.setStorageSync('token', this.data.res.data.token)
      globalData.identityType = Enums.identityType.salesman
    } else if (tokenIndex === 2) {
      globalData.token = this.data.res.data.partner_token
      wx.setStorageSync('token', this.data.res.data.partner_token)
      globalData.identityType = Enums.identityType.partner
      // globalData.identityType = 'partner'

      this.data.res.data.enterprise = JSON.parse(JSON.stringify(this.data.res.data.partner_enterprise_list))// 改动小 后期优化
    }

    globalData.isFirstLogin = this.data.res.data.is_first_login === 1
    showToast('登录成功', true)
    this.enterprise(this.data.res)
  },
  // 选择企业
  enterprise (res) {
    if (Array.isArray(res.data.enterprise)) {
      // 存储用户头像
      globalData.guide_avatar = res.data.avatar
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
        // 合伙人又给的直接是个对象  多个注意统一修改

        globalData.roleMark = currentEnterprise && currentEnterprise.guide_roles && currentEnterprise.guide_roles[0].role_mark

        choosoeEnterpri(
          {
            login_cache_key: res.data.login_cache_key,
            enterprise_id: currentEnterprise.enterprise_id
          },
          wx.redirectTo
        )
      } else {
        wx.redirectTo({
          url: '/pages/select-enterprise/select-enterprise'
        })
      }
    } else {
      showToast('返回企业格式错误' + typeof res.data.enterprise)
    }
  },

  // 页面分享设置
  onShareAppMessage () {
    return {
      title: 'share title',
      path: 'pages/select-role/select-role'
    }
  }
})
