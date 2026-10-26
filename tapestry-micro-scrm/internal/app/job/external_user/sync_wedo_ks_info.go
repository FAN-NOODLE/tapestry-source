package externalUserJob

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"io"
	"io/ioutil"
	"micro-scrm/internal/app"
	thirdCrmSftpDao "micro-scrm/internal/app/dao/third_crm_sftp"
	"micro-scrm/internal/app/http/request/external"
	"micro-scrm/internal/app/models"
	crmSftpFileModels "micro-scrm/internal/app/models/crm_sftp_file"
	externalUserService "micro-scrm/internal/app/service/external_user"
	planService "micro-scrm/internal/app/service/plan"
	thirdCrmSftpSvc "micro-scrm/internal/app/service/third_crm_sftp"
	"micro-scrm/internal/modules/marketingPromotion/storage"
	"micro-scrm/internal/modules/marketingPromotion/storage/cos"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sftp"
	"micro-scrm/internal/pkg/xjob"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	SFTPROOTKS = `/`
)

type SyncWedoKSInfo struct {
	name       string
	filePrefix string
}

func NewSyncWedoKSInfo() *SyncWedoKSInfo {
	return &SyncWedoKSInfo{
		name:       "同步外部联系人在wedo上的信息",
		filePrefix: "weclient_member_result_",
	}
}

func (job *SyncWedoKSInfo) Run(ctx ctx.Context, handler *xjob.JobHandler) error {

	var cond = map[string]interface{}{"brand_name": "ks", "sftp_type": 5}
	config, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetSFTPConfig(ctx, cond)
	if err != nil {
		log.Info(ctx).Msgf("%s::获取SFTP配置失败: %v", job.name, err)
	}

	if config.BrandName != "ks" {
		return nil
	}
	fileType := crmSftpFileModels.Member_Info_KS

	if client, clientErr := sftp.NewSFTPClient(config.ConfigDetail); clientErr != nil {
		log.Info(ctx).Msgf("%s::企业ID:%d, 连接SFTP失败: %v", job.name, config.EnterpriseID, clientErr)
		// return
	} else {
		if len(config.PathConfig.CdpExternalUser) == 0 {
			return nil
		}
		filePath := config.PathConfig.CdpExternalUser[0]
		// sftp读取文件
		fileLists, err := client.Client.ReadDir(SFTPROOTKS + filePath)
		if err != nil {
			log.Info(ctx).Msgf("%s::企业ID:%d, 获取目录下文件失败: %v", job.name, config.EnterpriseID, clientErr)
			return err
		}
		// 排序，顺序执行文件
		sort.Slice(fileLists, func(i, j int) bool {
			return fileLists[i].ModTime().Unix() < fileLists[i].ModTime().Unix()
		})
		// 查询已经在数据库的文件
		fileInfoList, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetExistsFile(ctx, config.EnterpriseID, []int{fileType})
		if err != nil {
			log.Info(ctx).Msgf("%s::企业ID:%d, 获取目录下文件失败: %v", job.name, config.EnterpriseID, clientErr)
			return err
		}
		fileInfoMap := make(map[string]crmSftpFileModels.CrmSftpDownloadedFileModel)
		if len(fileInfoList) > 0 {
			for _, v := range fileInfoList {
				value := v
				fileInfoMap[value.FileName] = value
			}
		}
		for _, fileInfo := range fileLists {
			fileName := fileInfo.Name()
			if len(fileInfoMap) > 0 {
				if info, ok := fileInfoMap[fileName]; ok {
					if info.IsExtract == 0 {
						job.pullLogData(ctx, config.EnterpriseID, client, config, info, fileName, fileType)
					}
				} else {
					job.pullLogData(ctx, config.EnterpriseID, client, config, crmSftpFileModels.CrmSftpDownloadedFileModel{}, fileName, fileType)
				}
			} else {
				job.pullLogData(ctx, config.EnterpriseID, client, config, crmSftpFileModels.CrmSftpDownloadedFileModel{}, fileName, fileType)
			}
		}
	}

	return nil
}

