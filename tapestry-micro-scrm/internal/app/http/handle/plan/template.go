package plan

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	"micro-scrm/internal/app/http/handle"
	planReq "micro-scrm/internal/app/http/request/plan"
	planService "micro-scrm/internal/app/service/plan"
)

type TemplateHandle struct {
	handle.BaseHandle
}

func NewTemplateHandle() TemplateHandle {
	return TemplateHandle{
		BaseHandle: handle.BaseHandle{},
	}
}

func (t TemplateHandle) List(c *gin.Context) {
	var request planReq.TemplateListsReq
	if err := c.ShouldBind(&request); err != nil {
		t.ResError(c, err)
		return
	}
	svc := planService.NewPlanTemplateSvc()
	result, total, err := svc.Lists(t.WarpContext(c), request)
	if err != nil {
		t.ResError(c, err)
		return
	}
	t.ResList(c, total, request.Page, request.Size, result)
}

// Create 任务模板的创建
func (t TemplateHandle) Create(c *gin.Context) {
	var request planReq.TemplateAddReq
	if err := c.ShouldBind(&request); err != nil {
		t.ResError(c, err)
		return
	}

	svc := planService.NewPlanTemplateSvc()
	err := svc.Create(t.WarpContext(c), request)
	if err != nil {
		t.ResError(c, err)
		return
	}
	t.ResOk(c)
}

// Update 任务模板的创建
func (t TemplateHandle) Update(c *gin.Context) {
	var request planReq.TemplateUpdateReq
	if err := c.ShouldBind(&request); err != nil {
		t.ResError(c, err)
		return
	}

	svc := planService.NewPlanTemplateSvc()
	err := svc.Update(t.WarpContext(c), request)
	if err != nil {
		t.ResError(c, err)
		return
	}
	t.ResOk(c)
}

func (t TemplateHandle) Status(c *gin.Context) {
	var request struct {
		Id           int64 `json:"id"`
		EnterpriseID int64 `json:"enterprise_id"`
		Status       int   `json:"status"`
	}
	if err := c.ShouldBind(&request); err != nil {
		t.ResError(c, err)
		return
	}
	// 任务模板的禁用状态，需要判断任务
	if request.Status != 1 {
		count, err2 := microScrmDao.NewPlanTaskDao().GetPlanCountByTemplateId(t.WarpContext(c), request.EnterpriseID, request.Id)
		if err2 != nil {
			return
		}
		if count > 0 {
			t.ResError(c, errorx.New("模板关联的任务正在进行中，暂不可禁用"))
			return
		}
	}
	err := microScrmDao.NewPlanTemplateDao().UpdateData(t.WarpContext(c), request.EnterpriseID,
		map[string]any{"id": request.Id},
		map[string]any{"status": request.Status})
	if err != nil {
		t.ResError(c, err)
		return
	}
	t.ResOk(c)
	return
}

func (t TemplateHandle) Details(c *gin.Context) {
	var request planReq.TemplateDetails
	if err := c.ShouldBind(&request); err != nil {
		t.ResError(c, err)
		return
	}
	svc := planService.NewPlanTemplateSvc()

	result, err := svc.Details(t.WarpContext(c), request)
	if err != nil {
		t.ResError(c, err)
		return
	}
	t.ResSuccess(c, result)
	return
}
