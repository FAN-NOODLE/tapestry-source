package microScrmModel

import "micro-scrm/internal/app/models"

const (
	InventorySkuUpcBatchStatusProcessing = 0
	InventorySkuUpcBatchStatusSuccess    = 1
	InventorySkuUpcBatchStatusFailed     = 2

	InventorySkuUpcBatchInactive = 0
	InventorySkuUpcBatchActive   = 1

	InventorySkuUpcMappingDisabled = 0
	InventorySkuUpcMappingEnabled  = 1
)

type InventorySkuUpcImportBatchModel struct {
	models.BaseModel
	BatchNo        string `json:"batch_no" gorm:"column:batch_no;type:varchar(64);not null;default:'';comment:import batch number"`
	EffectiveMonth string `json:"effective_month" gorm:"column:effective_month;type:char(7);not null;default:'';comment:YYYY-MM"`
	SourceFile     string `json:"source_file" gorm:"column:source_file;type:varchar(255);not null;default:'';comment:source CSV file name"`
	TotalRows      int    `json:"total_rows" gorm:"column:total_rows;type:int;not null;default:0;comment:total rows"`
	SuccessRows    int    `json:"success_rows" gorm:"column:success_rows;type:int;not null;default:0;comment:success rows"`
	FailedRows     int    `json:"failed_rows" gorm:"column:failed_rows;type:int;not null;default:0;comment:failed rows"`
	DuplicateRows  int    `json:"duplicate_rows" gorm:"column:duplicate_rows;type:int;not null;default:0;comment:duplicate rows"`
	UniqueSkuCount int    `json:"unique_sku_count" gorm:"column:unique_sku_count;type:int;not null;default:0;comment:unique SKU count"`
	UniqueUpcCount int    `json:"unique_upc_count" gorm:"column:unique_upc_count;type:int;not null;default:0;comment:unique UPC count"`
	Status         int    `json:"status" gorm:"column:status;type:tinyint;not null;default:0;comment:0 processing, 1 success, 2 failed"`
	IsActive       int    `json:"is_active" gorm:"column:is_active;type:tinyint;not null;default:0;comment:1 current active batch"`
	MappingEnabled int    `json:"mapping_enabled" gorm:"column:mapping_enabled;type:tinyint;not null;default:1;comment:1 use DB mapping, 0 bypass DB"`
	ErrorSummary   string `json:"error_summary" gorm:"column:error_summary;type:text;comment:error summary"`
}

func (InventorySkuUpcImportBatchModel) TableName() string {
	return "inventory_sku_upc_import_batch"
}
