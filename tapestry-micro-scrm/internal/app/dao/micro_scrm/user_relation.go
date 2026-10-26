package microScrmDao

import (
	"context"
	"fmt"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	"micro-scrm/internal/app/http/request"
	customerReq "micro-scrm/internal/app/http/request/customer"
	memberReq "micro-scrm/internal/app/http/request/member"
	openReq "micro-scrm/internal/app/http/request/open"
	planReq "micro-scrm/internal/app/http/request/plan"
	drainageContactResp "micro-scrm/internal/app/http/response/drainage_contact"
	"micro-scrm/internal/app/models"
	workExternalModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"strconv"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
)

type UserRelation struct {
	dao       dao.Dao
	tableName string
}

func NewUserRelation() *UserRelation {
	return &UserRelation{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: workExternalModel.WorkExternalUserRelationModel{}.TableName(),
	}
}

type BaCodeAndUnionid struct {
	UnionID string `json:"unionid"`
	SaCode  string `json:"sa_code"`
}

func (d *UserRelation) UpdateKeySa(c app.Ctx, enterpriseId int64, updateOldByUnionid []string, updateNew [][]string) error {

	tx := d.dao.WithContext(c).Begin()
	// 将文件中相关的客户专属导购关系标识清空
	res := tx.Table(d.tableName).Model(&workExternalModel.WorkExternalUserRelationModel{}).
		Where("enterprise_id=?", enterpriseId).
		Where("unionid in ?", updateOldByUnionid).
		Where("is_key_sa = 1").
		Update("is_key_sa", 0)
	if res.Error == nil {
		fmt.Printf("移除专属导购关系数量%d\n", res.RowsAffected)
	}

	// 将客户表中有KaySa和没有KaySa的导购分类出来
	var noKeySaUnionidSlice []string
	var existUpdateNew [][]string
	for _, row := range updateNew {
		if row[1] == "" {
			noKeySaUnionidSlice = append(noKeySaUnionidSlice, row[0])
		} else if row[0] != "" && row[1] != "" {
			existUpdateNew = append(existUpdateNew, row)
		}
	}
	userExternalTable := " micro_scrm_enterprise_" + fmt.Sprintf("%d", enterpriseId) + ".wx_user_corp_external "
	res = tx.Exec("UPDATE "+userExternalTable+" SET ba_code = '' WHERE unionid in ?", noKeySaUnionidSlice)
	if res.Error == nil {
		fmt.Printf("移除客户表BaCode数量%d\n", res.RowsAffected)
	}
	// 将文件中相关的客户标记专属导购
	res = tx.Table(d.tableName).Model(&workExternalModel.WorkExternalUserRelationModel{}).
		Where("enterprise_id=?", enterpriseId).
		Where("(unionid, user_id) in ?", existUpdateNew).
		Update("is_key_sa", 1)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	fmt.Printf("新增专属导购数量%d\n", res.RowsAffected)

	// 将客户表中的导购标识更新
	res = tx.Exec("UPDATE work_external_user_relation relation"+
		" LEFT JOIN "+userExternalTable+" external ON external.id=relation.work_external_user_id"+
		" SET external.ba_code = relation.`user_id`"+
		" WHERE relation.enterprise_id="+fmt.Sprintf("%d", enterpriseId)+
		" AND relation.unionid in ?"+
		" AND relation.is_key_sa=1", updateOldByUnionid)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	fmt.Printf("客户表专属导购更新数量%d\n", res.RowsAffected)

	return tx.Commit().Error
}

func (d *UserRelation) GetUserRelationByUnionId(c app.Ctx, enterpriseId int64, unionId string) (info workExternalModel.WorkExternalUserRelationModel, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"unionid":       unionId,
		"is_key_sa":     1,
		"deleted_at":    0,
	}).Find(&info).Error
	return
}

func (d *UserRelation) GetAllByUnionId(c app.Ctx, enterpriseId int64, unionid string, status []int) ([]workExternalModel.WorkExternalUserRelationModel, error) {
	var result []workExternalModel.WorkExternalUserRelationModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"unionid":       unionid}).
		Where("relation_status in (?)", status).
		Find(&result).Error
	return result, err
}

// GetUserSAInfoByUnionId 查询外部联系人专属SA信息
func (d *UserRelation) GetUserSAInfoByUnionId(c app.Ctx, enterpriseId int64, unionId string) (info workExternalModel.UserSaInfo, err error) {

	err = d.dao.WithContext(c).Table(d.tableName).
		Select("work_external_user_relation.member_id,work_external_user_relation.description,work_external_user_relation.user_id as sa_user_id,work_external_user_relation.is_key_sa,b.`name` as sa_user_name,d.id as department_id,d.`name` as department_name").
		Joins(" LEFT JOIN members b ON work_external_user_relation.member_id = b.id").
		Joins(" LEFT JOIN member_departments c ON b.id = c.member_id").
		Joins(" LEFT JOIN departments d ON c.department_id = d.id").
		Where(map[string]interface{}{
			"work_external_user_relation.enterprise_id": enterpriseId,
			"work_external_user_relation.unionid":       unionId,
			"work_external_user_relation.is_key_sa":     1,
			"work_external_user_relation.deleted_at":    0,
		}).First(&info).Error

	return

}

// GetUserSAInfoByUnionIds 查询外部联系人专属SA信息
func (d *UserRelation) GetUserSAInfoByUnionIds(c app.Ctx, enterpriseId int64, unionIds []string) (info []workExternalModel.UserSaInfo, err error) {

	err = d.dao.WithContext(c).Table(d.tableName).
		Select("work_external_user_relation.unionid,work_external_user_relation.member_id,work_external_user_relation.description,work_external_user_relation.user_id as sa_user_id,work_external_user_relation.is_key_sa,b.`name` as sa_user_name,d.id as department_id,d.`name` as department_name").
		Joins(" LEFT JOIN members b ON work_external_user_relation.member_id = b.id").
		Joins(" LEFT JOIN member_departments c ON b.id = c.member_id").
		Joins(" LEFT JOIN departments d ON c.department_id = d.id").
		Where(map[string]interface{}{
			"work_external_user_relation.enterprise_id": enterpriseId,
			"work_external_user_relation.is_key_sa":     1,
			"work_external_user_relation.deleted_at":    0,
		}).
		Where("work_external_user_relation.unionid in ?", unionIds).
		Find(&info).Error

	return

}

// GetRepeatCustomerRelationInfo 获取重复外部联系人id
func (d *UserRelation) GetRepeatCustomerRelationInfo(c app.Ctx, enterpriseId int64, ids []int64) (result []workExternalModel.RepeatCustomerRelation, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Where("work_external_user_relation.enterprise_id=?", enterpriseId).
		Where("work_external_user_relation.relation_status=0").
		Where("work_external_user_relation.work_external_user_id in ?", ids).
		Select("work_external_user_relation.work_external_user_id,count(1) as num").
		Group("work_external_user_relation.work_external_user_id").Find(&result).Error
	return
}

