import {
  getMarketingTaskDetail,
  marketingPlanSetStatus,
  marketingPlanGetVideo,
  reportShareApi,
  searchCustomerListApi,
  sendFinishApi
} from '@/api/plan-task-api'
import Retry from '@/lib/retry'
import { toDate } from '@/utils/tools'
// 步骤数据
import { groupStepData, friendStepData, questionData, friendQuestion } from './step-data'
// 批量保存图片到相册
import { wxSaveAuth, downloadFromUrls } from './download'
import { shareToExternalChat, shareToExternalContact, shareToExternalMoments } from '@/utils/wx.qy'
import { getCosCompressUrl } from '@/utils/upload/utils'

const { globalData } = getApp()
// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    /** 是否完成任务确认弹窗 */
    confirmShow: false,
    /** 子任务详情 */
    taskInfo: {},

    /** 素材内容 1文本；2图片；3语音；4视频；5位置；6链接,7 小程序 */
    contentType: 1,

    /** 小程序素材数据 */
    applets: {},

    /** 图片素材数据 */
    image: {},

    /** 链接素材数据 */
    link: {},

    /** 文本类型数据 */
    text: {},

    /** 视频素材数据 */
    video: {},

    /** 是否隐藏播放层 */
    videoOpcity: 100,

    /** 判断是否需要适配 */
    isIPhoneX: globalData.isIPhoneX,

    /** 是否企微环境 */
    wxwork: globalData.system.environment === 'wxwork',
    // 群发助手弹窗实例
    popupSteps: null,
    // todo：群发助手步骤数据
    stepList: groupStepData,
    // todo：是否禁止页面滚动
    isScroll: false,
    // 任务类型
    task_category: 0,
    showLoading: false,
    showQustionTips: false,
    questionData: questionData,
    friendQuestion: friendQuestion
  },

  /** 页面生命周期 */
  onLoad(option) {
    /** 格式化路由参数，由详情接口统一setData */
    const taskInfo = {
      parentId: Number(option.id),
      parentNum: option.num || 0
    }
    this.setData({
      taskInfo,
      option
    })
  },
  onShow() {
    this.getTaskDetail(this.data.taskInfo)
  },
  onUnload() {
    console.log('页面返回了')
    const pages = getCurrentPages()
    const prevPage = pages[pages.length - 2] // 上一个页面
    // 直接调用上一个页面的setData()方法，把数据存到上一个页面中去
    prevPage &&
      prevPage.setData({
        isBack: true
      })
    const eventChannel = this.getOpenerEventChannel()
    eventChannel.emit('update', this.data.taskInfo)
  },
  // 打开步骤弹窗
  openPopupSteps() {
    this.popupSteps.updatePopup()
  },
  // 禁止页面滚动
  onScroll(e) {
    this.setData({
      isScroll: e.detail
    })
  },
  /** 手动播放视频 */
  videoPlay() {
    const videoPlay = wx.createVideoContext('myVideo')
    videoPlay.play()
    this.setData({
      videoOpcity: 0
    })
  },

  /** 获取视频文件的media_id */
  getVideoMediaId(id) {
    return new Promise((resolve) => {
      marketingPlanGetVideo({ id: Number(id) })
        .then((res) => {
          resolve(res.data.media_id || '')
        })
        .catch((e) => console.error('获取media_id 失败'))
    })
  },

  /** 获取子任务详情 */
  getTaskDetail(taskInfo) {
    this.setData({
      showLoading: true
    })
    getMarketingTaskDetail({
      id: taskInfo.parentId
    })
      .then(async (res) => {
        if (res.code !== 200) return

        const { content_info, condition_info, content_type, task_category, ...rest } =
          res.data
        // 如果是视频文件，重新获取media_id
        if (content_type === 4 && content_info.video.media_id) {
          content_info.video.origin_media_id = content_info.video.media_id
          content_info.video.media_id = await this.getVideoMediaId(
            content_info.video.media_id
          )
        }
        // 群发助手步骤数据,如果是企微群发需要获取图片高度，其他不需要
        if (res.data.group_send_type === 1 && res.data.status_code === 0) {
        // todo：群发助手步骤数据
          if (task_category === 0) {
            this.getImageHeight(groupStepData).then((res) => {
              console.log(res, 'res')
              this.setData({
                stepList: res,
                showLoading: false
              })
              //  弹窗实例
              this.popupSteps = this.selectComponent('#popupSteps')
            }).catch((err) => {
              console.log(err, 'err')
            })
          } else {
            this.getImageHeight(friendStepData).then((res) => {
              console.log(res, 'res')
              this.setData({
                stepList: res,
                showLoading: false
              })
              //  弹窗实例
              this.popupSteps = this.selectComponent('#popupSteps')
            }).catch((err) => {
              console.log(err, 'err')
            })
          }
        } else {
          this.setData({
            showLoading: false
          })
        }
        // 子任务有效期
        rest.task_start_time_text = toDate(rest.task_start_time, 'yyyy-MM-dd hh:mm:ss').nowTime
        rest.task_end_time_text = toDate(rest.task_end_time, 'yyyy-MM-dd hh:mm:ss').nowTime
        // 组合统一参数
        this.setData({
          // 任务类型
          task_category,
          // 素材类型
          contentType: content_type,
          // 素材内容
          ...content_info,
          // 标签
          tags: condition_info.tag || [],
          // 详细信息
          taskInfo: {
            task_category,
            ...condition_info,
            ...taskInfo,
            ...rest
          }
        })

        this.trackAssetView(res.data)
        // 神策数据上报 -- 查看任务
        wx.track.trackTaskView({
          task_id: rest.task_id,
          task_status: rest.status_code,
          entry_point: this.data.option.from || 0,
          task_type: task_category
        })
      })
      .catch((err) => {
        this.setData({
          showLoading: false
        })
        console.error(err)
      })
  },
  // 获取图片高度
  getImageHeight(urlList) {
    // 使用promise.all获取图片高度
    return new Promise((resolve, reject) => {
      let promiseList = []
      urlList.forEach((item) => {
        promiseList.push(this.getImageInfo(item))
      })
      Promise.all(promiseList).then((res) => {
        resolve(res)
      }).catch((err) => {
        reject(err)
      })
    })
  },
  getImageInfo(itemObj) {
    return new Promise((resolve, reject) => {
      wx.getImageInfo({
        src: itemObj.description_image,
        success: (res) => {
          // 根据设备宽度和图片宽度计算图片高度
          const { system } = globalData
          const { windowWidth } = system
          const { width, height } = res
          const imageHeight = (windowWidth * height) / width
          // 设置urlList中的图片高度取整
          let newHeight = Math.ceil(imageHeight)
          if (newHeight > 500) {
            newHeight = Math.ceil(newHeight * 0.92)
          }
          itemObj.height = newHeight
          resolve(itemObj)
        },
        fail: (err) => {
          reject(err)
        }
      })
    })
  },
  /** 神策数据上报 -- 查看素材 */
  trackAssetView (data) {
    const {content_type, content_info, task_id} = data
    const { image, video, applets, link } = content_info

    switch (content_type) {
      case 2:
        // 图片素材
        image.imgList.forEach(item => {
          wx.track.trackAssetView({
            task_id: task_id,
            source_id: item.media_id,
            is_official: 1,
            asset_type_modify: 2,
            official_url: '非官网商品素材',
            official_upc: '非官网商品素材',
            official_product: '非官网商品素材',
            official_name: '非官网商品素材'
          })
        })
        break
      case 4:
        // 视频素材
        wx.track.trackAssetView({
          task_id: task_id,
          source_id: video.origin_media_id,
          is_official: 1,
          asset_type_modify: 5,
          official_url: '非官网商品素材',
          official_upc: '非官网商品素材',
          official_product: '非官网商品素材',
          official_name: '非官网商品素材'
        })
        break
      case 6:
        // 链接素材
        wx.track.trackAssetView({
          task_id: task_id,
          source_id: link.id || 0,
          is_official: 1,
          asset_type_modify: 3,
          official_url: '非官网商品素材',
          official_upc: '非官网商品素材',
          official_product: '非官网商品素材',
          official_name: '非官网商品素材'
        })
        break
      case 7:
        // 小程序素材
        wx.track.trackAssetView({
          task_id: task_id,
          source_id: applets.id,
          is_official: 1,
          asset_type_modify: 4,
          official_url: '非官网商品素材',
          official_upc: '非官网商品素材',
          official_product: '非官网商品素材',
          official_name: '非官网商品素材'
        })
        break

      default:
        break
    }
  },

  /** 打开群聊 */
  goChat(e) {
    //  任务过期不能打开群聊
    const { status_code } = this.data.taskInfo
    if (status_code === 2) return wx.showToast('该任务已过期，不能操作')
    wx.qy.openEnterpriseChat({
      chatId: e.currentTarget.dataset.chat,
      groupName: e.currentTarget.dataset.name,
      success: () => {
        console.log('打开会话成功')
      },
      fail: (err) => {
        wx.showToast('打开会话失败')
        console.error(err.errMsg, err)
      }
    })
  },

  /** 保存图片视频-点击事件 */
  savaImageVideo() {
    // 获取保存到相册权限
    wxSaveAuth().then(res => {
      // 保存逻辑
      this.savaImageList(this.data.contentType === 2 ? 'img' : 'video')
    }).catch(err => {
      console.log(err)
    })
  },

  /** 保存素材-辅助函数 */
  savaImageList(type) {
    let self = this
    if (type === 'img') {
      let imgSrc = self.data.image.imgList
      downloadFromUrls(imgSrc, 'image', self.data)
    } else if (type === 'video') {
      let video = [{
        pic_url: self.data.video.video_url
      }]
      downloadFromUrls(video, 'video', self.data)
    }
  },

  /** old 发送给顾客-点击事件 */
  async sendClient() {
    const { taskInfo } = this.data
    try {
      wx.showLoading({
        title: '加载中',
        mask: true
      })
      //  收集素材id
      const {assetIds, message, types} = this.createAttachments()
      // task_category 0 沟通任务 1 朋友圈任务
      if (taskInfo.dimension === 1) {
        // 群发助手发送群发消息 - 群
        await shareToExternalChat({
          ...message
        })
        this.reportShare()
        // 神策数据上报 -- 发送素材
        assetIds.forEach(child => {
          wx.track.trackAssetSend({
            asset_id: child,
            task_id: taskInfo.task_id,
            customer_id: '任务群发无法获取',
            operation_type: 0,
            is_official: 1,
            asset_type_modify: types[child],
            official_url: '非官网商品素材',
            official_upc: '非官网商品素材',
            official_product: '非官网商品素材',
            official_name: '非官网商品素材'
          })
        })
      } else {
        // 群发助手发送群发消息 - 顾客
        await shareToExternalContact({
          ...message
        })
        this.reportShare()
        // 神策数据上报 -- 发送素材
        assetIds.forEach(child => {
          wx.track.trackAssetSend({
            asset_id: child,
            task_id: taskInfo.task_id,
            customer_id: '任务群发无法获取',
            operation_type: 3,
            is_official: 1,
            asset_type_modify: types[child],
            official_url: '非官网商品素材',
            official_upc: '非官网商品素材',
            official_product: '非官网商品素材',
            official_name: '非官网商品素材'
          })
        })
        this.handleTrackTaskDetailOperation(9)
      }
    } catch (error) {
      console.log(error)
    } finally {
      wx.hideLoading()
    }
  },
  /** 创建消息内容 */
  createAttachments() {
    const { contentType, image, video, link, applets, text } =
      this.data
    const message = {
      text: {
        content: text.content
      },
      attachments: []
    }
    let assetIds = []
    let types = {}
    switch (contentType) {
      case 2:
        message.attachments = image.imgList.map((item) => {
          return {
            msgtype: 'image',
            image: {
              imgUrl: getCosCompressUrl({
                url: item.pic_url,
                maxWidth: 1080,
                maxHeight: 1440
              })
            }
          }
        })
        //  取出素材id
        assetIds = image.imgList.map((item) => {
          types[item.media_id] = 2
          return item.media_id
        })
        break
      case 4:
        message.attachments = [
          {
            msgtype: 'video',
            video: {
              mediaid: video.media_id
            }
          }
        ]
        //  取出素材id
        assetIds = [video.origin_media_id]
        types[video.media_id] = 5
        break
      case 6:
        message.attachments = [
          {
            msgtype: 'link',
            link: {
              url: link.link_url,
              title: link.title,
              imgUrl: link.cover_url,
              desc: link.description
            }
          }
        ]
        //  取出素材id
        assetIds = [link.id]
        types[link.id] = 3
        break
      case 7:
        message.attachments = [
          {
            msgtype: 'miniprogram',
            miniprogram: {
              appid: applets.app_id,
              title: applets.title,
              imgUrl: applets.cover_url,
              page: applets.path
            }
          }
        ]
        //  取出素材id
        assetIds = [applets.id]
        types[applets.id] = 4
        break
      default:
        break
    }
    return { message, assetIds, types }
  },
  /** new 企微沟通任务 转发顾客 */
  async shareToExternalContact() {
    try {
      // wx.showLoading({
      //   title: '加载中',
      //   mask: true
      // })

      const { id, task_id } = this.data.taskInfo
      const { assetIds, message, types } = this.createAttachments()
      const res = await searchCustomerListApi({
        id: id,
        page: 1,
        size: 1000,
        filterType: 2,
        customer_status: 0,
        search: {
          is_search: 0
        }
      })
      // 获取顾客列表
      const externalUserIds = res.data.list.map(item => item.work_external_user_id)
      await shareToExternalContact({
        ...message,
        externalUserIds
      })
      // 数据更新
      const unionIds = res.data.list.map(item => item.union_id)
      await sendFinishApi({
        task_allot_log_id: id,
        union_ids: unionIds
      })
      // 神策数据上报 -- 群发素材
      // wx.track.trackBatchMsgConfirm({
      //   asset_coverage: externalUserIds.length,
      //   asset_list: assetIds.map(item => item.toString()),
      //   customer_list: externalUserIds,
      //   operation_type: 0
      // })
      // 神策数据上报 -- 发送素材
      externalUserIds.forEach(item => {
        assetIds.forEach(child => {
          wx.track.trackAssetSend({
            asset_id: child,
            task_id: task_id,
            customer_id: item,
            operation_type: 0,
            is_official: 1,
            asset_type_modify: types[child],
            official_url: '非官网商品素材',
            official_upc: '非官网商品素材',
            official_product: '非官网商品素材',
            official_name: '非官网商品素材'
          })
        })
      })
      // 更新页面数据
      this.getTaskDetail(this.data.taskInfo)
    } catch (error) {
      console.log(error)
    } finally {
      wx.hideLoading()
    }
  },
  /** 转发朋友圈 */
  async shareToExternalMoments() {
    try {
      wx.showLoading({
        title: '加载中',
        mask: true
      })
      const { task_id } = this.data.taskInfo
      const { assetIds, message, types } = this.createAttachments()
      await shareToExternalMoments({
        ...message
      })
      this.reportShare()
      // 神策数据上报 -- 发送素材
      assetIds.forEach(child => {
        wx.track.trackAssetSend({
          asset_id: child,
          task_id: task_id,
          customer_id: '朋友圈群发无法获取',
          operation_type: 1,
          is_official: 1,
          asset_type_modify: types[child],
          official_url: '非官网商品素材',
          official_upc: '非官网商品素材',
          official_product: '非官网商品素材',
          official_name: '非官网商品素材'
        })
      })
    } catch (error) {
      console.log(error)
    } finally {
      wx.hideLoading()
    }
  },

  /**
   * 转发给顾客操作 上报
   */
  reportShare() {
    console.log('上报')
    // status_code 0待完成,1已完成,2已过期
    const { status_code, id } = this.data.taskInfo
    if (status_code === 1 || status_code === 2) {
      // 任务已完成 或者 任务过期 都不执行上报
      return
    }
    const retryRequest = new Retry(reportShareApi)
    //  入参待确认
    const params = { id }
    retryRequest
      .request(params)
      .then((res) => {
        // 操作完成后，重新查询子任务状态
        this.getTaskDetail(this.data.taskInfo)
        console.log(res)
      })
      .catch((err) => {
        console.error(err)
      })
  },
  /** 复制话术 */
  copyTrick () {
    let { content } = this.data.text
    if (content) {
      wx.showToast('文本内容已复制')
      wx.setClipboardData({
        data: content
      })
      this.handleTrackTaskDetailOperation(0)
    } else {
      wx.showToast('文本内容为空')
    }
  },

  /** 标记完成 */
  handleSign(e) {
    let self = this
    const { id, condition_id, dimension, status_code, task_id } = this.data.taskInfo
    const { group_chat_id, is_finish } = e.currentTarget.dataset.item

    // 状态不为0，不能标记
    if (is_finish) return
    // 已过期的任务，不能标记
    if (status_code === 2) return wx.showToast('该任务已过期，不能操作')
    wx.showModal({
      content: '是否标记完成',
      cancelText: '否',
      cancelColor: '#999999',
      confirmText: '是',
      confirmColor: '#F52349',
      success(res) {
        if (res.confirm) {
          const params = {
            status: 1,
            id,
            condition_id,
            dimension,
            group_chat_id
          }
          marketingPlanSetStatus(params)
            .then((res) => {
              if (res.code === 200) {
                self.getTaskDetail(self.data.taskInfo)
                // 神策数据上报 -- 完成任务
                // wx.track.trackTaskComplete({
                //   task_id: task_id
                // })
              }
            })
            .catch((err) => console.log(err))
        }
      }
    })
  },
  // 完成弹窗回调
  confirmFn(e) {
    const { detail } = e
    const self = this
    if (detail.type === 'ok') {
      wx.showNavigationBarLoading()
      const { id, condition_id, dimension, condition_type, group_chat_id, task_id } =
        this.data.taskInfo
      const params = {
        status: 1,
        id,
        condition_id,
        dimension,
        group_chat_id: 0
      }
      if (condition_type === 7) {
        params.group_chat_id = group_chat_id - 0
      }
      marketingPlanSetStatus(params)
        .then((res) => {
          wx.hideNavigationBarLoading()
          if (res.code === 200) {
            self.getTaskDetail(self.data.taskInfo)
            // 神策数据上报 -- 完成任务
            // wx.track.trackTaskComplete({
            //   task_id: task_id
            // })
          }
        })
        .catch((err) => {
          wx.hideNavigationBarLoading()
          console.log(err)
        })
      this.setData({
        confirmShow: false
      })
    } else {
      this.setData({
        confirmShow: false
      })
      console.log('用户点击取消')
    }
  },

  /** 任务完成 */
  fulfilTask() {
    const { status_code } = this.data.taskInfo
    if (status_code === 2) return wx.showToast('该任务已过期，不能操作')
    this.setData({
      confirmShow: true
    })
  },

  /** 查看标签客户名单 */
  gotoTaps(e) {
    // 跳转到标签客户名单
    // status(status_code) 0待完成,1已完成,2已过期
    // type 0-标签画像顾客名单;1-退群顾客名单;2-导入顾客名单的
    // category 0 沟通任务 1 朋友圈任务
    const { value } = e.currentTarget.dataset
    const { id, group_send_type, dimension, task_category } = this.data.taskInfo

    if (group_send_type === 2) {
      // 企微沟通任务
      wx.navigateTo({
        url: `/packages/task/pages/task-customer/task-customer?id=${id}&dimension=${dimension}&task_category=${task_category}`
      })
    } else {
      wx.navigateTo({
        url: `/packages/task/pages/task-tags/task-tags?id=${id}&dimension=${dimension}&task_category=${task_category}`
      })
    }
    this.handleTrackTaskDetailOperation(value)
  },

  /** 图片预览 */
  previewImage(e) {
    let index = e.target.dataset.index
    var urls = this.data.image.imgList.map((item) => {
      return item.pic_url
    })
    wx.previewImage(urls[index], urls)
  },
  //  问答弹窗
  openPopupQuestion() {
    this.setData({
      showQustionTips: true
    })
  },
  closeTips() {
    this.setData({
      showQustionTips: false
    })
  },
  handleShowTooltip(e) {
    const { selector } = e.currentTarget.dataset

    let content = '1、任务触达率=目标客户中联系状态为【已联系】顾客数/目标客户总数*100%<br/>2、任务状态更新会存在5～10分钟延迟'
    this.selectComponent('#mp-tooltip').show(this, selector, content)
  },
  /** 数据上报 */
  handleTrackTaskDetailOperation(type) {
    // 神策数据上报 -- 沟通子任务详情页操作
    wx.track.trackTaskDetailOperation({
      task_id: this.data.taskInfo.task_id,
      operation_type: type
    })
  }
})
