import { getOnlineOrderPage } from '@/api/order'
import { memberPositionApi } from '@/api/compass-api'

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    url: '', // 调换的h5地址
    loading: true, // 查询h5地址
    isError: false, // 群聊打开时为true,
    /** 路由参数 */
    options: {}
  },

  // 生命周期函数--监听页面加载
  onLoad (options) {
    this.setData({
      options
    })
    wx.hideShareMenu({
      withShareTicket: true,
      menus: ['shareAppMessage', 'shareTimeline']
    })
  },
  onShow () {
    // 重置状态
    this.setData({
      url: '',
      loading: true,
      isError: false
    })
    if (this.data.options.scene === 'list') {
      // 从导购个人中心打开，访问订单列表
      this.getOrderPage(1)
    } else {
      // 防止不会重新登录相当于触发wx.qy.login
      this.getMember().then(() => {
        // 通过聊天工具打开
        this.getUserId()
      }).catch((e) => console.error(e))
    }
  },

  /**
   * 获取身份信息
   */
  async getMember () {
    return new Promise(resolve => {
      memberPositionApi({ is_new: 0 }).then(res => {
        if (res.code === 200) {
          resolve()
        }
      }).catch(err => {
        console.log(err)
      })
    })
  },
  /** 获取用户ID */
  getUserId () {
    const _this = this
    wx.qy.getContext({
      success: (res) => {
        // 返回进入小程序的入口类型
        var entry = res.entry
        if (entry === 'single_chat_tools') {
          wx.qy.getCurExternalContact({
            success: (res) => {
              console.log('外部联系人返回res,=======>', res)
              // 返回当前外部联系人userId
              this.setData({
                user_id: res.userId
              })
              this.getOrderPage(2)
              // 神策数据上报 -- 点击侧边栏
              wx.track.trackClickItem({
                item_name: _this.route,
                user_external_id: res.userId
              })
            },
            fail: (err) => {
              console.error(err, 'getCurExternalContact===')
            }
          })
        } else {
          _this.setData({
            isError: true,
            loading: false
          })
          console.log('wx.qy.getCurExternalContact==')
        }
      },
      fail: function (err) {
        console.error('wx.qy.getContext 异常', err)
        _this.setData({
          isError: true,
          loading: false
        })
      }
    })
  },
  /**
   * 查询线上订单配置
   */
  getOrderPage (type) {
    wx.showNavigationBarLoading()
    let params = {
      type, // 1 在线订单(个人) 2 分享订单(聊天)
      openid: ''
    }
    // 2-聊天工具进入
    if (type === 2) {
      params.openid = this.data.user_id
    }
    getOnlineOrderPage(params)
      .then((res) => {
        this.setData({
          loading: false
        })
        wx.hideNavigationBarLoading()
        // 错误判断
        if (res.code !== 200) return wx.showToast(res.message)
        this.setData({
          url: res.data.url
        })
      })
      .catch((err) => {
        this.setData({
          loading: false
        })
        wx.hideNavigationBarLoading()
        wx.showToast(err.message)
      })
  },
  /**
   * h5 postMessage
   * @return
   */
  messageFn (e) {
    const {dataList} = e.detail
    console.log(dataList)
  }
})
