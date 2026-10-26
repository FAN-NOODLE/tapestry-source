package http

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/ba_kanban"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/errs"
	response "micro-scrm/internal/pkg/resoponse"
)

type consoleHandle struct {
	Name string
}

func NewConsoleHandle() *consoleHandle {
	return &consoleHandle{
		Name: "管理后台Handle",
	}
}

func (handle *consoleHandle) GetBaKanBanData(c *gin.Context) {
	var request ba_kanban.MiniBaKanBanDataRequest
	bindJonsErr := c.BindJSON(&request)
	if bindJonsErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	if data, dataErr := ba_kanban.NewBaKanBan(request).GetBaKanBanDataConsole(ctx.Wrap(c)); dataErr == nil {
		ginplus.ResSuccess(c, data)
		return
	} else {
		ginplus.ResError(c, dataErr)
	}
}

func (handle *consoleHandle) GetSalesDataDetail(c *gin.Context) {
	var request ba_kanban.MiniBaKanBanDataRequest
	bindJonsErr := c.BindJSON(&request)
	if bindJonsErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	if data, dataErr := ba_kanban.NewBaKanBan(request).GetSalesDataDetail(ctx.Wrap(c)); dataErr == nil {
		ginplus.ResSuccess(c, data)
		return
	} else {
		ginplus.ResError(c, dataErr)
	}
}

func (handle *consoleHandle) GetOperationDataDetail(c *gin.Context) {
	var request ba_kanban.MiniBaKanBanDataRequest
	bindJonsErr := c.BindJSON(&request)
	if bindJonsErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	if data, dataErr := ba_kanban.NewBaKanBan(request).GetOperationDataDetail(ctx.Wrap(c)); dataErr == nil {
		ginplus.ResSuccess(c, data)
		return
	} else {
		ginplus.ResError(c, dataErr)
	}
}

func GetList(c *gin.Context) {
	var request ba_kanban.BaPerformanceQueryReq
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	list, total, err := orm.Service.BaPerformance.GetPageList(request.EnterpriseId, request.IsPage, request.Page, request.Size, request.UserIds, request.Name)
	if err != nil {
		ginplus.ResError(c, errors.New("获取BA指标数据异常"))
		return
	}
	response.SuccessList(c, total, request.Page, request.Size, list)
	return
}

func Create(c *gin.Context) {
	var request ba_kanban.BaPerformanceSavaReq
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	var err error
	//删除之前导入数据，重新插入
	err = orm.Service.BaPerformance.Delete(request.EnterpriseId)
	if err != nil {
		ginplus.ResError(c, errors.New("BA指标数据删除异常"))
		return
	}
	err = orm.Service.BaPerformance.Create(request.EnterpriseId, request.ModelArray)
	if err != nil {
		ginplus.ResError(c, errors.New("BA指标导入异常"))
		return
	}
	err = orm.Service.BaPerformance.AddLog(request.EnterpriseId, request.ModelArray)
	if err != nil {
		ginplus.ResError(c, errors.New("BA指标导入记录异常"))
		return
	}

	ginplus.ResSuccess(c, struct{}{})
	return
}

func baSourceDepte(c *gin.Context) {
	params := ba_kanban.BaSourceDeptRequst{}
	err := c.BindJSON(&params)
	if err != nil {
		ginplus.ResError(c, errors.New("BA手机号不存在"))
		return
	}
	result := ba_kanban.BaSourceDepte(params)
	ginplus.ResSuccess(c, result)
	return
}

func (handle *consoleHandle) GetSalesDataDetailV1(c *gin.Context) {
	params := ba_kanban.BaKanBanDataRequest{}
	bindJonsErr := c.BindJSON(&params)
	if bindJonsErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	data, _ := ba_kanban.GetSalesDataDetailV1(params)
	ginplus.ResSuccess(c, data)
	return
}
func (handle *consoleHandle) GrowthRate(c *gin.Context) {
	params := ba_kanban.BaKanBanDataRequest{}
	bindJonsErr := c.BindJSON(&params)
	if bindJonsErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	data, _ := ba_kanban.GrowthRate(params)
	ginplus.ResSuccess(c, data)
	return
}
func (handle *consoleHandle) GetOperationDataDetailv1(c *gin.Context) {
	params := ba_kanban.BaKanBanDataRequest{}
	bindJonsErr := c.BindJSON(&params)
	if bindJonsErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	data, total := ba_kanban.GetOperationDataDetailv1(params)
	response.SuccessList(c, total, params.Page, params.PageSize, data)
	return

}
func (handle *consoleHandle) baKanbandateupload(c *gin.Context) {
	params := ba_kanban.BaKanDateUploadRequst{}
	bindJonsErr := c.BindJSON(&params)
	if bindJonsErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	if params.EnterpriseId != enterpriseID.(int64) {
		ginplus.ResError(c, errors.New("暂无权限上传该企业的数据"))
		return
	}
	err := ba_kanban.BaKanbandateupload(params)

	ginplus.ResSuccess(c, err)
	return

}
