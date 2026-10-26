CREATE TABLE `work_group_chat` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0',
  `chat_id` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '企业微信群聊id',
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '群聊名称',
  `type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '群类型 1:外部群 2内部群',
  `owner_member_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '群创建人id',
  `notice` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `create_time` int(10) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `created_at` int(10) NOT NULL DEFAULT '0',
  `updated_at` int(10) NOT NULL DEFAULT '0',
  `deleted_at` int(10) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `enterprise_id_chat_id_index` (`enterprise_id`,`chat_id`(191)),
  KEY `owner_member_id_index` (`owner_member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='企业微信群聊表';
