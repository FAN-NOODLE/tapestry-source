package department

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/request/department"
	cacheSvc "micro-scrm/internal/app/service/cache"
	departmentService "micro-scrm/internal/app/service/department"
	"micro-scrm/internal/pkg/sdk/boss"
)

type Handle struct {
	handle.BaseHandle
	departmentsService departmentService.DepartmentsService
}

func NewHandle() Handle {
	return Handle{
		departmentsService: departmentService.NewDepartmentsService(),
	}
}

func (h Handle) SaStoreList(c *gin.Context) {
	var request department.SAStoreListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	// 判断企业
	config, err := cacheSvc.NewRedisData().GetEnterpriseConfig(h.WarpContext(c), request.EnterpriseID, "store_name")
	if err != nil {
		h.ResError(c, err)
		return
	}
	isCoach := false
	if config.Name == "coach" {
		isCoach = true
	}

	total, result, err := h.departmentsService.GetSaStoreLists(h.WarpContext(c), request, isCoach)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, 0, 0, result)
}

func (h Handle) StoreList(c *gin.Context) {
	var request department.SAStoreListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	// 获取企业配置的门店父级部门id信息
	configData, err := boss.EnterpriseConfig(h.WarpContext(c), map[string]interface{}{
		"enterprise_id": request.EnterpriseID,
		"key":           "store_root_dept_id",
	})
	if err != nil || len(configData.Data) == 0 {
		h.ResError(c, err)
		return
	}
	storeParentThirdId := int(configData.Data[0].Value) //所有门店的父级企微部门id

	storeList, err := h.departmentsService.GetAllStore(h.WarpContext(c), request.EnterpriseID, storeParentThirdId)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, storeList)
}
