package v2

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"path"
	"strings"
	"sync"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/config"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise/plan_task"
	planTaskSftp "micro-scrm/internal/orm/service/plan_task"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/sftp"
	"micro-scrm/internal/pkg/uid"
	"micro-scrm/internal/pkg/workwechat"
	"micro-scrm/internal/pkg/xjob"
)

// 断言run
var _ runner = pullSFTPData{}

const (
	READ_NUM  = 2000
	WRITE_NUM = 500
	SFTP_ROOT = `/`
	FILED_NUM = 4
)

type pullSFTPData struct {
	baseJob
	aesCbc256

	name string
}

func NewPullSFTPData(config *config.Config) *pullSFTPData {
	return &pullSFTPData{
		name: fmt.Sprintf("%s::%s", planName, "拉取SFTP的文件数据"),
		aesCbc256: aesCbc256{
			key: config.CampaignToolService.AES_CBC_256.Key,
			iv:  config.CampaignToolService.AES_CBC_256.Iv,
		},
	}
}

func (job pullSFTPData) Run(ctx ctx.Context, handler *xjob.JobHandler) error {
	// 获取品牌的SFTP配置列表
	sftpConfig, err := job.getSFTPConfigList(planTaskSftp.SFTP_TYPE_PLAN_TASK)
	if err != nil {
		log.Info(ctx).Msgf("%s::获取SFTP配置失败: %v", job.name, err)
	}

	var wg sync.WaitGroup
	for _, config := range sftpConfig {
		wg.Add(1)
		go func(config micro_scrm.BrandConfigSftpDetail) {
			defer wg.Done()

			if client, clientErr := sftp.NewSFTPClient(config.ConfigDetail); clientErr != nil {
				log.Info(ctx).Msgf("%s::企业ID:%d, 连接SFTP失败: %v", job.name, config.EnterpriseID, clientErr)
				return
			} else {
				for _, filePath := range config.PathConfig.PlanTask {
					job.syncSFTPData(ctx, config.EnterpriseID, client, filePath)
				}
			}
		}(config)
	}
	wg.Wait()

	return nil
}

// 同步对应的计划任务的数据
func (job pullSFTPData) syncSFTPData(ctx ctx.Context, enterpriseID int64, sftp *sftp.SFTP, filePath string) {
	// 读取文件夹下的所有文件
	fileInfoList, fileInfoErr := sftp.Client.ReadDir(SFTP_ROOT + filePath)
	if fileInfoErr != nil {
		log.Info(ctx).Msgf("%s::读取SFTP文件时异常: %v", job.name, fileInfoErr)
		return
	}

	// 获取已经拉取的文件集合
	if downloadFileList, fileErr := orm.Service.PlanTask.GetCTDownloadedFileList(enterpriseID, plan_task.PTDF_FileType_CT); fileErr != nil {
		log.Info(ctx).Msgf("%s::读取SFTP文件时异常: %v", job.name, fileErr)
		return
	} else {
		mapFile := map[string]struct{}{}
		for _, file := range downloadFileList {
			mapFile[file] = struct{}{}
		}

		for _, ctFile := range fileInfoList {
			if _, ok := mapFile[ctFile.Name()]; !ok {
				tags, tagErr := job.splitFileName(ctFile.Name())
				if tagErr != nil {
					log.Info(ctx).Msgf("%s::读取SFTP文件,获取标签时异常: %v", job.name, tagErr)
					continue
				}

				// 新增文件
				newFileModel := plan_task.PlanTaskDownloadedFileModel{
					FileName:     ctFile.Name(),
					FileType:     plan_task.PTDF_FileType_CT,
					TagName:      tags.tagName,
					GroupTagName: tags.groupName,
					BaType:       tags.baType,
				}
				if insertErr := orm.Service.PlanTask.Scopes(model.BranchTable(&plan_task.PlanTaskDownloadedFileModel{}, enterpriseID)).
					Create(&newFileModel).Error; insertErr == nil {
					// 更新标签
					if createTagErr := job.createTag(ctx, enterpriseID, tags); createTagErr == nil {
						job.sftpFileToDb(ctx, enterpriseID, sftp, fmt.Sprintf("%s%s/%s", SFTP_ROOT, filePath, ctFile.Name()), tags.tagName, newFileModel.ID)
					}
				}

				if err := orm.Service.PlanTask.Scopes(model.BranchTable(&plan_task.PlanTaskDownloadedFileModel{}, enterpriseID)).
					Model(&newFileModel).Update("is_extract", plan_task.PTDF_IsExtract_SFTPData).
					Error; err != nil {
					log.Info(ctx).Msgf("%s::更新SFTP文件提取数据状态异常: %v", job.name, err)
					continue
				}
			}
		}
	}
}

// 读取SFTP人群包数据
func (job pullSFTPData) sftpFileToDb(ctx ctx.Context, enterpriseID int64, sftp *sftp.SFTP, srcFile, tagName string, fileID int64) {
	file, _ := sftp.Client.Open(srcFile)
	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1

	dataChan := make(chan string, READ_NUM)
	go func(csvReader *csv.Reader) {
		isTitle := true
		for {
			row, rowErr := csvReader.Read()
			if isTitle {
				isTitle = false
				continue
			}
			if rowErr != nil && rowErr != io.EOF {
				log.Info(ctx).Msgf("读取SFTP文件%s，逐行读取，错误为: %v", srcFile, rowErr)
			}

			if rowErr == io.EOF {
				close(dataChan)
				break
			}

			data := strings.Join(row, ",")
			dataChan <- data
		}
	}(csvReader)

	insertData := []string{}
	i := 1
	for row := range dataChan {
		if i%WRITE_NUM == 0 {
			fmt.Println("开始")
			job.writeToDB(ctx, enterpriseID, insertData, srcFile, tagName, fileID)
			insertData = []string{}
		}

		insertData = append(insertData, row)
		i++
	}

	if len(insertData) != 0 {
		job.writeToDB(ctx, enterpriseID, insertData, srcFile, tagName, fileID)
	}

}

