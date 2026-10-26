package customerDao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	customerReq "micro-scrm/internal/app/http/request/customer"
	"micro-scrm/internal/app/models/customer"
	"micro-scrm/internal/pkg/db/gormx"
	"strings"
)

type SowTemplate struct {
	dao   dao.Dao
	model *customerModels.SowTemplateModel
}

func NewSowTemplate() *SowTemplate {
	return &SowTemplate{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(customerModels.SowTemplateModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *SowTemplate) WithTx(tx *gorm.DB) *SowTemplate {
	d.dao.DB = tx
	return d
}

func (d *SowTemplate) GetOne(c app.Ctx, query interface{}, args ...interface{}) (customerModels.SowTemplateModel, error) {
	var result customerModels.SowTemplateModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *SowTemplate) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]customerModels.SowTemplateModel, error) {
	var result []customerModels.SowTemplateModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *SowTemplate) GetUnScopedMany(c app.Ctx, query interface{}, args ...interface{}) ([]customerModels.SowTemplateModel, error) {
	var result []customerModels.SowTemplateModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Unscoped().Find(&result).Error
	return result, err
}

func (d *SowTemplate) GetPagination(c app.Ctx, cond customerReq.SowTemplateListReq) ([]customerModels.SowTemplateModel, int64, error) {
	var (
		total int64
		resp  []customerModels.SowTemplateModel
		err   error
	)

	query := d.dao.WithContext(c).Model(d.model).
		Scopes(func(db *gorm.DB) *gorm.DB {
			// 模版名称
			if cond.Name != "" {
				searchName := strings.ReplaceAll(cond.Name, "%", "\\%")
				db.Where("name like ?", "%"+strings.Trim(searchName, " ")+"%")
			}
			// 上下架
			if cond.Status > 0 {
				db.Where("status = ?", cond.Status)
			}
			// 更新时间
			if cond.StartTime > 0 && cond.EndTime > 0 {
				db.Where("updated_at between ? and ?", cond.StartTime, cond.EndTime)
			}
			return db
		})
	query.Count(&total)

	offset := (cond.Page - 1) * cond.Size

	err = query.Preload("NoticeList").Order("created_at DESC").Order("id DESC").Limit(cond.Size).Offset(offset).Find(&resp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return resp, total, err
}

func (d *SowTemplate) GetDetail(c app.Ctx, enterpriseId int64, templateId int64) (customerModels.SowTemplateModel, error) {
	var result customerModels.SowTemplateModel
	err := d.dao.WithContext(c).Model(d.model).Preload("NoticeList").Where("enterprise_id = ? AND id=?", enterpriseId, templateId).Find(&result).Error
	return result, err
}

func (d *SowTemplate) CreateData(c app.Ctx, data customerModels.SowTemplateModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *SowTemplate) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *SowTemplate) DeleteData(c app.Ctx, query interface{}, args ...interface{}) error {
	var model customerModels.SowTemplateModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Delete(&model).Error
	return err
}

func (d *SowTemplate) SearchIds(c app.Ctx, query interface{}, args ...interface{}) ([]int64, error) {
	var result []int64
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Pluck("id", &result).Error
	return result, err
}

func (d *SowTemplate) GetTemplateList(c app.Ctx, query interface{}, args ...interface{}) ([]customerModels.SowTemplateModel, error) {
	var result []customerModels.SowTemplateModel
	err := d.dao.WithContext(c).Model(d.model).Preload("NoticeList").Where(query, args...).Find(&result).Error
	return result, err
}

func (d *SowTemplate) Increment(c app.Ctx, query interface{}, column string, num int64) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumn(column, gorm.Expr(fmt.Sprintf("%s + ?", column), num)).Error
	return err
}
