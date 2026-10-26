CREATE TABLE `drainage_welcomes` (
  `id` bigint(20) NOT NULL,
  `content` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '内容',,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业ID',
  `is_default` int(1) NOT NULL DEFAULT '0' COMMENT '是否是默认欢迎语 1是 0不是',
  `created_at` int(10) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_id` (`enterprise_id`) USING BTREE,
  KEY `created_at` (`created_at`) USING BTREE,
  KEY `updated_at` (`updated_at`) USING BTREE,
  KEY `deleted_at` (`deleted_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='引流模块欢迎语表';


CREATE TABLE `drainage_welcome_members` (
  `id` bigint(20) NOT NULL,
  `welcome_id` bigint(20) NOT NULL COMMENT '欢迎语表主键',
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业ID',
  `member_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '员工ID',
  `created_at` int(10) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `welcome_id` (`welcome_id`) USING BTREE,
  KEY `enterprise_id` (`enterprise_id`) USING BTREE,
  KEY `member_id` (`member_id`) USING BTREE,
  KEY `created_at` (`created_at`) USING BTREE,
  KEY `updated_at` (`updated_at`) USING BTREE,
  KEY `deleted_at` (`deleted_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='引流模块欢迎语关联员工表';

CREATE TABLE `drainage_contacts` (
  `id` bigint(20) NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '二维码名称',
  `channel_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '渠道id,标识企业微信state',
  `auto_friend` int(1) NOT NULL DEFAULT '0' COMMENT '是否自动成为好友 1 开启 0 不开启',
  `qr_code_url` varchar(1000) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '二维码地址',
  `member_count` int(10) NOT NULL DEFAULT '0' COMMENT '创建员工数',
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业ID',
  `sum_add_member_number` int(10) NOT NULL DEFAULT '0' COMMENT '累加添加人数',
  `created_at` int(10) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) NOT NULL DEFAULT '0' COMMENT '删除时间',
  `config_id` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '微信返回配置id',
  PRIMARY KEY (`id`),
  KEY `name` (`name`) USING BTREE,
  KEY `enterprise_id` (`enterprise_id`) USING BTREE,
  KEY `created_at` (`created_at`) USING BTREE,
  KEY `updated_at` (`updated_at`) USING BTREE,
  KEY `deleted_at` (`deleted_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='引流模块联系我表';

CREATE TABLE `drainage_contact_members` (
  `id` bigint(20) NOT NULL,
  `contact_id` bigint(20) NOT NULL COMMENT '联系我表主键',
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业ID',
  `member_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '员工ID',
  `created_at` int(10) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `contact_id` (`contact_id`) USING BTREE,
  KEY `enterprise_id` (`enterprise_id`) USING BTREE,
  KEY `member_id` (`member_id`) USING BTREE,
  KEY `created_at` (`created_at`) USING BTREE,
  KEY `updated_at` (`updated_at`) USING BTREE,
  KEY `deleted_at` (`deleted_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='引流模块联系我关联员工表';

CREATE TABLE `drainage_contact_channels` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL COMMENT '企业ID',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '渠道名称',
  `created_at` int(10) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `id` (`id`) USING BTREE,
  KEY `name` (`name`) USING BTREE,
  KEY `enterprise_id` (`enterprise_id`) USING BTREE,
  KEY `created_at` (`created_at`) USING BTREE,
  KEY `updated_at` (`updated_at`) USING BTREE,
  KEY `deleted_at` (`deleted_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='引流模块联系我渠道表';

CREATE TABLE `drainage_contact_channel_labels` (
  `id` bigint(20) NOT NULL,
  `channel_id` bigint(20) NOT NULL COMMENT '渠道id',
  `label_group_id` bigint(20) NOT NULL COMMENT '标签分组id',
  `label_id` bigint(20) NOT NULL COMMENT '标签id',
  `enterprise_id` bigint(20) NOT NULL COMMENT '企业ID',
  `wx_tag_group_id` varchar(100) NOT NULL COMMENT '企微标签分组ID',
  `wx_tag_id` varchar(100) NOT NULL COMMENT '企微标签ID',
  `created_at` int(10) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `channel_id` (`channel_id`) USING BTREE,
  KEY `label_group_id` (`label_group_id`) USING BTREE,
  KEY `enterprise_id` (`enterprise_id`) USING BTREE,
  KEY `created_at` (`created_at`) USING BTREE,
  KEY `updated_at` (`updated_at`) USING BTREE,
  KEY `deleted_at` (`deleted_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='引流模块联系我渠道标签表';

CREATE TABLE `work_external_user_relation_log` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业id',
  `work_external_user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业外部联系人信息表id',
  `member_id` bigint(11) NOT NULL DEFAULT '0' COMMENT '成员表id',
  `state` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '企业自定义的state参数',
  `relation_status` int(11) NOT NULL DEFAULT '0' COMMENT '外部联系人关系状态：0-正常，1-员工删除联系人，2-联系人删除员工，3-双方互相删除，4-关系转移中，5-员工离职关系解除',
  `created_at` int(11) NOT NULL DEFAULT '0',
  `updated_at` int(11) NOT NULL DEFAULT '0',
  `deleted_at` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `member_id_index` (`member_id`) USING BTREE,
  KEY `work_user_id_index` (`work_external_user_id`) USING BTREE,
  KEY `enterprise_id_index` (`enterprise_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='成员和外部联系人关系日志表';