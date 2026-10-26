package mysql

// 标签组数据实现

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/tag"
	"micro-scrm/internal/modules/tag/entity"
	"time"
)

// 标签组相关服务
var (
	_ tag.WorkTagGroupRepository = &workTagGroupRepo{}
)

type workTagGroupRepo struct {
	*repos.SQLRepo
}

// 检测标签名称是是否已经存在
func (g *workTagGroupRepo) IsExistTagName(c ctx.Context, enterpriseID int64, tagNameList []string) (bool, error) {
	stmt, tmpErr := g.DB.Template("t_work_external_tag/is-exist-by-name-list", &map[string]interface{}{
		"name_list": tagNameList,
	})
	if tmpErr != nil {
		return false, tmpErr
	}

	count := struct {
		Count int `db:"use_num"`
	}{}
	resErr := stmt.NamedGet(c, &count, map[string]interface{}{"enterprise_id": enterpriseID})
	if resErr != nil {
		return false, nil
	}

	return count.Count > 0, resErr
}

// 根据企业的ID，和新增的标签组名称，检测该标签组名称是否已经被使用
func (g *workTagGroupRepo) GetTagGroupInfoByName(c ctx.Context, request entity.CreateManualTagGroupParams) (bool, error) {
	stmt, tmpErr := g.DB.Template("t_work_external_tag_group/get-tag-group-by-name")
	if tmpErr != nil {
		return false, tmpErr
	}

	data := entity.WorkExternalTagGroupModel{}
	resErr := stmt.NamedGet(c, &data, request)
	if resErr == nil {
		return true, nil
	} else {
		return true, resErr
	}
}

// 新增标签组
func (g *workTagGroupRepo) AddTagGroup(c ctx.Context, param entity.WorkExternalTagGroupModel) (int64, error) {
	stmt, tmpErr := g.DB.Template("t_work_external_tag_group/insert")
	if tmpErr != nil {
		return 0, tmpErr
	}

	res, resErr := stmt.NamedExec(c, param)
	if resErr != nil {
		return 0, resErr
	}

	return res.LastInsertId()
}

// 新增标签组下的标签
func (g *workTagGroupRepo) AddTag(c ctx.Context, param entity.WorkExternalTagModel) error {
	stmt, tmpErr := g.DB.Template("t_work_external_tag/insert")
	if tmpErr != nil {
		return tmpErr
	}

	_, resErr := stmt.NamedExec(c, param)
	if resErr != nil {
		return resErr
	}

	return nil
}

// 检测标签组下是否有用户存在
func (g *workTagGroupRepo) IsExistExternalUser(c ctx.Context, param entity.DeleteManualTagGroupRequestParam) (*entity.WorkExternalTagGroupEntity, error) {
	stmt, tmpErr := g.DB.Template("t_work_external_tag_group/is-exist-external-user")
	if tmpErr != nil {
		return nil, tmpErr
	}

	data := entity.WorkExternalTagGroupEntity{}
	resErr := stmt.NamedGet(c, &data, param)
	if resErr != nil {
		return nil, resErr
	}
	return &data, nil
}

// 删除标签组
func (g *workTagGroupRepo) DeleteTagGroup(c ctx.Context, param entity.DeleteTagGroup) (int64, error) {
	stmt, tmpErr := g.DB.Template("t_work_external_tag_group/delete-tag-group", g.Cond2Args(param))
	if tmpErr != nil {
		return 0, tmpErr
	}

	param.DeletedAt = time.Now().Unix()
	res, resErr := stmt.NamedExec(c, param)
	if resErr != nil {
		return 0, resErr
	}
	return res.RowsAffected()
}
