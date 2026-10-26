package http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/ba_kanban"
	"micro-scrm/internal/pkg/errs"
)

type miniHandle struct {
	Name string
}

func NewMiniHandle() *miniHandle {
	return &miniHandle{
		Name: "小程序Handle",
	}
}

func (handle *miniHandle) GetBaKanBanData(c *gin.Context) {
	var request ba_kanban.MiniBaKanBanDataRequest
	bindJonsErr := c.BindJSON(&request)
	if bindJonsErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	if data, dataErr := ba_kanban.NewBaKanBan(request).GetBaKanBanDataMini(ctx.Wrap(c)); dataErr == nil {
		ginplus.ResSuccess(c, data)
		return
	} else {
		ginplus.ResError(c, dataErr)
	}
}

func (handle *miniHandle) GetBaKanBanDatav1(c *gin.Context) {
	params := ba_kanban.MiNiBaKanBanDataRequest{}
	bindJonsErr := c.BindJSON(&params)
	if bindJonsErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	data, _ := ba_kanban.GetBaKanBanDataMiniV1(params)
	ginplus.ResSuccess(c, data)
	return
}
