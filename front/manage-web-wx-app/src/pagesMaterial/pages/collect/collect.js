// Component Object
import { getUploadAuthApi, getCollectListApi, reportTimesApi, getMaterialNewNumApi, checkFirstApi, collectApi } from '@/api/material-api'
import { copyText, saveAllImageToPhotosAlbum, saveAllVideoToPhotosAlbum, saveImageToPhotosAlbum, saveVideoToPhotosAlbum } from '../../utils/util'
import { debounce } from '@/utils/tools'
import dayjs from 'dayjs'

Page({
  data: {
    /** tab列表 */
    tabs: [
      {
        title: '朋友圈',
        value: 7,
        num: 0,
        column: 1,
        gutter: 0,
        list: [],
        page: 1,
        size: 10,
        status: '',
        update: true
      },
      {
        title: '小红书',
        value: 3,
        num: 0,
        column: 1,
        gutter: 0,
        list: [],
        page: 1,
        size: 10,
        status: '',
        update: true
      },
      {
        title: '图片',
        value: 2,
        num: 0,
        column: 2,
        gutter: 30,
        list: [],
        page: 1,
        size: 10,
        status: '',
        update: true
      },
      {
        title: '视频',
        value: 5,
        num: 0,
        column: 2,
        gutter: 30,
        list: [],
        page: 1,
        size: 10,
        status: '',
        update: true
      },
      {
        title: '文本',
        value: 1,
        num: 0,
        column: 1,
        gutter: 0,
        list: [],
        page: 1,
        size: 10,
        status: '',
        update: true
      },
      {
        title: '官网商品',
        value: 8,
        num: 0,
        column: 2,
        gutter: 0,
        list: [],
        page: 1,
        size: 10,
        status: '',
        update: true
      },
      {
        title: '文件',
        value: 6,
        num: 0,
        column: 2,
        gutter: 30,
        list: [],
        page: 1,
        size: 10,
        status: '',
        update: true
      },
      {
        title: '小程序',
        value: 4,
        num: 0,
        column: 2,
        gutter: 30,
        list: [],
        page: 1,
        size: 10,
        status: '',
        update: true
      }
    ],
    /** tab下标 */
    active: 0,
    /** 宫格样式 */
    gridIndex: 0,
    /** 筛选条件 */
    query: {
      /** 关键字 */
      keyword: '',
      /** 分组 */
      source_group_id: 0,
      /** 标签 */
      category_ids: []
    },
    /** 筛选条件 文案 */
    filterInfo: [],
    /** 筛选条件 文案 */
    filterName: '',
    /** 选中列表 */
    selection: [],
    /** 宫格提示 */
    showFirstTips: false,
    /** 预览滚动提示 */
    showFirstPreviewTips: false,
    /** 上传权限 */
    canUpload: false,
    /** 当前操作的数据下标 */
    current: -1,
    /** 一次性下载最大数据 */
    maxlength: 200000,
    /** 预览下标 */
    previewIndex: 0,
    /** 加载状态 */
    loading: true,
    showModal: false
  },

  onLoad(options) {
    const { entryPoint, enterType } = options
    this.setData({
      entryPoint: entryPoint,
      enterType: enterType || 1
    })
    this.getData()
    this.checkFirst()
    this.getUploadAuth()

    // 神策数据上报 -- 离开素材库
    this.eventChannel = this.getOpenerEventChannel()
  },
  onShow() {
    this.getNewNumber()
  },
  onPageScroll({ scrollTop }) {
    this.selectComponent('#com-create-entry').scroll()
  },


  handleBack() {
    // 神策数据上报 -- 点击收藏入口
    wx.track.trackMyFavoriteEntranceClick({
      operation_type: '返回素材库',
      entry_point: '工作台',
      library_type: '品牌素材库'
    })
    wx.navigateBack({
      delta: 1
    })
  },

  /**
   * @description 获取上传素材权限
   */
  async getUploadAuth() {
    try {
      const res = await getUploadAuthApi()
      if (res.code === 200) {
        this.setData({
          canUpload: res.data.auth === 1
        })
      }
    } catch (error) {
      console.log(error)
    }
  },
  /**
   * @description 获取角标数量
   */
  async getNewNumber() {
    const { tabs } = this.data
    const key = 'new_num'
    const value = dayjs().format('YYYY-MM-DD')
    if (wx.getStorageSync(key) === value) return
    try {
      const res = await getMaterialNewNumApi()
      if (res.code === 200) {
        wx.setStorageSync(key, value)
        for (let i = 0; i < res.data.length; i++) {
          const {count_num, type} = res.data[i]
          const index = tabs.findIndex(item => item.value === type)
          this.setData({
            [`tabs[${index}].num`]: count_num
          })
        }
      }
    } catch (error) {
      console.log(error)
    }
  },
  /**
   * @description 是否首次进入
   */
  async checkFirst() {
    try {
      const res = await checkFirstApi()
      if (res.code === 200) {
        this.setData({
          showFirstTips: res.data.first === 0,
          showFirstPreviewTips: res.data.first === 0
        })
      }
    } catch (error) {
      console.log(error)
    }
  },
  /**
     * @description 获取列表
     */
  async getData(page = 1) {
    const { tabs, active, query, selection } = this.data
    const item = tabs[active]

    const params = {
      ...query,
      // size: item.size,
      size: 9999,
      page,
      type: item.value
    }
    this.setData({
      [`tabs[${active}].status`]: 'loading'
    })

    try {
      const res = await getCollectListApi(params)
      if (res.code === 200) {
        const ids = selection.map(child => child.id)
        // 设置选中素材
        const list = params.type === 8 ? res.data.list.map((item) => {
          item.merge_image = (item.brand_source_list || []).map(image => {
            return {
              ...image,
              id: image.resource_key,
              product_name: item.product_name,
              product_code: item.product_code,
              upc_code: item.upc_code
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
          : res.data.list.map(child => {
            child.checked = ids.includes(child.id)
            if (child.type === 2 && child.merge_id) {
              child.merge_image.forEach(image => {
                image.checked = ids.includes(image.id)
                image.is_top = child.is_top
              })
            }
            return child
          })
        item.page = page
        item.update = false
        item.list = (page === 1 ? [] : item.list).concat(list)
        item.status = res.data.total > res.data.page * res.data.size ? 'loadMore' : 'loaded'
        if (item.value === 2) {
          item._total = res.data.list.reduce((total, item) => {
            return total + (item.merge_id && item.merge_image.length > 0 ? item.merge_image.length : 1)
          }, 0)
        }
        item.total = res.data.total || item.list.length
        this.setData({
          [`tabs[${active}]`]: item,
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
    }
  },
  /**
   * @description 加载更多
   */
  loadMore() {
    const { tabs, active } = this.data
    const { status, page } = tabs[active]
    if (status === 'loadMore') {
      this.getData(page + 1)
    }
  },
  /**
   * @description 隐藏首次提示
   */
  handleHideTips() {
    this.setData({
      showFirstTips: false
    })
  },
  /**
   * @description 隐藏首次提示
   */
  handleHidePreviewTips() {
    this.setData({
      showFirstPreviewTips: false
    })
  },
  /**
   * @description 搜索输入
   */
  handleInput(e) {
    this.setData({
      'query.keyword': e.detail.value
    })
  },
  /**
   * @description 确认搜索
   */
  handleSearch: debounce(function() {
    this.setUpdateList()
    this.getData()
  }),
  /**
   * @description 列表更新状态
   */
  setUpdateList() {
    const { tabs } = this.data
    const data = {}
    tabs.forEach((item, index) => {
      data[`tabs[${index}].update`] = true
    })

    this.setData(data)
  },
  /**
   * @description 切换tab
   */
  handleTabChange(e) {
    const { index } = e.detail
    const { tabs } = this.data
    this.setData({
      active: index,
      [`tabs[${index}].num`]: 0
    }, () => {
      if (tabs[index].update) {
        this.getData()
      }
    })
  },
  /**
   * @description 显示筛选弹框
   */
  handleFilter() {
    this.selectComponent('#com-popup-filter').show()
  },
  /**
   * @description 确认筛选
   */
  handleFilterChange(e) {
    const { value } = e.detail
    const res = value.filter(item => item.id)
    const name = res.map(item => item.group_name).join(' > ')
    const id = res[res.length - 1] ? res[res.length - 1].id : 0

    this.setData({
      filterName: name,
      filterInfo: res.map(item => item.id),
      'query.source_group_id': id
    }, () => {
      this.setUpdateList()
      this.getData()
    })
  },
  /**
   * @description 清空筛选条件
   */
  handleFilterClear() {
    this.setData({
      filterName: '',
      filterInfo: [],
      'query.source_group_id': 0
    }, () => {
      this.setUpdateList()
      this.getData()
    })
  },
  /**
   * @description 标签筛选
   */
  handleTagChange(e) {
    const { value } = e.detail
    this.setData({
      'query.category_ids': value.map(item => item.id)
    }, () => {
      this.setUpdateList()
      this.getData()
    })
  },
  /**
   * @description 新建素材
   */
  handleCreate() {
    this.selectComponent('#com-popup-select').show()
  },
  /**
   * @description 切换宫格
   */
  handleGideChange(e) {
    const { value } = e.detail
    this.setData({
      gridIndex: value,
      'tabs[1].column': Number(value) === 0 ? 2 : 3,
      'tabs[1].gutter': Number(value) === 0 ? 30 : 20
    })
  },
  /**
   * @description 选择素材
   * @param {Object} e
   */
  handleSelect(e) {
    const { selection, maxlength } = this.data
    const { item, status, isMergeImage, multiple = false } = {...e.currentTarget.dataset, ...e.detail}
    console.log(item, status, isMergeImage)

    if (isMergeImage) {
      if (status === 2) {
        const ids = item.merge_image.map(val => val.id)
        // 如果是选中状态，取消选中
        this.data.selection = selection.filter(val => !ids.includes(val.id))
      } else {
        // 剩余可选择的数量
        const count = maxlength - selection.length
        if (count === 0) {
          wx.showToast('最多只能选择' + maxlength + '个素材')
          return
        }
        const ids = selection.map(val => val.id)
        // 如果是取消状态，选中
        const list = item.merge_image.filter(val => !ids.includes(val.id)).map(val => {
          return {
            is_top: item.is_top,
            type: item.type,
            ...val
          }
        })
        this.data.selection = [...selection, ...list].splice(0, maxlength)
      }
    } else {
      if (status === 2) {
        const index = selection.findIndex(val => val.id === item.id)
        this.data.selection.splice(index, 1)
      } else {
        if (selection.length === maxlength) {
          wx.showToast('最多只能选择' + maxlength + '个素材')
          return
        }
        this.data.selection.push(item)
      }
    }
    this.setData({
      selection: this.data.selection
    }, () => {
      if (!multiple) {
        this.updateSelectAll()
      }
    })
  },
  updateSelectAll() {
    const { selection, tabs } = this.data
    const data = {}
    const total = selection.reduce((total, item) => {
      if (!total[item.type]) {
        total[item.type] = 1
      } else {
        total[item.type] += 1
      }

      return total
    }, {})
    console.log(total)
    tabs.forEach((item, index) => {
      console.log(item.value, total[item.value] >= item.total && item.total > 0)
      if (item.value === 2) {
        data[`tabs[${index}].isAll`] = total[item.value] === item._total && item._total > 0
      } else {
        data[`tabs[${index}].isAll`] = total[item.value] === item.total && item.total > 0
      }
    })
    this.setData(data)
  },
  /**
   * @description 展示组合素材选择弹框
   */
  handleShowSelectPopup(e) {
    const { index } = e.currentTarget.dataset
    this.setData({
      current: index
    }, () => {
      this.selectComponent('#com-popup-save').show()
    })
  },
  /**
   * @description 组合图片保存单张
   */
  handleSave: debounce(async function(e) {
    const currentPage = getCurrentPages().pop()
    const data = {}
    const { index } = e.detail
    const { tabs, current, active } = this.data
    const item = tabs[active].list[current]
    wx.showLoading({
      title: '保存中'
    })
    try {
      const { original_url, url, id } = item.merge_id && item.merge_image.length > 0 ? item.merge_image[index] : item
      const filePath = original_url || url
      // 朋友圈素材-视频 || 视频素材
      const isVideo = (item.type === 7 && item.original_url) || item.type === 5
      const source_id = item.type === 7 ? item.id : id

      if (isVideo) {
        await saveVideoToPhotosAlbum(filePath)
      } else {
        await saveImageToPhotosAlbum(filePath)
      }

      // 朋友圈素材单个下载不记录数据
      if (item.type !== 7) {
        const reportRes = await reportTimesApi({
          ids: [id],
          type: 1
        })
        if (item.merge_id) {
          data[`tabs[${active}].list[${current}].merge_image[${index}].download_times`] = reportRes.data[item.merge_image[index].id].download_times
        } else {
          data[`tabs[${active}].list[${current}].download_times`] = reportRes.data[item.id].download_times
        }
        // 更新数据
        this.setData(data)
        // 神策数据上报 -- 保存素材
        wx.track.trackAssetSave({
          source_id: source_id,
          task_id: '素材库无任务ID',
          is_official: 1,
          asset_type_modify: item.type,
          official_url: '非官网商品素材',
          official_upc: '非官网商品素材',
          official_product: '非官网商品素材',
          official_name: '非官网商品素材',
          current_page: currentPage.route,
          entry_point: this.data.entryPoint || (this.data.enterType === 2 ? 3 : 1),
          library_type: this.data.enterType === 2 ? 3 : 1,
          is_top: !!item.is_top,
          is_favorite: !!item.is_collect
        })
      }

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
   * @description 保存素材：图片素材、视频素材、朋友圈素材
   */
  handleSaveItem: debounce(async function(e) {
    const currentPage = getCurrentPages().pop()
    const { item } = e.detail
    const { index } = e.currentTarget.dataset
    const { active } = this.data
    const promise = []
    const ids = []
    const text = [] // 提示文案
    wx.showLoading({
      title: '保存中',
      mask: true
    })

    switch (item.type) {
      case 2:
        promise.push(saveImageToPhotosAlbum(item.url))
        ids.push(item.id)
        text.push('图片已保存到相册')
        break
      case 5:
        // 保存视频
        promise.push(saveVideoToPhotosAlbum(item.url))
        ids.push(item.id)
        text.push('视频已保存到相册')
        break
      case 7:
        // 保存图片
        if (item.merge_id && item.merge_image.length > 0) {
          text.push('图片已保存到相册')
          promise.push(saveAllImageToPhotosAlbum(item.merge_image.map(item => item.url)))
        } else if (item.original_url) {
          text.push('视频已保存到相册')
          promise.push(saveVideoToPhotosAlbum(item.original_url))
        }
        ids.push(item.id)
        break

      default:
        break
    }

    try {
      // 保存逻辑
      await Promise.all(promise)
      // 复制文本
      if (item.text_content) {
        await copyText(item.text_content)
        text.push('文本已复制到剪贴板')
      }
      wx.showToast(text.join('，'))
      // 神策数据上报 -- 保存素材
      wx.track.trackAssetSave({
        source_id: item.id,
        task_id: '素材库无任务ID',
        is_official: 1,
        asset_type_modify: item.type,
        official_url: '非官网商品素材',
        official_upc: '非官网商品素材',
        official_product: '非官网商品素材',
        official_name: '非官网商品素材',
        current_page: currentPage.route,
        entry_point: this.data.entryPoint || (this.data.enterType === 2 ? 3 : 1),
        library_type: this.data.enterType === 2 ? 3 : 1,
        is_top: !!item.is_top,
        is_favorite: !!item.is_collect
      })
      const reportRes = await reportTimesApi({
        ids: [item.id],
        type: 1
      })
      const data = {
        [`tabs[${active}].list[${index}].download_times`]: reportRes.data[item.id].download_times
      }
      // 更新数据
      this.setData(data)
    } catch (error) {
      console.log(error)
      if (error.errMsg === 'saveVideoToPhotosAlbum:fail cancel') {
        wx.showToast('您已取消')
      } else if (error.errMsg.indexOf('permission') > -1 || error.errMsg.indexOf('deny') > -1) {
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
   * @description 数据变化时更新数据
   */
  handleItemChange(e) {
    const { index } = e.currentTarget.dataset
    const { active } = this.data
    this.setData({
      [`tabs[${active}].list[${index}]`]: e.detail
    })
  },
  /**
   * @description 更新选中列表
   */
  updateSelectedList(item) {
    const { selection } = this.data
    const is_group = item.merge_id && item.merge_image.length > 1
    const value = is_group ? item.merge_image : [item]

    for (let i = 0; i < value.length; i++) {
      const index = selection.findIndex(item => item.id === value[i].id)
      if (value[i].checked) {
        if (index === -1) {
          selection.push(value[i])
        }
      } else {
        if (index !== -1) {
          selection.splice(index, 1)
        }
      }
    }
    this.setData({
      selection
    })
  },
  /**
   * @description 批量保存
   */
  handleSaveAll: debounce(async function(e) {
    console.log(e.detail)
    const currentPage = getCurrentPages().pop()
    const { tabs } = this.data
    const selection = e.detail.list
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
      // 神策数据上报 -- 保存素材
      selection.map(item => {
        wx.track.trackAssetSave({
          source_id: item.id,
          task_id: '素材库无任务ID',
          is_official: 1,
          asset_type_modify: item.type,
          official_url: '非官网商品素材',
          official_upc: '非官网商品素材',
          official_product: '非官网商品素材',
          official_name: '非官网商品素材',
          current_page: currentPage.route,
          entry_point: this.data.entryPoint || (this.data.enterType === 2 ? 3 : 1),
          library_type: this.data.enterType === 2 ? 3 : 1,
          is_top: !!item.is_top,
          is_favorite: !!item.is_collect
        })
      })
      const reportRes = await reportTimesApi({
        ids: ids,
        type: 1
      })
      // 更新下载数量
      const data = {}

      for (let i = 1; i < 3; i++) {
        for (let j = 0; j < tabs[i].list.length; j++) {
          const item = tabs[i].list[j]
          if (item.merge_id && item.merge_image.length > 1) {
            for (let index = 0; index < item.merge_image.length; index++) {
              const child = item.merge_image[index]
              if (ids.includes(child.id)) {
                data[`tabs[${i}].list[${j}].merge_image[${index}].download_times`] = reportRes.data[item.merge_image[index].id].download_times
              }
            }
          } else {
            if (ids.includes(item.id)) {
              data[`tabs[${i}].list[${j}].download_times`] = reportRes.data[item.id].download_times
            }
          }
        }
      }
      console.log(data)
      this.setData(data, () => {
        // this.handleClearSelected()
      })
      this.setData({
        selection: this.data.selection.filter(item => !ids.includes(item.id))
      })

      this.selectComponent('#com-popup-save').hide()
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
    const { tabs } = this.data
    const data = {}
    // 清空列表选中状态
    for (let i = 0; i < tabs.length; i++) {
      // 只处理图片和视频
      tabs[i].list.forEach((item, index) => {
        if (item.checked) {
          data[`tabs[${i}].list[${index}].checked`] = false
        }
        if (item.type === 2 && item.merge_id) {
          item.merge_image.forEach((child, idx) => {
            if (child.checked) {
              data[`tabs[${i}].list[${index}].merge_image[${idx}].checked`] = false
            }
          })
        }
      })

      data[`tabs[${i}].isAll`] = false
    }

    this.setData({
      ...data,
      selection: []
    })
  },
  /**
   * @description 刷新列表
   */
  handleUpdate(e) {
    const { source_type } = e.detail
    const { tabs, active } = this.data
    const index = tabs.findIndex(item => item.value === source_type)
    if (index === active) {
      this.getData()
    } else {
      this.setData({
        [`tabs[${index}].update`]: true
      })
    }
  },
  /**
   * @description 显示属性筛选弹窗
   */
  handleTag() {
    this.selectComponent('#com-popup-tag').show()
  },

  /**
     * 全选
     */
  handleSelectAll() {
    const { tabs, active, selection } = this.data
    const list = tabs[active].list
    const value = !tabs[active].isAll
    this.setData({
      [`tabs[${active}].isAll`]: value
    })
    console.log(value)
    if (value) {
      // 全选
      list.forEach(item => {
        const status = this.getStatus(selection, item, 3)
        console.log('status', status)

        if (status === 0) {
          this.handleSelect({
            currentTarget: {
              dataset: {}
            },
            detail: {
              item,
              status: 0,
              isMergeImage: item.type === 2 && item.merge_id,
              multiple: true
            }
          })
        }
      })
    } else {
      // 取消全选
      list.forEach(item => {
        this.handleSelect({
          currentTarget: {
            dataset: {}
          },
          detail: {
            item,
            status: 2,
            isMergeImage: item.type === 2 && item.merge_id,
            multiple: true
          }
        })
      })
    }
  },
  getStatus(selection, item, entry) {
    var ids = selection.map(function(val) { return val.id })
    // 图片素材中有多张图片的，商品素材中有多张素材的
    var isMergeImage = entry === 3 ? item.type === 2 && item.merge_image && item.merge_image.length > 0 : (item.type === 2 || item.type === 8 || item.type === 7) && item.merge_image && item.merge_image.length > 0
    if (isMergeImage) {
      var checkedLength = item.merge_image.filter(function(val) { return ids.indexOf(val.id) !== -1 })
      if (checkedLength.length === item.merge_image.length) {
        return 2
      } else if (checkedLength.length > 0) {
        return 1
      } else {
        return 0
      }
    } else {
      if (ids.indexOf(item.id) !== -1) {
        return 2
      } else {
        return 0
      }
    }
  },

  /**
   * 取消收藏所有
   */
  async handleCollectAll() {
    const { selection, tabs } = this.data
    if (selection.length === 0) return
    // 收藏
    const params = {
      is_collect: 2
    }

    params.source_info = Object.values(selection.reduce((prev, cur) => {
      if (!prev[cur.type]) {
        prev[cur.type] = {
          type: cur.type,
          source_id: [cur.id],
          merge_id: [cur.merge_id],
          product_code: [cur.product_code]
        }
      } else {
        prev[cur.type].source_id.push(cur.id)
        prev[cur.type].merge_id.push(cur.merge_id)
        prev[cur.type].product_code.push(cur.product_code)
      }
      return prev
    }, {})).map(item => {
      if (item.type === 8) {
        return {
          type: item.type,
          product_code: item.product_code
        }
      } else if (item.type === 7) {
        return {
          type: item.type,
          merge_id: item.merge_id
        }
      } else {
        return {
          type: item.type,
          source_id: item.source_id
        }
      }
    })

    console.log(params)


    await collectApi(params)

    // 清空选中状态
    wx.showToast('已取消收藏')
    // 重新请求数据
    const data = {
      selection: []
    }

    tabs.forEach((item, index) => {
      data[`tabs[${index}].update`] = true
      data[`tabs[${index}].isAll`] = false
    })
    this.setData(data)
    this.getData()

    // trackAssetFavorite
    const currentPage = getCurrentPages().pop()
    console.log(111, selection)

    selection.map(item => {
      if (item.type === 8) {
        return {
          is_top: !!item.is_top,
          asset_id: item.id,
          is_official: 2,
          asset_type_modify: item.type,
          official_url: item.url,
          official_upc: item.upc_code,
          official_product: item.product_name,
          official_name: item.product_code
        }
      }
      return {
        is_top: !!item.is_top,
        asset_id: item.id,
        is_official: 1,
        asset_type_modify: item.type
      }
    }).forEach(item => {
      wx.track.trackAssetFavorite({
        current_page: currentPage.route,
        entry_point: this.data.entryPoint || (this.data.enterType === 2 ? 3 : 1),
        library_type: this.data.enterType === 2 ? 3 : 1,
        operation_position: 1, // 固定值
        operation_type: 2,
        official_url: '非官网商品素材',
        official_upc: '非官网商品素材',
        official_product: '非官网商品素材',
        official_name: '非官网商品素材',
        ...item
      })
    })
  },

  handleCollect() {
    const { selection } = this.data
    if (selection.length === 0) return
    this.setData({
      showModal: true
    })
  },
  handleCollectConfirm() {
    this.setData({
      showModal: false
    })
    this.selectComponent('#com-popup-save').hide()
    this.handleCollectAll()
  },
  handleCollectCancel() {
    this.setData({
      showModal: false
    })
  },
  async handleCollectOne(e) {
    const { item, index } = e.detail
    const { current, active } = this.data
    this.setData({
      [`tabs[${active}].list[${current}]`]: item
    })
    if (this.eventChannel.emit) {
      this.eventChannel.emit('update', {
        item
      })
    }
    if (item.merge_image.every(v => !v.is_collect)) {
      this.selectComponent('#com-popup-preview').hide()
    }

    await this.getData()
    if (item.merge_image.length) {
      const v = this.data.tabs[active].list[current]
      console.log(index, v.merge_image.length - 1)

      this.setData({
        previewIndex: index === v.merge_image.length ? v.merge_image.length - 1 : index
      })
    }
  },

  /**
   * @description 展示组合素材选择弹框
   */
  handleViewDetail(e) {
    console.log(e)
    wx.navigateTo({
      url: '/pagesMaterial/pages/select/select?enterType=' + 3,
      events: {
        update: (data) => {}
      },
      success: (result) => {
        result.eventChannel.emit('selection', {
          detail: e.detail,
          disabled: true,
          selection: [],
          externalIds: []
        })
      }
    })
  },

  handleRemoveItem(e) {
    console.log(e)
    if (this.eventChannel.emit) {
      this.eventChannel.emit('update', {
        item: e.detail
      })
    }
    this.getData()
  }

})
