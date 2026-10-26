// 获取全局应用程序实例对象
// const app = getApp();
import { userCrmOrderApi } from '@/api/member-search-api'
import { getCustomerDetailApi } from '@/api/customer-api'
// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    external_userid: '',
    unionId: '',
    expand: false,
    showSelect: false,
    selectIndex: 0,
    selectList: [{
      key: 'style',
      value: '款号'
    }, {
      key: 'sku_name',
      value: '商品名称'
    }, {
      key: 'ssociate',
      value: '员工ID'
    }, {
      key: 'store_name',
      value: '店铺'
    }],
    lastPurchaseTime: '',
    L12MPurchaseAmount: '',
    L12MPurchaseCount: '',
    userDetail: {},
    recordList: [],
    showLoading: false
  },
  onShow () {

  },
  // 生命周期函数--监听页面加载
  onLoad (options) {
    console.log(options)
    // TODO: onLoad
    if (options.unionId && options.external_userid) {
      this.setData({
        unionId: options.unionId,
        external_userid: options.external_userid
      })
      this.getCustomerDetail()
    }
  },

  // 页面分享设置
  onShareAppMessage () {
    return {
      title: 'share title',
      path: 'pages/consume-record/consume-record'
    }
  },
  /**
   * @description: 获取顾客详情,获取消费就
   */
  async getCustomerDetail () {
    const params = {
      external_userid: this.data.external_userid,
      union_id: this.data.unionId
    }
    this.setData({
      showLoading: true
    })
    wx.showNavigationBarLoading()
    Promise.all([getCustomerDetailApi(params), userCrmOrderApi(params)]).then(res => {
      this.setData({
        showLoading: false
      })
      wx.hideNavigationBarLoading()
      if (res[0].code !== 200) {
        wx.showToast(res[0].message)
        return
      }
      if (res[1].code !== 200) {
        wx.showToast(res[1].message)
        return
      }
      let recordList = res[1].data.records || []
      // 将nll数据替换成--
      recordList.forEach(item => {
        // 循环item key 为null的值替换成--
        for (let key in item) {
          if (item[key] === null) {
            item[key] = '--'
          }
        }
      })
      this.setData({
        userDetail: res[0].data,
        lastPurchaseTime: res[1].data.lastPurchaseTime,
        L12MPurchaseAmount: res[1].data.L12MPurchaseAmount,
        L12MPurchaseCount: res[1].data.L12MPurchaseCount,
        recordList: recordList
      })
    }).catch(err => {
      this.setData({
        showLoading: false
      })
      wx.hideNavigationBarLoading()
      console.log(err)
    })
  },
  /**
     * 展开折叠
     */
  expand () {
    this.setData({ expand: !this.data.expand })
  },
  /**
   * @description: 点击页面关闭下拉框
   */
  closeSelect (e) {
    if (!this.data.showSelect) return
    this.setData({
      showSelect: false
    })
  },
  /**
   * @description:显示下拉框
   */
  showSelectArea () {
    this.setData({
      showSelect: !this.data.showSelect
    })
  },
  /**
 * @description:选择下拉款
 */
  chooseSelect (e) {
    console.log(e)
    const {index} = e.currentTarget.dataset
    this.setData({
      selectIndex: index,
      showSelect: !this.data.showSelect
    })
  },
  /**
   * @description: input blur
   */
  inputBlur (e) {
    this.setData({
      searchValue: e.detail.value
    })
    this.searchRecordList()
  },
  /**
 * @description: 根据selectIndex判断选择的是哪个，用于模糊搜索recordList
 */
  async searchRecordList () {
    let {selectIndex, searchValue} = this.data
    console.log(searchValue)
    const itemKey = this.data.selectList[selectIndex].key
    // 先去接口查询最新数据
    try {
      // 加载中
      this.setData({
        showLoading: true
      })
      const params = {
        external_userid: this.data.external_userid,
        union_id: this.data.unionId
      }
      const res = await userCrmOrderApi(params)
      this.setData({
        showLoading: false
      })
      if (res.code === 200) {
        let recordList = res.data.records
        let list = []
        // 将nll数据替换成--
        recordList.forEach(item => {
          // 循环item key 为null的值替换成--
          for (let key in item) {
            if (item[key] === null) {
              item[key] = '--'
            }
          }
        })
        const searchRex = new RegExp(searchValue, 'i')
        list = recordList.filter(item => searchRex.test(item[itemKey] && item[itemKey]))
        this.setData({
          recordList: list
        })
      }
    } catch (error) {
      console.log(error)
      this.setData({
        showLoading: false
      })
    }
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
