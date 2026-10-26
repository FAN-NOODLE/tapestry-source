/**
 * @description 神策数据上报 -- 自定义事件上报
 * @class
 */
export default class SensorsTrack {
  /**
   * @constructor
   * @param {object} options
   * @param {object} options.sensors 神策实例
   * @param {Object<string, string>} options.keys 字段转换表
   * @param {boolean} options.open 是否开启上报
   */
  constructor (options) {
    this.sensors = options.sensors
    this.keys = options.keys
    this.open = options.open
  }

  /**
   * @description 数据转换
   * @param {Object<string, string>} data
   * @returns
   */
  transform (name, data) {
    const result = {}
    const keys = this.keys[name] || this.keys
    for (let key in data) {
      if (keys[key]) {
        const item = keys[key]
        if (!item || typeof item === 'string') {
          // 如果没有配置，或者配置项是字符串，直接使用传入字段
          result[keys[key]] = data[key]
        } else if (item instanceof Array) {
          // 如果配置项是一个数组
          result[key] = item[data[key]]
        } else if (typeof item === 'function') {
          // 如果配置项是一个方法
          result[key] = item(data[key])
        } else if (typeof item === 'object') {
          // 如果配置项是一个对象
          if (typeof item.value === 'object') {
            // 如果value是Array或者Object
            result[item.key] = item.value[data[key]]
          } else if (typeof item.value === 'function') {
            // 如果value是一个方法
            result[item.key] = item.value(data[key])
          } else {
            // 如果value是其他
            result[item.key] = data[key]
          }
        } else {
          console.log(`配置错误 {${key}}`)
        }
      } else {
        console.log(`未转换成功的属性key为 {${key}}, 请查看对应关系`)
      }
    }
    return result
  }

  /**
   * @description 事件上报，异常处理
   * @param {string} name 上报事件名称
   * @param {Object<string, string>} value 上报内容
   */
  trackCatch (name, value) {
    try {
      const data = this.transform(name, value)
      console.log('神策上报', name, data)
      if (!this.open) return
      this.sensors.track(name, data)
    } catch (error) {
      console.log('上报异常', error)
    }
  }


  /**
   * @typedef ITrackAppointmentSend
   * @property {string} reserve_number 邀约ID eg: 230108214932189
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   * @property {string?} store_name 门店名称 已全局注入，可以不填
   * @property {string?} store_code 门店编码 已全局注入，可以不填
   * @property {string} reserve_name 服务类型 eg: 美甲服务
   * @property {string} reserve_time 到店时间日期 eg: 2015/6/19
   * @property {string} time_period 到店时间时间段 eg: 17:51:21
   *
   * @description 神策数据上报 -- 发送到店邀约
   * @param {ITrackAppointmentSend} data
   * @example
   * wx.track.trackAppointmentSend({
   *  reserve_number: '230108214932189',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  reserve_name: '美甲服务',
   *  reserve_time: '2015/6/19',
   *  time_period: '17:51:21',
   * })
   */
  trackAppointmentSend (data) {
    this.trackCatch('AppointmentSend', data)
  }

  /**
   * @typedef ITrackAppointmentAccept
   * @property {string} reserve_number 邀约ID eg: 230108214932189
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   * @property {string?} store_name 门店名称 已全局注入，可以不填
   * @property {string?} store_code 门店编码 已全局注入，可以不填
   * @property {string} reserve_name 服务类型 eg: 美甲服务
   * @property {string} reserve_time 到店时间日期 eg: 2015/6/19
   * @property {string} time_period 到店时间时间段 eg: 17:51:21
   *
   * @description 神策数据上报 -- 邀约被接受
   * @param {ITrackAppointmentAccept} data
   * @example
   * wx.track.trackAppointmentAccept({
   *  reserve_number: '230108214932189',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  reserve_name: '美甲服务',
   *  reserve_time: '2015/6/19',
   *  time_period: '17:51:21',
   * })
   */
  trackAppointmentAccept (data) {
    this.trackCatch('AppointmentAccept', data)
  }

