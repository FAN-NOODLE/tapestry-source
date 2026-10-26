package drainagehttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/drainage"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/sdk/work"
)

var contactService *drainage.ContactService
var enterpriseService *enterprise.Service

func contactCreate(c *gin.Context) {
	var request drainage.ContactCreateRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	corpReq := enterprise.Corp{
		EnterpriseID: request.EnterpriseId,
	}
	enterpriseData, _ := enterpriseService.Detail(ctx.Wrap(c), corpReq)
	request.Domain = enterpriseData.Domain
	ID, queryError := contactService.ContactCreate(ctx.Wrap(c), request)
	if queryError != nil {
		ginplus.ResError(c, queryError)
		return
	}
	ginplus.ResSuccess(c, ID)
	return
}

func contactUpdate(c *gin.Context) {
	var request drainage.ContactCreateRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	contactDetailsRequest := drainage.ContactDetailsRequest{
		ID:           request.ID,
		EnterpriseId: request.EnterpriseId,
	}
	detailsData, _ := contactService.GetContact(ctx.Wrap(c), contactDetailsRequest)
	corpReq := enterprise.Corp{
		EnterpriseID: request.EnterpriseId,
	}
	enterpriseData, _ := enterpriseService.Detail(ctx.Wrap(c), corpReq)
	workImpl := new(work.QyWorkImpl).SetDomain(enterpriseData.Domain)
	res, updateWxExterErr := workImpl.UpdateWxEexternalContact(ctx.New(), work.ContactCreateRequest{
		ID:           request.ID,
		CodeCame:     request.CodeCame,
		ChannelId:    request.ChannelId,
		AutoFriend:   request.AutoFriend,
		MemberIds:    request.MemberIds,
		ThirdIds:     request.ThirdIds,
		EnterpriseId: request.EnterpriseId,
		Type:         work.DrainageContactsTypeMultiple,
		UserIds:      request.UserIds,
		State:        request.State,
	}, detailsData.ConfigId)

	if res == nil {
		e, ok := updateWxExterErr.(*errorx.Error)
		// 判断是不是企业微信错误码
		if ok && e.Code == 41054 {
			ginplus.ResError(c, errorx.NewErrorMsg(e.Code, "人员未激活"))
			return
		}
		if ok && e.Code == 41041 {
			ginplus.ResError(c, errorx.NewErrorMsg(e.Code, "总的人员数量不能超过100"))
			return
		}
		ginplus.ResError(c, updateWxExterErr)
		return
	}

	// 查询渠道是否存在
	channelRequest := drainage.ContactChannels{
		ID:           request.ChannelId,
		EnterpriseID: request.EnterpriseId,
	}
	detail, _ := drainageService.ChannelDetail(ctx.Wrap(c), &channelRequest)
	if detail == nil {
		ginplus.ResError(c, errs.ChannelNotExists)
		return
	}

	response, queryError := contactService.ContactUpdate(ctx.Wrap(c), request)
	if queryError != nil {
		ginplus.ResError(c, queryError)
		return
	}
	ginplus.ResSuccess(c, response)
}
func contactDetails(c *gin.Context) {
	var request drainage.ContactDetailsRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	response, queryError := contactService.ContactDetails(ctx.Wrap(c), request)
	if queryError != nil {
		ginplus.ResError(c, queryError)
		return
	}
	ginplus.ResSuccess(c, response)
}
func contactDelete(c *gin.Context) {
	var request drainage.ContactDeleteRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	contactDetailsRequest := drainage.ContactDetailsRequest{
		ID:           request.ID[0],
		EnterpriseId: request.EnterpriseId,
	}
	detailsData, _ := contactService.GetContact(ctx.Wrap(c), contactDetailsRequest)
	corpReq := enterprise.Corp{
		EnterpriseID: request.EnterpriseId,
	}
	enterpriseData, _ := enterpriseService.Detail(ctx.Wrap(c), corpReq)
	workImpl := new(work.QyWorkImpl).SetDomain(enterpriseData.Domain)
	res, deleteWxExterErr := workImpl.DeleteWxEexternalContact(ctx.New(), detailsData.ConfigId)
	if res.Errcode != 0 {
		ginplus.ResError(c, deleteWxExterErr)
		return
	}
	queryError := contactService.ContactDelete(ctx.Wrap(c), request)
	if queryError != nil {
		ginplus.ResError(c, queryError)
		return
	}
	ginplus.ResSuccess(c, []drainage.ContactDeleteRequest{})
}
func contactList(c *gin.Context) {
	var request drainage.ContactsListResponse
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	total, resp := contactService.List(ctx.Wrap(c), request)
	res := drainage.CommonList{}
	if resp == nil {
		res.List = []drainage.Contacts{}
	} else {
		res.List = *resp
	}
	res.Total = total
	res.Page = request.Page
	res.Size = request.Size
	ginplus.ResSuccess(c, res)
}

func contactQrCode(c *gin.Context) {
	var request drainage.QrCodeRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	file, queryError := contactService.DownLoadQrCode(ctx.Wrap(c), request)
	if queryError != nil {
		ginplus.ResError(c, queryError)
		return
	}
	ginplus.ResSuccess(c, file)
}
