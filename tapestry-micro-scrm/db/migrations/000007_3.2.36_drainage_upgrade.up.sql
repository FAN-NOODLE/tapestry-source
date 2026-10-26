ALTER TABLE drainage_welcomes ADD COLUMN `channel_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '渠道id,标识企业微信state' AFTER is_default;