  /**
   * @typedef ITrackAppointmentComplete
   * @property {string} reserve_number 邀约ID eg: 230108214932189
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   * @property {string?} store_name 门店名称 已全局注入，可以不填
   * @property {string?} store_code 门店编码 已全局注入，可以不填
   * @property {string} reserve_name 服务类型 eg: 美甲服务
   * @property {string} reserve_time 到店时间日期 eg: 2015/6/19
   * @property {string} time_period 到店时间时间段 eg: 17:51:21
   * @property {number} customer_status 顾客类型 1-新客，2-老客 eg: 1
   * @property {string} is_store 是否到店购买 1-是，2-否 eg: 1
   * @property {number} reserve_order_price 订单金额（单位：分） eg: 300000
   *
   * @description 神策数据上报 -- 邀约成功 后端脚本
   * @param {ITrackAppointmentComplete} data
   * @example
   * wx.track.trackAppointmentComplete({
   *  reserve_number: '230108214932189',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  reserve_name: '美甲服务',
   *  reserve_time: '2015/6/19',
   *  time_period: '17:51:21',
   *  customer_status: 1,
   *  is_store: 1,
   *  reserve_order_price: 300000,
   * })
   */
  trackAppointmentComplete (data) {
    this.trackCatch('AppointmentComplete', data)
  }

  /**
   * @typedef ITrackAppointmentModify
   * @property {string} reserve_number 邀约ID eg: 230108214932189
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   * @property {string?} store_name 门店名称 已全局注入，可以不填
   * @property {string?} store_code 门店编码 已全局注入，可以不填
   * @property {string} change_item 修改项 eg: 到店时间，服务类型，到店时间&服务类型
   * @property {string} modify_reserve_name 服务类型 eg: 美甲服务
   * @property {string} modify_reserve_time 到店时间日期 eg: 2015/6/19
   * @property {string} modify_time_period 到店时间时间段 eg: 17:51:21
   *
   * @description 神策数据上报 -- 修改邀约单
   * @param {ITrackAppointmentModify} data
   * @example
   * wx.track.trackAppointmentModify({
   *  reserve_number: '230108214932189',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  change_item: '到店时间，服务类型，到店时间&服务类型',
   *  modify_reserve_name: '美甲服务',
   *  modify_reserve_time: '2015/6/19',
   *  modify_time_period: '17:51:21',
   * })
   */
  trackAppointmentModify (data) {
    this.trackCatch('AppointmentModify', data)
  }

  /**
   * @typedef ITrackAppointmentCancel
   * @property {string} reserve_number 邀约ID eg: 230108214932189
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   * @property {string?} store_name 门店名称 已全局注入，可以不填
   * @property {string?} store_code 门店编码 已全局注入，可以不填
   * @property {string} reserve_name 服务类型 eg: 美甲服务
   * @property {string} reserve_time 到店时间日期 eg: 2015/6/19
   * @property {string} time_period 到店时间时间段 eg: 17:51:21
   * @property {string} cancel_type 门店编码 1-手动 2-匹配 3-自动 eg: 1
   *
   * @description 神策数据上报 -- 取消邀约单
   * @param {ITrackAppointmentCancel} data
   * @example
   * wx.track.trackAppointmentCancel({
   *  reserve_number: '230108214932189',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  reserve_name: '美甲服务',
   *  reserve_time: '2015/6/19',
   *  time_period: '17:51:21',
   *  cancel_type: '1',
   * })
   */
  trackAppointmentCancel (data) {
    this.trackCatch('AppointmentCancel', data)
  }

  /**
   * @typedef ITrackInvitationSend
   * @property {string} invite_record_id 邀请涵ID eg: 230108214932189
   * @property {string} invite_template_id 邀请涵模版ID eg: 230108214932189
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   * @property {string} reserve_name 服务类型 eg: 美甲服务
   * @property {string} entrance 入口 企微侧边栏 eg: 企微侧边栏
   * @property {string?} store_name 门店名称 已全局注入，可以不填
   * @property {string?} store_code 门店编码 已全局注入，可以不填
   *
   * @description 神策数据上报 -- 发送邀请函
   * @param {ITrackInvitationSend} data
   * @example
   * wx.track.trackInvitationSend({
   *  invite_record_id: '230108214932189',
   *  invite_template_id: '230108214932189',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  reserve_name: '美甲服务',
   *  entrance: '企微侧边栏',
   * })
   */
  trackInvitationSend (data) {
    this.trackCatch('InvitationSend', data)
  }

