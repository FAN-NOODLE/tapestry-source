UPDATE `source_link`
SET `link_url` = LEFT(`link_url`, 255),
    `cover_url` = LEFT(`cover_url`, 255)
WHERE CHAR_LENGTH(`link_url`) > 255 OR CHAR_LENGTH(`cover_url`) > 255;

ALTER TABLE `source_link`
  MODIFY COLUMN `link_url` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '链接URL',
  MODIFY COLUMN `cover_url` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '封面图URL';
