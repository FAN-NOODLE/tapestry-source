package user_event

import (
	"micro-scrm/internal/orm/model/micro_scrm"
)

func (f *UserEventService) CreateModel(model micro_scrm.UserEventModel) error {
	return f.Create(&model).Error
}
