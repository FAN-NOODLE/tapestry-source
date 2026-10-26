package reserve

import (
	"micro-scrm/internal/app/http/handle"
	reserveRequest "micro-scrm/internal/app/http/request/reserve"
	reserveService "micro-scrm/internal/app/service/reserve"

	"github.com/gin-gonic/gin"
)

type ServiceScoreHandle struct {
	handle.BaseHandle
	serviceScoreService *reserveService.ServiceScoreService
}

func NewServiceScoreHandle() ServiceScoreHandle {
	return ServiceScoreHandle{
		serviceScoreService: reserveService.NewServiceScoreService(),
	}
}

// List 用户评分列表
func (s ServiceScoreHandle) List(c *gin.Context) {

	var request reserveRequest.ServiceScoreListReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	result, total, err := s.serviceScoreService.GetList(s.WarpContext(c), request.EnterpriseId, request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResList(c, total, request.Page, request.Size, result)
	return
}
