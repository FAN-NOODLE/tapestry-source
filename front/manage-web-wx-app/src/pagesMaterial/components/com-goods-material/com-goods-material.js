import { saveAllImageToPhotosAlbum, saveAllVideoToPhotosAlbum, saveImageToPhotosAlbum, saveVideoToPhotosAlbum } from '../../utils/util'
import { debounce } from '@/utils/tools'
import { getCollectMaterialListApi, getMediaInfoApi, reportTimesApi } from '@/api/material-api'
import { getCurrentPage } from '@/components/common/common/utils'
import { getCosCompressUrl } from '@/utils/upload/utils'
// Component Object
Component({
  options: {
    virtualHost: false,
    multipleSlots: true,
    styleIsolation: 'shared'
  },
  properties: {
    enterType: {
      type: Number,
      value: 1
    },
    entryPoint: {
      type: Number,
      value: 0
    },
    maxlength: {
      type: Number,
      value: 20
    },
    externalIds: {
      type: Array,
      value: []
    }
  },
  data: {
    /** 查询方式切换 */
    tabs: {
      /** 下标  */
      active: 0,
      /** 列表 */
      list: [
        {
          title: '商品信息查询',
          value: 1
        },
        {
          title: '条形码查询',
          value: 'upc_code'
        }
      ]
    },
    /** 搜索类型 */
    types: {
      list: [
        {
          title: '商品款号',
          value: 'sku_item',
          placeholder: '商品款号'
        },
        {
          title: '商品名称',
          value: 'product_name',
          placeholder: '商品名称'
        }
      ],
      active: 0,
      show: false
    },
    /** 搜索值 */
    searchValue: '',
    /** 搜索状态 */
    searchStatus: false,
    /** 宫格 */
    gridIndex: 0,
    /** 列 */
    column: 2,
    /** 间距 */
    gutter: 30,
    /** 列表数据 */
    list: [],
    /** 当前操作的数据下标 */
    current: -1,
    /** 预览下标 */
    previewIndex: 0,
    /** 分页 */
    last_code: '',
    /** 每页条数 */
    size: 15,
    /** 加载状态 */
    status: '',
    /** 是否加载完成 */
    loading: true,
    /** 已选列表 */
    selection: []
  },

  methods: {
    /**
     * @description 获取列表
     */
    async getData() {
      const { last_code, size } = this.data
      const params = {
        size,
        last_code,
        ...this.getSearchParams()
      }
      this.setData({
        status: 'loading'
      })

      try {
        const res = await getCollectMaterialListApi(params)
        if (res.code === 200) {
          const data = res.data.length ? res.data : []
          // 设置选中素材
          const list = data.map((item) => {
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
              _item_type: currentImage ? currentImage.type : '',
              type: 8
            }
          })

          const newList = (!last_code ? [] : this.data.list).concat(list)

          this.setData({
            last_code: newList.length ? newList[newList.length - 1].id : '',
            list: newList,
            status: data.length === size ? 'loadMore' : 'loaded',
            loading: false
          })
        } else {
          this.setData({
            loading: false
          })
        }
      } catch (error) {
        console.log(error)
        this.setData({
          loading: false
        })
      } finally {
        wx.hideLoading()
      }
    },
    /**
     * @description 加载更多
     */
    loadMore() {
      const { status, page } = this.data
      if (status === 'loadMore') {
        this.getData(page + 1)
      }
    },
    getSearchParams() {
      const { tabs, types, searchValue } = this.data
      let key = tabs.list[tabs.active].value
      if (key === 1) {
        key = types.list[types.active].value
      }
      return {
        [key]: searchValue
      }
    },
    /**
     * 切换查询方式
     * @param {*} e
     */
    handleTabChange(e) {
      const { index } = e.detail
      this.setData({
        'tabs.active': index
      })
    },
    /**
     * 搜索类型切换
     */
    handleSearchTypeChange() {
      this.setData({
        'types.show': true
      })
    },
    /** 确定 */
    handleSearchTypeConfirm(e) {
      const { value } = e.detail
      this.setData({
        searchValue: '',
        'types.show': false,
        'types.active': value
      })
    },
    /** 关闭弹框 */
    handleSearchTypeClose() {
      this.setData({
        'types.show': false
      })
    },
    /**
     * 扫码
     */
    handleScan() {
      wx.scanCode({
        success: (res) => {
          this.setData({
            searchValue: res.result
          }, () => {
            this.handleSearch()
          })
        }
      })
    },
    /**
     * 输入
     */
    handleInput(e) {
      this.setData({
        searchValue: e.detail.value
      })
    },
    /**
     * 清空
     */
    handleClear() {
      this.setData({
        searchValue: ''
      })
    },
    /**
     * 搜索
     * @returns
     */
    handleSearch() {
      const { searchValue } = this.data
      if (!searchValue.trim()) return
      this.setData({
        last_code: '',
        searchStatus: true
      })
      wx.showLoading({
        title: '搜索中...',
        mask: true
      })

      this.getData()


      const { tabs, types } = this.data
      let key = tabs.list[tabs.active].value
      if (key === 1) {
        key = types.list[types.active].value
      }
      // 神策数据上报 -- 查询官网商品素材
      wx.track.trackAssetInquiry({
        inquiry_type: key,
        entry_point: this.data.entryPoint || (this.data.enterType === 2 ? 3 : 1)
      })
    },

    /**
     * @description 切换宫格
     */
    handleGideChange(e) {
      const { value } = e.detail
      this.setData({
        gridIndex: value,
        column: [2, 3][value],
        gutter: [30, 20][value]
      })
    },
    /**
     * @description 选择素材
     * @param {Object} e
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
              is_collect: item.is_collect,
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
          this.data.selection.push({...item})
        }
      }
      this.setData({
        selection: this.data.selection
      })
    },
    /**
     * @description 展示组合素材选择弹框
     */
    handleViewDetail(e) {
      const { index } = e.currentTarget.dataset
      const { selection, enterType, entryPoint, externalIds } = this.data
      this.setData({
        current: index
      }, () => {
        wx.navigateTo({
          url: '/pagesMaterial/pages/select/select?enterType=' + enterType + '&entryPoint=' + entryPoint,
          events: {
            update: (data) => {
              this.setData({
                selection: data.selection,
                [`list[${index}]`]: data.item
              })
            }
          },
          success: (result) => {
            result.eventChannel.emit('selection', {
              detail: e.detail,
              selection: selection,
              externalIds: externalIds
            })
          }
        })
      })
    },
    /**
     * @description 数据变化时更新数据
     */
    handleItemChange(e) {
      const { index } = e.currentTarget.dataset
      this.setData({
        [`list[${index}]`]: e.detail
      })
    },
    /**
     * @description 组合图片保存单张
     */
    handleSave: debounce(async function(e) {
      const currentPage = getCurrentPages().pop()
      console.log(e)
      const data = {}
      const { index, item } = e.detail
      const { current } = this.data
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
          entry_point: this.data.entryPoint || (this.data.enterType === 2 ? 3 : 1),
          library_type: this.data.enterType === 2 ? 3 : 1,
          is_top: false,
          is_favorite: !!item.is_collect
        })

        const reportRes = await reportTimesApi({
          resource_keys: [id],
          type: 1
        })
        data[`list[${current}].merge_image[${index}].download_times`] = reportRes.data[id].download_times

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
     * @description 预览
     */
    handlePreview(e) {
      const { previewIndex } = e.detail
      const { index } = e.currentTarget.dataset
      this.setData({
        current: index,
        previewIndex
      }, () => {
        this.selectComponent('#com-popup-preview').show()
      })
    },
    /**
     * @description 弹框关闭时，清除数据
     */
    handleClose() {
      this.setData({
        current: -1
      })
    },

    /**
     * @description 批量保存
     */
    handleSaveAll: debounce(async function() {
      const { selection, list } = this.data
      const currentPage = getCurrentPages().pop()
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
        console.log(this.data.entryPoint)

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
            entry_point: this.data.entryPoint || (this.data.enterType === 2 ? 3 : 1),
            library_type: this.data.enterType === 2 ? 3 : 1,
            is_top: false,
            is_favorite: !!item.is_collect
          })
        })
        const reportRes = await reportTimesApi({
          resource_keys: ids,
          type: 1
        })
        // 更新下载数量
        const data = {}

        for (let j = 0; j < list.length; j++) {
          const item = list[j]
          if (item.merge_image.length > 0) {
            for (let index = 0; index < item.merge_image.length; index++) {
              const child = item.merge_image[index]
              if (ids.includes(child.id)) {
                data[`list[${j}].merge_image[${index}].download_times`] = reportRes.data[item.merge_image[index].id].download_times
              }
            }
          }
        }
        console.log(data)
        this.setData(data, () => {
          this.handleClearSelected()
        })

        // this.selectComponent('#com-popup-save').hide()
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
     * @description 清空选中数据
     */
    handleClearSelected() {
      const { list } = this.data
      const data = {}
      list.forEach((item, index) => {
        if (item.checked) {
          data[`list[${index}].checked`] = false
        }
        if (item.merge_id && item.merge_image.length > 0) {
          item.merge_image.forEach((child, idx) => {
            if (child.checked) {
              data[`list[${index}].merge_image[${idx}].checked`] = false
            }
          })
        }
      })

      this.setData({
        ...data,
        selection: []
      })
    },

    /**
     * @description todo:群发接口
     */
    handleSendAll: debounce(async function() {
      const currentPage = getCurrentPages().pop()
      const { selection, externalIds } = this.data


      // externalIds.forEach(item => {
      //   selection.forEach(child => {
      //     wx.track.trackAssetSend({
      //       asset_id: child.id,
      //       task_id: '群发无法获取',
      //       customer_id: item,
      //       operation_type: 0,
      //       is_official: 2,
      //       asset_type_modify: child.type,
      //       official_url: child.url,
      //       official_upc: child.upc_code,
      //       official_product: child.product_name,
      //       official_name: child.product_code,
      //       current_page: currentPage.route,
      //       entry_point: this.data.enterType === 2 ? 3 : 1,
      //       library_type: this.data.enterType === 2 ? 3 : 1,
      //       is_top: false,
      //       is_favorite: !!child.is_collect
      //     })
      //   })
      // })
      // return

      wx.showLoading({
        title: '加载中...',
        mask: true
      })

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
          // 神策数据上报 -- 发送素材

          externalIds.forEach(item => {
            selection.forEach(child => {
              wx.track.trackAssetSend({
                asset_id: child.id,
                task_id: '群发无法获取',
                customer_id: item,
                operation_type: 0,
                is_official: 2,
                asset_type_modify: child.type,
                official_url: child.url,
                official_upc: child.upc_code,
                official_product: child.product_name,
                official_name: child.product_code,
                current_page: currentPage.route,
                entry_point: this.data.entryPoint || (this.data.enterType === 2 ? 3 : 1),
                library_type: this.data.enterType === 2 ? 3 : 1,
                is_top: false,
                is_favorite: !!child.is_collect
              })
            })
          })
          // todo:清空选中数据
          that.handleClearSelected()
          // 神策数据上报 -- 群发素材
          // wx.track.trackBatchMsgConfirm({
          //   asset_coverage: externalIds.length,
          //   asset_list: selection.map(item => item.id.toString()),
          //   customer_list: externalIds,
          //   operation_type: 0
          // })
          // wx.navigateTo({
          //   url: '/pages/customer-list/customer-list'
          // })
          const page = getCurrentPage()
          wx.navigateBack({
            delta: Number(page.options.back || 3)
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
    }),

    updateCurrent(item) {
      const { list } = this.data

      const index = list.findIndex(child => child.id === item.id)
      if (index > -1) {
        this.setData({
          [`list[${index}]`]: item
        })
      }
    },
    search(e) {
      this.setData({
        searchValue: e
      })
      this.handleSearch()
    }
  }
})
