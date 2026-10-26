package microScrmDao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	customerReq "micro-scrm/internal/app/http/request/customer"
	"micro-scrm/internal/app/http/request/external"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"strings"
)

type WxUserCorpExternal struct {
	dao       dao.Dao
	tableName string
}

func NewWxUserCorpExternal() *WxUserCorpExternal {
	return &WxUserCorpExternal{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.WxUserCorpExternalModel{}.TableName(),
	}
}

func (d *WxUserCorpExternal) GetOne(c app.Ctx, query interface{}, enterpriseId int64) (microScrmModel.WxUserCorpExternalModel, error) {
	var result microScrmModel.WxUserCorpExternalModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query).Find(&result).Error
	return result, err
}

func (d *WxUserCorpExternal) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.WxUserCorpExternalModel, error) {
	var result []microScrmModel.WxUserCorpExternalModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *WxUserCorpExternal) GetUnionIdsByOpenIds(c app.Ctx, enterpriseId int64, openIds []string) ([]string, error) {
	var unionIds []string
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("openid in ?", openIds).Pluck("unionid", &unionIds).Error
	return unionIds, err
}

func (d *WxUserCorpExternal) GetUnionIdMapByOpenIds(c app.Ctx, enterpriseId int64, openIds []string) (map[string]string, error) {
	var unionIdMap = make(map[string]string, 0)
	var list = make([]microScrmModel.WxUserCorpExternalModel, 0)
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("openid in ?", openIds).Find(&list).Error
	for _, item := range list {
		unionIdMap[item.Openid] = item.Unionid
	}
	return unionIdMap, err
}

func (d *WxUserCorpExternal) GetManyByCardLevelOrCltType(c app.Ctx, enterpriseId int64, cardLevel, cltType []string, workUserRelationIds []int64) ([]microScrmModel.WxUserCorpExternalModel, error) {
	var result []microScrmModel.WxUserCorpExternalModel

	err := d.dao.WithContext(c).
		Where("id", workUserRelationIds).
		Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Scopes(func(db *gorm.DB) *gorm.DB {
		var cardLevels []string
		if len(cardLevel) > 0 {
			if len(cardLevel) == 1 {
				switch cardLevel[0] {
				case customerReq.HIGH_GRADE:
					//db.Where("amount_total > 0")
					db.Where("clt_type <> '' ")
				case customerReq.POTENTIAL:
					db.Where("card_level = ? and amount_total = 0", customerReq.PROSPECT)
				case "O":
					cardLevel = append(cardLevel, "O", "")
					db.Where("card_level in ?", cardLevel)
				default:
					db.Where("card_level = ?", cardLevel[0])
				}
			} else {
				highGrade := false
				potential := false
				for _, v := range cardLevel {
					if v == customerReq.HIGH_GRADE {
						highGrade = true
					} else if v == customerReq.POTENTIAL {
						potential = true
					} else {
						if v == "O" { // 非会员 或为空
							cardLevels = append(cardLevels, "O", "")
						} else {
							cardLevels = append(cardLevels, v)
						}
					}
				}
				if len(cardLevels) > 0 {
					if !highGrade && !potential {
						db.Where("card_level in (?)", cardLevels)
					}
					if highGrade && potential {
						//db.Where("card_level in (?) or amount_total > 0 or (card_level = ? and amount_total = 0) ", cardLevels, customerReq.PROSPECT)
						db.Where("card_level in (?) or clt_type <> ''  or (card_level = ? and amount_total = 0) ", cardLevels, customerReq.PROSPECT)
					}
					if highGrade && !potential {
						//db.Where("card_level in (?) or amount_total > 0", cardLevels)
						db.Where("card_level in (?) or clt_type <> '' ", cardLevels)
					}
					if potential && !highGrade {
						db.Where("card_level in (?) or (card_level = ? and amount_total = 0)", cardLevels, customerReq.PROSPECT)
					}
				} else {
					if highGrade && potential {
						//db.Where("amount_total > 0 or (card_level = ? and amount_total = 0)", customerReq.PROSPECT)
						db.Where("clt_type <> ''  or (card_level = ? and amount_total = 0)", customerReq.PROSPECT)
					}
				}
			}
		}
		if len(cltType) > 0 {
			db.Where("clt_type in (?)", cltType)
		}
		return db
	}).Find(&result).Error

	return result, err
}

