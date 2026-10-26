package callbackSvc

import (
	"micro-scrm/internal/app/service"
)

type Callback struct {
	service.Base
}

func NewCallback() *Callback {
	return &Callback{}
}

func (c Callback) name() {

}
