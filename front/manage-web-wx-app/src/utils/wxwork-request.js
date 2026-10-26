import { sensorsLoginInit, setSensorsProfile } from '../lib/sensors/index'

const { host, apiRoot, env } = require('../lib/config')

const newApiRoot = '/gw-assistant/app/brand/'

/** 错误码 */
const ERRCODE = [652002]

/**
 * 发起请求
 * @param path 请求路径
 * @param options 请求配置参数，例如method、data...
 * @param _app 全局配置信息
 */
const request = (path, options, _app) => {
  const app = getApp() || _app
  return new Promise((resolve, reject) => {
    // 缓存需要企业id的接口
    if (options.data && options.data.enterprise_id === '') {
      _subscribeEnterpriseRefresh(app, () => {
        _requestFunc(resolve, reject, { app, path, options })
      })
    }
    // 判断是否需要校验session, load 空接口登录就是走这个方法、非load页面进入也是走的这个接口
    if (!app.globalData.checkSessionOK) {
      _ckSession(resolve, reject, { app, path, options })
      return
    }
    _requestFunc(resolve, reject, { app, path, options })
  })
}

const _ckSession = (resolve, reject, { app, path, options }) => {
  wx.qy.checkSession({
    success (res) {
      let token = wx.getStorageSync('token')
      if (app.globalData.token) {
        app.globalData.checkSessionOK = true
        app.globalData.token = token
        _requestFunc(resolve, reject, { app, path, options })
      } else {
        _checkFail(resolve, reject, { app, path, options })
      }
    },
    fail () {
      app.globalData.checkSessionOK = false
      app.globalData.token = ''
      wx.removeStorageSync('token')
      _checkFail(resolve, reject, { app, path, options })
    }
  })
}

/** session 过期失败处理，重新走登录逻辑 */
const _checkFail = (resolve, reject, { app, path, options }) => {
  if (!app.globalData.isCanLogin) {
    app.globalData.isCanLogin = true
    _login(resolve, reject, { app, path, options })
  } else {
    _subscribeTokenRefresh(app, () => {
      _requestFunc(resolve, reject, { app, path, options })
    })
  }
}

/** 企业微信登录 */
const _login = (resolve, reject, { app, path, options }) => {
  wx.qy.login({
    success: (res) => {
      _requestLogin(resolve, reject, { app, path, options }, res.code)
    },
    fail (res) {
      app.globalData.isCanLogin = false
      app.globalData.checkSessionOK = false
      wx.showToast('请求超时,请检查您的网络连接,稍后重试~')
      reject()
    }
  })
}