func (job *SyncWedoKSInfo) pullLogData(ctx ctx.Context, enterpriseId int64, sftp *sftp.SFTP, config micro_scrm.BrandConfigSftpDetail, fileInfo crmSftpFileModels.CrmSftpDownloadedFileModel, fileName string, fileType int) error {
	// 同步cdp的会员信息
	filePath := config.PathConfig.CdpExternalUser[0]
	fileId := fileInfo.ID

	if fileInfo.IsExtract != 0 {
		return nil
	}
	sftpFile, err := sftp.Client.Open(SFTPROOTKS + filePath + fileName)
	if err != nil {
		log.Info(ctx).Msgf("%s::读取SFTP目录下的文件失败: %v", SFTPROOTKS+filePath+fileName, err)
		return err
	}

	LocalFileName := "wedo_" + strconv.FormatInt(app.UID(), 10) + `.csv`
	LocalFilePath := `./storage/wedo/` + LocalFileName

	dstFile, err := os.Create(LocalFilePath)
	if err != nil {
		log.Info(ctx).Msgf("%s::创建文件失败: %v", LocalFilePath, err)
		return err
	}
	if _, err := sftpFile.WriteTo(dstFile); err != nil {
		log.Info(ctx).Msgf("%s::写入文件失败: %v", LocalFilePath, err)
		return err
	}
	// 文件记录
	if fileId == 0 {

		cs := storage.Init(cos.COS)
		path, err := cs.Upload(ctx, LocalFilePath, "third_sftp/ks_cdp"+LocalFileName)
		if err != nil {
			log.Info(ctx).Msgf("uploadCloud Error: %v", err)
			return err
		}
		fileId = app.UID()
		data := crmSftpFileModels.CrmSftpDownloadedFileModel{
			MicroBaseModel: models.CreateMicroBase(enterpriseId, fileId),
			FileName:       fileName,
			FileFullPath:   "",
			FileType:       int32(fileType),
			FileURL:        cs.GetURL(path),
			IsExtract:      1,
		}
		err = thirdCrmSftpDao.NewThirdCrmSftpFile().CreateData(ctx, data)
		if err != nil {
			return err
		}
	}

	databyte, _ := ioutil.ReadFile(LocalFilePath)
	reader := csv.NewReader(bytes.NewBuffer(databyte))
	externalSvc := externalUserService.NewWxUserCorpExternal()

	var fileChan = make(chan string)
	// 发送任务到channel
	go func(reader *csv.Reader) {
		var flag = true
		for {
			row, err := reader.Read()
			if flag {
				flag = false
				continue
			}
			if err != nil && err != io.EOF {
				log.Info(ctx).Msgf("%s-%s::文件读取错误: %v", LocalFilePath, row, err)
				fmt.Println(err.Error())
				continue
			}
			if err == io.EOF {
				close(fileChan)
				break
			}
			fileChan <- strings.Join(row, ",")
		}
	}(reader)
	// 启动goroutine池
	var wg threading.RoutineGroup
	poolSize := 5
	execNum := 0
	execErrNum := 0

	for i := 0; i < poolSize; i++ {

		wg.Run(func() {

			for v := range fileChan {

				var cdpInfoKS external.CdpExternalInfo
				splitInfo := strings.Split(v, ",")

				if len(splitInfo) >= 6 {

					cdpInfoKS.UnionId = splitInfo[0]
					cdpInfoKS.CustomerCode = splitInfo[1]
					cdpInfoKS.Name = splitInfo[2]

					cdpInfoKS.CustomerLevel = splitInfo[6]
					last36MPurchaseAmount, _ := strconv.ParseFloat(splitInfo[9], 64)
					cdpInfoKS.Last36MPurchaseAmount = int64(last36MPurchaseAmount * 100)
					if splitInfo[11] != "" {
						cdpInfoKS.LastPurchaseDate = int(helper.GetDateToTimes(splitInfo[11]))
					}
					cdpInfoKS.LastPurchaseStoreCode = splitInfo[10]
				}

				cdpInfoKS.UnionIdBindFlag = 1
				updateData := map[string]interface{}{
					"union_code":               cdpInfoKS.CustomerCode,
					"name":                     cdpInfoKS.Name,
					"card_level":               cdpInfoKS.CustomerLevel,
					"amount_total":             cdpInfoKS.Last36MPurchaseAmount,
					"last_purchase_date":       cdpInfoKS.LastPurchaseDate,
					"last_purchase_store_code": cdpInfoKS.LastPurchaseStoreCode,
				}

				err = externalSvc.SyncCdpInfo(ctx, enterpriseId, cdpInfoKS, updateData)
				if err != nil {
					// fmt.Printf("企业%v::ks同步wedo会员信息数据报错,用户%s,错误%s", enterpriseId, cdpInfoKS.UnionId, err.Error())
					execErrNum++
				}
				execNum++
			}

		})
	}

	// 等待所有任务完成
	wg.Wait()

	err = thirdCrmSftpSvc.NewThirdCrmSftp().UpdateFileExTract(ctx, config.EnterpriseID, fileId, 2)
	planService.NewTaskNotice().SftpSendNotice(ctx, planService.SFTPMessageData{
		EnterpriseName: config.BrandName,
		SFTPName:       "SFTP同步CDP会员信息",
		UploadRes:      helper.RuleThree(err != nil, "异常", "成功"),
	})
	fmt.Printf("企业%v::同步用户信息完成:共处理 %v 行,错误共 %v 行", enterpriseId, execNum, execErrNum)

	return nil
}
