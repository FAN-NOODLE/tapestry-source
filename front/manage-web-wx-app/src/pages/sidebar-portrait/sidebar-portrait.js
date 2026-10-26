import { memberPositionApi } from '@/api/compass-api'
/** 侧边栏打开客户画像、社群画像 loading页面 */
Page({
  data: {
    shadeShow: true,
    isError: false
  },
  async onLoad(options) {
    console.log('onLoad', options)
    this.data.options = options
    // 取消顶部返回按钮
    wx.hideHomeButton()
    await this.getMember()
    await this.judgeEntry()
  },
  /**
   * 获取身份信息
   */
  async getMember () {
    await memberPositionApi({ is_new: 0 }).then(res => {
      if (res.code === 200) {
        // role 3 企业导购 5 企业管理（目前只有3和5）
        this.setData({
          role: res.data.role
        })
      }
    }).catch(err => {
      console.log(err)
    })
  },
  /**
   * 判断入口
   */
  judgeEntry () {
    console.log('judgeEntry')
    const app = getApp()
    const { nav } = this.data.options
    const that = this
    wx.qy.getContext({
      success: res => {
        var entry = res.entry // 返回进入小程序的入口类型
        console.log('entry', entry)
        if (entry === 'single_chat_tools' || entry === 'contact_profile') {
          if (app.globalData.is_beta_user === 1) {
            wx.reLaunch({
              url: `/packages/customer/pages/detail/detail?nav=${Number(nav || 0) === 2 ? 4 : 0}`
            })
          } else {
            wx.reLaunch({
              url: `/pages/customer-detail/customer-detail?nav=${nav || 0}`
            })
          }
        } else if (entry === 'group_chat_tools') {
          wx.qy.getCurExternalChat({
            success: res => {
              wx.reLaunch({
                url: `/pagesA/group-portrait/group-portrait?chatId=${res.chatId}`
              })
            },
            fail: err => {
              that.setData({
                isError: true,
                shadeShow: false
              })
              console.error(err, 'getCurExternalContact===')
              throw new Error(err.errMsg)
            }
          })
        } else {
          console.log('wx.qy.getContext==')
        }
        that.setData({
          isError: false,
          shadeShow: false
        })
      },
      fail: function (err) {
        console.error('wx.qy.getContext 异常', err)
        that.setData({
          isError: true,
          shadeShow: false
        })
        throw new Error(err.errMsg)
      }
    })
  }
})
