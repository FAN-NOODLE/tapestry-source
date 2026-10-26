import sensors from 'sa-sdk-miniprogram'
import SensorsTrack from './track'
import dayjs from 'dayjs'

const { env, sensors_server_url, sensors_app_name } = require('@/lib/config')

/** 添加配置 */
export function setSensorsConfig () {
  // 配置初始化参数
  sensors.setPara({
    name: 'sensors',
    server_url: sensors_server_url,
    // 全埋点控制开关

    autoTrack: {
      appLaunch: true, // 默认为 true，false 则关闭 $MPLaunch 事件采集
      appShow: false, // 默认为 true，false 则关闭 $MPShow 事件采集
      appHide: false, // 默认为 true，false 则关闭 $MPHide 事件采集
      pageShow: true, // 默认为 true，false 则关闭 $MPViewScreen 事件采集
      pageShare: true, // 默认为 true，false 则关闭 $MPShare 事件采集
      mpClick: true, // 默认为 false，true 则开启 $MPClick 事件采集
      mpFavorite: false, // 默认为 true，false 则关闭 $MPAddFavorites 事件采集
      pageLeave: true // 默认为 false， true 则开启 $MPPageLeave事件采集
    },
    // 自定义渠道追踪参数，如 source_channel: ["custom_param"]
    source_channel: [],
    // 是否允许控制台打印查看埋点数据(建议开启查看)
    show_log: env === 'test',
    // 是否允许修改 onShareAppMessage 里 return 的 path，用来增加(登录 ID，分享层级，当前的 path)，在 app onShow 中自动获取这些参数来查看具体分享来源、层级等
    allow_amend_share_path: true
  })

  // 设置公共属性
  sensors.registerApp({
    // 应用名称 -- 不同企业名称不一样
    app_name() {
      return sensors_app_name
    }
  })
}

/**
 * @description 初始化
 */
export function setInit () {
  sensors.init()
}

/**
 * @description 设置user ID 并启动上报
 * @param {string} userId 用户ID
 */
export function sensorsLoginInit (userId) {
  console.log('神策上报，设置user ID 并启动上报', userId)
  // TODO 神策数据上报 写好之后记得打开
  // if (env === 'test') return
  sensors.login(userId)
}

/**
 * @description 设置用户属性
 * @param {object} data
 */
export function setSensorsProfile (data) {
  console.log('神策上报，设置用户属性', data)
  sensors.setProfile(data)
}


