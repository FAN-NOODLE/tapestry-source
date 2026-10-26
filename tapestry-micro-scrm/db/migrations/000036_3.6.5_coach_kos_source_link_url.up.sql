ALTER TABLE `source_link`
  MODIFY COLUMN `link_url` varchar(1000) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '链接URL',
  MODIFY COLUMN `cover_url` varchar(1000) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '封面图URL';
