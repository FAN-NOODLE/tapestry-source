import { uploadMediaApi, getTalkieImageSearchApi } from '@/api/material-api'
import { debounce } from '@/utils/tools'

const now = new Date()
const start = new Date(now)
start.setFullYear(now.getFullYear())
const end = new Date(now)
end.setFullYear(now.getFullYear() + 1)
// Page Object
Page({
  data: {
    empty_image: 'https://cdn-cos.coach.com.cn/uploads/2024/v4.2.0/icon_empty.png',
    title: {
      1: '文本素材',
      2: '图片素材',
      5: '视频素材',
      7: '朋友圈素材'
    },
    showTypeOptions: {
      list: [
        {
          title: '单张展示',
          limit: 20,
          value: 1
        },
        {
          title: '合并展示',
          limit: 9,
          value: 2
        }
      ],
      current: 0
    },
    form: {
      /** 素材的开始时间 */
      start_time: start.setMinutes(start.getMinutes() + 5),
      /** 素材的结束时间 */
      end_time: end.setHours(23, 59, 59, 999),
      /** 素材的标签 */
      category_ids: [],
      /** 素材的类型，1文本，2图片，5视频，7朋友圈 */
      source_type: 1,
      /** 素材的分组 */
      group_id: [],
      /** 素材的资源 */
      sources: [],
      /** 文本内容 */
      content: '',
      /** 展示类型: 1单张, 2合并 */
      show_type: 1,
      /** 文件类型: 1图片, 2视频 */
      file_type: 0,
      /** 视频名称 视频素材有效 */
      videoName: ''
    },
    /** 图片关联商品列表，key为图片id，info为图片信息，list为商品列表 */
    relatedMap: {
      // url: {
      //   status: 0, // 0-AI识别中 1-识别完成
      //   info: {},
      //   list: []
      // }
    },
    /** 图片关联商品列表，key为图片id，info为图片信息，list为商品列表 */
    relatedList: [],
    /** 预览商品图片 */
    previewInfo: {},
    showPreviewPopup: false
  },
  onLoad(options) {
    this.setData({
      options,
      'form.source_type': Number(options.source_type),
      'form.file_type': Number(options.file_type) || 0,
      'form.show_type': Number(options.file_type) === 1 ? 2 : 1
    })
  },
  setFormData(key, value) {
    // console.log('setFormData', key, value)
    this.setData({
      [`form.${key}`]: value
    })
  },
  handleGroupChange(e) {
    const { value } = e.detail
    this.setFormData('group_id', value.map(item => item.id))
  },
  handleTagChange(e) {
    const { value } = e.detail
    this.setFormData('category_ids', value.map(item => item.id))
  },
  handleTimeChange(e) {
    const { value } = e.detail
    this.setData({
      'form.start_time': value[0],
      'form.end_time': value[1]
    })
  },
  handleInputChange(e) {
    let { value } = e.detail
    if (value.length > 500) {
      value = value.slice(0, 500)
    }
    this.setFormData('content', value)
  },
  handleNameChange(e) {
    const { value } = e.detail
    this.setFormData('videoName', value)
    // this.setFormData('sources[0].name', value)
  },
  handleUploadChange(e) {
    const { type } = e.currentTarget.dataset
    const { value, source } = e.detail
    const { relatedMap } = this.data
    console.log('我加载了', source, value)
    const {sources} = this.data.form
    this.setFormData('sources', value.map((item, index) => {
      switch (source) {
        case 'select':
          Object.assign(item, sources[index] || {})
          item.id = item.id || Math.random().toString(36).substring(2)
          item.is_origin = item.is_origin || 0
          item.is_transcode = item.is_transcode || 0
          break
        case 'editName':
          if (relatedMap[item.url]) {
            relatedMap[item.url].info.name = item.name
          }
          break
        case 'useAI':
        case 'selectAllAI':
          item.is_ai = item.is_ai || 0
          // 图片资源，处理关联商品列表
          if (item.is_ai && item.url) {
          // AI图片，添加到关联商品列表中
            relatedMap[item.url] = relatedMap[item.url] || {
              status: 0,
              info: item,
              list: []
            }
            relatedMap[item.url].status = relatedMap[item.url].status === 2 ? 2 : (relatedMap[item.url].list.filter(item => item.add_type === 2).length > 0 ? 1 : 0)
          } else {
            // 非AI图片，且商品列表为空，删除该图片
            if (relatedMap[item.url]) {
              const res = relatedMap[item.url].list.filter(item => item.add_type === 1)
              if (res.length === 0) {
                delete relatedMap[item.url]
              } else {
                relatedMap[item.url].status = 1
                relatedMap[item.url].list = relatedMap[item.url].list.filter(item => item.add_type === 1)
              }
            }
          }
          break
      }
      return item
    }))

    // 启动AI识别关联商品
    if (['71', '2'].includes(type)) {
      // 朋友圈图片素材
      // 清除已删除的图片关联记录
      const urls = value.map(item => item.url).filter(item => item)
      Object.keys(relatedMap).forEach(key => {
        if (!urls.includes(key)) {
          delete relatedMap[key]
        }
      })
      this.setData({
        relatedMap,
        relatedList: Object.values(relatedMap)
      })
      this.startAiRecognition()
    } else {
      if (value.length === 0) {
        this.setData({
          relatedMap: {},
          relatedList: []
        })
      }
    }
  },
  /** 启动AI识别关联商品 */
  async startAiRecognition() {
    const { relatedMap } = this.data
    const keys = Object.keys(relatedMap)
    for (let i = 0; i < keys.length; i++) {
      const key = keys[i]
      if (relatedMap[key].status === 0) {
        // AI图片，启动AI识别
        try {
          const goodsList = await this.aiRecognition(key)
          const upcList = goodsList.map(item => item.upc)

          relatedMap[key].list = [
            ...goodsList,
            ...relatedMap[key].list.filter(item => item.add_type === 1 && !upcList.includes(item.upc))
          ]
          relatedMap[key].status = goodsList.length === 0 ? 2 : 1
        } catch (error) {
          console.error('AI识别失败', error)
          relatedMap[key].status = 2
        }

        this.setData({
          relatedMap,
          relatedList: Object.values(relatedMap)
        })
      }
    }
  },
  /**
   * AI识别关联商品
   * {
        upc: item.upc_code,
        sku: item.sku_item,
        color: item.color,
        url: item.image_url || '',
        name: item.product_name || '',
        add_type: 1, // 1.手动添加 2.AI识别
        score: '' // 置信度 （AI识别时有）
      }
   *  */
  async aiRecognition(key) {
    const res = await getTalkieImageSearchApi({
      image_url: key,
      num: 2
    })
    if (res.code !== 200) {
      return []
    }
    const list = res.data.length ? res.data || [] : []
    const result = []
    // 如果list数据不够5个，需要重复填充5个
    let index = 0
    if (list.length <= 2) {
      list.forEach(item => {
        result.push(...item.picList.map(goods => ({
          upc: goods.upc,
          sku: goods.sku,
          color: goods.color,
          url: goods.url || '',
          name: goods.name || '',
          add_type: 2, // 1.手动添加 2.AI识别
          score: String(goods.score) || '' // 置信度 （AI识别时有）
        })))
      })
    } else {
      while (list.length < 5) {
        const item = JSON.parse(JSON.stringify(list[index]))
        // 去除第一个元素
        item.picList.shift()
        list.push(item)
        index++
      }
      list.forEach((item) => {
        const goods = item.picList[0]
        result.push({
          upc: goods.upc,
          sku: goods.sku,
          color: goods.color,
          url: goods.url || '',
          name: goods.name || '',
          add_type: 2, // 1.手动添加 2.AI识别
          score: String(goods.score) || '' // 置信度 （AI识别时有）
        })
      })
    }

    result.sort((a, b) => Number(b.score) - Number(a.score))

    return result
  },
  /**
   * 手动选择商品
   * @param {*} e
   */
  handleSelectGoodsChange(e) {
    const { value } = e.detail
    const { key } = e.currentTarget.dataset
    const { relatedMap, form } = this.data
    if (key) {
      relatedMap[key].list = value.map(item => ({...item}))
      // relatedMap[key].status = 1
    } else {
      const item = form.sources[0]
      relatedMap[item.url] = {
        status: 1,
        info: item,
        list: value.map(item => ({...item}))
      }
    }
    this.setData({
      relatedMap,
      relatedList: Object.values(relatedMap)
    })
  },
  /**
   * 删除关联商品
   * @param {*} e
   */
  handleDeleteGoods(e) {
    const { key, index } = e.currentTarget.dataset
    const { relatedMap } = this.data
    relatedMap[key].list.splice(index, 1)
    if (relatedMap[key].list.length === 0) {
      delete relatedMap[key]
    }
    this.setData({
      relatedMap,
      relatedList: Object.values(relatedMap)
    })
  },

  handleCheck(e) {
    const { key, index } = e.currentTarget.dataset
    const { relatedMap } = this.data
    relatedMap[key].list[index].checked = !relatedMap[key].list[index].checked
    this.setData({
      relatedMap,
      relatedList: Object.values(relatedMap)
    })
  },

  handleShowTypeChange(e) {
    const { index } = e.currentTarget.dataset
    const { list } = this.data.showTypeOptions
    const { sources } = this.data.form
    const item = list[index]
    if (index === 1 && sources.length > item.limit) {
      wx.showToast(`${item.title}仅支持上传${item.limit}张图片`)
    }
    this.setData({
      'showTypeOptions.current': index
    })
    this.setFormData('show_type', list[index].value)
  },
  handleShowType() {
    this.selectComponent('#com-popup-tips').show()
  },
  handleSubmit: debounce(function() {
    const { showTypeOptions, relatedMap } = this.data
    const form = JSON.parse(JSON.stringify(this.data.form))
    // 数据校验
    if (!form.start_time || !form.end_time) {
      wx.showToast('请选择素材有效期')
      return
    }
    if (form.source_type === 5 && !form.videoName.trim()) {
      wx.showToast('请输入视频名称')
      return
    }
    const isEmpty = {
      1: !form.content,
      2: form.sources.length === 0,
      5: form.sources.length === 0,
      7: !form.content && form.sources.length === 0
    }[form.source_type]

    if (isEmpty) {
      wx.showToast('请添加素材内容')
      return
    }

    // 图片素材需要验证素材数量
    if (form.source_type === 2) {
      const { title, limit } = showTypeOptions.list[showTypeOptions.current]
      if (form.sources.length > limit) {
        wx.showToast(`${title}仅支持上传${limit}张图片`)
        return
      }
      if (form.show_type === 2 && form.sources.length < 2) {
        wx.showToast(`${title}最少上传${2}张图片`)
        return
      }
    }

    const hasUploadFile = form.sources.some(item => !item.url)
    if (hasUploadFile) {
      wx.showToast(`${form.file_type === 2 ? '视频' : '图片'}上传中，请稍后提交`)
      return
    }

    // 数据格式化
    form.group_id = form.group_id[form.group_id.length - 1] || -1
    form.start_time = parseInt(form.start_time / 1000)
    form.end_time = parseInt(form.end_time / 1000)

    // 格式化数据
    const errorImage = []
    form.sources = form.sources.map((item, index) => {
      // 过滤多余字段
      const file = {
        url: item.url,
        is_origin: item.is_origin,
        name: item.name,
        ext: item.ext,
        size: item.size,
        is_transcode: ['mp4', 'jpg', 'png'].includes(item.ext) ? 0 : 1,
        relation_goods: (relatedMap[item.url] ? relatedMap[item.url].list || [] : []).filter(item => item.checked)
      }
      // 朋友圈图片素材+图片组合素材添加排序
      if (form.show_type === 2) {
        file.sort = index
      }
      if (form.source_type === 2) {
        if (file.size > 30 * 1024 * 1024) {
          errorImage.push(index + 1)
        }
      }
      // 视频素材处理视频名称
      if (form.source_type === 5) {
        file.name = form.videoName
        file.cover_url = item.cover_url
      }
      // 朋友圈素材判断原图
      if (form.source_type === 7) {
        // 如果不使用原图，择直接压缩
        if (file.is_origin === 0) {
          file.url = item.compress_url || file.url
          file.cover_url = item.cover_url
        } else if (file.size > 30 * 1024 * 1024) {
          errorImage.push(index + 1)
        }
      }
      return file
    })

    // 朋友圈原图素材
    if (errorImage.length) {
      wx.showToast(`第${errorImage.join()}张图片尺寸超过限制`)
      return
    }

    wx.showLoading({
      title: '上传中',
      mask: true
    })
    this.uploadMethod(form)
    // console.log(form, errorImage); debugger
  }, 300, true),
  handleCancel() {
    wx.navigateBack({
      delta: 1
    })
  },
  handleChangeVideo() {
    const { is_origin } = this.data.form.sources[0]
    this.setData({
      'form.sources[0].is_origin': is_origin === 1 ? 0 : 1
    })
  },
  /**
   * @description: 视频上传方法
   * @return {*}
   */
  uploadMethod(form) {
    uploadMediaApi(form).then(res => {
      // 神策数据上报 -- 上传素材
      (res.data.length ? res.data : []).forEach((item) => {
        wx.track.trackAssetUpload({
          asset_id: item
        })
      })
      wx.showToast('上传成功', false, () => {
        try {
          const eventChannel = this.getOpenerEventChannel()
          eventChannel.emit('update', form)
        } catch (error) {
          console.log(error)
        }
        wx.navigateBack({
          delta: 1
        })
      }, true)
    }).catch(err => {
      console.log(err)
      wx.showToast('上传失败')
    })
  },
  /**
   * @description 显示tooltip
   */
  handleShowTooltip(e) {
    const { selector } = e.currentTarget.dataset
    this.selectComponent('#mp-tooltip').show(this, selector, '')
  },
  /**
   * @description 预览图片
   */
  handlePreview(e) {
    const { item } = e.currentTarget.dataset
    this.setData({
      previewInfo: item,
      showPreviewPopup: true
    })
  },
  /**
   * @description 预览图片
   */
  handlePreviewImage(e) {
    const { item } = e.currentTarget.dataset
    wx.previewImage(item.cover_path || item.cover_url || item.path)
  },
  /**
   * @description 关闭预览弹窗
   */
  handlePreviewClose() {
    this.setData({
      showPreviewPopup: false
    })
  }
})