/** 企业微信登录成功处理 - success */
const _requestLogin = (resolve, reject, { app, path, options }, code) => {
  // 定义基础信息
  const corpArr = {
    test: 'ww29202b3a45270ff3',
    prod: 'wx1b5619d5190a04e4'
  }
  let _path = 'common/qy-login'
  let $options = {
    method: 'POST',
    header: {
      'content-type': 'application/json',
      Accept: 'application/json',
      identifier: 1
    },
    data: {
      code: code,
      corp_id: app.globalData.corp_id || corpArr[env]
    }
  }

  // 向后端发起登录请求, 拿企业微信code 换去用户信息
  wx.request({
    url: `${host}${apiRoot}${_path}`,
    data: $options.data,
    method: $options.method,
    header: $options.header,
    success: ({ data, statusCode }) => {
      app.globalData.isCanLogin = false
      if (parseInt(statusCode) === 200 && data.code === 200) {
        // 登陆成功
        app.globalData.token = data.data.token
        app.globalData.mem_uid = data.data.mem_uid
        app.globalData.brand_code = data.data.brand_code
        wx.setStorageSync('token', data.data.token)
        console.log('data => ', _path, data)
        // todo:增加userid
        app.globalData.userid = data.data.userid
        if (Array.isArray(data.data.enterprise) && data.data.enterprise.length === 1) {
          // 存储用户头像
          // app.globalData.guide_avatar = data.data.avatar
          // 是否是beta用户
          app.globalData.is_beta_user = 1
          // 存储用户可选企业列表
          app.globalData.enterprise = data.data.enterprise
          // 存储缓存key
          app.login_cache_key = data.data.login_cache_key
          // 存储企业信息
          app.globalData.is_manager = data.data.is_manager
          // 调用选择企业
          const currentEnterprise = data.data.enterprise[0]
          app.globalData.logo = currentEnterprise.logo
          app.globalData.choiceEnterprise = currentEnterprise
          // 存储当前所在企业
          app.globalData.currentEnterprise = currentEnterprise
          // 存储当前角色(PS：店长、店员角色会在选择门店时调用member/detail接口重新赋值)
          app.globalData.roleMark = currentEnterprise && currentEnterprise.guide_roles && currentEnterprise.guide_roles[0].role_mark
          // 存储当前门店
          app.globalData.storeInfo = {
            store_city: data.data.store_city,
            store_code: data.data.store_code,
            store_inventory_name: data.data.store_inventory_name,
            store_name: data.data.store_name,
            store_phone: data.data.store_phone
          }
        }

        // 更改登录状态
        app.globalData.sessionKeyLogin = false
        app.globalData.checkSessionOK = true
        app.globalData.isCanLogin = false
        console.log('登录地址', path)

        // 如果非load 页面进入
        if (path) {
          _jumpPage()
        } else {
          _onRrefreshed(app)
          // 返回登陆数据
          if (ERRCODE.includes(data.code)) {
            wx.redirectTo({
              url: `/pages/unorganized/unorganized?code=${data.code}&appId=${data.data && data.data.shop && data.data.shop.app_id}&memUid=${data.data.mem_uid}`
            })
            return
          }
          resolve(data, statusCode)
        }
        // 神策数据上报 -- 开始
        sensorsLoginInit(data.data.userid)
        // 神策数据上报 -- 设置用户属性
        setSensorsProfile({
          store_name: data.data.store_name,
          store_code: data.data.store_code || data.data.store_name
        })
      } else {
        app.globalData.isCanLogin = false
        app.globalData.checkSessionOK = false
        wx.showToast('请求超时,请检查您的网络连接,稍后重试~')
        // todo：上报异常
        wx.report({
          title: 'coach小程序企业微信登录失败',
          url: `${host}${apiRoot}${_path}`,
          time: new Date().getTime(),
          header: $options.header,
          data: $options.data,
          errMsg: data
        })
        reject(data)
      }
    },
    fail: (err) => {
      // todo：上报异常
      wx.report({
        title: 'coach小程序企业微信登录失败',
        url: `${host}${apiRoot}${_path}`,
        time: new Date().getTime(),
        header: $options.header,
        data: $options.data,
        errMsg: err
      })
      reject(err)
    }
  })
}

const _requestFunc = (resolve, reject, { app, path, options }, again) => {
  // 拼接请求路径
  let url = ''
  if (options.apiVersion && options.apiVersion === 'brand') {
    url = `${host}${newApiRoot}${path}`
  } else {
    url = `${host}${apiRoot}${path}`
  }

  // 删除apiVersion，请求不需要apiVersion
  delete options.apiVersion

  // 参数格式化 - 单独处理
  const $options = _getOption(options, app)

  // load页面登陆接口拦截
  if (!path) return

  // 调用接口包含企业id的重置企业id的值
  if (options.data && options.data.enterprise_id === '') {
    options.data.enterprise_id = app.globalData.currentEnterprise.enterprise_id
  }

  wx.request({
    url: url,
    data: $options.data,
    method: $options.method,
    header: {
      Token: app.globalData.token
        ? `${app.globalData.token}`
        : wx.getStorageSync('token'),
      ...$options.header
    },
    success: ({ data, statusCode, header }) => {
      console.log('data => ', path, data)
      // header 头中返回授权失败或者授权时返回失败时需要重新登陆
      _checkSessionKey(app, header, data.code)
      if (parseInt(statusCode) === 200) {
        resolve(data)
      } else {
        wx.hideLoading()
        wx.showToast('请求超时,请检查您的网络连接,稍后重试~')
        // todo：上报异常
        wx.report({
          title: 'coach小程序请求失败',
          url: url,
          time: new Date().getTime(),
          header: JSON.stringify({
            Token: app.globalData.token
              ? `${app.globalData.token}`
              : wx.getStorageSync('token'),
            ...$options.header
          }),
          data: $options.data,
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
          url: url,
          time: new Date().getTime(),
          header: {
            Token: app.globalData.token
              ? `${app.globalData.token}`
              : wx.getStorageSync('token'),
            ...$options.header
          },
          data: $options.data,
          errMsg: err
        })
      }
      reject(err)
    }
  })
}

