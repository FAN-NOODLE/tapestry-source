
CREATE TABLE `chat_group_msg_task` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0',
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '群发任务名称',
  `type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '群发类型1-立即 2-定时',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '任务触达状态：1-定时未开始 2-未发送 3-部分未发送 4-已发送',
  `task_status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否已执行（1：未执行,2:执行中，3：已执行,4:执行失败,5:执行失败结束）默认1',
  `fail_msg` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '执行失败原因',
  `send_time` int(10) NOT NULL DEFAULT '0' COMMENT '定时发送时间 默认为0',
  `last_send_msg` int(10) NOT NULL DEFAULT '0' COMMENT '最新发送消息时间 默认为0',
  `latest_check_time` int(10) NOT NULL DEFAULT '0' COMMENT '最近的检查时间',
  `content` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '内容',
  `created_at` int(10) NOT NULL DEFAULT '0',
  `updated_at` int(10) NOT NULL DEFAULT '0',
  `deleted_at` int(10) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `enterprise_id_index` (`enterprise_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='群发任务表';

CREATE TABLE `chat_group_msg_task_member` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0',
  `msg_task_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '群发任务ID',
  `member_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '成员ID',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '发送状态：0默认(等待调用企微接口发送群发任务) 1-未发送(任务已下发给群主,等待群主发送群发消息) 2-已发送（群主已将任务发送至企业群）',
  `task_status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否已执行（1：未执行，2：已执行,3:执行失败）默认1',
  `fail_msg` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '执行失败原因',
  `msgid` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '企微msgid',
  `created_at` int(10) NOT NULL DEFAULT '0',
  `updated_at` int(10) NOT NULL DEFAULT '0',
  `deleted_at` int(10) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `enterprise_id_index` (`deleted_at`, `enterprise_id`, `msg_task_id`, `status`),
  KEY `member_id_index` (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='群发任务会员表';

CREATE TABLE `chat_group_msg_send_record` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0',
  `msg_task_id` bigint(20) NOT NULL COMMENT '群发任务ID',
  `msg_task_member_id` bigint(20) NOT NULL DEFAULT '0' COMMENT 'chat_group_msg_task_member表id',
  `chat_id` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '企微外部群id',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '发送状态 0-未发送 1-已发送 2-因客户不是好友导致发送失败 3-因客户已经收到其他群发消息导致发送失败',
  `send_time` int(10) NOT NULL DEFAULT '0' COMMENT '发送时间',
  `created_at` int(10) NOT NULL DEFAULT '0',
  `updated_at` int(10) NOT NULL DEFAULT '0',
  `deleted_at` int(10) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `enterprise_id_index` (`enterprise_id`, `msg_task_id`),
  KEY `msg_task_member_id` (`msg_task_member_id`),
  KEY `chat_id_index` (`chat_id`),
  KEY `send_time_index` (`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='群发结果记录表';