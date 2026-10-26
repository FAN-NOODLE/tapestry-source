package excitation

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"io"
	"io/ioutil"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	membersDao "micro-scrm/internal/app/dao/members"
	reserveDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	appJob "micro-scrm/internal/app/job"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	planService "micro-scrm/internal/app/service/plan"
	"micro-scrm/internal/modules/marketingPromotion/storage"
	"micro-scrm/internal/modules/marketingPromotion/storage/cos"
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

type SyncCs2k struct {
	Name string
}

func NewSyncCs2k() SyncCs2k {
	return SyncCs2k{Name: "同步cs2k文件"}
}

func (s SyncCs2k) Run(c ctx.Context, handler *xjob.JobHandler) error {
	enterpriseIds, err := appJob.GetEnterpriseIds(handler)
	if err != nil {
		return err
	}
	for _, id := range enterpriseIds {
		enterpriseId := id
		err := s.handle(c, enterpriseId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s SyncCs2k) handle(c ctx.Context, enterpriseId int64) error {
	info, err := microWorkDao.NewStoreSftpConfig().GetOne(c, map[string]interface{}{
		"brand":      "cs2k",
		"deleted_at": 0,
	})
	if err != nil {
		return err
	}
	port, _ := strconv.Atoi(info.Port)
	config := sftp.SSHConfig{
		Host:     info.Host,
		Port:     int64(port),
		UserName: info.User,
		Password: info.Password,
	}
	client, err := sftp.NewSFTPClient(config)
	if err != nil {
		return err
	}
	fileLists, err := client.Client.ReadDir(SFTPROOT + info.RootPath)
	if err != nil {
		return err
	}
	if err != nil {
		log.Info(c).Msgf("%s::企业ID:%d, 获取目录下文件失败: %v", s.Name, enterpriseId, err)
		return err
	}
	fileInfoList, err := reserveDao.NewReserveOrderDataAcxiomSftpLogDao().GetExistsFile(c, enterpriseId, 1, "cs2k")
	if err != nil {
		log.Info(c).Msgf("%s::企业ID:%d, 获取目录下文件失败: %v", s.Name, enterpriseId, err)
		return err
	}
	fileInfoMap := make(map[string]microScrmModel.ReserveOrderDataAcxiomSftpLogModel)
	if len(fileInfoList) > 0 {
		for _, v := range fileInfoList {
			value := v
			fileInfoMap[value.LocalFile] = value
		}
	}

	for _, fileInfo := range fileLists {
		if fileInfo.IsDir() {
			continue
		}
		fileName := SFTPROOT + info.RootPath + fileInfo.Name()
		if len(fileInfoMap) > 0 {
			if fileModel, ok := fileInfoMap[fileName]; ok {
				if fileModel.Status == 0 {
					err := s.pullLogData(c, enterpriseId, client, fileModel, fileName, info.Brand)
					if err != nil {
						return err
					}
				}
			} else {
				err := s.pullLogData(c, enterpriseId, client, microScrmModel.ReserveOrderDataAcxiomSftpLogModel{}, fileName, info.Brand)
				if err != nil {
					return err
				}
			}
		} else {
			err := s.pullLogData(c, enterpriseId, client, microScrmModel.ReserveOrderDataAcxiomSftpLogModel{}, fileName, info.Brand)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (s SyncCs2k) pullLogData(ctx ctx.Context, enterpriseId int64, client *sftp.SFTP, info microScrmModel.ReserveOrderDataAcxiomSftpLogModel, fileName, brandName string) error {
	sftpFile, err := client.Client.Open(fileName)
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
	if logId == 0 {
		// 上传cos
		cs := storage.Init(cos.COS)
		path, err := cs.Upload(ctx, LocalFilePath, "third_sftp/cs2k/"+LocalFileName)
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
		insertData := microScrmModel.ReserveOrderDataAcxiomSftpLogModel{
			MicroBaseModel: models.MicroBaseModel{
				BaseModel:    models.CreateBase(logId),
				EnterpriseID: enterpriseId,
			},
			FileType:  1,
			LocalFile: sftpFile.Name(),
			FileURL:   fileUrl,
			Status:    0,
			FileTime:  int(time.Now().Unix()),
			PushType:  "cs2k",
		}
		err = reserveDao.NewReserveOrderDataAcxiomSftpLogDao().CreateLog(ctx, insertData)
		if err != nil {
			log.Info(ctx).Msgf("创建cs2k报错:%v", err)
		}
	}
	reader := csv.NewReader(bytes.NewBuffer(databyte))
	// 存在本地数据
	var fileChan = make(chan string, 2000)
	go func(reader *csv.Reader) {
		var flag = true
		for {
			row, errs := reader.Read()
			if flag {
				flag = false
				continue
			}
			if errs != nil && errs != io.EOF {
				fmt.Println(errs)
				continue
			}
			if errs == io.EOF {
				close(fileChan)
				break
			}
			fileChan <- strings.Join(row, ",")
		}
	}(reader)
	MemberLists, err := membersDao.NewMembers().GetMany(ctx, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"deleted_at":    0,
	})
	if err != nil {
		return err
	}
	MemberMap := make(map[string]int64, 0)
	for _, list := range MemberLists {
		info := list
		MemberMap[info.Userid] = info.ID
	}
	err, departmentLists := departmentsDao.NewDepartments().GetAllDepartments(ctx, enterpriseId)
	if err != nil {
		return err
	}
	departmentMap := make(map[string]int64, 0)
	for _, item := range departmentLists {
		value := item
		departmentMap[value.Name] = value.ID
	}
	defaultFloatValue := float64(0)
	salesDao := reserveDao.NewMemberSalesVolume()
	for v := range fileChan {
		splitInfo := strings.Split(v, ",")
		if len(splitInfo) >= 12 {
			var priceArr [13]float64
			for i, s2 := range splitInfo {
				if i >= 3 {
					if s2 != "" {
						price, _ := strconv.ParseFloat(s2, 10)
						priceArr[i] = price
					} else {
						priceArr[i] = defaultFloatValue
					}
				}
			}
			saID := strings.TrimSpace(splitInfo[1])
			storeId := strings.TrimSpace(splitInfo[2])
			var memberId, departmentId int64
			if value, ok := MemberMap[saID]; ok {
				memberId = value
			}
			if item, ok := departmentMap[storeId]; ok {
				departmentId = item
			}
			dataStr := strings.TrimSpace(splitInfo[0])
			tmpTime, err := time.Parse("2006-01-02", dataStr)
			if err != nil {
				return err
			}
			startTime := tmpTime.Unix()
			// startTime := helper.StrToTime(dataStr).Unix()
			nowTime := time.Now().Unix()
			where := map[string]interface{}{
				"start_time": startTime,
				"member_id":  memberId,
				"store_code": storeId,
			}
			saleInfo, err := salesDao.GetOne(ctx, enterpriseId, where)
			if err != nil {
				log.Info(ctx).Msgf("cps的执行查询数据报错:%v", err)
				continue
			}
			if memberId > 0 && saleInfo.ID > 0 {

				_ = salesDao.UpdateData(ctx, enterpriseId, where, map[string]interface{}{
					"home_store":    priceArr[3],
					"ll_eshop":      priceArr[4],
					"ll_livestream": priceArr[5],
					"pp_livestream": priceArr[6],
					"pp_post":       priceArr[7],
					"remote_sales":  priceArr[8],
					"cloud_store":   priceArr[9],
					"o2o":           priceArr[10],
					"try_at_home":   priceArr[11],
					"link_to_ec":    helper.RuleThree(priceArr[12] > 0, priceArr[12], saleInfo.LinkToEc),
					"updated_at":    nowTime,
				})
			} else {
				MemberSale := microScrmModel.MemberSalesVolumeModel{
					MemberID:     memberId,
					DepartmentID: departmentId,
					StartTime:    startTime,
					Date:         splitInfo[0],
					StoreCode:    storeId,
					HomeStore:    priceArr[3],
					LlEshop:      priceArr[4],
					LlLivestream: priceArr[5],
					PpLivestream: priceArr[6],
					PpPost:       priceArr[7],
					RemoteSales:  priceArr[8],
					CloudStore:   priceArr[9],
					O2o:          priceArr[10],
					TryAtHome:    priceArr[11],
					LinkToEc:     priceArr[12],
					CreatedAt:    nowTime,
					UpdatedAt:    nowTime,
				}
				_ = salesDao.CreateInfo(ctx, enterpriseId, MemberSale)
			}

		}
	}

	err = reserveDao.NewReserveOrderDataAcxiomSftpLogDao().UpdateLogStatus(ctx, enterpriseId, logId)
	if err != nil {
		return err
	}
	planService.NewTaskNotice().SftpSendNotice(ctx, planService.SFTPMessageData{
		EnterpriseName: brandName,
		SFTPName:       "SFTP同步Cs2K文件",
		UploadRes:      helper.RuleThree(err != nil, "异常", "成功"),
	})
	return nil
}
