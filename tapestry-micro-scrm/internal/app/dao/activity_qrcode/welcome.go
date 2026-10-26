package activityQRCodeDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	activityQRCodeReq "micro-scrm/internal/app/http/request/activity_qrcode"
	activityQRCodeModel "micro-scrm/internal/app/models/activity_qrcode"
)

type WelcomeDao struct {
	dao.EnterpriseRepo[activityQRCodeModel.ActivityQRCodeWelcomeModel]
	model     *activityQRCodeModel.ActivityQRCodeWelcomeModel
	tableName string
}

func NewWelcomeDao() *WelcomeDao {
	return &WelcomeDao{
		EnterpriseRepo: dao.EnterpriseRepo[activityQRCodeModel.ActivityQRCodeWelcomeModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(activityQRCodeModel.ActivityQRCodeWelcomeModel),
		tableName: activityQRCodeModel.ActivityQRCodeWelcomeModel{}.TableName(),
	}
}

func (d *WelcomeDao) GetWelcomeList(c ctx.Context, req activityQRCodeReq.WelcomeListReq) ([]activityQRCodeModel.ActivityQRCodeWelcomeModel, int64, error) {
	var (
		result []activityQRCodeModel.ActivityQRCodeWelcomeModel
		total  int64
	)
	query := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, req.EnterpriseId)).Model(d.model).
		Where("rel_type = ?", activityQRCodeModel.WelcomeTypeCategory)
	if req.Content != "" {
		query.Where("content like ? ", "%"+req.Content+"%")
	}
	if req.CategoryID > 0 {
		query.Where("rel_id = ?", req.CategoryID)
	}
	offset := (req.Page - 1) * req.Size
	err := query.Count(&total).
		Order("created_at desc").
		Offset(offset).Limit(req.Size).
		Find(&result).Error
	return result, total, err
}

func (d *WelcomeDao) GetAllCategoryIds(c app.Ctx, enterpriseId int64, query any, args ...any) ([]int64, error) {
	var result []int64
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).Model(d.model).
		Where(query, args...).Pluck("rel_id", &result).Error
	return result, err
}
