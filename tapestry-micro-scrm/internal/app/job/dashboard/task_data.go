package dashboardJob

import (
	"encoding/csv"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"io/ioutil"
	"micro-scrm/internal/app"
	customerDao "micro-scrm/internal/app/dao/customer"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	appJob "micro-scrm/internal/app/job"
	cacheSvc "micro-scrm/internal/app/service/cache"
	planService "micro-scrm/internal/app/service/plan"
	thirdCrmSftpSvc "micro-scrm/internal/app/service/third_crm_sftp"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sftp"
	"micro-scrm/internal/pkg/xjob"
	"os"
	"path"
	"strconv"
	"time"
)

type TaskData struct {
	name       string
	startTime  int64
	endTime    int64
	sftpClient *sftp.SFTP
	remotePath string
}

func NewTaskData() *TaskData {
	// startTime, endTime := int64(0), int64(0) // 全量数据，设置startTime的值为0，增量数据取昨天的时间戳
	var startTime, endTime int64 = helper.GetBetweenByTime(1)
	return &TaskData{
		name:       "任务数据dashboard",
		startTime:  startTime,
		endTime:    endTime,
		sftpClient: nil,
		remotePath: "TRANSFERIN/",
	}
}

func (job *TaskData) Run(c ctx.Context, handler *xjob.JobHandler) error {
	enterpriseIds, err := appJob.GetEnterpriseIds(handler)
	if err != nil {
		return err
	}
	var wg threading.RoutineGroup
	for _, id := range enterpriseIds {
		tmpId := id
		wg.Run(func() {
			err = job.Handle(c, tmpId)
			if err != nil {
				app.Logger().Error(c).Msgf("error : %v", err)
			}
		})
	}
	wg.Wait()
	return nil
}

// Handle 同步dashboard数据
func (job *TaskData) Handle(c ctx.Context, enterpriseId int64) error {
	config, err := cacheSvc.NewRedisData().GetEnterpriseConfig(c, enterpriseId, "store_name")
	if err != nil {
		return err
	}
	if config.Name != "coach" {
		return nil
	}
	sftpConfigList, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetSFTPConfigList(micro_scrm.SftpTypeKUNXIN)
	if err != nil {
		app.Logger().Error(c).Msgf("%s::获取SFTP配置失败: %v", job.name, err)
	}
	for _, config := range sftpConfigList {
		if client, clientErr := sftp.NewSFTPClient(config.ConfigDetail); clientErr != nil {
			message := planService.SFTPMessageData{
				EnterpriseName: config.BrandName,
				SFTPName:       "KUNXIN的sftp连接失败",
				UploadRes:      fmt.Sprintf("失败,错误为：%s", clientErr.Error()),
			}
			planService.NewTaskNotice().SftpSendNotice(c, message)
			return clientErr
		} else {

			job.sftpClient = client
		}
	}
	// 主任务的dashboard
	err = job.createPlanData(c, enterpriseId)
	if err != nil {
		app.Logger().Error(c).Msgf("生成主任务表的dashboard报错，err:%v", err)
	}
	// 子任务的dashboard
	err = job.createPlanTaskData(c, enterpriseId)
	if err != nil {
		app.Logger().Error(c).Msgf("生成子任务表的dashboard报错，err:%v", err)
	}

	// 任务导购表数据
	err = job.createTaskAllotLogData(c, enterpriseId)
	if err != nil {
		app.Logger().Error(c).Msgf("生成任务导购表的dashboard报错，err:%v", err)
	}

	// 邀约记录表数据
	err = job.createReserveData(c, enterpriseId)
	if err != nil {
		app.Logger().Error(c).Msgf("生成邀约记录表的dashboard报错，err:%v", err)
	}

	// 邀约服务表数据
	err = job.createReserveServiceData(c, enterpriseId)
	if err != nil {
		app.Logger().Error(c).Msgf("生成邀约服务表的dashboard报错，err:%v", err)
	}

	// 播种表数据
	err = job.createSowTaskData(c, enterpriseId)
	if err != nil {
		app.Logger().Error(c).Msgf("生成播种表的dashboard报错，err:%v", err)
	}

	// 任务和销售的数据记录表，取updated_at值在当天的数据，因为acxiom给的销售记录是T-1的数据，给KUNXIN的数据要满足T-1，则需要给今天的销售记录
	// 任务销售表数据
	err = job.createTaskSalesData(c, enterpriseId)
	if err != nil {
		app.Logger().Error(c).Msgf("生成任务销售表的dashboard报错，err:%v", err)
	}

	// 邀约销售表数据
	err = job.createReserveSalesData(c, enterpriseId)
	if err != nil {
		app.Logger().Error(c).Msgf("生成邀约销售表的dashboard报错，err:%v", err)
	}

	return err
}

