package http

import (
	"errors"

	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/common"

	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/activityCampaign"
)

// signInDetail 签到数据列表
func signInList(c *gin.Context) {
	var request activityCampaign.SignReceiveInfo
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	info, count, err := activityCampaign.GetSignInInfo(enterpriseID.(int64), request.ID, request.StartTime,
		request.EndTime, request.Page.PageSize, request.Page.PageNum, request.Type)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var result micro_scrm_enterprise.NewResultList
	result.Count = count
	result.SignInList = info
	ginplus.ResSuccess(c, result)
}

func signInDetail(c *gin.Context) {
	var request activityCampaign.SignDetailInfo
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	detail, err := activityCampaign.GetSignInInfoDetail(enterpriseID.(int64), request.ID, request.StartTime,
		request.EndTime, request.UnionID, request.Page.PageSize, request.Page.PageNum)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, detail)
}

func activitiesH5Info(c *gin.Context) {
	params := activityCampaign.ActivitiesRequest{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if params.ActivityId == 0 {
		ginplus.ResError(c, errors.New("活动编号错误"))
		return
	}
	result := activityCampaign.ActivitiesH5Info(params)
	ginplus.ResSuccess(c, result)
}

// H5用户签到列表(小程序||H5)
func usersignInH5List(c *gin.Context) {
	var params activityCampaign.UsersignInRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if params.UnionId == "" {
		ginplus.ResError(c, errors.New("用户的union_id不能为空"))
		return
	}
	result := activityCampaign.H5UsersignInList(params)
	ginplus.ResSuccess(c, result)
}

func usersignInH5(c *gin.Context) {
	var params activityCampaign.UsersignInRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if params.UnionId == "" {
		ginplus.ResError(c, errors.New("用户的union_id不能为空"))
		return
	}
	if err := activityCampaign.UsersignInH5(ctx.Wrap(c), params); err == nil {
		ginplus.ResSuccess(c, err)
	} else {
		ginplus.ResError(c, err)
	}
}
func usersignInInfoH5(c *gin.Context) {
	var params activityCampaign.UsersignInRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if params.UnionId == "" {
		ginplus.ResError(c, errors.New("用户的union_id不能为空"))
		return
	}
	if params.SignInDate == "" {
		ginplus.ResError(c, errors.New("用户签到日期不能为空"))
		return
	}
	result := activityCampaign.UsersignInInfoH5(params)
	ginplus.ResSuccess(c, result)

}
