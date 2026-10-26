CREATE TABLE `fission_activity` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业id',
  `name` varchar(120) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '活动名称',
  `desc` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '活动描述',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '活动状态：1-上架中，2-已下架',
  `channel_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '渠道id',
  `start_time` int(10) NOT NULL DEFAULT '0' COMMENT '活动开始时间',
  `end_time` int(10) NOT NULL DEFAULT '0' COMMENT '活动结束时间',
  `exchange_end_time` int(10) NOT NULL DEFAULT '0' COMMENT '兑换奖品截至时间',
  `fission_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '裂变类型：1-添加员工，2-添加员工且入群',
  `cond_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '参与条件：1-添加任意员工，2-添加裂变员工',
  `mutual_type` tinyint(1) NOT NULL DEFAULT '2' COMMENT '互斥类型：1-互斥，2-不互斥',
  `drainage_welcomes_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '欢迎语id',
  `contact_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '活动默认联系我二维码id，drainage_contacts表id',
  `user_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '发起人数',
  `fission_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '好友裂变数',
  `run_away_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '流失人数',
  `join_group_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '群成员裂变数',
  `leave_group_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '群成员退群数',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `deleted_at_enterprise_id_idx` (`enterprise_id`, `deleted_at`),
  KEY `channel_id_index` (`channel_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='裂变活动表';

ALTER TABLE drainage_welcomes
  ADD `type` TINYINT(1) UNSIGNED NOT NULL DEFAULT '1' COMMENT '欢迎语类型：1-系统欢迎语，2-裂变活动欢迎语';

ALTER TABLE work_group_chat
  ADD `member_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '群人数';

CREATE TABLE `fission_assistant_records` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业id',
  `activity_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '活动id',
  `fission_user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '裂变用户关系表id',
  `assistant_user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '助力用户id,work_external_user表user_id',
  `add_member_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '添加好友记录：0-未添加，1-已添加，2-已删除',
  `add_member_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '好友添加时间',
  `join_chat_group_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '入群状态：0-未入群，1-已入群，2-已退群',
  `join_chat_group_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '入群时间',
  `score` int(10) NOT NULL DEFAULT '0' COMMENT '计分',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `deleted_enterprise_activity_index` (`deleted_at`,`enterprise_id`,`activity_id`),
  KEY `fission_user_id_index` (`fission_user_id`),
  KEY `assistant_user_member_idx` (`assistant_user_id`,`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='裂变助力记录表';


CREATE TABLE `fission_chat_group` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业id',
  `activity_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '活动id',
  `chat_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '群聊id',
  `qr_code` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '入群二维码',
  `sort` int(10) NOT NULL DEFAULT '0' COMMENT '排序权重，升序',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_activity_deleted_at_index` (`enterprise_id`,`activity_id`,`deleted_at`),
  KEY `chat_id_idx` (`chat_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='裂变活动社群表';


CREATE TABLE `fission_exchange_prize_records` (
  `id` bigint(20) NOT NULL DEFAULT '0',
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业id',
  `activity_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '活动id',
  `prize_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '奖品id',
  `fission_user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '裂变活动用户关系表id',
  `user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户id',
  `num` int(10) NOT NULL DEFAULT '0' COMMENT '兑换数量',
  `cost_assistant_num` int(10) NOT NULL DEFAULT '0' COMMENT '消耗的助力人数',
  `deliver_name` varchar(60) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '收件人姓名',
  `deliver_mobile` varchar(60) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '收件人手机号',
  `deliver_address` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '收件人地址',
  `deliver_area` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '收货地址省市区',
  `deliver_area_code` varchar(60) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '0' COMMENT '收货地址省市区编码',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_activity_index` (`enterprise_id`,`activity_id`,`deleted_at`),
  KEY `fission_user_id_index` (`fission_user_id`),
  KEY `user_id_idx` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='裂变活动奖品兑换记录';

CREATE TABLE `fission_members` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业id',
  `activity_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '活动id',
  `member_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '员工id',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_activity_deleted_member_index` (`enterprise_id`,`activity_id`,`deleted_at`,`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='裂变活动员工表';

CREATE TABLE `fission_mutual_relation` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业id',
  `activity_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '活动id',
  `mutual_activity_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '互斥的活动id',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_activity_deleted_index` (`enterprise_id`,`activity_id`,`deleted_at`),
  KEY `mutual_activity_index` (`mutual_activity_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='互动互斥关系表';

CREATE TABLE `fission_prize` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业id',
  `activity_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '活动id',
  `name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '奖品名称',
  `pic` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '奖品图片链接',
  `stocks` int(10) NOT NULL DEFAULT '0' COMMENT '奖品库存',
  `exchange_num` int(10) NOT NULL DEFAULT '0' COMMENT '已兑换奖品数',
  `need_num` int(10) NOT NULL DEFAULT '0' COMMENT '需要的裂变人数',
  `limit_num` int(10) NOT NULL DEFAULT '0' COMMENT '每人限兑换个数',
  `sort` int(10) NOT NULL DEFAULT '0' COMMENT '排序权重：升序',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_activity_deleted_index` (`enterprise_id`,`activity_id`,`deleted_at`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='裂变活动奖品';

CREATE TABLE `fission_user` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业id',
  `activity_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '活动id',
  `user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户id',
  `total_assistant_num` int(10) NOT NULL DEFAULT '0' COMMENT '累计助力人数',
  `cur_assistant_num` int(10) NOT NULL DEFAULT '0' COMMENT '当前剩余助力人数',
  `contact_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '联系我二维码id，drainage_contacts表id',
  `fission_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '好友裂变数',
  `run_away_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '流失人数',
  `join_group_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '群成员裂变数',
  `leave_group_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '群成员退群数',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_activity_user_idx` (`enterprise_id`,`activity_id`,`user_id`,`deleted_at`),
  KEY `contact_id_idx` (`contact_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='参与裂变活动用户表';

ALTER TABLE `drainage_contacts`
  ADD `type` TINYINT(1) UNSIGNED NOT NULL DEFAULT '1' COMMENT '欢迎语类型：1-系统联系我，2-裂变活动默认二维码,3-裂变活动发起人二维码';

CREATE TABLE `enterprise_config` (
    `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业ID',
    `hash` varchar(64) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT 'hash',
    `key` varchar(100) CHARACTER SET utf8 DEFAULT '' COMMENT '配置项',
    `value` varchar(1024) CHARACTER SET utf8 NOT NULL COMMENT '配置值',
    `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
    `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
    `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
    UNIQUE KEY `idx_enterprise` (`enterprise_id`,`key`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='企业信息配置表';

CREATE TABLE `fission_mp_user` (
    `id` bigint(20) unsigned NOT NULL,
    `open_id` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `union_id` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `nick_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '昵称',
    `sex` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '值为1时是男性，值为2时是女性，值为0时是未知',
    `province` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '省份',
    `city` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '城市',
    `country` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '国家',
    `img_url` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '头像',
    `privilege` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '用户特权信息',
    `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业id',
    `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
    `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
    `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `unq_open_id_enterprise` (`open_id`,`enterprise_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


ALTER TABLE fission_assistant_records
    ADD `member_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '添加的员工ID' AFTER assistant_user_id ;

CREATE TABLE `robot_user` (
    `id` bigint(20) unsigned NOT NULL,
    `nick_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '昵称',
    `sex` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '值为1时是男性，值为2时是女性，值为0时是未知',
    `img_url` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '头像',
    `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
    `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
    `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统假用户';


ALTER TABLE fission_assistant_records
    ADD `external_user_id` varchar(255) NOT NULL DEFAULT '' COMMENT '外部联系人userid' AFTER assistant_user_id ;

INSERT INTO `enterprise_config` (`enterprise_id`, `hash`, `key`, `value`, `created_at`
                                , `updated_at`, `deleted_at`)
SELECT `enterprise_id`, `hash`, 'official_info' AS `key`
     , '{"app_id":"wxcb73917918005dfe","secret":"REDACTED_FOR_SOURCE_SHARING"}' AS `value`, unix_timestamp(now()) AS created_at, unix_timestamp(now()) AS updated_at
     , 0 AS deleted_at
FROM enterprise_corp
WHERE enterprise_id NOT IN (
    SELECT enterprise_id
    FROM enterprise_config
    WHERE `key` = 'official_info'
      AND `deleted_at` != 0
);

ALTER TABLE work_group_chat
 MODIFY `notice` TEXT COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '群公告';