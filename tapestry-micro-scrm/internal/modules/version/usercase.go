package version

import (
	"encoding/json"
	"errors"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm/version"
	"micro-scrm/internal/pkg/uid"
)

const (
	YYYY_MM_DD = "2006-01-02"
)

func GetCurrentVersionInfo(ctx ctx.Context, param GetCurrentVersionInfoRequest) (data GetCurrentVersionInfoResponse, err error) {
	switch param.Type {
	case CurrentVersionInfoRequest_Type_Get:
		data, err = getCurrentVersionInfo(param)
		data.Type = CurrentVersionInfoRequest_Type_Get
	case CurrentVersionInfoRequest_Type_Record:
		// 检测用户ID
		data, err = getCurrentVersionInfoAndRecord(ctx, param)
		data.Type = CurrentVersionInfoRequest_Type_Record
	}

	if err != nil {
		log.Info(ctx).Msgf("获取版本信息异常，request: %v, error: %v", param, err)
		return
	}

	return
}

// 版本信息查询类型 0
func getCurrentVersionInfo(param GetCurrentVersionInfoRequest) (data GetCurrentVersionInfoResponse, err error) {
	versionInfo, err := orm.Service.Version.GetCurrentVersionInfo()
	if err != nil {
		return
	}

	data.VersionInfo.Them = versionInfo.Theme
	data.VersionInfo.Number = versionInfo.VersionNumber
	data.VersionInfo.Date = time.Unix(versionInfo.CreatedAt, 0).Format(YYYY_MM_DD)

	var versionDetail []VersionDetail
	if jsonErr := json.Unmarshal([]byte(versionInfo.VersionInfo), &versionDetail); jsonErr == nil {
		data.VersionInfo.Detail = versionDetail
	}

	return
}

// 版本信息查询类型 1
func getCurrentVersionInfoAndRecord(ctx ctx.Context, param GetCurrentVersionInfoRequest) (data GetCurrentVersionInfoResponse, err error) {
	versionInfo, err := orm.Service.Version.GetCurrentVersionInfo()
	if err != nil {
		return
	}

	// 检测企业ID、用户ID
	if 0 == param.EnterpriseID {
		err = errors.New("企业ID不能为空")
		return
	}

	// 检测是否已经有版本信息弹出记录
	isHas, isHasErr := orm.Service.Version.IsHasLatestPromptRecord(version.VersionPromptLogModel{
		MicroScrmBaseModel: model.MicroScrmBaseModel{
			EnterpriseID: param.EnterpriseID,
		},
		AdminID:         param.AdminID,
		VersionRecordID: versionInfo.ID,
	})

	if isHas && isHasErr == nil {
		data.IsHasShow = true
	}

	// 记录信息
	if err := orm.Service.Version.RecordVersionPrompt(version.VersionPromptLogModel{
		MicroScrmBaseModel: model.MicroScrmBaseModel{
			BaseModel: model.BaseModel{
				ID: uid.ID(ctx),
			},
			EnterpriseID: param.EnterpriseID,
		},
		AdminID:         param.AdminID,
		VersionRecordID: versionInfo.ID,
	}); err != nil {
		log.Info(ctx).Msgf("记录用户浏览版本信息异常，request: %v, error: %v", param, err)
	}
	vInfo, err := getCurrentVersionInfo(param)
	if err != nil {
		return
	}
	data.VersionInfo = vInfo.VersionInfo
	data.Type = param.Type

	return data, nil
}