  /**
   * @typedef ITrackInvitationAccept
   * @property {string} invite_record_id 邀请涵ID eg: 230108214932189
   * @property {string} invite_template_id 邀请涵模版ID eg: 230108214932189
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   * @property {string} reserve_name 服务类型 eg: 美甲服务
   * @property {string} reserve_time 到店时间日期 eg: 2015/6/19
   * @property {string} time_period 到店时间时间段 eg: 17:51:21
   * @property {string?} store_name 门店名称 已全局注入，可以不填
   * @property {string?} store_code 门店编码 已全局注入，可以不填
   *
   * @description 神策数据上报 -- 邀请函被接受
   * @param {ITrackInvitationAccept} data
   * @example
   * wx.track.trackInvitationAccept({
   *  invite_record_id: '230108214932189',
   *  invite_template_id: '230108214932189',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  reserve_name: '美甲服务',
   *  reserve_time: '2015/6/19',
   *  time_period: '17:51:21',
   * })
   */
  trackInvitationAccept (data) {
    this.trackCatch('InvitationAccept', data)
  }

  /**
   * @typedef ITrackTaskReceive
   * @property {string} task_id 任务ID eg: 123，321
   * @property {string} task_name 任务名称 eg: 任务1，任务2
   * @property {number} customer_total 任务覆盖顾客数 eg: 100
   * @property {string} start_time 任务发布时间 eg: 12月1日，12月2日
   *
   * @description 神策数据上报 -- 收到任务 后端
   * @param {ITrackTaskReceive} data
   * @example
   * wx.track.trackTaskReceive({
   *  task_id: '123',
   *  task_name: '任务2',
   *  customer_total: 100,
   *  start_time: '12月1日',
   * })
   */
  trackTaskReceive (data) {
    this.trackCatch('TaskReceive', data)
  }

  /**
   * @typedef ITrackTaskView
   * @property {string} task_id 任务ID eg: 123，321
   * @property {string} task_status 任务状态 eg: 0进行中，1完成，2未完成
   * @property {number} entry_point 任务入口 eg: 1任务中心主任务列表，0应用消息通知
   * @property {number} task_type 任务类型 eg: 0朋友圈任务、1沟通任务
   *
   * @description 神策数据上报 -- 查看任务
   * @param {ITrackTaskView} data
   * @example
   * wx.track.trackTaskView({
   *  task_id: '123',
   *  task_status: '已完成',
   *  entry_point: 0,
   *  task_type: 0,
   * })
   */
  trackTaskView (data) {
    this.trackCatch('TaskView', data)
  }

  /**
   * @typedef ITrackTaskComplete
   * @property {string} task_id 任务ID eg: 123，321
   *
   * @description 神策数据上报 -- 完成任务
   * @param {ITrackTaskComplete} data
   * @example
   * wx.track.trackTaskComplete({
   *  task_id: '123',
   * })
   */
  trackTaskComplete (data) {
    this.trackCatch('TaskComplete', data)
  }

  /**
   * @typedef ITrackTaskExpire
   * @property {string} task_id 任务ID eg: 123，321
   *
   * @description 神策数据上报 -- 任务过期 后端
   * @param {ITrackTaskExpire} data
   * @example
   * wx.track.trackTaskExpire({
   *  task_id: '123',
   * })
   */
  trackTaskExpire (data) {
    this.trackCatch('TaskExpire', data)
  }

