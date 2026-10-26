/* eslint-disable standard/no-callback-literal */
import { getCosTokenApi } from '@/api/material-api'

/**
 * 获取cos临时密钥
 * @param {object} options
 * @param {Function} callback 回调
 */
export const getAuthorization = async function (options, callback) {
  try {
    const res = await getCosTokenApi()
    console.log(1111, res)
    if (res.data) {
      callback({
        TmpSecretId: res.data.secret_id,
        TmpSecretKey: res.data.secret_key,
        XCosSecurityToken: res.data.session_token,
        StartTime: res.data.start_time, // 时间戳，单位秒，如：1580000000
        ExpiredTime: res.data.expired_time // 时间戳，单位秒，如：1580000900
      })
    } else {
      callback(false, res)
    }
  } catch (error) {
    callback(false, error)
  }
}
