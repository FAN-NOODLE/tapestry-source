package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
)

// SourceApi 素材 SDK
type SourceApi struct {
	BaseSDK
}

func NewSource(host string) *SourceApi {
	urlHost = host
	return &SourceApi{BaseSDK{c: &clients.HTTPClient{}}}
}

// SourceMiniProgramList 小程序素材--列表
func (api *SourceApi) SourceMiniProgramList(c ctx.Context, search *RequestSourceSearch) (int64, *[]SourceMini, error) {
	path := "/micro-scrm/v1/source/mini-program-list"
	url := urlHost + path
	res, err := api.Post(c, url, search)
	if err != nil {
		return 0, nil, err
	}
	result := &ResponseList{
		Data: List{
			List: &[]SourceMini{},
		},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != 200 {
		return 0, nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if err != nil {
		return 0, nil, err
	}

	return result.Data.Total, result.Data.List.(*[]SourceMini), err
}

// SourceMiniProgramList 视频素材--列表
func (api *SourceApi) SourceVideoList(c ctx.Context, search *RequestSourceSearch) (int64, *[]SourceVideo, error) {
	path := "/micro-scrm/v1/source/video-list"
	url := urlHost + path
	res, err := api.Post(c, url, search)
	if err != nil {
		return 0, nil, err
	}
	result := &ResponseList{
		Data: List{
			List: &[]SourceVideo{},
		},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return 0, nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if err != nil {
		return 0, nil, err
	}
	return result.Data.Total, result.Data.List.(*[]SourceVideo), err
}

// SourceMiniProgramDetail 小程序素材--详情
func (api *SourceApi) SourceMiniProgramDetail(c ctx.Context, vo *ID) (*SourceMini, error) {
	path := "/micro-scrm/v1/source/mini-program-detail"
	url := urlHost + path
	res, err := api.Post(c, url, vo)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &SourceMini{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if err != nil {
		return nil, err
	}

	return result.Data.(*SourceMini), err
}

// SourceVideoDetail 视频素材--详情
func (api *SourceApi) SourceVideoDetail(c ctx.Context, vo *ID) (*SourceVideo, error) {
	path := "/micro-scrm/v1/source/video-detail"
	url := urlHost + path
	res, err := api.Post(c, url, vo)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &SourceVideo{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if err != nil {
		return nil, err
	}

	if result.Data == nil {
		return nil, nil
	}
	return result.Data.(*SourceVideo), err
}

// SourceMiniProgramList 小程序素材--新增
func (api *SourceApi) SourceMiniProgramAdd(c ctx.Context, params *RequestSourceMini) (*SourceMini, error) {
	path := "/micro-scrm/v1/source/mini-program-add"
	url := urlHost + path
	res, err := api.Post(c, url, params)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &SourceMini{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if err != nil {
		return nil, err
	}
	return result.Data.(*SourceMini), err
}

// SourceMiniProgramList 小程序素材--更新
func (api *SourceApi) SourceMiniProgramUpdate(c ctx.Context, params *RequestSourceMini) (*ID, error) {

	path := "/micro-scrm/v1/source/mini-program-edit"
	url := urlHost + path
	res, err := api.Post(c, url, params)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &ID{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if err != nil {
		return nil, err
	}
	return result.Data.(*ID), err
}

// SourceMiniProgramList 视频素材--新增
func (api *SourceApi) SourceVideoAdd(c ctx.Context, params *RequestSourceVideo) (*SourceVideo, error) {

	path := "/micro-scrm/v1/source/video-add"
	url := urlHost + path
	res, err := api.Post(c, url, params)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &SourceVideo{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if err != nil {
		return nil, err
	}
	return result.Data.(*SourceVideo), err
}

// SourceMiniProgramList 视频素材--更新
func (api *SourceApi) SourceVideoUpdate(c ctx.Context, params *RequestSourceVideo) (*ID, error) {

	path := "/micro-scrm/v1/source/video-edit"
	url := urlHost + path
	res, err := api.Post(c, url, params)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &ID{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if err != nil {
		return nil, err
	}
	return result.Data.(*ID), err
}

// 图片素材新增
func (api *SourceApi) SourcesImageAdd(c ctx.Context, params *SourcesImageRequest) (*[]SourcesImageRequestAdd, error) {
	path := "/micro-scrm/v1/source/image-add"
	url := urlHost + path
	res, err := api.Post(c, url, params)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &[]SourcesImageRequestAdd{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	return result.Data.(*[]SourcesImageRequestAdd), err
}

// 文件素材新增
func (api *SourceApi) SourcesFileAdd(c ctx.Context, params *SourcesFileRequest) (*SourcesFileRequest, error) {
	path := "/micro-scrm/v1/source/file-add"
	url := urlHost + path
	res, err := api.Post(c, url, params)

	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &SourcesFileRequest{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if err != nil {
		return nil, err
	}
	return result.Data.(*SourcesFileRequest), err
}

// 文件素材列表
func (api *SourceApi) SourceFileList(c ctx.Context, customer *RequestList) (int64, *ResponseFileList, error) {
	path := "/micro-scrm/v1/source/file-list"
	url := urlHost + path
	res, err := api.Post(c, url, customer)
	if err != nil {
		return 0, nil, err
	}
	result := &ResponseList{
		Data: List{
			List: &ResponseFileList{},
		},
	}
	_, err = res.Read(clients.JSON(result))
	if err != nil {
		return 0, nil, err
	}
	if result.Code != StatusOK {
		return 0, nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	return result.Data.Total, result.Data.List.(*ResponseFileList), err
}

// 文件素材列表
func (api *SourceApi) SourceImageList(c ctx.Context, customer *RequestList) (int64, *ResponseImgList, error) {
	path := "/micro-scrm/v1/source/image-list"
	url := urlHost + path
	res, err := api.Post(c, url, customer)
	if err != nil {
		return 0, nil, err
	}
	result := &ResponseList{
		Data: List{
			List: &ResponseImgList{},
		},
	}
	_, err = res.Read(clients.JSON(result))
	if err != nil {
		return 0, nil, err
	}
	if result.Code != StatusOK {
		return 0, nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	return result.Data.Total, result.Data.List.(*ResponseImgList), err
}

// SourceTextSave 保存文本素材
func (api *SourceApi) SourceTextSave(c ctx.Context, params *RequestTextSave) (*ResponseTextDetail, error) {
	path := "/micro-scrm/v1/source/text-save"
	url := urlHost + path
	res, err := api.Post(c, url, params)

	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &ResponseTextDetail{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if err != nil {
		return nil, err
	}
	return result.Data.(*ResponseTextDetail), err
}

// SourceLinkSave 保存链接素材
func (api *SourceApi) SourceLinkSave(c ctx.Context, params *RequestLinkSave) (*ResponseLinkDetail, error) {
	path := "/micro-scrm/v1/source/link-save"
	url := urlHost + path
	res, err := api.Post(c, url, params)

	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &ResponseLinkDetail{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if err != nil {
		return nil, err
	}
	return result.Data.(*ResponseLinkDetail), err
}

// SourceTextList 文本素材列表
func (api *SourceApi) SourceTextList(c ctx.Context, customer *RequestTextList) (int64, *ResponseTextList, error) {
	path := "/micro-scrm/v1/source/text-list"
	url := urlHost + path
	res, err := api.Post(c, url, customer)
	if err != nil {
		return 0, nil, err
	}
	result := &ResponseList{
		Data: List{
			List: &ResponseTextList{},
		},
	}
	_, err = res.Read(clients.JSON(result))
	if err != nil {
		return 0, nil, err
	}
	if result.Code != StatusOK {
		return 0, nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	return result.Data.Total, result.Data.List.(*ResponseTextList), err
}

// SourceLinkList 链接素材列表
func (api *SourceApi) SourceLinkList(c ctx.Context, customer *RequestLinkList) (int64, *ResponseLinkList, error) {
	path := "/micro-scrm/v1/source/link-list"
	url := urlHost + path
	res, err := api.Post(c, url, customer)
	if err != nil {
		return 0, nil, err
	}
	result := &ResponseList{
		Data: List{
			List: &ResponseLinkList{},
		},
	}
	_, err = res.Read(clients.JSON(result))
	if err != nil {
		return 0, nil, err
	}
	if result.Code != StatusOK {
		return 0, nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	return result.Data.Total, result.Data.List.(*ResponseLinkList), err
}

// SourceTextDetail 文本素材详情
func (api *SourceApi) SourceTextDetail(c ctx.Context, vo *RequestSourceDetail) (*ResponseTextDetail, error) {
	path := "/micro-scrm/v1/source/text-detail"
	url := urlHost + path
	res, err := api.Post(c, url, vo)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &ResponseTextDetail{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if err != nil {
		return nil, err
	}

	return result.Data.(*ResponseTextDetail), err
}

// SourceLinkDetail 链接素材详情
func (api *SourceApi) SourceLinkDetail(c ctx.Context, vo *RequestSourceDetail) (*ResponseLinkDetail, error) {
	path := "/micro-scrm/v1/source/link-detail"
	url := urlHost + path
	res, err := api.Post(c, url, vo)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &ResponseLinkDetail{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if err != nil {
		return nil, err
	}

	return result.Data.(*ResponseLinkDetail), err
}

// 删除素材
func (api *SourceApi) SourceDelete(c ctx.Context, source *RequestSourceDelete) (*Rows, error) {
	path := "/micro-scrm/v1/source/delete"
	url := urlHost + path
	res, err := api.Post(c, url, source)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &Rows{},
	}
	_, err = res.Read(clients.JSON(result))
	if err != nil {
		return nil, err
	}
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	return result.Data.(*Rows), err
}

// 素材移动分组
func (api *SourceApi) MoveGroup(c ctx.Context, source *RequestGroupMove) (*Rows, error) {
	path := "/micro-scrm/v1/source/move-group"
	url := urlHost + path
	res, err := api.Post(c, url, source)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &Rows{},
	}
	_, err = res.Read(clients.JSON(result))
	if err != nil {
		return nil, err
	}
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	return result.Data.(*Rows), err
}
