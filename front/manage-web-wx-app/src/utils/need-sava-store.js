import { wxworkLoginApi, choiceEnterpriseApi, checkTokenApi } from '@/api/login-api'
import { getMemberDetailApi } from '@/api/store-api'

const { globalData } = getApp()

// 从非load、login页面进来，得调用此接口存一下所在当前企业、店长店员还得存当前门店
const needSavaStore = async function (store_id) {
  let resolve = null
  let reject = null
  const promise = new Promise((res, rej) => {
    resolve = res
    reject = rej
  })
  const wxwork = wx.getSystemInfoSync().environment === 'wxwork'
  if (globalData.currentEnterprise && globalData.currentEnterprise.enterprise_id !== '') {
    if (parseInt(store_id) && getCurrentPages() && getCurrentPages().length === 1) {
      globalData.roleMark = ''
      globalData.store = {
        store_id: ''
      }
      globalData.memberInfo = {
        isBoss: false,
        memberId: 0
      }
      await recordStore(store_id)
      resolve(true)
    } else {
      resolve(true)
    }
  } else if (wxwork) {
    wxworkLoginApi()
      .then(async res => {
        const params = await saveEnterprise(res)
        await recordEnterprise(params)
        if (parseInt(store_id)) {
          await recordStore(store_id)
          resolve(true)
        } else {
          resolve(true)
        }
      })
      .catch(err => {
        reject(err)
      })
  } else {
    // TODO 微信端目前不发模板消息,后期开发可以优化此处
    const token = wx.getStorageSync('token')
    if (token) {
      checkTokenApi({ token })
        .then(async res => {
          if (res.code === 200) {
            globalData.token = token
            const params = await saveEnterprise(res)
            await recordEnterprise(params)
            if (store_id) {
              await recordStore(store_id)
              resolve(true)
            } else {
              resolve(true)
            }
          } else {
            wx.removeStorageSync('token')
            globalData.token = ''
            wx.redirectTo({
              url: '/pages/login/login'
            })
            reject(false)
          }
        })
        .catch(err => {
          console.error('校验token', err)
          wx.redirectTo({
            url: '/pages/login/login'
          })
          reject(err)
        })
    } else {
      setTimeout(() => {
        wx.redirectTo({
          url: '/pages/login/login'
        })
      }, 1000)
      reject(false)
    }
  }
  return promise
}
// 存储企业信息到全局变量
const saveEnterprise = (res) => {
  let resolve = null
  let reject = null
  const promise = new Promise((res, rej) => {
    resolve = res
    reject = rej
  })
  if (Array.isArray(res.data.enterprise) && res.data.enterprise.length === 1) {
    // 存储用户可选企业列表
    globalData.enterprise = res.data.enterprise
    // 存储缓存key
    globalData.login_cache_key = res.data.login_cache_key
    // 调用选择企业
    const currentEnterprise = res.data.enterprise[0]
    globalData.logo = currentEnterprise.logo
    // 存储当前所在企业
    globalData.currentEnterprise = currentEnterprise
    // 存储当前角色(PS：店长、店员角色会在选择门店时调用member/detail接口重新赋值)
    globalData.roleMark = currentEnterprise && currentEnterprise.guide_roles && currentEnterprise.guide_roles[0].role_mark
    resolve({
      login_cache_key: res.data.login_cache_key,
      enterprise_id: currentEnterprise.enterprise_id
    })
  } else {
    reject(false)
  }
  return promise
}
// 此接口的功能是为了让php记录当前进入的企业
const recordEnterprise = (params) => {
  let resolve = null
  let reject = null
  const promise = new Promise((res, rej) => {
    resolve = res
    reject = rej
  })
  choiceEnterpriseApi(params).then(res => {
    if (res.code === 200) {
      resolve(true)
    } else {
      wx.showToast(res.message)
      reject(false)
    }
  }).catch(err => {
    wx.showToast(err.message)
    reject(false)
  })
  return promise
}
// 此接口的功能是为了存储当前门店信息
const recordStore = (store_id) => {
  let resolve = null
  let reject = null
  const promise = new Promise((res, rej) => {
    resolve = res
    reject = rej
  })
  getMemberDetailApi({
    store_id: store_id,
    login_cache_key: globalData.login_cache_key
  }).then(res => {
    if (res.code === 200) {
      if (res.code === 200) {
        const { is_boss, id, store_sun_code, guide_role, guide_permissions } = res.data
        globalData.memberInfo = {
          isBoss: is_boss === 1,
          memberId: id
        }
        // 选择门店时，店员/店长角色从这里重新拿取赋值
        globalData.roleMark = guide_role
        // 选择门店时，店员/店长角色的 展示模块权限 从这里重新拿取
        globalData.currentEnterprise.guide_permissions = guide_permissions
        globalData.store = res.data
        globalData.store.store_id = store_id
        globalData.store.qr_code = store_sun_code
      } else {
        wx.showToast(res.message)
      }
      resolve(true)
    } else {
      wx.showToast(res.message)
      reject(false)
    }
  }).catch(err => {
    wx.showToast(err.message)
    reject(false)
  })
  return promise
}
export {
  needSavaStore
}
