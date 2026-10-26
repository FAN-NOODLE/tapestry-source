package remindconfigthttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/members"
	"micro-scrm/internal/modules/remindconfig"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/errs"
	response "micro-scrm/internal/pkg/resoponse"
	"micro-scrm/internal/pkg/uid"
	"net/http"
)

var service *remindconfig.Service
var memberService *members.Service

// configInfo 获取提醒设置
func configInfo(c *gin.Context) {
	request := &remindconfig.ConfigInfoRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	configInfo, err := service.ConfigInfo(ctx.Wrap(c), request, memberService)
	if errorIgnore.Ignore(err) != nil {
		ginplus.ResError(c, errs.ConfigInfoErr)
		return
	}
	if configInfo == nil {
		configInfo = &remindconfig.ConfigInfoResponse{}
	}
	response.Success(c, configInfo)
	return
}

// configSave 保存提醒设置
func configSave(c *gin.Context) {
	request := &remindconfig.ConfigSaveRequest{}
	if err := c.ShouldBind(request); err != nil {
		ginplus.ResError(c, err)
		return
	}

	if request.ID == 0 {
		cc := ctx.Wrap(c)
		request.ID = uid.ID(cc)
		if request.NoticeType == remindconfig.OvertimeNotice {
			if request.ReceiveStart == "" || request.ReceiveEnd == "" {
				ginplus.ResError(c, errs.LostTimeErr)
				return
			}
			if request.OverTime == 0 || request.OverTime < remindconfig.MinOvertime || request.OverTime > remindconfig.MaxOvertime {
				ginplus.ResError(c, errs.OvertimeErr)
				return
			}
		}
		id, errAdd := service.ConfigAdd(cc, request)
		if errAdd != nil {
			ginplus.ResError(c, errAdd)
			return
		}
		responseData := customer.AddMessageResponse{ID: id}
		ginplus.ResSuccess(c, responseData)
		return
	} else {
		_, UpErr := service.ConfigSave(ctx.Wrap(c), request)
		if UpErr != nil {
			ginplus.ResError(c, UpErr)
			return
		}
		ginplus.ResSuccess(c, []struct{}{})
		return
	}
}
