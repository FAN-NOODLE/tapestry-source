import COS from 'cos-wx-sdk-v5'
import { getAuthorization } from './request'
import { setFilePath } from './utils'

const { region, bucket } = require('@/lib/config')

const uat = 'https://tapestry-uat-1308608494.cos.ap-shanghai.myqcloud.com/'
const uatcdn = 'https://cdn-cos-uat.coach.com.cn/'
const prod = 'https://tapestry-prod-1308608494.cos.ap-shanghai.myqcloud.com/'
const prodcdn = 'https://cdn-cos.coach.com.cn/'
const publicPrefixes = ['local-upload/', 'static/', 'uploads/']
const privateBucket = bucket === 'tapestry-prod-1308608494' ? 'tapestry-prod-private-1308608494' : 'tapestry-uat-private-1308608494'

function normalizeKey(filePath) {
  return (filePath || '').replace(/^\/+/, '')
}

function isPublicKey(filePath) {
  const key = normalizeKey(filePath)
  return publicPrefixes.some(prefix => key.indexOf(prefix) === 0)
}

function getBucket(filePath) {
  return isPublicKey(filePath) ? bucket : privateBucket
}

function formatFileUrl(location, filePath) {
  const url = `https://${location}`
  if (isPublicKey(filePath)) {
    return url.replace(uat, uatcdn).replace(prod, prodcdn)
  }
  return url
}


const cosConfig = {}

cosConfig.getAuthorization = getAuthorization
// cosConfig.SecretId = 'REDACTED_FOR_SOURCE_SHARING'
// cosConfig.SecretKey = 'REDACTED_FOR_SOURCE_SHARING'

// 初始化
const client = COS ? new COS(cosConfig) : null

/**
 *
 * cos 上传
 * @param {object} config
 * @param {object} config.originPath 需要上传的文件对象
 * @param {string} config.path 指定上传的服务器目录
 * @param {boolean} config.rename 是否重命名
 * @param {boolean} config.onProgress 上传进度回调
 * @returns
 */
export function uploadImage(config) {
  const { originPath } = config
  return new Promise((resolve, reject) => {
    const filePath = setFilePath(config)
    // 上传
    client.postObject(
      {
        Bucket: getBucket(filePath),
        Region: region,
        Key: filePath,
        FilePath: originPath, // 文件路径
        onProgress: function (progressData) {
          config.onProgress && config.onProgress(progressData)
        }
      },
      function (err, data) {
        if (err) {
          reject(err)
        } else {
          console.log(data)
          resolve({
            url: formatFileUrl(data.Location, filePath),
            ext: filePath.substr(filePath.lastIndexOf('.') + 1),
            name: filePath.substr(filePath.lastIndexOf('/') + 1).split('.')[0]
          })
        }
      },
    )
  })
}
