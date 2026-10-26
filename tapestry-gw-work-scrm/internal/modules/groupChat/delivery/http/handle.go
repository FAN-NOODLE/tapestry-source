package http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
	"gw-work-scrm/internal/pkg/validate_params"
)

var groupChatSdk *microScrm.GroupChatApi

// index 社群首页数据分析
func index(c *gin.Context) {
	req := &microScrm.GroupChatStaticsRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	req.MemberID = c.Value("memberId").(int64)
	list, err := groupChatSdk.Index(ctx.Wrap(c), *req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
}

// index 社群首页数据分析
func membersGuide(c *gin.Context) {
	req := &microScrm.GroupChatStaticsRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	req.MemberID = c.Value("memberId").(int64)
	req.UserID = c.Value("userId").(string)
	list, err := groupChatSdk.MembersGuide(ctx.Wrap(c), *req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
}

// members 群成员明细
func members(c *gin.Context) {
	req := &microScrm.GroupChatStaticsRequest{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	req.MemberID = c.Value("memberId").(int64)
	list, err := groupChatSdk.Members(ctx.Wrap(c), *req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
}

// getLastDataUpdateTime 数据更新时间
func getLastDataUpdateTime(c *gin.Context) {
	var req = microScrm.RequestBase{
		EnterpriseID: c.Value("enterpriseId").(int64),
	}
	resp, err := groupChatSdk.GetDataUpdateTime(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, resp)
	return
}

// getLastDataUpdateTime 数据更新时间
func getActivityInfo(c *gin.Context) {
	req := &microScrm.ActivityInfo{}
	ok := validate_params.ValidateParams(c, req)
	if !ok {
		return
	}
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	resp, err := groupChatSdk.ActivityInfo(ctx.Wrap(c), *req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, resp)
	return
}
