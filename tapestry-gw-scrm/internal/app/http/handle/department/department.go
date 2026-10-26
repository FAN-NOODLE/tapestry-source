package department

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	"gw-scrm/internal/app/http/request/department"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type Handle struct {
	handle.BaseHandle
}

func NewHandle() Handle {
	return Handle{}
}

// StoreLists 获取coach sa的门店列表，奥莱，零售
func (h Handle) StoreLists(c *gin.Context) {
	var request department.SAStoreListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId

	result, err := microscrm.GetSaStoreLists(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
	return
}
