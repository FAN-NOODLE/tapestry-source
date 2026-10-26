package microScrmModel

import "micro-scrm/internal/app/models"

type InventorySearchModel struct {
	models.MicroBaseModel
	MemberID    int64  `gorm:"column:member_id;type:bigint(20);default:0;comment:成员ID" json:"member_id" `
	KeywordType string `gorm:"column:keyword_type;type:varchar(255);default:''" json:"keyword_type"`
	Keywords    string `gorm:"column:keywords;type:varchar(255);default:''" json:"keywords"`
}

func (InventorySearchModel) TableName() string {
	return "inventory_search"
}

const (
	InventorySearchTypeSkuItem = 1
	InventorySearchTypeUpcCode = 2
	InventorySearchProductName = 3
)

const LimitNum = 15