/** 验证session 是否过期 */
const _checkSessionKey = (app, header, code) => {
  if (
    ((header && header['Auth-Status'] === 'false') ||
      parseInt(code) === 900001) &&
    !app.globalData.sessionKeyLogin
  ) {
    app.globalData.sessionKeyLogin = true

    var promise = new Promise((resolve, reject) => {
      app.globalData.checkSessionOK = false
      wx.removeStorageSync('token')
      app.globalData.isCanLogin = false
      if (app.globalData.sessionKeyLogin) {
        _checkFail(resolve, reject, { app })
      }
    })
    return promise
  }
}

// 调用接口不含enterprise_id的接口缓存
const _subscribeTokenRefresh = (app, cb) => {
  if (!app.globalData.refreshSubscribers) {
    app.globalData.refreshSubscribers = []
  }
  cb && app.globalData.refreshSubscribers.push(cb)
}
// 调用接口不含enterprise_id的接口执行
const _onRrefreshed = (app) => {
  app.globalData.refreshSubscribers.map((cb) => {
    cb && cb()
  })
}

// 调用接口含enterprise_id的接口缓存
const _subscribeEnterpriseRefresh = (app, cb) => {
  if (!app.globalData.refreshEnterpriseSubscribers) {
    app.globalData.refreshEnterpriseSubscribers = []
  }
  cb && app.globalData.refreshEnterpriseSubscribers.push(cb)
}

/** 更具特定规则处理options */
const _getOption = (options, app) => {
  // 初始化options
  let $options = {
    method: 'GET',
    header: {
      'content-type': 'application/json',
      Accept: 'application/json',
      identifier: 1
    },
    data: undefined,
    needToken: true,
    ...options
  }

  // 添加企业hash
  if (env === 'test' && app.globalData && app.globalData.currentEnterprise) {
    $options.header['Enterprise-Hash'] = app.globalData.currentEnterprise.hash
  }

  // 处理为空数据
  if ($options.data) {
    $options.data = Object.keys($options.data).reduce((result, key) => {
      if (
        $options.data[key] !== undefined &&
        $options.data[key] !== null &&
        (!Array.isArray($options.data[key]) || $options.data[key].length)
      ) {
        result[key] = $options.data[key]
      }
      return result
    }, {})
  }

  // 添加store_id
  const store_id = ($options.data && $options.data.store_id) || app.globalData.store.store_id || 0

  $options.data = store_id ? { ...$options.data, store_id } : $options.data

  return $options
}

const createModuleRequest = (module) => {
  // var self = this;
  return (path, ...rest) => {
    if (module) {
      return request(`${module}/${path}`, ...rest)
    } else {
      return request(`${path}`, ...rest)
    }
  }
}

/** 处理跳转逻辑，如果非首页进入，登录之后需要做重定向 */
const _jumpPage = () => {
  // 获取跳转路径和参数
  const { path = '', query = {} } = wx.getLaunchOptionsSync() || {}

  // tarBar 页面
  let canSwitchPage = [
    '/pages/task/task',
    '/pages/compass-data/compass-data',
    '/pages/my/my',
    '/pages/customer-list/customer-list'
  ]

  let url = (path.startsWith('/') ? '' : '/') + path
  let type = 'redirectTo'
  const page = (getCurrentPages() || []).pop() || {}
  // 打开的是tabBar页面
  if (canSwitchPage.includes(url)) {
    type = '/' + page.route === url ? 'reLaunch' : 'switchTab'
  }
  // todo: 不管是tabBar页面还是子页面，都需要把参数带上
  if (query) {
    // 打开的是子页面
    const queryKey = Object.keys(query)
    let str = ''
    queryKey.forEach((item) => {
      str += `${item}=${query[item]}&`
    })
    url = url + '?' + str.substring(0, str.length - 1)
  }
  console.log('登录后页面跳转>>', type, url)
  wx[type]({
    url,
    complete (msg) {
      console.log('登录后页面跳回调>>>', msg)
    }
  })
}

export { request, createModuleRequest }
