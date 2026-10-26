package microScrmDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	planReq "micro-scrm/internal/app/http/request/plan"
	membersModel "micro-scrm/internal/app/models/members"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type PlanImportDao struct {
	dao       dao.Dao
	tableName string
}

func NewPlanImportDao() *PlanImportDao {
	return &PlanImportDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.TaskImportFileModel{}.TableName(),
	}
}

func (d *PlanImportDao) FileDetection(c app.Ctx, enterpriseID int64, fileName string) (result int64, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseID)).Select("COUNT(1)").
		Where("file_name = ? AND deleted_at = 0 AND is_extract != 0", fileName).
		Find(&result).Error
	return
}

func (d *PlanImportDao) SaveImportFile(c app.Ctx, enterpriseID int64, info microScrmModel.TaskImportFileModel) (Id int64, err error) {
	res := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseID)).Create(&info)
	Id = info.ID
	err = res.Error
	return
}

func (d *PlanImportDao) GetImportFile(c app.Ctx, enterpriseID, ID int64) (data microScrmModel.TaskImportFileModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseID)).Where("id=?", ID).First(&data).Error
	return
}

// GetInsertImportDataInfoByUnionId 根据用户Unionid查询信息匹配插入
func (d *PlanImportDao) GetInsertImportDataInfoByUnionId(c app.Ctx, enterpriseID int64, UnionIdData []*planReq.UnionIDData, sendType int) (data []planReq.ImportInsertData, err error) {
	var unionIds []string
	for _, v := range UnionIdData {
		if v.UnionID != "" {
			unionIds = append(unionIds, v.UnionID)
		}
	}
	query := d.dao.WithContext(c).Model(&microScrmModel.WorkExternalUserRelationModel{}).
		Select("`user_id`,`unionid` as `union_id`").
		Where("enterprise_id=? and relation_status=0", enterpriseID).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if sendType == microScrmModel.SendTypeKeySa {
				db.Where("is_key_sa=1")
			}
			return db
		}).Where("unionid in ?", unionIds)
	if sendType == microScrmModel.SendTypeKeySa {
		query = query.Group("unionid")
	}
	err = query.Scan(&data).Error
	if len(data) > 0 {
		var userIds []string
		for _, val := range data {
			userIds = append(userIds, val.UserID)
		}
		externalMap, memberNameMap, err := d.getExternalAndMember(c, enterpriseID, unionIds, userIds)
		if err != nil {
			return nil, err
		}
		for i, val := range data {
			if external, ok := externalMap[val.UnionID]; ok {
				data[i].WorkExternalUserId = external.Openid
				data[i].UserName = external.Name
			}
			if name, ok := memberNameMap[val.UserID]; ok {
				data[i].Name = name
			}
		}
	}

	return
}

// GetInsertImportDataInfoBySaCodeAndUnionId 根据用户unionid和导购的sacode查询插入
func (d *PlanImportDao) GetInsertImportDataInfoBySaCodeAndUnionId(c app.Ctx, enterpriseID int64, UnionIdData []*planReq.UnionIDData) (data []planReq.ImportInsertData, err error) {
	var whereIn [][]interface{}
	var unionIds []string
	var userIds []string
	for _, v := range UnionIdData {
		if v.UnionID != "" && v.SaCode != "" {
			unionIds = append(unionIds, v.UnionID)
			userIds = append(userIds, v.SaCode)
			whereIn = append(whereIn, []interface{}{
				v.UnionID, v.SaCode,
			})
		}
	}
	if len(whereIn) == 0 {
		return
	}

	err = d.dao.WithContext(c).Model(&microScrmModel.WorkExternalUserRelationModel{}).
		Select("`user_id`,`unionid` as `union_id`").
		Where("enterprise_id=? and relation_status=0", enterpriseID).
		Where("(unionid, user_id) in ?", whereIn).
		Scan(&data).Error
	if len(data) > 0 {
		externalMap, memberNameMap, err := d.getExternalAndMember(c, enterpriseID, unionIds, userIds)
		if err != nil {
			return nil, err
		}
		for i, val := range data {
			if external, ok := externalMap[val.UnionID]; ok {
				data[i].WorkExternalUserId = external.Openid
				data[i].UserName = external.Name
			}
			if name, ok := memberNameMap[val.UserID]; ok {
				data[i].Name = name
			}
		}
	}

	return
}

func (d *PlanImportDao) getExternalAndMember(c app.Ctx, enterpriseID int64, unionIds, userIds []string) (map[string]microScrmModel.WxUserCorpExternalModel, map[string]string, error) {
	var externalMap = make(map[string]microScrmModel.WxUserCorpExternalModel)
	var memberNameMap = make(map[string]string)
	if len(unionIds) > 0 {
		var externalList []microScrmModel.WxUserCorpExternalModel
		err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.WxUserCorpExternalModel{}.TableName(), enterpriseID)).
			Where("unionid in ?", unionIds).
			Select("openid,unionid,name").
			Find(&externalList).Error
		if err != nil {
			return externalMap, memberNameMap, err
		}
		for _, model := range externalList {
			externalMap[model.Unionid] = model
		}
	}

	if len(userIds) > 0 {
		var memberList []membersModel.MemberModel
		err := d.dao.WithContext(c).Model(&membersModel.MemberModel{}).
			Where("enterprise_id=? AND userid in ?", enterpriseID, userIds).
			Select("userid,name").
			Find(&memberList).Error
		if err != nil {
			return externalMap, memberNameMap, err
		}
		for _, model := range memberList {
			memberNameMap[model.Userid] = model.Name
		}
	}

	return externalMap, memberNameMap, nil
}

// UpdateImportFileExtract 人群包匹配完毕之后将结果的文件上传cos,保存cos地址
func (d *PlanImportDao) UpdateImportFileExtract(c app.Ctx, enterpriseID, ID int64, path string) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseID)).Where("id=?", ID).Updates(map[string]interface{}{
		"is_extract":   1,
		"download_url": path,
	}).Error
}