func (d *UserRelation) UpdateRemark(c app.Ctx, enterpriseId int64, id, memberID int64, remark string) error {
	return d.dao.WithContext(c).Table(d.tableName).Model(&workExternalModel.WorkExternalUserRelationModel{}).
		Where("enterprise_id=?", enterpriseId).
		Where("work_external_user_id = ?", id).
		Where("member_id = ?", memberID).
		Update("remark", remark).Error
}

func (d *UserRelation) GetUserByUnionIdAndMemberIds(c app.Ctx, enterpriseId int64, unionId string, memberIDs []int64) ([]workExternalModel.WorkExternalUserRelationModel, error) {
	var result []workExternalModel.WorkExternalUserRelationModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id=?", enterpriseId).
		Where("unionId = ?", unionId).
		Where("member_id in (?)", memberIDs).
		Where("relation_status = 0").
		Find(&result).Error
	return result, err
}

func (d *UserRelation) GetMemberIdsByUnionId(c app.Ctx, enterpriseId int64, unionId string) ([]int64, error) {
	var result []int64
	err := d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"unionid":       unionId,
		"deleted_at":    0,
	}).Pluck("member_id", &result).Error
	return result, err
}

// GetUserList 根据条件筛选导购和外部联系人列表
func (d *UserRelation) GetUserList(c app.Ctx, request customerReq.ExternalUserListsReq) ([]workExternalModel.ExternalUserList, int64, error) {

	var (
		total  int64
		result []workExternalModel.ExternalUserList
		err    error
	)
	selectStr := "work_external_user_relation.id,work_external_user_relation.enterprise_id,work_external_user_relation.remark,work_external_user_relation.member_id,work_external_user_relation.work_external_user_id,work_external_user_relation.user_id as sa_user_id,work_external_user_relation.is_key_sa,work_external_user_relation.created_at," +
		"c.unionid,c.openid,c.euid,c.nick,c.name,c.avatar,c.birthday,c.union_code,c.mobile,c.card_level,c.clt_type,c.amount_total,c.customer_code,c.last_purchase_date,c.last_purchase_store_code,c.upgrade_need_amount,c.customer_type,"

	query := d.HandelQuerySql(c, request, selectStr)

	if request.IsGroup == true {
		// query.Order("work_external_user_relation.id desc")
		query.Group("work_external_user_relation.work_external_user_id")
	} else {
		// query.Order("work_external_user_relation.id desc")
		// query.Order("work_external_user_relation.work_external_user_id desc")
		query.Group("work_external_user_relation.id")
	}
	query.Select(strings.TrimRight(selectStr, ","))
	err = query.Count(&total).
		Offset((request.Page - 1) * request.Size).Limit(request.Size).
		Find(&result).Error

	return result, total, err

	/**
	query := d.dao.WithContext(c).Table(d.tableName).
		Joins(" LEFT JOIN "+userCorpExternalTable+" c ON work_external_user_relation.work_external_user_id = c.id").
		Where(map[string]interface{}{
			"work_external_user_relation.enterprise_id":   request.EnterpriseID,
			"work_external_user_relation.deleted_at":      0,
			"work_external_user_relation.relation_status": 0,
		}).
		Where("c.unionid <> '' ").
		Where("work_external_user_relation.member_id in (?)", request.MemberList)

	// 1【专属】指ks_sa所属顾问是自己的顾客
	// 2【非专属】指所属顾问不是自己的顾客
	if request.ManageType == 1 {
		query.Where("work_external_user_relation.is_key_sa = 1")
		// SM有选择的导购
		if len(request.Userids) > 0 {
			query.Where("work_external_user_relation.user_id in (?)", request.Userids)
		} else {
			// SM无选择的导购 或 SA
			query.Where("work_external_user_relation.user_id = (?)", request.Userid)
		}
	} else if request.ManageType == 2 {
		query.Where("work_external_user_relation.is_key_sa = 0")
	}
	if len(request.TagList) > 0 {
		selectStr += "b.external_tag_id,"
		query.Joins(" LEFT JOIN work_external_user_tag b ON work_external_user_relation.id = b.relation_id").
			Where("b.deleted_at = 0").Where(" b.external_tag_id in (?)  ", request.TagList)
	}
	if len(request.ImpressionTagList) > 0 {
		query.Joins(" LEFT JOIN impression_tag_relation d ON c.euid = d.external_user_id").
			Where("d.deleted_at = 0").Where(" d.tag_id in (?)  ", request.ImpressionTagList)
	}
	if len(request.CltType) > 0 {
		query.Where("c.clt_type in (?)", request.CltType)
	}
	if len(request.CardType) > 0 {
		var cardLevel []string
		if len(request.CardType) == 1 {
			switch request.CardType[0] {
			case customerReq.HIGH_GRADE:
				// query.Where("c.amount_total > 0")
				// 优质客 clt类型有值
				query.Where("c.clt_type <> '' ")
			case customerReq.POTENTIAL:
				query.Where("c.card_level = ? and c.amount_total = 0", customerReq.PROSPECT)
			case "O":
				cardLevel = append(cardLevel, "O", "")
				query.Where("c.card_level in ?", cardLevel)
			default:
				query.Where("c.card_level = ?", request.CardType[0])
			}
		} else {
			highGrade := false
			potential := false
			for _, v := range request.CardType {
				if v == customerReq.HIGH_GRADE {
					highGrade = true
				} else if v == customerReq.POTENTIAL {
					potential = true
				} else {
					if v == "O" { // 非会员 或为空
						cardLevel = append(cardLevel, "O", "")
					} else {
						cardLevel = append(cardLevel, v)
					}
				}
			}
			if len(cardLevel) > 0 {
				if !highGrade && !potential {
					query.Where("c.card_level in (?)", cardLevel)
				}
				if highGrade && potential {
					// query.Where("c.card_level in (?) or c.amount_total > 0 or (c.card_level = ? and c.amount_total = 0) ", cardLevel, customerReq.PROSPECT)
					query.Where("c.card_level in (?) or c.clt_type <> '' or (c.card_level = ? and c.amount_total = 0) ", cardLevel, customerReq.PROSPECT)
				}
				if highGrade && !potential {
					// query.Where("c.card_level in (?) or c.amount_total > 0", cardLevel)
					query.Where("c.card_level in (?) or c.clt_type <> '' ", cardLevel)
				}
				if potential && !highGrade {
					query.Where("c.card_level in (?) or (c.card_level = ? and c.amount_total = 0)", cardLevel, customerReq.PROSPECT)
				}
			} else {
				if highGrade && potential {
					// query.Where("c.amount_total > 0 or (c.card_level = ? and c.amount_total = 0)", customerReq.PROSPECT)
					query.Where("c.clt_type <> '' or (c.card_level = ? and c.amount_total = 0)", customerReq.PROSPECT)
				}
			}
		}
	}

	// 会员code
	if len(request.UnionCodes) > 0 {
		query.Where("c.union_code in (?)", request.UnionCodes)
	}
	// sw顾客类型
	if len(request.CustomerType) > 0 {
		customerType := request.CustomerType
		// vip客
		if helper.InArray(workExternalModel.VipTypeChoice, request.CustomerType) {
			query.Where("c.vip_type = (?)", workExternalModel.VipType)

			item := []int{workExternalModel.VipTypeChoice}
			customerType = helper.DelItems(request.CustomerType, item)
		}
		if len(customerType) > 0 {
			query.Where("c.customer_type in (?)", customerType)
		}
	}
	// sw会员最近购买时间
	if request.LastPurchaseStartTime > 0 && request.LastPurchaseEndTime > 0 {
		query.Where("c.last_purchase_date >= (?) and c.last_purchase_date < (?) ", request.LastPurchaseStartTime, request.LastPurchaseEndTime)
	}
	// sw会员生日月份时间
	if request.BirthdayStart > 0 && request.BirthdayEnd > 0 {
		var birthdayArr []string
		for i := request.BirthdayStart; i <= request.BirthdayEnd; i++ {
			birthdayArr = append(birthdayArr, helper.EncodeToString(strconv.Itoa(i)))
		}
		query.Where("c.birth_month in (?) ", birthdayArr)
	}

	if request.Name != "" {
		query.Where("work_external_user_relation.remark LIKE ? OR c.nick LIKE ?", "%"+request.Name+"%", "%"+request.Name+"%")
	}

	if request.IsGroup == true {
		// query.Order("work_external_user_relation.id desc")
		query.Group("work_external_user_relation.work_external_user_id")
	} else {
		// query.Order("work_external_user_relation.id desc")
		// query.Order("work_external_user_relation.work_external_user_id desc")
		query.Group("work_external_user_relation.id")
	}
	query.Select(strings.TrimRight(selectStr, ","))
	err = query.Count(&total).
		Offset((request.Page - 1) * request.Size).Limit(request.Size).
		Find(&result).Error

	return result, total, err
	*/

}

