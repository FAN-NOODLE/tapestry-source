import {
  getInviteConfigList,
  sendInvite,
  getInvitationInfo,
  recordInfo
} from '@/api/invitation-letter'
import { userInfoApi } from '@/api/guide-invite'
import { throttle } from '../../utils/tools'

const { globalData } = getApp()

Page({
  data: {
    /** 样式 */
    styleIndex: 1,
    canClick: false, // 创建按钮是否可以点击，防止多次点击
    /** 背景下标 */
    currentIndex: 0,
    chinese: ['一', '二', '三', '四', '五', '六', '七', '八', '九', '十'],
    /** 邀请函名称选择 */
    optionsChoose: {},
    /** 选择邀请函弹窗 */
    optionsShow: false,
    /** 是否出错 */
    isError: false,
    /** 是否加载 */
    showLoading: true,

    /** 邀请函模板详情 */
    detail: {},
    /** 是否显示邀请函模板详情 */
    isDetail: false,

    /** 专属顾问 */
    member_name: '--',

    appid: wx.getAccountInfoSync().miniProgram.appId, // appid,
    wxwork: globalData.system.environment === 'wxwork',

    /** 邀请函内容 */
    editValue: '',

    /** 是否隐藏button, 如果是从客户端直接重定向过来的，需要隐藏 */
    hiddenButton: false
  },

  onLoad ({ scene, invite_record_id }) {
    if (scene) {
      this.data.id = decodeURIComponent(scene).split('=')[1]
    }

    this.data.wxwork = wx.getSystemInfoSync().environment === 'wxwork'
    const hiddenButton = !!(this.data.wxwork && invite_record_id)
    if (!this.data.wxwork) {
      this.setData({
        isError: true,
        showLoading: false,
        errorTxt: this.data.id
          ? '请使用企业微信扫码查看'
          : '请通过私聊窗口发送邀请'
      })
      return
    }
    this.setData({
      member_name: globalData.currentEnterprise.account_name || '--',
      isDetail: !!this.data.id,
      hiddenButton,
      invite_record_id
    })
  },
  onHide () {
    // 清空用户身份。防止从安卓进入不刷新
    this.setData({
      userName: '',
      userInfo: {}
    })
  },
  onShow () {
    if (this.data.id) {
      this.getInvitationInfo()
    } else {
      // 存在ID的时候是后台预览，不需要请求顾客信息
      if (!this.data.hiddenButton) {
        this.getUserId()
        this.getInvitation()
      } else {
        // 隐藏返回按钮
        wx.hideHomeButton()
        this.getRecordInfo()
      }
    }
  },

  /** 保存邀请函内容 */
  editorInput (e) {
    const value = e.detail.value
    /** 防止自定义字数变化 */
    if (value.length > 150) return

    this.setData({
      editValue: value
    })
  },

  /** 客户端跳转到导购端重定向 获取邀请函信息 */
  getRecordInfo () {
    recordInfo({invite_record_id: Number(this.data.invite_record_id)})
      .then((res) => {
        if (res.code !== 200) return this.$message.error(res.message)
        let optionsChoose = {
          ...res.data,
          background_url: [res.data.background_url],
          content: res.data.send_content,
          share_url: JSON.parse(res.data.share_url)
        }

        this.setData({
          isError: false,
          showLoading: false,
          currentIndex: 0,
          styleIndex: optionsChoose.style_type,
          options: [optionsChoose],
          optionsChoose,
          editValue: optionsChoose.content
        })
      }).catch(e => console.error('获取接口失败'))
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
              _this.getUserInfo(res.userId)
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
            showLoading: false,
            errorTxt: '请通过私聊窗口发送邀请'
          })
          console.log('wx.qy.getCurExternalContact==', res)
        }
      },
      fail: function (err) {
        console.error('wx.qy.getContext 异常', err)
        console.log(_this, 'xxxxxx')
        // _this.getUserInfo('wmvjp0DgAA-vjnccwO0FzDLkE2WZZvdA')
        _this.setData({
          isError: true,
          showLoading: false,
          errorTxt: '获取用户ID失败'
        })
      }
    })
  },

  /** 获取用户信息 */
  async getUserInfo (external_userid) {
    try {
      const { code, data, message } = await userInfoApi({ external_userid })
      if (code === 200) {
        this.setData({
          userInfo: data,
          userName: data.name || data.nick
        })
      } else {
        wx.showToast(message)
      }
    } catch (error) {
      console.error(error, '获取用户userInfo错误')
    }
  },

  /** 获取邀请函内容 */
  async getInvitation () {
    wx.showLoading({
      title: '加载中'
    })
    try {
      const res = await getInviteConfigList()
      if (res.code === 200) {
        let optionsChoose = res.data[0]
        optionsChoose.share_url = JSON.parse(optionsChoose.share_url)
        optionsChoose.background_url = JSON.parse(optionsChoose.background_url)

        this.setData({
          options: res.data,
          optionsChoose,
          editValue: optionsChoose.content
        })
      } else {
        wx.showToast(res.message)
      }
    } catch (error) {
    } finally {
      this.setData({
        showLoading: false
      })
      wx.hideLoading()
    }
  },

  /** 获取邀请函模板详情 */
  async getInvitationInfo () {
    wx.showLoading({
      title: '加载中'
    })
    try {
      let params = {
        id: Number(this.data.id)
      }
      let res = await getInvitationInfo(params)
      if (res.code === 200) {
        let { data } = res
        try {
          data.share_url = JSON.parse(data.share_url)
          data.background_url = JSON.parse(data.background_url)
        } catch (error) {}
        this.setData({
          detail: data,
          showLoading: false,
          editValue: data.content
        })
      } else {
        wx.showToast(res.message)
      }
    } catch (error) {
      console.log(error)
    } finally {
      wx.hideLoading()
    }
  },

  /** 切换背景图 */
  handleChangePic (e) {
    const { type } = e.currentTarget.dataset
    if (type === '1') {
      if (this.data.currentIndex <= 0) return
      --this.data.currentIndex
    } else {
      if (
        this.data.currentIndex >=
        this.data.optionsChoose.background_url.length - 1
      ) {
        return
      }
      ++this.data.currentIndex
    }
    this.setData({
      currentIndex: this.data.currentIndex
    })
  },

  /** 切换样式 */
  handleChangeStyle (e) {
    const { type } = e.currentTarget.dataset
    this.setData({
      styleIndex: type
    })
  },

  /** 切换邀请函 */
  handleChangeOptions (e) {
    const { item } = e.currentTarget.dataset

    try {
      item.share_url = JSON.parse(item.share_url)
      item.background_url = JSON.parse(item.background_url)
    } catch (error) {}

    this.setData({
      optionsChoose: item,
      editValue: item.content,
      optionsShow: false
    })
  },

  /** 打开切换邀请函下拉 */
  handleOpenOptions () {
    this.setData({
      optionsShow: !this.data.optionsShow
    })
  },

  /** 发送邀请函 */
  sendingInvitation: throttle(function () {
    wx.showLoading({
      title: '加载中'
    })
    const { optionsChoose, userInfo, styleIndex, currentIndex, editValue } = this.data
    const params = {
      invite_template_id: optionsChoose.id,
      user_openid: userInfo.external_userid,
      style_type: styleIndex,
      background_url: optionsChoose.background_url[currentIndex],
      send_content: editValue
    }
    if (this.data.canClick) return
    this.setData({
      canClick: true
    })

    sendInvite(params)
      .then((res) => {
        wx.hideLoading()
        if (res.code === 200) {
          this.sendMessage(res.data.invite_record_id)
        } else {
          wx.showToast(res.message)
          this.setData({
            canClick: false
          })
        }
      })
      .catch((error) => {
        this.setData({
          canClick: false
        })
        console.log(error)
      })
  }, 1500),

  /** 发送 */
  sendMessage (invite_record_id) {
    const { appid, optionsChoose, userInfo } = this.data


    const self = this
    wx.qy.sendChatMessage({
      msgtype: 'miniprogram', // 消息类型，必填
      enterChat: true,
      miniprogram: {
        // 小程序的appid
        appid,
        // 小程序消息的title
        title: optionsChoose.share_title,
        imgUrl: `${optionsChoose.share_url[0]}?imageMogr2/thumbnail/!50p`,
        // 小程序消息打开后的路径，注意要以.html作为后缀，否则在微信端打开会提示找不到页面
        page: `pagesOrder/pages/invite-check/invite-check.html?id=${invite_record_id}`
      },
      success: function (res) {
        console.log(res, '++++')
        self.setData({
          canClick: false
        })
        // 神策数据上报 -- 发送邀请函
        wx.track.trackInvitationSend({
          invite_record_id: invite_record_id,
          invite_template_id: optionsChoose.id,
          user_external_id: userInfo.external_userid,
          entrance: 1
        })
      },
      fail: function (err) {
        console.log(err, '-----error')
        self.setData({
          canClick: false
        })
      }
    })
  }
})
