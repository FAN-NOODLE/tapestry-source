// 更新渠道表
ALTER TABLE marketing_strategy_channels MODIFY COLUMN channel_type tinyint(1) NOT NULL DEFAULT 1  COMMENT '渠道；1=短信，2=企微群发，3=自动发券，4=公众号';

// 新增内容部分
ALTER TABLE marketing_strategy_channels ADD COLUMN `content` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文本内容' AFTER `junior_id`;

ALTER TABLE marketing_strategy_records ADD COLUMN  `external_userid` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '外部联系人userid' AFTER `union_id`;

ALTER TABLE marketing_strategy_records ADD COLUMN `msg_id` varchar(64)  COLLATE utf8mb4_unicode_ci  NOT NULL DEFAULT '' COMMENT '群发返回msgid' AFTER `external_userid`;


// 企业微信群发记录表
CREATE TABLE `qy_msg_send_log` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业ID',
  `strategy_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '策略id',
  `channel_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '渠道id',
  `chat_type` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'single' COMMENT '群发任务类型 single:客户 group:客户群',
  `external_userid` text COLLATE utf8mb4_unicode_ci COMMENT '外部联系人userid集合',
  `sender` varchar(64) NOT NULL DEFAULT '' COMMENT '发送企业群发消息的成员userid',
  `msg_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '企业群发消息的id',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_id_index` (`enterprise_id`) USING BTREE,
  KEY `strategy_id_index` (`strategy_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='企业微信群发记录表';

// 企业微信群发结果表
CREATE TABLE `qy_msg_send_result` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业ID',
  `send_log_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '发送记录ID',
  `external_userid` varchar(64) NOT NULL DEFAULT '' COMMENT '外部联系人userid',
  `userid` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发送人',
  `chat_id` varchar(64) NOT NULL DEFAULT '' COMMENT '外部客户群id',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '发送状态（0-未发送 1-已发送 2-因客户不是好友导致发送失败 3-因客户已经收到其他群发消息导致发送失败）',
  `msg_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '企业群发消息的id',
  `send_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '发送时间',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `msg_id_index` (`msg_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='企业微信群发结果表';