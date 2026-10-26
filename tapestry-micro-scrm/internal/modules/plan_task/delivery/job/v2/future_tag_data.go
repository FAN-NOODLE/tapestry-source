package v2

import (
	"fmt"
	"sync"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm/work_external_user"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise/plan_task"
	planTaskSftp "micro-scrm/internal/orm/service/plan_task"
	"micro-scrm/internal/pkg/xjob"
)

var _ runner = futureTagData{}

const (
	DB_EveryTime_Num = 2000
	DB_Channel_Num   = 1000
	Read_Data_Nun    = 500
)

const (
	PlanTask_TagDataType_Member = 1 + iota
	PlanTask_TagDataType_UnionID
)

type userTagData struct {
	ExternalUserID     int64  `json:"external_user_id"`
	UnionID            string `json:"union_id"`
	WorkExternalUserID string `json:"work_external_user_id"`
	Mobile             string `json:"mobile"`
	MemberID           int64  `json:"member_id"`
	UserID             string `json:"user_id"`
	RelationID         int64  `json:"relation_id"`

	TagID   int64  `json:"tag_id"`
	TagName string `json:"tag_name"`
	WwTagID string `json:"ww_tag_id"`
	IsUnion int    `json:"is_union"`
}

type futureTagData struct {
	baseJob
	name string
}

func NewFutureTagData() *futureTagData {
	return &futureTagData{
		name: "即将打标签的数据",
	}
}

func (job futureTagData) Run(ctx ctx.Context, handler *xjob.JobHandler) error {
	// 获取配置了计划任务的企业
	configList, configErr := job.getSFTPConfigList(planTaskSftp.SFTP_TYPE_PLAN_TASK)
	if configErr != nil {
		log.Info(ctx).Msgf("%s::获取sftp企业配置异常: %v", job.name, configErr)
		return configErr
	}

	var wg sync.WaitGroup
	for _, config := range configList {
		wg.Add(1)
		go func(enterpriseID int64) {
			defer wg.Done()

			job.createBrandTagData(ctx, enterpriseID)
		}(config.EnterpriseID)
	}
	wg.Wait()

	return nil
}

// 处理每个品牌的人群包，生成打标签数据
func (job futureTagData) createBrandTagData(ctx ctx.Context, enterpriseID int64) {
	// 按文件再次划分
	fileInfoList := []plan_task.PlanTaskDownloadedFileModel{}
	if fileErr := orm.Service.PlanTask.Scopes(model.BranchTable(&plan_task.PlanTaskDownloadedFileModel{}, enterpriseID)).
		Where("is_extract = ?", plan_task.PTDF_IsExtract_SFTPData).
		Find(&fileInfoList).Error; fileErr != nil {
		log.Info(ctx).Msgf("%s::获取计划任务文件信息失败: %v", job.name, fileErr)
		return
	}

	var wg sync.WaitGroup
	for _, info := range fileInfoList {
		wg.Add(1)
		go func(fileInfo plan_task.PlanTaskDownloadedFileModel) {
			defer wg.Done()
			memberChan := make(chan micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel, DB_Channel_Num)
			unionChan := make(chan micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel, DB_Channel_Num)

			// 获取查询的基础数据，查询结束时，关闭memberChan, unionChan
			go job.getBaseData(ctx, enterpriseID, fileInfo.ID, memberChan, unionChan)

			job.makeTagData(ctx, enterpriseID, fileInfo, memberChan, unionChan)

			// sftp数据转换成预处理数据完成
			if err := orm.Service.PlanTask.Model(&fileInfo).
				Scopes(model.BranchTable(&fileInfo, enterpriseID)).
				Update("is_extract", plan_task.PTDF_IsExtract_PreData).Error; err != nil {
				log.Info(ctx).Msgf("%s::更新sftp的预处理数据，更新文件is_extract字段失败: %v", job.name, err)
			}
		}(info)
	}
	wg.Wait()
}

