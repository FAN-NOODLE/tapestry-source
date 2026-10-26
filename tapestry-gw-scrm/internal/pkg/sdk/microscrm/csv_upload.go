package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

const (
	//计划任务标签的客户导入匹配标签
	URL_CSV_UPLOAD        = "%s/micro-scrm/v1/plan/csvFile/upload"
	URL_SAVE_UPLOAD_INFO  = "%s/micro-scrm/v1/plan/upload/saveInfo"
	URL_FILE_DETECTION    = "%s/micro-scrm/v1/plan/upload/fileDetection"
	URL_UPDATE_IS_EXTRACT = "%s/micro-scrm/v1/plan/upload/update/is_extract"

	//客户分组的客户名单导入
	USER_GROUP_FILE_DETECTION    = "%s/micro-scrm/v1/user/group/fileDetection"
	USER_GROUP_SAVE_UPLOAD_INFO  = "%s/micro-scrm/v1/user/group/saveInfo"
	USER_GROUP_UPDATE_IS_EXTRACT = "%s/micro-scrm/v1/user/group/update/is_extract"

	//计划任务客户导入
	USER_PLAN_IMPORT_FILE_DETECTION    = "%s/micro-scrm/v1/plan/import/fileDetection"
	USER_PLAN_IMPORT_SAVE_UPLOAD_INFO  = "%s/micro-scrm/v1/plan/import/saveInfo"
	USER_PLAN_IMPORT_UPDATE_IS_EXTRACT = "%s/micro-scrm/v1/plan/import/update/is_extract"
)

func UploadCsvFile(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(URL_CSV_UPLOAD), request)
}

func DownloadFileInfo(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(URL_SAVE_UPLOAD_INFO), request)
}

func FileDetection(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(URL_FILE_DETECTION), request)
}

func UpdateISExtract(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(URL_UPDATE_IS_EXTRACT), request)
}

func UserGroupFileDetection(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(USER_GROUP_FILE_DETECTION), request)
}

func UserGroupDownloadFileInfo(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(USER_GROUP_SAVE_UPLOAD_INFO), request)
}

func UserGroupUpdateISExtract(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(USER_GROUP_UPDATE_IS_EXTRACT), request)
}



func PlanImportFileDetection(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(USER_PLAN_IMPORT_FILE_DETECTION), request)
}

func PlanImportDownloadFileInfo(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(USER_PLAN_IMPORT_SAVE_UPLOAD_INFO), request)
}


func PlanImportUpdateISExtract(ctx ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(USER_PLAN_IMPORT_UPDATE_IS_EXTRACT), request)
}
