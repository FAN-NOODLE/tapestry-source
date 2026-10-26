import {
  getOrderDetail,
  getDeliveryCompany,
  handleToDelivery
} from '@/api/order'
import { formatDate } from '@/utils/tools'

const { host, apiRoot, env } = require('@/lib/config')

const { globalData } = getApp()

Page({
  data: {
    /** 物流公司数据 */
    deliveryCompany: [],

    /** 发货详情 */
    detail: {},

    /** 订单ID */
    id: '',

    /** 自提或快递订单 */
    delivery_type: '',

    /** 是否修改 */
    isEdit: false,

    /** 是否正在加载 */
    isLoading: true
  },

  async onLoad ({ id, delivery_type, isEdit }) {
    this.setData({
      id,
      delivery_type,
      isEdit: isEdit === 'true'
    })

    if (this.data.isEdit) {
      wx.setNavigationBarTitle({
        title: '修改发货信息'
      })
    }

    await this.getDeliveryCompany()
    this.getDeliveryInfo()
  },

  /** 快递公司 */
  async getDeliveryCompany () {
    try {
      let res = await getDeliveryCompany()
      if (res.code === 200) {
        let deliveryCompany = []

        for (let key in res.data) {
          deliveryCompany.push({
            label: key,
            value: res.data[key]
          })
        }
        this.setData({
          deliveryCompany
        })
      } else {
        wx.showToast(res.message)
      }
    } catch (error) {}
  },

  /** 发货详情 */
  getDeliveryInfo () {
    wx.showLoading()
    let params = {
      id: this.data.id
    }
    getOrderDetail(params)
      .then((res) => {
        if (res.code === 200) {
          let { delivery_time, sales_order_no, store_name, pay_time } =
            res.data

          if (delivery_time) {
            delivery_time = formatDate(delivery_time * 1000, 'yyyy-MM-dd')
          }

          sales_order_no = sales_order_no.split('-')[1]

          this.setData({
            isLoading: false,
            startTime: formatDate(pay_time * 1000, 'yyyy-MM-dd'),
            endTime: formatDate(new Date().getTime(), 'yyyy-MM-dd'),
            detail: {
              logistics_no: res.data.logistics_no,
              logistics_com: res.data.logistics_com || '',
              logistics_name: res.data.logistics_name || '',
              tickets: res.data.tickets || '',
              delivery_time,
              sales_order_no,
              store_name
            }
          })
          wx.hideLoading()
        } else {
          wx.hideLoading()
          wx.showToast(res.message)
        }
      })
      .catch((error) => {
        console.log(error)
      })
  },

  /** 选择快递公司 */
  handleToChooseCompany (e) {
    const { value } = e.detail

    this.setData({
      'detail.logistics_name': this.data.deliveryCompany[value].value,
      'detail.logistics_com': this.data.deliveryCompany[value].label
    })
  },

  /** 输入快递单号 */
  handleInput (e) {
    let value = (e.detail.value || '').trim()

    const reg = /[\u4e00-\u9fa5]/gi

    this.setData({
      'detail.logistics_no': value.replace(reg, '')
    })
  },

  /** 扫描快递单号 */
  handleToScan () {
    wx.scanCode({
      onlyFromCamera: true,
      scanType: ['barCode'],
      success: (res) => {
        if (res.errMsg === 'scanCode:ok') {
          this.setData({
            'detail.logistics_no': res.result
          })
        }
      },
      fail: () => {},
      complete: () => {}
    })
  },

  /** 输入销售单号 */
  handleInputSalesOrderNo (e) {
    let value = (e.detail.value || '').trim()

    const reg = /[\u4e00-\u9fa5]/gi

    this.setData({
      'detail.sales_order_no': value.replace(reg, '')
    })
  },

  /** 选择录单日期 */
  bindJoinDateChange (e) {
    const { value } = e.detail

    this.setData({
      'detail.delivery_time': value
    })
  },

  /** 上传小票 */
  handleUpload () {
    const self = this

    wx.chooseImage({
      count: 1,
      sizeType: ['compressed'],
      sourceType: ['album', 'camera'],
      success: async (result) => {
        console.log('chooseImageSync', result)
        if (result.errMsg === 'chooseImage:ok') {
          let header = {
            'content-type': 'multipart/form-data'
            // 'Enterprise-Hash': 'ecd85280aa135bbd0108dd6aa424565a'// 临时设置
          }
          if (
            (env === 'beta' || env === 'test_beta' || env === 'test') &&
            globalData &&
            globalData.temporary_hash
          ) {
            header['Enterprise-Hash'] = globalData.temporary_hash
          }

          wx.showLoading({
            title: '上传中'
          })
          wx.uploadFile({
            url: host + apiRoot + 'partner-member/upload',
            filePath: result.tempFiles[0].path,
            name: 'file',
            header: header,
            success: function (res) {
              console.log('上传图片成功', res)
              let data = JSON.parse(res.data)
              if (res.errMsg === 'uploadFile:ok') {
                self.setData({
                  'detail.tickets': data.data.url
                })
              } else {
              }
            },
            fail: function () {
              console.log('上传图片失败')
            },
            complete: function () {
              wx.hideLoading()
            }
          })
        }
      },
      fail: (error) => {
        console.log(error)
      }
    })
  },

  /** 发货 */
  handleToDelivery () {
    const {
      logistics_no,
      logistics_name,
      logistics_com,
      tickets,
      delivery_time,
      sales_order_no,
      store_name
    } = this.data.detail

    // 判断是否可以提交
    const condition = this.data.delivery_type === 2
      ? (!logistics_name || !logistics_no || !tickets || !delivery_time || !sales_order_no)
      : (!tickets || !delivery_time || !sales_order_no)
    if (condition) {
      return
    }

    let params = {
      id: this.data.id,
      logistics_no,
      logistics_name,
      logistics_com,
      tickets,
      delivery_time: new Date(delivery_time).getTime() / 1000,
      sales_order_no: `${store_name}-${sales_order_no}`,
      delivery_type: this.data.delivery_type
    }

    handleToDelivery(params)
      .then((res) => {
        if (res.code === 200) {
          wx.showToast(
            `${this.data.isEdit ? '修改' : '发货'}成功`,
            false,
            () => {
              this.emitUpdateList()
              wx.navigateBack({ delta: 1 })
            }
          )
        } else {
          wx.showToast(res.message)
        }
      })
      .catch((error) => {
        console.log(error)
      })
  },
  /** 更新列表 */
  emitUpdateList () {
    try {
      const eventChannel = this.getOpenerEventChannel()
      eventChannel.emit('updateData', {})
    } catch (e) {
      console.log('updateData', e)
    }
  }
})
