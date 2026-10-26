import Enums from '../../utils/enums'

// 获取全局应用程序实例对象
const { globalData, choosoeEnterpri } = getApp()

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    title: 'select-enterprise',
    enterprise: [],
    isPartner: ''
  },

  // 生命周期函数--监听页面显示
  onShow () {
    if (globalData.identityType === Enums.identityType.partner) {
      wx.setNavigationBarTitle({
        title: '选择账户'
      })
      this.setData({ isPartner: globalData.identityType })
    }
    Array.isArray(globalData.enterprise) &&
      this.setData({ enterprise: globalData.enterprise })
    console.log('单一企业自动选择，多个企业跳转筛选开始')
  },
  select (e) {
    this.setData({
      selectIdx: e.currentTarget.dataset.idx
    })
    const currentEnterprise = this.data.enterprise[e.currentTarget.dataset.idx]
    globalData.logo = currentEnterprise.logo
    // 存储当前所在企业
    globalData.currentEnterprise = currentEnterprise
    // 门店信息全局变量初始化
    globalData.store = {
      store_id: ''
    }
    // 用户信息字段全局变量初始化
    globalData.memberInfo = {
      isBoss: false,
      memberId: 0
    }
    // 存储当前角色(PS：店长、店员角色会在选择门店时调用member/detail接口重新赋值)
    // 合伙人又给的直接是个对象   又改回去了
    // if (globalData.identityType === 'partner') {
    //   globalData.roleMark = currentEnterprise && currentEnterprise.guide_roles && currentEnterprise.guide_roles.role_mark
    // } else {
    globalData.roleMark = currentEnterprise && currentEnterprise.guide_roles && currentEnterprise.guide_roles[0].role_mark

    choosoeEnterpri(
      {
        login_cache_key: globalData.login_cache_key,
        enterprise_id: currentEnterprise.enterprise_id
      },
      wx.navigateTo
    )
  }
})
