CREATE TABLE `chat_member` (
  `id` bigint(20) unsigned NOT NULL,
  `member_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT 'members:ID',
  `last_chat_at` int(10) NOT NULL DEFAULT '0' COMMENT '最后聊天时间',
  `enterprise_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '企业ID',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_id_index` (`enterprise_id`) USING BTREE,
  KEY `member_id_index` (`member_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天员工表';


CREATE TABLE `chat_msg` (
  `id` bigint(20) unsigned NOT NULL,
  `msg_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '企微消息ID',
  `action` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '消息动作',
  `from_user` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发送用户',
  `from_user_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '发送用户渠道 1:员工,2:机器人,3:企业微信外部联系人,4:微信外部联系人,5:同企业',
  `to_list` text COLLATE utf8mb4_unicode_ci COMMENT '消息接收用户列表',
  `room_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '群聊消息的群ID',
  `msg_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '消息发送时间',
  `msg_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '消息发送时间ms',
  `type` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '类型',
  `content` json DEFAULT NULL COMMENT '内容',
  `seq` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '企业微信指针',
  `enterprise_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '企业ID',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_id_index` (`enterprise_id`) USING BTREE,
  KEY `from_user_index` (`from_user`) USING BTREE,
  KEY `type_index` (`type`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天消息表';


CREATE TABLE `chat_receive_user` (
  `id` bigint(20) unsigned NOT NULL DEFAULT '0',
  `chat_msg_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '消息表ID',
  `to_user` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '接收用户企微userid',
  `user_channel` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '接收用户渠道 1:员工,2:机器人,3:企业微信外部联系人,4:微信外部联系人,5:同企业',
  `enterprise_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '企业ID',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_id_index` (`enterprise_id`) USING BTREE,
  KEY `chat_msg_id_index` (`chat_msg_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='消息接收用户表';


CREATE TABLE `chat_remind_config` (
  `id` bigint(20) unsigned NOT NULL,
  `overtime` tinyint(2) NOT NULL DEFAULT '0' COMMENT '聊天记录超时 (单位):分钟',
  `receive_start` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '09:00:00' COMMENT '接待开始时间',
  `receive_end` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '18:00:00' COMMENT '接待结束时间',
  `notice_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '提醒开关类型 1-超时 2-敏感词',
  `notice_status` tinyint(1) NOT NULL DEFAULT '2' COMMENT '提醒开关 1-开启 2-关闭',
  `stat_txt` text COLLATE utf8mb4_unicode_ci COMMENT '设置员工字段 员工id(members:id)',
  `enterprise_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '企业ID',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_id_index` (`enterprise_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会话提醒设置';

CREATE TABLE `chat_sensitive` (
  `id` bigint(20) unsigned NOT NULL,
  `sensitive_text` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '敏感词',
  `sensitive_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '敏感词类型 1 名片 2 文本',
  `trigger_num` int(10) NOT NULL DEFAULT '0' COMMENT '总触发次数',
  `enterprise_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '企业ID',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_id_index` (`enterprise_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='敏感词';

CREATE TABLE `chat_settings` (
  `id` bigint(20) unsigned NOT NULL,
  `secret` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '会话存档secret',
  `message_public_key` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '消息加密公钥',
  `message_private_key` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '消息加密私钥',
  `limit` int(10) unsigned NOT NULL DEFAULT '100' COMMENT '一次拉多少条消息',
  `enterprise_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '企业ID',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_id_index` (`enterprise_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会话存档配置';

CREATE TABLE `chat_sensitive_records` (
  `id` bigint(20) unsigned NOT NULL,
  `chat_msg_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '消息id',
  `sensitive_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '敏感词id',
  `enterprise_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '企业ID',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='敏感词触发记录表';


CREATE TABLE `chat_overtime_records` (
  `id` bigint(20) unsigned NOT NULL,
  `chat_msg_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '消息id',
  `enterprise_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '企业ID',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_id_index` (`enterprise_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='员工会话超时记录表';


CREATE TABLE `chat_object` (
  `id` bigint(20) unsigned NOT NULL,
  `from_userid` varchar(64) NOT NULL DEFAULT '' COMMENT '发送方userid',
  `to_userid` varchar(64) NOT NULL DEFAULT '' COMMENT '接收方userid',
  `from_object_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '发送对象ID',
  `to_object_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '接收对象ID',
  `from_channel` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '发送用户渠道 1:员工,2:机器人,3:企业微信外部联系人,4:微信外部联系人,5:群',
  `to_channel` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '接收用户渠道 1:员工,2:机器人,3:企业微信外部联系人,4:微信外部联系人,5:群',
  `last_chat_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '最后聊天时间',
  `agree_status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '同意状态 1:默认同意 2:同意 3:不同意',
  `agree_status_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '同意状态变更时间',
  `is_reply` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否回复,1:是2:否',
  `enterprise_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '企业ID',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_id_index` (`enterprise_id`) USING BTREE,
  KEY `from_object_id_from_channel_index` (`from_object_id`,`from_channel`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天对象表';

CREATE TABLE `chat_msg_media` (
  `id` bigint(20) unsigned NOT NULL,
  `chat_msg_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '聊天记录ID',
  `sdk_file_id` text COMMENT '资源文件ID',
  `path` varchar(255) NOT NULL DEFAULT '' COMMENT '资源路径',
  `enterprise_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '企业ID',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_id_index` (`enterprise_id`) USING BTREE,
  KEY `chat_msg_id_index` (`chat_msg_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天记录关联资源表';