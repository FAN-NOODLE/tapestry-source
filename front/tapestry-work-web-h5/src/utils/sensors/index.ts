import sensors from 'sa-sdk-javascript'
import SensorsTrack from './track';
import dayjs from 'dayjs';
import Helper from '../helper'
import { API } from '../env-config';
import store from '@/store';
import pageleave from 'sa-sdk-javascript/dist/web/plugin/pageleave/index.es6.js';
import pageload from 'sa-sdk-javascript//dist/web/plugin/pageload/index.es6.js';

function getSensorsOptions() {
  const config = {
    /** 是否开启 */
    isOpen: true,
    /** 初始化配置 */
    initOptions: {
      /** 上报地址 */
      server_url: API.server_url,
      /** 上报日志 */
      show_log: true,
    },
    /** 自定全局属性 */
    globalOptions: {
      // title: () => document.title,
      /** 应用名称 */
      app_name() {
        // location.pathname.split('/').slice(3).join('/')
        const userInfo:IUserInfo = store.getters["user/getUserInfo"]
        const role = store.getters["user/getRole"]
        const isOnlineOrder = location.href.indexOf('online-order') !== -1
        // 网仓订单
        if(isOnlineOrder) {
          if(!userInfo.user_id) {
            return 'OnlineWarehouse_unknown'
          } else if(role === 5) {
            return 'OnlineWarehouse_Office'
          } else {
            return 'OnlineWarehouse_SA'
          }
        }

        const isTransmital = location.href.indexOf('communication') !== -1
        // 传讯
        if(isTransmital) {
          if(!userInfo.user_id) {
            return 'StoreTransmittal_unknown'
          } else if(userInfo.is_ro === 1) {
            return 'StoreTransmittal_RO'
          } else if(userInfo.is_ft === 1) {
            return 'StoreTransmittal_FT'
          } else {
            return 'StoreTransmittal_SA'
          }
        }

        const isVM = location.href.indexOf('display-task') !== -1
        // VM
        if(isVM) {
          if(!userInfo.user_id) {
            return 'WeStoreVM_unknown'
          } else if(userInfo.is_vm === 1) {
            return 'WeStoreVM_vm'
          } else {
            return 'WeStoreVM_sa'
          }
        }
        
        const isStoreVisit = location.href.indexOf('store-visit') !== -1
        // 巡店
        if(isStoreVisit) {
          console.log(1111, userInfo);
          if(!userInfo.user_id) {
            return 'StoreVisit_unknown'
          } else if(userInfo.is_store_visit_manager) {
            return 'StoreVisit_CorporateUser'
          } else if(userInfo.is_store_manager) {
            return 'StoreVisit_StoreManager'
          } else {
            return 'StoreVisit_SA'
          }
        }

        const isUser = location.pathname.split('/').slice(3).join('/').startsWith('user')
        if(isUser) {
          return 'Wecom_QRCode'
        }

        return API.app_name
      },
    },
  }
  console.log(config);
  
  return config
}

