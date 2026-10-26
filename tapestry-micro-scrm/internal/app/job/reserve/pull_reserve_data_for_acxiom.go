package reserveJob

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/shopspring/decimal"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"io"
	"io/ioutil"
	"micro-scrm/internal/app"
	reserveDao "micro-scrm/internal/app/dao/micro_scrm"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	reserveReq "micro-scrm/internal/app/http/request/reserve"
	"micro-scrm/internal/app/models"
	excitationModels "micro-scrm/internal/app/models/excitation"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	cacheSvc "micro-scrm/internal/app/service/cache"
	excitationSvc "micro-scrm/internal/app/service/excitation"
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
	SFTPROOT = `/`
)

type pullReserveDataForAcxiom struct {
	name string
}

func NewPullReserveDataForAcxiom() *pullReserveDataForAcxiom {
	return &pullReserveDataForAcxiom{
		name: "获取安客诚的订单数据，修改到店邀约单",
	}
}

func (job *pullReserveDataForAcxiom) Run(ctx ctx.Context, handler *xjob.JobHandler) error {

	sftpConfigList, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetSFTPConfigList(micro_scrm.SftpTypeAcxiom)
	if err != nil {
		log.Info(ctx).Msgf("%s::获取SFTP配置失败: %v", job.name, err)
	}

	for _, config := range sftpConfigList {

		if client, clientErr := sftp.NewSFTPClient(config.ConfigDetail); clientErr != nil {
			log.Info(ctx).Msgf("%s::企业ID:%d, 连接SFTP失败: %v", job.name, config.EnterpriseID, clientErr)
			// return
		} else {
			// 需要验证ks和sw，调用enterprise_config,PM要求三端企业都要有未完成提醒
			if len(config.PathConfig.Acxiom) == 0 {
				return nil
			}
			filePath := config.PathConfig.Acxiom[1]
			fileLists, err := client.Client.ReadDir(SFTPROOT + filePath)
			if err != nil {
				log.Info(ctx).Msgf("%s::企业ID:%d, 获取目录下文件失败: %v", job.name, config.EnterpriseID, clientErr)
				return err
			}
			fileInfoList, err := reserveDao.NewReserveOrderDataAcxiomSftpLogDao().GetExistsFile(ctx, config.EnterpriseID, 1, "pull_day_order")
			if err != nil {
				log.Info(ctx).Msgf("%s::企业ID:%d, 获取目录下文件失败: %v", job.name, config.EnterpriseID, clientErr)
				return err
			}
			fileInfoMap := make(map[string]microScrmModel.ReserveOrderDataAcxiomSftpLogModel)
			if len(fileInfoList) > 0 {
				for _, v := range fileInfoList {
					value := v
					fileInfoMap[value.LocalFile] = value
				}
			}
			sort.Slice(fileLists, func(i, j int) bool {
				return fileLists[i].ModTime().Unix() < fileLists[j].ModTime().Unix()
			})
			for _, fileInfo := range fileLists {
				fileName := SFTPROOT + filePath + fileInfo.Name()
				if len(fileInfoMap) > 0 {
					if info, ok := fileInfoMap[fileName]; ok {
						if info.Status == 0 {
							job.pullLogData(ctx, config.EnterpriseID, client, config, info, fileName)
						}
					} else {
						job.pullLogData(ctx, config.EnterpriseID, client, config, microScrmModel.ReserveOrderDataAcxiomSftpLogModel{}, fileName)
					}
				} else {
					job.pullLogData(ctx, config.EnterpriseID, client, config, microScrmModel.ReserveOrderDataAcxiomSftpLogModel{}, fileName)
				}
			}
		}

	}

	return nil

}

