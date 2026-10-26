package tagsDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	tagsReq "micro-scrm/internal/app/http/request/tags"
	crmSftpFileModels "micro-scrm/internal/app/models/crm_sftp_file"
	"micro-scrm/internal/app/models/micro_scrm"
	workExternalModel "micro-scrm/internal/app/models/micro_scrm"
	tagModels "micro-scrm/internal/app/models/tag"
	workModels "micro-scrm/internal/app/models/work"
	"micro-scrm/internal/pkg/db/gormx"
	"time"
)

type WorkTags struct {
	dao dao.Dao
}

func NewWorkTagsDao() *WorkTags {
	return &WorkTags{
		dao: dao.Dao{DB: gormx.GetInstance()},
	}
}

// WorkTagsGroupInfoByName 根据标签组名称判断是否存在
func (t WorkTags) WorkTagsGroupInfoByName(c app.Ctx, enterpriseId int64, name string) (tagsGroupInfo workModels.WorkExternalTagGroupModel, err error) {
	err = t.dao.WithContext(c).Model(&workModels.WorkExternalTagGroupModel{}).
		Where("enterprise_id=?", enterpriseId).
		Where("group_name=? and deleted_at=0", name).
		Find(&tagsGroupInfo).Error
	return

}

// DelWorkTagsGroup 删除企微标签组
func (t WorkTags) DelWorkTagsGroup(c app.Ctx, enterpriseId int64, id int64) error {
	return t.dao.WithContext(c).Model(&workModels.WorkExternalTagGroupModel{}).Delete("id=? and enterpriseId=?", id, enterpriseId).Error
}

// GetWorkTagsByName 根据标签名称查询是否有已设置的标签
func (t WorkTags) GetWorkTagsByName(c app.Ctx, enterpriseId int64, groupId int64, nameMap []string) (tagsList []workModels.WorkExternalTagModel, err error) {
	err = t.dao.WithContext(c).Model(&workModels.WorkExternalTagModel{}).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"group_id":      groupId,
			"name":          nameMap,
		}).Where("deleted_at=0").Find(&tagsList).Error
	return
}

// GetWorkTagsByThirdIds 根据第三方标签id查询已设置的标签
func (t WorkTags) GetWorkTagsByThirdIds(c app.Ctx, enterpriseId int64, groupId int64, thirdIds []string) (tagsList []workModels.WorkExternalTagModel, err error) {
	err = t.dao.WithContext(c).Model(&workModels.WorkExternalTagModel{}).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"group_id":      groupId,
			"third_tag_id":  thirdIds,
		}).Where("deleted_at=0").Find(&tagsList).Error
	return
}

func (t WorkTags) UpdateWorkTagsThirdByName(c app.Ctx, enterpriseId int64, groupId int64, Name string, ThirdId string) error {
	return t.dao.WithContext(c).Model(&workModels.WorkExternalTagModel{}).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"group_id":      groupId,
			"name":          Name,
		}).Where("deleted_at=0").Update("third_tag_id", ThirdId).Error

}

func (t WorkTags) UpdateWorkTagsNameByThird(c app.Ctx, enterpriseId int64, groupId int64, ThirdId string, Name string) error {
	return t.dao.WithContext(c).Model(&workModels.WorkExternalTagModel{}).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"group_id":      groupId,
			"third_tag_id":  ThirdId,
		}).Where("deleted_at=0").Update("name", Name).Error
}

func (t WorkTags) InsertWorkTags(c app.Ctx, InsertTagsInfo []workModels.WorkExternalTagModel) error {
	return t.dao.WithContext(c).Model(&workModels.WorkExternalTagModel{}).Create(&InsertTagsInfo).Error
}

func (t WorkTags) InsertWorkTagGroup(c app.Ctx, model workModels.WorkExternalTagGroupModel) (id int64, err error) {
	err = t.dao.WithContext(c).Model(&workModels.WorkExternalTagGroupModel{}).Create(&model).Error
	id = model.ID
	return id, err
}

