import { getCurrentPageUrl, getCurrentPageUrlWithArgs } from './index'

const { host, apiRoot, enterpriseid } = require('../../lib/config')

const app = getApp({ allowDefault: true })

let requestQueue = []

const request = ({ url, method }, options) => {
  return new Promise((resolve, reject) => {
    console.log('customer_login_token ========' + app.globalData.customer_login_token)
    // todo: 会存在token未过期 但是却没有token情况
    if (app.globalData.customer_login_token) {
      // todo:检查checksession是否到期
      wx.checkSession({
        success () {
          // session_key 未过期，并且在本生命周期一直有效
          console.log('session_key 未过期')
          requestReturn({ url, method }, options, resolve, reject)
        },
        fail () {
          // session_key 已经失效，需要重新执行登录流程
          console.log('session_key 已经失效')
          dealQueue({ url, method }, options, resolve, reject)
        }
      })
    } else {
      // 如果没有token，直接登录,会存在多个请求同时触发登录的情况，这里需要做一个队列
      dealQueue({ url, method }, options, resolve, reject)
    }
  })
}
/**
 * @description: 先登录后处理多并发请求逻辑
 * @return {*}
 */
const dealQueue = ({ url, method }, options, resolve, reject) => {
  requestQueue.push({ url, method, options, resolve, reject })
  if (requestQueue.length === 1) {
    getToken({ url, method }, options).then(() => {
      requestQueue.forEach((item) => {
        requestReturn(item, item.options, item.resolve, item.reject)
      })
      requestQueue = []
    }).catch((err) => {
      console.log(err, '====Err')
    })
  }
}
/**
 * @description: 请求return的数据
 * @return {*}
 */
const requestReturn = ({ url, method }, options, resolve, reject) => {
  console.log(url, '====', options)
  let path = `${host}${apiRoot}${url}`
  let header = {
    CustomerToken: app.globalData.customer_login_token,
    identifier: 1
  }

  wx.request({
    url: path,
    data: options,
    method: method,
    header: header,
    success: ({ data, statusCode }) => {
      if (parseInt(statusCode) === 200) {
        resolve(data)
      } else {
        reject(data)
      }
    },
    fail: (err) => {
      console.log(err, '====Err')
      reject(err)
    }
  })
}
/**
 * 获取首次打开小程序的页面,用当前参数刷新该页面。
 */
const reloadPage = () => {
  const currentPage = getCurrentPageUrl()
  const currentPageWithArgs = getCurrentPageUrlWithArgs()
  console.log('所在页面=>', currentPage, currentPageWithArgs)
  wx.redirectTo({
    url: currentPageWithArgs
  })
}
/**
 * 调用登录方法
 * @param {*} data
 * @param {*} option
 */
const getToken = (data, option) => {
  return new Promise((resolve, reject) => {
    wx.login({
      success: (res) => {
        let enterprise_id = Number(enterpriseid)
        const code = res.code
        if (code && enterprise_id) {
          let path = `${host}${apiRoot}reserve/login`
          // todo: 使用wx.request
          wx.request({
            url: path,
            data: {
              code: code,
              enterprise_id: enterprise_id
            },
            method: 'post',
            success: ({ data, statusCode }) => {
              if (parseInt(statusCode) === 200) {
                app.globalData.customer_login_token = data.data.token
                resolve(true)
              } else {
                reject(data)
                wx.showModal({
                  title: '提示',
                  content: '登陆失败，请重试',
                  showCancel: false,
                  success (res) {
                    if (res.confirm) {
                      reloadPage()
                    }
                  }
                })
              }
            },
            fail: (err) => {
              console.log(err, '====Err')
              reject(err)
              wx.showModal({
                title: '提示',
                content: '登陆失败，请重试',
                showCancel: false,
                success (res) {
                  if (res.confirm) {
                    reloadPage()
                  }
                }
              })
            }
          })
        }
      },
      fail (res) {
        console.error(res, '====失败')
        reject(res)
        wx.showModal({
          title: '提示',
          content: '登陆失败，请重试',
          showCancel: false,
          success (res) {
            if (res.confirm) {
              reloadPage()
            }
          }
        })
      }
    })
  })
}
export { request }
