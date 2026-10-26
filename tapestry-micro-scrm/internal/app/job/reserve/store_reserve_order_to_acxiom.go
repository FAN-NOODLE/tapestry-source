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
	"micro-scrm/internal/app/models"
	departmentsModels "micro-scrm/internal/app/models/departments"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	externalUserService "micro-scrm/internal/app/service/external_user"
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

type StoreReserveOrderToAcxiom struct {
	name string
}

func NewStoreReserveOrderToAcxiom() *StoreReserveOrderToAcxiom {
	return &StoreReserveOrderToAcxiom{
		name: "Daily 的SA成功邀约文件，所取当天日期-6 的成功邀约记录(Successful_Appointment)",
	}
}

func (job *StoreReserveOrderToAcxiom) Run(ctx ctx.Context, handler *xjob.JobHandler) error {
	sftpConfigList, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetSFTPConfigList(micro_scrm.SftpTypeAcxiom)
	if err != nil {
		log.Info(ctx).Msgf("%s::获取SFTP配置失败: %v", job.name, err)
	}

	for _, config := range sftpConfigList {

		if client, clientErr := sftp.NewSFTPClient(config.ConfigDetail); clientErr != nil {
			log.Info(ctx).Msgf("%s::企业ID:%d, 连接SFTP失败: %v", job.name, config.EnterpriseID, clientErr)
			// return
		} else {
			job.pullLogData(ctx, config.EnterpriseID, client, config)
		}

	}

	return nil
}

func (job *StoreReserveOrderToAcxiom) pullLogData(ctx ctx.Context, enterpriseId int64, sftp *sftp.SFTP, config micro_scrm.BrandConfigSftpDetail) error {
	if len(config.PathConfig.ChangYiReport) == 0 {
		return nil
	}

	startTime, endTime := helper.YesterdayPeriodOfTimestamp(6) // 6天前的successful
	// _, endTime := helper.PeriodOfTimestamp()             //当天的结束时间
	info, err := reserveDao.NewReserveOrderDataAcxiomSftpLogDao().GetInfo(ctx, enterpriseId, 0, startTime, "successful_appointment")
	if err != nil {
		return err
	}
	if info.Status == 1 {
		log.Info(ctx).Msgf("%s::SFTP文件已上传完毕", job.name)
		return nil
	}
	list, err := reserveSvc.NewOrderService().GetSuccessfulOrders(ctx, enterpriseId, startTime, endTime)
	if err != nil {
		fmt.Printf("获取邀约成功的预约单失败,报错%v\n", err)
		return err
	}

	remotePath := config.PathConfig.ChangYiReport[2] // IN/WeClient_Report_PROD/Succesful_Appointment/
	LocalFileName := "Successful_Appointment_" + helper.GetFormatDate(0, 2) + `.csv`
	LocalTextName := "Successful_Appointment_" + helper.GetFormatDate(0, 2) + `OK.txt`
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
		{"date", "storeCode", "saAccount", "memberID", "sendDate"},
	}
	var unionids []string
	for _, model := range list {
		unionids = append(unionids, model.UserUnionid)
	}
	externalLists, err := externalUserService.NewWxUserCorpExternal().GetWxUserCorpExternalByUnionIds(ctx, enterpriseId, unionids)
	if err != nil {
		log.Info(ctx).Msgf("%s::csv文件创建失败: %v", job.name, err)
		return err
	}
	externalMap := make(map[string]string)
	for _, external := range externalLists {
		if _, ok := externalMap[external.Unionid]; !ok {
			externalMap[external.Unionid] = external.UnionCode
		}
	}
	// 获取所有部门，departments表中name为门店code
	departments, err := departmentsDao.NewDepartments().GetAll(ctx, departmentsModels.DepartmentsModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: enterpriseId,
		}})
	if err != nil {
		return err
	}
	departmentsMap := make(map[int]string)
	for _, department := range departments {
		departmentsMap[department.ThirdID] = department.Name
	}
	w := csv.NewWriter(f) // 创建一个新写入文件流
	defer f.Close()
	uniqueMap := make(map[string]string)
	for _, v := range list {
		dateTimeStr := helper.GetFormatDateDay(int64(v.ReserveOrderTime))
		sendDataTimeStr := helper.GetFormatDateDay(v.CreatedAt)
		key := dateTimeStr + "-" + departmentsMap[int(v.StoreID)] + "-" + v.MemberUserID + "-" + externalMap[v.UserUnionid] + "-" + sendDataTimeStr
		if _, ok := uniqueMap[key]; !ok {
			uniqueMap[key] = key
			tmp := []string{
				dateTimeStr, departmentsMap[int(v.StoreID)], v.MemberUserID, externalMap[v.UserUnionid], sendDataTimeStr,
			}
			data = append(data, tmp)
		}
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

	// 将csv文件上传直cos存储
	cs := storage.Init(cos.COS)
	cosPath, err := cs.Upload(ctx, FilePath, "acxiom/"+LocalFileName)
	if err != nil {
		log.Info(ctx).Msgf("uploadCloud Error: %v", err)
		return err
	}
	// 删除本地文件
	err = os.Remove(FilePath)
	if err != nil {
		log.Info(ctx).Msgf("delete local file error：%s ,%v", FilePath, err)
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
		PushType:  "successful_appointment",
	}
	err = reserveDao.NewReserveOrderDataAcxiomSftpLogDao().CreateLog(ctx, insertData)
	if err != nil {
		log.Info(ctx).Msgf("推送数据到acxiom记录日志失败:%v", err)
	}
	return nil
}

func (job *StoreReserveOrderToAcxiom) uploadFile(sftp *sftp.SFTP, localFilePath string, remotePath string) error {
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