  /**
   * @typedef ITrackTaskDetailOperation
   * @property {string} task_id 任务ID eg: 123，321
   * @property {string} operation_type 任务ID eg: 复制文本、查看标签顾客名单、查看顾客名单、电话、发消息、写备注、勾选、全选、移入单独沟通、转发给顾客、选择顾客发送、移回群发顾客
   *
   * @description 神策数据上报 -- 沟通子任务详情页操作
   * @param {ITrackTaskDetailOperation} data
   * @example
   * wx.track.trackTaskDetailOperation({
   *  task_id: '123',
   *  operation_type: '复制文本',
   * })
   */
  trackTaskDetailOperation (data) {
    this.trackCatch('TaskDetailOperation', data)
  }

  /**
   * @typedef ITrackCommunicationSend
   * @property {string} task_id 任务ID eg: 123，321
   * @property {number} customer_number 任务ID eg: 1
   *
   * @description 神策数据上报 -- 沟通子任务详情页操作
   * @param {ITrackCommunicationSend} data
   * @example
   * wx.track.trackCommunicationSend({
   *  task_id: '123',
   *  customer_number: 1,
   * })
   */
  trackCommunicationSend (data) {
    this.trackCatch('CommunicationSend', data)
  }

  /**
   * @typedef ITrackAssetReceive
   * @property {string} source_id 素材ID eg: 123，321
   *
   * @description 神策数据上报 -- 收到素材 H5
   * @param {ITrackAssetReceive} data
   * @example
   * wx.track.trackAssetReceive({
   *  source_id: '123',
   * })
   */
  trackAssetReceive (data) {
    this.trackCatch('AssetReceive', data)
  }

  /**
   * @typedef ITrackAssetView
   * @property {string} source_id 素材ID eg: 123，321
   * @property {string} task_id 任务ID eg: 123，321
   *
   * @description 神策数据上报 -- 查看素材
   * @param {ITrackAssetView} data
   * @example
   * wx.track.trackAssetView({
   *  source_id: '123',
   *  task_id: '123',
   * })
   */
  trackAssetView (data) {
    this.trackCatch('AssetView', data)
  }

  /**
   * @typedef ITrackAssetSave
   * @property {string} source_id 素材ID eg: 123，321
   * @property {string} task_id 任务ID eg: 123，321
   *
   * @description 神策数据上报 -- 保存素材
   * @param {ITrackAssetSave} data
   * @example
   * wx.track.trackAssetSave({
   *  source_id: '123',
   *  task_id: '123',
   * })
   */
  trackAssetSave (data) {
    this.trackCatch('AssetSave', data)
  }

  /**
   * @typedef ITrackAssetSend
   * @property {string} asset_id 素材ID eg: 123，321
   * @property {string} customer_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   * @property {string} task_id 任务ID eg: 123，321
   * @property {string} operation_type 操作类型 ['群发', '转发朋友圈', '普通发送']
   */
  trackAssetSend (data) {
    this.trackCatch('AssetSend', data)
  }

  /**
   * @typedef ITrackAssetFavorite
   * @property {string} asset_id 素材ID eg: 123，321
   * @property {string} customer_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   * @property {string} task_id 任务ID eg: 123，321
   * @property {string} operation_type 操作类型 ['群发', '转发朋友圈', '普通发送']
   */
  trackAssetFavorite (data) {
    this.trackCatch('AssetFavorite', data)
  }

  /**
   * @typedef ITrackMyFavoriteEntranceClick
   * @property {string} asset_id 素材ID eg: 123，321
   * @property {string} customer_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   * @property {string} task_id 任务ID eg: 123，321
   * @property {string} operation_type 操作类型 ['群发', '转发朋友圈', '普通发送']
   */
  trackMyFavoriteEntranceClick (data) {
    this.trackCatch('MyFavoriteEntranceClick', data)
  }

  /**
   * @typedef ITrackReminderCreate
   * @property {string} sow_task_id 播种ID eg: 123，321
   * @property {string} sow_task_content 播种内容 eg: 给a客户打电话
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   * @property {string} entrance 入口 eg: 侧边栏，顾客详情
   *
   * @description 神策数据上报 -- 创建播种
   * @param {ITrackReminderCreate} data
   * @example
   * wx.track.trackReminderCreate({
   *  sow_task_id: '123',
   *  sow_task_content: '给a客户打电话',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  entrance: '侧边栏',
   * })
   */
  trackReminderCreate (data) {
    this.trackCatch('ReminderCreate', data)
  }

