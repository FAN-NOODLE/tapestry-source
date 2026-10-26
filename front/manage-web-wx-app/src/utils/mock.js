// 重写request结合mock.js
import Mock from 'mockjs'

const __request = wx.request

Object.defineProperty(wx, 'request', {writable: true})
wx.request = function (config) {
  if (typeof Mock._mocked[config.url] === 'undefined') {
    __request(config)
    return
  }
  const resTemplate = Mock._mocked[config.url].template
  const response = Mock.mock(resTemplate)
  if (typeof config.success === 'function') {
    config.success(response)
  }
  if (typeof config.complete === 'function') {
    config.complete(response)
  }
}

module.exports = Mock

