CREATE TABLE `enterprise_corp` (
    `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业ID',
    `corp_id` varchar(64) NOT NULL DEFAULT '' COMMENT 'corp_id',
    `domain` varchar(64) NOT NULL DEFAULT '' COMMENT 'domain',
    `hash` varchar(64) NOT NULL DEFAULT '' COMMENT 'hash',
    `tag_init_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '标签初始化状态 1=成功 2=失败',
    `reason` varchar(255) NOT NULL DEFAULT '' COMMENT '标签初始化结果理由',
    PRIMARY KEY `enterprise_id` (`enterprise_id`),
    KEY `corp_id` (`corp_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='企业信息';

CREATE TABLE `departments` (
    `id` bigint(20) NOT NULL,
    `parent_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '业务中父部门id',
    `third_id` int(10) NOT NULL DEFAULT '0' COMMENT '企业微信中的组织ID',
    `third_parent_id` int(10) NOT NULL DEFAULT '0' COMMENT '企业微信中的组织父ID',
    `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '部门名称',
    `name_en` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '英文名称',
    `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业ID',
    `order` int(10) NOT NULL DEFAULT '0' COMMENT '在父部门中的次序值',
    `created_at` int NOT NULL DEFAULT '0' COMMENT '数据创建时间',
    `updated_at` int NOT NULL DEFAULT '0' COMMENT '数据更新时间',
    `deleted_at` int NOT NULL DEFAULT '0' COMMENT '数据删除时间',
    PRIMARY KEY (`id`),
    KEY `parent_id_index` (`parent_id`) USING BTREE,
    KEY `enterprise_id` (`enterprise_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='企业微信组织结构';

CREATE TABLE `members` (
    `id` bigint(20) NOT NULL,
    `userid` varchar(64) NOT NULL DEFAULT '' COMMENT 'userid',
    `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '企业微信姓名',
    `mobile` varchar(15) NOT NULL DEFAULT '' COMMENT '企业微信手机号',
    `position` varchar(64) NOT NULL DEFAULT '' COMMENT '职务',
    `gender` tinyint(2) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '0' COMMENT '性别0=未知 1=男性 2=女性',
    `email` varchar(128) NOT NULL DEFAULT '' COMMENT '邮箱',
    `avatar` varchar(128) NOT NULL DEFAULT '' COMMENT '员工头像',
    `thumb_avatar` varchar(128) NOT NULL DEFAULT '' COMMENT '缩略图',
    `telephone` varchar(32) NOT NULL DEFAULT '' COMMENT '座机',
    `alias` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '别名',
    `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '激活状态: 0=刚创建 1=已激活，2=已禁用，4=未激活，5=退出企业',
    `qr_code` varchar(128) NOT NULL DEFAULT '' COMMENT '员工个人二维码',
    `address` varchar(255) NOT NULL DEFAULT '' COMMENT '地址',
    `open_userid` varchar(64) NOT NULL DEFAULT 'open_userid',
    `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业ID',
    `created_at` int NOT NULL DEFAULT '0' COMMENT '数据创建时间',
    `updated_at` int NOT NULL DEFAULT '0' COMMENT '数据更新时间',
    `deleted_at` int NOT NULL DEFAULT '0' COMMENT '数据删除时间',
    PRIMARY KEY (`id`),
    KEY `user_id_index` (`userid`) USING BTREE,
    KEY `mobile_index` (`mobile`) USING BTREE,
    KEY `email_index` (`email`) USING BTREE,
    KEY `enterprise_id` (`enterprise_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='企业微信员工';

CREATE TABLE `member_departments` (
    `id` bigint(20) NOT NULL,
    `member_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '成员ID',
    `department_id` bigint(20) NOT NULL DEFAULT  '0' COMMENT '组织ID',
    `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业ID',
    `created_at` int NOT NULL DEFAULT '0' COMMENT '数据创建时间',
    `updated_at` int NOT NULL DEFAULT '0' COMMENT '数据更新时间',
    `deleted_at` int NOT NULL DEFAULT '0' COMMENT '数据删除时间',
    PRIMARY KEY (`id`),
    KEY `department_id_index` (`department_id`) USING BTREE,
    KEY `member_id_index` (`member_id`) USING BTREE,
    KEY `enterprise_id` (`enterprise_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='企业微信成员、组织关系';

CREATE TABLE `department_update_records` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业ID',
  `admin_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '操作人ID',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '更新类型: 1=自动更新，2=手动更新',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态: 0=初始化创建,1=更新失败,2=更新成功,3=推送成功,4=开始异步处理',
  `request_id` varchar(100) NOT NULL DEFAULT '' COMMENT 'CMQ返回请求ID',
  `msg_id` varchar(100) NOT NULL DEFAULT '' COMMENT 'CMQ返回消息ID',
  `created_at` int(11) NOT NULL DEFAULT '0' COMMENT '数据创建时间',
  `updated_at` int(11) NOT NULL DEFAULT '0' COMMENT '数据更新时间',
  `deleted_at` int(11) NOT NULL DEFAULT '0' COMMENT '数据删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `admin_id_index` (`admin_id`) USING BTREE,
  KEY `enterprise_id_index` (`enterprise_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='企业微信组织结构更新记录';