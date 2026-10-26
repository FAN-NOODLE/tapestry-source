CREATE TABLE `customer_relation_guide` (
  `id` bigint(20) unsigned NOT NULL COMMENT '主键ID',
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业ID',
  `members_user_id` varchar(255) NOT NULL DEFAULT '' COMMENT '企业微信唯一标识(members表userid一致)',
  `members_id` bigint(20) NOT NULL DEFAULT 0 COMMENT 'members表id',
  `user_id` bigint(20) NOT NULL DEFAULT 0 COMMENT 'user表id',
  `third_user_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '第三方ID',
  `bind_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '绑定时间',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_i
  d_index` (`enterprise_id`) USING BTREE,
  KEY `user_id_index` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='用户绑定导购关系';


CREATE TABLE `wx_coupon` (
  `id` bigint(20) unsigned NOT NULL COMMENT '主键ID',
  `enterprise_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '企业ID',
  `coupon_id` bigint(20) unsigned DEFAULT '0' COMMENT '第三方优惠券id',
  `coupon_stock_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '0' COMMENT '优惠券批次',
  `coupon_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '优惠券名称',
  `coupon_value` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '优惠券面额（单位分）',
  `coupon_mininumn` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '优惠券使用最低额（单位分）',
  `coupon_status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '状态：1-启用，2-禁用',
  `coupon_stock_status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '批次状态: 1-未激活；2-审批中；4-已激活；8-已作废；16-中止发放',
  `coupon_total` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '优惠券数量',
  `use_day_limit` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0代表不限 领取后x天后过期,需要按批次时间计算',
  `max_quota` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '优惠券最大领取数(0：不限制)',
  `is_send_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '优惠券已发送数量',
  `coupon_budget` int(10) unsigned NOT NULL DEFAULT '1' COMMENT '优惠券预算额度',
  `scope` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '优惠范围GLOBAL-全场代金券,SINGLE-单品优惠',
  `favor_type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'COUPON- 代金券，需要走结算资金的充值型代金券,DISCOUNT- 优惠券，不走结算资金的免充值型优惠券',
  `start_time` int(10) unsigned NOT NULL DEFAULT '1' COMMENT '开始生效时间',
  `end_time` int(10) unsigned NOT NULL DEFAULT '1' COMMENT '生效结束时间',
  `is_overlay` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '是否叠加；1-可以叠加，2-不可叠加',
  `overlay_coupons` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '可叠加的优惠券',
  `channel` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '领券渠道：默认1=站内;2=站外',
  `is_all` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '可用范围(1=全场，2=单品)',
  `range` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '可用商品范围集合',
  `range_hide` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '可用商品(不展示)范围集合',
  `tag_goods` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '标签关联的商品',
  `tag_name` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标签名称',
  `tag_rule_type` tinyint(3) NOT NULL DEFAULT '0' COMMENT '标签商品结算规则（1：订单包含所选商品，2：订单包含所选商品，且所选商品金额达到券门槛）',
   `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_enterpriseId_couponStockId_deletedAt` (`enterprise_id`,`coupon_stock_id`,`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='微信优惠券表';


CREATE TABLE `coupon_user` (
  `id` bigint(20) unsigned NOT NULL COMMENT '主键ID',
  `enterprise_id` BIGINT(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '企业ID',
  `user_id` BIGINT(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT 'user表id',
  `union_id` VARCHAR(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '用户unionid',
  `effective` TINYINT(3) UNSIGNED NOT NULL DEFAULT '1' COMMENT '有效时间类型（1:无限制，2:固定时间段）',
  `order_no` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '订单号',
  `start_time` INT(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '有效期开始时间(时间戳-秒)',
  `end_time` INT(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '有效期结束时间(时间戳-秒)',
  `status` TINYINT(3) UNSIGNED NOT NULL DEFAULT '1' COMMENT '使用状态（1:待使用，2:使用中，3:已使用，4:已过期）',
  `coupon_user_id` BIGINT(20) UNSIGNED DEFAULT '0' COMMENT '第三方用户优惠券表ID',
  `coupon_id` BIGINT(20) UNSIGNED DEFAULT '0' COMMENT '第三方优惠券id',
  `source_id` BIGINT(20) UNSIGNED DEFAULT '0' COMMENT '来源ID',
  `source` TINYINT(3) UNSIGNED NOT NULL DEFAULT '1' COMMENT '来源：默认1=用户领取;2=队列同步;3=cdp;4=翻倍券活动;5=助力;6=助力抽奖;7=新拼团;8=新砍价;9=内购;10=外部关联优惠券;11=直播活动',
  `action_id` VARCHAR(60) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '触发发券的对象唯一标识：默认0，点赞id，翻倍任务id，支付订单号',
  `action_type` TINYINT(3) UNSIGNED NOT NULL COMMENT '1点赞、2完成翻倍、3支付成功、4补发、5助力成功、6助力失败、7帮助助力成功、8助力抽奖、9助力抽奖未参与者、10拼团发起或参与成功、11砍价、12直播订阅、13直播分享',
  `use_time` INT(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '优惠券时间使用时间',
  `created_at` INT(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` INT(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` INT(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_id_index` (`enterprise_id`),
  KEY `coupon_user_source_id_source_index` (`source_id`,`source`),
  KEY `coupon_user_user_id_platform_coupon_id_index` (`user_id`,`coupon_id`),
  KEY `idx_status_endTime_platformCouponId_deletedAt` (`status`,`end_time`,`coupon_id`,`deleted_at`),
  KEY `idx_platformCouponId_enterpriseId_deletedAt` (`coupon_id`,`enterprise_id`,`deleted_at`)
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户礼劵信息表';

CREATE TABLE `customer_remark_info` (
  `user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT 'user表id',
  `enterprise_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '企业id',
  `age` int(10) NOT NULL DEFAULT '0' COMMENT '年龄',
  `birthday` int(10) NOT NULL DEFAULT '0' COMMENT '生日',
  `birthday_date` int(10) NOT NULL DEFAULT '0' COMMENT '生日20200908',
  `gender` tinyint(1) NOT NULL DEFAULT '0' COMMENT '性别：0-未知 1-男性 2-女性',
  `mobile` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '电话',
  `email` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '邮件',
  `qq` varchar(16) NOT NULL DEFAULT '' COMMENT 'QQ号',
  `describe` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '描述',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`user_id`),
  KEY `enterprise_id_index` (`enterprise_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='客户详细信息备注表';

CREATE TABLE `shop_activity` (
  `id` bigint(20) unsigned NOT NULL DEFAULT '0',
  `activity_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '商城活动id',
  `source` varchar(60) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '活动数据来源：crs',
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '活动名称',
  `start_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '活动开始时间',
  `end_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '活动结束时间',
  `need_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '需要数量',
  `activity_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '活动类型',
  `activity_status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '活动状态 0:草稿 1：待上架 2：已上架 3：已下架 4：已结束 5：暂停 6：定时上架 7：定时下架',
  `enterprise_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '企业id',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_IDX` (`enterprise_id`,`activity_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商城活动规则主表';

CREATE TABLE `shop_activity_sku` (
  `id` bigint(20) unsigned NOT NULL,
  `shop_activity_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '商城活动主表ID',
  `enterprise_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '企业ID',
  `spu_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '商品spu_id',
  `sku_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '商品sku_id',
  `sku_status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态 1：启用 2：禁用',
  `activity_stock_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '库存类型 1：使用活动库存 2：使用商品库存',
  `activity_stock_set` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '活动设置库存',
  `sku_picture` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'sku图像',
  `sell_price` int(11) NOT NULL DEFAULT '0' COMMENT '生成活动时候的商品销售价格',
  `activity_price` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '活动价格',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `shop_activity_id_idx` (`shop_activity_id`),
  KEY `enterprise_id_idx` (`enterprise_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品活动关联表';


CREATE TABLE `goods_spu` (
  `id` bigint(20) unsigned NOT NULL,
  `spu_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT 'spu_id',
  `source_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '门店Id',
  `spu_name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '商品名称',
  `spu_sn` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '商品编号',
  `short_title` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '商品短标题',
  `sell_point` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '卖点',
  `sales` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '销量',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '商品状态 1待上架 2上架中 3已下架',
  `category_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '分类ID',
  `category_name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '分类名称',
  `updated_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商品更新时间',
  `enterprise_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '企业ID',
  `picture` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'spu图像',
  `default_attribute` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'SPU商品默认属性值集合',
  `create_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'spu创建时间',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `enterprise_id-index` (`enterprise_id`),
  KEY `spu_id_index` (`spu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品SPU表';


CREATE TABLE `goods_sku` (
  `id` bigint(20) unsigned NOT NULL,
  `sku_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT 'SKU ID',
  `spu_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT 'SPU ID',
  `source_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '门店Id',
  `material_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '物料ID',
  `sku_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'sku名称',
  `card_title` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '分享卡片标题',
  `poster_title` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '海报标题',
  `sell_point` varchar(225) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '卖点',
  `short_title` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '商品短标题',
  `sales` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '销量',
  `stock` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商品库存',
  `guide_price` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '总部指导价',
  `sign_price` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商品标牌价（单位：分）',
  `sell_price` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商品售卖价（单位：分）',
  `sale_type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '商品售卖类型（1:快递 2:自提 3:快递+自提）',
  `qrcode_url` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '二维码地址',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '商品状态 1待上架 2上架中 3已下架',
  `goods_weight` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '商品重量',
  `goods_volume` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '商品体积',
  `enterprise_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '企业ID',
  `sku_attribute` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'SKU商品属性值集合',
  `sku_picture` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'sku图像',
  `create_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'sku创建时间',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `spu_id-index` (`spu_id`),
  KEY `enid_sid_del_index` (`enterprise_id`,`source_id`,`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品SKU表';

ALTER TABLE `source_image` ADD COLUMN `media_id` varchar(128) NOT NULL DEFAULT '' COMMENT '企业微信临时素材ID' AFTER  `ext`;
ALTER TABLE `source_image` ADD COLUMN `media_expire_time` int(10) NOT NULL DEFAULT '' COMMENT '企业微信临时素材过期时间' AFTER  `media_id`;
ALTER TABLE `source_image` ADD COLUMN `media_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '企业微信临时素材更新状态0=未更新 1=更新成功 2=更新失败' AFTER  `media_expire_time`;
ALTER TABLE `source_image` ADD COLUMN `media_fail_reason` varchar(255) NOT NULL DEFAULT '' COMMENT '企业微信临时素材更新错误原因' AFTER  `media_status`;


ALTER TABLE `source_video` ADD COLUMN `media_id` varchar(128) NOT NULL DEFAULT '' COMMENT '企业微信临时素材ID' AFTER  `ext`;
ALTER TABLE `source_video` ADD COLUMN `media_expire_time` int(10) NOT NULL DEFAULT '' COMMENT '企业微信临时素材过期时间' AFTER  `media_id`;
ALTER TABLE `source_video` ADD COLUMN `media_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '企业微信临时素材更新状态0=未更新 1=更新成功 2=更新失败' AFTER  `media_expire_time`;
ALTER TABLE `source_video` ADD COLUMN `media_fail_reason` varchar(255) NOT NULL DEFAULT '' COMMENT '企业微信临时素材更新错误原因' AFTER  `media_status`;

ALTER TABLE `source_file` ADD COLUMN `media_id` varchar(128) NOT NULL DEFAULT '' COMMENT '企业微信临时素材ID' AFTER  `ext`;
ALTER TABLE `source_file` ADD COLUMN `media_expire_time` int(10) NOT NULL DEFAULT '' COMMENT '企业微信临时素材过期时间' AFTER  `media_id`;
ALTER TABLE `source_file` ADD COLUMN `media_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '企业微信临时素材更新状态0=未更新 1=更新成功 2=更新失败' AFTER  `media_expire_time`;
ALTER TABLE `source_file` ADD COLUMN `media_fail_reason` varchar(255) NOT NULL DEFAULT '' COMMENT '企业微信临时素材更新错误原因' AFTER  `media_status`;


ALTER TABLE `old_user` ADD COLUMN `sync` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '是否已同步　0=否　１=是';




ALTER TABLE `coupon_user` ADD COLUMN `open_id` varchar(255) utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '用户openid' AFTER  `enterprise_id`;