func (t WorkTags) GetTagInfoByThirdID(c app.Ctx, enterpriseID int64, thirdTagID string) (info workModels.WorkExternalTagModel, err error) {
	err = t.dao.WithContext(c).Model(&workModels.WorkExternalTagModel{}).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseID,
			"third_tag_id":  thirdTagID,
			"deleted_at":    0,
		}).First(&info).Error
	return
}

// UserTagsAndTagGroup 客户标签
func (t WorkTags) UserTagsAndTagGroup(c app.Ctx, enterpriseID int64, workExternalUserId int64, memberIds []int64) (info []tagsReq.UserTagsRes, total int64, err error) {

	// var weurRes workExternalModel.WorkExternalUserRelationModel
	type RelationIds struct {
		Id int64 `json:"Id"`
	}
	var weurIds []RelationIds
	err = t.dao.WithContext(c).Model(microScrmModel.WorkExternalUserRelationModel{}).
		// relationCkTable := fmt.Sprintf("micro_scrm_ck.work_external_user_relation")
		// err = gormx.GetCkInstance().WithContext(c).Table(relationCkTable).
		Select("work_external_user_relation.id").
		Where(map[string]interface{}{
			"work_external_user_relation.enterprise_id":         enterpriseID,
			"work_external_user_relation.work_external_user_id": workExternalUserId,
			"work_external_user_relation.member_id":             memberIds,
			"work_external_user_relation.deleted_at":            0,
		}).Find(&weurIds).Error
	if err != nil {
		return nil, total, err
	}
	var relationIds []int64
	for _, v := range weurIds {
		relationIds = append(relationIds, v.Id)
	}
	// externalUserTagCk := fmt.Sprintf("micro_scrm_ck.work_external_user_tag")
	query := t.dao.WithContext(c).Model(tagModels.WorkExternalUserTagModel{}).
		// query := gormx.GetCkInstance().WithContext(c).Table(externalUserTagCk).
		Select("work_external_user_tag.external_tag_id,wet.id,wet.group_id,wet.tag_id,wet.third_tag_description,wet.`name`,wet.type,wetg.id,wetg.group_name,wetg.is_third_tag ").
		Joins("LEFT JOIN micro_scrm.work_external_tag wet ON work_external_user_tag.external_tag_id = wet.id").
		Joins("LEFT JOIN micro_scrm.work_external_tag_group wetg ON wet.group_id = wetg.id ").
		Where(map[string]interface{}{
			"work_external_user_tag.enterprise_id": enterpriseID,
			"work_external_user_tag.relation_id":   relationIds,
			"work_external_user_tag.deleted_at":    0,
		})

	err = query.Count(&total).Group("wet.tag_id").Order("wet.id ASC").Find(&info).Error

	return info, total, err

}

func (t WorkTags) GetExternalIdByUnionId(c app.Ctx, enterpriseID int64, unionID string) (info workExternalModel.WxUserCorpExternalModel, err error) {
	err = t.dao.WithContext(c).Scopes(t.dao.BranchTable(workExternalModel.WxUserCorpExternalModel{}.TableName(), enterpriseID)).
		Where("unionid=?", unionID).
		Where("deleted_at=0").
		Find(&info).Error
	return
}

func (t WorkTags) AddExternalUserTag(c app.Ctx, info microScrmModel.WorkExternalUserTagModel) error {
	return t.dao.WithContext(c).Model(&microScrmModel.WorkExternalUserTagModel{}).Create(&info).Error
}

func (t WorkTags) BatchCreate(c app.Ctx, data []microScrmModel.WorkExternalUserTagModel) error {
	return t.dao.WithContext(c).Model(&microScrmModel.WorkExternalUserTagModel{}).CreateInBatches(&data, len(data)).Error
}

// DelExternalUserTag 外部联系人的标签软删除
func (t WorkTags) DelExternalUserTag(c app.Ctx, enterpriseID, tagID, relationID int64) error {
	return t.dao.WithContext(c).Model(&microScrmModel.WorkExternalUserTagModel{}).Where(map[string]interface{}{
		"enterprise_id":   enterpriseID,
		"external_tag_id": tagID,
		"relation_id":     relationID,
	}).Update("deleted_at", time.Now().Unix()).Error
}

