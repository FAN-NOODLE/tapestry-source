package externalUserJob

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
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
	"sync"
	"sync/atomic"
	"time"
)

const (
	SFTPROOT = `/`
)

type SyncCdpInfo struct {
	name       string
	filePrefix string
}

func NewSyncCdpInfo() *SyncCdpInfo {
	return &SyncCdpInfo{
		name:       "同步外部联系人在cdp上的信息",
		filePrefix: "weclient_member_result_",
	}
}

func (job *SyncCdpInfo) Run(ctx ctx.Context, _ *xjob.JobHandler) error {

	sftpConfigList, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetSFTPConfigList(micro_scrm.SftpTypeAcxiom)
	if err != nil {
		log.Info(ctx).Msgf("%s::获取SFTP配置失败: %v", job.name, err)
	}

	for _, config := range sftpConfigList {
		if config.BrandName != "coach" {
			continue
		}
		fileType := crmSftpFileModels.Member_Info
		if client, clientErr := sftp.NewSFTPClient(config.ConfigDetail); clientErr != nil {
			log.Info(ctx).Msgf("%s::企业ID:%d, 连接SFTP失败: %v", job.name, config.EnterpriseID, clientErr)
			// return
		} else {
			if len(config.PathConfig.CdpExternalUser) == 0 {
				return nil
			}
			filePath := config.PathConfig.CdpExternalUser[0]
			fileLists, err := client.Client.ReadDir(SFTPROOT + filePath)
			if err != nil {
				log.Info(ctx).Msgf("%s::企业ID:%d, 获取目录下文件失败: %v", job.name, config.EnterpriseID, clientErr)
				return err
			}
			// 排序，顺序执行文件
			sort.Slice(fileLists, func(i, j int) bool {
				return fileLists[i].ModTime().Unix() < fileLists[i].ModTime().Unix()
			})
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
							_ = job.pullLogData(ctx, config.EnterpriseID, client, config, info, fileName, fileType)
						}
					} else {
						_ = job.pullLogData(ctx, config.EnterpriseID, client, config, crmSftpFileModels.CrmSftpDownloadedFileModel{}, fileName, fileType)
					}
				} else {
					_ = job.pullLogData(ctx, config.EnterpriseID, client, config, crmSftpFileModels.CrmSftpDownloadedFileModel{}, fileName, fileType)
				}

			}

		}

	}

	return nil
}

