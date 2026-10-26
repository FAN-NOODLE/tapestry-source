package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// SensitiveApi 客户 SDK
type SensitiveApi struct {
	BaseSDK
}

func NewSensitive(host string) *SensitiveApi {
	urlHost = host
	return &SensitiveApi{BaseSDK{c: &clients.HTTPClient{}}}
}

// SensitiveList 敏感词列表
func (api *SensitiveApi) SensitiveList(c ctx.Context, sensitive *SensitiveListRequest) (int64, *SensitiveListResponse, error) {
	path := "/micro-scrm/v1/sensitive/list"
	url := urlHost + path
	res, err := api.Post(c, url, sensitive)
	if err != nil {
		return 0, nil, err
	}
	result := &ResponseList{
		Data: List{
			Total: 0,
			Page:  0,
			Size:  0,
			List:  &SensitiveListResponse{},
		},
	}
	_, err = res.Read(clients.JSON(result))
	return result.Data.Total, result.Data.List.(*SensitiveListResponse), err
}

// SensitiveSave 敏感词保存
func (api *SensitiveApi) SensitiveSave(c ctx.Context, sensitive *SensitiveSaveRequest) error {
	path := "/micro-scrm/v1/sensitive/save"
	url := urlHost + path
	res, err := api.Post(c, url, sensitive)
	if err != nil {
		return err
	}
	result := &struct{}{}
	_, err = res.Read(clients.JSON(result))
	return err
}

// SensitiveDelete 敏感词删除
func (api *SensitiveApi) SensitiveDelete(c ctx.Context, sensitive *SensitiveDeleteRequest) error {
	path := "/micro-scrm/v1/sensitive/delete"
	url := urlHost + path
	res, err := api.Post(c, url, sensitive)
	if err != nil {
		return err
	}
	result := &struct{}{}
	_, err = res.Read(clients.JSON(result))
	return err
}