func (d *WxUserCorpExternal) GetCount(c app.Ctx, enterpriseId int64) (count int64, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where(" unionid <> ' '").
		Count(&count).Error
	return
}

func (d *WxUserCorpExternal) GetList(c app.Ctx, enterpriseId int64, page, size int) ([]microScrmModel.WxUserCorpExternalModel, error) {
	var result []microScrmModel.WxUserCorpExternalModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where(" unionid <> ' '").
		Order("id Desc").
		Offset((page - 1) * size).Limit(size).
		Find(&result).Error
	return result, err
}

// GetListLimit 分批查询
func (d *WxUserCorpExternal) GetListLimit(c app.Ctx, enterpriseId int64, pageLimit int, beginID int64) ([]microScrmModel.WxUserCorpExternalModel, error) {
	var result []microScrmModel.WxUserCorpExternalModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("type = ? ", 1).
		Where("union_code = ' '").
		Where("unionid <> ' '").
		Where("id > ? ", beginID).
		Order("id ASC").
		Limit(pageLimit).
		Find(&result).Error

	return result, err
}

func (d *WxUserCorpExternal) FindOne(c app.Ctx, enterpriseId int64, openid string) (microScrmModel.WxUserCorpExternalModel, error) {
	var result microScrmModel.WxUserCorpExternalModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where(" openid = ? ", openid).
		Find(&result).Error
	return result, err
}

// GetInfoByUnionId 获取用户信息
func (d *WxUserCorpExternal) GetInfoByUnionId(c app.Ctx, enterpriseId int64, unionId string) (info microScrmModel.WxUserCorpExternalModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("unionid=?", unionId).Find(&info).Error
	return
}

// UpdateBirthByUnionId 更改用户的生日
func (d *WxUserCorpExternal) UpdateBirthByUnionId(c app.Ctx, enterpriseId int64, unionId string, month, day int) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("unionid", unionId).
		Updates(map[string]interface{}{
			"BirthMonth": month,
			"BirthDay":   day,
		}).Error
}

// SyncCdpInfo 同步cdp信息
func (d *WxUserCorpExternal) SyncCdpInfo(c app.Ctx, enterpriseId int64, info external.CdpExternalInfo, updateData interface{}) error {
	if info.UnionIdBindFlag != 1 {
		return nil
	}
	if info.UnionId == "" {
		return nil
	}
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("unionid", info.UnionId).
		Updates(updateData).Error
	//return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("unionid", info.UnionId).
	//	Updates(map[string]interface{}{
	//		"union_code":    info.MemberShip,
	//		"customer_code": info.CustomerCode,
	//		"customer_id":   info.ID,
	//		"amount_total":  info.Last12MPurchaseAmount,
	//		"card_level":    info.CustomerType,
	//		"clt_type":      info.ClientType,
	//		"name":          info.Name,
	//	}).Error
}

// GetWxUserCorpExternalByUnionIds 根据unionids获取外部联系人信息
func (d *WxUserCorpExternal) GetWxUserCorpExternalByUnionIds(c app.Ctx, enterpriseId int64, unionIds []string) (lists []microScrmModel.WxUserCorpExternalModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("unionid in ?", unionIds).Find(&lists).Error
	return
}

func (d *WxUserCorpExternal) GetWxUserCorpExternalByUnionCode(c app.Ctx, cond customerReq.CustomerListsReq) ([]microScrmModel.ExternalUserList, error) {
	//var result []external.WxUserCorpExternalInfo
	var result []microScrmModel.ExternalUserList

	//selectStr := "work_external_user_relation.id,work_external_user_relation.enterprise_id,work_external_user_relation.remark,work_external_user_relation.member_id,work_external_user_relation.work_external_user_id,work_external_user_relation.user_id as sa_user_id,work_external_user_relation.is_key_sa,work_external_user_relation.created_at," +
	//	"c.unionid,c.openid,c.euid,c.nick,c.name,c.avatar,c.birthday,c.union_code,c.mobile,c.card_level,c.clt_type,c.amount_total,"

	selectStr := "id as work_external_user_id,unionid,openid,euid,nick,name,avatar,birthday,union_code,mobile,card_level,clt_type,amount_total"

	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, cond.EnterpriseID)).
		Select(strings.TrimRight(selectStr, ",")).
		Where("union_code in (?)", cond.UnionCodes).Find(&result).Error

	return result, err
}

