package customerDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	customerReq "micro-scrm/internal/app/http/request/customer"
	planReq "micro-scrm/internal/app/http/request/plan"
	customerResp "micro-scrm/internal/app/http/response/customer"
	"micro-scrm/internal/app/models"
	"micro-scrm/internal/app/models/customer"
	"micro-scrm/internal/pkg/db/gormx"
)

type SowTask struct {
	dao   dao.Dao
	model *customerModels.SowTaskModel
}

func NewSowTask() *SowTask {
	return &SowTask{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(customerModels.SowTaskModel),
	}
}

func (d *SowTask) GetOne(c app.Ctx, query interface{}, args ...interface{}) (customerModels.SowTaskModel, error) {
	var result customerModels.SowTaskModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *SowTask) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]customerModels.SowTaskModel, error) {
	var result []customerModels.SowTaskModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *SowTask) GetList(c app.Ctx, page, size int, cond interface{}, option map[string]interface{}) ([]customerModels.SowTaskModel, int64, error) {
	var result []customerModels.SowTaskModel
	var total int64
	offset := (page - 1) * size
	query := d.dao.WithContext(c).Model(d.model).Where(cond)
	if _, ok := option["start_time"]; ok {
		query.Where("remind_time BETWEEN ? AND ?", option["start_time"], option["end_time"])
	}
	err := query.Count(&total).Offset(offset).Limit(size).
		Order("status asc").
		Order("remind_time desc").
		Find(&result).Error

	return result, total, err
}

func (d *SowTask) CreateData(c app.Ctx, data customerModels.SowTaskModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *SowTask) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).Updates(data).Error
	return err
}

func (d *SowTask) DeleteData(c app.Ctx, query interface{}) error {
	var model customerModels.SowTaskModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Delete(&model).Error
	return err
}

// GetMonthLists 获取每个月的列表
func (d *SowTask) GetMonthLists(c ctx.Context, request planReq.MonthListsReq, memberUserId string) ([]customerModels.SowTaskModel, error) {
	var result []customerModels.SowTaskModel
	query := d.dao.WithContext(c).Model(d.model).
		Where("enterprise_id", request.EnterpriseId).
		Where("member_user_id", memberUserId).
		Where("remind_time BETWEEN ? AND ?", request.StartTime, request.EndTime)
	err := query.Find(&result).Error
	return result, err
}

// GetPlanTipsCount 获取任务中心的指定时间内的未完成播种数
func (d *SowTask) GetPlanTipsCount(c ctx.Context, enterpriseId, startTime, endTime int64, memberId int64) (total int64, err error) {
	err = d.dao.WithContext(c).Model(d.model).
		Where("enterprise_id", enterpriseId).
		Where("member_id = ? and status = 1", memberId).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if startTime > 0 && endTime > 0 {
				db.Where("remind_time BETWEEN ? AND ?", startTime, endTime)
			}
			return db
		}).Count(&total).Error
	return
}

func (d *SowTask) GetSowTaskTipsCount(c ctx.Context, enterpriseId, workExternalUserId int64, memberUserIds []int64) (total int64, err error) {
	err = d.dao.WithContext(c).Model(d.model).
		Where("enterprise_id=?", enterpriseId).
		Where("member_id in  (?) and deleted_at = 0", memberUserIds).
		Where("work_external_user_id = ?", workExternalUserId).
		Count(&total).Error

	return
}

// GetSowTaskCountByMemberIds  获取导购指定时间内的未完成播种数
func (d *SowTask) GetSowTaskCountByMemberIds(c ctx.Context, enterpriseId, startTime, endTime int64, memberIds []int64) (total int64, err error) {
	err = d.dao.WithContext(c).Model(d.model).
		Where("enterprise_id", enterpriseId).
		Where("member_id in ? and status = 1", memberIds).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if startTime > 0 && endTime > 0 {
				db.Where("remind_time BETWEEN ? AND ?", startTime, endTime)
			}
			return db
		}).Count(&total).Error
	return
}

