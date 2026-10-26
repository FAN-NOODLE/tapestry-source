/** 客户端静默登录 */
import { getWxlogin } from '@/api/scratch-api'

const { enterpriseid } = require('@/lib/config')

/** 获取code */
export const login = (fn) => {
  wx.login({
    success: (res) => {
      getWxLogin(res.code, fn)
    },
    fail (res) {
      console.log('error', res)
    }
  })
}

/** 获取token */
const getWxLogin = async (code, fn) => {
  try {
    const params = {
      enterprise_id: enterpriseid,
      code
    }
    const res = await getWxlogin(params)
    if (res.code === 200) {
      wx.setStorageSync('customer_token', res.data.token)
      fn()
    }
  } catch (error) {}
}
