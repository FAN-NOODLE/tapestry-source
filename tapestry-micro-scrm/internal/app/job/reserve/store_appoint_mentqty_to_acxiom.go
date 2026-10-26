package reserveJob

import (
	"encoding/csv"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"io/ioutil"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	reserveDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	"micro-scrm/internal/app/models"
	departmentsModels "micro-scrm/internal/app/models/departments"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	reserveSvc "micro-scrm/internal/app/service/reserve"
	thirdCrmSftpSvc "micro-scrm/internal/app/service/third_crm_sftp"
	"micro-scrm/internal/modules/marketingPromotion/storage"
	"micro-scrm/internal/modules/marketingPromotion/storage/cos"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sftp"
	"micro-scrm/internal/pkg/xjob"
	"os"
	"path"
	"strconv"
)

type StoreAppointMENTQTYToAcxiom struct {
	name string
}

func NewStoreAppointMENTQTYToAcxiom() *StoreAppointMENTQTYToAcxiom {
	return &StoreAppointMENTQTYToAcxiom{
		name: "Monthly的Store成功邀约文件，每个财月结束后上传(STORE_APPOINTMENTQTY)",
	}
}

func (job *StoreAppointMENTQTYToAcxiom) Run(ctx ctx.Context, handler *xjob.JobHandler) error {
	sftpConfigList, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetSFTPConfigList(micro_scrm.SftpTypeAcxiom)
	if err != nil {
		log.Info(ctx).Msgf("%s::获取SFTP配置失败: %v", job.name, err)
	}

	for _, config := range sftpConfigList {

		if client, clientErr := sftp.NewSFTPClient(config.ConfigDetail); clientErr != nil {
			log.Info(ctx).Msgf("%s::企业ID:%d, 连接SFTP失败: %v", job.name, config.EnterpriseID, clientErr)
			// return
		} else {
			job.pushData(ctx, config.EnterpriseID, client, config)
		}
	}

	return nil
}

