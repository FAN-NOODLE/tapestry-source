package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/plan_task"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/common"
	response "micro-scrm/internal/pkg/resoponse"
)

type ActivityHandler struct {
	Name string
}

func NewActivityHandler() *ActivityHandler {
	return &ActivityHandler{
		Name: "报名活动",
	}
}

func (h ActivityHandler) VerifyUser(c *gin.Context) {
	var param plan_task.UserInfoReq
	err := c.ShouldBindJSON(&param)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	UpdateErr := plan_task.VerifyUser(enterpriseID.(int64), param.ActivityUserId)
	if UpdateErr != nil {
		ginplus.ResError(c, UpdateErr)
		return
	}
	ginplus.ResSuccess(c, true)
}

func (h ActivityHandler) GetUserInfo(c *gin.Context) {
	var param plan_task.UserInfoReq
	err := c.ShouldBindJSON(&param)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	data, getInfoErr := plan_task.GetUserInfo(enterpriseID.(int64), param.ActivityUserId)
	if getInfoErr != nil {
		ginplus.ResError(c, getInfoErr)
		return
	}
	ginplus.ResSuccess(c, data)
}

func (h ActivityHandler) CheckActivityStatus(c *gin.Context) {
	var param plan_task.CheckStatusReq
	err := c.ShouldBindJSON(&param)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	data, err := plan_task.CheckActivityStatus(param)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if data.ID != 0 {
		decodePhone, decodePhoneErr := common.AesDecrypt(data.ReceiverPhone)
		if decodePhoneErr != nil {
			ginplus.ResError(c, decodePhoneErr)
			return
		}
		data.ReceiverPhone = decodePhone

		decodeName, decodeNameErr := common.AesDecrypt(data.ReceiverName)
		if decodeNameErr != nil {
			ginplus.ResError(c, decodeNameErr)
			return
		}
		data.ReceiverName = decodeName

		ginplus.ResSuccess(c, plan_task.InsertUserResp{
			IsJoin: true,
			Data:   data,
		})
	} else {
		ginplus.ResSuccess(c, plan_task.InsertUserResp{
			IsJoin: false,
			Data:   micro_scrm_enterprise.OperationActivityUserModel{},
		})
	}
}

func (h ActivityHandler) GetActivityId(c *gin.Context) {
	var param plan_task.GetActivityIdReq
	err := c.ShouldBindJSON(&param)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	data, err := plan_task.GetActivityId(param.EnterpriseID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, data)
}

func (h ActivityHandler) InsertUserData(c *gin.Context) {
	var param plan_task.InsertUserReq
	err := c.ShouldBindJSON(&param)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	insertErr := plan_task.InsertUserData(ctx.Wrap(c), param)
	if insertErr != nil {
		ginplus.ResError(c, insertErr)
		return
	}
	ginplus.ResSuccess(c, plan_task.InsertUserResp{
		IsJoin: true,
	})
}

func (h ActivityHandler) ExportUserData(c *gin.Context) {
	var param plan_task.ExportUsreReq
	err := c.ShouldBindJSON(&param)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	data, ListErr := plan_task.ExportUserList(enterpriseID.(int64), param)
	if ListErr != nil {
		ginplus.ResError(c, ListErr)
		return
	}
	ginplus.ResSuccess(c, data)
}

func (h ActivityHandler) GetUserListByTpyeId(c *gin.Context) {
	var param plan_task.UserListReq
	err := c.ShouldBindJSON(&param)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	total, data, err := plan_task.GetUserListByTypeId(enterpriseID.(int64), param)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	response.SuccessList(c, total, param.Page, param.Size, data)
	return
}

func (h ActivityHandler) GetTypeList(c *gin.Context) {
	var param plan_task.BaseReq
	err := c.ShouldBindJSON(&param)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	total, data, err := plan_task.GetTypeList(enterpriseID.(int64), param)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	response.SuccessList(c, total, param.Page, param.Size, data)
	return
}
