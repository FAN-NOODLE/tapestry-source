import { wxworkLoginApi } from '@/api/login-api'

const { globalData } = getApp()

Page({
  data: {
    options: {},
    shadeShow: true
  },
  /**
   *
   * @param {*} options
   * @param {*} options.target 目标页面 字符串 materialList-素材列表
   * @param {*} options.query 查询参数 字符串 encodeURIComponent('from=TALK&current=1&sku_item={skucode或者skucode_颜色}')
   */
  onLoad(options) {
    console.log('页面参数', options)
    this.setData({
      options
    })
    this.handleNavigate()
  },

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
        this.handleNavigate()
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
  },

  handleNavigate() {
    const { target, query } = this.data.options
    const pages = {
      materialList: '/pagesMaterial/pages/list/list',
      goodsMaterialDetail: '/pagesMaterial/pages/select/select'
    }
    wx.redirectTo({
      url: pages[target] + '?' + decodeURIComponent(query)
    })
  }
})