const sensorsOptions = getSensorsOptions()
/**
 * @description 神策数据上报 -- 添加配置
*/
export function sensorsInit(){
  const { initOptions, globalOptions } = sensorsOptions;
  console.log(sensors);
  
  // 注册公共属性
  sensors.registerPage(globalOptions);
  // 初始化神策数据
  sensors.use(pageleave,{
    custom_props:{
      $title: function () {
        // console.log('神策数据上报页面标题', store.getters["storeVisit/getTitle"]);
        return store.getters["storeVisit/getTitle"]
      },
    },
    heartbeat_interval_time: 5,
    max_duration: 5 * 24 * 60 * 60,
    isCollectUrl: function(url){
      // url 为要采集页面浏览时长的页面地址。
      return true; // 采集
      // return false; // 不采集
    }
  });
  // 见下文初始化参数说明
  sensors.use(pageload,{max_duration:120});
  
  // 初始化设置
  sensors.init({
    heatmap: {
       //是否开启点击图，default 表示开启，自动采集 $WebClick 事件，可以设置 'not_collect' 表示关闭。
       clickmap: 'default',
       //是否开启触达图，not_collect 表示关闭，不会自动采集 $WebStay 事件，可以设置 'default' 表示开启。
       scroll_notice_map: 'default'
    },
    is_track_single_page: true,
    use_client_time: true,
    send_type:'beacon',
    ...initOptions
  });
  

  (sensors as any).registerPropertyPlugin({
    isMatchedWithFilter:function(data){
      const pathEnd = data.properties['$url'].split('?')[0].split('/').pop()
      // console.log(222, pathEnd);
      return data.event === "$pageview" && pathEnd === 'store-visit'
    },
    properties:function(data){
      data.properties['$title'] = '巡店管理';
    }
  });
  (sensors as any).registerPropertyPlugin({
    isMatchedWithFilter:function(data){
      const pathEnd = data.properties['$url'].split('?')[0].split('/').pop()
      // console.log(222, pathEnd);
      return data.event === "$pageview" && pathEnd === 'store-visit-dashboard'
    },
    properties:function(data){
      data.properties['$title'] = '巡店工具';
    }
  });
  
  (sensors as any).registerPropertyPlugin({
    isMatchedWithFilter:function(data){
      const pathEnd = data.properties['$url'].split('?')[0].split('/').pop()
      // console.log(1111, pathEnd);
      return data.event === "$pageview" && pathEnd === 'store-visit-record'
    },

    properties:function(data){
      data.properties['$title'] = '巡店记录';
    }
  });

  (sensors as any).registerPropertyPlugin({
    isMatchedWithFilter:function(data){
      const pathEnd = data.properties['$url'].split('?')[0].split('/').pop()
      // console.log(1111, pathEnd);
      return data.event === "$pageview" && pathEnd === 'store-visit-report'
    },

    properties:function(data){
      data.properties['$title'] = data.properties['$title'] === '巡访门店' ? '巡店计划' : data.properties['$title'];
    }
  });

  (sensors as any).registerPropertyPlugin({
    isMatchedWithFilter:function(data){
      const isUser = location.pathname.split('/').slice(3).join('/').startsWith('user')
      // console.log(1111, pathEnd);
      return isUser
    },

    properties:function(data){
      // 获取链接上的参数
      const unionid = Helper.getUrlParam('unionid');
      const state = Helper.getUrlParam('scene')
      const activityDetail= JSON.parse(localStorage.getItem('activityDetail') || '{}')
      data.properties['union_id'] = unionid
      data.properties['page_type'] = '企微智慧导购'
      data.properties['page_title'] = data.properties['$title']
      data.properties['campaign_name'] = activityDetail?.name
      data.properties['campaign_id'] = activityDetail?.id
      data.properties['valid_period_type'] = ['长期有效', '指定时间'][activityDetail?.time_type] || ''
      if(activityDetail?.time_type == 1) {
        data.properties['valid_period'] = `${dayjs(activityDetail?.start_time*1000).format('YYYY-MM-DD')}至${dayjs(activityDetail?.end_time*1000).format('YYYY-MM-DD')}`
      }
      data.properties['campaign_scope'] = [...(activityDetail?.scope_area_list||[]), ...(activityDetail?.scope_member_list||[])].map(item => item.name).join(',')

      activityDetail?.category_list.forEach(item => {
        if(item.query_str.indexOf(state) > -1) {
          data.properties['campaign_channel'] = item.name
          data.properties['qr_poster_code_id'] = state
        } else {
          item.point_list.forEach(point => {
            if(point.query_str.indexOf(state) > -1) {
              data.properties['campaign_channel'] = `${item.name}_${point.name}`
              data.properties['qr_poster_code_id'] = String(point.id)
            }
          });
        }
      });
      
    }
  });

  // 设置之后，SDK 就会自动收集页面浏览事件，以及设置初始来源。
  sensors.quick('autoTrack')
  return sensors
}

/**
 * @description 登录并上报
 * @param {string} userId 用户ID
 */
export function sensorsLogin(userId:string) {
  sensors.login(userId)
}

/**
 * @description 设置用户属性
 */
