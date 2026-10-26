package overtimehttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/overtime"
	"micro-scrm/internal/pkg/errs"
)

var service *overtime.Service

// overtimeList 获取聊天超时列表
func overtimeList(c *gin.Context) {
	var request overtime.ListRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	if request.Page == 0 {
		request.Page = 1
	}
	if request.Size == 0 {
		request.Size = 10
	}
	// 搜索客户姓名
	if request.SearchType == 1 {
		request.CustomerName = request.KeyWords
		request.KeyWords = ""
	}

	count, entity, err := service.SelectList(ctx.Wrap(c), request)
	res := overtime.PageListResponse{}
	if entity != nil {
		res.List = *entity
	} else {
		res.List = []overtime.RecordsListResponse{}
	}
	res.Total = count
	res.Page = request.Page
	res.Size = request.Size

	ginplus.ResSuccess(c, res)
}
