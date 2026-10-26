const {cdphost} = require('./config')

module.exports = {
  appKey: 'REDACTED_FOR_SOURCE_SHARING', // 应用唯一标识
  secretKey: 'REDACTED_FOR_SOURCE_SHARING', // API接口调用密钥
  appName: 'WeClient Pro', // 小程序名称
  appId: 'wxe4191c725f430b22', // 小程序appid
  useComponetSystem: true,
  server: cdphost || 'https://cdptest.vchangyi.com',
  showLog: false // 打印上报消息调试信息
}
