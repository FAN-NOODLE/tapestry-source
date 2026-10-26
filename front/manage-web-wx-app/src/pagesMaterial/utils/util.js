
/**
 * @description 复制文本到粘贴板
 * @param {string} text
 * @returns
 */
export function copyText(text) {
  return new Promise((resolve, reject) => {
    wx.setClipboardData({
      data: text,
      success: (res) => {
        wx.hideToast() // 强制关闭默认提示
        resolve(res)
      },
      fail: reject
    })
  })
}


/**
 * @description 保存图片到相册
 * @param {String} filePath
 * @returns
 */
export function saveImageToPhotosAlbum(filePath) {
  return new Promise((resolve, reject) => {
    wx.downloadFile({
      url: filePath,
      success: (res) => {
        console.log(res)
        wx.saveImageToPhotosAlbum({
          filePath: res.tempFilePath,
          success: resolve,
          fail: reject
        })
      },
      fail: reject
    })
  })
}

/**
 * @description 保存视频到相册
 * @param {String} filePath
 * @returns
 */
export function saveVideoToPhotosAlbum(filePath) {
  return new Promise((resolve, reject) => {
    wx.downloadFile({
      url: filePath,
      success: (res) => {
        wx.saveVideoToPhotosAlbum({
          filePath: res.tempFilePath,
          success: resolve,
          fail: reject
        })
      },
      fail: reject
    })
  })
}

/**
 * @description 保存所有图片到相册
 * @param {Array} list
 * @returns
 */
export function saveAllImageToPhotosAlbum(list) {
  return new Promise((resolve, reject) => {
    const promises = list.map(item => {
      return saveImageToPhotosAlbum(item)
    })
    Promise.all(promises).then(resolve).catch(reject)
  })
}

/**
 * @description 保存所有视频到相册
 * @param {Array} list
 * @returns
 */
export function saveAllVideoToPhotosAlbum(list) {
  return new Promise((resolve, reject) => {
    const promises = list.map(item => {
      return saveVideoToPhotosAlbum(item)
    })
    Promise.all(promises).then(resolve).catch(reject)
  })
}


/**
 * 下载并打开文档
 * @param {string} url 文档链接
 * @param {boolean} open 是否打开 默认打开
 * @returns
 */
export const downloadFile = (url, open = true) => {
  return new Promise((resolve, reject) => {
    const filePath = `${wx.env.USER_DATA_PATH}${decodeURI(
      url.substr(url.lastIndexOf('/')),
    )}`

    const fileType = filePath.substr(filePath.lastIndexOf('.') + 1)

    wx.downloadFile({
      url,
      filePath: filePath,
      success: function (res) {
        if (open) {
          wx.openDocument({
            showMenu: true,
            fileType: fileType,
            filePath: filePath,
            success: resolve,
            fail: reject
          })
        } else {
          resolve(res)
        }
      },
      fail: reject
    })
  })
}


/**
 * 下载并打开文档
 * @param {string} url 文档链接
 * @param {boolean} open 是否打开 默认打开
 * @returns
 */
export const downloadFile2 = (url) => {
  return new Promise((resolve, reject) => {
    wx.downloadFile({
      url,
      success: function (res) {
        resolve(res)
      },
      fail: reject
    })
  })
}

export function getRect(context, selector) {
  return new Promise(function (resolve) {
    wx.createSelectorQuery()
      .in(context)
      .select(selector)
      .boundingClientRect()
      .exec(function (rect) {
        if (rect === void 0) {
          rect = []
        }
        return resolve(rect[0])
      })
  })
}
export function getAllRect(context, selector) {
  return new Promise(function (resolve) {
    wx.createSelectorQuery()
      .in(context)
      .selectAll(selector)
      .boundingClientRect()
      .exec(function (rect) {
        if (rect === void 0) {
          rect = []
        }
        return resolve(rect[0])
      })
  })
}
/**
 *  获取相册权限
 */
export function wxSaveAuth() {
  return new Promise((resolve, reject) => {
    wx.getSetting({
      success(res) {
        resolve(res)
        if (!res.authSetting['scope.writePhotosAlbum']) {
          // 如果没有写入权限，则获取写入相册权限
          wx.authorize({
            scope: 'scope.writePhotosAlbum',
            success(writeRes) {
              resolve(writeRes)
            },
            fail(err) {
              reject(err)
              wx.hideLoading()
              // 用户拒绝授权
              wx.showModal({
                content: '检测到您没打开相册权限，是否去设置打开？',
                confirmText: '确认',
                cancelText: '取消',
                success(res) {
                  if (res.confirm) {
                    wx.openSetting({
                      success: (res) => {}
                    })
                  }
                }
              })
            }
          })
        } else {
          resolve()
        }
      },
      fail(e) {
        reject(e)
      }
    })
  })
}
