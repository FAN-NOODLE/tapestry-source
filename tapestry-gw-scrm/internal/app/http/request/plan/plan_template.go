package plan

type TemplateAddReq struct {
	BaseReq
	TemplateName        string                 `json:"template_name" binding:"required" comment:"模板名称不能为空"`
	TemplateDesc        string                 `json:"template_desc" binding:"required" comment:"模板描述不能为空"`
	TemplateContent     map[string]interface{} `json:"template_content" comment:"任务的条件内容"`
	TemplateContentType int32                  `json:"template_content_type" comment:"消息类型；1文本；2图片；3语音；4视频；5位置；6链接,7 小程序"`
	AdminId             int64                  `json:"admin_id"`
	GroupSendType       int                    `json:"group_send_type"`
	FinishRate          int                    `json:"finish_rate"`
	ALLowSelect         int                    `json:"allow_select"`
	OpenBackUp          int                    `json:"open_back_up"`
	TaskCategory        int                    `json:"task_category" comment:"0沟通任务 1朋友圈任务"`
	TemplateBack        TemplateBackUp         `json:"template_back"`
}
type TemplateBackUp struct {
	TemplateContent     map[string]interface{} `json:"template_content"`
	TemplateContentType int32                  `json:"template_content_type" comment:"消息类型；1文本；2图片；3语音；4视频；5位置；6链接,7 小程序"`
	StartTime           int64                  `json:"start_time"`
	EndTime             int64                  `json:"end_time"`
	GroupSendType       int                    `json:"group_send_type"`
	ALLowSelect         int                    `json:"allow_select"`
	FinishRate          int                    `json:"finish_rate"`
}

type TemplateUpdateReq struct {
	Id int64 `json:"id"`
	BaseReq
	TemplateName        string                 `json:"template_name" binding:"required" comment:"模板名称不能为空"`
	TemplateDesc        string                 `json:"template_desc" binding:"required" comment:"模板描述不能为空"`
	TemplateContent     map[string]interface{} `json:"template_content" comment:"任务的条件内容"`
	TemplateContentType int32                  `json:"template_content_type" comment:"消息类型；1文本；2图片；3语音；4视频；5位置；6链接,7 小程序"`
	AdminId             int64                  `json:"admin_id"`
	FinishRate          int                    `json:"finish_rate"`
	GroupSendType       int                    `json:"group_send_type"`
	ALLowSelect         int                    `json:"allow_select"`
	OpenBackUp          int                    `json:"open_back_up"`
	TaskCategory        int                    `json:"task_category" comment:"0沟通任务 1朋友圈任务"`
	TemplateBack        TemplateBackUp         `json:"template_back"`
}

type TemplateListsReq struct {
	BaseReq
	SearchType   int    `json:"search_type" json:"search_type"`
	SearchId     int64  `json:"search_id"`
	SearchName   string `json:"search_name"`
	TaskCategory int    `json:"task_category"`
	AdminName    string `json:"admin_name"`
	StartTime    int64  `json:"start_time"`
	EndTime      int64  `json:"end_time"`
	Page         int    `json:"page" json:"page"`
	Size         int    `json:"size" json:"size"`
}

type TemplateDetails struct {
	BaseReq
	Id int64 `json:"id"`
}
