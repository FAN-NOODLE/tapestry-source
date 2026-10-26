package work_external_user

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"gorm.io/gorm"
	planReq "micro-scrm/internal/app/http/request/plan"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/service"
	"micro-scrm/internal/pkg/db"
)

func (w *WorkExternalUserService) GetCustomerInfo(enterpriseID int64) (data []micro_scrm_enterprise.WxUserCorpExternalModel, err error) {

	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.WxUserCorpExternalModel{}, enterpriseID)).
		Scopes(service.ScopeCond_EQDelete).
		Where("euid IN ?", func(w *WorkExternalUserService) (data1 []int64) {
			w.Scopes(model.BranchTable(&micro_scrm_enterprise.EnterpriseUserModel{}, enterpriseID)).
				Where("grade = ?", "黑钻卡").Select("id").Find(&data1)

			return
		}(w)).
		Find(&data).Error

	//w.Raw("SELECT * FROM micro_scrm_enterprise_3314240094228544.wx_user_corp_external W INNER JOIN micro_scrm_enterprise_3314240094228544.enterprise_user E ON E.id = W.euid WHERE E.grade = \"黑钻卡\"").Scan(&data)

	return
}

func (w *WorkExternalUserService) GetAllCustomer(enterpriseID int64) (data []micro_scrm_enterprise.WxUserCorpExternalModel, err error) {

	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.WxUserCorpExternalModel{}, enterpriseID)).
		Scopes(service.ScopeCond_EQDelete).
		Find(&data).Error

	return
}

// GetWorkExternalUserInfoByWorkExternalUserID 根据企业微信中的外部联系人ID，获取外部联系人的信息
func (w *WorkExternalUserService) GetWorkExternalUserInfoByWorkExternalUserID(enterpriseID int64, workExternalUserID string) (data micro_scrm_enterprise.WxUserCorpExternalModel, err error) {
	err = w.Scopes(model.BranchTable(&data, enterpriseID)).
		Scopes(service.ScopeCond_EQDelete).
		Where("openid = ? ", workExternalUserID).
		First(&data).Error

	return
}

// GetUserIDByExternalUserIDList 根据系统的external_user_id集合，获取聚合ID
func (w *WorkExternalUserService) GetUserIDByExternalUserIDList(enterpriseID int64, externalUserIDList []int64) (userID []int64, err error) {
	m := micro_scrm_enterprise.WxUserCorpExternalModel{}

	if err = w.Model(m).Scopes(model.BranchTable(&m, enterpriseID), service.ScopeCond_EQDelete).
		Where("id IN ?", externalUserIDList).
		Pluck("euid", &userID).Error; err != nil {
		return
	}

	return
}

// GetExternalUserIDByExternalUserIDList 获取顾客的企微userid, work_external_work_id
func (w *WorkExternalUserService) GetExternalUserIDByExternalUserIDList(enterpriseID int64, externalUserIDList []int64) (userID []int64, err error) {
	m := micro_scrm_enterprise.WxUserCorpExternalModel{}

	if err = w.Model(m).Scopes(model.BranchTable(&m, enterpriseID), service.ScopeCond_EQDelete).
		Where("id IN ?", externalUserIDList).
		Pluck("openid", &userID).Error; err != nil {
		return
	}

	return
}

// GetWxUserCorpExternalInfoByUnionID 根据unionID获取用户的信息
func (w *WorkExternalUserService) GetWxUserCorpExternalInfoByUnionID(enterpriseID int64, unionID string) (data micro_scrm_enterprise.WxUserCorpExternalModel, err error) {
	err = w.Scopes(model.BranchTable(&data, enterpriseID), service.ScopeCond_EQDelete).Where("unionid = ?", unionID).First(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return data, nil
	}

	return
}

func (w *WorkExternalUserService) UpdateuserCorpExternalbaCodeByunionId(enterpriseID int64, unionID, memberCode, adminBa, gradeCodeName string) (err error) {
	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.WxUserCorpExternalModel{}, enterpriseID)).
		Where("unionid=?", unionID).Updates(map[string]interface{}{
		"union_code": memberCode,
		"ba_code":    adminBa,
		"grade":      gradeCodeName,
	}).Error
	return
}

// GetUserCorpInfoBYMemberIdAndUnionId 获取外部联系人的用户信息
func (w *WorkExternalUserService) GetUserCorpInfoBYMemberIdAndUnionId(enterpriseID int64, unionID []string, memberIds []int64) (data []micro_scrm_enterprise.WxUserCorpExternalModel, err error) {
	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.WxUserCorpExternalModel{}, enterpriseID)).
		Where("unionid in ?", unionID).Where("member_id in ? ", memberIds).Group("member_id").Find(&data).Error
	return data, err
}

