const { collectApi } = require('@/api/material-api')

// Component Object
Component({
  properties: {
    show: {
      type: Boolean,
      value: false
    },
    item: {
      type: Object,
      observer: 'init'
    },
    type: {
      type: Number,
      value: 0
    },
    current: {
      type: Number,
      value: 0
    },
    showTips: {
      type: Boolean,
      value: false
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
    from: {
      type: Number,
      value: 0
    },
    showCount: {
      type: Boolean,
      value: true
    },
    libraryType: {
      type: Number,
      value: 0
    }
  },
  data: {
    list: []
  },
  methods: {
    init() {
      const { item } = this.data
      if (!item) return
      const isMergeImage = item.merge_image && item.merge_image.length > 0
      let list = []
      // 如果有组合
      if (isMergeImage) {
        list = item.merge_image
      } else {
        list = [item]
      }
      this.setData({
        list
      })
    },
    show() {
      this.setData({
        show: true
      })
      this.init()

      if (this.data.item.type !== 5 && this.data.list[this.data.current].type !== 5) {
        // 视频素材 放大的时候不上报， 播放的时候再上报
        this.trackAssetView()
      }
    },
    hide() {
      this.setData({
        show: false
      })
      this.triggerEvent('close')
    },
    handleChange(e) {
      console.log(e)
      const { current, source } = e.detail
      if (source === 'touch') {
        this.setData({
          current
        })
        this.handleHideTips()
        if (this.data.item.type !== 7 && this.data.list[current].type !== 5) {
          this.trackAssetView()
        }
      }
    },
    handleHideTips() {
      this.setData({
        showTips: false
      })
      this.triggerEvent('hide')
    },
    /**
     * @description 选择
     */
    handleSelect(e) {
      const { status, item, isMergeImage } = e.currentTarget.dataset
      // 返回 当前列表 下标 是否选中
      this.triggerEvent('select', {
        item: {
          type: this.data.item.type,
          ...item
        },
        isMergeImage,
        status
      })
    },
    /**
     * @description 保存
     */
    handleSave() {
      const { current, item } = this.data
      this.triggerEvent('save', {
        index: current,
        item: item
      })
    },
    async handleCollect() {
      const { item, current, list, from, enterType, entryPoint, libraryType } = this.data
      console.log(item)
      const params = {
        source_info: [
          {
            type: item.type
          }
        ],
        is_collect: list[current].is_collect ? 2 : 1
      }
      if (item.merge_id && item.merge_image.length > 0) {
        params.source_info[0].source_id = [item.merge_image[current].id]
      } else {
        params.source_info[0].source_id = [item.id]
      }

      // 收藏
      const res = await collectApi(params)
      if (res.code === 200) {
        // 组合素材
        if (item.type === 2 && item.merge_id) {
          // 组合素材
          item.merge_image[current].is_collect = params.is_collect === 2 ? 0 : 1
          item.merge_image[current].collect_num = res.data.source_list[0].collect_num
          item.collect_num = res.data.collect_num
          item.is_collect = item.merge_image.every(v => v.is_collect) ? 1 : 0
        } else {
          // 单素材
          item.is_collect = params.is_collect === 2 ? 0 : 1
          item.collect_num = res.data.collect_num
        }

        wx.showToast(params.is_collect === 2 ? '已取消收藏' : '已收藏')
        console.log(11111, item)

        this.triggerEvent('collect', {
          index: current,
          item
        })
        const currentPage = getCurrentPages().pop()
        wx.track.trackAssetFavorite({
          current_page: currentPage.route,
          entry_point: entryPoint || {
            1: 1,
            2: 3,
            3: 1,
            4: 4
          }[Number(from) === 2 ? 4 : enterType],
          library_type: libraryType || {
            1: 1,
            2: 3,
            3: 1,
            4: 4
          }[Number(from) === 2 ? 4 : enterType],
          operation_position: 2, // 固定值
          operation_type: params.is_collect,
          official_url: '非官网商品素材',
          official_upc: '非官网商品素材',
          official_product: '非官网商品素材',
          official_name: '非官网商品素材',
          is_top: !!item.is_top,
          asset_id: params.source_info[0].source_id[0],
          is_official: 1,
          asset_type_modify: item.type
        })
      }
    },
    handlePlay() {
      this.trackAssetView()
    },
    /** 预览上报 */
    trackAssetView() {
      const pages = getCurrentPages()
      const currentPage = pages[pages.length - 1].route === 'pagesMaterial/pages/select/select' ? pages[pages.length - 2] : pages[pages.length - 1]
      const { current, list, enterType, entryPoint, item, from } = this.data
      const value = list[current]
      // console.log('神策', value)
      // console.log('神策', item)
      const params = {
        source_id: item.id,
        task_id: '素材库无任务ID',
        is_official: 1,
        asset_type_modify: item.type,
        official_url: '非官网商品素材',
        official_upc: '非官网商品素材',
        official_product: '非官网商品素材',
        official_name: '非官网商品素材',
        current_page: currentPage.route,
        entry_point: entryPoint || {
          1: 1,
          2: 3,
          3: 1,
          4: 4
        }[Number(from) === 2 ? 4 : enterType],
        library_type: {
          1: 1,
          2: 3,
          3: 1,
          4: 4
        }[Number(from) === 2 ? 4 : enterType],
        is_top: !!item.is_top,
        is_favorite: !!item.is_collect
      }

      if (item.type === 2) {
        // 图片素材
        params.source_id = value.id
        params.is_favorite = !!value.is_collect
      }

      if (item.type === 8) {
        // 商品素材
        params.source_id = value.id
        params.is_official = 2
        params.asset_type_modify = value.type
        params.official_url = value.url
        params.official_upc = value.upc_code
        params.official_product = value.product_code
        params.official_name = value.product_name
      }

      // 神策数据上报 -- 查看素材
      wx.track.trackAssetView(params)
    }
  }
})

