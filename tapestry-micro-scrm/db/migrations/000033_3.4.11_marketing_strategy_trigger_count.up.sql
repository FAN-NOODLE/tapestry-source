CREATE TABLE `marketing_strategy_trigger_count` (
  `strategy_id` bigint(20) NOT NULL COMMENT '策略id',
  `count` int(10) NOT NULL DEFAULT '0' COMMENT '可触达人数',
  PRIMARY KEY (`strategy_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

ALTER TABLE `marketing_strategies` ADD COLUMN  `history_status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '长期不限存量数据是否已执行，1=执行 0=未执行' AFTER `is_notice`;

