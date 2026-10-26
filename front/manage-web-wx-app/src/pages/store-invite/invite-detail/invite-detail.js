// 获取全局应用程序实例对象
import {
  guideTimeApi,
  serviceListApi,
  orderDetailApi,
  orderEditApi,
  orderCancelApi,
  orderConfirmStore,
  getStaffLists,
  changeStatus,
  commentShareApi
} from '../../../api/guide-invite'
import { toDate, formatWithArray } from '../../../utils/tools'
import { StatusList } from '../constants'

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    id: 0,
    serverIndex: -1, // 服务选择的下标
    serverList: [], // 服务列表
    reserveNumber: '', // 备注
    role: 0, // 导购角色
    statusList: StatusList,
    showLoading: true, // 页面loading
    dateModel: false, // 是否可以展示选择时间
    detail: {}, // 预约单详情
    selectMember: {}, // 店长可选------选择的员工
    selectServer: {}, // 店员可选------选择服务
    staffChange: false, // 店员是否改变了数据
    staffChangeFiled: [], // 店员改变的数据有哪些
    goShare: false, // 去分享了
    tipModel: false, // 选择时间是是否占用
    orderNumber: '',
    // orderInfo: {},
    // appid: require('../../../lib/config').appid,
    serviceScoreList: [],
    storeScoreList: [],
    staffLists: [],
    store_address: '',
    reserve_remark: '',
    isComment: false,
    is_store: 0, // 是否到店
    customer_status: '', // 顾客类型
    memberIndex: -1, // 导购选择的下标
    isConfirmModal: false, // 二次确认弹窗
    confirmDesc: '', // 确认文案标题
    confirmDescTwo: '', // 确认文案说明
    operationType: '', // 操作类型 - 修改/取消
    showBottomBtn: false, // 是否展示底部按钮
    showCommentView: false, // 是否展示评价相关内容
    timeSlot: [] // 时间段
  },
  // 生命周期函数--监听页面加载
  onLoad (options) {
    this.data.id = Number(options.id)
    this.data.reserveNumber = options.reserveNumber
    this.data.index = options.index// 从列表某个索引点击进入
    this.setData({
      role: Number(options.role || 3)
    })
    wx.hideShareMenu({
      withShareTicket: true,
      menus: ['shareAppMessage', 'shareTimeline']
    })
    wx.hideHomeButton()
  },
  onShow () {
    this.init()
  },
  // 设置星级评价显示
  setStarList (score) {
    // 评级星级
    const checkedStar = new Array(score).fill({
      checked: true
    })
    // 灰色星级
    const uncheckedStar = new Array(5 - score).fill({
      checked: false
    })

    const starList = checkedStar.concat(uncheckedStar)
    return starList
  },
  // 可编辑的内容: 门店地址/导购备注
  changeInput (e) {
    const key = e.currentTarget.dataset.key
    const value = e.detail.value
    this.setData({
      [key]: value
    })
    this.hasChange()
  },
  init () {
    let p = [
      this.getServerList(), // 获取服务列表
      this.getTime() // 获取一个月之内时间段预约的人数
    ]
    Promise.all(p)
      .then(() => {
        this.getDetail()
      })
      .catch((err) => {
        this.setData({
          showLoading: false
        })
        console.log(err, '出错')
      })
  },
  // 获取导购列表
  getStaffLists (unionId) {
    getStaffLists({ unionId })
      .then((res) => {
        if (res.code === 200) {
          // 有数据返回数组，无数据返回{}
          if (JSON.stringify(res.data) === '{}') {
            this.setData({
              staffLists: []
            })
            return
          }
          this.setData({
            staffLists: res.data
          })
        }
      })
      .catch((err) => {
        console.log(err)
      })
  },
  // 是否有改变
  hasChange () {
    const {
      detail,
      selectTime,
      selectServer,
      selectMember,
      store_address,
      reserve_remark
    } = this.data
    let staffChange = false
    // 时间/服务/导购/门店地址/导购备注
    if (
      selectServer.id !== detail.reserve_id ||
      selectTime.startTime !== detail.reserve_time ||
      selectMember.mobile !== detail.mobile ||
      store_address !== detail.store_address ||
      reserve_remark !== detail.reserve_remark
    ) {
      staffChange = true
    }

    // 保存变更的字段
    let staffChangeFiled = []
    if (selectTime.startTime !== detail.reserve_time) {
      staffChangeFiled.push('到店时间')
    }
    if (selectServer.id !== detail.reserve_id) {
      staffChangeFiled.push('服务类型')
    }
    this.setData({
      staffChange,
      staffChangeFiled
    })
  },
  // 显示确认弹窗
  showConfirm (e) {
    const dataset = e.currentTarget.dataset
    if (dataset && dataset.type) {
      if (dataset.type === 'update') {
        if (!this.data.staffChange) return
        this.setData({
          isConfirmModal: true,
          confirmDesc: '确认修改',
          confirmDescTwo: '确定修改当前预约单？',
          operationType: 'update'
        })
      } else {
        // 如果不是自己的邀请单，不能操作
        if (!this.data.detail.is_own) return
        this.setData({
          isConfirmModal: true,
          confirmDesc: '',
          confirmDescTwo: '确定取消当前预约单？',
          operationType: 'cancel'
        })
      }
    }
  },
  // 取消订单
  cancelOrder () {
    const { reserveNumber, id, detail } = this.data
    const params = {
      id,
      reserve_number: reserveNumber
    }
    wx.showLoading({
      title: '取消中'
    })
    orderCancelApi(params)
      .then((res) => {
        wx.hideLoading()
        if (res.code === 200) {
          this.setData({
            isConfirmModal: false
          })
          wx.showToast('取消成功')
          this.getDetail()
          // 神策数据上报 -- 取消邀约单
          wx.track.trackAppointmentCancel({
            reserve_number: detail.reserve_number,
            user_external_id: detail.user_external_id,
            reserve_name: detail.serve_name,
            reserve_time: detail.reserve_time,
            time_period: detail.reserve_time,
            store_name: detail.store_name,
            store_code: detail.store_code || detail.store_name,
            cancel_type: 1
          })
          return
        }
        wx.showToast(res.message)
      })
      .catch((err) => {
        wx.hideLoading()
        console.log(err)
      })
  },
  // 确认修改
  updateOrder () {
    if (!this.data.staffChange) return
    const { selectTime, selectServer, id, store_address, reserve_remark, staffChangeFiled, detail } =
      this.data
    const params = {
      id,
      reserve_id: selectServer.id,
      reserve_time: selectTime.startTime,
      reserve_end_time: selectTime.endTime,
      store_address,
      reserve_remark,
      mobile: this.data.selectMember.mobile || '' // 所选中的导购的 手机号
    }
    wx.showLoading({
      title: '加载中'
    })
    orderEditApi(params)
      .then((res) => {
        wx.hideLoading()
        if (res.code === 200) {
          wx.showToast('修改成功')
          this.setData({
            isConfirmModal: false
          })
          this.getDetail()
          // 神策数据上报 -- 修改邀约单
          wx.track.trackAppointmentModify({
            reserve_number: detail.reserve_number,
            user_external_id: detail.user_external_id,
            change_item: staffChangeFiled.join('&'),
            modify_reserve_name: selectServer.reserve_name,
            modify_reserve_time: selectTime.startTime,
            modify_time_period: selectTime.startTime,
            store_code: detail.store_code || detail.store_name,
            store_name: detail.store_name
          })
          return
        }
        wx.showToast(res.message)
      })
      .catch((err) => {
        wx.hideLoading()
        console.log(err)
      })
  },
  // 弹窗点击确认
  sureOperation () {
    if (this.data.operationType === 'update') {
      this.updateOrder()
    } else {
      this.cancelOrder()
    }
  },
  tipShowModel () {
    this.setData({
      tipModel: false
    })
  },
  // 确认顾客到店
  orderConfirmStore () {
    const { reserveNumber, id } = this.data
    const params = {
      id,
      reserve_number: reserveNumber
    }
    wx.showLoading({
      title: '加载中'
    })
    orderConfirmStore(params)
      .then((res) => {
        wx.hideLoading()
        if (res.code === 200) {
          wx.showToast('确认成功')
          this.getDetail()
          return
        }
        wx.showToast(res.message)
      })
      .catch((err) => {
        wx.hideLoading()
        console.log(err)
      })
  },
  getTime () {
    return new Promise((resolve, reject) => {
      guideTimeApi()
        .then((res) => {
          if (res.code === 200) {
            this.setData({
              timeObj: res.data,
              showTime: true
            })
            resolve()
          }
          reject(res.message)
        })
        .catch((err) => {
          reject(err)
        })
    })
  },
  // 关闭二次确认弹窗弹窗
  closeTipsModel () {
    this.setData({
      isConfirmModal: false
    })
  },
  // 设置导购回显
  setStaff (info) {
    // 查询导购是否存在
    const memberIndex = this.data.staffLists.findIndex(
      (item) => item.mobile === info.mobile
    )
    let selectMember = {}
    // 兼容详情服务id被删之后的展示
    if (memberIndex === -1) {
      selectMember = {
        mobile: info.mobile || '',
        members_name: info.member_name
      }
    } else {
      selectMember = this.data.staffLists[memberIndex]
    }
    return { selectMember, memberIndex }
  },
  // 设置服务回显
  setServer (info) {
    // 查询服务是否存在
    const serverIndex = this.data.serverList.findIndex(
      (item) => item.id === info.reserve_id
    )
    let selectServer = {}
    // 兼容详情服务id被删之后的展示
    if (serverIndex === -1) {
      selectServer = {
        id: info.reserve_id,
        reserve_img: info.serve_img,
        reserve_name: info.serve_name
      }
    } else {
      selectServer = this.data.serverList[serverIndex]
    }
    return { selectServer, serverIndex }
  },
  // 获取邀约单详情信息
  getDetail () {
    const { id, reserveNumber } = this.data
    const params = {
      id,
      reserve_number: reserveNumber
    }
    orderDetailApi(params)
      .then((res) => {
        if (res.code === 200) {
          const orderInfo = res.data

          // 根据 user_unionid 获取导购列表
          this.getStaffLists(res.data.user_unionid)

          const date = new Date(orderInfo.reserve_time * 1000)
          // todo: 到店开始结束时间
          let select_hours = ''
          if (orderInfo.reserve_end_time === 0) {
            select_hours = `${date.getHours()}:00 - ${date.getHours() + 1}:00`
          } else {
            const endDate = new Date(orderInfo.reserve_end_time * 1000)
            select_hours = `${date.getHours()}:00 - ${endDate.getHours()}:00`
          }
          const select_name = toDate(
            date,
            'yyyy-MM-dd'
          ).nowTime
          const selectTime = {
            select_name,
            inviteNum: 0,
            name: select_hours,
            startTime: orderInfo.reserve_time
          }
          // 处理评价星级显示
          const serviceScoreList = this.setStarList(res.data.service_score) // 活动安排
          const storeScoreList = this.setStarList(res.data.store_score) // 整体服务

          res.data.created_at_str = toDate(
            res.data.created_at,
            'yyyy-MM-dd hh:mm'
          ).nowTime
          res.data.reserve_order_time_str = toDate(
            res.data.reserve_order_time,
            'yyyy-MM-dd'
          ).nowTime
          res.data.reserve_order_price =
            res.data.reserve_order_price &&
            '¥' + formatWithArray(res.data.reserve_order_price / 100)

          const { selectServer, serverIndex } = this.setServer(orderInfo)
          const { selectMember, memberIndex } = this.setStaff(orderInfo)
          this.setData({
            serverIndex,
            memberIndex,
            selectTime,
            selectMember,
            selectServer,
            showLoading: false,
            staffChange: false,
            detail: res.data,
            serviceScoreList,
            storeScoreList,
            // 条件： is_comment-2 未评价 & status-1 邀约成功 / status-2 邀约完成
            isComment:
              orderInfo.is_comment === 2 &&
              (orderInfo.status === 1 || orderInfo.status === 2),
            // 条件： status-1 邀约成功 / status-2 邀约完成  & is_comment-1 已评价
            showCommentView:
              (orderInfo.status === 1 || orderInfo.status === 2) &&
              orderInfo.is_comment === 1,
            store_address: orderInfo.store_address,
            reserve_remark: orderInfo.reserve_remark,
            is_store: orderInfo.is_store,
            customer_status: orderInfo.customer_status,
            // 状态为 3-已取消/4-已过期时，不显示页面底部按钮
            showBottomBtn: orderInfo.status !== 3 && orderInfo.status !== 4
          })
          this.emitUpdateList()
          return
        }
        this.setData({
          showLoading: false
        })
      })
      .catch((err) => {
        console.log(err)
        this.setData({
          showLoading: false
        })
      })
  },
  // 获取服务列表
  getServerList () {
    const params = {
      page: 1,
      size: 9999
    }
    return new Promise((resolve, reject) => {
      serviceListApi(params)
        .then((res) => {
          if (res.code === 200) {
            this.setData({
              serverList: res.data.list
            })
            resolve()
          }
          reject(res.message)
        })
        .catch((err) => {
          reject(err)
        })
    })
  },
  chooseTime () {
    // todo:判断有没有选择服务
    if (!this.data.selectServer.id) {
      wx.showToast('请选择服务类型')
      return
    }
    const { detail } = this.data
    // 如果是进行中并且是自己的邀请单才可以操作时间
    if (detail.status === 0 && detail.is_own) {
      this.setData({
        dateModel: true,
        timeSlot: this.data.selectServer.time_slot
      })
    }
  },
  closeModel () {
    this.setData({
      dateModel: false
    })
  },
  // 提交选择时间
  submit (e) {
    console.log(e)
    this.setData(
      {
        selectTime: e.detail,
        tipModel: e.detail.inviteNum !== 0
      },
      () => {
        this.hasChange()
        this.closeModel()
      }
    )
  },
  changeServer (e) {
    const serverList = this.data.serverList
    this.setData({
      serverIndex: Number(e.detail.value),
      selectServer: serverList[Number(e.detail.value)]
    })
    this.hasChange()
  },
  changeMember (e) {
    const staffLists = this.data.staffLists
    this.setData({
      memberIndex: Number(e.detail.value),
      selectMember: staffLists[Number(e.detail.value)]
    })
    this.hasChange()
  },
  // 是否到店购买 & 选择顾客类型
  onChangeStatus (e) {
    if (e.detail.value) {
      // 更新接口信息
      if (e.currentTarget.dataset.item && e.currentTarget.dataset.type) {
        const { item, type } = e.currentTarget.dataset

        const params = {
          id: item.id,
          reserve_number: item.reserve_number,
          [type]: e.detail.value && Number(e.detail.value)
        }
        changeStatus(params)
          .then((res) => {
            if (res.code === 200) {
              this.setData({
                [type]: e.detail.value
              })
            } else {
              // todo: 提示错误信息
              wx.showToast(res.message)
              if (type === 'is_store') {
                this.setData({
                  is_store: this.data.is_store
                })
              } else {
                this.setData({
                  customer_status: this.data.customer_status
                })
              }
            }
          })
          .catch((err) => {
            console.log(err)
          })
      }
    }
    return false
  },

  /** 页面分享配置 */
  onShareAppMessage (e) {
    const { id, reserveNumber } = this.data
    let title = '邀约评价'
    let url = `pagesOrder/pages/store-invite/comment/comment?id=${id}&reserve_number=${reserveNumber}`
    // 记录邀请评价日志
    const params = {
      id: id, // 主键id
      reserve_number: reserveNumber // 预约单号
    }
    commentShareApi(params)
    return {
      title: title,
      path: url,
      imageUrl: 'https://cdn-cos.coach.com.cn/uploads/2022/5/23/60f86724d69e71c10979b4e5bb17fa23.png'
    }
  },

  // 通知父组件更新状态
  emitUpdateList () {
    try {
      const eventChannel = this.getOpenerEventChannel()
      eventChannel.emit('acceptDataFromOpenedPage', { index: this.data.index, status: this.data.detail.status })
    } catch (e) {
      console.log('getOpenerEventChannel', e)
    }
  },
  /**
   * 图片加载失败
   */
  handleError(e) {
    const { index, goodsIndex } = e.currentTarget.dataset
    this.setData({
      [`detail.reserve_order_sales[${index}].reserve_order_goods[${goodsIndex}].goods_url`]: 'https://cdn-cos.coach.com.cn/uploads/2024/v4.2.0/icon_empty.png'
    })
  }
})
