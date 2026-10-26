CREATE TABLE `customers_del_message` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '开启开关',
  `send_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '类别',
  `stat_txt` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '设置员工字段json',
  `qy_stat_txt` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '设置员工字段json 企业微信相关',
  `created_at` int(10) NOT NULL DEFAULT '0',
  `updated_at` int(10) NOT NULL DEFAULT '0',
  `deleted_at` int(10) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='删除客户消息配置表';

CREATE TABLE `work_external_user` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL COMMENT 'user主表id',
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业id',
  `union_id` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '用户unionid',
  `external_userid` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '外部联系人userid',
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '用户名称',
  `position` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '用户职称',
  `avatar` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '头像',
  `corp_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '企业简称',
  `corp_full_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '企业全称',
  `type` int(10) NOT NULL COMMENT '外部联系人的类型，1表示该外部联系人是微信用户，2表示该外部联系人是企业微信用户',
  `gender` int(10) NOT NULL COMMENT '外部联系人性别 0-未知 1-男性 2-女性',
  `external_profile` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '外部联系人的自定义展示信息,JSON',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_id_index` (`enterprise_id`),
  KEY `user_id_index` (`user_id`),
  KEY `union_id_index` (`union_id`),
  KEY `enterprise_id_external_userid_index` (`enterprise_id`,`external_userid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='外部联系人信息表';

CREATE TABLE `work_external_user_relation` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业id',
  `work_external_user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业外部联系人信息表id',
  `member_id` bigint(11) NOT NULL DEFAULT '0' COMMENT '成员表id',
  `remark` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '成员设置的外部联系人备注',
  `description` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '成员设置的外部联系人描述',
  `createtime` int(11) NOT NULL DEFAULT '0' COMMENT '关系建立时间',
  `remark_corp_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '成员备注的外部联系人企业名称',
  `remark_mobiles` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '成员备注的外部联系人手机号',
  `add_way` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '该成员添加此外部联系人的来源',
  `state` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '企业自定义的state参数',
  `oper_userid` varchar(255) NOT NULL DEFAULT '0' COMMENT '发起添加的userid，如果成员主动添加，为成员的userid；如果是客户主动添加，则为客户的外部联系人userid；如果是内部成员共享/管理员分配，则为对应的成员/管理员userid',
  `relation_status` int(11) NOT NULL DEFAULT '0' COMMENT '外部联系人关系状态：0-正常，1-员工删除联系人，2-联系人删除员工，3-双方互相删除，4-关系转移中，5-成员离职关系解除',
  `member_deleted_time` int(11) NOT NULL DEFAULT '0' COMMENT '成员主动删除关系时间',
  `external_deleted_time` int(11) NOT NULL DEFAULT '0' COMMENT '外部联系人删除关系时间',
  `created_at` int(11) NOT NULL DEFAULT '0',
  `updated_at` int(11) NOT NULL DEFAULT '0',
  `deleted_at` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `enterprise_id_index` (`enterprise_id`),
  KEY `work_external_user_id_index` (`work_external_user_id`),
  KEY `external_deleted_time_index` (`external_deleted_time`),
  KEY `member_deleted_time_index` (`member_deleted_time`),
  KEY `relation_status_index` (`relation_status`),
  KEY `createtime_index` (`createtime`),
  KEY `member_id_index` (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='成员和外部联系人关系表';

CREATE TABLE `relation_callback_records` (
    `id` bigint(20) NOT NULL,
    `to_user_name` varchar (32)  NOT NULL DEFAULT '' COMMENT '企业ID',
    `from_user_name` varchar (32)  NOT NULL DEFAULT '' COMMENT '此事件该值固定为sys，表示该消息由系统生成',
    `create_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '消息创建时间 （整型）',
    `msg_type` varchar (32)  NOT NULL DEFAULT '' COMMENT '消息的类型，此时固定为event',
    `event` varchar (32)  NOT NULL DEFAULT '' COMMENT '事件的类型',
    `change_type` varchar (32)  NOT NULL DEFAULT '' COMMENT '此时固定为add_external_contact',
    `user_id` varchar (32)  NOT NULL DEFAULT '' COMMENT '企业服务人员的UserID',
    `external_user_id` varchar (32)  NOT NULL DEFAULT '' COMMENT '外部联系人的userid，注意不是企业成员的帐号',
    `state` varchar (32)  NOT NULL DEFAULT '' COMMENT '添加此用户的「联系我」方式配置的state参数，可用于识别添加此用户的渠道',
    `welcome_code` varchar (128)  NOT NULL DEFAULT '' COMMENT '欢迎语welcome code',
    `fail_reason` varchar (32)  NOT NULL DEFAULT '' COMMENT '失败理由',
    `chat_id` varchar (32)  NOT NULL DEFAULT '' COMMENT '群ID',
    PRIMARY KEY (`id`),
    KEY `to_user_name_index` (`to_user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='企业微信外部联系人回调消息记录表';

