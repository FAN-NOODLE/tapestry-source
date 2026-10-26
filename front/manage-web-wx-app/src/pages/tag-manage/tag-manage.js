import { addCustomerUsertagApi, deleteCustomerUsertagApi, getCustomerUsertagApi } from '@/api/customer-api'
// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    showDeleteModal: false, // 控制删除二次确认弹框显隐
    errorText: '', // 错误提示文案
    modalShow: false // 控制新增标签弹框显隐
  },

  // 生命周期函数--监听页面加载
  onLoad () {
    this.getCustomerUsertag()
  },
  // 获取筛选自定义标签
  getCustomerUsertag () {
    getCustomerUsertagApi({
      page: 1,
      size: 10000
    }).then(res => {
      const { code, data } = res
      if (code === 200) {
        this.setData({
          guideLabel: data
        })
      } else {
        wx.showToast(res.message)
      }
    }).catch(err => {
      wx.showToast(err.message)
    })
  },
  // 添加自定义标签
  addCustomerUsertag () {
    if (this.isClick) {
      return
    }
    this.isClick = true
    this.setData({
      currenInput: this.data.currenInput.trim()
    })
    addCustomerUsertagApi({
      name: this.data.currenInput.trim()
    }).then(res => {
      if (res.code === 200) {
        wx.showToast('标签新增成功')
        this.setData({
          modalShow: false,
          currenInput: ''
        })
        this.isClick = false
        this.getCustomerUsertag()
      } else {
        this.isClick = false
        // 标签已存在
        if (res.code === 820000) {
          this.setData({
            errorText: res.message
          })
        } else {
          wx.showToast(res.message)
        }
      }
    }).catch(err => {
      this.isClick = false
      // 标签已存在
      if (err.code === 820000) {
        this.setData({
          errorText: err.message
        })
      } else {
        wx.showToast(err.message)
      }
    })
  },
  // 打开删除确认弹框
  openDeleteModal (e) {
    this.data.currentDeleteId = e.currentTarget.dataset.id
    this.setData({
      showDeleteModal: true
    })
  },
  // 删除确认弹框取消按钮回调
  closeDeleteModel () {
    this.setData({
      showDeleteModal: false
    })
  },
  // 删除自定义标签
  deleteCustomerUsertag () {
    deleteCustomerUsertagApi({
      id: this.data.currentDeleteId
    }).then(res => {
      if (res.code === 200) {
        wx.showToast('标签已删除')
        this.setData({
          showDeleteModal: false
        })
        this.getCustomerUsertag()
      } else {
        wx.showToast(res.message)
      }
    }).catch(err => {
      wx.showToast(err.message)
    })
  },
  // 输入标签
  tagChange (e) {
    this.setData({
      errorText: '',
      currenInput: e.detail.value
    })
  },
  // 修改页面定位
  changePosition1 () {
    this.setData({
      bottom: '80px'
    })
  },
  changePosition2 () {
    this.setData({
      bottom: '0px'
    })
  },
  // 打开新增标签弹框
  openModal () {
    this.setData({
      modalShow: true,
      errorText: '',
      currenInput: ''
    })
  },
  // 关闭新增标签弹框
  closeModal () {
    this.setData({
      modalShow: false
    })
  }
})
