import { reportTimesApi, getMediaInfoApi, collectApi } from '@/api/material-api'
import { copyText, downloadFile, saveAllImageToPhotosAlbum, saveVideoToPhotosAlbum } from '../../utils/util'
import dayjs from 'dayjs'
import { debounce } from '@/utils/tools'
import { shareToExternalMoments } from '@/utils/wx.qy'
import { getCosCompressUrl } from '@/utils/upload/utils'

// Component Object
Component({
  externalClasses: ['custom-class'],
  options: {
    virtualHost: false,
    multipleSlots: true,
    styleIsolation: 'apply-shared'
  },
  behaviors: [],
  properties: {
    customStyle: {
      type: String
    },
    item: {
      type: Object,
      observer: 'init'
    },
    selection: {
      type: Array,
      value: []
    },
    /** 群发还是保存  1保存 2群发 */
    enterType: {
      type: Number,
      value: 1
    },
    entryPoint: {
      type: Number,
      value: 0
    },
    showCount: {
      type: Boolean,
      value: true
    }
  },
  data: {
    default_avatar: 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_logo.png',
    empty_image: 'https://cdn-cos.coach.com.cn/uploads/2024/v4.2.0/icon_empty.png',
    video_cover_url: 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/video_cover.png',
    // 合并素材数量
    merge_count: 0,
    forwarding_times: 0,
    download_times: 0,
    is_group: false,
    // 文件类型图标
    fileType: [
      '',
      'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/file/1.png',
      'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/file/2.png',
      'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/file/3.png',
      'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/file/4.png',
      'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/file/5.png'
    ]
  },
  methods: {
    /**
     * @description 初始化：图片素材、视频素材
     */
    init() {
      const { item } = this.data
      switch (item.type) {
        case 2:
        case 8:
          this.computeImageInfo(item)
          break
        case 5:
          this.computeVideoInfo(item)
          break
        case 7:
          this.computeInfo(item)
          break

        default:
          break
      }
    },
    /**
     * @description 计算图片素材信息
     */
    computeImageInfo(item) {
      const is_group = (item.merge_id && item.merge_image.length > 1) || (item.type === 8 && item.merge_image.length > 0)
      if (is_group) {
        const merge_count = item.merge_image.length
        const [forwarding_times, download_times] = item.merge_image.reduce((prev, current) => {
          prev[0] += current.forwarding_times
          prev[1] += current.download_times
          return prev
        }, [0, 0])
        this.setData({
          merge_count,
          forwarding_times,
          download_times,
          is_group
        })
      } else {
        this.setData({
          forwarding_times: item.forwarding_times || 0,
          download_times: item.download_times || 0,
          is_group,
          is_half_checked: false
        })
      }
    },
    computeVideoInfo(item) {
      const { cover_url, url } = item
      const video_cover_url = /^http/.test(cover_url) && cover_url !== url ? cover_url : 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/video_cover.png'
      this.setData({
        video_cover_url
      })
    },
    /**
     * @description 计算朋友圈素材
     */
    computeInfo(item) {
      let btnText = []
      if (item.merge_image.length > 0) {
        btnText.push('保存图片')
      } else if (item.original_url) {
        btnText.push('保存视频')
      }

      if (item.text_content) {
        btnText.push('复制文本')
      }
      // 获取视频封面图
      const video_cover_url = item.video_cover_url !== item.original_url && /^http/.test(item.video_cover_url) ? item.video_cover_url : 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/video_cover.png'
      this.setData({
        video_cover_url,
        btnText: btnText.join(' + '),
        start_time: dayjs(item.start_time * 1000).format('YYYY/MM/DD HH:mm:ss')
      })
    },
    /**
     * @description 预览素材: 图片素材、视频素材、朋友圈素材
     */
    handlePreview: debounce(function(e) {
      const currentPage = getCurrentPages().pop()
      const { index } = e.currentTarget.dataset
      const { item, enterType } = this.data
      switch (item.type) {
        case 2:
          const isMergeImage = item.merge_image.length > 1
          if (enterType === 2 && isMergeImage) {
            this.handleSavePopup()
            return
          }
          this.triggerEvent('preview', {previewIndex: index, item})

          break
        case 6:
          // 文件预览
          wx.showLoading({
            title: '打开中'
          })
          downloadFile(item.url, true).then(() => {
            wx.hideLoading()
          }).catch(() => {
            wx.showToast('企微暂不支持该格式预览')
          })
          // 神策数据上报 -- 查看素材
          wx.track.trackAssetView({
            source_id: item.id,
            task_id: '素材库无任务ID',
            is_official: 1,
            asset_type_modify: 6,
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
          break
        case 8:
          if (enterType === 2) {
            this.handleViewDetail()
          } else {
            if (item.merge_image.length === 0) return
            this.triggerEvent('preview', {previewIndex: index, item})
          }
          break

        case 7:
          this.triggerEvent('preview', {previewIndex: index, item})
          break
        case 5:
          this.triggerEvent('preview', {previewIndex: index, item})
          break
      }
    }),
    /**
     * @description 保存素材：图片素材、视频素材、朋友圈素材
     */
    handleSave() {
      const { item } = this.data
      this.triggerEvent('save', {
        item,
        index: -1
      })
    },
    /**
     * @description 分组图片保存选择
     */
    handleSavePopup: debounce(function() {
      const { item } = this.data
      this.triggerEvent('select-save', item)
    }),
    /**
     * @description 选中单张素材：图片素材、视频素材
     */
    handleSelect(e) {
      const { item } = this.data
      const { status, isMergeImage } = e.currentTarget.dataset

      this.triggerEvent('select', {
        item,
        status,
        isMergeImage
      })
    },
    /**
     * @description 复制素材：文本素材、朋友圈素材的文本内容
     */
    handleCopy: debounce(async function(e) {
      const currentPage = getCurrentPages().pop()
      const { text } = e.currentTarget.dataset
      const { item } = this.data
      try {
        await copyText(text)
        wx.showToast('文本已复制到剪贴板')
        // 文本素材的复制需要上报
        if (item.type === 1) {
          const reportRes = await reportTimesApi({
            ids: [item.id],
            type: 3
          })
          if (reportRes.code === 200) {
            item.copy_times = reportRes.data[item.id].copy_times
          } else {
            item.copy_times++
          }
          this.triggerEvent('change', item)
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
        }
      } catch (error) {
        console.log(error)
      }
    }),
    /**
     * @description 分享素材：朋友圈素材
     */
    handleShare: debounce(async function() {
      const currentPage = getCurrentPages().pop()
      wx.showLoading({
        title: '转发中'
      })
      const { item } = this.data
      const params = {}
      // 文本内容
      if (item.text_content) {
        params.text = {
          content: item.text_content
        }
      }

      if (item.merge_image.length > 0) {
        // 图片附件
        params.attachments = item.merge_image.map(item => {
          return {
            msgtype: 'image', // 消息类型，必填
            image: {
              imgUrl: getCosCompressUrl({
                url: item.url,
                maxWidth: 1080,
                maxHeight: 1440
              }) // 图片的imgUrl,跟图片mediaid填其中一个即可
            }
          }
        })
      } else if (item.original_url) {
        // 获取视频的media ID
        const mediaInfo = await getMediaInfoApi({
          url: item.original_url,
          source_type: 5,
          type: 1
        })
        if (mediaInfo.code === 200) {
          params.attachments = [
            {
              msgtype: 'video', // 消息类型，必填
              video: {
                mediaid: mediaInfo.data.media_id // 视频的素材id
              }
            }
          ]
        }
      }

      console.log(params)
      try {
        const res = await shareToExternalMoments(params)
        wx.hideLoading()
        if (res.errMsg === 'qy__shareToExternalMoments:ok') {
          // 神策数据上报 -- 发送素材 分享朋友圈
          wx.track.trackAssetSend({
            asset_id: item.id,
            task_id: '转发到朋友圈无法获取',
            customer_id: '转发到朋友圈无法获取',
            operation_type: 1,
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
            type: 2
          })
          if (reportRes.code === 200) {
            item.forwarding_times = reportRes.data[item.id].forwarding_times
          } else {
            item.forwarding_times++
          }
          this.triggerEvent('change', item)
        }
      } catch (error) {
        console.log(error)
        wx.hideLoading()
      }
    }, 300, true),
    /**
     * 查看详情
     */
    handleViewDetail() {
      this.triggerEvent('viewDetail', this.data.item)
    },
    async handleCollect() {
      const { item } = this.data
      console.log(item)
      const params = {
        source_info: [
          {
            type: item.type
          }
        ],
        is_collect: item.is_collect ? 2 : 1
      }
      const trackList = []

      switch (item.type) {
        case 1:
          params.source_info[0].source_id = [item.id]
          trackList.push({
            is_top: !!item.is_top,
            asset_id: item.id,
            is_official: 1,
            asset_type_modify: item.type
          })
          break
        case 2:
          if (item.merge_id) {
            const images = item.merge_image.filter(v => {
              if (params.is_collect === 1) {
                return v.is_collect === 0
              } else {
                return v.is_collect === 1
              }
            })
            trackList.push(...images.map(v => ({
              is_top: !!item.is_top,
              asset_id: v.id,
              is_official: 1,
              asset_type_modify: item.type
            })))
            params.source_info[0].source_id = images.map(v => v.id)
          } else {
            params.source_info[0].source_id = [item.id]
            trackList.push({
              is_top: !!item.is_top,
              asset_id: item.id,
              is_official: 1,
              asset_type_modify: item.type
            })
          }

          params.count = params.source_info[0].source_id.length < item.merge_image.length ? 0 : 1
          break
        case 3:
          params.source_info[0].source_id = [item.id]
          trackList.push({
            is_top: !!item.is_top,
            asset_id: item.id,
            is_official: 1,
            asset_type_modify: item.type
          })
          break
        case 4:
          params.source_info[0].source_id = [item.id]
          trackList.push({
            is_top: !!item.is_top,
            asset_id: item.id,
            is_official: 1,
            asset_type_modify: item.type
          })
          break
        case 5:
          params.source_info[0].source_id = [item.id]
          trackList.push({
            is_top: !!item.is_top,
            asset_id: item.id,
            is_official: 1,
            asset_type_modify: item.type
          })
          break
        case 6:
          params.source_info[0].source_id = [item.id]
          trackList.push({
            is_top: !!item.is_top,
            asset_id: item.id,
            is_official: 1,
            asset_type_modify: item.type
          })
          break
        case 7:
          params.source_info[0].merge_id = [item.merge_id]
          trackList.push({
            is_top: !!item.is_top,
            asset_id: item.id,
            is_official: 1,
            asset_type_modify: item.type
          })
          break
        case 8:
          params.source_info[0].product_code = [item.product_code]
          trackList.push({
            is_top: !!item.is_top,
            asset_id: item.id,
            is_official: 2,
            asset_type_modify: item.type,
            official_url: item.url,
            official_upc: item.upc_code,
            official_product: item.product_name,
            official_name: item.product_code
          })
          break

        default:
          break
      }

      // 收藏
      const res = await collectApi(params)
      if (res.code === 200) {
        item.is_collect = params.is_collect === 2 ? 0 : 1
        item.collect_num = res.data.collect_num
        // 组合素材
        if (item.type === 2 && item.merge_id) {
          item.merge_image.forEach(image => {
            const v = res.data.source_list.find(k => k.source_id === image.id)
            if (v) {
              image.is_collect = params.is_collect === 2 ? 0 : 1
              image.collect_num = v.collect_num
            }
          })
        }
        if (item.type === 7 && item.merge_image.length) {
          item.merge_image.forEach(image => {
            image.is_collect = params.is_collect === 2 ? 0 : 1
          })
        }
        if (item.type === 8 && item.merge_image.length) {
          item.merge_image.forEach(image => {
            image.is_collect = params.is_collect === 2 ? 0 : 1
          })
        }
        wx.showToast(params.is_collect === 2 ? '已取消收藏' : '已收藏')

        this.triggerEvent('change', item)
        this.triggerEvent('collect', item)

        const currentPage = getCurrentPages().pop()
        trackList.forEach(v => {
          // 神策数据上报 -- 收藏素材
          wx.track.trackAssetFavorite({
            current_page: currentPage.route,
            entry_point: this.data.entryPoint || (this.data.enterType === 2 ? 3 : 1),
            library_type: this.data.enterType === 2 ? 3 : 1,
            operation_position: 1, // 固定值
            operation_type: params.is_collect,
            official_url: '非官网商品素材',
            official_upc: '非官网商品素材',
            official_product: '非官网商品素材',
            official_name: '非官网商品素材',
            ...v
          })
        })
      }
    }

  }
})

