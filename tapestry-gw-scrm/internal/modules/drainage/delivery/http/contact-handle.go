package drainage_http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/drainage"
	"gw-scrm/internal/pkg/common"
	scrm_errs "gw-scrm/internal/pkg/errs"
	"gw-scrm/internal/pkg/sdk/microscrm"
	"net"
	"net/url"
)

func contactList(c *gin.Context) {
	var request drainage.ContactsListResponse
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	request.EnterpriseId = tokenInfo.EnterpriseId
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	response, err := microscrm.ContactList(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}
func contactCreate(c *gin.Context) {
	var request drainage.ContactCreateRequest
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	request.EnterpriseId = tokenInfo.EnterpriseId
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	request.Type = 1
	response, err := microscrm.ContactCreate(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}
func contactUpdate(c *gin.Context) {
	var request drainage.ContactCreateRequest
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	request.EnterpriseId = tokenInfo.EnterpriseId
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	response, err := microscrm.ContactUpdate(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

func contactDelete(c *gin.Context) {
	var request drainage.ContactDeleteRequest
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	request.EnterpriseId = tokenInfo.EnterpriseId
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	response, err := microscrm.ContactDelete(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

func contactDetails(c *gin.Context) {
	var request drainage.ContactDetailsRequest
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	request.EnterpriseId = tokenInfo.EnterpriseId
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	response, err := microscrm.ContactDetails(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

func downloadQrCode(c *gin.Context) {
	var request drainage.QrCodeRequest
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	request.EnterpriseId = tokenInfo.EnterpriseId
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}

	// 检测图片的链接是否为ip
	urlObj, urlErr := url.Parse(request.ImageUrl)
	if urlErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	if ip := net.ParseIP(urlObj.Host); ip != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}

	response, err := microscrm.DownloadQrCode(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}
