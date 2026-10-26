package modulesSvc

import (
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/modules/customer/respository/elastic"
)

type CustomerService struct {
	service.Base
	es *elastic.EsRepo
}

func NewCustomerService() *CustomerService {
	return &CustomerService{
		es: elastic.NewEsRepo(),
	}
}

func (s CustomerService) EsRepo() *elastic.EsRepo {
	return s.es
}
