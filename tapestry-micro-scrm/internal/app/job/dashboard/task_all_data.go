package dashboardJob

import (
	"encoding/csv"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	customerDao "micro-scrm/internal/app/dao/customer"
	dashboardDao "micro-scrm/internal/app/dao/dashboard"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	appJob "micro-scrm/internal/app/job"
	"micro-scrm/internal/app/models"
	dashboardModel "micro-scrm/internal/app/models/dashboard"
	cacheSvc "micro-scrm/internal/app/service/cache"
	planService "micro-scrm/internal/app/service/plan"
	"micro-scrm/internal/modules/marketingPromotion/storage"
	"micro-scrm/internal/modules/marketingPromotion/storage/cos"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/xjob"
	"os"
	"strconv"
	"time"
)

type TaskAllData struct {
	name      string
	startTime int64
	endTime   int64
}

func NewTaskAllData() *TaskAllData {
	startTime, endTime := int64(0), int64(0) // 全量数据，设置startTime的值为0，增量数据取昨天的时间戳
	return &TaskAllData{
		name:      "任务数据dashboard",
		startTime: startTime,
		endTime:   endTime,
	}
}

func (job *TaskAllData) Run(c ctx.Context, handler *xjob.JobHandler) error {
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
func (job *TaskAllData) Handle(c ctx.Context, enterpriseId int64) error {
	config, err := cacheSvc.NewRedisData().GetEnterpriseConfig(c, enterpriseId, "store_name")
	if err != nil {
		return err
	}
	if config.Name != "coach" {
		return nil
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

	return err
}

// 主任务表数据
func (job *TaskAllData) createPlanData(c ctx.Context, enterpriseId int64) error {
	// 生成主任务表的数据dashboard

	many, err := planDao.NewPlanTaskDao().GetPlanMany(c, enterpriseId, "updated_at>? ", job.startTime)

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
func (job *TaskAllData) createPlanTaskData(c ctx.Context, enterpriseId int64) error {
	// 生成子任务表的数据dashboard
	many, err := planDao.NewPlanTaskDao().GetMany(c, enterpriseId, "updated_at>? ", job.startTime)
	if err != nil {
		return err
	}
	data := [][]string{
		{"id", "plan_id", "task_category", "name", "remark", "task_type", "plan_send_time", "updated_at"},
	}
	for _, v := range many {
		tmp := []string{
			strconv.FormatInt(v.ID, 10), strconv.FormatInt(v.PlanID, 10), strconv.Itoa(v.TaskCategory), v.Name, v.Remark, strconv.Itoa(int(v.TaskType)), strconv.FormatInt(v.PlanSendTime, 10), strconv.FormatInt(v.UpdatedAt, 10),
		}
		data = append(data, tmp)
	}
	filePrefix := "plan_tasks"
	err = job.createCsv(c, enterpriseId, filePrefix, data)
	return err
}

// 任务导购表数据
func (job *TaskAllData) createTaskAllotLogData(c ctx.Context, enterpriseId int64) error {
	// 生成任务导购task_allot_log表的数据dashboard
	data := [][]string{
		{"plan_id", "task_id", "member_user_id", "end_time", "is_finished", "customer_num", "updated_at", "contacted_count", "moment_look_count", "no_contacted_count", "current_contacted_count"},
	}

	id := int64(0)
	for {
		many, err := planDao.NewPlanTaskAllotLogDao().GetManyByIdPagination(c, enterpriseId, "id>?", id)
		if err != nil {
			return err
		}
		index := len(many)
		if index == 0 {
			break
		}
		id = many[index-1].ID
		for _, v := range many {
			customerNum := strconv.Itoa(v.CustomerNum)
			contactedCount := strconv.Itoa(v.ContactedCount)
			noContactedCount := v.CustomerNum - v.ContactedCount
			momentLookCount := v.MomentLookCount
			tmp := []string{
				strconv.FormatInt(v.PlanId, 10), strconv.FormatInt(v.TaskId, 10), v.UserId, strconv.Itoa(v.EndTime), strconv.Itoa(v.IsFinished), customerNum, strconv.FormatInt(v.UpdatedAt, 10), contactedCount, strconv.Itoa(momentLookCount), strconv.Itoa(noContactedCount), "0",
			}
			data = append(data, tmp)
		}
	}

	filePrefix := "task_allot_log"
	err := job.createCsv(c, enterpriseId, filePrefix, data)
	return err
}

// 邀约记录表数据
func (job *TaskAllData) createReserveData(c ctx.Context, enterpriseId int64) error {
	data := [][]string{
		{"reserve_number", "reserve_id", "reserve_time", "store_id", "status", "member_user_id", "reserve_order_price", "reserve_order_time", "is_store", "is_confirm", "updated_at", "created_at"},
	}
	id := int64(0)
	for {
		many, err := planDao.NewReservesOrder().GetManyByIdPagination(c, enterpriseId, "id > ?", id)
		if err != nil {
			return err
		}
		index := len(many)
		if index == 0 {
			break
		}
		id = many[index-1].ID
		for _, v := range many {
			tmp := []string{
				v.ReserveNumber, strconv.FormatInt(v.ReserveID, 10), strconv.Itoa(v.ReserveTime),
				strconv.FormatInt(v.StoreID, 10), strconv.Itoa(int(v.Status)), v.MemberUserID,
				strconv.Itoa(v.ReserveOrderPrice), strconv.Itoa(v.ReserveOrderTime), strconv.Itoa(int(v.IsStore)),
				strconv.Itoa(int(v.IsConfirm)), strconv.FormatInt(v.UpdatedAt, 10), strconv.FormatInt(v.CreatedAt, 10),
			}
			data = append(data, tmp)
		}
	}
	filePrefix := "reserve_order"
	err := job.createCsv(c, enterpriseId, filePrefix, data)
	return err
}

// 邀约服务表数据
func (job *TaskAllData) createReserveServiceData(c ctx.Context, enterpriseId int64) error {

	many, err := planDao.NewReserveServiceDao().GetMany(c, enterpriseId, "updated_at>=? ", job.startTime)
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
func (job *TaskAllData) createSowTaskData(c ctx.Context, enterpriseId int64) error {
	data := [][]string{
		{"id", "member_user_id", "status", "created_at", "remind_time", "updated_at"},
	}

	id := int64(0)
	for {
		many, err := customerDao.NewSowTask().GetManyByIdPagination(c, "enterprise_id=? and id>?", enterpriseId, id)
		if err != nil {
			return err
		}
		index := len(many)
		if index == 0 {
			break
		}
		id = many[index-1].ID
		for _, v := range many {
			tmp := []string{
				strconv.FormatInt(v.ID, 10), v.MemberUserID, strconv.Itoa(int(v.Status)), strconv.FormatInt(v.CreatedAt, 10),
				strconv.FormatInt(v.RemindTime, 10), strconv.FormatInt(v.UpdatedAt, 10),
			}
			data = append(data, tmp)
		}
	}

	filePrefix := "sow_task"
	err := job.createCsv(c, enterpriseId, filePrefix, data)
	return err
}

// 生成csv文件
func (job *TaskAllData) createCsv(c ctx.Context, enterpriseId int64, filePrefix string, data [][]string) error {
	var cosUrl string
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
	// 记录上传日志
	cs := storage.Init(cos.COS)
	uploadFile := fmt.Sprintf("uploads/dashboard/%s/%s", time.Now().Format("2006/01/02"), localFileName)
	path, err := cs.Upload(c, filePath, uploadFile)
	cosUrl = cs.GetURL(path)
	nowTime := time.Now().Unix()
	logData := dashboardModel.DashboardLogModel{
		MicroBaseModel: models.MicroBaseModel{
			BaseModel: models.BaseModel{
				ID:        app.UID(),
				CreatedAt: nowTime,
				UpdatedAt: nowTime,
			},
			EnterpriseID: enterpriseId,
		},
		FileType: filePrefix,
		FileName: localFileName,
		FileURL:  cosUrl,
	}
	err = dashboardDao.NewLog().CreateData(c, logData)
	// 上传异常通知
	num := len(data)
	planService.NewTaskNotice().NfsUploadNotice(c, enterpriseId, localFileName, cosUrl, num)
	return err
}
