package drainageContactDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	drainageContactGroupReq "micro-scrm/internal/app/http/request/drainage_contact"
	"micro-scrm/internal/app/models/drainage"
)

type DrainageContactGroupDao struct {
	dao.EnterpriseRepo[drainageContactModel.DrainageContactGroupModel]
	tableName string
	model     *drainageContactModel.DrainageContactGroupModel
}

func NewDrainageContactGroupDao() *DrainageContactGroupDao {
	return &DrainageContactGroupDao{
		EnterpriseRepo: dao.EnterpriseRepo[drainageContactModel.DrainageContactGroupModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		tableName: drainageContactModel.DrainageContactGroupModel{}.TableName(),
		model:     new(drainageContactModel.DrainageContactGroupModel),
	}
}

func (d *DrainageContactGroupDao) GetOnePreloadContactAndMember(c app.Ctx, enterpriseId, ID int64) (drainageContactModel.DrainageContactGroupModel, error) {
	var result drainageContactModel.DrainageContactGroupModel
	err := d.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, enterpriseId)).
		Where("id = ?", ID).Preload("Contacts.Members").Find(&result).Error
	return result, err
}

func (d *DrainageContactGroupDao) GetOnePreloadContacts(c app.Ctx, enterpriseId, ID int64) (drainageContactModel.DrainageContactGroupModel, error) {
	var result drainageContactModel.DrainageContactGroupModel
	err := d.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, enterpriseId)).
		Where("id = ?", ID).Preload("Contacts").Find(&result).Error
	return result, err
}

func (d *DrainageContactGroupDao) GetOnePreloadWelcomes(c app.Ctx, enterpriseId, ID int64) (drainageContactModel.DrainageContactGroupModel, error) {
	var result drainageContactModel.DrainageContactGroupModel
	err := d.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, enterpriseId)).
		Where("id = ?", ID).
		Preload("Welcomes", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(d.BranchTable(drainageContactModel.DrainageContactsWelcomesModel{}.TableName(), enterpriseId)).
				Where("deleted_at = 0")
		}).
		Find(&result).Error
	return result, err
}

func (d *DrainageContactGroupDao) GetOneBySourceType(c app.Ctx, enterpriseId int64, sourceType int) (drainageContactModel.DrainageContactGroupModel, error) {
	var result drainageContactModel.DrainageContactGroupModel
	err := d.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, enterpriseId)).
		Where("source_type = ? AND deleted_at = 0", sourceType).
		Order("id asc").
		Limit(1).
		Find(&result).Error
	return result, err
}

func (d *DrainageContactGroupDao) GetPage(c ctx.Context, req drainageContactGroupReq.GroupListReq) ([]drainageContactModel.DrainageContactGroupModel, int64, error) {
	var (
		result []drainageContactModel.DrainageContactGroupModel
		total  int64
	)
	query := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, req.EnterpriseId)).Model(d.model)
	query = query.Where("source_type = ?", drainageContactModel.DrainageContactGroupSourceConsole)
	if req.ChannelId > 0 {
		query = query.Where("channel_id =  ? ", req.ChannelId)
	}
	if req.Type > 0 {
		query = query.Where("type =  ? ", req.Type)
	}
	if req.CreateTime.StartTime > 0 {
		query = query.Where("created_at >= ? ", req.CreateTime.StartTime)
	}
	if req.CreateTime.EndTime > 0 {
		query = query.Where("created_at <= ? ", req.CreateTime.EndTime)
	}
	if req.Name != "" {
		query = query.Where("name like ? ", "%"+req.Name+"%")
	}
	offset := (req.Page - 1) * req.Size
	err := query.Count(&total).
		Order("created_at desc").
		Preload("Contacts", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at = 0")
		}).
		Preload("Welcomes", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(d.BranchTable(drainageContactModel.DrainageContactsWelcomesModel{}.TableName(), req.EnterpriseId)).
				Where("deleted_at = 0")
		}).
		Offset(offset).Limit(req.Size).
		Find(&result).Error
	return result, total, err
}
