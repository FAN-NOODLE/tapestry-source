package sourcehttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-work-scrm/internal/modules/source"
	"gw-work-scrm/internal/pkg/errs"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
	"gw-work-scrm/internal/pkg/validate_params"
)

var sourceSdk *microScrm.SourceApi

// sourceList 文本素材
func sourceList(c *gin.Context) {
	req := &microScrm.SourceListRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	isOk := false
	for _, v := range source.TypeSourceArr {
		if v == req.Type {
			isOk = true
		}
	}
	if !isOk {
		ginplus.ResError(c, errs.SourceTypeErr)
		return
	}
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	if req.Page == 0 {
		req.Page = 1
	}

	if req.Size == 0 {
		req.Size = 10
	}
	memberInfo := c.Value("memberInfo").(*microScrm.Member)
	req.MemberUserId = memberInfo.Userid
	req.MemberID = memberInfo.ID
	list, err := sourceSdk.SourceList(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
}

// sourceGroupList 素材分组列表
func sourceGroupList(c *gin.Context) {
	req := &microScrm.SourceGroupListRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	list, err := sourceSdk.SourceGroupList(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
}

// makeMedia 素材分组列表
func makeMedia(c *gin.Context) {
	req := &microScrm.SourceMediaRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	list, err := sourceSdk.SourceMakeMedia(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
}