// 自定事件全局挂载
wx.track = new SensorsTrack({
  sensors,
  // TODO 神策数据上报 写好之后记得打开
  open: true,
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
      value: val => dayjs(val * 1000).format('YYYY/MM/DD')
    },
    /** 到店时间时间段 */
    time_period: {
      key: 'service_timeslot',
      value: val => dayjs(val * 1000).format('HH:mm:ss')
    },
    /** 门店名称 */
    // store_name: 'store_name',
    /** 门店编码 */
    // store_code: 'store_code',

    /** 顾客类型 */
    customer_status: {
      key: 'customer_type',
      value: ['', '新客', '老客']
    },
    /** 是否到店购买 */
    is_store: {
      key: 'is_instore',
      value: ['', '是', '否']
    },
    /** 订单金额 */
    reserve_order_price: 'order_amount',
    /** 修改项 */
    change_item: 'modified_item',
    /** 修改后服务类型 */
    modify_reserve_name: 'modified_service_type',
    /** 修改后到店时间日期 */
    modify_reserve_time: {
      key: 'modified_service_date',
      value: val => dayjs(val * 1000).format('YYYY/MM/DD')
    },
    /** 修改后到店时间时间段 */
    modify_time_period: {
      key: 'modified_service_timeslot',
      value: val => dayjs(val * 1000).format('HH:mm:ss')
    },
    /** 取消类型 */
    cancel_type: {
      key: 'cancellation_type',
      value: ['', '手动']
    },
    /** 邀请函ID */
    invite_record_id: {
      key: 'invite_id',
      value: (val) => val + ''
    },
    /** 邀请函模版ID */
    invite_template_id: {
      key: 'invite_template_id',
      value: (val) => val + ''
    },
    /** 入口 */
    entrance: {
      key: 'entry_point',
      value: ['', '企微侧边栏', '顾客详情', '首页邀约榜', '我的到店邀约记录(导购)', '我的邀约记录(顾客)', '到店邀约']
    },
    /** 任务ID */
    task_id: {
      key: 'task_id',
      value: (val) => val + ''
    },
    /** 任务名称 */
    task_name: 'task_name',
    /** 任务覆盖顾客数 */
    customer_total: 'task_coverage',
    /** 任务发布时间 */
    start_time: 'task_starttime',
    /** 素材ID */
    source_id: {
      key: 'asset_id',
      value: (val) => val + ''
    },
    /** 播种ID */
    sow_task_id: {
      key: 'reminder_id',
      value: (val) => val + ''
    },
    /** 播种内容 */
    sow_task_content: 'reminder_content',
    template_name: 'template_name',
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
      value: ([start, end]) => `${dayjs(start * 1000).format('YYYYMMDD')}-${dayjs(end * 1000).format('YYYYMMDD')}`
    },
    item_name: {
      key: 'feature',
      value: {
        'pages/invite-config/invite-config': '邀请函',
        'pages/store-invite/create/create': '到店邀约',
        'pages/online-order/online-order': '线上订单',
        'pages/customer-detail/customer-detail': '客户画像'
      }
    },
    /** 分配ID */
    inherit_id: {
      key: 'assignment_id',
      value: (val) => val + ''
    },
    /** 分配前导购ID */
    handover_userid: 'handover_userid',
    /** 分配后导购ID */
    takeover_userid: {
      key: 'takeover_userid',
      value: val => val + ''
    },
    /** 顾客筛选 */
    manage_type: 'by_customertype',
    /** 员工筛选 */
    member_list: {
      key: 'by_saname',
      value: val => val.join()
    },
    /** 会员卡等 */
    card_type: {
      key: 'by_viplevel',
      value: val => val.join()
    },
    /** 印象标签 */
    impression_tag_list: {
      key: 'impression_tag',
      value: val => val.join()
    },
    /** 客户标签 */
    tag_list: {
      key: 'customer_tag',
      value: val => val.join()
    },
    /** CLT类型 */
    clt_type: {
      key: 'clt_type',
      value: val => val.join()
    },
    /** 搜索原始词 */
    original_word: 'search_key_word',
    /** 评价结果详情 */
    evaluation_result: 'review_details',
    /** 活动安排星级 */
    activity_star: 'activity_rating',
    /** 整体服务星级 */
    service_star: 'service_rating',
    AppointmentModify: {
      /** 邀约ID */
      reserve_number: 'appointment_id',
      /** 顾客ID */
      user_external_id: 'customer_id',
      /** 修改项 */
      change_item: 'modified_item',
      /** 修改后服务类型 */
      modify_reserve_name: 'modified_service_type',
      /** 修改后到店时间日期 */
      modify_reserve_time: {
        key: 'modified_service_date',
        // value: val => val
        value: val => dayjs(val * 1000).format('YYYY/MM/DD')
      },
      /** 修改后到店时间时间段 */
      modify_time_period: {
        key: 'modified_service_timeslot',
        // value: val => val
        value: val => dayjs(val * 1000).format('HH:mm:ss')
      }
      /** 门店名称 */
      // store_name: 'store_name',
      /** 门店编码 */
      // store_code: 'store_code'
    },
    /** 分配顾客 */
    AssignCustomer: {
      /** 分配ID */
      inherit_id: {
        key: 'assignment_id',
        value: (val) => val + ''
      },
      /** 顾客ID */
      user_external_id: 'customer_id',
      /** 分配前导购ID */
      handover_userid: 'handover_userid',
      /** 分配后导购ID */
      takeover_userid: {
        key: 'takeover_userid',
        value: val => val + ''
      }
    },
    /** 顾客被分配 */
    CustomerAssigned: {
      /** 分配ID */
      inherit_id: {
        key: 'assignment_id',
        value: (val) => val + ''
      },
      /** 顾客ID */
      user_external_id: 'customer_id',
      /** 分配前导购ID */
      handover_userid: 'handover_userid',
      /** 分配后导购ID */
      takeover_userid: {
        key: 'takeover_userid',
        value: val => val + ''
      }
    },
    /** 顾客筛选 */
    Filter: {
      current_page: {
        key: 'current_page',
        value: [
          '好友数据', // 0
          '社群数据', // 1
          '播种计划', // 2
          '邀约计划', // 3
          '顾客列表', // 4
          '顾客分配', // 5
          '创建群发', // 6
          '任务群发', // 7
          '子任务详情-群发顾客', // 8
          '子任务详情-单独沟通', // 9
          '任务列表' // 10
        ]
      },
      button_name: 'button_name',
      filter_type_sv: {
        key: 'filter_type_sv',
        value: ['时间', '顾客', '状态', '任务类型', '任务状态', '门店', '导购']
      },
      filter_option: {
        key: 'filter_option',
        value: val => JSON.stringify(val)
      }
    },
    /** 顾客搜索 */
    Search: {
      current_page: {
        key: 'current_page',
        value: {
          0: '任务计划',
          1: '播种计划',
          2: '邀约计划',
          3: '顾客列表',
          4: '子任务详情-群发顾客',
          5: '子任务详情-单独沟通'
        }
      },
      search_type: {
        key: 'search_type',
        value: ['任务', '顾客', '播种', '邀约']
      },
      search_key_word: 'search_key_word'
    },
    /** 顾客详情操作 */
    ClickDetails: {
      user_external_id: 'customer_id',
      operation_type: {
        key: 'operation_type',
        value: [
          '点击发消息', // 0
          '点击消费记录', // 1
          '点击卡券详情', // 2
          '点击会员信息', // 3
          '点击邀约记录', // 4
          '点击播种', // 5
          '点击绑定关系', // 6
          '点击编辑备注', // 7
          '点击编辑客户描述', // 8
          '点击待办任务', // 9
          '点去去完成', // 10
          '点击查看三个月内全部任务', // 11
          '点击拨打电话' // 12
        ]
      }
    },
    /** 上传素材 */
    AssetUpload: {
      asset_id: {
        key: 'asset_id',
        value: val => val + ''
      }
    },
    /** 点击创建群发 */
    BatchMsgCreate: {
      is_filter: 'is_filter',
      entry_point: 'entry_point',
      filter_option: {
        key: 'filter_option',
        value: val => JSON.stringify(val)
      }
    },
    /** 素材群发页操作 */
    BatchMsgOperation: {
      operation_type: {
        key: 'operation_type',
        value: ['仅看好友', '全选', '取消全选', '取消群发', '勾选好友', '勾除好友', '选择群发素材']
      }
    },
    /** 选择群发素材 */
    BatchMsgAssetSelect: {
      operation_type: {
        key: 'operation_type',
        value: ['去素材库选素材', '自定义素材']
      }
    },
    /** 群发素材 */
    BatchMsgConfirm: {
      asset_coverage: 'asset_coverage',
      asset_list: 'asset_list',
      customer_list: 'customer_list',
      operation_type: {
        key: 'operation_type',
        value: ['素材库素材', '自定义素材']
      }
    },
    /** 查看任务 */
    TaskView: {
      /** 任务ID */
      task_id: 'task_id',
      /** 任务状态 */
      task_status: ['待完成', '已完成', '未完成'],
      /** 任务入口 */
      entry_point: ['应用消息通知', '任务中心主任务列表'],
      /** 任务类型 */
      task_type: ['沟通任务', '朋友圈任务']
    },
    /** 沟通子任务详情页操作 */
    TaskDetailOperation: {
      /** 任务ID */
      task_id: 'task_id',
      operation_type: {
        key: 'operation_type',
        value: {
          0: '复制文本',
          1: '查看标签顾客名单',
          2: '查看顾客名单',
          3: '电话',
          4: '发消息',
          5: '写备注',
          6: '加好友',
          61: '勾选',
          62: '取消勾选',
          71: '全选',
          72: '取消全选',
          8: '移入单独沟通',
          9: '转发给顾客',
          10: '标记已联系',
          11: '移回群发顾客'
        }
      }
    },
    /** 沟通任务发送 */
    CommunicationSend: {
      task_id: 'task_id',
      customer_number: 'customer_number'
    },
    /** 会员查询 */
    InquireMember: {
      /** 会员查询入口 */
      entry_point: (value) => ['', '我的页面', '顾客列表'][value],
      /** 会员查询条件 */
      inquiry_condition: 'inquiry_condition'
    },
    /** 会员查询结果 */
    InquireMemberResult: {
      /** 会员查询入口 */
      entry_point: (value) => ['', '我的页面', '顾客列表'][value],
      /** 会员查询结果个数 */
      total: 'result_number'
    },
    /** 查看订单记录 */
    ViewOrderHistory: {
      entry_point: (value) => ({
        1: '顾客详情',
        2: '到店邀约记录'
      }[value])
    },
    /** 进入素材库 */
    AssetLibrary: {
      entry_point: (value) => ({
        1: '工作台',
        2: '1V1聊天窗口',
        3: '顾客列表',
        4: '库存查询',
        5: '图片识别'
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
      entry_point: (value) => ({
        1: '工作台',
        2: '1V1聊天窗口',
        3: '顾客列表',
        4: '库存查询',
        5: '图片识别'
      }[value]),
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
        3: '素材库',
        5: '图片识别'
      }[value]),
      inquiry_type: (value) => ({
        0: 'skuitem',
        1: '条形码查询',
        2: '商品库存查询'
      }[value]),
      search_type: 'search_type',
      is_have_historical_search_word: 'is_have_historical_search_word'
    },
    /** 跨店库存查询 */
    OtherStockInquiry: {
      product_id: 'product_id',
      product_name: 'product_name'
    },
    /** 查看邀约记录 */
    ViewAppointmentHistory: {
      entry_point: (value) => ({
        3: '首页邀约榜',
        4: '我的到店邀约记录(导购)',
        5: '我的邀约记录(顾客)',
        6: '到店邀约',
        7: '顾客详情点击邀约记录'
      }[value])
    },
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
        7: '朋友圈素材',
        8: '官网商品'
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
          3: '顾客列表群发素材',
          4: '库存查询',
          5: '图片识别'
        }[val])
      },
      entry_point: {
        key: 'entry_point',
        value: (val) => ({
          1: '工作台',
          2: '1V1聊天窗口',
          3: '顾客列表',
          4: '库存查询',
          5: '图片识别'
        }[val])
      },
      current_page: {
        key: 'current_page',
        value: (val) => ({
          'pagesMaterial/pages/list/list': '素材库',
          'pagesMaterial/pages/send-list/send-list': '素材库',
          'pagesMaterial/pages/select/select': '素材库',
          'pagesMaterial/pages/collect/collect': '我的收藏'
        }[val])
      },
      is_favorite: 'is_favorite'
    },
    /** 保存素材 */
    AssetSave: {
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
          3: '顾客列表群发素材',
          4: '库存查询',
          5: '图片识别'
        }[val])
      },
      entry_point: {
        key: 'entry_point',
        value: (val) => ({
          1: '工作台',
          2: '1V1聊天窗口',
          3: '顾客列表',
          4: '库存查询',
          5: '图片识别'
        }[val])
      },
      current_page: {
        key: 'current_page',
        value: (val) => ({
          'pagesMaterial/pages/list/list': '素材库',
          'pagesMaterial/pages/send-list/send-list': '素材库',
          'pagesMaterial/pages/select/select': '素材库',
          'pagesMaterial/pages/collect/collect': '我的收藏'
        }[val])
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
          3: '顾客列表群发素材',
          4: '库存查询',
          5: '图片识别'
        }[val])
      },
      entry_point: {
        key: 'entry_point',
        value: (val) => ({
          1: '工作台',
          2: '1V1聊天窗口',
          3: '顾客列表',
          4: '库存查询',
          5: '图片识别'
        }[val])
      },
      current_page: {
        key: 'current_page',
        value: (val) => ({
          'pagesMaterial/pages/list/list': '素材库',
          'pagesMaterial/pages/send-list/send-list': '素材库',
          'pagesMaterial/pages/collect/collect': '我的收藏',
          'pagesMaterial/pages/select/select': '素材库'
        }[val])
      },
      is_favorite: 'is_favorite'
    },
    AssetFavorite: {
      current_page: {
        key: 'current_page',
        value: (val) => ({
          'pagesMaterial/pages/list/list': '素材库',
          'pagesMaterial/pages/send-list/send-list': '素材库',
          'pagesMaterial/pages/collect/collect': '我的收藏'
        }[val])
      },
      is_top: 'is_top',
      entry_point: {
        key: 'entry_point',
        value: (val) => ({
          1: '工作台',
          2: '1V1聊天窗口',
          3: '顾客列表',
          4: '库存查询',
          5: '图片识别'
        }[val])
      },
      library_type: {
        key: 'library_type',
        value: (val) => ({
          1: '品牌素材库',
          2: '素材库',
          3: '顾客列表群发素材',
          4: '库存查询',
          5: '图片识别'
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
    CustomerOperation: {
      current_page: 'current_page',
      operation_type: 'operation_type'
    },
    CommunicateOperation: {
      current_page: 'current_page',
      operation_type: 'operation_type'
    },
    TaskListOperation: {
      current_page: 'current_page',
      operation_type: 'operation_type'
    },
    FilterDetails: {
      current_page: {
        key: 'current_page',
        value: [
          '好友数据', // 0
          '社群数据', // 1
          '播种计划', // 2
          '邀约计划', // 3
          '顾客列表', // 4
          '顾客分配', // 5
          '创建群发', // 6
          '任务群发', // 7
          '子任务详情-群发顾客', // 8
          '子任务详情-单独沟通', // 9
          '任务列表' // 10
        ]
      },
      button_name: 'button_name',
      filter_type_sv: {
        key: 'filter_type_sv',
        value: ['时间', '顾客', '状态', '任务类型', '任务状态', '门店', '导购']
      },
      tag_level_1: 'tag_level_1',
      tag_level_2: 'tag_level_2',
      tag_value: 'tag_value'
    },
    ViewLabel: {
      customer_id: 'customer_id',
      label_entrance: {
        key: 'label_entrance',
        value: (val) => {
          return {
            '-1': '其他',
            'sidebar': '1V1聊天框-客户画像',
            'sidebar-1': '1V1聊天框-播种',
            'pages/customer-list/customer-list': '顾客列表',
            'packages/allocation/pages/customer-select-list/customer-select-list': '顾客分配-选择顾客',
            'packages/allocation/pages/customer-list/customer-list': '顾客分配',
            'pages/consume-record/consume-record': '消费记录',
            'packages/customer/pages/create-group-post-beta/create-group-post-beta': '创建群发',
            'packages/customer/pages/create-group-post/create-group-post': '创建群发',
            'pages/coupon-list/coupon-list': '邀约详情',
            'pages/online-order/online-order': '线上订单',
            'packages/customer/pages/member-search-result/member-search-result': '会员查询结果',
            'packages/task/pages/task-customer/task-customer': '子任务顾客列表',
            'packages/customer/pages/communication-list/index': '沟通名单列表'
          }[val]
        }
      }
    },
    CommodityClick: {
      customer_id: 'customer_id',
      commodity_name: 'commodity_name',
      product_UPC: 'product_UPC',
      product_code: 'product_code',
      mkt_rank: 'mkt_rank',
      mkt_id: 'mkt_id'
    }
  }
})
