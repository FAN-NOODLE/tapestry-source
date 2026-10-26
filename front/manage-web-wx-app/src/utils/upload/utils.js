import dayjs from 'dayjs'
/**
 * 获取文件目录
 *
 * @returns 项目名称_端/环境/功能/YYYYMMDD
 */
export const setFolderName = () => {
  const date = dayjs().format('YYYY/MM/DD')
  return ['uploads', date].join('/')
}

/**
 * 重置
 * @param {string} ext 文件后缀
 * @param {boolean} rename 是否重命名
 * @returns
 */
export const resetFileName = (ext) => {
  const chars = 'ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678'
  let fileName = ''
  for (let i = 0; i < 32; i++) {
    fileName += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  return fileName + ext
}

/**
 * 设置服务器文件目录
 * @param {object} param
 * @param {object} param.originPath 需要上传的文件路径
 * @param {string} param.path 指定上传的服务器目录
 * @param {string} param.folder 指定上传的服务器目录
 * @param {boolean} param.rename 是否重命名
 * @returns {string}
 */
export const setFilePath = ({ originPath, path, folder, rename = true }) => {
  const fileName = resetFileName(
    originPath.substr(originPath.lastIndexOf('.'))
  )
  const folderName = path || setFolderName(folder)
  console.log(originPath, originPath.substr(originPath.lastIndexOf('.')))
  return `${folderName}/${fileName}`
}


/**
 * COS 图片压缩 URL 生成
 * @param url 原COS图片地址
 * @param quality 画质 1-100 默认70
 * @param maxWidth 限制最大宽度
 * @param format 输出格式 jpg/png/webp
 */
export function getCosCompressUrl(opt) {
  const { url, quality = 60, maxWidth, maxHeight, format = 'jpg' } = opt
  const base = url.split('?')[0]
  const params = [`imageMogr2/quality/${quality}`, `format/${format}`]

  if (maxWidth && !maxHeight) {
    params.push(`thumbnail/${maxWidth}x`)
  } else if (maxWidth && maxHeight) {
    params.push(`thumbnail/${maxWidth}x${maxHeight}`)
  } else if (maxHeight && !maxWidth) {
    params.push(`thumbnail/x${maxHeight}`)
  }

  return `${base}?${params.join('/')}`
}
