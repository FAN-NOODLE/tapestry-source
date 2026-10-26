ALTER TABLE `source_link_kos`
  ADD COLUMN `source_url` varchar(1000) NOT NULL DEFAULT '' COMMENT '小红书原始帖子地址' AFTER `post_id`;
