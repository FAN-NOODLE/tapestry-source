package dao

import (
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	membersModel "micro-scrm/internal/app/models/members"
	"micro-scrm/internal/pkg/db/gormx"
)

type Demo struct {
	dao   Dao
	model *membersModel.MemberModel
}

func NewDemo() *Demo {
	return &Demo{
		dao:   Dao{DB: gormx.GetInstance()},
		model: new(membersModel.MemberModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *Demo) WithTx(tx *gorm.DB) *Demo {
	d.dao.DB = tx
	return d
}

func (d *Demo) GetOne(c app.Ctx, query interface{}, args ...interface{}) (membersModel.MemberModel, error) {
	var result membersModel.MemberModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *Demo) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]membersModel.MemberModel, error) {
	var result []membersModel.MemberModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *Demo) CreateData(c app.Ctx, data membersModel.MemberModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *Demo) CreateInBatchesData(c app.Ctx, data []membersModel.MemberModel) error {
	err := d.dao.WithContext(c).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *Demo) SaveData(c app.Ctx, data membersModel.MemberModel) error {
	err := d.dao.WithContext(c).Model(d.model).Save(&data).Error
	return err
}

func (d *Demo) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *Demo) DeleteData(c app.Ctx, query interface{}, args ...interface{}) error {
	var model membersModel.MemberModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Delete(&model).Error
	return err
}

func (d *Demo) ForceDelete(c app.Ctx, query interface{}, args ...interface{}) error {
	var model membersModel.MemberModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Unscoped().Delete(&model).Error
	return err
}

func (d *Demo) Increment(c app.Ctx, query interface{}, column string, num int64) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumn(column, gorm.Expr(fmt.Sprintf("%s + ?", column), num)).Error
	return err
}

func (d *Demo) Decrement(c app.Ctx, query interface{}, column string, num int64) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumn(column, gorm.Expr(fmt.Sprintf("%s - ?", column), num)).Error
	return err
}
