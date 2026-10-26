package drainage

import "micro-scrm/internal/orm/model/micro_scrm/welcome"

type ExternalTypeIdRequest struct {
	TypeId       int64 `json:"type_id" required:"类型ID不能为空"`
	EnterpriseId int64 `json:"enterprise_id" required:"企业ID不能为空"`
}

type ExternalTypeAllRequest struct {
	EnterpriseId int64 `json:"enterprise_id" required:"企业ID不能为空"`
}

type ExternalWelcomeRequest struct {
	ID               int64                    `json:"id"`
	EnterpriseId     int64                    `json:"enterprise_id"`
	Content          []map[string]interface{} `json:"content"`
	ConditionContent map[string][]string      `json:"condition_content" comment:"条件集合json格式"`
}

type WelcomeDelRequest struct {
	EnterpriseId int64 `json:"enterprise_id" required:"企业ID不能为空"`
	ID           int64 `json:"id" required:"数据ID不能为空"`
}

type WelcomeDetailNewRequest struct {
	EnterpriseId int64 `json:"enterprise_id" required:"企业ID不能为空"`
	ID           int64 `json:"id" required:"数据ID不能为空"`
}

type WelcomeAllRequest struct {
	EnterpriseId int64  `json:"enterprise_id" required:"企业ID不能为空"`
	Page         int    `json:"page"`
	Size         int    `json:"size"`
	KeyWords     string `json:"key_words"`
}

type WelcomeAllResponse struct {
	welcome.DrainageWelcomesModel
	ContentText      []map[string]interface{} `json:"content_text"`
	TypeName         string                   `json:"type_name"`
	ConditionContent map[string][]any         `json:"condition_content"`
}
