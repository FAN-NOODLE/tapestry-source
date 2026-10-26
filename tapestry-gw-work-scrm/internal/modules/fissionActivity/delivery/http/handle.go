package fissionActivityhttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-work-scrm/internal/pkg/errs"
	"gw-work-scrm/internal/pkg/sdk/boss"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
	"net/http"
)

var FissionActivitySdk *microScrm.FissionActivityApi

// fissionActivityList 裂变活动列表
func fissionActivityList(c *gin.Context) {
	var req microScrm.FissionActivityListRequest
	jsonError := c.BindJSON(&req)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	req.EnterpriseId = c.Value("enterpriseId").(int64)
	enterpriseHash := c.Value("enterprise_hash").(string)

	if enterpriseHash == "" {
		// 查询企业hash
		enterpriseReq := boss.EnterpriseDetailRequest{
			EnterpriseID: req.EnterpriseId,
		}
		enterpriseInfo, err := boss.EnterpriseDetail(ctx.Wrap(c), &enterpriseReq)
		if err != nil {
			ginplus.ResError(c, errs.TokenEnterpriseIdErr)
			return
		}
		enterpriseHash = enterpriseInfo.Hash
	}
	req.EnterpriseHash = enterpriseHash
	list, err := FissionActivitySdk.FissionActivityList(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
}
