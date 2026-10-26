package open

import (
	"errors"

	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/middleware"
	openReq "micro-scrm/internal/app/http/request/open"
	membersSvc "micro-scrm/internal/app/service/member"
	openSvc "micro-scrm/internal/app/service/open"
)

type Source struct {
	handle.BaseHandle
	sourceSvc *openSvc.SourceService
}

func NewSource() Source {
	return Source{
		sourceSvc: openSvc.NewSourceService(),
	}
}

// LinkList 获取对外链接列表
func (h Source) LinkList(c *gin.Context) {
	var request openReq.SearchListRequest
	request.EnterpriseID = c.GetInt64(middleware.EnterpriseId)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if request.UserId != "" {
		var err error
		request.ThirdIds, err = membersSvc.NewMembersService().GetThirdIdsInfoByUserId(h.WarpContext(c), request.EnterpriseID, request.UserId)
		if err != nil {
			return
		}
	}
	if request.Page == 0 {
		request.Page = 1
	}
	if request.Size == 0 {
		request.Size = 10
	}
	// 获取链接列表
	lists, total, err := h.sourceSvc.LinkList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResList(c, total, request.Page, request.Size, lists)
	return
}

// SourceTree 获取素材分组树
func (h Source) SourceTree(c *gin.Context) {
	var request openReq.SourceTreeRequest
	request.EnterpriseID = c.GetInt64(middleware.EnterpriseId)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	tree, err := h.sourceSvc.SourceTree(request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, tree)
}

// SourceList   获取素材列表
func (h Source) SourceList(c *gin.Context) {
	var request openReq.SourceOpenReq
	request.EnterpriseID = c.GetInt64(middleware.EnterpriseId)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if request.Page == 0 {
		request.Page = 1
	}
	if request.Size == 0 {
		request.Size = 10
	}

	textList, total, err := h.sourceSvc.SourceList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResList(c, total, request.Page, request.Size, textList)
	return
}

// fillThirdIds 根据 UserId 填充 ThirdIds（权限过滤）
func fillThirdIds(c app.Ctx, request *openReq.SearchListRequest) error {
	if request.UserId != "" && len(request.ThirdIds) == 0 {
		var err error
		request.ThirdIds, err = membersSvc.NewMembersService().GetThirdIdsInfoByUserId(c, request.EnterpriseID, request.UserId)
		return err
	}
	return nil
}

// TextList 获取文本素材列表
func (h Source) TextList(c *gin.Context) {
	var request openReq.SearchListRequest
	request.EnterpriseID = c.GetInt64(middleware.EnterpriseId)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if err := fillThirdIds(h.WarpContext(c), &request); err != nil {
		h.ResError(c, err)
		return
	}
	if request.Page == 0 {
		request.Page = 1
	}
	if request.Size == 0 {
		request.Size = 10
	}

	list, total, err := h.sourceSvc.TextList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResList(c, total, request.Page, request.Size, list)
}

// ImageList 获取图片素材列表
func (h Source) ImageList(c *gin.Context) {
	var request openReq.SearchListRequest
	request.EnterpriseID = c.GetInt64(middleware.EnterpriseId)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if err := fillThirdIds(h.WarpContext(c), &request); err != nil {
		h.ResError(c, err)
		return
	}
	if request.Page == 0 {
		request.Page = 1
	}
	if request.Size == 0 {
		request.Size = 10
	}

	list, total, err := h.sourceSvc.ImageList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResList(c, total, request.Page, request.Size, list)
}

// VideoList 获取视频素材列表
func (h Source) VideoList(c *gin.Context) {
	var request openReq.SearchListRequest
	request.EnterpriseID = c.GetInt64(middleware.EnterpriseId)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if err := fillThirdIds(h.WarpContext(c), &request); err != nil {
		h.ResError(c, err)
		return
	}
	if request.Page == 0 {
		request.Page = 1
	}
	if request.Size == 0 {
		request.Size = 10
	}

	list, total, err := h.sourceSvc.VideoList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResList(c, total, request.Page, request.Size, list)
}

// FriendList 获取朋友圈素材列表
func (h Source) FriendList(c *gin.Context) {
	var request openReq.SearchListRequest
	request.EnterpriseID = c.GetInt64(middleware.EnterpriseId)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if err := fillThirdIds(h.WarpContext(c), &request); err != nil {
		h.ResError(c, err)
		return
	}
	if request.Page == 0 {
		request.Page = 1
	}
	if request.Size == 0 {
		request.Size = 10
	}

	list, total, err := h.sourceSvc.FriendList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResList(c, total, request.Page, request.Size, list)
}

// SearchAll 合并检索：通过商品条件一次性检索所有素材
func (h Source) SearchAll(c *gin.Context) {
	var request openReq.SearchAllRequest
	request.EnterpriseID = c.GetInt64(middleware.EnterpriseId)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	// 校验：至少要有一个商品检索条件
	if request.ProductSku == "" && request.ProductName == "" {
		h.ResError(c, errors.New("至少需要提供一个商品检索条件(product_sku/product_name)"))
		return
	}

	if request.UserId != "" {
		var err error
		request.ThirdIds, err = membersSvc.NewMembersService().GetThirdIdsInfoByUserId(h.WarpContext(c), request.EnterpriseID, request.UserId)
		if err != nil {
			h.ResError(c, err)
			return
		}
	}

	result, err := h.sourceSvc.SearchAll(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)
}