// HandelQuerySql 公共查询会员列表sql
func (d *UserRelation) HandelQuerySql(c app.Ctx, request customerReq.ExternalUserListsReq, selectStr string) (query *gorm.DB) {

	userCorpExternalTable := d.dao.GetTable("wx_user_corp_external", request.EnterpriseID)

	query = d.dao.WithContext(c).Table(d.tableName).
		Scopes(func(db *gorm.DB) *gorm.DB {
			db.Joins(" LEFT JOIN " + userCorpExternalTable + " c ON work_external_user_relation.work_external_user_id = c.id").
				Where(map[string]interface{}{
					"work_external_user_relation.enterprise_id":   request.EnterpriseID,
					"work_external_user_relation.deleted_at":      0,
					"work_external_user_relation.relation_status": 0,
				}).
				Where("c.unionid <> '' ")

			if len(request.MemberList) > 0 {
				db.Where("work_external_user_relation.member_id in (?)", request.MemberList)
			}

			// 1【专属】指ks_sa所属顾问是自己的顾客
			// 2【非专属】指所属顾问不是自己的顾客
			if request.ManageType == 1 {
				db.Where("work_external_user_relation.is_key_sa = 1")
				db.Where("work_external_user_relation.user_id in (?)", request.Userids)

			} else if request.ManageType == 2 {
				db.Where("work_external_user_relation.is_key_sa = 0")
			}
			// 企微标签
			if len(request.TagList) > 0 || len(request.ChatTag) > 0 {
				selectStr += "b.external_tag_id,"
				db.Joins(" LEFT JOIN work_external_user_tag b ON work_external_user_relation.id = b.relation_id").Where("b.deleted_at = 0")
			}
			if len(request.TagList) > 0 {
				db.Where(" b.external_tag_id in (?)  ", request.TagList)
			}
			if len(request.ChatTag) > 0 {
				db.Where(" b.external_tag_id in (?)  ", request.ChatTag).Where("work_external_user_relation.member_id = ?", request.MemberId)
			}

			// 印象标签
			if len(request.ImpressionTagList) > 0 {
				db.Joins(" LEFT JOIN impression_tag_relation d ON c.euid = d.external_user_id").
					Where("d.deleted_at = 0").Where(" d.tag_id in (?)  ", request.ImpressionTagList)
			}
			// clt类型
			if len(request.CltType) > 0 {
				db.Where("c.clt_type in (?)", request.CltType)
			}
			// 会员卡等
			if len(request.CardType) > 0 {
				var cardLevel []string
				if len(request.CardType) == 1 {
					switch request.CardType[0] {
					case customerReq.HIGH_GRADE:
						// 优质客 clt类型有值
						db.Where("c.clt_type <> '' ")
					case customerReq.POTENTIAL:
						db.Where("c.card_level = ? and c.amount_total = 0", customerReq.PROSPECT)
					case "O":
						cardLevel = append(cardLevel, "O", "")
						db.Where("c.card_level in ?", cardLevel)
					default:
						db.Where("c.card_level = ?", request.CardType[0])
					}
				} else {
					highGrade := false
					potential := false
					for _, v := range request.CardType {
						if v == customerReq.HIGH_GRADE {
							highGrade = true
						} else if v == customerReq.POTENTIAL {
							potential = true
						} else {
							if v == "O" { // 非会员 或为空
								cardLevel = append(cardLevel, "O", "")
							} else {
								cardLevel = append(cardLevel, v)
							}
						}
					}
					if len(cardLevel) > 0 {
						if !highGrade && !potential {
							db.Where("c.card_level in (?)", cardLevel)
						}
						if highGrade && potential {
							db.Where("c.card_level in (?) or c.clt_type <> '' or (c.card_level = ? and c.amount_total = 0) ", cardLevel, customerReq.PROSPECT)
						}
						if highGrade && !potential {
							db.Where("c.card_level in (?) or c.clt_type <> '' ", cardLevel)
						}
						if potential && !highGrade {
							db.Where("c.card_level in (?) or (c.card_level = ? and c.amount_total = 0)", cardLevel, customerReq.PROSPECT)
						}
					} else {
						if highGrade && potential {
							db.Where("c.clt_type <> '' or (c.card_level = ? and c.amount_total = 0)", customerReq.PROSPECT)
						}
					}
				}
			}
			// 会员卡号code
			if len(request.UnionCodes) > 0 {
				db.Where("c.union_code in (?)", request.UnionCodes)
			}

			// sw顾客类型
			if len(request.CustomerType) > 0 {
				customerType := request.CustomerType
				// vip客
				if helper.InArray(workExternalModel.VipTypeChoice, request.CustomerType) {
					db.Where("c.vip_type = (?)", workExternalModel.VipType)

					item := []int{workExternalModel.VipTypeChoice}
					customerType = helper.DelItems(request.CustomerType, item)
				}
				if len(customerType) > 0 {
					db.Where("c.customer_type in (?)", customerType)
				}
			}
			// sw会员最近购买时间
			if request.LastPurchaseStartTime > 0 && request.LastPurchaseEndTime > 0 {
				db.Where("c.last_purchase_date >= (?) and c.last_purchase_date < (?) ", request.LastPurchaseStartTime, request.LastPurchaseEndTime)
			}
			// sw会员生日月份时间
			if request.BirthdayStart > 0 && request.BirthdayEnd > 0 {
				birthdayArr := helper.GetBirthdayArr(request.BirthdayStart, request.BirthdayEnd)
				db.Where("c.birth_month in (?) ", birthdayArr)
			}

			if request.Name != "" {
				db.Where("work_external_user_relation.remark LIKE ? OR c.name LIKE ?", "%"+request.Name+"%", "%"+request.Name+"%")
			}

			return db
		})

	return
}