func (job *pullReserveDataForAcxiom) pullLogData(ctx ctx.Context, enterpriseId int64, sftp *sftp.SFTP, config micro_scrm.BrandConfigSftpDetail, info microScrmModel.ReserveOrderDataAcxiomSftpLogModel, fileName string) error {
	// 需要验证ks和sw，调用enterprise_config,PM要求三端企业都要有未完成提醒
	// if len(config.PathConfig.Acxiom) == 0 {
	//	return nil
	// }
	// filePath := config.PathConfig.Acxiom[1]
	// zeroTime := helper.GetDateMorningTimestamp(helper.GetNowDate())
	// info, err := reserveDao.NewReserveOrderDataAcxiomSftpLogDao().GetInfo(ctx, enterpriseId, 1, zeroTime, "")
	// if err != nil {
	//	return err
	// }
	// if info.ID > 0 && info.Status == 1 {
	//	log.Info(ctx).Msgf("%s::读取SFTP文件已处理完毕", job.name)
	//	return nil
	// }
	// fileInfoList, fileInfoErr := sftp.Client.ReadDir(SFTPROOT + filePath)
	// if fileInfoErr != nil {
	//	log.Info(ctx).Msgf("%s::读取SFTP文件时异常: %v", job.name, fileInfoErr)
	//	fmt.Println(fileInfoList)
	//	return fileInfoErr
	// }
	// if len(fileInfoList) == 0 {
	//	log.Info(ctx).Msgf("%s::该SFTP目录下没有文件: %v"+SFTPROOT+filePath, job.name, fileInfoErr)
	//	return fileInfoErr
	// }
	// sort.Slice(fileInfoList, func(i, j int) bool {
	//	return fileInfoList[i].ModTime().Unix() > fileInfoList[j].ModTime().Unix()
	// })

	// if fileInfoList[0].ModTime().Unix() < zeroTime {
	//	log.Info(ctx).Msgf("%s::该SFTP目录下%s天的文件还未同步过来: %v"+SFTPROOT+filePath, job.name, helper.GetNowDate(), fileInfoErr)
	//	return fileInfoErr
	// }
	// file := fileInfoList[0]
	enterpriseConfig, err := cacheSvc.NewRedisData().GetEnterpriseConfig(ctx, enterpriseId, "store_name")
	if err != nil {
		return err
	}
	if enterpriseConfig.Name == "coach" {
		return err
	}
	isSw := false
	if enterpriseConfig.Name == "sw" {
		isSw = true
	}
	hourTime := time.Now().Hour()
	sftpFile, err := sftp.Client.Open(fileName)
	if err != nil {
		log.Info(ctx).Msgf("%s::读取SFTP目录下的文件失败: %v", fileName, err)
		return err
	}

	LocalFileName := strconv.FormatInt(app.UID(), 10) + `.csv`
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

	databyte, _ := ioutil.ReadFile(LocalFilePath)
	logId := info.ID
	updatedAt := info.CreatedAt
	if logId == 0 {
		// 上传cos
		updatedAt = time.Now().Unix()
		cs := storage.Init(cos.COS)
		path, err := cs.Upload(ctx, LocalFilePath, "third_sftp/reserve_for_acxiom/"+LocalFileName)
		if err != nil {
			log.Info(ctx).Msgf("uploadCloud Error: %v", err)
			return err
		}
		// 删除本地文件
		err = os.Remove(LocalFilePath)
		if err != nil {
			log.Info(ctx).Msgf("delete local file error：%s ,%v", LocalFilePath, err)
		}
		fileUrl := cs.GetURL(path)
		logId = app.UIDContext(ctx)
		localFile := LocalFileName
		if hourTime > 12 {
			localFile = sftpFile.Name()
		}
		insertData := microScrmModel.ReserveOrderDataAcxiomSftpLogModel{
			MicroBaseModel: models.MicroBaseModel{
				BaseModel:    models.CreateBase(logId),
				EnterpriseID: enterpriseId,
			},
			FileType:  1,
			LocalFile: localFile,
			FileURL:   fileUrl,
			Status:    0,
			FileTime:  int(time.Now().Unix()),
			PushType:  "pull_day_order",
		}
		err = reserveDao.NewReserveOrderDataAcxiomSftpLogDao().CreateLog(ctx, insertData)
		if err != nil {
			log.Info(ctx).Msgf("推送数据到acxiom记录日志失败:%v", err)
		}
	}
	reader := csv.NewReader(bytes.NewBuffer(databyte))
	// 存在本地数据
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
	execNum := 0
	for v := range fileChan {
		var acxiomSftp reserveReq.AcxiomSftpInfo
		splitInfo := strings.Split(v, ",")
		if len(splitInfo) >= 5 {
			acxiomSftp.ReserveNumber = splitInfo[0]
			acxiomSftp.ReserveOrderTime = 0
			if splitInfo[1] != "" {
				acxiomSftp.ReserveOrderTime = int(helper.GetDateMorningTimestamp(splitInfo[1]))
			}
			floatPrice, _ := strconv.ParseFloat(splitInfo[2], 10)
			acxiomSftp.ReserveOrderPrice = floatPrice
			// acxiomSftp.ReserveOrderSn, _ = strconv.Atoi(splitInfo[3]) // 交易单号，目前coach和sw都不提供
			reserveOrderStatus, _ := strconv.Atoi(splitInfo[4]) // coach的下标4是匹配状态
			if isSw {
				reserveOrderStatus, _ = strconv.Atoi(splitInfo[5]) // sw的下标5是匹配状态
			}
			acxiomSftp.ReserveOrderStatus = int32(reserveOrderStatus)
		}
		reserveOrderPrice := decimal.NewFromFloat(acxiomSftp.ReserveOrderPrice).Mul(decimal.NewFromInt(100))
		updateMap := map[string]interface{}{
			"updated_at":           updatedAt,
			"reserve_order_time":   acxiomSftp.ReserveOrderTime,
			"reserve_order_price":  reserveOrderPrice,
			"reserve_order_status": acxiomSftp.ReserveOrderStatus,
		}
		if acxiomSftp.ReserveOrderStatus == 1 { // 安客诚订单匹配成功，将邀约单状态改为邀约成功
			updateMap["status"] = microScrmModel.FINISH
			updateMap["is_pay"] = microScrmModel.PAY_SUCCESS
			if isSw {
				updateMap["is_store"] = microScrmModel.STORE_YES
			}
			reserveDao.NewReservesOrder().Update(ctx, enterpriseId, map[string]interface{}{
				"reserve_number": acxiomSftp.ReserveNumber,
			}, updateMap)
			orderInfo, _ := reserveDao.NewReservesOrder().GetOrderInfo(ctx, enterpriseId, map[string]interface{}{
				"reserve_number": acxiomSftp.ReserveNumber,
			})
			if orderInfo.ID > 0 && hourTime > 12 {
				excitationSvc.NewPointsProduceService().ProduceQueue(ctx, excitationReq.PointsProduceSave{
					MemberId:     orderInfo.MemberID,
					Type:         excitationModels.RuleCategoryId4,
					Value:        1,
					UniqueHash:   acxiomSftp.ReserveNumber + fmt.Sprintf("-%d-%d", orderInfo.MemberID, excitationModels.RuleCategoryId4),
					EnterpriseId: enterpriseId,
				})
			}
		}
		if acxiomSftp.ReserveOrderStatus == 3 { // 安客成订单匹配过期，需要将邀约单改为已过期, 只处理邀约单为进行中的改为已过期，邀约完成和邀约成功的不更改邀约单状态
			updateMap["status"] = microScrmModel.OVERTIME
			reserveDao.NewReservesOrder().Update(ctx, enterpriseId, map[string]interface{}{
				"reserve_number": acxiomSftp.ReserveNumber,
			}, updateMap)
		}
		if acxiomSftp.ReserveOrderStatus == 2 { // 安客诚订单匹配取消，需要将订单修改到上次状态
			orderInfo, err := reserveDao.NewReservesOrder().GetOrderInfo(ctx, enterpriseId, map[string]interface{}{
				"reserve_number": acxiomSftp.ReserveNumber,
			})
			if err != nil {
				continue
			}
			if orderInfo.IsConfirm == microScrmModel.CONFIRM_YES {
				updateMap["status"] = microScrmModel.WAITFINISH
			} else {
				updateMap["status"] = microScrmModel.RUNNING
			}
			reserveDao.NewReservesOrder().Update(ctx, enterpriseId, map[string]interface{}{
				"reserve_number": acxiomSftp.ReserveNumber,
			}, updateMap)
		}
		execNum++
	}
	log.Info(ctx).Msgf("同步安客诚数据完毕，共执行:%d行记录", execNum)
	fmt.Printf("%s---同步安客诚数据完毕，共执行:%d行记录\n", fileName, execNum)
	if hourTime > 12 {
		reserveDao.NewReserveOrderDataAcxiomSftpLogDao().UpdateLogStatus(ctx, enterpriseId, logId)
		planService.NewTaskNotice().SftpSendNotice(ctx, planService.SFTPMessageData{
			EnterpriseName: enterpriseConfig.Name,
			SFTPName:       "sftp邀约单更新cdp数据",
			UploadRes:      helper.RuleThree(err != nil, "异常", "成功"),
		})
	}

	return nil
}
