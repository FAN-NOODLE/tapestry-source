package transfer_customer

import "gorm.io/gorm"

type TransferCustomerService struct {
	*gorm.DB
	ServiceName string
}

// NewTransferCustomerService 创建服务
func NewTransferCustomerService(g *gorm.DB) *TransferCustomerService {
	return &TransferCustomerService{
		DB:          g,
		ServiceName: "transfer_customer_service",
	}
}
