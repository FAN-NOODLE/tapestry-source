// 获取全局应用程序实例对象

import { saveRegisterApi, getCodeApi, getMobileApi } from '@/api/login-api'

const { host, apiRoot, env } = require('@/lib/config')
// const app = getApp();
const { globalData, showToast } = getApp()

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    title: 'register',
    switchChecked: false,
    status: 1,
    worker: {},
    partner_info: {},
    phone: {
      msg: '', // '手机号不能为空' || '手机号格式错误' || '手机号码错误，查无此账号'
      value: '', // 手机号值
      error: false // 是否错误
    },
    code: {
      msg: '', // 错误提示
      value: '', // code值
      error: '', // 是否错误
      get: false, // 是否已经点击获取验证码
      click: false, // 是否可点击
      ms: 60 // 倒计时
    },
    iv: '', // 后台微信昵称授权后需要此字段，phone授权不传--赵鹏飞
    encryptedData: '', // 后台微信昵称授权后需要此字段
    timer: null,
    defaultAvatar: 'https://img-crs.vchangyi.com/store-manage-wxapp/sy_register_photo.png'
  },

  // 生命周期函数--监听页面加载
  onLoad (option) {
    // TODO: onLoad
    console.log(globalData.version)
    if (Object.keys(option).length > 0) {
      this.setData({
        status: Number(option.status),
        worker: globalData.register ? globalData.register.worker : {},
        'phone.value': globalData.register ? globalData.register.partner_info.partner_mobile : '',
        partner_info: globalData.register ? globalData.register.partner_info : {},
        partner_id: globalData.register && Number(option.status) === 2 ? globalData.register.partner_info.id : ''
      })
    }
  },
  onShow () {

  },
  // input设置手机号
  inputGetPhone (e) {
    this.setData({
      'phone.value': e.detail.value,
      'phone.error': false,
      'phone.msg': ''
    })
  },
  // input获取code
  getInputCode (e) {
    this.setData({
      'code.value': e.detail.value,
      'code.msg': '',
      'code.error': false
    })
  },
  // 合伙人称呼
  getPartnerName (e) {
    this.setData({
      'partner_info.partner_name': e.detail.value.replace(/\s+/g, '')
    })
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
    } else if (this.data.code.value.length < 4) {
      msg = '验证码格式错误'
    }
    this.setData({
      'code.msg': msg,
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
      'code.error': false,
      'code.value': '',
      'code.msg': ''
    })
    try {
      // TODO 测试方法需删除
      // const getCodeApi = () => {
      //   return ({code: 200})
      // }
      const res = await getCodeApi({ mobile: this.data.phone.value, module: 'PARTNER_REGISTER' })
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
  // 注册
  sendRegister () {
    this.verifyCode()
    if (this.data.phone.error) {
      return
    }
    if (this.data.code.error && this.data.status === 1) {
      return
    }
    if (!this.data.partner_info.avatar) {
      // showToast('请上传头像')
      // return
      this.data.partner_info.avatar = this.data.defaultAvatar
    }
    if (!this.data.phone.value) {
      showToast('请先填写手机号')
      return
    }
    if (!this.data.partner_info.partner_name) {
      showToast('请先填写称呼')
      return
    }
    const data = {
      avatar: this.data.partner_info.avatar,
      partner_mobile: this.data.phone.value,
      partner_name: this.data.partner_info.partner_name,
      sms_code: this.data.code.value,
      mobile_sensitive: Number(!this.data.switchChecked),
      cache_key: globalData.register.cache_key,
      scene: globalData.register.scene,
      partner_id: this.data.partner_id,
      iv: this.data.iv || '',
      encryptedData: this.data.encryptedData || ''
    }

    saveRegisterApi(data)
      .then(res => {
        if (res.code === 200) {
          console.log(res)
          showToast('注册成功')

          globalData.token = res.data.token
          wx.setStorageSync('token', res.data.token)
          wx.redirectTo({
            url: '/pages/load/load'
          })
        } else {
          // 定义需要特殊处理的code
          showToast(res.message || '注册失败')
        }
      })
      .catch(err => {
        console.error('1', err)
      })
  },
  // 获取用户信息
  getUserInfo (e) {
    this.setData({
      'partner_info.partner_name': e.detail.userInfo.nickName,
      encryptedData: e.detail.encryptedData,
      iv: e.detail.iv
    })
    if (!this.data.partner_info.avatar) {
      this.setData({
        'partner_info.avatar': e.detail.userInfo.avatarUrl
      })
    }
  },

  // 授权手机号
  bindGetPhoneNumber (e) {
    // 如果是组件内的button走授权，将组件内的授权按钮事件对象原封不动的带过来
    if (e.detail && e.detail.button_e) {
      e = e.detail.button_e
    }
    if (e.detail && e.detail.errMsg === 'getPhoneNumber:ok') {
      wx.showNavigationBarLoading()

      // 回头换成 获取手机号接口
      getMobileApi({
        encryptedData: e.detail.encryptedData,
        iv: e.detail.iv,
        cache_key: globalData.register ? globalData.register.cache_key : ''
      })
        .then(res => {
          wx.hideNavigationBarLoading()
          if (res.code === 200) {
            console.log(res)
            this.setData({
              'phone.value': res.data.data.phoneNumber
            })
          } else {
            // 定义需要特殊处理的code
            let errorCode = {
              900001: () => {
                showToast('本次授权已过期')
              }
            }
            wx.codeException(res, errorCode)
          }
        })
        .catch(err => {
          wx.hideNavigationBarLoading()
          // if (Number(err.code) === 900001) {
          //   showToast("本次授权已过期");
          //   return;
          // }
          console.error('我的-getPhoneNumber:', err)
        })
    }
  },
  // 手机号展示给导购开关
  switchChange () {
    this.setData({
      switchChecked: !this.data.switchChecked
    })
  },
  // 上传相关
  // (1)选择本地、拍照上传图片
  chooseImageTap: function (e) {
    let _this = this
    wx.showActionSheet({
      itemList: ['本地上传', '拍照上传'],
      itemColor: '#000000',
      success: function (res) {
        console.log('打开ActionSheet成功')
        // if (!res.cancel) {
        //   if (res.tapIndex === 0) {
        //     _this.chooseWxImage('album')
        //   } else if (res.tapIndex === 1) {
        //     _this.chooseWxImage('camera')
        //   }
        // }
        if (res.tapIndex === 0) {
          _this.chooseWxImage('album')
        } else if (res.tapIndex === 1) {
          _this.chooseWxImage('camera')
        }
      },
      fail: function (e) {
        console.log('打开ActionSheet失败')
        console.error(e)
      }
    })
  },
  // (2)图片选择后的回调方法
  chooseWxImage: function (type) {
    let _this = this
    let count = 1
    wx.chooseImage({
      count: count,
      sizeType: ['original', 'compressed'],
      sourceType: [type],
      success: function (res) {
        console.log('选择图片成功')
        wx.showLoading({
          title: '图片上传中...',
          mask: true
        })
        let str = true
        _this.upImg(res.tempFilePaths[0], str)
      },
      fail: function (e) {
        console.log('选择图片失败')
        console.error(e)
      }
    })
  },
  // (3)图片选择后的,上传到服务器
  upImg (path, str) {
    let _this = this
    const header = {
      'content-type': 'multipart/form-data'
      // 'Enterprise-Hash': 'ecd85280aa135bbd0108dd6aa424565a'// 临时设置
    }
    if ((env === 'beta' || env === 'test_beta' || env === 'test') && globalData && globalData.temporary_hash) {
      header['Enterprise-Hash'] = globalData.temporary_hash
    }
    wx.uploadFile({
      url: host + apiRoot + 'partner-member/upload',
      filePath: path,
      name: 'file',
      header: header,
      success: function (res) {
        console.log('上传图片成功', res)
        let data = JSON.parse(res.data)
        if (res.errMsg === 'uploadFile:ok') {
          _this.setData({
            'partner_info.avatar': data.data.url
          })
          if (str) {
            wx.hideLoading()
          }
        } else {
          wx.hideLoading()
          // showMessage({
          //   title: _this.data.imgUploadFailText,
          //   content: _this.data.plsAgainUploadTex
          // })
        }
      },
      fail: function () {
        console.log('上传图片失败')
        wx.hideLoading()
        // showMessage({
        //   title: _this.data.imgUploadFailText,
        //   content: _this.data.plsAgainUploadTex
        // })
      }
    })
  }
})
