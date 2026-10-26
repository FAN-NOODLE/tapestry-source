
import { uploadImage } from './cos'

var isArray = Array.isArray

function isString(value) {
  return typeof value === 'string'
}
function isObject(value) {
  var type = typeof value
  return value != null && (type === 'object' || type === 'function')
}
/**
 * 上传
 */
export async function upload(files, config = {}) {
  let tempFiles = []

  if (!files) {
    throw new Error('请选择图片')
  }
  // 如果是一个对象
  if (isObject(files)) {
    tempFiles.push(files)
  }
  // 如果是一个字符串
  if (isString(files)) {
    tempFiles.push({ path: files })
  }
  // 如果是一个数组
  if (isArray(files)) {
    tempFiles = files
  }

  let failCount = {}
  let count = 0
  const len = tempFiles.length
  let result = []

  // wx.showLoading({
  //   title: `上传中${count}/${len}`,
  //   mask: false
  // })

  // eslint-disable-next-line no-async-promise-executor
  return new Promise(async (resolve, reject) => {
    while (tempFiles.length >= 1) {
      let item = tempFiles.shift()
      if (item) {
        try {
          // eslint-disable-next-line no-await-in-loop
          const res = await uploadImage({
            originPath: item.path,
            ...config
          })
          result.push({
            url: res.url,
            name: res.name,
            ext: res.ext,
            size: item.size
          })
          // wx.showLoading({
          //   title: `上传中${count++}/${len}`,
          //   mask: false
          // })
        } catch (error) {
          // 如果保存失败，添加到数组中，重新保存
          tempFiles.push(item)
          // 记录每一张图错误的次数
          if (failCount[item]) {
            failCount[item]++
          } else {
            failCount[item] = 1
          }
          // 如果错误次数超过五次终止操作
          if (failCount[item] > 3) {
            // wx.showToast(`上传成功${count},上传失败${len - count}`)
            reject(new Error('上传失败'))
            break
          }
        }
      }
      if (tempFiles.length === 0) {
        // wx.showToast('上传成功')
        resolve(result)
        break
      }
    }
  })
}