// 主任务表数据
func (job *TaskData) createPlanData(c ctx.Context, enterpriseId int64) error {
	// 生成主任务表的数据dashboard
	many, err := planDao.NewPlanTaskDao().GetPlanMany(c, enterpriseId, "updated_at>=? and updated_at<=?", job.startTime, job.endTime)
	if err != nil {
		return err
	}
	data := [][]string{
		{"id", "name", "task_num", "updated_at", "start_time", "end_time", "remark"},
	}
	for _, v := range many {
		tmp := []string{
			strconv.FormatInt(v.ID, 10), v.Name, strconv.Itoa(v.TaskNum), strconv.FormatInt(v.UpdatedAt, 10), strconv.Itoa(v.StartTime), strconv.Itoa(v.EndTime), v.Remark,
		}
		data = append(data, tmp)
	}
	filePrefix := "plans"
	err = job.createCsv(c, enterpriseId, filePrefix, data)
	return err
}

// 子任务表数据
func (job *TaskData) createPlanTaskData(c ctx.Context, enterpriseId int64) error {
	// 生成子任务表的数据dashboard
	many, err := planDao.NewPlanTaskDao().GetMany(c, enterpriseId, "updated_at>=? and updated_at<=?", job.startTime, job.endTime)
	if err != nil {
		return err
	}
	data := [][]string{
		{"id", "plan_id", "task_category", "name", "remark", "task_type", "plan_send_time", "updated_at", "task_start_time", "task_end_time"},
	}
	for _, v := range many {
		tmp := []string{
			strconv.FormatInt(v.ID, 10), strconv.FormatInt(v.PlanID, 10), strconv.Itoa(v.TaskCategory), v.Name, v.Remark, strconv.Itoa(int(v.TaskType)), strconv.FormatInt(v.PlanSendTime, 10), strconv.FormatInt(v.UpdatedAt, 10), strconv.FormatInt(v.TaskStartTime, 10), strconv.FormatInt(v.TaskEndTime, 10),
		}
		data = append(data, tmp)
	}
	filePrefix := "plan_tasks"
	err = job.createCsv(c, enterpriseId, filePrefix, data)
	return err
}

// 任务导购表数据
func (job *TaskData) createTaskAllotLogData(c ctx.Context, enterpriseId int64) error {
	// 生成任务导购task_allot_log表的数据dashboard
	data := [][]string{
		{"plan_id", "task_id", "member_user_id", "end_time", "is_finished", "customer_num", "updated_at", "contacted_count", "moment_look_count", "no_contacted_count", "current_contacted_count"},
	}

	many, err := planDao.NewPlanTaskAllotLogDao().GetMany(c, enterpriseId, "updated_at>=? and updated_at<=?", job.startTime, job.endTime)
	if err != nil {
		return err
	}
	for _, v := range many {
		customerNum := strconv.Itoa(v.CustomerNum)
		contactedCount := strconv.Itoa(v.ContactedCount)
		noContactedCount := v.CustomerNum - v.ContactedCount
		momentLookCount := v.MomentLookCount
		currentCount, _ := planDao.NewPlanTaskConditionDataDao().GetDistinctUserCount(c, enterpriseId, v.PlanId, "task_allot_log_id=? and finished_at between ? and ?", v.ID, job.startTime, job.endTime)
		currentNum := strconv.Itoa(int(currentCount))
		tmp := []string{
			strconv.FormatInt(v.PlanId, 10), strconv.FormatInt(v.TaskId, 10), v.UserId, strconv.Itoa(v.EndTime), strconv.Itoa(v.IsFinished), customerNum, strconv.FormatInt(v.UpdatedAt, 10), contactedCount, strconv.Itoa(momentLookCount), strconv.Itoa(noContactedCount), currentNum,
		}
		data = append(data, tmp)
	}
	filePrefix := "task_allot_log"
	err = job.createCsv(c, enterpriseId, filePrefix, data)
	return err
}

