package notificationDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	notificationReq "micro-scrm/internal/app/http/request/notification"
	notificationResp "micro-scrm/internal/app/http/response/notification"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type InfoDao struct {
	dao       dao.Dao
	model     *microScrmModel.NotificationInfoModel
	tableName string
}

func NewInfoDao() *InfoDao {
	return &InfoDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		model:     new(microScrmModel.NotificationInfoModel),
		tableName: microScrmModel.NotificationInfoModel{}.TableName(),
	}
}

func (d *InfoDao) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.NotificationInfoModel, error) {
	var result microScrmModel.NotificationInfoModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchDataBase(enterpriseId, d.tableName)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *InfoDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.NotificationInfoModel, error) {
	var result []microScrmModel.NotificationInfoModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchDataBase(enterpriseId, d.tableName)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *InfoDao) CreateData(c app.Ctx, enterpriseId int64, data microScrmModel.NotificationInfoModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchDataBase(enterpriseId, d.tableName)).Model(d.model).Create(&data).Error
	return err
}
func (d *InfoDao) SaveData(c app.Ctx, enterpriseId int64, data microScrmModel.NotificationInfoModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchDataBase(enterpriseId, d.tableName)).Model(d.model).Save(&data).Error
	return err
}

func (d *InfoDao) UpdateData(c app.Ctx, enterpriseId int64, query interface{}, data interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchDataBase(enterpriseId, d.tableName)).Model(&d.model).Where(query).Updates(data).Error
	return err
}

func (d *InfoDao) UpdateDataOnly(c app.Ctx, enterpriseId int64, query interface{}, data interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchDataBase(enterpriseId, d.tableName)).Where(query).Updates(data).Error
	return err
}

func (d *InfoDao) DeleteData(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) error {
	var model microScrmModel.NotificationInfoModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchDataBase(enterpriseId, d.tableName)).Model(d.model).Where(query, args...).Delete(&model).Error
	return err
}

// GetList 获取FT&RO创建的传讯列表
func (d *InfoDao) GetList(c ctx.Context, enterpriseId int64, req notificationReq.SearchListReq, infoManager notificationResp.ManagerInfo) (int64, []microScrmModel.NotificationInfoModel, error) {
	result := make([]microScrmModel.NotificationInfoModel, 0)
	var total int64
	whereMap := make(map[string]interface{}, 0)
	if req.Status >= 0 {
		whereMap["status"] = req.Status
	}
	if req.Type > 0 {
		whereMap["type"] = req.Type
	}
	if req.DepartmentName != "" {
		whereMap["department_name"] = req.DepartmentName
	}
	if (infoManager.IsFT == notificationReq.IsFT || infoManager.IsRo == notificationReq.IsRo) && req.IsAll != notificationReq.IsAll {
		if req.MemberId > 0 {
			whereMap["member_id"] = req.MemberId
		} else {
			whereMap["admin_id"] = req.AdminID
		}
	}
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(whereMap)
	if req.UpdatedStartTime > 0 && req.UpdatedEndTime > 0 {
		query.Where("updated_at BETWEEN ? AND ?", req.UpdatedStartTime, req.UpdatedEndTime)
	}
	if req.ReviewStartTime > 0 && req.ReviewEndTime > 0 {
		query.Where("review_time BETWEEN ? AND ?", req.ReviewStartTime, req.ReviewEndTime)
	}
	if req.Title != "" {
		query.Where("title Like ?", "%"+req.Title+"%")
	}
	if infoManager.IsRo == notificationReq.IsRo && req.IsAll == notificationReq.IsAll {
		query.Where("status >= ?", microScrmModel.NotificationStatusChecking).Where("status!=?", microScrmModel.NotificationStatusRecall)
	}
	if req.NoRelease == notificationReq.NoRelease {
		query.Where("merge_release", microScrmModel.NotificationMergeNotSend)
	}
	offset := (req.Page - 1) * req.Size
	err := query.Count(&total).Order("updated_at desc").Offset(offset).Limit(req.Size).Find(&result).Error
	return total, result, err
}

func (d *InfoDao) GetDepartments(c ctx.Context, enterpriseId int64) ([]microScrmModel.NotificationDepartmentModel, error) {
	result := make([]microScrmModel.NotificationDepartmentModel, 0)
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.NotificationDepartmentModel{}.TableName(), enterpriseId)).Where("status=1").Find(&result).Error
	return result, err
}

func (d *InfoDao) GetManagerInfo(c ctx.Context, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.NotificationManagerModel, error) {
	var result microScrmModel.NotificationManagerModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchDataBase(enterpriseId, microScrmModel.NotificationManagerModel{}.TableName())).Where(query, args).Find(&result).Error
	return result, err
}

func (d *InfoDao) GetAllManager(c ctx.Context, enterpriseId int64) ([]microScrmModel.NotificationManagerModel, error) {
	var result []microScrmModel.NotificationManagerModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchDataBase(enterpriseId, microScrmModel.NotificationManagerModel{}.TableName())).Find(&result).Error
	return result, err
}

func (d *InfoDao) UpdateDataColumns(c app.Ctx, enterpriseId int64, query interface{}, data interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchDataBase(enterpriseId, d.tableName)).Model(&d.model).Where(query).UpdateColumns(data).Error
	return err
}
