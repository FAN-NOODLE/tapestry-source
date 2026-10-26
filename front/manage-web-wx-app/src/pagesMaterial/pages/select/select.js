// Page Object
import { saveAllImageToPhotosAlbum, saveAllVideoToPhotosAlbum, saveImageToPhotosAlbum, saveVideoToPhotosAlbum, wxSaveAuth } from '../../utils/util'
import { debounce } from '@/utils/tools'
import { getCollectMaterialListApi, getMediaInfoApi, reportTimesApi } from '@/api/material-api'
import { getCosCompressUrl } from '@/utils/upload/utils'

Page({
  data: {
    video_cover_url: 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/video_cover.png',
    list: [],
    isHalf: false,
    isSelectAll: false,
    selection: [],
    maxlength: 20,
    showExpand: false,
    isExpand: true,
    current: -1,
    previewIndex: 0,
    externalIds: [],
    isLoading: true
  },
  // options(Object)
  onLoad: function(options) {
    if (options.from === 'TALK') {
      options.entryPoint = 5
      options.enterType = 1
      this.setData({
        libraryType: 5
      })
    }

    const enterType = options.enterType ? Number(options.enterType) : 1
    const entryPoint = options.entryPoint ? Number(options.entryPoint) : 0
    this.setData({
      options,
      enterType: enterType,
      maxlength: enterType === 1 ? 20 : 9,
      entryPoint: entryPoint
    })
    if (options.upc_code) {
      this.getDetail()
    } else {
      const eventChannel = this.getOpenerEventChannel()
      eventChannel.on('selection', data => {
        this.setData({
          isLoading: false,
          disabled: data.disabled,
          detail: data.detail,
          selection: data.selection,
          externalIds: data.externalIds
        })
      })
    }
  },
  onUnload() {
    const eventChannel = this.getOpenerEventChannel()
    if (eventChannel.emit && this.data.detail) {
      eventChannel.emit('update', {
        selection: this.data.selection,
        item: this.data.detail
      })
    }
  },
  /**
   * 获取详情
   */
  async getDetail() {
    const { options } = this.data
    const res = await getCollectMaterialListApi({
      last_code: '',
      upc_code: options.upc_code,
      store_code: options.store_code
    })
    if (res.code === 200 && res.data.length) {
      const detail = res.data.map((item) => {
        item.merge_image = (item.brand_source_list || []).map(image => {
          return {
            ...image,
            id: image.resource_key,
            product_name: item.product_name,
            product_code: item.product_code,
            upc_code: item.upc_code,
            is_collect: item.is_collect
          }
        })
        const currentImage = item.merge_image.find(image => image.is_main === 1) || item.merge_image[0]
        return {
          ...item,
          url: currentImage ? currentImage.url : '',
          name: currentImage ? currentImage.name : item.product_code,
          type: 8
        }
      })[0]
      detail.sku_tag_price_text = Number(detail.sku_tag_price).toLocaleString('zh-CN', {
        style: 'currency',
        currency: 'CNY'
      })
      detail.store_price_text = Number(detail.store_price).toLocaleString('zh-CN', {
        style: 'currency',
        currency: 'CNY'
      })
      this.setData({
        detail,
        isLoading: false
      })
      this.setShowExpand()
    } else {
      // wx.showToast('未查询到该商品')
      this.setData({
        isLoading: false
      })
    }
  },
  /**
   * 展开收起按钮
   */
  setShowExpand() {
    setTimeout(() => {
      wx.nextTick(() => {
        wx.createSelectorQuery()
          .in(this)
          .select('#desc')
          .boundingClientRect()
          .exec(rect => {
            if (rect[0]) {
              this.setData({
                showExpand: rect[0].height > 200
              })
            }
          })
      })
    }, 100)
  },
  /**
   * 展开
   */
  handleExpand() {
    this.setData({
      isExpand: !this.data.isExpand
    })
  },
  /**
   * 库存查询
   */
  handleSearch() {
    const params = {
      searchValue: this.data.detail.upc_code,
      active: 1
    }
    // 转成字符串
    const paramsStr = encodeURIComponent(JSON.stringify(params))
    wx.navigateTo({
      url: `/pagesStock/pages/list/list?from=2&searchParams=${paramsStr}`
    })

    // 神策数据上报 -- 库存查询
    wx.track.trackStockInquiry({
      entry_point: this.data.entryPoint === 5 ? 5 : this.data.enterType,
      inquiry_type: 2
    })
  },
  /**
   * @description 选择
   */
  handleSelect(e) {
    const { selection, maxlength } = this.data
    const { item, status, isMergeImage } = {...e.currentTarget.dataset, ...e.detail}

    if (isMergeImage) {
      // 剩余可选择的数量
      const count = maxlength - selection.length
      if (status !== 2 && count === 0) {
        wx.showToast('最多只能选择' + maxlength + '个素材')
        return
      }

      if (status === 2) {
        const ids = item.merge_image.map(val => val.id)
        // 如果是选中状态，取消选中
        this.data.selection = selection.filter(val => !ids.includes(val.id))
      } else {
        const ids = selection.map(val => val.id)
        // 如果是取消状态，选中
        const list = item.merge_image.filter(val => !ids.includes(val.id)).map(val => {
          return {
            type: item.type,
            ...val
          }
        })
        this.data.selection = [...selection, ...list].splice(0, maxlength)
      }
    } else {
      if (status === 0 && selection.length === maxlength) {
        wx.showToast('最多只能选择' + maxlength + '个素材')
        return
      }
      if (status === 2) {
        const index = selection.findIndex(val => val.id === item.id)
        this.data.selection.splice(index, 1)
      } else {
        this.data.selection.push(item)
      }
    }
    this.setData({
      selection: this.data.selection
    })
  },
  /**
   * @description: 预览图片
   */
  previewImage(e) {
    // const currentPage = getCurrentPages().pop()
    const { index } = e.currentTarget.dataset
    this.setData({
      previewIndex: index
    })
    this.selectComponent('#com-popup-save-preview').show()
  },
  /**
   * @description: 预览保存
   * @return {*}
   */
  handleSave: debounce(async function(e) {
    const pages = getCurrentPages()
    const currentPage = pages[pages.length - 2] || pages[0]
    const { options, enterType, entryPoint, libraryType } = this.data
    console.log(e)
    const data = {}
    const { index, item } = e.detail
    wx.showLoading({
      title: '保存中'
    })
    try {
      const { type, url, id } = item.merge_image[index]
      // 朋友圈素材-视频 || 视频素材
      const isVideo = type === 5
      isVideo ? await saveVideoToPhotosAlbum(url) : await saveImageToPhotosAlbum(url)
      // 神策数据上报 -- 保存素材
      wx.track.trackAssetSave({
        source_id: id,
        task_id: '素材库无任务ID',
        is_official: 2,
        asset_type_modify: type,
        official_url: url,
        official_upc: item.upc_code,
        official_product: item.product_name,
        official_name: item.product_code,
        current_page: currentPage.route,
        entry_point: entryPoint || {
          1: 1,
          2: 3,
          3: 1,
          4: 4
        }[Number(options.from) === 2 ? 4 : enterType],
        library_type: libraryType || {
          1: 1,
          2: 3,
          3: 1,
          4: 4
        }[Number(options.from) === 2 ? 4 : enterType],
        is_top: false,
        is_favorite: !!this.data.detail.is_collect
      })
      const reportRes = await reportTimesApi({
        resource_keys: [id],
        type: 1
      })
      data[`detail.merge_image[${index}].download_times`] = reportRes.data[id].download_times

      // 更新数据
      this.setData(data)

      wx.showToast('保存成功')
    } catch (error) {
      console.log(error)
      if (error.errMsg.indexOf('permission') > -1 || error.errMsg.indexOf('deny') > -1) {
        wx.hideLoading()
        wx.showModal({
          title: '保存图片需要您的相册权限',
          content: '打开企微相册权限，可保存到相册哦~',
          showCancel: false,
          cancelText: '好',
          cancelColor: '#000000',
          confirmText: '好',
          confirmColor: '#000000'
        })
      } else {
        wx.showToast('保存失败')
      }
    }
  }),
  /**
   * @description 批量保存
   */
  handleSaveAll: debounce(async function() {
    const pages = getCurrentPages()
    const currentPage = pages[pages.length - 2] || pages[0]
    // await wxSaveAuth()
    const { selection, detail, options, enterType, entryPoint, libraryType } = this.data
    if (selection.length === 0) return
    const videoList = selection.filter(item => item.type === 5).map(item => item.url)
    const imageList = selection.filter(item => item.type !== 5).map(item => item.url)
    const ids = selection.map(item => item.id)
    wx.showLoading({
      title: '保存中',
      mask: true
    })

    try {
      if (videoList.length) {
        await saveAllVideoToPhotosAlbum(videoList)
      }
      if (imageList.length) {
        await saveAllImageToPhotosAlbum(imageList)
      }

      selection.map(item => {
        // 神策数据上报 -- 保存素材
        wx.track.trackAssetSave({
          source_id: item.id,
          task_id: '素材库无任务ID',
          is_official: 2,
          asset_type_modify: item.type,
          official_url: item.url,
          official_upc: item.upc_code,
          official_product: item.product_name,
          official_name: item.product_code,
          current_page: currentPage.route,
          entry_point: entryPoint || {
            1: 1,
            2: 3,
            3: 1,
            4: 4
          }[Number(options.from) === 2 ? 4 : enterType],
          library_type: libraryType || {
            1: 1,
            2: 3,
            3: 1,
            4: 4
          }[Number(options.from) === 2 ? 4 : enterType],
          is_top: false,
          is_favorite: !!detail.is_collect
        })
      })
      const res = await reportTimesApi({
        resource_keys: ids,
        type: 1
      })
      const data = {
        selection: []
      }
      ids.forEach((id) => {
        const index = detail.merge_image.findIndex(item => item.id + '' === id + '')
        data[`detail.merge_image[${index}].download_times`] = res.data[id].download_times
      })
      this.setData(data)
      wx.showToast('保存成功')
    } catch (error) {
      console.log(error)
      if (error.errMsg.indexOf('permission') > -1 || error.errMsg.indexOf('deny') > -1) {
        wx.hideLoading()
        wx.showModal({
          title: '保存图片需要您的相册权限',
          content: '打开企微相册权限，可保存到相册哦~',
          showCancel: false,
          cancelText: '好',
          cancelColor: '#000000',
          confirmText: '好',
          confirmColor: '#000000'
        })
      } else {
        wx.showToast('保存失败')
      }
    }
  }),
  /**
   * 清空已选
   */
  handleClearSelected() {
    this.setData({
      selection: []
    })
  },

  /**
   * @description todo:群发接口
   */
  handleSendAll: debounce(async function() {
    const pages = getCurrentPages()
    const currentPage = pages[pages.length - 2] || pages[0]
    wx.showLoading({
      title: '加载中...',
      mask: true
    })
    const { selection, externalIds, options, enterType, entryPoint, libraryType } = this.data

    const that = this
    if (selection.length === 0) return
    // todo:筛选出所有视频素材
    const videoList = selection.filter(item => item.type === 5)
    const imageList = selection.filter(item => item.type === 2)
    const types = {}
    const params = {
      externalUserIds: externalIds,
      text: {
        content: ''
      },
      attachments: [
        ...imageList.map(item => {
          types[item.id] = 2
          return {
            msgtype: 'image',
            image: {
              imgUrl: getCosCompressUrl({
                url: item.url,
                maxWidth: 1080,
                maxHeight: 1440
              })
            }
          }
        })
      ]
    }
    // todo:获取视频meida_id
    if (videoList.length > 0) {
      const res = await Promise.all(videoList.map(item => {
        let params = {
          url: item.url,
          type: 0,
          source_type: 5
        }
        types[item.id] = 5
        return getMediaInfoApi(params)
      }))

      params.attachments = params.attachments.concat(...res.map(item => {
        return {
          msgtype: 'video',
          video: {
            mediaid: item.data.media_id
          }
        }
      }))
    }
    console.log(params.attachments)
    wx.qy.shareToExternalContact({
      ...params,
      success: res => {
        console.log(res)
        // todo:清空选中数据
        that.handleClearSelected()
        // 神策数据上报 -- 群发素材
        wx.track.trackBatchMsgConfirm({
          asset_coverage: externalIds.length,
          asset_list: selection.map(item => item.id.toString()),
          customer_list: externalIds,
          operation_type: 0
        })
        // 神策数据上报 -- 发送素材
        externalIds.forEach(item => {
          selection.forEach(child => {
            wx.track.trackAssetSend({
              asset_id: child.id,
              task_id: '群发无法获取',
              customer_id: item,
              operation_type: 0,
              is_official: 2,
              asset_type_modify: types[child.id],
              official_url: child.url,
              official_upc: child.upc_code,
              official_product: child.product_name,
              official_name: child.product_code,
              current_page: currentPage.route,
              entry_point: entryPoint || {
                1: 1,
                2: 3,
                3: 1,
                4: 4
              }[Number(options.from) === 2 ? 4 : enterType],
              library_type: libraryType || {
                1: 1,
                2: 3,
                3: 1,
                4: 4
              }[Number(options.from) === 2 ? 4 : enterType],
              is_top: !!child.is_top,
              is_favorite: !!this.data.detail.is_collect
            })
          })
        })
        // wx.navigateTo({
        //   url: '/pages/customer-list/customer-list'
        // })
        wx.navigateBack({
          delta: 2
        })
      },
      fail: function (res) {
        console.log(res)
      },
      complete: function (res) {
        console.log('complete', res)
        wx.hideLoading()
      }
    })

    this.selectComponent('#com-popup-save').hide()
  })
})
