// 获取全局应用程序实例对象
import { memberSearchPickerApi } from '@/api/member-search-api'
import { toDate, debounce } from '@/utils/tools'
// const app = getApp()
function getDefaultFormData() {
  return JSON.parse(JSON.stringify({
    /** 顾客姓名 */
    name: '',
    /** 顾客手机号 */
    phoneNumber: '',
    /** 会员卡号 */
    memberIds: '',
    /** 顾客类型 */
    type: '',
    /** 会员等级 */
    vipLevel: '',
    /** L12M消费记录 */
    priceRange: '',
    /** 购买日期 开始 */
    firstPurchase: '',
    /** 购买日期 结束 */
    lastPurchase: '',
    /** 产品类别 */
    department: '',
    /** 款号 */
    style: '',
    /** 产品描述 */
    collection: '',
    /** 性别 */
    gender: '',
    /** 生日月 */
    birthMonth: ''
  }))
}
// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    /** 筛选选项 */
    pickerData: {},
    /** 时间选择 */
    isShowDate: false,
    /** 显示更多筛选条件 */
    showMore: false,
    /** 顾客类型 下标 */
    typeIndex: -1,
    /** 会员等级 下标 */
    vipLevelIndex: -1,
    /** L12M消费记录 下标 */
    priceRangeIndex: -1,
    /** 产品类别 下标 */
    departmentIndex: -1,
    /** 性别 下标 */
    genderIndex: -1,
    /** 生日月 下标 */
    birthMonthIndex: -1,
    /** 表单 */
    form: getDefaultFormData(),
    /** 是否可以搜索 */
    disabled: true,
    currentDate: []
  },
  /**
   *
   * @param {Object} options
   * @param {Number} options.from  1-个人中心；2-顾客列表
  */
  onLoad (options) {
    this.setData({
      options
    })
    this.getMemberSearchPicker()
  },
  // 页面分享设置
  onShareAppMessage () {
    return {
      title: '会员查询',
      path: 'packages/customer/pages/member-search/member-search'
    }
  },
  /**
   * @description: 展开收起
   */
  showMoreSearch () {
    this.setData({
      showMore: !this.data.showMore
    })
  },
  /**
   * @description: 获取会员查询筛选条件接口
   */
  async getMemberSearchPicker () {
    wx.showLoading({
      title: '加载中'
    })
    const res = await memberSearchPickerApi()
    if (res.code === 200) {
      this.setData({
        pickerData: res.data
      })
      wx.hideLoading()
    }
  },
  /**
   * @description: input输入框
   * @param {*} e
   */
  handleInputChange (e) {
    const { key } = e.currentTarget.dataset
    const { value } = e.detail
    this.setData({
      [`form.${key}`]: value
    })
    this.check()
  },
  /**
   * @description: 失去焦点数据校验
   */
  handleBlur (e) {
    const { value } = e.detail
    const { key } = e.currentTarget.dataset
    // 手机号验证
    if (key === 'phoneNumber') {
      const reg = /^1[3456789]\d{9}$/
      if (!reg.test(value)) {
        wx.showToast('请输入正确的电话号码')
        this.setData({
          [`form.${key}`]: ''
        })
      }
    }
  },
  /**
   * @description 选择变更
   */
  handlePickerChange(e) {
    const { value } = e.detail
    const { key } = e.currentTarget.dataset
    const { pickerData } = this.data
    this.setData({
      [`${key}Index`]: value,
      [`form.${key}`]: pickerData[key][value].key
    })
    this.check()
  },
  /**
   * 打开时间筛选
   */
  handleSelectDate () {
    this.setData({ isShowDate: true })
  },
  /**
   * @description 时间选择
   */
  handleSelectDateConfirm (e) {
    console.log('时间选择', e)
    const { value } = e.detail
    this.setData({
      'form.firstPurchase': value[0] ? toDate(value[0], 'yyyy-MM').nowTime : '',
      'form.lastPurchase': value[1] ? toDate(value[1], 'yyyy-MM').nowTime : '',
      isShowDate: false,
      currentDate: value
    })
    this.check()
  },
  /**
   * @description 时间选择关闭
   */
  handleSelectDateCancel () {
    this.setData({ isShowDate: false })
  },
  /**
   * @description: 重置
   */
  handleReset () {
    this.setData({
      form: getDefaultFormData(),
      /** 顾客类型 下标 */
      typeIndex: -1,
      /** 会员等级 下标 */
      vipLevelIndex: -1,
      /** L12M消费记录 下标 */
      priceRangeIndex: -1,
      /** 产品类别 下标 */
      departmentIndex: -1,
      /** 性别 下标 */
      genderIndex: -1,
      /** 生日月 下标 */
      birthMonthIndex: -1,
      currentDate: []
    })
    this.check()
  },
  /**
   * @description: 前往顾客列表，传递参数
   */
  handleSearchConfirm: debounce(function() {
    const { form, options } = this.data
    // 将参数转成Object
    const params = {
      ...form,
      firstPurchase: form.firstPurchase ? form.firstPurchase.replace(/-/g, '') : '',
      lastPurchase: form.lastPurchase ? form.lastPurchase.replace(/-/g, '') : ''
    }
    // 去除空参数
    for (const key in params) {
      if (params[key] === '') {
        delete params[key]
      }
    }
    // 如果params为空，提示
    if (Object.keys(params).length === 0) {
      wx.showToast('请先输入搜索内容哦~')
      return
    }
    // 转成字符串
    const paramsStr = JSON.stringify(params)
    wx.navigateTo({
      url: `/packages/customer/pages/member-search-result/member-search-result?from=${options.from}&searchParams=${paramsStr}`
    })
  }),
  check() {
    const {form} = this.data

    this.setData({
      disabled: Object.values(form).filter(Boolean).length === 0
    })
  }
})
