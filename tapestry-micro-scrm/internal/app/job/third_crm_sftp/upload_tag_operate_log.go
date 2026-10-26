package thirdCrmSftpJob

import (
	"encoding/csv"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"io/ioutil"
	planService "micro-scrm/internal/app/service/plan"
	tagsSvc "micro-scrm/internal/app/service/tags"
	thirdCrmSftpSvc "micro-scrm/internal/app/service/third_crm_sftp"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sftp"
	"micro-scrm/internal/pkg/xjob"
	"os"
	"path"
	"strconv"
)

type uploadTagOperateLog struct {
	name    string
	tagsSvc *tagsSvc.ImpressionTagLog
	sftpSvc *thirdCrmSftpSvc.ThirdCrmSftp
}

func NewUploadTagOperateLog() *uploadTagOperateLog {
	return &uploadTagOperateLog{
		name:    fmt.Sprint("上传外部联系人标签操作记录到sftp"),
		tagsSvc: tagsSvc.NewImpressionTagLog(),
		sftpSvc: thirdCrmSftpSvc.NewThirdCrmSftp(),
	}
}

func (job uploadTagOperateLog) Run(ctx ctx.Context, handler *xjob.JobHandler) error {
	sftpConfigList, err := job.sftpSvc.GetSFTPConfigList(micro_scrm.SftpTypeAcxiom)
	if err != nil {
		log.Info(ctx).Msgf("%s::获取SFTP配置失败: %v", job.name, err)
	}
	// var wg sync.WaitGroup
	for _, config := range sftpConfigList {
		// wg.Add(1)
		// go func(config micro_scrm.BrandConfigSftpDetail) {
		//
		//	defer wg.Done()

		if client, clientErr := sftp.NewSFTPClient(config.ConfigDetail); clientErr != nil {
			log.Info(ctx).Msgf("%s::企业ID:%d, 连接SFTP失败: %v", job.name, config.EnterpriseID, clientErr)
			// return
		} else {
			job.pushLogData(ctx, config.EnterpriseID, client, config)
		}
		// }(sftpConfig)
	}
	// wg.Wait()

	return nil

}

// 同步对应的计划任务的数据
func (job *uploadTagOperateLog) pushLogData(ctx ctx.Context, enterpriseID int64, sftp *sftp.SFTP, config micro_scrm.BrandConfigSftpDetail) {
	list, err := job.tagsSvc.GetTagLog(ctx, enterpriseID)
	if err != nil || len(list) == 0 {
		return
	}
	localFilePath := `./storage/crm_csv/impression_tags_` + helper.GetNowDate() + "_" + strconv.FormatInt(helper.GetNowDateTimestamp(), 10) + ".csv"
	f, err := os.Create(localFilePath) // 创建文件
	if err != nil {
		panic(err)
	}

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	data := [][]string{
		{"sa_code", "customer_unionid", "action_type", "tag_id", "unix_time"},
	}

	w := csv.NewWriter(f) // 创建一个新的写入文件流
	defer f.Close()
	var lastDataCreatedAt int64 // 定义最后一条数据的时间，推送成功后软删除日志
	for _, v := range list {
		tmp := []string{
			v.OperateUserId, v.ExternalUnionId, fmt.Sprintf("%d", v.ActionType), v.ThirdTagID, fmt.Sprintf("%d", v.CreatedAt),
		}
		data = append(data, tmp)
		lastDataCreatedAt = v.CreatedAt
	}

	w.WriteAll(data) // 写入数据
	w.Flush()

	fmt.Println("csv日志文件写入成功" + localFilePath + ",日志数量:" + strconv.Itoa(len(list)))

	err = job.uploadFile(sftp, localFilePath, config.PathConfig.TagLog)
	if err != nil {
		log.Info(ctx).Msgf("%s::sftp文件写入失败: %v", job.name, err)
		return
	}
	log.Info(ctx).Msgf("%s::写入sftp文件成功: %s", job.name, localFilePath)

	if delErr := job.tagsSvc.DeleteList(ctx, enterpriseID, lastDataCreatedAt); delErr != nil {
		log.Info(ctx).Msgf("%s::已经推送的日志软删除失败: %v", job.name, delErr)
		return
	}
	planService.NewTaskNotice().SftpSendNotice(ctx, planService.SFTPMessageData{
		EnterpriseName: config.BrandName,
		SFTPName:       "SFTP同步印象标签至CDP",
		UploadRes:      helper.RuleThree(err != nil, "异常", "成功"),
	})

}

func (job *uploadTagOperateLog) uploadFile(sftp *sftp.SFTP, localFilePath string, remotePath string) error {
	srcFile, err := os.Open(localFilePath)
	defer srcFile.Close()
	if err != nil {
		return err
	}
	var remoteFileName = path.Base(localFilePath)
	dstFile, err := sftp.Client.Create(path.Join(remotePath, remoteFileName))
	defer dstFile.Close()
	if err != nil {
		return err
	}
	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		return err
	}
	dstFile.Write(ff)
	return nil
}