// 邀约记录表数据
func (job *TaskData) createReserveData(c ctx.Context, enterpriseId int64) error {

	getMany, err := departmentsDao.NewDepartments().GetMany(c, map[string]interface{}{"enterprise_id": enterpriseId})
	if err != nil {
		return err
	}
	storeNameMap := make(map[int]string)
	for _, model := range getMany {
		storeNameMap[model.ThirdID] = model.Name
	}
	data := [][]string{
		{"reserve_number", "reserve_id", "reserve_time", "store_code", "status", "member_user_id", "is_store", "is_confirm", "updated_at", "created_at"},
	}
	endTime := time.Now().Unix()
	many, err := planDao.NewReservesOrder().GetMany(c, enterpriseId, "updated_at>=? and updated_at<=?", job.startTime, endTime)
	if err != nil {
		return err
	}
	for _, v := range many {
		storeCode := storeNameMap[int(v.StoreID)]
		tmp := []string{
			v.ReserveNumber, strconv.FormatInt(v.ReserveID, 10), strconv.Itoa(v.ReserveTime),
			storeCode, strconv.Itoa(int(v.Status)), v.MemberUserID, strconv.Itoa(int(v.IsStore)),
			strconv.Itoa(int(v.IsConfirm)), strconv.FormatInt(v.UpdatedAt, 10), strconv.FormatInt(v.CreatedAt, 10),
		}
		data = append(data, tmp)
	}

	filePrefix := "reserve_order"
	err = job.createCsv(c, enterpriseId, filePrefix, data)
	return err
}

// 邀约服务表数据
func (job *TaskData) createReserveServiceData(c ctx.Context, enterpriseId int64) error {
	many, err := planDao.NewReserveServiceDao().GetMany(c, enterpriseId, "updated_at>=? and updated_at<=?", job.startTime, job.endTime)
	if err != nil {
		return err
	}
	data := [][]string{
		{"reserve_id", "reserve_name", "updated_at"},
	}
	for _, v := range many {
		tmp := []string{
			strconv.FormatInt(v.ID, 10), v.ReserveName, strconv.FormatInt(v.UpdatedAt, 10),
		}
		data = append(data, tmp)
	}
	filePrefix := "reserve_service"
	err = job.createCsv(c, enterpriseId, filePrefix, data)
	return err
}

// 播种表数据
func (job *TaskData) createSowTaskData(c ctx.Context, enterpriseId int64) error {
	data := [][]string{
		{"id", "member_user_id", "status", "created_at", "remind_time", "updated_at"},
	}

	many, err := customerDao.NewSowTask().GetMany(c, "enterprise_id=? and updated_at>=? and updated_at<=?", enterpriseId, job.startTime, job.endTime)
	if err != nil {
		return err
	}
	for _, v := range many {
		tmp := []string{
			strconv.FormatInt(v.ID, 10), v.MemberUserID, strconv.Itoa(int(v.Status)), strconv.FormatInt(v.CreatedAt, 10),
			strconv.FormatInt(v.RemindTime, 10), strconv.FormatInt(v.UpdatedAt, 10),
		}
		data = append(data, tmp)
	}

	filePrefix := "sow_task"
	err = job.createCsv(c, enterpriseId, filePrefix, data)
	return err
}

