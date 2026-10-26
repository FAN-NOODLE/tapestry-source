package source

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/app/http/handle"
	"gw-scrm/internal/app/http/request/source"
	"gw-scrm/internal/pkg/common"
	errs "gw-scrm/internal/pkg/errs"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type Handle struct {
	handle.BaseHandle
	sourceSdk microscrm.SourceApi
}

func NewHandle() Handle {
	return Handle{
		sourceSdk: microscrm.NewSourceApi(),
	}
}

// TextList 文本素材列表
func (s Handle) TextList(c *gin.Context) {
	var request source.RequestTextList
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	request.AdminID = adminId

	res, err := s.sourceSdk.TextList(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

// ImageList 图片素材列表
func (s Handle) ImageList(c *gin.Context) {
	var request source.RequestList
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	request.AdminID = adminId

	res, err := s.sourceSdk.ImageList(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

// LinkList 链接素材列表
func (s Handle) LinkList(c *gin.Context) {
	var request source.RequestLinkList
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	request.AdminID = adminId

	res, err := s.sourceSdk.LinkList(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

// MiniProgramList 小程序素材列表
func (s Handle) MiniProgramList(c *gin.Context) {
	var request source.RequestSourceSearch
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	request.AdminID = adminId

	res, err := s.sourceSdk.MiniProgramList(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

// VideoList 视频素材列表
func (s Handle) VideoList(c *gin.Context) {
	var request source.RequestSourceSearch
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	request.AdminID = adminId

	res, err := s.sourceSdk.VideoList(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

// FileList 文本素材列表
func (s Handle) FileList(c *gin.Context) {
	var request source.RequestList
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	request.AdminID = adminId

	res, err := s.sourceSdk.FileList(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}

	s.ResSuccess(c, res)
	return
}

// SourceDetail 各类素材详情
func (s Handle) SourceDetail(c *gin.Context) {
	var request source.RequestDetail
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID

	var (
		res interface{}
		err error
	)
	cc := ctx.Wrap(c)

	switch request.Type {

	// 类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件)
	case 1:
		res, err = s.sourceSdk.TextDetail(cc, request)
	case 2:
		res, err = s.sourceSdk.ImageDetail(cc, request)
	case 3:
		res, err = s.sourceSdk.LinkDetail(cc, request)
	case 4:
		res, err = s.sourceSdk.MiniProgramDetail(cc, request)
	case 5:
		res, err = s.sourceSdk.VideoDetail(cc, request)
	case 6:
		res, err = s.sourceSdk.FileDetail(cc, request)
	case 7:
		res, err = s.sourceSdk.FriendDetail(cc, request)

	default:
		s.ResError(c, errs.ParamsErr)
		return
	}
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

// SourceUpdate 各类素材编辑保存
func (s Handle) SourceUpdate(c *gin.Context) {
	var request source.RequestUpdate
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}
	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	request.AdminId = adminId

	var (
		res interface{}
		err error
	)
	cc := ctx.Wrap(c)

	switch request.SourceType {

	// 类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件 7 朋友圈)
	case 1:
		res, err = s.sourceSdk.TextUpdate(cc, request)
	case 2:
		res, err = s.sourceSdk.ImageUpdate(cc, request)
	case 3:
		res, err = s.sourceSdk.LinkUpdate(cc, request)
	case 4:
		res, err = s.sourceSdk.MiniProgramUpdate(cc, request)
	case 5:
		res, err = s.sourceSdk.VideoUpdate(cc, request)
	case 6:
		res, err = s.sourceSdk.FileUpdate(cc, request)
	case 7:
		res, err = s.sourceSdk.FriendUpdate(cc, request)
	default:
		s.ResError(c, errs.ParamsErr)
		return
	}

	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

func (s Handle) RemoveCategory(c *gin.Context) {
	var request source.RequestRemoveCategory
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID

	res, err := s.sourceSdk.RemoveCategory(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}

	s.ResSuccess(c, res)
	return
}

func (s Handle) UpdateStatus(c *gin.Context) {

	var request source.RequestUpdateStatus
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	res, err := s.sourceSdk.UpdateStatus(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

// ConsoleList 后台管理查询各素材列表
func (s Handle) ConsoleList(c *gin.Context) {
	var request source.RequestConsoleList

	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	request.AdminID = adminId

	res, err := s.sourceSdk.ConsoleList(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

func (s Handle) GetMediaInfo(c *gin.Context) {
	var request source.RequestGetMediaInfo
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	res, err := s.sourceSdk.GetMediaInfo(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

func (s Handle) CompressFile(c *gin.Context) {
	var request struct {
		FilePath string `json:"file_path" binding:"required" required:"图片地址必填"`
		FileType string `json:"file_type" binding:"required" required:"资源类型必填"`
	}

	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	res, err := s.sourceSdk.CompressFile(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

func (s Handle) SetTop(c *gin.Context) {
	var request source.SourceSetTop
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	res, err := s.sourceSdk.SetTop(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
}