func (d *UserRelation) GetOne(c app.Ctx, query interface{}) (workExternalModel.WorkExternalUserRelationModel, error) {
	var result workExternalModel.WorkExternalUserRelationModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query).Find(&result).Error
	return result, err
}

func (d *UserRelation) GetOneUnscoped(c app.Ctx, query interface{}) (workExternalModel.WorkExternalUserRelationModel, error) {
	var result workExternalModel.WorkExternalUserRelationModel
	err := d.dao.WithContext(c).Unscoped().Table(d.tableName).Where(query).Find(&result).Error
	return result, err
}

func (d *UserRelation) CreateData(c app.Ctx, data workExternalModel.WorkExternalUserRelationModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).Create(&data).Error
	return err
}

func (d *UserRelation) CreateInBatchesData(c app.Ctx, data []workExternalModel.WorkExternalUserRelationModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *UserRelation) SaveData(c app.Ctx, data workExternalModel.WorkExternalUserRelationModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).Save(&data).Error
	return err
}

func (d *UserRelation) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]workExternalModel.WorkExternalUserRelationModel, error) {
	var result []workExternalModel.WorkExternalUserRelationModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *UserRelation) GetCustomStoreWorkerByUnionId(c app.Ctx, request openReq.CustomerStoreWorkListReq) ([]workExternalModel.WorkExternalUserRelationModel, int64, error) {
	var result []workExternalModel.WorkExternalUserRelationModel
	var total int64
	offset := (request.Page - 1) * request.Size
	err := d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"enterprise_id":   request.EnterpriseId,
			"unionid":         request.UnionId,
			"relation_status": 0,
		}).
		Group("member_id").
		Count(&total).
		Offset(offset).Limit(request.Size).
		Preload("Members.MemberDepartments.Department").
		Find(&result).Error
	return result, total, err
}

func (d *UserRelation) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Table(d.tableName).Where(query).Model(&workExternalModel.WorkExternalUserRelationModel{}).Updates(data).Error
	return err
}

// GetUnionidsBySearchLists 筛选出有关标签的unionid
func (d *UserRelation) GetUnionidsBySearchLists(c app.Ctx, enterpriseId int64, searchList request.SearchList, unionids []string) (data []customerReq.UserUnionIdAndMemberId, err error) {
	userCorpExternalTable := d.dao.GetTable("wx_user_corp_external", enterpriseId)
	query := d.dao.WithContext(c).Table(d.tableName).Select("work_external_user_relation.unionid,work_external_user_relation.member_id").Scopes(func(db *gorm.DB) *gorm.DB {
		if len(searchList.TagList) > 0 || len(searchList.ImpressionTagList) > 0 || len(searchList.CardType) > 0 || len(searchList.CltType) > 0 || searchList.BirthdayStart > 0 || searchList.LastPurchaseStartTime > 0 || len(searchList.ChatTag) > 0 {
			db.Joins(" LEFT JOIN " + userCorpExternalTable + " c ON work_external_user_relation.work_external_user_id = c.id")

			if len(searchList.CardType) > 0 {
				var cardLevel []string
				if len(searchList.CardType) == 1 {
					switch searchList.CardType[0] {
					case customerReq.HIGH_GRADE:
						// db.Where("c.amount_total > 0")
						// 优质客是 clt标签 有值
						db.Where("c.clt_type <> '' ")
					case customerReq.POTENTIAL:
						db.Where("c.card_level = ? and c.amount_total = 0", customerReq.PROSPECT)
					case "O":
						cardLevel = append(cardLevel, "O", "")
						db.Where("c.card_level in ?", cardLevel)
					default:
						db.Where("c.card_level = ?", searchList.CardType[0])
					}
				} else {
					highGrade := false
					potential := false
					for _, v := range searchList.CardType {
						if v == customerReq.HIGH_GRADE {
							highGrade = true
						} else if v == customerReq.POTENTIAL {
							potential = true
						} else {
							if v == "O" { // 非会员 或为空
								cardLevel = append(cardLevel, "O", "")
							} else {
								cardLevel = append(cardLevel, v)
							}
						}
					}
					if len(cardLevel) > 0 {
						if !highGrade && !potential {
							db.Where("c.card_level in (?)", cardLevel)
						}
						if highGrade && potential {
							// db.Where("c.card_level in (?) or c.amount_total > 0 or (c.card_level = ? and c.amount_total = 0) ", cardLevel, customerReq.PROSPECT)
							db.Where("c.card_level in (?) or c.clt_type <> ''  or (c.card_level = ? and c.amount_total = 0) ", cardLevel, customerReq.PROSPECT)
						}
						if highGrade && !potential {
							// db.Where("c.card_level in (?) or c.amount_total > 0", cardLevel)
							db.Where("c.card_level in (?) or c.clt_type <> '' ", cardLevel)
						}
						if potential && !highGrade {
							db.Where("c.card_level in (?) or (c.card_level = ? and c.amount_total = 0)", cardLevel, customerReq.PROSPECT)
						}
					} else {
						if highGrade && potential {
							// db.Where("c.amount_total > 0 or (c.card_level = ? and c.amount_total = 0)", customerReq.PROSPECT)
							db.Where("c.clt_type <> ''  or (c.card_level = ? and c.amount_total = 0)", customerReq.PROSPECT)
						}
					}
				}
			}
			if len(searchList.CltType) > 0 {
				db.Where("c.clt_type in (?)", searchList.CltType)
			}
			// sw会员最近购买时间
			if searchList.LastPurchaseStartTime > 0 && searchList.LastPurchaseEndTime > 0 {
				db.Where("c.last_purchase_date >= (?) and c.last_purchase_date < (?) ", searchList.LastPurchaseStartTime, searchList.LastPurchaseEndTime)
			}
			// sw会员生日月份时间
			if searchList.BirthdayStart > 0 && searchList.BirthdayEnd > 0 {
				var birthdayArr []string
				for i := searchList.BirthdayStart; i <= searchList.BirthdayEnd; i++ {
					birthdayArr = append(birthdayArr, helper.EncodeToString(strconv.Itoa(i)))
				}
				db.Where("c.birth_month in (?) ", birthdayArr)
			}

			if len(searchList.TagList) > 0 || len(searchList.ChatTag) > 0 {
				db.Joins(" LEFT JOIN work_external_user_tag b ON work_external_user_relation.id = b.relation_id").Where("b.deleted_at = 0")
			}
			if len(searchList.TagList) > 0 {
				db.Where(" b.external_tag_id in (?)  ", searchList.TagList)
			}
			if len(searchList.ChatTag) > 0 {
				db.Where(" b.external_tag_id in (?)  ", searchList.ChatTag)
			}

			if len(searchList.ImpressionTagList) > 0 {
				db.Joins(" LEFT JOIN impression_tag_relation d ON c.euid = d.external_user_id").Where(" d.tag_id in (?)  ", searchList.ImpressionTagList).Where("d.deleted_at = 0")
			}
		}

		// 1【专属】指所属顾问是自己的顾客
		// 2【非专属】指所属顾问不是自己的顾客
		if searchList.ManageType == 1 {
			db.Where("work_external_user_relation.is_key_sa = 1")
		} else if searchList.ManageType == 2 {
			db.Where("work_external_user_relation.is_key_sa <> 1")
		}

		return db
	}).Where(map[string]interface{}{
		"work_external_user_relation.enterprise_id":   enterpriseId,
		"work_external_user_relation.deleted_at":      0,
		"work_external_user_relation.relation_status": 0,
		// "work_external_user_relation.member_deleted_time":   0,
		// "work_external_user_relation.external_deleted_time": 0,
	})
	err = query.Where("work_external_user_relation.unionid in ?", unionids).Find(&data).Error
	return

}

