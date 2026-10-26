package dashboardSvc

import (
	"micro-scrm/internal/app/service"
)

type ReserveService struct {
	service.Base
}

func NewReserveService() *ReserveService {
	return &ReserveService{}
}

func (s ReserveService) DataKanban() (any, error) {
	return nil, nil
}

func (s ReserveService) DataList() (any, int64, error) {
	return nil, 0, nil
}