  /**
   * @typedef ITrackReminderDelete
   * @property {string} sow_task_id 播种ID eg: 123，321
   * @property {string} sow_task_content 播种内容 eg: 给a客户打电话
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   *
   * @description 神策数据上报 -- 删除播种
   * @param {ITrackReminderDelete} data
   * @example
   * wx.track.trackReminderDelete({
   *  sow_task_id: '123',
   *  sow_task_content: '给a客户打电话',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   * })
   */
  trackReminderDelete (data) {
    this.trackCatch('ReminderDelete', data)
  }

  /**
   * @typedef ITrackReminderComplete
   * @property {string} sow_task_id 播种ID eg: 123，321
   * @property {string} sow_task_content 播种内容 eg: 给a客户打电话
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   *
   * @description 神策数据上报 -- 播种完成
   * @param {ITrackReminderComplete} data
   * @example
   * wx.track.trackReminderComplete({
   *  sow_task_id: '123',
   *  sow_task_content: '给a客户打电话',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   * })
   */
  trackReminderComplete (data) {
    this.trackCatch('ReminderComplete', data)
  }

  /**
   * @typedef ITrackViewRanking
   * @property {string} ranking_cycle 榜单周期 eg: 月榜，年榜
   * @property {string} ranking_type 榜单类型 eg: 1明星榜，2邀约榜，3顾客榜，4活跃榜
   * @property {string} time_range 时间段（月份） eg: 20220304-20230203
   *
   * @description 神策数据上报 -- 查看绩效榜单
   * @param {ITrackViewRanking} data
   * @example
   * wx.track.trackViewRanking({
   *  ranking_cycle: '123',
   *  ranking_type: '月榜',
   *  time_range: '20220304-20230203',
   * })
   */
  trackViewRanking (data) {
    this.trackCatch('ViewRanking', data)
  }

  /**
   * @typedef ITrackClickItem
   * @property {string} item_name 功能 eg: 直播，客户详情，素材库
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   *
   * @description 神策数据上报 -- 点击侧边栏
   * @param {ITrackClickItem} data
   * @example
   * wx.track.trackClickItem({
   *  item_name: 'pages/index/index',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   * })
   */
  trackClickItem (data) {
    this.trackCatch('ClickItem', data)
  }

  /**
   * @typedef ITrackInquireMember
   * @property {Object<string,string>} query_condition 查询条件 eg: {}
   *
   * @description 神策数据上报 -- 会员查询
   * @param {ITrackInquireMember} data
   * @example
   * wx.track.trackInquireMember({
   *  entry_type:'',
   *  query_condition: {
   *    page: 1,
   *    size: 10
   *  },
   * })
   */
  trackInquireMember (data) {
    this.trackCatch('InquireMember', data)
  }

  /**
   * @typedef ITrackInquireMemberResult
   * @property {number} total 结果个数 eg: 100
   *
   * @description 神策数据上报 -- 会员查询结果
   * @param {ITrackInquireMemberResult} data
   * @example
   * wx.track.trackInquireMemberResult({
   *  entry_type:'',
   *  total: 100,
   * })
   */
  trackInquireMemberResult (data) {
    this.trackCatch('InquireMemberResult', data)
  }

  /**
   * @typedef ITrackAssignCustomer
   * @property {string} inherit_id 分配ID eg: 1，2
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   * @property {string} handover_userid 分配前导购ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   * @property {string} takeover_userid 分配后导购ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   *
   * @description 神策数据上报 -- 分配顾客（SM）后端
   * @param {ITrackAssignCustomer} data
   * @example
   * wx.track.trackAssignCustomer({
   *  inherit_id: 100,
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  handover_userid: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  takeover_userid: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   * })
   */
  trackAssignCustomer (data) {
    this.trackCatch('AssignCustomer', data)
  }