// GetUserTagInfoByUnionIds 获取顾客信息，带标签
func (d *UserRelation) GetUserTagInfoByUnionIds(c app.Ctx, userId string, unionids []string, request planReq.LeaveGroupChatReq) (result []customerReq.ExternalUserTagsInfoListsRes, err error) {
	var data []customerReq.ExternalUserTagsInfoRes
	enterpriseId := request.EnterpriseId
	status := request.Status
	page := request.Page
	size := request.Size
	offset := (page - 1) * size
	userCorpExternalTable := d.dao.GetTable("wx_user_corp_external", enterpriseId)
	// 获取用户信息
	err = d.dao.WithContext(c).Table(userCorpExternalTable).Select("id,avatar,nick,name,unionid,openid as work_external_user_id").Where(map[string]interface{}{
		"unionid": unionids,
	}).Offset(offset).Limit(size).Find(&data).Error
	if err != nil {
		return nil, err
	}
	externalUserInfo := make(map[string]customerReq.ExternalUserTagsInfoListsRes, 0)
	if len(data) == 0 {
		return nil, err
	}
	for _, info := range data {
		externalUserInfo[info.Unionid] = customerReq.ExternalUserTagsInfoListsRes{
			ExternalUserTagsInfoRes: info,
			TgName:                  []string{},
		}
	}
	type unionIdsFriend struct {
		ID      int64    `json:"id"`
		Unionid string   `json:"unionid"`
		TgLists []string `json:"tg_lists"`
	}
	var unionIdsFriends []unionIdsFriend

	err = d.dao.WithContext(c).Table(d.tableName).Select("id,unionid").Scopes(func(db *gorm.DB) *gorm.DB {
		if status == 2 { // 获取非好友关系的用户
			db.Where("relation_status <> 0")
		} else { // 全部的筛选或者是好友关系的筛选
			db.Where("relation_status = 0")
		}
		return db
	}).Where(map[string]interface{}{
		"user_id":    userId,
		"unionid":    unionids,
		"deleted_at": 0,
	}).Offset(offset).Limit(size).Find(&unionIdsFriends).Error
	if err != nil {
		return result, err
	}
	unionIdsFriendsMap := make(map[string]unionIdsFriend, 0)
	if status == 1 || status == 0 { // 好友关系或者全部的查询
		relationFriendUnionIds := make(map[string]int64, 0)
		if len(unionIdsFriends) > 0 {
			var relationIds []int64
			relationFriendsTagLists := make(map[int64]unionIdsFriend, 0)
			for _, item := range unionIdsFriends {
				relationIds = append(relationIds, item.ID)
				relationFriendsTagLists[item.ID] = item
				unionIdsFriendsMap[item.Unionid] = item
			}
			var tagLists []struct {
				RelationId int64  `json:"relation_id"`
				Name       string `json:"name"`
			}
			err = d.dao.WithContext(c).Table("work_external_user_tag").Select("work_external_user_tag.relation_id,tg.name").Joins("LEFT JOIN work_external_tag as tg ON tg.id=work_external_user_tag.external_tag_id").Where(map[string]interface{}{
				"work_external_user_tag.enterprise_id": enterpriseId,
				"relation_id":                          relationIds,
				"work_external_user_tag.deleted_at":    0,
				"tg.deleted_at":                        0,
			}).Find(&tagLists).Error
			if err != nil {
				return result, err
			}
			unionIdsFriendsTagLists := make(map[int64][]string, 0)

			for _, list := range tagLists {
				unionIdsFriendsTagLists[list.RelationId] = append(unionIdsFriendsTagLists[list.RelationId], list.Name)
			}
			for k, tags := range unionIdsFriendsTagLists {
				info := relationFriendsTagLists[k]
				info.TgLists = tags
				relationFriendsTagLists[k] = info
				relationFriendUnionIds[info.Unionid] = info.ID
			}

			for k, value := range relationFriendUnionIds {
				unionIdsFriendsMap[k] = relationFriendsTagLists[value]
			}
		}
	}
	if (status == 2 || status == 1) && len(unionIdsFriends) == 0 { // 如果没有查到数据，则返回空
		return result, err
	}
	for s, res := range externalUserInfo {
		if len(unionIdsFriendsMap) == 0 { // 如果没有内容，则所有用户的用户都是非好友关系的
			res.Status = 2
			result = append(result, res)
		} else {
			if item, ok := unionIdsFriendsMap[s]; !ok { // 非好友关系
				res.Status = 2
				result = append(result, res)
			} else {
				res.Status = 1
				res.TgName = item.TgLists
				result = append(result, res)
			}
		}

	}
	return result, err

	// todo 下面代码是sql查询
	// str := fmt.Sprintf("if(max(IF(work_external_user_relation.user_id='%s',1,0)+IF(work_external_user_relation.relation_status=0,1,0)=2,1,0) as `status`,c.avatar,c.nick,tg.name as tg_name,c.openid as external_user_id,c.unionid", userId)
	// if status == 2 {
	//	str = fmt.Sprintf("min(IF(work_external_user_relation.user_id!='%s',1,0)) as `status`,c.avatar,c.nick,'' as tg_name,c.openid as external_user_id,c.unionid", userId)
	// }
	// query := d.dao.WithContext(c).Table(d.tableName).
	//	Select(str).
	//	Joins(" LEFT JOIN " + userCorpExternalTable + " c ON work_external_user_relation.work_external_user_id = c.id").
	//	Joins(" LEFT JOIN work_external_user_tag b ON work_external_user_relation.id = b.relation_id").
	//	Joins(" LEFT JOIN work_external_tag as tg ON tg.id=b.external_tag_id").
	//	Where(map[string]interface{}{
	//		"work_external_user_relation.enterprise_id": enterpriseId,
	//		"work_external_user_relation.deleted_at":    0,
	//		"work_external_user_relation.unionid":       unionids,
	//	})
	// if status == 0 {
	//	err = query.Group("c.unionid").Find(&data).Error
	// } else {
	//	err = query.Group("c.unionid").Having("status=1").Find(&data).Error
	// }

	// return

}

