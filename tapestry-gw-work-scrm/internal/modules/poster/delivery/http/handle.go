package http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-work-scrm/internal/pkg/errs"
	"gw-work-scrm/internal/pkg/sdk/boss"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
)

// posterSdk 请求API
var posterSdk *microScrm.PosterApi

// 下载海报
func poster(c *gin.Context) {
	enterpriseId, ok := c.Value("enterpriseId").(int64)
	enterpriseHash := c.Value("enterprise_hash").(string)
	memberInfo := c.Value("memberInfo").(*microScrm.Member)
	if !ok {
		ginplus.ResError(c, errs.TokenEnterpriseIdErr)
		return
	}
	if enterpriseHash == "" {
		// 查询企业hash
		enterpriseReq := boss.EnterpriseDetailRequest{
			EnterpriseID: enterpriseId,
		}
		enterpriseInfo, err := boss.EnterpriseDetail(ctx.Wrap(c), &enterpriseReq)
		if err != nil {
			ginplus.ResError(c, errs.TokenEnterpriseIdErr)
			return
		}
		enterpriseHash = enterpriseInfo.Hash
	}

	request := new(microScrm.PosterRequest)
	if err := c.BindJSON(request); err != nil {
		ginplus.ResError(c, err)
		return
	}

	request.EnterpriseID = enterpriseId
	request.EnterpriseHash = enterpriseHash
	request.MemUserID = memberInfo.Userid
	request.MemberName = memberInfo.Name
	request.MemberAvatar = memberInfo.Avatar
	result, err := posterSdk.Poster(ctx.Wrap(c), *request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, result)
}
