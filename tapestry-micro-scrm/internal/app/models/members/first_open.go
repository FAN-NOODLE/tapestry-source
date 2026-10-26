package membersModel

import "micro-scrm/internal/app/models"

const IsFirstOpen = 1

type FirstOpenModel struct {
	models.MicroBaseModel
	MarkType int    `json:"type" gorm:"column:type;type:tinyint(1);default:0;comment:业务类型 1 三方库存查询复制按钮"`
	Mobile   string `json:"mobile" gorm:"column:mobile;type:varchar(50);comment:导购电话"`
}

// TableName returns the table name of the BrandConfigSftp model
func (FirstOpenModel) TableName() string {
	return "first_open"
}
