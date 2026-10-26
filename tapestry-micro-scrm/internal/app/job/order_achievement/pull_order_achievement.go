package orderAchievementJob

import (
	"bytes"
	"encoding/csv"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"io"
	"io/ioutil"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	thirdCrmSftpSvc "micro-scrm/internal/app/service/third_crm_sftp"
	"micro-scrm/internal/modules/marketingPromotion/storage"
	"micro-scrm/internal/modules/marketingPromotion/storage/cos"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sftp"
	"micro-scrm/internal/pkg/xjob"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	SFTPROOT = `/`
)

type PullOrderAchievement struct {
	name  string
	brand string
}

func NewPullOrderAchievement() *PullOrderAchievement {
	return &PullOrderAchievement{
		name:  "同步SA订单信息",
		brand: "",
	}
}

func (job *PullOrderAchievement) Run(ctx ctx.Context, handler *xjob.JobHandler) error {

	specsArgs := handler.Runner().Specs()
	typeValue := specsArgs[0]
	typeArr := []int{micro_scrm.SftpTypeWedo, micro_scrm.SftpTypeAcxiom}
	num, _ := strconv.Atoi(typeValue)
	if !helper.InArray(num, typeArr) {
		log.Info(ctx).Msgf("无效的配置参数")
		return nil
	}
	sftpConfigList, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetSFTPConfigList(num)

	if err != nil {
		log.Info(ctx).Msgf("%s::获取SFTP配置失败: %v", job.name, err)
	}
	for _, config := range sftpConfigList {
		job.brand = config.BrandName
		// 连接sftp
		if client, clientErr := sftp.NewSFTPClient(config.ConfigDetail); clientErr != nil {
			log.Info(ctx).Msgf("%s::企业ID:%d, 连接SFTP失败: %v", job.name, config.EnterpriseID, clientErr)
		} else {
			// 配置文件目录
			filePath := config.PathConfig.OrderAchievement
			if filePath == "" {
				continue
			}

			// sftp上的目录地址 fileLists
			fileLists, err := client.Client.ReadDir(SFTPROOT + filePath)
			if err != nil {
				log.Info(ctx).Msgf("%s::企业ID:%d, 获取目录下文件失败: %v", job.name, config.EnterpriseID, clientErr)
				return err
			}

			// 拉取数据库同步过的记录
			fileLogList, err := microScrmDao.NewOrderAchievementSftpLog().GetLogList(ctx, config.EnterpriseID)
			if err != nil {
				log.Info(ctx).Msgf("%s::企业ID:%d, 获取数据库记录失败: %v", job.name, config.EnterpriseID, clientErr)
				return err
			}

			// 把日志转化成map
			fileInfoMap := make(map[string]microScrmModel.OrderAchievementSftpLogModel)
			if len(fileLogList) > 0 {
				for _, v := range fileLogList {
					fileInfoMap[v.FileName] = v
				}
			}

			// 匹配文件处理业务
			for _, fileInfo := range fileLists {
				fileName := fileInfo.Name()
				filePath := SFTPROOT + filePath + fileInfo.Name()
				// 如果有日志记录
				if len(fileLogList) > 0 {
					// 判断有没有处理过
					if info, ok := fileInfoMap[fileName]; ok {
						if info.Status == 0 {
							job.pullLogData(ctx, config.EnterpriseID, client, filePath, fileName, info)
						}
					} else {
						// 不存在map里就是没有处理过的
						job.pullLogData(ctx, config.EnterpriseID, client, filePath, fileName, microScrmModel.OrderAchievementSftpLogModel{})
					}
				} else {
					// 没有处理过
					job.pullLogData(ctx, config.EnterpriseID, client, filePath, fileName, microScrmModel.OrderAchievementSftpLogModel{})
				}
			}
		}
	}

	return nil
}

