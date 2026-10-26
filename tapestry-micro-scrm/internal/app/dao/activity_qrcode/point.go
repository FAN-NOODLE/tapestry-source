package activityQRCodeDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	activityQRCodeReq "micro-scrm/internal/app/http/request/activity_qrcode"
	activityQRCodeModel "micro-scrm/internal/app/models/activity_qrcode"
)

type PointDao struct {
	dao.EnterpriseRepo[activityQRCodeModel.ActivityQRCodePointModel]
	model     *activityQRCodeModel.ActivityQRCodePointModel
	tableName string
}

func NewPointDao() *PointDao {
	return &PointDao{
		EnterpriseRepo: dao.EnterpriseRepo[activityQRCodeModel.ActivityQRCodePointModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(activityQRCodeModel.ActivityQRCodePointModel),
		tableName: activityQRCodeModel.ActivityQRCodePointModel{}.TableName(),
	}
}

func (d *PointDao) GetStatisticList(c ctx.Context, req activityQRCodeReq.StatisticsReq, activityIds []int64) ([]activityQRCodeModel.ActivityQRCodePointModel, []int64, error) {
	var (
		result   []activityQRCodeModel.ActivityQRCodePointModel
		pointIds []int64
	)
	query := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, req.EnterpriseId)).Model(d.model).
		Where("activity_id in ?", activityIds)
	if req.CategoryID > 0 {
		query.Where("category_id = ? ", req.CategoryID)
	}
	if req.PointNameStr != "" {
		query.Where("name like ? ", "%"+req.PointNameStr+"%")
	}
	offset := (req.Page - 1) * req.Size
	query.Session(&gorm.Session{}).Pluck("id", &pointIds)
	err := query.Order("activity_id desc,category_id asc,is_point asc").
		Offset(offset).Limit(req.Size).
		Find(&result).Error
	return result, pointIds, err
}
