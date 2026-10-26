CREATE TABLE `sms_sign` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业ID',
  `third_sign_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '第三方签名ID',
  `third_sign_status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '签名状态:1:审核中、2:审核失败、3:审核成功',
  `third_sign_type` bigint(20) NOT NULL DEFAULT '0' COMMENT '签名类型，0表示国内短信，1表示国际短信',
  `sign_name` varchar(20) NOT NULL DEFAULT '' COMMENT '签名名称',
  `created_at` int(11) NOT NULL DEFAULT '0' COMMENT '数据创建时间',
  `updated_at` int(11) NOT NULL DEFAULT '0' COMMENT '数据更新时间',
  `deleted_at` int(11) NOT NULL DEFAULT '0' COMMENT '数据删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='短信模板表';

CREATE TABLE `sms_template` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业ID',
  `admin_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '操作人ID',
  `template_id` bigint(1) NOT NULL DEFAULT '0' COMMENT '第三方短信模板ID',
  `template_name` varchar(80) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '模板名称',
  `template_sign_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '模板签名',
  `message_template` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '短信模板',
  `template_content` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '短信插入参数后内容体',
  `template_html` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '前端需要的html',
  `message_str` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '前端列表展示短信内容体',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '审核状态（1：未审核、2：审核中、3：审核失败、4：审核成功、5：已赋值）',
  `fail_reason` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '审核失败原因',
  `type` tinyint(1) NOT NULL DEFAULT '2',
  `created_at` int(11) NOT NULL DEFAULT '0' COMMENT '数据创建时间',
  `updated_at` int(11) NOT NULL DEFAULT '0' COMMENT '数据更新时间',
  `deleted_at` int(11) NOT NULL DEFAULT '0' COMMENT '数据删除时间',
  PRIMARY KEY (`id`),
  KEY `admin_id_index` (`admin_id`) USING BTREE,
  KEY `enterprise_id_index` (`enterprise_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='短信模板表';

CREATE TABLE `sms_send_log` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业ID',
  `user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `sms_template_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '短信模板主键',
  `channel_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '短信触发渠道（1：营销策略）',
  `relation_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '关联ID',
  `channel_value_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '渠道数据Id',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '短信发送状态（1：成功 2：失败）',
  `created_at` int(11) NOT NULL DEFAULT '0' COMMENT '数据创建时间',
  `updated_at` int(11) NOT NULL DEFAULT '0' COMMENT '数据更新时间',
  `deleted_at` int(11) NOT NULL DEFAULT '0' COMMENT '数据删除时间',
  PRIMARY KEY (`id`),
  KEY `sms_send_log_sms_template_id_index` (`sms_template_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='短信发送记录表';


CREATE TABLE `marketing_strategies` (
    `id` bigint(20) NOT NULL,
    `strategy_name` varchar(64) NOT NULL DEFAULT '' COMMENT '策略名称',
    `customer_group_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '客户分组ID',
    `enterprise_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '企业ID',
    `admin_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '操作人ID',
    `strategy_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '策略类型: 1=客户策略，2=内置策略',
    `strategy_method` tinyint(1) NOT NULL DEFAULT 1 COMMENT '执行方式: 1=长期执行，2=立即执行，3=定时执行',
    `start_time` int NOT NULL DEFAULT 0 COMMENT '执行开始时间',
    `end_time` int NOT NULL DEFAULT 0 COMMENT '执行结束时间',
    `strategy_policy` tinyint(1) NOT NULL DEFAULT 1 COMMENT '触达规则：1=通知一个渠道，2=通知所有渠道',
    `strategy_status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态: 1=未开始，2=投放中，3=暂停中，4=已结束',
    `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态: 1=禁用，2=启用',
    `is_finished` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否完成: 1=未完成，2=已完成',
    `setting_day` tinyint(1) NOT NULL DEFAULT 7 COMMENT '设置默认天数',
    `is_notice` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否已通知发送短信提醒',
    `mobile` varchar(20) NOT NULL DEFAULT '' COMMENT '联系方式',
    `created_at` int NOT NULL DEFAULT 0 COMMENT '数据创建时间',
    `updated_at` int NOT NULL DEFAULT 0 COMMENT '数据更新时间',
    `touched_at` int NOT NULL DEFAULT 0 COMMENT '触达时间',
    `deleted_at` int NOT NULL DEFAULT 0 COMMENT '数据删除时间',
    PRIMARY KEY (`id`),
    KEY `customer_group_id_index` (`customer_group_id`) USING BTREE,
    KEY `enterprise_id_index` (`enterprise_id`) USING BTREE,
    KEY `admin_id_index` (`admin_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='营销策略表';

CREATE TABLE `marketing_strategy_channels` (
    `id` bigint(20) NOT NULL,
    `strategy_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '策略id',
    `channel_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '渠道；1=短信',
    `junior_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '从属id',
    `order` int NOT NULL DEFAULT 1 COMMENT '排序:默认1  （此处倒序排列）',
    `created_at` int NOT NULL DEFAULT 0 COMMENT '数据创建时间',
    `updated_at` int NOT NULL DEFAULT 0 COMMENT '数据更新时间',
    `deleted_at` int NOT NULL DEFAULT 0 COMMENT '数据删除时间',
    PRIMARY KEY (`id`),
    KEY `strategy_id_index` (`strategy_id`) USING BTREE,
    KEY `junior_id_index` (`junior_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='策略渠道表';


CREATE TABLE `marketing_strategy_conditions` (
    `strategy_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '策略id',
    `condition` int NOT NULL DEFAULT 1 COMMENT '触发条件;1=不限，2=访问，3=加购，4=下单，5=支付',
    `inner_condition` int NOT NULL DEFAULT 1 COMMENT '内置触发条件;1=不设置条件，2=购买金额，3=购买商品',
    `assign_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '指定类型；1=全部，2=指定',
    `relation_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '关联类型；1=首页，2=商品，3=直播，4=活动页，5=活动',
    `belong_type` tinyint(1) NOT NULL DEFAULT 0 COMMENT '所属类型',
    `start_point` int NOT NULL DEFAULT 0 COMMENT '起始点',
    `end_point` int NOT NULL DEFAULT 0 COMMENT '结束点',
    `time_interval` int NOT NULL DEFAULT 0 COMMENT '间隔时间（之后执行）',
    `time_type` int NOT NULL DEFAULT 1 COMMENT '时间类型1=秒，2=分，3=时，4=天',
    `time` int NOT NULL DEFAULT 0 COMMENT '间隔时间/秒',
    `created_at` int NOT NULL DEFAULT 0 COMMENT '数据创建时间',
    `updated_at` int NOT NULL DEFAULT 0 COMMENT '数据更新时间',
    `deleted_at` int NOT NULL DEFAULT 0 COMMENT '数据删除时间',
    PRIMARY KEY (`strategy_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='策略执行条件表';

CREATE TABLE `marketing_trigger_lists` (
    `id` bigint(20) NOT NULL,
    `list_name` varchar(64) NOT NULL DEFAULT '' COMMENT '列表名称',
    `desc` varchar(100) NOT NULL DEFAULT '' COMMENT '列表描述',
    `created_at` int NOT NULL DEFAULT 0 COMMENT '数据创建时间',
    `updated_at` int NOT NULL DEFAULT 0 COMMENT '数据更新时间',
    `deleted_at` int NOT NULL DEFAULT 0 COMMENT '数据删除时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='策略触发列表';

CREATE TABLE `marketing_strategy_trigger_relations` (
    `id` bigint(20) NOT NULL,
    `strategy_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '策略id',
    `relation_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '关联id',
    `created_at` int NOT NULL DEFAULT 0 COMMENT '数据创建时间',
    `updated_at` int NOT NULL DEFAULT 0 COMMENT '数据更新时间',
    `deleted_at` int NOT NULL DEFAULT 0 COMMENT '数据删除时间',
    PRIMARY KEY (`id`),
    KEY `strategy_id_index` (`strategy_id`) USING BTREE,
    KEY `relation_id_index` (`relation_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='策略关联关系表';

CREATE TABLE `marketing_strategy_records` (
    `id` bigint(20) NOT NULL,
    `user_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '用户id',
    `enterprise_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '企业id',
    `union_id` varchar(255) NOT NULL DEFAULT "" COMMENT 'union_id',
    `strategy_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '策略id',
    `channel_type` bigint(20) NOT NULL DEFAULT 0 COMMENT '渠道类型,1：短信；2：企微群发；3：自动发券；4：公众号；',
    `relation_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '触发模板id',
    `channel_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '渠道id',
    `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '执行状态；1=未执行，2=执行成功，3=执行失败',
    `user_name` varchar(50) NOT NULL DEFAULT "" COMMENT '用户名称',
    `mobile` varchar(20) NOT NULL DEFAULT "" COMMENT '用户电话',
    `address_user_name` varchar(50) NOT NULL DEFAULT "" COMMENT '联系人名称',
    `order_id` varchar(255) NOT NULl DEFAULT "" COMMENT '订单号',
    `start_time` int NOT NULL DEFAULT 0 COMMENT '开始时间',
    `end_time` int NOT NULL DEFAULT 0 COMMENT '结束时间',
    `created_at` int NOT NULL DEFAULT 0 COMMENT '数据创建时间',
    `updated_at` int NOT NULL DEFAULT 0 COMMENT '数据更新时间',
    `touched_at` int NOT NULL DEFAULT 0 COMMENT '触达时间',
    `deleted_at` int NOT NULL DEFAULT 0 COMMENT '数据删除时间',
    PRIMARY KEY (`id`),
    KEY `strategy_id_index` (`strategy_id`) USING BTREE,
    KEY `user_id_index` (`user_id`) USING BTREE,
    KEY `union_id_index` (`union_id`) USING BTREE,
    KEY `relation_id_index` (`relation_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='策略执行记录表';
