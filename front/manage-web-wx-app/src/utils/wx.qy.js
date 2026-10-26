
const systemInfo = wx.getSystemInfoSync()

function isWxWork() {
  const wxwork = systemInfo.environment === 'wxwork'
  if (!wxwork) {
    wx.showToast(
      '联系功能仅支持企业微信登录时使用，请从企业微信打开WeClient Pro后再试'
    )
    return true
  }
  return false
}

function handleApiCall(apiFunction, options) {
  if (isWxWork()) {
    return Promise.reject(new Error(`${apiFunction.name} is not supported in wxwork`))
  }

  return new Promise((resolve, reject) => {
    apiFunction({
      ...options,
      success: (res) => {
        options.success && options.success(res)
        resolve(res)
      },
      fail: (res) => {
        options.fail && options.fail(res)
        reject(res)
      },
      complete: (res) => {
        options.complete && options.complete(res)
      }
    })
  })
}

/**
 * @description 打开个人信息页
 * @link https://developer.work.weixin.qq.com/document/path/93567
 * @param {*} options
 * @returns
 */
export function openUserProfile(options = {}) {
  return handleApiCall(wx.qy.openUserProfile, options)
}

/**
 * @description 打开个人信息页
 * @link https://developer.work.weixin.qq.com/document/path/91519
 * @param {*} options
 * @returns
 */
export function openEnterpriseChat(options = {}) {
  return handleApiCall(wx.qy.openEnterpriseChat, options)
}
/**
 * @description 分享到外部联系人
 * @link https://developer.work.weixin.qq.com/document/path/93571
 * @param {*} options
 * @returns
 */
export function shareToExternalContact(options = {}) {
  return handleApiCall(wx.qy.shareToExternalContact, options)
}
/**
 * @description 该接口支持在小程序环境中，具有客户联系权限的企业成员将文本内容和附件传递到客户群群发、发送到客户群
 * @link https://developer.work.weixin.qq.com/document/path/93572
 * @param {*} options
 * @returns
 */
export function shareToExternalChat(options = {}) {
  return handleApiCall(wx.qy.shareToExternalChat, options)
}
/**
 * @description 具有客户联系权限的企业成员，可通过该接口将文本内容和附件传递到客户朋友圈。当前暂仅支持在企业微信内调用。
 * @link https://developer.work.weixin.qq.com/document/path/94608
 * @param {*} options
 * @returns
 */
export function shareToExternalMoments(options = {}) {
  return handleApiCall(wx.qy.shareToExternalMoments, options)
}

