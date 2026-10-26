// 获取历史第n个页面栈
const getHistoryPage = (n = 1) => {
  const pages = getCurrentPages() || []
  const historyPage =
    pages.length - n - 1 >= 0 ? pages[pages.length - (n + 1)].route : ''
  return String(historyPage)
}
// 获取scene参数
const getScene = param => {
  param = decodeURIComponent(param)
  if (param.indexOf('=') > -1) {
    return param.split('=')[1]
  }
}
// 弹框
const _showMessage = wx.showModal
const showMessage = option => {
  const {
    title = '',
    content = '',
    isReject = false,
    cancelText = '取消',
    confirmText = '确定',
    confirmColor = '#dd1d21',
    resolve,
    reject,
    complete,
    fail
  } = option
  _showMessage({
    title: title,
    content: content || '网络繁忙，请稍后再试',
    showCancel: isReject,
    cancelText: cancelText,
    confirmText: !content ? '好的' : confirmText,
    confirmColor: confirmColor,
    success (res) {
      if (res.confirm) {
        resolve && resolve()
      } else if (res.cancel) {
        reject && reject()
      }
    },
    fail () {
      fail && fail()
    },
    complete () {
      complete && complete()
    }
  })
}
const _showToast = wx.showToast
// toast
const showToast = (title, is, cb, mask = false) => {
  _showToast({
    title,
    icon: is ? 'success' : 'none',
    duration: 2000,
    image: '',
    mask: mask
  })
  cb &&
  setTimeout(() => {
    cb()
  }, 2000)
}
// 预览图片(开发者注意，企业微信-ios机型，在图片放大完后，不触发onshow生命周期)
const _previewImage = wx.previewImage
const previewImage = (src, more, text = '导购码暂未生成') => {
  const { globalData } = getApp()
  if (!src) {
    showToast(text)
    return false
  }
  globalData.isReviewImg = true
  _previewImage({
    current: src, // 当前显示图片的http链接
    urls: more || [src] // 需要预览的图片http链接列表
  })
}
// 保存图片到本地
const saveImg = url => {
  return new Promise((resolve, reject) => {
    wx.downloadFile({
      url,
      success (res) {
        wx.saveImageToPhotosAlbum({
          filePath: res.tempFilePath,
          success () {
            resolve()
          },
          fail (err) {
            reject(err)
          }
        })
      },
      fail (err) {
        reject(err)
      }
    })
  })
}

const includes = (arr, value) => arr.indexOf(value) !== -1

if (wx) {
  Object.defineProperties(wx, {
    showToast: { writable: true },
    showMessage: { writable: true },
    previewImage: { writable: true }
  })
  wx.getHistoryPage = getHistoryPage
  wx.getScene = getScene
  wx.saveImg = saveImg
  wx.showToast = showToast
  wx.showMessage = showMessage
  wx.previewImage = previewImage
  wx.includes = includes
}
