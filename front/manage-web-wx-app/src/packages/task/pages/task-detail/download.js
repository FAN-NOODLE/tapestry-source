/**
 *  获取相册权限
 */
export function wxSaveAuth() {
  return new Promise((resolve, reject) => {
    wx.getSetting({
      success(res) {
        if (!res.authSetting['scope.writePhotosAlbum']) {
          // 如果没有写入权限，则获取写入相册权限
          wx.authorize({
            scope: 'scope.writePhotosAlbum',
            success() {
              resolve()
            },
            fail(err) {
              reject(err)
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

/**
 * 多文件下载并且保存
 * @param {Array} fileList 网络图片数组
 * @param {String} type 文件类型
 */
export function downloadFromUrls(fileList, type, data) {
  wx.showLoading({
    title: '保存中...'
  })
  let saveTask = []
  getTempPath(0, fileList, type, function (res) {
    console.log(res, data.taskInfo)
    const { result, index } = res
    saveTask.push(result)
    if (type === 'image') {
      // 神策数据上报 -- 保存素材
      wx.track.trackAssetSave({
        source_id: fileList[index].media_id || 0,
        task_id: data.taskInfo.task_id,
        is_official: 1,
        asset_type_modify: 2,
        official_url: '非官网商品素材',
        official_upc: '非官网商品素材',
        official_product: '非官网商品素材',
        official_name: '非官网商品素材'
      })
    } else {
      // 神策数据上报 -- 保存素材
      wx.track.trackAssetSave({
        source_id: data.video.origin_media_id,
        task_id: data.taskInfo.task_id,
        is_official: 1,
        asset_type_modify: 5,
        official_url: '非官网商品素材',
        official_upc: '非官网商品素材',
        official_product: '非官网商品素材',
        official_name: '非官网商品素材'
      })
    }
    if (saveTask.length === fileList.length) {
      // 全部保存成功
      wx.hideLoading()
      copyTrick(data)
    }
  })
}
/**
 *
 * @param {String} url 单张网络图片
 */
// 下载内容,临时文件路径
function getTempPath(i, fileList, type, callback) {
  wx.downloadFile({
    url: fileList[i].pic_url,
    success: function (res) {
      var temp = res.tempFilePath
      if (type === 'image') {
        wx.saveImageToPhotosAlbum({
          filePath: temp,
          success: function (res) {
            callback &&
              // eslint-disable-next-line standard/no-callback-literal
              callback({
                result: res,
                index: i
              })
            // 一张保存成功后，递归调用，实现多张图片保存
            if (i < fileList.length - 1) {
              getTempPath(i + 1, fileList, type, callback)
            }
          },
          fail: function (err) {
            console.log(err)
            wx.hideLoading()
          }
        })
      } else {
        wx.saveVideoToPhotosAlbum({
          filePath: temp,
          success: function (res) {
            callback && callback(res)
            console.log(res)
          },
          fail: function (err) {
            console.log(err)
            wx.hideLoading()
          }
        })
      }
    },
    fail: function (err) {
      console.log(err)
    }
  })
}
/**
 * @description: 复制素材文本
 * @param {*} data
 * @return {*}
 */
function copyTrick (data) {
  let { content } = data.text
  if (content) {
    wx.setClipboardData({
      data: content,
      success() {
        wx.showToast('文本已复制到剪贴板,素材已保存到相册')// 神策数据上报 -- 保存素材
        wx.track.trackAssetSave({
          source_id: 0,
          task_id: data.taskInfo.task_id,
          is_official: 1,
          asset_type_modify: 1,
          official_url: '非官网商品素材',
          official_upc: '非官网商品素材',
          official_product: '非官网商品素材',
          official_name: '非官网商品素材'
        })
      }
    })
  } else {
    wx.showToast('文本内容为空')
  }
}
