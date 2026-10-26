
import { getBrandMaterialListApi } from '@/api/material-api'

const { debounce } = require('@/utils/tools')

Component({
  properties: {
    show: {
      type: Boolean,
      value: false
    },
    value: {
      type: Array,
      value: []
    },
    disabled: {
      type: Boolean,
      value: false
    }
  },
  observers: {
    value: function (newVal, oldVal) {
      if (!newVal || newVal.length === 0) {
        this.setData({
          selection: [],
          selectionMap: {}
        })
        return
      }
      const list = newVal.filter(item => item.add_type === 1)
      this.setData({
        selection: list,
        selectionMap: list.reduce((prev, cur) => {
          prev[cur.upc] = cur
          return prev
        }, {})
      })
    },
    show: function (newVal, oldVal) {
      if (newVal) {
        this.setData({
          list: [],
          searchValue: '',
          'types.active': 0
        })
      }
    }
  },
  data: {
    empty_image: 'https://cdn-cos.coach.com.cn/uploads/2024/v4.2.0/icon_empty.png',
    list: [],
    /** 搜索类型 */
    types: {
      list: [
        {
          title: '商品款号',
          value: 'sku_color',
          placeholder: '商品款号'
        },
        {
          title: '商品名称',
          value: 'product_name',
          placeholder: '商品名称'
        },
        {
          title: '条形码',
          value: 'upc_code',
          placeholder: '条形码'
        }
      ],
      active: 0,
      show: false
    },
    searchValue: '',
    selection: [],
    selectionMap: {}
  },
  methods: {
    handleShow() {
      const { disabled } = this.data
      if (disabled) {
        wx.showToast('请先上传素材')
        return
      }
      this.setData({
        show: true
      })
    },
    handleClose() {
      this.setData({
        show: false
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
    /** 搜索 */
    handleInput(e) {
      const { value } = e.detail
      this.setData({
        searchValue: value
      })
      if (value.trim().length > 2) {
        this.handleSearch()
      }
    },
    /** 清除搜索 */
    handleClearSearch() {
      this.setData({
        searchValue: ''
      })
    },
    /**
     * 扫码
     */
    handleScan() {
      wx.scanCode({
        success: (res) => {
          this.setData({
            searchValue: res.result,
            isScan: true,
            'types.active': 2
          }, () => {
            this.handleSearch()
          })
        }
      })
    },
    /**
     * 搜索
     * @returns
     */
    handleSearch: debounce(function() {
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
    }, 500, false),

    /**
     * @description 获取列表
     */
    async getData() {
      const { last_code, size } = this.data

      try {
        const params = {
          size,
          last_code,
          include_none_image_product: true,
          ...this.getSearchParams()
        }
        this.setData({
          status: 'loading',
          isScan: false
        })
        const res = await getBrandMaterialListApi(params)
        if (res.code === 200) {
          const data = res.data.length ? res.data : []
          // 设置选中素材
          const list = data.map((item) => {
            item.image_url = (item.brand_source_list.find(image => image.is_main === 1) || item.brand_source_list[0] || {}).url || ''
            return item
          })

          const newList = (!last_code ? [] : this.data.list).concat(list)

          this.setData({
            last_code: newList.length ? newList[newList.length - 1].id : '',
            list: newList,
            status: data.length > 0 ? 'loadMore' : 'loaded',
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
     * @description 获取搜索参数
     */
    getSearchParams() {
      const { types, searchValue, isScan } = this.data
      const key = isScan ? 'upc_code' : types.list[types.active].value
      return {
        [key]: searchValue
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
    transform(item) {
      return {
        upc: item.upc_code,
        sku: item.sku_item,
        color: item.color,
        url: item.image_url || '',
        name: item.product_name || '',
        add_type: 1, // 1.手动添加 2.AI识别
        score: '', // 置信度 （AI识别时有）
        checked: true
      }
    },

    handleSelect(e) {
      const { selectionMap } = this.data
      const { item } = e.currentTarget.dataset

      const aiGoods = this.data.value.filter(item => item.add_type === 2)
      const manualGoods = this.data.value.filter(item => item.add_type === 1)
      if (manualGoods.length >= 10 && !selectionMap[item.upc_code]) {
        wx.showToast('最多只能添加10个手动添加商品')
        return
      }
      const upcList = aiGoods.map(item => item.upc)
      // 手动添加的素材，不重复添加
      if (upcList.includes(item.upc_code)) {
        wx.showToast(`${item.sku_item}_${item.color}已存在AI识别商品，不重复添加`)
        return
      }

      if (selectionMap[item.upc_code]) {
        delete selectionMap[item.upc_code]
      } else {
        selectionMap[item.upc_code] = this.transform(item)
      }
      this.setData({
        selectionMap,
        selection: Object.values(selectionMap)
      })

      this.triggerEvent('change', {
        value: [...aiGoods, ...(this.data.selection.filter(item => !upcList.includes(item.upc)))]
      })
    }

  }
})
