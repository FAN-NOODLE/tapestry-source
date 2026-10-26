import { getCodeApi, loginApi } from '@/api/login-api'
import { debounce } from '../../utils/tools'
import Enums from '@/utils/enums'

const { showToast, globalData, choosoeEnterpri } = getApp()
// const { env } = require('../../lib/config')

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    title: 'login',
    phone: {
      msg: '', // '手机号不能为空' || '手机号格式错误' || '手机号码错误，查无此账号'
      value: '', // 手机号值
      fouce: false, // 是否获得焦点
      error: false // 是否错误
    },
    code: {
      focus: false, // 是否获得焦点
      msg: '', // 错误提示
      value: '', // code值
      error: '', // 是否错误
      get: false, // 是否已经点击获取验证码
      click: false, // 是否可点击
      ms: 60 // 倒计时
    },
    version: globalData.version
  },

  // 生命周期函数--监听页面加载
  onLoad () {
    this.preCode = ''
  },
  // 输入手机号获得焦点
  phoneFoucs () {
    const { phone, code } = this.data
    // 700009 一分钟内多次点击特殊处理
    phone.value =
      this.data.phone.error &&
      this.data.phone.value.length === 11 &&
      this.getCodeCode !== 700009
        ? ''
        : this.data.phone.value
    phone.msg = ''
    phone.error = false
    phone.fouce = true
    if (this.code !== 652001 && this.code !== 700007) {
      code.error = false
      code.value = ''
      code.msg = ''
    }
    this.setData({
      phone,
      code
    })
  },
  // 验证码框获得焦点
  focusIpt () {
    if (this.data.phone.error || this.verifyPhone()) return
    if (!this.data.code.get) {
      showToast('请先获取验证码')
      return
    }
    this.setData({
      'code.focus': true,
      'code.error': false,
      'code.value':
        this.data.code.error && this.data.code.value.length === 6
          ? ''
          : this.data.code.value,
      'code.msg': ''
    })
    // debugger
  },
  // input设置手机号
  inputGetPhone (e) {
    this.setData({
      'phone.value': e.detail.value
    })
  },
  // blur校验手机号
  getPhone () {
    this.verifyPhone()
    this.setData({
      'phone.fouce': false
    })
  },
  // input获取code
  getInputCode (e) {
    // debugger
    const { code } = this.data
    if (this.preCode.length > e.detail.value.length) {
      code.msg = ''
      code.error = false
    }
    this.setData({
      'code.value': e.detail.value,
      'code.msg': code.msg,
      'code.error': code.error
    })
    this.preCode = e.detail.value
    this.login()
  },

  // 选择企业前 先选择选择 角色
  selectRole () {
    wx.redirectTo({
      url: '/pages/select-role/select-role'
    })
  },
  // 登陆
  login: debounce(async function () {
    const { error, value, get } = this.data.code
    if (!error && value.length === 6 && get) {
      try {
        wx.showLoading({
          title: '登录中',
          mask: true
        })
        globalData.identityType = '' // 清身份信息
        let res = await loginApi({
          code: this.data.code.value,
          mobile: this.data.phone.value
        })
        // res = loginData
        wx.hideLoading()
        if (res.code === 200) {
          // 用户是否选择角色1:导购2:合伙人3选择
          // debugger
          if (res.data.user_type + '' === '1') {
            globalData.token = res.data.token
            wx.setStorageSync('token', res.data.token)
            globalData.isFirstLogin = res.data.is_first_login === 1
            showToast('登录成功', true)
            this.enterprise(res)
            globalData.identityType = Enums.identityType.salesman
          } else if (res.data.user_type + '' === '2') {
            globalData.token = res.data.partner_token
            wx.setStorageSync('token', res.data.partner_token)
            globalData.isFirstLogin = res.data.is_first_login === 1
            globalData.identityType = Enums.identityType.partner
            showToast('登录成功', true)
            res.data.enterprise = JSON.parse(
              JSON.stringify(res.data.partner_enterprise_list)
            ) // 改动小 后期优化
            this.enterprise(res)
          } else if (res.data.user_type + '' === '3') {
            // token备份 去筛选页面 选择
            globalData.tokenObj = res
            this.selectRole()
          }

          // 登录200 必有token
        } else {
          // | 652001 | 验证码不正确|
          // | 652003 | 此用户非店长|
          // | 700017 | 验证码已经过期 |
          this.code = res.code
          if (
            res.code === 652001 ||
            res.code === 700007 ||
            res.code === 700017
          ) {
            this.setData({
              'code.msg': res.message,
              'code.error': !!res.message
            })
          } else {
            this.setData({
              'phone.msg': res.message,
              'phone.error': !!res.message
            })
          }
        }
      } catch (error) {
        console.error(error)
        wx.hideLoading()
      }
    }
  }, 500),
  // 选择企业
  // 坑 存在多个注意统一修改
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

        globalData.roleMark =
          currentEnterprise &&
          currentEnterprise.guide_roles &&
          currentEnterprise.guide_roles[0].role_mark
        // debugger
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
  // 校验手机号
  verifyPhone () {
    const phone = this.data.phone.value
    let msg = phone.length === 0 ? '手机号不能为空' : ''
    !msg &&
      (msg = /^1(2|3|4|5|6|7|8|9)\d{9}$/.test(phone) ? '' : '手机号格式错误')
    this.setData({
      phone: {
        msg: msg, // '手机号不能为空' || '手机号格式错误' || '手机号码错误，查无此账号'
        value: phone,
        error: !!msg
      }
    })
    return !!msg
  },
  // 校验验证码
  verifyCode () {
    let msg = this.data.code.msg
    if (this.data.code.value.length === 0) {
      msg = '验证码不能为空'
    } else if (this.data.code.value.length < 6) {
      msg = '验证码格式错误'
    }
    this.setData({
      'code.msg': msg,
      'code.focus': false,
      'code.error': !!msg,
      'code.value': this.code === 700005 ? '' : this.data.code.value
    })
  },
  // 获取验证码
  async getCode () {
    if (this.data.code.click) return
    // 若手机号错误无法获取验证码
    if (this.data.phone.error || this.verifyPhone()) {
      return
    }
    this.data.code.click = true
    this.setData({
      'code.focus': false,
      'code.error': false,
      'code.value': '',
      'code.msg': ''
    })
    try {
      // TODO 测试方法需删除
      // const getCodeApi = () => {
      //   return ({code: 200})
      // }
      const res = await getCodeApi({
        mobile: this.data.phone.value,
        corp_id: globalData.corp_id
      })
      this.setData({
        'code.click': true,
        'code.get': true
      })
      if (res.code === 200) {
        showToast('已发送验证码')
        this.timer = setInterval(() => {
          this.setData({
            'code.ms': this.data.code.ms - 1
          })
          if (this.data.code.ms <= 0) {
            clearInterval(this.timer)
            this.setData({
              'code.click': false,
              'code.ms': 60
            })
          }
        }, 1000)
      } else {
        this.getCodeCode = res.code
        // showToast(res.message)
        this.setData({
          'code.click': false,
          'code.get': false,
          'phone.error': !!res.message,
          'phone.msg': res.message
        })
      }
    } catch (error) {
      console.error(error)
      this.setData({
        'code.click': false,
        'code.get': false
      })
    }
  },
  goRecords (e) {
    let { type } = e.currentTarget.dataset
    wx.navigateTo({
      url: `/pages/records-link/records-link?type=${type}`
    })
  }
})
