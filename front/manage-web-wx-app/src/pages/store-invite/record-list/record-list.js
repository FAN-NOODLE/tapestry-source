import {
  orderListApi,
  guideTimeApi,
  getStoreLists,
  getStaffLists,
  commentShareApi,
  verifyShareApi
} from '@/api/guide-invite'
import { memberPositionApi } from '@/api/compass-api'
import { StatusList, sourceList } from '../constants'

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    title: 'record-list',
    list: [],
    selectDate: '',
    statusIndex: -1, // 状态下标
    pageLoading: true,
    isLoading: false, // 分页loading
    page: 1,
    hasMore: true,
    role: 0,
    timeObj: {},
    showTime: false,
    storeLists: [], // 门店列表
    appid: require('../../../lib/config').appid,
    staffLists: [], // 导购列表
    listOperationObj: {}, // 邀约列表子组件参数
    statusList: StatusList,
    storeIndex: -1, // 门店下标
    memberIndex: -1, // 导购下标
    sourceList: sourceList
  },

  // 生命周期函数--监听页面加载
  async onLoad (option) {
    console.log('UUUnion_id: ', option)
    if (option.union_id) {
      this.setData({
        union_id: option.union_id
      })
    }

    // 默认选中状态
    if (option.statusIndex) {
      this.setData({
        statusIndex: this.getCheckedIndex(Number(option.statusIndex), StatusList, 'id')
      })
    }
    const eventChannel = this.getOpenerEventChannel()
    console.log(eventChannel)
    // 监听acceptDataFromOpenerPage事件，获取上一页面通过eventChannel传送到当前页面的数据
    if (option.reserve_numbers) {
      option.reserve_numbers = await new Promise((resolve, reject) => {
        eventChannel.on('accept', function(data) {
          resolve(data)
        })
      })

      this.setData({
        member_id: Number(option.member_id),
        reserve_numbers: option.reserve_numbers
      })
    }

    // 传递下去，默认选中导购, 默认选中时间
    this.init(option)
  },
  onShow () {
    // 获取实例
    this.filterModel = this.selectComponent('#filterModel')
  },

  /** 打开更多筛选弹窗 */
  openMoreFilter () {
    const {statusIndex, statusList, sourceIndex, sourceList} = this.data
    const params = {status: -1, order_origin: 0}
    // 添加状态id
    if (statusIndex > -1) {
      params.status = statusList[Number(statusIndex)].id
    }
    // 添加来源id
    if (sourceIndex > -1) {
      params.order_origin = sourceList[Number(sourceIndex)].id
    }
    this.filterModel.updatePopup(true, params)
  },
  reset (pullDown = false) {
    this.setData(
      {
        pageLoading: !pullDown,
        list: [],
        page: 1
      },
      () => {
        this.getList(pullDown)
      }
    )
  },
  /** 更多筛选提交 */
  filterModelSubmit (e) {
    const {statusList, sourceList} = this.data
    const {status, order_origin} = e.detail

    // 获取index
    const statusIndex = statusList.findIndex(item => item.id === status)
    const sourceIndex = sourceList.findIndex(item => item.id === order_origin)

    this.setData({
      statusIndex,
      sourceIndex
    })

    this.reset()
  },
  searchReset () {
    this.setData({
      pageLoading: true,
      page: 1,
      selectDate: '',
      statusIndex: -1,
      storeIndex: -1,
      memberIndex: -1,
      sourceIndex: -1,
      list: []
    })
    this.getList()
    this.getStaffLists({})
  },

  /**
   * 设置默认选中状态
   * @param optionValue 路由参数
   * @param arr 需要选中的list
   * @param attribute 属性
   */
  getCheckedIndex (optionValue, arr, attribute) {
    return arr.findIndex(item => item[attribute] === optionValue)
  },

  // 获取邀约记录列表
  getList (pullDown = false) {
    const {
      statusIndex,
      statusList,
      page,
      staffLists,
      storeLists,
      selectDate,
      union_id,
      storeIndex,
      memberIndex,
      sourceIndex,
      sourceList
    } = this.data

    const storeId = []
    const membersId = []
    if (memberIndex > -1 && staffLists.length > 0) {
      membersId.push(staffLists[Number(memberIndex)].members_id)
    }
    if (storeIndex > -1 && storeLists.length > 0) {
      storeId.push(storeLists[Number(storeIndex)].store_id)
    }
    const params = {
      page,
      member_ids: membersId,
      store_ids: storeId,
      size: 10
    }
    // 从我的页面进入有次参数
    if (union_id) {
      params.unionid = union_id
    }

    if (selectDate) {
      const date = selectDate.replace(/-/g, '/')
      params.start_time = new Date(`${date} 00:00:00`).getTime() / 1000
    }
    if (statusIndex > -1) {
      params.status = statusList[Number(statusIndex)].id
    }
    if (sourceIndex > -1) {
      params.order_origin = sourceList[Number(sourceIndex)].id
    }
    if (this.data.reserve_numbers) {
      params.member_ids.push(this.data.member_id)
      params.reserve_numbers = this.data.reserve_numbers
    }
    console.log(
      'params',
      selectDate,
      statusIndex,
      memberIndex,
      storeIndex,
      params
    )
    this.setData({
      isLoading: true
    })
    orderListApi(params)
      .then((res) => {
        pullDown && wx.stopPullDownRefresh()
        if (res.code !== 200) {
          this.setData({
            pageLoading: false,
            isLoading: false
          })
          if (this.data.page > 1) this.data.page -= 1
          return
        }
        let { total, page, size } = res.data
        let list = res.data.list || []
        // todo:判断时间是否在一个月之内
        list = list.map((item) => {
          item.isOneMonth = this.isOneMonth(item.updated_at * 1000)
          return item
        })
        console.log('list', list)
        this.setData({
          list: [...this.data.list, ...list],
          hasMore: list.length === size,
          pageLoading: false,
          isLoading: false,
          listOperationObj: {
            role: this.data.role
          }
        })
      })
      .catch((err) => {
        wx.showToast(err.message)
        this.setData({
          pageLoading: false,
          isLoading: false
        })
        if (this.data.page > 1) this.data.page -= 1
        pullDown && wx.stopPullDownRefresh()
      })
  },
  // todo:判断时间是否在一个月之内
  isOneMonth (time) {
    const now = new Date().getTime()
    const oneMonth = 30 * 24 * 60 * 60 * 1000
    return now - time < oneMonth
  },
  // 获取门店列表
  getStoreLists () {
    return new Promise((resolve, reject) => {
      getStoreLists()
        .then((res) => {
          if (res.code === 200) {
            if (JSON.stringify(res.data) === '{}') {
              this.setData({
                storeLists: []
              })
            } else {
              this.setData({
                storeLists: res.data
              })
            }

            resolve()
          }
          reject(res.message)
        })
        .catch((err) => {
          reject(err)
        })
    })
  },
  // 获取导购列表
  getStaffLists (option) {
    return new Promise((resolve, reject) => {
      getStaffLists({
        store_id: option.store_id
      })
        .then((res) => {
          if (res.code === 200) {
            if (JSON.stringify(res.data) === '{}') {
              this.setData({
                staffLists: []
              })
            } else {
              this.setData({
                staffLists: res.data,
                memberIndex: this.getCheckedIndex(Number(option.member_id), res.data, 'members_id')
              })
            }

            resolve()
          }
          reject(res.message)
        })
        .catch((err) => {
          reject(err)
        })
    })
  },
  // 获取时间列表
  getTime (option) {
    return new Promise((resolve, reject) => {
      guideTimeApi()
        .then((res) => {
          if (res.code === 200) {
            this.setData({
              timeObj: res.data,
              showTime: true,
              selectDate: option.time
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
  init (option) {
    let p = [
      this.getStaffLists(option), // 获取导购列表
      this.getStoreLists(), // 获取门店列表
      this.getTime(option), // 获取一个月之内时间段预约的人数
      this.getMember() // 查询店长身份，子组件需要使用
    ]
    Promise.all(p)
      .then(() => {
        this.setData({
          pageLoading: true,
          page: 1,
          list: []
        })
        this.getList()
      })
      .catch((err) => {
        this.setData({
          pageLoading: false
        })
        console.log(err, '出错')
      })
  },
  // 筛选条件change
  changeKey (e) {
    console.log(e)
    const key = e.currentTarget.dataset.key
    const value = e.detail.value
    this.setData({
      [key]: value
    })
    // 根据门店获取导购数据
    if (key === 'storeIndex') {
      this.getStaffLists({
        store_id: this.data.storeLists[value].store_id
      })
    }
    this.reset()
  },
  // 列表滑动到底部 - 上滑加载更多
  onReachBottom () {
    if (this.data.hasMore) {
      this.data.page += 1
      this.getList()
    }
  },
  /**
   * 下拉刷新
   */
  onPullDownRefresh () {
    console.log('我下拉刷新了')
    this.reset(true) // 请求当前页面数据
  },
  // 获取身份信息
  async getMember () {
    await memberPositionApi({ is_new: 0 })
      .then((res) => {
        if (res.code === 200) {
          // 企业导购: role-3 企业管理: role-5（目前只有3和5）
          this.setData({
            role: res.data.role
          })
        }
      })
      .catch((err) => {
        console.log(err)
      })
  },
  /**
   * 子组件与列表状态不一致，按当前筛选条件，从第一页刷新
   */
  statusChange (e) {
    // this.reset()
    const { list, statusIndex } = this.data
    const { index, item } = e.detail
    // 判断当前是否有状态筛选，有则剔除当前数据，无则修改当前数据
    if (statusIndex === -1) {
      this.setData({
        [`list[${index}].status`]: item.status
      })
    } else {
      list.splice(index, 1)
      this.setData({
        list
      })
    }
  },

  /** 判断是否导购自己的邀请单 */
  verifyShareSend (id, reserve_number) {
    return new Promise(resolve => {
      const params = {
        id,
        reserve_number
      }
      verifyShareApi(params).then(res => {
        if (res.code !== 200) {
          resolve(false)
          return wx.showToast(res.message)
        }
        resolve(true)
      }).catch(e => {
        resolve(false)
        console.error(e, 'verifyShareSend 出错')
      })
    })
  },
  /** 页面分享配置 */
  async onShareAppMessage (options) {
    if (options.from === 'button') {
      const { id, reserve_number, is_comment } = options.target.dataset.item
      console.log(11111, options.target.dataset.item)
      const isOwn = await this.verifyShareSend(id, reserve_number)
      // 判断是否导购自己的单子
      if (!isOwn) return
      if (is_comment === 1) return

      let title = '邀约评价'
      let url = `pagesOrder/pages/store-invite/comment/comment?id=${id}&reserve_number=${reserve_number}`

      // 记录邀请评价日志
      const params = {
        id, // 主键id
        reserve_number // 预约单号
      }
      commentShareApi(params)

      return {
        title: title,
        path: url,
        imageUrl: 'https://cdn-cos.coach.com.cn/uploads/2022/5/23/60f86724d69e71c10979b4e5bb17fa23.png'
      }
    }
  }
})
