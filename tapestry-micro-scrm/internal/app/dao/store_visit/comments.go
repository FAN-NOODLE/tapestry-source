package storeVisitDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	storeVisitReq "micro-scrm/internal/app/http/request/store_visit"
	storeVisitModel "micro-scrm/internal/app/models/store_visit"
)

type CommentDao struct {
	dao.EnterpriseRepo[storeVisitModel.StoreVisitCommentModel]
	model     *storeVisitModel.StoreVisitCommentModel
	tableName string
}

func NewCommentDao() *CommentDao {
	return &CommentDao{
		EnterpriseRepo: dao.EnterpriseRepo[storeVisitModel.StoreVisitCommentModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(storeVisitModel.StoreVisitCommentModel),
		tableName: storeVisitModel.StoreVisitCommentModel{}.TableName(),
	}
}

func (d *CommentDao) GetListPreloadMedia(c app.Ctx, req storeVisitReq.CommentListReq) ([]storeVisitModel.StoreVisitCommentModel, int64, error) {
	var (
		result []storeVisitModel.StoreVisitCommentModel
		total  int64
	)
	offset := (req.Page - 1) * req.Size
	err := d.WithContext(c).Scopes(d.BranchDataBase(req.EnterpriseId, d.model.TableName())).
		Model(d.model).
		Preload("MediaList", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(d.BranchTable(storeVisitModel.StoreVisitMediaModel{}.TableName(), req.EnterpriseId)).
				Where("rel_type = ? AND deleted_at = 0", storeVisitModel.Comments)
		}).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(d.BranchTable(storeVisitModel.StoreVisitCommentModel{}.TableName(), req.EnterpriseId)).
				Order("created_at desc").Preload("MediaList", func(db *gorm.DB) *gorm.DB {
				return db.Scopes(d.BranchTable(storeVisitModel.StoreVisitMediaModel{}.TableName(), req.EnterpriseId)).
					Where("rel_type = ? AND deleted_at = 0", storeVisitModel.Comments)
			})
		}).
		Where("visit_id = ? AND parent_id = 0", req.ID).Count(&total).
		Offset(offset).Limit(req.Size).Order("created_at desc").Find(&result).Error
	return result, total, err
}
