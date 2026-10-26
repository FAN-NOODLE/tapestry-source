package openReq

import planReq "micro-scrm/internal/app/http/request/plan"

type PlanTemplateListsReq struct {
	EnterpriseId int64 `json:"enterprise_id"`
}

type AddPlanReq struct {
	PlanName       string `json:"plan_name" binding:"required"  comment:"任务名称"`
	PlanRemark     string `json:"plan_remark" binding:"required"  comment:"任务描述"`
	PlanTemplateId int64  `json:"plan_template_id"  comment:"任务模板id"`
	PlanSendType   int    `json:"plan_send_type"  comment:"任务推送类型，1 客户和keysa关系，2 好友关系"`
	PlanStartTime  string `json:"plan_start_time"  comment:"任务开始时间"`
	PlanEndTime    string `json:"plan_end_time" comment:"任务结束时间"`
	PlanSendTime   string `json:"plan_send_time"  comment:"任务的发送时间"`
	PlanType       int    `json:"plan_type"  comment:"1持续任务，2定时任务"`
}

type CreatePlanByContentReq struct {
	PlanName      string `json:"plan_name" binding:"required"`
	PlanRemark    string `json:"plan_remark" binding:"required"`
	PlanStartTime string `json:"plan_start_time" binding:"required"`
	PlanEndTime   string `json:"plan_end_time" binding:"required"`
}

type TaskCrowdReq struct {
	TaskId         int64                 `json:"task_id" binding:"required" `
	CrowdId        string                `json:"crowd_id" binding:"required" `
	PlanSendType   int                   `json:"plan_send_type" binding:"required" `
	PlanType       int                   `json:"plan_type" binding:"required" `
	PlanSendTime   string                `json:"plan_send_time" binding:"required"  comment:"任务的发送时间"`
	TaskStartTime  string                `json:"task_start_time" binding:"required" `
	TaskEndTime    string                `json:"task_end_time" binding:"required" `
	TaskRemark     string                `json:"task_remark" binding:"required" comment:"子任务的描述"`
	TaskName       string                `json:"task_name" binding:"required" comment:"子任务名称"`
	PlanTemplateId int64                 `json:"plan_template_id" binding:"required"  comment:"任务模板id"`
	FlowId         string                `json:"flowId" comment:"自动流ID"`
	Campaign       string                `json:"campaign" comment:"营销活动编码"`
	TaskTag        string                `json:"task_tag" comment:"朋友圈任务标签"`
	Crowd          []planReq.UnionIDData `json:"crowd" comment:"人群包"`
}

type CreateTaskReq struct {
	PlanId              int64                  `json:"plan_id" binding:"required"`
	CrowdId             string                 `json:"crowd_id" binding:"required"`
	TaskStartTime       string                 `json:"task_start_time" binding:"required"`
	TaskEndTime         string                 `json:"task_end_time" binding:"required"`
	TaskRemark          string                 `json:"task_remark" binding:"required"`
	TaskName            string                 `json:"task_name" binding:"required"`
	FlowId              string                 `json:"flowId"`
	Campaign            string                 `json:"campaign"`
	TaskTag             string                 `json:"task_tag"`
	TemplateContent     map[string]interface{} `json:"template_content"`
	TemplateContentType int32                  `json:"template_content_type"`
	TaskCategory        int                    `json:"task_category"`
	GroupSendType       int                    `json:"group_send_type"`
	FinishRate          int                    `json:"finish_rate"`
	ALLowSelect         int                    `json:"allow_select"`
	Crowd               []planReq.UnionIDData  `json:"crowd"`
}

type PushTaskCrowdReq struct {
	TaskId  int64                 `json:"task_id" binding:"required"`
	CrowdId string                `json:"crowd_id" binding:"required"`
	Crowd   []planReq.UnionIDData `json:"crowd" binding:"required"`
}

type ConditionUserInfo struct {
	CustomerId         string `json:"customer_id" comment:"用户的会员卡号，plan_send_type=1时，unionid和customer_id两个字段必须传1个"`
	Unionid            string `json:"unionid" comment:"用户的unionid，plan_send_type=2时，unionid必传"`
	SaCode             string `json:"sa_code" comment:"导购的工号，pland_send_type=1时，sa_code字段必传"`
	WorkExternalUserId string `json:"work_external_user_id"`
}
