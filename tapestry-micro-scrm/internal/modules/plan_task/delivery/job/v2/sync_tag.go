package v2

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gorm.io/gorm"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise/plan_task"
	"micro-scrm/internal/orm/service"
	planTaskSftp "micro-scrm/internal/orm/service/plan_task"
	"micro-scrm/internal/pkg/xjob"
	"sync"
)

/**
 * 同步计划任务标签 到 enterprise_user 表
 * 根据plan_task_downloaded_file.is_extract = 3，
 * id对应的task_condition_data_log.tag_id|union_id, 更新enterprise_user.tags,进行批量更新
 **/

var _ runner = externalUserTagging{}

type syncTag struct {
	baseJob
	name         string
	tagDataLimit int
}

func NewSyncTag() *syncTag {
	return &syncTag{
		name:         "计划任务::同步用户标签",
		tagDataLimit: 2000,
	}
}

func (job syncTag) Run(ctx ctx.Context, handler *xjob.JobHandler) error {
	// 获取配置了计划任务的企业
	configList, configErr := job.getSFTPConfigList(planTaskSftp.SFTP_TYPE_PLAN_TASK)
	if configErr != nil {
		log.Error(ctx).Msgf("%s::获取sftp企业配置异常: %v", job.name, configErr)
		return configErr
	}

	// 每个品牌进行并发
	var wg sync.WaitGroup
	for _, config := range configList {
		wg.Add(1)
		go func(enterpriseID int64) {
			defer wg.Done()
			job.fileData(ctx, enterpriseID)
		}(config.EnterpriseID)
	}
	wg.Wait()

	return nil
}

// 分文件并发
func (job syncTag) fileData(ctx ctx.Context, enterpriseID int64) error {
	// 获取需要处理的文件数据
	dataList := []plan_task.PlanTaskDownloadedFileModel{}
	if dataListErr := orm.Service.PlanTask.Scopes(model.BranchTable(&plan_task.PlanTaskDownloadedFileModel{}, enterpriseID), service.ScopeCond_EQDelete).Where(map[string]interface{}{
		"is_extract": plan_task.PTDF_IsExtract_Tag,
	}).Find(&dataList).Error; dataListErr != nil {
		log.Error(ctx).Msgf("%s::获取要更新的文件数据失败: %v", job.name, dataListErr)
	}

	var wg sync.WaitGroup
	for _, data := range dataList {
		wg.Add(1)
		go func(data plan_task.PlanTaskDownloadedFileModel) {
			defer wg.Done()

			job.syncTag(ctx, enterpriseID, data)
		}(data)
	}
	wg.Wait()

	return nil
}

// 同步标签
func (job syncTag) syncTag(ctx ctx.Context, enterpriseID int64, fileData plan_task.PlanTaskDownloadedFileModel) {
	// 获取标签信息
	tag := tagMsg{}
	if tagErr := orm.Service.PlanTask.Table(fmt.Sprintf("%s AS T", micro_scrm.WorkExternalTagModel{}.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %s G ON G.`id` = T.`group_id`", micro_scrm.WorkExternalTagGroupModel{}.TableName())).
		Select("T.id, T.name, T.tag_id").
		Where(map[string]interface{}{
			"T.deleted_at": 0,
			"T.name":       fileData.TagName,
			"G.deleted_at": 0,
			"G.group_name": fileData.GroupTagName,
		}).
		Scan(&tag).Error; tagErr != nil || 0 == tag.ID {
		return
	}

Label:
	for {
		unionIDList := []string{}
		if unionIDListErr := orm.Service.PlanTask.Scopes(model.BranchTable(&plan_task.TaskConditionDataLogModel{}, enterpriseID), service.ScopeCond_EQDelete).
			Model(&plan_task.TaskConditionDataLogModel{}).
			Where(map[string]interface{}{
				"is_success": plan_task.TCDL_IsSuccess_Yes,
				"is_has_tag": plan_task.TCDL_IsHasTag_No,
				"tag_id":     tag.ID,
			}).
			Limit(job.tagDataLimit).Pluck("union_id", &unionIDList).Error; unionIDListErr != nil {
			log.Error(ctx).Msgf("%s::获取需要更新的标签的用户union_id失败: %v", job.name, unionIDListErr)
			continue
		}

		if 0 == len(unionIDList) {
			orm.Service.PlanTask.Scopes(model.BranchTable(&plan_task.PlanTaskDownloadedFileModel{}, enterpriseID), service.ScopeCond_EQDelete).
				Model(&plan_task.PlanTaskDownloadedFileModel{}).Where("id = ?", fileData.ID).Update("is_extract", plan_task.PTDF_IsExtract_SyncTag)

			break Label
		}

		// 更新标签
		euidList := []int64{}
		if idListErr := orm.Service.PlanTask.Scopes(model.BranchTable(&micro_scrm_enterprise.WxUserCorpExternalModel{}, enterpriseID)).
			Model(&micro_scrm_enterprise.WxUserCorpExternalModel{}).
			Where("unionid IN (?)", unionIDList).
			Pluck("euid", &euidList).Error; idListErr != nil {
			log.Error(ctx).Msgf("%s::获取需要更新的标签的用户euid失败: %v", job.name, idListErr)
		}

		orm.Service.PlanTask.Transaction(func(tx *gorm.DB) error {
			if updateErr := orm.Service.PlanTask.Scopes(model.BranchTable(&micro_scrm_enterprise.EnterpriseUserModel{}, enterpriseID)).
				Model(&micro_scrm_enterprise.EnterpriseUserModel{}).
				Where("id in (?)", euidList).
				Update("tags", gorm.Expr(fmt.Sprintf("IF(IFNULL(tags, \"\") = \"\", JSON_ARRAY(%d),JSON_ARRAY_APPEND(tags, '$',%d))", tag.ID, tag.ID))).Error; updateErr != nil {
				return updateErr
			}
			if updateErr := orm.Service.PlanTask.Scopes(model.BranchTable(&plan_task.TaskConditionDataLogModel{}, enterpriseID)).
				Model(&plan_task.TaskConditionDataLogModel{}).
				Where("union_id in (?)", unionIDList).Where(map[string]interface{}{
				"is_success": plan_task.TCDL_IsSuccess_Yes,
				"is_has_tag": plan_task.TCDL_IsHasTag_No,
				"tag_id":     tag.ID,
			}).Update("is_has_tag", plan_task.TCDL_IsHasTag_Yes).Error; updateErr != nil {
				return updateErr
			}

			return nil
		})
	}
}
