package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
)

// SourceGroupApi 素材分组 SDK
type SourceGroupApi struct {
	BaseSDK
}

func NewSourceGroup(host string) *SourceGroupApi {
	urlHost = host
	return &SourceGroupApi{BaseSDK{c: &clients.HTTPClient{}}}
}

// List 分组--列表
func (api *SourceGroupApi) List(c ctx.Context, params *SourceGroup) (*[]SourceGroup, error) {
	path := "/micro-scrm/v1/source-group/list"
	url := urlHost + path
	res, err := api.Post(c, url, params)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &[]SourceGroup{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if err != nil || result.Data == nil {
		return nil, err
	}
	return result.Data.(*[]SourceGroup), err
}

// Add 分组--新增
func (api *SourceGroupApi) Add(c ctx.Context, params *SourceGroup) (*ID, error) {
	path := "/micro-scrm/v1/source-group/add"
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
	if err != nil || result.Data == nil {
		return nil, err
	}
	return result.Data.(*ID), err
}

// Edit 分组--编辑
func (api *SourceGroupApi) Edit(c ctx.Context, params *SourceGroup) (*ID, error) {
	path := "/micro-scrm/v1/source-group/edit"
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
	if err != nil || result.Data == nil {
		return nil, err
	}
	return result.Data.(*ID), err
}

// BatchDel 删除素材分组
func (api SourceGroupApi) BatchDel(c ctx.Context, ids []int64) (*Rows, error) {
	params := &struct {
		IDs []int64 `json:"ids"  validate:"required" required:"素材分组IDs不能为空"`
	}{
		IDs: ids,
	}
	path := "/micro-scrm/v1/source-group/delete"
	url := urlHost + path
	res, err := api.Post(c, url, params)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &Rows{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if err != nil || result.Data == nil {
		return nil, err
	}
	return result.Data.(*Rows), err
}
