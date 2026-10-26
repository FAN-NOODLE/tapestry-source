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

type InfoMergeMemberDao struct {
	dao       dao.Dao
	model     *microScrmModel.NotificationInfoMergeMemberModel
	tableName string
}

func NewInfoMergeMemberDao() *InfoMergeMemberDao {
	return &InfoMergeMemberDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		model:     new(microScrmModel.NotificationInfoMergeMemberModel),
		tableName: microScrmModel.NotificationInfoMergeMemberModel{}.TableName(),
	}
}

func (d InfoMergeMemberDao) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.NotificationInfoMergeMemberModel, error) {
	var result microScrmModel.NotificationInfoMergeMemberModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d InfoMergeMemberDao) GetListNum(c app.Ctx, enterpriseId int64, mergeIds []int64) (map[int64]notificationResp.InfoMergeMemberNumResp, error) {
	var list []notificationResp.InfoMergeMemberNumResp
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).
		Where("merge_id in ?", mergeIds).
		Select("merge_id,count(1) as send_num,sum(view_num) as view_num, count(IF(view_num>0,true,null)) as view_people_num").
		Group("merge_id").
		Find(&list).Error
	var numMap = make(map[int64]notificationResp.InfoMergeMemberNumResp)
	for _, info := range list {
		numMap[info.MergeId] = info
	}
	return numMap, err
}

func (d InfoMergeMemberDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.NotificationInfoMergeMemberModel, error) {
	var result []microScrmModel.NotificationInfoMergeMemberModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d InfoMergeMemberDao) GetMemberIds(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]int64, error) {
	var result []int64
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Pluck("member_id", &result).Error
	return result, err
}

func (d InfoMergeMemberDao) CreateData(c app.Ctx, enterpriseId int64, data microScrmModel.NotificationInfoMergeMemberModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Create(&data).Error
	return err
}

func (d InfoMergeMemberDao) CreateInBatchesData(c app.Ctx, enterpriseId int64, data []microScrmModel.NotificationInfoMergeMemberModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d InfoMergeMemberDao) UpdateData(c app.Ctx, enterpriseId int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d InfoMergeMemberDao) DeleteData(c app.Ctx, enterpriseId int64, query interface{}) error {
	var model microScrmModel.NotificationInfoMergeMemberModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).Delete(&model).Error
	return err
}
func (d InfoMergeMemberDao) UnscopedDeleteData(c app.Ctx, enterpriseId int64, query interface{}) error {
	var model microScrmModel.NotificationInfoMergeMemberModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).Unscoped().Delete(&model).Error
	return err
}
func (d InfoMergeMemberDao) Increment(c app.Ctx, enterpriseId int64, query interface{}, column string, num int64) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).UpdateColumn(column, gorm.Expr(fmt.Sprintf("%s + ?", column), num)).Error
	return err
}
func (d InfoMergeMemberDao) GetCount(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (total int64, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Count(&total).Error
	return total, err
}
