import { getMaterialListApi, getMaterialNewNumApi, checkFirstApi, getMediaInfoApi } from '@/api/material-api'
import { getCurrentPage } from '@/components/common/common/utils'
import { debounce } from '@/utils/tools'
import { getCosCompressUrl } from '@/utils/upload/utils'
import dayjs from 'dayjs'

// Page Object
Component({
  options: {
    virtualHost: false,
    multipleSlots: true,
    styleIsolation: 'shared'
  },
  properties: {
    externalIds: {
      type: Array,
      value: []
    }
  },
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
    selection: [],
    /** 选中文本列表 */
    selectedTextList: [],
    /** 选中图片视频列表 */
    selectedMediaList: [],
    /** 宫格提示 */
    showFirstTips: false,
    /** 预览滚动提示 */
    showFirstPreviewTips: false,
    /** 当前操作的数据下标 */
    current: -1,
    /** 一次性下载最大文本数据 */
    textMax: 1,
    /** 一次性下载最大图片视频数据 */
    mediaMax: 9,
    /** 预览下标 */
    previewIndex: 0,
    /** 加载状态 */
    loading: true,
    /** 已选文本 */
    textNum: 0,
    /** 已选图片或者视频 */
    mediaNum: 0
  },

  lifetimes: {
    attached() {
      this.getData()
      this.checkFirst()
    },
    detached() {
      // 在组件实例被从页面节点树移除时执行
    }
  },
  pageLifetimes: {
    // 组件所在页面的生命周期函数
    show() {
      // 在组件实例进入页面节点树时执行
      this.getNewNumber()
      const page = getCurrentPage()
      console.log(1111, page)
    },
    hide() {},
    resize() {}
  },
  methods: {

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
      const { tabs, active, query } = this.data
      const item = tabs[active]

      const params = {
        ...query,
        size: item.size,
        page,
        type: item.value
      }
      this.setData({
        [`tabs[${active}].status`]: 'loading'
      })

      try {
        const res = await getMaterialListApi(params)
        if (res.code === 200) {
        // todo:合并文本和图片视频
          const selectedList = [...this.data.selectedTextList, ...this.data.selectedMediaList]
          const ids = selectedList.map(child => child.id)
          // 设置选中素材
          const list = res.data.list.map(child => {
            child.checked = ids.includes(child.id)
            if (child.type === 2 && child.merge_id) {
              child.merge_image.forEach(image => {
                image.checked = ids.includes(image.id)
                image.is_top = child.is_top
              })
            } else if (child.type === 7 && child.merge_image.length) {
              child.merge_image.forEach(image => {
                image.is_top = child.is_top
                image.is_collect = child.is_collect
              })
            }
            return child
          })
          item.page = page
          item.update = false
          item.list = (page === 1 ? [] : item.list).concat(list)
          item.status = res.data.total > res.data.page * res.data.size ? 'loadMore' : 'loaded'
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
   * @description 切换宫格
   */
    handleGideChange(e) {
      const { value } = e.detail
      const { active } = this.data
      this.setData({
        gridIndex: value,
        [`tabs[${active}].column`]: Number(value) === 0 ? 2 : 3,
        [`tabs[${active}].gutter`]: Number(value) === 0 ? 30 : 20
      })
      this.handleHideTips()
    },
    /**
   * @description 选择素材
   * @param {Object} e
   * @param {Object} e.detail
   */
    handleSelect(e) {
      const { textNum, mediaNum, textMax, mediaMax, selection } = this.data
      const { item, status, isMergeImage } = {...e.currentTarget.dataset, ...e.detail}
      if (status === 2) {
        if (isMergeImage) {
        // 取消选中
          const ids = item.merge_image.map(val => val.id)
          // 如果是选中状态，取消选中
          this.data.selection = selection.filter(val => !ids.includes(val.id))
        } else {
          const index = selection.findIndex(val => val.id === item.id)
          this.data.selection.splice(index, 1)
        }
      } else {
        // 选中
        // 朋友圈素材 文本 | 文本素材
        if (textNum === textMax && ((item.type === 7 && item.text_content) || item.type === 1)) {
          wx.showToast(`最多可选${textMax}个文本素材`)
          return
        }
        // 朋友圈素材 视频/图片 | 图片素材 ｜ 视频素材
        if (mediaMax === mediaNum && ((item.type === 7 && (item.original_url || (item.merge_id && item.merge_image.length > 0))) || item.type === 2 || item.type === 5 || item.type === 3)) {
          wx.showToast(`最多可选${mediaMax}个图片/视频/链接素材`)
          return
        }


        if (isMergeImage) {
          const ids = selection.map(val => val.id)
          // 如果是取消状态，选中
          const list = item.merge_image.filter(val => !ids.includes(val.id)).map(val => {
            return {
              _type: 2,
              ...val
            }
          })
          // 如果有文本素材，添加文本内容
          if (item.text_content) {
            selection.unshift({
              _type: 1,
              ...item,
              content: item.text_content
            })
          }

          // 添加选中的素材
          this.data.selection = [...selection, ...list].splice(0, item.text_content || textNum === 1 ? textMax + mediaMax : mediaMax)
        } else {
          if (item.type === 3) {
            // 链接素材
            this.data.selection.push({
              _type: 3,
              ...item
            })
          } else if (item.content || item.text_content) {
            this.data.selection.unshift({
              _type: 1,
              ...item
            })
          } else {
            let type = item.type
            if (type === 7 && item.merge_id && item.merge_image.length > 0) {
              type = 2
            } else if (type === 7 && item.original_url) {
              type = 5
            }

            this.data.selection.push({
              _type: type,
              url: item.original_url,
              ...item
            })
          }
        }
      }

      // 更新数据
      this.setData({
        selection: this.data.selection
      }, () => {
        this.computedSelection()
      })
    },
    /**
   * 更新选中列表
   */
    computedSelection() {
      const {selection} = this.data
      const selectedTextList = selection.filter(item => item._type === 1)
      const selectedMediaList = selection.filter(item => item._type !== 1)

      this.setData({
        selectedTextList,
        selectedMediaList,
        textNum: selectedTextList.length,
        mediaNum: selectedMediaList.length
      })
    },
    /**
   * @description 展示组合素材选择弹框
   */
    handleShowGroupSelect(e) {
      const { index } = e.currentTarget.dataset
      this.setData({
        current: index
      }, () => {
        this.selectComponent('#com-popup-save').show()
      })
    },

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
   * @description todo:群发接口
   */
    handleSendAll: debounce(async function() {
      const currentPage = getCurrentPages().pop()
      const { selectedTextList, selectedMediaList, externalIds } = this.data
      // console.log(selectedTextList, selectedMediaList)
      // externalIds.forEach(item => {
      //   selectedTextList.forEach(child => {
      //     wx.track.trackAssetSend({
      //       asset_id: child.id,
      //       task_id: '群发无法获取',
      //       customer_id: item,
      //       operation_type: 0,
      //       is_official: 1,
      //       asset_type_modify: child._type,
      //       official_url: '非官网商品素材',
      //       official_upc: '非官网商品素材',
      //       official_product: '非官网商品素材',
      //       official_name: '非官网商品素材',
      //       current_page: currentPage.route,
      //       entry_point: 3,
      //       library_type: 3,
      //       is_top: !!child.is_top,
      //       is_favorite: !!child.is_collect
      //     })
      //   })
      //   selectedMediaList.forEach(child => {
      //     wx.track.trackAssetSend({
      //       asset_id: child.id,
      //       task_id: '群发无法获取',
      //       customer_id: item,
      //       operation_type: 0,
      //       is_official: 1,
      //       asset_type_modify: child._type,
      //       official_url: '非官网商品素材',
      //       official_upc: '非官网商品素材',
      //       official_product: '非官网商品素材',
      //       official_name: '非官网商品素材',
      //       current_page: currentPage.route,
      //       entry_point: 3,
      //       library_type: 3,
      //       is_top: !!child.is_top,
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
      if (selectedTextList.length === 0 && selectedMediaList.length === 0) return
      // todo:筛选出所有视频素材
      const videoList = selectedMediaList.filter(item => item._type === 5)
      const imageList = selectedMediaList.filter(item => item._type === 2)
      const linkList = selectedMediaList.filter(item => item._type === 3)
      const types = {}
      const params = {
        externalUserIds: externalIds,
        text: {
          content: selectedTextList[0] ? selectedTextList[0].content : ''
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
      // todo:获取链接素材
      if (linkList.length > 0) {
        params.attachments = params.attachments.concat(...linkList.map(item => {
          types[item.id] = 3
          return {
            msgtype: 'link',
            link: {
              title: item.title,
              url: item.link_url,
              desc: item.description,
              imgUrl: item.cover_url
            }
          }
        }))
      }
      console.log(params.attachments)
      wx.qy.shareToExternalContact({
        ...params,
        success: function (res) {
          console.log(res)
          // todo:清空选中数据
          that.handleClearSelected()
          // 神策数据上报 -- 群发素材
          wx.track.trackBatchMsgConfirm({
            asset_coverage: externalIds.length,
            asset_list: selectedMediaList.map(item => item.id.toString()),
            customer_list: externalIds,
            operation_type: 0
          })
          // 神策数据上报 -- 发送素材

          externalIds.forEach(item => {
            selectedTextList.forEach(child => {
              wx.track.trackAssetSend({
                asset_id: child.id,
                task_id: '群发无法获取',
                customer_id: item,
                operation_type: 0,
                is_official: 1,
                asset_type_modify: child._type,
                official_url: '非官网商品素材',
                official_upc: '非官网商品素材',
                official_product: '非官网商品素材',
                official_name: '非官网商品素材',
                current_page: currentPage.route,
                entry_point: 3,
                library_type: 3,
                is_top: !!child.is_top,
                is_favorite: !!child.is_collect
              })
            })
            selectedMediaList.forEach(child => {
              wx.track.trackAssetSend({
                asset_id: child.id,
                task_id: '群发无法获取',
                customer_id: item,
                operation_type: 0,
                is_official: 1,
                asset_type_modify: child._type,
                official_url: '非官网商品素材',
                official_upc: '非官网商品素材',
                official_product: '非官网商品素材',
                official_name: '非官网商品素材',
                current_page: currentPage.route,
                entry_point: 3,
                library_type: 3,
                is_top: !!child.is_top,
                is_favorite: !!child.is_collect
              })
            })
          })
          const page = getCurrentPage()
          wx.navigateBack({
            delta: Number(page.options.back || 2)
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
    /**
   * @description 清空选中数据
   */
    handleClearSelected() {
      this.setData({
        selection: [],
        selectedMediaList: [],
        selectedTextList: [],
        textNum: 0,
        mediaNum: 0
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
    }
  }
})

