package microScrmDao

import (
	"errors"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	reserveRequest "micro-scrm/internal/app/http/request/reserve"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"time"
)

type ReserveServiceDao struct {
	dao       dao.Dao
	tableName string
}

func NewReserveServiceDao() *ReserveServiceDao {
	return &ReserveServiceDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.ReserveServiceModel{}.TableName(),
	}
}

func (d *ReserveServiceDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.ReserveServiceModel, error) {
	var result []microScrmModel.ReserveServiceModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *ReserveServiceDao) Create(c app.Ctx, enterpriseId int64, info microScrmModel.ReserveServiceModel) (id int64, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&info).Error
	return info.ID, err
}

func (d *ReserveServiceDao) Update(c app.Ctx, enterpriseId int64, request reserveRequest.ReserveServiceUpdateReq) (err error) {
	updateMap := map[string]interface{}{
		"reserve_name": request.ReserveName,
		"reserve_img":  request.ReserveImg,
		"reserve_time": request.ReserveTime,
		// "end_time":     int(helper.GetDateToTimes(request.EndTime)),
		"end_time":     request.EndTime,
		"start_time":   request.StartTime,
		"service_type": request.ServiceType,
		"service_role": request.ServiceRole,
	}
	if len(request.TimeSlot) > 0 {
		updateMap["time_slot"] = helper.MarshalToString(request.TimeSlot)
	}
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where(map[string]interface{}{
			"id": request.ID,
		}).
		Updates(updateMap).Error
}

func (d *ReserveServiceDao) Delete(c app.Ctx, enterpriseId int64, request reserveRequest.ReserveServiceDeleteReq) (err error) {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where(map[string]interface{}{
			"id": request.ID,
		}).
		Updates(map[string]interface{}{
			"deleted_at": time.Now().Unix(),
		}).Error
}

// GetList 后台管理的服务配置列表
func (d *ReserveServiceDao) GetList(c app.Ctx, enterpriseId int64, listReq reserveRequest.ReserveServiceListReq) (list []microScrmModel.ReserveServiceModel, total int64, err error) {
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("deleted_at = 0").Scopes(func(db *gorm.DB) *gorm.DB {
		if (listReq.ReserveName) != "" {
			db.Where("reserve_name  LIKE ?  ", "%"+listReq.ReserveName+"%")
		}
		if listReq.ReserveStatus >= 0 {
			db.Where("reserve_status  = ?  ", listReq.ReserveStatus)
		}
		if listReq.StartTime > 0 && listReq.EndTime > 0 {
			db.Where("start_time <=? and end_time>=? and start_time>0", listReq.StartTime, listReq.EndTime)
		}

		return db
	})

	err = query.Count(&total).
		Order("created_at DESC").
		Offset((listReq.Page - 1) * listReq.Size).
		Limit(listReq.Size).
		Find(&list).Error

	return list, total, err

}

func (d *ReserveServiceDao) GetListAll(c app.Ctx, enterpriseId int64, listReq reserveRequest.ReserveServiceListReq) (list []microScrmModel.ReserveServiceModel, total int64, err error) {
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Unscoped()

	if (listReq.ReserveName) != "" {
		query.Where("reserve_name = ?", listReq.ReserveName)
	}

	if len(listReq.ReserveIds) > 0 {
		query.Where("id in (?)", listReq.ReserveIds)
	}

	err = query.Count(&total).
		Order("created_at DESC").
		Find(&list).Error

	return list, total, err

}

func (d *ReserveServiceDao) GetInfoById(c app.Ctx, enterpriseId, id int64) (microScrmModel.ReserveServiceModel, error) {
	var info microScrmModel.ReserveServiceModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Unscoped().Where(id).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return info, err
}

// CreateServiceRelation 服务配置的权限范围配置
func (d *ReserveServiceDao) CreateServiceRelation(c app.Ctx, enterpriseId, serviceId int64, roleRules []reserveRequest.RoleRules) error {
	var createInfoArr []microScrmModel.ReserveServiceRelationModel
	if len(roleRules) > 0 {
		for _, v := range roleRules {
			createInfo := microScrmModel.ReserveServiceRelationModel{
				BaseModel: models.CreateBase(app.UIDContext(c)),
				ServiceID: serviceId,
				StoreType: v.StoreType,
			}
			if v.IsDP {
				createInfo.StoreID = v.ThirdId
				createInfo.StoreName = v.Name
				createInfo.IsStore = 1
				createInfo.DepartmentId = v.ID
			} else {
				createInfo.MemberID = v.ID
				createInfo.MemberName = v.Name
				createInfo.UserId = v.UserId
			}
			createInfoArr = append(createInfoArr, createInfo)
		}
		return d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.ReserveServiceRelationModel{}.TableName(), enterpriseId)).Create(&createInfoArr).Error
	}
	return nil
}

