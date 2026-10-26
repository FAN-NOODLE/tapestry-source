package v2

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"gorm.io/gorm"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise/plan_task"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm/work_external_user"
	planTaskSftp "micro-scrm/internal/orm/service/plan_task"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/workwechat"
	"micro-scrm/internal/pkg/xjob"
)

// 计划任务 - 外部联系人打标签

var _ runner = externalUserTagging{}

type externalUserTagging struct {
	baseJob
	name           string
	queryLimit     int
	chanLen        int
	controlChanLen int
}

type tagMsg struct {
	ID    int64
	Name  string
	TagID string
}

func NewExternalUserTagging() *externalUserTagging {
	return &externalUserTagging{
		name:           "用户打企微标签",
		queryLimit:     2000,
		chanLen:        1000,
		controlChanLen: 30,
	}
}

func (job externalUserTagging) Run(ctx ctx.Context, handler *xjob.JobHandler) error {
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

			// 查询查询历史人员标签，历史人员标签不需要重新打标
			job.FindHistoryUserTag(ctx, enterpriseID)

			dataListChan := make(chan plan_task.TaskConditionDataLogModel, job.chanLen)
			pollingFlag := make(chan struct{})
			go job.getBaseData(ctx, enterpriseID, pollingFlag, dataListChan)
			job.makeTag(ctx, enterpriseID, pollingFlag, dataListChan)
		}(config.EnterpriseID)
	}
	wg.Wait()

	return nil
}

// 打标签
func (job externalUserTagging) makeTag(ctx ctx.Context, enterpriseID int64, pollingFlag chan<- struct{}, dataListChan <-chan plan_task.TaskConditionDataLogModel) {
	controlChan := make(chan struct{}, job.controlChanLen-1) // job.controlChanLen 量的并发

	dataList := []plan_task.TaskConditionDataLogModel{}
	var wg sync.WaitGroup
	// 获取已经拉完数据的标签
	wwObj, wwObjErr := workwechat.NewWorkWeChat(ctx, enterpriseID)
	if wwObjErr != nil {
		log.Info(ctx).Msgf("%s::获取品牌企微OBJ异常: %v", job.name, wwObjErr)
		return
	}

	// 成功和失败
	successDataList := []plan_task.TaskConditionDataLogModel{}
	failureDataList := []plan_task.TaskConditionDataLogModel{}

	for data := range dataListChan {
		select {
		case controlChan <- struct{}{}:
			dataList = append(dataList, data)
		default:
			dataList = append(dataList, data)
			// 打标签
			for _, tagData := range dataList {
				wg.Add(1)
				select {
				case <-time.Tick(10 * time.Millisecond):
					go func(tagInfo plan_task.TaskConditionDataLogModel) {
						defer wg.Done()

						params := work.UserTagHandlerRequest{
							Userid:         tagInfo.UserId,
							ExternalUserid: tagInfo.WorkExternalUserId,
							AddTag:         []string{tagInfo.WwTagId},
						}
						if _, err := wwObj.ManagerCustomerTag(params); err == nil {
							successDataList = append(successDataList, tagInfo)
						} else {
							log.Info(ctx).Msgf("%s::打企微标签%v失败: %v", job.name, params, err)
							failureDataList = append(failureDataList, tagInfo)
						}
					}(tagData)
				}
			}
			wg.Wait()

			// 更新wesocial的标签关系数据
			if 0 != len(successDataList) {
				idList := []int64{}
				insertData := []work_external_user.WorkExternalUserTagModel{}
				sourceList := map[string]struct{}{}
				for _, logInfo := range successDataList {
					sourceKey := fmt.Sprintf("%d%d", logInfo.TagId, logInfo.RelationId)
					if _, ok := sourceList[sourceKey]; ok == true {
						continue
					} else {
						sourceList[sourceKey] = struct{}{}
					}
					// 判断是否已经添加标签
					var count int64
					if countErr := orm.Service.PlanTask.Model(&work_external_user.WorkExternalUserTagModel{}).Where(map[string]interface{}{
						"external_tag_id": logInfo.TagId,
						"relation_id":     logInfo.RelationId,
						"deleted_at":      0,
					}).Count(&count).Error; countErr != nil {
						log.Info(ctx).Msgf("%s::查询异常: %v", job.name, countErr)
						continue
					}
					// 判断是否存在
					idList = append(idList, logInfo.ID)

					insertData = append(insertData, work_external_user.WorkExternalUserTagModel{
						MicroScrmBaseModel: model.MicroScrmBaseModel{
							BaseModel: model.BaseModel{
								ID: <-NewUUID().ID,
							},
							EnterpriseID: enterpriseID,
						},
						RelationID:    logInfo.RelationId,
						ExternalTagID: logInfo.TagId,
					})
				}

				orm.Service.PlanTask.Transaction(func(tx *gorm.DB) error {
					// 1、更新task_condition_data_log.is_success = 2
					if updateErr := orm.Service.PlanTask.Scopes(model.BranchTable(&plan_task.TaskConditionDataLogModel{}, enterpriseID)).
						Model(&plan_task.TaskConditionDataLogModel{}).Where("id IN (?)", idList).
						Update("is_success", plan_task.TCDL_IsSuccess_Yes).Error; updateErr != nil {
						return updateErr
					}
					// 2、建立标签关系
					if 0 != len(insertData) {
						if createErr := orm.Service.PlanTask.Create(&insertData).Error; createErr != nil {
							return createErr
						}
					}

					return nil
				})
			}

			// 失败的数据
			if 0 != len(failureDataList) {
				// 1、task_condition_data_log.tag_num + 1
				idList := []int64{}
				for _, log := range failureDataList {
					idList = append(idList, log.ID)
				}

				orm.Service.PlanTask.Scopes(model.BranchTable(&plan_task.TaskConditionDataLogModel{}, enterpriseID)).
					Model(&plan_task.TaskConditionDataLogModel{}).Where("id IN (?)", idList).
					Update("tag_num", gorm.Expr("tag_num + 1"))
			}

			// 初始化数据
			{
				successDataList = []plan_task.TaskConditionDataLogModel{}
				failureDataList = []plan_task.TaskConditionDataLogModel{}

				dataList = []plan_task.TaskConditionDataLogModel{}
				controlChan = make(chan struct{}, job.controlChanLen-1)
			}

		}
	}
}

