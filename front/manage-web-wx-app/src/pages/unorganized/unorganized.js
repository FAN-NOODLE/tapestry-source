// 获取全局应用程序实例对象
import Enums from '../../utils/enums'
import { getPartnerContactApi } from '../../api/login-api'

const { globalData } = getApp()
// 创建页面实例对象
Page({
  data: {
    confirmModelShow: false,
    ispartner: '',
    isPartnerDisabled: false
  },
  onLoad (options) {
    this.setData({
      code: options.code,
      appId: options.appId || '',
      memUid: options.memUid || '',
      path: options.path || '',
      ispartner: globalData.identityType === Enums.identityType.partner
    })
    if (globalData.identityType === Enums.identityType.partner) {
      if (globalData.globalPartnerStatus === 0) {
        options.code = '666661'
        this.setData({
          code: options.code,
          isPartnerDisabled: true
        })
      } else if (globalData.currentEnterprise.partner_status === 0) {
        // options.code = '666662'
        // this.setData({
        //   code: options.code
        // })
      }
    }
    this.dealRole(options.code)
  },
  onShow () {
    if (globalData.identityType === Enums.identityType.partner) {
      wx.setNavigationBarTitle({
        title: '品牌合伙人'
      })
    }
    const { extraData } = globalData.appOnShowoption && globalData.appOnShowoption.referrerInfo
    // 因为场景参数全局可能要用，所以不能随意清理，为了防止死循环，因此加个needReload标记
    if (extraData && extraData.authenticationSuccess && globalData.needReload) {
      globalData.needReload = false
      wx.redirectTo({
        url: '/pages/load/load'
      })
    }
  },
  // 关闭二次确认弹框
  closeConfirmModel () {
    this.setData({
      confirmModelShow: false
    })
    if (this.data.code === '656005') {
      return
    }
    if (this.data.code === '666662') {
      this.onPartnerAccountDisabledClose()
    }
    this.gotoNext()
  },
  // 进入下一步
  gotoNext () {
    // 小程序正常运行模式
    // 可选门店的角色(店员、店长)
    let canChooseStoreRole = ['guide_store_member', 'guide_store_manager']
    let currentRole = globalData.roleMark
    if (canChooseStoreRole.some((item) => item === currentRole)) {
      wx.redirectTo({
        url: '/pages/select-store/select-store'
      })
    } else {
      if (globalData.identityType === Enums.identityType.partner) {
        console.log('跳转至品牌合伙人')
      } else {
        wx.switchTab({
          url: '/pages/compass-data/compass-data'
        })
      }
    }
  },
  onPartnerAccountDisabledClose () {
    if (globalData.enterprise && globalData.enterprise.length > 1) {
      wx.redirectTo({
        url: '/pages/select-enterprise/select-enterprise'
      })
    } else {
      wx.removeStorageSync('token')
      globalData.token = ''
      wx.redirectTo({
        url: '/pages/login/login'
      })
    }
  },
  // 弹框确定按钮回调
  onModelConfirm () {
    switch (this.data.code) {
      // 企业未创建或发布小程序商城
      case '656005':
        this.setData({
          confirmModelShow: false
        })
        break
      // 无导购组织
      case '652002':
        this.setData({
          confirmModelShow: false
        })
        break
      // 导购未完成商城身份认证
      case '656006':
        globalData.needReload = true
        wx.navigateToMiniProgram({
          appId: this.data.appId,
          // envVersion: 'develop',
          path: `pages/authentication-success/authentication-success?mem_uid=${this.data.memUid}`
        })
        this.setData({
          confirmModelShow: false
        })
        this.gotoNext()
        break
      // 同656006  参数后台 带入路径  非后台code
      case '66666':
        console.log('去商城认证')
        globalData.needReload = true
        const path = this.data.path ? decodeURIComponent(this.data.path) : ''
        wx.navigateToMiniProgram({
          appId: this.data.appId,
          // path: `pages/authentication-success/authentication-success?mem_uid=${this.data.memUid}`
          path: path
        })
        this.setData({
          confirmModelShow: false
        })
        this.gotoNext()
        break
      case '666661':
        this.gotoNext()
        break
      case '666662':
        // eslint-disable-next-line promise/catch-or-return
        getPartnerContactApi({ msgType: 1 }).then(() => {
          this.onPartnerAccountDisabledClose()
        }).catch((e) => {
          wx.showToast(e.message)
        })
        break
      default:
        wx.redirectTo({
          url: '/pages/load/load'
        })
        break
    }
  },
  // 不同情况展示弹框处理
  dealRole (code) {
    switch (code) {
      // 企业未创建或发布小程序商城
      case '656005':
        this.setData({
          dialogText: '您的企业还未开启商城小程序，商城开通后即可完整使用WeClient Pro功能',
          sureText: '我知道了',
          confirmModelShow: true,
          isSingle: true,
          isAuthent: true
        })
        break
      // 无导购组织
      case '652002':
        this.setData({
          dialogText: '未识别到您的导购组织，请向企业管理员申请导购身份后即可产生和查看业务数据',
          sureText: '我知道了',
          confirmModelShow: true,
          isSingle: true,
          isAuthent: true
        })
        break
      // 导购未完成商城身份认证
      case '656006':
        this.setData({
          dialogText: '您还未完成导购身份认证，认证后才可准确记录您的推广业绩',
          sureText: '去认证',
          confirmModelShow: true,
          isSingle: true,
          isAuthent: false
        })
        break
      case '66666':
        this.setData({
          dialogText: '您还需要认证微信账号',
          dialogTwoText: `如未认证您的${this.data.ispartner ? '佣金订单' : '推广业绩'}将无法正常统计`,
          sureText: '去认证',
          confirmModelShow: true,
          isSingle: true,
          isAuthent: false
        })
        break
      case '666661':
        this.setData({
          dialogText: '您的合伙品牌已经关闭分销业务，您的分销订单将不再产生新的佣金，已产生的佣金将继续结算直至完成提现',
          sureText: '我知道了',
          confirmModelShow: true,
          isSingle: false,
          isAuthent: true,
          ispartner: globalData.identityType === Enums.identityType.partner
        })
        break
      // case '666662':
      //   this.setData({
      //     dialogText: '您当前的品牌合伙人账号已禁用，请联系品牌导购恢复账号',
      //     sureText: '联系导购',
      //     confirmModelShow: true,
      //     isSingle: false,
      //     isPartner: globalData.identityType === Enums.identityType.partner
      //   })
      //   break
      default:
        wx.redirectTo({
          url: '/pages/load/load'
        })
        break
    }
  }
})