// DeleteServiceRelation 删除有关的服务配置权限
func (d *ReserveServiceDao) DeleteServiceRelation(c app.Ctx, enterpriseId, serviceId int64) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.ReserveServiceRelationModel{}.TableName(), enterpriseId)).Where("service_id=?", serviceId).Delete(&microScrmModel.ReserveServiceRelationModel{}).Error
}

// GetServiceRole 获取服务配置有关的权限范围
func (d *ReserveServiceDao) GetServiceRole(c app.Ctx, enterpriseId, serviceId int64) (data []microScrmModel.ReserveServiceRelationModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.ReserveServiceRelationModel{}.TableName(), enterpriseId)).Where("service_id=?", serviceId).Find(&data).Error
	return
}

// GetListByMember 导购获取服务配置
func (d *ReserveServiceDao) GetListByMember(c app.Ctx, enterpriseId int64, memberId int64, storeIds []int64) ([]microScrmModel.ReserveServiceModel, error) {
	nowTime := time.Now().Unix()
	var result, data1, data2 []microScrmModel.ReserveServiceModel
	var err error
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("reserve_status=0").Where("service_role=0").Scopes(func(db *gorm.DB) *gorm.DB {
		db.Where("service_type=1 OR (service_type=0 AND end_time>? AND start_time<?)", nowTime, nowTime)
		return db
	}).Find(&data1).Error
	if len(storeIds) > 0 || memberId != 0 { // 这里是查询配置权限范围的服务
		relationModel := d.dao.GetTable(microScrmModel.ReserveServiceRelationModel{}.TableName(), enterpriseId)
		err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("reserve_service.reserve_status=0").Select("reserve_service.*").
			Joins("LEFT JOIN " + relationModel + " as b ON b.service_id=reserve_service.id").Where("b.deleted_at=0").Scopes(func(db *gorm.DB) *gorm.DB {
			if memberId != 0 && len(storeIds) > 0 {
				db.Where("(member_id=? or store_id in ?) ", memberId, storeIds)

			}
			if memberId != 0 && len(storeIds) == 0 {
				db.Where("member_id=? ", memberId)
			}
			if memberId == 0 && len(storeIds) > 0 {
				db.Where("store_id in ?  ", storeIds)
			}
			return db
		}).Group("reserve_service.Id").Find(&data2).Error
	}
	result = append(result, data1...)
	result = append(result, data2...)
	return result, err
}

// UpdateServeStatus 更新服务预约的状态
func (d *ReserveServiceDao) UpdateServeStatus(c app.Ctx, enterpriseId, endTime int64) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Unscoped().Where("service_type=0 and end_time<?", endTime).Update("reserve_status", microScrmModel.FINISH).Error

}

// GetListsByServiceIds 根据服务配置的ids,获取服务配置的权限范围
func (d *ReserveServiceDao) GetListsByServiceIds(c app.Ctx, enterpriseId int64, serviceIds []int64) (map[int64][]reserveRequest.RoleRules, error) {
	var result []microScrmModel.ReserveServiceModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where(map[string]interface{}{
			"id":         serviceIds,
			"deleted_at": 0,
		}).Find(&result).Error
	if err != nil {
		return nil, err
	}
	var res = make(map[int64][]reserveRequest.RoleRules)
	for _, v := range result {
		if v.ServiceRole == 1 {
			roles, err := d.GetServiceRole(c, enterpriseId, v.ID)
			if err != nil {
				continue
			}
			for _, role := range roles {
				item := reserveRequest.RoleRules{
					ID:                 role.MemberID,
					ThirdId:            role.StoreID,
					StoreType:          role.StoreType,
					MemSubscribeStatus: 1,
					UserId:             role.UserId,
				}
				if role.IsStore == 0 { // 不是门店，name为导购名称
					item.IsDP = false
					item.Name = role.MemberName
				} else {
					item.Name = role.StoreName
					item.IsDP = true
					item.ID = role.DepartmentId
				}
				res[v.ID] = append(res[v.ID], item)
			}
		}
	}
	return res, err
}