  /**
   * @typedef ITrackCustomerAssigned
   * @property {string} inherit_id 分配ID eg: 1，2
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   * @property {string} takeover_userid 分配后导购ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   *
   * @description 神策数据上报 -- 顾客被分配（SA）后端
   * @param {ITrackCustomerAssigned} data
   * @example
   * wx.track.trackCustomerAssigned({
   *  inherit_id: 100,
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  takeover_userid: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   * })
   */
  trackCustomerAssigned (data) {
    this.trackCatch('CustomerAssigned', data)
  }

  /**
   * @typedef ITrackViewAppointmentHistory
   * @property {string} entrance 入口 eg: 点击到店邀约记录，点击【查看预约记录】，点击提醒
   *
   * @description 神策数据上报 -- 查看邀约记录
   * @param {ITrackViewAppointmentHistory} data
   * @example
   * wx.track.trackViewAppointmentHistory({
   *  entrance: 100,
   * })
   */
  trackViewAppointmentHistory (data) {
    this.trackCatch('ViewAppointmentHistory', data)
  }

  /**
   * @typedef ITrackViewCustomerList
   * @property {string} entrance 入口 eg: 点击到店邀约记录，点击【查看预约记录】，点击提醒
   *
   * @description 神策数据上报 -- 查看顾客列表
   * @param {ITrackViewCustomerList} data
   * @example
   * wx.track.trackViewCustomerList({
   *  entrance: 100,
   * })
   */
  trackViewCustomerList (data) {
    this.trackCatch('ViewCustomerList', data)
  }

  /**
   * @typedef ITrackFilterCustomer
   * @property {?string} manage_type 顾客筛选 eg: 0-全部 1-专属 2-非专属
   * @property {?Array} member_list 员工筛选 eg: xx,YY
   * @property {?string} card_type 会员卡等 eg: 注册会员
   * @property {?Array} impression_tag_list 印象标签 eg: xx,YY
   * @property {?Array} tag_list 客户标签 eg: xx,YY
   * @property {?string} clt_type CLT类型 eg: top client
   *
   * @description 神策数据上报 -- 筛选顾客 弃用
   * @param {ITrackFilterCustomer} data
   * @example
   * wx.track.trackFilterCustomer({
   *  manage_type: '全部',
   *  member_list: ['xxx'],
   *  card_type: '注册会员',
   *  impression_tag_list: ['xxx'],
   *  tag_list: ['xxx'],
   *  clt_type: 'xxxx',
   * })
   */
  trackFilterCustomer (data) {
    this.trackCatch('FilterCustomer', data)
  }

  /**
   * @typedef ITrackSearchCustomer
   * @property {string} original_word 搜索原始词 eg: ls
   * @property {number} total 搜索结果数 eg: 20
   *
   * @description 神策数据上报 -- 搜索顾客 弃用
   * @param {ITrackSearchCustomer} data
   * @example
   * wx.track.trackSearchCustomer({
   *  original_word: '全部',
   *  total: 20,
   * })
   */
  trackSearchCustomer (data) {
    this.trackCatch('SearchCustomer', data)
  }

  /**
   * @typedef ITrackCallCustomer
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   *
   * @description 神策数据上报 -- 拨打顾客电话
   * @param {ITrackCallCustomer} data
   * @example
   * wx.track.trackCallCustomer({
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   * })
   */
  trackCallCustomer (data) {
    this.trackCatch('CallCustomer', data)
  }

  /**
   * @typedef ITrackViewLabel
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   *
   * @description 神策数据上报 -- 查看标签
   * @param {ITrackViewLabel} data
   * @example
   * wx.track.trackViewLabel({
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   * })
   */
  trackViewLabel (data) {
    this.trackCatch('ViewLabel', data)
  }

