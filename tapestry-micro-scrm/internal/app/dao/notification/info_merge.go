package notificationDao

import (
	"errors"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	notificationReq "micro-scrm/internal/app/http/request/notification"
	notificationResp "micro-scrm/internal/app/http/response/notification"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"strings"
)

type InfoMergeDao struct {
	dao       dao.Dao
	model     *microScrmModel.NotificationInfoMergeModel
	tableName string
}

func NewInfoMergeDao() *InfoMergeDao {
	return &InfoMergeDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		model:     new(microScrmModel.NotificationInfoMergeModel),
		tableName: microScrmModel.NotificationInfoMergeModel{}.TableName(),
	}
}

func (d InfoMergeDao) GetOneData(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.NotificationInfoMergeModel, error) {
	var result microScrmModel.NotificationInfoMergeModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Preload("ContentList", func(db *gorm.DB) *gorm.DB {
		return db.Scopes(d.dao.BranchTable(microScrmModel.NotificationInfoMergeContentModel{}.TableName(), enterpriseId)).Order("id asc")
	}).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}
func (d InfoMergeDao) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.NotificationInfoMergeModel, error) {
	var result microScrmModel.NotificationInfoMergeModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d InfoMergeDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.NotificationInfoMergeModel, error) {
	var result []microScrmModel.NotificationInfoMergeModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d InfoMergeDao) GetPagination(c app.Ctx, cond notificationReq.InfoMergeListReq) ([]microScrmModel.NotificationInfoMergeModel, int64, error) {
	var (
		total int64
		resp  []microScrmModel.NotificationInfoMergeModel
		err   error
	)

	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, cond.EnterpriseId)).Model(d.model).
		Scopes(func(db *gorm.DB) *gorm.DB {
			// 规则名称
			if cond.Type > 0 {
				db.Where("type = ?", cond.Type)
			}
			if cond.Title != "" {
				db.Where("title like ?", "%"+strings.Trim(cond.Title, " ")+"%")
			}
			if cond.MemberName != "" {
				db.Where("member_name like ?", "%"+strings.Trim(cond.MemberName, " ")+"%")
			}
			// 获得时间
			if cond.StartTime > 0 && cond.EndTime > 0 {
				db.Where("updated_at between ? and ?", cond.StartTime, cond.EndTime)
			}
			return db
		})
	query.Count(&total)

	offset := (cond.Page - 1) * cond.Size

	err = query.Order("updated_at desc").Limit(cond.Size).Offset(offset).Find(&resp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return resp, total, err
}

func (d InfoMergeDao) CreateData(c app.Ctx, enterpriseId int64, data microScrmModel.NotificationInfoMergeModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Create(&data).Error
	return err
}

func (d InfoMergeDao) CreateInBatchesData(c app.Ctx, enterpriseId int64, data []microScrmModel.NotificationInfoMergeModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d InfoMergeDao) UpdateData(c app.Ctx, enterpriseId int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d InfoMergeDao) DeleteData(c app.Ctx, enterpriseId int64, query interface{}) error {
	var model microScrmModel.NotificationInfoMergeModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).Delete(&model).Error
	return err
}

func (d InfoMergeDao) GetShowList(c app.Ctx, cond notificationReq.InfoMergeListReq) ([]notificationResp.InfoMergeListByH5Resp, int64, error) {
	var (
		total int64
		resp  []notificationResp.InfoMergeListByH5Resp
		err   error
	)
	query := d.dao.WithContext(c).Table(d.tableName).Scopes(d.dao.BranchTable(d.tableName, cond.EnterpriseId)).Model(d.model).
		Joins("left join "+d.dao.GetTable(microScrmModel.NotificationInfoMergeMemberModel{}.TableName(), cond.EnterpriseId)+" ON notification_info_merge_member.merge_id = notification_info_merge.id").
		Select("notification_info_merge.id,notification_info_merge.type,notification_info_merge.title,notification_info_merge.cover,notification_info_merge.send_time,notification_info_merge.send_ro_name as ro_name,notification_info_merge.pdf_url,notification_info_merge_member.view_num").
		Where("notification_info_merge_member.member_id = ?", cond.MemberId)
	if cond.Type > 0 {
		query = query.Where("notification_info_merge.type = ?", cond.Type)
	}
	if cond.Title != "" {
		query = query.Where("notification_info_merge.title like ?", "%"+cond.Title+"%")
	}
	query.Count(&total)

	offset := (cond.Page - 1) * cond.Size

	err = query.Order("notification_info_merge_member.created_at DESC").Limit(cond.Size).Offset(offset).Find(&resp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return resp, total, err
}
