package thirdCrmSftpJob

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/app"
	thirdCrmSftpDao "micro-scrm/internal/app/dao/third_crm_sftp"
	"micro-scrm/internal/app/models"
	crmSftpFileModels "micro-scrm/internal/app/models/crm_sftp_file"
	planService "micro-scrm/internal/app/service/plan"
	thirdCrmSftpSvc "micro-scrm/internal/app/service/third_crm_sftp"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sftp"
	"micro-scrm/internal/pkg/xjob"
	"strconv"
	"strings"
	"sync"
)

const (
	SFTP_ROOT = `/`
)

type pullSFTPData struct {
	thirdDao *thirdCrmSftpDao.ThirdCrmSftpFile
	sftpSvc  *thirdCrmSftpSvc.ThirdCrmSftp
	name     string
}

func NewPullSFTPData() *pullSFTPData {
	return &pullSFTPData{
		thirdDao: thirdCrmSftpDao.NewThirdCrmSftpFile(),
		sftpSvc:  thirdCrmSftpSvc.NewThirdCrmSftp(),
		name:     fmt.Sprint("拉取SFTP的文件数据"),
	}
}

func (job pullSFTPData) Run(ctx ctx.Context, handler *xjob.JobHandler) error {
	// 获取品牌的SFTP配置列表
	sftpConfigList, err := job.sftpSvc.GetSFTPConfigList(micro_scrm.SftpTypeAcxiom)
	if err != nil {
		log.Info(ctx).Msgf("%s::获取SFTP配置失败: %v", job.name, err)
	}

	var wg sync.WaitGroup
	for _, sftpConfig := range sftpConfigList {
		wg.Add(1)
		// go func(config micro_scrm.BrandConfigSftpDetail) {
		go func(config micro_scrm.BrandConfigSftpDetail) {

			defer wg.Done()

			if client, clientErr := sftp.NewSFTPClient(config.ConfigDetail); clientErr != nil {
				log.Info(ctx).Msgf("%s::企业ID:%d, 连接SFTP失败: %v", job.name, config.EnterpriseID, clientErr)
				return
			} else {
				for _, filePath := range config.PathConfig.PlanTask {
					job.syncSFTPData(ctx, config.EnterpriseID, client, filePath, config.BrandName)
				}
			}
		}(sftpConfig)
	}
	wg.Wait()

	return nil
}

// 同步对应的计划任务的数据
func (job pullSFTPData) syncSFTPData(ctx ctx.Context, enterpriseID int64, sftp *sftp.SFTP, filePath, brandName string) {
	// 读取文件夹下的所有文件
	fileInfoList, err := sftp.Client.ReadDir(SFTP_ROOT + filePath)
	if err != nil {
		log.Info(ctx).Msgf("%s::读取SFTP文件时异常: %v", job.name, err)
		fmt.Println(fileInfoList)
		return
	}
	if len(fileInfoList) == 0 {
		log.Info(ctx).Msgf("%s::该SFTP目录下没有文件: %v"+SFTP_ROOT+filePath, job.name, err)
		return
	}

	// 获取已经拉取的文件集合
	fileNameConfig := make(map[int]string)
	fileNameConfig[crmSftpFileModels.TAG_GROUP] = "tag_group_20"
	fileNameConfig[crmSftpFileModels.TAG] = "tag_20"
	fileNameConfig[crmSftpFileModels.TAG_CUSTOMER_REL] = "customer_tag"
	fileNameConfig[crmSftpFileModels.KAY_SA] = "customer_key_sa"
	var fileType []int
	for index := range fileNameConfig {
		fileType = append(fileType, index)
	}
	existsFileList, err := job.sftpSvc.GetExistsFile(ctx, enterpriseID, fileType)
	if err != nil {
		log.Info(ctx).Msgf("%s::读取SFTP文件时异常: %v", job.name, err)
		return
	}
	mapFile := map[string]struct{}{}
	for _, file := range existsFileList {
		mapFile[file.FileName] = struct{}{}
	}
	var insertFileData []crmSftpFileModels.CrmSftpDownloadedFileModel
	for _, ctFile := range fileInfoList {
		if _, ok := mapFile[ctFile.Name()]; !ok {
			for fType, str := range fileNameConfig {
				if strings.Index(ctFile.Name(), str) == 0 {
					fileUrl, err := job.sftpSvc.DecryptGpg(ctx, SFTP_ROOT+filePath+ctFile.Name(), sftp)
					if err != nil {
						return
					}
					// 新增文件
					newFileModel := crmSftpFileModels.CrmSftpDownloadedFileModel{
						FileName:       ctFile.Name(),
						FileFullPath:   SFTP_ROOT + filePath + ctFile.Name(),
						FileType:       int32(fType),
						FileURL:        fileUrl,
						MicroBaseModel: models.CreateMicroBase(enterpriseID, app.UID()),
					}
					insertFileData = append(insertFileData, newFileModel)
					break
				}
			}
		}
	}
	if len(insertFileData) == 0 {
		log.Info(ctx).Msgf("%s::该SFTP目录下没有新增文件"+SFTP_ROOT+filePath, job.name)
		return
	}

	if err := job.thirdDao.SaveFile(ctx, insertFileData); err != nil {
		log.Info(ctx).Msgf("%s::写入文件记录失败: %v", job.name, err)
		return
	}
	planService.NewTaskNotice().SftpSendNotice(ctx, planService.SFTPMessageData{
		EnterpriseName: brandName,
		SFTPName:       "SFTP同步CDP标签文件",
		UploadRes:      helper.RuleThree(err != nil, "异常", "成功"),
	})
	fmt.Println("文件拉取成功" + SFTP_ROOT + filePath + ",文件数量:" + strconv.Itoa(len(insertFileData)))

}
