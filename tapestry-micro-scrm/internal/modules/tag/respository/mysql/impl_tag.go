package mysql

import (
	"errors"
	"gorm.io/gorm"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"micro-scrm/internal/modules/tag"
	"micro-scrm/internal/modules/tag/entity"
)

type workTagRepo struct {
	*repos.SQLRepo
}

var (
	_ tag.WorkRepository = &workTagRepo{}
)

// 新增标签
func (tag *workTagRepo) AddTag(c ctx.Context, data entity.WorkExternalTagModel) (int64, error) {
	stmt, tmpErr := tag.DB.Template("t_work_external_tag/insert")
	if tmpErr != nil {
		return 0, tmpErr
	}

	res, resErr := stmt.NamedExec(c, data)
	if resErr != nil {
		return 0, resErr
	}

	return res.RowsAffected()
}

// 根据企业ID和标签组ID，获取标签组信息
func (tag *workTagRepo) GetTagGroupInfoByID(c ctx.Context, param entity.CreateManualTagParam) (*entity.WorkExternalTagGroupModel, error) {
	stmt, tmpErr := tag.DB.Template("t_work_external_tag_group/get-tag-group")
	if tmpErr != nil {
		return nil, tmpErr
	}

	data := entity.WorkExternalTagGroupModel{}
	resErr := stmt.NamedGet(c, &data, param)
	if resErr != nil {
		return nil, resErr
	}

	return &data, nil
}

// 根据企业ID和标签ID，更新标签的名称
func (tag *workTagRepo) UpdateTagNameByID(c ctx.Context, param entity.UpdateManualTagParam) error {
	stmt, tmpErr := tag.DB.Template("t_work_external_tag/update-tag-name-by-id")
	if tmpErr != nil {
		return tmpErr
	}

	res, resErr := stmt.NamedExec(c, map[string]interface{}{
		"enterprise_id": param.EnterpriseID,
		"tag_name":      param.NewTagName,
		"id":            param.TagID,
		"updated_at":    time.Now().Unix(),
	})
	if resErr != nil {
		return resErr
	}

	num, _ := res.RowsAffected()
	if num <= 0 {
		return errors.New("标签名称更新失败")
	}

	return nil
}

// IsExistTagByName 根据企业ID和标签名称，判断标签是否存在
func (tag *workTagRepo) IsExistTagByName(c ctx.Context, param entity.UpdateManualTagParam) (bool, error) {
	where := map[string]interface{}{
		"enterprise_id": param.EnterpriseID,
		"name":          param.NewTagName,
		"deleted_at":    0,
	}

	var model micro_scrm.WorkExternalTagModel

	tagOrm := gormx.GetInstance().Where(where)
	if param.TagID != 0 {
		// 如果存在id，则不考虑当前id这条是否重复
		tagOrm = tagOrm.Where("id != ?", param.TagID)
	}
	err := tagOrm.First(&model).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	return model.ID > 0, err
}

// 根据标签ID，获取标签信息
func (tag *workTagRepo) GetTagInfoByID(c ctx.Context, params entity.UpdateManualTagParam) (*entity.WorkExternalTagModel, error) {
	stmt, tmpErr := tag.DB.Template("t_work_external_tag/select-tag-info-by-tag-id")
	if tmpErr != nil {
		return nil, tmpErr
	}

	data := entity.WorkExternalTagModel{}
	resErr := stmt.NamedGet(c, &data, params)
	if resErr != nil {
		return nil, resErr
	}

	return &data, nil
}

// 根据标签组ID，获取组中标签
// 标签信息需要统计每个标签使用的人数
func (tag *workTagRepo) QueryTagListByGroupId(c ctx.Context, params entity.TagParam) ([]entity.WorkExternalTagEntity, error) {
	stmt, tmpErr := tag.DB.Template("t_work_external_tag/select-tag-use-num-list", tag.Cond2Args(params))
	if tmpErr != nil {
		return nil, tmpErr
	}

	tagList := make([]entity.WorkExternalTagEntity, 0)
	resErr := stmt.NamedSelect(c, &tagList, params)
	if resErr != nil {
		return nil, resErr
	}

	return tagList, nil
}

// 获取企业下的所有的标签组信息
func (tag *workTagRepo) GetTagGroupListByEnterpriseId(c ctx.Context, enterpriseID int64) ([]entity.WorkExternalTagGroupEntity, error) {
	stmt, tmpErr := tag.DB.Template("t_work_external_tag_group/select-group-list-by-enterprise-id")
	if tmpErr != nil {
		return nil, tmpErr
	}

	tagGroupList := make([]entity.WorkExternalTagGroupEntity, 0)
	resErr := stmt.Get(c, tagGroupList, enterpriseID)
	if resErr != nil {
		return nil, resErr
	}

	return tagGroupList, nil
}

// 根据标签组/标签名称，模糊匹配标签组信息
func (tag *workTagRepo) MatchTagGroupByGroupName(c ctx.Context, params entity.SearchTagParam) ([]entity.WorkExternalTagGroupEntity, error) {
	templateSql := "t_work_external_user_tag/select-tag-group-use-num-info"
	// if params.Type == 3 {
	//	templateSql = "t_work_external_user_tag/select-tag-and-group-list-type-keword"
	// }
	stmt, temErr := tag.DB.Template(templateSql, tag.Cond2Args(params))
	if temErr != nil {
		return nil, temErr
	}

	res, resErr := stmt.NamedQuery(c, map[string]interface{}{})
	if resErr != nil {
		return nil, resErr
	}
	defer res.Close()

	dataList := []entity.WorkExternalTagGroupEntity{}
	plusErr := sqlplus.ScanAll(res, &dataList)
	if plusErr != nil {
		return nil, plusErr
	}

	return dataList, nil
}

// 检测标签ID集合下是否存在外部联系人
func (tag *workTagRepo) CheckTagHasExternalUser(c ctx.Context, params entity.TagOperationParam) ([]entity.WorkExternalTagEntity, error) {
	mapping := map[string]interface{}{
		"enterprise_id": params.EnterpriseID,
		"tag_id_list":   params.TagId,
	}

	stmt, tmpErr := tag.DB.Template("t_work_external_user_tag/check-has-external-user", &mapping)
	if tmpErr != nil {
		return nil, tmpErr
	}

	data := []entity.WorkExternalTagEntity{}
	resErr := stmt.NamedSelect(c, &data, mapping)
	if resErr != nil {
		return nil, resErr
	}

	return data, nil
}

// 根据标签ID集合，删除标签信息
func (tag *workTagRepo) DeleteTag(c ctx.Context, params entity.TagOperationParam) (int64, error) {
	mapping := map[string]interface{}{
		"enterprise_id": params.EnterpriseID,
		"tag_id_list":   []int64{params.TagId},
		"deleted_at":    time.Now().Unix(),
	}

	stmt, tmpErr := tag.DB.Template("t_work_external_tag/delete-tag", &mapping)
	if tmpErr != nil {
		return 0, tmpErr
	}

	res, resErr := stmt.NamedExec(c, mapping)
	if resErr != nil {
		return 0, resErr
	}

	return res.RowsAffected()
}
