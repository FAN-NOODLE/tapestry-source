package templateReq

type PageVO struct {
	Page int `json:"page" binding:"min=1" comment:"页码"`
	Size int `json:"size" binding:"min=10,max=99999999" comment:"分页数"`
}

type RequestBase struct {
	EnterpriseID int64 `json:"enterprise_id" binding:"required"`
}

type TemplateBaseInfo struct {
	ID                 int      `json:"id"`
	TemplateName       string   `json:"template_name" binding:"required,min=1,max=20" comment:"模板名称"`
	TemplatePictureURL []string `json:"template_picture_url" binding:"required" comment:"背景图"`
	TemplateContent    string   `json:"template_content" binding:"required,min=1,max=50" comment:"文案"`
}

type TemplateInfo struct {
	TemplateBaseInfo
	UpdatedAt int64 `json:"updated_at"`
}

type ReqList struct {
	RequestBase
	PageVO
	TemplateName string `json:"template_name"`
}

type ReqCreate struct {
	RequestBase
	TemplateBaseInfo
}

type RequestInfo struct {
	RequestBase
	ID int `json:"id" binding:"required"`
}

type ResponseList struct {
	Total int64 `json:"total"`
	List  []TemplateInfo
}

type ResponseCheck struct {
	Check bool `json:"check"`
}
