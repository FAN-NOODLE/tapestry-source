package taghttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
	"gw-work-scrm/internal/pkg/validate_params"
)

var tagSdk *microScrm.TagApi

// tagList 添加的员工
func tagList(c *gin.Context) {
	req := microScrm.TagRequestList{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	list, err := tagSdk.TagList(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
}

// customerTagList 添加的员工
func customerTagList(c *gin.Context) {
	req := &microScrm.TagRequestList{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	list, err := tagSdk.CustomerManualList(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
}
