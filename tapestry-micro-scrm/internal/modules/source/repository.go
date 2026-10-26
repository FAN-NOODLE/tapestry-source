// Package callback 对接Open 平台
package source

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

//go:generate mockgen -source=${GOFILE} -destination=mocks/mock_${GOFILE}

// Repository 存储定义
type Repository interface {
	repos.IRepository
	// TextSave 文本素材保存
	TextSave(ctx ctx.Context, request RequestTextSave) (int64, error)
	// LinkSave 链接素材保存
	LinkSave(ctx ctx.Context, request RequestLinkSave) (int64, error)
	// TextDetail 文本素材详情
	TextDetail(ctx ctx.Context, request *RequestSourceDetail) (*ResponseTextDetail, error)
	// LinkDetail 链接素材详情
	LinkDetail(ctx ctx.Context, request *RequestSourceDetail) (*ResponseLinkDetail, error)

	AddSourcesImage(ctx ctx.Context, data *SourcesImageRequest) (int64, error)
	AddSourcesFile(ctx ctx.Context, data *SourcesFileRequest) (int64, error)
	//图片素材列表
	SourcesImageList(ctx ctx.Context, data *RequestList) (int64, *[]ResponseImgList, error)
	//文件素材列表
	SourcesFileList(ctx ctx.Context, data *RequestList) (int64, *[]ResponseFileList, error)
	// SourceMiniProgramAdd 小程序素材--新增
	SourceMiniProgramAdd(ctx ctx.Context, data *RequestSourceMini) (int64, error)
	// SourceMiniProgramUpdate 小程序素材--更新
	SourceMiniProgramUpdate(ctx ctx.Context, data *RequestSourceMini) (int64, error)
	// SourceVideoAdd 视频素材--新增
	SourceVideoAdd(ctx ctx.Context, data *RequestSourceVideo) (int64, error)
	// SourceVideoUpdate 视频素材--更新
	SourceVideoUpdate(ctx ctx.Context, data *RequestSourceVideo) (int64, error)
	SourceVideoList(cc ctx.Context, search *RequestSourceSearch) (int64, *[]SourceVideo, error)
	SourceMiniProgramList(cc ctx.Context, search *RequestSourceSearch) (int64, *[]SourceMini, error)
	// 删除素材
	SourceDelete(c ctx.Context, data *RequestSourceDelete) (int64, error)
	// 移动分组
	MoveGroup(c ctx.Context, data *RequestGroupMove) (int64, error)
	// 更新分组
	UpdateGroup(c ctx.Context, data *RequestGroupMove) (int64, error)
	// 更新素材
	UpdateSourceMedia(c ctx.Context, req UpdateSourceMediaRequest) error
	InitSourceMedia(c ctx.Context, enterpriseID int64, id []int64, sourceType string) error
	SourceMediaDetail(ctx ctx.Context, enterpriseID int64, sourceID int64, sourceType string) (*MediaItemDetail, error)
}
