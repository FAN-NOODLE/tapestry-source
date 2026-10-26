package microScrmDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	openReq "micro-scrm/internal/app/http/request/open"
	planReq "micro-scrm/internal/app/http/request/plan"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
)

type PlanTemplateDao struct {
	dao       dao.Dao
	tableName string
}

func NewPlanTemplateDao() *PlanTemplateDao {
	return &PlanTemplateDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.PlanTaskTemplateModel{}.TableName(),
	}
}

func (d *PlanTemplateDao) Create(c ctx.Context, enterpriseId int64, data microScrmModel.PlanTaskTemplateModel) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&data).Error
}

func (d *PlanTemplateDao) UpdateData(c app.Ctx, enterpriseId int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query).UpdateColumns(data).Error
	return err
}

func (d *PlanTemplateDao) Lists(c ctx.Context, cond planReq.TemplateListsReq, adminId int64) ([]microScrmModel.PlanTaskTemplateModel, int64, error) {

	var (
		resp  []microScrmModel.PlanTaskTemplateModel
		total int64
		err   error
	)
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, cond.EnterpriseId)).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if cond.SearchType == 1 && cond.SearchName != "" {
				db.Where("template_name like ?", "%"+cond.SearchName+"%")
			}
			if cond.SearchType == 2 && cond.SearchId != 0 {
				db.Where("id = ?", cond.SearchId)
			}
			if cond.TaskCategory >= 0 {
				db.Where("task_category = ?", cond.TaskCategory)
			}
			if adminId > 0 {
				db.Where("admin_id = ?", adminId)
			}

			if cond.StartTime > 0 && cond.EndTime > 0 {
				db.Where("created_at between ? and ?", cond.StartTime, cond.EndTime)
			}
			return db
		})

	query.Count(&total)

	offset := (cond.Page - 1) * cond.Size
	err = query.Order("id DESC").Order("created_at DESC").Limit(cond.Size).Offset(offset).Find(&resp).Error

	return resp, total, err
}

// MaTemplateLists MA任务的三方接口获取模板列表
func (d *PlanTemplateDao) MaTemplateLists(c ctx.Context, request openReq.PlanTemplateListsReq) (result []microScrmModel.MaTemplateLists, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseId)).Where(map[string]interface{}{
		"status": 1,
	}).Select("id,template_name,template_desc,group_send_type,task_category").Find(&result).Error
	return
}

// GetTemplateInfo 获取任务模板信息
func (d *PlanTemplateDao) GetTemplateInfo(c ctx.Context, enterpriseId, id int64) (result microScrmModel.PlanTaskTemplateModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
		"id": id,
	}).Find(&result).Error
	return
}

// CreateBack 创建备份模板
func (d *PlanTemplateDao) CreateBack(c ctx.Context, enterpriseId int64, data microScrmModel.PlanTaskTemplateBackupModel) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.PlanTaskTemplateBackupModel{}.TableName(), enterpriseId)).Create(&data).Error
}

// GetTemplateBackUpInfo 获取备份模板的信息
func (d *PlanTemplateDao) GetTemplateBackUpInfo(c ctx.Context, enterpriseId, templateId int64) (result microScrmModel.PlanTaskTemplateBackupModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.PlanTaskTemplateBackupModel{}.TableName(), enterpriseId)).Where(map[string]interface{}{
		"template_id": templateId,
	}).Find(&result).Error
	return
}

func (d *PlanTemplateDao) DeleteBack(c ctx.Context, enterpriseId, id int64) error {
	var model microScrmModel.PlanTaskTemplateBackupModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.PlanTaskTemplateBackupModel{}.TableName(), enterpriseId)).Where(map[string]interface{}{
		"id": id,
	}).Delete(&model).Error
	return err
}

func (d *PlanTemplateDao) UpdateBack(c ctx.Context, enterpriseId int64, data map[string]interface{}, query interface{}, args ...interface{}) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.PlanTaskTemplateBackupModel{}.TableName(), enterpriseId)).Where(query, args...).Model(&microScrmModel.PlanTaskTemplateBackupModel{}).Updates(&data).Error
}
