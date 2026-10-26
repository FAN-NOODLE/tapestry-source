const { globalData } = getApp()
const getLastYear = () => {
  let date = new Date()
  date.setHours(0, 0, 0, 0)
  return date.setFullYear(date.getFullYear() - 1) / 1000
}

Page({
  // 页面的初始数据
  data: {
    isIpx: globalData.isIpx,
    typeList: [
      {
        type: 1,
        title: ''
      },
      {
        type: 2,
        title: '选择指定导购'
      },
      {
        type: 3,
        title: '选择时间'
      },
      {
        type: 4,
        title: ''
      },
      {
        type: 5,
        title: '选择指定门店'
      }
    ],
    role: ''
  },

  // 生命周期函数--监听页面加载
  onLoad (options) {
    const {
      type = 1,
      role,
      showDate,
      mautiple,
      isBa,
      selectDateText
    } = options
    // 选择门店/部门时，页面传参赋值
    if (parseInt(type) === 1 || parseInt(type) === 4 || parseInt(type) === 5) {
      this.setData({
        hasIds: options.hasIds ? JSON.parse(options.hasIds) : [],
        mautiple: mautiple !== 'false'
      })
    }
    // 选择导购时，页面传参赋值
    if (parseInt(type) === 2) {
      this.setData({
        storeIds: options.storeIds ? JSON.parse(options.storeIds) : [],
        dpIds: options.dpIds ? JSON.parse(options.dpIds) : [],
        hasChooseGuideIds: options.hasChooseGuideIds
          ? JSON.parse(options.hasChooseGuideIds)
          : [],
        role,
        isBa: isBa !== 'false',
        mautiple: mautiple !== 'false'
      })
    }
    // 选择时间时，页面传参赋值
    if (parseInt(type) === 3) {
      globalData[selectDateText || 'selectDate'] = []
      this.setData({
        startTime: options.startTime
          ? options.startTime
          : new Date(new Date().setHours(0, 0, 0, 0)) / 1000,
        endTime: options.endTime ? options.endTime : 0,
        showDate: showDate !== 'false',
        selectDateText: selectDateText || 'selectDate',
        min: options.min || getLastYear(),
        max: options.max || new Date().setHours(23, 59, 59, 0) / 1000
      })
      console.log('onload时间', this.data.startTime, this.data.endTime)
    }
    setTimeout(() => {
      this.setData({
        currentSelectType: this.data.typeList[type - 1]
      })
      wx.setNavigationBarTitle({
        title: this.data.typeList[type - 1].title
      })
    }, 50)
  }
})
