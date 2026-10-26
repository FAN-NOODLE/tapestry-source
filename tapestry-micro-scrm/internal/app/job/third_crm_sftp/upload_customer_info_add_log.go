package thirdCrmSftpJob

import (
	"encoding/csv"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"io/ioutil"
	customerSvc "micro-scrm/internal/app/service/customer"
	planService "micro-scrm/internal/app/service/plan"
	thirdCrmSftpSvc "micro-scrm/internal/app/service/third_crm_sftp"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sftp"
	"micro-scrm/internal/pkg/xjob"
	"os"
	"path"
	"strconv"
	"time"
)

type uploadCustomerInfoAddLog struct {
	name            string
	customerInfoSvc *customerSvc.CustomerInfoAdd
	sftpSvc         *thirdCrmSftpSvc.ThirdCrmSftp
}

func NewUploadCustomerInfoAddLog() *uploadCustomerInfoAddLog {
	return &uploadCustomerInfoAddLog{
		name:            fmt.Sprint("上传外部联系人标签操作记录到sftp"),
		customerInfoSvc: customerSvc.NewCustomerInfoAdd(),
		sftpSvc:         thirdCrmSftpSvc.NewThirdCrmSftp(),
	}
}

func (job uploadCustomerInfoAddLog) Run(ctx ctx.Context, handler *xjob.JobHandler) error {
	sftpConfigList, err := job.sftpSvc.GetSFTPConfigList(micro_scrm.SftpTypeAcxiom)
	if err != nil {
		log.Info(ctx).Msgf("%s::获取SFTP配置失败: %v", job.name, err)
	}
	// var wg sync.WaitGroup
	for _, config := range sftpConfigList {
		if config.PathConfig.CustomerLog == "" {
			continue
		}
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
func (job *uploadCustomerInfoAddLog) pushLogData(ctx ctx.Context, enterpriseID int64, sftp *sftp.SFTP, config micro_scrm.BrandConfigSftpDetail) {
	list, err := job.customerInfoSvc.GetListByTime(ctx, enterpriseID)
	if err != nil || len(list) == 0 {
		return
	}
	localFilePath := `./storage/crm_csv/customer_Info_` + time.Now().Format("20060102") + "_" + strconv.FormatInt(helper.GetNowDateTimestamp(), 10) + ".csv"
	f, err := os.Create(localFilePath) // 创建文件
	if err != nil {
		panic(err)
	}

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	data := [][]string{
		{"sa_code", "customer_unionid", "shoe_size", "birthday", "address", "unix_time"},
	}

	w := csv.NewWriter(f) // 创建一个新的写入文件流
	defer f.Close()
	for _, v := range list {
		tmp := []string{
			v.OperateUserId, v.ExternalUnionId, v.ShoeSize, fmt.Sprintf("%d-%d", v.Month, v.Day), v.Address, fmt.Sprintf("%d", v.UpdatedAt),
		}
		data = append(data, tmp)
	}

	w.WriteAll(data) // 写入数据
	w.Flush()

	fmt.Println("csv日志文件写入成功" + localFilePath + ",日志数量:" + strconv.Itoa(len(list)))
	err = job.uploadFile(sftp, localFilePath, config.PathConfig.CustomerLog)
	if err != nil {
		log.Info(ctx).Msgf("%s::sftp文件写入失败: %v", job.name, err)
		return
	}
	defer planService.NewTaskNotice().SftpSendNotice(ctx, planService.SFTPMessageData{
		EnterpriseName: config.BrandName,
		SFTPName:       "SFTP同步顾客信息上传CDP",
		UploadRes:      helper.RuleThree(err != nil, "异常", "成功"),
	})
	log.Info(ctx).Msgf("%s::写入sftp文件成功: %s", job.name, localFilePath)

}

func (job *uploadCustomerInfoAddLog) uploadFile(sftp *sftp.SFTP, localFilePath string, remotePath string) error {
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
