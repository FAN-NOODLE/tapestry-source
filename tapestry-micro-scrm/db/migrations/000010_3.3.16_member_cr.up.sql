// 商城会员表
CREATE TABLE `shop_members` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0',
  `union_id` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户id',
  `number` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '编号',
  `level_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '等级id',
  `available_integral_total` int(11) NOT NULL DEFAULT '0' COMMENT '可用积分总数量',
  `year_integral_total` int(11) NOT NULL DEFAULT '0' COMMENT '本年积分总数量',
  `history_integral_total` int(11) NOT NULL DEFAULT '0' COMMENT '历史积分总数量',
  `qr_code_url` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '二维码链接',
  `bar_code_url` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '条形码链接',
  `open_time` int(11) NOT NULL DEFAULT '0' COMMENT '当前等级开通时间',
  `overdue_time` int(11) NOT NULL DEFAULT '0' COMMENT '当前等级过期时间',
  `created_at` int(10) NOT NULL DEFAULT '0',
  `updated_at` int(10) NOT NULL DEFAULT '0',
  `deleted_at` int(10) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `enterprise_id_index` (`enterprise_id`),
  KEY `union_id_index` (`union_id`),
  KEY `level_id_index` (`level_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商城会员表';


// 会员等级表
CREATE TABLE `shop_member_level` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0',
  `level_id` bigint(20) NOT NULL COMMENT '等级ID',
  `level` tinyint(1) NOT NULL DEFAULT '1' COMMENT '系统内置1-10级',
  `level_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '等级名称',
  `created_at` int(10) NOT NULL DEFAULT '0',
  `updated_at` int(10) NOT NULL DEFAULT '0',
  `deleted_at` int(10) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `enterprise_id_index` (`enterprise_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员等级表';