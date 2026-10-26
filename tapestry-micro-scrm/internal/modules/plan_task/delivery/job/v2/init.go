package v2

import (
	"encoding/json"
	"errors"
	"strings"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/config"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/sftp"
	"micro-scrm/internal/pkg/xjob"
)

const planName = "计划任务"

type planTaskJob struct {
	name              string
	PullSFTPData      *pullSFTPData
	MakeData          *futureTagData
	MakeTag           *externalUserTagging
	PushTask          *taskPushData
	PushCommunityTask *taskCommunityPushData
	SyncTag           *syncTag
}

type runner interface {
	Run(ctx.Context, *xjob.JobHandler) error
}

func NewPlanTaskJob() *planTaskJob {
	return &planTaskJob{
		name:              "计划任务",
		PullSFTPData:      NewPullSFTPData(config.CfgData),
		MakeData:          NewFutureTagData(),
		MakeTag:           NewExternalUserTagging(),
		PushTask:          NewtaskPushData(),
		PushCommunityTask: NewTaskCommunityPushData(),
		SyncTag:           NewSyncTag(),
	}
}

type baseJob struct{}

type tagInfo struct {
	groupName string
	tagName   string
	baType    uint8
}

type aesCbc256 struct {
	key string
	iv  string
}

// 默认标签组
const (
	default_Tag_Group  = `CT-TAG-GROUP`
	group_Split_Flag   = `_Delimit2ed_`
	fixed_FileName_Str = `_ct_segment_node_to_wesocial_`
	ba_Relation_Str    = `+ba+`
	ba_Manager_Str     = `+sba+`
)

const (
	ba_Type_Relation = 1 + iota
	ba_Type_Manager
)

// 获取SFTP配置（只获取配置了计划任务的品牌）
func (job baseJob) getSFTPConfigList(sftpType int) (configList []micro_scrm.BrandConfigSftpDetail, err error) {
	dataList, dataErr := orm.Service.PlanTask.GetBrandSFTPConfig(sftpType)
	if dataErr != nil {
		err = dataErr
		return
	}

	for _, data := range dataList {
		var sftpConfig sftp.SSHConfig
		var pathConfig micro_scrm.SFTPPathConfig
		_ = json.Unmarshal([]byte(data.SftpConfig), &sftpConfig)
		_ = json.Unmarshal([]byte(data.SftpPath), &pathConfig)

		configList = append(configList, micro_scrm.BrandConfigSftpDetail{
			ConfigDetail: sftpConfig,
			PathConfig:   pathConfig,
			BrandConfigSftpModel: micro_scrm.BrandConfigSftpModel{
				BrandCode:  data.BrandCode,
				BrandName:  data.BrandName,
				SftpConfig: data.SftpConfig,
				SftpPath:   data.SftpPath,
				MicroScrmBaseModel: model.MicroScrmBaseModel{
					EnterpriseID: data.EnterpriseID,
					BaseModel: model.BaseModel{
						ID:        data.ID,
						CreatedAt: data.CreatedAt,
						UpdatedAt: data.UpdatedAt,
					},
				},
			},
		})
	}

	return
}

// 拆分SFTP文件名字为标签信息
func (job baseJob) splitFileName(fileName string) (info tagInfo, err error) {
	if strings.Contains(fileName, ba_Manager_Str) {
		info.baType = ba_Type_Manager
	} else {
		info.baType = ba_Type_Relation
	}

	// 自定义标签组
	// lrl_ct_segment_node_to_wesocial_2109会员日+低客单用户_Delimit2ed_2109会员日+低客单用户_20210918180113564.csv
	if strings.Contains(fileName, group_Split_Flag) {
		splitList := strings.Split(fileName, group_Split_Flag)
		if len(splitList) >= 2 && strings.Contains(splitList[1], "+") {
			// 拆分标签
			splitTagGroupList := strings.Split(splitList[1], "+")
			// 标签组
			if splitTagGroupList[0] != "" {
				info.groupName = splitTagGroupList[0]
			} else {
				err = errors.New("文件名定义标签组名称异常")
				return
			}
			// 标签
			if len(splitTagGroupList) < 2 || splitTagGroupList[1] == "" {
				err = errors.New("文件名定义标签名称异常")
				return
			}
			splitTagList := strings.Split(splitTagGroupList[1], "_")
			if splitTagList[0] != "" {
				info.tagName = splitTagList[0]
				return
			} else {
				err = errors.New("文件名定义标签名称异常")
				return
			}
		} else {
			tagStrInfo := strings.Split(splitList[1], "_")
			if len(tagStrInfo) >= 2 {
				info.tagName = strings.Join(tagStrInfo[0:len(tagStrInfo)-1], "_")
				info.groupName = default_Tag_Group
				return
			} else {
				err = errors.New("文件名定义标签名称异常")
				return
			}
		}
	}

	// 默认标签组
	// hr_ct_segment_node_to_wesocial_白金卡入群邀请_20210902110201400.csv
	if strings.Contains(fileName, fixed_FileName_Str) {
		splitList := strings.Split(fileName, fixed_FileName_Str)
		if len(splitList) >= 2 && strings.Contains(splitList[1], "_") {
			splitTagList := strings.Split(splitList[1], "_")

			info.tagName = strings.Join(splitTagList[0:len(splitTagList)-1], "_")
			info.groupName = default_Tag_Group

			return
		}
	}

	err = errors.New("文件名称不符合规范")
	return
}