// 生成csv文件
func (job *TaskData) createCsv(c ctx.Context, enterpriseId int64, filePrefix string, data [][]string) error {
	// var cosUrl string
	localFilePath := `./storage/dashboard/`
	localFileName := filePrefix + "_" + helper.GetFormatDate(0, 2) + ".csv"
	filePath := localFilePath + localFileName
	f, err := os.Create(filePath) // 创建文件
	if err != nil {
		log.Info(c).Msgf("%s::csv文件创建失败: %v", job.name, err)
		return err
	}
	_, err = f.WriteString("\xEF\xBB\xBF")
	if err != nil {
		return err
	} // 写入UTF-8 BOM
	w := csv.NewWriter(f) // 创建一个新写入文件流
	err = w.WriteAll(data)
	w.Flush()
	err = job.uploadFile(job.sftpClient, filePath, job.remotePath)

	planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
		EnterpriseName: "coach",
		SFTPName:       "KUNXIN对接" + localFileName,
		UploadRes:      helper.RuleThree(err != nil, "失败", "成功"),
	})

	return err
}

func (job *TaskData) uploadFile(sftp *sftp.SFTP, localFilePath string, remotePath string) error {
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

func (job *TaskData) CreateDir(timeStr string) (string, error) {
	dirPath := "TRANSFERIN/" + timeStr
	err := job.sftpClient.Client.Mkdir(dirPath) // 创建目录
	if err != nil {
		return "", err
	}
	err = job.sftpClient.Client.Chmod(dirPath, 0777)
	realPath, err := job.sftpClient.Client.RealPath(dirPath)
	if err != nil {
		return "", err
	}
	return realPath + "/", nil
}

func (job *TaskData) createTaskSalesData(c ctx.Context, enterpriseId int64) error {
	data := [][]string{
		{"plan_id", "task_id", "union_id", "user_id", "store_code", "order_price", "order_no", "origin_order_no", "order_pay_time", "updated_at", "order_status", "created_at"},
	}
	todayStart, _ := helper.PeriodOfTimestamp()
	startTime := todayStart - 86400
	endTime := time.Now().Unix()
	many, err := planDao.NewTaskConditionDataSaleDao().GetMany(c, enterpriseId, "updated_at>=? and updated_at<=? and order_status in (1,2) ", startTime, endTime)
	if err != nil {
		return err
	}
	for _, v := range many {
		tmp := []string{
			strconv.FormatInt(v.PlanID, 10), strconv.FormatInt(v.TaskID, 10), v.UnionID,
			v.UserID, v.StoreCode, strconv.Itoa(v.OrderPrice), v.OrderNo, v.OriginOrderNo, strconv.FormatInt(int64(v.OrderPayTime), 10),
			strconv.FormatInt(v.UpdatedAt, 10), strconv.FormatInt(int64(v.OrderStatus), 10), strconv.FormatInt(v.CreatedAt, 10),
		}
		data = append(data, tmp)
	}

	filePrefix := "task_condition_data_sale"
	err = job.createCsv(c, enterpriseId, filePrefix, data)
	return err
}

func (job *TaskData) createReserveSalesData(c ctx.Context, enterpriseId int64) error {
	data := [][]string{
		{"reserve_number", "user_id", "unionid", "store_code", "order_no", "origin_order_no", "order_price", "order_pay_time", "status", "updated_at", "created_at"},
	}
	todayStart, _ := helper.PeriodOfTimestamp()
	startTime := todayStart - 86400
	endTime := time.Now().Unix()
	sales, err := planDao.NewReservesOrderSales().GetReserveOrderSales(c, enterpriseId, "updated_at>=? and updated_at<=? and order_type in (1,2)", startTime, endTime)
	if err != nil {
		return err
	}
	for _, v := range sales {
		tmp := []string{
			v.ReserveNumber, v.MemberUserID, v.UserUnionid, v.StoreCode, v.RelationOrderID, v.OriginOrderID, strconv.FormatInt(int64(v.OrderPrice), 10),
			strconv.FormatInt(int64(v.PayTime), 10), strconv.FormatInt(int64(v.OrderType), 10), strconv.FormatInt(v.UpdatedAt, 10), strconv.FormatInt(v.CreatedAt, 10),
		}
		data = append(data, tmp)
	}
	filePrefix := "reserve_order_sale"
	err = job.createCsv(c, enterpriseId, filePrefix, data)
	return err
}
