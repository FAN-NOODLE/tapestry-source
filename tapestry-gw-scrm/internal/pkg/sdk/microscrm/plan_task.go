package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	planTask "gw-scrm/internal/modules/plan_task"
	"gw-scrm/internal/pkg/sdk"
	"mime/multipart"
)

const (
	URL_PLANTASK_INFO         = "%s/micro-scrm/v1/plan_tasks/info"
	URL_PLANTASK_ADD          = "%s/micro-scrm/v1/plan_tasks/add"
	URL_PLANTASK_UPDATE       = "%s/micro-scrm/v1/plan_tasks/update"
	URL_PLANTASK_DELETE       = "%s/micro-scrm/v1/plan_tasks/delete"
	URL_PLAN_LIST             = "%s/micro-scrm/v1/plan/list"
	URL_TASK_LIST             = "%s/micro-scrm/v1/task/list"
	URL_PLANTASK_PlanTaskInfo = "%s/micro-scrm/v1/plan_tasks/info"
	URL_PLAN_PLAN_ADD         = "%s/micro-scrm/v1/plan_task/add_plan"
	URL_PLAN_PLAN_EDIT        = "%s/micro-scrm/v1/plan_task/plan/edit"
	URL_PLAN_PLAN_REMOVE      = "%s/micro-scrm/v1/plan_task/remove"
	URL_PLAN_PLAN_COPY        = "%s/micro-scrm/v1/plan_task/plan/copy"
	URL_PLAN_PLAN_STATUS      = "%s/micro-scrm/v1/plan_task/plan/update_status"

	ULR_PlANTASK_CONDITIONS_COMPLETION      = "%s/micro-scrm/v1/plan_tasks/conditions/completion"
	ULR_PlANTASK_CONDITIONS_USER_COMPLETION = "%s/micro-scrm/v1/plan_tasks/conditions/user_completion"

	ULR_EXPORT_USER_DATA       = "%s/micro-scrm/v1/activity/export-user-data"
	ULR_ACTIVITY_LIST          = "%s/micro-scrm/v1/activity/get-type-list"
	ULR_GET_USER_LIST          = "%s/micro-scrm/v1/activity/get-user-list"
	ULR_GET_ACTIVITY_USER_INFO = "%s/micro-scrm/v1/activity/get-user-info"
	ULR_VERIFY_USER_INFO       = "%s/micro-scrm/v1/activity/verify-user"
	URL_CSV_FILE_UPLOAD        = "%s/micro-scrm/v1/plan/csvFile/upload"

	URL_CUSTOMER_DATA          = "%s/micro-scrm/v1/plan/subtask/customer"
	URL_CUSTOMER_LIST_EXPORT   = "%s/micro-scrm/v1/plan/subtask/customer/list_export"
	URL_CUSTOMER_DETAIL_EXPORT = "%s/micro-scrm/v1/plan/subtask/customer/detail_export"
	URL_GROUP_CHAT_DATA        = "%s/micro-scrm/v1/plan/subtask/group_task"
	URL_GROUP_CHAT_LIST_DATA   = "%s/micro-scrm/v1/plan/subtask/group_task_list"
	URL_CUSTOMER_LIST_DATA     = "%s/micro-scrm/v1/plan/subtask/customer_task_list"
	URL_GROUP_CHAT_DATA_EXPORT = "%s/micro-scrm/v1/plan/subtask/group_task_export"
	URL_SUBTASK_UPLOAD_USER    = "%s/micro-scrm/v1/plan/subtask/upload-users"
)

func PlanTaskInfo(ctx ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(URL_PLANTASK_INFO), request)
	return response, err
}
func PlanTasksAdd(ctx ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(URL_PLANTASK_ADD), request)
	return response, err
}
func PlanTasksUpdate(ctx ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(URL_PLANTASK_UPDATE), request)
	return response, err
}
func PlanTasksDelete(ctx ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(URL_PLANTASK_DELETE), request)
	return response, err
}
func PlanList(ctx ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(URL_PLAN_LIST), request)
	return response, err
}
func TaskList(ctx ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(URL_TASK_LIST), request)
	return response, err
}

func PlanAdd(ctx ctx.Context, info *planTask.PlanRequest) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(URL_PLAN_PLAN_ADD), info)
}

func PlanEdit(ctx ctx.Context, info *planTask.ReceivePlanModel) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(URL_PLAN_PLAN_EDIT), info)
}

func PlanRemove(ctx ctx.Context, planId planTask.PlanId) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(URL_PLAN_PLAN_REMOVE), planId)
}

func PlanCopy(ctx ctx.Context, plan *planTask.ReceivePlanModel) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(URL_PLAN_PLAN_COPY), plan)
}
func PlanStatus(ctx ctx.Context, info *planTask.PlanStatus) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(URL_PLAN_PLAN_STATUS), info)
}

func PlanTasksConditionsCompletion(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(ULR_PlANTASK_CONDITIONS_COMPLETION), request)
}
func PlanTasksConditionsUserCompletion(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(ULR_PlANTASK_CONDITIONS_USER_COMPLETION), request)
}

func ExportUserData(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(ULR_EXPORT_USER_DATA), request)
}

func GetActivityList(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(ULR_ACTIVITY_LIST), request)
}

func GetUserList(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(ULR_GET_USER_LIST), request)
}

func GetActivityUserInfo(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(ULR_GET_ACTIVITY_USER_INFO), request)
}

func VerifyActivityUserInfo(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(ULR_VERIFY_USER_INFO), request)
}

func CSVUpload(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(ULR_VERIFY_USER_INFO), request)
}

func CustomerData(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(URL_CUSTOMER_DATA), request)
}

func CustomerListDataExport(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(URL_CUSTOMER_LIST_EXPORT), request)
}
func CustomerDetailDataExport(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(URL_CUSTOMER_DETAIL_EXPORT), request)
}

func GroupChatTaskData(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(URL_GROUP_CHAT_DATA), request)
}

func GroupChatTaskListData(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(URL_GROUP_CHAT_LIST_DATA), request)
}

func CustomerTaskListData(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(URL_CUSTOMER_LIST_DATA), request)
}

func GroupChatTaskDataExport(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(URL_GROUP_CHAT_DATA_EXPORT), request)
}

func UploadSubTaskUsers(ctx ctx.Context, params map[string]string, upFileName string, files ...*multipart.FileHeader) (interface{}, error) {
	return sdk.SendFormData(ctx, GetUrl(URL_SUBTASK_UPLOAD_USER), params, upFileName, files...)
}