// 获取品牌需要打标签打基础数据
func (job externalUserTagging) getBaseData(ctx ctx.Context, enterpriseID int64, pollingFlag <-chan struct{}, dataListChan chan<- plan_task.TaskConditionDataLogModel) {
	defer close(dataListChan)

	// 获取需要打标的标签名称
	tagList := []plan_task.PlanTaskDownloadedFileModel{}
	if tagErr := orm.Service.PlanTask.Scopes(model.BranchTable(&plan_task.PlanTaskDownloadedFileModel{}, enterpriseID)).Where("is_extract = ?", plan_task.PTDF_IsExtract_PreData).Find(&tagList).Error; tagErr != nil {
		log.Info(ctx).Msgf("%s::获取标签名称失败", job.name)
		return
	}

	for _, info := range tagList {
		// 获取标签的企微ID
		tag := tagMsg{}
		if tagErr := orm.Service.PlanTask.Table(fmt.Sprintf("%s AS T", micro_scrm.WorkExternalTagModel{}.TableName())).
			Joins(fmt.Sprintf("LEFT JOIN %s G ON G.`id` = T.`group_id`", micro_scrm.WorkExternalTagGroupModel{}.TableName())).
			Select("T.id, T.name, T.tag_id").
			Where(map[string]interface{}{
				"T.deleted_at": 0,
				"T.name":       info.TagName,
				"G.deleted_at": 0,
				"G.group_name": info.GroupTagName,
			}).
			Scan(&tag).Error; tagErr != nil || 0 == tag.ID {
			continue
		}

		dataList := []plan_task.TaskConditionDataLogModel{}

		// 监听数据的储量
		isFinish := make(chan struct{})
		go job.listenLogDataNum(enterpriseID, tag, isFinish)

	Label:
		for i := 1; ; i++ {
			if dataErr := orm.Service.PlanTask.Scopes(model.BranchTable(&plan_task.TaskConditionDataLogModel{}, enterpriseID)).Where(map[string]interface{}{
				"deleted_at": 0,
				"is_success": plan_task.TCDL_IsSuccess_No,
				"tag_id":     tag.ID,
			}).Where("tag_num < 5").Offset((i - 1) * job.queryLimit).Limit(job.queryLimit).Find(&dataList).Error; dataErr != nil {
				log.Info(ctx).Msgf("%s::查询基础数据失败", job.name)
				continue
			}

			select {
			case <-isFinish:
				// 更新plan_task_downloaded_file.is_extract = 3
				if updateFileErr := orm.Service.PlanTask.Scopes(model.BranchTable(&plan_task.PlanTaskDownloadedFileModel{}, enterpriseID)).
					Model(&info).Update("is_extract", plan_task.PTDF_IsExtract_Tag); updateFileErr != nil {
					log.Info(ctx).Msgf("%s::更新plan_task_downloaded_file.is_extract = 3异常: %v", job.name, updateFileErr)
				}
				break Label
			default:
				if 0 == len(dataList) {
					i = 0
				}
			}

			for _, data := range dataList {
				dataListChan <- data
			}

			// 初始化dataList
			dataList = []plan_task.TaskConditionDataLogModel{}
		}

		// 处理结束数据
		{
			close(isFinish)
		}
	}
}