func (job futureTagData) makeTagData(ctx ctx.Context, enterpriseID int64, fileInfo plan_task.PlanTaskDownloadedFileModel, memberChan, unionChan <-chan micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 1
		dataList := []micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel{}
		for data := range unionChan {
			dataList = append(dataList, data)
			if i%Read_Data_Nun == 0 {
				job.makeTagDataBa(ctx, enterpriseID, fileInfo, dataList, PlanTask_TagDataType_UnionID)
				dataList = []micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel{}
			}
			i++
		}

		if 0 != len(dataList) {
			job.makeTagDataBa(ctx, enterpriseID, fileInfo, dataList, PlanTask_TagDataType_UnionID)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 1
		dataList := []micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel{}
		for data := range memberChan {
			dataList = append(dataList, data)
			if i%Read_Data_Nun == 0 {
				job.makeTagDataBa(ctx, enterpriseID, fileInfo, dataList, PlanTask_TagDataType_Member)
				dataList = []micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel{}
			}
			i++
		}

		if 0 != len(dataList) {
			job.makeTagDataBa(ctx, enterpriseID, fileInfo, dataList, PlanTask_TagDataType_Member)
		}
	}()
	wg.Wait()
}

func (job futureTagData) makeTagDataBa(ctx ctx.Context, enterpriseID int64, fileInfo plan_task.PlanTaskDownloadedFileModel, dataList []micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel, dataType uint8) {
	switch fileInfo.BaType {
	case plan_task.PTDF_BaType_Relation:
		job.makeTagData_RelationBA(ctx, enterpriseID, fileInfo, dataList, dataType)
	case plan_task.PTDF_BaType_Manager:
		job.makeTagData_ManagerBA(ctx, enterpriseID, fileInfo, dataList, dataType)
	}
}