func (job *PullOrderAchievement) pullLogData(ctx ctx.Context, enterpriseId int64, sftp *sftp.SFTP, filePath string, fileName string, info microScrmModel.OrderAchievementSftpLogModel) error {

	// 文件上传到cos
	err, localFilePath, id := job.fileToCos(ctx, sftp, filePath, fileName, enterpriseId, info)
	// 读取文件处理业务
	if err != nil {
		log.Info(ctx).Msgf("%s::读取SFTP目录下的文件失败: %v", fileName, err)
		return err
	}

	// 读取本地文件
	dataByte, _ := ioutil.ReadFile(localFilePath)
	reader := csv.NewReader(bytes.NewBuffer(dataByte))

	// 添加管道
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
				return
			}
			if err == io.EOF {
				close(fileChan)
				break
			}
			fileChan <- strings.Join(row, ",")
		}
	}(reader)
	// 消费管道
	fileNumber := 0
	successNumber := 0

	for v := range fileChan {
		// 每行数据
		splitInfo := strings.Split(v, ",")
		if job.brand == "ks" {
			splitInfo[0] = strings.Trim(splitInfo[0], " ")
			// KS传过来的文件 导购编码多出了前缀 000,需要去掉
			if splitInfo[0] != "000489350" { // 只有这一个用户确实有000，排除她
				splitInfo[0] = strings.Replace(splitInfo[0], "000", "", 1)
			}
		}
		fileNumber++

		if len(splitInfo) == 6 || len(splitInfo) == 5 {
			// 整理数据
			memberId, _ := membersDao.NewMembers().GetMemberIdByUserId(ctx, enterpriseId, splitInfo[0])
			var res bool
			// ks
			if len(splitInfo) == 5 {
				// 订单总数
				orderTotal, _ := strconv.Atoi(splitInfo[3])
				// 订单总金额
				orderPriceTotal, _ := strconv.ParseFloat(splitInfo[4], 64)
				data := microScrmModel.OrderAchievementModel{
					MemberID:        memberId,
					Userid:          splitInfo[0],
					StoreCode:       splitInfo[1],
					OrderTime:       helper.GetDateMorningTimestamp(splitInfo[2]),
					OrderTotal:      int64(orderTotal),
					OrderPriceTotal: int64(orderPriceTotal * 100),
					FileNumber:      int64(fileNumber),
				}
				res = microScrmDao.NewOrderAchievement().InsertAchievement(ctx, enterpriseId, data)
			}
			// sw
			if len(splitInfo) == 6 {
				// 订单总数

				// 订单总金额
				orderPriceTotal, _ := strconv.ParseFloat(splitInfo[5], 64)
				// 实付订单数量
				payOrderTotal, _ := strconv.Atoi(splitInfo[3])
				// 退单总数
				refundOrderTotal, _ := strconv.Atoi(splitInfo[4])
				orderTotal := payOrderTotal + refundOrderTotal
				data := microScrmModel.OrderAchievementModel{
					MemberID:         memberId,
					Userid:           splitInfo[0],
					StoreCode:        splitInfo[1],
					OrderTime:        helper.GetDateMorningTimestamp(splitInfo[2]),
					OrderTotal:       int64(orderTotal),
					OrderPriceTotal:  int64(orderPriceTotal * 100),
					PayOrderTotal:    int64(payOrderTotal),
					RefundOrderTotal: int64(refundOrderTotal),
					FileNumber:       int64(fileNumber),
				}
				res = microScrmDao.NewOrderAchievement().InsertAchievement(ctx, enterpriseId, data)
			}

			if res {
				successNumber++
			} else {
				log.Info(ctx).Msgf("%s::处理失败文件编号: %v", fileName, fileNumber)
			}
		} else {
			log.Info(ctx).Msgf("表头不正确")
			break
		}
	}

	// 更新日志记录
	err = microScrmDao.NewOrderAchievementSftpLog().UpdateLog(ctx, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"id":            id,
	}, microScrmModel.OrderAchievementSftpLogModel{
		Status:            1,
		SyncSuccessNumber: int64(successNumber),
		FileNumber:        int64(fileNumber),
	})
	if err != nil {
		log.Info(ctx).Msgf("%s::更新日志失败: %v", fileName, err)
		return err
	}

	err = job.delLocalFile(ctx, localFilePath)
	if err != nil {
		log.Info(ctx).Msgf("%s::删除本地文件失败: %v", fileName, err)
		return err
	}

	return nil
}

func (job *PullOrderAchievement) fileToCos(ctx ctx.Context, sftp *sftp.SFTP, filePath string, fileName string, enterpriseId int64, info microScrmModel.OrderAchievementSftpLogModel) (error, string, int64) {

	var (
		fileUrl       string
		localFilePath string
		logId         int64
	)

	// 读取文件
	sftpFile, err := sftp.Client.Open(filePath)
	if err != nil {
		log.Info(ctx).Msgf("%s::读取SFTP目录下的文件失败: %v", fileName, err)
		return err, localFilePath, logId
	}

	// 创建本地文件
	localFileName := strconv.FormatInt(app.UID(), 10) + `.csv`
	localFilePath = `./storage/upload/` + localFileName
	dstFile, err := os.Create(localFilePath)
	if err != nil {
		log.Info(ctx).Msgf("%s::创建文件失败: %v", localFilePath, err)
		return err, localFilePath, logId
	}
	if _, err := sftpFile.WriteTo(dstFile); err != nil {
		log.Info(ctx).Msgf("%s::写入文件失败: %v", localFilePath, err)
		return err, localFilePath, logId
	}

	// 上传cos
	cs := storage.Init(cos.COS)
	path, err := cs.Upload(ctx, localFilePath, "third_sftp/achievement"+localFileName)
	if err != nil {
		log.Info(ctx).Msgf("uploadCloud Error: %v", err)
		return err, localFilePath, logId
	}

	fileUrl = cs.GetURL(path)

	// 写日志记录
	if info.ID > 0 {
		// 删除原来的日志记录
		microScrmDao.NewOrderAchievementSftpLog().DelLog(ctx, map[string]interface{}{
			"id": info.ID,
		})
	}
	logId = app.UIDContext(ctx)
	insertData := microScrmModel.OrderAchievementSftpLogModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: enterpriseId,
		},
		ID:                logId,
		FileName:          fileName,
		FilePath:          localFilePath,
		FileURL:           fileUrl,
		Status:            0,
		FileTime:          time.Now().Unix(),
		FileNumber:        0,
		SyncSuccessNumber: 0,
	}
	err = microScrmDao.NewOrderAchievementSftpLog().CreateLog(ctx, insertData)
	if err != nil {
		log.Info(ctx).Msgf("拉取记录日志失败:%v", err)
		return err, localFilePath, logId
	}

	dstFile.Close()

	return nil, localFilePath, logId
}

func (job *PullOrderAchievement) delLocalFile(ctx ctx.Context, localFilePath string) error {
	// 删除本地文件
	err := os.Remove(localFilePath)
	if err != nil {
		log.Info(ctx).Msgf("delete local file error：%s ,%v", localFilePath, err)
		return err
	}

	return nil
}
