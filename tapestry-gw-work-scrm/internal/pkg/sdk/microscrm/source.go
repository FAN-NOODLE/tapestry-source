package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gw-work-scrm/internal/app/http/request/source"
)

// SourceApi 员工 SDK
type SourceApi struct {
	BaseSDK
}

// NewScource 返回API实例
func NewScource(host string) *SourceApi {
	urlHost = host
	return &SourceApi{BaseSDK{c: &clients.HTTPClient{}}}
}

// SourceList 素材列表
func (s SourceApi) SourceList(ctx ctx.Context, req *SourceListRequest) (interface{}, error) {
	path := "/micro-scrm/v1/source/" + req.Type + "-list"
	url := urlHost + path
	res, err := s.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &ResponseList{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

// SourceGroupList 素材分组列表
func (s SourceApi) SourceGroupList(ctx ctx.Context, req *SourceGroupListRequest) (interface{}, error) {
	path := "/micro-scrm/v1/source-group/list"
	url := urlHost + path
	res, err := s.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

// SourceMakeMedia 生成临时素材
func (s SourceApi) SourceMakeMedia(ctx ctx.Context, req *SourceMediaRequest) (interface{}, error) {
	path := "/micro-scrm/v1/source/make-media"
	url := urlHost + path
	res, err := s.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

// GetUploadCredential 导购获取上传的凭证
func (s SourceApi) GetUploadCredential(ctx ctx.Context) (interface{}, error) {
	path := "/micro-scrm/v1/source/credential"
	url := urlHost + path
	res, err := s.Post(ctx, url, nil)
	if err != nil {
		return nil, err
	}

	result := &Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

// UploadMedia 导购上传素材
func (s SourceApi) UploadMedia(ctx ctx.Context, req *source.UploadMediaRequest) (interface{}, error) {
	path := "/micro-scrm/v1/source/upload-media"
	url := urlHost + path
	res, err := s.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (s SourceApi) Compress(ctx ctx.Context, req *source.CompressReq) (interface{}, error) {
	path := "/micro-scrm/v1/source/compress"
	url := urlHost + path
	res, err := s.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

// CheckFirst 验证是否用户首次打开素材交互提示
func (s SourceApi) CheckFirst(ctx ctx.Context, req *source.CheckAuth) (interface{}, error) {

	path := "/micro-scrm/v1/source/check-first"
	url := urlHost + path
	res, err := s.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err

}

func (s SourceApi) ReportTimes(ctx ctx.Context, req source.ReportTimesReq) (interface{}, error) {
	path := "/micro-scrm/v1/source/report-times"
	url := urlHost + path
	res, err := s.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (s SourceApi) BrandSearchList(ctx ctx.Context, req source.BrandListReq) (interface{}, error) {
	path := "/micro-scrm/v1/source/brand/search-list"
	url := urlHost + path
	res, err := s.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (s SourceApi) BrandGetCacheInfo(ctx ctx.Context, req source.BrandCacheReq) (interface{}, error) {
	path := "/micro-scrm/v1/source/brand/get-cache-info"
	url := urlHost + path
	res, err := s.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (s SourceApi) CollectEdit(ctx ctx.Context, req source.CollectReq) (interface{}, error) {
	path := "/micro-scrm/v1/source/collect/edit"
	url := urlHost + path
	res, err := s.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (s SourceApi) BrandCollectSearchList(ctx ctx.Context, req source.BrandListReq) (interface{}, error) {
	path := "/micro-scrm/v1/source/brand/collect/search-list"
	url := urlHost + path
	res, err := s.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}
