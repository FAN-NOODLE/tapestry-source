package source

import "micro-scrm/internal/app/models"

type BaozunCollectModel struct {
	models.MicroBaseModel
	ProductCode string `gorm:"column:product_code;type:varchar(255);default:''" json:"product_code"`
	CollectNum  int    `gorm:"column:collect_num;default:0" json:"collect_num"`
}

func (BaozunCollectModel) TableName() string {
	return "source_baozun_collect"
}
