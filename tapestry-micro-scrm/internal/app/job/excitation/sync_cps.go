package excitation

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/tealeg/xlsx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
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

type SyncCps struct {
	Name string
}

func NewSyncCps() SyncCps {
	return SyncCps{Name: "同步cps文件"}
}

func (s SyncCps) Run(c ctx.Context, handler *xjob.JobHandler) error {
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

func (s SyncCps) handle(c ctx.Context, enterpriseId int64) error {
	info, err := microWorkDao.NewStoreSftpConfig().GetOne(c, map[string]interface{}{
		"brand":      "cps",
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
	fileInfoList, err := reserveDao.NewReserveOrderDataAcxiomSftpLogDao().GetExistsFile(c, enterpriseId, 1, "cps")
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
		ext := helper.GetFileExt(fileInfo.Name())
		if ext != ".xlsx" {
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

func (s SyncCps) pullLogData(ctx ctx.Context, enterpriseId int64, client *sftp.SFTP, info microScrmModel.ReserveOrderDataAcxiomSftpLogModel, fileName, brandName string) error {
	sftpFile, err := client.Client.Open(fileName)
	if err != nil {
		log.Info(ctx).Msgf("%s::读取SFTP目录下的文件失败: %v", fileName, err)
		return err
	}

	LocalFileName := strconv.FormatInt(app.UID(), 10) + helper.GetFileExt(fileName)
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
		path, err := cs.Upload(ctx, LocalFilePath, "third_sftp/cps/"+LocalFileName)
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
			PushType:  "cps",
		}
		err = reserveDao.NewReserveOrderDataAcxiomSftpLogDao().CreateLog(ctx, insertData)
		if err != nil {
			log.Info(ctx).Msgf("创建cps报错:%v", err)
		}
	}

	zr, err := zip.NewReader(bytes.NewReader(databyte), int64(len(databyte)))
	if err != nil {
		return err
	}
	fileData, err := xlsx.ReadZipReader(zr)
	if err != nil {
		return err
	}
	data, _ := readXlsxData(fileData.Sheets)
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
	salesDao := reserveDao.NewMemberSalesVolume()
	for _, v := range data {
		value := v
		saID := strings.TrimSpace(value.SaId)
		storeId := strings.TrimSpace(value.StoreId)
		var memberId, departmentId int64
		if value, ok := MemberMap[saID]; ok {
			memberId = value
		}
		if item, ok := departmentMap[storeId]; ok {
			departmentId = item
		}
		dataStr := value.Data
		tmpTime, err := time.Parse("2006-01-02", dataStr)
		if err != nil {
			return err
		}
		startTime := tmpTime.Unix()
		// startTime := helper.StrToTime(strings.TrimSpace(dataStr)).Unix()
		linkToEc := float64(0)
		if value.LinkToEc != "" {
			linkToEc, _ = strconv.ParseFloat(value.LinkToEc, 10)
		}
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
		nowTime := time.Now().Unix()
		if memberId > 0 && saleInfo.ID > 0 {
			_ = salesDao.UpdateData(ctx, enterpriseId, where, map[string]interface{}{
				"link_to_ec": linkToEc,
				"updated_at": nowTime,
			})
		} else {
			MemberSale := microScrmModel.MemberSalesVolumeModel{
				MemberID:     memberId,
				DepartmentID: departmentId,
				StartTime:    startTime,
				Date:         dataStr,
				StoreCode:    storeId,
				LinkToEc:     linkToEc,
				CreatedAt:    nowTime,
				UpdatedAt:    nowTime,
			}
			_ = salesDao.CreateInfo(ctx, enterpriseId, MemberSale)
		}
	}

	err = reserveDao.NewReserveOrderDataAcxiomSftpLogDao().UpdateLogStatus(ctx, enterpriseId, logId)
	if err != nil {
		return err
	}
	planService.NewTaskNotice().SftpSendNotice(ctx, planService.SFTPMessageData{
		EnterpriseName: brandName,
		SFTPName:       "SFTP同步Cps文件",
		UploadRes:      helper.RuleThree(err != nil, "异常", "成功"),
	})
	return nil
}

func readXlsxData(sheets []*xlsx.Sheet) ([]*microScrmModel.CpsRecord, error) {
	riMin := 1

	var cpsRecords []*microScrmModel.CpsRecord
	for index, sheet := range sheets {
		if index > 0 {
			return nil, errors.New("文件格式错误")
		}
		for ri, row := range sheet.Rows {
			if ri < riMin {
				continue
			}
			if len(row.Cells) == 0 {
				continue
			}
			var cpsRecord = &microScrmModel.CpsRecord{}

			if row.Cells[0].String() != "" {
				cpsRecord.Data = row.Cells[0].String()
				cpsRecord.SaId = row.Cells[1].String()
				cpsRecord.StoreId = row.Cells[2].String()
				cpsRecord.LinkToEc = row.Cells[3].String()
				cpsRecords = append(cpsRecords, cpsRecord)
			}

		}
	}
	return cpsRecords, nil
}
