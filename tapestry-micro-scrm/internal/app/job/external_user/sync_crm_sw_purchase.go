package externalUserJob

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"io"
	"io/ioutil"
	"micro-scrm/internal/app"
	externalSearchDao "micro-scrm/internal/app/dao/external_search"
	swCrmCustomerDao "micro-scrm/internal/app/dao/sw_crm_customer"
	thirdCrmSftpDao "micro-scrm/internal/app/dao/third_crm_sftp"
	"micro-scrm/internal/app/models"
	crmSftpFileModels "micro-scrm/internal/app/models/crm_sftp_file"
	swCrmCustomerModel "micro-scrm/internal/app/models/sw_crm_customer"
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

type SyncCrmSWPurchase struct {
	name string
}

func NewSyncCrmSWPurchase() *SyncCrmSWPurchase {
	return &SyncCrmSWPurchase{
		name: "同步外部联系人在crm上的信息",
	}
}

func (job *SyncCrmSWPurchase) Run(ctx ctx.Context, handler *xjob.JobHandler) error {

	var cond = map[string]interface{}{"brand_name": "sw", "sftp_type": 4}
	config, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetSFTPConfig(ctx, cond)
	if err != nil {
		log.Info(ctx).Msgf("%s::获取SFTP配置失败: %v", job.name, err)
	}

	fileTypes := []int{crmSftpFileModels.Customer_Purchase_SW, crmSftpFileModels.Customer_Extension_SW}
	if client, clientErr := sftp.NewSFTPClient(config.ConfigDetail); clientErr != nil {
		log.Info(ctx).Msgf("%s::企业ID:%d, 连接SFTP失败: %v", job.name, config.EnterpriseID, clientErr)
		return nil
	} else {
		if len(config.PathConfig.SwCrmPurchase) == 0 {
			return nil
		}
		// 查询在数据库的文件
		fileInfoList, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetExistsFile(ctx, config.EnterpriseID, fileTypes)
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
		for _, filePath := range config.PathConfig.SwCrmPurchase {
			// sftp读取文件
			fileLists, err := client.Client.ReadDir(SFTPROOTSW + filePath)
			if err != nil {
				log.Info(ctx).Msgf("%s::企业ID:%d, 获取目录%s文件失败: %v", job.name, config.EnterpriseID, SFTPROOTSW+filePath, err)
				continue
			}
			// 排序，顺序执行文件
			sort.Slice(fileLists, func(i, j int) bool {
				return fileLists[i].ModTime().Unix() < fileLists[i].ModTime().Unix()
			})
			for _, fileInfo := range fileLists {
				fileName := fileInfo.Name()
				if len(fileInfoMap) > 0 {
					if _, ok := fileInfoMap[fileName]; !ok {
						_ = job.pullLogData(ctx, config.EnterpriseID, client, config, fileName, filePath)
					}
				} else {
					_ = job.pullLogData(ctx, config.EnterpriseID, client, config, fileName, filePath)
				}
			}
		}
		err = externalSearchDao.NewExternalExtraCkSw().MergeData(ctx, config.EnterpriseID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (job *SyncCrmSWPurchase) pullLogData(ctx ctx.Context, enterpriseId int64, sftp *sftp.SFTP, config micro_scrm.BrandConfigSftpDetail, fileName string, filePath string) error {
	// 同步crm的会员信息
	sftpFile, err := sftp.Client.Open(SFTPROOTSW + filePath + fileName)
	if err != nil {
		log.Info(ctx).Msgf("%s::读取SFTP目录下的文件失败: %v", SFTPROOTSW+filePath+fileName, err)
		return err
	}
	LocalFileName := "crm_" + strconv.FormatInt(app.UID(), 10) + `.csv`
	LocalFilePath := `./storage/acxiom/` + LocalFileName

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
	cs := storage.Init(cos.COS)
	path, err := cs.Upload(ctx, LocalFilePath, "third_sftp/sw_crm"+LocalFileName)
	if err != nil {
		log.Info(ctx).Msgf("uploadCloud Error: %v", err)
		return err
	}
	var fileTypeMap = map[string]int{
		"/out/SW/customer_purchase/":  crmSftpFileModels.Customer_Purchase_SW,
		"/out/SW/customer_extension/": crmSftpFileModels.Customer_Extension_SW,
	}
	var fileType int
	if tmpFileType, ok := fileTypeMap[filePath]; ok {
		fileType = tmpFileType
	}
	fileId := app.UID()
	data := crmSftpFileModels.CrmSftpDownloadedFileModel{
		MicroBaseModel: models.CreateMicroBase(enterpriseId, fileId),
		FileName:       fileName,
		FileFullPath:   filePath + fileName,
		FileType:       int32(fileType),
		FileURL:        cs.GetURL(path),
		IsExtract:      1,
	}
	err = thirdCrmSftpDao.NewThirdCrmSftpFile().CreateData(ctx, data)
	if err != nil {
		return err
	}

	// 读取文件
	dataByte, _ := ioutil.ReadFile(LocalFilePath)
	reader := csv.NewReader(bytes.NewBuffer(dataByte))

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
	defer planService.NewTaskNotice().SftpSendNotice(ctx, planService.SFTPMessageData{
		EnterpriseName: config.BrandName,
		SFTPName:       "SFTP同步SW-Crm会员消费信息" + fileName,
		UploadRes:      helper.RuleThree(err != nil, "异常", "成功"),
	})

	execNum := 0
	execErrNum := 0
	if fileType == crmSftpFileModels.Customer_Purchase_SW {
		var addSlice = make([]swCrmCustomerModel.CrmCustomerPurchaseModel, 0)
		for v := range fileChan {
			splitInfo := strings.Split(v, ",")
			if len(splitInfo) < 6 || splitInfo[1] == "" {
				continue
			}
			addSlice = append(addSlice, swCrmCustomerModel.CrmCustomerPurchaseModel{
				BaseModel:    models.CreateBase(app.UID()),
				CustomerID:   splitInfo[0],
				UnionID:      splitInfo[1],
				SaCode:       splitInfo[2],
				CategoryCode: splitInfo[3],
				PurchaseCode: splitInfo[4],
				PurchaseTime: splitInfo[5],
			})
			if len(addSlice) >= 1000 {
				err := job.savePurchaseData(ctx, enterpriseId, addSlice)
				if err != nil {
					execErrNum += len(addSlice)
				}
				execNum += len(addSlice)
				addSlice = make([]swCrmCustomerModel.CrmCustomerPurchaseModel, 0)
			}
		}
		err := job.savePurchaseData(ctx, enterpriseId, addSlice)
		if err != nil {
			execErrNum += len(addSlice)
		}
		execNum += len(addSlice)
	} else if fileType == crmSftpFileModels.Customer_Extension_SW {
		var addSlice = make([]swCrmCustomerModel.CrmCustomerExtensionModel, 0)
		for v := range fileChan {
			splitInfo := strings.Split(v, ",")
			if len(splitInfo) < 6 || splitInfo[1] == "" {
				continue
			}
			amount, _ := strconv.ParseFloat(splitInfo[4], 64)
			addSlice = append(addSlice, swCrmCustomerModel.CrmCustomerExtensionModel{
				BaseModel:         models.CreateBase(app.UID()),
				CustomerID:        splitInfo[0],
				UnionID:           splitInfo[1],
				PriceTag:          splitInfo[2],
				AgeTag:            splitInfo[3],
				R12PurchaseAmount: amount,
				AtValues:          splitInfo[5],
			})
			if len(addSlice) >= 1000 {
				err := job.saveExtensionData(ctx, enterpriseId, addSlice)
				if err != nil {
					execErrNum += len(addSlice)
				}
				execNum += len(addSlice)
				addSlice = make([]swCrmCustomerModel.CrmCustomerExtensionModel, 0)
			}
		}
		err := job.saveExtensionData(ctx, enterpriseId, addSlice)
		if err != nil {
			execErrNum += len(addSlice)
		}
		execNum += len(addSlice)
	}

	err = thirdCrmSftpSvc.NewThirdCrmSftp().UpdateFileExTract(ctx, config.EnterpriseID, fileId, 2)
	fmt.Printf("企业%v::同步用%s完成:共处理 %v 行,错误共 %v 行", enterpriseId, fileName, execNum, execErrNum)

	return nil
}

func (job *SyncCrmSWPurchase) savePurchaseData(ctx ctx.Context, enterpriseId int64, addSlice []swCrmCustomerModel.CrmCustomerPurchaseModel) error {
	_ = swCrmCustomerDao.NewPurchaseDao().CreateInBatchesData(ctx, enterpriseId, addSlice)
	return externalSearchDao.NewExternalExtraCkSw().SavePurchaseData(ctx, enterpriseId, addSlice)
}

func (job *SyncCrmSWPurchase) saveExtensionData(ctx ctx.Context, enterpriseId int64, addSlice []swCrmCustomerModel.CrmCustomerExtensionModel) error {
	_ = swCrmCustomerDao.NewExtensionDao().CreateInBatchesData(ctx, enterpriseId, addSlice)
	_ = externalSearchDao.NewExternalTagCkSw().SaveAtValuesData(ctx, enterpriseId, addSlice)
	return externalSearchDao.NewExternalExtraCkSw().SaveExtensionData(ctx, enterpriseId, addSlice)
}
