package templateDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	templateModels "micro-scrm/internal/app/models/template"
	"micro-scrm/internal/pkg/db/gormx"
)

type Template struct {
	dao       dao.Dao
	tableName string
}

func NewTemplate() *Template {
	return &Template{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: templateModels.TemplateModel{}.TableName(),
	}
}

// WithTx 切换db连接,用于事物
func (d *Template) WithTx(tx *gorm.DB) *Template {
	d.dao.DB = tx
	return d
}

// GetLists 获取模板列表
func (d *Template) GetLists(c app.Ctx, templateName string, enterpriseID int64, offset, limit int) (data []templateModels.TemplateModel, err error) {
	buildQuery := d.dao.WithContext(c).Table(d.tableName).Where("enterprise_id=?", enterpriseID)
	if templateName != "" {
		buildQuery.Where("template_name like ? ", "%"+templateName+"%")
	}
	err = buildQuery.Where("deleted_at=0").Order("id desc").Offset(offset).Limit(limit).Find(&data).Error
	return
}

func (d *Template) CreateInfo(c app.Ctx, info templateModels.TemplateModel) error {
	return d.dao.WithContext(c).Table(d.tableName).Create(&info).Error
}

func (d *Template) UpdateInfo(c app.Ctx, data templateModels.TemplateModel, id int) error {
	return d.dao.WithContext(c).Table(d.tableName).Where(id).Updates(&data).Error
}

func (d *Template) GetCount(c app.Ctx, templateName string, enterpriseID int64) (Total int64, err error) {
	buildQuery := d.dao.WithContext(c).Table(d.tableName).Where("enterprise_id=?", enterpriseID)
	if templateName != "" {
		buildQuery.Where("template_name like ? ", "%"+templateName+"%")
	}
	err = buildQuery.Where("deleted_at=0").Count(&Total).Error
	return
}

func (d *Template) Detail(c app.Ctx, id int) (data templateModels.TemplateModel, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).Where(id).First(&data).Error
	return
}

func (d *Template) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Table(d.tableName).Where(query).UpdateColumns(data).Error
	return err
}