// GetUserSAByExternalUserIds 根据work_external_user_id 查询外部联系人的SA
func (d *UserRelation) GetUserSAByExternalUserIds(c app.Ctx, enterpriseId, workExternalUserId int64) (info []workExternalModel.UserSaData, err error) {

	err = d.dao.WithContext(c).Table(d.tableName).
		Select("work_external_user_relation.member_id, work_external_user_relation.unionid,   work_external_user_relation.user_id as sa_user_id, work_external_user_relation.is_key_sa, "+
			" m.userid, m.qr_code, m.`name` as sa_user_name").
		Joins(" LEFT JOIN members m ON work_external_user_relation.member_id = m.id").
		Where(map[string]interface{}{
			"work_external_user_relation.enterprise_id": enterpriseId,
			"work_external_user_relation.deleted_at":    0,
			"m.status":                                  1,
			"m.deleted_at":                              0,
		}).
		Where("work_external_user_relation.work_external_user_id = ?", workExternalUserId).
		Order("work_external_user_relation.is_key_sa desc").
		Find(&info).Error

	return

}

// GetUserSAListByExternalUserIds 根据work_external_user_id 查询外部联系人的SA列表
func (d *UserRelation) GetUserSAListByExternalUserIds(c app.Ctx, enterpriseId, workExternalUserId int64, memberIds []int64) (list []workExternalModel.UserSaListData, err error) {

	// 外部联系人全部的导购关系
	query := d.dao.WithContext(c).Table(d.tableName).
		Select("work_external_user_relation.member_id, work_external_user_relation.unionid,   work_external_user_relation.user_id as sa_user_id, work_external_user_relation.is_key_sa, work_external_user_relation.createtime as created_at, " +
			" b.userid, b.thumb_avatar, b.qr_code, b.`name` as sa_user_name," +
			"d.third_id, d.`name` as department_name").
		Joins(" LEFT JOIN members b ON work_external_user_relation.member_id = b.id").
		Joins(" LEFT JOIN member_departments c ON b.id = c.member_id").
		Joins(" LEFT JOIN departments d ON c.department_id = d.id").
		Where(map[string]interface{}{
			"work_external_user_relation.enterprise_id":         enterpriseId,
			"work_external_user_relation.deleted_at":            0,
			"work_external_user_relation.relation_status":       0,
			"work_external_user_relation.work_external_user_id": workExternalUserId,
			"b.status":     1,
			"b.deleted_at": 0,
			"c.deleted_at": 0,
			"d.deleted_at": 0,
		})

	if len(memberIds) > 0 {
		query.Where("b.id in (?)", memberIds)
	}

	err = query.Order("work_external_user_relation.is_key_sa desc").
		Order("work_external_user_relation.createtime desc").
		Group("b.id").
		Find(&list).Error

	return

}

func (d *UserRelation) GetUserSARelation(c app.Ctx, enterpriseId, memberId, workExternalUserID int64) (info workExternalModel.UserSaInfo, err error) {

	err = d.dao.WithContext(c).Table(d.tableName).
		Select("work_external_user_relation.id,work_external_user_relation.member_id,work_external_user_relation.user_id as sa_user_id, work_external_user_relation.remark").
		Joins(" LEFT JOIN members b ON work_external_user_relation.member_id = b.id").
		Where(map[string]interface{}{
			"work_external_user_relation.enterprise_id":         enterpriseId,
			"work_external_user_relation.work_external_user_id": workExternalUserID,
			"work_external_user_relation.member_id":             memberId,
			"work_external_user_relation.relation_status":       0,
			"work_external_user_relation.deleted_at":            0,
		}).Find(&info).Error

	return

}

// GetUserSARelations 查询顾客和导购的关系
func (d *UserRelation) GetUserSARelations(c app.Ctx, enterpriseId, memberId int64, workExternalUserIDs []int64) (list []workExternalModel.UserSaInfo, err error) {

	err = d.dao.WithContext(c).Table(d.tableName).
		Select("work_external_user_relation.id,work_external_user_relation.member_id,work_external_user_relation.work_external_user_id, work_external_user_relation.user_id as sa_user_id").
		Joins(" LEFT JOIN members b ON work_external_user_relation.member_id = b.id").
		Where(map[string]interface{}{
			"work_external_user_relation.enterprise_id":   enterpriseId,
			"work_external_user_relation.member_id":       memberId,
			"work_external_user_relation.relation_status": 0,
			"work_external_user_relation.deleted_at":      0,
		}).
		Where("work_external_user_relation.work_external_user_id in ?", workExternalUserIDs).
		Find(&list).Error

	return

}

// GetUserKeySARelations 查询顾客的key_sa
func (d *UserRelation) GetUserKeySARelations(c app.Ctx, enterpriseId int64, workExternalUserIDs []int64) (list []workExternalModel.UserSaInfo, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Select("id,member_id,work_external_user_id, user_id as sa_user_id,is_key_sa").
		Where(map[string]interface{}{
			"enterprise_id":   enterpriseId,
			"relation_status": 0,
			"deleted_at":      0,
			"is_key_sa":       1,
		}).
		Where("work_external_user_id in ?", workExternalUserIDs).
		Find(&list).Error

	return
}

