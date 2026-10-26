CREATE TABLE IF NOT EXISTS `inventory_sku_upc_mapping` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'id',
  `batch_no` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'import batch number',
  `sku_item` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'SKU item',
  `upc_code` VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'UPC code',
  `effective_month` CHAR(7) NOT NULL DEFAULT '' COMMENT 'YYYY-MM',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT 'created timestamp',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT 'updated timestamp',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT 'deleted timestamp',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_batch_sku_upc` (`batch_no`, `sku_item`, `upc_code`),
  KEY `idx_sku_active` (`sku_item`, `deleted_at`),
  KEY `idx_upc_active` (`upc_code`, `deleted_at`),
  KEY `idx_batch` (`batch_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='inventory SKU UPC mapping';

CREATE TABLE IF NOT EXISTS `inventory_sku_upc_import_batch` (
  `id` BIGINT NOT NULL COMMENT 'id',
  `batch_no` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'import batch number',
  `effective_month` CHAR(7) NOT NULL DEFAULT '' COMMENT 'YYYY-MM',
  `source_file` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'source CSV file name',
  `total_rows` INT NOT NULL DEFAULT 0 COMMENT 'total rows',
  `success_rows` INT NOT NULL DEFAULT 0 COMMENT 'success rows',
  `failed_rows` INT NOT NULL DEFAULT 0 COMMENT 'failed rows',
  `duplicate_rows` INT NOT NULL DEFAULT 0 COMMENT 'duplicate rows',
  `unique_sku_count` INT NOT NULL DEFAULT 0 COMMENT 'unique SKU count',
  `unique_upc_count` INT NOT NULL DEFAULT 0 COMMENT 'unique UPC count',
  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0 processing, 1 success, 2 failed',
  `is_active` TINYINT NOT NULL DEFAULT 0 COMMENT '1 current active batch',
  `mapping_enabled` TINYINT NOT NULL DEFAULT 1 COMMENT '1 use DB mapping, 0 bypass DB',
  `error_summary` TEXT NULL COMMENT 'error summary',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT 'created timestamp',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT 'updated timestamp',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT 'deleted timestamp',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_batch_no` (`batch_no`),
  KEY `idx_active_enabled` (`is_active`, `mapping_enabled`, `status`, `deleted_at`),
  KEY `idx_effective_month` (`effective_month`, `status`, `deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='inventory SKU UPC import batch';
