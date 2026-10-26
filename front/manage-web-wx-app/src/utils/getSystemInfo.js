/**
 * 获取系统信息
 * @returns {Wxapp.getSystemInfoSync.Result}
 */
export const getSystemInfoSync = () => {
  return wx.getSystemInfoSync()
}
export const systemInfo = getSystemInfoSync()

/**
 * 获取安全区
 * @returns {Wxapp.General.SafeAreaResult}
 */
export const getSafeArea = () => {
  return systemInfo.safeArea
}

/**
 * 企业微信环境
 */
export const wxwork = systemInfo.environment === 'wxwork'

/**
 * 获取导航栏信息
 * @returns
 */
export function getNavInfo() {
  const { right, width } = wx.getMenuButtonBoundingClientRect()
  const isAndroid = systemInfo.platform === 'android'
  // 胶囊宽度
  const menuWidth = (systemInfo.windowWidth - right) * 2 + width
  // 状态栏的高度
  const navHeight = isAndroid ? 48 : 44
  // 导航栏的高度
  const totalHeight = navHeight + systemInfo.statusBarHeight
  // 中间位置宽度
  const middleWidth = systemInfo.windowWidth - 2 * menuWidth
  return {
    menuWidth, // 胶囊的宽度占位
    height: navHeight, // 状态栏的高度
    width: middleWidth, // 中间内容的宽度
    statusBarHeight: systemInfo.statusBarHeight, // 导航栏的高度
    totalHeight
  }
}
