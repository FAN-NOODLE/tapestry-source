package microshop

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

var client = sdk.NewHttpClient(sdk.ParseJson{})

//ActivityPages 活动页
type ActivityPages struct {
	List  []ActivityPage `json:"list"`
	Total int64          `json:"total"`
	Page  int64          `json:"page"`
	Size  int64          `json:"size"`
}

type ActivityPage struct {
	ID           int64  `json:"id"`
	TemplateName string `json:"template_name"`
}

//ActivityPagesResponse 活动页返回
type ActivityPagesResponse struct {
	sdk.ResponseBase
	Data ActivityPages `json:"data"`
}

//ActivityPagesRequest 活动页请求
type ActivityPagesRequest struct {
	TemplateType string `json:"template_type"`
	Page         int64  `json:"page"`
	Size         int64  `json:"size"`
	EnterpriseID int64  `json:"enterprise_id"`
	KeyType      string `json:"key_type"`
	Keywords     string `json:"keywords"`
}

//GetActivitiesPages 获取活动页列表
func GetActivitiesPages(c ctx.Context, request ActivityPagesRequest) (*ActivityPagesResponse, error) {
	var response ActivityPagesResponse
	err := client.Post(c, GetUrl(urlEnum.ActivityPages), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
