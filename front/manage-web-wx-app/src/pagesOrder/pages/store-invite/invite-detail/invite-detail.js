// 获取全局应用程序实例对象
import { orderDetailApi, orderEditApi } from '../../../api/guide-invite'
import { toDate } from '@/utils/tools'
import { StatusList } from '../constants'
import { sensorsLoginInit } from '@/lib/sensors'

/** 时间补0 */
let tf = function (i) {
  return (i < 10 ? '0' : '') + i
}

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    id: 0,
    serverIndex: -1, // 服务选择的下标
    serverList: [], // 服务列表
    reserveNumber: '', // 备注
    statusList: StatusList,
    showLoading: true, // 页面loading
    dateModel: false, // 是否可以展示选择时间
    detail: {}, // 邀约单详情
    selectMember: {}, // 店长可选------选择的员工
    selectServer: {}, // 店员可选------选择服务
    memberList: [], // 店长可选------用户的手机号
    ownerChange: false, // 店长是否改变分配的员工
    staffChange: false, // 店员是否改变了数据
    showModel: false, // 店员完成服务弹框
    // canShare: true, // 是否向用户分享评论页面
    goShare: false, // 去分享了
    tipModel: false, // 选择时间是是否占用
    showFinishBtn: false, // 展示核销按钮
    writeOffModel: false, // 展示核销弹框
    orderNumber: '',
    // orderInfo: {},
    isComment: false,
    showBottomBtn: false, // 是否展示底部按钮
    showConfirmModal: false, // 二次确认弹窗
    confirmDesc: '确认',
    confirmDescTwo: '确定提交该邀约单',
    user_remark: '', // 用户备注信息
    serviceScoreList: [], // 活动安排星级
    storeScoreList: [] // 整体服务星级
  },
  // 生命周期函数--监听页面加载
  onLoad (options) {
    this.data.id = Number(options.id)
    this.data.reserveNumber = options.reserve_number
    wx.hideShareMenu({
      withShareTicket: true,
      menus: ['shareAppMessage', 'shareTimeline']
    })
    wx.hideHomeButton()
  },
  onShow () {
    this.getDetail()
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
  // (邀约单确认) -备注信息填写
  changeInput (e) {
    const key = e.currentTarget.dataset.key
    const value = e.detail.value
    this.setData({
      [key]: value
    })
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
  // 获取邀约详情信息
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
          if (res.data.status === 0) {
            wx.showShareMenu({
              withShareTicket: true,
              menus: ['shareAppMessage', 'shareTimeline']
            })
          }

          const date = new Date(orderInfo.reserve_time * 1000)
          // todo: 到店开始结束时间
          let select_hours = ''
          if (orderInfo.reserve_end_time === 0) {
            select_hours = `${date.getHours()}:00 - ${date.getHours() + 1}:00`
          } else {
            const endDate = new Date(orderInfo.reserve_end_time * 1000)
            select_hours = `${date.getHours()}:00 - ${endDate.getHours()}:00`
          }
          const select_name = `${date.getFullYear()}-${
            tf(date.getMonth() + 1)
          }-${tf(date.getDate())}`
          const selectTime = {
            select_name,
            inviteNum: 0,
            isFinish: false,
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
            '¥' + res.data.reserve_order_price / 100

          const { selectServer, serverIndex } = this.setServer(orderInfo)
          this.setData({
            serverIndex,
            selectTime,
            selectServer,
            showFinishBtn: orderInfo.reserve_time * 1000 < new Date().getTime(),
            showLoading: false,
            ownerChange: false,
            staffChange: false,
            serviceScoreList,
            storeScoreList,
            detail: res.data,
            selectMember: {
              mobile: res.data.mobile || '',
              member_name: res.data.member_name
            },
            // 条件： is_comment-1 已评价 & status-1 邀约成功 / status-2 邀约完成
            isComment:
              orderInfo.is_comment === 1 &&
              (orderInfo.status === 1 || orderInfo.status === 2),
            // 状态为 3-已取消/4-已过期时，不显示页面底部按钮
            showBottomBtn: orderInfo.status !== 3 && orderInfo.status !== 4,
            user_remark: orderInfo.user_remark
          })
          try {
            // 通知父组件更新状态
            const eventChannel = this.getOpenerEventChannel()
            eventChannel.emit('acceptDataFromOpenedPage', {
              status: this.data.detail.status
            })
          } catch (e) {
            console.log('getOpenerEventChannel', e)
          }

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
  // 确认邀约单信息
  confirmOrder () {
    this.setData({
      showConfirmModal: true
    })
  },
  // 关闭确认弹窗
  cancelModel () {
    this.setData({
      showConfirmModal: false
    })
  },
  // 确认邀约单
  sureOperation () {
    const { id, reserveNumber, user_remark, detail } = this.data
    const params = {
      id,
      user_remark,
      reserve_number: reserveNumber,
      reserve_time: detail.reserve_time,
      reserve_end_time: detail.reserve_end_time,
      reserve_id: detail.reserve_id,
      status: 1
    }
    wx.showLoading({
      title: '加载中'
    })
    this.cancelModel()
    orderEditApi(params)
      .then((res) => {
        wx.hideLoading()
        if (res.code === 200) {
          wx.showToast('确认成功')
          this.getDetail()
          // 神策数据上报 -- 开始
          sensorsLoginInit(detail.member_user_id)
          // 神策数据上报 -- 邀约被接受
          wx.track.trackAppointmentAccept({
            reserve_number: detail.reserve_number,
            user_external_id: detail.user_external_id,
            reserve_name: detail.serve_name,
            reserve_time: detail.reserve_time,
            time_period: detail.reserve_time,
            store_name: detail.store_name,
            store_code: detail.store_name
          })
          return
        }

        wx.showToast({
          icon: 'none',
          title: res.message
        })
      })
      .catch((err) => {
        wx.hideLoading()
        console.log(err)
      })
  },
  // 评价/查看 页面跳转
  commentService (item) {
    const { id, reserve_number } =
      item.currentTarget && item.currentTarget.dataset.item
    wx.redirectTo({
      url: `/pagesOrder/pages/store-invite/comment/comment?id=${id}&reserve_number=${reserve_number}&notShare=true`
    })
  }
})