// GetListForCond 外部联系人分页条件查询
func (d *WxUserCorpExternal) GetListForCond(c app.Ctx, request customerReq.ExternalUserListsReq, query interface{}, args ...interface{}) ([]microScrmModel.ExternalUserList, int64, error) {
	var (
		total  int64
		result []microScrmModel.ExternalUserList
		err    error
	)

	//selectStr := "work_external_user_relation.id,work_external_user_relation.enterprise_id,work_external_user_relation.remark,work_external_user_relation.member_id,work_external_user_relation.work_external_user_id,work_external_user_relation.user_id as sa_user_id,work_external_user_relation.is_key_sa,work_external_user_relation.created_at," +
	//	"c.unionid,c.openid,c.euid,c.nick,c.name,c.avatar,c.birthday,c.union_code,c.mobile,c.card_level,c.clt_type,c.amount_total,c.customer_code,c.last_purchase_date,c.last_purchase_store_code,c.upgrade_need_amount,c.customer_type,"
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseID)).
		Select("id as work_external_user_id, unionid, openid, euid, nick, name, avatar, birthday, union_code, mobile, card_level, clt_type, amount_total, customer_code, last_purchase_date, last_purchase_store_code, upgrade_need_amount, customer_type").
		Where("unionid <> ' '").
		Where(query, args...).
		Order("id ASC").
		Count(&total).
		Offset((request.Page - 1) * request.Size).Limit(request.Size).
		Find(&result).Error

	return result, total, err
}

func (d *WxUserCorpExternal) UpdateData(c app.Ctx, enterpriseId int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query).UpdateColumns(data).Error
	return err
}

func (d *WxUserCorpExternal) SaveData(c app.Ctx, enterpriseId int64, data microScrmModel.WxUserCorpExternalModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&data).Error
	return err
}

func (d *WxUserCorpExternal) GetIdsBySearch(c app.Ctx, cond customerReq.SearchUserIdsReq) (ids []int64, err error) {
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, cond.EnterpriseID)).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if cond.Name != "" {
				db.Where("name Like ?", "%"+cond.Name+"%")
			}
			if cond.UnionCode != "" {
				db.Where("union_code = ?", cond.UnionCode)
			}
			return db
		})
	err = query.Pluck("id", &ids).Error
	return
}

// GetInfoByIds  根据id获取外部联系人信息
func (d *WxUserCorpExternal) GetInfoByIds(c app.Ctx, enterpriseId int64, ids []int64) (lists []microScrmModel.WxUserCorpExternalModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("id in ?", ids).Find(&lists).Error
	return
}

// GetOpenIdsByUnionIds   根据unionid获取openid
func (d *WxUserCorpExternal) GetOpenIdsByUnionIds(c app.Ctx, enterpriseId int64, unionIds []string) (openIds []string, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("unionid in ?", unionIds).Pluck("openid", &openIds).Error
	return
}

func (d *WxUserCorpExternal) BatchUpdateByUnionIds(c context.Context, enterpriseId int64, mappings []microScrmModel.UnionidMapModel, nowUnix int64) (int64, error) {
	if len(mappings) == 0 {
		return 0, nil
	}

	// 获取当前批次的最小和最大ID
	minID := mappings[0].ID
	maxID := mappings[len(mappings)-1].ID

	query := `
        UPDATE ` + fmt.Sprintf("%s%d.%s", app.Cfg().DB.DBPrefix, enterpriseId, d.tableName) + ` AS t1
        INNER JOIN ` + fmt.Sprintf("%s%d.%s", app.Cfg().DB.DBPrefix, enterpriseId, "unionid_map") + ` AS t2 ON t1.unionid = t2.old_unionid
        SET t1.unionid = t2.new_unionid, t1.updated_at = ?
        WHERE t2.id >= ? AND t2.id <= ?`

	db := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Exec(query, nowUnix, minID, maxID)
	return db.RowsAffected, db.Error
}
