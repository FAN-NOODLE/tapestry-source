import {getCodeApi} from '@/api/login-api'
import {changeSmsApi} from '@/api/store-api'

Component({
  properties: {
    show: {
      type: Boolean,
      value: false
    }
  },
  data: {
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
    }
  },
  observers: {
    show () {
      this.setData({
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
        }
      })
    }
  },

  /*    系统事件开始   */
  created () {
  },
  attached () {
  },
  ready () {
    this.preCode = ''
  },
  /*    系统事件结束   */

  methods: {
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
        phone, code
      })
    },
    // 验证码框获得焦点
    focusIpt () {
      if (this.data.phone.error || this.verifyPhone()) return
      if (!this.data.code.get) {
        wx.showToast('请先获取验证码')
        return
      }
      this.setData({
        'code.focus': true,
        'code.error': false,
        'code.value':
          this.data.code.error &&
          this.data.code.value.length === 4
            ? ''
            : this.data.code.value,
        'code.msg': ''
      })
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
      const { error, value, get } = this.data.code
      if (!error && value.length === 4 && get) this.verify()
    },
    // 校验手机号
    verifyPhone () {
      const phone = this.data.phone.value
      let msg = phone.length === 0 ? '手机号不能为空' : ''
      !msg &&
      (msg = /^1(2|3|4|5|6|7|8|9)\d{9}$/.test(phone) ? '' : '手机号码格式不正确')
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
        const res = await getCodeApi({
          mobile: this.data.phone.value,
          module: 'CHECK_SMS_MOBILE'
        })
        this.setData({
          'code.click': true,
          'code.get': true
        })
        if (res.code === 200) {
          wx.showToast('已发送验证码')
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
    verify () {
      const { phone: { error: phoneError, value: sms_mobile }, code: { error: codeError, value: code } } = this.data
      const condition = [phoneError, codeError]

      if (condition.every(item => !item)) {
        wx.showNavigationBarLoading()

        changeSmsApi({
          code,
          sms_mobile
        }).then(res => {
          if (res.code === 200) {
            this.setData({
              complete: true
            })
          } else {
            this.setData({
              complete: false
            })
            // | 652001 | 验证码不正确|
            // | 652003 | 此用户非店长|
            // | 700017 | 验证码已经过期 |
            this.code = res.code
            if (res.code === 652001 || res.code === 700007 || res.code === 700017) {
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
          wx.hideNavigationBarLoading()
        }).catch(err => {
          this.setData({
            complete: false
          })
          wx.hideNavigationBarLoading()
          wx.showToast(err.message)
        })
      }
    },
    sure () {
      const { phone: { error: phoneError, value: sms_mobile }, code: { error: codeError, value: code } } = this.data
      const condition = [phoneError, codeError, !sms_mobile, !code]

      if (condition.every(item => !item)) {
        wx.showToast('修改手机号码成功')
        wx.event.$broadcast('modal-slot-confirm')
      }
    },
    confirm () {
      const { phone: { value: sms_mobile } } = this.data
      this.triggerEvent('confirm', { sms_mobile })
    },
    cancel () {
      wx.event.$broadcast('modal-slot-close')
    },
    close () {
      this.triggerEvent('close')
    }
  }
})
