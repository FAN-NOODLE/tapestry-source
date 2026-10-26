// 获取全局应用程序实例对象
// const app = getApp();
import { userCrmOrderApi } from '@/api/member-search-api'
// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    external_userid: '',
    unionId: '',
    showSelect: false,
    list: [],
    meta: {
      status: ''
    },
    showPicker: false,
    pickerValue: 0,
    pickerOptions: [
      {
        value: 'sku_name',
        label: '商品名称'
      },
      {
        value: 'style',
        label: '款号'
      },
      {
        value: 'ssociate',
        label: '员工ID'
      },
      {
        value: 'store_name',
        label: '店铺'
      }
    ]
  },
  onShow () {

  },
  // 生命周期函数--监听页面加载
  /**
   * @param {*} options
   * @param {*} options.external_userid
   * @param {*} options.unionId
   */
  onLoad (options) {
    this.setData({
      options
    })
    this.getData()
  },
  async getData () {
    try {
      // 加载中
      const {
        meta: { status },
        options
      } = this.data

      if (status === 'loading') return
      this.setData({
        'meta.status': 'loading'
      })
      const params = {
        external_userid: options.external_userid,
        union_id: options.unionId
      }
      const res = await userCrmOrderApi(params)

      if (res.code === 200) {
        this.setData({
          list: this.format(res.data.records || []),
          'meta.status': 'loaded'
        })
      }
    } catch (error) {
      this.setData({
        'meta.status': 'loadError'
      })
    }
  },

  format(list) {
    const { searchValue, pickerOptions, pickerValue } = this.data
    const key = pickerOptions[pickerValue].value
    return list
      .filter(item => {
        if (!searchValue) return true
        const searchRex = new RegExp(searchValue, 'i')
        return searchRex.test(item[key] || '')
      })
      .map(item => {
        for (let key in item) {
          if (item[key] === null) {
            item[key] = '--'
          }
        }
        return item
      })
  },

  /**
   * @description: 显示弹框
   */
  handleShowPicker () {
    this.setData({
      showPicker: true
    })
  },
  handlePickerCancel() {
    this.setData({
      showPicker: false
    })
  },
  handlePickerConfirm(e) {
    this.setData({
      showPicker: false,
      pickerValue: e.detail.index
    })
  },
  /**
   * @description: input blur
   */
  handleSearch (e) {
    this.setData({
      searchValue: e.detail.value
    })
    this.getData()
  },
  /**
 * @description:图片加载失败
 */
  skuImgError (e) {
    let index = e.currentTarget.dataset.index
    let recordList = this.data.recordList
    recordList[index].sku_image = `https://cdn-cos.coach.com.cn/uploads/2023/v3.2.0.1/default.png?time=${new Date().getTime()}`
    this.setData({
      recordList: recordList
    })
  }
})
