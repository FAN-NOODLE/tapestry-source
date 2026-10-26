package reserve

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
)

type OrderHandle struct {
	handle.BaseHandle
}

func NewOrderHandle() OrderHandle {
	return OrderHandle{}
}

func (h OrderHandle) CreateOrder(c *gin.Context) {

}

func (h OrderHandle) OrderInfo(c *gin.Context) {

}
