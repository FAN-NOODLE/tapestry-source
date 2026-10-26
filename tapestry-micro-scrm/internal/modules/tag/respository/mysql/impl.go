package mysql

import (
	"database/sql"
	"fmt"
	"micro-scrm/internal/pkg/db/gormx"
	"time"

	"github.com/jmoiron/sqlx"

	"micro-scrm/internal/modules/tag/entity"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

type tagRepo struct {
	*repos.SQLRepo
}

// 根据企业ID，获取标签信息，包含标签组信息
func (t *tagRepo) GetTagInfoList(c ctx.Context, enterpriseID int64, tagType, isThirdTag int) ([]entity.WorkTagInfoEntity, error) {
	mapping := make(map[string]interface{})
	mapping["enterprise_id"] = enterpriseID
	mapping["type"] = tagType
	mapping["is_third_tag"] = isThirdTag
	template := "t_work_external_tag/select-tag-and-group-list"
	// if mapping["type"] == 3 {
	//	template = "t_work_external_tag/select-tag-and-group-list-type"
	// }
	stmt, tmpErr := t.DB.Template(template, &mapping)
	if tmpErr != nil {
		return nil, tmpErr
	}

	tagList := make([]entity.WorkTagInfoEntity, 0)
	resErr := stmt.NamedSelect(c, &tagList, mapping)
	if resErr != nil {
		return nil, resErr
	}

	return tagList, nil
}

// QueryTags 通过id查询标签
func (t *tagRepo) QueryTags(c ctx.Context, IDList []int64) ([]*entity.TagEntity, error) {
	mapping := make(map[string]interface{})
	mapping["id_list"] = IDList
	stmt, err := t.DB.Template("t_tag/query-by-ids", &mapping)
	if err != nil {
		return nil, err
	}
	var list = make([]*entity.TagEntity, 0)
	err = stmt.NamedSelect(c, &list, mapping)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// QueryTags 通过id查询标签
func (t *tagRepo) QueryEnterpriseTags(c ctx.Context, enterpriseID int64) ([]*entity.TagEntity, error) {
	mapping := make(map[string]interface{})
	mapping["enterprise_id"] = enterpriseID
	stmt, err := t.DB.Template("t_tag/query-by-enterprise-id", &mapping)
	if err != nil {
		return nil, err
	}
	var list = make([]*entity.TagEntity, 0)
	err = stmt.NamedSelect(c, &list, mapping)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// QueryTags 通过id查询标签
func (t *tagRepo) QueryTagGroups(c ctx.Context, IDList []int64) ([]*entity.TagGroupEntity, error) {
	mapping := make(map[string]interface{})
	mapping["id_list"] = IDList
	stmt, err := t.DB.Template("t_tag_group/query-by-ids", &mapping)
	if err != nil {
		return nil, err
	}
	var list = make([]*entity.TagGroupEntity, 0)
	err = stmt.NamedSelect(c, &list, mapping)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// QueryUserTagGroups 通过用户分组id查询用户分组
func (t *tagRepo) QueryUserTagGroups(c ctx.Context, IDList []int64) ([]*entity.UserGroupEntity, error) {
	mapping := make(map[string]interface{})
	mapping["id_list"] = IDList
	stmt, err := t.DB.Template("t_user_group/query-by-ids", &mapping)
	if err != nil {
		return nil, err
	}
	var list = make([]*entity.UserGroupEntity, 0)
	err = stmt.NamedSelect(c, &list, mapping)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// QueryTagGroupsOfEnterprise 查询企业下所有标签组
func (t *tagRepo) QueryEnterpriseTagGroups(c ctx.Context, enterpriseID int64) ([]*entity.TagGroupEntity, error) {
	mapping := make(map[string]interface{})
	mapping["enterprise_id"] = enterpriseID
	stmt, err := t.DB.Template("t_tag_group/query-by-enterprise-id", &mapping)
	if err != nil {
		return nil, err
	}
	var list = make([]*entity.TagGroupEntity, 0)
	err = stmt.NamedSelect(c, &list, mapping)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// AddTagGroup 添加标签分组
func (t *tagRepo) AddTagGroup(c ctx.Context, data []*entity.TagGroupEntity) (int64, error) {
	stmt, err := t.DB.Template("t_tag_group/insert")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// AddTag 添加标签数据
func (t *tagRepo) AddTag(c ctx.Context, data []*entity.TagEntity) (int64, error) {
	stmt, err := t.DB.Template("t_tag/insert")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// AddUserGroup 添加用户分组
func (t *tagRepo) AddUserGroup(c ctx.Context, data []*entity.UserGroupEntity) (int64, error) {
	stmt, err := t.DB.Template("t_user_group/insert")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// GetTagList 获取标签列表
func (t *tagRepo) GetTagList(c ctx.Context, enterpriseID int64, search string) ([]*entity.TagEntity, error) {
	mapping := make(map[string]interface{})
	mapping["search"] = search
	mapping["enterprise_id"] = enterpriseID
	stmt, err := t.DB.Template("t_tag/query-list", &mapping)
	if err != nil {
		return nil, err
	}
	var list = make([]*entity.TagEntity, 0)
	err = stmt.NamedSelect(c, &list, mapping)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetTagListPage 分页获取标签列表
func (t *tagRepo) GetTagListPage(c ctx.Context, req *entity.GetTagListRequest) (int64, []*entity.TagEntity, error) {
	var tpl = "t_tag/query-list"
	offset := (req.Page - 1) * req.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(req.Size)}
	option := &repos.FindPageOption{
		TplPrefix: tpl,
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]*entity.TagEntity{}
		},
	}
	total, list, err := t.FindPage(c, pageVo, req, option)
	resList := list.(*[]*entity.TagEntity)
	return total, *resList, err
}

// GetUserGroupList 获取用户分组列表
func (t *tagRepo) GetUserGroupList(c ctx.Context, enterpriseID int64, search string) ([]*entity.UserGroupEntity, error) {
	mapping := make(map[string]interface{})
	mapping["search"] = search
	mapping["enterprise_id"] = enterpriseID
	stmt, err := t.DB.Template("t_user_group/query-list", &mapping)
	if err != nil {
		return nil, err
	}
	var list = make([]*entity.UserGroupEntity, 0)
	err = stmt.NamedSelect(c, &list, mapping)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (t *tagRepo) GetUserGroupListPage(c ctx.Context, req *entity.GetTagListRequest) (int64, []*entity.UserGroupEntity, error) {
	var tpl = "t_user_group/query-list"
	offset := (req.Page - 1) * req.Size
	pageVo := repos.PageVO{Offset: offset, Limit: req.Size}
	var userGroupEntity []*entity.UserGroupEntity
	option := &repos.FindPageOption{
		TplPrefix: tpl,
		Fn: func(context ctx.Context, i int64) interface{} {
			return &userGroupEntity
		},
	}
	total, _, err := t.FindPage(c, pageVo, req, option)
	return total, userGroupEntity, err
}

// UpdateTagGroup 更新标签组数据
func (t *tagRepo) UpdateTagGroup(c ctx.Context, data *entity.TagGroupEntity) (int64, error) {
	stmt, err := t.DB.Template("t_tag_group/update")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// DeleteTagGroup 删除标签组
func (t *tagRepo) DeleteTagGroup(c ctx.Context, enterpriseID int64, tagGroupIDList []int64) (int64, error) {
	mapping := make(map[string]interface{})
	mapping["group_id_list"] = tagGroupIDList
	mapping["enterprise_id"] = enterpriseID
	mapping["deleted_at"] = time.Now().Unix()
	stmt, err := t.DB.Template("t_tag_group/delete-group-id-list", &mapping)
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, mapping)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// UpdateTag 更新标签数据
func (t *tagRepo) UpdateTag(c ctx.Context, data *entity.TagEntity) (int64, error) {
	stmt, err := t.DB.Template("t_tag/update")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// DeleteTag 根据标签id删除标签数据
func (t *tagRepo) DeleteTag(c ctx.Context, enterpriseID int64, tagIDList []int64) (int64, error) {
	mapping := make(map[string]interface{})
	mapping["tag_id_list"] = tagIDList
	mapping["enterprise_id"] = enterpriseID
	mapping["deleted_at"] = time.Now().Unix()
	stmt, err := t.DB.Template("t_tag/delete-tag-id-list", &mapping)
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, mapping)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// UpdateUserGroup 更新用户分组数据
func (t *tagRepo) UpdateUserGroup(c ctx.Context, data *entity.UserGroupEntity) (int64, error) {
	stmt, err := t.DB.Template("t_user_group/update")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// DeleteUserGroup 删除用户分组
func (t *tagRepo) DeleteUserGroup(c ctx.Context, enterpriseID int64, userGroupIDList []int64) (int64, error) {
	mapping := make(map[string]interface{})
	mapping["user_group_id_list"] = userGroupIDList
	mapping["enterprise_id"] = enterpriseID
	mapping["deleted_at"] = time.Now().Unix()
	stmt, err := t.DB.Template("t_user_group/delete-group-id-list", &mapping)
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, mapping)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// WorkExternalTagList 企业微信标签
func (t *tagRepo) WorkExternalTagList(cc ctx.Context, enterpriseID int64, relationID int64) (*[]entity.WorkExternalTag, error) {
	stmt, err := t.DB.Template("t_work_external_user_tag/select-work-external-tag")
	if err != nil {
		return nil, err
	}
	list := &[]entity.WorkExternalTag{}
	sErr := stmt.Select(cc, list, enterpriseID, relationID, 0)
	if sErr != nil && sErr != sql.ErrNoRows {
		return nil, sErr
	}
	return list, nil
}

// UserGroupByTag 根据标签查询客户分组
func (t *tagRepo) UserGroupByTag(c ctx.Context, enterpriseID int64, tagIDList []int64) ([]entity.UserGroupDetailTagInfo, error) {
	dataList := make([]entity.UserGroupDetailTagInfo, 0)

	stmt, err := t.DB.Template("t_user_group_detail/user-group-info-by-tag-id-list", &map[string]interface{}{
		"tag_id_list": tagIDList,
	})
	if err != nil {
		return dataList, err
	}

	// include_type = 1 为包含标签
	resErr := stmt.NamedSelect(c, &dataList, map[string]interface{}{"enterprise_id": enterpriseID, "include_type": 1})
	if resErr != nil {
		return dataList, resErr
	}

	return dataList, nil
}

func (t *tagRepo) GetAllUsedExternalTags(ctx ctx.Context, enterpriseID int64) ([]int64, error) {
	// todo 慢sql
	list := make([]int64, 0)
	sqlStr := "SELECT DISTINCT(UT.external_tag_id) FROM micro_scrm_ck.work_external_user_tag UT GLOBAL INNER JOIN micro_scrm_ck.work_external_user_relation UR ON UR.id = UT.relation_id WHERE UT.enterprise_id = ? AND UT.deleted_at = 0 and UR.relation_status = 0 and UR.deleted_at = 0"
	err := gormx.GetCkInstance().WithContext(ctx).Raw(sqlStr, enterpriseID).Error
	// err := t.DB.Select(cc, &list, sqlStr, enterpriseID)
	if err != nil {
		return list, err
	}
	return list, nil
}

func (t *tagRepo) GetAllUsedExternalGroupByTags(cc ctx.Context, enterpriseID int64, tagIDs []int64) ([]entity.UserGroupDetailTagInfo, error) {
	list := make([]entity.UserGroupDetailTagInfo, 0)
	query := "select UGD.`user_group_id` AS user_group_id,UGD.`tag_group_id` AS tag_group_id,UGD.`tag_id`,UGD.`include_type` from `user_group_detail` UGD INNER JOIN `user_group` UG ON UG.`id` = UGD.`user_group_id` where UGD.`enterprise_id` = ? AND UGD.`tag_id` IN (?) AND UGD.`include_type`=1 AND UGD.`deleted_at` = 0 AND UG.`deleted_at` = 0"
	query, args, err := sqlx.In(query, enterpriseID, tagIDs)
	if err != nil {
		return nil, err
	}
	err = t.DB.Select(cc, &list, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return list, err
	}
	return list, nil
}

func (t *tagRepo) QueryWorkExternalTagList(c ctx.Context, tagIDs []int64, relationID int64) ([]entity.WorkExternalTag, error) {
	list := make([]entity.WorkExternalTag, 0)
	if relationID > 0 {
		query := fmt.Sprintf("select a.tag_id,a.`name` as tag_name from work_external_tag a,work_external_user_tag b where a.id = b.external_tag_id AND a.id in(?) AND b.relation_id = %d and a.deleted_at=0 group by a.tag_id", relationID)
		query, args, err := sqlx.In(query, tagIDs)
		if err != nil {
			return nil, err
		}
		err = t.DB.Select(c, &list, query, args...)
		if err != nil && err != sql.ErrNoRows {
			return list, err
		}
	} else {
		query := "select tag_id,`name` as tag_name from work_external_tag where id in(?) and deleted_at=0"
		query, args, err := sqlx.In(query, tagIDs)
		if err != nil {
			return nil, err
		}
		err = t.DB.Select(c, &list, query, args...)
		if err != nil && err != sql.ErrNoRows {
			return list, err
		}
	}
	return list, nil
}
