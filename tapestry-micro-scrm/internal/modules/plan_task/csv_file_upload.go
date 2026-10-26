package plan_task

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	config2 "micro-scrm/internal/config"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise/plan_task"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"
	"micro-scrm/internal/pkg/workwechat"
)

func SaveUploadInfo(data DownLoadFileInfo, enterpriseID int64) (ID int64, err error) {
	var newInfo plan_task.PlanTaskDownloadedFileModel
	newInfo.CreatedAt = time.Now().Unix()
	newInfo.UpdatedAt = time.Now().Unix()
	newInfo.FileName = data.FileName
	newInfo.FileType = plan_task.PTDF_FileType_WS

	newInfo.GroupTagName = data.TagGroupName
	newInfo.TagName = data.TagName
	newInfo.FileURL = data.FileUrl
	if baType, err := strconv.ParseUint(data.BAType, 10, 64); err != nil {
		return 0, errors.New("BA类型有误，请确认！")
	} else {
		if baType != 1 && baType != 2 {
			return 0, errors.New("BA类型有误，请确认！")
		}
		newInfo.BaType = uint8(baType)
	}

	id, err := orm.Service.PlanTask.CreateUploadInfo(enterpriseID, newInfo)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func CSVUploadLogin(data GetFileInfo, enterpriseID int64) (err error) {
	var newInfo = make([]micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel, 0)
	for _, field := range data.File {
		var newData micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel
		newData.CreatedAt = time.Now().Unix()
		newData.UpdatedAt = time.Now().Unix()
		newData.MemberCode = field.Ccode
		newData.FileID = data.ID
		newData.CronFileName = data.FileName
		KEY := config2.CfgData.CampaignToolService.AES_CBC_256.Key
		IV := config2.CfgData.CampaignToolService.AES_CBC_256.Iv
		if data.UnionIDEncrypt == "1" {
			if field.Unionid != "" {
				if union_id, unionErr := base64.StdEncoding.DecodeString(field.Unionid); unionErr == nil {
					decrypt, unionErr := common.AES_CBC_Decrypt2(union_id, []byte(KEY), []byte(IV))
					if unionErr != nil {
						orm.Service.PlanTask.UpdatedFileInfo(enterpriseID, data.ID)
						return unionErr
					}
					newData.UnionID = string(decrypt)
				} else {
					return errors.New("密文解析异常，请确认字段是否加密！")
				}
			} else {
				newData.UnionID = ""
			}
		} else {
			if field.Unionid != "" {
				newData.UnionID = field.Unionid
			} else {
				newData.UnionID = ""
			}
		}
		if field.Ccode == "" && newData.UnionID == "" {
			continue
		}
		newData.GradeCode = field.GradeCode
		newData.AdminBA = field.AdminBA
		newInfo = append(newInfo, newData)
	}
	if len(newInfo) < 1 {
		return nil
	}
	if err := orm.Service.PlanTask.CreateCSVInfo(enterpriseID, newInfo); err != nil {
		return err
	}
	return nil
}

func GetFileNameDetection(enterpriseID int64, fileName string) (result int64, err error) {
	return orm.Service.PlanTask.GetFileInfo(enterpriseID, fileName)
}

func GetUpdateInfo(enterpriseID, ID int64) (err error) {
	return orm.Service.PlanTask.UpdateISExtract(enterpriseID, ID)
}

// 创建企微标签
func CreateTag(ctx ctx.Context, enterpriseID int64, info TagInfo) error {
	tagGroup, tagGroupIsExist, tagGroupErr := orm.Service.Tag.CheckTagGroupName(enterpriseID, info.GroupName)
	if tagGroupErr != nil {
		return tagGroupErr
	}
	// 插入新的标签
	obj, err := workwechat.NewWorkWeChat(ctx, enterpriseID)
	if err != nil {
		return err
	}

	if true != tagGroupIsExist {
		tagResult, tagErr := obj.AddTag(work.AddCorpTagRequest{
			GroupName: info.GroupName,
			TagList: []work.AddCorpTag{
				work.AddCorpTag{
					Name: info.TagName,
				},
			},
		})

		if tagErr != nil {
			return errors.New(fmt.Sprintf("写入标签组信息错误: %v", tagErr))
		}

		// 跟新标签组，标签
		tagGroup, _, tagGroupErr = orm.Service.Tag.CreateTagGroup2(micro_scrm.WorkExternalTagGroupModel{
			MicroScrmBaseModel: model.MicroScrmBaseModel{
				BaseModel: model.BaseModel{
					ID: uid.ID(ctx),
				},
				EnterpriseID: enterpriseID,
			},
			GroupName:  info.GroupName,
			GroupID:    tagResult.TagGroup.GroupID,
			Order:      tagResult.TagGroup.Order,
			CreateTime: int(tagResult.TagGroup.CreateTime),
		})

		if tagGroupErr != nil {
			return errors.New(fmt.Sprintf("写入标签组信息错误: %v", tagGroupErr))
		}

		var insertTagData []micro_scrm.WorkExternalTagModel
		for _, tag := range tagResult.TagGroup.TagList {
			insertTagData = append(insertTagData, micro_scrm.WorkExternalTagModel{
				MicroScrmBaseModel: model.MicroScrmBaseModel{
					BaseModel: model.BaseModel{
						ID: uid.ID(ctx),
					},
					EnterpriseID: enterpriseID,
				},
				Name:       tag.Name,
				TagID:      tag.ID,
				CreateTime: int(tag.CreateTime),
				Order:      tag.Order,
				GroupID:    tagGroup.ID,
				Type:       micro_scrm.WET_Type_CT,
			})
		}
		_, _, tagErr = orm.Service.Tag.CreateTag(insertTagData)
		if tagErr != nil {
			return errors.New(fmt.Sprintf("写入标签组信息错误: %v", tagGroupErr))
		}
	}

	_, tagIsExist, tagErr := orm.Service.Tag.CheckTagName(enterpriseID, tagGroup.ID, info.TagName)
	if tagErr != nil {
		return errors.New(fmt.Sprintf("写入标签组信息错误: %v", tagGroupErr))
	}
	if !tagIsExist {
		tagResult, tagResultErr := obj.AddTag(work.AddCorpTagRequest{
			GroupID:   tagGroup.GroupID,
			GroupName: info.GroupName,
			TagList: []work.AddCorpTag{
				work.AddCorpTag{
					Name: info.TagName,
				},
			},
		})
		if tagResultErr != nil {
			return errors.New(fmt.Sprintf("创建企微对象失败: %v", tagGroupErr))
		}
		// 更新标签
		var insertTagData []micro_scrm.WorkExternalTagModel
		for _, tag := range tagResult.TagGroup.TagList {
			insertTagData = append(insertTagData, micro_scrm.WorkExternalTagModel{
				MicroScrmBaseModel: model.MicroScrmBaseModel{
					BaseModel: model.BaseModel{
						ID: uid.ID(ctx),
					},
					EnterpriseID: enterpriseID,
				},
				Name:       tag.Name,
				TagID:      tag.ID,
				CreateTime: int(tag.CreateTime),
				Order:      tag.Order,
				GroupID:    tagGroup.ID,
				Type:       micro_scrm.WET_Type_CT,
			})
		}
		_, _, tagErr = orm.Service.Tag.CreateTag(insertTagData)
		if tagErr != nil {
			return errors.New(fmt.Sprintf("写入标签组信息错误: %v", tagGroupErr))
		}
	}

	return nil
}
