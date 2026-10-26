package display

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	"micro-scrm/internal/app/http/request/display"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
)

type TaskDao struct {
	dao       dao.Dao
	model     *microScrmModel.DisplayTaskModel
	tableName string
}

func NewTaskDao() *TaskDao {
	return &TaskDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		model:     new(microScrmModel.DisplayTaskModel),
		tableName: microScrmModel.DisplayTaskModel{}.TableName(),
	}
}

func (d *TaskDao) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.DisplayTaskModel, error) {
	var result microScrmModel.DisplayTaskModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *TaskDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.DisplayTaskModel, error) {
	var result []microScrmModel.DisplayTaskModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Find(&result).Error
	return result, err
}

// CreateData 新增任务创建数据
func (d *TaskDao) CreateData(c app.Ctx, enterpriseId int64, data microScrmModel.DisplayTaskModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&data).Error
	return err
}

// GetList 管理后台的任务列表查询
func (d *TaskDao) GetList(c app.Ctx, request display.SearchReq) (int64, []microScrmModel.DisplayTaskModel, error) {
	var result []microScrmModel.DisplayTaskModel
	var total int64
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseId))
	if request.DisplayName != "" {
		query.Where("display_name LIKE ?", "%"+request.DisplayName+"%")
	}
	if request.StartTime > 0 && request.EndTime > 0 {
		query.Where("start_time >= ? and end_time <=?", request.StartTime, request.EndTime)
	}
	if request.CreatedStartAt > 0 && request.CreatedEndAt > 0 {
		query.Where("created_at between ? and ? ", request.CreatedStartAt, request.CreatedEndAt)
	}
	if request.Status >= 0 {
		query.Where("status=?", request.Status)
	}
	if request.ShowType >= 0 {
		query.Where("show_type=?", request.ShowType)
	}
	offset := (request.Page - 1) * request.Size
	err := query.Count(&total).Order("id desc").Offset(offset).Limit(request.Size).Find(&result).Error
	return total, result, err
}

func (d *TaskDao) UpdateData(c app.Ctx, enterpriseId int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query).UpdateColumns(data).Error
	return err
}

// CreateDisplayTaskRelation 服务配置的权限范围配置
func (d *TaskDao) CreateDisplayTaskRelation(c app.Ctx, enterpriseId, displayTaskId int64, roleRules []display.StoreDepartmentInfo) error {
	var createInfoArr []microScrmModel.DisplayTaskRelationModel
	if len(roleRules) > 0 {
		for _, v := range roleRules {
			createInfo := microScrmModel.DisplayTaskRelationModel{
				BaseModel:     models.CreateBase(app.UID()),
				DisplayTaskID: displayTaskId,
				StoreType:     int32(v.StoreType),
				StoreID:       v.DpThirdID,
				StoreName:     v.DpStoreName,
				DepartmentId:  int64(v.DepartmentID),
				IsStore:       1,
			}
			createInfoArr = append(createInfoArr, createInfo)
		}
		return d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskRelationModel{}.TableName(), enterpriseId)).Create(&createInfoArr).Error
	}
	return nil
}

// GetDisplayStore 获取VM的指定门店数据
func (d *TaskDao) GetDisplayStore(c app.Ctx, enterpriseId, displayTaskId int64) (data []microScrmModel.DisplayTaskRelationModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskRelationModel{}.TableName(), enterpriseId)).Where("display_task_id=?", displayTaskId).Find(&data).Error
	return
}

// DeleteDisplayStore 删除门店数据
func (d *TaskDao) DeleteDisplayStore(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) error {
	var model microScrmModel.DisplayTaskRelationModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskRelationModel{}.TableName(), enterpriseId)).Where(query, args...).Delete(&model).Error
	return err
}

func (d *TaskDao) GetStoreTaskList(c ctx.Context, enterpriseId int64, storeId, finishStatus int, ids []int) ([]microScrmModel.DisplayTaskModel, error) {
	var result []microScrmModel.DisplayTaskModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
		"id": ids,
	}).Preload("TaskAllotInfo", func(db *gorm.DB) *gorm.DB {
		return db.Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotModel{}.TableName(), enterpriseId)).Where(map[string]interface{}{
			"store_id":      storeId,
			"finish_status": finishStatus,
		})
	}).Find(&result).Error
	return result, err
}

func (d *TaskDao) GetSpecialList(c ctx.Context, enterpriseId int64, query interface{}, args ...interface{}) ([]string, error) {
	var data []string
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayManagerModel{}.TableName(), enterpriseId)).Where(query, args...).Pluck("user_id", &data).Error
	return data, err
}

func (d *TaskDao) GetSpecialInfo(c ctx.Context, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.DisplayManagerModel, error) {
	var data microScrmModel.DisplayManagerModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayManagerModel{}.TableName(), enterpriseId)).Where(query, args...).Find(&data).Error
	return data, err
}