// 下一个财月的第二天会上传上一个财月的数据
func (job *StoreAppointMENTQTYToAcxiom) pushData(ctx ctx.Context, enterpriseId int64, sftp *sftp.SFTP, config micro_scrm.BrandConfigSftpDetail) error {
	if len(config.PathConfig.ChangYiReport) == 0 {
		return nil
	}

	// 根据当前时间的零点获取上个财月的信息,
	nowTime := helper.GetNowDateTimestamp()
	zeroTime := helper.GetZeroDateByParams(0)
	fiscalInfo, err := microWorkDao.NewFiscalMonth().GetFiscalMonth(ctx, enterpriseId, nowTime)
	if err != nil {
		return err
	}
	if fiscalInfo.StartTime == zeroTime { // 当天的0点和所在财月的开始时间一致，则不执行
		return nil
	}
	if fiscalInfo.StartTime+(86400*7) != zeroTime {
		log.Info(ctx).Msgf("%s::每个财月的第7天执行", job.name)
		return nil
	}
	fiscalPreMonthEndTime := fiscalInfo.StartTime - 1 // 算出上一个财月的结束时间
	fiscalPreInfo, err := microWorkDao.NewFiscalMonth().GetFiscalMonthByEndTime(ctx, enterpriseId, fiscalPreMonthEndTime)
	if err != nil {
		return err
	}
	startTime, endTime := fiscalPreInfo.StartTime, fiscalPreInfo.EndTime

	info, err := reserveDao.NewReserveOrderDataAcxiomSftpLogDao().GetInfo(ctx, enterpriseId, 0, startTime, "STORE_APPOINTMENTQTY")
	if err != nil {
		return err
	}
	if info.Status == 1 {
		log.Info(ctx).Msgf("%s::读取SFTP文件已处理完毕", job.name)
		return nil
	}

	list, err := reserveSvc.NewOrderService().GetStoreAppointmentLists(ctx, enterpriseId, startTime, endTime)

	if err != nil {
		fmt.Printf("获取Monthly的Store成功邀约,报错%v\n", err)
		return err
	}

	remotePath := config.PathConfig.ChangYiReport[1] // IN/WeClient_Report_PROD/STOREAPPOINTMENTQTY/
	LocalFileName := "STORE_APPOINTMENTQTY_" + helper.GetFormatDate(0, 2) + `.csv`
	LocalTextName := "STORE_APPOINTMENTQTY_" + helper.GetFormatDate(0, 2) + `OK.txt`
	FilePath := `./storage/acxiom/` + LocalFileName
	TextPath := `./storage/acxiom/` + LocalTextName
	t, _ := os.Create(TextPath) // 创建验证上传OK的文件
	defer t.Close()

	f, err := os.Create(FilePath) // 创建文件
	if err != nil {
		log.Info(ctx).Msgf("%s::csv文件创建失败: %v", job.name, err)
		return err
	}

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	data := [][]string{
		{"year", "period", "storeCode", "appointmentQTY"},
	}
	period := strconv.FormatInt(fiscalPreInfo.Period, 10)
	year := strconv.FormatInt(fiscalPreInfo.FiscalYear, 10)
	// 获取所有部门，departments表中name为门店code
	departments, err := departmentsDao.NewDepartments().GetAll(ctx, departmentsModels.DepartmentsModel{
		MicroBaseModel: models.MicroBaseModel{
			BaseModel: models.BaseModel{
				DeletedAt: 0,
			},
			EnterpriseID: enterpriseId,
		},
	})
	if err != nil {
		return err
	}
	departmentsMap := make(map[int]string)
	for _, department := range departments {
		departmentsMap[department.ThirdID] = department.Name
	}
	w := csv.NewWriter(f) // 创建一个新写入文件流
	defer f.Close()
	for _, v := range list {
		appointMentQTY := strconv.FormatInt(v.AppointmentQTY, 10)

		tmp := []string{
			year, period, departmentsMap[int(v.StoreId)], appointMentQTY,
		}
		data = append(data, tmp)
	}

	err = w.WriteAll(data)
	if err != nil {
		return err
	} // 写入数据
	w.Flush()
	fmt.Println("csv日志文件写入成功" + FilePath + ",日志数量:" + strconv.Itoa(len(list)))
	err = job.uploadFile(sftp, FilePath, remotePath)
	if err != nil {
		log.Info(ctx).Msgf("%s::sftp文件写入失败: %v", job.name, err)
		return err
	}
	err = job.uploadFile(sftp, TextPath, remotePath)
	if err != nil {
		log.Info(ctx).Msgf("%s::sftp验证OK文件写入失败: %v", job.name, err)
		return err
	}
	err = os.Remove(FilePath)
	if err != nil {
		log.Info(ctx).Msgf("delete local file error：%s ,%v", FilePath, err)
	}
	// 将csv文件上传直cos存储
	cs := storage.Init(cos.COS)
	cosPath, err := cs.Upload(ctx, FilePath, "acxiom/"+LocalFileName)
	if err != nil {
		log.Info(ctx).Msgf("uploadCloud Error: %v", err)
		return err
	}
	_ = os.Remove(TextPath)
	fileUrl := cs.GetURL(cosPath)
	insertData := microScrmModel.ReserveOrderDataAcxiomSftpLogModel{
		MicroBaseModel: models.MicroBaseModel{
			BaseModel:    models.CreateBase(app.UIDContext(ctx)),
			EnterpriseID: enterpriseId,
		},
		FileType:  0,
		LocalFile: LocalFileName,
		FileURL:   fileUrl,
		Status:    1,
		FileTime:  int(startTime),
		PushType:  "STORE_APPOINTMENTQTY",
		Year:      fiscalPreInfo.FiscalYear,
		Period:    fiscalPreInfo.Period,
	}
	err = reserveDao.NewReserveOrderDataAcxiomSftpLogDao().CreateLog(ctx, insertData)
	if err != nil {
		log.Info(ctx).Msgf("推送数据到acxiom记录日志失败:%v", err)
	}
	return nil
}

func (job *StoreAppointMENTQTYToAcxiom) uploadFile(sftp *sftp.SFTP, localFilePath string, remotePath string) error {
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
