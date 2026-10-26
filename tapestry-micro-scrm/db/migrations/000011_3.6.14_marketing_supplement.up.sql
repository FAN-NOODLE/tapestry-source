// 商城用户动作表
CREATE TABLE `shop_user_action` (
  `id` bigint(20) NOT NULL,
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0',
  `union_id` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户union_id',
  `open_id`  varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户open_id',
  `action_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '行为类型:1-浏览,2-加购,3-下单,4-支付',
  `is_touch` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已触发',
  `created_at` int(10) NOT NULL DEFAULT '0',
  `updated_at` int(10) NOT NULL DEFAULT '0',
  `deleted_at` int(10) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `enterprise_id_index` (`enterprise_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商城用户动作表';