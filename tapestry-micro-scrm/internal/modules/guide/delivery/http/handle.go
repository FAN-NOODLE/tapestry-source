package guideHttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/guide"
	"micro-scrm/internal/modules/members"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/errs"
)

var guideService *guide.Service
var membersService *members.Service

// SyncGuide 同步导购绑定关系相关接口
func SyncGuide(c *gin.Context) {
	method := c.GetHeader(common.DDCMethodHeader)
	var err error
	switch method {
	case common.DdcEventEnum.UserBindingGuideAdd.String(),
		common.DdcEventEnum.UserBindingGuideUpdate.String(),
		common.DdcEventEnum.UserBindingGuideDelete.String():
		var request guide.SyncUserBingGuideRequest
		err = c.BindJSON(&request)
		if err != nil {
			ginplus.ResError(c, errs.ParamsErr)
			return
		}
		_, err = guideService.SyncUserBingGuide(ctx.Wrap(c), &request.DataUserBingGuide.UserBingGuideList, common.DdcEvent(request.Method))

	default:
		ginplus.ResError(c, errs.MethodHeaderErr)
		return
	}

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}
