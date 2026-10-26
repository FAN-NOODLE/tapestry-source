import { customerLoginApi } from '@/api/store-invite'

const { host, apiRoot } = require('../lib/config')

const app = getApp()
const request = ({ url, method }, options) => {
  return new Promise((resolve, reject) => {
    console.log(url, '====', options)
    let path = `${host}${apiRoot}${url}`
    let header = {
      CustomerToken: app.globalData.customer_token || wx.getStorageSync('customer_token') || '',
      identifier: 1
    }
    wx.request({
      url: path,
      data: options,
      method: method,
      header: header,
      success: ({ data, statusCode }) => {
        if (parseInt(statusCode) === 200) {
          console.log('登录成功', data.code)
          if ([652006, 652004].indexOf(data.code) !== -1) {
            // customerLoginApi()
            console.log(url, method, options, '000999')
            getToken({ url, method }, options)
            return
          }
          // todo: 如果是系统错误调到缺省页
          if (data.message === '系统错误') {
            wx.navigateTo({
              url: '/pages/preserve-page/preserve-page'
            })
            return
          }
          resolve(data)
        } else {
          // todo：上报异常
          wx.report({
            title: 'coach小程序请求失败',
            url: path,
            time: new Date().getTime(),
            header: header,
            data: options,
            errMsg: data
          })
          reject(data)
        }
      },
      fail: (err) => {
        // todo:判断微信请求状态码是否在600000-600006
        if (err.errno >= 600000 && err.errno <= 600006) {
          wx.showToast('请求超时,请检查您的网络连接,稍后重试~')
        } else {
        // todo：上报异常
          wx.report({
            title: 'coach小程序请求失败',
            url: path,
            time: new Date().getTime(),
            header: header,
            data: options,
            errMsg: err
          })
        }
        reject(err)
      }
    })
  })
}

const getToken = (data, option) => {
  wx.login({
    success: res => {
      console.log(res, '9999999')
      // this.getRegisterHash(res.code, option)
      const enterprise_id = app.globalData.customer_enterprise_id
      const code = res.code
      if (code && enterprise_id) {
        // eslint-disable-next-line promise/catch-or-return
        customerLoginApi({ code, enterprise_id }).then(
          res => {
            // app.globalData.customer_token = res.data.token
            // wx.setStorageSync('customer_token', res.data.token)
            // request(data, option)
            try {
              console.error(app.globalData)
              console.log('customerLoginApi接口获取token的res=====》', res)
              console.log('customerLoginApi接口获取token的tokentokentoken=====》', res.data.token)
              app.globalData.customer_token = res.data.token
              wx.setStorageSync('customer_token', res.data.token)
              // request(data, option)
              const customerInfo = app.globalData.customerOption
              let url = `/${customerInfo.path}?`
              const keys = Object.keys(customerInfo.query)
              keys.forEach((item, index) => {
                if (index !== 0) {
                  url += `&${item}=${customerInfo.query[item]}`
                } else {
                  url += `${item}=${customerInfo.query[item]}`
                }
              })
              wx.redirectTo({
                url
              })
            } catch (error) {
              console.error(error)
              console.error(app.globalData)
            }
          }
        )
      }
    },
    fail (res) {
      console.error(res, '====')
    }
  })
}
export { request }