  /**
   * @typedef ITrackAppointmentReview
   * @property {string} reserve_number 邀约ID eg: 230108214932189
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   * @property {string} evaluation_result 评价结果详情 eg: 评价结果详情
   * @property {string} activity_star 活动安排星级 eg: 1，2
   * @property {string} service_star 整体服务星级 eg: 1，2
   *
   * @description 神策数据上报 -- 邀约被评价
   * @param {ITrackAppointmentReview} data
   * @example
   * wx.track.trackAppointmentReview({
   *  reserve_number: '230108214932189',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  evaluation_result: '评价结果详情',
   *  activity_star: '5',
   *  service_star: '5',
   * })
   */
  trackAppointmentReview (data) {
    this.trackCatch('AppointmentReview', data)
  }

  /**
   * @typedef ITrackSignIn
   * @property {string} reserve_number 邀约ID eg: 230108214932189
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   * @property {string?} store_name 门店名称 已全局注入，可以不填
   * @property {string?} store_code 门店编码 已全局注入，可以不填
   * @property {string} reserve_name 服务类型 eg: 美甲服务
   * @property {string} reserve_time 到店时间日期 eg: 2015/6/19
   * @property {string} time_period 到店时间时间段 eg: 17:51:21
   *
   * @description 神策数据上报 -- 顾客到店签到
   * @param {ITrackSignIn} data
   * @example
   * wx.track.trackSignIn({
   *  reserve_number: '230108214932189',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  reserve_name: '美甲服务',
   *  reserve_time: '2015/6/19',
   *  time_period: '17:51:21',
   * })
   */
  trackSignIn (data) {
    this.trackCatch('SignIn', data)
  }

  /**
   * @typedef ITrackFilte
   * @property {string} current_page 当前页面
   * @property {string} filte_type 筛选类型
   * @property {string?} filte_option 筛选条件
   *
   * @description 神策数据上报 -- 顾客列表筛选
   * @param {ITrackFilte} data
   * @example
   * wx.track.trackFilter({
   *  current_page: '230108214932189',
   *  filter_type: '筛选类型',
   *  filter_option: '美甲服务',
   * })
   */
  trackFilter (data) {
    this.trackCatch('Filter', data)
  }

  /**
   * @typedef ITrackSearch
   * @property {string} current_page 当前页面
   * @property {string} search_type 搜索类型
   * @property {string} search_key_word 搜索原始词
   * @property {number} result_number 搜索结果个数 弃用
   *
   * @description 神策数据上报 -- 顾客列表搜索
   * @param {ITrackSearch} data
   * @example
   * wx.track.trackSearch({
   *  current_page: '230108214932189',
   *  search_type: '搜索类型',
   *  search_key_word: '搜索原始词',
   * })
   */
  trackSearch (data) {
    this.trackCatch('Search', data)
  }

  /**
   * @typedef ITrackClickDetails
   * @property {string} user_external_id 顾客ID
   * @property {string} operation_type 操作类型
   *
   * @description 神策数据上报 -- 顾客详情页操作
   * @param {ITrackClickDetails} data
   * @example
   * wx.track.trackClickDetails({
   *  user_external_id: '230108214932189',
   *  operation_type: '操作类型',
   * })
   */
  trackClickDetails (data) {
    this.trackCatch('ClickDetails', data)
  }

