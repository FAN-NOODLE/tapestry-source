package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
)

// MiniProgramApi 小程序 SDK
type MiniProgramApi struct {
	BaseSDK
}

func NewMiniProgramApi(host string) *MiniProgramApi {
	urlHost = host
	return &MiniProgramApi{BaseSDK{c: &clients.HTTPClient{}}}
}

// List 小程序列表
func (api MiniProgramApi) List(c ctx.Context, params *RequestPrograms) (int64, *ProgramList, error) {
	path := "/micro-scrm/v1/mini-program/list"
	url := urlHost + path
	res, err := api.Post(c, url, params)
	if err != nil {
		return 0, nil, err
	}
	result := &ResponseList{
		Data: List{
			List: &ProgramList{},
		},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return 0, nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if err != nil || result.Data.Total == 0 {
		return 0, nil, err
	}
	return result.Data.Total, result.Data.List.(*ProgramList), err
}

// Add 新增小程序
func (api MiniProgramApi) Add(c ctx.Context, params *MiniProgram) (*ID, error) {
	path := "/micro-scrm/v1/mini-program/add"
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

// Edit 修改
func (api MiniProgramApi) Edit(c ctx.Context, params *MiniProgram) (*ID, error) {
	path := "/micro-scrm/v1/mini-program/edit"
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

// Del 删除素材分组
func (api MiniProgramApi) Del(c ctx.Context, id int64) (*Rows, error) {
	params := &struct {
		ID int64 `json:"id"`
	}{
		ID: id,
	}
	path := "/micro-scrm/v1/mini-program/del"
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
