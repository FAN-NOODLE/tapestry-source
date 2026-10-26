package open

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/middleware"
	openReq "micro-scrm/internal/app/http/request/open"
	customerSvc "micro-scrm/internal/app/service/customer"
	openSvc "micro-scrm/internal/app/service/open"
)

type WecomCommon struct {
	handle.BaseHandle
	wecomCommonSvc *openSvc.WecomCommonService
}

func NewWecomCommon() *WecomCommon {
	return &WecomCommon{
		wecomCommonSvc: openSvc.NewWecomCommonService(),
	}
}

// ConvertOpenUserIdToUserId 将企微代开发应用的open_userid通过自建应用转成userid
func (h *WecomCommon) ConvertOpenUserIdToUserId(c *gin.Context) {
	var request openReq.ConvertOpenUserIdRequest
	request.EnterpriseID = c.GetInt64(middleware.EnterpriseId)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if request.SourceAgentId == 0 {
		h.ResError(c, errorx.New("source_agentid is empty"))
		return
	}
	if len(request.OpenUseridList) == 0 {
		h.ResError(c, errorx.New("open_userid_list is empty"))
		return
	}
	userId, err := h.wecomCommonSvc.ConvertOpenUserIdToUserId(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, userId)
	return
}

// GetInfoBySourceExternal 通过企微代开发应用的open_external_userid获取外部联系人信息
func (h *WecomCommon) GetInfoBySourceExternal(c *gin.Context) {
	var request openReq.ConvertOpenExternalUserIdRequest
	request.EnterpriseID = c.GetInt64(middleware.EnterpriseId)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if request.SourceAgentId == 0 {
		h.ResError(c, errorx.New("source_agentid is empty"))
		return
	}
	if request.ExternalUserid == "" {
		h.ResError(c, errorx.New("external_userid is empty"))
		return
	}
	res, err := h.wecomCommonSvc.ConvertOpenExternalUseridToExternalUserid(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	detail, err := customerSvc.NewCustomer().CustomerDetail(h.WarpContext(c), request.EnterpriseID, res.ExternalUserid)
	if err != nil {
		return
	}

	h.ResSuccess(c, map[string]interface{}{
		"external_userid": detail.Openid,
		"unionid":         detail.Unionid,
		"nick":            detail.Nick,
		"avatar":          detail.Avatar,
	})
	return
}

// ConvertOpenExternalUseridToExternalUserid 将企微代开发应用的open_external_userid通过自建应用转成external_userid
func (h *WecomCommon) ConvertOpenExternalUseridToExternalUserid(c *gin.Context) {
	var request openReq.ConvertOpenExternalUserIdRequest
	request.EnterpriseID = c.GetInt64(middleware.EnterpriseId)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if request.SourceAgentId == 0 {
		h.ResError(c, errorx.New("source_agentid is empty"))
		return
	}
	if request.ExternalUserid == "" {
		h.ResError(c, errorx.New("external_userid is empty"))
		return
	}
	externalUserId, err := h.wecomCommonSvc.ConvertOpenExternalUseridToExternalUserid(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, externalUserId)
	return
}