// DelExternalUserTags 外部联系人的标签软删除
func (t WorkTags) DelExternalUserTags(c app.Ctx, enterpriseID, relationID int64, tagIds []int64) error {
	return t.dao.WithContext(c).Model(&microScrmModel.WorkExternalUserTagModel{}).Where(map[string]interface{}{
		"enterprise_id":   enterpriseID,
		"external_tag_id": tagIds,
		"relation_id":     relationID,
	}).Update("deleted_at", time.Now().Unix()).Error
}

func (t WorkTags) TagsAuth(c app.Ctx, enterpriseID int64) (total int64, err error) {
	err = t.dao.WithContext(c).Model(&crmSftpFileModels.BrandConfigSftpModel{}).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseID,
			"sftp_type":     4,
			"deleted_at":    0,
		}).Count(&total).Error
	return
}

// GetExternalUserTagInfo 获取外部联系人的标签信息
func (t WorkTags) GetExternalUserTagInfo(c app.Ctx, enterpriseID, tagID, relationID int64) (data microScrmModel.WorkExternalUserTagModel, err error) {
	err = t.dao.WithContext(c).Model(&microScrmModel.WorkExternalUserTagModel{}).
		Where(map[string]interface{}{
			"enterprise_id":   enterpriseID,
			"external_tag_id": tagID,
			"relation_id":     relationID,
		}).Find(&data).Error
	return
}

// GetExternalTagList 企微标签列表
func (t WorkTags) GetExternalTagList(c app.Ctx, enterpriseID int64) (data []workModels.WorkExternalTagModel, err error) {
	err = t.dao.WithContext(c).Model(&workModels.WorkExternalTagModel{}).WithContext(c).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseID,
			"type":          1,
		}).
		Where("third_tag_id > 0").
		Order("`order` ASC").
		Find(&data).Error
	return
}

func (t WorkTags) GetExternalTagIds(c app.Ctx, enterpriseId int64, tagIds []string) ([]int64, error) {
	var result []int64
	err := t.dao.WithContext(c).Model(&workModels.WorkExternalTagModel{}).WithContext(c).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"type":          1,
		}).
		Where("tag_id in (?)", tagIds).
		Pluck("id", &result).Error
	return result, err
}

// GetChoiceExternalTagList 获取小程序端筛选条件的标签列表
func (t WorkTags) GetChoiceExternalTagList(c app.Ctx, enterpriseID int64) (data []workModels.WorkExternalTagModel, err error) {
	query := t.dao.WithContext(c).Model(&workModels.WorkExternalTagModel{}).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseID,
			"type":          workModels.WorkExternalTagTypeQy,
		})
	err = query.Find(&data).Error
	return
}

func (t WorkTags) GetListByIds(c app.Ctx, enterpriseID int64, ids []int64) ([]workModels.WorkExternalTagModel, error) {
	var result []workModels.WorkExternalTagModel
	err := t.dao.WithContext(c).Model(&workModels.WorkExternalTagModel{}).WithContext(c).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseID,
			"id":            ids,
		}).
		Find(&result).Error
	return result, err
}

func (t WorkTags) DelTagByQyId(c app.Ctx, enterpriseID int64, qyTagId string) error {
	return t.dao.WithContext(c).Model(&workModels.WorkExternalTagModel{}).Where(map[string]interface{}{
		"enterprise_id": enterpriseID,
		"tag_id":        qyTagId,
	}).Update("deleted_at", time.Now().Unix()).Error
}

func (t WorkTags) GetTagInfoByQyId(c app.Ctx, enterpriseID int64, qyTagId string) (workModels.WorkExternalTagModel, error) {
	var result workModels.WorkExternalTagModel
	err := t.dao.WithContext(c).Model(&workModels.WorkExternalTagModel{}).Where(map[string]interface{}{
		"enterprise_id": enterpriseID,
		"tag_id":        qyTagId,
	}).Find(&result).Error
	return result, err
}

func (t WorkTags) GetMany(c app.Ctx, query interface{}, args ...interface{}) (data []workModels.WorkExternalTagModel, err error) {
	err = t.dao.WithContext(c).Model(&workModels.WorkExternalTagModel{}).Where(query, args...).Find(&data).Error
	return
}