  /**
   * @typedef ITrackAssetUpload
   * @property {string} asset_id 素材ID
   *
   * @description 神策数据上报 -- 上传素材
   * @param {ITrackAssetUpload} data
   * @example
   * wx.track.trackAssetUpload({
   *  asset_id: '230108214932189',
   * })
   */
  trackAssetUpload (data) {
    this.trackCatch('AssetUpload', data)
  }
  /**
   * @description 神策数据上报 -- 点击创建群发
   * @property {bool} is_filter 是否筛选
   * @property {string} entry_point 入口
   * @property {string} filter_option 筛选选项
   */
  trackBatchMsgCreate (data) {
    this.trackCatch('BatchMsgCreate', data)
  }
  /**
   * @description 神策数据上报 -- 点击顾客群发页操作
   * @property {string} operation_type 操作类型
   */
  trackBatchMsgOperation (data) {
    this.trackCatch('BatchMsgOperation', data)
  }
  /**
   * @description 神策数据上报 -- 点击顾客群发页操作
   * @property {string} operation_type 操作类型
   */
  trackBatchMsgAssetSelect (data) {
    this.trackCatch('BatchMsgAssetSelect', data)
  }
  /**
   * @description 神策数据上报 -- 群发素材
   * @property {number} asset_coverage 素材覆盖人数
   * @property {array} asset_list 素材列表
   * @property {array} customer_list 顾客列表
   * @property {string} operation_type 操作类型
   */
  trackBatchMsgConfirm (data) {
    this.trackCatch('BatchMsgConfirm', data)
  }
  /**
   * @description 神策数据上报 -- 查看订单记录
   * @property {string} entry_point 入口
   */
  trackViewOrderHistory (data) {
    this.trackCatch('ViewOrderHistory', data)
  }
  /**
   * @description 神策数据上报 -- 进入素材库 +H5
   * @property {string} entry_point 入口
   * @property {string} library_type 素材库类型
   */
  trackAssetLibrary (data) {
    this.trackCatch('AssetLibrary', data)
  }
  /**
   * @description 神策数据上报 -- 查询官网商品素材 +H5
   * @property {string} inquiry_type 查询方式
   */
  trackAssetInquiry (data) {
    this.trackCatch('AssetInquiry', data)
  }
  /**
   * @description 神策数据上报 -- 库存查询 +H5
   * @property {string} entry_point 入口
   * @property {string} inquiry_type 查询方式
   */
  trackStockInquiry (data) {
    this.trackCatch('StockInquiry', data)
  }
  /**
   * @description 神策数据上报 -- 跨店库存查询
   * @property {string} product_id 商品ID
   * @property {string} product_name 商品名称
   */
  trackOtherStockInquiry (data) {
    this.trackCatch('OtherStockInquiry', data)
  }
  /**
   * @description 神策数据上报 -- 任务列表页浏览
   */
  trackViewTaskList (data) {
    this.trackCatch('ViewTaskList', data)
  }
  /**
   * @description 神策数据上报 -- 沟通名单页浏览
   */
  trackViewCommunicate (data) {
    this.trackCatch('ViewCommunicate', data)
  }
  /**
   * @description 神策数据上报 -- 打包名单浮层显示
   */
  trackViewPackagePage (data) {
    this.trackCatch('ViewPackagePage', data)
  }
  /**
   * @description 神策数据上报 -- 名单配置详情弹窗
   */
  trackViewConfigureDetail (data) {
    this.trackCatch('ViewConfigureDetail', data)
  }
  /**
   * @description 神策数据上报 -- 顾客列表页操作
   * @property {string} current_page 当前页面
   * @property {string} operation_type 操作类型
   */
  trackCustomerOperation (data) {
    this.trackCatch('CustomerOperation', data)
  }
  /**
   * @description 神策数据上报 -- 沟通名单页操作
   * @property {string} current_page 当前页面
   * @property {string} operation_type 操作类型
   */
  trackCommunicateOperation (data) {
    this.trackCatch('CommunicateOperation', data)
  }
  /**
   * @description 神策数据上报 -- 任务列表页操作
   * @property {string} current_page 当前页面
   * @property {string} operation_type 操作类型
   */
  trackTaskListOperation (data) {
    this.trackCatch('TaskListOperation', data)
  }
  /**
   * @description 神策数据上报 -- 顾客列表筛选详情
   * @property {string} current_page 当前页面
   * @property {string} button_name 按钮名称
   * @property {string} filter_name 筛选类型
   * @property {string} tag_level_1 标签分类1
   * @property {string} tag_level_2 标签分类2
   * @property {string} tag_value 标签值
   */
  trackFilterDetails (data) {
    this.trackCatch('FilterDetails', data)
  }
  /**
   * @description 神策数据上报 -- 顾客商品推荐商品点击
   * @property {string} customer_id 顾客ID
   * @property {string} commodity_name 商品名称
   * @property {string} product_code 商品code
   * @property {string} product_UPC 商品upccode
   * @property {string} mkt_rank 排序
   * @property {string} mkt_id 图片链接
   */
  trackCommodityClick (data) {
    this.trackCatch('CommodityClick', data)
  }
}
