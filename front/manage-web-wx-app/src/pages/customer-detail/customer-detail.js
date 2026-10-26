import { getCustomerDetailApi, userDetailTipsNumApi, updateRemarkApi } from '@/api/customer-api'
import { memberPositionApi } from '@/api/compass-api'

const { globalData } = getApp()

const defaultNavList = [
  {
    id: 0,
    label: '会员信息'
  },
  {
    id: 1,
    label: '邀约记录',
    count: 0
  },
  {
    id: 2,
    label: '播种',
    count: 0
  },
  {
    id: 3,
    label: '绑定关系'
  }
]
// 创建页面实例对象
Page({
  data: {
    openid: null,
    userDetail: {}, // 用户信息
    navList: defaultNavList,
    navId: 0,
    role: 3, // role 3 企业导购 5 企业管理（目前只有3和5）
    wxwork: globalData.system.environment === 'wxwork',
    imageHost: 'https://cdn-cos.coach.com.cn/uploads/2022/v3.0.2/weclient/',
    isShowEdit: false, // 是否显示编辑弹窗
    editValue: '',
    isRelation: false,
    /** 判断是否从企微侧边栏打开 */
    sidebar: false,
    /** 是否显示错误提示 */
    isError: false,
    /** 判断是否从任务计划顾客名单 */
    isTask: false,
    /** 用户所在的索引，反向查询用户数据 */
    customerIndex: 0,
    // 顾客id
    customer_id: ''
  },
  onLoad (options) {
    let { openid, memberUserid, memberId, isRelation, nav, isTask, index, customer_id } = options
    // 子任务客户名单页面过来，需要记录客户索引
    if (isTask === '1') {
      index = Number(index)
    }
    this.setData({
      options,
      openid,
      memberUserid,
      memberId: +memberId,
      isRelation: +isRelation === 1,
      navId: +nav || 0,
      isTask: isTask === '1', // 值为1代表从子任务名单跳转来，需要关联子任务操作
      customerIndex: index,
      customer_id: customer_id
    })
    // 从任务计划进来 获取参数
  },
  async onShow () {
    await this.getMember()
    await this.getNum()
  },
  onHide () {
    // 重置openId, 重置页面
    if (this.data.sidebar) {
      this.setData({
        openid: null
      })
    }
    defaultNavList[1].count = 0
    defaultNavList[2].count = 0
  },
  /**
   * 获取用户ID
   */
  getUserId () {
    const _this = this
    return new Promise((resolve, reject) => {
      wx.qy.getContext({
        success: res => {
          var entry = res.entry // 返回进入小程序的入口类型

          // 如果从企微侧边栏播种入口打开,
          if (entry === 'group_chat_tools' && _this.data.navId === 2) {
            _this.setData({ isError: true })
            reject(new Error('从群聊进入'))
          }
          if (entry === 'single_chat_tools' || entry === 'contact_profile') {
            wx.qy.getCurExternalContact({
              success: res => {
                console.log('外部联系人返回res,=======>', res)
                _this.setData({ openid: res.userId, isRelation: true, sidebar: true })// 返回当前外部联系人userId
                resolve(res.userId)
                // 取消顶部返回按钮
                wx.hideHomeButton()
                // 神策数据上报 -- 点击侧边栏
                wx.track.trackClickItem({
                  item_name: _this.route,
                  user_external_id: res.userId
                })
              },
              fail: err => {
                console.error(err, 'getCurExternalContact===')
                reject(err)
              }
            })
          } else {
            console.log('wx.qy.getCurExternalContact==')
            reject(Error)
          }
        },
        fail: function (err) {
          console.error('wx.qy.getContext 异常', err)
          reject(err)
        }
      })
    })
  },
  /**
   * 获取顾客详情
   */
  async getCustomerDetail () {
    const { memberUserid, memberId, openid, isRelation, sidebar, customer_id, options } = this.data
    let params = {}
    // 先判断work_external_user_id 这边是openid
    if (openid) {
      // 判断是否好友关系
      params = isRelation ? {
        external_userid: openid
      } : {
        member_userid: memberUserid,
        member_id: memberId,
        external_userid: openid
      }
      // 判断 customer_id
    } else if (customer_id) {
      // 判断是否好友关系
      params = isRelation ? {
        customer_code: customer_id
      } : {
        member_userid: memberUserid,
        member_id: memberId,
        customer_code: customer_id
      }
      // 如果两个都没有，调用getUserId
    } else {
      let id = ''
      try {
        id = await this.getUserId()
        console.log('侧边栏ID', id)
        this.setData({
          openid: id
        })
      } catch (e) {
        console.log(e)
        return
      }

      params = {
        external_userid: id
      }
    }

    const pages = getCurrentPages()
    // 神策数据上报 -- 查看标签
    wx.track.trackViewLabel({
      customer_id: params.external_userid,
      label_entrance: this.data.sidebar ? (Number(options.nav) ? 'sidebar-1' : 'sidebar') : pages[pages.length - 2] ? pages[pages.length - 2].route : '-1'
    })

    params.is_work = sidebar
    console.log('侧边栏参数', params)
    wx.showNavigationBarLoading()
    try {
      const { code, data, message } = await getCustomerDetailApi(params)
      wx.hideNavigationBarLoading()

      if (code !== 200) {
        wx.showToast(message)
        return
      }

      const tagsList = []

      if (data.user_tags && data.user_tags.length) {
        data.user_tags.forEach(item => {
          if (item.tags_list && item.tags_list.length) {
            item.tags_list.forEach((tag) => {
              tag.isShowPop = false
              tagsList.push(tag)
            })
          }
        })
      }

      data.tagsList = tagsList

      this.setData({ userDetail: {...data, isRelation: this.data.isRelation} })
    } catch (error) {
      console.log(error)
      wx.hideNavigationBarLoading()
    }
  },
  /**
   * 获取顾客邀约和播种数量
   */
  async getNum () {
    try {
      if (!this.data.openid) return
      const { data, code, message } = await userDetailTipsNumApi({ external_userid: this.data.openid })

      if (code !== 200) {
        wx.showToast(message)
        return
      }

      const { reserve_num = 0, sow_task_num = 0 } = data
      const { navList } = this.data

      navList[1].count = reserve_num
      navList[2].count = sow_task_num

      this.setData({ navList })
    } catch (error) {
      console.log(error)
    }
  },
  /**
   * 打开外部联系人对话框
   */
  onOpenExternal (e) {
    if (!this.data.wxwork) {
      wx.showToast('联系功能仅支持企业微信登录时使用，请从企业微信打开WeClient Pro后再试')
      return
    }
    let user_id = this.data.userDetail.external_userid
    console.log(user_id)
    wx.qy.openEnterpriseChat({
      userIds: '',
      externalUserIds: user_id,
      groupName: '',
      success: function (res) {
        // 回调
      },
      fail: function (err) {
        // 失败处理
        console.log(err)
      }
    })
    // 神策数据上报 -- 顾客详情页操作
    this.trackClickDetails(0)
    // 从任务场景查看的客户详情，有拨打电话操作 联动任务计划
    if (this.data.isTask) {
      const eventChannel = this.getOpenerEventChannel()
      eventChannel.emit('acceptDataFromOpenedPage', {
        customerIndex: this.data.customerIndex,
        key: 'is_contact'
      })
    }
  },
  /**
   * 获取身份信息
   */
  async getMember () {
    await memberPositionApi({ is_new: 0 }).then(async res => {
      if (res.code === 200) {
        // role 3 企业导购 5 企业管理（目前只有3和5）
        this.setData({
          role: res.data.role
        })
        // todo:确保首个接口不是调用企微接口
        await this.getCustomerDetail()
      }
    }).catch(err => {
      console.log(err)
    })
  },
  /**
   * 拨打用户电话
   */
  callMember () {
    wx.makePhoneCall({
      phoneNumber: this.data.userDetail.mobile,
      success: () => {
        // 神策数据上报 -- 拨打顾客电话
        wx.track.trackCallCustomer({
          user_external_id: this.data.userDetail.external_userid
        })
        // 神策数据上报 -- 顾客详情页操作
        this.trackClickDetails(12)
      }
    })

    // 从任务场景查看的客户详情，有拨打电话操作 联动任务计划
    if (this.data.isTask) {
      const eventChannel = this.getOpenerEventChannel()
      eventChannel.emit('acceptDataFromOpenedPage', {
        customerIndex: this.data.customerIndex,
        key: 'is_call'
      })
    }
  },
  /**
   * 点击导航栏
   */
  tapNav (e) {
    const {id} = e.currentTarget.dataset
    this.setData({ navId: +id })
    // 神策数据上报 -- 顾客详情页操作
    this.trackClickDetails(+id + 3)
    if (Number(id) === 1) {
      // 神策数据上报 -- 查看邀约记录
      wx.track.trackViewAppointmentHistory({
        entry_point: 7
      })
    }
  },
  /**
   * 打开编辑弹窗
   */
  openEdit () {
    this.setData({ isShowEdit: true })
    // 神策数据上报 -- 顾客详情页操作
    this.trackClickDetails(7)
  },
  /**
   * 标记弹窗输入
   */
  nameInput (e) {
    this.setData({ editValue: e.detail.value })
  },
  /**
   * 触底翻页
   */
  onReachBottom () {
    const { navId } = this.data
    let id = ''

    // 预约列表
    if (navId === 1) {
      id = '#invite-list'
    } else if (navId === 2) {
      id = '#sow-list'
    }

    if (id) {
      this.selectComponent(id).getList(true)
    }
  },
  /**
   * 关闭编辑弹窗
   */
  close () {
    this.setData({ isShowEdit: false })
  },
  /**
   * 确认修改备注名
   */
  async updateRemark () {
    console.log('确认修改备注名')
    const { editValue, openid } = this.data
    if (!editValue.trim()) {
      wx.showToast('请输入备注名')
      return
    }

    const params = {
      external_userid: openid,
      remark: editValue
    }
    try {
      const { data, code, message } = await updateRemarkApi(params)

      if (code !== 200) {
        wx.showToast(message)
        return
      }

      if (data.status === 'OK') {
        wx.showToast('备注修改成功')
        this.getCustomerDetail()
      }
    } catch (error) {
      console.log(error)
    }
    this.setData({ isShowEdit: false })
  },
  /**
     * 去购物清单页面
     */
  goShoppingList () {
    console.log('去购物清单页面')
    const url = `/pages/consume-record/consume-record?unionId=${this.data.userDetail.union_id}&external_userid=${this.data.openid}`
    wx.navigateTo({ url })
    // 神策数据上报 -- 顾客详情页操作
    this.trackClickDetails(1)
  },
  /**
   * 去卡券详情页面
   */
  goCardDetail () {
    console.log('去卡券详情页面')
    const {union_id, membership} = this.data.userDetail
    const params = {
      union_id: union_id || '',
      customer_code: membership || ''
    }
    const url = `/pages/coupon-list/coupon-list?info=${JSON.stringify(params)}`
    wx.navigateTo({ url })
    // 神策数据上报 -- 顾客详情页操作
    this.trackClickDetails(2)
  },
  /**
 * @description: 复制会员卡号
 */
  copyMemberShip () {
    const { membership } = this.data.userDetail
    wx.setClipboardData({
      data: membership,
      success: function (res) {
        wx.showToast('复制成功')
      }
    })
  },
  // 神策数据上报 -- 顾客详情页操作
  trackClickDetails(type) {
    const { external_userid } = this.data.userDetail
    wx.track.trackClickDetails({
      user_external_id: external_userid,
      operation_type: type
    })
  }

})
