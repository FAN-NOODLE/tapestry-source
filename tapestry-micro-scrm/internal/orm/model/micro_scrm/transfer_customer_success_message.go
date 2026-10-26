package micro_scrm

import "micro-scrm/internal/orm/model"

type TransferCustomerSuccessMessageModel struct {
	model.MicroScrmBaseModel

	Content string `json:"content" gorm:"column:content;type:text;not null;comment:欢迎语内容"` // 欢迎语内容
}

func (TransferCustomerSuccessMessageModel) TableName() string {
	return "transfer_customer_success_message"
}
