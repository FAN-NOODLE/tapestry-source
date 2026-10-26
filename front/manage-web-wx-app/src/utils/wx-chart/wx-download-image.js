/* 下载图片 */
const fsm = wx.getFileSystemManager() // 文件管理器
const TEMP_NAME = 'temporary_file_'
const TEMP_DIR = wx.env.USER_DATA_PATH
export function downloadImage (url) {
  if (!url) {
    return Promise.reject(new Error('图片被程序猿拐跑啦，请再试一次'))
  }
  let regexp = /^https.*?/
  let regexpBase64 = /^data:image*?/

  if (regexp.test(url)) {
    return new Promise((resolve, reject) => {
      wx.downloadFile({
        url: url, // 仅为示例，并非真实的资源
        success (res) {
          if (res.statusCode === 200) {
            resolve(res.tempFilePath)
          } else {
            reject(res)
          }
        },
        fail (err) {
          console.error(err)
          reject(err)
        }
      })
    })
  } else if (regexpBase64.test(url)) {
    let FILE_BASE_NAME = `${TEMP_NAME}${Date.now().toString(36)}` // 自定义文件名
    let buffer = wx.base64ToArrayBuffer(url.replace(/data:image\/(png|jpeg|jpg);base64,/, '')) // 转为Buffer数据
    const filePath = `${TEMP_DIR}/${FILE_BASE_NAME}.png` // 文件路径
    return new Promise((resolve, reject) => {
      fsm.writeFile({
        // 把图片从本来拿到
        filePath,
        data: buffer,
        encoding: 'binary',
        success (res) {
          if (res.errMsg === 'writeFile:ok') {
            resolve(filePath)
          }
        },
        fail (err) {
          reject(err)
        }
      })
    })
  } else {
    return url
  }
}