func (job pullSFTPData) writeToDB(ctx ctx.Context, enterpriseID int64, insertData []string, srcFile, tagName string, fileID int64) {
	if 0 == len(insertData) {
		return
	}

	// 获取Campaign-Tool推送数据时间
	if tagName == "" {
		return
	}
	filename := strings.Replace(path.Base(srcFile), path.Ext(srcFile), "", -1)
	fileArr := strings.Split(filename, "_")
	cronTime := fileArr[len(fileArr)-1]

	newData := []micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel{}
	for _, data := range insertData {
		arr := strings.Split(data, ",")

		// SFTP文件的字段为: c_code, union_id, grade_code, admin_ba
		if FILED_NUM > len(arr) {
			continue
		}

		// 判断c_code, union_id是否为空
		isEmptyCode := false
		isEmptyUnionID := false

		if "" == arr[0] {
			isEmptyCode = true
		}

		unionID := ""
		if "" == arr[1] {
			isEmptyUnionID = true
		} else {
			if unionId, unionErr := base64.StdEncoding.DecodeString(arr[1]); unionErr == nil {
				if unionID = string(common.AES_CBC_Decrypt(unionId, []byte(job.aesCbc256.key), []byte(job.aesCbc256.iv))); "" == unionID {
					isEmptyUnionID = true
				}
			}
		}

		if isEmptyCode && isEmptyUnionID {
			continue
		}

		newData = append(newData, micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel{
			MemberCode:   arr[0], // 会员卡号
			UnionID:      unionID,
			GradeCode:    arr[2],
			AdminBA:      arr[3],
			CronFileName: srcFile,
			CreateTime:   cronTime,
			TagName:      tagName,
			FileID:       fileID,
		})

		if len(insertData) == 0 {
			break
		}
	}
	// Campaign-Tool 落表
	if len(newData) > 0 {
		if err := orm.Service.PlanTask.InsertCampaignToolPlanTaskData(enterpriseID, newData); err != nil {
			log.Info(ctx).Msgf("CampaignTool数据写入异常: %v", err)
		}
	}
}

// 创建企微标签
func (job pullSFTPData) createTag(ctx ctx.Context, enterpriseID int64, info tagInfo) error {
	tagGroup, tagGroupIsExist, tagGroupErr := orm.Service.Tag.CheckTagGroupName(enterpriseID, info.groupName)
	if tagGroupErr != nil {
		log.Info(ctx).Msgf("获取标签组信息错误: %v", tagGroupErr)
		return tagGroupErr
	}

	// 插入新的标签
	obj, err := workwechat.NewWorkWeChat(ctx, enterpriseID)
	if err != nil {
		log.Info(ctx).Msgf("创建企微对象失败: %v", err)
	}

	if true != tagGroupIsExist {
		tagResult, tagErr := obj.AddTag(work.AddCorpTagRequest{
			GroupName: info.groupName,
			TagList: []work.AddCorpTag{
				{
					Name: info.tagName,
				},
			},
		})

		if tagErr != nil {
			log.Info(ctx).Msgf("创建企微对象失败: %v", tagErr)
			return tagErr
		}

		// 跟新标签组，标签
		tagGroup, _, tagGroupErr = orm.Service.Tag.CreateTagGroup2(micro_scrm.WorkExternalTagGroupModel{
			MicroScrmBaseModel: model.MicroScrmBaseModel{
				BaseModel: model.BaseModel{
					ID: uid.ID(ctx),
				},
				EnterpriseID: enterpriseID,
			},
			GroupName:  info.groupName,
			GroupID:    tagResult.TagGroup.GroupID,
			Order:      tagResult.TagGroup.Order,
			CreateTime: int(tagResult.TagGroup.CreateTime),
		})

		if tagGroupErr != nil {
			log.Info(ctx).Msgf("写入标签组信息错误: %v", tagGroupErr)
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
			log.Info(ctx).Msgf("写入标签信息错误: %v", tagGroupErr)
		}
	}

	_, tagIsExist, tagErr := orm.Service.Tag.CheckTagName(enterpriseID, tagGroup.ID, info.tagName)
	if tagErr != nil {
		log.Info(ctx).Msgf("获取标签信息错误: %v", tagGroupErr)
		return tagGroupErr
	}
	if !tagIsExist {
		tagResult, tagResultErr := obj.AddTag(work.AddCorpTagRequest{
			GroupName: info.groupName,
			TagList: []work.AddCorpTag{
				{
					Name: info.tagName,
				},
			},
		})
		if tagResultErr != nil {
			log.Info(ctx).Msgf("创建企微对象失败: %v", tagResultErr)
			return tagResultErr
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
			log.Info(ctx).Msgf("写入标签信息错误: %v", tagGroupErr)
		}
	}

	return nil
}