// GetUserRelationIds 查询顾客的关系id
func (d *UserRelation) GetUserRelationIds(c app.Ctx, enterpriseId int64, workExternalUserIDs []int64) ([]int64, error) {
	var relationIds []int64
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("id").
		Where(map[string]interface{}{
			"enterprise_id":   enterpriseId,
			"relation_status": 0,
		}).
		Where("work_external_user_id in ?", workExternalUserIDs).
		Scan(&relationIds).Error

	return relationIds, err
}

// GetAddFriendsCount 获取指定时间的新增好友数
func (d *UserRelation) GetAddFriendsCount(c app.Ctx, enterpriseId int64, startTime int64, endTime int64, memberIds []int64) (total int64, err error) {
	whereMap := map[string]interface{}{
		"enterprise_id":   enterpriseId,
		"relation_status": 0,
	}
	if len(memberIds) > 0 {
		whereMap["member_id"] = memberIds
	}
	err = d.dao.WithContext(c).Table(d.tableName).Where(whereMap).Where("created_at BETWEEN ? AND ?", startTime, endTime).Count(&total).Error
	return
}

// GetALLFriendsCount 获取累计的好友数
func (d *UserRelation) GetALLFriendsCount(c app.Ctx, enterpriseId, endTime int64, memberIds []int64) (total int64, err error) {
	whereMap := map[string]interface{}{
		"enterprise_id":   enterpriseId,
		"relation_status": 0,
	}
	if len(memberIds) > 0 {
		whereMap["member_id"] = memberIds
	}
	err = d.dao.WithContext(c).Table(d.tableName).Where(whereMap).Where("created_at <= ? ", endTime).Count(&total).Error
	return
}

// GetLossFriendsCount 获取指定时间的流失好友数
func (d *UserRelation) GetLossFriendsCount(c app.Ctx, enterpriseId int64, startTime int64, endTime int64, memberIds []int64) (total int64, err error) {
	whereMap := map[string]interface{}{
		"enterprise_id": enterpriseId,
	}
	if len(memberIds) > 0 {
		whereMap["member_id"] = memberIds
	}
	err = d.dao.WithContext(c).Table(d.tableName).Where(whereMap).Where("updated_at BETWEEN ? AND ? AND relation_status!= 0", startTime, endTime).Count(&total).Error
	return
}

// GetGroupMemberLossFriendsCount  获取指定时间的流失好友数
func (d *UserRelation) GetGroupMemberLossFriendsCount(c app.Ctx, enterpriseId int64, startTime int64, endTime int64) (res []models.BehaviorList, err error) {
	err = d.dao.WithContext(c).
		Table(d.tableName).
		Where("enterprise_id = ?", enterpriseId).
		Where("updated_at BETWEEN ? AND ? AND relation_status > 0", startTime, endTime).
		Group("member_id").
		Select("member_id, count(1) as total").
		Scan(&res).Error
	return
}

func (d *UserRelation) GetUserListForUnionCode(c app.Ctx, request customerReq.CustomerListsReq) ([]workExternalModel.ExternalUserList, error) {

	var (
		result []workExternalModel.ExternalUserList
		err    error
	)
	userCorpExternalTable := d.dao.GetTable("wx_user_corp_external", request.EnterpriseID)
	selectStr := "work_external_user_relation.id,work_external_user_relation.enterprise_id,work_external_user_relation.remark,work_external_user_relation.member_id,work_external_user_relation.work_external_user_id,work_external_user_relation.user_id as sa_user_id,work_external_user_relation.is_key_sa,work_external_user_relation.created_at," +
		"c.unionid,c.openid,c.euid,c.nick,c.name,c.avatar,c.birthday,c.union_code,c.mobile,c.card_level,c.clt_type,c.amount_total,"

	query := d.dao.WithContext(c).Table(d.tableName).
		Joins(" LEFT JOIN "+userCorpExternalTable+" c ON work_external_user_relation.work_external_user_id = c.id").
		Where(map[string]interface{}{
			"work_external_user_relation.enterprise_id":   request.EnterpriseID,
			"work_external_user_relation.deleted_at":      0,
			"work_external_user_relation.relation_status": 0,
		}).
		Where("c.unionid <> '' ").
		Where("work_external_user_relation.member_id in (?)", request.MemberList).
		Where("c.union_code in (?)", request.UnionCodes)

	query.Select(strings.TrimRight(selectStr, ","))
	err = query.Find(&result).Error

	return result, err

}

// SwCustomers 顾客榜单查询 TODO
func (d *UserRelation) SwCustomers(c app.Ctx, request memberReq.SwCustomersListReq, membersID []int64) ([]memberReq.SwCustomersListReqRes, int64, error) {

	var (
		total  int64
		result []memberReq.SwCustomersListReqRes
		err    error
	)

	offset := (request.Page - 1) * request.Size
	startTime := request.FinancialTime[0]
	endTime := request.FinancialTime[1]

	userCorpExternalTable := d.dao.GetTable("wx_user_corp_external", request.EnterpriseId)
	selectStr := "count(CASE WHEN c.customer_type = 1 THEN 1 ELSE NULL END) as new_customers," +
		"count(CASE WHEN c.customer_type = 2 THEN 1 ELSE NULL END) as old_customers," +
		"count(CASE WHEN c.customer_type = 3 THEN 1 ELSE NULL END) as sleep_customers," +
		"count(CASE WHEN c.card_level <> '' THEN 1 ELSE NULL END) as vip_customers," +
		"work_external_user_relation.member_id"

	query := d.dao.WithContext(c).Table(d.tableName).
		Joins(" LEFT JOIN "+userCorpExternalTable+" c ON work_external_user_relation.work_external_user_id = c.id").
		Where(map[string]interface{}{
			"work_external_user_relation.enterprise_id":   request.EnterpriseId,
			"work_external_user_relation.deleted_at":      0,
			"c.deleted_at":                                0,
			"work_external_user_relation.relation_status": 0,
		}).
		Where("c.unionid <> '' ").
		Where("c.customer_type_updated_at BETWEEN ? AND ?", startTime, endTime).
		Group("work_external_user_relation.member_id")

	if len(membersID) > 0 {
		query.Where("work_external_user_relation.member_id in ?", membersID)
		query.Count(&total)
	}

	query.Select(selectStr).
		Order("new_customers desc").
		Order("old_customers desc").
		Order("vip_customers desc").
		Order("work_external_user_relation.member_id asc")

	err = query.Offset(offset).Limit(request.Size).Find(&result).Error

	return result, total, err
}

