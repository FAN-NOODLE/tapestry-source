package transfer_customer

import "micro-scrm/internal/orm/model"

type TransferCustomerWelcomeLanguageModel struct {
	model.BranchModel

	Content string `gorm:"column:content;NOT NULL" json:"content"` // 欢迎语内容

}

func (TransferCustomerWelcomeLanguageModel) TableName() string {
	return "transfer_customer_welcome_language"
}
