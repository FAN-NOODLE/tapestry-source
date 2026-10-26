package http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/tag"
	"micro-scrm/internal/modules/tag/entity"
	"micro-scrm/internal/pkg/errs"
)

// 小程序端接口

// 获取企业的自定义标签详情
func getTagGroupDetailList(c *gin.Context) {
	// 绑定请求参数
	var request entity.MiniBaseParam
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	result, resErr := tag.MiniTag.GetTagGroupDetail(tag.MiniTag{}, request)
	if resErr != nil {
		ginplus.ResError(c, resErr)
		return
	}

	ginplus.ResSuccess(c, result)
}
