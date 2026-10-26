package notificationDao

import (
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	notificationResp "micro-scrm/internal/app/http/response/notification"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type InfoMemberDao struct {
	dao       dao.Dao
	model     *microScrmModel.NotificationInfoMemberModel
	tableName string
}

func NewInfoMemberDao() *InfoMemberDao {
	return &InfoMemberDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		model:     new(microScrmModel.NotificationInfoMemberModel),
		tableName: microScrmModel.NotificationInfoMemberModel{}.TableName(),
	}
}

func (d InfoMemberDao) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.NotificationInfoMemberModel, error) {
	var result microScrmModel.NotificationInfoMemberModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d InfoMemberDao) GetInfoReadNum(c app.Ctx, enterpriseId int64, infoId int64) (notificationResp.InfoViewNumResp, error) {
	var readInfo notificationResp.InfoViewNumResp
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).
		Where("info_id = ?", infoId).
		Select("info_id,count(1) as view_people_num,sum(view_num) as view_num").
		Group("info_id").
		Find(&readInfo).Error

	return readInfo, err
}

func (d InfoMemberDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.NotificationInfoMemberModel, error) {
	var result []microScrmModel.NotificationInfoMemberModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d InfoMemberDao) GetMemberIds(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]int64, error) {
	var result []int64
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Pluck("member_id", &result).Error
	return result, err
}

func (d InfoMemberDao) CreateData(c app.Ctx, enterpriseId int64, data microScrmModel.NotificationInfoMemberModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Create(&data).Error
	return err
}

func (d InfoMemberDao) CreateInBatchesData(c app.Ctx, enterpriseId int64, data []microScrmModel.NotificationInfoMemberModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d InfoMemberDao) UpdateData(c app.Ctx, enterpriseId int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d InfoMemberDao) DeleteData(c app.Ctx, enterpriseId int64, query interface{}) error {
	var model microScrmModel.NotificationInfoMemberModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).Delete(&model).Error
	return err
}
func (d InfoMemberDao) UnscopedDeleteData(c app.Ctx, enterpriseId int64, query interface{}) error {
	var model microScrmModel.NotificationInfoMemberModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).Unscoped().Delete(&model).Error
	return err
}
func (d InfoMemberDao) Increment(c app.Ctx, enterpriseId int64, query interface{}, column string, num int64) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).UpdateColumn(column, gorm.Expr(fmt.Sprintf("%s + ?", column), num)).Error
	return err
}
func (d InfoMemberDao) GetCount(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (total int64, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Count(&total).Error
	return total, err
}
