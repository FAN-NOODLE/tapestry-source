
interface ISensorsTrackOptions {
  /** 是否开启上报 */
  open: boolean;
  /** 神策实例 */
  sensors: any;
  /** 字段转换表 */
  keys: any;
}

/**
 * @description 神策数据上报 -- 自定义事件上报
 * @class
 */
export default class SensorsTrack {
  sensors: any;
  keys: any;
  open: boolean;
  /**
   * @constructor
   * @param {ISensorsTrackOptions} options
   * @param {object} options.sensors 神策实例
   * @param {Object<string, string>} options.keys 字段转换表
   * @param {boolean} options.open 是否开启上报
   */
  constructor (options: ISensorsTrackOptions) {
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
        result[key] = data[key]
      }
    }
    return result
  }

  /**
   * @description 事件上报，异常处理
   * @param {string} name 上报事件名称
   * @param {Object<string, string>} value 上报内容
   */
  trackCatch (name: string, value: any) {
    try {
      const data = this.transform(name,value)
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
   * this.$track.trackAppointmentSend({
   *  reserve_number: '230108214932189',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  reserve_name: '美甲服务',
   *  reserve_time: '2015/6/19',
   *  time_period: '17:51:21',
   * })
   */
  trackAppointmentSend (data: any) {
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
   * this.$track.trackAppointmentAccept({
   *  reserve_number: '230108214932189',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  reserve_name: '美甲服务',
   *  reserve_time: '2015/6/19',
   *  time_period: '17:51:21',
   * })
   */
  trackAppointmentAccept (data: any) {
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
   * @description 神策数据上报 -- 邀约成功
   * @param {ITrackAppointmentComplete} data
   * @example
   * this.$track.trackAppointmentComplete({
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
  trackAppointmentComplete (data: any) {
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
   * this.$track.trackAppointmentModify({
   *  reserve_number: '230108214932189',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  change_item: '到店时间，服务类型，到店时间&服务类型',
   *  modify_reserve_name: '美甲服务',
   *  modify_reserve_time: '2015/6/19',
   *  modify_time_period: '17:51:21',
   * })
   */
  trackAppointmentModify (data: any) {
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
   * this.$track.trackAppointmentCancel({
   *  reserve_number: '230108214932189',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  reserve_name: '美甲服务',
   *  reserve_time: '2015/6/19',
   *  time_period: '17:51:21',
   *  cancel_type: '1',
   * })
   */
  trackAppointmentCancel (data: any) {
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
   * this.$track.trackInvitationSend({
   *  invite_record_id: '230108214932189',
   *  invite_template_id: '230108214932189',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  reserve_name: '美甲服务',
   *  entrance: '企微侧边栏',
   * })
   */
  trackInvitationSend (data: any) {
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
   * this.$track.trackInvitationAccept({
   *  invite_record_id: '230108214932189',
   *  invite_template_id: '230108214932189',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  reserve_name: '美甲服务',
   *  reserve_time: '2015/6/19',
   *  time_period: '17:51:21',
   * })
   */
  trackInvitationAccept (data: any) {
    this.trackCatch('InvitationAccept', data)
  }

  /**
   * @typedef ITrackTaskReceive
   * @property {string} task_id 任务ID eg: 123，321
   * @property {string} task_name 任务名称 eg: 任务1，任务2
   * @property {number} customer_total 任务覆盖顾客数 eg: 100
   * @property {string} start_time 任务发布时间 eg: 12月1日，12月2日
   *
   * @description 神策数据上报 -- 收到任务
   * @param {ITrackTaskReceive} data
   * @example
   * this.$track.trackTaskReceive({
   *  task_id: '123',
   *  task_name: '任务2',
   *  customer_total: 100,
   *  start_time: '12月1日',
   * })
   */
  trackTaskReceive (data: any) {
    this.trackCatch('TaskReceive', data)
  }

  /**
   * @typedef ITrackTaskView
   * @property {string} task_id 任务ID eg: 123，321
   * @property {string} task_status 任务状态 eg: 0进行中，1完成，2未完成
   *
   * @description 神策数据上报 -- 查看任务
   * @param {ITrackTaskView} data
   * @example
   * this.$track.trackTaskView({
   *  task_id: '123',
   *  task_status: '0',
   * })
   */
  trackTaskView (data: any) {
    this.trackCatch('TaskView', data)
  }

  /**
   * @typedef ITrackTaskComplete
   * @property {string} task_id 任务ID eg: 123，321
   *
   * @description 神策数据上报 -- 完成任务
   * @param {ITrackTaskComplete} data
   * @example
   * this.$track.trackTaskComplete({
   *  task_id: '123',
   * })
   */
  trackTaskComplete (data: any) {
    this.trackCatch('TaskComplete', data)
  }

  /**
   * @typedef ITrackTaskExpire
   * @property {string} task_id 任务ID eg: 123，321
   *
   * @description 神策数据上报 -- 任务过期
   * @param {ITrackTaskExpire} data
   * @example
   * this.$track.trackTaskExpire({
   *  task_id: '123',
   * })
   */
  trackTaskExpire (data: any) {
    this.trackCatch('TaskExpire', data)
  }

  /**
   * @typedef ITrackAssetReceive
   * @property {string} source_id 素材ID eg: 123，321
   *
   * @description 神策数据上报 -- 收到素材
   * @param {ITrackAssetReceive} data
   * @example
   * this.$track.trackAssetReceive({
   *  source_id: '123',
   * })
   */
  trackAssetReceive (data: any) {
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
   * this.$track.trackAssetView({
   *  source_id: '123',
   *  task_id: '123',
   * })
   */
  trackAssetView (data: any) {
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
   * this.$track.trackAssetSave({
   *  source_id: '123',
   *  task_id: '123',
   * })
   */
  trackAssetSave (data: any) {
    this.trackCatch('AssetSave', data)
  }

  /**
   * @typedef ITrackAssetSend
   * @property {string} source_id 素材ID eg: 123，321
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   * @property {string} task_id 任务ID eg: 123，321
   *
   * @description 神策数据上报 -- 发送素材
   * @param {ITrackAssetSend} data
   * @example
   * this.$track.trackAssetSend({
   *  source_id: '123',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  task_id: '123',
   * })
   */
  trackAssetSend (data: any) {
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
   * this.$track.trackReminderCreate({
   *  sow_task_id: '123',
   *  sow_task_content: '给a客户打电话',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  entrance: '侧边栏',
   * })
   */
  trackReminderCreate (data: any) {
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
   * this.$track.trackReminderDelete({
   *  sow_task_id: '123',
   *  sow_task_content: '给a客户打电话',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   * })
   */
  trackReminderDelete (data: any) {
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
   * this.$track.trackReminderComplete({
   *  sow_task_id: '123',
   *  sow_task_content: '给a客户打电话',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   * })
   */
  trackReminderComplete (data: any) {
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
   * this.$track.trackViewRanking({
   *  ranking_cycle: '123',
   *  ranking_type: '月榜',
   *  time_range: '20220304-20230203',
   * })
   */
  trackViewRanking (data: any) {
    this.trackCatch('ViewRanking', data)
  }

  /**
   * @typedef ITrackClickItem
   * @property {string} feature 功能 eg: 直播，客户详情，素材库
   * @property {string} customer_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   *
   * @description 神策数据上报 -- 点击侧边栏
   * @param {ITrackClickItem} data
   * @example
   * this.$track.trackClickItem({
   *  feature: 'pages/index/index',
   *  customer_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   * })
   */
  trackClickItem (data: any) {
    this.trackCatch('ClickItem', data)
  }

  /**
   * @typedef ITrackInquireMember
   * @property {Object<string,string>} query_condition 查询条件 eg: {}
   *
   * @description 神策数据上报 -- 会员查询
   * @param {ITrackInquireMember} data
   * @example
   * this.$track.trackInquireMember({
   *  query_condition: {
   *    page: 1,
   *    size: 10
   *  },
   * })
   */
  trackInquireMember (data: any) {
    this.trackCatch('InquireMember', data)
  }

  /**
   * @typedef ITrackInquireMemberResult
   * @property {number} total 结果个数 eg: 100
   *
   * @description 神策数据上报 -- 会员查询结果
   * @param {ITrackInquireMemberResult} data
   * @example
   * this.$track.trackInquireMemberResult({
   *  total: 100,
   * })
   */
  trackInquireMemberResult (data: any) {
    this.trackCatch('InquireMemberResult', data)
  }

  /**
   * @typedef ITrackAssignCustomer
   * @property {string} inherit_id 分配ID eg: 1，2
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   * @property {string} handover_userid 分配前导购ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   * @property {string} takeover_userid 分配后导购ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   *
   * @description 神策数据上报 -- 分配顾客（SM）
   * @param {ITrackAssignCustomer} data
   * @example
   * this.$track.trackAssignCustomer({
   *  inherit_id: 100,
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  handover_userid: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  takeover_userid: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   * })
   */
  trackAssignCustomer (data: any) {
    this.trackCatch('AssignCustomer', data)
  }

  /**
   * @typedef ITrackCustomerAssigned
   * @property {string} inherit_id 分配ID eg: 1，2
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   * @property {string} takeover_userid 分配后导购ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   *
   * @description 神策数据上报 -- 顾客被分配（SA）
   * @param {ITrackCustomerAssigned} data
   * @example
   * this.$track.trackCustomerAssigned({
   *  inherit_id: 100,
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  takeover_userid: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   * })
   */
  trackCustomerAssigned (data: any) {
    this.trackCatch('CustomerAssigned', data)
  }

  /**
   * @typedef ITrackViewAppointmentHistory
   * @property {string} entrance 入口 eg: 点击到店邀约记录，点击【查看预约记录】，点击提醒
   *
   * @description 神策数据上报 -- 查看邀约记录
   * @param {ITrackViewAppointmentHistory} data
   * @example
   * this.$track.trackViewAppointmentHistory({
   *  entrance: 100,
   * })
   */
  trackViewAppointmentHistory (data: any) {
    this.trackCatch('ViewAppointmentHistory', data)
  }

  /**
   * @typedef ITrackViewCustomerList
   * @property {string} entrance 入口 eg: 点击到店邀约记录，点击【查看预约记录】，点击提醒
   *
   * @description 神策数据上报 -- 查看顾客列表
   * @param {ITrackViewCustomerList} data
   * @example
   * this.$track.trackViewCustomerList({
   *  entrance: 100,
   * })
   */
  trackViewCustomerList (data: any) {
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
   * @description 神策数据上报 -- 筛选顾客
   * @param {ITrackFilterCustomer} data
   * @example
   * this.$track.trackFilterCustomer({
   *  manage_type: '全部',
   *  member_list: ['xxx'],
   *  card_type: '注册会员',
   *  impression_tag_list: ['xxx'],
   *  tag_list: ['xxx'],
   *  clt_type: 'xxxx',
   * })
   */
  trackFilterCustomer (data: any) {
    this.trackCatch('FilterCustomer', data)
  }

  /**
   * @typedef ITrackSearchCustomer
   * @property {string} original_word 搜索原始词 eg: ls
   * @property {number} total 搜索结果数 eg: 20
   *
   * @description 神策数据上报 -- 搜索顾客
   * @param {ITrackSearchCustomer} data
   * @example
   * this.$track.trackSearchCustomer({
   *  original_word: '全部',
   *  total: 20,
   * })
   */
  trackSearchCustomer (data: any) {
    this.trackCatch('SearchCustomer', data)
  }

  /**
   * @typedef ITrackCallCustomer
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   *
   * @description 神策数据上报 -- 拨打顾客电话
   * @param {ITrackCallCustomer} data
   * @example
   * this.$track.trackCallCustomer({
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   * })
   */
  trackCallCustomer (data: any) {
    this.trackCatch('CallCustomer', data)
  }

  /**
   * @typedef ITrackViewLabel
   * @property {string} user_external_id 顾客ID eg: wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA
   *
   * @description 神策数据上报 -- 查看标签
   * @param {ITrackViewLabel} data
   * @example
   * this.$track.trackViewLabel({
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   * })
   */
  trackViewLabel (data: any) {
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
   * this.$track.trackAppointmentReview({
   *  reserve_number: '230108214932189',
   *  user_external_id: 'wm9hx_EQAAqi1Jw61IZHlRb4bWwDbGiA',
   *  evaluation_result: '评价结果详情',
   *  activity_star: '5',
   *  service_star: '5',
   * })
   */
  trackAppointmentReview (data: any) {
    this.trackCatch('AppointmentReview', data)
  }

  /**
   * @typedef ITrackAssetUpload
   * @property {string} asset_id 素材ID
   *
   * @description 神策数据上报 -- 上传素材
   * @param {ITrackAssetUpload} data
   * @example
   * this.$track.trackAssetUpload({
   *  asset_id: '230108214932189',
   * })
   */
  trackAssetUpload (data: any) {
    this.trackCatch('AssetUpload', data)
  }
  /**
   * @typedef ITrackEnterTransmittal
   * @property {string} transmittal_entrance 传讯入口
   * 
   * @description 神策数据上报 -- 进入传讯 
   * 1、从企业微信工作台点击图标进入 
   * 2、从消息列表点击进入
   * @param {ITrackEnterTransmittal} data
   * @example
   * this.$track.trackEnterTransmittal({
   *  transmittal_entrance: '企微工作台',
   * })
   */
  trackEnterTransmittal (data: any) {
    this.trackCatch('EnterTransmittal', data)
  }
  /**
   * @typedef ITrackViewTransmittal
   * @property {string} transmittal_title 传讯标题
   * @property {string} transmittal_publisher 传讯发布者
   * @property {string} release_time 发布时间
   * @property {string} transmittal_entrance 入口 
   * FT-我的传讯-查看；
   * RO-传讯审核-去审核；
   * 企业传讯，
   * 企业微信消息
   * 
   * @description 神策数据上报 -- 浏览传讯内容 
   * 1、从企业微信工作台点击图标进入 
   * 2、从消息列表点击进入
   * @param {ITrackViewTransmittal} data
   */
  trackViewTransmittal (data: any) {
    this.trackCatch('ViewTransmittal', data)
  }
  /**
   * @typedef ITrackClickTransmittalPDF
   * 
   * @description 神策数据上报 -- 点击下载周传讯pdf 
   * @param {ITrackClickTransmittalPDF} data
   */
  trackClickTransmittalPDF (data: any) {
    this.trackCatch('ClickTransmittalPDF', data)
  }
  /**
   * @typedef ITrackClickTransmittalDownload
   * 
   * @property {string} transmittal_type 传讯类型
   * 
   * @description 神策数据上报 -- 点击下载传讯文件
   * @param {ITrackClickTransmittalDownload} data
   */
  trackClickTransmittalDownload (data: any) {
    this.trackCatch('ClickTransmittalDownload', data)
  }
  /**
   * @typedef ITrackDownloadTransmittal
   * @property {string} transmittal_title 传讯标题
   * @property {string} transmittal_publisher 传讯发布者
   * @property {string} release_time 发布时间
   * 
   * @description 神策数据上报 -- 下载传讯附件
   * @param {ITrackDownloadTransmittal} data
   */
  trackDownloadTransmittal (data: any) {
    this.trackCatch('DownloadTransmittal', data)
  }
  /**
   * @typedef ITrackWithdrawTransmittal
   * @property {string} transmittal_title 传讯标题
   * @property {string} transmittal_publisher 传讯发布者
   * @property {string} release_time 发布时间
   * 
   * @description 神策数据上报 -- 撤回传讯
   * @param {ITrackWithdrawTransmittal} data
   */
  trackWithdrawTransmittal (data: any) {
    this.trackCatch('WithdrawTransmittal', data)
  }
  /**
   * @typedef ITrackVoidTransmittal
   * @property {string} transmittal_title 传讯标题
   * @property {string} transmittal_publisher 传讯发布者
   * @property {string} release_time 发布时间
   * 
   * @description 神策数据上报 -- 作废传讯
   * @param {ITrackVoidTransmittal} data
   */
  trackVoidTransmittal (data: any) {
    this.trackCatch('VoidTransmittal', data)
  }
  /**
   * @typedef ITrackAuditTransmittal
   * @property {string} transmittal_title 传讯标题
   * @property {string} transmittal_publisher 传讯发布者
   * @property {string} release_time 发布时间
   * @property {string} audit_type 审核类型
   * 
   * @description 神策数据上报 -- 审核传讯
   * @param {ITrackAuditTransmittal} data
   */
  trackAuditTransmittal (data: any) {
    this.trackCatch('AuditTransmittal', data)
  }
  /**
   * @typedef ITrackClickAllTransmittalDownload
   * @property {string} transmittal_type 传讯类型
   * 
   * @description 神策数据上报 -- 一键下载传讯及全部附件
   * @param {ITrackClickAllTransmittalDownload} data
   */
  trackClickAllTransmittalDownload (data: any) {
    this.trackCatch('ClickAllTransmittalDownload', data)
  }
  /**
   * @typedef ITrackClickAllAttachmentDownload
   * @property {string} transmittal_type 传讯类型
   * 
   * @description 神策数据上报 -- 下载传讯及全部附件
   * @param {ITrackClickAllAttachmentDownload} data
   */
  trackClickAllAttachmentDownload (data: any) {
    this.trackCatch('ClickAllAttachmentDownload', data)
  }
  /**
   * @typedef ITrackSearchClick
   * @property {string} tab_name 
   * 
   * @description 神策数据上报 -- 点击搜索框
   * @param {ITrackSearchClick} data
   */
  trackSearchClick (data: any) {
    this.trackCatch('SearchClick', data)
  }
  /**
   * @typedef ITrackFilter
   * @property {string} tab_name
   * @property {string} filter_type_sv
   * @property {string} filter_option
   * 
   * @description 神策数据上报 -- 点击搜索框
   * @param {ITrackFilter} data
   */
  trackFilter (data: any) {
    this.trackCatch('Filter', data)
  }

  /**
   * @typedef ITrackEnterDisplayTask
   * @property {string} vm_entrance 入口
   * 
   * @description 神策数据上报 -- 进入vm陈列任务
   * @param {ITrackEnterDisplayTask} data
   */
  trackEnterDisplayTask (data: any) {
    this.trackCatch('EnterDisplayTask', data)
  }
  /**
   * @typedef ITrackViewStoreDisplay
   * @property {string} dispalytask_title 任务标题
   * @property {string} displaytask_complete_state 完成状态
   * @property {string} displaytask_state 任务状态
   * @property {string} display_store 门店编码
   * 
   * @description 神策数据上报 -- 查看他店陈列图
   * @param {ITrackViewStoreDisplay} data
   */
  trackViewStoreDisplay (data: any) {
    this.trackCatch('ViewStoreDisplay', data)
  }
  /**
   * @typedef ITrackViewDisplayTask
   * @property {string} dispalytask_title 任务标题
   * @property {string} displaytask_complete_state 完成状态
   * @property {string} displaytask_state 任务状态
   * @property {string} display_store 门店编码
   * 
   * @description 神策数据上报 -- 查看陈列任务
   * @param {ITrackViewDisplayTask} data
   */
  trackViewDisplayTask (data: any) {
    this.trackCatch('ViewDisplayTask', data)
  }
  /**
   * @typedef ITrackSubmitDisplayTask
   * @property {string} dispalytask_title 任务标题
   * @property {string} displaytask_complete_state 完成状态
   * @property {string} displaytask_state 任务状态
   * @property {string} is_resubmit 是否重新提交
   * 
   * @description 神策数据上报 -- 提交陈列任务
   * @param {ITrackSubmitDisplayTask} data
   */
  trackSubmitDisplayTask (data: any) {
    this.trackCatch('SubmitDisplayTask', data)
  }
  /**
   * @typedef ITrackCourseLearing
   * @property {string} course_title 课程标题
   * 
   * @description 神策数据上报 -- 课程学习
   * @param {ITrackCourseLearing} data
   */
  trackCourseLearing (data: any) {
    this.trackCatch('CourseLearing', data)
  }
  /**
   * @typedef ITrackDisplayReminder
   * @property {string} reminder_type 提醒类型
   * 
   * @description 神策数据上报 -- 陈列提醒（vm视角）
   * @param {ITrackDisplayReminder} data
   */
  trackDisplayReminder (data: any) {
    this.trackCatch('DisplayReminder', data)
  }
  /**
   * @typedef ITrackMarkPass
   * @property {string} display_store 门店编码
   * @property {string} display_score 陈列评分
   * 
   * @description 神策数据上报 -- 操作陈列审核通过（vm视角）
   * @param {ITrackMarkPass} data
   */
  trackMarkPass (data: any) {
    this.trackCatch('MarkPass', data)
  }
  /**
   * @typedef ITrackMarkNotPass
   * @property {string} display_store 门店编码
   * 
   * @description 神策数据上报 -- 操作陈列审核不通过（vm视角）
   * @param {ITrackMarkNotPass} data
   */
  trackMarkNotPass (data: any) {
    this.trackCatch('MarkNotPass', data)
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
   * @description 神策数据上报 -- 添加商品
   * @typedef ITrackAddItems
   * @property {string} add_method 添加方式
   * @property {string} commodity_name 商品名称
   * @property {string} sku_id SKUID
   * @property {string} customer_id 顾客ID
   */
  trackAddItems (data) {
    this.trackCatch('AddItems', data)
  }

  /**
   * @description 神策数据上报 -- 提交浮窗点击
   * @typedef ITrackSubmitPopupClick
   * @property {string} sales_order_id 销售单号
   * @property {string} has_reject_risk 是否有拒单风险
   * @property {string} risk_type 风险类型
   * @property {string} button_name 按钮名称
   * @property {string} customer_id 顾客ID
   */
  trackSubmitPopupClick (data) {
    this.trackCatch('SubmitPopupClick', data)
  }

  /**
   * @description 神策数据上报 -- 订单生成
   * @typedef ITrackOrderCreated
   * @property {string} sales_order_id 销售单号
   * @property {string} has_reject_risk 是否有拒单风险
   * @property {string} risk_type 风险类型
   * @property {string} customer_id 顾客ID
   */
  trackOrderCreated (data) {
    this.trackCatch('OrderCreated', data)
  }
  /**
   * @description 神策数据上报 -- 订单生成明细
   * @typedef ITrackOrderCreatedDetail
   * @property {string} sales_order_id 销售单号
   * @property {string} commodity_name 商品名称
   * @property {string} sku_id SKUID
   * @property {string} commodity_quantity 商品数量
   * @property {string} inventory_inquiry_status 库存查询状态
   * @property {string} customer_id 顾客ID
   */
  trackOrderCreatedDetail (data) {
    this.trackCatch('OrderCreatedDetail', data)
  }

  /**
   * @description 神策数据上报 -- 发送顾客弹框点击
   * @typedef ITrackToCustomerPopupClick
   * @property {string} sales_order_id 销售单号
   * @property {string} pop_name 弹框名称
   * @property {string} button_name 按钮名称
   * @property {string} customer_id 顾客ID
   */
  trackToCustomerPopupClick (data) {
    this.trackCatch('ToCustomerPopupClick', data)
  }

  /**
   * @description 神策数据上报 -- 订单筛选
   * @typedef ITrackOrderFilter
   * @property {string} current_page 当前页面
   * @property {string} filter_type_sv 筛选类型
   * @property {string} filter_option 筛选条件
   */
  trackOrderFilter (data) {
    this.trackCatch('OrderFilter', data)
  }

  /**
   * @description 神策数据上报 -- tab切换
   * @typedef ITrackTabSwitch
   * @property {string} tab_name tab名称
   */
  trackTabSwitch (data) {
    this.trackCatch('TabSwitch', data)
  }
  
  /**
   * @description 神策数据上报 -- 订单列表点击
   * @typedef ITrackOrderListClick
   * @property {string} sales_order_id 销售单号
   * @property {string} order_status 订单状态
   * @property {string} order_creator 订单创建人
   * @property {string} order_belonged_store 订单所属部门
   * @property {string} customer_id 顾客ID
   */
  trackOrderListClick (data) {
    this.trackCatch('OrderListClick', data)
  }
  
  /**
   * @description 神策数据上报 -- 订单详情点击
   * @typedef ITrackOrderDetailClick
   * @property {string} sales_order_id 销售单号
   * @property {string} order_status 订单状态
   * @property {string} order_creator 订单创建人
   * @property {string} order_belonged_store 订单所属部门
   * @property {string} customer_id 顾客ID
   * @property {string} button_name 按钮名称
   */
  trackOrderDetailClick (data) {
    this.trackCatch('OrderDetailClick', data)
  }
  
  /**
   * @description 神策数据上报 -- 导出当前列表
   * @typedef ITrackExportOrderList
   * @property {string} period 时间段
   * @property {string} store_name 门店
   * @property {string} sales_name 导购
   * @property {string} original_words 搜索原始词
   * @property {string} order_status 订单类型
   */
  trackExportOrderList (data) {
    this.trackCatch('ExportOrderList', data)
  }
  
  /**
   * @description 神策数据上报 -- 进入网仓
   * @typedef ITrackEnterOnlineWarehouse
   * @property {string} OnlineWarehouse_entrance 入口
   * @property {string} message_type 消息类型
   */
  trackEnterOnlineWarehouse (data) {
    this.trackCatch('EnterOnlineWarehouse', data)
  }
  
  
  /**
   * @typedef ITrackCalendarPageClick
   * @property {string} button_type 操作类型
   * @property {string} button_name 按钮名称
   * 
   * @description 神策数据上报 -- 日历页面点击
   * @param {ITrackCalendarPageClick} data
   * 
   */
  trackCalendarPageClick (data) {
    this.trackCatch('CalendarPageClick', data)
  }
  
  /**
   * @description 神策数据上报 -- 选择门店
   * @typedef ITrackChooseStore
   * @property {string} choose_scene 选择场景
   * @property {string} choose_entrance 选择入口
   * @property {string} seek_method 查找形式
   * @property {string} search_word 搜索词条
   * @property {string} filter_option 筛选项
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店编吗
   */
  trackChooseStore (data) {
    this.trackCatch('ChooseStore', data)
  }
  
  /**
   * @description 神策数据上报 -- 选择员工
   * @typedef ITrackChooseStaff
   * @property {string} choose_scene 选择场景
   * @property {string} choose_entrance 选择入口
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店编吗
   * @property {string} seek_method 查找形式
   * @property {string} search_word 搜索词条
   * @property {string} filter_option 筛选项
   * @property {string} SV_id 巡店ID
   * @property {string} visit_time 巡店时间
   * @property {string} staff_name 员工名字/邮箱
   * @property {string} staff_id 员工id
   * @property {string} staff_department 员工部门
   */
  trackChooseStaff (data) {
    this.trackCatch('ChooseStaff', data)
  }

  /**
   * @description 神策数据上报 -- 发布巡店计划
   * @typedef ITrackSubmitAppointment
   * @property {string} store_name_list 门店名称列表
   * @property {string} store_code_list 门店ID列表
   * @property {string} visit_time_list 巡店时间列表
   * @property {string} is_specific_time 是否具体时间
   * @property {string} is_SM_notified 是否通知门店店长
   * @property {string} sync_staff_list 同步人员列表
   */
  trackSubmitAppointment (data) {
    this.trackCatch('SubmitAppointment', data)
  }

  /**
   * @description 神策数据上报 -- 发布巡店计划详情
   * @typedef ITrackSubmitAppointmentDetail
   * @property {string} SV_id 巡店ID
   * @property {string} store_name 门店名称列表
   * @property {string} store_code 门店ID列表
   * @property {string} visit_time 巡店时间列表
   * @property {string} is_specific_time 是否具体时间
   * @property {string} is_SM_notified 是否通知门店店长
   * @property {string} sync_staff_list 同步人员列表
   */
  trackSubmitAppointmentDetail (data) {
    this.trackCatch('SubmitAppointmentDetail', data)
  }
  

  /**
   * @description 神策数据上报 -- 巡店计划浏览
   * @typedef ITrackStoreVisitAppointmentView
   * @property {string} SV_id 巡店ID
   * @property {string} SV_initiator 巡店发起人
   * @property {string} SV_initiator_id 巡店发起人工号
   * @property {string} SV_initiator_department 巡店发起人所在部门
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店ID
   * @property {string} visit_time 巡店时间
   * @property {string} is_synced 是否已同步给相关人员
   * 
   */
  trackStoreVisitAppointmentView (data) {
    this.trackCatch('StoreVisitAppointmentView', data)
  }

  /**
   * @description 神策数据上报 -- 巡店计划点击
   * @typedef ITrackStoreVisitAppointmentClick
   * @property {string} SV_id 巡店ID
   * @property {string} SV_initiator 巡店发起人
   * @property {string} SV_initiator_id 巡店发起人工号
   * @property {string} SV_initiator_department 巡店发起人所在部门
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店ID
   * @property {string} visit_time 巡店时间
   * @property {string} button_name 按钮名称
   * 
   */
  trackStoreVisitAppointmentClick (data) {
    this.trackCatch('StoreVisitAppointmentClick', data)
  }

  /**
   * @description 神策数据上报 -- 联系员工
   * @typedef ITrackFindContact
   * @property {string} SV_id 巡店ID
   * @property {string} SV_initiator 巡店发起人
   * @property {string} SV_initiator_id 巡店发起人工号
   * @property {string} SV_initiator_department 巡店发起人所在部门
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店ID
   * @property {string} visit_time 巡店时间
   * @property {string} role_type 按钮名称
   * @property {string} staff_name 员工姓名
   * @property {string} staff_id 员工ID
   * @property {string} staff_department 员工部门
   * 
   */
  trackFindContact (data) {
    this.trackCatch('FindContact', data)
  }

  /**
   * @description 神策数据上报 -- 同步巡店计划
   * @typedef ITrackSyncStoreVisit
   * @property {string} SV_id 巡店ID
   * @property {string} SV_initiator 巡店发起人
   * @property {string} SV_initiator_id 巡店发起人工号
   * @property {string} SV_initiator_department 巡店发起人所在部门
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店ID
   * @property {string} visit_time 巡店时间
   * @property {string} is_new_added 是否新增
   * @property {string} sync_staff_list 同步人员列表
   * 
   */
  trackSyncStoreVisit (data) {
    this.trackCatch('SyncStoreVisit', data)
  }

  /**
   * @description 神策数据上报 -- 开始巡店
   * @typedef ITrackStartStoreVisit
   * @property {string} choose_entrance 选择入口
   * @property {string} SV_id 巡店ID
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店ID
   * @property {string} visit_time 巡店时间
   * @property {string} button_name 按钮名称
   * 
   */
  trackStartStoreVisit (data) {
    this.trackCatch('StartStoreVisit', data)
  }

  /**
   * @description 神策数据上报 -- 打卡
   * @typedef ITrackCheckIn
   * @property {string} previous_action 前向来源
   * @property {string} SV_id 巡店ID
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店ID
   * @property {string} visit_time 巡店时间
   * 
   */
  trackCheckIn (data) {
    this.trackCatch('CheckIn', data)
  }

  /**
   * @description 神策数据上报 -- 提交巡店记录
   * @typedef ITrackSubmitReport
   * @property {string} SV_id 巡店ID
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店ID
   * @property {string} visit_time 巡店时间
   * @property {string} is_synced 是否已同步给相关人员
   * @property {string} feedback_quantity 反馈数量
   * @property {string} feedback_category_list 反馈分类列表
   * @property {string} button_name 按钮名称
   * 
   */
  trackSubmitReport (data) {
    this.trackCatch('SubmitReport', data)
  }

  /**
   * @description 神策数据上报 -- 提交巡店记录详情
   * @typedef ITrackSubmitReportDetail
   * @property {string} SV_id 巡店ID
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店ID
   * @property {string} visit_time 巡店时间
   * @property {string} is_synced 是否已同步给相关人员
   * @property {string} feedback_id 反馈ID
   * @property {string} feedback_category 反馈分类
   * @property {string} is_customized 是否自定义
   * @property {string} is_picture 是否上传照片/视频
   * @property {string} is_text 是否输入文字
   * @property {string} button_name 按钮名称
   * 
   */
  trackSubmitReportDetail (data) {
    this.trackCatch('SubmitReportDetail', data)
  }

  /**
   * @description 神策数据上报 -- 编辑报告概览
   * @typedef ITrackEditReportSummary
   * @property {string} SV_id 巡店ID
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店ID
   * @property {string} visit_time 巡店时间
   * @property {string} is_synced 是否已同步给相关人员
   * 
   * @property {string} summary_type 概览类型
   * @property {string} edit_type 编辑类型
   * 
   */
  trackEditReportSummary (data) {
    this.trackCatch('EditReportSummary', data)
  }

  /**
   * @description 神策数据上报 -- 删除反馈
   * @typedef ITrackDeleteFeedback
   * @property {string} SV_id 巡店ID
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店ID
   * @property {string} visit_time 巡店时间
   * @property {string} is_synced 是否已同步给相关人员
   * 
   * @property {string} feedback_category 反馈分类
   * @property {string} is_customized 是否自定义
   * @property {string} is_picture 是否上传照片/视频
   * @property {string} is_text 是否输入文字
   * 
   */
  trackDeleteFeedback (data) {
    this.trackCatch('DeleteFeedback', data)
  }

  /**
   * @description 神策数据上报 -- 推送报告
   * @typedef ITrackdSyncReport
   * @property {string} SV_id 巡店ID
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店ID
   * @property {string} visit_time 巡店时间
   * @property {string} is_synced 是否已同步给相关人员
   * 
   * @property {string} is_new_added 是否新增
   * @property {string} is_sa_visible 是否门店人员可见
   * @property {string} is_relevant_visible 是否巡店相关人员可见
   * @property {string} sync_staff_list 同步人员列表
   * 
   */
  trackSyncReport (data) {
    this.trackCatch('SyncReport', data)
  }

  /**
   * @description 神策数据上报 -- 下载报告
   * @typedef ITrackDownloadReport
   * @property {string} SV_id 巡店ID
   * @property {string} SV_initiator 巡店发起人
   * @property {string} SV_initiator_id 巡店发起人工号
   * @property {string} SV_initiator_department 巡店发起人所在部门
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店ID
   * @property {string} visit_time 巡店时间
   * @property {string} report_time 报告时间
   * 
   */
  trackDownloadReport (data) {
    this.trackCatch('DownloadReport', data)
  }
  

  /**
   * @description 神策数据上报 -- 评论
   * @typedef ITrackComments
   * @property {string} SV_id 巡店ID
   * @property {string} SV_initiator 巡店发起人
   * @property {string} SV_initiator_id 巡店发起人工号
   * @property {string} SV_initiator_department 巡店发起人所在部门
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店ID
   * @property {string} visit_time 巡店时间
   * @property {string} report_time 报告时间
   * 
   * @property {string} comment_type 评论类型
   * @property {string} comment_id 评论ID
   * @property {string} is_picture 是否上传照片/视频
   * @property {string} is_text 是否输入文字
   * 
   * 
   */
  trackComments (data) {
    this.trackCatch('Comments', data)
  }

  /**
   * @description 神策数据上报 -- 巡店筛选
   * @typedef ITrackSVFilter
   * 
   * @property {string} current_page 当前页面
   * @property {string} filter_type_sv 筛选类型
   * @property {string} filter_option 筛选项
   * 
   */
  trackSVFilter (data) {
    this.trackCatch('SVFilter', data)
  }

  /**
   * @description 神策数据上报 -- 看板点击
   * @typedef ITrackDashboardClick
   * 
   * @property {string} button_type 操作类型
   * @property {string} view_dimension 查看维度
   * @property {string} period 统计范围
   * @property {string} only_initiated_oneself 是否仅看我发起的巡店
   * @property {string} SV_cnt 巡店数
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店ID
   * 
   */
  trackDashboardClick (data) {
    this.trackCatch('DashboardClick', data)
  }

  /**
   * @description 神策数据上报 -- 巡店记录导出
   * @typedef ITrackSVExport
   * 
   * @property {string} view_dimension 查看维度
   * @property {string} period 统计范围
   * @property {string} only_initiated_oneself 是否仅看我发起的巡店
   * @property {string} SV_cnt 巡店数
   * 
   * 
   */
  trackSVExport (data) {
    this.trackCatch('SVExport', data)
  }

  /**
   * @typedef ITrackSVPlanClick
   * @property {string} SV_id 巡店ID
   * @property {string} SV_initiator 巡店发起人
   * @property {string} SV_initiator_id 巡店发起人工号
   * @property {string} SV_initiator_department 巡店发起人所在部门
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店ID
   * @property {string} visit_time 巡店时间
   * @property {string} report_time 报告时间
   * 
   * @property {string} plan_status 记录状态
   * @property {string} button_name 按钮名称
   * 
   * 
   * @description 神策数据上报 -- 巡店记录点击
   * @param {ITrackSVPlanClick} data
   */
  trackSVPlanClick (data) {
    this.trackCatch('SVPlanClick', data)
  }
  

  /**
   * @typedef ITrackAddSAPageView
   * @property {string} display_campaign_name 活动名称
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店编码
   * @property {string} SA_name 导购名称
   * @property {string} SA_id 导购ID
   * @property {string} QR_code_id 二维码ID
   * @property {string} leading_words 添加引导语
   * @property {string} SA_added_status 用户状态
   * @property {string} campaign_status 活动状态
   * 
   * 
   * @description 神策数据上报 -- 添加专属顾问页浏览
   * @param {ITrackAddSAPageView} data
   */
  trackAddSAPageView (data) {
    this.trackCatch('AddSAPageView', data)
  }
  

  /**
   * @typedef ITrackScanQRCode
   * @property {string} display_campaign_name 活动名称
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店编码
   * @property {string} SA_name 导购名称
   * @property {string} SA_id 导购ID
   * @property {string} QR_code_id 二维码ID
   * @property {string} leading_words 添加引导语
   * @property {string} SA_added_status 用户状态
   * @property {string} campaign_status 活动状态
   * 
   * 
   * @description 神策数据上报 -- 添加专属顾问页浏览
   * @param {ITrackScanQRCode} data
   */
  trackScanQRCode (data) {
    this.trackCatch('ScanQRCode', data)
  }

  /**
   * @typedef ITrackChooseOtherSA
   * @property {string} display_campaign_name 活动名称
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店编码
   * @property {string} SA_name 导购名称
   * @property {string} SA_id 导购ID
   * @property {string} QR_code_id 二维码ID
   * @property {string} leading_words 添加引导语
   * @property {string} SA_added_status 用户状态
   * @property {string} campaign_status 活动状态
   * 
   * 
   * @description 神策数据上报 -- 添加专属顾问页浏览
   * @param {ITrackChooseOtherSA} data
   */
  trackChooseOtherSA (data) {
    this.trackCatch('ChooseOtherSA', data)
  }

  /**
   * @typedef ITrackStoreListView
   * @property {string} display_campaign_name 活动名称
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店编码
   * @property {string} SA_name 导购名称
   * @property {string} SA_id 导购ID
   * @property {string} QR_code_id 二维码ID
   * @property {string} leading_words 添加引导语
   * @property {string} SA_added_status 用户状态
   * @property {string} campaign_status 活动状态
   * 
   * 
   * @description 神策数据上报 -- 添加专属顾问页浏览
   * @param {ITrackStoreListView} data
   */
  trackStoreListView (data) {
    this.trackCatch('StoreListView', data)
  }

  /**
   * @typedef ITrackLBSAuth
   * @property {string} display_campaign_name 活动名称
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店编码
   * @property {string} SA_name 导购名称
   * @property {string} SA_id 导购ID
   * @property {string} QR_code_id 二维码ID
   * @property {string} leading_words 添加引导语
   * @property {string} SA_added_status 用户状态
   * @property {string} campaign_status 活动状态
   * 
   * 
   * @description 神策数据上报 -- 添加专属顾问页浏览
   * @param {ITrackLBSAuth} data
   */
  trackLBSAuth (data) {
    this.trackCatch('LBSAuth', data)
  }

  /**
   * @typedef ITrackAddSAButtonClick
   * @property {string} display_campaign_name 活动名称
   * @property {string} store_name 门店名称
   * @property {string} store_code 门店编码
   * @property {string} SA_name 导购名称
   * @property {string} SA_id 导购ID
   * @property {string} QR_code_id 二维码ID
   * @property {string} leading_words 添加引导语
   * @property {string} SA_added_status 用户状态
   * @property {string} campaign_status 活动状态
   * 
   * 
   * @description 神策数据上报 -- 添加专属顾问页浏览
   * @param {ITrackAddSAButtonClick} data
   */
  trackAddSAButtonClick (data) {
    this.trackCatch('AddSAButtonClick', data)
  }
  
}