func (d *UserRelation) GetUserRelationByUnionIdUserId(c app.Ctx, enterpriseId int64, unionId, userId string) (workExternalModel.WorkExternalUserRelationModel, error) {
	var result workExternalModel.WorkExternalUserRelationModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"unionid":       unionId,
		"user_id":       userId,
		"deleted_at":    0,
	}).Find(&result).Error
	return result, err
}

// GetInfoByUnionIds 查询外部联系人关联员工
func (d *UserRelation) GetInfoByUnionIds(c app.Ctx, enterpriseId int64, unionIds []string) ([]workExternalModel.WorkExternalUserRelationModel, error) {
	var result []workExternalModel.WorkExternalUserRelationModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"enterprise_id":   enterpriseId,
			"unionid":         unionIds,
			"relation_status": 0,
			"deleted_at":      0,
		}).
		Find(&result).Error

	return result, err

}

// GetInfoByUnionIdsOrCustomers 通过UnionIds或顾客标识获取关系
func (d *UserRelation) GetInfoByUnionIdsOrCustomers(c app.Ctx, enterpriseId int64, unionIds []string, externalUserIds []int64) ([]workExternalModel.WorkExternalUserRelationModel, error) {
	var result []workExternalModel.WorkExternalUserRelationModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id", enterpriseId).
		Where("relation_status", 0).
		Where("deleted_at", 0).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if len(unionIds) == 0 && len(externalUserIds) > 0 {
				db.Where("work_external_user_id in (?)", externalUserIds)
			}
			if len(unionIds) > 0 && len(externalUserIds) == 0 {
				db.Where("unionid in (?)", unionIds)
			}
			if len(unionIds) > 0 && len(externalUserIds) > 0 {
				db.Where("unionid in (?) or work_external_user_id in (?)", unionIds, externalUserIds)
			}
			return db
		}).
		Find(&result).Error

	return result, err
}

// GetAllByUnionIdUserIdMap 通过 unionid,user_id组合的切片获取所有数据
func (d *UserRelation) GetAllByUnionIdUserIdMap(c app.Ctx, enterpriseId int64, unionIdUserIdSlice [][]any) ([]workExternalModel.WorkExternalUserRelationModel, error) {
	var result []workExternalModel.WorkExternalUserRelationModel
	err := d.dao.WithContext(c).Table(d.tableName).Where("enterprise_id = ? and (unionid,user_id) in (?)", enterpriseId, unionIdUserIdSlice).Find(&result).Error
	return result, err
}

// GetManyPaginationById 根据ID分页获取数据
func (d *UserRelation) GetManyPaginationById(c ctx.Context, query interface{}, args ...interface{}) ([]workExternalModel.WorkExternalUserRelationModel, error) {
	var result []workExternalModel.WorkExternalUserRelationModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query, args...).Order("id asc").Limit(1000).Find(&result).Error
	return result, err
}

// GetManyUnionIds 获取unionid
func (d *UserRelation) GetManyUnionIds(c app.Ctx, query interface{}, args ...interface{}) ([]string, error) {
	var result []string
	err := d.dao.WithContext(c).Table(d.tableName).Model(&workExternalModel.WorkExternalUserRelationModel{}).Where(query, args...).Pluck("unionid", &result).Error
	return result, err
}

func (d *UserRelation) BatchDeleteExternalRelationByMemberIds(c app.Ctx, enterpriseID int64, lastUpdateAt int64, memberIds []int64) error {
	err := d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id=? AND updated_at<? AND member_id in ?", enterpriseID, lastUpdateAt, memberIds).
		UpdateColumns(map[string]interface{}{
			"relation_status": workExternalModel.WorkUserRelationStatusDeleteEachOther,
			"updated_at":      time.Now().Unix(),
		}).Error
	return err
}

func (d *UserRelation) GetRelInfoByUnionIdAndMemberIdsWithOrder(c app.Ctx, enterpriseId int64, unionId string, memberIds []int64) ([]workExternalModel.WorkExternalUserRelationModel, error) {
	var result []workExternalModel.WorkExternalUserRelationModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id=?", enterpriseId).
		Where("unionid = ?", unionId).
		Where("member_id in (?)", memberIds).
		Order("is_key_sa desc,createtime desc").
		Find(&result).Error
	return result, err
}

func (d *UserRelation) BatchUpdateByUnionIds(c context.Context, enterpriseId int64, mappings []workExternalModel.UnionidMapModel, nowUnix int64) (int64, error) {
	if len(mappings) == 0 {
		return 0, nil
	}

	// 获取当前批次的最小和最大ID
	minID := mappings[0].ID
	maxID := mappings[len(mappings)-1].ID

	query := `
       UPDATE micro_scrm.` + d.tableName + ` AS t1
        INNER JOIN ` + fmt.Sprintf("%s%d.%s", app.Cfg().DB.DBPrefix, enterpriseId, "unionid_map") + ` AS t2 ON t1.unionid = t2.old_unionid
        SET t1.unionid = t2.new_unionid, t1.updated_at = ?
        WHERE t1.enterprise_id = ?
        AND t2.id >= ? AND t2.id <= ?`

	db := d.dao.WithContext(c).Table(d.tableName).Exec(query, nowUnix, enterpriseId, minID, maxID)
	return db.RowsAffected, db.Error
}

func (d *UserRelation) GetManyLimit(c app.Ctx, enterpriseId int64, id int64, limit int) ([]workExternalModel.WorkExternalUserRelationModel, error) {
	var result []workExternalModel.WorkExternalUserRelationModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"enterprise_id":   enterpriseId,
			"relation_status": 0,
		}).
		Where("unionid != ?", "").
		Where("id > ?", id).
		Order("id asc").
		Limit(limit).
		Find(&result).Error
	return result, err
}

func (d *UserRelation) GetRemarkByUnionIdAndMemberId(c app.Ctx, enterpriseId int64, unionId string, memberId int64) (string, error) {
	var remark string
	err := d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id=?", enterpriseId).
		Where("unionid = ?", unionId).
		Where("member_id = ?", memberId).
		Pluck("remark", &remark).Error
	return remark, err
}

// GetCountMapByStateSlice   获取累计的好友数
func (d *UserRelation) GetCountMapByStateSlice(c app.Ctx, enterpriseId int64, stateSlice []string) (map[string]int64, error) {
	var resList = make([]drainageContactResp.FriendsNumberResp, 0)
	err := d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id":   enterpriseId,
		"relation_status": 0,
		"state":           stateSlice,
	}).Select("count(1) as total,state").Group("state").Scan(&resList).Error
	var resMap = make(map[string]int64)
	for _, resp := range resList {
		resMap[resp.State] = resp.Total
	}
	return resMap, err
}
