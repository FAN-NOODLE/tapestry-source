package reserve

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	reserveRequest "micro-scrm/internal/app/http/request/reserve"
	reserveService "micro-scrm/internal/app/service/reserve"
	"sort"
)

type ReserveServiceHandle struct {
	handle.BaseHandle
	reserveServiceService *reserveService.ReserveServiceService
}

func NewReserveServiceHandle() ReserveServiceHandle {
	return ReserveServiceHandle{
		reserveServiceService: reserveService.NewReserveServiceService(),
	}
}

func (s ReserveServiceHandle) Create(c *gin.Context) {

	var req reserveRequest.ReserveServiceCreateReq
	if err := c.ShouldBind(&req); err != nil {
		s.ResError(c, err)
		return
	}

	err := s.reserveServiceService.Create(s.WarpContext(c), req.EnterpriseId, req)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResOk(c)

}

func (s ReserveServiceHandle) Update(c *gin.Context) {

	var req reserveRequest.ReserveServiceUpdateReq
	if err := c.ShouldBind(&req); err != nil {
		s.ResError(c, err)
		return
	}

	err := s.reserveServiceService.Update(s.WarpContext(c), req.EnterpriseId, req)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResOk(c)

}

func (s ReserveServiceHandle) Delete(c *gin.Context) {

	var req reserveRequest.ReserveServiceDeleteReq
	if err := c.ShouldBind(&req); err != nil {
		s.ResError(c, err)
		return
	}

	err := s.reserveServiceService.Delete(s.WarpContext(c), req.EnterpriseId, req)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResOk(c)

}

func (s ReserveServiceHandle) List(c *gin.Context) {

	var request reserveRequest.ReserveServiceListReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	result, total, err := s.reserveServiceService.GetList(s.WarpContext(c), request.EnterpriseId, request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResList(c, total, request.Page, request.Size, result)
	return
}

// CheckDelete 检查预约服务是否可以删除
func (s ReserveServiceHandle) CheckDelete(c *gin.Context) {
	var request reserveRequest.ReserveServiceDeleteReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}
	status, err := s.reserveServiceService.CheckDelete(s.WarpContext(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, map[string]interface{}{
		"status": status,
	})

}

// ListsByMember 导购获取服务配置列表
func (s ReserveServiceHandle) ListsByMember(c *gin.Context) {
	var request reserveRequest.ReserveServiceListByMemberReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	result, err := s.reserveServiceService.GetListByMember(s.WarpContext(c), request.EnterpriseId, request.Mobile)
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt > result[j].CreatedAt
	})
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResList(c, 0, 1, 10, result)
	return
}

func (s ReserveServiceHandle) Info(c *gin.Context) {
	var request reserveRequest.ReserveServiceDeleteReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}
	info, err := s.reserveServiceService.ServiceInfo(s.WarpContext(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, info)

}