// GetUserCorpInfoBYMemberIdAndBirth  获取外部联系人的用户信息
func (w *WorkExternalUserService) GetUserCorpInfoBYMemberIdAndBirth(enterpriseID int64, birth_month, birth_day int, memberIds []int64) (data []planReq.BirthdayUserInfo, err error) {
	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.WxUserCorpExternalModel{}, enterpriseID)).
		Where(map[string]interface{}{
			"wx_user_corp_external.birth_month": birth_month,
			"wx_user_corp_external.birth_day":   birth_day,
			"wr.enterprise_id":                  enterpriseID,
			"wr.member_id":                      memberIds,
			"wr.relation_status":                0,
			"wr.is_key_sa":                      1,
		}).Joins("LEFT JOIN work_external_user_relation AS wr ON wr.unionid=wx_user_corp_external.unionid").
		Select("wx_user_corp_external.unionid,wx_user_corp_external.openid,wr.member_id").
		Group("wr.`member_id`").Find(&data).Error
	return data, err
}

// UpdateUserCorpExternalByEuId 更新字段
func (w *WorkExternalUserService) UpdateUserCorpExternalByEuId(enterpriseID, euId int64, unionCode, birthday, userName, mobile string, birthYear, birthMonth, birthDay int, cardLevel, cltType string, amountTotal float32, customerCode string) (err error) {
	userMap := make(map[string]interface{})
	userMap["union_code"] = unionCode
	//userMap["name"] = userName
	//userMap["mobile"] = mobile
	userMap["card_level"] = cardLevel
	userMap["clt_type"] = cltType
	userMap["amount_total"] = amountTotal * 100
	userMap["customer_code"] = customerCode
	if birthday != "" {
		userMap["birthday"] = birthday
		userMap["birth_year"] = birthYear
		userMap["birth_month"] = birthMonth
		userMap["birth_day"] = birthDay
	}
	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.WxUserCorpExternalModel{}, enterpriseID)).
		Where("euid = ?", euId).
		Updates(userMap).Error
	return
}

// CHHandleExternalUserLog 处理数据到clickHose
func (w *WorkExternalUserService) CHHandleExternalUserLog(c ctx.Context, EnterpriseID int64, Unionid, Birthday, UnionCode, Name, Mobile string) error {
	log.Debug(c).Msgf("====================同步用户信息日志=======================")
	ckdb := db.ClickHouse()
	var resId int64
	sql := fmt.Sprintf("select id  from micro_scrm_enterprise_%d.wx_user_corp_external  WHERE unionid = ?", EnterpriseID)
	err := ckdb.DB.Get(c, &resId, sql, Unionid)
	if err != nil {
		log.Error(c).Msgf("查询数据失败 err: %+v", err)
	}
	if resId > 0 {
		// update
		updt := fmt.Sprintf("ALTER TABLE  micro_scrm_enterprise_%d.wx_user_corp_external_local update  union_code ='%s',name ='%s' WHERE unionid = '%s'", EnterpriseID, UnionCode, Name, Unionid)

		err := ckdb.Tx(c, func(tx *sqlplus.Tx) error {
			_, err := tx.Exec(c, updt)

			return err
		})
		if err != nil {
			log.Error(c).Msgf("更新数据失败 err: %+v", err)
		}

	} else {
		// insert
		data := micro_scrm_enterprise.WxUserCorpExternalModelInfo{}
		err := w.Scopes(model.BranchTable(&micro_scrm_enterprise.WxUserCorpExternalModel{}, EnterpriseID)).
			Select("`id`, `appid`, `openid`, `unionid`, `nick`, `name`, `position`, `avatar`, `birthday`, `union_code`, `point`, `grade`, `store_code`, `ba_code`, `corp_name`, `corp_full_name`, `type`, `sex`, `euid`, `external_profile`, `created_at`, `updated_at`, `deleted_at`").
			Where("unionid = ?", Unionid).First(&data).Error
		if err != nil {
			log.Error(c).Msgf("查询数据失败 unionID : %+v, err: %+v", Unionid, err)

		}

		inst := fmt.Sprintf("insert into micro_scrm_enterprise_%d.wx_user_corp_external_local (`id`, `appid`, `openid`, `unionid`, `nick`, `name`, `position`, `avatar`, `birthday`, `union_code`, `point`, `grade`, `store_code`, `ba_code`, `corp_name`, `corp_full_name`, `type`, `sex`, `euid`, `external_profile`, `created_at`, `updated_at`, `deleted_at`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", EnterpriseID)
		err = ckdb.Tx(c, func(tx *sqlplus.Tx) error {
			_, err = tx.Exec(c, inst, data.ID, data.Appid, data.Openid, data.Unionid, data.Nick, data.Name, data.Position, data.Avatar, data.Birthday, data.UnionCode, data.Point, data.Grade, data.StoreCode, data.BaCode, data.CorpName, data.CorpFullName, int8(data.Type), int8(data.Sex), data.Euid, data.ExternalProfile, data.CreatedAt, data.UpdatedAt, data.DeletedAt)
			return err
		})

		if err != nil {
			log.Error(c).Msgf("写入数据失败 err: %+v", err)
		}

	}
	return nil
}
