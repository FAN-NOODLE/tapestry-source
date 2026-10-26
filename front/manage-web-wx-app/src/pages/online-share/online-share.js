// pages/online-share/online-share.js
Page({
  /**
   * 页面的初始数据
   */
  data: {
    fail: '',
    popShow: false,
    /** H5消息页面url */
    link: '',
    /** H5消息标题 */
    title: '',
    /** H5消息摘要 */
    desc: '',
    /** H5消息封面图片URL */
    imgUrl: ''
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad (options) {
    let optionString = options.options ? JSON.parse(options.options) : {}
    /**
     *   参数名 必选 类型 说明
     *   type 是 string 分享方式(Wechat/WeClient)
     *   title 是 string 分享标题 建议 14 个字(标点符号,字母,汉字都算一个字)
     *   smallcover 是 string 封面 限制内容大小不超过 32KB,200*200
     *   bigcover 是 string 封面 限制内容大小不超过 100KB,700*400
     *   description 是 string 描述 长度不超过 26 个字(标点符号,字母,汉字都算一个字)
     *   url 是 string 心愿清单的 url
     *   orderno 是 string 订单编号
     */
    // 无参数，则阻止分享
    if (!Object.keys(optionString).length) {
      this.setData({
        popShow: true,
        fail: '订单缺少参数，无法发送'
      })
      return
    }
    this.setData({
      link: optionString.url,
      title: optionString.title,
      desc: optionString.description,
      imgUrl: optionString.smallcover
    })
    this.shareH5()
  },
  /**
   * 发送订单
   */
  shareH5 () {
    const {link, desc, title, imgUrl} = this.data
    const _this = this
    wx.qy.sendChatMessage({
      msgtype: 'news', // 消息类型，必填
      enterChat: false, // 为true时表示发送完成之后顺便进入会话，仅移动端3.1.10及以上版本支持该字段
      news: {
        link, // H5消息页面url 必填
        title, // H5消息标题
        desc, // H5消息摘要
        imgUrl // H5消息封面图片URL
      },
      success: function (res) {
        console.log('success', res)
        _this.setData({
          fail: '',
          popShow: true
        })
      },
      fail: function (res) {
        console.log('fail', res)
        // TODO: 捕获取消操作
        // _this.setData({
        //   fail: '发送失败',
        //   popShow: true
        // })
        // 跳回上一页
        wx.navigateBack({
          delta: 1
        })
      }
    })
  }
})