// 处理数据，每次读取100条处理(关系BA)
func (job futureTagData) makeTagData_RelationBA(ctx ctx.Context, enterpriseID int64, fileInfo plan_task.PlanTaskDownloadedFileModel, dataList []micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel, dataType uint8) {
	tag := struct {
		ID    int64
		Name  string
		TagID string
	}{}
	if tagErr := orm.Service.PlanTask.Table(fmt.Sprintf("%s AS T", micro_scrm.WorkExternalTagModel{}.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %s G ON G.`id` = T.`group_id`", micro_scrm.WorkExternalTagGroupModel{}.TableName())).
		Select("T.id, T.name, T.tag_id").
		Where(map[string]interface{}{
			"T.deleted_at": 0,
			"T.name":       fileInfo.TagName,
			"G.deleted_at": 0,
			"G.group_name": fileInfo.GroupTagName,
		}).
		Scan(&tag).Error; tagErr != nil {
		log.Info(ctx).Msgf("%s::获取标签信息异常: %v", job.name, tagErr)
	}

	userData := []userTagData{}
	dbPrefix := fmt.Sprintf("%s%d.", model.BranchModel{}.DBPrefix(), enterpriseID)
	fields := []string{
		"WUCE.`id` AS enternal_user_id",
		"WUCE.`unionid` AS union_id",
		"WUCE.`openid` AS work_external_user_id",
		"EU.`mobile`",
		"WEUR.`member_id`",
		"M.`userid` AS user_id",
		"WEUR.`id` AS relation_id",
		fmt.Sprintf("'%s' AS tag_name", tag.Name),
		fmt.Sprintf("'%s' AS ww_tag_id", tag.TagID),
		fmt.Sprintf("%d AS tag_id", tag.ID),
		"1 AS is_union",
	}

	whereMap := map[string]interface{}{
		"WEUR.relation_status": 0,
		"WEUR.deleted_at":      0,
		"WEUR.enterprise_id":   enterpriseID,
	}

	switch dataType {
	case PlanTask_TagDataType_Member:
		codeList := []string{}
		for _, code := range dataList {
			if "" == code.MemberCode {
				continue
			}
			codeList = append(codeList, code.MemberCode)
		}
		if 0 == len(codeList) {
			return
		}
		whereMap["WUCE.union_code"] = codeList

	case PlanTask_TagDataType_UnionID:
		unionIDList := []string{}
		for _, code := range dataList {
			if "" == code.UnionID {
				continue
			}
			unionIDList = append(unionIDList, code.UnionID)
		}
		if 0 == len(unionIDList) {
			return
		}
		whereMap["WUCE.unionid"] = unionIDList

	default:
		log.Info(ctx).Msgf("%s::打标签类型错误, 错误类型为: %d", job.name, dataType)
		return
	}

	if dataErr := orm.Service.PlanTask.Table(work_external_user.WorkExternalUserRelationModel{}.TableName() + " AS WEUR").Select(fields).
		Joins("LEFT JOIN " + dbPrefix + micro_scrm_enterprise.WxUserCorpExternalModel{}.TableName() + " WUCE ON WUCE.`id` = WEUR.`work_external_user_id`").
		Joins("LEFT JOIN " + dbPrefix + micro_scrm_enterprise.EnterpriseUserModel{}.TableName() + " EU ON EU.`id` = WUCE.`euid`").
		Joins("LEFT JOIN " + micro_scrm.MemberModel{}.TableName() + " M ON M.`id` = WEUR.`member_id`").
		Where(whereMap).Find(&userData).Error; dataErr != nil {
		log.Info(ctx).Msgf("%s::查询用户的关系数据时异常: %v", job.name, dataErr)
		return
	}

	conditionDataList := []plan_task.TaskConditionDataLogModel{}
	for _, item := range userData {
		conditionDataList = append(conditionDataList, plan_task.TaskConditionDataLogModel{
			TagName:            item.TagName,
			WwTagId:            item.WwTagID,
			UnionId:            item.UnionID,
			WorkExternalUserId: item.WorkExternalUserID,
			IsUnion:            item.IsUnion,
			Mobile:             item.Mobile,
			UserId:             item.UserID,
			MemberId:           item.MemberID,
			TagId:              item.TagID,
			RelationId:         item.RelationID,
		})
	}

	if 0 != len(conditionDataList) {
		if createErr := orm.Service.PlanTask.Scopes(model.BranchTable(&plan_task.TaskConditionDataLogModel{}, enterpriseID)).Create(&conditionDataList).Error; createErr != nil {
			log.Info(ctx).Msgf("%s::查询标签数据失败: %v", job.name, createErr)
		}
	}
}

// 处理数据，管理BA
func (job futureTagData) makeTagData_ManagerBA(ctx ctx.Context, enterpriseID int64, fileInfo plan_task.PlanTaskDownloadedFileModel, dataList []micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel, dataType uint8) {
	tag := struct {
		ID    int64
		Name  string
		TagID string
	}{}
	if tagErr := orm.Service.PlanTask.Table(fmt.Sprintf("%s AS T", micro_scrm.WorkExternalTagModel{}.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %s G ON G.`id` = T.`group_id`", micro_scrm.WorkExternalTagGroupModel{}.TableName())).
		Select("T.id, T.name, T.tag_id").
		Where(map[string]interface{}{
			"T.deleted_at": 0,
			"T.name":       fileInfo.TagName,
			"G.deleted_at": 0,
			"G.group_name": fileInfo.GroupTagName,
		}).
		Scan(&tag).Error; tagErr != nil {
		log.Info(ctx).Msgf("%s::获取标签信息异常: %v", job.name, tagErr)
	}

	userData := []userTagData{}
	dbPrefix := fmt.Sprintf("%s%d.", model.BranchModel{}.DBPrefix(), enterpriseID)
	fields := []string{
		"WUCE.`id` AS enternal_user_id",
		"WUCE.`unionid` AS union_id",
		"WUCE.`openid` AS work_external_user_id",
		"EU.`mobile`",
		"WEUR.`member_id`",
		"M.`userid` AS user_id",
		"WEUR.`id` AS relation_id",
		fmt.Sprintf("'%s' AS tag_name", tag.Name),
		fmt.Sprintf("'%s' AS ww_tag_id", tag.TagID),
		fmt.Sprintf("%d AS tag_id", tag.ID),
		"IF(ISNULL(WEUR.`id`), 0, 1) AS is_union",
	}

	whereMap := map[string]interface{}{
		"M.enterprise_id": enterpriseID,
	}

	filedMain := "union_id"
	filedChild := "unionid"
	conditionMap := map[string]interface{}{}
	switch dataType {
	case PlanTask_TagDataType_Member:
		filedMain = "member_code"
		filedChild = "union_code"

		codeList := []string{}
		for _, code := range dataList {
			if "" == code.MemberCode {
				continue
			}
			codeList = append(codeList, code.MemberCode)
		}
		if 0 == len(codeList) {
			return
		}
		conditionMap["member_code"] = codeList

	case PlanTask_TagDataType_UnionID:
		unionIDList := []string{}
		for _, code := range dataList {
			if "" == code.UnionID {
				continue
			}
			unionIDList = append(unionIDList, code.UnionID)
		}
		if 0 == len(unionIDList) {
			return
		}
		conditionMap["union_id"] = unionIDList
	default:
		log.Info(ctx).Msgf("%s::打标签类型错误, 错误类型为: %d", job.name, dataType)
		return
	}

	if dataErr := orm.Service.PlanTask.Table("(?) AS CTPTCD", orm.Service.PlanTask.
		Scopes(model.BranchTable(&micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel{}, enterpriseID)).
		Select(filedMain, "admin_ba").Where(conditionMap).
		Where(map[string]interface{}{
			"file_id": fileInfo.ID,
		})).Select(fields).
		Joins("INNER JOIN "+dbPrefix+micro_scrm_enterprise.WxUserCorpExternalModel{}.TableName()+" WUCE ON WUCE."+filedChild+" = CTPTCD."+filedMain).
		Joins("INNER JOIN "+dbPrefix+micro_scrm_enterprise.EnterpriseUserModel{}.TableName()+" EU ON EU.`id` = WUCE.`euid`").
		Joins("INNER JOIN "+dbPrefix+micro_scrm_enterprise.MembersModel{}.TableName()+" BM ON BM.`mars_id` = CTPTCD.`admin_ba` AND CTPTCD.`admin_ba` <> ''").
		Joins("INNER JOIN "+micro_scrm.MemberModel{}.TableName()+" M ON M.`userid` = BM.`userid`").
		Joins("LEFT JOIN "+work_external_user.WorkExternalUserRelationModel{}.TableName()+" WEUR ON WEUR.`member_id` = M.`id` AND WEUR.`work_external_user_id` = WUCE.`id` AND `WEUR`.`deleted_at` = 0 "+
			"AND `WEUR`.`enterprise_id` = ?  AND `WEUR`.`relation_status` = 0", enterpriseID).
		Where(whereMap).
		Find(&userData).Error; dataErr != nil {
		log.Info(ctx).Msgf("%s::查询用户的关系数据时异常: %v", job.name, dataErr)
		return
	}

	conditionDataList := []plan_task.TaskConditionDataLogModel{}
	for _, item := range userData {
		conditionDataList = append(conditionDataList, plan_task.TaskConditionDataLogModel{
			TagName:            item.TagName,
			WwTagId:            item.WwTagID,
			UnionId:            item.UnionID,
			WorkExternalUserId: item.WorkExternalUserID,
			IsUnion:            item.IsUnion,
			Mobile:             item.Mobile,
			UserId:             item.UserID,
			MemberId:           item.MemberID,
			TagId:              item.TagID,
			RelationId:         item.RelationID,
		})
	}

	if 0 != len(conditionDataList) {
		if createErr := orm.Service.PlanTask.Scopes(model.BranchTable(&plan_task.TaskConditionDataLogModel{}, enterpriseID)).Create(&conditionDataList).Error; createErr != nil {
			log.Info(ctx).Msgf("%s::查询标签数据失败: %v", job.name, createErr)
		}
	}
}

// 获取基础数据
// 1、member_code
// 2、union_id
func (job futureTagData) getBaseData(ctx ctx.Context, enterpriseID int64, fileID int64, memberChan, unionChan chan micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel) {
	defer func() {
		close(memberChan)
		close(unionChan)
	}()

	baseData := []micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel{}
	i := 0
	for ; ; i++ {
		if baseErr := orm.Service.PlanTask.
			Scopes(model.BranchTable(&micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel{}, enterpriseID)).
			Where(map[string]interface{}{
				"file_id": fileID,
			}).Offset(DB_EveryTime_Num * i).Limit(DB_EveryTime_Num).Find(&baseData).Error; baseErr != nil {
			log.Info(ctx).Msgf("%s::读取打标签的基础数据异常: %v", job.name, baseErr)
			continue
		}

		if 0 == len(baseData) {
			break
		}

		for _, data := range baseData {
			if "" == data.UnionID {
				if "" == data.MemberCode {
					continue
				}
				memberChan <- data
			} else {
				unionChan <- data
			}
		}

		// 初始化
		baseData = []micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel{}
	}
}
