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
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
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
	"time"
)

const (
	SFTPROOTSW = `/`
)

type SyncCdpSWInfo struct {
	name       string
	filePrefix string
}

func NewSyncCdpSWInfo() *SyncCdpSWInfo {
	return &SyncCdpSWInfo{
		name:       "同步外部联系人在cdp上的信息",
		filePrefix: "weclient_member_result_",
	}
}

func (job *SyncCdpSWInfo) Run(ctx ctx.Context, handler *xjob.JobHandler) error {

	var cond = map[string]interface{}{"brand_name": "sw", "sftp_type": 4}
	config, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetSFTPConfig(ctx, cond)
	if err != nil {
		log.Info(ctx).Msgf("%s::获取SFTP配置失败: %v", job.name, err)
	}

	if config.BrandName != "sw" {
		return nil
	}
	fileType := crmSftpFileModels.Member_Info_SW

	if client, clientErr := sftp.NewSFTPClient(config.ConfigDetail); clientErr != nil {
		log.Info(ctx).Msgf("%s::企业ID:%d, 连接SFTP失败: %v", job.name, config.EnterpriseID, clientErr)
		// return
	} else {
		if len(config.PathConfig.CdpExternalUser) == 0 {
			return nil
		}
		filePath := config.PathConfig.CdpExternalUser[0]
		// sftp读取文件
		fileLists, err := client.Client.ReadDir(SFTPROOTSW + filePath)
		// 排序，顺序执行文件
		sort.Slice(fileLists, func(i, j int) bool {
			return fileLists[i].ModTime().Unix() < fileLists[i].ModTime().Unix()
		})
		if err != nil {
			log.Info(ctx).Msgf("%s::企业ID:%d, 获取目录下文件失败: %v", job.name, config.EnterpriseID, clientErr)
			return err
		}

		// 查询在数据库的文件
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

func (job *SyncCdpSWInfo) pullLogData(ctx ctx.Context, enterpriseId int64, sftp *sftp.SFTP, config micro_scrm.BrandConfigSftpDetail, fileInfo crmSftpFileModels.CrmSftpDownloadedFileModel, fileName string, fileType int) error {
	// 同步cdp的会员信息
	filePath := config.PathConfig.CdpExternalUser[0]
	fileId := fileInfo.ID
	// 0未处理 1处理中 2处理完成 3处理失败
	if fileInfo.IsExtract != 0 {
		return nil
	}
	sftpFile, err := sftp.Client.Open(SFTPROOTSW + filePath + fileName)
	if err != nil {
		log.Info(ctx).Msgf("%s::读取SFTP目录下的文件失败: %v", SFTPROOTSW+filePath+fileName, err)
		return err
	}

	LocalFileName := "cdp_" + strconv.FormatInt(app.UID(), 10) + `.csv`
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
	if fileId == 0 {

		cs := storage.Init(cos.COS)
		path, err := cs.Upload(ctx, LocalFilePath, "third_sftp/sw_cdp"+LocalFileName)
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

	// 读取文件
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
	defer planService.NewTaskNotice().SftpSendNotice(ctx, planService.SFTPMessageData{
		EnterpriseName: config.BrandName,
		SFTPName:       "SFTP同步cdp会员信息",
		UploadRes:      helper.RuleThree(err != nil, "异常", "成功"),
	})
	// 启动goroutine池
	var wg threading.RoutineGroup

	poolSize := 5
	execNum := 0
	execErrNum := 0

	// 顾客类型（ 新客、老客、VIP客, 沉睡客）
	// 定义改为 数字标识  1新客、2老客、3VIP客、4沉睡客; 如果有多个中间 英文竖线隔开（|）
	// 生日或购买日期格式  年-月-日  yyyy-m-d

	for i := 0; i < poolSize; i++ {
		wg.Run(func() {

			for v := range fileChan {

				var cdpInfoSW external.CdpExternalInfo
				splitInfo := strings.Split(v, ",")

				if len(splitInfo) >= 9 {
					if splitInfo[2] == "" {
						continue
					}

					cdpInfoSW.CustomerCode = splitInfo[0]
					cdpInfoSW.CustomerLevel = splitInfo[1]
					cdpInfoSW.UnionId = splitInfo[2]
					Last12MPurchaseAmount, _ := strconv.ParseFloat(splitInfo[3], 64)
					cdpInfoSW.Last12MPurchaseAmount = int64(Last12MPurchaseAmount * 100)
					// 购买日期
					if splitInfo[4] != "" {
						var tempLastPurchaseDate string
						tempLastPurchaseDate = splitInfo[4]
						var tempLastPurchaseDateArr = strings.Split(tempLastPurchaseDate, "-")
						tempLastPurchaseDateYear, _ := strconv.Atoi(tempLastPurchaseDateArr[0])
						tempLastPurchaseDateMonth, _ := strconv.Atoi(tempLastPurchaseDateArr[1])
						tempLastPurchaseDateDay, _ := strconv.Atoi(tempLastPurchaseDateArr[2])
						lastPurchaseDate := time.Date(tempLastPurchaseDateYear, time.Month(tempLastPurchaseDateMonth), tempLastPurchaseDateDay, 0, 0, 0, 0, time.Local).Unix()
						cdpInfoSW.LastPurchaseDate = int(lastPurchaseDate)

					}
					cdpInfoSW.LastPurchaseStoreCode = splitInfo[5]
					upgradeNeedAmount, _ := strconv.ParseFloat(splitInfo[6], 64)
					cdpInfoSW.UpgradeNeedAmount = int64(upgradeNeedAmount * 100)
					// 匹配顾客类型
					tempCustomerTypeStr := splitInfo[7]
					if strings.Contains(tempCustomerTypeStr, "|") == true {
						var tempCustomerTypeArr []int
						tempCustomerTypeArrStr := strings.Split(tempCustomerTypeStr, "|")
						for _, value := range tempCustomerTypeArrStr {
							tempValue, _ := strconv.Atoi(value)
							tempCustomerTypeArr = append(tempCustomerTypeArr, tempValue)
						}
						if helper.InArray(microScrmModel.VipTypeChoice, tempCustomerTypeArr) == true {

							cdpInfoSW.VipType = microScrmModel.VipType

							item := []int{microScrmModel.VipTypeChoice}
							customerType := helper.DelItems(tempCustomerTypeArr, item)

							cdpInfoSW.CustomerType = customerType[0]

						}
					} else {
						cdpInfoSW.CustomerType, _ = strconv.Atoi(tempCustomerTypeStr)
					}

					// 生日截取和加密
					tempBirthday := splitInfo[8]
					if tempBirthday != "" {
						var tempBirthdayArr = strings.Split(tempBirthday, "-")
						cdpInfoSW.BirthYear = helper.EncodeToString(tempBirthdayArr[0])
						cdpInfoSW.BirthMonth = helper.EncodeToString(tempBirthdayArr[1])
						cdpInfoSW.BirthDay = helper.EncodeToString(tempBirthdayArr[2])
						cdpInfoSW.Birthday = helper.EncodeToString(tempBirthday)
					}

					cdpInfoSW.UnionIdBindFlag = 1
				}

				updateData := map[string]interface{}{
					"customer_code":            cdpInfoSW.CustomerCode,
					"card_level":               cdpInfoSW.CustomerLevel,
					"amount_total":             cdpInfoSW.Last12MPurchaseAmount,
					"last_purchase_date":       cdpInfoSW.LastPurchaseDate,
					"last_purchase_store_code": cdpInfoSW.LastPurchaseStoreCode,
					"upgrade_need_amount":      cdpInfoSW.UpgradeNeedAmount,
					"customer_type":            cdpInfoSW.CustomerType,
					"customer_type_updated_at": time.Now().Unix(),
					"birthday":                 cdpInfoSW.Birthday,
					"birth_year":               cdpInfoSW.BirthYear,
					"birth_month":              cdpInfoSW.BirthMonth,
					"birth_day":                cdpInfoSW.BirthDay,
					"vip_type":                 cdpInfoSW.VipType,
				}

				err := externalSvc.SyncCdpInfo(ctx, enterpriseId, cdpInfoSW, updateData)
				if err != nil {
					// fmt.Printf("企业%v::sw同步cdp会员信息数据报错,用户%s,错误%s", enterpriseId, cdpInfoSW.UnionId, err.Error())
					execErrNum++
				}
				execNum++
			}

		})
	}

	// 等待所有任务完成
	wg.Wait()

	err = thirdCrmSftpSvc.NewThirdCrmSftp().UpdateFileExTract(ctx, config.EnterpriseID, fileId, 2)
	fmt.Printf("企业%v::同步用户信息完成:共处理 %v 行,错误共 %v 行", enterpriseId, execNum, execErrNum)

	return nil
}