func (job *SyncCdpInfo) pullLogData(c ctx.Context, enterpriseId int64, sftp *sftp.SFTP, config micro_scrm.BrandConfigSftpDetail, fileInfo crmSftpFileModels.CrmSftpDownloadedFileModel, fileName string, fileType int) error {
	// 同步cdp的会员信息
	if len(config.PathConfig.CdpExternalUser) == 0 {
		return nil
	}
	filePath := config.PathConfig.CdpExternalUser[0]
	fileId := fileInfo.ID

	if fileInfo.IsExtract != 0 {
		return nil
	}
	// timeStr := helper.GetFormatDate(0, 2)
	// fileName := job.filePrefix + timeStr + ".csv"
	sftpFile, err := sftp.Client.Open(SFTPROOT + filePath + fileName)
	if err != nil {
		log.Info(c).Msgf("%s::读取SFTP目录下的文件失败: %v", SFTPROOT+filePath+fileName, err)
		return err
	}

	LocalFileName := "cdp_" + strconv.FormatInt(app.UID(), 10) + `.csv`
	LocalFilePath := `./storage/acxiom/` + LocalFileName

	dstFile, err := os.Create(LocalFilePath)
	if err != nil {
		log.Info(c).Msgf("%s::创建文件失败: %v", LocalFilePath, err)
		return err
	}
	if _, err := sftpFile.WriteTo(dstFile); err != nil {
		log.Info(c).Msgf("%s::写入文件失败: %v", LocalFilePath, err)
		return err
	}
	cs := storage.Init(cos.COS)
	path, err := cs.Upload(c, LocalFilePath, "third_sftp/cdp_info"+LocalFileName)
	if err != nil {
		log.Info(c).Msgf("uploadCloud Error: %v", err)
		return err
	}
	if fileId == 0 {
		fileId = app.UID()
		data := crmSftpFileModels.CrmSftpDownloadedFileModel{
			MicroBaseModel: models.CreateMicroBase(enterpriseId, fileId),
			FileName:       fileName,
			FileFullPath:   "",
			FileType:       int32(fileType),
			FileURL:        cs.GetURL(path),
			IsExtract:      1,
		}
		err := thirdCrmSftpDao.NewThirdCrmSftpFile().CreateData(c, data)
		if err != nil {
			return err
		}
	}

	databyte, _ := ioutil.ReadFile(LocalFilePath)
	reader := csv.NewReader(bytes.NewBuffer(databyte))
	externalSvc := externalUserService.NewWxUserCorpExternal()

	var fileChan = make(chan string, 2000)
	go func(reader *csv.Reader) {
		var flag = true
		for {
			row, err := reader.Read()
			if flag {
				flag = false
				continue
			}
			if err != nil && err != io.EOF {
				log.Info(c).Msgf("%s-%s::文件读取错误: %v", LocalFilePath, row, err)
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
	nowTime := time.Now().Unix()
	p, _ := ants.NewPool(5)
	defer p.Release()
	execNum := 0
	var execErrNum int64
	wg := sync.WaitGroup{}
	for v := range fileChan {
		var cdpInfo external.CdpExternalInfo
		splitInfo := strings.Split(v, ",")
		if len(splitInfo) >= 10 {
			cdpInfo.ID = splitInfo[0]
			cdpInfo.CustomerCode = splitInfo[1]
			cdpInfo.CustomerLevel = splitInfo[2]
			cdpInfo.UnionId = splitInfo[3]
			cdpInfo.Name = splitInfo[4]
			amount12m, _ := strconv.ParseFloat(splitInfo[5], 64)
			cdpInfo.Last12MPurchaseAmount = int64(amount12m) * 100
			amount36m, _ := strconv.ParseFloat(splitInfo[5], 64)
			cdpInfo.Last36MPurchaseAmount = int64(amount36m) * 100
			cdpInfo.MemberShip = splitInfo[7]
			cdpInfo.ClientType = splitInfo[8]
			cdpInfo.UnionIdBindFlag, _ = strconv.Atoi(splitInfo[9])
			if splitInfo[10] != "" {
				tmpTime, _ := time.Parse("2006-01-02", strings.TrimSpace(splitInfo[10]))
				cdpInfo.LastPurchaseDate = int(tmpTime.Unix())
			} else {
				cdpInfo.LastPurchaseDate = 0
			}

			cdpInfo.BirthDay = helper.EncodeToString(splitInfo[11])
			cdpInfo.BirthMonth = helper.EncodeToString(splitInfo[12])
			cdpInfo.BirthYear = helper.EncodeToString(splitInfo[13])
			cdpInfo.Birthday = helper.EncodeToString(splitInfo[13] + splitInfo[12] + splitInfo[11])
			// 下个月有积分将到期 0：不存在下个月积分过期 1：存在下个月积分过期
			cdpInfo.RemindExpiredPoints, _ = strconv.Atoi(splitInfo[14])
			cdpInfo.AgeTag = splitInfo[15]
		}

		updateData := map[string]interface{}{
			"union_code":               cdpInfo.MemberShip,
			"customer_code":            cdpInfo.CustomerCode,
			"customer_id":              cdpInfo.ID,
			"amount_total":             cdpInfo.Last12MPurchaseAmount,
			"card_level":               cdpInfo.CustomerLevel,
			"clt_type":                 cdpInfo.ClientType,
			"name":                     cdpInfo.Name,
			"last_purchase_date":       cdpInfo.LastPurchaseDate,
			"birthday":                 cdpInfo.Birthday,
			"birth_year":               cdpInfo.BirthYear,
			"birth_month":              cdpInfo.BirthMonth,
			"birth_day":                cdpInfo.BirthDay,
			"point":                    cdpInfo.RemindExpiredPoints, // 由于该字段暂时不用，先用point字段存储
			"age_tag":                  cdpInfo.AgeTag,
			"customer_type_updated_at": time.Now().Unix(), //增加一个会员信息更新时间
			"updated_at":               nowTime,
		}
		wg.Add(1)
		fn := func(cc ctx.Context, enterpriseId int64, cdpInfo external.CdpExternalInfo, updateData map[string]interface{}) func() {
			defer wg.Done()
			return func() {
				err = externalSvc.SyncCdpInfo(cc, enterpriseId, cdpInfo, updateData)
				if err != nil {
					// 此处不能打印，会报错
					// fmt.Printf("企业%v::cocah同步cdp会员信息数据报错,用户%s,错误%s", enterpriseId, cdpInfo.UnionId, err.Error())
					atomic.AddInt64(&execErrNum, 1)
				}
			}
		}
		_ = p.Submit(fn(c, enterpriseId, cdpInfo, updateData))
		execNum++
	}
	wg.Wait()

	err = thirdCrmSftpSvc.NewThirdCrmSftp().UpdateFileExTract(c, config.EnterpriseID, fileId, 2)
	planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
		EnterpriseName: config.BrandName,
		SFTPName:       "SFTP同步cdp会员信息",
		UploadRes:      helper.RuleThree(err != nil, "异常", "成功"),
	})
	fmt.Printf("企业%v::同步用户信息完成:共处理 %v 行,错误共 %v 行", enterpriseId, execNum, execErrNum)
	return err
}