func (job externalUserTagging) listenLogDataNum(enterpriseID int64, tag tagMsg, isFinish chan<- struct{}) {
	var count int64
Label:
	for {
		select {
		case <-time.Tick(10 * time.Second):
			if countErr := orm.Service.PlanTask.Scopes(model.BranchTable(&plan_task.TaskConditionDataLogModel{}, enterpriseID)).Where(map[string]interface{}{
				"deleted_at": 0,
				"is_success": plan_task.TCDL_IsSuccess_No,
				"tag_id":     tag.ID,
			}).Where("tag_num < 5").Count(&count).Error; countErr == nil {
				if 0 == count {
					isFinish <- struct{}{}
					break Label
				}
			}
		}
	}
}

type UUID struct {
	ID <-chan int64
}

var uuid *UUID
var uuidOnce sync.Once
var onceMutex sync.Once

func NewUUID() *UUID {
	uuidOnce.Do(func() {
		uuid = &UUID{
			ID: timeToUUID(),
		}
	})

	return uuid
}

func timeToUUID() <-chan int64 {
	var ch = make(chan int64)
	var id int64
	onceMutex.Do(
		func() {
			timeString := "799" + fmt.Sprintf("%03v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100)) + fmt.Sprintf("%s", strconv.FormatInt(time.Now().Unix(), 10))
			id, _ = strconv.ParseInt(timeString, 10, 64)
		},
	)
	go func() {
		for {
			id += 1
			ch <- id
		}
	}()

	return ch
}
func (job externalUserTagging) FindHistoryUserTag(ctx ctx.Context, enterpriseID int64) {
	// 获取需要打标的标签名称
	tagList := []plan_task.PlanTaskDownloadedFileModel{}
	if tagErr := orm.Service.PlanTask.Scopes(model.BranchTable(&plan_task.PlanTaskDownloadedFileModel{}, enterpriseID)).Where("is_extract = ?", plan_task.PTDF_IsExtract_PreData).Find(&tagList).Error; tagErr != nil {
		log.Info(ctx).Msgf("%s::获取标签名称失败", job.name)
		return
	}

	for _, info := range tagList {
		// 获取标签的企微ID
		tag := tagMsg{}
		if tagErr := orm.Service.PlanTask.Table(fmt.Sprintf("%s AS T", micro_scrm.WorkExternalTagModel{}.TableName())).
			Joins(fmt.Sprintf("LEFT JOIN %s G ON G.`id` = T.`group_id`", micro_scrm.WorkExternalTagGroupModel{}.TableName())).
			Select("T.id, T.name, T.tag_id").
			Where(map[string]interface{}{
				"T.deleted_at": 0,
				"T.name":       info.TagName,
				"G.deleted_at": 0,
				"G.group_name": info.GroupTagName,
			}).
			Scan(&tag).Error; tagErr != nil || 0 == tag.ID {
			return
		}
		for page := 0; ; page++ {
			relationIds := make([]int64, 0)
			if err := orm.Service.PlanTask.Model(&work_external_user.WorkExternalUserTagModel{}).Where(map[string]interface{}{
				"external_tag_id": tag.TagID,
				"enterprise_id":   enterpriseID,
				"deleted_at":      0,
			}).Select("relation_id").Order("relation_id").Offset(job.queryLimit * page).
				Limit(job.queryLimit).Find(&relationIds).Error; err != nil {
				page--
				continue
			} else {
				if len(relationIds) > 0 {
					orm.Service.PlanTask.Scopes(model.BranchTable(&plan_task.TaskConditionDataLogModel{}, enterpriseID)).
						Model(&plan_task.TaskConditionDataLogModel{}).
						Where("tag_id = ? AND is_success=? AND is_has_tag=? AND deleted_at=? AND  relation_id IN (?)",
							tag.TagID, plan_task.TCDL_IsSuccess_No, plan_task.TCDL_IsHasTag_No, 0, relationIds).
						Updates(map[string]interface{}{
							"is_success": plan_task.TCDL_IsSuccess_Yes,
							"is_has_tag": plan_task.TCDL_IsHasTag_Yes,
						})
					if len(relationIds) < job.queryLimit {
						break
					}
				} else {
					break
				}
			}
		}

	}
}
