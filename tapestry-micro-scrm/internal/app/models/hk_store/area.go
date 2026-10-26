package hkStoreModel

import "micro-scrm/internal/app/models"

type AreaModel struct {
	models.BaseModel
	Type       int         `json:"type" gorm:"column:type;type:int(10);not null;default:0;comment:1 HK/MO 4 store"`
	Name       string      `json:"name" gorm:"column:name;type:varchar(100);not null;comment:name"`
	ParentID   int         `json:"parent_id" gorm:"column:parent_id;type:int(10);not null;default:0;comment:parent id"`
	Path       string      `json:"path" gorm:"column:path;type:varchar(255);not null;comment:path"`
	ThirdID    int         `json:"third_id" gorm:"column:third_id;type:int(10);not null;default:0;comment:wechat department id"`
	IsVirtual  int         `json:"is_virtual" gorm:"column:is_virtual;type:int(10);not null;default:0;comment:is virtual"`
	FiscalYear string      `json:"fiscal_year" gorm:"column:fiscal_year;type:varchar(255);not null;comment:fiscal year"`
	ClosedAt   int64       `json:"closed_at" gorm:"column:closed_at;type:int(10);not null;default:0;comment:closed time"`
	Sort       int64       `json:"sort" gorm:"column:sort;type:int(10);not null;default:0;comment:sort"`
	Children   []AreaModel `json:"children" gorm:"foreignKey:ParentID"`
}

func (AreaModel) TableName() string {
	return "hkmo_area"
}