export function setSensorsProfile (data:object) {
  console.log('神策上报，设置用户属性', data)
  sensors.setProfile(data)
}
export const sensorsTrack = new SensorsTrack({
  sensors,
  open: sensorsOptions.isOpen,
  keys: {
    /** 邀约ID */
    reserve_number: 'appointment_id',
    /** 顾客ID */
    user_external_id: 'customer_id',
    /** 服务类型 */
    reserve_name: 'service_type',
    /** 服务类型 */
    serve_name: 'service_type',
    /** 到店时间日期 */
    reserve_time: {
      key: 'service_date',
      value: (val:number) => dayjs(val * 1000).format('YYYY/MM/DD')
    },
    /** 到店时间时间段 */
    time_period: {
      key: 'service_timeslot',
      value: (val:number) => dayjs(val * 1000).format('HH:mm:ss')
    },
    /** 门店名称 */
    // store_name: 'store_name',
    /** 门店编码 */
    // store_code: 'store_code',

    /** 顾客类型 */
    customer_status: 'customer_type',
    /** 是否到店购买 */
    is_store: 'is_instore',
    /** 订单金额 */
    reserve_order_price: 'order_amount',
    /** 修改项 */
    change_item: 'modified_item',
    /** 修改后服务类型 */
    modify_reserve_name: 'modified_service_type',
    /** 修改后到店时间日期 */
    modify_reserve_time: 'modified_service_date',
    /** 修改后到店时间时间段 */
    modify_time_period: 'modified_service_timeslot',
    /** 取消类型 */
    cancel_type: {
      key: 'cancellation_type',
      value: ['', '手动']
    },
    /** 邀请函ID */
    invite_record_id: {
      key: 'invite_id',
      value: (val:number) => val + ''
    },
    /** 邀请函模版ID */
    invite_template_id: {
      key: 'invite_template_id',
      value: (val:number) => val + ''
    },
    /** 入口 */
    entrance: {
      key: 'entry_point',
      value: ['', '企微侧边栏', '顾客详情', '客户邀约']
    },
    /** 任务ID */
    task_id: 'task_id',
    /** 任务名称 */
    task_name: 'task_name',
    /** 任务覆盖顾客数 */
    customer_total: 'task_coverage',
    /** 任务发布时间 */
    start_time: 'task_starttime',
    /** 任务状态 */
    task_status: ['待完成', '已完成', '未完成'],
    /** 素材ID */
    asset_id: {
      key: 'asset_id',
      value: (val:number) => val + ''
    },
    /** 播种ID */
    sow_task_id: {
      key: 'reminder_id',
      value: (val:number) => val + ''
    },
    /** 播种内容 */
    sow_task_content: 'reminder_content',
    /** 榜单周期 */
    ranking_cycle: {
      key: 'ranking_periodicity',
      value: {
        month: '月榜',
        year: '年榜'
      }
    },
    /** 榜单类型 */
    ranking_type: ['', '明星榜', '邀约榜', '顾客榜', '活跃榜'],
    /** 时间段（月份） */
    time_range: {
      key: 'ranking_daterange',
      value: ([start, end]:number[]) => `${dayjs(start * 1000).format('YYYYMMDD')}-${dayjs(end * 1000).format('YYYYMMDD')}`
    },
    /** 会员查询条件 */
    query_condition: 'inquiry_condition',
    /** 会员查询结果个数 */
    total: 'result_number',
    /** 分配ID */
    inherit_id: 'inherit_id',
    /** 分配前导购ID */
    handover_userid: 'handover_userid',
    /** 分配后导购ID */
    takeover_userid: 'takeover_userid',
    /** 顾客筛选 */
    manage_type: 'by_customertype',
    /** 员工筛选 */
    member_list: {
      key: 'by_saname',
      value: (val:string[]) => val.join()
    },
    /** 会员卡等 */
    card_type: {
      key: 'by_viplevel',
      value: (val:string[]) => val.join()
    },
    /** 印象标签 */
    impression_tag_list: {
      key: 'impression_tag',
      value: (val:string[]) => val.join()
    },
    /** 客户标签 */
    tag_list: {
      key: 'customer_tag',
      value: (val:string[]) => val.join()
    },
    /** CLT类型 */
    clt_type: {
      key: 'clt_type',
      value: (val:string[]) => val.join()
    },
    /** 搜索原始词 */
    original_word: 'search_key_word',
    /** 评价结果详情 */
    evaluation_result: 'review_details',
    /** 活动安排星级 */
    activity_star: 'activity_rating',
    /** 整体服务星级 */
    service_star: 'service_rating',
    
    /** 查看素材 */
    AssetView: {
      task_id: (val) => val + '',
      source_id: {
        key: 'asset_id',
        value: (val) => val + ''
      },
      is_official: (val) => ({
        1: '普通素材',
        2: '官网商品素材'
      }[val]),
      asset_type_modify: (val) => ({
        1: '文本素材',
        2: '图片素材',
        3: '链接素材',
        4: '小程序素材',
        5: '视频素材',
        6: '文件素材',
        7: '朋友圈素材'
      }[val]),
      official_url: 'official_url',
      official_upc: 'official_upc',
      official_product: 'official_product',
      official_name: 'official_name',


      is_top: 'is_top',
      library_type: {
        key: 'library_type',
        value: (val) => ({
          1: '品牌素材库',
          2: '素材库',
          3: '顾客列表群发素材'
        }[val])
      },
      entry_point: {
        key: 'entry_point',
        value: (val) => ({
          1: '工作台',
          2: '1V1聊天窗口',
          3: '顾客列表'
        }[val])
      },
      current_page: {
        key: 'current_page',
        value: (val) => ({
          '/material-library-coach': '素材库',
          '/my-collect': '我的收藏',
          '/material-detail': '官网素材详情',
        }[val] || val)
      },
      is_favorite: 'is_favorite'
    },
    /** 发送素材 */
    AssetSend: {
      asset_id: 'asset_id',
      customer_id: 'customer_id',
      task_id: 'task_id',
      operation_type: {
        key: 'operation_type',
        value: ['群发', '转发朋友圈', '普通发送', '沟通任务转发']
      },
      is_official: (val) => ({
        1: '普通素材',
        2: '官网商品素材'
      }[val]),
      asset_type_modify: (val) => ({
        1: '文本素材',
        2: '图片素材',
        3: '链接素材',
        4: '小程序素材',
        5: '视频素材',
        6: '文件素材',
        7: '朋友圈素材'
      }[val]),
      official_url: 'official_url',
      official_upc: 'official_upc',
      official_product: 'official_product',
      official_name: 'official_name',

      is_top: 'is_top',
      library_type: {
        key: 'library_type',
        value: (val) => ({
          1: '品牌素材库',
          2: '素材库',
          3: '顾客列表群发素材'
        }[val])
      },
      entry_point: {
        key: 'entry_point',
        value: (val) => ({
          1: '工作台',
          2: '1V1聊天窗口',
          3: '顾客列表'
        }[val])
      },
      current_page: {
        key: 'current_page',
        value: (val) => ({
          '/material-library-coach': '素材库',
          '/my-collect': '我的收藏',
          '/material-detail': '官网素材详情',
        }[val] || val)
      },
      is_favorite: 'is_favorite'
    },
    AssetFavorite: {
      current_page: {
        key: 'current_page',
        value: (val) => ({
          '/material-library-coach': '素材库',
          '/my-collect': '我的收藏',
          '/material-detail': '官网素材详情',
        }[val] || val)
      },
      is_top: 'is_top',
      entry_point: {
        key: 'entry_point',
        value: (val) => ({
          1: '工作台',
          2: '1V1聊天窗口',
          3: '顾客列表'
        }[val])
      },
      library_type: {
        key: 'library_type',
        value: (val) => ({
          1: '品牌素材库',
          2: '素材库',
          3: '顾客列表群发素材'
        }[val])
      },
      operation_position: {
        key: 'operation_position',
        value: (val) => ({
          1: '列表',
          2: '放大'
        }[val])
      },
      asset_id: {
        key: 'asset_id',
        value: (val) => val + ''
      },
      operation_type: {
        key: 'operation_type',
        value: (val) => ({
          1: '收藏',
          2: '取消收藏'
        }[val])
      },
      is_official: (val) => ({
        1: 'KOL素材',
        2: '官网商品素材'
      }[val]),
      asset_type_modify: (val) => ({
        1: '文本素材',
        2: '图片素材',
        3: '链接素材',
        4: '小程序素材',
        5: '视频素材',
        6: '文件素材',
        7: '朋友圈素材',
        8: '官网商品'
      }[val]),
      official_url: 'official_url',
      official_upc: 'official_upc',
      official_product: 'official_product',
      official_name: 'official_name'
    },
    MyFavoriteEntranceClick: {
      operation_type: 'operation_type',
      entry_point: 'entry_point',
      library_type: 'library_type'
    },
    /** 进入传讯 */
    EnterTransmittal: {
      transmittal_entrance: {
        key: 'transmittal_entrance',
        value: {
          1: '企微工作台',
          2: '企业微信消息'
        }
      }
    },
    /** 浏览传讯内容 */
    ViewTransmittal: {
      transmittal_title: 'transmittal_rubric',
      transmittal_publisher: 'transmittal_publisher',
      release_time: {
        key: 'release_time',
        value: (val) => val ? dayjs(val*1000).format('YYYY.MM.DD HH:mm:ss') : '待提交'
      },
      transmittal_entrance: {
        key: 'transmittal_entrance',
        value: {
          ft: '我的传讯-查看',
          ro: '传讯审核-去审核',
          sa: '企业传讯',
          4: '企业微信消息'
        }
      }
    },
    /** 下载传讯附件 */
    DownloadTransmittal: {
      transmittal_title: 'transmittal_rubric',
      transmittal_publisher: 'transmittal_publisher',
      release_time: {
        key: 'release_time',
        value: (val) => val ? dayjs(val*1000).format('YYYY.MM.DD HH:mm:ss') : '待提交'
      },
    },
    /** 撤回传讯 */
    WithdrawTransmittal: {
      transmittal_title: 'transmittal_rubric',
      transmittal_publisher: 'transmittal_publisher',
      release_time: {
        key: 'release_time',
        value: (val) => val ? dayjs(val*1000).format('YYYY.MM.DD HH:mm:ss') : '待提交'
      },
    },
    /** 作废传讯 */
    VoidTransmittal: {
      transmittal_title: 'transmittal_rubric',
      transmittal_publisher: 'transmittal_publisher',
      release_time: {
        key: 'release_time',
        value: (val) => val ? dayjs(val*1000).format('YYYY.MM.DD HH:mm:ss') : '待提交'
      },
    },
    /** 审核传讯 */
    AuditTransmittal: {
      transmittal_title: 'transmittal_rubric',
      transmittal_publisher: 'transmittal_publisher',
      release_time: {
        key: 'release_time',
        value: (val) => val ? dayjs(val*1000).format('YYYY.MM.DD HH:mm:ss') : '待提交'
      },
      audit_type: {
        key: 'audit_type',
        value: {
          0: '不通过',
          1: '通过',
        }
      }
    },
    /** 进入vm陈列任务 */
    EnterDisplayTask: {
      vm_entrance: {
        key: 'vm_entrance',
        value: {
          1: '企微工作台',
          2: '企微消息'
        }
      }
    },
    /** 查看他店陈列图 */
    ViewStoreDisplay: {
      dispalytask_title: 'dispalytask_title',
      displaytask_complete_state: {
        key: 'displaytask_complete_state',
        value: {
          0: '待提交',
          1: '待审核',
          2: '未通过',
          3: '已通过',
        }
      },
      displaytask_state: {
        key: 'displaytask_state',
        value: {
          0: '待开始',
          1: '进行中',
          2: '已结束',
        }
      },
      display_store: 'display_store',
    },
    /** 查看陈列任务 */
    ViewDisplayTask: {
      dispalytask_title: 'dispalytask_title',
      displaytask_complete_state: {
        key: 'displaytask_complete_state',
        value: {
          0: '待提交',
          1: '待审核',
          2: '未通过',
          3: '已通过',
        }
      },
      displaytask_state: {
        key: 'displaytask_state',
        value: {
          0: '待开始',
          1: '进行中',
          2: '已结束',
        }
      },
      display_store: 'display_store'
    },
    /** 提交陈列任务 */
    SubmitDisplayTask: {
      dispalytask_title: 'dispalytask_title',
      displaytask_complete_state: {
        key: 'displaytask_complete_state',
        value: {
          0: '待提交',
          1: '待审核',
          2: '未通过',
          3: '已通过',
        }
      },
      displaytask_state: {
        key: 'displaytask_state',
        value: {
          0: '待开始',
          1: '进行中',
          2: '已结束',
        }
      },
      is_resubmit: {
        key: 'is_resubmit',
        value: {
          0: '否',
          1: '是',
        }
      },
    },
    /** 去学习 */
    CourseLearing: {
      course_title: 'course_title',
    },
    /** 陈列提醒（vm视角） */
    DisplayReminder: {
      reminder_type: {
        key: 'reminder_type',
        value: {
          1: '一键提醒',
          2: '点击去提醒'
        }
      },
    },
    /** 操作陈列审核通过（vm视角） */
    MarkPass: {
      display_store: 'display_store',
      display_score: 'display_score',
    },
    /** 操作陈列审核不通过（vm视角） */
    MarkNotPass: {
      display_store: 'display_store',
    },
    /** 进入素材库 */
    AssetLibrary: {
      entry_point: (value) => ({
        1: '工作台',
        2: '1V1聊天窗口',
        3: '顾客列表',
        4: '库存查询'
      }[value]),
      library_type: (value) => ({
        1: '品牌素材库',
        2: '素材库',
        3: '顾客列表群发素材',
        4: '库存查询'
      }[value])
    },
    /** 查询官网商品素材 */
    AssetInquiry: {
      inquiry_type: (value) => ({
        'sku_item': 'skuitem',
        'product_name': '商品名称',
        'upc_code': '条形码查询'
      }[value])
    },
    /** 库存查询 */
    StockInquiry: {
      entry_point: (value) => ({
        0: '我的-库存查询',
        1: '品牌素材库',
        2: '顾客列表群发素材',
        3: '素材库'
      }[value]),
      inquiry_type: (value) => ({
        0: 'skuitem',
        1: '条形码查询',
        2: '商品库存查询'
      }[value])
    },
    /** 一键下载传讯及所有附件 */
    ClickAllTransmittalDownload: {
      transmittal_type: (value) => ({
        1: '周传讯',
        2: '转货传讯',
        3: '紧急传讯',
      }[value]),
    },
    /** 下载传讯及所有附件 */
    ClickAllAttachmentDownload: {
      transmittal_type: (value) => ({
        1: '周传讯',
        2: '转货传讯',
        3: '紧急传讯',
      }[value]),
    },
    /** 下载传讯文件 */
    ClickTransmittalDownload: {
      transmittal_type: (value) => ({
        1: '周传讯',
        2: '转货传讯',
        3: '紧急传讯',
      }[value]),
    },
    /** 搜索点击 */
    SearchClick: {
      tab_name: 'tab_name',
      current_page: 'current_page',
    },
    /** 筛选 */
    Filter: {
      tab_name: 'tab_name',
      filter_type_sv: 'filter_type_sv',
      filter_option: 'filter_option',
    },
    ClickItem: {
      feature: {
        key: 'feature',
        value: {
          '/customer-portrait': '客户画像',
          '/invite-config': '邀请函',
          '/store-invited': '到店邀约',
          '/material-library': '素材库',
          '/online-order-create': '网仓库存',
        }
      },
      customer_id: 'customer_id'
    },
    AddItems: {
      add_method: 'add_method',
      commodity_name: 'commodity_name',
      sku_id: 'sku_id',
      customer_id: 'customer_id',
    },
    SubmitPopupClick: {
      sales_order_id: 'sales_order_id',
      has_reject_risk: 'has_reject_risk',
      risk_type: 'risk_type',
      button_name: 'button_name',
      customer_id: 'customer_id',
    },
    OrderCreated: {
      sales_order_id: 'sales_order_id',
      has_reject_risk: 'has_reject_risk',
      risk_type: 'risk_type',
      customer_id: 'customer_id',
    },
    OrderCreatedDetail: {
      sales_order_id: 'sales_order_id',
      commodity_name: 'commodity_name',
      sku_id: 'sku_id',
      commodity_quantity: 'commodity_quantity',
      inventory_inquiry_status: 'inventory_inquiry_status',
      customer_id: 'customer_id',
    },
    ToCustomerPopupClick: {
      sales_order_id: 'sales_order_id',
      pop_name: 'pop_name',
      button_name: 'button_name',
      customer_id: 'customer_id',
    },
    EnterOnlineWarehouse: {
      OnlineWarehouse_entrance: 'OnlineWarehouse_entrance',
      message_type: 'message_type',
    },
    OrderFilter: {
      current_page: 'current_page',
      filter_type: 'filter_type_sv',
      filter_option: 'filter_option',
    },
    TabSwitch: {
      tab_name: 'tab_name',
    },
    OrderListClick: {
      sales_order_id: 'sales_order_id',
      order_status: 'order_status',
      order_creator: 'order_creator',
      order_belonged_store: 'order_belonged_store',
      customer_id: 'customer_id',
    },
    OrderDetailClick: {
      sales_order_id: 'sales_order_id',
      order_status: 'order_status',
      order_creator: 'order_creator',
      order_belonged_store: 'order_belonged_store',
      customer_id: 'customer_id',
      button_name: 'button_name',
    },
    ExportOrderList: {
      period: 'period',
      store_name: 'store_name',
      sales_name: 'sales_name',
      original_words: 'original_words',
      order_status: 'order_status',
    },
    
  }
})