func (d *SowTask) GetTopList(c app.Ctx, enterpriseId int64, inviteTop int, startTime, endTime int64) ([]models.TopList, error) {
	var result []models.TopList
	var lastTopTotal int64
	orm := d.dao.WithContext(c).Model(d.model).
		Select("COUNT(1) AS total").
		Where("enterprise_id = ?", enterpriseId).
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Group("member_id").
		Order("total desc")
	// 获取top以及和top最后一位数量相同的其他导购
	newOrm := orm.Session(&gorm.Session{NewDB: true}).WithContext(c)
	err := orm.Offset(inviteTop-1).Limit(1).Pluck("total", &lastTopTotal).Error
	err = newOrm.Select("member_id,COUNT(1) AS total").
		Having("total >= ?", lastTopTotal).Scan(&result).Error
	return result, err
}

// GetManyByIdPagination 根据id分页查询
func (d *SowTask) GetManyByIdPagination(c app.Ctx, query interface{}, args ...interface{}) ([]customerModels.SowTaskModel, error) {
	var result []customerModels.SowTaskModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Limit(1000).Find(&result).Error
	return result, err
}

func (d *SowTask) GetMemberIdNumByTemplateIds(c app.Ctx, enterpriseId int64, templateIds []int64) ([]customerResp.SowSaNum, error) {
	var result []customerResp.SowSaNum
	err := d.dao.WithContext(c).Model(d.model).
		Select("template_id,count(DISTINCT member_id) as sa_num").
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"template_id":   templateIds,
		}).Group("template_id").Find(&result).Error
	return result, err
}

func (d *SowTask) GetBoardList(c app.Ctx, params customerReq.SowAllTaskListReq) ([]customerModels.SowTaskModel, int64, error) {
	var result []customerModels.SowTaskModel
	var total int64
	offset := (params.Page - 1) * params.Size
	query := d.dao.WithContext(c).Model(d.model).Where("enterprise_id=?", params.EnterpriseID).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if params.TemplateId > 0 {
				db.Where("template_id = ?", params.TemplateId)
			}
			if params.TemplateType == customerReq.SowTaskListTemplateCustom {
				db.Where("template_id = 0")
			} else if params.TemplateType == customerReq.SowTaskListTemplateTemplate {
				db.Where("template_id > 0")
			}
			if params.SowStatus > 0 {
				db.Where("status = ?", params.SowStatus)
			}
			if len(params.TemplateIds) > 0 {
				db.Where("template_id in ?", params.TemplateIds)
			}
			if len(params.ExternalIds) > 0 {
				db.Where("work_external_user_id in ?", params.ExternalIds)
			}
			if len(params.MemberIds) > 0 {
				db.Where("member_id in ?", params.MemberIds)
			}
			if params.StartCreatedAt > 0 {
				db.Where("created_at > ?", params.StartCreatedAt)
			}
			if params.EndCreatedAt > 0 {
				db.Where("created_at < ?", params.EndCreatedAt)
			}
			if params.StartNoticeAt > 0 {
				db.Where("remind_time > ?", params.StartNoticeAt)
			}
			if params.EndNoticeAt > 0 {
				db.Where("remind_time < ?", params.EndNoticeAt)
			}
			if params.StartFinishedAt > 0 {
				db.Where("status=? and updated_at > ?", customerModels.SowTaskStatusFinished, params.StartFinishedAt)
			}
			if params.EndFinishedAt > 0 {
				db.Where("updated_at < ?", params.EndFinishedAt)
			}
			return db
		})
	err := query.Count(&total).Offset(offset).Limit(params.Size).
		Order("created_at desc").
		Find(&result).Error

	return result, total, err
}

func (d *SowTask) CreateManyData(c app.Ctx, data []customerModels.SowTaskModel) error {
	err := d.dao.WithContext(c).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}
